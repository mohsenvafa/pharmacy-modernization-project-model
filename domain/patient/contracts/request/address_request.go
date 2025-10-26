package request

type AddressCreateRequest struct {
	Line1 string `json:"line1" validate:"required,min=1,max=100"`
	Line2 string `json:"line2" validate:"omitempty,max=100"`
	City  string `json:"city" validate:"required,min=1,max=50"`
	State string `json:"state" validate:"required,min=2,max=2"`
	Zip   string `json:"zip" validate:"required,len=5,numeric"`
}
