package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "time"
    "github.com/gorilla/mux"
    "shifa/internal/models"
    "shifa/internal/service"
    "shifa/internal/repository"
)

type AppointmentHandler struct {
    appointmentService *service.AppointmentService
}

func NewAppointmentHandler(appointmentService *service.AppointmentService) *AppointmentHandler {
    return &AppointmentHandler{appointmentService: appointmentService}
}



func (h *AppointmentHandler) CreateAppointment(w http.ResponseWriter, r *http.Request) {
    var appointment models.Appointment
    if err := json.NewDecoder(r.Body).Decode(&appointment); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    createdAppointment, err := h.appointmentService.CreateAppointment(r.Context(), &appointment)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(createdAppointment)
}

func (h *AppointmentHandler) GetAppointment(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    appointmentID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid appointment ID", http.StatusBadRequest)
        return
    }

    appointment, err := h.appointmentService.GetAppointment(r.Context(), appointmentID)
	if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(appointment)
}

func (h *AppointmentHandler) UpdateAppointment(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    appointmentID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid appointment ID", http.StatusBadRequest)
        return
    }

    var appointment models.Appointment
    if err := json.NewDecoder(r.Body).Decode(&appointment); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    appointment.ID = appointmentID
    updatedAppointment, err := h.appointmentService.UpdateAppointment(r.Context(), &appointment)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(updatedAppointment)
}

func (h *AppointmentHandler) DeleteAppointment(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    appointmentID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid appointment ID", http.StatusBadRequest)
        return
    }

    err = h.appointmentService.DeleteAppointment(r.Context(), appointmentID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// File: internal/api/handlers/appointment_handler.go
func (h *AppointmentHandler) ListAppointments(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // Parse query parameters
    offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    if limit == 0 {
        limit = 10 // default limit
    }

    // Create filter from query parameters
    filter := repository.AppointmentFilter{
        Status:       r.URL.Query().Get("status"),
        ProviderType: r.URL.Query().Get("providerType"),
    }

    // Parse date filters if provided
    if startDateStr := r.URL.Query().Get("startDate"); startDateStr != "" {
        startDate, err := time.Parse("2006-01-02", startDateStr)
        if err == nil {
            filter.StartDate = &startDate
        }
    }
    if endDateStr := r.URL.Query().Get("endDate"); endDateStr != "" {
        endDate, err := time.Parse("2006-01-02", endDateStr)
        if err == nil {
            filter.EndDate = &endDate
        }
    }

    appointments, err := h.appointmentService.ListAppointments(ctx, filter, offset, limit)
    if err != nil {
        http.Error(w, "Failed to list appointments", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(appointments)
}

func (h *AppointmentHandler) GetAppointmentsByProvider(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    providerID, err := strconv.Atoi(vars["providerId"])
    if err != nil {
        http.Error(w, "Invalid provider ID", http.StatusBadRequest)
        return
    }

    providerType := r.URL.Query().Get("type")
    if providerType != "doctor" && providerType != "home_care_provider" {
        http.Error(w, "Invalid provider type", http.StatusBadRequest)
        return
    }

    appointments, err := h.appointmentService.GetAppointmentsByProvider(r.Context(), providerID, providerType)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(appointments)
}


func (h *AppointmentHandler) ListAppointmentsByPatient(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    patientID, err := strconv.Atoi(vars["patientId"])
    if err != nil {
        http.Error(w, "Invalid patient ID", http.StatusBadRequest)
        return
    }

    // Parse pagination parameters
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
    if limit == 0 {
        limit = 10 // default limit
    }

    appointments, err := h.appointmentService.GetAppointmentsByPatient(r.Context(), patientID, limit, offset)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(appointments)
}