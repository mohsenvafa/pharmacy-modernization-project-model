package helper

import (
	"net/http"

	"pharmacy-modernization-project-model/internal/bind"
)

// UIFormError represents a form validation error for UI
type UIFormError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ConvertFieldErrorsToUIErrors converts bind.FieldError to UIFormError
func ConvertFieldErrorsToUIErrors(fieldErrors []bind.FieldError) map[string]string {
	errors := make(map[string]string)
	for _, fe := range fieldErrors {
		errors[fe.Field] = getUIErrorMessage(fe)
	}
	return errors
}

// getUIErrorMessage converts validation tags to user-friendly messages
func getUIErrorMessage(fe bind.FieldError) string {
	switch fe.Tag {
	case "required":
		return "This field is required"
	case "min":
		if fe.Param != "" {
			return "Value must be at least " + fe.Param + " characters long"
		}
		return "Value is too short"
	case "max":
		if fe.Param != "" {
			return "Value must be no more than " + fe.Param + " characters long"
		}
		return "Value is too long"
	case "email":
		return "Please enter a valid email address"
	case "numeric":
		return "Please enter a valid number"
	case "len":
		if fe.Param != "" {
			return "Value must be exactly " + fe.Param + " characters"
		}
		return "Value has incorrect length"
	case "oneof":
		if fe.Param != "" {
			return "Value must be one of: " + fe.Param
		}
		return "Value is not in the allowed list"
	case "dob":
		return "Please enter a valid date of birth (YYYY-MM-DD). Date cannot be in the future or more than 150 years ago."
	default:
		return "Invalid value"
	}
}

// WriteUIError writes an error response for UI handlers
func WriteUIError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write([]byte(`<html><body><h1>Error</h1><p>` + message + `</p></body></html>`))
}

// WriteUINotFound writes a 404 response for UI handlers
func WriteUINotFound(w http.ResponseWriter, message string) {
	WriteUIError(w, message, http.StatusNotFound)
}

// WriteUIInternalError writes a 500 response for UI handlers
func WriteUIInternalError(w http.ResponseWriter, message string) {
	WriteUIError(w, message, http.StatusInternalServerError)
}
