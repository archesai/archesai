package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// SessionsRepository handles session persistence
type SessionsRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, entity *SessionEntity) (*SessionEntity, error)
	Get(ctx context.Context, id uuid.UUID) (*SessionEntity, error)
	Update(ctx context.Context, id uuid.UUID, entity *SessionEntity) (*SessionEntity, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, params ListSessionsParams) ([]*SessionEntity, int64, error)

	// Additional operations
	GetByToken(ctx context.Context, token string) (*SessionEntity, error)
	DeleteByToken(ctx context.Context, token string) error
	DeleteByUser(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
}

// SessionEntity represents a session in the database
type SessionEntity struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	Token          string
	OrganizationID uuid.UUID
	AuthMethod     string
	AuthProvider   string
	IPAddress      string
	UserAgent      string
	ExpiresAt      time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ListSessionsParams represents parameters for listing sessions
type ListSessionsParams struct {
	Page   PageQuery
	UserID *uuid.UUID
	Sort   *string
}

// PageQuery represents pagination parameters
type PageQuery struct {
	Limit  *int
	Offset *int
}
