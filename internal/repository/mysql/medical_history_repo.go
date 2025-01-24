// File: internal/repository/mysql/medical_history_repo.go

package mysql

import (
	"context"
	"database/sql"
	"errors"

	"shifa/internal/models"
)

type MedicalHistoryRepo struct {
    db *sql.DB
}

func NewMedicalHistoryRepo(db *sql.DB) *MedicalHistoryRepo {
    return &MedicalHistoryRepo{db: db}
}

// Create inserts a new medical history record into the database
func (r *MedicalHistoryRepo) Create(ctx context.Context, history *models.MedicalHistory) error {
	query := `
		INSERT INTO medical_history (patient_id, condition_name, diagnosis_date, treatment, is_current)
		VALUES (?, ?, ?, ?, ?)
	`
	
	result, err := r.db.ExecContext(ctx, query, 
		history.PatientID, history.ConditionName, history.DiagnosisDate, history.Treatment, history.IsCurrent)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	history.ID = int(id)
	return nil
}

// GetByID retrieves a medical history record by its ID
func (r *MedicalHistoryRepo) GetByID(ctx context.Context, id int) (*models.MedicalHistory, error) {
	query := `
		SELECT id, patient_id, condition_name, diagnosis_date, treatment, is_current
		FROM medical_history
		WHERE id = ?
	`

	var history models.MedicalHistory
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&history.ID, &history.PatientID, &history.ConditionName,
		&history.DiagnosisDate, &history.Treatment, &history.IsCurrent,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("medical history record not found")
		}
		return nil, err
	}

	return &history, nil
}

// Update updates an existing medical history record
func (r *MedicalHistoryRepo) Update(ctx context.Context, history *models.MedicalHistory) error {
	query := `
		UPDATE medical_history
		SET condition_name = ?, diagnosis_date = ?, treatment = ?, is_current = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		history.ConditionName, history.DiagnosisDate, history.Treatment,
		history.IsCurrent, history.ID)

	return err
}

// Delete removes a medical history record from the database
func (r *MedicalHistoryRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM medical_history WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// Change this method name
func (r *MedicalHistoryRepo) GetByPatientID(ctx context.Context, patientID int) ([]*models.MedicalHistory, error) {
    query := `
        SELECT id, patient_id, condition_name, diagnosis_date, treatment, is_current
        FROM medical_history
        WHERE patient_id = ?
        ORDER BY diagnosis_date DESC
    `

    rows, err := r.db.QueryContext(ctx, query, patientID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var histories []*models.MedicalHistory
    for rows.Next() {
        var history models.MedicalHistory
        err := rows.Scan(
            &history.ID, &history.PatientID, &history.ConditionName,
            &history.DiagnosisDate, &history.Treatment, &history.IsCurrent,
        )
        if err != nil {
            return nil, err
        }
        histories = append(histories, &history)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return histories, nil
}