package sqlite

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mikasa/skk/internal/domain"
	"github.com/mikasa/skk/internal/repository"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type userDAO struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	Username     string    `gorm:"not null;uniqueIndex"`
	Name         string    `gorm:"not null"`
	Email        string    `gorm:"not null;uniqueIndex"`
	PasswordHash string    `gorm:"not null"`
	Age          int       `gorm:"not null"`
	Avatar       string
	Bio          string
	Phone        string
	City         string
	Website      string
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
}

func (userDAO) TableName() string {
	return "users"
}

func (u *userDAO) toDomain() *domain.User {
	return &domain.User{
		ID:           u.ID,
		Username:     u.Username,
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Age:          u.Age,
		Avatar:       u.Avatar,
		Bio:          u.Bio,
		Phone:        u.Phone,
		City:         u.City,
		Website:      u.Website,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func fromDomain(u *domain.User) *userDAO {
	return &userDAO{
		ID:           u.ID,
		Username:     u.Username,
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Age:          u.Age,
		Avatar:       u.Avatar,
		Bio:          u.Bio,
		Phone:        u.Phone,
		City:         u.City,
		Website:      u.Website,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(dbPath string) (*UserRepository, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get underlying db: %w", err)
	}
	sqlDB.SetMaxOpenConns(1)

	if err := db.AutoMigrate(&userDAO{}); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return &UserRepository{db: db}, nil
}

func (r *UserRepository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return fmt.Errorf("get underlying db: %w", err)
	}
	return sqlDB.Close()
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	dao := fromDomain(user)
	result := r.db.WithContext(ctx).Create(dao)
	if result.Error != nil {
		if isUniqueViolation(result.Error) {
			return nil, fmt.Errorf("create user: %w", repository.ErrExists)
		}
		return nil, fmt.Errorf("create user: %w", result.Error)
	}
	return dao.toDomain(), nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	var dao userDAO
	result := r.db.WithContext(ctx).First(&dao, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get user %d: %w", id, repository.ErrNotFound)
		}
		return nil, fmt.Errorf("get user %d: %w", id, result.Error)
	}
	return dao.toDomain(), nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var dao userDAO
	result := r.db.WithContext(ctx).Where("username = ?", username).First(&dao)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("get user by username %q: %w", username, repository.ErrNotFound)
		}
		return nil, fmt.Errorf("get user by username %q: %w", username, result.Error)
	}
	return dao.toDomain(), nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	dao := fromDomain(user)
	result := r.db.WithContext(ctx).Model(&userDAO{}).Where("id = ?", user.ID).Updates(map[string]any{
		"name":        dao.Name,
		"email":       dao.Email,
		"age":         dao.Age,
		"updated_at":  dao.UpdatedAt,
	})
	if result.Error != nil {
		if isUniqueViolation(result.Error) {
			return nil, fmt.Errorf("update user %d: %w", user.ID, repository.ErrExists)
		}
		return nil, fmt.Errorf("update user %d: %w", user.ID, result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("update user %d: %w", user.ID, repository.ErrNotFound)
	}
	return dao.toDomain(), nil
}

func (r *UserRepository) UpdateProfile(ctx context.Context, user *domain.User) (*domain.User, error) {
	dao := fromDomain(user)
	result := r.db.WithContext(ctx).Model(&userDAO{}).Where("id = ?", user.ID).Updates(map[string]any{
		"name":       dao.Name,
		"bio":        dao.Bio,
		"phone":      dao.Phone,
		"city":       dao.City,
		"website":    dao.Website,
		"avatar":     dao.Avatar,
		"updated_at": dao.UpdatedAt,
	})
	if result.Error != nil {
		return nil, fmt.Errorf("update profile %d: %w", user.ID, result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("update profile %d: %w", user.ID, repository.ErrNotFound)
	}
	return dao.toDomain(), nil
}

func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&userDAO{}, id)
	if result.Error != nil {
		return fmt.Errorf("delete user %d: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("delete user %d: %w", id, repository.ErrNotFound)
	}
	return nil
}

func (r *UserRepository) List(ctx context.Context, offset, limit int) ([]*domain.User, error) {
	var daos []userDAO
	result := r.db.WithContext(ctx).Order("id").Limit(limit).Offset(offset).Find(&daos)
	if result.Error != nil {
		return nil, fmt.Errorf("list users: %w", result.Error)
	}

	users := make([]*domain.User, len(daos))
	for i := range daos {
		users[i] = daos[i].toDomain()
	}
	return users, nil
}

func isUniqueViolation(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}
