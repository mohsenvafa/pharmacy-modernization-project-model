// internal/validators/custom_validators.go
package validators

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// RegisterCustomValidators registers all custom validation rules
func RegisterCustomValidators(validate *validator.Validate) {
	// Register DOB validator
	validate.RegisterValidation("dob", validateDOB)
}

// validateDOB validates date of birth fields
func validateDOB(fl validator.FieldLevel) bool {
	dobStr := fl.Field().String()

	// Check if empty (handled by required tag)
	if dobStr == "" {
		return true
	}

	// Parse the date
	dob, err := time.Parse("2006-01-02", dobStr)
	if err != nil {
		return false
	}

	// Check if DOB is not in the future
	if dob.After(time.Now()) {
		return false
	}

	// Check if patient is not too old (reasonable limit)
	age := time.Now().Year() - dob.Year()
	return age <= 150
}
