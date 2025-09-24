package iris_billing

type GetInvoiceResponse struct {
	ID             string  `json:"id"`
	PrescriptionID string  `json:"prescription_id"`
	Amount         float64 `json:"amount"`
	Status         string  `json:"status"`
}
