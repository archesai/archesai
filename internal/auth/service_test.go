package auth

import (
	"context"
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

const (
	// testPassword is a common test password used across tests
	testPassword = "SecurePass123!"
)

// Note: Using MockUsersRepository from handler_test.go to avoid duplication

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
}

// TestService_RefreshToken tests token refresh
func TestService_RefreshToken(t *testing.T) {
	t.Run("successful refresh", func(t *testing.T) {
		service, _, mockUsersRepo := createTestService(t)

		userID := uuid.New()
		testUser := &users.User{
			Id:    userID,
			Email: "test@example.com",
			Name:  "Test User",
		}

		// Add user to mock repository
		mockUsersRepo.users[userID] = testUser

		// Create a valid refresh token
		claims := &Claims{
			UserID: userID,
			Email:  "test@example.com",
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
			},
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
}
