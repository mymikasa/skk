package sqlite

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mikasa/skk/internal/domain"
	"github.com/mikasa/skk/internal/repository"
)

func newTestRepo(t *testing.T) *UserRepository {
	t.Helper()
	repo, err := NewUserRepository(":memory:")
	if err != nil {
		t.Fatalf("create repo: %v", err)
	}
	t.Cleanup(func() { repo.Close() })
	return repo
}

func newUser(name, email string, age int) *domain.User {
	now := time.Now().Truncate(time.Second)
	return &domain.User{
		Username:     email,
		Name:         name,
		Email:        email,
		PasswordHash: "$2a$12$dummyhash",
		Age:          age,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func TestUserRepository_Create(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	user := newUser("Alice", "alice@example.com", 25)
	created, err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.ID == 0 {
		t.Error("Create() ID should not be zero")
	}
	if created.Name != "Alice" {
		t.Errorf("Create() Name = %q, want %q", created.Name, "Alice")
	}
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	repo.Create(ctx, newUser("A", "dup@example.com", 20))
	_, err := repo.Create(ctx, newUser("B", "dup@example.com", 21))
	if !errors.Is(err, repository.ErrExists) {
		t.Errorf("Create() error = %v, want ErrExists", err)
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	created, _ := repo.Create(ctx, newUser("Bob", "bob@example.com", 30))
	got, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if got.Name != "Bob" {
		t.Errorf("GetByID() Name = %q, want %q", got.Name, "Bob")
	}
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, 999)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("GetByID() error = %v, want ErrNotFound", err)
	}
}

func TestUserRepository_GetByUsername(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	created, _ := repo.Create(ctx, newUser("Bob", "bob@example.com", 30))
	got, err := repo.GetByUsername(ctx, created.Username)
	if err != nil {
		t.Fatalf("GetByUsername() error = %v", err)
	}
	if got.Name != "Bob" {
		t.Errorf("GetByUsername() Name = %q, want %q", got.Name, "Bob")
	}
}

func TestUserRepository_GetByUsername_NotFound(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	_, err := repo.GetByUsername(ctx, "nonexistent")
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("GetByUsername() error = %v, want ErrNotFound", err)
	}
}

func TestUserRepository_Update(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	created, _ := repo.Create(ctx, newUser("Charlie", "charlie@example.com", 28))
	created.Name = "Charles"
	created.UpdatedAt = time.Now().Truncate(time.Second)

	updated, err := repo.Update(ctx, created)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if updated.Name != "Charles" {
		t.Errorf("Update() Name = %q, want %q", updated.Name, "Charles")
	}
}

func TestUserRepository_Update_NotFound(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	user := newUser("Ghost", "ghost@example.com", 1)
	user.ID = 999
	_, err := repo.Update(ctx, user)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("Update() error = %v, want ErrNotFound", err)
	}
}

func TestUserRepository_Update_DuplicateEmail(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	repo.Create(ctx, newUser("A", "a@example.com", 20))
	created, _ := repo.Create(ctx, newUser("B", "b@example.com", 21))
	created.Email = "a@example.com"
	_, err := repo.Update(ctx, created)
	if !errors.Is(err, repository.ErrExists) {
		t.Errorf("Update() error = %v, want ErrExists", err)
	}
}

func TestUserRepository_UpdateProfile(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	created, _ := repo.Create(ctx, newUser("Alice", "alice@example.com", 25))
	created.Bio = "Go developer"
	created.City = "Beijing"
	created.Phone = "+86-138xxxx"
	created.Website = "https://example.com"
	created.UpdatedAt = time.Now().Truncate(time.Second)

	updated, err := repo.UpdateProfile(ctx, created)
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

func TestUserRepository_UpdateProfile_NotFound(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	user := newUser("Ghost", "ghost@example.com", 1)
	user.ID = 999
	_, err := repo.UpdateProfile(ctx, user)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("UpdateProfile() error = %v, want ErrNotFound", err)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	created, _ := repo.Create(ctx, newUser("Dave", "dave@example.com", 35))

	if err := repo.Delete(ctx, created.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err := repo.GetByID(ctx, created.ID)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("GetByID() after delete error = %v, want ErrNotFound", err)
	}
}

func TestUserRepository_Delete_NotFound(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	err := repo.Delete(ctx, 999)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("Delete() error = %v, want ErrNotFound", err)
	}
}

func TestUserRepository_List(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	repo.Create(ctx, newUser("A", "a@example.com", 20))
	repo.Create(ctx, newUser("B", "b@example.com", 21))
	repo.Create(ctx, newUser("C", "c@example.com", 22))

	users, err := repo.List(ctx, 0, 2)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(users) != 2 {
		t.Errorf("List() returned %d users, want 2", len(users))
	}

	users2, err := repo.List(ctx, 2, 2)
	if err != nil {
		t.Fatalf("List() page 2 error = %v", err)
	}
	if len(users2) != 1 {
		t.Errorf("List() page 2 returned %d users, want 1", len(users2))
	}
}

func TestUserRepository_List_Empty(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()

	users, err := repo.List(ctx, 0, 10)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(users) != 0 {
		t.Errorf("List() returned %d users, want 0", len(users))
	}
}
