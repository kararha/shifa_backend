package service

import (
    "context"
    "shifa/internal/models"
    "shifa/internal/repository/mysql"
    "time"
)

// SystemLogService handles the logging actions
type SystemLogService struct {
    repo *mysql.SystemLogRepo
}

// NewSystemLogService creates a new instance of SystemLogService
func NewSystemLogService(repo *mysql.SystemLogRepo) *SystemLogService {
    return &SystemLogService{repo: repo}
}

// LogAction logs an action with the current timestamp
func (s *SystemLogService) LogAction(ctx context.Context, log *models.SystemLog) error {
    log.Timestamp = time.Now()
    return s.repo.Create(ctx, log)
}

// GetUserLogs retrieves logs by user ID
func (s *SystemLogService) GetUserLogs(ctx context.Context, userID int) ([]*models.SystemLog, error) {
    return s.repo.GetByUserID(ctx, userID)
}

// GetActionLogs retrieves logs by action type
func (s *SystemLogService) GetActionLogs(ctx context.Context, actionType string) ([]*models.SystemLog, error) {
    return s.repo.GetByActionType(ctx, actionType)
}
