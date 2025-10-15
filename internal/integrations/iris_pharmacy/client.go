package iris_pharmacy

import "context"

// PharmacyClient defines the interface for interacting with the IRIS pharmacy API
type PharmacyClient interface {
	GetPrescription(ctx context.Context, prescriptionID string) (*PrescriptionResponse, error)
}
