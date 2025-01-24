// // File: internal/repository/mysql/service_type_repo.go

// package mysql

// import (
// 	"context"
// 	"database/sql"
// 	"errors"

// 	"shifa/internal/models"
// )

// // ServiceTypeRepo represents the MySQL repository for service type-related database operations
// type ServiceTypeRepo struct {
// 	db *sql.DB
// }

// // NewServiceTypeRepo creates a new ServiceTypeRepo instance
// func NewServiceTypeRepo(db *sql.DB) *ServiceTypeRepo {
// 	return &ServiceTypeRepo{db: db}
// }

// // Create inserts a new service type into the database
// func (r *ServiceTypeRepo) Create(ctx context.Context, serviceType *models.ServiceType) error {
// 	query := `
// 		INSERT INTO service_types (name, description, is_home_care)
// 		VALUES (?, ?, ?)
// 	`
	
// 	result, err := r.db.ExecContext(ctx, query, 
// 		serviceType.Name, serviceType.Description, serviceType.IsHomeCare)
// 	if err != nil {
// 		return err
// 	}

// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		return err
// 	}

// 	serviceType.ID = int(id)
// 	return nil
// }

// // GetByID retrieves a service type by its ID
// func (r *ServiceTypeRepo) GetByID(ctx context.Context, id int) (*models.ServiceType, error) {
// 	query := `
// 		SELECT id, name, description, is_home_care
// 		FROM service_types
// 		WHERE id = ?
// 	`

// 	var serviceType models.ServiceType
// 	err := r.db.QueryRowContext(ctx, query, id).Scan(
// 		&serviceType.ID, &serviceType.Name, &serviceType.Description, &serviceType.IsHomeCare,
// 	)

// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, errors.New("service type not found")
// 		}
// 		return nil, err
// 	}

// 	return &serviceType, nil
// }

// // Update updates an existing service type's information
// func (r *ServiceTypeRepo) Update(ctx context.Context, serviceType *models.ServiceType) error {
// 	query := `
// 		UPDATE service_types
// 		SET name = ?, description = ?, is_home_care = ?
// 		WHERE id = ?
// 	`

// 	_, err := r.db.ExecContext(ctx, query,
// 		serviceType.Name, serviceType.Description, serviceType.IsHomeCare, serviceType.ID)

// 	return err
// }

// // Delete removes a service type from the database
// func (r *ServiceTypeRepo) Delete(ctx context.Context, id int) error {
// 	query := `DELETE FROM service_types WHERE id = ?`
// 	_, err := r.db.ExecContext(ctx, query, id)
// 	return err
// }

// // List retrieves a list of service types with optional pagination
// func (r *ServiceTypeRepo) List(ctx context.Context, offset, limit int) ([]*models.ServiceType, error) {
// 	query := `
// 		SELECT id, name, description, is_home_care
// 		FROM service_types
// 		ORDER BY id
// 		LIMIT ? OFFSET ?
// 	`

// 	rows, err := r.db.QueryContext(ctx, query, limit, offset)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var serviceTypes []*models.ServiceType
// 	for rows.Next() {
// 		var serviceType models.ServiceType
// 		err := rows.Scan(
// 			&serviceType.ID, &serviceType.Name, &serviceType.Description, &serviceType.IsHomeCare,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		serviceTypes = append(serviceTypes, &serviceType)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return serviceTypes, nil
// }



// File: internal/repository/mysql/service_type_repository.go
package mysql

import (
    "context"
    "database/sql"
    "errors"
    "shifa/internal/repository" // Import the repository interface
)

// Ensure ServiceTypeRepository implements repository.ServiceTypeRepository
var _ repository.ServiceTypeRepository = &ServiceTypeRepository{}

type ServiceTypeRepository struct {
    db *sql.DB
}

// NewServiceTypeRepository creates a new ServiceTypeRepository instance
func NewServiceTypeRepository(db *sql.DB) *ServiceTypeRepository {
    return &ServiceTypeRepository{db: db}
}

// Create inserts a new service type into the database
func (r *ServiceTypeRepository) Create(ctx context.Context, serviceType *repository.ServiceType) error {
    query := `
        INSERT INTO service_types (name, description, is_home_care)
        VALUES (?, ?, ?)
    `
    
    result, err := r.db.ExecContext(ctx, query, 
        serviceType.Name, serviceType.Description, serviceType.IsHomeCare)
    if err != nil {
        return err
    }

    id, err := result.LastInsertId()
    if err != nil {
        return err
    }

    serviceType.ID = int(id)
    return nil
}

// GetByID retrieves a service type by its ID
func (r *ServiceTypeRepository) GetByID(ctx context.Context, id int) (*repository.ServiceType, error) {
    query := `
        SELECT id, name, description, is_home_care
        FROM service_types
        WHERE id = ?
    `

    var serviceType repository.ServiceType
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &serviceType.ID, &serviceType.Name, &serviceType.Description, &serviceType.IsHomeCare,
    )

    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, errors.New("service type not found")
        }
        return nil, err
    }

    return &serviceType, nil
}

// GetAll retrieves all service types
func (r *ServiceTypeRepository) GetAll(ctx context.Context) ([]*repository.ServiceType, error) {
    query := `
        SELECT id, name, description, is_home_care
        FROM service_types
        ORDER BY id
    `

    rows, err := r.db.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var serviceTypes []*repository.ServiceType
    for rows.Next() {
        var serviceType repository.ServiceType
        err := rows.Scan(
            &serviceType.ID, &serviceType.Name, &serviceType.Description, &serviceType.IsHomeCare,
        )
        if err != nil {
            return nil, err
        }
        serviceTypes = append(serviceTypes, &serviceType)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return serviceTypes, nil
}

// Update updates an existing service type's information
func (r *ServiceTypeRepository) Update(ctx context.Context, serviceType *repository.ServiceType) error {
    query := `
        UPDATE service_types
        SET name = ?, description = ?, is_home_care = ?
        WHERE id = ?
    `

    result, err := r.db.ExecContext(ctx, query,
        serviceType.Name, serviceType.Description, serviceType.IsHomeCare, serviceType.ID)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return errors.New("no service type found with the given ID")
    }

    return nil
}

// Delete removes a service type from the database
func (r *ServiceTypeRepository) Delete(ctx context.Context, id int) error {
    query := `DELETE FROM service_types WHERE id = ?`
    
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return err
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return errors.New("no service type found with the given ID")
    }

    return nil
}