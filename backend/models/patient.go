package models

type Patient struct {
	ID        string `json:"id" bson:"_id,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	DOB       string `json:"dob"`
	Gender    string `json:"gender"`
}
