package request

type PatientListQueryRequest struct {
	Limit       int    `form:"limit" validate:"omitempty,min=1,max=100"`
	Offset      int    `form:"offset" validate:"omitempty,min=0"`
	PatientName string `form:"patientName" validate:"omitempty,min=3"`
	BirthDate   string `form:"birthDate" validate:"omitempty"`
	State       string `form:"state" validate:"omitempty,min=1"`
}
