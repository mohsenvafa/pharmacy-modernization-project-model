package validation

import (
	"errors"
	"fmt"
	"strings"

	"pharmacy-modernization-project-model/internal/bind"
	"pharmacy-modernization-project-model/internal/graphql/generated"

	"github.com/go-playground/validator/v10"
)

// GraphQLValidationError represents a validation error for GraphQL responses
type GraphQLValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// GraphQLValidationErrors represents multiple validation errors
type GraphQLValidationErrors struct {
	Errors []GraphQLValidationError `json:"errors"`
}

// Error implements the error interface
func (e *GraphQLValidationErrors) Error() string {
	var messages []string
	for _, err := range e.Errors {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(messages, ", ")
}

// ConvertBindErrorsToGraphQLErrors converts bind.FieldError to GraphQLValidationError
func ConvertBindErrorsToGraphQLErrors(fieldErrors []bind.FieldError) *GraphQLValidationErrors {
	var errors []GraphQLValidationError
	for _, fe := range fieldErrors {
		errors = append(errors, GraphQLValidationError{
			Field:   fe.Field,
			Message: getGraphQLErrorMessage(fe),
		})
	}
	return &GraphQLValidationErrors{Errors: errors}
}

// getGraphQLErrorMessage converts validation tags to user-friendly messages
func getGraphQLErrorMessage(fe bind.FieldError) string {
	switch fe.Tag {
	case "required":
		return "This field is required"
	case "min":
		if fe.Param != "" {
			return fmt.Sprintf("Value must be at least %s", fe.Param)
		}
		return "Value is too short"
	case "max":
		if fe.Param != "" {
			return fmt.Sprintf("Value must be no more than %s", fe.Param)
		}
		return "Value is too long"
	case "email":
		return "Please enter a valid email address"
	case "numeric":
		return "Please enter a valid number"
	case "len":
		if fe.Param != "" {
			return fmt.Sprintf("Value must be exactly %s characters", fe.Param)
		}
		return "Value has incorrect length"
	case "oneof":
		if fe.Param != "" {
			return fmt.Sprintf("Value must be one of: %s", fe.Param)
		}
		return "Value is not in the allowed list"
	case "dob":
		return "Please enter a valid date of birth (YYYY-MM-DD). Date cannot be in the future or more than 150 years ago."
	case "alphanum":
		return "Value can only contain letters, numbers, hyphens, and underscores"
	default:
		return "Invalid value"
	}
}

// Validation structs for GraphQL input types
// These mirror the generated types but add validation tags

// CreatePatientInputValidation represents validated input for creating a patient
type CreatePatientInputValidation struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Dob   string `json:"dob" validate:"required,dob"`
	Phone string `json:"phone" validate:"required,min=10,max=15"`
	State string `json:"state" validate:"required,min=2,max=50"`
}

// UpdatePatientInputValidation represents validated input for updating a patient
type UpdatePatientInputValidation struct {
	Name  *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Dob   *string `json:"dob,omitempty" validate:"omitempty,dob"`
	Phone *string `json:"phone,omitempty" validate:"omitempty,min=10,max=15"`
	State *string `json:"state,omitempty" validate:"omitempty,min=2,max=50"`
}

// CreatePrescriptionInputValidation represents validated input for creating a prescription
type CreatePrescriptionInputValidation struct {
	PatientID string `json:"patientID" validate:"required,min=1,max=50,alphanum"`
	Drug      string `json:"drug" validate:"required,min=2,max=100"`
	Dose      string `json:"dose" validate:"required,min=1,max=50"`
	Status    string `json:"status" validate:"required,oneof=DRAFT ACTIVE PAUSED COMPLETED"`
}

// UpdatePrescriptionInputValidation represents validated input for updating a prescription
type UpdatePrescriptionInputValidation struct {
	Drug   *string `json:"drug,omitempty" validate:"omitempty,min=2,max=100"`
	Dose   *string `json:"dose,omitempty" validate:"omitempty,min=1,max=50"`
	Status *string `json:"status,omitempty" validate:"omitempty,oneof=DRAFT ACTIVE PAUSED COMPLETED"`
}

