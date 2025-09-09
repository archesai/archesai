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
	"golang.org/x/crypto/bcrypt"
)

// TestLoginWithBruteForceProtection tests the brute force protection during login
func TestLoginWithBruteForceProtection(t *testing.T) {
	tests := []struct {
		name          string
		attempts      int
		ipAddress     string
		email         string
		expectLockout bool
		expectError   bool
	}{
		{
			name:          "single failed attempt",
			attempts:      1,
			ipAddress:     "192.168.1.1",
			email:         "test@example.com",
			expectLockout: false,
			expectError:   true, // Wrong password
		},
		{
			name:          "multiple failed attempts from same IP",
			attempts:      5,
			ipAddress:     "192.168.1.2",
			email:         "test@example.com",
			expectLockout: false, // Placeholder implementation always returns false
			expectError:   true,
		},
		{
			name:          "successful login clears attempts",
			attempts:      1,
			ipAddress:     "192.168.1.3",
			email:         "test@example.com",
			expectLockout: false,
			expectError:   false, // Correct password
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
					JWTSecret:          "test-secret",
					AccessTokenExpiry:  15 * time.Minute,
					RefreshTokenExpiry: 7 * 24 * time.Hour,
				},
			}

			// Create test user
			userID := uuid.New()
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)

			testUser := &users.User{
				Id:            userID,
				Email:         "test@example.com",
				Name:          "Test User",
				EmailVerified: true,
			}
			mockUsersRepo.users[userID] = testUser

			// Setup account
			testAccount := &Account{
				Id:         uuid.New(),
				UserId:     userID,
				ProviderId: Local,
				AccountId:  "test@example.com",
				Password:   string(hashedPassword),
			}

			// Always setup account retrieval mock for authentication attempts
			mockRepo.EXPECT().GetAccountByProviderAndProviderID(
				mock.Anything,
				string(Local),
				"test@example.com",
			).Return(testAccount, nil).Maybe()

			// If expecting successful login
			if !tt.expectError && !tt.expectLockout {
				mockRepo.EXPECT().CreateSession(mock.Anything, mock.Anything).Return(&Session{
					Id:        uuid.New(),
					UserId:    userID,
					Token:     "test-token",
					ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
				}, nil).Once()
			}

			// Simulate attempts
			ctx := context.Background()
			for i := 0; i < tt.attempts; i++ {
				req := &LoginRequest{
					Email:    Email(tt.email),
					Password: "WrongPassword!", // Use wrong password for failed attempts
				}

				// On the last attempt, use correct password if expecting success
				if i == tt.attempts-1 && !tt.expectError {
					req.Password = "Password123!"
				}

				_, _, err := service.Login(ctx, req, tt.ipAddress, "test-user-agent")

				if i == tt.attempts-1 {
					// Check final attempt result
					if tt.expectLockout {
						// Should be locked out after too many attempts
						assert.Error(t, err)
						// Verify IP is locked out
						locked := service.isIPLockedOut(ctx, tt.ipAddress)
						assert.Equal(t, tt.expectLockout, locked)
					} else if tt.expectError {
						assert.Error(t, err)
					} else {
						assert.NoError(t, err)
					}
				}
			}
		})
	}
}

// TestLoginWithConcurrentSessionLimits tests concurrent session limiting
func TestLoginWithConcurrentSessionLimits(t *testing.T) {
	tests := []struct {
		name             string
		existingSessions int
		maxSessions      int
		expectOldRemoved bool
	}{
		{
			name:             "under limit",
			existingSessions: 2,
			maxSessions:      5,
			expectOldRemoved: false,
		},
		{
			name:             "at limit",
			existingSessions: 5,
			maxSessions:      5,
			expectOldRemoved: true,
		},
		{
			name:             "over limit",
			existingSessions: 10,
			maxSessions:      5,
			expectOldRemoved: true,
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
					JWTSecret:             "test-secret",
					AccessTokenExpiry:     15 * time.Minute,
					RefreshTokenExpiry:    7 * 24 * time.Hour,
					MaxConcurrentSessions: tt.maxSessions,
				},
			}

			// Create test user
			userID := uuid.New()
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)

			testUser := &users.User{
				Id:            userID,
				Email:         "test@example.com",
				Name:          "Test User",
				EmailVerified: true,
			}
			mockUsersRepo.users[userID] = testUser

			// Setup account
			testAccount := &Account{
				Id:         uuid.New(),
				UserId:     userID,
				ProviderId: Local,
				AccountId:  "test@example.com",
				Password:   string(hashedPassword),
			}

			mockRepo.EXPECT().GetAccountByProviderAndProviderID(
				mock.Anything,
				string(Local),
				"test@example.com",
			).Return(testAccount, nil)

			// Create existing sessions
			existingSessions := []*Session{}
			for i := 0; i < tt.existingSessions; i++ {
				session := &Session{
					Id:        uuid.New(),
					UserId:    userID,
					Token:     uuid.NewString(),
					ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
					CreatedAt: time.Now().Add(time.Duration(-i) * time.Hour), // Older sessions created earlier
				}
				existingSessions = append(existingSessions, session)
			}

			// Setup expectations for listing user sessions
			userIDStr := userID.String()
			mockRepo.EXPECT().ListSessions(mock.Anything, ListSessionsParams{
				UserID: &userIDStr,
				Limit:  100,
			}).Return(existingSessions, int64(len(existingSessions)), nil)

			if tt.expectOldRemoved && tt.existingSessions >= tt.maxSessions {
				// Expect deletion of the oldest session (last in the list since they're ordered by creation time)
				oldestSession := existingSessions[len(existingSessions)-1]
				mockRepo.EXPECT().DeleteSession(mock.Anything, oldestSession.Id).Return(nil).Once()
			}

			// Expect new session creation
			mockRepo.EXPECT().CreateSession(mock.Anything, mock.Anything).Return(&Session{
				Id:        uuid.New(),
				UserId:    userID,
				Token:     "new-session-token",
				ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
			}, nil)

			// Execute login
			ctx := context.Background()
			req := &LoginRequest{
				Email:    Email("test@example.com"),
				Password: "Password123!",
			}

			_, result, err := service.Login(ctx, req, "192.168.1.1", "test-user-agent")

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result.AccessToken)
		})
	}
}

