package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mikasa/skk/internal/domain"
	"github.com/mikasa/skk/internal/middleware"
	"github.com/mikasa/skk/internal/repository"
	"github.com/mikasa/skk/internal/service"
)

type mockAuthService struct {
	users      map[int64]*domain.User
	nextID     int64
	tokens     map[int64]string // user_id -> token
	err        error
	authSvc    *service.AuthService // real auth service for token generation
}

func newMockAuthService() *mockAuthService {
	return &mockAuthService{
		users:  make(map[int64]*domain.User),
		nextID: 1,
		tokens: make(map[int64]string),
	}
}

func (m *mockAuthService) Register(_ context.Context, username, email, password, name string) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, u := range m.users {
		if u.Username == username {
			return nil, fmt.Errorf("register: %w", repository.ErrExists)
		}
		if u.Email == email {
			return nil, fmt.Errorf("register: %w", repository.ErrExists)
		}
	}
	user := &domain.User{
		ID:       m.nextID,
		Username: username,
		Name:     name,
		Email:    email,
	}
	m.users[m.nextID] = user
	m.nextID++
	return user, nil
}

func (m *mockAuthService) Login(_ context.Context, username, password string) (string, *domain.User, error) {
	if m.err != nil {
		return "", nil, m.err
	}
	for _, u := range m.users {
		if u.Username == username {
			token := fmt.Sprintf("token-%d", u.ID)
			return token, u, nil
		}
	}
	return "", nil, fmt.Errorf("login: %w", repository.ErrInvalidCredentials)
}

func (m *mockAuthService) GetCurrentUser(_ context.Context, userID int64) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	user, ok := m.users[userID]
	if !ok {
		return nil, fmt.Errorf("get current user: %w", repository.ErrNotFound)
	}
	return user, nil
}

func setupAuthRouter(svc AuthService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := NewAuthHandler(svc)
	handler.RegisterRoutes(r)
	return r
}

func TestAuthHandler_Register_Success(t *testing.T) {
	r := setupAuthRouter(newMockAuthService())

	body, _ := json.Marshal(RegisterRequest{
		Username: "zhangsan",
		Email:    "zhangsan@example.com",
		Password: "password123",
		Name:     "张三",
	})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}
}

func TestAuthHandler_Register_MissingFields(t *testing.T) {
	r := setupAuthRouter(newMockAuthService())

	body, _ := json.Marshal(map[string]string{
		"username": "zhangsan",
	})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestAuthHandler_Login_Success(t *testing.T) {
	mockSvc := newMockAuthService()
	r := setupAuthRouter(mockSvc)

	mockSvc.Register(context.Background(), "zhangsan", "zhangsan@example.com", "password123", "张三")

	body, _ := json.Marshal(LoginRequest{
		Username: "zhangsan",
		Password: "password123",
	})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp LoginResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Token == "" {
		t.Error("Login() token should not be empty")
	}
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	r := setupAuthRouter(newMockAuthService())

	body, _ := json.Marshal(LoginRequest{
		Username: "nonexistent",
		Password: "password123",
	})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestAuthHandler_Me_WithAuth(t *testing.T) {
	mockSvc := newMockAuthService()
	mockSvc.Register(context.Background(), "zhangsan", "zhangsan@example.com", "password123", "张三")

	gin.SetMode(gin.TestMode)
	r := gin.New()
	authHandler := NewAuthHandler(mockSvc)
	authHandler.RegisterRoutes(r)

	// Override /auth/me to use a middleware that sets user_id
	// For testing, we'll add a test route that mimics the auth middleware
	r.GET("/test/me", func(c *gin.Context) {
		c.Set(middleware.UserIDKey, int64(1))
	}, authHandler.Me)

	req := httptest.NewRequest(http.MethodGet, "/test/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestAuthHandler_Me_NoAuth(t *testing.T) {
	r := setupAuthRouter(newMockAuthService())

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}
