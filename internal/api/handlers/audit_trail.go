// File: api/handlers/audit_trail.go
package handlers

import (
    "encoding/json"
    "net/http"
    "shifa/internal/service"
    "strconv"

    "github.com/gorilla/mux"
)

type AuditTrailHandler struct {
    auditService *service.AuditTrailService
}

func NewAuditTrailHandler(auditService *service.AuditTrailService) *AuditTrailHandler {
    return &AuditTrailHandler{auditService: auditService}
}

func (h *AuditTrailHandler) GetChangeHistory(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    tableName := vars["tableName"]
    recordID, err := strconv.Atoi(vars["recordID"])
    if err != nil {
        http.Error(w, "Invalid record ID", http.StatusBadRequest)
        return
    }

    history, err := h.auditService.GetChangeHistory(tableName, recordID)
    if err != nil {
        http.Error(w, "Failed to retrieve change history", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(history)
}
