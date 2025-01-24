package models

import (
	"time"
)

type Doctor struct {
	UserID            int       `json:"user_id"`
	Specialty         string    `json:"specialty"`
	ServiceTypeID     int       `json:"service_type_id"`
	LicenseNumber     string    `json:"license_number"`
	ExperienceYears   int       `json:"experience_years"`
	Qualifications    string    `json:"qualifications"`
	Achievements      string    `json:"achievements"`
	Bio               string    `json:"bio"`
	ProfilePictureURL string    `json:"profile_picture_url"`
	ConsultationFee   float64   `json:"consultation_fee"`
	Rating            float64   `json:"rating"`
	IsVerified        bool      `json:"is_verified"`
	IsAvailable       bool      `json:"is_available"`
	Status            string    `json:"status"`
	Latitude          float64   `json:"latitude"`
	Longitude         float64   `json:"longitude"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}


// Add this new struct
type DoctorFilter struct {
	Specialty     string  `json:"specialty"`
	ServiceTypeID int     `json:"service_type_id"`
	MinRating     float64 `json:"min_rating"`
}