package service

import (
    "context"
    "errors"
    "time"
    "github.com/sirupsen/logrus"
    "shifa/internal/repository"
)

// Define the PaymentService interface
type PaymentService interface {
    CreatePayment(ctx context.Context, payment *repository.Payment) error
    GetPaymentByID(ctx context.Context, id int) (*repository.Payment, error)
    UpdatePaymentStatus(ctx context.Context, id int, status string) error
    GetPaymentByConsultationID(ctx context.Context, consultationID int) (*repository.Payment, error)
    GetPaymentByHomeCareVisitID(ctx context.Context, homeCareVisitID int) (*repository.Payment, error)
    ProcessRefund(ctx context.Context, paymentID int) error
}

type paymentService struct {
    paymentRepo repository.PaymentRepository
    logger      *logrus.Entry // Change to *logrus.Entry
}

// NewPaymentService constructor
func NewPaymentService(paymentRepo repository.PaymentRepository, logger *logrus.Logger) PaymentService {
    // Create a service-specific logger with context
    serviceLogger := logger.WithField("service", "payment_service")
    
    return &paymentService{
        paymentRepo: paymentRepo,
        logger:      serviceLogger, // Use the *logrus.Entry
    }
}

func (s *paymentService) CreatePayment(ctx context.Context, payment *repository.Payment) error {
    if payment.Amount <= 0 {
        return errors.New("invalid payment amount")
    }

    payment.Status = "pending"
    now := time.Now()
    payment.PaymentDate = &now

    err := s.paymentRepo.Create(ctx, payment)
    if err != nil {
        s.logger.Error("Failed to create payment", "error", err)
        return errors.New("failed to process payment")
    }

    return nil
}

func (s *paymentService) GetPaymentByID(ctx context.Context, id int) (*repository.Payment, error) {
    payment, err := s.paymentRepo.GetByID(ctx, id)
    if err != nil {
        s.logger.Error("Failed to get payment", "error", err, "paymentID", id)
        return nil, errors.New("payment not found")
    }
    return payment, nil
}

func (s *paymentService) UpdatePaymentStatus(ctx context.Context, id int, status string) error {
    err := s.paymentRepo.UpdateStatus(ctx, id, status)
    if err != nil {
        s.logger.Error("Failed to update payment status", 
            "error", err, 
            "paymentID", id, 
            "status", status)
        return errors.New("failed to update payment status")
    }
    return nil
}

func (s *paymentService) GetPaymentByConsultationID(ctx context.Context, consultationID int) (*repository.Payment, error) {
    payment, err := s.paymentRepo.GetByConsultationID(ctx, consultationID)
    if err != nil {
        s.logger.Error("Failed to get payment by consultation ID", 
            "error", err, 
            "consultationID", consultationID)
        return nil, errors.New("payment not found")
    }
    return payment, nil
}

func (s *paymentService) GetPaymentByHomeCareVisitID(ctx context.Context, homeCareVisitID int) (*repository.Payment, error) {
    payment, err := s.paymentRepo.GetByHomeCareVisitID(ctx, homeCareVisitID)
    if err != nil {
        s.logger.Error("Failed to get payment by home care visit ID", 
            "error", err, 
            "homeCareVisitID", homeCareVisitID)
        return nil, errors.New("payment not found")
    }
    return payment, nil
}

func (s *paymentService) ProcessRefund(ctx context.Context, paymentID int) error {
    payment, err := s.paymentRepo.GetByID(ctx, paymentID)
    if err != nil {
        s.logger.Error("Failed to get payment for refund", 
            "error", err, 
            "paymentID", paymentID)
        return errors.New("payment not found")
    }

    if payment.Status != "paid" {
        return errors.New("payment is not eligible for refund")
    }

    // Implement refund logic here (e.g., integrating with a payment gateway)

    payment.Status = "refunded"
    now := time.Now()
    payment.RefundDate = &now

    err = s.paymentRepo.UpdateStatus(ctx, paymentID, payment.Status)
    if err != nil {
        s.logger.Error("Failed to update payment status after refund", 
            "error", err, 
            "paymentID", paymentID)
        return errors.New("failed to process refund")
    }

    return nil
}