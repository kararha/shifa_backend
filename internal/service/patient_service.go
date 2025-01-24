// package service

// import (
//     "context"
//     "errors"
//     "fmt"
//     "shifa/internal/models"
//     "shifa/internal/repository"
//     "github.com/sirupsen/logrus"
// )

// // PatientService handles operations related to patients
// type PatientService struct {
//     patientRepo repository.PatientRepository
//     logger      *logrus.Logger
// }

// // NewPatientService creates a new instance of PatientService
// func NewPatientService(patientRepo repository.PatientRepository, logger *logrus.Logger) *PatientService {
//     return &PatientService{
//         patientRepo: patientRepo,
//         logger:      logger,
//     }
// }

// // RegisterPatient validates and registers a new patient
// func (s *PatientService) RegisterPatient(ctx context.Context, patient models.Patient) error {
//     if err := s.validatePatient(patient); err != nil {
//         s.logger.WithError(err).Error("Invalid patient data")
//         return err
//     }

//     repoPatient := convertToRepoPatient(patient)

//     if err := s.patientRepo.Create(ctx, &repoPatient); err != nil {
//         s.logger.WithError(err).Error("Failed to create patient")
//         return fmt.Errorf("failed to register patient: %w", err)
//     }

//     s.logger.Infof("Patient registered successfully: %v", patient)
//     return nil
// }

// // GetPatientById retrieves a patient by ID
// func (s *PatientService) GetPatientById(ctx context.Context, patientId int) (models.Patient, error) {
//     repoPatient, err := s.patientRepo.GetByID(ctx, patientId)
//     if err != nil {
//         s.logger.WithError(err).Errorf("Failed to get patient by ID: %d", patientId)
//         return models.Patient{}, fmt.Errorf("failed to get patient: %w", err)
//     }
//     return convertToModelPatient(*repoPatient), nil
// }

// // UpdatePatient validates and updates an existing patient
// func (s *PatientService) UpdatePatient(ctx context.Context, patient models.Patient) error {
//     if err := s.validatePatient(patient); err != nil {
//         s.logger.WithError(err).Error("Invalid patient data")
//         return err
//     }

//     repoPatient := convertToRepoPatient(patient)

//     if err := s.patientRepo.Update(ctx, &repoPatient); err != nil {
//         s.logger.WithError(err).Errorf("Failed to update patient: %v", patient)
//         return fmt.Errorf("failed to update patient: %w", err)
//     }

//     s.logger.Infof("Patient updated successfully: %v", patient)
//     return nil
// }

// // ListPatients retrieves all patients
// func (s *PatientService) ListPatients(ctx context.Context) ([]models.Patient, error) {
//     repoPatients, err := s.patientRepo.GetAll(ctx)
//     if err != nil {
//         s.logger.WithError(err).Error("Failed to list patients")
//         return nil, fmt.Errorf("failed to list patients: %w", err)
//     }

//     patients := make([]models.Patient, len(repoPatients))
//     for i, repoPatient := range repoPatients {
//         patients[i] = convertToModelPatient(*repoPatient)
//     }
//     return patients, nil
// }

// // validatePatient checks if the patient model is valid
// func (s *PatientService) validatePatient(patient models.Patient) error {
//     if patient.UserID == 0 {
//         return errors.New("user ID is required")
//     }
//     if patient.Phone == "" {
//         return errors.New("patient phone is required")
//     }
//     return nil
// }

// // convertToRepoPatient converts models.Patient to repository.Patient
// func convertToRepoPatient(patient models.Patient) repository.Patient {
//     return repository.Patient{
//         UserID: patient.UserID,
//         // Map other fields as necessary
//     }
// }

// // convertToModelPatient converts repository.Patient to models.Patient
// func convertToModelPatient(repoPatient repository.Patient) models.Patient {
//     return models.Patient{
//         UserID: repoPatient.UserID,
//         // Map other fields as necessary
//     }
// }

// // DeletePatient deletes a patient by ID
// func (s *PatientService) DeletePatient(ctx context.Context, patientID int) error {
//     err := s.patientRepo.Delete(ctx, patientID)
//     if err != nil {
//         s.logger.WithError(err).Errorf("Failed to delete patient with ID: %d", patientID)
//         return fmt.Errorf("failed to delete patient: %w", err)
//     }
//     s.logger.Infof("Patient deleted successfully with ID: %d", patientID)
//     return nil
// }

package service

import (
    "context"
    "errors"
    "fmt"
    "shifa/internal/models"
    "shifa/internal/repository"
    "github.com/sirupsen/logrus"
)

type PatientService struct {
    patientRepo repository.PatientRepository
    logger      *logrus.Logger
}

func NewPatientService(patientRepo repository.PatientRepository, logger *logrus.Logger) *PatientService {
    return &PatientService{
        patientRepo: patientRepo,
        logger:      logger,
    }
}

func (s *PatientService) RegisterPatient(ctx context.Context, patient models.Patient) error {
    if err := s.validatePatient(patient); err != nil {
        s.logger.WithError(err).Error("Invalid patient data")
        return err
    }

    if err := s.patientRepo.Create(ctx, &patient); err != nil {
        s.logger.WithError(err).Error("Failed to create patient")
        return fmt.Errorf("failed to register patient: %w", err)
    }

    s.logger.Infof("Patient registered successfully: %v", patient)
    return nil
}

func (s *PatientService) GetPatientByUserID(ctx context.Context, userID int) (models.Patient, error) {
    patient, err := s.patientRepo.GetByUserID(ctx, userID)
    if err != nil {
        s.logger.WithError(err).Errorf("Failed to get patient by user ID: %d", userID)
        return models.Patient{}, fmt.Errorf("failed to get patient: %w", err)
    }
    return *patient, nil
}

func (s *PatientService) UpdatePatient(ctx context.Context, patient models.Patient) error {
    if err := s.validatePatient(patient); err != nil {
        s.logger.WithError(err).Error("Invalid patient data")
        return err
    }

    if err := s.patientRepo.Update(ctx, &patient); err != nil {
        s.logger.WithError(err).Errorf("Failed to update patient: %v", patient)
        return fmt.Errorf("failed to update patient: %w", err)
    }

    s.logger.Infof("Patient updated successfully: %v", patient)
    return nil
}

func (s *PatientService) DeletePatient(ctx context.Context, userID int) error {
    err := s.patientRepo.Delete(ctx, userID)
    if err != nil {
        s.logger.WithError(err).Errorf("Failed to delete patient with user ID: %d", userID)
        return fmt.Errorf("failed to delete patient: %w", err)
    }
    s.logger.Infof("Patient deleted successfully with user ID: %d", userID)
    return nil
}

func (s *PatientService) ListPatients(ctx context.Context, offset, limit int) ([]models.Patient, error) {
    repoPatients, err := s.patientRepo.List(ctx, offset, limit)
    if err != nil {
        s.logger.WithError(err).Error("Failed to list patients")
        return nil, fmt.Errorf("failed to list patients: %w", err)
    }

    patients := make([]models.Patient, len(repoPatients))
    for i, repoPatient := range repoPatients {
        patients[i] = *repoPatient
    }
    return patients, nil
}

func (s *PatientService) validatePatient(patient models.Patient) error {
    if patient.UserID == 0 {
        return errors.New("user ID is required")
    }
    if patient.Phone == "" {
        return errors.New("patient phone is required")
    }
    if patient.DateOfBirth.IsZero() {
        return errors.New("date of birth is required")
    }
    if patient.Gender == "" {
        return errors.New("gender is required")
    }
    return nil
}