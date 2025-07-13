package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/gin-gonic/gin"

	"filmhub/pkg/config"
	"filmhub/pkg/database"
	"filmhub/pkg/logger"
	jwt "filmhub/pkg/login"

	"filmhub/internal/handler"
	"filmhub/internal/repository"
	"filmhub/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	// Initialize logger (dependency-injected)
	log := logger.New(cfg.AppEnv, cfg.SentryDSN)
	defer logger.Sync(log)

	// Initialize database
	pool, err := database.NewPostgresPool()
	if err != nil {
		log.Warnf("Failed to connect to database: %v", err)
		log.Warn("Starting server without database connection...")
	}
	if pool != nil {
		defer pool.Close()
		// apply migrations
		sqlDB, err := sql.Open("pgx", pool.Config().ConnString())
		if err == nil {
			defer sqlDB.Close()
			if err := database.ApplyMigrations(sqlDB, log); err != nil {
				log.Warnf("migration error: %v", err)
			}
		} else {
			log.Warnf("sql open error: %v", err)
		}
	}

	// Ensure we have a working connection pool before proceeding.
	if pool == nil {
		log.Warn("Using mock repositories (no database connection)")
		// TODO: inject mock repositories for offline mode instead of exiting.
		return
	}

	// Initialize repositories with pooled connection
	filmRepo := repository.NewFilmRepository(pool)
	userRepo := repository.NewUserRepository(pool)

	// Initialize services
	filmService := service.NewFilmService(filmRepo)
	authService := service.NewAuthService(userRepo)

	reviewRepo := repository.NewReviewRepository(pool)
	reviewService := service.NewReviewService(reviewRepo)

	// Initialize handlers
	filmHandler := handler.NewFilmHandler(filmService)
	authHandler := handler.NewAuthHandler(authService)
	reviewHandler := handler.NewReviewHandler(reviewService)

	// Setup router (Gin in release mode for prod.)
	if cfg.AppEnv == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(gin.Recovery())

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
		auth.POST("/films/:id/reviews", reviewHandler.CreateReview)
		auth.GET("/films/:id/reviews", reviewHandler.ListReviews)
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
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Info("Server exited")
}
