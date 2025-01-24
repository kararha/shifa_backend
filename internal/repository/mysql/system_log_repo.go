package mysql

import (
    "context"
    "database/sql"
    "shifa/internal/models"
)

// SystemLogRepo represents the MySQL repository for system log-related database operations
type SystemLogRepo struct {
    db *sql.DB
}

// NewSystemLogRepo creates a new SystemLogRepo instance
func NewSystemLogRepo(db *sql.DB) *SystemLogRepo {
    return &SystemLogRepo{db: db}
}

// Create inserts a new system log entry into the database
func (r *SystemLogRepo) Create(ctx context.Context, log *models.SystemLog) error {
    query := `INSERT INTO system_logs (user_id, user_type, action_type, action_description, entity_type, entity_id, old_value, new_value, ip_address, user_agent, additional_info)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
    result, err := r.db.ExecContext(ctx, query, log.UserID, log.UserType, log.ActionType, log.ActionDescription,
        log.EntityType, log.EntityID, log.OldValue, log.NewValue,
        log.IPAddress, log.UserAgent, log.AdditionalInfo)
    if err != nil {
        return err
    }
    id, err := result.LastInsertId()
    if err != nil {
        return err
    }
    log.ID = id
    return nil
}

// GetByUserID retrieves logs by user ID
func (r *SystemLogRepo) GetByUserID(ctx context.Context, userID int) ([]*models.SystemLog, error) {
    query := `SELECT id, timestamp, user_id, user_type, action_type, action_description, entity_type, entity_id, old_value, new_value, ip_address, user_agent, additional_info
              FROM system_logs WHERE user_id = ?`
    rows, err := r.db.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var logs []*models.SystemLog
    for rows.Next() {
        var log models.SystemLog
        err := rows.Scan(&log.ID, &log.Timestamp, &log.UserID, &log.UserType, &log.ActionType,
            &log.ActionDescription, &log.EntityType, &log.EntityID, &log.OldValue,
            &log.NewValue, &log.IPAddress, &log.UserAgent, &log.AdditionalInfo)
        if err != nil {
            return nil, err
        }
        logs = append(logs, &log)
    }
    return logs, nil
}

// GetByActionType retrieves logs by action type
func (r *SystemLogRepo) GetByActionType(ctx context.Context, actionType string) ([]*models.SystemLog, error) {
    query := `SELECT id, timestamp, user_id, user_type, action_type, action_description, entity_type, entity_id, old_value, new_value, ip_address, user_agent, additional_info
              FROM system_logs WHERE action_type = ?`
    rows, err := r.db.QueryContext(ctx, query, actionType)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    var logs []*models.SystemLog
    for rows.Next() {
        var log models.SystemLog
        err := rows.Scan(&log.ID, &log.Timestamp, &log.UserID, &log.UserType, &log.ActionType,
            &log.ActionDescription, &log.EntityType, &log.EntityID, &log.OldValue,
            &log.NewValue, &log.IPAddress, &log.UserAgent, &log.AdditionalInfo)
        if err != nil {
            return nil, err
        }
        logs = append(logs, &log)
    }
    return logs, nil
}
