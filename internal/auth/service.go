package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

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

// UpdateUserRequest represents a user update request
type UpdateUserRequest = UpdateUserJSONBody

// Claims represents JWT token claims
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	jwt.RegisteredClaims
}

// TokenResponse represents authentication token response
type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int64     `json:"expires_in"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// Service handles authentication operations
type Service struct {
	repo      AuthRepository
	jwtSecret []byte
	logger    *slog.Logger
	config    Config
}

// Config holds authentication configuration
type Config struct {
	JWTSecret          string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	SessionTokenExpiry time.Duration
	BCryptCost         int
	MaxLoginAttempts   int
	LockoutDuration    time.Duration
}

// NewService creates a new authentication service
func NewService(repo AuthRepository, config Config, logger *slog.Logger) *Service {
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
		jwtSecret: []byte(config.JWTSecret),
		logger:    logger,
		config:    config,
	}
}

// Register creates a new user account
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*User, *TokenResponse, error) {
	// Check if user already exists
	existingUser, err := s.repo.GetUserByEmail(ctx, string(req.Email))
	if err == nil && existingUser != nil {
		return nil, nil, ErrUserExists
	}

	// Hash the password
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		s.logger.Error("failed to hash password", "error", err)
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create new user with embedded User
	now := time.Now()
	user := &User{
		Id:            uuid.New(),
		Email:         req.Email,
		Name:          req.Name,
		EmailVerified: false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Save user to database - repository expects User
	createdEntity, err := s.repo.CreateUser(ctx, user)
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
		_ = s.repo.DeleteUser(ctx, user.Id)
		return nil, nil, fmt.Errorf("failed to create account: %w", err)
	}

	// Generate tokens
	tokens, err := s.generateTokens(user)
	if err != nil {
		s.logger.Error("failed to generate tokens", "error", err)
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Create session
	sessionNow := time.Now()
	session := &Session{
		Id:        uuid.New(),
		UserId:    user.Id.String(),
		Token:     tokens.RefreshToken,
		ExpiresAt: sessionNow.Add(s.config.SessionTokenExpiry).Format(time.RFC3339),
		CreatedAt: sessionNow,
		UpdatedAt: sessionNow,
		// Required fields with empty defaults
		ActiveOrganizationId: "",
		IpAddress:            "",
		UserAgent:            "",
	}

	_, err = s.repo.CreateSession(ctx, session)
	if err != nil {
		s.logger.Error("failed to create session", "error", err)
		return nil, nil, fmt.Errorf("failed to create session: %w", err)
	}

	s.logger.Info("user signed up successfully", "user_id", user.Id.String())
	return user, tokens, nil
}

// Login authenticates a user
func (s *Service) Login(ctx context.Context, req *LoginRequest, ipAddress, userAgent string) (*User, *TokenResponse, error) {
	// Get user by email
	userEntity, err := s.repo.GetUserByEmail(ctx, string(req.Email))
	if err != nil {
		s.logger.Warn("user not found", "email", req.Email)
		return nil, nil, ErrInvalidCredentials
	}

	// Get the user's local account to verify password
	account, err := s.repo.GetAccountByProviderAndProviderID(ctx, string(Local), string(req.Email))
	if err != nil {
		s.logger.Warn("account not found", "email", req.Email)
		return nil, nil, ErrInvalidCredentials
	}

	// Verify password
	if account.Password != "" {
		if err := s.verifyPassword(req.Password, account.Password); err != nil {
			s.logger.Warn("invalid password", "user_id", userEntity.Id.String())
			return nil, nil, ErrInvalidCredentials
		}
	}

	user := userEntity

	// Generate tokens
	tokens, err := s.generateTokens(user)
	if err != nil {
		s.logger.Error("failed to generate tokens", "error", err)
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Create session
	sessionNow := time.Now()
	session := &Session{
		Id:                   uuid.New(),
		UserId:               user.Id.String(),
		Token:                tokens.RefreshToken,
		ExpiresAt:            sessionNow.Add(s.config.SessionTokenExpiry).Format(time.RFC3339),
		CreatedAt:            sessionNow,
		UpdatedAt:            sessionNow,
		ActiveOrganizationId: "",
		IpAddress:            ipAddress,
		UserAgent:            userAgent,
	}

	_, err = s.repo.CreateSession(ctx, session)
	if err != nil {
		s.logger.Error("failed to create session", "error", err)
		return nil, nil, fmt.Errorf("failed to create session: %w", err)
	}

	s.logger.Info("user signed in successfully", "user_id", userEntity.Id.String())
	return userEntity, tokens, nil
}

// Logout invalidates a user session
func (s *Service) Logout(ctx context.Context, token string) error {
	// Get session by token
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

// ValidateToken validates a JWT token
func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
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
	// Validate refresh token
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Get user (convert uuid.UUID to uuid.UUID)
	userUUID, _ := uuid.Parse(claims.UserID.String())
	entity, err := s.repo.GetUserByID(ctx, userUUID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Generate new tokens
	tokens, err := s.generateTokens(entity)
	if err != nil {
		s.logger.Error("failed to generate tokens", "error", err)
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return tokens, nil
}

// generateTokens generates access and refresh tokens
func (s *Service) generateTokens(user *User) (*TokenResponse, error) {
	now := time.Now()
	expiresAt := now.Add(s.config.AccessTokenExpiry)

	// Create access token claims
	accessClaims := &Claims{
		UserID: user.Id,
		Email:  string(user.Email),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "archesai",
			Subject:   user.Id.String(),
			ID:        uuid.New().String(),
		},
	}

	// Create access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Create refresh token claims
	refreshClaims := &Claims{
		UserID: user.Id,
		Email:  string(user.Email),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.config.RefreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "archesai",
			Subject:   user.Id.String(),
			ID:        uuid.New().String(),
		},
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
		ExpiresAt:    expiresAt,
	}, nil
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

// GetUser retrieves a user by ID
func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
	entity, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

// UpdateUser updates user information
func (s *Service) UpdateUser(ctx context.Context, id uuid.UUID, req *UpdateUserRequest) (*User, error) {
	// Get existing user
	entity, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Email != "" {
		entity.Email = Email(req.Email)
	}
	if req.Image != "" {
		entity.Image = req.Image
	}
	entity.UpdatedAt = time.Now()

	// Save updated user
	updatedEntity, err := s.repo.UpdateUser(ctx, id, entity)
	if err != nil {
		return nil, err
	}

	return updatedEntity, nil
}

// DeleteUser deletes a user
func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteUser(ctx, id)
}

// ListUsers lists users with pagination
func (s *Service) ListUsers(ctx context.Context, limit, offset int32) ([]*User, error) {
	params := ListUsersParams{
		Limit:  int(limit),
		Offset: int(offset),
	}
	entities, _, err := s.repo.ListUsers(ctx, params)
	if err != nil {
		return nil, err
	}

	// Convert entities to domain users
	users := make([]*User, len(entities))
	for i, entity := range entities {
		users[i] = entity
	}
	return users, nil
}

// GetUserSessions retrieves all sessions for a user
func (s *Service) GetUserSessions(_ context.Context, _ uuid.UUID) ([]*Session, error) {
	// TODO: Add ListSessionsByUser query to auth.sql
	// For now, return empty slice
	return []*Session{}, nil
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
	return s.repo.DeleteUserSessions(ctx, userID)
}
