package domain

import "time"

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Age          int       `json:"age"`
	Avatar       string    `json:"avatar,omitempty"`
	Bio          string    `json:"bio,omitempty"`
	Phone        string    `json:"phone,omitempty"`
	City         string    `json:"city,omitempty"`
	Website      string    `json:"website,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
