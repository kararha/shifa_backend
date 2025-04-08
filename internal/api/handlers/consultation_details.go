package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"shifa/internal/models"
	"shifa/internal/service"

	"github.com/gorilla/mux"
)

// ConsultationDetailsHandler handles HTTP requests for consultation details
type ConsultationDetailsHandler struct {
	detailsService *service.ConsultationDetailsService
}

// NewConsultationDetailsHandler creates a new ConsultationDetailsHandler
func NewConsultationDetailsHandler(detailsService *service.ConsultationDetailsService) *ConsultationDetailsHandler {
	return &ConsultationDetailsHandler{detailsService: detailsService}
}

// CreateDetails handles creating new consultation details
func (h *ConsultationDetailsHandler) CreateDetails(w http.ResponseWriter, r *http.Request) {
	var details models.ConsultationDetails
	if err := json.NewDecoder(r.Body).Decode(&details); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.detailsService.CreateDetails(r.Context(), details)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(details)
}

// GetDetails handles retrieving consultation details by ID
func (h *ConsultationDetailsHandler) GetDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	detailsID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid details ID", http.StatusBadRequest)
		return
	}

	details, err := h.detailsService.GetDetailsByID(r.Context(), detailsID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(details)
}

// GetDetailsByConsultation handles retrieving consultation details by consultation ID
func (h *ConsultationDetailsHandler) GetDetailsByConsultation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	consultationID, err := strconv.Atoi(vars["consultationId"])
	if err != nil {
		http.Error(w, "Invalid consultation ID", http.StatusBadRequest)
		return
	}

	details, err := h.detailsService.GetDetailsByConsultationID(r.Context(), consultationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(details)
}

// UpdateDetails handles updating existing consultation details
func (h *ConsultationDetailsHandler) UpdateDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	detailsID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid details ID", http.StatusBadRequest)
		return
	}

	var details models.ConsultationDetails
	if err := json.NewDecoder(r.Body).Decode(&details); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	details.ID = detailsID

	err = h.detailsService.UpdateDetails(r.Context(), details)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(details)
}

// DeleteDetails handles deleting consultation details
func (h *ConsultationDetailsHandler) DeleteDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	detailsID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid details ID", http.StatusBadRequest)
		return
	}

	err = h.detailsService.DeleteDetails(r.Context(), detailsID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
