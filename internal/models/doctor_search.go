package models

// DoctorSearchParams represents search parameters for doctors
type DoctorSearchParams struct {
    SearchTerm    string   `json:"search_term"`
    Specialty     string   `json:"specialty"`
    MinRating     *float64 `json:"min_rating"`
    LocationLat   *float64 `json:"location_lat"`
    LocationLng   *float64 `json:"location_lng"`
    RadiusKm      *int     `json:"radius_km"`
    ServiceTypeID *int     `json:"service_type_id"`
}

// DoctorSearchResult represents the result of a doctor search
type DoctorSearchResult struct {
    UserID            int     `json:"user_id"`
    FirstName         string  `json:"first_name"`
    LastName          string  `json:"last_name"`
    Email            string  `json:"email"`
    Phone            string  `json:"phone"`
    Specialty        string  `json:"specialty"`
    ServiceTypeID    int     `json:"service_type_id"`
    Rating           float64 `json:"rating"`
    DistanceKm       *float64 `json:"distance_km,omitempty"`
    IsAvailable      bool    `json:"is_available"`
    Status           string  `json:"status"`
    // Add other relevant fields
}
