package auth

import (
	"context"

	"github.com/archesai/archesai/apps/studio/generated/core"
	"github.com/archesai/archesai/pkg/database"
)

// MagicLinkToken is an alias for core.MagicLinkToken
type MagicLinkToken = core.MagicLinkToken

// Session is an alias for core.Session
type Session = core.Session

// User is an alias for core.User
type User = core.User

// Account is an alias for core.Account
type Account = core.Account

// SessionRepository is an alias for database.Repository[Session]
type SessionRepository = database.Repository[Session]

// SessionAuthProvider is an alias for core.SessionAuthProvider
type SessionAuthProvider = core.SessionAuthProvider

// ErrUserNotFound is an alias for core.ErrUserNotFound
var ErrUserNotFound = core.ErrUserNotFound

// AccountProviderLocal is an alias for core.AccountProviderLocal
var AccountProviderLocal = core.AccountProviderLocal

// NewUser creates a new User instance with the provided details
func NewUser(email string, _ bool, _ *string, _ string) (*User, error) {
	return &User{
		Email: email,
	}, nil
}

// UserRepository handles user persistence
type UserRepository interface {
	database.Repository[User]

	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserBySessionID(ctx context.Context, sessionID string) (*User, error)
}

// AccountRepository handles account persistence
type AccountRepository interface {
	database.Repository[Account]

	GetAccountByProvider(
		ctx context.Context,
		provider, accountIdentifier string,
	) (*Account, error)
	ListAccountsByUserID(ctx context.Context, userID string) ([]*Account, error)
}
