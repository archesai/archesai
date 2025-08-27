package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user entity
type User struct {
	ID            uuid.UUID `json:"id" db:"id"`
	Email         string    `json:"email" db:"email"`
	Name          string    `json:"name" db:"name"`
	PasswordHash  string    `json:"-" db:"password_hash"`
	EmailVerified bool      `json:"email_verified" db:"email_verified"`
	Image         *string   `json:"image,omitempty" db:"image"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// SignUpRequest represents a sign-up request
type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required"`
}

// SignInRequest represents a sign-in request
type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdateUserRequest represents a user update request
type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty"`
	Image *string `json:"image,omitempty"`
}