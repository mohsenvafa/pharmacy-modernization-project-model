package iris_billing

import "context"

// BillingClient defines the interface for interacting with the IRIS billing API
type BillingClient interface {
	// Invoice operations
	GetInvoice(ctx context.Context, prescriptionID string) (*InvoiceResponse, error)
	CreateInvoice(ctx context.Context, req CreateInvoiceRequest) (*CreateInvoiceResponse, error)
	AcknowledgeInvoice(ctx context.Context, invoiceID string, req AcknowledgeInvoiceRequest) (*AcknowledgeInvoiceResponse, error)

	// Payment operations
	GetInvoicePayment(ctx context.Context, invoiceID string) (*InvoicePaymentResponse, error)
}
