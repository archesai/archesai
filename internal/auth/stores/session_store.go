package stores

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/archesai/archesai/internal/auth"
	genericcache "github.com/archesai/archesai/internal/cache"
)

// SessionStore handles session operations with caching.
type SessionStore struct {
	repo  auth.SessionsRepository
	cache genericcache.Cache[auth.SessionEntity]
	ttl   time.Duration
}

// NewSessionStore creates a new session store.
func NewSessionStore(
	repo auth.SessionsRepository,
	cache genericcache.Cache[auth.SessionEntity],
	ttl time.Duration,
) *SessionStore {
	if ttl == 0 {
		ttl = 30 * 24 * time.Hour // 30 days default
	}
	return &SessionStore{
		repo:  repo,
		cache: cache,
		ttl:   ttl,
	}
}

// Create creates a new session.
func (s *SessionStore) Create(
	ctx context.Context,
	userID uuid.UUID,
	metadata map[string]interface{},
) (*auth.Session, error) {
	// Extract metadata fields with defaults
	orgID, _ := metadata["organization_id"].(uuid.UUID)
	authMethod, _ := metadata["auth_method"].(string)
	authProvider, _ := metadata["auth_provider"].(string)
	ipAddress, _ := metadata["ip_address"].(string)
	userAgent, _ := metadata["user_agent"].(string)

	// Generate secure session token
	token, err := generateSessionToken()
	if err != nil {
		return nil, fmt.Errorf("generate session token: %w", err)
	}

	// Create session entity
	session := &auth.SessionEntity{
		ID:             uuid.New(),
		UserID:         userID,
		Token:          token,
		OrganizationID: orgID,
		ExpiresAt:      time.Now().Add(s.ttl),
		IPAddress:      ipAddress,
		UserAgent:      userAgent,
		AuthMethod:     authMethod,
		AuthProvider:   authProvider,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Store in database
	created, err := s.repo.Create(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("create session in db: %w", err)
	}

	// Store in cache with TTL
	if s.cache != nil {
		_ = s.cache.Set(ctx, created.ID.String(), created, s.ttl)
	}

	return convertToAuthSession(created), nil
}

// Get retrieves a session by ID.
func (s *SessionStore) Get(ctx context.Context, sessionID uuid.UUID) (*auth.Session, error) {
	// Try cache first
	if s.cache != nil {
		cached, err := s.cache.Get(ctx, sessionID.String())
		if err == nil && cached != nil {
			// Validate expiry
			if !s.isSessionExpired(cached) {
				return convertToAuthSession(cached), nil
			}
			// If expired, delete from cache
			_ = s.cache.Delete(ctx, sessionID.String())
		}
	}

	// Fallback to database
	session, err := s.repo.Get(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Validate expiry
	if s.isSessionExpired(session) {
		// Clean up expired session
		_ = s.repo.Delete(ctx, sessionID)
		return nil, auth.ErrSessionExpired
	}

	// Update cache
	if s.cache != nil && session != nil {
		_ = s.cache.Set(ctx, session.ID.String(), session, s.ttl)
	}

	return convertToAuthSession(session), nil
}

// GetByToken retrieves a session by token.
func (s *SessionStore) GetByToken(ctx context.Context, token string) (*auth.Session, error) {
	// Cache doesn't support token-based lookup, use database directly
	session, err := s.repo.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Validate expiry
	if s.isSessionExpired(session) {
		// Clean up expired session
		_ = s.repo.Delete(ctx, session.ID)
		return nil, auth.ErrSessionExpired
	}

	// Update cache by ID
	if s.cache != nil && session != nil {
		_ = s.cache.Set(ctx, session.ID.String(), session, s.ttl)
	}

	return convertToAuthSession(session), nil
}

// Update updates session metadata.
func (s *SessionStore) Update(
	ctx context.Context,
	sessionID uuid.UUID,
	updates *auth.Session,
) (*auth.Session, error) {
	sessionsUpdate := &auth.SessionEntity{
		UpdatedAt:      updates.UpdatedAt,
		ExpiresAt:      updates.ExpiresAt,
		IPAddress:      updates.IPAddress,
		UserAgent:      updates.UserAgent,
		AuthMethod:     updates.AuthMethod,
		AuthProvider:   updates.AuthProvider,
		OrganizationID: updates.OrganizationID,
	}

	// Update in database
	updated, err := s.repo.Update(ctx, sessionID, sessionsUpdate)
	if err != nil {
		return nil, err
	}

	// Update cache
	if s.cache != nil && updated != nil {
		_ = s.cache.Set(ctx, updated.ID.String(), updated, s.ttl)
	}

	return convertToAuthSession(updated), nil
}

// Delete removes a session.
func (s *SessionStore) Delete(ctx context.Context, sessionID uuid.UUID) error {
	// Delete from database
	if err := s.repo.Delete(ctx, sessionID); err != nil {
		return err
	}

	// Delete from cache
	if s.cache != nil {
		_ = s.cache.Delete(ctx, sessionID.String())
	}

	return nil
}

// DeleteByToken removes a session by token.
func (s *SessionStore) DeleteByToken(ctx context.Context, token string) error {
	// Get session first
	session, err := s.GetByToken(ctx, token)
	if err != nil {
		return err
	}

	// Delete by ID
	return s.Delete(ctx, session.ID)
}

// DeleteByUser removes all sessions for a user.
func (s *SessionStore) DeleteByUser(ctx context.Context, userID uuid.UUID) error {
	// Delete from database
	if err := s.repo.DeleteByUser(ctx, userID); err != nil {
		return err
	}

	// Cache entries will expire naturally with TTL
	return nil
}

// DeleteByUserID removes all sessions for a user (alias for DeleteByUser).
func (s *SessionStore) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return s.DeleteByUser(ctx, userID)
}

// List returns all active sessions for a user (interface method).
func (s *SessionStore) List(ctx context.Context, userID uuid.UUID) ([]*auth.Session, error) {
	return s.ListByUser(ctx, userID)
}

// ListByUser returns all active sessions for a user.
func (s *SessionStore) ListByUser(ctx context.Context, userID uuid.UUID) ([]*auth.Session, error) {
	// For now, use database directly
	params := auth.ListSessionsParams{
		Page: auth.PageQuery{},
	}

	dbSessions, _, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	// Filter for user and non-expired sessions
	var activeSessions []*auth.Session
	for _, session := range dbSessions {
		if session.UserID == userID && !s.isSessionExpired(session) {
			activeSessions = append(activeSessions, convertToAuthSession(session))
		}
	}

	return activeSessions, nil
}

// RefreshSession extends the expiry of a session.
func (s *SessionStore) RefreshSession(
	ctx context.Context,
	sessionID uuid.UUID,
) (*auth.Session, error) {
	session, err := s.Get(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	// Update expiry
	session.ExpiresAt = time.Now().Add(s.ttl)
	session.UpdatedAt = time.Now()

	// Update in database
	updated, err := s.Update(ctx, sessionID, session)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

// CleanupExpired removes all expired sessions.
func (s *SessionStore) CleanupExpired(ctx context.Context) error {
	return s.repo.DeleteExpired(ctx)
}

// Validate checks if a session is valid.
func (s *SessionStore) Validate(ctx context.Context, token string) (*auth.Session, error) {
	session, err := s.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Update last activity
	session.UpdatedAt = time.Now()
	updated, err := s.Update(ctx, session.ID, session)
	if err != nil {
		// Log error but don't fail validation
		return session, nil
	}

	return updated, nil
}

// Helper functions

func (s *SessionStore) isSessionExpired(session *auth.SessionEntity) bool {
	if session == nil {
		return true
	}
	return time.Now().After(session.ExpiresAt)
}

func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func convertToAuthSession(s *auth.SessionEntity) *auth.Session {
	if s == nil {
		return nil
	}
	return &auth.Session{
		ID:             s.ID,
		UserID:         s.UserID,
		Token:          s.Token,
		OrganizationID: s.OrganizationID,
		ExpiresAt:      s.ExpiresAt,
		IPAddress:      s.IPAddress,
		UserAgent:      s.UserAgent,
		AuthMethod:     s.AuthMethod,
		AuthProvider:   s.AuthProvider,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
	}
}
