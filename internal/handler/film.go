package handler

import (
	"filmhub/internal/models"
	"filmhub/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FilmHandler struct {
	service *service.FilmService
}

func NewFilmHandler(service *service.FilmService) *FilmHandler {
	return &FilmHandler{service: service}
}

func (h *FilmHandler) CreateFilm(c *gin.Context) {
	var req models.FilmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.service.CreateFilm(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create film"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *FilmHandler) GetFilm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid film ID"})
		return
	}

	film, err := h.service.GetFilm(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "film not found"})
		return
	}

	c.JSON(http.StatusOK, film)
}

func (h *FilmHandler) SearchFilms(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search query required"})
		return
	}

	films, err := h.service.SearchFilms(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search failed"})
		return
	}

	c.JSON(http.StatusOK, films)
}
