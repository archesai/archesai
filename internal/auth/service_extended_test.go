package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/archesai/archesai/internal/logger"
	"github.com/archesai/archesai/internal/users"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// TestService_NewServiceWithCache tests service creation with cache
func TestService_NewServiceWithCache(t *testing.T) {
	t.Run("creates service with cache and session manager", func(t *testing.T) {
		mockRepo := NewMockRepository(t)
		mockUsersRepo := NewMockUsersRepository()
		mockCache := NewMockCache(t)
		log := logger.NewTest()
		config := Config{
			JWTSecret:          "test-secret",
			AccessTokenExpiry:  15 * time.Minute,
			RefreshTokenExpiry: 7 * 24 * time.Hour,
		}

		service := NewServiceWithCache(mockRepo, mockUsersRepo, mockCache, config, log)

		assert.NotNil(t, service)
		assert.NotNil(t, service.cache)
		assert.NotNil(t, service.sessionManager)
		assert.Equal(t, mockRepo, service.repo)
		assert.Equal(t, mockUsersRepo, service.usersRepo)
	})
}

// TestService_GetUserSessions tests getting all sessions for a user
func TestService_GetUserSessions(t *testing.T) {
	t.Run("returns user sessions", func(t *testing.T) {
		service, mockRepo, _ := createTestService(t)
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

		result, err := service.GetUserSessions(context.Background(), userID)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, sessions, result)
	})

	t.Run("returns error on repository failure", func(t *testing.T) {
		service, mockRepo, _ := createTestService(t)
		userID := uuid.New()
		userIDStr := userID.String()

		mockRepo.EXPECT().ListSessions(mock.Anything, ListSessionsParams{
			UserID: &userIDStr,
			Limit:  100,
		}).Return(nil, int64(0), errors.New("db error"))

		result, err := service.GetUserSessions(context.Background(), userID)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

// TestService_RevokeSession tests revoking a specific session
func TestService_RevokeSession(t *testing.T) {
	t.Run("successfully revokes session", func(t *testing.T) {
		service, mockRepo, _ := createTestService(t)
		sessionID := uuid.New()

		mockRepo.EXPECT().DeleteSession(mock.Anything, sessionID).Return(nil)

		err := service.RevokeSession(context.Background(), sessionID)

		assert.NoError(t, err)
	})

	t.Run("returns error on repository failure", func(t *testing.T) {
		service, mockRepo, _ := createTestService(t)
		sessionID := uuid.New()

		mockRepo.EXPECT().DeleteSession(mock.Anything, sessionID).Return(errors.New("db error"))

		err := service.RevokeSession(context.Background(), sessionID)

		assert.Error(t, err)
	})
}

// TestService_CleanupExpiredSessions tests cleanup of expired sessions
func TestService_CleanupExpiredSessions(t *testing.T) {
	t.Run("successfully cleans up expired sessions", func(t *testing.T) {
		service, mockRepo, _ := createTestService(t)

		mockRepo.EXPECT().DeleteExpiredSessions(mock.Anything).Return(nil)

		err := service.CleanupExpiredSessions(context.Background())

		assert.NoError(t, err)
	})

	t.Run("returns error on repository failure", func(t *testing.T) {
		service, mockRepo, _ := createTestService(t)

		mockRepo.EXPECT().DeleteExpiredSessions(mock.Anything).Return(errors.New("db error"))

		err := service.CleanupExpiredSessions(context.Background())

		assert.Error(t, err)
	})
}

// TestService_DeleteUserSessions tests deleting all sessions for a user
func TestService_DeleteUserSessions(t *testing.T) {
	t.Run("successfully deletes all user sessions", func(t *testing.T) {
		service, mockRepo, _ := createTestService(t)
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

		service := NewServiceWithCache(mockRepo, mockUsersRepo, mockCache, config, log)
		userID := uuid.New()

		// When using session manager, it should delegate to it
		mockRepo.EXPECT().DeleteUserSessions(mock.Anything, userID).Return(nil)
		mockCache.EXPECT().DeleteUserSessions(mock.Anything, userID).Return(nil).Maybe()

		err := service.DeleteUserSessions(context.Background(), userID)

		assert.NoError(t, err)
	})
}

// TestService_ListUserSessions tests listing user sessions
func TestService_ListUserSessions(t *testing.T) {
	t.Run("without session manager uses repository directly", func(t *testing.T) {
		service, mockRepo, _ := createTestService(t)
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

		service := NewServiceWithCache(mockRepo, mockUsersRepo, mockCache, config, log)
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

		// Session manager will call repository internally
		mockRepo.EXPECT().ListSessions(mock.Anything, ListSessionsParams{
			UserID: &userIDStr,
			Limit:  100,
		}).Return(sessions, int64(1), nil)

		result, err := service.ListUserSessions(context.Background(), userID)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
	})
}

// TestService_ValidateSessionWithoutManager tests session validation without session manager
func TestService_ValidateSessionWithoutManager(t *testing.T) {
	t.Run("validates non-expired session", func(t *testing.T) {
		service, mockRepo, _ := createTestService(t)
		token := testValidToken
		session := &Session{
			Id:        uuid.New(),
			Token:     token,
			UserId:    uuid.New(),
			ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
		}

		mockRepo.EXPECT().GetSessionByToken(mock.Anything, token).Return(session, nil)

		result, err := service.ValidateSession(context.Background(), token)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, token, result.Token)
	})

	t.Run("rejects expired session", func(t *testing.T) {
		service, mockRepo, _ := createTestService(t)
		token := "expired-token"
		session := &Session{
			Id:        uuid.New(),
			Token:     token,
			UserId:    uuid.New(),
			ExpiresAt: time.Now().Add(-time.Hour).Format(time.RFC3339),
		}

		mockRepo.EXPECT().GetSessionByToken(mock.Anything, token).Return(session, nil)

		result, err := service.ValidateSession(context.Background(), token)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

// TestService_RegisterErrors tests error cases in registration
func TestService_RegisterErrors(t *testing.T) {
	t.Run("handles user creation error", func(t *testing.T) {
		service, _, mockUsersRepo := createTestService(t)

		// User doesn't exist but creation fails
		mockUsersRepo.err = errors.New("db error")

		req := RegisterRequest{
			Email:    "test@example.com",
			Password: testPassword,
			Name:     "Test User",
		}

		_, result, err := service.Register(context.Background(), &req)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("handles account creation error", func(t *testing.T) {
		service, mockRepo, mockUsersRepo := createTestService(t)

		// User creation succeeds but account creation fails
		mockUsersRepo.err = nil
		mockRepo.EXPECT().CreateAccount(mock.Anything, mock.AnythingOfType("*auth.Account")).Return(nil, errors.New("db error"))

		req := RegisterRequest{
			Email:    "test@example.com",
			Password: testPassword,
			Name:     "Test User",
		}

		_, result, err := service.Register(context.Background(), &req)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

// TestService_LoginErrors tests error cases in login
func TestService_LoginErrors(t *testing.T) {
	t.Run("handles account not found", func(t *testing.T) {
		service, _, mockUsersRepo := createTestService(t)

		// First call to GetUserByEmail
		mockUsersRepo.err = users.ErrUserNotFound

		req := LoginRequest{
			Email:    openapi_types.Email("test@example.com"),
			Password: testPassword,
		}

		_, result, err := service.Login(context.Background(), &req, "127.0.0.1", "test-agent")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})

	t.Run("handles invalid password", func(t *testing.T) {
		service, mockRepo, mockUsersRepo := createTestService(t)
		userID := uuid.New()
		email := "test@example.com"

		// User exists
		testUser := &users.User{
			Id:    userID,
			Email: openapi_types.Email(email),
			Name:  "Test User",
		}
		mockUsersRepo.users[userID] = testUser
		mockUsersRepo.err = nil

		// Account exists with different password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("different-password"), bcrypt.DefaultCost)
		mockRepo.EXPECT().GetAccountByProviderAndProviderID(mock.Anything, string(Local), email).Return(&Account{
			Id:         uuid.New(),
			UserId:     userID,
			ProviderId: Local,
			AccountId:  email,
			Password:   string(hashedPassword),
		}, nil)

		req := LoginRequest{
			Email:    openapi_types.Email(email),
			Password: testPassword, // Wrong password
		}

		_, result, err := service.Login(context.Background(), &req, "127.0.0.1", "test-agent")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})
}

// TestService_LogoutErrors tests error cases in logout
func TestService_LogoutErrors(t *testing.T) {
	t.Run("handles session not found", func(t *testing.T) {
		service, mockRepo, _ := createTestService(t)
		token := "nonexistent-token"

		mockRepo.EXPECT().GetSessionByToken(mock.Anything, token).Return(nil, ErrSessionNotFound)

		err := service.Logout(context.Background(), token)

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidToken)
	})

	t.Run("handles delete error", func(t *testing.T) {
		service, mockRepo, _ := createTestService(t)
		token := testValidToken
		sessionID := uuid.New()

		mockRepo.EXPECT().GetSessionByToken(mock.Anything, token).Return(&Session{
			Id:    sessionID,
			Token: token,
		}, nil)
		mockRepo.EXPECT().DeleteSession(mock.Anything, sessionID).Return(errors.New("db error"))

		err := service.Logout(context.Background(), token)

		assert.Error(t, err)
	})
}
