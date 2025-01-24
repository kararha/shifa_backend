package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
    "shifa/internal/models"
    "shifa/internal/service"
)

type DoctorAvailabilityHandler struct {
    service *service.DoctorAvailabilityService
}

func NewDoctorAvailabilityHandler(service *service.DoctorAvailabilityService) *DoctorAvailabilityHandler {
    return &DoctorAvailabilityHandler{service: service}
}

func (h *DoctorAvailabilityHandler) SetAvailability(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    doctorID, err := strconv.Atoi(vars["doctorId"])
    if err != nil {
        http.Error(w, "Invalid doctor ID", http.StatusBadRequest)
        return
    }

    var availability models.DoctorAvailability
    if err := json.NewDecoder(r.Body).Decode(&availability); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    availability.DoctorID = doctorID

    if err := h.service.SetAvailability(availability); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "Availability set successfully"})
}

func (h *DoctorAvailabilityHandler) GetAvailability(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    doctorID, err := strconv.Atoi(vars["doctorId"])
    if err != nil {
        http.Error(w, "Invalid doctor ID", http.StatusBadRequest)
        return
    }

    availabilities, err := h.service.ListAvailabilityByDoctor(r.Context(), doctorID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(availabilities)
}

func (h *DoctorAvailabilityHandler) UpdateAvailability(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid availability ID", http.StatusBadRequest)
        return
    }

    var availability models.DoctorAvailability
    if err := json.NewDecoder(r.Body).Decode(&availability); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    availability.ID = id

    if err := h.service.Update(r.Context(), availability); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "Availability updated successfully"})
}

func (h *DoctorAvailabilityHandler) DeleteAvailability(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid availability ID", http.StatusBadRequest)
        return
    }

    if err := h.service.Delete(r.Context(), id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}