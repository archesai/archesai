package auth

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/archesai/archesai/internal/accounts"
	"github.com/archesai/archesai/internal/logger"
	"github.com/archesai/archesai/internal/sessions"
	"github.com/archesai/archesai/internal/users"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

const (
	testPassword   = "SecurePass123!"
	testValidToken = "valid-token-123"
)

// createAuthTestService helper function to create a service with mocks for authentication tests
func createAuthTestService(t *testing.T) (*Service, *MockAccountsRepository, *MockSessionsRepository, *MockUsersRepository) {
	t.Helper()

	mockAccountsRepo := NewMockAccountsRepository(t)
	mockSessionsRepo := NewMockSessionsRepository(t)
	mockUsersRepo := NewMockUsersRepository(t)
	log := logger.NewTest()
	config := Config{
		JWTSecret:          "test-secret",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
		BCryptCost:         10,
	}

	// Use mock cache for tests
	cache := NewMockSessionsCache(t)
	// Setup cache expectations - Set is called when creating sessions
	cache.EXPECT().Set(mock.Anything, mock.AnythingOfType("*sessions.Session"), mock.AnythingOfType("time.Duration")).Return(nil).Maybe()
	cache.EXPECT().GetByToken(mock.Anything, mock.AnythingOfType("string")).Return(nil, nil).Maybe()
	cache.EXPECT().Get(mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, nil).Maybe()
	cache.EXPECT().Delete(mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil).Maybe()
	service := NewService(mockAccountsRepo, mockSessionsRepo, mockUsersRepo, cache, config, log)
	return service, mockAccountsRepo, mockSessionsRepo, mockUsersRepo
}

// TestService_ValidatePassword tests password validation
func TestService_ValidatePassword(t *testing.T) {
	service, _, _, _ := createAuthTestService(t)

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
			password:    strings.Repeat("A", 129) + "aB1!",
			wantErr:     true,
			errContains: "not exceed 128 characters",
		},
		{
			name:        "missing uppercase",
			password:    "password123!",
			wantErr:     true,
			errContains: "uppercase letter",
		},
		{
			name:        "missing lowercase",
			password:    "PASSWORD123!",
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
			errContains: "at least",
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
	name := "Test User"

	t.Run("successful registration", func(t *testing.T) {
		service, mockAccountsRepo, mockSessionsRepo, mockUsersRepo := createAuthTestService(t)

		// Setup expectations
		// First call to GetByEmail should return not found
		mockUsersRepo.EXPECT().GetByEmail(mock.Anything, string(email)).Return(nil, users.ErrUserNotFound).Once()

		// Expect user creation
		createdUserID := uuid.New()
		mockUsersRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(u *users.User) bool {
			return u.Email == email && u.Name == name
		})).RunAndReturn(func(_ context.Context, u *users.User) (*users.User, error) {
			u.Id = createdUserID
			u.CreatedAt = time.Now()
			u.UpdatedAt = time.Now()
			return u, nil
		}).Once()

		mockAccountsRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(a *accounts.Account) bool {
			return a.ProviderId == accounts.Local && a.UserId == createdUserID
		})).Return(&accounts.Account{
			Id:         uuid.New(),
			UserId:     createdUserID,
			ProviderId: accounts.Local,
			AccountId:  string(email),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}, nil)

		sessionID := uuid.New()
		mockSessionsRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*sessions.Session")).Return(&sessions.Session{
			Id:                   sessionID,
			Token:                "test-token",
			UserId:               createdUserID,
			ActiveOrganizationId: uuid.New(),
			ExpiresAt:            time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}, nil)

		// Add expectation for UpdateSession which is called to set ActiveOrganizationId
		mockSessionsRepo.EXPECT().Update(mock.Anything, sessionID, mock.AnythingOfType("*sessions.Session")).Return(&sessions.Session{
			Id:                   sessionID,
			Token:                "test-token",
			UserId:               createdUserID,
			ActiveOrganizationId: uuid.New(),
			ExpiresAt:            time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}, nil)

		req := RegisterRequest{
			Email:    email,
			Password: testPassword,
			Name:     name,
		}

		user, tokenResp, err := service.Register(context.Background(), &req)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotNil(t, tokenResp)
		assert.NotEmpty(t, tokenResp.AccessToken)
		assert.NotEmpty(t, tokenResp.RefreshToken)

		// Verify expectations were met
		mockUsersRepo.AssertExpectations(t)
		mockAccountsRepo.AssertExpectations(t)
		mockSessionsRepo.AssertExpectations(t)
	})

	t.Run("user already exists", func(t *testing.T) {
		service, _, _, mockUsersRepo := createAuthTestService(t)

		existingUser := &users.User{
			Id:    uuid.New(),
			Email: email,
			Name:  "Existing User",
		}

		// Setup mock to return existing user
		mockUsersRepo.EXPECT().GetByEmail(mock.Anything, string(email)).Return(existingUser, nil).Once()

		req := RegisterRequest{
			Email:    email,
			Password: testPassword,
			Name:     "Test User",
		}

		_, result, err := service.Register(context.Background(), &req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrUserExists)
		mockUsersRepo.AssertExpectations(t)
	})

	t.Run("handles user creation error", func(t *testing.T) {
		service, _, _, mockUsersRepo := createAuthTestService(t)

		// User doesn't exist but creation fails
		mockUsersRepo.EXPECT().GetByEmail(mock.Anything, string(email)).Return(nil, users.ErrUserNotFound).Once()

		createErr := errors.New("database error")
		mockUsersRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*users.User")).Return(nil, createErr).Once()

		req := RegisterRequest{
			Email:    email,
			Password: testPassword,
			Name:     name,
		}

		_, result, err := service.Register(context.Background(), &req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to create user")
		mockUsersRepo.AssertExpectations(t)
	})

	t.Run("handles account creation error", func(t *testing.T) {
		service, mockAccountsRepo, _, mockUsersRepo := createAuthTestService(t)

		// User doesn't exist
		mockUsersRepo.EXPECT().GetByEmail(mock.Anything, string(email)).Return(nil, users.ErrUserNotFound).Once()

		// User creation succeeds
		createdUserID := uuid.New()
		mockUsersRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*users.User")).RunAndReturn(func(_ context.Context, u *users.User) (*users.User, error) {
			u.Id = createdUserID
			u.CreatedAt = time.Now()
			u.UpdatedAt = time.Now()
			return u, nil
		}).Once()

		// Account creation fails
		accountErr := errors.New("account creation failed")
		mockAccountsRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*accounts.Account")).Return(nil, accountErr).Once()

		// Expect user deletion due to rollback
		mockUsersRepo.EXPECT().Delete(mock.Anything, createdUserID).Return(nil).Once()

		req := RegisterRequest{
			Email:    email,
			Password: testPassword,
			Name:     name,
		}

		_, result, err := service.Register(context.Background(), &req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to create account")
		mockUsersRepo.AssertExpectations(t)
		mockAccountsRepo.AssertExpectations(t)
	})
}

