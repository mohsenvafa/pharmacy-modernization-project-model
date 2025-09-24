package iris_pharmacy

import "context"

type Client interface {
	GetPrescription(ctx context.Context, prescriptionID string) (GetPrescriptionResponse, error)
}
