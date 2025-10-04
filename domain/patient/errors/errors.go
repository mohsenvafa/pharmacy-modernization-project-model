package errors

import (
	"errors"

	platformErrors "pharmacy-modernization-project-model/internal/platform/errors"
)

// Domain-specific errors for patient operations
var (
	ErrPatientNotFound   = errors.New("patient not found")
	ErrInvalidPatient    = errors.New("invalid patient data")
	ErrDuplicatePatient  = errors.New("patient already exists")
	ErrInvalidPatientID  = errors.New("invalid patient ID")
	ErrPatientIDRequired = errors.New("patient ID is required")
)

// Re-export platform errors for convenience
type ValidationError = platformErrors.ValidationError
type RecordNotFoundError = platformErrors.RecordNotFoundError
type DuplicateRecordError = platformErrors.DuplicateRecordError
type BusinessLogicError = platformErrors.BusinessLogicError
type ConfigurationError = platformErrors.ConfigurationError
type AuthorizationError = platformErrors.AuthorizationError
type ExternalServiceError = platformErrors.ExternalServiceError
type RateLimitError = platformErrors.RateLimitError

// Convenience functions for creating platform errors
var (
	NewValidationError      = platformErrors.NewValidationError
	NewRecordNotFoundError  = platformErrors.NewRecordNotFoundError
	NewDuplicateRecordError = platformErrors.NewDuplicateRecordError
	NewBusinessLogicError   = platformErrors.NewBusinessLogicError
	NewConfigurationError   = platformErrors.NewConfigurationError
	NewAuthorizationError   = platformErrors.NewAuthorizationError
	NewExternalServiceError = platformErrors.NewExternalServiceError
	NewRateLimitError       = platformErrors.NewRateLimitError
)
