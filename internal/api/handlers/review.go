package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"shifa/internal/models"
	"shifa/internal/service"
)

type ReviewHandler struct {
	service service.ReviewService
}

func NewReviewHandler(service service.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	var review models.Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Call the service to create the review
	err := h.service.CreateReview(r.Context(), &review)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the created review (or just a success message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(review) // Optionally return the review
}

func (h *ReviewHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	review, err := h.service.GetReviewByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review)
}

func (h *ReviewHandler) ListReviews(w http.ResponseWriter, r *http.Request) {
    // Parse doctorID
    doctorIDStr := r.URL.Query().Get("doctorId")
    doctorID, err := strconv.Atoi(doctorIDStr)
    if err != nil {
        http.Error(w, "Invalid doctor ID", http.StatusBadRequest)
        return
    }

    // Parse pagination parameters
    pageStr := r.URL.Query().Get("page")
    pageSizeStr := r.URL.Query().Get("pageSize")

    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1 // Default to first page if not specified or invalid
    }

    pageSize, err := strconv.Atoi(pageSizeStr)
    if err != nil || pageSize < 1 {
        pageSize = 10 // Default to 10 items per page if not specified or invalid
    }

    reviews, err := h.service.GetReviewsByDoctorID(r.Context(), doctorID, page, pageSize)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(reviews)
}

func (h *ReviewHandler) UpdateReview(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
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
	err = h.service.UpdateReview(r.Context(), &review)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(review) // Optionally return the updated review
}

func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteReview(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
