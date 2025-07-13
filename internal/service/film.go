package service

import (
	"context"
	"errors"
	"fmt"
	"filmhub/internal/models"

	"github.com/jackc/pgx/v5"
)

// ErrFilmNotFound returned when the film can't be located in storage.
var ErrFilmNotFound = errors.New("film not found")

// FilmRepo describes storage operations required by FilmService. This allows
// us to inject mocks in tests and keeps the service agnostic of the concrete
// repository implementation.
type FilmRepo interface {
	CreateFilm(ctx context.Context, film *models.FilmRequest) (int, error)
	GetFilmByID(ctx context.Context, id int) (*models.Film, error)
	SearchFilms(ctx context.Context, query string) ([]models.Film, error)
}

type FilmService struct {
	repo FilmRepo
}

func NewFilmService(repo FilmRepo) *FilmService {
	return &FilmService{repo: repo}
}

func (s *FilmService) CreateFilm(ctx context.Context, film *models.FilmRequest) (int, error) {
	return s.repo.CreateFilm(ctx, film)
}

func (s *FilmService) GetFilm(ctx context.Context, id int) (*models.Film, error) {
	film, err := s.repo.GetFilmByID(ctx, id)
	if err != nil {
		// Map storage no-row error to domain-level error.
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrFilmNotFound
		}
		return nil, fmt.Errorf("get film: %w", err)
	}
	return film, nil
}

func (s *FilmService) SearchFilms(ctx context.Context, query string) ([]models.Film, error) {
	films, err := s.repo.SearchFilms(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("search films: %w", err)
	}
	return films, nil
}
