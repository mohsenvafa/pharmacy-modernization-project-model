// internal/validators/custom_validators.go
package validators

import (
	"pharmacy-modernization-project-model/internal/validators/validation_logic"

	"github.com/go-playground/validator/v10"
)

// RegisterCustomValidators registers all custom validation rules
func RegisterCustomValidators(validate *validator.Validate) {
	// Register DOB validator using validation logic
	validate.RegisterValidation("dob", validation_logic.ValidateDOB)
}
