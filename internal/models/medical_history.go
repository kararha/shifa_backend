// File: internal/models/medical_history.go
package models

import "time"

type MedicalHistory struct {
	ID             int       `json:"id" db:"id"`
	PatientID      int       `json:"patient_id" db:"patient_id"`
	ConditionName  string    `json:"condition_name" db:"condition_name"`
	DiagnosisDate  time.Time `json:"diagnosis_date" db:"diagnosis_date"`
	Treatment      string    `json:"treatment" db:"treatment"`
	IsCurrent      bool      `json:"is_current" db:"is_current"`
}