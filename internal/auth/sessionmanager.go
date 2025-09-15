package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// SessionManager handles session operations with Redis caching
type SessionManager struct {
	repo  Repository
	cache Cache
	ttl   time.Duration
}

// NewSessionManager creates a new session manager
func NewSessionManager(repo Repository, cache Cache, ttl time.Duration) *SessionManager {
	if ttl == 0 {
		ttl = 30 * 24 * time.Hour // 30 days default
	}
	return &SessionManager{
		repo:  repo,
		cache: cache,
		ttl:   ttl,
	}
}

// CreateSession creates a new session and stores it in both database and Redis
func (sm *SessionManager) CreateSession(ctx context.Context, userID, orgID uuid.UUID, ipAddress, userAgent string) (*Session, error) {
	// Generate secure session token
	token, err := generateSecureToken()
	if err != nil {
		return nil, fmt.Errorf("generate session token: %w", err)
	}

	// Create session entity
	session := &Session{
		Id:                   uuid.New(),
		UserId:               userID,
		Token:                token,
		ActiveOrganizationId: orgID,
		ExpiresAt:            time.Now().Add(sm.ttl).Format(time.RFC3339),
		IpAddress:            ipAddress,
		UserAgent:            userAgent,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	// Store in database first
	created, err := sm.repo.CreateSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("create session in db: %w", err)
	}

	// Store in Redis cache with TTL
	if sm.cache != nil {
		// Store by ID
		_ = sm.cache.SetSession(ctx, created, sm.ttl)

		// Store by token for quick lookup
		_ = sm.cache.SetSessionByToken(ctx, created.Token, created, sm.ttl)

		// Store user session index for listing
		userSessionKey := fmt.Sprintf("user:%s:session:%s", userID.String(), created.Id.String())
		_ = sm.storeUserSessionIndex(ctx, userSessionKey, created.Id, sm.ttl)
	}

	return created, nil
}

// GetSession retrieves a session by ID, checking cache first
func (sm *SessionManager) GetSession(ctx context.Context, sessionID uuid.UUID) (*Session, error) {
	// Try cache first
	if sm.cache != nil {
		cached, err := sm.cache.GetSession(ctx, sessionID)
		if err == nil && cached != nil {
			// Validate expiry
			if !sm.isSessionExpired(cached) {
				return cached, nil
			}
			// If expired, delete from cache
			_ = sm.cache.DeleteSession(ctx, sessionID)
		}
	}

	// Fallback to database
	session, err := sm.repo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Validate expiry
	if sm.isSessionExpired(session) {
		// Clean up expired session
		_ = sm.repo.DeleteSession(ctx, sessionID)
		return nil, ErrSessionExpired
	}

	// Update cache
	if sm.cache != nil && session != nil {
		_ = sm.cache.SetSession(ctx, session, sm.ttl)
	}

	return session, nil
}

// GetSessionByToken retrieves a session by token, checking cache first
func (sm *SessionManager) GetSessionByToken(ctx context.Context, token string) (*Session, error) {
	// Try cache first
	if sm.cache != nil {
		cached, err := sm.cache.GetSessionByToken(ctx, token)
		if err == nil && cached != nil {
			// Validate expiry
			if !sm.isSessionExpired(cached) {
				return cached, nil
			}
			// If expired, delete from cache
			_ = sm.cache.DeleteSessionByToken(ctx, token)
			_ = sm.cache.DeleteSession(ctx, cached.Id)
		}
	}

	// Fallback to database
	session, err := sm.repo.GetSessionByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Validate expiry
	if sm.isSessionExpired(session) {
		// Clean up expired session
		_ = sm.repo.DeleteSession(ctx, session.Id)
		return nil, ErrSessionExpired
	}

	// Update cache
	if sm.cache != nil && session != nil {
		_ = sm.cache.SetSessionByToken(ctx, token, session, sm.ttl)
	}

	return session, nil
}

// UpdateSession updates session metadata (like last activity)
func (sm *SessionManager) UpdateSession(ctx context.Context, sessionID uuid.UUID, updates *Session) (*Session, error) {
	// Update in database
	updated, err := sm.repo.UpdateSession(ctx, sessionID, updates)
	if err != nil {
		return nil, err
	}

	// Update cache
	if sm.cache != nil && updated != nil {
		_ = sm.cache.SetSession(ctx, updated, sm.ttl)
		if updated.Token != "" {
			_ = sm.cache.SetSessionByToken(ctx, updated.Token, updated, sm.ttl)
		}
	}

	return updated, nil
}

