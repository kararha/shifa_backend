// File: internal/repository/mysql/user_repo.go

package mysql

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"shifa/internal/repository"
	"shifa/internal/models"
)



type UserRepo struct {
	db *sql.DB
}

// Ensure UserRepo implements the UserRepository interface
var _ repository.UserRepository = (*UserRepo)(nil)

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Create inserts a new user into the database
func (r *UserRepo) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, name, role)
		VALUES (?, ?, ?, ?)
	`
	
	result, err := r.db.ExecContext(ctx, query, user.Email, user.PasswordHash, user.Name, user.Role)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}

// GetByID retrieves a user by their ID
func (r *UserRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, name, role, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetByEmail retrieves a user by their email address
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, password_hash, name, role, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	var user models.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Name,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// Update updates an existing user's information
func (r *UserRepo) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET email = ?, password_hash = ?, name = ?, role = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, user.Email, user.PasswordHash, user.Name, user.Role, time.Now(), user.ID)
	return err
}

// Delete removes a user from the database
func (r *UserRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List retrieves a list of users with optional pagination
func (r *UserRepo) List(ctx context.Context, offset, limit int) ([]*models.User, error) {
	query := `
		SELECT id, email, password_hash, name, role, created_at, updated_at
		FROM users
		ORDER BY id
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.Name,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}