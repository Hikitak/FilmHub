package models

import "time"

type Film struct {
	ID          int       `json:"id" example:"1" description:"Уникальный идентификатор фильма"`
	Title       string    `json:"title" validate:"required" example:"The Matrix" description:"Название фильма"`
	Description string    `json:"description" validate:"required" example:"Sci-fi action movie about virtual reality" description:"Описание фильма"`
	ReleaseDate time.Time `json:"release_date" example:"1999-03-31T00:00:00Z" description:"Дата выхода фильма"`
	Rating      float32   `json:"rating" example:"8.7" description:"Рейтинг фильма"`
	CreatedAt   time.Time `json:"created_at" example:"2023-01-01T00:00:00Z" description:"Дата создания записи"`
}

type FilmRequest struct {
	Title       string    `json:"title" validate:"required" example:"The Matrix" description:"Название фильма"`
	Description string    `json:"description" validate:"required" example:"Sci-fi action movie about virtual reality" description:"Описание фильма"`
	ReleaseDate time.Time `json:"release_date" example:"1999-03-31T00:00:00Z" description:"Дата выхода фильма"`
}

type Review struct {
	ID        int       `json:"id" example:"1" description:"Уникальный идентификатор отзыва"`
	FilmID    int       `json:"film_id" validate:"required" example:"1" description:"ID фильма"`
	UserID    int       `json:"user_id" example:"1" description:"ID пользователя"`
	Rating    int       `json:"rating" validate:"required,min=1,max=10" example:"8" description:"Оценка от 1 до 10"`
	Comment   string    `json:"comment" example:"Отличный фильм!" description:"Комментарий к отзыву"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z" description:"Дата создания отзыва"`
}
