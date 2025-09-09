package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	testValidToken = "valid-token"
)

// TestSessionManager_DeleteSessionByToken tests deleting session by token
func TestSessionManager_DeleteSessionByToken(t *testing.T) {
	ctx := context.Background()
	token := "delete-by-token"
	sessionID := uuid.New()
	userID := uuid.New()

	t.Run("successful deletion", func(t *testing.T) {
		mockRepo := NewMockRepository(t)
		mockCache := NewMockCache(t)
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		session := &Session{
			Id:        sessionID,
			UserId:    userID,
			Token:     token,
			ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
		}

		// Expectations for GetSessionByToken
		mockCache.EXPECT().GetSessionByToken(ctx, token).Return(session, nil)

		// Expectations for DeleteSession (called internally)
		mockCache.EXPECT().GetSession(ctx, sessionID).Return(session, nil).Maybe()
		mockRepo.EXPECT().DeleteSession(ctx, sessionID).Return(nil)
		mockCache.EXPECT().DeleteSession(ctx, sessionID).Return(nil).Maybe()
		mockCache.EXPECT().DeleteSessionByToken(ctx, token).Return(nil).Maybe()

		err := sm.DeleteSessionByToken(ctx, token)

		assert.NoError(t, err)
	})

	t.Run("session not found", func(t *testing.T) {
		mockRepo := NewMockRepository(t)
		mockCache := NewMockCache(t)
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		// Session not in cache
		mockCache.EXPECT().GetSessionByToken(ctx, token).Return(nil, ErrCacheMiss)
		// Session not in database either
		mockRepo.EXPECT().GetSessionByToken(ctx, token).Return(nil, ErrSessionNotFound)

		err := sm.DeleteSessionByToken(ctx, token)

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrSessionNotFound)
	})
}

