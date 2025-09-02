// Package repositories defines repository interfaces for the auth domain.
package repositories

import (
	"context"

	"github.com/archesai/archesai/internal/domains/auth/entities"
	"github.com/google/uuid"
)

// Repository defines the interface for auth data persistence
type Repository interface {
	// User operations
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	CreateUser(ctx context.Context, user *entities.User) error
	UpdateUser(ctx context.Context, user *entities.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ListUsers(ctx context.Context, limit, offset int32) ([]*entities.User, error)

	// Session operations
	CreateSession(ctx context.Context, session *entities.Session) error
	GetSessionByToken(ctx context.Context, token string) (*entities.Session, error)
	GetSessionByID(ctx context.Context, id uuid.UUID) (*entities.Session, error)
	UpdateSession(ctx context.Context, session *entities.Session) error
	DeleteSession(ctx context.Context, id uuid.UUID) error
	DeleteUserSessions(ctx context.Context, userID uuid.UUID) error
	DeleteExpiredSessions(ctx context.Context) error
}
