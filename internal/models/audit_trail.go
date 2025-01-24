// File: models/audit_trail.go
package models

import (
    "database/sql/driver"
    "encoding/json"
    "errors" 
    "time"
)

type AuditTrail struct {
    ID            int       `json:"id"`
    TableName     string    `json:"table_name"`
    RecordID      int       `json:"record_id"`
    Action        string    `json:"action"`
    ChangedFields JSON      `json:"changed_fields"`
    ChangedBy     int       `json:"changed_by"`
    ChangedAt     time.Time `json:"changed_at"`
}

// JSON is a custom type to handle JSON in the database
type JSON map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
    return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
    b, ok := value.([]byte)
    if !ok {
        return errors.New("type assertion to []byte failed")
    }
    return json.Unmarshal(b, &j)
}
