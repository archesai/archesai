// Package auth provides authentication and authorization functionality.
// It includes user management, session handling, JWT token generation,
// and middleware for protecting routes.
package auth

//go:generate go tool oapi-codegen --config=../../types.codegen.yaml --package auth --include-tags Auth,Sessions,Accounts ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../server.codegen.yaml --package auth --include-tags Auth,Sessions,Accounts ../../api/openapi.bundled.yaml

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/archesai/archesai/internal/database/postgresql"
	"github.com/archesai/archesai/internal/email"
	"github.com/archesai/archesai/internal/users"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ContextKey is a type for context keys
type ContextKey string

const (
	// UserContextKey is the context key for the authenticated user
	UserContextKey ContextKey = "user"
	// ClaimsContextKey is the context key for JWT claims
	ClaimsContextKey ContextKey = "claims"
	// SessionTokenContextKey is the context key for session token
	SessionTokenContextKey ContextKey = "session_token"
)

// Domain errors
var (
	// ErrUserNotFound is returned when a user is not found
	ErrUserNotFound = errors.New("user not found")
	// ErrInvalidCredentials is returned when credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrUserExists is returned when a user already exists
	ErrUserExists = errors.New("user already exists")
	// ErrInvalidPassword is returned for invalid passwords
	ErrInvalidPassword = errors.New("invalid password")
	// ErrSessionNotFound is returned when a session is not found
	ErrSessionNotFound = errors.New("session not found")
	// ErrSessionExpired is returned when a session has expired
	ErrSessionExpired = errors.New("session expired")
	// ErrAccountNotFound is returned when an account is not found
	ErrAccountNotFound = errors.New("account not found")
	// ErrInvalidToken is returned for invalid tokens
	ErrInvalidToken = errors.New("invalid token")
	// ErrTokenExpired is returned when a token has expired
	ErrTokenExpired = errors.New("token expired")
	// ErrUnauthorized is returned for unauthorized access
	ErrUnauthorized = errors.New("unauthorized")
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

// SetDatabaseQueries sets the database queries for the service
func (s *Service) SetDatabaseQueries(queries *postgresql.Queries) {
	s.dbQueries = queries
}

// SetEmailService sets the email service for the service
func (s *Service) SetEmailService(emailService *email.Service) {
	s.emailService = emailService
}

// SetAPIKeyService sets the API key service for authentication
func (s *Service) SetAPIKeyService(apiKeyService *APIKeyService) {
	s.apiKeyService = apiKeyService
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

// Register creates a new user account
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*users.User, *TokenResponse, error) {
	// Check if user already exists
	existingUser, err := s.usersRepo.GetUserByEmail(ctx, string(req.Email))
	if err == nil && existingUser != nil {
		return nil, nil, ErrUserExists
	}

	// Validate password strength
	if err := s.validatePassword(req.Password); err != nil {
		return nil, nil, fmt.Errorf("password validation failed: %w", err)
	}

	// Hash the password
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		s.logger.Error("failed to hash password", "error", err)
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create new user with embedded User
	now := time.Now()
	user := &users.User{
		Id:            uuid.New(),
		Email:         req.Email,
		Name:          req.Name,
		EmailVerified: false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Save user to database - repository expects User
	createdEntity, err := s.usersRepo.CreateUser(ctx, user)
	if err != nil {
		s.logger.Error("failed to create user", "error", err)
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}
	// Update user with created entity (in case DB added fields)
	user = createdEntity

	// Create local account with password
	account := &Account{
		Id:         uuid.New(),
		UserId:     user.Id,
		ProviderId: Local,
		AccountId:  string(user.Email), // Use email as account ID for local auth
		Password:   hashedPassword,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	_, err = s.repo.CreateAccount(ctx, account)
	if err != nil {
		s.logger.Error("failed to create account", "error", err)
		// Try to clean up the created user
		_ = s.usersRepo.DeleteUser(ctx, user.Id)
		return nil, nil, fmt.Errorf("failed to create account: %w", err)
	}

	// Generate email verification token if email service is configured
	if s.emailService != nil && s.dbQueries != nil {
		verificationToken, err := s.generateVerificationToken()
		if err != nil {
			s.logger.Error("failed to generate verification token", "error", err)
			// Continue without email verification
		} else {
			// Store verification token in database
			_, err = s.dbQueries.CreateVerificationToken(ctx, postgresql.CreateVerificationTokenParams{
				Id:         uuid.New(),
				Identifier: string(user.Email),
				Value:      verificationToken,
				ExpiresAt:  time.Now().Add(24 * time.Hour), // Token expires in 24 hours
			})
			if err != nil {
				s.logger.Error("failed to store verification token", "error", err)
				// Continue without email verification
			} else {
				// Send verification email
				err = s.emailService.SendVerificationEmail(ctx, string(user.Email), user.Name, verificationToken)
				if err != nil {
					s.logger.Error("failed to send verification email", "error", err)
					// Continue - user can request resend later
				}
			}
		}
	}

	// Generate tokens
	tokens, err := s.generateTokens(user)
	if err != nil {
		s.logger.Error("failed to generate tokens", "error", err)
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Create session - use SessionManager if available
	var session *Session
	if s.sessionManager != nil {
		session, err = s.sessionManager.CreateSession(ctx, user.Id, uuid.Nil, "", "")
		if err != nil {
			s.logger.Error("failed to create session", "error", err)
			return nil, nil, fmt.Errorf("failed to create session: %w", err)
		}
		// Update session with refresh token
		session.Token = tokens.RefreshToken
		_, err = s.sessionManager.UpdateSession(ctx, session.Id, session)
		if err != nil {
			s.logger.Error("failed to update session token", "error", err)
		}
	} else {
		// Fallback to direct repository
		sessionNow := time.Now()
		session = &Session{
			Id:        uuid.New(),
			UserId:    user.Id,
			Token:     tokens.RefreshToken,
			ExpiresAt: sessionNow.Add(s.config.SessionTokenExpiry).Format(time.RFC3339),
			CreatedAt: sessionNow,
			UpdatedAt: sessionNow,
			// Required fields with empty defaults
			ActiveOrganizationId: uuid.Nil,
			IpAddress:            "",
			UserAgent:            "",
		}
		_, err = s.repo.CreateSession(ctx, session)
		if err != nil {
			s.logger.Error("failed to create session", "error", err)
			return nil, nil, fmt.Errorf("failed to create session: %w", err)
		}
	}

	s.logger.Info("user signed up successfully", "user_id", user.Id.String())
	return user, tokens, nil
}

// Login authenticates a user
func (s *Service) Login(ctx context.Context, req *LoginRequest, ipAddress, userAgent string) (*users.User, *TokenResponse, error) {
	// Check if IP is locked out due to brute force attempts
	if s.config.MaxLoginAttempts > 0 {
		if s.isIPLockedOut(ctx, ipAddress) {
			s.logger.Warn("IP address locked out due to brute force attempts", "ip", ipAddress)
			return nil, nil, fmt.Errorf("too many failed login attempts, try again later")
		}
	}

	// Get user by email
	userEntity, err := s.usersRepo.GetUserByEmail(ctx, string(req.Email))
	if err != nil {
		// Track failed attempt
		s.trackFailedLoginAttempt(ctx, ipAddress, string(req.Email))
		s.logger.Warn("user not found", "email", req.Email)
		return nil, nil, ErrInvalidCredentials
	}

	// Get the user's local account to verify password
	account, err := s.repo.GetAccountByProviderAndProviderID(ctx, string(Local), string(req.Email))
	if err != nil {
		// Track failed attempt
		s.trackFailedLoginAttempt(ctx, ipAddress, string(req.Email))
		s.logger.Warn("account not found", "email", req.Email)
		return nil, nil, ErrInvalidCredentials
	}

	// Verify password
	if account.Password != "" {
		if err := s.verifyPassword(req.Password, account.Password); err != nil {
			// Track failed attempt
			s.trackFailedLoginAttempt(ctx, ipAddress, string(req.Email))
			s.logger.Warn("invalid password", "user_id", userEntity.Id.String(), "ip", ipAddress)
			return nil, nil, ErrInvalidCredentials
		}
	}

	user := userEntity

	// Clear any failed login attempts on successful authentication
	s.clearFailedAttempts(ctx, ipAddress, string(req.Email))

	// Check concurrent session limits
	if s.config.MaxConcurrentSessions > 0 {
		activeSessions, err := s.ListUserSessions(ctx, user.Id)
		if err == nil && len(activeSessions) >= s.config.MaxConcurrentSessions {
			// Remove oldest session if limit reached
			s.logger.Info("concurrent session limit reached, removing oldest session",
				"user_id", user.Id,
				"max_sessions", s.config.MaxConcurrentSessions,
				"active_sessions", len(activeSessions))

			// Find and remove the oldest session
			if len(activeSessions) > 0 {
				oldestSession := activeSessions[0]
				for _, session := range activeSessions {
					if session.CreatedAt.Before(oldestSession.CreatedAt) {
						oldestSession = session
					}
				}
				_ = s.repo.DeleteSession(ctx, oldestSession.Id)
			}
		}
	}

	// Generate tokens with extended refresh token if remember me is enabled
	var tokens *TokenResponse
	if req.RememberMe {
		// Use extended refresh token expiry for remember me
		extendedConfig := s.config
		extendedConfig.RefreshTokenExpiry = 30 * 24 * time.Hour // 30 days for remember me
		tokens, err = s.generateTokensWithConfig(user, extendedConfig)
		s.logger.Info("generating extended session for remember me", "user_id", user.Id)
	} else {
		tokens, err = s.generateTokens(user)
	}
	if err != nil {
		s.logger.Error("failed to generate tokens", "error", err)
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Create session - use SessionManager if available
	var session *Session
	if s.sessionManager != nil {
		// TODO: Get organization ID from user's default org
		session, err = s.sessionManager.CreateSession(ctx, user.Id, uuid.Nil, ipAddress, userAgent)
		if err != nil {
			s.logger.Error("failed to create session", "error", err)
			return nil, nil, fmt.Errorf("failed to create session: %w", err)
		}
		// Update session with refresh token
		session.Token = tokens.RefreshToken
		_, err = s.sessionManager.UpdateSession(ctx, session.Id, session)
		if err != nil {
			s.logger.Error("failed to update session token", "error", err)
		}
	} else {
		// Fallback to direct repository
		sessionNow := time.Now()
		session = &Session{
			Id:                   uuid.New(),
			UserId:               user.Id,
			Token:                tokens.RefreshToken,
			ExpiresAt:            sessionNow.Add(s.config.SessionTokenExpiry).Format(time.RFC3339),
			CreatedAt:            sessionNow,
			UpdatedAt:            sessionNow,
			ActiveOrganizationId: uuid.Nil, // TODO: Set proper organization ID
			IpAddress:            ipAddress,
			UserAgent:            userAgent,
		}
		_, err = s.repo.CreateSession(ctx, session)
		if err != nil {
			s.logger.Error("failed to create session", "error", err)
			return nil, nil, fmt.Errorf("failed to create session: %w", err)
		}
	}

	s.logger.Info("user signed in successfully", "user_id", userEntity.Id.String())
	return userEntity, tokens, nil
}

// Logout invalidates a user session
func (s *Service) Logout(ctx context.Context, token string) error {
	// Use SessionManager if available
	if s.sessionManager != nil {
		err := s.sessionManager.DeleteSessionByToken(ctx, token)
		if err != nil {
			s.logger.Error("failed to delete session", "error", err)
			return ErrInvalidToken
		}
		s.logger.Info("user signed out successfully")
		return nil
	}

	// Fallback to direct repository
	session, err := s.repo.GetSessionByToken(ctx, token)
	if err != nil {
		return ErrInvalidToken
	}

	// Delete session
	if err := s.repo.DeleteSession(ctx, session.Id); err != nil {
		s.logger.Error("failed to delete session", "error", err)
		return fmt.Errorf("failed to delete session: %w", err)
	}

	s.logger.Info("user signed out successfully", "user_id", session.UserId)
	return nil
}

// GetUserByID retrieves a user by their ID
func (s *Service) GetUserByID(ctx context.Context, userID uuid.UUID) (*users.User, error) {
	return s.usersRepo.GetUser(ctx, userID)
}
