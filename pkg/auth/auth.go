package auth

import (
	"github.com/archesai/archesai/pkg/auth/models"
	"github.com/archesai/archesai/pkg/auth/repositories"
)

const (
	bindHost = "0.0.0.0"
)

// MagicLinkToken is an alias for the magic link token model.
type MagicLinkToken = models.MagicLinkToken

// User is an alias for the user model.
type User = models.User

// Session is an alias for the session model.
type Session = models.Session

// Account is an alias for the account model.
type Account = models.Account

// SessionAuthProvider is an alias for the session auth provider model.
type SessionAuthProvider = models.SessionAuthProvider

// SessionRepository is an alias for the session repository interface.
type SessionRepository = repositories.SessionRepository

// UserRepository is an alias for the user repository interface.
type UserRepository = repositories.UserRepository

// AccountRepository is an alias for the account repository interface.
type AccountRepository = repositories.AccountRepository

// NewUser creates a new user with the given parameters.
func NewUser(email string, emailVerified bool, image *string, name string) (*models.User, error) {
	return models.NewUser(email, emailVerified, image, name)
}

// ErrUserNotFound is returned when a user is not found.
var ErrUserNotFound = models.ErrUserNotFound

// AccountProviderLocal is the constant for local authentication provider.
const AccountProviderLocal = models.AccountProviderLocal

// Config holds the configuration for the authentication service.
type Config struct {
	API      *APIConfig
	Platform *PlatformConfig
	Auth     *ProvidersConfig
}

// PlatformConfig holds the configuration for the platform URL.
type PlatformConfig struct {
	URL *string
}

// APIConfig holds the configuration for the API server.
type APIConfig struct {
	Host string
	Port uint16
}

// ProvidersConfig holds the configuration for authentication providers.
type ProvidersConfig struct {
	Local     *LocalAuthConfig
	Google    *OAuthProviderConfig
	Github    *OAuthProviderConfig
	Microsoft *OAuthProviderConfig
}

// OAuthProviderConfig holds the configuration for an OAuth provider.
type OAuthProviderConfig struct {
	Enabled      bool
	ClientID     *string
	ClientSecret *string
	RedirectURL  *string
}

// LocalAuthConfig holds the configuration for local authentication.
type LocalAuthConfig struct {
	Enabled   bool
	JWTSecret string
}
