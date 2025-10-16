package iris_billing

// InvoiceResponse represents a billing invoice from IRIS
type InvoiceResponse struct {
	ID             string  `json:"id"`
	PrescriptionID string  `json:"prescription_id"`
	Amount         float64 `json:"amount"`
	Status         string  `json:"status"`
	CreatedAt      string  `json:"created_at,omitempty"`
	UpdatedAt      string  `json:"updated_at,omitempty"`
}

// CreateInvoiceRequest represents a request to create an invoice
type CreateInvoiceRequest struct {
	PrescriptionID string  `json:"prescription_id"`
	Amount         float64 `json:"amount"`
	Description    string  `json:"description,omitempty"`
}

// CreateInvoiceResponse represents the response from creating an invoice
type CreateInvoiceResponse struct {
	InvoiceResponse
}

// AcknowledgeInvoiceRequest represents a request to acknowledge an invoice
type AcknowledgeInvoiceRequest struct {
	AcknowledgedBy string `json:"acknowledged_by"`
	Notes          string `json:"notes,omitempty"`
}

// AcknowledgeInvoiceResponse represents the response from acknowledging an invoice
type AcknowledgeInvoiceResponse struct {
	InvoiceResponse
}

// InvoicePaymentResponse represents payment details for an invoice
type InvoicePaymentResponse struct {
	InvoiceID     string  `json:"invoice_id"`
	PaymentID     string  `json:"payment_id"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"payment_method"`
	Status        string  `json:"status"`
	PaidAt        string  `json:"paid_at,omitempty"`
}

// InvoiceListResponse represents a list of invoices for a patient
type InvoiceListResponse struct {
	PatientID string            `json:"patient_id"`
	Invoices  []InvoiceResponse `json:"invoices"`
	Total     int               `json:"total"`
}
