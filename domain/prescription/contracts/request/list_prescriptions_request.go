package request

// PrescriptionListQueryRequest represents filters accepted by the prescriptions listing endpoint.
type PrescriptionListQueryRequest struct {
	Status string `form:"status" validate:"omitempty,oneof=Active Pending Completed Cancelled"`
	Limit  int    `form:"limit" validate:"omitempty,min=1,max=100"`
	Offset int    `form:"offset" validate:"omitempty,min=0"`
}
