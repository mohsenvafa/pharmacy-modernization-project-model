package model

import "time"

type Status string

const (
	Draft     Status = "Draft"
	Active    Status = "Active"
	Paused    Status = "Paused"
	Completed Status = "Completed"
)

type Prescription struct {
	ID        string    `json:"id" bson:"_id"`
	PatientID string    `json:"patient_id" bson:"patient_id"`
	Drug      string    `json:"drug" bson:"drug"`
	Dose      string    `json:"dose" bson:"dose"`
	Status    Status    `json:"status" bson:"status"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
