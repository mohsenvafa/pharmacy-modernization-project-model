package request

// PatientListPageRequest represents query parameters for patient list page
type PatientListPageRequest struct {
	Page        int    `form:"page" validate:"omitempty,min=1"`
	PatientName string `form:"patientName" validate:"omitempty,min=1"`
	BirthDate   string `form:"birthDate" validate:"omitempty"`
	State       string `form:"state" validate:"omitempty,min=1"`
}

// PatientComponentRequest represents query parameters for patient-related components
type PatientComponentRequest struct {
	PatientID string `form:"patientId" validate:"required,min=1"`
}
