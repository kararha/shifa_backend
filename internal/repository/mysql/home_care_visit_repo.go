package mysql

import (
    "context"
    "database/sql"
    "time"
    "shifa/internal/models"
)

// HomeVisitRepo represents the MySQL repository for home care visit-related database operations
type HomeVisitRepo struct {
    db *sql.DB
}

// NewHomeVisitRepo creates a new HomeVisitRepo instance
func NewHomeVisitRepo(db *sql.DB) *HomeVisitRepo {
    return &HomeVisitRepo{db: db}
}

// Create inserts a new home care visit into the database
func (r *HomeVisitRepo) Create(ctx context.Context, visit *models.HomeCareVisit) error {
    query := `
        INSERT INTO home_care_visits (appointment_id, address, latitude, longitude, 
            duration_hours, special_requirements, status)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `

    _, err := r.db.ExecContext(ctx, query, 
        visit.AppointmentID, visit.Address, visit.Latitude, visit.Longitude,
        visit.DurationHours, visit.SpecialRequirements, visit.Status)
    return err
}

// GetByID retrieves a home care visit by its ID
func (r *HomeVisitRepo) GetByID(ctx context.Context, id int) (*models.HomeCareVisit, error) {
    query := `
        SELECT id, appointment_id, address, latitude, longitude, 
            duration_hours, special_requirements, status
        FROM home_care_visits
        WHERE id = ?
    `

    var visit models.HomeCareVisit
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &visit.ID, &visit.AppointmentID, &visit.Address, &visit.Latitude, &visit.Longitude,
        &visit.DurationHours, &visit.SpecialRequirements, &visit.Status)
    if err != nil {
        return nil, err
    }

    return &visit, nil
}

// List retrieves a list of home care visits with optional filtering
func (r *HomeVisitRepo) List(ctx context.Context, filter models.HomeCareVisitFilter) ([]models.HomeCareVisit, error) {
    query := `
        SELECT id, appointment_id, address, latitude, longitude, 
            duration_hours, special_requirements, status
        FROM home_care_visits
        WHERE 1=1
    `
    var args []interface{}

    if filter.AppointmentID != 0 {
        query += " AND appointment_id = ?"
        args = append(args, filter.AppointmentID)
    }

    if filter.Status != "" {
        query += " AND status = ?"
        args = append(args, filter.Status)
    }

    rows, err := r.db.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var visits []models.HomeCareVisit // Change to slice of values
    for rows.Next() {
        var visit models.HomeCareVisit
        err := rows.Scan(&visit.ID, &visit.AppointmentID, &visit.Address, &visit.Latitude, &visit.Longitude,
            &visit.DurationHours, &visit.SpecialRequirements, &visit.Status)
        if err != nil {
            return nil, err
        }
        visits = append(visits, visit) // Append value instead of pointer
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return visits, nil
}

// Delete removes a home care visit from the database by its ID
func (r *HomeVisitRepo) Delete(ctx context.Context, id int) error {
    query := `DELETE FROM home_care_visits WHERE id = ?`
    
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return sql.ErrNoRows // Or handle it differently if needed
    }

    return nil
}


// GetByDateRange retrieves home care visits within the specified date range
func (r *HomeVisitRepo) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.HomeCareVisit, error) {
    query := `
        SELECT id, appointment_id, address, latitude, longitude, 
            duration_hours, special_requirements, status
        FROM home_care_visits
        WHERE appointment_date BETWEEN ? AND ?
        ORDER BY appointment_date
    `

    rows, err := r.db.QueryContext(ctx, query, startDate, endDate)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var visits []models.HomeCareVisit
    for rows.Next() {
        var visit models.HomeCareVisit
        err := rows.Scan(&visit.ID, &visit.AppointmentID, &visit.Address, &visit.Latitude, &visit.Longitude,
            &visit.DurationHours, &visit.SpecialRequirements, &visit.Status)
        if err != nil {
            return nil, err
        }
        visits = append(visits, visit) // Append value instead of pointer
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return visits, nil
}


// GetByPatientID retrieves home care visits for a specific patient ID
func (r *HomeVisitRepo) GetByPatientID(ctx context.Context, patientID int) ([]models.HomeCareVisit, error) {
    query := `
        SELECT hv.id, hv.appointment_id, hv.address, hv.latitude, hv.longitude, 
               hv.duration_hours, hv.special_requirements, hv.status
        FROM home_care_visits hv
        INNER JOIN appointments a ON hv.appointment_id = a.id
        WHERE a.patient_id = ?
        ORDER BY hv.id
    `

    rows, err := r.db.QueryContext(ctx, query, patientID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var visits []models.HomeCareVisit
    for rows.Next() {
        var visit models.HomeCareVisit
        if err := rows.Scan(&visit.ID, &visit.AppointmentID, &visit.Address, &visit.Latitude, &visit.Longitude,
            &visit.DurationHours, &visit.SpecialRequirements, &visit.Status); err != nil {
            return nil, err
        }
        visits = append(visits, visit)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return visits, nil
}

// GetByProviderID retrieves home care visits for a specific provider ID
func (r *HomeVisitRepo) GetByProviderID(ctx context.Context, providerID int) ([]models.HomeCareVisit, error) {
    query := `
        SELECT hv.id, hv.appointment_id, hv.address, hv.latitude, hv.longitude, 
               hv.duration_hours, hv.special_requirements, hv.status
        FROM home_care_visits hv
        INNER JOIN appointments a ON hv.appointment_id = a.id
        WHERE a.provider_id = ?
        ORDER BY hv.id
    `

    rows, err := r.db.QueryContext(ctx, query, providerID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var visits []models.HomeCareVisit
    for rows.Next() {
        var visit models.HomeCareVisit
        err := rows.Scan(&visit.ID, &visit.AppointmentID, &visit.Address, &visit.Latitude, &visit.Longitude,
            &visit.DurationHours, &visit.SpecialRequirements, &visit.Status)
        if err != nil {
            return nil, err
        }
        visits = append(visits, visit) // Append value instead of pointer
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return visits, nil
}


// Update modifies an existing home care visit in the database
func (r *HomeVisitRepo) Update(ctx context.Context, visit *models.HomeCareVisit) error {
    query := `
        UPDATE home_care_visits
        SET appointment_id = ?, address = ?, latitude = ?, longitude = ?, 
            duration_hours = ?, special_requirements = ?, status = ?
        WHERE id = ?
    `

    _, err := r.db.ExecContext(ctx, query, 
        visit.AppointmentID, visit.Address, visit.Latitude, visit.Longitude,
        visit.DurationHours, visit.SpecialRequirements, visit.Status, visit.ID)
    return err
}