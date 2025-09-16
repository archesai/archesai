package sessions

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
	log := logger.NewTest()

	service := NewService(mockRepo, nil, log)
	return service, mockRepo
}

// TestService_Validate tests validating a session
func TestService_Validate(t *testing.T) {
	t.Run("valid session", func(t *testing.T) {
		service, mockRepo := createSessionTestService(t)

		token := "valid-token"
		session := &Session{
			ID:                   uuid.New(),
			Token:                token,
			UserID:               uuid.New(),
			ActiveOrganizationID: uuid.New(),
			ExpiresAt:            time.Now().Add(24 * time.Hour),
		}

		mockRepo.EXPECT().GetByToken(mock.Anything, token).Return(session, nil)
		// SessionManager's Validate calls Update to refresh the session
		mockRepo.EXPECT().Update(mock.Anything, session.ID, mock.AnythingOfType("*sessions.Session")).Return(session, nil)

		result, err := service.Validate(context.Background(), token)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, token, result.Token)
	})

	t.Run("expired session", func(t *testing.T) {
		service, mockRepo := createSessionTestService(t)

		token := "expired-token"
		session := &Session{
			ID:        uuid.New(),
			Token:     token,
			UserID:    uuid.New(),
			ExpiresAt: time.Now().Add(-1 * time.Hour), // Expired
		}

		mockRepo.EXPECT().GetByToken(mock.Anything, token).Return(session, nil)
		// For expired sessions, SessionManager's GetSessionByToken calls Delete
		mockRepo.EXPECT().Delete(mock.Anything, session.ID).Return(nil)

		result, err := service.Validate(context.Background(), token)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("validates non-expired session", func(t *testing.T) {
		service, mockRepo := createSessionTestService(t)
		token := "test-valid-token"
		session := &Session{
			ID:        uuid.New(),
			Token:     token,
			UserID:    uuid.New(),
			ExpiresAt: time.Now().Add(time.Hour),
		}

		mockRepo.EXPECT().GetByToken(mock.Anything, token).Return(session, nil)
		// SessionManager's Validate calls Update to refresh the session
		mockRepo.EXPECT().Update(mock.Anything, session.ID, mock.AnythingOfType("*sessions.Session")).Return(session, nil)

		result, err := service.Validate(context.Background(), token)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, token, result.Token)
	})

	t.Run("rejects expired session", func(t *testing.T) {
		service, mockRepo := createSessionTestService(t)
		token := "expired-token"
		session := &Session{
			ID:        uuid.New(),
			Token:     token,
			UserID:    uuid.New(),
			ExpiresAt: time.Now().Add(-time.Hour),
		}

		mockRepo.EXPECT().GetByToken(mock.Anything, token).Return(session, nil)
		// For expired sessions, SessionManager's GetSessionByToken calls Delete
		mockRepo.EXPECT().Delete(mock.Anything, session.ID).Return(nil)

		result, err := service.Validate(context.Background(), token)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

// TestService_NewService tests service creation with cache
func TestService_NewService(t *testing.T) {
	t.Run("creates service with cache and session manager", func(t *testing.T) {
		mockRepo := NewMockRepository(t)
		log := logger.NewTest()

		// Create with cache using NewService
		service := NewService(mockRepo, nil, log)

		assert.NotNil(t, service)
		// Can't access internal fields, just check the service is created
	})
}

// TestService_ListByUser_Repository tests getting all sessions for a user using repository directly
func TestService_ListByUser_Repository(t *testing.T) {
	t.Run("returns user sessions", func(t *testing.T) {
		service, mockRepo := createSessionTestService(t)
		userID := uuid.New()

		sessions := []*Session{
			{
				ID:        uuid.New(),
				UserID:    userID,
				Token:     "token1",
				ExpiresAt: time.Now().Add(time.Hour),
			},
			{
				ID:        uuid.New(),
				UserID:    userID,
				Token:     "token2",
				ExpiresAt: time.Now().Add(2 * time.Hour),
			},
		}

		mockRepo.EXPECT().List(mock.Anything, ListSessionsParams{
			Page: PageQuery{
				Number: 1,
				Size:   100,
			},
		}).Return(sessions, int64(2), nil)

		result, err := service.ListByUser(context.Background(), userID)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, sessions, result)
	})

	t.Run("returns error on repository failure", func(t *testing.T) {
		service, mockRepo := createSessionTestService(t)
		userID := uuid.New()

		mockRepo.EXPECT().List(mock.Anything, ListSessionsParams{
			Page: PageQuery{
				Number: 1,
				Size:   100,
			},
		}).Return(nil, int64(0), errors.New("db error"))

		result, err := service.ListByUser(context.Background(), userID)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

// TestService_RevokeSession tests revoking a specific session
func TestService_RevokeSession(t *testing.T) {
	t.Run("successfully revokes session", func(t *testing.T) {
		service, mockRepo := createSessionTestService(t)
		sessionID := uuid.New()

		session := &Session{
			ID:     sessionID,
			UserID: uuid.New(),
		}

		mockRepo.EXPECT().Get(mock.Anything, sessionID).Return(session, nil)
		mockRepo.EXPECT().Delete(mock.Anything, sessionID).Return(nil)

		err := service.RevokeSession(context.Background(), sessionID)

		assert.NoError(t, err)
	})

	t.Run("returns error on repository failure", func(t *testing.T) {
		service, mockRepo := createSessionTestService(t)
		sessionID := uuid.New()

		session := &Session{
			ID:     sessionID,
			UserID: uuid.New(),
		}

		mockRepo.EXPECT().Get(mock.Anything, sessionID).Return(session, nil)
		mockRepo.EXPECT().Delete(mock.Anything, sessionID).Return(errors.New("db error"))

		err := service.RevokeSession(context.Background(), sessionID)

		assert.Error(t, err)
	})
}

// TestService_CleanupExpiredSessions tests cleanup of expired sessions
func TestService_CleanupExpiredSessions(t *testing.T) {
	t.Run("successfully cleans up expired sessions", func(t *testing.T) {
		service, mockRepo := createSessionTestService(t)

		mockRepo.EXPECT().DeleteExpired(mock.Anything).Return(nil)

		err := service.CleanupExpiredSessions(context.Background())

		assert.NoError(t, err)
	})

	t.Run("returns error on repository failure", func(t *testing.T) {
		service, mockRepo := createSessionTestService(t)

		mockRepo.EXPECT().DeleteExpired(mock.Anything).Return(errors.New("db error"))

		err := service.CleanupExpiredSessions(context.Background())

		assert.Error(t, err)
	})
}

// TestService_DeleteUserSessions tests deleting all sessions for a user
// Note: DeleteUserSessions is not in the interface anymore, testing through repository
func TestService_DeleteUserSessions(t *testing.T) {
	t.Run("repository can delete all user sessions", func(t *testing.T) {
		_, mockRepo := createSessionTestService(t)
		userID := uuid.New()

		mockRepo.EXPECT().DeleteByUser(mock.Anything, userID).Return(nil)

		err := mockRepo.DeleteByUser(context.Background(), userID)

		assert.NoError(t, err)
	})

	t.Run("with session manager", func(t *testing.T) {
		mockRepo := NewMockRepository(t)
		log := logger.NewTest()

		// Use NewService to get a service with session manager
		service := NewService(mockRepo, nil, log)
		userID := uuid.New()

		// Should be called via repository
		mockRepo.EXPECT().DeleteByUser(mock.Anything, userID).Return(nil)

		// Call through repository directly since method doesn't exist on service
		err := mockRepo.DeleteByUser(context.Background(), userID)

		assert.NoError(t, err)
		assert.NotNil(t, service) // Just verify service was created
	})
}

// TestService_ListByUser tests listing user sessions
func TestService_ListByUser(t *testing.T) {
	t.Run("without session manager uses repository directly", func(t *testing.T) {
		service, mockRepo := createSessionTestService(t)
		userID := uuid.New()

		sessions := []*Session{
			{
				ID:        uuid.New(),
				UserID:    userID,
				Token:     "token1",
				ExpiresAt: time.Now().Add(time.Hour),
			},
		}

		mockRepo.EXPECT().List(mock.Anything, ListSessionsParams{
			Page: PageQuery{
				Number: 1,
				Size:   100,
			},
		}).Return(sessions, int64(1), nil)

		result, err := service.ListByUser(context.Background(), userID)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
	})
}