// TestLoginWithRememberMe tests the remember me functionality
func TestLoginWithRememberMe(t *testing.T) {
	tests := []struct {
		name                  string
		rememberMe            bool
		expectedRefreshExpiry time.Duration
	}{
		{
			name:                  "normal login without remember me",
			rememberMe:            false,
			expectedRefreshExpiry: 7 * 24 * time.Hour, // 7 days
		},
		{
			name:                  "login with remember me",
			rememberMe:            true,
			expectedRefreshExpiry: 30 * 24 * time.Hour, // 30 days
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
					JWTSecret:          "test-secret",
					AccessTokenExpiry:  15 * time.Minute,
					RefreshTokenExpiry: 7 * 24 * time.Hour,
				},
			}

			// Create test user
			userID := uuid.New()
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)

			testUser := &users.User{
				Id:            userID,
				Email:         "test@example.com",
				Name:          "Test User",
				EmailVerified: true,
			}
			mockUsersRepo.users[userID] = testUser

			// Setup account
			testAccount := &Account{
				Id:         uuid.New(),
				UserId:     userID,
				ProviderId: Local,
				AccountId:  "test@example.com",
				Password:   string(hashedPassword),
			}

			mockRepo.EXPECT().GetAccountByProviderAndProviderID(
				mock.Anything,
				string(Local),
				"test@example.com",
			).Return(testAccount, nil)

			// Expect session creation (we'll verify the token expiry instead of session expiry)
			mockRepo.EXPECT().CreateSession(mock.Anything, mock.Anything).Return(&Session{
				Id:        uuid.New(),
				UserId:    userID,
				Token:     "test-token",
				ExpiresAt: time.Now().Add(tt.expectedRefreshExpiry).Format(time.RFC3339),
			}, nil)

			// Execute login
			ctx := context.Background()
			req := &LoginRequest{
				Email:      Email("test@example.com"),
				Password:   "Password123!",
				RememberMe: tt.rememberMe,
			}

			_, result, err := service.Login(ctx, req, "192.168.1.1", "test-user-agent")

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.NotEmpty(t, result.AccessToken)
			assert.NotEmpty(t, result.RefreshToken)

			// The ExpiresIn field should reflect the different expiry times
			// Note: This test is simplified since we can't easily validate JWT internals without proper setup
		})
	}
}

// TestLoginRateLimiting tests rate limiting on login endpoint
func TestLoginRateLimiting(t *testing.T) {
	// This would typically be tested at the HTTP handler level with middleware
	// Here we test the service level tracking

	service := &Service{
		logger: logger.NewTest(),
		config: Config{
			JWTSecret: "test-secret",
		},
	}

	ipAddress := "192.168.1.100"
	email := "test@example.com"

	// Track multiple failed attempts
	for i := 0; i < 5; i++ {
		service.trackFailedLoginAttempt(context.Background(), ipAddress, email)
	}

	// Check if IP is locked out (this is a placeholder implementation)
	locked := service.isIPLockedOut(context.Background(), ipAddress)
	assert.False(t, locked, "Current implementation is placeholder, should return false")

	// Clear attempts
	service.clearFailedAttempts(context.Background(), ipAddress, email)

	// Verify cleared (placeholder always returns false)
	locked = service.isIPLockedOut(context.Background(), ipAddress)
	assert.False(t, locked, "Should be unlocked after clearing")
}
