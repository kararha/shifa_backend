// package service

// import (
//     "context"  // Add context import
//     "errors"
//     "fmt"
//     "shifa/internal/models"
//     "shifa/internal/repository"

//     "github.com/sirupsen/logrus"
// )

// type ServiceTypeService struct {
//     serviceTypeRepo repository.ServiceTypeRepository
//     logger         *logrus.Logger
// }

// func NewServiceTypeService(serviceTypeRepo repository.ServiceTypeRepository, logger *logrus.Logger) *ServiceTypeService {
//     return &ServiceTypeService{
//         serviceTypeRepo: serviceTypeRepo,
//         logger:         logger,
//     }
// }

// // CreateServiceType now includes context and handles proper type conversion
// func (s *ServiceTypeService) CreateServiceType(ctx context.Context, serviceType models.ServiceType) error {
//     // Validate the service type model
//     if err := s.validateServiceType(serviceType); err != nil {
//         s.logger.WithError(err).Error("Invalid service type data")
//         return err
//     }

//     // Convert models.ServiceType to repository.ServiceType
//     repoServiceType := &repository.ServiceType{
//         // Map the fields accordingly
//         Name:        serviceType.Name,
//         Description: serviceType.Description,
//         // Add other fields as needed
//     }

//     // Create the service type with context
//     if err := s.serviceTypeRepo.Create(ctx, repoServiceType); err != nil {
//         s.logger.WithError(err).Error("Failed to create service type")
//         return fmt.Errorf("failed to create service type: %w", err)
//     }

//     s.logger.Infof("Service type created successfully: %v", serviceType)
//     return nil
// }

// // ListServiceTypes now includes context
// func (s *ServiceTypeService) ListServiceTypes(ctx context.Context) ([]models.ServiceType, error) {
//     serviceTypes, err := s.serviceTypeRepo.GetAll(ctx)
//     if err != nil {
//         s.logger.WithError(err).Error("Failed to list service types")
//         return nil, fmt.Errorf("failed to list service types: %w", err)
//     }

//     // Convert repository.ServiceType slice to models.ServiceType slice
//     result := make([]models.ServiceType, len(serviceTypes))
//     for i, st := range serviceTypes {
//         result[i] = models.ServiceType{
//             // Map the fields accordingly
//             Name:        st.Name,
//             Description: st.Description,
//             // Add other fields as needed
//         }
//     }

//     return result, nil
// }

// func (s *ServiceTypeService) validateServiceType(serviceType models.ServiceType) error {
//     if serviceType.Name == "" {
//         return errors.New("service type name is required")
//     }

//     if serviceType.Description == "" {
//         return errors.New("service type description is required")
//     }

//     // Add any other validation logic here

//     return nil
// }

// // GetServiceTypeByID retrieves a service type by ID
// func (s *ServiceTypeService) GetServiceTypeByID(ctx context.Context, id int) (*models.ServiceType, error) {
//     serviceType, err := s.serviceTypeRepo.GetByID(ctx, id)
//     if err != nil {
//         s.logger.WithError(err).Error("Failed to get service type by ID")
//         return nil, fmt.Errorf("failed to get service type: %w", err)
//     }

//     // Convert repository.ServiceType to models.ServiceType
//     return &models.ServiceType{
//         Name:        serviceType.Name,
//         Description: serviceType.Description,
//         // Map other fields
//     }, nil
// }

// // UpdateServiceType updates an existing service type
// func (s *ServiceTypeService) UpdateServiceType(ctx context.Context, serviceType models.ServiceType) error {
//     // Validate the service type model
//     if err := s.validateServiceType(serviceType); err != nil {
//         s.logger.WithError(err).Error("Invalid service type data for update")
//         return err
//     }

//     // Convert models.ServiceType to repository.ServiceType
//     repoServiceType := &repository.ServiceType{
//         Name:        serviceType.Name,
//         Description: serviceType.Description,
//         // Map other fields
//     }

//     // Update the service type
//     if err := s.serviceTypeRepo.Update(ctx, repoServiceType); err != nil {
//         s.logger.WithError(err).Error("Failed to update service type")
//         return fmt.Errorf("failed to update service type: %w", err)
//     }

//     s.logger.Infof("Service type updated successfully: %v", serviceType)
//     return nil
// }

