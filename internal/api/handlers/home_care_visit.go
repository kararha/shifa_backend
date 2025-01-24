// package handlers

// import (
//     "encoding/json"
//     "net/http"
//     "strconv"
//     "time"
//     "errors"
//     "github.com/gorilla/mux"
//     "shifa/internal/models"
//     "shifa/internal/service"
// )

// type HomeCareVisitHandler struct {
//     homeCareVisitService *service.HomeCareVisitService
// }

// func NewHomeCareVisitHandler(service *service.HomeCareVisitService) *HomeCareVisitHandler {
//     return &HomeCareVisitHandler{
//         homeCareVisitService: service,
//     }
// }

// // ScheduleHomeCareVisit handles creating a new home care visit
// func (h *HomeCareVisitHandler) ScheduleHomeCareVisit(w http.ResponseWriter, r *http.Request) {
//     var visit models.HomeCareVisit
//     if err := json.NewDecoder(r.Body).Decode(&visit); err != nil {
//         http.Error(w, "Invalid request body", http.StatusBadRequest)
//         return
//     }

//     err := h.homeCareVisitService.ScheduleHomeCareVisit(r.Context(), &visit)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     w.Header().Set("Content-Type", "application/json")
//     w.WriteHeader(http.StatusCreated)
//     json.NewEncoder(w).Encode(visit)
// }

// // GetHomeCareVisit handles retrieving a specific home care visit
// func (h *HomeCareVisitHandler) GetHomeCareVisit(w http.ResponseWriter, r *http.Request) {
//     vars := mux.Vars(r)
//     visitID, err := strconv.Atoi(vars["id"])
//     if err != nil {
//         http.Error(w, "Invalid visit ID", http.StatusBadRequest)
//         return
//     }

//     visit, err := h.homeCareVisitService.GetVisitDetails(r.Context(), visitID)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(visit)
// }

// // UpdateHomeCareVisit handles updating an existing home care visit
// func (h *HomeCareVisitHandler) UpdateHomeCareVisit(w http.ResponseWriter, r *http.Request) {
//     vars := mux.Vars(r)
//     visitID, err := strconv.Atoi(vars["id"])
//     if err != nil {
//         http.Error(w, "Invalid visit ID", http.StatusBadRequest)
//         return
//     }

//     var visit models.HomeCareVisit
//     if err := json.NewDecoder(r.Body).Decode(&visit); err != nil {
//         http.Error(w, "Invalid request body", http.StatusBadRequest)
//         return
//     }
//     visit.ID = visitID

//     err = h.homeCareVisitService.UpdateHomeCareVisit(r.Context(), &visit)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(visit)
// }

// // DeleteHomeCareVisit handles deleting a home care visit
// func (h *HomeCareVisitHandler) DeleteHomeCareVisit(w http.ResponseWriter, r *http.Request) {
//     vars := mux.Vars(r)
//     visitID, err := strconv.Atoi(vars["id"])
//     if err != nil {
//         http.Error(w, "Invalid visit ID", http.StatusBadRequest)
//         return
//     }

//     err = h.homeCareVisitService.DeleteHomeCareVisit(r.Context(), visitID)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     w.WriteHeader(http.StatusNoContent)
// }

// // ListHomeCareVisits handles listing home care visits with filters
// func (h *HomeCareVisitHandler) ListHomeCareVisits(w http.ResponseWriter, r *http.Request) {
//     filter := models.HomeCareVisitFilter{
//         // Parse query parameters and populate filter
//         PatientID:  parseIntParam(r, "patientId"),
//         ProviderID: parseIntParam(r, "providerId"),
//         Status:     r.URL.Query().Get("status"),
//         StartDate:  parseTimeParam(r, "startDate"),
//         EndDate:    parseTimeParam(r, "endDate"),
//     }

//     visits, err := h.homeCareVisitService.ListHomeCareVisits(r.Context(), filter)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(visits)
// }

// // GetHomeCareVisitsByPatient handles retrieving visits for a specific patient
// func (h *HomeCareVisitHandler) GetHomeCareVisitsByPatient(w http.ResponseWriter, r *http.Request) {
//     vars := mux.Vars(r)
//     patientID, err := strconv.Atoi(vars["patientId"])
//     if err != nil {
//         http.Error(w, "Invalid patient ID", http.StatusBadRequest)
//         return
//     }

//     visits, err := h.homeCareVisitService.GetHomeCareVisitsByPatient(r.Context(), patientID)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(visits)
// }

// // GetHomeCareVisitsByProvider handles retrieving visits for a specific provider
// func (h *HomeCareVisitHandler) GetHomeCareVisitsByProvider(w http.ResponseWriter, r *http.Request) {
//     vars := mux.Vars(r)
//     providerID, err := strconv.Atoi(vars["providerId"])
//     if err != nil {
//         http.Error(w, "Invalid provider ID", http.StatusBadRequest)
//         return
//     }

//     visits, err := h.homeCareVisitService.GetHomeCareVisitsByProvider(r.Context(), providerID)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(visits)
// }

