package service

import (
	"context"
	"filmhub/internal/models"
	"filmhub/internal/repository"
)

type FilmService struct {
	repo *repository.FilmRepository
}

func NewFilmService(repo *repository.FilmRepository) *FilmService {
	return &FilmService{repo: repo}
}

func (s *FilmService) CreateFilm(ctx context.Context, film *models.FilmRequest) (int, error) {
	return s.repo.CreateFilm(ctx, film)
}

func (s *FilmService) GetFilm(ctx context.Context, id int) (*models.Film, error) {
	return s.repo.GetFilmByID(ctx, id)
}

func (s *FilmService) SearchFilms(ctx context.Context, query string) ([]models.Film, error) {
	return s.repo.SearchFilms(ctx, query)
}
