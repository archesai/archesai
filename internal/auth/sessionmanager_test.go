package auth

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testTokenConst = "test-token"
)

func TestSessionManager_CreateSession(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	orgID := uuid.New()
	ipAddress := "192.168.1.1"
	userAgent := "Test Browser"

	t.Run("successful creation with cache", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockCache := NewMockCache(t)
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		// Expectations
		mockRepo.EXPECT().CreateSession(ctx, mock.AnythingOfType("*auth.Session")).Return(&Session{
			Id:                   uuid.New(),
			UserId:               userID,
			Token:                "generated-token",
			ActiveOrganizationId: orgID,
			ExpiresAt:            time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			IpAddress:            ipAddress,
			UserAgent:            userAgent,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}, nil)

		mockCache.EXPECT().SetSession(ctx, mock.AnythingOfType("*auth.Session"), 24*time.Hour).Return(nil).Maybe()
		mockCache.EXPECT().SetSessionByToken(ctx, mock.AnythingOfType("string"), mock.AnythingOfType("*auth.Session"), 24*time.Hour).Return(nil).Maybe()

		// Execute
		session, err := sm.CreateSession(ctx, userID, orgID, ipAddress, userAgent)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, userID, session.UserId)
		assert.Equal(t, orgID, session.ActiveOrganizationId)
		assert.Equal(t, ipAddress, session.IpAddress)
		assert.Equal(t, userAgent, session.UserAgent)
		assert.NotEmpty(t, session.Token)
	})

	t.Run("successful creation without cache", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		sm := NewSessionManager(mockRepo, nil, 24*time.Hour)

		// Expectations
		mockRepo.EXPECT().CreateSession(ctx, mock.AnythingOfType("*auth.Session")).Return(&Session{
			Id:                   uuid.New(),
			UserId:               userID,
			Token:                "generated-token",
			ActiveOrganizationId: orgID,
			ExpiresAt:            time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			IpAddress:            ipAddress,
			UserAgent:            userAgent,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}, nil)

		// Execute
		session, err := sm.CreateSession(ctx, userID, orgID, ipAddress, userAgent)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, userID, session.UserId)
	})
}

func TestSessionManager_GetSession(t *testing.T) {
	ctx := context.Background()
	sessionID := uuid.New()
	userID := uuid.New()

	validSession := &Session{
		Id:        sessionID,
		UserId:    userID,
		Token:     testTokenConst,
		ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	expiredSession := &Session{
		Id:        sessionID,
		UserId:    userID,
		Token:     "expired-token",
		ExpiresAt: time.Now().Add(-time.Hour).Format(time.RFC3339),
		CreatedAt: time.Now().Add(-2 * time.Hour),
		UpdatedAt: time.Now().Add(-2 * time.Hour),
	}

	t.Run("cache hit with valid session", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockCache := NewMockCache(t)
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		// Expectations
		mockCache.EXPECT().GetSession(ctx, sessionID).Return(validSession, nil)

		// Execute
		session, err := sm.GetSession(ctx, sessionID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, sessionID, session.Id)
		assert.Equal(t, userID, session.UserId)
	})

	t.Run("cache miss falls back to database", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockCache := NewMockCache(t)
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		// Expectations
		mockCache.EXPECT().GetSession(ctx, sessionID).Return(nil, ErrCacheMiss)
		mockRepo.EXPECT().GetSession(ctx, sessionID).Return(validSession, nil)
		mockCache.EXPECT().SetSession(ctx, validSession, 24*time.Hour).Return(nil).Maybe()

		// Execute
		session, err := sm.GetSession(ctx, sessionID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, sessionID, session.Id)
	})

	t.Run("expired session in cache", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockCache := NewMockCache(t)
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		// Expectations
		mockCache.EXPECT().GetSession(ctx, sessionID).Return(expiredSession, nil)
		mockCache.EXPECT().DeleteSession(ctx, sessionID).Return(nil).Maybe()
		// Fallback to database after cache cleanup
		mockRepo.EXPECT().GetSession(ctx, sessionID).Return(expiredSession, nil)
		mockRepo.EXPECT().DeleteSession(ctx, sessionID).Return(nil)

		// Execute
		session, err := sm.GetSession(ctx, sessionID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, ErrSessionExpired, err)
		assert.Nil(t, session)
	})
}

func TestSessionManager_GetSessionByToken(t *testing.T) {
	ctx := context.Background()
	token := testTokenConst
	sessionID := uuid.New()
	userID := uuid.New()

	validSession := &Session{
		Id:        sessionID,
		UserId:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("successful retrieval from cache", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockCache := NewMockCache(t)
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		// Expectations
		mockCache.EXPECT().GetSessionByToken(ctx, token).Return(validSession, nil)

		// Execute
		session, err := sm.GetSessionByToken(ctx, token)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, token, session.Token)
		assert.Equal(t, userID, session.UserId)
	})

	t.Run("cache miss with database fallback", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockCache := NewMockCache(t)
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		// Expectations
		mockCache.EXPECT().GetSessionByToken(ctx, token).Return(nil, ErrCacheMiss)
		mockRepo.EXPECT().GetSessionByToken(ctx, token).Return(validSession, nil)
		mockCache.EXPECT().SetSessionByToken(ctx, token, validSession, 24*time.Hour).Return(nil).Maybe()

		// Execute
		session, err := sm.GetSessionByToken(ctx, token)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, token, session.Token)
	})
}

