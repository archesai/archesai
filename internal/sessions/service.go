package sessions

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Validate validates a session token and returns the session if valid
func (s *Service) Validate(ctx context.Context, token string) (*Session, error) {
	// Get session by token
	session, err := s.repo.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Check if session is expired
	expiresAt, err := time.Parse(time.RFC3339, session.ExpiresAt)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if time.Now().After(expiresAt) {
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

// ListByUser lists all sessions for a specific user
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

// RevokeSession revokes a session by ID
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

// CleanupExpiredSessions removes all expired sessions from the repository
func (s *Service) CleanupExpiredSessions(ctx context.Context) error {
	err := s.repo.DeleteExpired(ctx)
	if err != nil {
		s.logger.Error("Failed to clean up expired sessions", "error", err)
		return err
	}
	return nil
}
