package handler

import (
	"filmhub/internal/models"
	"filmhub/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

type registerRequest struct {
	Username string `json:"username" binding:"required" example:"john_doe" description:"Имя пользователя"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com" description:"Email пользователя"`
	Password string `json:"password" binding:"required,min=6" example:"password123" description:"Пароль (минимум 6 символов)"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john@example.com" description:"Email пользователя"`
	Password string `json:"password" binding:"required" example:"password123" description:"Пароль пользователя"`
}

type loginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJyb2xlIjoidXNlciIsImV4cCI6MTYzNTQ5NjAwMH0.example" description:"JWT токен для авторизации"`
}

// @Summary Регистрация пользователя
// @Description Регистрирует нового пользователя в системе
// @Tags auth
// @Accept json
// @Produce json
// @Param user body registerRequest true "Данные пользователя"
// @Success 201 {object} map[string]interface{} "Пользователь успешно зарегистрирован"
// @Failure 400 {object} map[string]interface{} "Ошибка валидации"
// @Failure 500 {object} map[string]interface{} "Внутренняя ошибка сервера"
// @Router /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	if err := h.service.Register(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

// @Summary Авторизация пользователя
// @Description Авторизует пользователя и возвращает JWT токен
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body loginRequest true "Данные для входа"
// @Success 200 {object} loginResponse "Успешная авторизация"
// @Failure 400 {object} map[string]interface{} "Ошибка валидации"
// @Failure 401 {object} map[string]interface{} "Неверные учетные данные"
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
