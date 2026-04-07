package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mikasa/skk/internal/domain"
	"github.com/mikasa/skk/internal/middleware"
	"github.com/mikasa/skk/internal/repository"
)

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"required,min=0,max=150"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"required,min=0,max=150"`
}

type UpdateProfileRequest struct {
	Name    string `json:"name"`
	Bio     string `json:"bio"`
	Phone   string `json:"phone"`
	City    string `json:"city"`
	Website string `json:"website"`
}

type UserService interface {
	Create(ctx context.Context, name, email string, age int) (*domain.User, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	Update(ctx context.Context, id int64, name, email string, age int) (*domain.User, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, page, size int) ([]*domain.User, error)
	GetProfile(ctx context.Context, id int64) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID int64, name, bio, phone, city, website string) (*domain.User, error)
	UpdateAvatar(ctx context.Context, userID int64, avatarPath string) (*domain.User, error)
}

var allowedAvatarExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

const maxAvatarSize = 2 << 20 // 2 MB

type UserHandler struct {
	svc UserService
}

func NewUserHandler(svc UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) RegisterPublicRoutes(rg *gin.RouterGroup) {
	rg.POST("", h.Create)
	rg.GET("/:id", h.GetByID)
	rg.PUT("/:id", h.Update)
	rg.DELETE("/:id", h.Delete)
	rg.GET("", h.List)
	rg.GET("/:id/profile", h.GetProfile)
}

func (h *UserHandler) RegisterProtectedRoutes(rg *gin.RouterGroup) {
	rg.PUT("/me/profile", h.UpdateProfile)
	rg.POST("/me/avatar", h.UpdateAvatar)
}

func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.svc.Create(c.Request.Context(), req.Name, req.Email, req.Age)
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.svc.Update(c.Request.Context(), id, req.Name, req.Email, req.Age)
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		mapError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *UserHandler) List(c *gin.Context) {
	page := parseIntDefault(c.Query("page"), 1)
	size := parseIntDefault(c.Query("size"), 20)

	users, err := h.svc.List(c.Request.Context(), page, size)
	if err != nil {
		mapError(c, err)
		return
	}
	if users == nil {
		users = []*domain.User{}
	}
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.svc.GetProfile(c.Request.Context(), id)
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.svc.UpdateProfile(c.Request.Context(), userID, req.Name, req.Bio, req.Phone, req.City, req.Website)
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateAvatar(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "avatar file required"})
		return
	}
	defer file.Close()

	if header.Size > maxAvatarSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "avatar file too large (max 2MB)"})
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedAvatarExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file type, allowed: jpg, jpeg, png, gif, webp"})
		return
	}

	avatarDir := "data/avatars"
	filename := fmt.Sprintf("%d%s", userID, ext)
	savePath := filepath.Join(avatarDir, filename)

	if err := c.SaveUploadedFile(header, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save avatar"})
		return
	}

	avatarPath := fmt.Sprintf("/avatars/%s", filename)
	user, err := h.svc.UpdateAvatar(c.Request.Context(), userID, avatarPath)
	if err != nil {
		mapError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func parseID(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func parseIntDefault(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 {
		return def
	}
	return n
}

func mapError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, repository.ErrExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
