// File: service/audit_trail_service.go
package service

import (
    "shifa/internal/models"
    "shifa/internal/repository/mysql"
    "time"
)

type AuditTrailService struct {
    repo *mysql.AuditTrailRepo
}

func NewAuditTrailService(repo *mysql.AuditTrailRepo) *AuditTrailService {
    return &AuditTrailService{repo: repo}
}

func (s *AuditTrailService) LogChange(tableName string, recordID int, action string, changedFields models.JSON, changedBy int) error {
    audit := &models.AuditTrail{
        TableName:     tableName,
        RecordID:      recordID,
        Action:        action,
        ChangedFields: changedFields,
        ChangedBy:     changedBy,
        ChangedAt:     time.Now(),
    }
    return s.repo.Create(audit)
}

func (s *AuditTrailService) GetChangeHistory(tableName string, recordID int) ([]*models.AuditTrail, error) {
    return s.repo.GetByTableAndRecord(tableName, recordID)
}