func TestSessionManager_UpdateSession(t *testing.T) {
	ctx := context.Background()
	sessionID := uuid.New()
	userID := uuid.New()

	updates := &Session{
		Id:        sessionID,
		UserId:    userID,
		Token:     "updated-token",
		UpdatedAt: time.Now(),
	}

	t.Run("successful update with cache", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockCache := NewMockCache(t)
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		updatedSession := &Session{
			Id:        sessionID,
			UserId:    userID,
			Token:     "updated-token",
			ExpiresAt: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			UpdatedAt: time.Now(),
		}

		// Expectations
		mockRepo.EXPECT().UpdateSession(ctx, sessionID, updates).Return(updatedSession, nil)
		mockCache.EXPECT().SetSession(ctx, updatedSession, 24*time.Hour).Return(nil).Maybe()
		mockCache.EXPECT().SetSessionByToken(ctx, "updated-token", updatedSession, 24*time.Hour).Return(nil).Maybe()

		// Execute
		session, err := sm.UpdateSession(ctx, sessionID, updates)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, "updated-token", session.Token)
	})
}

func TestSessionManager_DeleteSession(t *testing.T) {
	ctx := context.Background()
	sessionID := uuid.New()
	userID := uuid.New()
	token := testTokenConst

	existingSession := &Session{
		Id:        sessionID,
		UserId:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
	}

	t.Run("successful deletion with cache", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockCache := NewMockCache(t)
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		// Expectations
		mockCache.EXPECT().GetSession(ctx, sessionID).Return(existingSession, nil).Maybe()
		mockRepo.EXPECT().GetSession(ctx, sessionID).Return(existingSession, nil).Maybe()
		mockRepo.EXPECT().DeleteSession(ctx, sessionID).Return(nil)
		mockCache.EXPECT().DeleteSession(ctx, sessionID).Return(nil).Maybe()
		mockCache.EXPECT().DeleteSessionByToken(ctx, token).Return(nil).Maybe()

		// Execute
		err := sm.DeleteSession(ctx, sessionID)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("successful deletion without cache", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		sm := NewSessionManager(mockRepo, nil, 24*time.Hour)

		// Expectations
		mockRepo.EXPECT().GetSession(ctx, sessionID).Return(existingSession, nil).Maybe()
		mockRepo.EXPECT().DeleteSession(ctx, sessionID).Return(nil)

		// Execute
		err := sm.DeleteSession(ctx, sessionID)

		// Assert
		assert.NoError(t, err)
	})
}

func TestSessionManager_ListUserSessions(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	validSession1 := &Session{
		Id:        uuid.New(),
		UserId:    userID,
		Token:     "token1",
		ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
		CreatedAt: time.Now(),
	}

	validSession2 := &Session{
		Id:        uuid.New(),
		UserId:    userID,
		Token:     "token2",
		ExpiresAt: time.Now().Add(2 * time.Hour).Format(time.RFC3339),
		CreatedAt: time.Now(),
	}

	expiredSession := &Session{
		Id:        uuid.New(),
		UserId:    userID,
		Token:     "expired-token",
		ExpiresAt: time.Now().Add(-time.Hour).Format(time.RFC3339),
		CreatedAt: time.Now().Add(-2 * time.Hour),
	}

	t.Run("returns only active sessions", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockCache := NewMockCache(t)
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		userIDStr := userID.String()
		params := ListSessionsParams{
			UserID: &userIDStr,
			Limit:  100,
		}

		// Expectations
		mockRepo.EXPECT().ListSessions(ctx, params).Return([]*Session{
			validSession1,
			validSession2,
			expiredSession,
		}, int64(3), nil)

		// Execute
		sessions, err := sm.ListUserSessions(ctx, userID)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, sessions, 2) // Only 2 active sessions
		assert.Contains(t, sessions, validSession1)
		assert.Contains(t, sessions, validSession2)
		assert.NotContains(t, sessions, expiredSession)
	})
}

