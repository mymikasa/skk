package service

import (
	"context"
	"errors"
	"testing"

	"github.com/mikasa/skk/internal/repository"
)

func TestAuthService_Register_Success(t *testing.T) {
	svc := NewAuthService(newMockRepo())
	user, err := svc.Register(context.Background(), "zhangsan", "zhangsan@example.com", "password123", "张三")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if user.ID == 0 {
		t.Error("Register() ID should not be zero")
	}
	if user.Username != "zhangsan" {
		t.Errorf("Register() Username = %q, want %q", user.Username, "zhangsan")
	}
	if user.PasswordHash == "" {
		t.Error("Register() PasswordHash should be set")
	}
}

func TestAuthService_Register_ShortUsername(t *testing.T) {
	svc := NewAuthService(newMockRepo())
	_, err := svc.Register(context.Background(), "ab", "a@test.com", "password123", "")
	if err == nil {
		t.Fatal("Register() expected error for short username")
	}
}

func TestAuthService_Register_EmptyUsername(t *testing.T) {
	svc := NewAuthService(newMockRepo())
	_, err := svc.Register(context.Background(), "", "a@test.com", "password123", "")
	if err == nil {
		t.Fatal("Register() expected error for empty username")
	}
}

func TestAuthService_Register_InvalidEmail(t *testing.T) {
	svc := NewAuthService(newMockRepo())
	_, err := svc.Register(context.Background(), "user", "not-email", "password123", "")
	if err == nil {
		t.Fatal("Register() expected error for invalid email")
	}
}

func TestAuthService_Register_ShortPassword(t *testing.T) {
	svc := NewAuthService(newMockRepo())
	_, err := svc.Register(context.Background(), "user", "a@test.com", "12345", "")
	if err == nil {
		t.Fatal("Register() expected error for short password")
	}
}

func TestAuthService_Register_DuplicateUsername(t *testing.T) {
	repo := newMockRepo()
	svc := NewAuthService(repo)
	svc.Register(context.Background(), "zhangsan", "a@test.com", "password123", "")
	_, err := svc.Register(context.Background(), "zhangsan", "b@test.com", "password123", "")
	if !errors.Is(err, repository.ErrExists) {
		t.Errorf("Register() error = %v, want ErrExists", err)
	}
}

func TestAuthService_Register_DefaultName(t *testing.T) {
	svc := NewAuthService(newMockRepo())
	user, err := svc.Register(context.Background(), "zhangsan", "a@test.com", "password123", "")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if user.Name != "zhangsan" {
		t.Errorf("Register() Name = %q, want username as default", user.Name)
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	svc := NewAuthService(newMockRepo())
	svc.Register(context.Background(), "zhangsan", "a@test.com", "password123", "张三")

	token, user, err := svc.Login(context.Background(), "zhangsan", "password123")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if token == "" {
		t.Error("Login() token should not be empty")
	}
	if user.Username != "zhangsan" {
		t.Errorf("Login() Username = %q, want %q", user.Username, "zhangsan")
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	svc := NewAuthService(newMockRepo())
	svc.Register(context.Background(), "zhangsan", "a@test.com", "password123", "")

	_, _, err := svc.Login(context.Background(), "zhangsan", "wrongpassword")
	if !errors.Is(err, repository.ErrInvalidCredentials) {
		t.Errorf("Login() error = %v, want ErrInvalidCredentials", err)
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	svc := NewAuthService(newMockRepo())
	_, _, err := svc.Login(context.Background(), "nonexistent", "password123")
	if !errors.Is(err, repository.ErrInvalidCredentials) {
		t.Errorf("Login() error = %v, want ErrInvalidCredentials", err)
	}
}

func TestAuthService_ParseToken_Valid(t *testing.T) {
	svc := NewAuthService(newMockRepo())
	svc.Register(context.Background(), "zhangsan", "a@test.com", "password123", "")
	token, _, _ := svc.Login(context.Background(), "zhangsan", "password123")

	claims, err := svc.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken() error = %v", err)
	}
	if claims.UserID != 1 {
		t.Errorf("ParseToken() UserID = %d, want 1", claims.UserID)
	}
	if claims.Username != "zhangsan" {
		t.Errorf("ParseToken() Username = %q, want %q", claims.Username, "zhangsan")
	}
}

func TestAuthService_ParseToken_Invalid(t *testing.T) {
	svc := NewAuthService(newMockRepo())
	_, err := svc.ParseToken("invalid-token")
	if err == nil {
		t.Fatal("ParseToken() expected error for invalid token")
	}
}

func TestAuthService_GetCurrentUser_Success(t *testing.T) {
	svc := NewAuthService(newMockRepo())
	created, _ := svc.Register(context.Background(), "zhangsan", "a@test.com", "password123", "张三")

	user, err := svc.GetCurrentUser(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetCurrentUser() error = %v", err)
	}
	if user.Username != "zhangsan" {
		t.Errorf("GetCurrentUser() Username = %q, want %q", user.Username, "zhangsan")
	}
}

func TestAuthService_GetCurrentUser_NotFound(t *testing.T) {
	svc := NewAuthService(newMockRepo())
	_, err := svc.GetCurrentUser(context.Background(), 999)
	if !errors.Is(err, repository.ErrNotFound) {
		t.Errorf("GetCurrentUser() error = %v, want ErrNotFound", err)
	}
}
