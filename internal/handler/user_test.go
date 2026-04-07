package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mikasa/skk/internal/domain"
	"github.com/mikasa/skk/internal/repository"
)

type mockService struct {
	users  map[int64]*domain.User
	nextID int64
	err    error
}

func newMockService() *mockService {
	return &mockService{
		users:  make(map[int64]*domain.User),
		nextID: 1,
	}
}

func (m *mockService) Create(_ context.Context, name, email string, age int) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	now := time.Now()
	user := &domain.User{
		ID: m.nextID, Name: name, Email: email, Age: age,
		CreatedAt: now, UpdatedAt: now,
	}
	m.users[m.nextID] = user
	m.nextID++
	return user, nil
}

func (m *mockService) GetByID(_ context.Context, id int64) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	user, ok := m.users[id]
	if !ok {
		return nil, fmt.Errorf("get user: %w", repository.ErrNotFound)
	}
	return user, nil
}

func (m *mockService) Update(_ context.Context, id int64, name, email string, age int) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	user, ok := m.users[id]
	if !ok {
		return nil, fmt.Errorf("update user: %w", repository.ErrNotFound)
	}
	user.Name = name
	user.Email = email
	user.Age = age
	return user, nil
}

func (m *mockService) Delete(_ context.Context, id int64) error {
	if m.err != nil {
		return m.err
	}
	if _, ok := m.users[id]; !ok {
		return fmt.Errorf("delete user: %w", repository.ErrNotFound)
	}
	delete(m.users, id)
	return nil
}

func (m *mockService) List(_ context.Context, page, size int) ([]*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	var all []*domain.User
	for _, u := range m.users {
		all = append(all, u)
	}
	return all, nil
}

func (m *mockService) GetProfile(_ context.Context, id int64) (*domain.User, error) {
	return m.GetByID(context.Background(), id)
}

func (m *mockService) UpdateProfile(_ context.Context, userID int64, name, bio, phone, city, website string) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	user, ok := m.users[userID]
	if !ok {
		return nil, fmt.Errorf("update profile: %w", repository.ErrNotFound)
	}
	if name != "" {
		user.Name = name
	}
	user.Bio = bio
	user.Phone = phone
	user.City = city
	user.Website = website
	return user, nil
}

func (m *mockService) UpdateAvatar(_ context.Context, userID int64, avatarPath string) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	user, ok := m.users[userID]
	if !ok {
		return nil, fmt.Errorf("update avatar: %w", repository.ErrNotFound)
	}
	user.Avatar = avatarPath
	return user, nil
}

func setupRouter(svc UserService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := NewUserHandler(svc)
	handler.RegisterPublicRoutes(r.Group("/users"))
	return r
}

func TestHandler_Create_Success(t *testing.T) {
	r := setupRouter(newMockService())

	body, _ := json.Marshal(CreateUserRequest{
		Name: "Alice", Email: "alice@example.com", Age: 25,
	})
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}
}

func TestHandler_Create_InvalidJSON(t *testing.T) {
	r := setupRouter(newMockService())

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_GetByID_Success(t *testing.T) {
	ms := newMockService()
	r := setupRouter(ms)

	created, _ := ms.Create(context.Background(), "Alice", "alice@example.com", 25)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%d", created.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestHandler_GetByID_NotFound(t *testing.T) {
	r := setupRouter(newMockService())

	req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestHandler_GetByID_InvalidID(t *testing.T) {
	r := setupRouter(newMockService())

	req := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_Update_Success(t *testing.T) {
	ms := newMockService()
	r := setupRouter(ms)

	created, _ := ms.Create(context.Background(), "Alice", "alice@example.com", 25)

	body, _ := json.Marshal(UpdateUserRequest{
		Name: "Alice Updated", Email: "alice2@example.com", Age: 26,
	})
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/users/%d", created.ID), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestHandler_Update_NotFound(t *testing.T) {
	r := setupRouter(newMockService())

	body, _ := json.Marshal(UpdateUserRequest{
		Name: "Alice", Email: "alice@example.com", Age: 25,
	})
	req := httptest.NewRequest(http.MethodPut, "/users/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestHandler_Delete_Success(t *testing.T) {
	ms := newMockService()
	r := setupRouter(ms)

	created, _ := ms.Create(context.Background(), "Alice", "alice@example.com", 25)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%d", created.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestHandler_Delete_NotFound(t *testing.T) {
	r := setupRouter(newMockService())

	req := httptest.NewRequest(http.MethodDelete, "/users/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestHandler_List_Success(t *testing.T) {
	ms := newMockService()
	r := setupRouter(ms)

	ms.Create(context.Background(), "A", "a@test.com", 20)
	ms.Create(context.Background(), "B", "b@test.com", 21)

	req := httptest.NewRequest(http.MethodGet, "/users?page=1&size=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var users []domain.User
	json.Unmarshal(w.Body.Bytes(), &users)
	if len(users) != 2 {
		t.Errorf("List() returned %d users, want 2", len(users))
	}
}

func TestHandler_List_Empty(t *testing.T) {
	r := setupRouter(newMockService())

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestHandler_InternalError(t *testing.T) {
	ms := &mockService{err: errors.New("unexpected"), users: make(map[int64]*domain.User)}
	r := setupRouter(ms)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestHandler_GetProfile_Success(t *testing.T) {
	ms := newMockService()
	r := setupRouter(ms)

	created, _ := ms.Create(context.Background(), "Alice", "alice@example.com", 25)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%d/profile", created.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestHandler_GetProfile_NotFound(t *testing.T) {
	r := setupRouter(newMockService())

	req := httptest.NewRequest(http.MethodGet, "/users/999/profile", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}
