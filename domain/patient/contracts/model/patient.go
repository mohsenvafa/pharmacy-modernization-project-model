package model

import "time"

type Patient struct {
	ID        string    `json:"id" bson:"_id"`
	Name      string    `json:"name" bson:"name"`
	DOB       time.Time `json:"dob" bson:"dob"`
	Phone     string    `json:"phone" bson:"phone"`
	State     string    `json:"state" bson:"state"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
