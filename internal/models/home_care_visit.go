// File: internal/models/home_care_visit.go
package models

// import "time"

type HomeCareVisit struct {
    ID                  int     `json:"id" db:"id"`
    PatientID           int     `json:"patient_id" db:"patient_id"`
    ProviderID          int     `json:"provider_id" db:"provider_id"`
    Address             string  `json:"address" db:"address"`
    Latitude            float64 `json:"latitude" db:"latitude"`
    Longitude           float64 `json:"longitude" db:"longitude"`
    DurationHours       float64 `json:"duration_hours" db:"duration_hours"`
    SpecialRequirements string  `json:"special_requirements" db:"special_requirements"`
    Status              string  `json:"status" db:"status"`
}

type HomeCareVisitFilter struct {
    PatientID     int       `json:"patient_id"`
    ProviderID    int       `json:"provider_id"`
    Status        string    `json:"status"`
}