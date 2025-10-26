package request

// PrescriptionPathVars represents path parameters for prescription endpoints
type PrescriptionPathVars struct {
	PrescriptionID string `path:"prescriptionID" validate:"required,min=1"`
}
