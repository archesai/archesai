// Package auth provides authentication and authorization functionality.
// It includes user management, session handling, JWT token generation,
// and middleware for protecting routes.
package auth

//go:generate go tool oapi-codegen --config=models.cfg.yaml ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=server.cfg.yaml ../../api/openapi.bundled.yaml

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

// ContextKey is a type for context keys
type ContextKey string

const (
	// UserContextKey is the context key for the authenticated user
	UserContextKey ContextKey = "user"
	// ClaimsContextKey is the context key for JWT claims
	ClaimsContextKey ContextKey = "claims"
)

// Domain errors
var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")
	// ErrInvalidCredentials is returned when credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrUserExists is returned when a user already exists
	ErrUserExists = errors.New("user already exists")
)

// AuthRepository combines the generated Repository with additional methods
type AuthRepository interface {
	Repository

	// Additional User methods
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)

	// Additional Session methods
	GetSessionByToken(ctx context.Context, token string) (*Session, error)
	DeleteExpiredSessions(ctx context.Context) error
	DeleteUserSessions(ctx context.Context, userID uuid.UUID) error

	// Additional Account methods
	GetAccountByProviderAndProviderID(ctx context.Context, provider, providerID string) (*Account, error)
	GetAccountsByUserID(ctx context.Context, userID uuid.UUID) ([]*Account, error)
}