// // DeleteServiceType deletes a service type by ID
// func (s *ServiceTypeService) DeleteServiceType(ctx context.Context, id int) error {
//     if err := s.serviceTypeRepo.Delete(ctx, id); err != nil {
//         s.logger.WithError(err).Error("Failed to delete service type")
//         return fmt.Errorf("failed to delete service type: %w", err)
//     }

//     s.logger.Infof("Service type with ID %d deleted successfully", id)
//     return nil
// }


package service

import (
    "context"
    "errors"
    "fmt"
    "shifa/internal/repository"

    "github.com/sirupsen/logrus"
)

type ServiceTypeService struct {
    serviceTypeRepo repository.ServiceTypeRepository
    logger         *logrus.Logger
}

func NewServiceTypeService(serviceTypeRepo repository.ServiceTypeRepository, logger *logrus.Logger) *ServiceTypeService {
    return &ServiceTypeService{
        serviceTypeRepo: serviceTypeRepo,
        logger:         logger,
    }
}

// CreateServiceType handles the creation of a new service type
func (s *ServiceTypeService) CreateServiceType(ctx context.Context, serviceType *repository.ServiceType) error {
    // Validate the service type
    if err := s.validateServiceType(serviceType); err != nil {
        s.logger.WithError(err).Error("Invalid service type data")
        return err
    }

    // Create the service type
    if err := s.serviceTypeRepo.Create(ctx, serviceType); err != nil {
        s.logger.WithError(err).Error("Failed to create service type")
        return fmt.Errorf("failed to create service type: %w", err)
    }

    s.logger.Infof("Service type created successfully: %v", serviceType)
    return nil
}

// ListServiceTypes retrieves all service types
func (s *ServiceTypeService) ListServiceTypes(ctx context.Context) ([]*repository.ServiceType, error) {
    serviceTypes, err := s.serviceTypeRepo.GetAll(ctx)
    if err != nil {
        s.logger.WithError(err).Error("Failed to list service types")
        return nil, fmt.Errorf("failed to list service types: %w", err)
    }

    return serviceTypes, nil
}

// validateServiceType performs validation on the service type
func (s *ServiceTypeService) validateServiceType(serviceType *repository.ServiceType) error {
    if serviceType == nil {
        return errors.New("service type cannot be nil")
    }

    if serviceType.Name == "" {
        return errors.New("service type name is required")
    }

    if len(serviceType.Name) > 255 {
        return errors.New("service type name is too long")
    }

    if len(serviceType.Description) > 1000 {
        return errors.New("service type description is too long")
    }

    return nil
}

// GetServiceTypeByID retrieves a service type by ID
func (s *ServiceTypeService) GetServiceTypeByID(ctx context.Context, id int) (*repository.ServiceType, error) {
    serviceType, err := s.serviceTypeRepo.GetByID(ctx, id)
    if err != nil {
        s.logger.WithError(err).WithField("id", id).Error("Failed to get service type by ID")
        return nil, fmt.Errorf("failed to get service type: %w", err)
    }

    return serviceType, nil
}

// UpdateServiceType updates an existing service type
func (s *ServiceTypeService) UpdateServiceType(ctx context.Context, serviceType *repository.ServiceType) error {
    // Validate the service type model
    if err := s.validateServiceType(serviceType); err != nil {
        s.logger.WithError(err).Error("Invalid service type data for update")
        return err
    }

    // Update the service type
    if err := s.serviceTypeRepo.Update(ctx, serviceType); err != nil {
        s.logger.WithError(err).Error("Failed to update service type")
        return fmt.Errorf("failed to update service type: %w", err)
    }

    s.logger.Infof("Service type updated successfully: %v", serviceType)
    return nil
}

// DeleteServiceType deletes a service type by ID
func (s *ServiceTypeService) DeleteServiceType(ctx context.Context, id int) error {
    if id <= 0 {
        return errors.New("invalid service type ID")
    }

    if err := s.serviceTypeRepo.Delete(ctx, id); err != nil {
        s.logger.WithError(err).WithField("id", id).Error("Failed to delete service type")
        return fmt.Errorf("failed to delete service type: %w", err)
    }

    s.logger.Infof("Service type with ID %d deleted successfully", id)
    return nil
}