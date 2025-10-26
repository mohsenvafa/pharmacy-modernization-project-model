package iris_billing

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// MockClient implements BillingClient with in-memory mock data
type MockClient struct {
	invoices          map[string]InvoiceResponse
	invoicesByPatient map[string][]InvoiceResponse
	payments          map[string]InvoicePaymentResponse
	logger            *zap.Logger
}

// NewMockClient creates a new mock billing client
func NewMockClient(logger *zap.Logger) *MockClient {
	return &MockClient{
		invoices:          make(map[string]InvoiceResponse),
		invoicesByPatient: make(map[string][]InvoiceResponse),
		payments:          make(map[string]InvoicePaymentResponse),
		logger:            logger,
	}
}

// SeedInvoice adds a mock invoice (useful for testing)
func (c *MockClient) SeedInvoice(invoice InvoiceResponse) {
	c.invoices[invoice.PrescriptionID] = invoice
}

// SeedPayment adds a mock payment (useful for testing)
func (c *MockClient) SeedPayment(payment InvoicePaymentResponse) {
	c.payments[payment.InvoiceID] = payment
}

// GetInvoice retrieves a mock invoice for a given prescription ID
func (c *MockClient) GetInvoice(ctx context.Context, prescriptionID string) (*InvoiceResponse, error) {
	if invoice, ok := c.invoices[prescriptionID]; ok {
		c.logger.Debug("mock invoice found",
			zap.String("invoice_id", invoice.ID),
			zap.String("status", invoice.Status),
		)
		return &invoice, nil
	}

	c.logger.Warn("mock invoice not found, returning default")

	// Return a default invoice for unknown prescription IDs
	defaultInvoice := &InvoiceResponse{
		PrescriptionID: prescriptionID,
		Status:         "unbilled",
		Amount:         0.0,
	}

	return defaultInvoice, nil
}

// GetInvoicesByPatientID retrieves all mock invoices for a given patient ID
func (c *MockClient) GetInvoicesByPatientID(ctx context.Context, patientID string) (*InvoiceListResponse, error) {
	if invoices, ok := c.invoicesByPatient[patientID]; ok {
		c.logger.Debug("mock invoices found for patient",
			zap.Int("count", len(invoices)),
		)
		return &InvoiceListResponse{
			PatientID: patientID,
			Invoices:  invoices,
			Total:     len(invoices),
		}, nil
	}

	c.logger.Warn("no mock invoices found for patient, returning empty list")

	// Return an empty list for unknown patient IDs
	return &InvoiceListResponse{
		PatientID: patientID,
		Invoices:  []InvoiceResponse{},
		Total:     0,
	}, nil
}

// CreateInvoice creates a mock invoice
func (c *MockClient) CreateInvoice(ctx context.Context, req CreateInvoiceRequest) (*CreateInvoiceResponse, error) {
	invoice := InvoiceResponse{
		ID:             fmt.Sprintf("mock-invoice-%s", req.PrescriptionID),
		PrescriptionID: req.PrescriptionID,
		Amount:         req.Amount,
		Status:         "pending",
	}

	c.invoices[req.PrescriptionID] = invoice

	c.logger.Debug("mock invoice created",
		zap.String("invoice_id", invoice.ID),
		zap.Float64("amount", req.Amount),
	)

	response := &CreateInvoiceResponse{
		InvoiceResponse: invoice,
	}

	return response, nil
}

// AcknowledgeInvoice acknowledges a mock invoice
func (c *MockClient) AcknowledgeInvoice(ctx context.Context, invoiceID string, req AcknowledgeInvoiceRequest) (*AcknowledgeInvoiceResponse, error) {
	// Find invoice by ID
	for prescID, invoice := range c.invoices {
		if invoice.ID == invoiceID {
			invoice.Status = "acknowledged"
			c.invoices[prescID] = invoice

			c.logger.Debug("mock invoice acknowledged",
				zap.String("invoice_id", invoiceID),
			)

			response := &AcknowledgeInvoiceResponse{
				InvoiceResponse: invoice,
			}

			return response, nil
		}
	}

	c.logger.Warn("mock invoice not found for acknowledgement")

	return nil, fmt.Errorf("invoice not found")
}

// GetInvoicePayment retrieves mock payment details
func (c *MockClient) GetInvoicePayment(ctx context.Context, invoiceID string) (*InvoicePaymentResponse, error) {
	if payment, ok := c.payments[invoiceID]; ok {
		c.logger.Debug("mock payment found",
			zap.String("payment_id", payment.PaymentID),
		)
		return &payment, nil
	}

	c.logger.Warn("mock payment not found, returning default")

	// Return default payment
	defaultPayment := &InvoicePaymentResponse{
		InvoiceID:     invoiceID,
		PaymentID:     fmt.Sprintf("mock-payment-%s", invoiceID),
		Status:        "pending",
		PaymentMethod: "mock",
	}

	return defaultPayment, nil
}

// Verify MockClient implements BillingClient
var _ BillingClient = (*MockClient)(nil)
