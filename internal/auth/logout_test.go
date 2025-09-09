package auth

import (
	"context"
	"testing"
	"time"

	"github.com/archesai/archesai/internal/logger"
	"github.com/archesai/archesai/internal/users"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const testSessionToken = "test-session-token"

// TestLogoutFlow tests the complete logout flow
func TestLogoutFlow(t *testing.T) {
	tests := []struct {
		name          string
		sessionToken  string
		sessionExists bool
		expectError   bool
		errorMessage  string
	}{
		{
			name:          "successful logout",
			sessionToken:  "valid-session-token",
			sessionExists: true,
			expectError:   false,
		},
		{
			name:          "logout with invalid session",
			sessionToken:  "invalid-session-token",
			sessionExists: false,
			expectError:   true,
			errorMessage:  "invalid token",
		},
		{
			name:         "logout without session token",
			sessionToken: "",
			expectError:  true,
			errorMessage: "invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := NewMockRepository(t)
			mockUsersRepo := NewMockUsersRepository()

			service := &Service{
				repo:      mockRepo,
				usersRepo: mockUsersRepo,
				logger:    logger.NewTest(),
				config: Config{
					JWTSecret: "test-secret",
				},
			}

			ctx := context.Background()

			// Always setup GetSessionByToken mock since Logout always calls it
			if tt.sessionExists {
				// Setup mock for getting session
				testSession := &Session{
					Id:        uuid.New(),
					UserId:    uuid.New(),
					Token:     tt.sessionToken,
					ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
				}
				mockRepo.EXPECT().GetSessionByToken(mock.Anything, tt.sessionToken).Return(testSession, nil)
				mockRepo.EXPECT().DeleteSession(mock.Anything, testSession.Id).Return(nil)
			} else {
				// Session doesn't exist or invalid token
				mockRepo.EXPECT().GetSessionByToken(mock.Anything, tt.sessionToken).Return(nil, ErrSessionNotFound)
			}

			// Execute logout
			err := service.Logout(ctx, tt.sessionToken)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMessage != "" {
					assert.Contains(t, err.Error(), tt.errorMessage)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestDeleteUserSessions tests deleting all user sessions
func TestDeleteUserSessions(t *testing.T) {
	tests := []struct {
		name         string
		userID       uuid.UUID
		sessionCount int
		deleteError  error
		expectError  bool
	}{
		{
			name:         "logout all sessions successfully",
			userID:       uuid.New(),
			sessionCount: 5,
			deleteError:  nil,
			expectError:  false,
		},
		{
			name:         "logout with no sessions",
			userID:       uuid.New(),
			sessionCount: 0,
			deleteError:  nil,
			expectError:  false,
		},
		{
			name:         "logout fails on delete",
			userID:       uuid.New(),
			sessionCount: 3,
			deleteError:  assert.AnError,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := NewMockRepository(t)
			mockUsersRepo := NewMockUsersRepository()

			service := &Service{
				repo:      mockRepo,
				usersRepo: mockUsersRepo,
				logger:    logger.NewTest(),
				config: Config{
					JWTSecret: "test-secret",
				},
			}

			// Setup expectations
			mockRepo.EXPECT().DeleteUserSessions(mock.Anything, tt.userID).Return(tt.deleteError)

			// Execute
			ctx := context.Background()
			err := service.DeleteUserSessions(ctx, tt.userID)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestRevokeSession tests revoking a specific session
func TestRevokeSession(t *testing.T) {
	tests := []struct {
		name        string
		sessionID   uuid.UUID
		deleteError error
		expectError bool
	}{
		{
			name:        "revoke session successfully",
			sessionID:   uuid.New(),
			deleteError: nil,
			expectError: false,
		},
		{
			name:        "revoke non-existent session",
			sessionID:   uuid.New(),
			deleteError: ErrSessionNotFound,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := NewMockRepository(t)
			mockUsersRepo := NewMockUsersRepository()

			service := &Service{
				repo:      mockRepo,
				usersRepo: mockUsersRepo,
				logger:    logger.NewTest(),
				config: Config{
					JWTSecret: "test-secret",
				},
			}

			// Setup expectations
			mockRepo.EXPECT().DeleteSession(mock.Anything, tt.sessionID).Return(tt.deleteError)

			// Execute
			ctx := context.Background()
			err := service.RevokeSession(ctx, tt.sessionID)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestLogoutWithSessionManager tests logout when using session manager
func TestLogoutWithSessionManager(t *testing.T) {
	// This test would require a mock SessionManager
	// For now, we test the basic flow without SessionManager

	mockRepo := NewMockRepository(t)
	mockUsersRepo := NewMockUsersRepository()

	service := &Service{
		repo:      mockRepo,
		usersRepo: mockUsersRepo,
		logger:    logger.NewTest(),
		config: Config{
			JWTSecret: "test-secret",
		},
		// sessionManager is nil, so it will use the repository directly
	}

	sessionToken := testSessionToken
	ctx := context.WithValue(context.Background(), SessionTokenContextKey, sessionToken)

	// Setup mock expectations
	testSession := &Session{
		Id:        uuid.New(),
		UserId:    uuid.New(),
		Token:     sessionToken,
		ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
	}

	mockRepo.EXPECT().GetSessionByToken(mock.Anything, sessionToken).Return(testSession, nil)
	mockRepo.EXPECT().DeleteSession(mock.Anything, testSession.Id).Return(nil)

	// Execute
	err := service.Logout(ctx, sessionToken)

	// Assert
	assert.NoError(t, err)
}

// TestLogoutIdempotency tests that logout can be called multiple times safely
func TestLogoutIdempotency(t *testing.T) {
	mockRepo := NewMockRepository(t)
	mockUsersRepo := NewMockUsersRepository()

	service := &Service{
		repo:      mockRepo,
		usersRepo: mockUsersRepo,
		logger:    logger.NewTest(),
		config: Config{
			JWTSecret: "test-secret",
		},
	}

	sessionToken := testSessionToken
	ctx := context.WithValue(context.Background(), SessionTokenContextKey, sessionToken)

	// First logout - session exists
	testSession := &Session{
		Id:        uuid.New(),
		UserId:    uuid.New(),
		Token:     sessionToken,
		ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
	}

	mockRepo.EXPECT().GetSessionByToken(mock.Anything, sessionToken).Return(testSession, nil).Once()
	mockRepo.EXPECT().DeleteSession(mock.Anything, testSession.Id).Return(nil).Once()

	err := service.Logout(ctx, sessionToken)
	assert.NoError(t, err)

	// Second logout - session no longer exists
	mockRepo.EXPECT().GetSessionByToken(mock.Anything, sessionToken).Return(nil, ErrSessionNotFound).Once()

	err = service.Logout(ctx, sessionToken)
	assert.Error(t, err)
	// The service returns ErrInvalidToken when session is not found
	assert.Equal(t, ErrInvalidToken, err)
}

// TestLogoutCleansUpRelatedData tests that logout properly cleans up related data
func TestLogoutCleansUpRelatedData(t *testing.T) {
	mockRepo := NewMockRepository(t)
	mockUsersRepo := NewMockUsersRepository()

	service := &Service{
		repo:      mockRepo,
		usersRepo: mockUsersRepo,
		logger:    logger.NewTest(),
		config: Config{
			JWTSecret: "test-secret",
		},
	}

	userID := uuid.New()
	sessionToken := testSessionToken
	ctx := context.WithValue(context.Background(), SessionTokenContextKey, sessionToken)

	// Setup test user
	testUser := &users.User{
		Id:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}
	mockUsersRepo.users[userID] = testUser

	// Setup session
	testSession := &Session{
		Id:        uuid.New(),
		UserId:    userID,
		Token:     sessionToken,
		ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Expect session retrieval and deletion
	mockRepo.EXPECT().GetSessionByToken(mock.Anything, sessionToken).Return(testSession, nil)
	mockRepo.EXPECT().DeleteSession(mock.Anything, testSession.Id).Return(nil)

	// Execute logout
	err := service.Logout(ctx, sessionToken)
	assert.NoError(t, err)

	// Verify session is deleted (would be checked by expecting the DeleteSession call)
	// In a real implementation, we might also verify:
	// - Cache is cleared
	// - Active session count is decremented
	// - Audit log is created
	// - Events are published
}
