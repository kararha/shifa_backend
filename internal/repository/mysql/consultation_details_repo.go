// File: internal/repository/mysql/consultation_details_repo.go

package mysql

import (
	"context"
	"database/sql"

	"shifa/internal/models"
)

// ConsultationDetailsRepo represents the MySQL repository for consultation details-related database operations
type ConsultationDetailsRepo struct {
	db *sql.DB
}

// NewConsultationDetailsRepo creates a new ConsultationDetailsRepo instance
func NewConsultationDetailsRepo(db *sql.DB) *ConsultationDetailsRepo {
	return &ConsultationDetailsRepo{db: db}
}

// Create inserts new consultation details into the database
func (r *ConsultationDetailsRepo) Create(ctx context.Context, details *models.ConsultationDetails) error {
	query := `
		INSERT INTO consultation_details (consultation_id, request_details, symptoms, diagnosis, prescription, notes)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	
	result, err := r.db.ExecContext(ctx, query, 
		details.ConsultationID, details.RequestDetails, details.Symptoms,
		details.Diagnosis, details.Prescription, details.Notes)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	details.ID = int(id)
	return nil
}

// GetByID retrieves consultation details by their ID
func (r *ConsultationDetailsRepo) GetByID(ctx context.Context, id int) (*models.ConsultationDetails, error) {
	query := `
		SELECT id, consultation_id, request_details, symptoms, diagnosis, prescription, notes
		FROM consultation_details
		WHERE id = ?
	`

	var details models.ConsultationDetails
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&details.ID, &details.ConsultationID, &details.RequestDetails,
		&details.Symptoms, &details.Diagnosis, &details.Prescription, &details.Notes,
	)

	if err != nil {
		return nil, err
	}

	return &details, nil
}

// GetByConsultationID retrieves consultation details by the consultation ID
func (r *ConsultationDetailsRepo) GetByConsultationID(ctx context.Context, consultationID int) (*models.ConsultationDetails, error) {
	query := `
		SELECT id, consultation_id, request_details, symptoms, diagnosis, prescription, notes
		FROM consultation_details
		WHERE consultation_id = ?
	`

	var details models.ConsultationDetails
	err := r.db.QueryRowContext(ctx, query, consultationID).Scan(
		&details.ID, &details.ConsultationID, &details.RequestDetails,
		&details.Symptoms, &details.Diagnosis, &details.Prescription, &details.Notes,
	)

	if err != nil {
		return nil, err
	}

	return &details, nil
}

// Update updates existing consultation details in the database
func (r *ConsultationDetailsRepo) Update(ctx context.Context, details *models.ConsultationDetails) error {
	query := `
		UPDATE consultation_details
		SET request_details = ?, symptoms = ?, diagnosis = ?, prescription = ?, notes = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		details.RequestDetails, details.Symptoms, details.Diagnosis,
		details.Prescription, details.Notes, details.ID)

	return err
}

// Delete removes consultation details from the database
func (r *ConsultationDetailsRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM consultation_details WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}