// PatientQueryValidation represents validated input for patient queries
type PatientQueryValidation struct {
	ID string `json:"id" validate:"required,min=1,max=50,alphanum"`
}

// PatientsQueryValidation represents validated input for patients list query
type PatientsQueryValidation struct {
	Query  *string `json:"query,omitempty" validate:"omitempty,min=3,max=100"`
	Limit  *int    `json:"limit,omitempty" validate:"omitempty,min=1,max=100"`
	Offset *int    `json:"offset,omitempty" validate:"omitempty,min=0"`
}

// PrescriptionQueryValidation represents validated input for prescription queries
type PrescriptionQueryValidation struct {
	ID string `json:"id" validate:"required,min=1,max=50,alphanum"`
}

// PrescriptionsQueryValidation represents validated input for prescriptions list query
type PrescriptionsQueryValidation struct {
	Status *string `json:"status,omitempty" validate:"omitempty,oneof=DRAFT ACTIVE PAUSED COMPLETED"`
	Limit  *int    `json:"limit,omitempty" validate:"omitempty,min=1,max=100"`
	Offset *int    `json:"offset,omitempty" validate:"omitempty,min=0"`
}

// ValidateGraphQLInput validates a GraphQL input using the bind validation system
func ValidateGraphQLInput[T any](input T) (T, *GraphQLValidationErrors) {
	// Use the bind validator directly
	validator := bind.Validator()
	if err := validator.Struct(input); err != nil {
		fieldErrors := convertValidatorErrors(err)
		return input, ConvertBindErrorsToGraphQLErrors(fieldErrors)
	}
	return input, nil
}

// convertValidatorErrors converts validator.ValidationErrors to bind.FieldError slice
func convertValidatorErrors(err error) []bind.FieldError {
	var ferrs []bind.FieldError
	if err == nil {
		return ferrs
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			ferrs = append(ferrs, bind.FieldError{
				Field: fe.Field(),
				Tag:   fe.Tag(),
				Param: fe.Param(),
			})
		}
		return ferrs
	}
	// generic
	return []bind.FieldError{{Field: "", Tag: "invalid", Message: err.Error()}}
}

// Helper function to convert generated input to validation struct
func ConvertCreatePatientInput(input generated.CreatePatientInput) CreatePatientInputValidation {
	return CreatePatientInputValidation{
		Name:  input.Name,
		Dob:   input.Dob.Format("2006-01-02"), // Convert time.Time to string for validation
		Phone: input.Phone,
		State: input.State,
	}
}

func ConvertUpdatePatientInput(input generated.UpdatePatientInput) UpdatePatientInputValidation {
	result := UpdatePatientInputValidation{}

	if input.Name != nil {
		result.Name = input.Name
	}
	if input.Dob != nil {
		dobStr := input.Dob.Format("2006-01-02")
		result.Dob = &dobStr
	}
	if input.Phone != nil {
		result.Phone = input.Phone
	}
	if input.State != nil {
		result.State = input.State
	}

	return result
}

func ConvertCreatePrescriptionInput(input generated.CreatePrescriptionInput) CreatePrescriptionInputValidation {
	return CreatePrescriptionInputValidation{
		PatientID: input.PatientID,
		Drug:      input.Drug,
		Dose:      input.Dose,
		Status:    string(input.Status),
	}
}

func ConvertUpdatePrescriptionInput(input generated.UpdatePrescriptionInput) UpdatePrescriptionInputValidation {
	result := UpdatePrescriptionInputValidation{}

	if input.Drug != nil {
		result.Drug = input.Drug
	}
	if input.Dose != nil {
		result.Dose = input.Dose
	}
	if input.Status != nil {
		statusStr := string(*input.Status)
		result.Status = &statusStr
	}

	return result
}
