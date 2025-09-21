package model

import "time"

type Status string
const (
	Draft Status = "Draft"
	Active Status = "Active"
	Paused Status = "Paused"
	Completed Status = "Completed"
)

type Prescription struct {
	ID string
	PatientID string
	Drug string
	Dose string
	Status Status
	CreatedAt time.Time
}
