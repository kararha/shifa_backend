package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "shifa/internal/models"
    "shifa/internal/service"
)

type PatientHandler struct {
    patientService *service.PatientService
}

func NewPatientHandler(patientService *service.PatientService) *PatientHandler {
    return &PatientHandler{patientService: patientService}
}

func (h *PatientHandler) CreatePatient(w http.ResponseWriter, r *http.Request) {
    var patient models.Patient
    if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err := h.patientService.RegisterPatient(r.Context(), patient)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(patient)
}

func (h *PatientHandler) GetPatient(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    patient, err := h.patientService.GetPatientByUserID(r.Context(), userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(patient)
}

func (h *PatientHandler) UpdatePatient(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    var patient models.Patient
    if err := json.NewDecoder(r.Body).Decode(&patient); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    patient.UserID = userID
    err = h.patientService.UpdatePatient(r.Context(), patient)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(patient)
}

func (h *PatientHandler) DeletePatient(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    err = h.patientService.DeletePatient(r.Context(), userID)
    if err != nil {
        http. Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (h *PatientHandler) ListPatients(w http.ResponseWriter, r *http.Request) {
    offset, limit := 0, 10 // Default values for pagination
    if r.URL.Query().Get("offset") != "" {
        offset, _ = strconv.Atoi(r.URL.Query().Get("offset"))
    }
    if r.URL.Query().Get("limit") != "" {
        limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
    }

    patients, err := h.patientService.ListPatients(r.Context(), offset, limit)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(patients)
}