// package mysql

// import (
//     "context"
//     "database/sql"
//     "errors"
//     "shifa/internal/repository"
//     "shifa/internal/models"
// )

// type mysqlReviewRepo struct {
//     db *sql.DB
// }

// func NewMySQLReviewRepo(db *sql.DB) repository.ReviewRepository {
//     return &mysqlReviewRepo{db: db}
// }

// // Create - matches the interface requirement
// func (r *mysqlReviewRepo) Create(ctx context.Context, review *repository.Review) error {
//     query := `INSERT INTO reviews (patient_id, review_type, consultation_id, home_care_visit_id, 
//               doctor_id, home_care_provider_id, rating, comment, created_at)
//               VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
    
//     _, err := r.db.ExecContext(ctx, query, 
//         review.PatientID, 
//         review.ReviewType, 
//         review.ConsultationID, 
//         review.HomeCareVisitID,
//         review.DoctorID, 
//         review.HomeCareProviderID, 
//         review.Rating, 
//         review.Comment, 
//         review.CreatedAt)
//     return err
// }

// // GetByID - matches the interface requirement
// func (r *mysqlReviewRepo) GetByID(ctx context.Context, id int) (*repository.Review, error) {
//     query := `SELECT id, patient_id, review_type, consultation_id, home_care_visit_id, 
//               doctor_id, home_care_provider_id, rating, comment, created_at
//               FROM reviews WHERE id = ?`
    
//     var review repository.Review
//     err := r.db.QueryRowContext(ctx, query, id).Scan(
//         &review.ID, 
//         &review.PatientID, 
//         &review.ReviewType, 
//         &review.ConsultationID, 
//         &review.HomeCareVisitID,
//         &review.DoctorID, 
//         &review.HomeCareProviderID, 
//         &review.Rating, 
//         &review.Comment, 
//         &review.CreatedAt,
//     )
    
//     if err != nil {
//         if err == sql.ErrNoRows {
//             return nil, errors.New("review not found")
//         }
//         return nil, err
//     }
    
//     return &review, nil
// }

// // GetByDoctorID - matches the interface requirement
// func (r *mysqlReviewRepo) GetByDoctorID(ctx context.Context, doctorID int, limit, offset int) ([]*repository.Review, error) {
//     query := `SELECT id, patient_id, review_type, consultation_id, home_care_visit_id, 
//               doctor_id, home_care_provider_id, rating, comment, created_at
//               FROM reviews WHERE doctor_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
    
//     rows, err := r.db.QueryContext(ctx, query, doctorID, limit, offset)
//     if err != nil {
//         return nil, err
//     }
//     defer rows.Close()
    
//     var reviews []*repository.Review
//     for rows.Next() {
//         var review repository.Review
//         err := rows.Scan(
//             &review.ID, 
//             &review.PatientID, 
//             &review.ReviewType, 
//             &review.ConsultationID, 
//             &review.HomeCareVisitID,
//             &review.DoctorID, 
//             &review.HomeCareProviderID, 
//             &review.Rating, 
//             &review.Comment, 
//             &review.CreatedAt,
//         )
//         if err != nil {
//             return nil, err
//         }
//         reviews = append(reviews, &review)
//     }
    
//     return reviews, nil
// }

// // GetByHomeCareProviderID - matches the interface requirement
// func (r *mysqlReviewRepo) GetByHomeCareProviderID(ctx context.Context, providerID int, limit, offset int) ([]*repository.Review, error) {
//     query := `SELECT id, patient_id, review_type, consultation_id, home_care_visit_id, 
//               doctor_id, home_care_provider_id, rating, comment, created_at
//               FROM reviews WHERE home_care_provider_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
    
//     rows, err := r.db.QueryContext(ctx, query, providerID, limit, offset)
//     if err != nil {
//         return nil, err
//     }
//     defer rows.Close()
    
//     var reviews []*repository.Review
//     for rows.Next() {
//         var review repository.Review
//         err := rows.Scan(
//             &review.ID, 
//             &review.PatientID, 
//             &review.ReviewType, 
//             &review.ConsultationID, 
//             &review.HomeCareVisitID,
//             &review.DoctorID, 
//             &review.HomeCareProviderID, 
//             &review.Rating, 
//             &review.Comment, 
//             &review.CreatedAt,
//         )
//         if err != nil {
//             return nil, err
//         }
//         reviews = append(reviews, &review)
//     }
    
//     return reviews, nil
// }

// // GetReviewByID - new method required by interface
// func (r *mysqlReviewRepo) GetReviewByID(ctx context.Context, id int) (*models.Review, error) {
//     review, err := r.GetByID(ctx, id)
//     if err != nil {
//         return nil, err
//     }
    
//     // Convert repository.Review to models.Review
//     return &models.Review{
//         ID:                 review.ID,
//         PatientID:         review.PatientID,
//         ReviewType:        review.ReviewType,
//         ConsultationID:    review.ConsultationID,
//         HomeCareVisitID:   review.HomeCareVisitID,
//         DoctorID:          review.DoctorID,
//         HomeCareProviderID: review.HomeCareProviderID,
//         Rating:            review.Rating,
//         Comment:           review.Comment,
//         CreatedAt:         review.CreatedAt,
//     }, nil
// }

// // GetReviewsByDoctorID - new method required by interface
// func (r *mysqlReviewRepo) GetReviewsByDoctorID(ctx context.Context, doctorID int) ([]models.Review, error) {
//     reviews, err := r.GetByDoctorID(ctx, doctorID, 100, 0) // Using default limit
//     if err != nil {
//         return nil, err
//     }

//     // Convert []*repository.Review to []models.Review
//     modelReviews := make([]models.Review, len(reviews))
//     for i, review := range reviews {
//         modelReviews[i] = models.Review{
//             ID:                 review.ID,
//             PatientID:         review.PatientID,
//             ReviewType:        review.ReviewType,
//             ConsultationID:    review.ConsultationID,
//             HomeCareVisitID:   review.HomeCareVisitID,
//             DoctorID:          review.DoctorID,
//             HomeCareProviderID: review.HomeCareProviderID,
//             Rating:            review.Rating,
//             Comment:           review.Comment,
//             CreatedAt:         review.CreatedAt,
//         }
//     }
    
//     return modelReviews, nil
// }

// // UpdateReview - matches the interface requirement
// func (r *mysqlReviewRepo) UpdateReview(ctx context.Context, review *models.Review) error {
//     query := `UPDATE reviews SET rating = ?, comment = ? WHERE id = ?`
//     _, err := r.db.ExecContext(ctx, query, review.Rating, review.Comment, review.ID)
//     return err
// }

// // DeleteReview - matches the interface requirement
// func (r *mysqlReviewRepo) DeleteReview(ctx context.Context, id int) error {
//     query := `DELETE FROM reviews WHERE id = ?`
// 	_, err := r.db.ExecContext(ctx, query, id)
//     return err
// }



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
    WHERE home_care_provider_id = ? 
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