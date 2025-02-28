// File: internal/repository/mysql/doctor_availability_repo.go

package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"shifa/internal/models"
)

// DoctorAvailabilityRepo represents the MySQL repository for doctor availability-related database operations
type DoctorAvailabilityRepo struct {
	db *sql.DB
}

// NewDoctorAvailabilityRepo creates a new DoctorAvailabilityRepo instance
func NewDoctorAvailabilityRepo(db *sql.DB) *DoctorAvailabilityRepo {
	return &DoctorAvailabilityRepo{db: db}
}

// Create inserts a new doctor availability record into the database
func (r *DoctorAvailabilityRepo) Create(ctx context.Context, availability *models.DoctorAvailability) error {
	query := `
		INSERT INTO doctor_availability (doctor_id, day_of_week, start_time, end_time)
		VALUES (?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		availability.DoctorID, availability.DayOfWeek, availability.StartTime, availability.EndTime)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	availability.ID = int(id)
	return nil
}

// GetByID retrieves a doctor availability record by its ID

func (r *DoctorAvailabilityRepo) GetByDoctorID(ctx context.Context, doctorID int) ([]*models.DoctorAvailability, error) {
	query := `
		SELECT id, doctor_id, day_of_week, start_time, end_time
		FROM doctor_availability
		WHERE doctor_id = ?
	`
	var availability models.DoctorAvailability
	var startStr, endStr string // Changed from []byte to string
	err := r.db.QueryRowContext(ctx, query, doctorID).Scan(
		&availability.ID, &availability.DoctorID, &availability.DayOfWeek, &startStr, &endStr,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("doctor availability not found")
		}
		return nil, err
	}
	availability.StartTime, err = time.Parse("15:04:05", startStr)
	if err != nil {
		return nil, err
	}
	availability.EndTime, err = time.Parse("15:04:05", endStr)
	if err != nil {
		return nil, err
	}
	return []*models.DoctorAvailability{&availability}, nil
}

// Update updates an existing doctor availability record
func (r *DoctorAvailabilityRepo) Update(ctx context.Context, availability *models.DoctorAvailability) error {
	query := `
		UPDATE doctor_availability
		SET doctor_id = ?, day_of_week = ?, start_time = ?, end_time = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		availability.DoctorID, availability.DayOfWeek, availability.StartTime,
		availability.EndTime, availability.ID)

	return err
}

// Delete removes a doctor availability record from the database
func (r *DoctorAvailabilityRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM doctor_availability WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// ListByDoctorID retrieves all availability records for a specific doctor
func (r *DoctorAvailabilityRepo) ListByDoctorID(ctx context.Context, doctorID int) ([]*models.DoctorAvailability, error) {
	query := `
		SELECT id, doctor_id, day_of_week, start_time, end_time
		FROM doctor_availability
		WHERE doctor_id = ?
		ORDER BY day_of_week, start_time
	`
	rows, err := r.db.QueryContext(ctx, query, doctorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var availabilities []*models.DoctorAvailability
	for rows.Next() {
		var availability models.DoctorAvailability
		var startStr, endStr string // Changed from []byte to string
		err := rows.Scan(
			&availability.ID, &availability.DoctorID, &availability.DayOfWeek, &startStr, &endStr,
		)
		if err != nil {
			return nil, err
		}
		availability.StartTime, err = time.Parse("15:04:05", startStr)
		if err != nil {
			return nil, err
		}
		availability.EndTime, err = time.Parse("15:04:05", endStr)
		if err != nil {
			return nil, err
		}
		availabilities = append(availabilities, &availability)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return availabilities, nil
}

// ListAllAvailability retrieves all availability records
func (r *DoctorAvailabilityRepo) ListAllAvailability(ctx context.Context) ([]*models.DoctorAvailability, error) {
	query := `
		SELECT id, doctor_id, day_of_week, start_time, end_time
		FROM doctor_availability
		ORDER BY doctor_id, day_of_week, start_time
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var availabilities []*models.DoctorAvailability
	for rows.Next() {
		var availability models.DoctorAvailability
		var startStr, endStr string // Changed from []byte to string
		err := rows.Scan(&availability.ID, &availability.DoctorID, &availability.DayOfWeek, &startStr, &endStr)
		if err != nil {
			return nil, err
		}
		availability.StartTime, err = time.Parse("15:04:05", startStr)
		if err != nil {
			return nil, err
		}
		availability.EndTime, err = time.Parse("15:04:05", endStr)
		if err != nil {
			return nil, err
		}
		availabilities = append(availabilities, &availability)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return availabilities, nil
}