// TestService_Login tests user login
func TestService_Login(t *testing.T) {
	email := openapi_types.Email("test@example.com")
	userID := uuid.New()

	t.Run("successful login", func(t *testing.T) {
		service, mockAccountsRepo, mockSessionsRepo, mockUsersRepo := createAuthTestService(t)

		// Setup user
		user := &users.User{
			Id:            userID,
			Email:         email,
			Name:          "Test User",
			EmailVerified: true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		mockUsersRepo.EXPECT().GetByEmail(mock.Anything, string(email)).Return(user, nil).Once()

		// Setup account with hashed password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(testPassword), bcrypt.DefaultCost)
		account := &accounts.Account{
			Id:         uuid.New(),
			UserId:     userID,
			ProviderId: accounts.Local,
			AccountId:  string(email),
			Password:   string(hashedPassword),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		mockAccountsRepo.EXPECT().GetByProviderId(mock.Anything, string(accounts.Local), string(email)).Return(account, nil).Once()

		// Setup session creation
		sessionID := uuid.New()
		mockSessionsRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*sessions.Session")).Return(&sessions.Session{
			Id:                   sessionID,
			Token:                "test-token",
			UserId:               userID,
			ActiveOrganizationId: uuid.New(),
			ExpiresAt:            time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}, nil).Once()

		// Add expectation for UpdateSession which is called to set ActiveOrganizationId
		mockSessionsRepo.EXPECT().Update(mock.Anything, sessionID, mock.AnythingOfType("*sessions.Session")).Return(&sessions.Session{
			Id:                   sessionID,
			Token:                "test-token",
			UserId:               userID,
			ActiveOrganizationId: uuid.New(),
			ExpiresAt:            time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}, nil).Once()

		req := LoginRequest{
			Email:    email,
			Password: testPassword,
		}

		loginUser, tokenResp, err := service.Login(context.Background(), &req, "127.0.0.1", "test-agent")

		assert.NoError(t, err)
		assert.NotNil(t, loginUser)
		assert.NotNil(t, tokenResp)
		assert.Equal(t, email, loginUser.Email)
		assert.NotEmpty(t, tokenResp.AccessToken)
		assert.NotEmpty(t, tokenResp.RefreshToken)

		mockUsersRepo.AssertExpectations(t)
		mockAccountsRepo.AssertExpectations(t)
		mockSessionsRepo.AssertExpectations(t)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		service, _, _, mockUsersRepo := createAuthTestService(t)

		// User not found
		mockUsersRepo.EXPECT().GetByEmail(mock.Anything, string(email)).Return(nil, users.ErrUserNotFound).Once()

		req := LoginRequest{
			Email:    email,
			Password: testPassword,
		}

		_, _, err := service.Login(context.Background(), &req, "127.0.0.1", "test-agent")

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidCredentials)
		mockUsersRepo.AssertExpectations(t)
	})

	t.Run("handles account not found", func(t *testing.T) {
		service, mockAccountsRepo, _, mockUsersRepo := createAuthTestService(t)

		// User exists
		user := &users.User{
			Id:    userID,
			Email: email,
			Name:  "Test User",
		}
		mockUsersRepo.EXPECT().GetByEmail(mock.Anything, string(email)).Return(user, nil).Once()

		// Account not found
		mockAccountsRepo.EXPECT().GetByProviderId(mock.Anything, string(accounts.Local), string(email)).Return(nil, errors.New("not found")).Once()

		req := LoginRequest{
			Email:    email,
			Password: testPassword,
		}

		_, _, err := service.Login(context.Background(), &req, "127.0.0.1", "test-agent")

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidCredentials)
		mockUsersRepo.AssertExpectations(t)
		mockAccountsRepo.AssertExpectations(t)
	})

	t.Run("handles invalid password", func(t *testing.T) {
		service, mockAccountsRepo, _, mockUsersRepo := createAuthTestService(t)

		// User exists
		user := &users.User{
			Id:    userID,
			Email: email,
			Name:  "Test User",
		}
		mockUsersRepo.EXPECT().GetByEmail(mock.Anything, string(email)).Return(user, nil).Once()

		// Account with different password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("different_password"), bcrypt.DefaultCost)
		account := &accounts.Account{
			Id:         uuid.New(),
			UserId:     userID,
			ProviderId: accounts.Local,
			AccountId:  string(email),
			Password:   string(hashedPassword),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		mockAccountsRepo.EXPECT().GetByProviderId(mock.Anything, string(accounts.Local), string(email)).Return(account, nil).Once()

		req := LoginRequest{
			Email:    email,
			Password: testPassword,
		}

		_, _, err := service.Login(context.Background(), &req, "127.0.0.1", "test-agent")

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrInvalidCredentials)
		mockUsersRepo.AssertExpectations(t)
		mockAccountsRepo.AssertExpectations(t)
	})
}

