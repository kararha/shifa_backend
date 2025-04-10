// File: internal/models/consultation.go
package models

import "time"

type Consultation struct {
	ID          int      `json:"id" db:"id"`
	PatientID   int      `json:"patient_id" db:"patient_id"`
	DoctorID    int      `json:"doctor_id" db:"doctor_id"`
	Status      string   `json:"status" db:"status"`
	StartedAt   NullTime `json:"started_at" db:"started_at"`
	CompletedAt NullTime `json:"completed_at" db:"completed_at"`
	Fee         float64  `json:"fee" db:"fee"`
}

// ConsultationFilter represents the filtering options for consultations
type ConsultationFilter struct {
	PatientID         int
	DoctorID          int
	Status            string
	StartDateFrom     time.Time
	StartDateTo       time.Time
	CompletedDateFrom time.Time
	CompletedDateTo   time.Time
}
