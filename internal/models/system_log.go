package models

import (
    "time"
)

// JSON type definition (assuming you have a custom JSON type defined)


type SystemLog struct {
    ID                 int64     `json:"id"`
    Timestamp          time.Time `json:"timestamp"`
    UserID             int       `json:"user_id"`
    UserType           string    `json:"user_type"`
    ActionType         string    `json:"action_type"`
    ActionDescription  string    `json:"action_description"`
    EntityType         string    `json:"entity_type"`
    EntityID           int       `json:"entity_id"`
    OldValue           JSON      `json:"old_value"`
    NewValue           JSON      `json:"new_value"`
    IPAddress          string    `json:"ip_address"`
    UserAgent          string    `json:"user_agent"`
    AdditionalInfo     JSON      `json:"additional_info"`
}

// SystemLogFilter represents the filter options for system logs
type SystemLogFilter struct {
    UserID     int64
    UserType   string
    ActionType string
    EntityType string
    StartDate  time.Time
    EndDate    time.Time
}