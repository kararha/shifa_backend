// File: internal/models/service_type.go
package models

type ServiceType struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	IsHomeCare  bool   `json:"is_home_care" db:"is_home_care"`
}