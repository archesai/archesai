package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/archesai/archesai/internal/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// createSessionTestService helper function to create a service with mocks for session tests
func createSessionTestService(t *testing.T) (*Service, *MockRepository) {
	t.Helper()

	mockRepo := NewMockRepository(t)
	mockUsersRepo := NewMockUsersRepository()
	log := logger.NewTest()
	config := Config{
		JWTSecret:          "test-secret",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
		BCryptCost:         10,
	}

	// Use NoOpCache for tests
	cache := NewNoOpCache()
	service := NewService(mockRepo, mockUsersRepo, cache, config, log)
	return service, mockRepo
}

// TestService_ValidateSession tests validating a session
func TestService_ValidateSession(t *testing.T) {
	t.Run("valid session", func(t *testing.T) {
		service, mockRepo := createSessionTestService(t)

		token := "valid-token"
		session := &Session{
			Id:                   uuid.New(),
			Token:                token,
			UserId:               uuid.New(),
			ActiveOrganizationId: uuid.New(),
			ExpiresAt:            time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		}

		mockRepo.EXPECT().GetSessionByToken(mock.Anything, token).Return(session, nil)
		// SessionManager's ValidateSession calls UpdateSession to refresh the session
		mockRepo.EXPECT().UpdateSession(mock.Anything, session.Id, mock.AnythingOfType("*auth.Session")).Return(session, nil)

		result, err := service.ValidateSession(context.Background(), token)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, token, result.Token)
	})

	t.Run("expired session", func(t *testing.T) {
		service, mockRepo := createSessionTestService(t)

		token := "expired-token"
		session := &Session{
			Id:        uuid.New(),
			Token:     token,
			UserId:    uuid.New(),
			ExpiresAt: time.Now().Add(-1 * time.Hour).Format(time.RFC3339), // Expired
		}

		mockRepo.EXPECT().GetSessionByToken(mock.Anything, token).Return(session, nil)
		// For expired sessions, SessionManager's GetSessionByToken calls DeleteSession
		mockRepo.EXPECT().DeleteSession(mock.Anything, session.Id).Return(nil)

		result, err := service.ValidateSession(context.Background(), token)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("validates non-expired session", func(t *testing.T) {
		service, mockRepo := createSessionTestService(t)
		token := testValidToken
		session := &Session{
			Id:        uuid.New(),
			Token:     token,
			UserId:    uuid.New(),
			ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
		}

		mockRepo.EXPECT().GetSessionByToken(mock.Anything, token).Return(session, nil)
		// SessionManager's ValidateSession calls UpdateSession to refresh the session
		mockRepo.EXPECT().UpdateSession(mock.Anything, session.Id, mock.AnythingOfType("*auth.Session")).Return(session, nil)

		result, err := service.ValidateSession(context.Background(), token)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, token, result.Token)
	})

	t.Run("rejects expired session", func(t *testing.T) {
		service, mockRepo := createSessionTestService(t)
		token := "expired-token"
		session := &Session{
			Id:        uuid.New(),
			Token:     token,
			UserId:    uuid.New(),
			ExpiresAt: time.Now().Add(-time.Hour).Format(time.RFC3339),
		}

		mockRepo.EXPECT().GetSessionByToken(mock.Anything, token).Return(session, nil)
		// For expired sessions, SessionManager's GetSessionByToken calls DeleteSession
		mockRepo.EXPECT().DeleteSession(mock.Anything, session.Id).Return(nil)

		result, err := service.ValidateSession(context.Background(), token)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

// TestService_NewServiceWithCache tests service creation with cache
func TestService_NewServiceWithCache(t *testing.T) {
	t.Run("creates service with cache and session manager", func(t *testing.T) {
		mockRepo := NewMockRepository(t)
		mockUsersRepo := NewMockUsersRepository()
		log := logger.NewTest()
		config := Config{
			JWTSecret:          "test-secret",
			AccessTokenExpiry:  15 * time.Minute,
			RefreshTokenExpiry: 7 * 24 * time.Hour,
		}

		// Create with cache
		cache := NewNoOpCache()
		service := NewService(mockRepo, mockUsersRepo, cache, config, log)

		assert.NotNil(t, service)
		assert.NotNil(t, service.sessionManager)
		assert.Equal(t, cache, service.cache)
	})
}

// TestService_ListUserSessions_Repository tests getting all sessions for a user using repository directly
func TestService_ListUserSessions_Repository(t *testing.T) {
	t.Run("returns user sessions", func(t *testing.T) {
		service, mockRepo := createSessionTestService(t)
		userID := uuid.New()
		userIDStr := userID.String()

		sessions := []*Session{
			{
				Id:        uuid.New(),
				UserId:    userID,
				Token:     "token1",
				ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
			},
			{
				Id:        uuid.New(),
				UserId:    userID,
				Token:     "token2",
				ExpiresAt: time.Now().Add(2 * time.Hour).Format(time.RFC3339),
			},
		}

		mockRepo.EXPECT().ListSessions(mock.Anything, ListSessionsParams{
			UserID: &userIDStr,
			Limit:  100,
		}).Return(sessions, int64(2), nil)

		result, err := service.ListUserSessions(context.Background(), userID)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, sessions, result)
	})

	t.Run("returns error on repository failure", func(t *testing.T) {
		service, mockRepo := createSessionTestService(t)
		userID := uuid.New()
		userIDStr := userID.String()

		mockRepo.EXPECT().ListSessions(mock.Anything, ListSessionsParams{
			UserID: &userIDStr,
			Limit:  100,
		}).Return(nil, int64(0), errors.New("db error"))

		result, err := service.ListUserSessions(context.Background(), userID)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

// TestService_RevokeSession tests revoking a specific session
func TestService_RevokeSession(t *testing.T) {
	t.Run("successfully revokes session", func(t *testing.T) {
		service, mockRepo := createSessionTestService(t)
		sessionID := uuid.New()

		mockRepo.EXPECT().DeleteSession(mock.Anything, sessionID).Return(nil)

		err := service.RevokeSession(context.Background(), sessionID)

		assert.NoError(t, err)
	})

	t.Run("returns error on repository failure", func(t *testing.T) {
		service, mockRepo := createSessionTestService(t)
		sessionID := uuid.New()

		mockRepo.EXPECT().DeleteSession(mock.Anything, sessionID).Return(errors.New("db error"))

		err := service.RevokeSession(context.Background(), sessionID)

		assert.Error(t, err)
	})
}

// TestService_CleanupExpiredSessions tests cleanup of expired sessions
func TestService_CleanupExpiredSessions(t *testing.T) {
	t.Run("successfully cleans up expired sessions", func(t *testing.T) {
		service, mockRepo := createSessionTestService(t)

		mockRepo.EXPECT().DeleteExpiredSessions(mock.Anything).Return(nil)

		err := service.CleanupExpiredSessions(context.Background())

		assert.NoError(t, err)
	})

	t.Run("returns error on repository failure", func(t *testing.T) {
		service, mockRepo := createSessionTestService(t)

		mockRepo.EXPECT().DeleteExpiredSessions(mock.Anything).Return(errors.New("db error"))

		err := service.CleanupExpiredSessions(context.Background())

		assert.Error(t, err)
	})
}

// TestService_DeleteUserSessions tests deleting all sessions for a user
func TestService_DeleteUserSessions(t *testing.T) {
	t.Run("successfully deletes all user sessions", func(t *testing.T) {
		service, mockRepo := createSessionTestService(t)
		userID := uuid.New()

		mockRepo.EXPECT().DeleteUserSessions(mock.Anything, userID).Return(nil)

		err := service.DeleteUserSessions(context.Background(), userID)

		assert.NoError(t, err)
	})

	t.Run("with session manager delegates to it", func(t *testing.T) {
		mockRepo := NewMockRepository(t)
		mockUsersRepo := NewMockUsersRepository()
		mockCache := NewMockCache(t)
		log := logger.NewTest()
		config := Config{
			JWTSecret:          "test-secret",
			AccessTokenExpiry:  15 * time.Minute,
			RefreshTokenExpiry: 7 * 24 * time.Hour,
		}

		service := NewService(mockRepo, mockUsersRepo, mockCache, config, log)
		userID := uuid.New()

		// Should be called via session manager
		mockRepo.EXPECT().DeleteUserSessions(mock.Anything, userID).Return(nil)
		mockCache.EXPECT().DeleteUserSessions(mock.Anything, userID).Return(nil)

		err := service.DeleteUserSessions(context.Background(), userID)

		assert.NoError(t, err)
	})
}

// TestService_ListUserSessions tests listing user sessions
func TestService_ListUserSessions(t *testing.T) {
	t.Run("without session manager uses repository directly", func(t *testing.T) {
		service, mockRepo := createSessionTestService(t)
		userID := uuid.New()
		userIDStr := userID.String()

		sessions := []*Session{
			{
				Id:        uuid.New(),
				UserId:    userID,
				Token:     "token1",
				ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
			},
		}

		mockRepo.EXPECT().ListSessions(mock.Anything, ListSessionsParams{
			UserID: &userIDStr,
			Limit:  100,
		}).Return(sessions, int64(1), nil)

		result, err := service.ListUserSessions(context.Background(), userID)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
	})
}
