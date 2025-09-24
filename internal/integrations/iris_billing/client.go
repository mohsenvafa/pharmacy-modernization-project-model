package iris_billing

import "context"

type Client interface {
	GetInvoice(ctx context.Context, prescriptionID string) (GetInvoiceResponse, error)
}
