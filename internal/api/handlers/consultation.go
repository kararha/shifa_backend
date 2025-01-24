// package handlers

// import (
//     "encoding/json"
//     "net/http"
//     "strconv"
//     "github.com/gorilla/mux"
//     "shifa/internal/models"
//     "shifa/internal/service"
// )

// type ConsultationHandler struct {
//     consultationService *service.ConsultationService
// }

// func NewConsultationHandler(consultationService *service.ConsultationService) *ConsultationHandler {
//     return &ConsultationHandler{consultationService: consultationService}
// }

// // StartConsultation handles starting a new consultation
// func (h *ConsultationHandler) StartConsultation(w http.ResponseWriter, r *http.Request) {
//     var consultation models.Consultation
//     if err := json.NewDecoder(r.Body).Decode(&consultation); err != nil {
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }

//     err := h.consultationService.StartConsultation(r.Context(), consultation)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     w.Header().Set("Content-Type", "application/json")
//     w.WriteHeader(http.StatusCreated)
//     json.NewEncoder(w).Encode(consultation)
// }

// // CompleteConsultation handles completing an existing consultation
// func (h *ConsultationHandler) CompleteConsultation(w http.ResponseWriter, r *http.Request) {
//     vars := mux.Vars(r)
//     consultationID, err := strconv.Atoi(vars["id"])
//     if err != nil {
//         http.Error(w, "Invalid consultation ID", http.StatusBadRequest)
//         return
//     }

//     var consultation models.Consultation
//     if err := json.NewDecoder(r.Body).Decode(&consultation); err != nil {
//         http.Error(w, err.Error(), http.StatusBadRequest)
//         return
//     }
//     consultation.ID = consultationID

//     err = h.consultationService.CompleteConsultation(r.Context(), consultation)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(consultation)
// }

// // // RegisterRoutes registers all the consultation routes
// // func (h *ConsultationHandler) RegisterRoutes(router *mux.Router) {
// //     router.HandleFunc("/consultations/start", h.StartConsultation).Methods("POST")
// //     router.HandleFunc("/consultations/{id}/complete", h.CompleteConsultation).Methods("PUT")
// // }



package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
    "shifa/internal/models"
    "shifa/internal/service"
)

type ConsultationHandler struct {
    consultationService *service.ConsultationService
}

func NewConsultationHandler(consultationService *service.ConsultationService) *ConsultationHandler {
    return &ConsultationHandler{consultationService: consultationService}
}

// StartConsultation handles starting a new consultation
func (h *ConsultationHandler) StartConsultation(w http.ResponseWriter, r *http.Request) {
    var consultation models.Consultation
    if err := json.NewDecoder(r.Body).Decode(&consultation); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err := h.consultationService.StartConsultation(r.Context(), consultation)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(consultation)
}

// CompleteConsultation handles completing an existing consultation
func (h *ConsultationHandler) CompleteConsultation(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    consultationID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid consultation ID", http.StatusBadRequest)
        return
    }

    var consultation models.Consultation
    if err := json.NewDecoder(r.Body).Decode(&consultation); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    consultation.ID = consultationID

    err = h.consultationService.CompleteConsultation(r.Context(), consultation)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(consultation)
}

// GetConsultation handles retrieving a single consultation by ID
func (h *ConsultationHandler) GetConsultation(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    consultationID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid consultation ID", http.StatusBadRequest)
        return
    }

    consultation, err := h.consultationService.GetByID(r.Context(), consultationID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(consultation)
}

// ListConsultations handles retrieving a list of consultations with optional filtering
func (h *ConsultationHandler) ListConsultations(w http.ResponseWriter, r *http.Request) {
    // Parse query parameters for filtering
    query := r.URL.Query()
    
    filter := models.ConsultationFilter{
        PatientID:        parseInt(query.Get("patient_id")),
        DoctorID:         parseInt(query.Get("doctor_id")),
        AppointmentID:    parseInt(query.Get("appointment_id")),
        ConsultationType: query.Get("consultation_type"),
        Status:           query.Get("status"),
    }

    // Parse pagination parameters
    page := parseInt(query.Get("page"))
    if page < 1 {
        page = 1
    }
    limit := parseInt(query.Get("limit"))
    if limit < 1 {
        limit = 10 // default limit
    }
    offset := (page - 1) * limit

    consultations, err := h.consultationService.List(r.Context(), filter, offset, limit)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(consultations)
}

// UpdateConsultation handles updating an existing consultation
func (h *ConsultationHandler) UpdateConsultation(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    consultationID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid consultation ID", http.StatusBadRequest)
        return
    }

    var consultation models.Consultation
    if err := json.NewDecoder(r.Body).Decode(&consultation); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    consultation.ID = consultationID

    err = h.consultationService.Update(r.Context(), consultation)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(consultation)
}

// DeleteConsultation handles deleting a consultation
func (h *ConsultationHandler) DeleteConsultation(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    consultationID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid consultation ID", http.StatusBadRequest)
        return
    }

    err = h.consultationService.Delete(r.Context(), consultationID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// Helper function to parse string to int
func parseInt(s string) int {
    if s == "" {
        return 0
    }
    i, err := strconv.Atoi(s)
    if err != nil {
        return 0
    }
    return i
}