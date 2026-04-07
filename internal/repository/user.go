package repository

import (
	"context"
	"errors"

	"github.com/mikasa/skk/internal/domain"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrExists         = errors.New("already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) (*domain.User, error)
	UpdateProfile(ctx context.Context, user *domain.User) (*domain.User, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, offset, limit int) ([]*domain.User, error)
}
