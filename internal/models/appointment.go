// File: internal/models/appointment.go
package models

import (
    "encoding/json"
    "time"

)

// CustomTime is a wrapper around time.Time that formats only the time portion
type CustomTime time.Time

// MarshalJSON implements the json.Marshaler interface
func (ct CustomTime) MarshalJSON() ([]byte, error) {
    t := time.Time(ct)
    return json.Marshal(t.Format("15:04:05"))
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
    var timeStr string
    if err := json.Unmarshal(data, &timeStr); err != nil {
        return err
    }

    t, err := time.Parse("15:04:05", timeStr)
    if err != nil {
        return err
    }

    *ct = CustomTime(t)
    return nil
}

// Time converts CustomTime back to time.Time
func (ct CustomTime) Time() time.Time {
    return time.Time(ct)
}

// IsZero reports whether t represents the zero time instant
func (ct CustomTime) IsZero() bool {
    return time.Time(ct).IsZero()
}

// Before reports whether the time instant ct is before u
func (ct CustomTime) Before(u CustomTime) bool {
    return time.Time(ct).Before(time.Time(u))
}

// Appointment represents an appointment entity
type Appointment struct {
	ID                  int       `json:"id" db:"id"`
	PatientID           int       `json:"patient_id" db:"patient_id"`
	ProviderType        string    `json:"provider_type" db:"provider_type"`
	DoctorID            *int      `json:"doctor_id,omitempty" db:"doctor_id"` // Optional
	HomeCareProviderID  *int      `json:"home_care_provider_id,omitempty" db:"home_care_provider_id"` // Optional
	AppointmentDate     time.Time `json:"appointment_date" db:"appointment_date"`
	StartTime           CustomTime `json:"start_time" db:"start_time"`
	EndTime             CustomTime `json:"end_time" db:"end_time"`
	Status              string    `json:"status" db:"status"`
	CancellationReason  *string   `json:"cancellation_reason,omitempty" db:"cancellation_reason"` // Optional
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}