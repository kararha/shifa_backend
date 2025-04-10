package handlers

import (
	// "context"
	"encoding/json"
	"net/http"
	"shifa/internal/models"
	"shifa/internal/repository"
	"shifa/internal/service"
	"strconv"

	"github.com/gorilla/mux"
)

type PaymentHandler struct {
	paymentService service.PaymentService // Correctly use the interface type
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var payment models.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check that at least one reference (consultation_id or home_care_visit_id) is provided
	if payment.ConsultationID <= 0 && payment.HomeCareVisitID <= 0 {
		http.Error(w, "either consultation_id or home_care_visit_id must be provided", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Convert models.Payment to repository.Payment
	repoPayment := repository.Payment{
		Amount:          payment.Amount,
		Status:          payment.Status,
		PaymentDate:     payment.PaymentDate,
		RefundDate:      payment.RefundDate,
		ConsultationID:  payment.ConsultationID,
		HomeCareVisitID: payment.HomeCareVisitID,
	}

	err := h.paymentService.CreatePayment(ctx, &repoPayment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update the model with all values from repository payment,
	// including those set by the service (ID, status, payment_date)
	payment.ID = repoPayment.ID
	payment.Status = repoPayment.Status
	payment.PaymentDate = repoPayment.PaymentDate
	payment.RefundDate = repoPayment.RefundDate

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payment)
}

func (h *PaymentHandler) GetPayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	paymentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	payment, err := h.paymentService.GetPaymentByID(ctx, paymentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}

func (h *PaymentHandler) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	paymentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}
	var payment models.Payment
	if err := json.NewDecoder(r.Body).Decode(&payment); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	payment.ID = paymentID
	ctx := r.Context()
	err = h.paymentService.UpdatePaymentStatus(ctx, paymentID, payment.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}

// In handlers/payment_handler.go, add these methods:

func (h *PaymentHandler) GetPaymentByConsultationID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	consultationID, err := strconv.Atoi(vars["consultationId"])
	if err != nil {
		http.Error(w, "Invalid consultation ID", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	payment, err := h.paymentService.GetPaymentByConsultationID(ctx, consultationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}

func (h *PaymentHandler) GetPaymentByHomeCareVisitID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	visitID, err := strconv.Atoi(vars["visitId"])
	if err != nil {
		http.Error(w, "Invalid home care visit ID", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	payment, err := h.paymentService.GetPaymentByHomeCareVisitID(ctx, visitID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)
}

func (h *PaymentHandler) ProcessRefund(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	paymentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	err = h.paymentService.ProcessRefund(ctx, paymentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Refund processed successfully"})
}
