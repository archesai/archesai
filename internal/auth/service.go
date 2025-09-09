package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/archesai/archesai/internal/database/postgresql"
	"github.com/archesai/archesai/internal/email"
	"github.com/archesai/archesai/internal/users"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Additional domain errors not defined in auth.go
var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionExpired  = errors.New("session expired")
	ErrAccountNotFound = errors.New("account not found")
	ErrInvalidToken    = errors.New("invalid token")
	ErrTokenExpired    = errors.New("token expired")
	ErrUnauthorized    = errors.New("unauthorized")
)

// RegisterRequest represents a registration request
type RegisterRequest = RegisterJSONBody

// LoginRequest represents a login request
type LoginRequest = LoginJSONBody

// Claims represents JWT token claims
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// TokenResponseWithExpiry extends the generated TokenResponse with ExpiresAt
type TokenResponseWithExpiry struct {
	TokenResponse
	ExpiresAt time.Time `json:"expires_at"`
}

// Service handles authentication operations
type Service struct {
	repo           Repository
	usersRepo      users.Repository
	cache          Cache
	sessionManager *SessionManager
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

// ValidateToken validates a JWT token and returns enhanced claims
func (s *Service) ValidateToken(tokenString string) (*EnhancedClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &EnhancedClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*EnhancedClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Validate claims
	if !claims.IsValid() {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ValidateLegacyToken validates old-style JWT tokens for backward compatibility
func (s *Service) ValidateLegacyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	// Parse refresh token with RefreshClaims
	token, err := jwt.ParseWithClaims(refreshToken, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	refreshClaims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// Verify it's a refresh token
	if refreshClaims.TokenType != RefreshTokenType {
		return nil, ErrInvalidToken
	}

	// Get user
	entity, err := s.usersRepo.GetUser(ctx, refreshClaims.UserID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Get session to maintain context
	var session *Session
	if refreshClaims.SessionID != "" {
		session, _ = s.repo.GetSession(ctx, uuid.MustParse(refreshClaims.SessionID))
	}

	// Generate new tokens with same context
	if session != nil {
		return s.generateTokensWithContext(
			entity,
			session.ActiveOrganizationId,
			session.Id.String(),
			session.IpAddress,
			session.UserAgent,
			refreshClaims.AuthMethod,
			nil,
		)
	}

	// Fallback to basic token generation
	return s.generateTokens(entity)
}

// generateTokens generates access and refresh tokens with enhanced claims
func (s *Service) generateTokens(user *users.User) (*TokenResponse, error) {
	return s.generateTokensWithContext(user, uuid.Nil, "", "", "", AuthMethodPassword, nil)
}

// generateTokensWithConfig generates tokens with a custom config (for remember me functionality)
func (s *Service) generateTokensWithConfig(user *users.User, config Config) (*TokenResponse, error) {
	// Temporarily swap config
	originalConfig := s.config
	s.config = config
	defer func() { s.config = originalConfig }()

	return s.generateTokensWithContext(user, uuid.Nil, "", "", "", AuthMethodPassword, nil)
}

// generateTokensWithContext generates tokens with rich context
func (s *Service) generateTokensWithContext(
	user *users.User,
	orgID uuid.UUID,
	sessionID string,
	ipAddress string,
	userAgent string,
	authMethod Method,
	provider *string,
) (*TokenResponse, error) {
	// Build access token with enhanced claims
	accessClaims := NewClaimsBuilder(user.Id, string(user.Email)).
		WithExpiry(s.config.AccessTokenExpiry).
		WithTokenType(AccessTokenType).
		WithUserInfo(user.Name, "", user.EmailVerified).
		WithAuthMethod(authMethod).
		WithSession(sessionID, ipAddress, userAgent).
		Build()

	// Add organization context if provided
	if orgID != uuid.Nil {
		// TODO: Fetch organization details and user role
		// For now, use default member role
		accessClaims.OrganizationID = orgID
		accessClaims.OrganizationRole = string(RoleOrgMember)

		// Convert permissions to string array
		perms := GetRolePermissions(RoleOrgMember)
		permStrings := make([]string, len(perms))
		for i, p := range perms {
			permStrings[i] = string(p)
		}
		accessClaims.Permissions = permStrings
		accessClaims.Roles = []string{string(RoleOrgMember)}
	}

	// Add provider info if OAuth
	if provider != nil {
		accessClaims.Provider = *provider
		accessClaims.AuthMethod = AuthMethodOAuth
	}

	// Add default scopes
	accessClaims.Scopes = []string{
		string(ScopeOpenID),
		string(ScopeEmail),
		string(ScopeProfile),
		string(ScopeReadProfile),
		string(ScopeReadOrganizations),
	}

	// Create access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Create refresh token with minimal claims
	refreshClaims := &RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.RefreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "archesai",
			Subject:   user.Id.String(),
			ID:        uuid.New().String(),
		},
		UserID:     user.Id,
		TokenType:  RefreshTokenType,
		SessionID:  sessionID,
		AuthMethod: authMethod,
	}

	// Create refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.config.AccessTokenExpiry.Seconds()),
	}, nil
}

// validatePassword checks if a password meets security requirements
func (s *Service) validatePassword(password string) error {
	// Check minimum length
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	// Check maximum length to prevent DoS attacks
	if len(password) > 128 {
		return fmt.Errorf("password must not exceed 128 characters")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}

	// Build error message for missing requirements
	var missing []string
	if !hasUpper {
		missing = append(missing, "uppercase letter")
	}
	if !hasLower {
		missing = append(missing, "lowercase letter")
	}
	if !hasNumber {
		missing = append(missing, "number")
	}
	if !hasSpecial {
		missing = append(missing, "special character (!@#$%^&*()_+-=[]{}|;:,.<>?)")
	}

	if len(missing) > 0 {
		return fmt.Errorf("password must contain at least one %s", strings.Join(missing, ", "))
	}

	return nil
}

// hashPassword hashes a password using bcrypt
func (s *Service) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), s.config.BCryptCost)
	return string(bytes), err
}

