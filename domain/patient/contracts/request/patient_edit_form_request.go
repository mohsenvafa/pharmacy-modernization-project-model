package request

// PatientEditFormRequest represents form data for editing a patient
type PatientEditFormRequest struct {
	Name  string `form:"name" validate:"required,min=2,max=100"`
	Phone string `form:"phone" validate:"required,min=10,max=15"`
	DOB   string `form:"dob" validate:"required,dob"`
	State string `form:"state" validate:"required,min=1,max=50"`
}
