// File: internal/models/consultation_details.go
package models

type ConsultationDetails struct {
	ID              int     `json:"id" db:"id"`
	ConsultationID  int     `json:"consultation_id" db:"consultation_id"`
	RequestDetails  string  `json:"request_details" db:"request_details"`
	Symptoms        string  `json:"symptoms" db:"symptoms"`
	Diagnosis       string  `json:"diagnosis" db:"diagnosis"`
	Prescription    string  `json:"prescription" db:"prescription"`
	Notes           string  `json:"notes" db:"notes"`
}
