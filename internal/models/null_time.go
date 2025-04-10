package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

// NullTime represents a time.Time that can be NULL in the database
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface for NullTime
func (nt *NullTime) Scan(value interface{}) error {
	var nullTime sql.NullTime
	if err := nullTime.Scan(value); err != nil {
		return err
	}

	nt.Time = nullTime.Time
	nt.Valid = nullTime.Valid
	return nil
}

// Value implements the driver.Valuer interface
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// MarshalJSON implements the json.Marshaler interface
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nt.Time)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (nt *NullTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		nt.Valid = false
		return nil
	}

	var t time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	nt.Time = t
	nt.Valid = true
	return nil
}
