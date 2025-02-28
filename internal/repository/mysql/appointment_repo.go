// File: internal/repository/mysql/appointment_repo.go

package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"shifa/internal/models"
	"shifa/internal/repository"
	"time"
)

// AppointmentFilter represents the filtering options for appointments
type AppointmentFilter struct {
	PatientID          int
	ProviderType       string
	DoctorID           int
	HomeCareProviderID int
	StartDate          time.Time
	EndDate            time.Time
	Status             string
}

// AppointmentRepo represents the MySQL repository for appointment-related database operations
type AppointmentRepo struct {
	db *sql.DB
}

// GetByProviderID retrieves appointments for a specific provider (doctor or home care provider)
func (r *AppointmentRepo) GetByProviderID(ctx context.Context, providerID int, providerType string) ([]*models.Appointment, error) {
	// Debug logging
	fmt.Printf("Getting appointments for provider ID: %d, type: %s\n", providerID, providerType)

	query := `
        SELECT a.id, a.patient_id, a.provider_type, a.doctor_id, a.home_care_provider_id,
            a.appointment_date, TIME_FORMAT(a.start_time, '%H:%i:%s') as start_time, 
            TIME_FORMAT(a.end_time, '%H:%i:%s') as end_time, 
            a.status, a.cancellation_reason, a.created_at, a.updated_at,
            COALESCE(u.name, 'Unknown Patient') as patient_name
        FROM appointments a
        LEFT JOIN users u ON u.id = a.patient_id
        WHERE `

	if providerType == "doctor" {
		query += "a.doctor_id = ?"
	} else if providerType == "home_care_provider" {
		query += "a.home_care_provider_id = ?"
	} else {
		return nil, fmt.Errorf("invalid provider type: %s", providerType)
	}

	query += " ORDER BY a.appointment_date DESC, a.start_time ASC"

	// Debug the final query
	fmt.Printf("Executing query: %s\n", query)

	rows, err := r.db.QueryContext(ctx, query, providerID)
	if err != nil {
		return nil, fmt.Errorf("failed to query appointments: %w", err)
	}
	defer rows.Close()

	var appointments []*models.Appointment
	for rows.Next() {
		var apt models.Appointment
		var startTimeStr, endTimeStr string
		var doctorID, homeCareProviderID sql.NullInt64
		var cancellationReason, patientName sql.NullString

		err := rows.Scan(
			&apt.ID,
			&apt.PatientID,
			&apt.ProviderType,
			&doctorID,
			&homeCareProviderID,
			&apt.AppointmentDate,
			&startTimeStr,
			&endTimeStr,
			&apt.Status,
			&cancellationReason,
			&apt.CreatedAt,
			&apt.UpdatedAt,
			&patientName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan appointment: %w", err)
		}

		// Handle nullable fields
		if doctorID.Valid {
			id := int(doctorID.Int64)
			apt.DoctorID = &id
		}
		if homeCareProviderID.Valid {
			id := int(homeCareProviderID.Int64)
			apt.HomeCareProviderID = &id
		}
		if cancellationReason.Valid {
			apt.CancellationReason = &cancellationReason.String
		}
		if patientName.Valid {
			apt.PatientName = patientName.String
		}

		// Parse time strings
		startTime, _ := time.Parse("15:04:05", startTimeStr)
		endTime, _ := time.Parse("15:04:05", endTimeStr)
		apt.StartTime = models.CustomTime(startTime)
		apt.EndTime = models.CustomTime(endTime)

		appointments = append(appointments, &apt)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating appointments: %w", err)
	}

	return appointments, nil
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
			doctorID           sql.NullInt64
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

func (r *AppointmentRepo) GetAppointmentsByPatient(patientID int) ([]models.Appointment, error) {
	query := `
        SELECT id, patient_id, doctor_id, home_care_provider_id, appointment_date, 
               status, created_at, updated_at, service_type_id, notes
        FROM appointments 
        WHERE patient_id = ?
        ORDER BY appointment_date DESC
    `

	rows, err := r.db.Query(query, patientID)
	if err != nil {
		return nil, fmt.Errorf("error querying appointments: %v", err)
	}
	defer rows.Close()

	var appointments []models.Appointment
	for rows.Next() {
		var apt models.Appointment
		err := rows.Scan(
			&apt.ID,
			&apt.PatientID,
			&apt.DoctorID,
			&apt.HomeCareProviderID,
			&apt.AppointmentDate,
			&apt.Status,
			&apt.CreatedAt,
			&apt.UpdatedAt,
			&apt.ServiceTypeID,
			&apt.CancellationReason,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning appointment row: %v", err)
		}
		appointments = append(appointments, apt)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating appointment rows: %v", err)
	}

	return appointments, nil
}