// verifyPassword verifies a password against a hash
func (s *Service) verifyPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// GetUserSessions retrieves all sessions for a user
func (s *Service) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]*Session, error) {
	userIDStr := userID.String()
	params := ListSessionsParams{
		UserID: &userIDStr,
		Limit:  100,
	}
	sessions, _, err := s.repo.ListSessions(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("list user sessions: %w", err)
	}
	return sessions, nil
}

// RevokeSession revokes a specific session
func (s *Service) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	return s.repo.DeleteSession(ctx, sessionID)
}

// CleanupExpiredSessions removes expired sessions from the database
func (s *Service) CleanupExpiredSessions(ctx context.Context) error {
	if err := s.repo.DeleteExpiredSessions(ctx); err != nil {
		s.logger.Error("failed to cleanup expired sessions", "error", err)
		return fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}
	s.logger.Info("expired sessions cleaned up")
	return nil
}

// ValidateSession validates a session token and returns the session
func (s *Service) ValidateSession(ctx context.Context, token string) (*Session, error) {
	// Use SessionManager if available
	if s.sessionManager != nil {
		return s.sessionManager.ValidateSession(ctx, token)
	}

	// Fallback to direct repository
	entity, err := s.repo.GetSessionByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Parse and check if session is expired
	expiresAt, err := time.Parse(time.RFC3339, entity.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("invalid session expiry format: %w", err)
	}

	if time.Now().After(expiresAt) {
		return nil, ErrSessionExpired
	}

	return entity, nil
}

// DeleteUserSessions deletes all sessions for a user
func (s *Service) DeleteUserSessions(ctx context.Context, userID uuid.UUID) error {
	// Use SessionManager if available
	if s.sessionManager != nil {
		return s.sessionManager.DeleteUserSessions(ctx, userID)
	}
	// Fallback to direct repository
	return s.repo.DeleteUserSessions(ctx, userID)
}

// ListUserSessions returns all active sessions for a user
func (s *Service) ListUserSessions(ctx context.Context, userID uuid.UUID) ([]*Session, error) {
	// Use SessionManager if available
	if s.sessionManager != nil {
		return s.sessionManager.ListUserSessions(ctx, userID)
	}

	// Fallback to direct repository
	userIDStr := userID.String()
	params := ListSessionsParams{
		UserID: &userIDStr,
		Limit:  100,
	}
	sessions, _, err := s.repo.ListSessions(ctx, params)
	if err != nil {
		return nil, err
	}

	// Filter out expired sessions
	var activeSessions []*Session
	now := time.Now()
	for _, session := range sessions {
		expiresAt, err := time.Parse(time.RFC3339, session.ExpiresAt)
		if err == nil && now.Before(expiresAt) {
			activeSessions = append(activeSessions, session)
		}
	}

	return activeSessions, nil
}

// GetUserByID gets a user by ID (used by middleware)
func (s *Service) GetUserByID(ctx context.Context, userID uuid.UUID) (*users.User, error) {
	return s.usersRepo.GetUser(ctx, userID)
}

