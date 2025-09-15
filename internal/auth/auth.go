// Package auth provides authentication and authorization functionality.
// It includes user management, session handling, JWT token generation,
// and middleware for protecting routes.
package auth

import (
	"log/slog"
	"time"

	"github.com/archesai/archesai/internal/database/postgresql"
	"github.com/archesai/archesai/internal/email"
	"github.com/archesai/archesai/internal/users"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest represents a registration request
type RegisterRequest = RegisterJSONBody

// LoginRequest represents a login request
type LoginRequest = LoginJSONBody

// Service handles authentication operations
type Service struct {
	repo           Repository
	usersRepo      users.Repository
	cache          Cache
	sessionManager *SessionManager
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

// NewService creates a new authentication service
func NewService(repo Repository, usersRepo users.Repository, config Config, logger *slog.Logger) *Service {
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

	// Note: API Key Service must be set separately using SetAPIKeyService
	// as it requires its own repository
	return &Service{
		repo:      repo,
		usersRepo: usersRepo,
		jwtSecret: []byte(config.JWTSecret),
		logger:    logger,
		config:    config,
	}
}

// NewServiceWithCache creates a new auth service with Redis cache support
func NewServiceWithCache(repo Repository, usersRepo users.Repository, cache Cache, config Config, logger *slog.Logger) *Service {
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

	// Create session manager if cache is provided
	var sessionManager *SessionManager
	if cache != nil {
		sessionManager = NewSessionManager(repo, cache, config.SessionTokenExpiry)
	}

	// Note: API Key Service must be set separately using SetAPIKeyService
	// as it requires its own repository
	return &Service{
		repo:           repo,
		usersRepo:      usersRepo,
		cache:          cache,
		sessionManager: sessionManager,
		jwtSecret:      []byte(config.JWTSecret),
		logger:         logger,
		config:         config,
	}
}