// DeleteSession removes a session from both database and cache
func (sm *SessionManager) DeleteSession(ctx context.Context, sessionID uuid.UUID) error {
	// Get session first to get the token
	session, _ := sm.GetSession(ctx, sessionID)

	// Delete from database
	if err := sm.repo.DeleteSession(ctx, sessionID); err != nil {
		return err
	}

	// Delete from cache
	if sm.cache != nil {
		_ = sm.cache.DeleteSession(ctx, sessionID)
		if session != nil && session.Token != "" {
			_ = sm.cache.DeleteSessionByToken(ctx, session.Token)
		}
		if session != nil {
			// Remove from user session index
			userSessionKey := fmt.Sprintf("user:%s:session:%s", session.UserId.String(), sessionID.String())
			_ = sm.removeUserSessionIndex(ctx, userSessionKey)
		}
	}

	return nil
}

// DeleteSessionByToken removes a session by token
func (sm *SessionManager) DeleteSessionByToken(ctx context.Context, token string) error {
	// Get session first
	session, err := sm.GetSessionByToken(ctx, token)
	if err != nil {
		return err
	}

	// Delete by ID
	return sm.DeleteSession(ctx, session.Id)
}

// DeleteUserSessions removes all sessions for a user
func (sm *SessionManager) DeleteUserSessions(ctx context.Context, userID uuid.UUID) error {
	// Delete from database
	if err := sm.repo.DeleteUserSessions(ctx, userID); err != nil {
		return err
	}

	// Delete from cache
	if sm.cache != nil {
		_ = sm.cache.DeleteUserSessions(ctx, userID)
	}

	return nil
}

// ListUserSessions returns all active sessions for a user
func (sm *SessionManager) ListUserSessions(ctx context.Context, userID uuid.UUID) ([]*Session, error) {
	// For now, use database directly
	// In a future enhancement, we could maintain a session index in Redis
	userIDStr := userID.String()
	params := ListSessionsParams{
		UserID: &userIDStr,
		Limit:  100,
	}
	sessions, _, err := sm.repo.ListSessions(ctx, params)
	if err != nil {
		return nil, err
	}

	// Filter out expired sessions
	var activeSessions []*Session
	for _, session := range sessions {
		if !sm.isSessionExpired(session) {
			activeSessions = append(activeSessions, session)
		}
	}

	return activeSessions, nil
}

// RefreshSession extends the expiry of a session
func (sm *SessionManager) RefreshSession(ctx context.Context, sessionID uuid.UUID) (*Session, error) {
	session, err := sm.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Update expiry
	newExpiry := time.Now().Add(sm.ttl)
	session.ExpiresAt = newExpiry.Format(time.RFC3339)
	session.UpdatedAt = time.Now()

	// Update in database
	updated, err := sm.repo.UpdateSession(ctx, sessionID, session)
	if err != nil {
		return nil, err
	}

	// Update cache with new TTL
	if sm.cache != nil && updated != nil {
		_ = sm.cache.SetSession(ctx, updated, sm.ttl)
		if updated.Token != "" {
			_ = sm.cache.SetSessionByToken(ctx, updated.Token, updated, sm.ttl)
		}
	}

	return updated, nil
}

// CleanupExpiredSessions removes all expired sessions
func (sm *SessionManager) CleanupExpiredSessions(ctx context.Context) error {
	return sm.repo.DeleteExpiredSessions(ctx)
}

// ValidateSession checks if a session is valid and not expired
func (sm *SessionManager) ValidateSession(ctx context.Context, token string) (*Session, error) {
	session, err := sm.GetSessionByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if sm.isSessionExpired(session) {
		_ = sm.DeleteSession(ctx, session.Id)
		return nil, ErrSessionExpired
	}

	// Update last activity by updating the UpdatedAt field
	session.UpdatedAt = time.Now()
	updated, err := sm.UpdateSession(ctx, session.Id, session)
	if err != nil {
		// Log error but don't fail validation
		return session, nil
	}

	return updated, nil
}

// Helper methods

func (sm *SessionManager) isSessionExpired(session *Session) bool {
	if session == nil {
		return true
	}

	expiresAt, err := time.Parse(time.RFC3339, session.ExpiresAt)
	if err != nil {
		return true
	}

	return time.Now().After(expiresAt)
}

func (sm *SessionManager) storeUserSessionIndex(_ context.Context, _ string, _ uuid.UUID, _ time.Duration) error {
	// This would need a custom Redis implementation to maintain a set of session IDs per user
	// For now, we'll rely on the database for listing
	return nil
}

func (sm *SessionManager) removeUserSessionIndex(_ context.Context, _ string) error {
	// This would need a custom Redis implementation
	return nil
}

// generateSecureToken generates a cryptographically secure random token
func generateSecureToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
