// internal/repository/mysql/review_repository.go
package mysql

import (
	"context"
	"database/sql"
	"errors"
	"shifa/internal/models"
	"shifa/internal/repository"
)

type mysqlReviewRepo struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) repository.ReviewRepository {
	return &mysqlReviewRepo{db: db}
}

func (r *mysqlReviewRepo) Create(ctx context.Context, review *repository.Review) error {
	query := `INSERT INTO reviews (
        patient_id, review_type, consultation_id, home_care_visit_id, 
        doctor_id, home_care_provider_id, rating, comment, created_at
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW())`

	result, err := r.db.ExecContext(ctx, query,
		review.PatientID,
		review.ReviewType,
		review.ConsultationID,
		review.HomeCareVisitID,
		review.DoctorID,
		review.HomeCareProviderID,
		review.Rating,
		review.Comment,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	review.ID = int(id)
	return nil
}

func (r *mysqlReviewRepo) GetByID(ctx context.Context, id int) (*repository.Review, error) {
	query := `SELECT 
        id, patient_id, review_type, consultation_id, home_care_visit_id, 
        doctor_id, home_care_provider_id, rating, comment, created_at
    FROM reviews WHERE id = ?`

	var review repository.Review
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&review.ID,
		&review.PatientID,
		&review.ReviewType,
		&review.ConsultationID,
		&review.HomeCareVisitID,
		&review.DoctorID,
		&review.HomeCareProviderID,
		&review.Rating,
		&review.Comment,
		&review.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("review not found")
		}
		return nil, err
	}

	return &review, nil
}

func (r *mysqlReviewRepo) GetReviewByID(ctx context.Context, id int) (*models.Review, error) {
	query := `SELECT 
        id, patient_id, review_type, consultation_id, home_care_visit_id, 
        doctor_id, home_care_provider_id, rating, comment, created_at
    FROM reviews WHERE id = ?`

	var review models.Review
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&review.ID,
		&review.PatientID,
		&review.ReviewType,
		&review.ConsultationID,
		&review.HomeCareVisitID,
		&review.DoctorID,
		&review.HomeCareProviderID,
		&review.Rating,
		&review.Comment,
		&review.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("review not found")
		}
		return nil, err
	}

	return &review, nil
}

func (r *mysqlReviewRepo) GetByDoctorID(ctx context.Context, doctorID int, limit, offset int) ([]*repository.Review, error) {
	query := `SELECT 
        id, patient_id, review_type, consultation_id, home_care_visit_id, 
        doctor_id, home_care_provider_id, rating, comment, created_at
    FROM reviews 
    WHERE doctor_id = ? 
    ORDER BY created_at DESC 
    LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, doctorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*repository.Review
	for rows.Next() {
		var review repository.Review
		err := rows.Scan(
			&review.ID,
			&review.PatientID,
			&review.ReviewType,
			&review.ConsultationID,
			&review.HomeCareVisitID,
			&review.DoctorID,
			&review.HomeCareProviderID,
			&review.Rating,
			&review.Comment,
			&review.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, &review)
	}

	return reviews, nil
}

func (r *mysqlReviewRepo) GetReviewsByDoctorID(ctx context.Context, doctorID int) ([]models.Review, error) {
	query := `SELECT 
        id, patient_id, review_type, consultation_id, home_care_visit_id, 
        doctor_id, home_care_provider_id, rating, comment, created_at
    FROM reviews 
    WHERE doctor_id = ? 
    ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, doctorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		err := rows.Scan(
			&review.ID,
			&review.PatientID,
			&review.ReviewType,
			&review.ConsultationID,
			&review.HomeCareVisitID,
			&review.DoctorID,
			&review.HomeCareProviderID,
			&review.Rating,
			&review.Comment,
			&review.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}

func (r *mysqlReviewRepo) GetByHomeCareProviderID(ctx context.Context, providerID int, limit, offset int) ([]*repository.Review, error) {
	query := `SELECT 
        id, patient_id, review_type, consultation_id, home_care_visit_id, 
        doctor_id, home_care_provider_id, rating, comment, created_at
    FROM reviews 
    WHERE home_care_provider_id = ? AND review_type = 'home_care'
    ORDER BY created_at DESC 
    LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, providerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*repository.Review
	for rows.Next() {
		var review repository.Review
		err := rows.Scan(
			&review.ID,
			&review.PatientID,
			&review.ReviewType,
			&review.ConsultationID,
			&review.HomeCareVisitID,
			&review.DoctorID,
			&review.HomeCareProviderID,
			&review.Rating,
			&review.Comment,
			&review.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, &review)
	}

	return reviews, nil
}

func (r *mysqlReviewRepo) UpdateReview(ctx context.Context, review *models.Review) error {
	query := `UPDATE reviews 
    SET rating = ?, comment = ? 
    WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, review.Rating, review.Comment, review.ID)
	return err
}

func (r *mysqlReviewRepo) DeleteReview(ctx context.Context, id int) error {
	query := `DELETE FROM reviews WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
