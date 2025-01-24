// package handlers

// import (
// 	// "context"
// 	"encoding/json"
// 	"net/http"
// 	"strconv"

// 	"github.com/gorilla/mux"
// 	"shifa/internal/models"
// 	"shifa/internal/service"
// )

// type ServiceTypeHandler struct {
// 	stService *service.ServiceTypeService
// }

// func NewServiceTypeHandler(stService *service.ServiceTypeService) *ServiceTypeHandler {
// 	return &ServiceTypeHandler{stService: stService}
// }

// // CreateServiceType handles creating a new service type
// func (h *ServiceTypeHandler) CreateServiceType(w http.ResponseWriter, r *http.Request) {
// 	var serviceType models.ServiceType
// 	if err := json.NewDecoder(r.Body).Decode(&serviceType); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// Create service type with context
// 	err := h.stService.CreateServiceType(r.Context(), serviceType)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// }

// // GetServiceType retrieves a service type by ID
// func (h *ServiceTypeHandler) GetServiceType(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	serviceTypeID, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, "Invalid service type ID", http.StatusBadRequest)
// 		return
// 	}

// 	serviceType, err := h.stService.GetServiceTypeByID(r.Context(), serviceTypeID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(serviceType)
// }

// // UpdateServiceType handles updating a service type
// func (h *ServiceTypeHandler) UpdateServiceType(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	serviceTypeID, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, "Invalid service type ID", http.StatusBadRequest)
// 		return
// 	}

// 	var serviceType models.ServiceType
// 	if err := json.NewDecoder(r.Body).Decode(&serviceType); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	serviceType.ID = serviceTypeID
// 	err = h.stService.UpdateServiceType(r.Context(), serviceType)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// }

// // DeleteServiceType deletes a service type by ID
// func (h *ServiceTypeHandler) DeleteServiceType(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	serviceTypeID, err := strconv.Atoi(vars["id"])
// 	if err != nil {
// 		http.Error(w, "Invalid service type ID", http.StatusBadRequest)
// 		return
// 	}

// 	err = h.stService.DeleteServiceType(r.Context(), serviceTypeID)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.WriteHeader(http.StatusNoContent)
// }

// // ListServiceTypes lists all service types
// func (h *ServiceTypeHandler) ListServiceTypes(w http.ResponseWriter, r *http.Request) {
// 	serviceTypes, err := h.stService.ListServiceTypes(r.Context())
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(serviceTypes)
// }



// File: internal/api/handlers/service_type_handler.go
package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    "shifa/internal/repository"
    "shifa/internal/service"
)

type ServiceTypeHandler struct {
    service *service.ServiceTypeService
}

func NewServiceTypeHandler(service *service.ServiceTypeService) *ServiceTypeHandler {
    return &ServiceTypeHandler{service: service}
}

func (h *ServiceTypeHandler) CreateServiceType(w http.ResponseWriter, r *http.Request) {
    var serviceType repository.ServiceType
    if err := json.NewDecoder(r.Body).Decode(&serviceType); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if err := h.service.CreateServiceType( r.Context(), &serviceType); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(serviceType)
}

func (h *ServiceTypeHandler) ListServiceTypes(w http.ResponseWriter, r *http.Request) {
    serviceTypes, err := h.service.ListServiceTypes(r.Context())
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(serviceTypes)
}

func (h *ServiceTypeHandler) GetServiceType(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    serviceType, err := h.service.GetServiceTypeByID(r.Context(), id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(serviceType)
}

func (h *ServiceTypeHandler) UpdateServiceType(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    var serviceType repository.ServiceType
    if err := json.NewDecoder(r.Body).Decode(&serviceType); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    serviceType.ID = id

    if err := h.service.UpdateServiceType(r.Context(), &serviceType); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (h *ServiceTypeHandler) DeleteServiceType(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(mux.Vars(r)["id"])
    if err != nil {
        http.Error(w, "Invalid ID", http.StatusBadRequest)
        return
    }

    if err := h.service.DeleteServiceType(r.Context(), id); err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}