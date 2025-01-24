// File: internal/models/consultation.go
package models

import "time"

type Consultation struct {
    ID               int       `json:"id" db:"id"`
    PatientID        int       `json:"patient_id"`
    DoctorID         int       `json:"doctor_id"`
    AppointmentID    int       `json:"appointment_id" db:"appointment_id"`
    ConsultationType string    `json:"consultation_type" db:"consultation_type"`
    Status           string    `json:"status" db:"status"`
    StartedAt        time.Time `json:"started_at" db:"started_at"`
    CompletedAt      time.Time `json:"completed_at" db:"completed_at"`
    Fee              float64   `json:"fee" db:"fee"`
}

// ConsultationFilter represents the filtering options for consultations
type ConsultationFilter struct {
    PatientID         int
    DoctorID          int
    AppointmentID     int
    ConsultationType  string
    Status            string
    StartDateFrom     time.Time
    StartDateTo       time.Time
    CompletedDateFrom time.Time
    CompletedDateTo   time.Time
}