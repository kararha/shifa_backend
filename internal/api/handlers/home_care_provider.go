// handlers/home_care_provider_handler.go
package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "github.com/gorilla/mux"
    "shifa/internal/models"
    "shifa/internal/service"
)

type HomeCareProviderHandler struct {
    hcpService *service.HomeCareProviderService
}

func NewHomeCareProviderHandler(hcpService *service.HomeCareProviderService) *HomeCareProviderHandler {
    return &HomeCareProviderHandler{hcpService: hcpService}
}

// CreateHomeCareProvider handles the creation of a new home care provider
func (h *HomeCareProviderHandler) CreateHomeCareProvider(w http.ResponseWriter, r *http.Request) {
    var hcp models.HomeCareProvider
    if err := json.NewDecoder(r.Body).Decode(&hcp); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err := h.hcpService.CreateHomeCareProvider(&hcp)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(hcp)
}

// GetHomeCareProvider handles retrieving a single home care provider
func (h *HomeCareProviderHandler) GetHomeCareProvider(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    hcp, err := h.hcpService.GetHomeCareProviderByID(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(hcp)
}

// UpdateHomeCareProvider handles updating an existing home care provider
func (h *HomeCareProviderHandler) UpdateHomeCareProvider(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var hcp models.HomeCareProvider
    if err := json.NewDecoder(r.Body).Decode(&hcp); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    hcp.UserID = id
    err = h.hcpService.UpdateHomeCareProvider(&hcp)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(hcp)
}

// DeleteHomeCareProvider handles deletion of a home care provider
func (h *HomeCareProviderHandler) DeleteHomeCareProvider(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    err = h.hcpService.DeleteHomeCareProvider(id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

// ListHomeCareProviders handles listing all home care providers
func (h *HomeCareProviderHandler) ListHomeCareProviders(w http.ResponseWriter, r *http.Request) {
    hcps, err := h.hcpService.ListHomeCareProviders()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(hcps)
}

// SearchHomeCareProviders handles searching for home care providers
func (h *HomeCareProviderHandler) SearchHomeCareProviders(w http.ResponseWriter, r *http.Request) {
    serviceTypeID, err := strconv.Atoi(r.URL.Query().Get("service_type_id"))
    if err != nil {
        http.Error(w, "Invalid service type ID", http.StatusBadRequest)
        return
    }

    limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
    if err != nil {
        http.Error(w, "Invalid limit", http.StatusBadRequest)
        return
    }

    offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
        http.Error(w, "Invalid offset", http.StatusBadRequest)
        return
    }

    providers, err := h.hcpService.GetProvidersByServiceType(serviceTypeID, limit, offset)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(providers)
}

// GetHomeCareProviderByUserID handles retrieving a provider by user ID
func (h *HomeCareProviderHandler) GetHomeCareProviderByUserID(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := strconv.Atoi(vars["user_id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    provider, err := h.hcpService.GetHomeCareProviderByUserID(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(provider)
}

// Helper method to send JSON response
func sendJSONResponse(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if data != nil {
        if err := json.NewEncoder(w).Encode(data); err != nil {
            http.Error(w, "Error encoding response", http.StatusInternalServerError)
        }
    }
}

// Helper method to send error response
func sendErrorResponse(w http.ResponseWriter, status int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// RegisterRoutes registers all the home care provider routes
func (h *HomeCareProviderHandler) RegisterRoutes(router *mux.Router) {
    router.HandleFunc("/providers", h.CreateHomeCareProvider).Methods("POST")
    router.HandleFunc("/providers", h.ListHomeCareProviders).Methods("GET")
    router.HandleFunc("/providers/{id}", h.GetHomeCareProvider).Methods("GET")
    router.HandleFunc("/providers/{id}", h.UpdateHomeCareProvider).Methods("PUT")
    router.HandleFunc("/providers/{id}", h.DeleteHomeCareProvider).Methods("DELETE")
    router.HandleFunc("/providers/user/{user_id}", h.GetHomeCareProviderByUserID).Methods("GET")
    router.HandleFunc("/providers/search", h.SearchHomeCareProviders).Methods("GET")
}