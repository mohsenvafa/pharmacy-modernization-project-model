package errors

import (
	"errors"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

// ErrorType represents the type of repository error
type ErrorType string

const (
	ErrorTypeNotFound      ErrorType = "NotFound"
	ErrorTypeDuplicateKey  ErrorType = "DuplicateKey"
	ErrorTypeTimeout       ErrorType = "Timeout"
	ErrorTypeNetworkError  ErrorType = "NetworkError"
	ErrorTypeValidation    ErrorType = "Validation"
	ErrorTypeConnection    ErrorType = "Connection"
	ErrorTypeDatabaseError ErrorType = "DatabaseError"
	ErrorTypeUnknown       ErrorType = "Unknown"
)

// RepositoryError represents a repository-specific error
type RepositoryError struct {
	Type    ErrorType
	Message string
	Err     error
	Context map[string]interface{}
}

func (e *RepositoryError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *RepositoryError) Unwrap() error {
	return e.Err
}

// IsNotFound checks if the error is a "not found" error
func (e *RepositoryError) IsNotFound() bool {
	return e.Type == ErrorTypeNotFound
}

// IsDuplicateKey checks if the error is a duplicate key error
func (e *RepositoryError) IsDuplicateKey() bool {
	return e.Type == ErrorTypeDuplicateKey
}

// IsTimeout checks if the error is a timeout error
func (e *RepositoryError) IsTimeout() bool {
	return e.Type == ErrorTypeTimeout
}

// IsNetworkError checks if the error is a network error
func (e *RepositoryError) IsNetworkError() bool {
	return e.Type == ErrorTypeNetworkError
}

// IsRetryable checks if the error is retryable
func (e *RepositoryError) IsRetryable() bool {
	return e.Type == ErrorTypeTimeout || e.Type == ErrorTypeNetworkError || e.Type == ErrorTypeConnection
}

// NewRepositoryError creates a new repository error
func NewRepositoryError(errorType ErrorType, message string, err error) *RepositoryError {
	return &RepositoryError{
		Type:    errorType,
		Message: message,
		Err:     err,
		Context: make(map[string]interface{}),
	}
}

// WithContext adds context to the error
func (e *RepositoryError) WithContext(key string, value interface{}) *RepositoryError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// HandleMongoError converts MongoDB errors to RepositoryError
func HandleMongoError(operation string, err error) *RepositoryError {
	if err == nil {
		return nil
	}

	// Handle specific MongoDB errors
	switch {
	case err == mongo.ErrNoDocuments:
		return NewRepositoryError(ErrorTypeNotFound, "Document not found", err).
			WithContext("operation", operation)

	case mongo.IsDuplicateKeyError(err):
		return NewRepositoryError(ErrorTypeDuplicateKey, "Duplicate key error", err).
			WithContext("operation", operation)

	case mongo.IsTimeout(err):
		return NewRepositoryError(ErrorTypeTimeout, "Operation timed out", err).
			WithContext("operation", operation)

	case mongo.IsNetworkError(err):
		return NewRepositoryError(ErrorTypeNetworkError, "Network error occurred", err).
			WithContext("operation", operation)

	case isConnectionError(err):
		return NewRepositoryError(ErrorTypeConnection, "Connection error", err).
			WithContext("operation", operation)

	case isValidationError(err):
		return NewRepositoryError(ErrorTypeValidation, "Validation error", err).
			WithContext("operation", operation)

	default:
		return NewRepositoryError(ErrorTypeDatabaseError, "Database operation failed", err).
			WithContext("operation", operation)
	}
}

// isConnectionError checks if the error is a connection error
func isConnectionError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	connectionErrors := []string{
		"connection refused",
		"connection reset",
		"connection timeout",
		"no connection",
		"connection lost",
		"server selection error",
		"topology closed",
	}

	for _, connErr := range connectionErrors {
		if strings.Contains(errStr, connErr) {
			return true
		}
	}

	return false
}

// isValidationError checks if the error is a validation error
func isValidationError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	validationErrors := []string{
		"validation failed",
		"invalid",
		"required",
		"constraint",
		"schema",
	}

	for _, valErr := range validationErrors {
		if strings.Contains(errStr, valErr) {
			return true
		}
	}

	return false
}

// IsNotFoundError checks if the error is a not found error
func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	var repoErr *RepositoryError
	if errors.As(err, &repoErr) {
		return repoErr.IsNotFound()
	}

	var notFoundErr RecordNotFoundError
	if errors.As(err, &notFoundErr) {
		return true
	}

	return false
}
