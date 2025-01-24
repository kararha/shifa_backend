package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
    "shifa/internal/models"
    "shifa/internal/service"
)

type MedicalHistoryHandler struct {
    service *service.MedicalHistoryService
}

func NewMedicalHistoryHandler(service *service.MedicalHistoryService) *MedicalHistoryHandler {
    return &MedicalHistoryHandler{service: service}
}

func (h *MedicalHistoryHandler) CreateMedicalHistory(w http.ResponseWriter, r *http.Request) {
    var medicalHistory models.MedicalHistory
    if err := json.NewDecoder(r.Body).Decode(&medicalHistory); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    created, err := h.service.CreateMedicalHistory(r.Context(), &medicalHistory)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(created)
}

func (h *MedicalHistoryHandler) GetMedicalHistory(w http.ResponseWriter, r *http.Request) {
    patientIDStr := r.URL.Query().Get("patient_id")
    patientID, err := strconv.Atoi(patientIDStr)
    if err != nil {
        http.Error(w, "Invalid patient ID", http.StatusBadRequest)
        return
    }

    histories, err := h.service.GetMedicalHistoryByPatientID(r.Context(), patientID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(histories)
}

func (h *MedicalHistoryHandler) UpdateMedicalHistory(w http.ResponseWriter, r *http.Request) {
    var medicalHistory models.MedicalHistory
    if err := json.NewDecoder(r.Body).Decode(&medicalHistory); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    updated, err := h.service.UpdateMedicalHistory(r.Context(), &medicalHistory)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(updated)
} 


func (h *MedicalHistoryHandler) DeleteMedicalHistory(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idStr := vars["id"]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid medical history ID", http.StatusBadRequest)
        return
    }

    // Decode the request body to get additional context or confirmation (optional)
    var deleteRequest struct {
        Reason string `json:"reason,omitempty"`
    }
    
    // Only try to decode if there's a body
    if r.Body != nil && r.ContentLength > 0 {
        if err := json.NewDecoder(r.Body).Decode(&deleteRequest); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }
    }

    // Call service method to delete
    if err := h.service.DeleteMedicalHistory(r.Context(), id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Set content type and return success response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusNoContent)
    
    // Optionally, you can return a confirmation message
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Medical history deleted successfully",
        "id":      idStr,
    })
}