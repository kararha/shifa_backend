// File: internal/models/doctor_availability.go
package models

import "time"

type DoctorAvailability struct {
	ID        int       `json:"id" db:"id"`
	DoctorID  int       `json:"doctor_id" db:"doctor_id"`
	DayOfWeek int       `json:"day_of_week" db:"day_of_week"`
	StartTime time.Time `json:"start_time" db:"start_time"`
	EndTime   time.Time `json:"end_time" db:"end_time"`
}