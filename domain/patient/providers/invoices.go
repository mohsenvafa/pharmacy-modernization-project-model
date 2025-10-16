package providers

import (
	"context"

	irisbilling "pharmacy-modernization-project-model/internal/integrations/iris_billing"
)

type PatientInvoiceProvider interface {
	GetInvoicesByPatientID(ctx context.Context, patientID string) (*irisbilling.InvoiceListResponse, error)
}
