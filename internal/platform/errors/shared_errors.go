package errors

import (
	"errors"
	"fmt"
)

// ValidationError represents a validation error with field details
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation failed: %s", e.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field string, value interface{}, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}

// RecordNotFoundError represents a record not found error
type RecordNotFoundError struct {
	Type string
	ID   string
}

func (e RecordNotFoundError) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("%s with ID '%s' not found", e.Type, e.ID)
	}
	return fmt.Sprintf("%s not found", e.Type)
}

// NewRecordNotFoundError creates a new record not found error
func NewRecordNotFoundError(recordType, id string) RecordNotFoundError {
	return RecordNotFoundError{
		Type: recordType,
		ID:   id,
	}
}

// DuplicateRecordError represents a duplicate record error
type DuplicateRecordError struct {
	Type string
	ID   string
}

func (e DuplicateRecordError) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("%s with ID '%s' already exists", e.Type, e.ID)
	}
	return fmt.Sprintf("%s already exists", e.Type)
}

// NewDuplicateRecordError creates a new duplicate record error
func NewDuplicateRecordError(recordType, id string) DuplicateRecordError {
	return DuplicateRecordError{
		Type: recordType,
		ID:   id,
	}
}

// BusinessLogicError represents a business logic violation
type BusinessLogicError struct {
	Operation string
	Reason    string
}

func (e BusinessLogicError) Error() string {
	return fmt.Sprintf("business logic error in %s: %s", e.Operation, e.Reason)
}

// NewBusinessLogicError creates a new business logic error
func NewBusinessLogicError(operation, reason string) BusinessLogicError {
	return BusinessLogicError{
		Operation: operation,
		Reason:    reason,
	}
}

// ConfigurationError represents a configuration error
type ConfigurationError struct {
	Component string
	Setting   string
	Message   string
}

func (e ConfigurationError) Error() string {
	return fmt.Sprintf("configuration error in %s.%s: %s", e.Component, e.Setting, e.Message)
}

// NewConfigurationError creates a new configuration error
func NewConfigurationError(component, setting, message string) ConfigurationError {
	return ConfigurationError{
		Component: component,
		Setting:   setting,
		Message:   message,
	}
}

// AuthorizationError represents an authorization error
type AuthorizationError struct {
	Resource string
	Action   string
	Reason   string
}

func (e AuthorizationError) Error() string {
	return fmt.Sprintf("authorization failed for %s on %s: %s", e.Action, e.Resource, e.Reason)
}

// NewAuthorizationError creates a new authorization error
func NewAuthorizationError(resource, action, reason string) AuthorizationError {
	return AuthorizationError{
		Resource: resource,
		Action:   action,
		Reason:   reason,
	}
}

// ExternalServiceError represents an error from an external service
type ExternalServiceError struct {
	Service   string
	Operation string
	Message   string
}

func (e ExternalServiceError) Error() string {
	return fmt.Sprintf("external service error from %s during %s: %s", e.Service, e.Operation, e.Message)
}

// NewExternalServiceError creates a new external service error
func NewExternalServiceError(service, operation, message string) ExternalServiceError {
	return ExternalServiceError{
		Service:   service,
		Operation: operation,
		Message:   message,
	}
}

// RateLimitError represents a rate limiting error
type RateLimitError struct {
	Resource string
	Limit    int
	Window   string
}

func (e RateLimitError) Error() string {
	return fmt.Sprintf("rate limit exceeded for %s: %d requests per %s", e.Resource, e.Limit, e.Window)
}

// NewRateLimitError creates a new rate limit error
func NewRateLimitError(resource string, limit int, window string) RateLimitError {
	return RateLimitError{
		Resource: resource,
		Limit:    limit,
		Window:   window,
	}
}

// Common domain-specific errors that can be used across domains
var (
	ErrIDRequired    = errors.New("ID is required")
	ErrNameRequired  = errors.New("name is required")
	ErrEmailRequired = errors.New("email is required")
	ErrPhoneRequired = errors.New("phone is required")
	ErrInvalidEmail  = errors.New("invalid email format")
	ErrInvalidPhone  = errors.New("invalid phone format")
	ErrInvalidID     = errors.New("invalid ID format")
	ErrEmptyField    = errors.New("field cannot be empty")
	ErrInvalidFormat = errors.New("invalid format")
	ErrOutOfRange    = errors.New("value is out of allowed range")
	ErrTooShort      = errors.New("value is too short")
	ErrTooLong       = errors.New("value is too long")
)

