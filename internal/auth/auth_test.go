package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/archesai/archesai/internal/logger"
	"github.com/archesai/archesai/internal/users"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// Test helper function to create a service with mocks
func createTestService(t *testing.T) (*Service, *MockRepository, *MockUsersRepository) {
	t.Helper()

	mockRepo := NewMockRepository(t)
	mockUsersRepo := NewMockUsersRepository()
	log := logger.NewTest()
	config := Config{
		JWTSecret:          "test-secret",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
	}

	service := NewService(mockRepo, mockUsersRepo, config, log)
	return service, mockRepo, mockUsersRepo
}

// TestService_ValidatePassword tests password validation
func TestService_ValidatePassword(t *testing.T) {
	service, _, _ := createTestService(t)

	tests := []struct {
		name        string
		password    string
		wantErr     bool
		errContains string
	}{
		{
			name:     "valid password",
			password: "SecurePass123!",
			wantErr:  false,
		},
		{
			name:        "too short",
			password:    "Pass1!",
			wantErr:     true,
			errContains: "at least 8 characters",
		},
		{
			name:        "too long",
			password:    string(make([]byte, 129)),
			wantErr:     true,
			errContains: "not exceed 128 characters",
		},
		{
			name:        "missing uppercase",
			password:    "securepass123!",
			wantErr:     true,
			errContains: "uppercase letter",
		},
		{
			name:        "missing lowercase",
			password:    "SECUREPASS123!",
			wantErr:     true,
			errContains: "lowercase letter",
		},
		{
			name:        "missing number",
			password:    "SecurePass!",
			wantErr:     true,
			errContains: "number",
		},
		{
			name:        "missing special character",
			password:    "SecurePass123",
			wantErr:     true,
			errContains: "special character",
		},
		{
			name:        "multiple missing requirements",
			password:    "password",
			wantErr:     true,
			errContains: "uppercase letter, number, special character",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validatePassword(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestService_Register tests user registration
func TestService_Register(t *testing.T) {
	email := openapi_types.Email("test@example.com")

	t.Run("successful registration", func(t *testing.T) {
		service, mockRepo, mockUsersRepo := createTestService(t)

		// Setup expectations - MockUsersRepository is a manual mock
		// First call to GetUserByEmail should return not found
		mockUsersRepo.err = nil // Will return ErrUserNotFound since user doesn't exist in map

		mockRepo.EXPECT().CreateAccount(mock.Anything, mock.MatchedBy(func(a *Account) bool {
			return a.ProviderId == Local
		})).Return(&Account{
			Id:         uuid.New(),
			UserId:     uuid.New(),
			ProviderId: Local,
			AccountId:  string(email),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}, nil)

		mockRepo.EXPECT().CreateSession(mock.Anything, mock.AnythingOfType("*auth.Session")).Return(&Session{
			Id:                   uuid.New(),
			Token:                "test-token",
			UserId:               uuid.New(),
			ActiveOrganizationId: uuid.New(),
			ExpiresAt:            time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}, nil)

		// Execute
		req := RegisterRequest{
			Email:    email,
			Password: testPassword,
			Name:     "Test User",
		}

		_, result, err := service.Register(context.Background(), &req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.AccessToken)
		assert.NotEmpty(t, result.RefreshToken)
	})

	t.Run("user already exists", func(t *testing.T) {
		service, _, mockUsersRepo := createTestService(t)

		existingUser := &users.User{
			Id:    uuid.New(),
			Email: email,
			Name:  "Existing User",
		}

		// Add existing user to mock repository
		mockUsersRepo.users[existingUser.Id] = existingUser

		req := RegisterRequest{
			Email:    email,
			Password: testPassword,
			Name:     "Test User",
		}

		_, result, err := service.Register(context.Background(), &req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrUserExists)
	})

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

// TestService_Login tests user login
func TestService_Login(t *testing.T) {
	email := openapi_types.Email("test@example.com")
	userID := uuid.New()

	t.Run("successful login", func(t *testing.T) {
		service, mockRepo, mockUsersRepo := createTestService(t)

		// Create test user with hashed password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
		testUser := &users.User{
			Id:            userID,
			Email:         email,
			Name:          "Test User",
			EmailVerified: true,
		}

		// Add user to mock repository
		mockUsersRepo.users[userID] = testUser

		// Setup expectations
		mockRepo.EXPECT().GetAccountByProviderAndProviderID(mock.Anything, string(Local), string(email)).Return(&Account{
			Id:         uuid.New(),
			UserId:     userID,
			ProviderId: Local,
			AccountId:  string(email),
			Password:   string(hashedPassword),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}, nil)

		mockRepo.EXPECT().CreateSession(mock.Anything, mock.AnythingOfType("*auth.Session")).Return(&Session{
			Id:                   uuid.New(),
			Token:                "test-token",
			UserId:               userID,
			ActiveOrganizationId: uuid.New(),
			ExpiresAt:            time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}, nil)

		// Execute
		req := LoginRequest{
			Email:    email,
			Password: testPassword,
		}

		_, result, err := service.Login(context.Background(), &req, "127.0.0.1", "test-agent")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.AccessToken)
		assert.NotEmpty(t, result.RefreshToken)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		service, mockRepo, mockUsersRepo := createTestService(t)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("different-password"), bcrypt.DefaultCost)
		testUser := &users.User{
			Id:    userID,
			Email: email,
			Name:  "Test User",
		}

		// Add user to mock repository
		mockUsersRepo.users[userID] = testUser

		mockRepo.EXPECT().GetAccountByProviderAndProviderID(mock.Anything, string(Local), string(email)).Return(&Account{
			Id:         uuid.New(),
			UserId:     userID,
			ProviderId: Local,
			AccountId:  string(email),
			Password:   string(hashedPassword),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}, nil)

		req := LoginRequest{
			Email:    email,
			Password: testPassword,
		}

		_, result, err := service.Login(context.Background(), &req, "127.0.0.1", "test-agent")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	})

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

// TestService_RefreshToken tests token refresh
func TestService_RefreshToken(t *testing.T) {
	t.Run("successful refresh", func(t *testing.T) {
		service, mockRepo, mockUsersRepo := createTestService(t)

		userID := uuid.New()
		sessionID := uuid.New()
		testUser := &users.User{
			Id:    userID,
			Email: "test@example.com",
			Name:  "Test User",
		}

		// Add user to mock repository
		mockUsersRepo.users[userID] = testUser

		// Setup mock for GetSession
		session := &Session{
			Id:                   sessionID,
			UserId:               userID,
			Token:                "session-token",
			ActiveOrganizationId: uuid.New(),
			ExpiresAt:            time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			IpAddress:            "192.168.1.1",
			UserAgent:            "Test Agent",
		}
		mockRepo.EXPECT().GetSession(mock.Anything, sessionID).Return(session, nil).Maybe()

		// Create a valid refresh token
		claims := &RefreshClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				Subject:   userID.String(),
			},
			UserID:     userID,
			TokenType:  RefreshTokenType,
			SessionID:  sessionID.String(),
			AuthMethod: AuthMethodPassword,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		refreshToken, _ := token.SignedString([]byte("test-secret"))

		// Execute
		result, err := service.RefreshToken(context.Background(), refreshToken)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.AccessToken)
		assert.NotEmpty(t, result.RefreshToken)
	})
}

// TestService_Logout tests user logout
func TestService_Logout(t *testing.T) {
	t.Run("successful logout", func(t *testing.T) {
		service, mockRepo, _ := createTestService(t)

		token := "logout-token"
		sessionID := uuid.New()
		session := &Session{
			Id:        sessionID,
			UserId:    uuid.New(),
			Token:     token,
			ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
		}

		mockRepo.EXPECT().GetSessionByToken(mock.Anything, token).Return(session, nil)
		mockRepo.EXPECT().DeleteSession(mock.Anything, sessionID).Return(nil)

		err := service.Logout(context.Background(), token)

		assert.NoError(t, err)
	})

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

// TestService_ValidateSession tests validating a session
func TestService_ValidateSession(t *testing.T) {
	t.Run("valid session", func(t *testing.T) {
		service, mockRepo, _ := createTestService(t)

		token := "valid-token"
		session := &Session{
			Id:                   uuid.New(),
			Token:                token,
			UserId:               uuid.New(),
			ActiveOrganizationId: uuid.New(),
			ExpiresAt:            time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		}

		mockRepo.EXPECT().GetSessionByToken(mock.Anything, token).Return(session, nil)

		result, err := service.ValidateSession(context.Background(), token)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, token, result.Token)
	})

	t.Run("expired session", func(t *testing.T) {
		service, mockRepo, _ := createTestService(t)

		token := "expired-token"
		session := &Session{
			Id:        uuid.New(),
			Token:     token,
			UserId:    uuid.New(),
			ExpiresAt: time.Now().Add(-1 * time.Hour).Format(time.RFC3339), // Expired
		}

		mockRepo.EXPECT().GetSessionByToken(mock.Anything, token).Return(session, nil)

		result, err := service.ValidateSession(context.Background(), token)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

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
