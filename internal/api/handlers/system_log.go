// File: internal/api/handlers/system_log.go

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"shifa/internal/models"
	"shifa/internal/service"
)

type SystemLogHandler struct {
	service *service.SystemLogService
}

func NewSystemLogHandler(service *service.SystemLogService) *SystemLogHandler {
	return &SystemLogHandler{service: service}
}

func (h *SystemLogHandler) LogAction(w http.ResponseWriter, r *http.Request) {
	var log models.SystemLog
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.service.LogAction(r.Context(), &log)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(log)
}

func (h *SystemLogHandler) GetUserLogs(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	logs, err := h.service.GetUserLogs(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func (h *SystemLogHandler) GetActionLogs(w http.ResponseWriter, r *http.Request) {
	actionType := r.URL.Query().Get("actionType")
	if actionType == "" {
		http.Error(w, "Action type is required", http.StatusBadRequest)
		return
	}

	logs, err := h.service.GetActionLogs(r.Context(), actionType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}