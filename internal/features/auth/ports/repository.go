package ports

import (
	"context"

	"github.com/archesai/archesai/internal/features/auth/domain"
	"github.com/google/uuid"
)

// Repository defines the interface for auth data persistence
type Repository interface {
	// User operations
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ListUsers(ctx context.Context, limit, offset int32) ([]*domain.User, error)

	// Session operations
	CreateSession(ctx context.Context, session *domain.Session) error
	GetSessionByToken(ctx context.Context, token string) (*domain.Session, error)
	GetSessionByID(ctx context.Context, id uuid.UUID) (*domain.Session, error)
	UpdateSession(ctx context.Context, session *domain.Session) error
	DeleteSession(ctx context.Context, id uuid.UUID) error
	DeleteUserSessions(ctx context.Context, userID uuid.UUID) error
	DeleteExpiredSessions(ctx context.Context) error
}