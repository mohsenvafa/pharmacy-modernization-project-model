package request

// PrescriptionListQueryRequest represents filters accepted by the prescriptions listing endpoint.
type PrescriptionListQueryRequest struct {
	Status string `json:"status"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}
