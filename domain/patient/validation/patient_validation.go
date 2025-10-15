package validation

import (
	"regexp"
	"time"

	"pharmacy-modernization-project-model/domain/patient/contracts/model"
	tools "pharmacy-modernization-project-model/internal/helper"
)

// PatientFormData represents the form data for patient validation
type PatientFormData struct {
	ID     string
	Name   string
	Phone  string
	DOB    string // Format: YYYY-MM-DD
	State  string
	Errors map[string]string
}

// ValidatePatientForm validates patient form data and returns validation errors
func ValidatePatientForm(formData PatientFormData) map[string]string {
	errors := make(map[string]string)

	// Validate name
	if formData.Name == "" {
		errors["name"] = "Name is required"
	} else if len(formData.Name) < 2 {
		errors["name"] = "Name must be at least 2 characters"
	} else if len(formData.Name) > 100 {
		errors["name"] = "Name must be less than 100 characters"
	}

	// Validate phone
	if formData.Phone == "" {
		errors["phone"] = "Phone number is required"
	} else if !isValidPhone(formData.Phone) {
		errors["phone"] = "Please enter a valid phone number (e.g., (555) 123-4567)"
	}

	// Validate DOB
	if formData.DOB == "" {
		errors["dob"] = "Date of birth is required"
	} else {
		dob, err := time.Parse("2006-01-02", formData.DOB)
		if err != nil {
			errors["dob"] = "Please enter a valid date"
		} else {
			// Check if DOB is not in the future
			if dob.After(time.Now()) {
				errors["dob"] = "Date of birth cannot be in the future"
			}
			// Check if patient is not too old (reasonable limit)
			age := tools.CalculateAge(dob)
			if age > 150 {
				errors["dob"] = "Please enter a valid date of birth"
			}
		}
	}

	// Validate state
	if formData.State == "" {
		errors["state"] = "State is required"
	} else if len(formData.State) > 50 {
		errors["state"] = "State must be less than 50 characters"
	}

	return errors
}

// ValidatePatient validates a Patient model and returns validation errors
func ValidatePatient(patient model.Patient) map[string]string {
	errors := make(map[string]string)

	// Validate name
	if patient.Name == "" {
		errors["name"] = "Name is required"
	} else if len(patient.Name) < 2 {
		errors["name"] = "Name must be at least 2 characters"
	} else if len(patient.Name) > 100 {
		errors["name"] = "Name must be less than 100 characters"
	}

	// Validate phone
	if patient.Phone == "" {
		errors["phone"] = "Phone number is required"
	} else if !isValidPhone(patient.Phone) {
		errors["phone"] = "Please enter a valid phone number (e.g., (555) 123-4567)"
	}

	// Validate DOB
	if patient.DOB.IsZero() {
		errors["dob"] = "Date of birth is required"
	} else {
		// Check if DOB is not in the future
		if patient.DOB.After(time.Now()) {
			errors["dob"] = "Date of birth cannot be in the future"
		}
		// Check if patient is not too old (reasonable limit)
		age := tools.CalculateAge(patient.DOB)
		if age > 150 {
			errors["dob"] = "Please enter a valid date of birth"
		}
	}

	// Validate state
	if patient.State == "" {
		errors["state"] = "State is required"
	} else if len(patient.State) > 50 {
		errors["state"] = "State must be less than 50 characters"
	}

	return errors
}

// isValidPhone validates phone number format
func isValidPhone(phone string) bool {
	// Simple phone validation - allows various formats
	// Remove all non-digit characters
	digits := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")

	// Check if we have 10 digits (US phone number)
	if len(digits) == 10 {
		return true
	}

	// Check if we have 11 digits starting with 1 (US phone number with country code)
	if len(digits) == 11 && digits[0] == '1' {
		return true
	}

	return false
}
