package domain

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"golang.org/x/crypto/bcrypt"
)

// Service handles authentication operations
type Service struct {
	repo      Repository
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
func NewService(repo Repository, config Config, logger *slog.Logger) *Service {
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
	existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, nil, ErrUserExists
	}

	// Hash the password
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		s.logger.Error("failed to hash password", "error", err)
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create new user with embedded UserEntity
	now := time.Now()
	user := &User{
		UserEntity: UserEntity{
			Id:            uuid.New(),
			Email:         openapi_types.Email(req.Email),
			Name:          req.Name,
			EmailVerified: false,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		PasswordHash: hashedPassword,
	}

	// Save user to database
	if err := s.repo.CreateUser(ctx, user); err != nil {
		s.logger.Error("failed to create user", "error", err)
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
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
		SessionEntity: SessionEntity{
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
		},
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		s.logger.Error("failed to create session", "error", err)
		return nil, nil, fmt.Errorf("failed to create session: %w", err)
	}

	s.logger.Info("user signed up successfully", "user_id", user.Id.String())
	return user, tokens, nil
}

// Login authenticates a user
func (s *Service) Login(ctx context.Context, req *LoginRequest, ipAddress, userAgent string) (*User, *TokenResponse, error) {
	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Warn("user not found", "email", req.Email)
		return nil, nil, ErrInvalidCredentials
	}

	// Verify password
	if err := s.verifyPassword(req.Password, user.PasswordHash); err != nil {
		s.logger.Warn("invalid password", "user_id", user.Id.String())
		return nil, nil, ErrInvalidCredentials
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
		SessionEntity: SessionEntity{
			Id:                   uuid.New(),
			UserId:               user.Id.String(),
			Token:                tokens.RefreshToken,
			ExpiresAt:            sessionNow.Add(s.config.SessionTokenExpiry).Format(time.RFC3339),
			CreatedAt:            sessionNow,
			UpdatedAt:            sessionNow,
			ActiveOrganizationId: "",
			IpAddress:            ipAddress,
			UserAgent:            userAgent,
		},
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		s.logger.Error("failed to create session", "error", err)
		return nil, nil, fmt.Errorf("failed to create session: %w", err)
	}

	s.logger.Info("user signed in successfully", "user_id", user.Id.String())
	return user, tokens, nil
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
	user, err := s.repo.GetUserByID(ctx, userUUID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Generate new tokens
	tokens, err := s.generateTokens(user)
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
	return s.repo.GetUserByID(ctx, id)
}

// UpdateUser updates user information
func (s *Service) UpdateUser(ctx context.Context, id uuid.UUID, req *UpdateUserRequest) (*User, error) {
	// Get existing user
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Image != nil {
		user.Image = *req.Image
	}
	user.UpdatedAt = time.Now()

	// Save updated user
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes a user
func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteUser(ctx, id)
}

// ListUsers lists users with pagination
func (s *Service) ListUsers(ctx context.Context, limit, offset int32) ([]*User, error) {
	return s.repo.ListUsers(ctx, limit, offset)
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
	session, err := s.repo.GetSessionByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Parse and check if session is expired
	expiresAt, err := time.Parse(time.RFC3339, session.ExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("invalid session expiry format: %w", err)
	}

	if time.Now().After(expiresAt) {
		return nil, ErrSessionExpired
	}

	return session, nil
}

// DeleteUserSessions deletes all sessions for a user
func (s *Service) DeleteUserSessions(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DeleteUserSessions(ctx, userID)
}
