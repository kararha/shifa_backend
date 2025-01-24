// File: internal/models/review.go
package models

import "time"

type Review struct {
	ID                  int       `json:"id" db:"id"`
	PatientID           int       `json:"patient_id" db:"patient_id"`
	ReviewType          string    `json:"review_type" db:"review_type"`
	ConsultationID      *int      `json:"consultation_id,omitempty" db:"consultation_id"`
	HomeCareVisitID     *int      `json:"home_care_visit_id,omitempty" db:"home_care_visit_id"`
	DoctorID            *int      `json:"doctor_id,omitempty" db:"doctor_id"`
	HomeCareProviderID  *int      `json:"home_care_provider_id,omitempty" db:"home_care_provider_id"`
	Rating              int       `json:"rating" db:"rating"`
	Comment             string    `json:"comment" db:"comment"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}