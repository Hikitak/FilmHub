package main

import (
	"context"
	"filmhub/pkg/database"
	"log"

	pgx "github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	pool, err := database.NewPostgresPool()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	conn, err := pgx.Connect(context.Background(), pool.Config().ConnString())
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}
	defer conn.Close(context.Background())

	// Создание таблицы пользователей
	_, err = conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(255) NOT NULL UNIQUE,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'user',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}
	log.Println("Users table created successfully")

	// Создание таблицы фильмов
	_, err = conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS films (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			release_date DATE,
			rating DECIMAL(3,2) DEFAULT 0.0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create films table: %v", err)
	}
	log.Println("Films table created successfully")

	// Создание таблицы отзывов
	_, err = conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS reviews (
			id SERIAL PRIMARY KEY,
			film_id INTEGER REFERENCES films(id) ON DELETE CASCADE,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			rating INTEGER CHECK (rating >= 1 AND rating <= 10),
			comment TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create reviews table: %v", err)
	}
	log.Println("Reviews table created successfully")

	log.Println("All migrations completed successfully")
}
