package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikasa/skk/internal/domain"
	"github.com/mikasa/skk/internal/middleware"
	"github.com/mikasa/skk/internal/repository"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  *domain.User `json:"user"`
}

type AuthService interface {
	Register(ctx context.Context, username, email, password, name string) (*domain.User, error)
	Login(ctx context.Context, username, password string) (string, *domain.User, error)
	GetCurrentUser(ctx context.Context, userID int64) (*domain.User, error)
}

type AuthHandler struct {
	svc AuthService
}

func NewAuthHandler(svc AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.GET("/me", h.Me)
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.svc.Register(c.Request.Context(), req.Username, req.Email, req.Password, req.Name)
	if err != nil {
		mapAuthError(c, err)
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.svc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		mapAuthError(c, err)
		return
	}
	c.JSON(http.StatusOK, LoginResponse{Token: token, User: user})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.svc.GetCurrentUser(c.Request.Context(), userID)
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func mapAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, repository.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
