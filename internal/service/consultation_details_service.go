package service

import (
	"context"
	"errors"
	"fmt"
	"shifa/internal/models"
	"shifa/internal/repository"

	"github.com/sirupsen/logrus"
)

// ConsultationDetailsService manages consultation details business logic
type ConsultationDetailsService struct {
	detailsRepo repository.ConsultationDetailsRepository
	logger      *logrus.Logger
}

// NewConsultationDetailsService creates a new ConsultationDetailsService
func NewConsultationDetailsService(detailsRepo repository.ConsultationDetailsRepository, logger *logrus.Logger) *ConsultationDetailsService {
	return &ConsultationDetailsService{
		detailsRepo: detailsRepo,
		logger:      logger,
	}
}

// CreateDetails creates new consultation details
func (s *ConsultationDetailsService) CreateDetails(ctx context.Context, details models.ConsultationDetails) error {
	if err := s.validateDetails(details); err != nil {
		s.logger.WithError(err).Error("Invalid consultation details data")
		return fmt.Errorf("validation error: %w", err)
	}

	if err := s.detailsRepo.Create(ctx, &details); err != nil {
		s.logger.WithError(err).Error("Failed to create consultation details")
		return fmt.Errorf("failed to create consultation details: %w", err)
	}

	s.logger.Infof("Consultation details created successfully: ID=%d", details.ID)
	return nil
}

// GetDetailsByID retrieves consultation details by ID
func (s *ConsultationDetailsService) GetDetailsByID(ctx context.Context, id int) (*models.ConsultationDetails, error) {
	details, err := s.detailsRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get consultation details")
		return nil, fmt.Errorf("failed to get consultation details: %w", err)
	}
	return details, nil
}

// GetDetailsByConsultationID retrieves consultation details by consultation ID
func (s *ConsultationDetailsService) GetDetailsByConsultationID(ctx context.Context, consultationID int) (*models.ConsultationDetails, error) {
	details, err := s.detailsRepo.GetByConsultationID(ctx, consultationID)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get consultation details by consultation ID")
		return nil, fmt.Errorf("failed to get consultation details by consultation ID: %w", err)
	}
	return details, nil
}

// UpdateDetails updates existing consultation details
func (s *ConsultationDetailsService) UpdateDetails(ctx context.Context, details models.ConsultationDetails) error {
	if err := s.validateDetails(details); err != nil {
		s.logger.WithError(err).Error("Invalid consultation details data")
		return fmt.Errorf("validation error: %w", err)
	}

	if err := s.detailsRepo.Update(ctx, &details); err != nil {
		s.logger.WithError(err).Error("Failed to update consultation details")
		return fmt.Errorf("failed to update consultation details: %w", err)
	}

	s.logger.Infof("Consultation details updated successfully: ID=%d", details.ID)
	return nil
}

// DeleteDetails deletes consultation details
func (s *ConsultationDetailsService) DeleteDetails(ctx context.Context, id int) error {
	if err := s.detailsRepo.Delete(ctx, id); err != nil {
		s.logger.WithError(err).Error("Failed to delete consultation details")
		return fmt.Errorf("failed to delete consultation details: %w", err)
	}

	s.logger.Infof("Consultation details deleted successfully: ID=%d", id)
	return nil
}

// validateDetails validates consultation details
func (s *ConsultationDetailsService) validateDetails(details models.ConsultationDetails) error {
	if details.ConsultationID == 0 {
		return errors.New("consultation ID is required")
	}
	return nil
}
