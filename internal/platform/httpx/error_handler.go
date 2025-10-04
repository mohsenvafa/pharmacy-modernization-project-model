package httpx

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"go.uber.org/zap"

	platformErrors "pharmacy-modernization-project-model/internal/platform/errors"
)

// ErrorHandler provides centralized error handling for HTTP responses
type ErrorHandler struct {
	logger *zap.Logger
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(logger *zap.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

// APIError represents a standardized API error response
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// HandleError handles different types of errors and returns appropriate HTTP responses
func (eh *ErrorHandler) HandleError(w http.ResponseWriter, err error) {
	// Handle custom error types
	switch {
	// Handle ValidationError
	case func() bool {
		var validationErr platformErrors.ValidationError
		return errors.As(err, &validationErr)
	}():
		var validationErr platformErrors.ValidationError
		errors.As(err, &validationErr)
		eh.writeError(w, http.StatusBadRequest, APIError{
			Code:    "validation_error",
			Message: validationErr.Error(),
			Details: validationErr.Field,
		})

	// Handle RecordNotFoundError
	case func() bool {
		var notFoundErr platformErrors.RecordNotFoundError
		return errors.As(err, &notFoundErr)
	}():
		var notFoundErr platformErrors.RecordNotFoundError
		errors.As(err, &notFoundErr)
		eh.writeError(w, http.StatusNotFound, APIError{
			Code:    "record_not_found",
			Message: notFoundErr.Error(),
			Details: notFoundErr.Type,
		})

	// Handle DuplicateRecordError
	case func() bool {
		var duplicateErr platformErrors.DuplicateRecordError
		return errors.As(err, &duplicateErr)
	}():
		var duplicateErr platformErrors.DuplicateRecordError
		errors.As(err, &duplicateErr)
		eh.writeError(w, http.StatusConflict, APIError{
			Code:    "duplicate_record",
			Message: duplicateErr.Error(),
			Details: duplicateErr.Type,
		})

	// Handle BusinessLogicError
	case func() bool {
		var businessErr platformErrors.BusinessLogicError
		return errors.As(err, &businessErr)
	}():
		var businessErr platformErrors.BusinessLogicError
		errors.As(err, &businessErr)
		eh.writeError(w, http.StatusUnprocessableEntity, APIError{
			Code:    "business_logic_error",
			Message: businessErr.Error(),
			Details: businessErr.Operation,
		})

	// Handle ConfigurationError
	case func() bool {
		var configErr platformErrors.ConfigurationError
		return errors.As(err, &configErr)
	}():
		var configErr platformErrors.ConfigurationError
		errors.As(err, &configErr)
		eh.logger.Error("Configuration error", zap.Error(err))
		eh.writeError(w, http.StatusInternalServerError, APIError{
			Code:    "configuration_error",
			Message: "Configuration error",
		})

	// Handle AuthorizationError
	case func() bool {
		var authErr platformErrors.AuthorizationError
		return errors.As(err, &authErr)
	}():
		var authErr platformErrors.AuthorizationError
		errors.As(err, &authErr)
		eh.writeError(w, http.StatusForbidden, APIError{
			Code:    "authorization_error",
			Message: authErr.Error(),
			Details: authErr.Resource,
		})

	// Handle ExternalServiceError
	case func() bool {
		var serviceErr platformErrors.ExternalServiceError
		return errors.As(err, &serviceErr)
	}():
		var serviceErr platformErrors.ExternalServiceError
		errors.As(err, &serviceErr)
		eh.logger.Error("External service error", zap.Error(err))
		eh.writeError(w, http.StatusBadGateway, APIError{
			Code:    "external_service_error",
			Message: "External service temporarily unavailable",
		})

	// Handle RateLimitError
	case func() bool {
		var rateLimitErr platformErrors.RateLimitError
		return errors.As(err, &rateLimitErr)
	}():
		var rateLimitErr platformErrors.RateLimitError
		errors.As(err, &rateLimitErr)
		eh.writeError(w, http.StatusTooManyRequests, APIError{
			Code:    "rate_limit_exceeded",
			Message: rateLimitErr.Error(),
			Details: rateLimitErr.Resource,
		})

	// Handle common domain errors
	case errors.Is(err, platformErrors.ErrIDRequired):
		eh.writeError(w, http.StatusBadRequest, APIError{
			Code:    "id_required",
			Message: "ID is required",
		})
	case errors.Is(err, platformErrors.ErrNameRequired):
		eh.writeError(w, http.StatusBadRequest, APIError{
			Code:    "name_required",
			Message: "Name is required",
		})
	case errors.Is(err, platformErrors.ErrEmailRequired):
		eh.writeError(w, http.StatusBadRequest, APIError{
			Code:    "email_required",
			Message: "Email is required",
		})
	case errors.Is(err, platformErrors.ErrPhoneRequired):
		eh.writeError(w, http.StatusBadRequest, APIError{
			Code:    "phone_required",
			Message: "Phone is required",
		})
	case errors.Is(err, platformErrors.ErrInvalidEmail):
		eh.writeError(w, http.StatusBadRequest, APIError{
			Code:    "invalid_email",
			Message: "Invalid email format",
		})
	case errors.Is(err, platformErrors.ErrInvalidPhone):
		eh.writeError(w, http.StatusBadRequest, APIError{
			Code:    "invalid_phone",
			Message: "Invalid phone format",
		})
	case errors.Is(err, platformErrors.ErrInvalidID):
		eh.writeError(w, http.StatusBadRequest, APIError{
			Code:    "invalid_id",
			Message: "Invalid ID format",
		})

	// Handle generic validation errors
	case strings.Contains(err.Error(), "validation"):
		eh.writeError(w, http.StatusBadRequest, APIError{
			Code:    "validation_error",
			Message: "Validation failed",
		})

	// Handle generic not found errors
	case strings.Contains(err.Error(), "not found"):
		eh.writeError(w, http.StatusNotFound, APIError{
			Code:    "not_found",
			Message: "Resource not found",
		})

	// Handle generic duplicate errors
	case strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "duplicate"):
		eh.writeError(w, http.StatusConflict, APIError{
			Code:    "duplicate",
			Message: "Resource already exists",
		})

	// Default case - internal server error
	default:
		eh.logger.Error("Internal server error", zap.Error(err))
		eh.writeError(w, http.StatusInternalServerError, APIError{
			Code:    "internal_error",
			Message: "Internal server error",
		})
	}
}

// writeError writes an error response to the HTTP response writer
func (eh *ErrorHandler) writeError(w http.ResponseWriter, statusCode int, apiError APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Simple JSON response (in a real app, you'd use a JSON encoder)
	jsonResponse := fmt.Sprintf(`{"code":"%s","message":"%s","details":"%s"}`,
		apiError.Code, apiError.Message, apiError.Details)
	w.Write([]byte(jsonResponse))
}

// WriteError is a convenience function for writing error responses
func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	// Extract logger from request context if available
	logger := zap.NewNop() // Default no-op logger
	if ctxLogger := r.Context().Value("logger"); ctxLogger != nil {
		if l, ok := ctxLogger.(*zap.Logger); ok {
			logger = l
		}
	}

	errorHandler := NewErrorHandler(logger)
	errorHandler.HandleError(w, err)
}
