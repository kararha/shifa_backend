// File: internal/models/home_care_visit.go
package models

import "time"

type HomeCareVisit struct {
	ID                  int     `json:"id" db:"id"`
	AppointmentID       int     `json:"appointment_id" db:"appointment_id"`
	Address             string  `json:"address" db:"address"`
	Latitude            float64 `json:"latitude" db:"latitude"`
	Longitude           float64 `json:"longitude" db:"longitude"`
	DurationHours       float64 `json:"duration_hours" db:"duration_hours"`
	SpecialRequirements string  `json:"special_requirements" db:"special_requirements"`
	Status              string  `json:"status" db:"status"`
	PatientID  			int
    ProviderID		    int
    VisitDate 		    time.Time
}

type HomeCareVisitFilter struct {
    AppointmentID int       `json:"appointment_id"`
    PlanID        int       `json:"plan_id"`
    PatientID     int       `json:"patient_id"`     // Added
    ProviderID    int       `json:"provider_id"`    // Added
    StartDate     time.Time `json:"start_date"`
    EndDate       time.Time `json:"end_date"`
    Status        string    `json:"status"`
}