package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "shifa/internal/models"
    "shifa/internal/service"
    "github.com/gorilla/mux"
    "github.com/sirupsen/logrus"
)

type HomeCareVisitHandler struct {
    service *service.HomeCareVisitService
    logger  *logrus.Logger
}

func NewHomeCareVisitHandler(service *service.HomeCareVisitService, logger *logrus.Logger) *HomeCareVisitHandler {
    return &HomeCareVisitHandler{
        service: service,
        logger:  logger,
    }
}

func (h *HomeCareVisitHandler) ScheduleHomeCareVisit(w http.ResponseWriter, r *http.Request) {
    h.logger.Info("Handling schedule home care visit request")

    var visit models.HomeCareVisit
    if err := json.NewDecoder(r.Body).Decode(&visit); err != nil {
        h.logger.WithError(err).Error("Failed to decode request body")
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    err := h.service.ScheduleHomeCareVisit(r.Context(), &visit)
    if err != nil {
        h.logger.WithError(err).Error("Failed to schedule home care visit")
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    h.logger.WithField("visit_id", visit.ID).Info("Successfully scheduled home care visit")
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(visit)
}

func (h *HomeCareVisitHandler) GetHomeCareVisit(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        h.logger.WithError(err).Error("Invalid visit ID")
        http.Error(w, "Invalid visit ID", http.StatusBadRequest)
        return
    }

    h.logger.WithField("id", id).Info("Fetching home care visit")
    visit, err := h.service.GetVisitDetails(r.Context(), id)
    if err != nil {
        h.logger.WithError(err).WithField("id", id).Error("Failed to get home care visit")
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    h.logger.WithField("id", id).Info("Successfully fetched home care visit")
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(visit)
}

func (h *HomeCareVisitHandler) ListHomeCareVisits(w http.ResponseWriter, r *http.Request) {
    h.logger.Info("Handling list home care visits request")

    filter := models.HomeCareVisitFilter{
        PatientID:  parseIntParam(r, "patient_id"),
        ProviderID: parseIntParam(r, "provider_id"),
        Status:     r.URL.Query().Get("status"),
    }

    visits, err := h.service.ListHomeCareVisits(r.Context(), filter)
    if err != nil {
        h.logger.WithError(err).Error("Failed to list home care visits")
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    h.logger.WithField("count", len(visits)).Info("Successfully listed home care visits")
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(visits)
}

func (h *HomeCareVisitHandler) UpdateHomeCareVisit(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        h.logger.WithError(err).Error("Invalid visit ID")
        http.Error(w, "Invalid visit ID", http.StatusBadRequest)
        return
    }

    var visit models.HomeCareVisit
    if err := json.NewDecoder(r.Body).Decode(&visit); err != nil {
        h.logger.WithError(err).Error("Failed to decode request body")
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    visit.ID = id // Ensure the ID is set for the update
    err = h.service.UpdateHomeCareVisit(r.Context(), &visit)
    if err != nil {
        h.logger.WithError(err).WithField("id", id).Error("Failed to update home care visit")
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    h.logger.WithField("id", id).Info("Successfully updated home care visit")
    w.WriteHeader(http.StatusNoContent)
}

func (h *HomeCareVisitHandler) DeleteHomeCareVisit(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        h.logger.WithError(err).Error("Invalid visit ID")
        http.Error(w, "Invalid visit ID", http.StatusBadRequest)
        return
    }

    err = h.service.DeleteHomeCareVisit(r.Context(), id)
    if err != nil {
        h.logger.WithError(err).WithField("id", id).Error("Failed to delete home care visit")
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    h.logger.WithField("id", id).Info("Successfully deleted home care visit")
    w.WriteHeader(http.StatusNoContent)
}

// Helper function to parse string to int
func parseIntParam(r *http.Request, param string) int {
    value := r.URL.Query().Get(param)
    if value == "" {
        return 0
    }
    
    intValue, err := strconv.Atoi(value)
    if err != nil {
        return 0
    }
    return intValue
}