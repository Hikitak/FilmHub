package main

import (
	"context"
	"database/sql"
	"filmhub/internal/handler"
	"filmhub/internal/repository"
	"filmhub/internal/service"
	"filmhub/pkg/database"
	jwt "filmhub/pkg/login"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"filmhub/pkg/logger"

	"github.com/gin-gonic/gin"
	pgx "github.com/jackc/pgx/v5"
)

func main() {
	// Initialize logger
	logger.Init()
	defer logger.Sync()

	// Initialize database
	pool, err := database.NewPostgresPool()
	if err != nil {
		logger.Log.Warnf("Failed to connect to database: %v", err)
		logger.Log.Warn("Starting server without database connection...")

	}
	if pool != nil {
		defer pool.Close()
		// apply migrations
		sqlDB, err := sql.Open("pgx", pool.Config().ConnString())
		if err == nil {
			defer sqlDB.Close()
			if err := database.ApplyMigrations(sqlDB, logger.Log); err != nil {
				logger.Log.Warnf("migration error: %v", err)
			}
		} else {
			logger.Log.Warnf("sql open error: %v", err)
		}
	}

	// Get connection for repositories
	var conn *pgx.Conn
	if pool != nil {
		conn, err = pgx.Connect(context.Background(), pool.Config().ConnString())
		if err != nil {
			logger.Log.Warnf("Failed to get database connection: %v", err)
			logger.Log.Warn("Starting server without database connection...")
		}
		if conn != nil {
			defer conn.Close(context.Background())
		}
	}

	// Initialize repositories
	var filmRepo *repository.FilmRepository
	var userRepo repository.UserRepository
	if conn != nil {
		filmRepo = repository.NewFilmRepository(conn)
		userRepo = repository.NewUserRepository(conn)
	} else {
		logger.Log.Warn("Using mock repositories (no database connection)")
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
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Log.Info("Server exited")
}
