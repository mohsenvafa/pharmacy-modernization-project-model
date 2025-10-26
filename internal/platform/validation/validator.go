package validation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	platformErrors "pharmacy-modernization-project-model/internal/platform/errors"
)

// Validator provides common validation functions
type Validator struct{}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateRequired validates that a field is not empty
func (v *Validator) ValidateRequired(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return platformErrors.NewValidationError(field, value, "field is required")
	}
	return nil
}

// ValidatePhone validates phone number format (basic validation)
func (v *Validator) ValidatePhone(field, phone string) error {
	if phone == "" {
		return platformErrors.NewValidationError(field, phone, "phone is required")
	}

	// Remove all non-digit characters for validation
	digitsOnly := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
	if len(digitsOnly) < 10 || len(digitsOnly) > 15 {
		return platformErrors.NewValidationError(field, phone, "phone number must be between 10-15 digits")
	}
	return nil
}

// ValidateLength validates string length
func (v *Validator) ValidateLength(field, value string, min, max int) error {
	length := utf8.RuneCountInString(value)
	if length < min {
		return platformErrors.NewValidationError(field, value, fmt.Sprintf("must be at least %d characters long", min))
	}
	if length > max {
		return platformErrors.NewValidationError(field, value, fmt.Sprintf("must be no more than %d characters long", max))
	}
	return nil
}

// ValidateID validates ID format (alphanumeric with optional hyphens/underscores)
func (v *Validator) ValidateID(field, id string) error {
	if id == "" {
		return platformErrors.NewValidationError(field, id, "ID is required")
	}

	idRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !idRegex.MatchString(id) {
		return platformErrors.NewValidationError(field, id, "ID can only contain letters, numbers, hyphens, and underscores")
	}

	if len(id) < 1 || len(id) > 50 {
		return platformErrors.NewValidationError(field, id, "ID must be between 1-50 characters")
	}
	return nil
}

// ValidateRange validates numeric range
func (v *Validator) ValidateRange(field string, value, min, max int) error {
	if value < min {
		return platformErrors.NewValidationError(field, value, fmt.Sprintf("must be at least %d", min))
	}
	if value > max {
		return platformErrors.NewValidationError(field, value, fmt.Sprintf("must be no more than %d", max))
	}
	return nil
}

// ValidatePositive validates that a number is positive
func (v *Validator) ValidatePositive(field string, value int) error {
	if value <= 0 {
		return platformErrors.NewValidationError(field, value, "must be positive")
	}
	return nil
}

// ValidateNonNegative validates that a number is non-negative
func (v *Validator) ValidateNonNegative(field string, value int) error {
	if value < 0 {
		return platformErrors.NewValidationError(field, value, "must be non-negative")
	}
	return nil
}

// ValidateOneOf validates that a value is one of the allowed values
func (v *Validator) ValidateOneOf(field string, value interface{}, allowedValues ...interface{}) error {
	for _, allowed := range allowedValues {
		if value == allowed {
			return nil
		}
	}
	return platformErrors.NewValidationError(field, value, fmt.Sprintf("must be one of: %v", allowedValues))
}

// ValidateRegex validates that a string matches a regex pattern
func (v *Validator) ValidateRegex(field, value, pattern, message string) error {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return platformErrors.NewValidationError(field, value, "invalid validation pattern")
	}

	if !regex.MatchString(value) {
		return platformErrors.NewValidationError(field, value, message)
	}
	return nil
}
