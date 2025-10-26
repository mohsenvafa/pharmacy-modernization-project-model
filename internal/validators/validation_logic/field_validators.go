// internal/validators/validation_logic/field_validators.go
package validation_logic

import (
	"fmt"
	"strings"
)

// FieldError represents a validation error for a specific field
type FieldError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Param   string `json:"param,omitempty"`
	Message string `json:"message,omitempty"`
}

// Error implements the error interface
func (e *FieldError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Field != "" {
		return e.Field + ": " + e.Tag
	}
	return e.Tag
}

// ValidateRequired validates that a field is not empty
func ValidateRequired(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return &FieldError{Field: field, Tag: "required", Message: "field is required"}
	}
	return nil
}

// ValidateID validates ID format (alphanumeric with optional hyphens/underscores)
func ValidateID(field, id string) error {
	if id == "" {
		return &FieldError{Field: field, Tag: "required", Message: "ID is required"}
	}

	// Check for valid ID format (alphanumeric with optional hyphens/underscores)
	if len(id) < 1 || len(id) > 50 {
		return &FieldError{Field: field, Tag: "len", Param: "1,50", Message: "ID must be between 1-50 characters"}
	}

	// Basic alphanumeric check (can be enhanced)
	for _, r := range id {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_') {
			return &FieldError{Field: field, Tag: "alphanum", Message: "ID can only contain letters, numbers, hyphens, and underscores"}
		}
	}
	return nil
}

// ValidateLength validates string length
func ValidateLength(field, value string, min, max int) error {
	length := len(value)
	if length < min {
		return &FieldError{Field: field, Tag: "min", Param: fmt.Sprintf("%d", min), Message: fmt.Sprintf("must be at least %d characters long", min)}
	}
	if length > max {
		return &FieldError{Field: field, Tag: "max", Param: fmt.Sprintf("%d", max), Message: fmt.Sprintf("must be no more than %d characters long", max)}
	}
	return nil
}

// ValidateOneOf validates that a value is one of the allowed values
func ValidateOneOf(field string, value interface{}, allowedValues ...interface{}) error {
	for _, allowed := range allowedValues {
		if value == allowed {
			return nil
		}
	}
	return &FieldError{Field: field, Tag: "oneof", Param: fmt.Sprintf("%v", allowedValues), Message: fmt.Sprintf("must be one of: %v", allowedValues)}
}

// ValidatePhone validates phone number format (basic validation)
func ValidatePhone(field, phone string) error {
	if phone == "" {
		return &FieldError{Field: field, Tag: "required", Message: "phone is required"}
	}

	// Remove all non-digit characters for validation
	digitsOnly := strings.ReplaceAll(phone, " ", "")
	digitsOnly = strings.ReplaceAll(digitsOnly, "-", "")
	digitsOnly = strings.ReplaceAll(digitsOnly, "(", "")
	digitsOnly = strings.ReplaceAll(digitsOnly, ")", "")
	digitsOnly = strings.ReplaceAll(digitsOnly, "+", "")

	if len(digitsOnly) < 10 || len(digitsOnly) > 15 {
		return &FieldError{Field: field, Tag: "len", Param: "10,15", Message: "phone number must be between 10-15 digits"}
	}
	return nil
}
