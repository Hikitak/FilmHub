package repository

import (
    "context"

    "filmhub/internal/models"
    "github.com/jackc/pgx/v5/pgxpool"
)

// ReviewRepository defines storage behaviour for film reviews.
type ReviewRepository struct {
    db *pgxpool.Pool
}

func NewReviewRepository(db *pgxpool.Pool) *ReviewRepository {
    return &ReviewRepository{db: db}
}

func (r *ReviewRepository) CreateReview(ctx context.Context, review *models.Review) (int, error) {
    var id int
    err := r.db.QueryRow(ctx,
        `INSERT INTO reviews (film_id, user_id, rating, comment) VALUES ($1, $2, $3, $4) RETURNING id`,
        review.FilmID, review.UserID, review.Rating, review.Comment,
    ).Scan(&id)
    return id, err
}

func (r *ReviewRepository) ListReviewsByFilm(ctx context.Context, filmID int) ([]models.Review, error) {
    rows, err := r.db.Query(ctx,
        `SELECT id, film_id, user_id, rating, comment, created_at FROM reviews WHERE film_id = $1 ORDER BY created_at DESC`,
        filmID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var reviews []models.Review
    for rows.Next() {
        var rv models.Review
        if err := rows.Scan(&rv.ID, &rv.FilmID, &rv.UserID, &rv.Rating, &rv.Comment, &rv.CreatedAt); err != nil {
            return nil, err
        }
        reviews = append(reviews, rv)
    }
    return reviews, nil
} 