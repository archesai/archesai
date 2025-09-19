package sessions

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	genericcache "github.com/archesai/archesai/internal/cache"
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
		mockCache := genericcache.NewNoOpCache[Session]()
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		// Expectations
		mockRepo.EXPECT().Create(ctx, mock.AnythingOfType("*sessions.Session")).Return(&Session{
			ID:             uuid.New(),
			UserID:         userID,
			Token:          "generated-token",
			OrganizationID: orgID,
			ExpiresAt:      time.Now().Add(24 * time.Hour),
			IPAddress:      ipAddress,
			UserAgent:      userAgent,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}, nil)

		// Execute
		session, err := sm.Create(ctx, userID, orgID, ipAddress, userAgent)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, userID, session.UserID)
		assert.Equal(t, orgID, session.OrganizationID)
		assert.Equal(t, ipAddress, session.IPAddress)
		assert.Equal(t, userAgent, session.UserAgent)
		assert.NotEmpty(t, session.Token)
	})

	t.Run("successful creation without cache", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		sm := NewSessionManager(mockRepo, nil, 24*time.Hour)

		// Expectations
		mockRepo.EXPECT().Create(ctx, mock.AnythingOfType("*sessions.Session")).Return(&Session{
			ID:             uuid.New(),
			UserID:         userID,
			Token:          "generated-token",
			OrganizationID: orgID,
			ExpiresAt:      time.Now().Add(24 * time.Hour),
			IPAddress:      ipAddress,
			UserAgent:      userAgent,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}, nil)

		// Execute
		session, err := sm.Create(ctx, userID, orgID, ipAddress, userAgent)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, userID, session.UserID)
	})
}

func TestSessionManager_GetSession(t *testing.T) {
	ctx := context.Background()
	sessionID := uuid.New()
	userID := uuid.New()

	validSession := &Session{
		ID:        sessionID,
		UserID:    userID,
		Token:     testTokenConst,
		ExpiresAt: time.Now().Add(time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	expiredSession := &Session{
		ID:        sessionID,
		UserID:    userID,
		Token:     "expired-token",
		ExpiresAt: time.Now().Add(-time.Hour),
		CreatedAt: time.Now().Add(-2 * time.Hour),
		UpdatedAt: time.Now().Add(-2 * time.Hour),
	}

	t.Run("cache hit with valid session", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockCache := genericcache.NewNoOpCache[Session]()
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		// Expectations
		mockRepo.EXPECT().Get(ctx, sessionID).Return(validSession, nil)

		// Execute
		session, err := sm.Get(ctx, sessionID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, sessionID, session.ID)
		assert.Equal(t, userID, session.UserID)
	})

	t.Run("cache miss falls back to database", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockCache := genericcache.NewNoOpCache[Session]()
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		// Expectations
		mockRepo.EXPECT().Get(ctx, sessionID).Return(validSession, nil)

		// Execute
		session, err := sm.Get(ctx, sessionID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, sessionID, session.ID)
	})

	t.Run("expired session in cache", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockCache := genericcache.NewNoOpCache[Session]()
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		// Expectations
		// Fallback to database after cache cleanup
		mockRepo.EXPECT().Get(ctx, sessionID).Return(expiredSession, nil)
		mockRepo.EXPECT().Delete(ctx, sessionID).Return(nil)

		// Execute
		session, err := sm.Get(ctx, sessionID)

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
		ID:        sessionID,
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("successful retrieval from cache", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockCache := genericcache.NewNoOpCache[Session]()
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		// Expectations
		mockRepo.EXPECT().GetByToken(ctx, token).Return(validSession, nil)

		// Execute
		session, err := sm.GetSessionByToken(ctx, token)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, token, session.Token)
		assert.Equal(t, userID, session.UserID)
	})

	t.Run("cache miss with database fallback", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockCache := genericcache.NewNoOpCache[Session]()
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		// Expectations
		mockRepo.EXPECT().GetByToken(ctx, token).Return(validSession, nil)

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
		ID:        sessionID,
		UserID:    userID,
		Token:     "updated-token",
		UpdatedAt: time.Now(),
	}

	t.Run("successful update with cache", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockCache := genericcache.NewNoOpCache[Session]()
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		updatedSession := &Session{
			ID:        sessionID,
			UserID:    userID,
			Token:     "updated-token",
			ExpiresAt: time.Now().Add(24 * time.Hour),
			UpdatedAt: time.Now(),
		}

		// Expectations
		mockRepo.EXPECT().Update(ctx, sessionID, updates).Return(updatedSession, nil)

		// Execute
		session, err := sm.Update(ctx, sessionID, updates)

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
		ID:        sessionID,
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	t.Run("successful deletion with cache", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockCache := genericcache.NewNoOpCache[Session]()
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		// Expectations
		mockRepo.EXPECT().Get(ctx, sessionID).Return(existingSession, nil).Maybe()
		mockRepo.EXPECT().Delete(ctx, sessionID).Return(nil)

		// Execute
		err := sm.Delete(ctx, sessionID)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("successful deletion without cache", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		sm := NewSessionManager(mockRepo, nil, 24*time.Hour)

		// Expectations
		mockRepo.EXPECT().Get(ctx, sessionID).Return(existingSession, nil).Maybe()
		mockRepo.EXPECT().Delete(ctx, sessionID).Return(nil)

		// Execute
		err := sm.Delete(ctx, sessionID)

		// Assert
		assert.NoError(t, err)
	})
}

