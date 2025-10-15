package iris_pharmacy

// PrescriptionResponse represents a prescription from IRIS pharmacy system
type PrescriptionResponse struct {
	ID           string `json:"id"`
	PatientID    string `json:"patient_id"`
	Drug         string `json:"drug"`
	Dose         string `json:"dose"`
	Status       string `json:"status"`
	PharmacyName string `json:"pharmacy_name"`
	PharmacyType string `json:"pharmacy_type"`
}
