// File: internal/models/patient.go
package models

import "time"

type Patient struct {
	UserID                int       `json:"user_id" db:"user_id"`
	Name                  string    `json:"name" db:"name"`
	DateOfBirth           time.Time `json:"date_of_birth" db:"date_of_birth"`
	Gender                string    `json:"gender" db:"gender"`
	Phone                 string    `json:"phone" db:"phone"`
	Address               string    `json:"address" db:"address"`
	EmergencyContactName  string    `json:"emergency_contact_name" db:"emergency_contact_name"`
	EmergencyContactPhone string    `json:"emergency_contact_phone" db:"emergency_contact_phone"`
}