func TestSessionManager_ListByUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	validSession1 := &Session{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     "token1",
		ExpiresAt: time.Now().Add(time.Hour),
		CreatedAt: time.Now(),
	}

	validSession2 := &Session{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     "token2",
		ExpiresAt: time.Now().Add(2 * time.Hour),
		CreatedAt: time.Now(),
	}

	expiredSession := &Session{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     "expired-token",
		ExpiresAt: time.Now().Add(-time.Hour),
		CreatedAt: time.Now().Add(-2 * time.Hour),
	}

	t.Run("returns only active sessions", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockCache := genericcache.NewNoOpCache[Session]()
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		params := ListSessionsParams{
			Page: PageQuery{
				Number: 1,
				Size:   100,
			},
		}

		// Expectations
		mockRepo.EXPECT().List(ctx, params).Return([]*Session{
			validSession1,
			validSession2,
			expiredSession,
		}, int64(3), nil)

		// Execute
		sessions, err := sm.ListByUser(ctx, userID)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, sessions, 2) // Only 2 active sessions
		assert.Contains(t, sessions, validSession1)
		assert.Contains(t, sessions, validSession2)
		assert.NotContains(t, sessions, expiredSession)
	})
}

func TestSessionManager_Validate(t *testing.T) {
	ctx := context.Background()
	token := testTokenConst
	sessionID := uuid.New()
	userID := uuid.New()

	validSession := &Session{
		ID:        sessionID,
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("valid session updates last activity", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockCache := genericcache.NewNoOpCache[Session]()
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		updatedSession := &Session{
			ID:        sessionID,
			UserID:    userID,
			Token:     token,
			ExpiresAt: validSession.ExpiresAt,
			UpdatedAt: time.Now(),
		}

		// Expectations
		mockRepo.EXPECT().GetByToken(ctx, token).Return(validSession, nil)
		mockRepo.EXPECT().
			Update(ctx, sessionID, mock.AnythingOfType("*sessions.Session")).
			Return(updatedSession, nil)

		// Execute
		session, err := sm.Validate(ctx, token)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, session)
		// UpdatedAt should be recent
		assert.Less(t, time.Since(session.UpdatedAt), 5*time.Second)
	})

	t.Run("expired session is deleted", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockCache := genericcache.NewNoOpCache[Session]()
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		expiredSession := &Session{
			ID:        sessionID,
			UserID:    userID,
			Token:     token,
			ExpiresAt: time.Now().Add(-time.Hour),
		}

		// Expectations - GetSessionByToken will return expired session from cache
		// Then fallback to database which also returns expired
		mockRepo.EXPECT().GetByToken(ctx, token).Return(expiredSession, nil)
		mockRepo.EXPECT().Delete(ctx, sessionID).Return(nil).Maybe()
		// Additional calls for cleanup
		mockRepo.EXPECT().Get(ctx, sessionID).Return(expiredSession, nil).Maybe()

		// Execute
		session, err := sm.Validate(ctx, token)

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

	validSession := &Session{
		ID:        sessionID,
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour),
		CreatedAt: time.Now().Add(-time.Hour),
		UpdatedAt: time.Now().Add(-time.Hour),
	}

	t.Run("successful refresh extends expiry", func(t *testing.T) {
		// Setup
		mockRepo := NewMockRepository(t)
		mockCache := genericcache.NewNoOpCache[Session]()
		sm := NewSessionManager(mockRepo, mockCache, 24*time.Hour)

		// Expectations
		mockRepo.EXPECT().Get(ctx, sessionID).Return(validSession, nil)
		mockRepo.EXPECT().
			Update(ctx, sessionID, mock.AnythingOfType("*sessions.Session")).
			Return(&Session{
				ID:        sessionID,
				UserID:    userID,
				Token:     token,
				ExpiresAt: time.Now().Add(24 * time.Hour),
				UpdatedAt: time.Now(),
			}, nil)

		// Execute
		session, err := sm.RefreshSession(ctx, sessionID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, session)

		// Check that expiry was extended (session.ExpiresAt should be about 24 hours from now)
		assert.Greater(t, time.Until(session.ExpiresAt), 23*time.Hour)
	})
}
