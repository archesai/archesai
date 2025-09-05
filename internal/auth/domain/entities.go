// Package domain contains the auth domain business logic and entities
package domain

import (
	"errors"
	"time"

	"github.com/archesai/archesai/internal/auth/generated/api"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Domain errors
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrSessionNotFound    = errors.New("session not found")
	ErrSessionExpired     = errors.New("session expired")
	ErrAccountNotFound    = errors.New("account not found")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// User extends the generated API UserEntity with auth-specific fields
type User struct {
	api.UserEntity
	PasswordHash string `json:"-"` // Never expose password hash
}

// Account extends the generated API AccountEntity with auth-specific fields
type Account struct {
	api.AccountEntity
	Password string `json:"-"` // Store hashed password for local accounts
}

// Session extends the generated API SessionEntity
type Session struct {
	api.SessionEntity
}

// Claims represents JWT token claims
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
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

// Tokens contains authentication tokens
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// TokenResponse represents authentication token response
type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// NewAccount creates a new account from the API entity
func NewAccount(entity api.AccountEntity) *Account {
	return &Account{AccountEntity: entity}
}

// NewSession creates a new session from the API entity
func NewSession(entity api.SessionEntity) *Session {
	return &Session{SessionEntity: entity}
}
