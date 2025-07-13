package handler

import (
    "filmhub/internal/models"
    "filmhub/internal/service"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

type ReviewHandler struct {
    service *service.ReviewService
}

func NewReviewHandler(s *service.ReviewService) *ReviewHandler {
    return &ReviewHandler{service: s}
}

// CreateReview godoc
// @Summary Добавить отзыв
// @Tags reviews
// @Accept json
// @Produce json
// @Param id path int true "ID фильма"
// @Param review body models.Review true "Отзыв"
// @Success 201 {object} map[string]int {"id":1}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /films/{id}/reviews [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
    userIDVal, exists := c.Get("user_id")
    if !exists {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
    userID, ok := userIDVal.(int)
    if !ok {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }
    filmID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid film id"})
        return
    }
    var req models.Review
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    req.FilmID = filmID
    req.UserID = userID
    id, err := h.service.CreateReview(c.Request.Context(), &req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"id": id})
}

// ListReviews godoc
// @Summary Список отзывов фильма
// @Tags reviews
// @Produce json
// @Param id path int true "ID фильма"
// @Success 200 {array} models.Review
// @Router /films/{id}/reviews [get]
func (h *ReviewHandler) ListReviews(c *gin.Context) {
    filmID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid film id"})
        return
    }
    reviews, err := h.service.ListReviews(c.Request.Context(), filmID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, reviews)
} 