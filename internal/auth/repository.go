package auth

import (
	"context"

	"github.com/google/uuid"
)

// User extends UserEntity with auth-specific fields
type User struct {
	UserEntity
	PasswordHash string `json:"-"` // Never expose password hash
}

// Session extends SessionEntity
type Session struct {
	SessionEntity
}

// Account extends AccountEntity
type Account struct {
	AccountEntity
}

// Repository defines the interface for auth data persistence
// This interface is defined in the domain package (hexagonal architecture)
type Repository interface {
	// User operations
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ListUsers(ctx context.Context, limit, offset int32) ([]*User, error)

	// Session operations
	CreateSession(ctx context.Context, session *Session) error
	GetSessionByToken(ctx context.Context, token string) (*Session, error)
	GetSessionByID(ctx context.Context, id uuid.UUID) (*Session, error)
	UpdateSession(ctx context.Context, session *Session) error
	DeleteSession(ctx context.Context, id uuid.UUID) error
	DeleteUserSessions(ctx context.Context, userID uuid.UUID) error
	DeleteExpiredSessions(ctx context.Context) error
}
