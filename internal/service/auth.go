package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mikasa/skk/internal/domain"
	"github.com/mikasa/skk/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      repository.UserRepository
	jwtSecret []byte
}

func NewAuthService(repo repository.UserRepository) *AuthService {
	secret := os.Getenv("SKK_JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-me"
	}
	return &AuthService{
		repo:      repo,
		jwtSecret: []byte(secret),
	}
}

type TokenClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *AuthService) Register(ctx context.Context, username, email, password, name string) (*domain.User, error) {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)
	name = strings.TrimSpace(name)

	if err := validateUsername(username); err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}
	if err := validateEmail(email); err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}
	if err := validatePassword(password); err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}
	if name == "" {
		name = username
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, fmt.Errorf("register: hash password: %w", err)
	}

	now := time.Now()
	user := &domain.User{
		Username:     username,
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	created, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}
	return created, nil
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, *domain.User, error) {
	username = strings.TrimSpace(username)

	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", nil, fmt.Errorf("login: %w", repository.ErrInvalidCredentials)
		}
		return "", nil, fmt.Errorf("login: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, fmt.Errorf("login: %w", repository.ErrInvalidCredentials)
	}

	token, err := s.generateToken(user.ID, user.Username)
	if err != nil {
		return "", nil, fmt.Errorf("login: generate token: %w", err)
	}

	return token, user, nil
}

func (s *AuthService) GetCurrentUser(ctx context.Context, userID int64) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get current user: %w", err)
	}
	return user, nil
}

func (s *AuthService) generateToken(userID int64, username string) (string, error) {
	claims := TokenClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) ParseToken(tokenStr string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func validateUsername(username string) error {
	if username == "" {
		return errors.New("username must not be empty")
	}
	if len(username) < 3 || len(username) > 32 {
		return errors.New("username must be between 3 and 32 characters")
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}
