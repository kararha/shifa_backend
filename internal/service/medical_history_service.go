package service

import (
    "context"
    "errors"
    "fmt"
    "shifa/internal/models"
    "shifa/internal/repository"
    "github.com/sirupsen/logrus"
)

type MedicalHistoryService struct {
    medicalHistoryRepo repository.MedicalHistoryRepository
    logger             *logrus.Logger
}

func NewMedicalHistoryService(
    medicalHistoryRepo repository.MedicalHistoryRepository, 
    logger *logrus.Logger,
) *MedicalHistoryService {
    return &MedicalHistoryService{
        medicalHistoryRepo: medicalHistoryRepo,
        logger:             logger,
    }
}

func (s *MedicalHistoryService) CreateMedicalHistory(
    ctx context.Context, 
    history *models.MedicalHistory,
) (*models.MedicalHistory, error) {
    if err := s.validateMedicalHistory(*history); err != nil {
        s.logger.WithError(err).Error("Invalid medical history data")
        return nil, err
    }

    if err := s.medicalHistoryRepo.Create(ctx, history); err != nil {
        s.logger.WithError(err).Error("Failed to create medical history")
        return nil, fmt.Errorf("failed to create medical history: %w", err)
    }

    return history, nil
}

func (s *MedicalHistoryService) GetMedicalHistoryByPatientID(
    ctx context.Context, 
    patientID int,
) ([]models.MedicalHistory, error) {
    histories, err := s.medicalHistoryRepo.GetByPatientID(ctx, patientID)
    if err != nil {
        s.logger.WithError(err).Errorf("Failed to get medical histories for patient ID: %d", patientID)
        return nil, fmt.Errorf("failed to get medical histories: %w", err)
    }

    // Convert []*models.MedicalHistory to []models.MedicalHistory if needed
    modelHistories := make([]models.MedicalHistory, len(histories))
    for i, history := range histories {
        modelHistories[i] = *history
    }

    return modelHistories, nil
}

func (s *MedicalHistoryService) UpdateMedicalHistory(
    ctx context.Context, 
    history *models.MedicalHistory,
) (*models.MedicalHistory, error) {
    if err := s.validateMedicalHistory(*history); err != nil {
        s.logger.WithError(err).Error("Invalid medical history data")
        return nil, err
    }

    if err := s.medicalHistoryRepo.Update(ctx, history); err != nil {
        s.logger.WithError(err).Error("Failed to update medical history")
        return nil, fmt.Errorf("failed to update medical history: %w", err)
    }

    return history, nil
}

func (s *MedicalHistoryService) DeleteMedicalHistory(
    ctx context.Context, 
    id int,
) error {
    if err := s.medicalHistoryRepo.Delete(ctx, id); err != nil {
        s.logger.WithError(err).Errorf("Failed to delete medical history with ID: %d", id)
        return fmt.Errorf("failed to delete medical history: %w", err)
    }
    return nil
}

func (s *MedicalHistoryService) validateMedicalHistory(history models.MedicalHistory) error {
    if history.PatientID == 0 {
        return errors.New("patient ID is required")
    }
    if history.ConditionName == "" {
        return errors.New("condition name is required")
    }
    if history.DiagnosisDate.IsZero() {
        return errors.New("diagnosis date is required")
    }
    return nil
}