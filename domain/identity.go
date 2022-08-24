package domain

import "time"

// Identity is a user identity.
type Identity struct {
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	DOB       time.Time `json:"dob,string"`
}
