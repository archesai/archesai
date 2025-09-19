package sessions

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Service implements the business logic.
type Service struct {
	repo      Repository
	jwtSecret string
	logger    *slog.Logger
}

// NewService creates a new service implementation.
func NewService(repo Repository, jwtSecret string, logger *slog.Logger) *Service {
	return &Service{
		repo:      repo,
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

// Validate validates a session token and returns the session if valid.
func (s *Service) Validate(ctx context.Context, token string) (*Session, error) {
	// Get session by token
	session, err := s.repo.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		// Clean up expired session
		_ = s.repo.Delete(ctx, session.ID)
		return nil, ErrSessionExpired
	}

	// Update last activity timestamp
	session.UpdatedAt = time.Now()
	updated, err := s.repo.Update(ctx, session.ID, session)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

// ListByUser lists all sessions for a specific user.
func (s *Service) ListByUser(ctx context.Context, _ uuid.UUID) ([]*Session, error) {
	params := ListSessionsParams{
		Page: PageQuery{
			Number: 1,
			Size:   100,
		},
	}

	sessions, _, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

// RevokeSession revokes a session by ID.
func (s *Service) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	// Check if entity exists first
	_, err := s.repo.Get(ctx, sessionID)
	if err != nil {
		return err
	}

	// Delete the entity
	err = s.repo.Delete(ctx, sessionID)
	if err != nil {
		return err
	}

	return nil
}

// CleanupExpiredSessions removes all expired sessions from the repository.
func (s *Service) CleanupExpiredSessions(ctx context.Context) error {
	err := s.repo.DeleteExpired(ctx)
	if err != nil {
		s.logger.Error("Failed to clean up expired sessions", "error", err)
		return err
	}
	return nil
}

// GenerateAccessToken generates a JWT access token for a user
func (s *Service) GenerateAccessToken(userID uuid.UUID, sessionID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub":        userID.String(),
		"session_id": sessionID.String(),
		"type":       "access",
		"exp":        time.Now().Add(1 * time.Hour).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken generates a JWT refresh token for a user
func (s *Service) GenerateRefreshToken(userID uuid.UUID, sessionID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub":        userID.String(),
		"session_id": sessionID.String(),
		"type":       "refresh",
		"exp":        time.Now().Add(30 * 24 * time.Hour).Unix(),
		"iat":        time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return tokenString, nil
}

// CreateSessionWithMethod creates a new session with specific auth method
func (s *Service) CreateSessionWithMethod(
	ctx context.Context,
	userID uuid.UUID,
	organizationID *uuid.UUID,
	ipAddress string,
	userAgent string,
	authMethod string,
	authProvider string,
) (*Session, error) {
	// Generate session token
	token := uuid.New().String()

	// Handle nullable organizationID
	var orgID uuid.UUID
	if organizationID != nil {
		orgID = *organizationID
	} else {
		orgID = uuid.Nil
	}

	// Create session
	session := &Session{
		ID:             uuid.New(),
		UserID:         userID,
		OrganizationID: orgID,
		Token:          token,
		IPAddress:      ipAddress,
		UserAgent:      userAgent,
		ExpiresAt:      time.Now().Add(30 * 24 * time.Hour), // 30 days
		AuthMethod:     authMethod,
		AuthProvider:   authProvider,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Store in database
	created, err := s.repo.Create(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return created, nil
}
