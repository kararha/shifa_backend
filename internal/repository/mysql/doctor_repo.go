// File: internal/repository/mysql/doctor_repo.go

package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"shifa/internal/models"
)

// DoctorRepo represents the MySQL repository for doctor-related database operations
type DoctorRepo struct {
	db *sql.DB
}

// NewDoctorRepo creates a new DoctorRepo instance
func NewDoctorRepo(db *sql.DB) *DoctorRepo {
	return &DoctorRepo{db: db}
}

// Create inserts a new doctor into the database
func (r *DoctorRepo) Create(ctx context.Context, doctor *models.Doctor) error {
	query := `
		INSERT INTO doctors (user_id, specialty, service_type_id, license_number, experience_years, 
			qualifications, achievements, bio, profile_picture_url, consultation_fee, rating, 
			is_verified, is_available, status, latitude, longitude)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	_, err := r.db.ExecContext(ctx, query, 
		doctor.UserID, doctor.Specialty, doctor.ServiceTypeID, doctor.LicenseNumber,
		doctor.ExperienceYears, doctor.Qualifications, doctor.Achievements, doctor.Bio,
		doctor.ProfilePictureURL, doctor.ConsultationFee, doctor.Rating, doctor.IsVerified,
		doctor.IsAvailable, doctor.Status, doctor.Latitude, doctor.Longitude)
	
	return err
}

// GetByUserID retrieves a doctor by their user ID
func (r *DoctorRepo) GetByUserID(ctx context.Context, userID int) (*models.Doctor, error) {
	query := `
		SELECT user_id, specialty, service_type_id, license_number, experience_years, 
			qualifications, achievements, bio, profile_picture_url, consultation_fee, rating, 
			is_verified, is_available, status, latitude, longitude
		FROM doctors
		WHERE user_id = ?
	`

	var doctor models.Doctor
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&doctor.UserID, &doctor.Specialty, &doctor.ServiceTypeID, &doctor.LicenseNumber,
		&doctor.ExperienceYears, &doctor.Qualifications, &doctor.Achievements, &doctor.Bio,
		&doctor.ProfilePictureURL, &doctor.ConsultationFee, &doctor.Rating, &doctor.IsVerified,
		&doctor.IsAvailable, &doctor.Status, &doctor.Latitude, &doctor.Longitude,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("doctor not found")
		}
		return nil, err
	}

	return &doctor, nil
}

// Update updates an existing doctor's information
func (r *DoctorRepo) Update(ctx context.Context, doctor *models.Doctor) error {
	query := `
		UPDATE doctors
		SET specialty = ?, service_type_id = ?, license_number = ?, experience_years = ?,
			qualifications = ?, achievements = ?, bio = ?, profile_picture_url = ?,
			consultation_fee = ?, rating = ?, is_verified = ?, is_available = ?,
			status = ?, latitude = ?, longitude = ?
		WHERE user_id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		doctor.Specialty, doctor.ServiceTypeID, doctor.LicenseNumber, doctor.ExperienceYears,
		doctor.Qualifications, doctor.Achievements, doctor.Bio, doctor.ProfilePictureURL,
		doctor.ConsultationFee, doctor.Rating, doctor.IsVerified, doctor.IsAvailable,
		doctor.Status, doctor.Latitude, doctor.Longitude, doctor.UserID)

	return err
}

// Delete removes a doctor from the database
func (r *DoctorRepo) Delete(ctx context.Context, userID int) error {
	query := `DELETE FROM doctors WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

// List retrieves a list of doctors with optional filtering and pagination
func (r *DoctorRepo) List(ctx context.Context, filter models.DoctorFilter, offset, limit int) ([]*models.Doctor, error) {
	query := `
		SELECT user_id, specialty, service_type_id, license_number, experience_years, 
			qualifications, achievements, bio, profile_picture_url, consultation_fee, rating, 
			is_verified, is_available, status, latitude, longitude
		FROM doctors
		WHERE 1=1
	`
	var args []interface{}

	if filter.Specialty != "" {
		query += " AND specialty = ?"
		args = append(args, filter.Specialty)
	}

	if filter.ServiceTypeID != 0 {
		query += " AND service_type_id = ?"
		args = append(args, filter.ServiceTypeID)
	}

	if filter.MinRating > 0 {
		query += " AND rating >= ?"
		args = append(args, filter.MinRating)
	}

	query += " ORDER BY rating DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doctors []*models.Doctor
	for rows.Next() {
		var doctor models.Doctor
		err := rows.Scan(
			&doctor.UserID, &doctor.Specialty, &doctor.ServiceTypeID, &doctor.LicenseNumber,
			&doctor.ExperienceYears, &doctor.Qualifications, &doctor.Achievements, &doctor.Bio,
			&doctor.ProfilePictureURL, &doctor.ConsultationFee, &doctor.Rating, &doctor.IsVerified,
			&doctor.IsAvailable, &doctor.Status, &doctor.Latitude, &doctor.Longitude,
		)
		if err != nil {
			return nil, err
		}
		doctors = append(doctors, &doctor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return doctors, nil
}

// Add this new method to the DoctorRepo struct
func (r *DoctorRepo) SearchDoctors(ctx context.Context, params models.DoctorSearchParams) ([]models.DoctorSearchResult, error) {
    query := `
        CALL search_doctors(?, ?, ?, ?, ?, ?, ?);
    `
    
    rows, err := r.db.QueryContext(ctx, query,
        params.SearchTerm,
        params.Specialty,
        params.MinRating,
        params.LocationLat,
        params.LocationLng,
        params.RadiusKm,
        params.ServiceTypeID)
    if (err != nil) {
        return nil, fmt.Errorf("failed to execute search query: %w", err)
    }
    defer rows.Close()

    var results []models.DoctorSearchResult
    for rows.Next() {
        var result models.DoctorSearchResult
        var distanceKm sql.NullFloat64

        err := rows.Scan(
            &result.UserID,
            &result.FirstName,
            &result.LastName,
            &result.Email,
            &result.Phone,
            &result.Specialty,
            &result.ServiceTypeID,
            &result.Rating,
            &distanceKm,
            &result.IsAvailable,
            &result.Status,
        )
        if err != nil {
            return nil, fmt.Errorf("failed to scan row: %w", err)
        }

        if distanceKm.Valid {
            result.DistanceKm = &distanceKm.Float64
        }

        results = append(results, result)
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating over rows: %w", err)
    }

    return results, nil
}