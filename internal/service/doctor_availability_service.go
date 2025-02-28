package service

import (
	"context"
	"errors"
	"fmt"
	"shifa/internal/models"
	"shifa/internal/repository"

	"github.com/sirupsen/logrus"
)

type DoctorAvailabilityService struct {
	doctorAvailabilityRepo repository.DoctorAvailabilityRepository
	logger                 *logrus.Logger
}

func NewDoctorAvailabilityService(doctorAvailabilityRepo repository.DoctorAvailabilityRepository, logger *logrus.Logger) *DoctorAvailabilityService {
	return &DoctorAvailabilityService{
		doctorAvailabilityRepo: doctorAvailabilityRepo,
		logger:                 logger,
	}
}

// SetAvailability creates a new doctor availability record
func (s *DoctorAvailabilityService) SetAvailability(availability models.DoctorAvailability) error {
	// Validate the doctor availability model
	if err := s.validateDoctorAvailability(availability); err != nil {
		s.logger.WithError(err).Error("Invalid doctor availability data")
		return err
	}

	ctx := context.Background()
	if err := s.doctorAvailabilityRepo.Create(ctx, &availability); err != nil {
		s.logger.WithError(err).Error("Failed to set doctor availability")
		return fmt.Errorf("failed to set doctor availability: %w", err)
	}

	s.logger.Infof("Doctor availability set successfully: %v", availability)
	return nil
}

// ListAvailabilityByDoctor retrieves all availability slots for a given doctor
func (s *DoctorAvailabilityService) ListAvailabilityByDoctor(ctx context.Context, doctorID int) ([]*models.DoctorAvailability, error) {
	availabilities, err := s.doctorAvailabilityRepo.GetByDoctorID(ctx, doctorID)
	if err != nil {
		s.logger.WithError(err).Errorf("Failed to list doctor availability for doctor ID: %d", doctorID)
		return nil, fmt.Errorf("failed to list doctor availability: %w", err)
	}
	return availabilities, nil
}

// Add new service method to list all availability slots
func (s *DoctorAvailabilityService) ListAllAvailability(ctx context.Context) ([]*models.DoctorAvailability, error) {
	availabilities, err := s.doctorAvailabilityRepo.ListAllAvailability(ctx)
	if err != nil {
		s.logger.WithError(err).Error("Failed to list all doctor availability")
		return nil, err
	}
	return availabilities, nil
}

// Update updates an existing doctor availability record
func (s *DoctorAvailabilityService) Update(ctx context.Context, availability models.DoctorAvailability) error {
	if err := s.validateDoctorAvailability(availability); err != nil {
		s.logger.WithError(err).Error("Invalid doctor availability data")
		return err
	}

	if err := s.doctorAvailabilityRepo.Update(ctx, &availability); err != nil {
		s.logger.WithError(err).Error("Failed to update doctor availability")
		return fmt.Errorf("failed to update doctor availability: %w", err)
	}

	s.logger.Infof("Doctor availability updated successfully: %v", availability)
	return nil
}

// Delete removes a doctor availability record
func (s *DoctorAvailabilityService) Delete(ctx context.Context, id int) error {
	if err := s.doctorAvailabilityRepo.Delete(ctx, id); err != nil {
		s.logger.WithError(err).Errorf("Failed to delete doctor availability with ID: %d", id)
		return fmt.Errorf("failed to delete doctor availability: %w", err)
	}

	s.logger.Infof("Doctor availability deleted successfully: %d", id)
	return nil
}

// validateDoctorAvailability validates the doctor availability data
func (s *DoctorAvailabilityService) validateDoctorAvailability(availability models.DoctorAvailability) error {
	if availability.DoctorID == 0 {
		return errors.New("doctor ID is required")
	}

	if availability.StartTime.IsZero() || availability.EndTime.IsZero() {
		return errors.New("start time and end time are required")
	}

	if availability.StartTime.After(availability.EndTime) {
		return errors.New("start time must be before end time")
	}

	if availability.DayOfWeek < 0 || availability.DayOfWeek > 6 {
		return errors.New("day of week must be between 0 and 6")
	}

	return nil
}
