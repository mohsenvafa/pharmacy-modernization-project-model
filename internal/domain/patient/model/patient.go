package model

import "time"

type Patient struct {
	ID    string
	Name  string
	DOB   time.Time
	Phone string
	CreatedAt time.Time
}
