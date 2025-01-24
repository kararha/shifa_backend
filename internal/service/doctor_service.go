package service  

import (  
    "context"
    "errors"  
    "fmt"  
    "shifa/internal/models"  
    "shifa/internal/repository/mysql"
    "github.com/sirupsen/logrus"  
)  

type DoctorService struct {
    doctorRepo *mysql.DoctorRepo
    logger     *logrus.Logger
}

// NewDoctorService now takes the concrete mysql.DoctorRepo
func NewDoctorService(doctorRepo *mysql.DoctorRepo, logger *logrus.Logger) *DoctorService {  
    return &DoctorService{  
        doctorRepo: doctorRepo,  
        logger:     logger,  
    }  
}  

func (s *DoctorService) RegisterDoctor(doctor models.Doctor) error {
    // Validate the doctor model
    if err := s.validateDoctor(doctor); err != nil {
        s.logger.WithError(err).Error("Invalid doctor data")
        return err
    }

    // Create a context
    ctx := context.Background()

    // Register the doctor with context
    if err := s.doctorRepo.Create(ctx, &doctor); err != nil {
        s.logger.WithError(err).Error("Failed to create doctor")
        return fmt.Errorf("failed to register doctor: %w", err)
    }

    s.logger.Infof("Doctor registered successfully: %v", doctor)
    return nil
}

func (s *DoctorService) GetByID(userID int) (models.Doctor, error) {
    ctx := context.Background()
    
    doctorPtr, err := s.doctorRepo.GetByUserID(ctx, userID)
    if err != nil {
        s.logger.WithError(err).Errorf("Failed to get doctor by ID: %d", userID)
        return models.Doctor{}, fmt.Errorf("failed to get doctor: %w", err)
    }

    // Check if doctorPtr is nil
    if doctorPtr == nil {
        return models.Doctor{}, fmt.Errorf("doctor not found")
    }

    return *doctorPtr, nil
}

func (s *DoctorService) UpdateDoctor(doctor models.Doctor) error {
    // Validate the doctor model
    if err := s.validateDoctor(doctor); err != nil {
        s.logger.WithError(err).Error("Invalid doctor data")
        return err
    }

    ctx := context.Background()
    
    if err := s.doctorRepo.Update(ctx, &doctor); err != nil {
        s.logger.WithError(err).Errorf("Failed to update doctor: %v", doctor)
        return fmt.Errorf("failed to update doctor: %w", err)
    }

    s.logger.Infof("Doctor updated successfully: %v", doctor)
    return nil
}

func (s *DoctorService) DeleteDoctor(userID int) error {
    ctx := context.Background()
    
    if err := s.doctorRepo.Delete(ctx, userID); err != nil {
        s.logger.WithError(err).Errorf("Failed to delete doctor with ID: %d", userID)
        return fmt.Errorf("failed to delete doctor: %w", err)
    }

    s.logger.Infof("Doctor deleted successfully: %d", userID)
    return nil
}

func (s *DoctorService) ListDoctorsBySpecialty(specialty string, limit, offset int) ([]models.Doctor, error) {
    ctx := context.Background()
    
    filter := models.DoctorFilter{
        Specialty: specialty,
    }
    
    doctorPtrs, err := s.doctorRepo.List(ctx, filter, offset, limit)
    if err != nil {
        s.logger.WithError(err).Errorf("Failed to list doctors by specialty: %s", specialty)
        return nil, fmt.Errorf("failed to list doctors: %w", err)
    }

    // Convert []*models.Doctor to []models.Doctor
    doctors := make([]models.Doctor, len(doctorPtrs))
    for i, doctorPtr := range doctorPtrs {
        doctors[i] = *doctorPtr
    }

    return doctors, nil
}

func (s *DoctorService) ListDoctorsByServiceType(serviceTypeID int, limit, offset int) ([]models.Doctor, error) {
    ctx := context.Background()
    
    filter := models.DoctorFilter{
        ServiceTypeID: serviceTypeID,
    }
    
    doctorPtrs, err := s.doctorRepo.List(ctx, filter, offset, limit)
    if err != nil {
        s.logger.WithError(err).Errorf("Failed to list doctors by service type: %d", serviceTypeID)
        return nil, fmt.Errorf("failed to list doctors: %w", err)
    }

    // Convert []*models.Doctor to []models.Doctor
    doctors := make([]models.Doctor, len(doctorPtrs))
    for i, doctorPtr := range doctorPtrs {
        doctors[i] = *doctorPtr
    }

    return doctors, nil
}

func (s *DoctorService) validateDoctor(doctor models.Doctor) error {  
    if doctor.Specialty == "" {  
        return errors.New("doctor specialty is required")  
    }  

    // Add any other validation logic here  
    if doctor.UserID == 0 {
        return errors.New("user ID is required")
    }

    return nil  
}

// Add this new method to the DoctorService struct
func (s *DoctorService) SearchDoctors(params models.DoctorSearchParams) ([]models.DoctorSearchResult, error) {
    ctx := context.Background()
    
    // Validate search parameters
    if params.RadiusKm != nil && *params.RadiusKm <= 0 {
        return nil, errors.New("radius must be positive")
    }
    if params.MinRating != nil && (*params.MinRating < 0 || *params.MinRating > 5) {
        return nil, errors.New("rating must be between 0 and 5")
    }

    // Both latitude and longitude must be provided if either is provided
    if (params.LocationLat != nil && params.LocationLng == nil) || 
       (params.LocationLat == nil && params.LocationLng != nil) {
        return nil, errors.New("both latitude and longitude must be provided for location search")
    }

    results, err := s.doctorRepo.SearchDoctors(ctx, params)
    if err != nil {
        s.logger.WithError(err).Error("Failed to search doctors")
        return nil, fmt.Errorf("failed to search doctors: %w", err)
    }

    return results, nil
}