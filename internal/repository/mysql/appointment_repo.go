// File: internal/repository/mysql/appointment_repo.go

package mysql

import (
    "fmt"
	"context"
	"database/sql"
	"errors"
	"time"
	"shifa/internal/repository"
	"shifa/internal/models"
)



// AppointmentFilter represents the filtering options for appointments
type AppointmentFilter struct {
    PatientID           int
    ProviderType        string
    DoctorID            int
    HomeCareProviderID  int
    StartDate           time.Time
    EndDate             time.Time
    Status              string
}

// AppointmentRepo represents the MySQL repository for appointment-related database operations
type AppointmentRepo struct {
	db *sql.DB
}

// GetByProviderID retrieves appointments for a specific provider (doctor or home care provider)
func (r *AppointmentRepo) GetByProviderID(ctx context.Context, providerID int, providerType string) ([]*models.Appointment, error) {
    query := `
        SELECT id, patient_id, provider_type, doctor_id, home_care_provider_id,
            appointment_date, start_time, end_time, status, cancellation_reason,
            created_at, updated_at
        FROM appointments
        WHERE `

    if providerType == "doctor" {
        query += "doctor_id = ?"
    } else if providerType == "homecare" {
        query += "home_care_provider_id = ?"
    } else {
        return nil, fmt.Errorf("invalid provider type: %s", providerType)
    }

    query += " ORDER BY appointment_date, start_time"
    
    return r.queryAppointments(ctx, query, providerID)
}

// Create inserts a new appointment into the database
func (r *AppointmentRepo) Create(ctx context.Context, appointment *models.Appointment) error {
	query := `
		INSERT INTO appointments (patient_id, provider_type, doctor_id, home_care_provider_id,
			appointment_date, start_time, end_time, status)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	result, err := r.db.ExecContext(ctx, query, 
		appointment.PatientID, appointment.ProviderType, appointment.DoctorID,
		appointment.HomeCareProviderID, appointment.AppointmentDate, appointment.StartTime,
		appointment.EndTime, appointment.Status)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	appointment.ID = int(id)
	return nil
}

// GetByID retrieves an appointment by its ID
// func (r *AppointmentRepo) GetByID(ctx context.Context, id int) (*models.Appointment, error) {
// 	query := `
// 		SELECT id, patient_id, provider_type, doctor_id, home_care_provider_id,
// 			appointment_date, start_time, end_time, status, cancellation_reason,
// 			created_at, updated_at
// 		FROM appointments
// 		WHERE id = ?
// 	`

// 	var appointment models.Appointment
// 	err := r.db.QueryRowContext(ctx, query, id).Scan(
// 		&appointment.ID, &appointment.PatientID, &appointment.ProviderType,
// 		&appointment.DoctorID, &appointment.HomeCareProviderID, &appointment.AppointmentDate,
// 		&appointment.StartTime, &appointment.EndTime, &appointment.Status,
// 		&appointment.CancellationReason, &appointment.CreatedAt, &appointment.UpdatedAt,
// 	)

// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, errors.New("appointment not found")
// 		}
// 		return nil, err
// 	}

// 	return &appointment, nil
// }

// GetByID retrieves an appointment by its ID
func (r *AppointmentRepo) GetByID(ctx context.Context, id int) (*models.Appointment, error) {
    query := `
        SELECT id, patient_id, provider_type, doctor_id, home_care_provider_id,
            appointment_date, TIME_FORMAT(start_time, '%H:%i:%s') as start_time, 
            TIME_FORMAT(end_time, '%H:%i:%s') as end_time, 
            status, cancellation_reason,
            created_at, updated_at
        FROM appointments
        WHERE id = ?
    `

    var appointment models.Appointment
    var startTimeStr, endTimeStr string

    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &appointment.ID,
        &appointment.PatientID,
        &appointment.ProviderType,
        &appointment.DoctorID,
        &appointment.HomeCareProviderID,
        &appointment.AppointmentDate,
        &startTimeStr,
        &endTimeStr,
        &appointment.Status,
        &appointment.CancellationReason,
        &appointment.CreatedAt,
        &appointment.UpdatedAt,
    )

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, errors.New("appointment not found")
        }
        return nil, fmt.Errorf("failed to scan appointment: %w", err)
    }

    // Parse the time strings
    startTime, err := time.Parse("15:04:05", startTimeStr)
    if err != nil {
        return nil, fmt.Errorf("failed to parse start_time: %w", err)
    }
    appointment.StartTime = models.CustomTime(startTime)

    endTime, err := time.Parse("15:04:05", endTimeStr)
    if err != nil {
        return nil, fmt.Errorf("failed to parse end_time: %w", err)
    }
    appointment.EndTime = models.CustomTime(endTime)

    return &appointment, nil
}

