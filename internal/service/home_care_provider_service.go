package service

import (
	"context"
	"errors"
	"fmt"
	"shifa/internal/models"
	"shifa/internal/repository"

	"github.com/sirupsen/logrus"
)

type HomeCareProviderService struct {
	homeCareRepo repository.HomeCareProviderRepository
	logger       *logrus.Logger
}

// Constructor
func NewHomeCareProviderService(repo repository.HomeCareProviderRepository, logger *logrus.Logger) *HomeCareProviderService {
	return &HomeCareProviderService{
		homeCareRepo: repo,
		logger:       logger,
	}
}

func (s *HomeCareProviderService) CreateHomeCareProvider(provider *models.HomeCareProvider) error {
	if err := s.validateHomeCareProvider(*provider); err != nil {
		s.logger.WithError(err).Error("Invalid home care provider data")
		return err
	}

	ctx := context.Background()
	err := s.homeCareRepo.Create(ctx, provider)
	if err != nil {
		s.logger.WithError(err).Error("Failed to create home care provider")
		return fmt.Errorf("failed to create home care provider: %w", err)
	}

	return nil
}

func (s *HomeCareProviderService) GetHomeCareProviderByID(id int) (*models.HomeCareProvider, error) {
	ctx := context.Background()
	provider, err := s.homeCareRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.WithError(err).Errorf("Failed to get home care provider with ID: %d", id)
		return nil, fmt.Errorf("failed to get home care provider: %w", err)
	}
	return provider, nil
}

func (s *HomeCareProviderService) GetHomeCareProviderByUserID(userID int) (*models.HomeCareProvider, error) {
	ctx := context.Background()
	provider, err := s.homeCareRepo.GetByUserID(ctx, userID)
	if err != nil {
		s.logger.WithError(err).Errorf("Failed to get home care provider with UserID: %d", userID)
		return nil, fmt.Errorf("failed to get home care provider: %w", err)
	}
	return provider, nil
}

func (s *HomeCareProviderService) UpdateHomeCareProvider(provider *models.HomeCareProvider) error {
	if err := s.validateHomeCareProvider(*provider); err != nil {
		s.logger.WithError(err).Error("Invalid home care provider data")
		return err
	}

	ctx := context.Background()
	err := s.homeCareRepo.Update(ctx, provider)
	if err != nil {
		s.logger.WithError(err).Error("Failed to update home care provider")
		return fmt.Errorf("failed to update home care provider: %w", err)
	}

	return nil
}

func (s *HomeCareProviderService) DeleteHomeCareProvider(id int) error {
	ctx := context.Background()
	err := s.homeCareRepo.Delete(ctx, id)
	if err != nil {
		s.logger.WithError(err).Error("Failed to delete home care provider")
		return fmt.Errorf("failed to delete home care provider: %w", err)
	}
	return nil
}

func (s *HomeCareProviderService) ListHomeCareProviders() ([]models.HomeCareProvider, error) {
	ctx := context.Background()
	providers, err := s.homeCareRepo.GetAll(ctx)
	if err != nil {
		s.logger.WithError(err).Error("Failed to list home care providers")
		return nil, fmt.Errorf("failed to list home care providers: %w", err)
	}
	return providers, nil
}

func (s *HomeCareProviderService) GetProvidersByServiceType(serviceTypeID int, limit, offset int) ([]*models.HomeCareProvider, error) {
	ctx := context.Background()
	providers, err := s.homeCareRepo.GetByServiceType(ctx, serviceTypeID, limit, offset)
	if err != nil {
		s.logger.WithError(err).Error("Failed to get providers by service type")
		return nil, fmt.Errorf("failed to get providers by service type: %w", err)
	}
	return providers, nil
}

// SearchProviders searches for providers based on a query string
func (s *HomeCareProviderService) SearchProviders(query string) ([]*models.HomeCareProvider, error) {
	if query == "" {
		return nil, errors.New("search query cannot be empty")
	}

	ctx := context.Background()
	providers, err := s.homeCareRepo.Search(ctx, query)
	if err != nil {
		s.logger.WithError(err).Error("Failed to search providers")
		return nil, fmt.Errorf("failed to search providers: %w", err)
	}

	return providers, nil
}

func (s *HomeCareProviderService) validateHomeCareProvider(provider models.HomeCareProvider) error {
	if provider.UserID == 0 {
		return errors.New("home care provider user ID is required")
	}
	if provider.ServiceTypeID == 0 {
		return errors.New("home care provider service type ID is required")
	}
	if provider.ExperienceYears < 0 {
		return errors.New("home care provider experience years cannot be negative")
	}
	if provider.HourlyRate < 0 {
		return errors.New("home care provider hourly rate cannot be negative")
	}
	if provider.Qualifications == "" {
		return errors.New("qualifications are required")
	}
	if provider.Bio == "" {
		return errors.New("bio is required")
	}
	return nil
}