// TestSessionManager_DeleteUserSessions tests deleting all sessions for a user
func TestSessionManager_DeleteUserSessions(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	t.Run("successful deletion with cache", func(t *testing.T) {
		mockRepo := NewMockRepository(t)
		mockCache := NewMockCache(t)
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		mockRepo.EXPECT().DeleteUserSessions(ctx, userID).Return(nil)
		mockCache.EXPECT().DeleteUserSessions(ctx, userID).Return(nil).Maybe()

		err := sm.DeleteUserSessions(ctx, userID)

		assert.NoError(t, err)
	})

	t.Run("successful deletion without cache", func(t *testing.T) {
		mockRepo := NewMockRepository(t)
		sm := NewSessionManager(mockRepo, nil, 24*time.Hour)

		mockRepo.EXPECT().DeleteUserSessions(ctx, userID).Return(nil)

		err := sm.DeleteUserSessions(ctx, userID)

		assert.NoError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := NewMockRepository(t)
		mockCache := NewMockCache(t)
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		mockRepo.EXPECT().DeleteUserSessions(ctx, userID).Return(errors.New("db error"))

		err := sm.DeleteUserSessions(ctx, userID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

// TestSessionManager_CleanupExpiredSessions tests cleanup of expired sessions
func TestSessionManager_CleanupExpiredSessions(t *testing.T) {
	ctx := context.Background()

	t.Run("successful cleanup", func(t *testing.T) {
		mockRepo := NewMockRepository(t)
		mockCache := NewMockCache(t)
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		mockRepo.EXPECT().DeleteExpiredSessions(ctx).Return(nil)

		err := sm.CleanupExpiredSessions(ctx)

		assert.NoError(t, err)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := NewMockRepository(t)
		sm := NewSessionManager(mockRepo, nil, 24*time.Hour)

		mockRepo.EXPECT().DeleteExpiredSessions(ctx).Return(errors.New("cleanup failed"))

		err := sm.CleanupExpiredSessions(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cleanup failed")
	})
}

// TestSessionManager_EdgeCases tests edge cases and error conditions
func TestSessionManager_EdgeCases(t *testing.T) {
	ctx := context.Background()

	t.Run("NewSessionManager with zero TTL uses default", func(t *testing.T) {
		mockRepo := NewMockRepository(t)
		sm := NewSessionManager(mockRepo, nil, 0)

		assert.NotNil(t, sm)
		assert.Equal(t, 30*24*time.Hour, sm.ttl)
	})

	t.Run("CreateSession with repository error", func(t *testing.T) {
		mockRepo := NewMockRepository(t)
		sm := NewSessionManager(mockRepo, nil, 24*time.Hour)

		mockRepo.EXPECT().CreateSession(ctx, mock.AnythingOfType("*auth.Session")).Return(nil, errors.New("db error"))

		session, err := sm.CreateSession(ctx, uuid.New(), uuid.New(), "127.0.0.1", "browser")

		assert.Error(t, err)
		assert.Nil(t, session)
		assert.Contains(t, err.Error(), "create session in db")
	})

	t.Run("GetSession with nil cache handles gracefully", func(t *testing.T) {
		mockRepo := NewMockRepository(t)
		sm := NewSessionManager(mockRepo, nil, 24*time.Hour)
		sessionID := uuid.New()

		validSession := &Session{
			Id:        sessionID,
			UserId:    uuid.New(),
			Token:     "token",
			ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
		}

		mockRepo.EXPECT().GetSession(ctx, sessionID).Return(validSession, nil)

		session, err := sm.GetSession(ctx, sessionID)

		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, sessionID, session.Id)
	})

	t.Run("UpdateSession with nil updates", func(t *testing.T) {
		mockRepo := NewMockRepository(t)
		sm := NewSessionManager(mockRepo, nil, 24*time.Hour)
		sessionID := uuid.New()

		updatedSession := &Session{
			Id:        sessionID,
			UserId:    uuid.New(),
			Token:     "updated",
			UpdatedAt: time.Now(),
		}

		mockRepo.EXPECT().UpdateSession(ctx, sessionID, mock.AnythingOfType("*auth.Session")).Return(updatedSession, nil)

		result, err := sm.UpdateSession(ctx, sessionID, &Session{})

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("ValidateSession with update error still returns session", func(t *testing.T) {
		mockRepo := NewMockRepository(t)
		mockCache := NewMockCache(t)
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		token := testValidToken
		sessionID := uuid.New()
		validSession := &Session{
			Id:        sessionID,
			UserId:    uuid.New(),
			Token:     token,
			ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
			UpdatedAt: time.Now().Add(-time.Hour),
		}

		mockCache.EXPECT().GetSessionByToken(ctx, token).Return(validSession, nil)
		mockRepo.EXPECT().UpdateSession(ctx, sessionID, mock.AnythingOfType("*auth.Session")).Return(nil, errors.New("update failed"))

		session, err := sm.ValidateSession(ctx, token)

		// Should still return the session even if update fails
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, token, session.Token)
	})
}

// TestSessionManager_IsSessionExpired tests session expiry validation
func TestSessionManager_IsSessionExpired(t *testing.T) {
	mockRepo := NewMockRepository(t)
	sm := NewSessionManager(mockRepo, nil, 24*time.Hour)

	t.Run("nil session is expired", func(t *testing.T) {
		assert.True(t, sm.isSessionExpired(nil))
	})

	t.Run("invalid expiry format is expired", func(t *testing.T) {
		session := &Session{
			ExpiresAt: "invalid-date",
		}
		assert.True(t, sm.isSessionExpired(session))
	})

	t.Run("past date is expired", func(t *testing.T) {
		session := &Session{
			ExpiresAt: time.Now().Add(-time.Hour).Format(time.RFC3339),
		}
		assert.True(t, sm.isSessionExpired(session))
	})

	t.Run("future date is not expired", func(t *testing.T) {
		session := &Session{
			ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
		}
		assert.False(t, sm.isSessionExpired(session))
	})
}

// TestGenerateSecureToken tests token generation
func TestGenerateSecureToken(t *testing.T) {
	t.Run("generates unique tokens", func(t *testing.T) {
		token1, err1 := generateSecureToken()
		token2, err2 := generateSecureToken()

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEmpty(t, token1)
		assert.NotEmpty(t, token2)
		assert.NotEqual(t, token1, token2)
	})

	t.Run("generates tokens of expected length", func(t *testing.T) {
		token, err := generateSecureToken()

		assert.NoError(t, err)
		// Base64 encoding of 32 bytes should be 44 characters (with padding)
		// or 43 without padding for URL encoding
		assert.GreaterOrEqual(t, len(token), 43)
	})
}

// TestSessionManager_ListUserSessionsErrors tests error cases
func TestSessionManager_ListUserSessionsErrors(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	t.Run("handles repository error", func(t *testing.T) {
		mockRepo := NewMockRepository(t)
		sm := NewSessionManager(mockRepo, nil, 24*time.Hour)

		userIDStr := userID.String()
		mockRepo.EXPECT().ListSessions(ctx, ListSessionsParams{
			UserID: &userIDStr,
			Limit:  100,
		}).Return(nil, int64(0), errors.New("db error"))

		sessions, err := sm.ListUserSessions(ctx, userID)

		assert.Error(t, err)
		assert.Nil(t, sessions)
	})
}

// TestSessionManager_RefreshSessionErrors tests refresh error cases
func TestSessionManager_RefreshSessionErrors(t *testing.T) {
	ctx := context.Background()
	sessionID := uuid.New()

	t.Run("session not found", func(t *testing.T) {
		mockRepo := NewMockRepository(t)
		mockCache := NewMockCache(t)
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		mockCache.EXPECT().GetSession(ctx, sessionID).Return(nil, ErrCacheMiss)
		mockRepo.EXPECT().GetSession(ctx, sessionID).Return(nil, ErrSessionNotFound)

		session, err := sm.RefreshSession(ctx, sessionID)

		assert.Error(t, err)
		assert.Nil(t, session)
	})

	t.Run("update fails", func(t *testing.T) {
		mockRepo := NewMockRepository(t)
		mockCache := NewMockCache(t)
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		existingSession := &Session{
			Id:        sessionID,
			UserId:    uuid.New(),
			Token:     "token",
			ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
		}

		mockCache.EXPECT().GetSession(ctx, sessionID).Return(existingSession, nil)
		mockRepo.EXPECT().UpdateSession(ctx, sessionID, mock.AnythingOfType("*auth.Session")).Return(nil, errors.New("update failed"))

		session, err := sm.RefreshSession(ctx, sessionID)

		assert.Error(t, err)
		assert.Nil(t, session)
	})
}
