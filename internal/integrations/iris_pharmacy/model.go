package iris_pharmacy

type GetPrescriptionResponse struct {
	ID           string `json:"id"`
	PatientID    string `json:"patient_id"`
	Drug         string `json:"drug"`
	Dose         string `json:"dose"`
	Status       string `json:"status"`
	PharmacyName string `json:"pharmacy_name"`
	PharmacyType string `json:"pharmacy_type"`
}
