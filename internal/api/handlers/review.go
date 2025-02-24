package handlers

import (
	"encoding/json"
	"net/http"
	"shifa/internal/models"
	"shifa/internal/service"
	"strconv"

	"github.com/gorilla/mux"
)

type ReviewHandler struct {
	reviewService service.ReviewService
}

func NewReviewHandler(service service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: service}
}

// ListReviews handles both doctor and provider reviews
func (h *ReviewHandler) ListReviews(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	doctorIDStr := r.URL.Query().Get("doctor_id")
	providerIDStr := r.URL.Query().Get("home_care_provider_id")

	// Set default page and pageSize
	page := 1
	pageSize := 10

	var reviews interface{}
	var err error

	if doctorIDStr != "" {
		doctorID, err := strconv.Atoi(doctorIDStr)
		if err != nil {
			http.Error(w, "Invalid doctor ID", http.StatusBadRequest)
			return
		}
		reviews, err = h.reviewService.GetReviewsByDoctorID(r.Context(), doctorID, page, pageSize)
	} else if providerIDStr != "" {
		providerID, err := strconv.Atoi(providerIDStr)
		if err != nil {
			http.Error(w, "Invalid provider ID", http.StatusBadRequest)
			return
		}
		reviews, err = h.reviewService.GetReviewsByHomeCareProviderID(r.Context(), providerID, page, pageSize)
	} else {
		http.Error(w, "Missing doctor_id or home_care_provider_id parameter", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	var review models.Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.reviewService.CreateReview(r.Context(), &review); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review)
}

func (h *ReviewHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	review, err := h.reviewService.GetReviewByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(review)
}

func (h *ReviewHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	var review models.Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	review.ID = id
	if err := h.reviewService.UpdateReview(r.Context(), &review); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(review)
}

func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	if err := h.reviewService.DeleteReview(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
