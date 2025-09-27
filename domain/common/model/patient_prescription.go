package model

import "time"

type PatientPrescription struct {
	ID        string
	Drug      string
	Dose      string
	Status    string
	CreatedAt time.Time
}
