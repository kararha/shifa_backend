package service

import (
    "context"
    "errors"
    "fmt"
    // "time"
    "shifa/internal/models"
    "shifa/internal/repository"
    "github.com/sirupsen/logrus" // Make sure to import Logrus
)

type HomeCareVisitService struct {
    Repo   repository.HomeCareVisitRepository
    Logger *logrus.Logger // Using Logrus directly
}

func NewHomeCareVisitService(repo repository.HomeCareVisitRepository, logger *logrus.Logger) *HomeCareVisitService {
    return &HomeCareVisitService{
        Repo:   repo,
        Logger: logger,
    }
}

// ScheduleHomeCareVisit creates a new home care visit
func (s *HomeCareVisitService) ScheduleHomeCareVisit(ctx context.Context, visit *models.HomeCareVisit) error {
    if visit.PatientID == 0 || visit.ProviderID == 0 || visit.Address == "" {
        s.Logger.Warn("Invalid home care visit data", logrus.Fields{"visit": visit})
        return errors.New("invalid home care visit data: missing required fields")
    }

    // Set default status if not provided
    if visit.Status == "" {
        visit.Status = "scheduled"
    }

    err := s.Repo.Create(ctx, visit)
    if err != nil {
        s.Logger.Error("Failed to schedule home care visit", logrus.Fields{"error": err})
        return fmt.Errorf("failed to schedule home care visit: %w", err)
    }

    s.Logger.Info("Home care visit scheduled successfully", logrus.Fields{"visitID": visit.ID})
    return nil
}

// GetVisitDetails retrieves a specific home care visit
func (s *HomeCareVisitService) GetVisitDetails(ctx context.Context, visitID int) (*models.HomeCareVisit, error) {
    visit, err := s.Repo.GetByID(ctx, visitID)
    if err != nil {
        s.Logger.Error("Failed to get home care visit details", logrus.Fields{"error": err, "visitID": visitID})
        return nil, fmt.Errorf("failed to get home care visit details: %w", err)
    }
    return visit, nil
}

// UpdateHomeCareVisit updates an existing home care visit
func (s *HomeCareVisitService) UpdateHomeCareVisit(ctx context.Context, visit *models.HomeCareVisit) error {
    if visit.ID == 0 {
        s.Logger.Warn("Invalid visit ID for update", logrus.Fields{"visit": visit})
        return errors.New("invalid visit ID")
    }

    err := s.Repo.Update(ctx, visit)
    if err != nil {
        s.Logger.Error("Failed to update home care visit", logrus.Fields{"error": err, "visitID": visit.ID})
        return fmt.Errorf("failed to update home care visit: %w", err)
    }

    s.Logger.Info("Home care visit updated successfully", logrus.Fields{"visitID": visit.ID})
    return nil
}

// DeleteHomeCareVisit deletes a home care visit
func (s *HomeCareVisitService) DeleteHomeCareVisit(ctx context.Context, id int) error {
    err := s.Repo.Delete(ctx, id)
    if err != nil {
        s.Logger.Error("Failed to delete home care visit", logrus.Fields{"error": err, "visitID": id})
        return fmt.Errorf("failed to delete home care visit: %w", err)
    }
    s.Logger.Info("Home care visit deleted successfully", logrus.Fields{"visitID": id})
    return nil
}

// ListHomeCareVisits returns filtered home care visits
func (s *HomeCareVisitService) ListHomeCareVisits(ctx context.Context, filter models.HomeCareVisitFilter) ([]models.HomeCareVisit, error) {
    visits, err := s.Repo.List(ctx, filter)
    if err != nil {
        s.Logger.Error("Failed to list home care visits", logrus.Fields{"error": err})
        return nil, fmt.Errorf("failed to list home care visits: %w", err)
    }
    s.Logger.Info("Retrieved home care visits", logrus.Fields{"count": len(visits)})
    return visits, nil
}