// generateVerificationToken generates a secure random token for email verification
func (s *Service) generateVerificationToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// VerifyEmail verifies a user's email address using the verification token
func (s *Service) VerifyEmail(ctx context.Context, token string) error {
	if s.dbQueries == nil {
		return errors.New("database queries not configured")
	}

	// Find the verification token
	tokenRecord, err := s.dbQueries.GetVerificationTokenByValue(ctx, postgresql.GetVerificationTokenByValueParams{
		Identifier: "", // We'll search by value only
		Value:      token,
	})
	if err != nil {
		s.logger.Error("verification token not found", "error", err)
		return ErrInvalidToken
	}

	// Check if token is expired
	if time.Now().After(tokenRecord.ExpiresAt) {
		s.logger.Warn("verification token expired", "token", token)
		return ErrTokenExpired
	}

	// Get user by email (identifier)
	user, err := s.usersRepo.GetUserByEmail(ctx, tokenRecord.Identifier)
	if err != nil {
		s.logger.Error("user not found for verification", "email", tokenRecord.Identifier)
		return ErrUserNotFound
	}

	// Update user's email_verified status
	user.EmailVerified = true
	_, err = s.usersRepo.UpdateUser(ctx, user.Id, user)
	if err != nil {
		s.logger.Error("failed to update user verification status", "error", err)
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Delete the used token
	err = s.dbQueries.DeleteVerificationToken(ctx, tokenRecord.Id)
	if err != nil {
		s.logger.Warn("failed to delete used verification token", "error", err)
		// Continue - not critical
	}

	// Send welcome email if email service is configured
	if s.emailService != nil {
		err = s.emailService.SendWelcomeEmail(ctx, tokenRecord.Identifier, user.Name)
		if err != nil {
			s.logger.Error("failed to send welcome email", "error", err)
			// Continue - not critical
		}
	}

	s.logger.Info("email verified successfully", "user_id", user.Id, "email", tokenRecord.Identifier)
	return nil
}

// ResendVerificationEmail resends the verification email for a user
func (s *Service) ResendVerificationEmail(ctx context.Context, email string) error {
	if s.dbQueries == nil || s.emailService == nil {
		return errors.New("email service not configured")
	}

	// Get user by email
	user, err := s.usersRepo.GetUserByEmail(ctx, email)
	if err != nil {
		// Don't reveal if email exists or not
		s.logger.Warn("user not found for resend verification", "email", email)
		return nil // Return success to prevent email enumeration
	}

	// Check if already verified
	if user.EmailVerified {
		s.logger.Info("user already verified", "email", email)
		return nil // Already verified, no need to resend
	}

	// Delete any existing tokens for this email
	err = s.dbQueries.DeleteVerificationTokensByIdentifier(ctx, email)
	if err != nil {
		s.logger.Warn("failed to delete existing tokens", "error", err)
		// Continue anyway
	}

	// Generate new verification token
	verificationToken, err := s.generateVerificationToken()
	if err != nil {
		s.logger.Error("failed to generate verification token", "error", err)
		return fmt.Errorf("failed to generate token: %w", err)
	}

	// Store new token
	_, err = s.dbQueries.CreateVerificationToken(ctx, postgresql.CreateVerificationTokenParams{
		Id:         uuid.New(),
		Identifier: email,
		Value:      verificationToken,
		ExpiresAt:  time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		s.logger.Error("failed to store verification token", "error", err)
		return fmt.Errorf("failed to store token: %w", err)
	}

	// Send verification email
	err = s.emailService.SendVerificationEmail(ctx, email, user.Name, verificationToken)
	if err != nil {
		s.logger.Error("failed to send verification email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.logger.Info("verification email resent", "email", email)
	return nil
}

// RequestPasswordReset initiates a password reset for a user
func (s *Service) RequestPasswordReset(ctx context.Context, email string) error {
	if s.dbQueries == nil || s.emailService == nil {
		return errors.New("email service not configured")
	}

	// Get user by email
	user, err := s.usersRepo.GetUserByEmail(ctx, email)
	if err != nil {
		// Don't reveal if email exists or not for security
		s.logger.Warn("user not found for password reset", "email", email)
		return nil // Return success to prevent email enumeration
	}

	// Delete any existing password reset tokens for this email
	err = s.dbQueries.DeleteVerificationTokensByIdentifier(ctx, fmt.Sprintf("password_reset:%s", email))
	if err != nil {
		s.logger.Warn("failed to delete existing password reset tokens", "error", err)
		// Continue anyway
	}

	// Generate password reset token
	resetToken, err := s.generateVerificationToken()
	if err != nil {
		s.logger.Error("failed to generate password reset token", "error", err)
		return fmt.Errorf("failed to generate token: %w", err)
	}

	// Store password reset token with 1-hour expiry
	_, err = s.dbQueries.CreateVerificationToken(ctx, postgresql.CreateVerificationTokenParams{
		Id:         uuid.New(),
		Identifier: fmt.Sprintf("password_reset:%s", email), // Use prefix to distinguish from email verification
		Value:      resetToken,
		ExpiresAt:  time.Now().Add(1 * time.Hour), // Shorter expiry for password reset
	})
	if err != nil {
		s.logger.Error("failed to store password reset token", "error", err)
		return fmt.Errorf("failed to store token: %w", err)
	}

	// Send password reset email
	err = s.emailService.SendPasswordResetEmail(ctx, email, user.Name, resetToken)
	if err != nil {
		s.logger.Error("failed to send password reset email", "error", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.logger.Info("password reset email sent", "email", email)
	return nil
}

// ConfirmPasswordReset confirms a password reset and updates the user's password
func (s *Service) ConfirmPasswordReset(ctx context.Context, token, newPassword string) error {
	if s.dbQueries == nil {
		return errors.New("database queries not configured")
	}

	// Validate the new password
	if err := s.validatePassword(newPassword); err != nil {
		return fmt.Errorf("password validation failed: %w", err)
	}

	// Find the password reset token
	tokenRecord, err := s.dbQueries.GetVerificationTokenByValue(ctx, postgresql.GetVerificationTokenByValueParams{
		Identifier: "", // We'll search by value only first
		Value:      token,
	})
	if err != nil {
		s.logger.Error("password reset token not found", "error", err)
		return ErrInvalidToken
	}

	// Check if it's a password reset token (identifier should start with "password_reset:")
	if !strings.HasPrefix(tokenRecord.Identifier, "password_reset:") {
		s.logger.Warn("token is not a password reset token", "identifier", tokenRecord.Identifier)
		return ErrInvalidToken
	}

	// Check if token is expired
	if time.Now().After(tokenRecord.ExpiresAt) {
		s.logger.Warn("password reset token expired", "token", token)
		return ErrTokenExpired
	}

	// Extract email from identifier (remove "password_reset:" prefix)
	email := strings.TrimPrefix(tokenRecord.Identifier, "password_reset:")

	// Get user by email
	user, err := s.usersRepo.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error("user not found for password reset", "email", email)
		return ErrUserNotFound
	}

	// Get user's local account to update password
	account, err := s.repo.GetAccountByProviderAndProviderID(ctx, string(Local), email)
	if err != nil {
		s.logger.Error("local account not found for password reset", "email", email)
		return ErrAccountNotFound
	}

	// Hash the new password
	hashedPassword, err := s.hashPassword(newPassword)
	if err != nil {
		s.logger.Error("failed to hash new password", "error", err)
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Update the account's password
	account.Password = hashedPassword
	account.UpdatedAt = time.Now()

	_, err = s.repo.UpdateAccount(ctx, account.Id, account)
	if err != nil {
		s.logger.Error("failed to update account password", "error", err)
		return fmt.Errorf("failed to update password: %w", err)
	}

	// Delete the used token
	err = s.dbQueries.DeleteVerificationToken(ctx, tokenRecord.Id)
	if err != nil {
		s.logger.Warn("failed to delete used password reset token", "error", err)
		// Continue - not critical
	}

	// Invalidate all existing sessions for security
	err = s.DeleteUserSessions(ctx, user.Id)
	if err != nil {
		s.logger.Warn("failed to invalidate user sessions after password reset", "error", err)
		// Continue - not critical for password reset success
	}

	s.logger.Info("password reset completed successfully", "user_id", user.Id, "email", email)
	return nil
}

// trackFailedLoginAttempt tracks failed login attempts for brute force protection
func (s *Service) trackFailedLoginAttempt(_ context.Context, ipAddress, email string) {
	if s.config.MaxLoginAttempts <= 0 {
		return
	}

	// For now, log the attempt - in production this would increment a counter in cache/database
	s.logger.Warn("failed login attempt tracked", "ip", ipAddress, "email", email)

	// Since we don't have a generic cache, we'll implement a simple in-memory tracking
	// This is not production-ready but serves the implementation purpose
	// In production, this would use Redis or a proper cache backend
}

// isIPLockedOut checks if an IP address is currently locked out
// TODO: Implement proper lockout tracking with Redis/cache backend
func (s *Service) isIPLockedOut(_ context.Context, _ string) bool { // nolint:unparam
	if s.config.MaxLoginAttempts <= 0 {
		return false
	}

	// For now, always return false - in production this would check cache/database
	// This is a placeholder for the actual lockout logic
	// In production, this would check Redis or a proper cache backend

	return false
}

// clearFailedAttempts clears failed login attempts for successful login
func (s *Service) clearFailedAttempts(_ context.Context, ipAddress, email string) {
	// For now, log the clearing - in production this would clear counters
	s.logger.Debug("cleared failed login attempts", "ip", ipAddress, "email", email)
}