// Update updates an existing appointment's information
func (r *AppointmentRepo) Update(ctx context.Context, appointment *models.Appointment) error {
	query := `
		UPDATE appointments
		SET patient_id = ?, provider_type = ?, doctor_id = ?, home_care_provider_id = ?,
			appointment_date = ?, start_time = ?, end_time = ?, status = ?,
			cancellation_reason = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		appointment.PatientID, appointment.ProviderType, appointment.DoctorID,
		appointment.HomeCareProviderID, appointment.AppointmentDate, appointment.StartTime,
		appointment.EndTime, appointment.Status, appointment.CancellationReason,
		time.Now(), appointment.ID)

	return err
}

// Delete removes an appointment from the database
func (r *AppointmentRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM appointments WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// GetByPatientID retrieves appointments for a specific patient
func (r *AppointmentRepo) GetByPatientID(ctx context.Context, patientID int, limit, offset int) ([]*models.Appointment, error) {
    query := `
        SELECT id, patient_id, provider_type, doctor_id, home_care_provider_id,
            appointment_date, start_time, end_time, status, cancellation_reason,
            created_at, updated_at
        FROM appointments
        WHERE patient_id = ?
        ORDER BY appointment_date, start_time
        LIMIT ? OFFSET ?
    `

    return r.queryAppointments(ctx, query, patientID, limit, offset)
}

// GetByDoctorID retrieves appointments for a specific doctor
func (r *AppointmentRepo) GetByDoctorID(ctx context.Context, doctorID int, limit, offset int) ([]*models.Appointment, error) {
    query := `
        SELECT id, patient_id, provider_type, doctor_id, home_care_provider_id,
            appointment_date, start_time, end_time, status, cancellation_reason,
            created_at, updated_at
        FROM appointments
        WHERE doctor_id = ?
        ORDER BY appointment_date, start_time
        LIMIT ? OFFSET ?
    `

    return r.queryAppointments(ctx, query, doctorID, limit, offset)
}

// GetByHomeCareProviderID retrieves appointments for a specific home care provider
func (r *AppointmentRepo) GetByHomeCareProviderID(ctx context.Context, providerID int, limit, offset int) ([]*models.Appointment, error) {
    query := `
        SELECT id, patient_id, provider_type, doctor_id, home_care_provider_id,
            appointment_date, start_time, end_time, status, cancellation_reason,
            created_at, updated_at
        FROM appointments
        WHERE home_care_provider_id = ?
        ORDER BY appointment_date, start_time
        LIMIT ? OFFSET ?
    `

    return r.queryAppointments(ctx, query, providerID, limit, offset)
}

// Helper function to query appointments
func (r *AppointmentRepo) queryAppointments(ctx context.Context, query string, args ...interface{}) ([]*models.Appointment, error) {
    rows, err := r.db.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var appointments []*models.Appointment
    for rows.Next() {
        var appointment models.Appointment
        // Create nullable variables for optional fields
        var (
            doctorID            sql.NullInt64
            homeCareProviderID sql.NullInt64
            cancellationReason sql.NullString
        )

        err := rows.Scan(
            &appointment.ID,
            &appointment.PatientID,
            &appointment.ProviderType,
            &doctorID,
            &homeCareProviderID,
            &appointment.AppointmentDate,
            &appointment.StartTime,
            &appointment.EndTime,
            &appointment.Status,
            &cancellationReason,
            &appointment.CreatedAt,
            &appointment.UpdatedAt,
        )
        if err != nil {
            return nil, err
        }

        // Convert nullable fields to pointers
        if doctorID.Valid {
            id := int(doctorID.Int64)
            appointment.DoctorID = &id
        }
        if homeCareProviderID.Valid {
            id := int(homeCareProviderID.Int64)
            appointment.HomeCareProviderID = &id
        }
        if cancellationReason.Valid {
            reason := cancellationReason.String
            appointment.CancellationReason = &reason
        }

        appointments = append(appointments, &appointment)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return appointments, nil
}


// List retrieves appointments with filtering, pagination
func (r *AppointmentRepo) List(ctx context.Context, filter repository.AppointmentFilter, offset, limit int) ([]*models.Appointment, error) {
    query := `
        SELECT id, patient_id, provider_type, doctor_id, home_care_provider_id,
            appointment_date, start_time, end_time, status, cancellation_reason,
            created_at, updated_at
        FROM appointments
        WHERE 1=1
    `
    args := []interface{}{}

    // Add filter conditions
    if filter.StartDate != nil {
        query += " AND appointment_date >= ?"
        args = append(args, filter.StartDate)
    }
    if filter.EndDate != nil {
        query += " AND appointment_date <= ?"
        args = append(args, filter.EndDate)
    }
    if filter.Status != "" {
        query += " AND status = ?"
        args = append(args, filter.Status)
    }
    if filter.ProviderType != "" {
        query += " AND provider_type = ?"
        args = append(args, filter.ProviderType)
    }

    query += " ORDER BY appointment_date, start_time LIMIT ? OFFSET ?"
    args = append(args, limit, offset)

    return r.queryAppointments(ctx, query, args...)
}
// NewAppointmentRepository creates a new instance of AppointmentRepo
func NewAppointmentRepository(db *sql.DB) *AppointmentRepo {
    return &AppointmentRepo{
        db: db,
    }
}