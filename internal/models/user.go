package models

import (
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Password     string    `json:"password"` 
	PasswordHash string    `json:"-"`
	Name         string    `json:"name"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Role string

const (
	RoleDoctor            Role = "doctor"
	RolePatient           Role = "patient"
	RoleHomeCareProvider  Role = "home_care_provider"
	RoleAdmin             Role = "admin"
)
