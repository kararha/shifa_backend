// File: internal/repository/mysql/consultation_repo.go

package mysql

import (
	"context"
	"database/sql"
	"errors"

	// "time"

	"shifa/internal/models"
)

// ConsultationRepo represents the MySQL repository for consultation-related database operations
type ConsultationRepo struct {
	db *sql.DB
}

// NewConsultationRepo creates a new ConsultationRepo instance
func NewConsultationRepo(db *sql.DB) *ConsultationRepo {
	return &ConsultationRepo{db: db}
}

// Create inserts a new consultation into the database
func (r *ConsultationRepo) Create(ctx context.Context, consultation *models.Consultation) error {
	query := `
		INSERT INTO consultations (patient_id, doctor_id, status, started_at, completed_at, fee)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		consultation.PatientID, consultation.DoctorID,
		consultation.Status,
		sql.NullTime{Time: consultation.StartedAt.Time, Valid: consultation.StartedAt.Valid},
		sql.NullTime{Time: consultation.CompletedAt.Time, Valid: consultation.CompletedAt.Valid},
		consultation.Fee)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	consultation.ID = int(id)
	return nil
}

// GetByID retrieves a consultation by its ID
func (r *ConsultationRepo) GetByID(ctx context.Context, id int) (*models.Consultation, error) {
	query := `
		SELECT id, patient_id, doctor_id, status, started_at, completed_at, fee
		FROM consultations
		WHERE id = ?
	`

	var consultation models.Consultation
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&consultation.ID, &consultation.PatientID, &consultation.DoctorID,
		&consultation.Status, &consultation.StartedAt, &consultation.CompletedAt,
		&consultation.Fee,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("consultation not found")
		}
		return nil, err
	}

	return &consultation, nil
}

// GetByAppointmentID retrieves a consultation by its appointment ID
func (r *ConsultationRepo) GetByAppointmentID(ctx context.Context, appointmentID int) (*models.Consultation, error) {
	query := `
        SELECT id, patient_id, doctor_id, status, started_at, completed_at, fee
        FROM consultations
        WHERE appointment_id = ?
        LIMIT 1
    `

	var consultation models.Consultation
	err := r.db.QueryRowContext(ctx, query, appointmentID).Scan(
		&consultation.ID, &consultation.PatientID, &consultation.DoctorID,
		&consultation.Status, &consultation.StartedAt, &consultation.CompletedAt,
		&consultation.Fee,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("consultation not found")
		}
		return nil, err
	}

	return &consultation, nil
}

// Update updates an existing consultation's information
func (r *ConsultationRepo) Update(ctx context.Context, consultation *models.Consultation) error {
	query := `
		UPDATE consultations
		SET status = ?, started_at = ?, completed_at = ?, fee = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		consultation.Status,
		sql.NullTime{Time: consultation.StartedAt.Time, Valid: consultation.StartedAt.Valid},
		sql.NullTime{Time: consultation.CompletedAt.Time, Valid: consultation.CompletedAt.Valid},
		consultation.Fee, consultation.ID)

	return err
}

// Delete removes a consultation from the database
func (r *ConsultationRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM consultations WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List retrieves a list of consultations with optional filtering and pagination
func (r *ConsultationRepo) List(ctx context.Context, filter models.ConsultationFilter, offset, limit int) ([]*models.Consultation, error) {
	query := `
        SELECT id, patient_id, doctor_id,
               status, started_at, completed_at, fee
        FROM consultations
        WHERE 1=1
    `
	var args []interface{}

	if filter.PatientID != 0 {
		query += " AND patient_id = ?"
		args = append(args, filter.PatientID)
	}

	if filter.DoctorID != 0 {
		query += " AND doctor_id = ?"
		args = append(args, filter.DoctorID)
	}

	if filter.Status != "" {
		query += " AND status = ?"
		args = append(args, filter.Status)
	}

	if !filter.StartDateFrom.IsZero() {
		query += " AND started_at >= ?"
		args = append(args, filter.StartDateFrom)
	}

	if !filter.StartDateTo.IsZero() {
		query += " AND started_at <= ?"
		args = append(args, filter.StartDateTo)
	}

	if !filter.CompletedDateFrom.IsZero() {
		query += " AND completed_at >= ?"
		args = append(args, filter.CompletedDateFrom)
	}

	if !filter.CompletedDateTo.IsZero() {
		query += " AND completed_at <= ?"
		args = append(args, filter.CompletedDateTo)
	}

	query += " ORDER BY started_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consultations []*models.Consultation
	for rows.Next() {
		var consultation models.Consultation
		err := rows.Scan(
			&consultation.ID, &consultation.PatientID, &consultation.DoctorID,
			&consultation.Status, &consultation.StartedAt, &consultation.CompletedAt,
			&consultation.Fee,
		)
		if err != nil {
			return nil, err
		}
		consultations = append(consultations, &consultation)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return consultations, nil
}
