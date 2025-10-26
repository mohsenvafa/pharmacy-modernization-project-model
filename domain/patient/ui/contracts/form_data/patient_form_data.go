package form_data

// PatientFormData represents the form data for patient validation
type PatientFormData struct {
	ID     string
	Name   string
	Phone  string
	DOB    string // Format: YYYY-MM-DD
	State  string
	Errors map[string]string
}
