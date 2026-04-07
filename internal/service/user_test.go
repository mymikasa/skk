package service

import (
	"context"
	"errors"
	"testing"

	"github.com/mikasa/skk/internal/domain"
	"github.com/mikasa/skk/internal/repository"
)

type mockRepo struct {
	users  map[int64]*domain.User
	nextID int64
	err    error
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		users:  make(map[int64]*domain.User),
		nextID: 1,
	}
}

func (m *mockRepo) Create(_ context.Context, user *domain.User) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, existing := range m.users {
		if existing.Email == user.Email {
			return nil, repository.ErrExists
		}
		if user.Username != "" && existing.Username == user.Username {
			return nil, repository.ErrExists
		}
	}
	user.ID = m.nextID
	m.nextID++
	m.users[user.ID] = user
	return user, nil
}

func (m *mockRepo) GetByID(_ context.Context, id int64) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	user, ok := m.users[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return user, nil
}

func (m *mockRepo) GetByUsername(_ context.Context, username string) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, u := range m.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, repository.ErrNotFound
}

func (m *mockRepo) Update(_ context.Context, user *domain.User) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, existing := range m.users {
		if existing.Email == user.Email && existing.ID != user.ID {
			return nil, repository.ErrExists
		}
	}
	if _, ok := m.users[user.ID]; !ok {
		return nil, repository.ErrNotFound
	}
	m.users[user.ID] = user
	return user, nil
}

func (m *mockRepo) UpdateProfile(_ context.Context, user *domain.User) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	if _, ok := m.users[user.ID]; !ok {
		return nil, repository.ErrNotFound
	}
	m.users[user.ID] = user
	return user, nil
}

func (m *mockRepo) Delete(_ context.Context, id int64) error {
	if m.err != nil {
		return m.err
	}
	if _, ok := m.users[id]; !ok {
		return repository.ErrNotFound
	}
	delete(m.users, id)
	return nil
}

func (m *mockRepo) List(_ context.Context, offset, limit int) ([]*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	var all []*domain.User
	for _, u := range m.users {
		all = append(all, u)
	}
	if offset >= len(all) {
		return nil, nil
	}
	end := offset + limit
	if end > len(all) {
		end = len(all)
	}
	return all[offset:end], nil
}

func TestService_Create_Success(t *testing.T) {
	svc := New(newMockRepo())
	user, err := svc.Create(context.Background(), "Alice", "alice@example.com", 25)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if user.ID == 0 {
		t.Error("Create() ID should not be zero")
	}
}

func TestService_Create_EmptyName(t *testing.T) {
	svc := New(newMockRepo())
	_, err := svc.Create(context.Background(), "", "alice@example.com", 25)
	if err == nil {
		t.Fatal("Create() expected error for empty name")
	}
}

func TestService_Create_InvalidEmail(t *testing.T) {
	svc := New(newMockRepo())
	_, err := svc.Create(context.Background(), "Alice", "not-an-email", 25)
	if err == nil {
		t.Fatal("Create() expected error for invalid email")
	}
}

func TestService_Create_InvalidAge(t *testing.T) {
	tests := []struct {
		name string
		age  int
	}{
		{"negative", -1},
		{"too large", 200},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := New(newMockRepo())
			_, err := svc.Create(context.Background(), "Alice", "alice@example.com", tt.age)
			if err == nil {
				t.Errorf("Create() expected error for age %d", tt.age)
			}
		})
	}
}

func TestService_Create_DuplicateEmail(t *testing.T) {
	svc := New(newMockRepo())
	svc.Create(context.Background(), "Alice", "alice@example.com", 25)
	_, err := svc.Create(context.Background(), "Bob", "alice@example.com", 30)
	if !errors.Is(err, repository.ErrExists) {
		t.Errorf("Create() error = %v, want ErrExists", err)
	}
}

func TestService_GetByID_Success(t *testing.T) {
	repo := newMockRepo()
	svc := New(repo)
	created, _ := svc.Create(context.Background(), "Alice", "alice@example.com", 25)
	got, err := svc.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got.Name != "Alice" {
		t.Errorf("GetByID() Name = %q, want %q", got.Name, "Alice")
	}
}

