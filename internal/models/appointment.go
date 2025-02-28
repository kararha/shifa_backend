// File: internal/models/appointment.go
package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
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

	// Try parsing as ISO 8601 first
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		// Fallback to time-only format
		t, err = time.Parse("15:04:05", timeStr)
		if err != nil {
			return err
		}
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

// Add the following methods to handle SQL value conversion
func (ct CustomTime) Value() (driver.Value, error) {
	return time.Time(ct).Format("15:04:05"), nil
}

func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*ct = CustomTime(v)
		return nil
	case string:
		t, err := time.Parse("15:04:05", v)
		if err != nil {
			return err
		}
		*ct = CustomTime(t)
		return nil
	}
	return fmt.Errorf("cannot scan %T into CustomTime", value)
}

// Appointment represents an appointment entity
type Appointment struct {
	ID                 int        `json:"id"`
	PatientID          int        `json:"patient_id"`
	DoctorID           *int       `json:"doctor_id,omitempty"`
	HomeCareProviderID *int       `json:"home_care_provider_id,omitempty"`
	ServiceTypeID      int        `json:"service_type_id"`
	AppointmentDate    time.Time  `json:"appointment_date"`
	StartTime          CustomTime `json:"start_time"`
	EndTime            CustomTime `json:"end_time"`
	Status             string     `json:"status"`
	CancellationReason *string    `json:"cancellation_reason,omitempty"`
	ProviderType       string     `json:"provider_type"`
	PatientName        string     `json:"patient_name,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}
