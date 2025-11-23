package auth

import (
	"context"

	"github.com/archesai/archesai/apps/studio/generated/core/models"
	"github.com/archesai/archesai/pkg/database"
)

// SessionRepository handles session persistence
type SessionRepository interface {
	database.CRUDRepository[models.Session]
}

// UserRepository handles user persistence
type UserRepository interface {
	database.CRUDRepository[models.User]

	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserBySessionID(ctx context.Context, sessionID string) (*models.User, error)
}

// AccountRepository handles account persistence
type AccountRepository interface {
	database.CRUDRepository[models.Account]

	GetAccountByProvider(
		ctx context.Context,
		provider, accountIdentifier string,
	) (*models.Account, error)
	ListAccountsByUserID(ctx context.Context, userID string) ([]*models.Account, error)
}