func TestService_GetByID_NotFound(t *testing.T) {
	svc := New(newMockRepo())
	_, err := svc.GetByID(context.Background(), 999)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("GetByID() error = %v, want ErrNotFound", err)
	}
}

func TestService_Update_Success(t *testing.T) {
	repo := newMockRepo()
	svc := New(repo)
	created, _ := svc.Create(context.Background(), "Alice", "alice@example.com", 25)
	updated, err := svc.Update(context.Background(), created.ID, "Alice Updated", "alice2@example.com", 26)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if updated.Name != "Alice Updated" {
		t.Errorf("Update() Name = %q, want %q", updated.Name, "Alice Updated")
	}
}

func TestService_Update_NotFound(t *testing.T) {
	svc := New(newMockRepo())
	_, err := svc.Update(context.Background(), 999, "Alice", "alice@example.com", 25)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("Update() error = %v, want ErrNotFound", err)
	}
}

func TestService_Update_DuplicateEmail(t *testing.T) {
	repo := newMockRepo()
	svc := New(repo)
	svc.Create(context.Background(), "A", "a@test.com", 20)
	created, _ := svc.Create(context.Background(), "B", "b@test.com", 21)
	_, err := svc.Update(context.Background(), created.ID, "B", "a@test.com", 21)
	if !errors.Is(err, repository.ErrExists) {
		t.Errorf("Update() error = %v, want ErrExists", err)
	}
}

func TestService_Delete_Success(t *testing.T) {
	repo := newMockRepo()
	svc := New(repo)
	created, _ := svc.Create(context.Background(), "Alice", "alice@example.com", 25)
	if err := svc.Delete(context.Background(), created.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
}

func TestService_Delete_NotFound(t *testing.T) {
	svc := New(newMockRepo())
	err := svc.Delete(context.Background(), 999)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}

func TestService_List_Success(t *testing.T) {
	repo := newMockRepo()
	svc := New(repo)
	svc.Create(context.Background(), "A", "a@test.com", 20)
	svc.Create(context.Background(), "B", "b@test.com", 21)
	svc.Create(context.Background(), "C", "c@test.com", 22)

	users, err := svc.List(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(users) != 2 {
		t.Errorf("List() returned %d users, want 2", len(users))
	}
}

func TestService_List_Defaults(t *testing.T) {
	svc := New(newMockRepo())
	users, err := svc.List(context.Background(), 0, 0)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if users != nil {
		t.Errorf("List() expected nil for empty repo, got %v", users)
	}
}

func TestService_UpdateProfile_Success(t *testing.T) {
	repo := newMockRepo()
	svc := New(repo)
	created, _ := svc.Create(context.Background(), "Alice", "alice@example.com", 25)

	updated, err := svc.UpdateProfile(context.Background(), created.ID, "Alice Smith", "Go developer", "+86-138xxxx", "Beijing", "https://example.com")
	if err != nil {
		t.Fatalf("UpdateProfile() error = %v", err)
	}
	if updated.Bio != "Go developer" {
		t.Errorf("UpdateProfile() Bio = %q, want %q", updated.Bio, "Go developer")
	}
	if updated.City != "Beijing" {
		t.Errorf("UpdateProfile() City = %q, want %q", updated.City, "Beijing")
	}
}

func TestService_UpdateProfile_NotFound(t *testing.T) {
	svc := New(newMockRepo())
	_, err := svc.UpdateProfile(context.Background(), 999, "name", "bio", "phone", "city", "website")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("UpdateProfile() error = %v, want ErrNotFound", err)
	}
}

func TestService_UpdateAvatar_Success(t *testing.T) {
	repo := newMockRepo()
	svc := New(repo)
	created, _ := svc.Create(context.Background(), "Alice", "alice@example.com", 25)

	updated, err := svc.UpdateAvatar(context.Background(), created.ID, "/avatars/1.jpg")
	if err != nil {
		t.Fatalf("UpdateAvatar() error = %v", err)
	}
	if updated.Avatar != "/avatars/1.jpg" {
		t.Errorf("UpdateAvatar() Avatar = %q, want %q", updated.Avatar, "/avatars/1.jpg")
	}
}