// TestService_RefreshToken tests token refresh functionality
func TestService_RefreshToken(t *testing.T) {
	t.Run("successful refresh", func(t *testing.T) {
		// RefreshToken requires parsing JWT tokens which requires integration tests
		t.Skip("RefreshToken requires valid JWT - covered by integration tests")
	})
}

// TestService_Logout tests user logout
func TestService_Logout(t *testing.T) {
	sessionID := uuid.New()
	userID := uuid.New()

	t.Run("successful logout", func(t *testing.T) {
		service, _, mockSessionsRepo, _ := createAuthTestService(t)

		// Setup session
		session := &sessions.Session{
			Id:        sessionID,
			Token:     testValidToken,
			UserId:    userID,
			ExpiresAt: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		// SessionManager.DeleteSessionByToken calls GetSessionByToken, then Delete(id), then internally Delete calls Get(id)
		mockSessionsRepo.EXPECT().GetByToken(mock.Anything, testValidToken).Return(session, nil).Once()
		mockSessionsRepo.EXPECT().Get(mock.Anything, sessionID).Return(session, nil).Once()
		mockSessionsRepo.EXPECT().Delete(mock.Anything, sessionID).Return(nil).Once()

		err := service.Logout(context.Background(), testValidToken)

		assert.NoError(t, err)
		mockSessionsRepo.AssertExpectations(t)
	})

	t.Run("handles session not found", func(t *testing.T) {
		service, _, mockSessionsRepo, _ := createAuthTestService(t)

		mockSessionsRepo.EXPECT().GetByToken(mock.Anything, testValidToken).Return(nil, sessions.ErrSessionNotFound).Once()

		err := service.Logout(context.Background(), testValidToken)

		assert.Error(t, err)
		// Service returns ErrInvalidToken when session is not found
		assert.ErrorIs(t, err, ErrInvalidToken)
		mockSessionsRepo.AssertExpectations(t)
	})

	t.Run("handles delete error", func(t *testing.T) {
		service, _, mockSessionsRepo, _ := createAuthTestService(t)

		// Setup session
		session := &sessions.Session{
			Id:        sessionID,
			Token:     testValidToken,
			UserId:    userID,
			ExpiresAt: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		mockSessionsRepo.EXPECT().GetByToken(mock.Anything, testValidToken).Return(session, nil).Once()
		mockSessionsRepo.EXPECT().Get(mock.Anything, sessionID).Return(session, nil).Once()

		deleteErr := errors.New("delete failed")
		mockSessionsRepo.EXPECT().Delete(mock.Anything, sessionID).Return(deleteErr).Once()

		err := service.Logout(context.Background(), testValidToken)

		assert.Error(t, err)
		// The error is wrapped as ErrInvalidToken by the service
		assert.ErrorIs(t, err, ErrInvalidToken)
		mockSessionsRepo.AssertExpectations(t)
	})
}