// // GetHomeCareVisitsByDateRange handles retrieving visits within a date range
// func (h *HomeCareVisitHandler) GetHomeCareVisitsByDateRange(w http.ResponseWriter, r *http.Request) {
//     startDate := parseTimeParam(r, "startDate")
//     endDate := parseTimeParam(r, "endDate")

//     if startDate.IsZero() || endDate.IsZero() {
//         http.Error(w, "Invalid date parameters", http.StatusBadRequest)
//         return
//     }

//     visits, err := h.homeCareVisitService.GetHomeCareVisitsByDateRange(r.Context(), startDate, endDate)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(visits)
// }

// // CancelHomeCareVisit handles canceling a home care visit
// func (h *HomeCareVisitHandler) CancelHomeCareVisit(w http.ResponseWriter, r *http.Request) {
//     vars := mux.Vars(r)
//     visitID, err := strconv.Atoi(vars["id"])
//     if err != nil {
//         http.Error(w, "Invalid visit ID", http.StatusBadRequest)
//         return
//     }

//     err = h.homeCareVisitService.CancelVisit(r.Context(), visitID)
//     if err != nil {
//         http.Error(w, err.Error(), http.StatusInternalServerError)
//         return
//     }

//     w.WriteHeader(http.StatusOK)
// }

// // RegisterRoutes registers all the home care visit routes
// func (h *HomeCareVisitHandler) RegisterRoutes(router *mux.Router) {
//     router.HandleFunc("/home-care-visits", h.ScheduleHomeCareVisit).Methods("POST")
//     router.HandleFunc("/home-care-visits", h.ListHomeCareVisits).Methods("GET")
//     router.HandleFunc("/home-care-visits/{id}", h.GetHomeCareVisit).Methods("GET")
//     router.HandleFunc("/home-care-visits/{id}", h.UpdateHomeCareVisit).Methods("PUT")
//     router.HandleFunc("/home-care-visits/{id}", h.DeleteHomeCareVisit).Methods("DELETE")
//     router.HandleFunc("/home-care-visits/{id}/cancel", h.CancelHomeCareVisit).Methods("PUT")
//     router.HandleFunc("/home-care-visits/patient/{patientId}", h.GetHomeCareVisitsByPatient).Methods("GET")
//     router.HandleFunc("/home-care-visits/provider/{providerId}", h.GetHomeCareVisitsByProvider).Methods("GET")
//     router.HandleFunc("/home-care-visits/date-range", h.GetHomeCareVisitsByDateRange).Methods("GET")
// }

// // Helper functions
// func parseIntParam(r *http.Request, param string) int {
//     value := r.URL.Query().Get(param)
//     if value == "" {
//         return 0
//     }
    
//     intValue, err := strconv.Atoi(value)
//     if err != nil {
//         return 0
//     }
//     return intValue
// }

// func parseTimeParam(r *http.Request, param string) time.Time {
//     value := r.URL.Query().Get(param)
//     if value == "" {
//         return time.Time{}
//     }
    
//     // Try parsing with different formats
//     layouts := []string{
//         time.RFC3339,
//         "2006-01-02",
//         "2006-01-02 15:04:05",
//     }

//     for _, layout := range layouts {
//         if t, err := time.Parse(layout, value); err == nil {
//             return t
//         }
//     }
    
//     return time.Time{}
// }

// // Response wrapper functions for consistent API responses
// func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
//     response, err := json.Marshal(payload)
//     if err != nil {
//         w.WriteHeader(http.StatusInternalServerError)
//         w.Write([]byte("Internal Server Error"))
//         return
//     }

//     w.Header().Set("Content-Type", "application/json")
//     w.WriteHeader(code)
//     w.Write(response)
// }

// func respondWithError(w http.ResponseWriter, code int, message string) {
//     respondWithJSON(w, code, map[string]string{"error": message})
// }

// // Middleware for common handler operations
// func (h *HomeCareVisitHandler) withLogging(next http.HandlerFunc) http.HandlerFunc {
//     return func(w http.ResponseWriter, r *http.Request) {
//         // Log request details
//         // You can add your logging logic here
//         next.ServeHTTP(w, r)
//     }
// }

// func (h *HomeCareVisitHandler) withAuth(next http.HandlerFunc) http.HandlerFunc {
//     return func(w http.ResponseWriter, r *http.Request) {
//         // Add authentication check logic here
//         // For example:
//         // token := r.Header.Get("Authorization")
//         // if !validateToken(token) {
//         //     respondWithError(w, http.StatusUnauthorized, "Unauthorized")
//         //     return
//         // }
//         next.ServeHTTP(w, r)
//     }
// }

// // Additional helper methods for validation
// func validateVisitRequest(visit *models.HomeCareVisit) error {
//     if visit.PatientID == 0 {
//         return errors.New("patient ID is required")
//     }
//     if visit.ProviderID == 0 {
//         return errors.New("provider ID is required")
//     }
//     if visit.VisitDate.IsZero() {
//         return errors.New("visit date is required")
//     }
//     return nil
// }



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
        AppointmentID: parseIntParam(r, "appointmentId"),
        Status:        r.URL.Query().Get("status"),
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