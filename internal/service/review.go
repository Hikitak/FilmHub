package service

import (
    "context"
    "fmt"

    "filmhub/internal/models"
)

// ReviewRepo describes repository dependencies for reviews.
type ReviewRepo interface {
    CreateReview(ctx context.Context, review *models.Review) (int, error)
    ListReviewsByFilm(ctx context.Context, filmID int) ([]models.Review, error)
}

type ReviewService struct {
    repo ReviewRepo
}

func NewReviewService(r ReviewRepo) *ReviewService {
    return &ReviewService{repo: r}
}

func (s *ReviewService) CreateReview(ctx context.Context, review *models.Review) (int, error) {
    id, err := s.repo.CreateReview(ctx, review)
    if err != nil {
        return 0, fmt.Errorf("create review: %w", err)
    }
    return id, nil
}

func (s *ReviewService) ListReviews(ctx context.Context, filmID int) ([]models.Review, error) {
    reviews, err := s.repo.ListReviewsByFilm(ctx, filmID)
    if err != nil {
        return nil, fmt.Errorf("list reviews: %w", err)
    }
    return reviews, nil
} 