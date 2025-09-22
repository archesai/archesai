package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/archesai/archesai/internal/auth"
	"github.com/archesai/archesai/internal/auth/stores"
	genericcache "github.com/archesai/archesai/internal/cache"
)

const (
	testTokenConst = "test-token"
)

func TestSessionStore_CreateSession(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	orgID := uuid.New()
	ipAddress := "192.168.1.1"
	userAgent := "Test Browser"

	t.Run("successful creation with cache", func(t *testing.T) {
		// Setup
		mockRepo := auth.NewMockSessionsRepository(t)
		mockCache := genericcache.NewNoOpCache[auth.SessionEntity]()
		store := stores.NewSessionStore(mockRepo, mockCache, 24*time.Hour)

		// Expectations
		mockRepo.EXPECT().
			Create(ctx, mock.AnythingOfType("*auth.SessionEntity")).
			Return(&auth.SessionEntity{
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
		metadata := map[string]interface{}{
			"organization_id": orgID,
			"ip_address":      ipAddress,
			"user_agent":      userAgent,
		}
		session, err := store.Create(ctx, userID, metadata)

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
		mockRepo := auth.NewMockSessionsRepository(t)
		store := stores.NewSessionStore(mockRepo, nil, 24*time.Hour)

		// Expectations
		mockRepo.EXPECT().
			Create(ctx, mock.AnythingOfType("*auth.SessionEntity")).
			Return(&auth.SessionEntity{
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
		metadata := map[string]interface{}{
			"organization_id": orgID,
			"ip_address":      ipAddress,
			"user_agent":      userAgent,
		}
		session, err := store.Create(ctx, userID, metadata)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, userID, session.UserID)
	})
}

func TestSessionStore_GetSession(t *testing.T) {
	ctx := context.Background()
	sessionID := uuid.New()
	userID := uuid.New()

	validSession := &auth.SessionEntity{
		ID:        sessionID,
		UserID:    userID,
		Token:     testTokenConst,
		ExpiresAt: time.Now().Add(time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	expiredSession := &auth.SessionEntity{
		ID:        sessionID,
		UserID:    userID,
		Token:     "expired-token",
		ExpiresAt: time.Now().Add(-time.Hour),
		CreatedAt: time.Now().Add(-2 * time.Hour),
		UpdatedAt: time.Now().Add(-2 * time.Hour),
	}

	t.Run("cache hit with valid session", func(t *testing.T) {
		// Setup
		mockRepo := auth.NewMockSessionsRepository(t)
		mockCache := genericcache.NewMemoryCache[auth.SessionEntity]()
		store := stores.NewSessionStore(mockRepo, mockCache, 24*time.Hour)

		// Pre-populate cache
		_ = mockCache.Set(ctx, sessionID.String(), validSession, time.Hour)

		// Execute
		session, err := store.Get(ctx, sessionID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, sessionID, session.ID)
		assert.Equal(t, userID, session.UserID)
	})

	t.Run("cache miss falls back to database", func(t *testing.T) {
		// Setup
		mockRepo := auth.NewMockSessionsRepository(t)
		mockCache := genericcache.NewMemoryCache[auth.SessionEntity]()
		store := stores.NewSessionStore(mockRepo, mockCache, 24*time.Hour)

		// Expectations
		mockRepo.EXPECT().Get(ctx, sessionID).Return(validSession, nil)

		// Execute
		session, err := store.Get(ctx, sessionID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, sessionID, session.ID)
	})

	t.Run("expired session returns error", func(t *testing.T) {
		// Setup
		mockRepo := auth.NewMockSessionsRepository(t)
		mockCache := genericcache.NewMemoryCache[auth.SessionEntity]()
		store := stores.NewSessionStore(mockRepo, mockCache, 24*time.Hour)

		// Expectations
		mockRepo.EXPECT().Get(ctx, sessionID).Return(expiredSession, nil)
		// Expired session gets deleted
		mockRepo.EXPECT().Delete(ctx, sessionID).Return(nil)

		// Execute
		session, err := store.Get(ctx, sessionID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, session)
	})
}

func TestSessionStore_GetSessionByToken(t *testing.T) {
	ctx := context.Background()
	token := testTokenConst
	sessionID := uuid.New()
	userID := uuid.New()

	validSession := &auth.SessionEntity{
		ID:        sessionID,
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("successful retrieval from database", func(t *testing.T) {
		// Setup
		mockRepo := auth.NewMockSessionsRepository(t)
		mockCache := genericcache.NewMemoryCache[auth.SessionEntity]()
		store := stores.NewSessionStore(mockRepo, mockCache, 24*time.Hour)

		// Expectations - since cache is empty, it will call the repo
		mockRepo.EXPECT().GetByToken(ctx, token).Return(validSession, nil)

		// Execute
		session, err := store.GetByToken(ctx, token)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, token, session.Token)
		assert.Equal(t, userID, session.UserID)
	})

	t.Run("cache miss with database fallback", func(t *testing.T) {
		// Setup
		mockRepo := auth.NewMockSessionsRepository(t)
		mockCache := genericcache.NewMemoryCache[auth.SessionEntity]()
		store := stores.NewSessionStore(mockRepo, mockCache, 24*time.Hour)

		// Expectations
		mockRepo.EXPECT().GetByToken(ctx, token).Return(validSession, nil)

		// Execute
		session, err := store.GetByToken(ctx, token)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, token, session.Token)
	})
}

func TestSessionStore_UpdateSession(t *testing.T) {
	ctx := context.Background()
	sessionID := uuid.New()
	userID := uuid.New()

	updates := &auth.Session{
		ID:        sessionID,
		UserID:    userID,
		Token:     "updated-token",
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	t.Run("successful update with cache", func(t *testing.T) {
		// Setup
		mockRepo := auth.NewMockSessionsRepository(t)
		mockCache := genericcache.NewMemoryCache[auth.SessionEntity]()
		store := stores.NewSessionStore(mockRepo, mockCache, 24*time.Hour)

		updatedSessionEntity := &auth.SessionEntity{
			ID:        sessionID,
			UserID:    userID,
			Token:     "updated-token",
			ExpiresAt: time.Now().Add(24 * time.Hour),
			UpdatedAt: time.Now(),
		}

		// Expectations - Update method expects SessionEntity internally
		mockRepo.EXPECT().
			Update(ctx, sessionID, mock.AnythingOfType("*auth.SessionEntity")).
			Return(updatedSessionEntity, nil)

		// Execute
		session, err := store.Update(ctx, sessionID, updates)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, "updated-token", session.Token)
	})
}

func TestSessionStore_DeleteSession(t *testing.T) {
	ctx := context.Background()
	sessionID := uuid.New()
	userID := uuid.New()
	token := testTokenConst

	existingSession := &auth.SessionEntity{
		ID:        sessionID,
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	t.Run("successful deletion with cache", func(t *testing.T) {
		// Setup
		mockRepo := auth.NewMockSessionsRepository(t)
		mockCache := genericcache.NewMemoryCache[auth.SessionEntity]()
		store := stores.NewSessionStore(mockRepo, mockCache, 24*time.Hour)

		// Pre-populate cache
		_ = mockCache.Set(ctx, sessionID.String(), existingSession, time.Hour)

		// Expectations
		mockRepo.EXPECT().Delete(ctx, sessionID).Return(nil)

		// Execute
		err := store.Delete(ctx, sessionID)

		// Assert
		assert.NoError(t, err)

		// Verify cache was cleared
		cached, _ := mockCache.Get(ctx, sessionID.String())
		assert.Nil(t, cached)
	})

	t.Run("successful deletion without cache", func(t *testing.T) {
		// Setup
		mockRepo := auth.NewMockSessionsRepository(t)
		store := stores.NewSessionStore(mockRepo, nil, 24*time.Hour)

		// Expectations
		mockRepo.EXPECT().Delete(ctx, sessionID).Return(nil)

		// Execute
		err := store.Delete(ctx, sessionID)

		// Assert
		assert.NoError(t, err)
	})
}

func TestSessionStore_ListByUser(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	validSession1 := &auth.SessionEntity{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     "token1",
		ExpiresAt: time.Now().Add(time.Hour),
		CreatedAt: time.Now(),
	}

	validSession2 := &auth.SessionEntity{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     "token2",
		ExpiresAt: time.Now().Add(2 * time.Hour),
		CreatedAt: time.Now(),
	}

	expiredSession := &auth.SessionEntity{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     "expired-token",
		ExpiresAt: time.Now().Add(-time.Hour),
		CreatedAt: time.Now().Add(-2 * time.Hour),
	}

	t.Run("returns only active sessions", func(t *testing.T) {
		// Setup
		mockRepo := auth.NewMockSessionsRepository(t)
		mockCache := genericcache.NewMemoryCache[auth.SessionEntity]()
		store := stores.NewSessionStore(mockRepo, mockCache, 24*time.Hour)

		// The implementation actually passes empty params
		params := auth.ListSessionsParams{
			Page: auth.PageQuery{},
		}

		// Expectations
		mockRepo.EXPECT().List(ctx, params).Return([]*auth.SessionEntity{
			validSession1,
			validSession2,
			expiredSession,
		}, int64(3), nil)

		// Execute
		sessions, err := store.ListByUser(ctx, userID)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, sessions, 2) // Only 2 active sessions (expired filtered out)
		// Check that we got the valid sessions
		var foundToken1, foundToken2 bool
		for _, s := range sessions {
			if s.Token == "token1" {
				foundToken1 = true
			}
			if s.Token == "token2" {
				foundToken2 = true
			}
		}
		assert.True(t, foundToken1)
		assert.True(t, foundToken2)
	})
}

func TestSessionStore_Validate(t *testing.T) {
	ctx := context.Background()
	token := testTokenConst
	sessionID := uuid.New()
	userID := uuid.New()

	validSession := &auth.SessionEntity{
		ID:        sessionID,
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("valid session updates last activity", func(t *testing.T) {
		// Setup
		mockRepo := auth.NewMockSessionsRepository(t)
		mockCache := genericcache.NewMemoryCache[auth.SessionEntity]()
		store := stores.NewSessionStore(mockRepo, mockCache, 24*time.Hour)

		updatedSession := &auth.SessionEntity{
			ID:        sessionID,
			UserID:    userID,
			Token:     token,
			ExpiresAt: validSession.ExpiresAt,
			UpdatedAt: time.Now(),
			CreatedAt: validSession.CreatedAt,
		}

		// Expectations
		mockRepo.EXPECT().GetByToken(ctx, token).Return(validSession, nil)
		mockRepo.EXPECT().
			Update(ctx, sessionID, mock.AnythingOfType("*auth.SessionEntity")).
			Return(updatedSession, nil)

		// Execute
		session, err := store.Validate(ctx, token)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, session)
		// UpdatedAt should be recent
		assert.Less(t, time.Since(session.UpdatedAt), 5*time.Second)
	})

	t.Run("expired session returns error", func(t *testing.T) {
		// Setup
		mockRepo := auth.NewMockSessionsRepository(t)
		mockCache := genericcache.NewMemoryCache[auth.SessionEntity]()
		store := stores.NewSessionStore(mockRepo, mockCache, 24*time.Hour)

		expiredSession := &auth.SessionEntity{
			ID:        sessionID,
			UserID:    userID,
			Token:     token,
			ExpiresAt: time.Now().Add(-time.Hour),
			CreatedAt: time.Now().Add(-2 * time.Hour),
			UpdatedAt: time.Now().Add(-2 * time.Hour),
		}

		// Expectations
		mockRepo.EXPECT().GetByToken(ctx, token).Return(expiredSession, nil)
		// Expired session gets deleted
		mockRepo.EXPECT().Delete(ctx, sessionID).Return(nil)

		// Execute
		session, err := store.Validate(ctx, token)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, session)
	})
}

func TestSessionStore_RefreshSession(t *testing.T) {
	ctx := context.Background()
	sessionID := uuid.New()
	userID := uuid.New()
	token := testTokenConst

	validSession := &auth.SessionEntity{
		ID:        sessionID,
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour),
		CreatedAt: time.Now().Add(-time.Hour),
		UpdatedAt: time.Now().Add(-time.Hour),
	}

	t.Run("successful refresh extends expiry", func(t *testing.T) {
		// Setup
		mockRepo := auth.NewMockSessionsRepository(t)
		mockCache := genericcache.NewMemoryCache[auth.SessionEntity]()
		store := stores.NewSessionStore(mockRepo, mockCache, 24*time.Hour)

		// Expectations
		mockRepo.EXPECT().Get(ctx, sessionID).Return(validSession, nil)
		mockRepo.EXPECT().
			Update(ctx, sessionID, mock.AnythingOfType("*auth.SessionEntity")).
			Return(&auth.SessionEntity{
				ID:        sessionID,
				UserID:    userID,
				Token:     token,
				ExpiresAt: time.Now().Add(24 * time.Hour),
				UpdatedAt: time.Now(),
				CreatedAt: validSession.CreatedAt,
			}, nil)

		// Execute
		session, err := store.RefreshSession(ctx, sessionID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, session)

		// Check that expiry was extended (session.ExpiresAt should be about 24 hours from now)
		assert.Greater(t, time.Until(session.ExpiresAt), 23*time.Hour)
	})
}
