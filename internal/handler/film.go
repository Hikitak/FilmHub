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

// @Summary Создание фильма
// @Description Создает новый фильм (требует авторизации)
// @Tags films
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param film body models.FilmRequest true "Данные фильма"
// @Success 201 {object} models.Film "Фильм успешно создан"
// @Failure 400 {object} map[string]interface{} "Ошибка валидации"
// @Failure 401 {object} map[string]interface{} "Не авторизован"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /films [post]
func (h *FilmHandler) CreateFilm(c *gin.Context) {
	var req models.FilmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.service.CreateFilm(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// @Summary Получение фильма по ID
// @Description Возвращает информацию о фильме по его ID
// @Tags films
// @Accept json
// @Produce json
// @Param id path int true "ID фильма"
// @Success 200 {object} models.Film "Информация о фильме"
// @Failure 400 {object} map[string]interface{} "Неверный ID"
// @Failure 404 {object} map[string]interface{} "Фильм не найден"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /films/{id} [get]
func (h *FilmHandler) GetFilm(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid film ID"})
		return
	}

	film, err := h.service.GetFilm(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Film not found"})
		return
	}

	c.JSON(http.StatusOK, film)
}

// @Summary Поиск фильмов
// @Description Ищет фильмы по названию или описанию
// @Tags films
// @Accept json
// @Produce json
// @Param query query string false "Поисковый запрос"
// @Success 200 {array} models.Film "Список найденных фильмов"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /films [get]
func (h *FilmHandler) SearchFilms(c *gin.Context) {
	query := c.Query("query")

	films, err := h.service.SearchFilms(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, films)
}
