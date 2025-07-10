package main

import (
	"context"
	"filmhub/internal/handler"
	"filmhub/internal/repository"
	"filmhub/internal/service"
	"filmhub/pkg/database"
	"filmhub/pkg/jwt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	_ = godotenv.Load()

	// Initialize database
	pool, err := database.NewPostgresPool()
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
		log.Println("Starting server without database connection...")
		// Можно запустить сервер без БД для тестирования
		// return
	}
	if pool != nil {
		defer pool.Close()
	}

	// Get connection for repositories
	var conn *pgx.Conn
	if pool != nil {
		conn, err = pgx.Connect(context.Background(), pool.Config().ConnString())
		if err != nil {
			log.Printf("Warning: Failed to get database connection: %v", err)
			log.Println("Starting server without database connection...")
		}
		if conn != nil {
			defer conn.Close(context.Background())

			// Run migrations if database is available
			log.Println("Running database migrations...")
			if err := runMigrations(conn); err != nil {
				log.Printf("Warning: Failed to run migrations: %v", err)
			} else {
				log.Println("Migrations completed successfully")
			}
		}
	}

	// Initialize repositories
	var filmRepo *repository.FilmRepository
	var userRepo repository.UserRepository
	if conn != nil {
		filmRepo = repository.NewFilmRepository(conn)
		userRepo = repository.NewUserRepository(conn)
	} else {
		log.Println("Using mock repositories (no database connection)")
		// Здесь можно добавить mock репозитории
		return
	}

	// Initialize services
	filmService := service.NewFilmService(filmRepo)
	authService := service.NewAuthService(userRepo)

	// Initialize handlers
	filmHandler := handler.NewFilmHandler(filmService)
	authHandler := handler.NewAuthHandler(authService)

	// Setup router
	router := gin.Default()

	// Public routes
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)
	router.GET("/films", filmHandler.SearchFilms)
	router.GET("/films/:id", filmHandler.GetFilm)

	// Protected routes (require JWT)
	auth := router.Group("/")
	auth.Use(jwt.AuthMiddleware())
	{
		auth.POST("/films", filmHandler.CreateFilm)
	}

	// Start server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func runMigrations(conn *pgx.Conn) error {
	ctx := context.Background()

	// Создание таблицы пользователей
	_, err := conn.Exec(ctx, `
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
		return err
	}

	// Создание таблицы фильмов
	_, err = conn.Exec(ctx, `
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
		return err
	}

	// Создание таблицы отзывов
	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS reviews (
			id SERIAL PRIMARY KEY,
			film_id INTEGER REFERENCES films(id) ON DELETE CASCADE,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			rating INTEGER CHECK (rating >= 1 AND rating <= 10),
			comment TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}
