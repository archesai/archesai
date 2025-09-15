package auth

import (
	"github.com/archesai/archesai/internal/accounts"
	"github.com/archesai/archesai/internal/sessions"
	"github.com/archesai/archesai/internal/users"
)

// AccountsRepository is an alias for accounts.Repository
type AccountsRepository interface {
	accounts.Repository
}

// SessionsRepository is an alias for sessions.Repository
type SessionsRepository interface {
	sessions.Repository
}

// SessionsCache is an alias for sessions.Cache
type SessionsCache interface {
	sessions.Cache
}

// UsersRepository is an alias for users.Repository
type UsersRepository interface {
	users.Repository
}