func TestSessionManager_ValidateSession(t *testing.T) {
	ctx := context.Background()
	token := testTokenConst
	sessionID := uuid.New()
	userID := uuid.New()

	validSession := &Session{
		Id:        sessionID,
		UserId:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("valid session updates last activity", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockCache := NewMockCache(t)
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		updatedSession := &Session{
			Id:        sessionID,
			UserId:    userID,
			Token:     token,
			ExpiresAt: validSession.ExpiresAt,
			UpdatedAt: time.Now(),
		}

		// Expectations
		mockCache.EXPECT().GetSessionByToken(ctx, token).Return(validSession, nil)
		mockRepo.EXPECT().UpdateSession(ctx, sessionID, mock.AnythingOfType("*auth.Session")).Return(updatedSession, nil)
		mockCache.EXPECT().SetSession(ctx, updatedSession, 24*time.Hour).Return(nil).Maybe()
		mockCache.EXPECT().SetSessionByToken(ctx, token, updatedSession, 24*time.Hour).Return(nil).Maybe()

		// Execute
		session, err := sm.ValidateSession(ctx, token)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, session)
		// UpdatedAt should be recent
		assert.True(t, time.Since(session.UpdatedAt) < 5*time.Second)
	})

	t.Run("expired session is deleted", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockCache := NewMockCache(t)
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		expiredSession := &Session{
			Id:        sessionID,
			UserId:    userID,
			Token:     token,
			ExpiresAt: time.Now().Add(-time.Hour).Format(time.RFC3339),
		}

		// Expectations - GetSessionByToken will return expired session from cache
		mockCache.EXPECT().GetSessionByToken(ctx, token).Return(expiredSession, nil)
		mockCache.EXPECT().DeleteSessionByToken(ctx, token).Return(nil).Maybe()
		mockCache.EXPECT().DeleteSession(ctx, sessionID).Return(nil).Maybe()
		// Then fallback to database which also returns expired
		mockRepo.EXPECT().GetSessionByToken(ctx, token).Return(expiredSession, nil)
		mockRepo.EXPECT().DeleteSession(ctx, sessionID).Return(nil).Maybe()
		// Additional calls for cleanup
		mockCache.EXPECT().GetSession(ctx, sessionID).Return(expiredSession, nil).Maybe()
		mockRepo.EXPECT().GetSession(ctx, sessionID).Return(expiredSession, nil).Maybe()

		// Execute
		session, err := sm.ValidateSession(ctx, token)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, ErrSessionExpired, err)
		assert.Nil(t, session)
	})
}

func TestSessionManager_RefreshSession(t *testing.T) {
	ctx := context.Background()
	sessionID := uuid.New()
	userID := uuid.New()
	token := testTokenConst

	existingSession := &Session{
		Id:        sessionID,
		UserId:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
		CreatedAt: time.Now().Add(-time.Hour),
		UpdatedAt: time.Now().Add(-time.Hour),
	}

	t.Run("successful refresh extends expiry", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockCache := NewMockCache(t)
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		// Expectations
		mockCache.EXPECT().GetSession(ctx, sessionID).Return(existingSession, nil)
		mockRepo.EXPECT().UpdateSession(ctx, sessionID, mock.AnythingOfType("*auth.Session")).Return(&Session{
			Id:        sessionID,
			UserId:    userID,
			Token:     token,
			ExpiresAt: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			UpdatedAt: time.Now(),
		}, nil)
		mockCache.EXPECT().SetSession(ctx, mock.AnythingOfType("*auth.Session"), 24*time.Hour).Return(nil).Maybe()
		mockCache.EXPECT().SetSessionByToken(ctx, token, mock.AnythingOfType("*auth.Session"), 24*time.Hour).Return(nil).Maybe()

		// Execute
		session, err := sm.RefreshSession(ctx, sessionID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, session)

		// Check that expiry was extended (session.ExpiresAt should be about 24 hours from now)
		newExpiry, _ := time.Parse(time.RFC3339, session.ExpiresAt)
		assert.True(t, time.Until(newExpiry) > 23*time.Hour)
	})
}

// Test OAuth state management functions
func TestOAuthStateManagement(t *testing.T) {
	// OAuth state is managed in-memory, not through cache
	sessionMgr := &SessionManager{}
	ctx := context.Background()

	// Reset the global state store for testing
	oauthStateStore.Lock()
	oauthStateStore.states = make(map[string]oauthState)
	oauthStateStore.Unlock()

	t.Run("store and retrieve OAuth state", func(t *testing.T) {
		state := "test-state-123"
		provider := "google"
		redirectURI := "http://localhost:8080/callback"
		ttl := 10 * time.Minute

		// Store state
		err := sessionMgr.StoreOAuthState(ctx, state, provider, redirectURI, ttl)
		assert.NoError(t, err)

		// Retrieve redirect URI
		retrievedURI, err := sessionMgr.GetOAuthRedirectURI(ctx, state)
		assert.NoError(t, err)
		assert.Equal(t, redirectURI, retrievedURI)

		// Delete state
		err = sessionMgr.DeleteOAuthState(ctx, state)
		assert.NoError(t, err)

		// Verify state is deleted
		_, err = sessionMgr.GetOAuthRedirectURI(ctx, state)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "state not found")
	})

	t.Run("retrieve non-existent state", func(t *testing.T) {
		_, err := sessionMgr.GetOAuthRedirectURI(ctx, "non-existent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "state not found")
	})

	t.Run("expired state is not returned", func(t *testing.T) {
		state := "expire-test"
		provider := "github"
		redirectURI := "http://localhost:8080/callback"
		ttl := 1 * time.Millisecond // Very short TTL

		// Store state
		err := sessionMgr.StoreOAuthState(ctx, state, provider, redirectURI, ttl)
		assert.NoError(t, err)

		// Wait for expiration
		time.Sleep(2 * time.Millisecond)

		// Try to retrieve - should fail
		_, err = sessionMgr.GetOAuthRedirectURI(ctx, state)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "state not found")
	})
}
