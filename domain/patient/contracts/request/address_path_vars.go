package request

// AddressPathVars represents path parameters for address endpoints
type AddressPathVars struct {
	PatientID string `path:"patientID" validate:"required,min=1"`
	AddressID string `path:"addressID" validate:"required,min=1"`
}
