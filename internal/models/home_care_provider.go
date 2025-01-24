// File: internal/models/home_care_provider.go
package models

import "time"

type HomeCareProvider struct {
	UserID            int       `json:"user_id" db:"user_id"`
	ServiceTypeID     int       `json:"service_type_id" db:"service_type_id"`
	ExperienceYears   int       `json:"experience_years" db:"experience_years"`
	Qualifications    string    `json:"qualifications" db:"qualifications"`
	Bio               string    `json:"bio" db:"bio"`
	ProfilePictureURL string    `json:"profile_picture_url" db:"profile_picture_url"`
	HourlyRate        float64   `json:"hourly_rate" db:"hourly_rate"`
	Rating            float64   `json:"rating" db:"rating"`
	IsVerified        bool      `json:"is_verified" db:"is_verified"`
	IsAvailable       bool      `json:"is_available" db:"is_available"`
	Status            string    `json:"status" db:"status"`
	Latitude          float64   `json:"latitude" db:"latitude"`
	Longitude         float64   `json:"longitude" db:"longitude"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}



// Add this new struct
type HomeCareProviderFilter struct {
	ServiceTypeID   int     `json:"service_type_id"`
	MinRating       float64 `json:"min_rating"`
	MaxHourlyRate   float64 `json:"max_hourly_rate"`
}