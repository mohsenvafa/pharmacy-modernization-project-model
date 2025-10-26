package request

// PatientPathVars represents path parameters for patient endpoints
type PatientPathVars struct {
	PatientID string `path:"patientID" validate:"required,min=1"`
}
