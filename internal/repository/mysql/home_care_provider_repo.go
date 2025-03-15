// File: internal/repository/mysql/home_care_provider_repo.go

package mysql

import (
	"context"
	"database/sql"
	"errors"
	"shifa/internal/models"
)

// HomeCareProviderRepo represents the MySQL repository for home care provider-related database operations
type HomeCareProviderRepo struct {
	db *sql.DB
}

// NewHomeCareProviderRepo creates a new HomeCareProviderRepo instance
func NewHomeCareProviderRepo(db *sql.DB) *HomeCareProviderRepo {
	return &HomeCareProviderRepo{db: db}
}

// GetAll retrieves all home care providers from the database
func (r *HomeCareProviderRepo) GetAll(ctx context.Context) ([]models.HomeCareProvider, error) {
	query := `
		SELECT user_id, service_type_id, experience_years, qualifications, bio, 
			profile_picture_url, hourly_rate, rating, is_verified, is_available, 
			status, latitude, longitude
		FROM home_care_providers
		ORDER BY rating DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []models.HomeCareProvider
	for rows.Next() {
		var provider models.HomeCareProvider
		err := rows.Scan(
			&provider.UserID, &provider.ServiceTypeID, &provider.ExperienceYears,
			&provider.Qualifications, &provider.Bio, &provider.ProfilePictureURL,
			&provider.HourlyRate, &provider.Rating, &provider.IsVerified,
			&provider.IsAvailable, &provider.Status, &provider.Latitude, &provider.Longitude,
		)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return providers, nil
}

// GetByID retrieves a home care provider by their ID
func (r *HomeCareProviderRepo) GetByID(ctx context.Context, id int) (*models.HomeCareProvider, error) {
	query := `
		SELECT user_id, service_type_id, experience_years, qualifications, bio, 
			profile_picture_url, hourly_rate, rating, is_verified, is_available, 
			status, latitude, longitude
		FROM home_care_providers
		WHERE user_id = ?
	`

	var provider models.HomeCareProvider
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&provider.UserID, &provider.ServiceTypeID, &provider.ExperienceYears,
		&provider.Qualifications, &provider.Bio, &provider.ProfilePictureURL,
		&provider.HourlyRate, &provider.Rating, &provider.IsVerified,
		&provider.IsAvailable, &provider.Status, &provider.Latitude, &provider.Longitude,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("home care provider not found")
		}
		return nil, err
	}

	return &provider, nil
}

// GetByServiceType retrieves home care providers by service type with pagination
func (r *HomeCareProviderRepo) GetByServiceType(ctx context.Context, serviceTypeID int, limit, offset int) ([]*models.HomeCareProvider, error) {
	query := `
		SELECT user_id, service_type_id, experience_years, qualifications, bio, 
			profile_picture_url, hourly_rate, rating, is_verified, is_available, 
			status, latitude, longitude
		FROM home_care_providers
		WHERE service_type_id = ?
		ORDER BY rating DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, serviceTypeID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []*models.HomeCareProvider
	for rows.Next() {
		var provider models.HomeCareProvider
		err := rows.Scan(
			&provider.UserID, &provider.ServiceTypeID, &provider.ExperienceYears,
			&provider.Qualifications, &provider.Bio, &provider.ProfilePictureURL,
			&provider.HourlyRate, &provider.Rating, &provider.IsVerified,
			&provider.IsAvailable, &provider.Status, &provider.Latitude, &provider.Longitude,
		)
		if err != nil {
			return nil, err
		}
		providers = append(providers, &provider)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return providers, nil
}

// Create inserts a new home care provider into the database
func (r *HomeCareProviderRepo) Create(ctx context.Context, provider *models.HomeCareProvider) error {
	query := `
		INSERT INTO home_care_providers (user_id, service_type_id, experience_years, 
			qualifications, bio, profile_picture_url, hourly_rate, rating, 
			is_verified, is_available, status, latitude, longitude)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		provider.UserID, provider.ServiceTypeID, provider.ExperienceYears,
		provider.Qualifications, provider.Bio, provider.ProfilePictureURL,
		provider.HourlyRate, provider.Rating, provider.IsVerified,
		provider.IsAvailable, provider.Status, provider.Latitude, provider.Longitude)

	return err
}

// GetByUserID retrieves a home care provider by their user ID
func (r *HomeCareProviderRepo) GetByUserID(ctx context.Context, userID int) (*models.HomeCareProvider, error) {
	query := `
		SELECT user_id, service_type_id, experience_years, qualifications, bio, 
			profile_picture_url, hourly_rate, rating, is_verified, is_available, 
			status, latitude, longitude
		FROM home_care_providers
		WHERE user_id = ?
	`

	var provider models.HomeCareProvider
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&provider.UserID, &provider.ServiceTypeID, &provider.ExperienceYears,
		&provider.Qualifications, &provider.Bio, &provider.ProfilePictureURL,
		&provider.HourlyRate, &provider.Rating, &provider.IsVerified,
		&provider.IsAvailable, &provider.Status, &provider.Latitude, &provider.Longitude,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("home care provider not found")
		}
		return nil, err
	}

	return &provider, nil
}

// Update updates an existing home care provider's information
func (r *HomeCareProviderRepo) Update(ctx context.Context, provider *models.HomeCareProvider) error {
	query := `
		UPDATE home_care_providers
		SET service_type_id = ?, experience_years = ?, qualifications = ?, 
			bio = ?, profile_picture_url = ?, hourly_rate = ?, rating = ?, 
			is_verified = ?, is_available = ?, status = ?, latitude = ?, longitude = ?
		WHERE user_id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		provider.ServiceTypeID, provider.ExperienceYears, provider.Qualifications,
		provider.Bio, provider.ProfilePictureURL, provider.HourlyRate, provider.Rating,
		provider.IsVerified, provider.IsAvailable, provider.Status,
		provider.Latitude, provider.Longitude, provider.UserID)

	return err
}

// Delete removes a home care provider from the database
func (r *HomeCareProviderRepo) Delete(ctx context.Context, userID int) error {
	query := `DELETE FROM home_care_providers WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

// Search searches for providers based on a query string
func (r *HomeCareProviderRepo) Search(ctx context.Context, query string) ([]*models.HomeCareProvider, error) {
	searchQuery := `
        SELECT user_id, service_type_id, experience_years, qualifications, bio, 
            profile_picture_url, hourly_rate, rating, is_verified, is_available, 
            status, latitude, longitude
        FROM home_care_providers
        WHERE 
            LOWER(qualifications) LIKE LOWER(?) OR
            LOWER(bio) LIKE LOWER(?)
        ORDER BY rating DESC
        LIMIT 10
    `

	searchParam := "%" + query + "%"
	rows, err := r.db.QueryContext(ctx, searchQuery, searchParam, searchParam)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []*models.HomeCareProvider
	for rows.Next() {
		var provider models.HomeCareProvider
		err := rows.Scan(
			&provider.UserID, &provider.ServiceTypeID, &provider.ExperienceYears,
			&provider.Qualifications, &provider.Bio, &provider.ProfilePictureURL,
			&provider.HourlyRate, &provider.Rating, &provider.IsVerified,
			&provider.IsAvailable, &provider.Status, &provider.Latitude, &provider.Longitude,
		)
		if err != nil {
			return nil, err
		}
		providers = append(providers, &provider)
	}

	return providers, rows.Err()
}
