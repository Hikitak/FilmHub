package handler

import (
    "bytes"
    "context"
    "encoding/json"
    "filmhub/internal/models"
    "filmhub/internal/service"
    jwtpkg "filmhub/pkg/login"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
)

type stubFilmRepo struct{}

func (stubFilmRepo) CreateFilm(_ context.Context, _ *models.FilmRequest) (int, error) { return 1, nil }
func (stubFilmRepo) GetFilmByID(_ context.Context, id int) (*models.Film, error) { return &models.Film{ID: id, Title: "Test", Description: "",}, nil }
func (stubFilmRepo) SearchFilms(_ context.Context, _ string) ([]models.Film, error) { return []models.Film{}, nil }

func TestAuthAndCreateFilmRoute(t *testing.T) {
    gin.SetMode(gin.TestMode)
    jwtpkg.Init("testsecret")
    // Prepare token
    token, _ := jwtpkg.GenerateToken(1, "admin")

    // stub film repository and service
    filmSvc := service.NewFilmService(stubFilmRepo{})
    filmHandler := NewFilmHandler(filmSvc)

    r := gin.Default()
    // protected group
    authGroup := r.Group("/")
    authGroup.Use(jwtpkg.AuthMiddleware())
    {
        authGroup.POST("/films", filmHandler.CreateFilm)
    }

    // Create film request through protected route
    filmBody, _ := json.Marshal(models.FilmRequest{Title: "Test", Description: "",})
    req := httptest.NewRequest(http.MethodPost, "/films", bytes.NewReader(filmBody))
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")
    resp := httptest.NewRecorder()
    r.ServeHTTP(resp, req)
    if resp.Code != http.StatusCreated {
        t.Fatalf("expected 201 created, got %d", resp.Code)
    }
} 