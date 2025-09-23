package model

import "time"

type Patient struct {
	ID        string
	Name      string
	DOB       time.Time
	Phone     string
	State     string
	CreatedAt time.Time
}
