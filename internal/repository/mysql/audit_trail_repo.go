// File: repository/mysql/audit_trail_repo.go
package mysql

import (
    "database/sql"
    "encoding/json"
    "shifa/internal/models"
)

type AuditTrailRepo struct {
    db *sql.DB
}

func NewAuditTrailRepo(db *sql.DB) *AuditTrailRepo {
    return &AuditTrailRepo{db: db}
}

func (r *AuditTrailRepo) Create(audit *models.AuditTrail) error {
    query := `
        INSERT INTO audit_trail (table_name, record_id, action, changed_fields, changed_by, changed_at)
        VALUES (?, ?, ?, ?, ?, ?)
    `
    changedFields, err := json.Marshal(audit.ChangedFields)
    if err != nil {
        return err
    }
    _, err = r.db.Exec(query, audit.TableName, audit.RecordID, audit.Action, changedFields, audit.ChangedBy, audit.ChangedAt)
    return err
}

func (r *AuditTrailRepo) GetByTableAndRecord(tableName string, recordID int) ([]*models.AuditTrail, error) {
    query := `
        SELECT id, table_name, record_id, action, changed_fields, changed_by, changed_at
        FROM audit_trail
        WHERE table_name = ? AND record_id = ?
        ORDER BY changed_at DESC
    `
    rows, err := r.db.Query(query, tableName, recordID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var audits []*models.AuditTrail
    for rows.Next() {
        var audit models.AuditTrail
        var changedFields []byte
        err := rows.Scan(&audit.ID, &audit.TableName, &audit.RecordID, &audit.Action, &changedFields, &audit.ChangedBy, &audit.ChangedAt)
        if err != nil {
            return nil, err
        }
        err = json.Unmarshal(changedFields, &audit.ChangedFields)
        if err != nil {
            return nil, err
        }
        audits = append(audits, &audit)
    }
    return audits, nil
}