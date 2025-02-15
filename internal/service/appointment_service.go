package service

import (
    "context"
    "errors"
    "fmt"
    
    "shifa/internal/models"
    "shifa/internal/repository"
    "github.com/sirupsen/logrus"
)

type AppointmentService struct {
    appointmentRepo      repository.AppointmentRepository
    doctorRepo          repository.DoctorRepository
    homeCareProviderRepo repository.HomeCareProviderRepository
    logger              *logrus.Logger
}

func NewAppointmentService(
    appointmentRepo repository.AppointmentRepository,
    doctorRepo repository.DoctorRepository,
    homeCareProviderRepo repository.HomeCareProviderRepository,
    logger *logrus.Logger,
) *AppointmentService {
    return &AppointmentService{
        appointmentRepo:      appointmentRepo,
        doctorRepo:          doctorRepo,
        homeCareProviderRepo: homeCareProviderRepo,
        logger:              logger,
    }
}

func (s *AppointmentService) validateAppointment(appointment *models.Appointment) error {
    if appointment.PatientID == 0 {
        return errors.New("patient ID is required")
    }
    if appointment.ProviderType == "" {
        return errors.New("provider type is required")
    }
    if appointment.DoctorID == nil && appointment.HomeCareProviderID == nil {
        return errors.New("either doctor ID or home care provider ID is required")
    }
    if appointment.AppointmentDate.IsZero() {
        return errors.New("appointment date is required")
    }
    if appointment.StartTime.IsZero() {
        return errors.New("start time is required")
    }
    if appointment.EndTime.IsZero() {
        return errors.New("end time is required")
    }
    if appointment.EndTime.Before(appointment.StartTime) {
        return errors.New("end time must be after start time")
    }
    if appointment.Status == "" {
        return errors.New("status is required")
    }
    return nil
}

func (s *AppointmentService) CreateAppointment(ctx context.Context, appointment *models.Appointment) (*models.Appointment, error) {
    // Validate basic appointment data
    if err := s.validateAppointment(appointment); err != nil {
        return nil, err
    }

    // Validate provider exists based on provider type
    if appointment.ProviderType == "doctor" && appointment.DoctorID != nil {
        // Check if doctor exists
        doctor, err := s.doctorRepo.GetByID(ctx, *appointment.DoctorID)
        if (err != nil) {
            return nil, fmt.Errorf("invalid doctor_id: %w", err)
        }
        if !doctor.IsAvailable || doctor.Status != "active" {
            return nil, fmt.Errorf("doctor is not available")
        }
    } else if appointment.ProviderType == "home_care_provider" && appointment.HomeCareProviderID != nil {
        // Check if home care provider exists
        provider, err := s.homeCareProviderRepo.GetByID(ctx, *appointment.HomeCareProviderID)
        if err != nil {
            return nil, fmt.Errorf("invalid home_care_provider_id: %w", err)
        }
        if !provider.IsAvailable || provider.Status != "active" {
            return nil, fmt.Errorf("home care provider is not available")
        }
    }

    err := s.appointmentRepo.Create(ctx, appointment)
    if err != nil {
        s.logger.WithError(err).Error("Failed to create appointment")
        return nil, fmt.Errorf("failed to create appointment: %w", err)
    }
    return appointment, nil
}

func (s *AppointmentService) GetAppointment(ctx context.Context, id int) (*models.Appointment, error) {
    appointment, err := s.appointmentRepo.GetByID(ctx, id)
    if err != nil {
        s.logger.WithError(err).Errorf("Failed to get appointment with ID: %d", id)
        return nil, fmt.Errorf("failed to get appointment with ID %d: %w", id, err)
    }
    return appointment, nil
}

func (s *AppointmentService) UpdateAppointment(ctx context.Context, appointment *models.Appointment) (*models.Appointment, error) {
    if err := s.validateAppointment(appointment); err != nil {
        return nil, err
    }
    err := s.appointmentRepo.Update(ctx, appointment)
    if err != nil {
        s.logger.WithError(err).Errorf("Failed to update appointment with ID: %d", appointment.ID)
        return nil, fmt.Errorf("failed to update appointment with ID %d: %w", appointment.ID, err)
    }
    return appointment, nil
}

func (s *AppointmentService) DeleteAppointment(ctx context.Context, id int) error {
    err := s.appointmentRepo.Delete(ctx, id)
    if err != nil {
        s.logger.WithError(err).Errorf("Failed to delete appointment with ID: %d", id)
        return fmt.Errorf("failed to delete appointment with ID %d: %w", id, err)
    }
    return nil
}

func (s *AppointmentService) ListAppointments(ctx context.Context, filter repository.AppointmentFilter, offset, limit int) ([]*models.Appointment, error) {
    appointments, err := s.appointmentRepo.List(ctx, filter, offset, limit)
    if err != nil {
        s.logger.WithError(err).Error("Failed to list appointments")
        return nil, fmt.Errorf("failed to list appointments: %w", err)
    }
    return appointments, nil
}

func (s *AppointmentService) ListAppointmentsByPatient(ctx context.Context, patientID, limit, offset int) ([]*models.Appointment, error) {
    appointments, err := s.appointmentRepo.GetByPatientID(ctx, patientID, limit, offset)
    if err != nil {
        s.logger.WithError(err).Errorf("Failed to get appointments for patient ID: %d", patientID)
        return nil, fmt.Errorf("failed to get appointments for patient ID %d: %w", patientID, err)
    }
    return appointments, nil
}

func (s *AppointmentService) GetAppointmentsByProvider(ctx context.Context, providerID int, providerType string) ([]*models.Appointment, error) {
    appointments, err := s.appointmentRepo.GetByProviderID(ctx, providerID, providerType)
    if err != nil {
        s.logger.WithError(err).Errorf("Failed to get appointments for provider ID: %d, type: %s", providerID, providerType)
        return nil, fmt.Errorf("failed to get appointments for provider ID %d, type %s: %w", providerID, providerType, err)
    }
    return appointments, nil
}
