package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/mikasa/skk/internal/domain"
	"github.com/mikasa/skk/internal/repository"
)

type Service struct {
	repo repository.UserRepository
}

func New(repo repository.UserRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, name, email string, age int) (*domain.User, error) {
	if err := validateName(name); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	if err := validateEmail(email); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	if err := validateAge(age); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	now := time.Now()
	user := &domain.User{
		Name:      strings.TrimSpace(name),
		Email:     strings.TrimSpace(email),
		Age:       age,
		CreatedAt: now,
		UpdatedAt: now,
	}

	created, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}

func (s *Service) Update(ctx context.Context, id int64, name, email string, age int) (*domain.User, error) {
	if err := validateName(name); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	if err := validateEmail(email); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	if err := validateAge(age); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	existing.Name = strings.TrimSpace(name)
	existing.Email = strings.TrimSpace(email)
	existing.Age = age
	existing.UpdatedAt = time.Now()

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}
	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

func (s *Service) List(ctx context.Context, page, size int) ([]*domain.User, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	offset := (page - 1) * size

	users, err := s.repo.List(ctx, offset, size)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return users, nil
}

// GetProfile returns a user's public profile.
func (s *Service) GetProfile(ctx context.Context, id int64) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get profile: %w", err)
	}
	return user, nil
}

// UpdateProfile updates the authenticated user's profile fields.
func (s *Service) UpdateProfile(ctx context.Context, userID int64, name, bio, phone, city, website string) (*domain.User, error) {
	existing, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("update profile: %w", err)
	}

	if name != "" {
		existing.Name = strings.TrimSpace(name)
	}
	existing.Bio = bio
	existing.Phone = phone
	existing.City = city
	existing.Website = website
	existing.UpdatedAt = time.Now()

	updated, err := s.repo.UpdateProfile(ctx, existing)
	if err != nil {
		return nil, fmt.Errorf("update profile: %w", err)
	}
	return updated, nil
}

// UpdateAvatar updates the user's avatar path.
func (s *Service) UpdateAvatar(ctx context.Context, userID int64, avatarPath string) (*domain.User, error) {
	existing, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("update avatar: %w", err)
	}

	existing.Avatar = avatarPath
	existing.UpdatedAt = time.Now()

	updated, err := s.repo.UpdateProfile(ctx, existing)
	if err != nil {
		return nil, fmt.Errorf("update avatar: %w", err)
	}
	return updated, nil
}

func validateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name must not be empty")
	}
	return nil
}

func validateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email must not be empty")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("email format invalid")
	}
	return nil
}

func validateAge(age int) error {
	if age < 0 || age > 150 {
		return errors.New("age must be between 0 and 150")
	}
	return nil
}
