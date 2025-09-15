package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// RevokeSession revokes a specific session
func (s *Service) RevokeSession(ctx context.Context, sessionID uuid.UUID) error {
	return s.repo.DeleteSession(ctx, sessionID)
}

// CleanupExpiredSessions removes all expired sessions
func (s *Service) CleanupExpiredSessions(ctx context.Context) error {
	// Use SessionManager if available
	if s.sessionManager != nil {
		return s.sessionManager.CleanupExpiredSessions(ctx)
	}

	// Fallback to direct repository
	return s.repo.DeleteExpiredSessions(ctx)
}

// ValidateSession validates a session token
func (s *Service) ValidateSession(ctx context.Context, token string) (*Session, error) {
	// Use SessionManager if available
	if s.sessionManager != nil {
		return s.sessionManager.ValidateSession(ctx, token)
	}

	// Fallback to direct repository
	session, err := s.repo.GetSessionByToken(ctx, token)
	if err != nil {
		return nil, ErrSessionNotFound
	}

	// Check if session is expired
	if session.ExpiresAt != "" {
		expiresAt, err := time.Parse(time.RFC3339, session.ExpiresAt)
		if err == nil && time.Now().After(expiresAt) {
			return nil, ErrSessionExpired
		}
	}

	return session, nil
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

// ListUserSessions lists all sessions for a user
func (s *Service) ListUserSessions(ctx context.Context, userID uuid.UUID) ([]*Session, error) {
	// Use SessionManager if available
	if s.sessionManager != nil {
		return s.sessionManager.ListUserSessions(ctx, userID)
	}

	// Fallback to direct repository
	userIDStr := userID.String()
	sessions, _, err := s.repo.ListSessions(ctx, ListSessionsParams{
		UserID: &userIDStr,
		Limit:  100,
	})
	if err != nil {
		s.logger.Error("failed to list user sessions", "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}

	// Filter out expired sessions
	var activeSessions []*Session
	now := time.Now()
	for _, session := range sessions {
		if session.ExpiresAt != "" {
			expiresAt, err := time.Parse(time.RFC3339, session.ExpiresAt)
			if err == nil && now.After(expiresAt) {
				// Skip expired session
				continue
			}
		}
		activeSessions = append(activeSessions, session)
	}

	return activeSessions, nil
}
