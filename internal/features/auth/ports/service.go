package ports

import (
	"context"

	"github.com/archesai/archesai/internal/features/auth/domain"
	"github.com/google/uuid"
)

// Service defines the interface for authentication business logic
type Service interface {
	// Authentication
	SignUp(ctx context.Context, req *domain.SignUpRequest) (*domain.User, *domain.TokenResponse, error)
	SignIn(ctx context.Context, req *domain.SignInRequest, ipAddress, userAgent string) (*domain.User, *domain.TokenResponse, error)
	SignOut(ctx context.Context, token string) error
	RefreshToken(ctx context.Context, refreshToken string) (*domain.TokenResponse, error)
	ValidateToken(tokenString string) (*domain.Claims, error)

	// User management
	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req *domain.UpdateUserRequest) (*domain.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ListUsers(ctx context.Context, limit, offset int32) ([]*domain.User, error)

	// Session management
	GetUserSessions(ctx context.Context, userID uuid.UUID) ([]*domain.Session, error)
	RevokeSession(ctx context.Context, sessionID uuid.UUID) error
	CleanupExpiredSessions(ctx context.Context) error
}