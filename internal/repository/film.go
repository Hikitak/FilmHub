package repository

import (
	"context"
	"filmhub/internal/models"

	pgx "github.com/jackc/pgx/v5"
)

type FilmRepository struct {
	db *pgx.Conn
}

func NewFilmRepository(db *pgx.Conn) *FilmRepository {
	return &FilmRepository{db: db}
}

func (r *FilmRepository) CreateFilm(ctx context.Context, film *models.FilmRequest) (int, error) {
	var id int
	err := r.db.QueryRow(ctx,
		`INSERT INTO films (title, description, release_date) 
         VALUES ($1, $2, $3) RETURNING id`,
		film.Title, film.Description, film.ReleaseDate).Scan(&id)
	return id, err
}

func (r *FilmRepository) GetFilmByID(ctx context.Context, id int) (*models.Film, error) {
	var film models.Film
	err := r.db.QueryRow(ctx,
		`SELECT id, title, description, release_date, rating, created_at 
         FROM films WHERE id = $1`, id).Scan(
		&film.ID, &film.Title, &film.Description, &film.ReleaseDate, &film.Rating, &film.CreatedAt,
	)
	return &film, err
}

func (r *FilmRepository) SearchFilms(ctx context.Context, query string) ([]models.Film, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, title, description, release_date, rating, created_at 
         FROM films WHERE title ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%'`,
		query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var films []models.Film
	for rows.Next() {
		var film models.Film
		if err := rows.Scan(
			&film.ID, &film.Title, &film.Description, &film.ReleaseDate, &film.Rating, &film.CreatedAt,
		); err != nil {
			return nil, err
		}
		films = append(films, film)
	}

	return films, nil
}
