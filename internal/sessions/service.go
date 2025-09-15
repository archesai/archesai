package sessions

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

// Service provides session business logic
type Service struct {
	repo           Repository
	cache          Cache
	sessionManager *SessionManager
	logger         *slog.Logger
}

// NewService creates a new session service
func NewService(repo Repository, cache Cache, logger *slog.Logger) *Service {
	sessionManager := NewSessionManager(repo, cache, 0) // Use default TTL
	return &Service{
		repo:           repo,
		cache:          cache,
		sessionManager: sessionManager,
		logger:         logger,
	}
}

// CreateSession creates a new user session
func (s *Service) CreateSession(ctx context.Context, userID uuid.UUID, organizationID uuid.UUID, ipAddress, userAgent string, _ bool) (*Session, string, error) {
	session, err := s.sessionManager.Create(ctx, userID, organizationID, ipAddress, userAgent)
	if err != nil {
		return nil, "", err
	}
	return session, session.Token, nil
}

// DeleteSession deletes a session
func (s *Service) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	return s.sessionManager.Delete(ctx, sessionID)
}

// FindSessions finds sessions for a user with pagination
func (s *Service) FindSessions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Session, int64, error) {
	userIDStr := userID.String()
	params := ListSessionsParams{
		UserID: &userIDStr,
		Limit:  limit,
		Offset: offset,
	}
	return s.repo.List(ctx, params)
}

// FindSessionByID finds a session by ID
func (s *Service) FindSessionByID(ctx context.Context, sessionID uuid.UUID) (*Session, error) {
	return s.sessionManager.Get(ctx, sessionID)
}
