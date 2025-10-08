package model

type Address struct {
	ID        string `json:"id" bson:"_id"`
	PatientID string `json:"patient_id" bson:"patient_id"`
	Line1     string `json:"line1" bson:"line1"`
	Line2     string `json:"line2" bson:"line2"`
	City      string `json:"city" bson:"city"`
	State     string `json:"state" bson:"state"`
	Zip       string `json:"zip" bson:"zip"`
}
