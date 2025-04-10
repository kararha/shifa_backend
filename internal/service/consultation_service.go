package service

import (
	"context"
	"errors"
	"fmt"
	"shifa/internal/models"
	"shifa/internal/repository"
	"time"

	"github.com/sirupsen/logrus"
)

type ConsultationRepository interface {
	Create(ctx context.Context, consultation *models.Consultation) error
	GetByID(ctx context.Context, id int) (*models.Consultation, error)
	GetByAppointmentID(ctx context.Context, appointmentID int) (*models.Consultation, error)
	Update(ctx context.Context, consultation *models.Consultation) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, filter models.ConsultationFilter, offset, limit int) ([]*models.Consultation, error)
}

type ConsultationService struct {
	consultationRepo repository.ConsultationRepository
	logger           *logrus.Logger
}

func NewConsultationService(consultationRepo repository.ConsultationRepository, logger *logrus.Logger) *ConsultationService {
	return &ConsultationService{
		consultationRepo: consultationRepo,
		logger:           logger,
	}
}

func (s *ConsultationService) StartConsultation(ctx context.Context, consultation models.Consultation) error {
	if err := s.validateConsultation(consultation); err != nil {
		s.logger.WithError(err).Error("Invalid consultation data")
		return fmt.Errorf("validation error: %w", err)
	}

	consultation.Status = "in_progress"
	now := time.Now()
	consultation.StartedAt = models.NullTime{Time: now, Valid: true}

	// Directly use the models.Consultation
	if err := s.consultationRepo.Create(ctx, &consultation); err != nil {
		s.logger.WithError(err).Error("Failed to start consultation")
		return fmt.Errorf("failed to start consultation: %w", err)
	}

	s.logger.Infof("Consultation started successfully: ID=%d", consultation.ID)
	return nil
}

func (s *ConsultationService) CompleteConsultation(ctx context.Context, consultation models.Consultation) error {
	existing, err := s.consultationRepo.GetByID(ctx, consultation.ID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to fetch consultation")
		return fmt.Errorf("failed to fetch consultation: %w", err)
	}

	if existing.Status != "in_progress" {
		return errors.New("consultation must be in progress to complete")
	}

	consultation.Status = "completed"
	now := time.Now()
	consultation.CompletedAt = models.NullTime{Time: now, Valid: true}

	// Directly use the models.Consultation
	if err := s.consultationRepo.Update(ctx, &consultation); err != nil {
		s.logger.WithError(err).Error("Failed to complete consultation")
		return fmt.Errorf("failed to complete consultation: %w", err)
	}

	s.logger.Infof("Consultation completed successfully: ID=%d", consultation.ID)
	return nil
}

func (s *ConsultationService) GetByID(ctx context.Context, id int) (*models.Consultation, error) {
	consultation, err := s.consultationRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get consultation")
		return nil, fmt.Errorf("failed to get consultation: %w", err)
	}
	return consultation, nil // No need to convert, already *models.Consultation
}

func (s *ConsultationService) List(ctx context.Context, filter models.ConsultationFilter, offset, limit int) ([]models.Consultation, error) {
	consultations, err := s.consultationRepo.List(ctx, filter, offset, limit) // This returns []*models.Consultation
	if err != nil {
		s.logger.WithError(err).Error("Failed to list consultations")
		return nil, fmt.Errorf("failed to list consultations: %w", err)
	}

	// Convert []*models.Consultation to []models.Consultation
	result := make([]models.Consultation, len(consultations))
	for i, c := range consultations {
		result[i] = *c // Dereference the pointer to get the value
	}

	return result, nil // Return the converted slice
}

func (s *ConsultationService) Update(ctx context.Context, consultation models.Consultation) error {
	if err := s.validateConsultation(consultation); err != nil {
		s.logger.WithError(err).Error("Invalid consultation data")
		return fmt.Errorf("validation error: %w", err)
	}

	// Directly use the models.Consultation
	if err := s.consultationRepo.Update(ctx, &consultation); err != nil {
		s.logger.WithError(err).Error("Failed to update consultation")
		return fmt.Errorf("failed to update consultation: %w", err)
	}

	s.logger.Infof("Consultation updated successfully: ID=%d", consultation.ID)
	return nil
}

func (s *ConsultationService) Delete(ctx context.Context, id int) error {
	if err := s.consultationRepo.Delete(ctx, id); err != nil {
		s.logger.WithError(err).Error("Failed to delete consultation")
		return fmt.Errorf("failed to delete consultation: %w", err)
	}

	s.logger.Infof("Consultation deleted successfully: ID=%d", id)
	return nil
}

func (s *ConsultationService) validateConsultation(consultation models.Consultation) error {
	if consultation.PatientID == 0 {
		return errors.New("patient ID is required")
	}
	if consultation.DoctorID == 0 {
		return errors.New("doctor ID is required")
	}
	return nil
}
