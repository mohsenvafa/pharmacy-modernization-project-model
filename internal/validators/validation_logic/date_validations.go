// internal/validators/validation_logic/date_validations.go
package validation_logic

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// ValidateDOB validates date of birth fields
func ValidateDOB(fl validator.FieldLevel) bool {
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
