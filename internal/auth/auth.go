// Package auth provides authentication and authorization functionality.
// It includes user management, session handling, JWT token generation,
// and middleware for protecting routes.
package auth

import (
	"log/slog"
	"time"

	"github.com/archesai/archesai/internal/accounts"
	"github.com/archesai/archesai/internal/database/postgresql"
	"github.com/archesai/archesai/internal/email"
	"github.com/archesai/archesai/internal/sessions"
	"github.com/archesai/archesai/internal/users"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest represents a registration request
type RegisterRequest = RegisterJSONBody

// LoginRequest represents a login request
type LoginRequest = LoginJSONBody

// Service handles authentication operations
type Service struct {
	accountsRepo   accounts.Repository
	sessionsRepo   sessions.Repository
	usersRepo      users.Repository
	cache          sessions.Cache
	sessionManager *sessions.SessionManager
	apiKeyService  *APIKeyService // API key management
	jwtSecret      []byte
	logger         *slog.Logger
	config         Config
	dbQueries      *postgresql.Queries
	emailService   *email.Service
}

// Config holds authentication configuration
type Config struct {
	JWTSecret             string
	AccessTokenExpiry     time.Duration
	RefreshTokenExpiry    time.Duration
	SessionTokenExpiry    time.Duration
	BCryptCost            int
	MaxLoginAttempts      int
	LockoutDuration       time.Duration
	MaxConcurrentSessions int // Maximum concurrent sessions per user (0 = unlimited)
}

// NewService creates a new authentication service with cache support
// If cache is nil, a NoOpCache will be used as fallback
func NewService(accountsRepo accounts.Repository, sessionsRepo sessions.Repository, usersRepo users.Repository, cache sessions.Cache, config Config, logger *slog.Logger) *Service {
	// Use NoOpCache if no cache provided
	if cache == nil {
		cache = sessions.NewNoOpCache()
	}
	if config.AccessTokenExpiry == 0 {
		config.AccessTokenExpiry = 15 * time.Minute
	}
	if config.RefreshTokenExpiry == 0 {
		config.RefreshTokenExpiry = 7 * 24 * time.Hour
	}
	if config.SessionTokenExpiry == 0 {
		config.SessionTokenExpiry = 30 * 24 * time.Hour
	}
	if config.BCryptCost == 0 {
		config.BCryptCost = bcrypt.DefaultCost
	}

	// Always create session manager with the cache (might be NoOpCache)
	sessionManager := sessions.NewSessionManager(sessionsRepo, cache, config.SessionTokenExpiry)

	return &Service{
		accountsRepo:   accountsRepo,
		sessionsRepo:   sessionsRepo,
		usersRepo:      usersRepo,
		cache:          cache,
		sessionManager: sessionManager,
		jwtSecret:      []byte(config.JWTSecret),
		logger:         logger,
		config:         config,
	}
}
