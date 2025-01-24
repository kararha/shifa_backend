// File: internal/repository/mysql/patient_repo.go

package mysql

import (
	"context"
	"database/sql"
	"errors"

	"shifa/internal/models"
)

// PatientRepo represents the MySQL repository for patient-related database operations
type PatientRepo struct {
	db *sql.DB
}

// NewPatientRepo creates a new PatientRepo instance
func NewPatientRepo(db *sql.DB) *PatientRepo {
	return &PatientRepo{db: db}
}

// Create inserts a new patient into the database
func (r *PatientRepo) Create(ctx context.Context, patient *models.Patient) error {
	query := `
		INSERT INTO patients (user_id, date_of_birth, gender, phone, address, 
			emergency_contact_name, emergency_contact_phone)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	
		_, err := r.db.ExecContext(ctx, query, 
			patient.UserID, patient.DateOfBirth, patient.Gender, patient.Phone,
			patient.Address, patient.EmergencyContactName, patient.EmergencyContactPhone)
		
		return err
	}
// GetByUserID retrieves a patient by their user ID
func (r *PatientRepo) GetByUserID(ctx context.Context, userID int) (*models.Patient, error) {
	query := `
		SELECT p.user_id, p.date_of_birth, p.gender, p.phone, p.address, 
			p.emergency_contact_name, p.emergency_contact_phone, u.name
		FROM patients p
		JOIN users u ON p.user_id = u.id
		WHERE p.user_id = ?
	`

	var patient models.Patient
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&patient.UserID, &patient.DateOfBirth, &patient.Gender, &patient.Phone,
		&patient.Address, &patient.EmergencyContactName, &patient.EmergencyContactPhone,
		&patient.Name,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("patient not found")
		}
		return nil, err
	}

	return &patient, nil
}

// Update updates an existing patient's information
func (r *PatientRepo) Update(ctx context.Context, patient *models.Patient) error {
	query := `
		UPDATE patients
		SET date_of_birth = ?, gender = ?, phone = ?, address = ?,
			emergency_contact_name = ?, emergency_contact_phone = ?
		WHERE user_id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		patient.DateOfBirth, patient.Gender, patient.Phone, patient.Address,
		patient.EmergencyContactName, patient.EmergencyContactPhone, patient.UserID)

	return err
}

// Delete removes a patient from the database
func (r *PatientRepo) Delete(ctx context.Context, userID int) error {
	query := `DELETE FROM patients WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

// List retrieves a list of patients with optional pagination
func (r *PatientRepo) List(ctx context.Context, offset, limit int) ([]*models.Patient, error) {
	query := `
		SELECT p.user_id, u.name, p.date_of_birth, p.gender, p.phone, p.address, 
			p.emergency_contact_name, p.emergency_contact_phone
		FROM patients p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.user_id
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var patients []*models.Patient
	for rows.Next() {
		var patient models.Patient
		err := rows.Scan(
			&patient.UserID, &patient.Name, &patient.DateOfBirth, &patient.Gender, &patient.Phone,
			&patient.Address, &patient.EmergencyContactName, &patient.EmergencyContactPhone,
		)
		if err != nil {
			return nil, err
		}
		patients = append(patients, &patient)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return patients, nil
}