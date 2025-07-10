package models

import "time"

type Film struct {
	ID          int       `json:"id"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	ReleaseDate time.Time `json:"release_date"`
	Rating      float32   `json:"rating"`
	CreatedAt   time.Time `json:"created_at"`
}

type FilmRequest struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	ReleaseDate time.Time `json:"release_date"`
}

type Review struct {
	ID        int       `json:"id"`
	FilmID    int       `json:"film_id" validate:"required"`
	UserID    int       `json:"user_id"`
	Rating    int       `json:"rating" validate:"required,min=1,max=10"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}
