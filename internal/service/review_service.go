package service

import (
    "context"
    "errors"
    "shifa/internal/models"
    "shifa/internal/repository"
    // "shifa/pkg/logger"
    "github.com/sirupsen/logrus" // Import logrus
)

type ReviewService interface {
    CreateReview(ctx context.Context, review *models.Review) error
    GetReviewsByDoctorID(ctx context.Context, doctorID int, page, pageSize int) ([]*models.Review, error)
    GetReviewByID(ctx context.Context, id int) (*models.Review, error)
    GetReviewsByHomeCareProviderID(ctx context.Context, providerID int, page, pageSize int) ([]*models.Review, error)
    UpdateReview(ctx context.Context, review *models.Review) error
    DeleteReview(ctx context.Context, id int) error
    CalculateAverageRating(ctx context.Context, doctorID int) (float64, error)
}

type reviewService struct {
    reviewRepo repository.ReviewRepository
    logger     *logrus.Logger // Change to *logrus.Logger
}

func NewReviewService(repo repository.ReviewRepository, log *logrus.Logger) ReviewService {
    return &reviewService{
        reviewRepo: repo,
        logger:     log,
    }
}

func (s *reviewService) CreateReview(ctx context.Context, review *models.Review) error {
    if review.Rating < 1 || review.Rating > 5 {
        return errors.New("invalid rating")
    }

    repoReview := convertToRepoReview(review)
    err := s.reviewRepo.Create(ctx, repoReview)
    if err != nil {
        s.logger.Error("Failed to create review", "error", err)
        return errors.New("failed to submit review")
    }

    return nil
}

func (s *reviewService) GetReviewByID(ctx context.Context, id int) (*models.Review, error) {
    repoReview, err := s.reviewRepo.GetByID(ctx, id)
    if err != nil {
        s.logger.Error("Failed to get review", "error", err, "reviewID", id)
        return nil, errors.New("review not found")
    }
    return convertToModelReview(repoReview), nil
}

func (s *reviewService) GetReviewsByDoctorID(ctx context.Context, doctorID int, page, pageSize int) ([]*models.Review, error) {
    offset := (page - 1) * pageSize
    repoReviews, err := s.reviewRepo.GetByDoctorID(ctx, doctorID, pageSize, offset)
    if err != nil {
        s.logger.Error("Failed to get reviews by doctor ID", "error", err, "doctorID", doctorID)
        return nil, errors.New("failed to fetch reviews")
    }
    return convertToModelReviews(repoReviews), nil
}

func (s *reviewService) GetReviewsByHomeCareProviderID(ctx context.Context, providerID int, page, pageSize int) ([]*models.Review, error) {
    offset := (page - 1) * pageSize
    repoReviews, err := s.reviewRepo.GetByHomeCareProviderID(ctx, providerID, pageSize, offset)
    if err != nil {
        s.logger.Error("Failed to get reviews by home care provider ID", "error", err, "providerID", providerID)
        return nil, errors.New("failed to fetch reviews")
    }
    return convertToModelReviews(repoReviews), nil
}

func (s *reviewService) UpdateReview(ctx context.Context, review *models.Review) error {
    if review.Rating < 1 || review.Rating > 5 {
        return errors.New("invalid rating")
    }

    err := s.reviewRepo.UpdateReview(ctx, review)
    if err != nil {
        s.logger.Error("Failed to update review", "error", err, "reviewID", review.ID)
        return errors.New("failed to update review")
    }

    return nil
}

func (s *reviewService) DeleteReview(ctx context.Context, id int) error {
    err := s.reviewRepo.DeleteReview(ctx, id)
    if err != nil {
        s.logger.Error("Failed to delete review", "error", err, "reviewID", id)
        return errors.New("failed to delete review")
    }
    return nil
}

func (s *reviewService) CalculateAverageRating(ctx context.Context, doctorID int) (float64, error) {
    // Implementation here
    return 4.5, nil
}

// Conversion functions
func convertToRepoReview(modelReview *models.Review) *repository.Review {
    return &repository.Review{
        ID:                  modelReview.ID,
        PatientID:           modelReview.PatientID,
        ReviewType:          modelReview.ReviewType,
        ConsultationID:      modelReview.ConsultationID,
        HomeCareVisitID:     modelReview.HomeCareVisitID,
        DoctorID:            modelReview.DoctorID,
        HomeCareProviderID:  modelReview.HomeCareProviderID,
        Rating:              modelReview.Rating,
        Comment:             modelReview.Comment,
        CreatedAt:           modelReview.CreatedAt,
    }
}

func convertToModelReview(repoReview *repository.Review) *models.Review {
    return &models.Review{
        ID:                  repoReview.ID,
        PatientID:           repoReview.PatientID,
        ReviewType:          repoReview.ReviewType,
        ConsultationID:      repoReview.ConsultationID,
        HomeCareVisitID:     repoReview.HomeCareVisitID,
        DoctorID:            repoReview.DoctorID,
        HomeCareProviderID:  repoReview.HomeCareProviderID,
        Rating:              repoReview.Rating,
        Comment:             repoReview.Comment,
        CreatedAt:           repoReview.CreatedAt,
    }
}

func convertToModelReviews(repoReviews []*repository.Review) []*models.Review {
    modelReviews := make([]*models.Review, len(repoReviews))
    for i, repoReview := range repoReviews {
        modelReviews[i] = convertToModelReview(repoReview)
    }
    return modelReviews
}