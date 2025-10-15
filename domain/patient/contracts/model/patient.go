package model

import "time"

type Patient struct {
	ID        string     `json:"id" bson:"_id"`
	Name      string     `json:"name" bson:"name"`
	DOB       time.Time  `json:"dob" bson:"dob"`
	Phone     string     `json:"phone" bson:"phone"`
	State     string     `json:"state" bson:"state"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	EditBy    *string    `json:"edit_by,omitempty" bson:"edit_by,omitempty"`
	EditTime  *time.Time `json:"edit_time,omitempty" bson:"edit_time,omitempty"`
}
