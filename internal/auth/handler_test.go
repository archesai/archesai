package auth

import (
	"context"
	"testing"
	"time"

	"github.com/archesai/archesai/internal/accounts"
	"github.com/archesai/archesai/internal/logger"
	"github.com/archesai/archesai/internal/sessions"
	"github.com/archesai/archesai/internal/users"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		request        RegisterRequestObject
		mockUser       *users.User
		mockTokens     *TokenResponse
		mockError      error
		expectedStatus int
	}{
		{
			name: "successful registration",
			request: RegisterRequestObject{
				Body: &RegisterJSONRequestBody{
					Email:    "test@example.com",
					Password: "Password123!",
					Name:     "Test User",
				},
			},
			mockUser: &users.User{
				Id:            uuid.New(),
				Email:         openapi_types.Email("test@example.com"),
				Name:          "Test User",
				EmailVerified: false,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			mockTokens: &TokenResponse{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				TokenType:    "Bearer",
				ExpiresIn:    3600,
			},
			mockError:      nil,
			expectedStatus: 201,
		},
		{
			name: "user already exists",
			request: RegisterRequestObject{
				Body: &RegisterJSONRequestBody{
					Email:    "existing@example.com",
					Password: "Password123!",
					Name:     "Existing User",
				},
			},
			mockUser:       nil,
			mockTokens:     nil,
			mockError:      ErrUserExists,
			expectedStatus: 401, // API spec returns 401 for existing user
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock repositories
			mockAccountsRepo := NewMockAccountsRepository(t)
			mockSessionsRepo := NewMockSessionsRepository(t)
			mockUsersRepo := NewMockUsersRepository(t)
			cache := NewMockSessionsCache(t)

			// Setup cache expectations
			cache.EXPECT().Set(mock.Anything, mock.AnythingOfType("*sessions.Session"), mock.AnythingOfType("time.Duration")).Return(nil).Maybe()
			cache.EXPECT().GetByToken(mock.Anything, mock.AnythingOfType("string")).Return(nil, nil).Maybe()
			cache.EXPECT().Get(mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, nil).Maybe()
			cache.EXPECT().Delete(mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil).Maybe()

			// Setup config
			config := Config{
				JWTSecret:          "test-secret",
				AccessTokenExpiry:  15 * time.Minute,
				RefreshTokenExpiry: 7 * 24 * time.Hour,
				BCryptCost:         10,
			}

			// Create service using NewService
			mockService := NewService(mockAccountsRepo, mockSessionsRepo, mockUsersRepo, cache, config, logger.NewTest())

			// Setup users repository mock and expectations
			switch tt.name {
			case "user already exists":
				// Mock GetByEmail to return existing user
				existingUser := &users.User{
					Id:    uuid.New(),
					Email: openapi_types.Email("existing@example.com"),
					Name:  "Existing User",
				}
				mockUsersRepo.EXPECT().GetByEmail(mock.Anything, "existing@example.com").Return(existingUser, nil)
			case "successful registration":
				// Setup expectations for successful registration
				mockUsersRepo.EXPECT().GetByEmail(mock.Anything, "test@example.com").Return(nil, users.ErrUserNotFound)
				mockUsersRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*users.User")).RunAndReturn(func(_ context.Context, u *users.User) (*users.User, error) {
					u.Id = uuid.New()
					u.CreatedAt = time.Now()
					u.UpdatedAt = time.Now()
					return u, nil
				})

				sessionID := uuid.New()
				mockAccountsRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*accounts.Account")).Return(&accounts.Account{
					Id:         uuid.New(),
					UserId:     uuid.New(),
					ProviderId: accounts.Local,
					AccountId:  "test@example.com",
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}, nil)

				mockSessionsRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*sessions.Session")).Return(&sessions.Session{
					Id:                   sessionID,
					Token:                "test-token",
					UserId:               uuid.New(),
					ActiveOrganizationId: uuid.New(),
					ExpiresAt:            time.Now().Add(time.Hour).Format(time.RFC3339),
					CreatedAt:            time.Now(),
					UpdatedAt:            time.Now(),
				}, nil)

				// Add expectation for UpdateSession which is called to set ActiveOrganizationId
				mockSessionsRepo.EXPECT().Update(mock.Anything, sessionID, mock.AnythingOfType("*sessions.Session")).Return(&sessions.Session{
					Id:                   sessionID,
					Token:                "test-token",
					UserId:               uuid.New(),
					ActiveOrganizationId: uuid.New(),
					ExpiresAt:            time.Now().Add(time.Hour).Format(time.RFC3339),
					CreatedAt:            time.Now(),
					UpdatedAt:            time.Now(),
				}, nil)
			}

			handler := NewHandler(mockService, logger.NewTest())

			// Execute
			ctx := context.Background()
			response, err := handler.Register(ctx, tt.request)

			// Assert
			if tt.mockError == ErrUserExists {
				// Should return 401 response, not an error
				if err != nil {
					t.Errorf("Expected no error for user exists case, got %v", err)
				}
				if _, ok := response.(Register401ApplicationProblemPlusJSONResponse); !ok {
					t.Errorf("Expected Register401ApplicationProblemPlusJSONResponse, got %T", response)
				}
			} else if tt.mockError != nil {
				// Other errors should return as errors
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				// Success case
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if _, ok := response.(Register201JSONResponse); !ok {
					t.Errorf("Expected Register201JSONResponse, got %T", response)
				}
			}
		})
	}
}

func TestHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		request        LoginRequestObject
		setupMocks     func(*testing.T) *Service
		expectedError  bool
		expectedStatus int
	}{
		{
			name: "successful login",
			request: LoginRequestObject{
				Body: &LoginJSONRequestBody{
					Email:    "test@example.com",
					Password: "Password123!",
				},
			},
			setupMocks: func(t *testing.T) *Service {
				mockAccountsRepo := NewMockAccountsRepository(t)
				mockSessionsRepo := NewMockSessionsRepository(t)
				mockUsersRepo := NewMockUsersRepository(t)
				cache := NewMockSessionsCache(t)

				// Setup cache expectations
				cache.EXPECT().Set(mock.Anything, mock.AnythingOfType("*sessions.Session"), mock.AnythingOfType("time.Duration")).Return(nil).Maybe()
				cache.EXPECT().GetByToken(mock.Anything, mock.AnythingOfType("string")).Return(nil, nil).Maybe()
				cache.EXPECT().Get(mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, nil).Maybe()
				cache.EXPECT().Delete(mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil).Maybe()

				userID := uuid.New()
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)

				// Setup account retrieval by provider and provider ID
				mockAccountsRepo.EXPECT().GetByProviderId(mock.Anything, "local", "test@example.com").Return(&accounts.Account{
					Id:         uuid.New(),
					UserId:     userID,
					ProviderId: accounts.Local,
					AccountId:  "test@example.com",
					Password:   string(hashedPassword),
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}, nil)

				sessionID := uuid.New()
				// Setup session creation
				mockSessionsRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(&sessions.Session{
					Id:                   sessionID,
					UserId:               userID,
					Token:                "test-token",
					ActiveOrganizationId: uuid.New(),
					ExpiresAt:            time.Now().Add(time.Hour).Format(time.RFC3339),
					CreatedAt:            time.Now(),
					UpdatedAt:            time.Now(),
				}, nil)

				// Add expectation for UpdateSession which is called to set ActiveOrganizationId
				mockSessionsRepo.EXPECT().Update(mock.Anything, sessionID, mock.AnythingOfType("*sessions.Session")).Return(&sessions.Session{
					Id:                   sessionID,
					Token:                "test-token",
					UserId:               userID,
					ActiveOrganizationId: uuid.New(),
					ExpiresAt:            time.Now().Add(time.Hour).Format(time.RFC3339),
					CreatedAt:            time.Now(),
					UpdatedAt:            time.Now(),
				}, nil)

				// Setup user in users repo
				user := &users.User{
					Id:    userID,
					Email: openapi_types.Email("test@example.com"),
					Name:  "Test User",
				}
				mockUsersRepo.EXPECT().GetByEmail(mock.Anything, "test@example.com").Return(user, nil)

				config := Config{
					JWTSecret:          "test-secret",
					AccessTokenExpiry:  15 * time.Minute,
					RefreshTokenExpiry: 7 * 24 * time.Hour,
					BCryptCost:         10,
				}

				service := NewService(mockAccountsRepo, mockSessionsRepo, mockUsersRepo, cache, config, logger.NewTest())

				return service
			},
			expectedError:  false,
			expectedStatus: 200,
		},
		{
			name: "user not found",
			request: LoginRequestObject{
				Body: &LoginJSONRequestBody{
					Email:    "nonexistent@example.com",
					Password: "Password123!",
				},
			},
			setupMocks: func(t *testing.T) *Service {
				mockAccountsRepo := NewMockAccountsRepository(t)
				mockSessionsRepo := NewMockSessionsRepository(t)
				mockUsersRepo := NewMockUsersRepository(t)
				cache := NewMockSessionsCache(t)

				// Setup cache expectations
				cache.EXPECT().Set(mock.Anything, mock.AnythingOfType("*sessions.Session"), mock.AnythingOfType("time.Duration")).Return(nil).Maybe()
				cache.EXPECT().GetByToken(mock.Anything, mock.AnythingOfType("string")).Return(nil, nil).Maybe()
				cache.EXPECT().Get(mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, nil).Maybe()
				cache.EXPECT().Delete(mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil).Maybe()

				// The user doesn't exist in the users repository
				mockUsersRepo.EXPECT().GetByEmail(mock.Anything, "nonexistent@example.com").Return(nil, users.ErrUserNotFound)

				config := Config{
					JWTSecret:          "test-secret",
					AccessTokenExpiry:  15 * time.Minute,
					RefreshTokenExpiry: 7 * 24 * time.Hour,
					BCryptCost:         10,
				}

				service := NewService(mockAccountsRepo, mockSessionsRepo, mockUsersRepo, cache, config, logger.NewTest())

				return service
			},
			expectedError:  false, // Returns 401 response, not error
			expectedStatus: 401,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			service := tt.setupMocks(t)
			handler := NewHandler(service, logger.NewTest())

			// Execute
			ctx := context.Background()
			response, err := handler.Login(ctx, tt.request)

			// Assert
			if tt.expectedError && err == nil {
				t.Errorf("Expected error, got nil")
			} else if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			switch tt.expectedStatus {
			case 200:
				if _, ok := response.(Login200JSONResponse); !ok && err == nil {
					t.Errorf("Expected Login200JSONResponse, got %T", response)
				}
			case 401:
				if _, ok := response.(Login401ApplicationProblemPlusJSONResponse); !ok && err == nil {
					t.Errorf("Expected Login401ApplicationProblemPlusJSONResponse, got %T", response)
				}
			}
		})
	}
}

func TestHandler_Logout(t *testing.T) {
	tests := []struct {
		name          string
		request       LogoutRequestObject
		setupMocks    func(*testing.T) *Service
		contextToken  string
		expectedError bool
	}{
		{
			name:    "successful logout",
			request: LogoutRequestObject{},
			setupMocks: func(t *testing.T) *Service {
				mockAccountsRepo := NewMockAccountsRepository(t)
				mockSessionsRepo := NewMockSessionsRepository(t)
				mockUsersRepo := NewMockUsersRepository(t)
				cache := NewMockSessionsCache(t)

				// Setup cache expectations
				cache.EXPECT().Set(mock.Anything, mock.AnythingOfType("*sessions.Session"), mock.AnythingOfType("time.Duration")).Return(nil).Maybe()
				cache.EXPECT().GetByToken(mock.Anything, mock.AnythingOfType("string")).Return(nil, nil).Maybe()
				cache.EXPECT().Get(mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, nil).Maybe()
				cache.EXPECT().Delete(mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil).Maybe()

				// Setup get session and delete session
				sessionID := uuid.New()
				testSession := &sessions.Session{
					Id:        sessionID,
					UserId:    uuid.New(),
					Token:     "test-token",
					ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
				}
				mockSessionsRepo.EXPECT().GetByToken(mock.Anything, "test-token").Return(testSession, nil)
				// SessionManager's DeleteSession calls Get first
				mockSessionsRepo.EXPECT().Get(mock.Anything, sessionID).Return(testSession, nil)
				mockSessionsRepo.EXPECT().Delete(mock.Anything, sessionID).Return(nil)

				config := Config{
					JWTSecret:          "test-secret",
					AccessTokenExpiry:  15 * time.Minute,
					RefreshTokenExpiry: 7 * 24 * time.Hour,
					BCryptCost:         10,
				}

				service := NewService(mockAccountsRepo, mockSessionsRepo, mockUsersRepo, cache, config, logger.NewTest())

				return service
			},
			contextToken:  "test-token",
			expectedError: false,
		},
		{
			name:    "session not found",
			request: LogoutRequestObject{},
			setupMocks: func(t *testing.T) *Service {
				mockAccountsRepo := NewMockAccountsRepository(t)
				mockSessionsRepo := NewMockSessionsRepository(t)
				mockUsersRepo := NewMockUsersRepository(t)
				cache := NewMockSessionsCache(t)

				// Setup cache expectations
				cache.EXPECT().Set(mock.Anything, mock.AnythingOfType("*sessions.Session"), mock.AnythingOfType("time.Duration")).Return(nil).Maybe()
				cache.EXPECT().GetByToken(mock.Anything, mock.AnythingOfType("string")).Return(nil, nil).Maybe()
				cache.EXPECT().Get(mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, nil).Maybe()
				cache.EXPECT().Delete(mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil).Maybe()

				// Setup get session returns not found
				mockSessionsRepo.EXPECT().GetByToken(mock.Anything, "non-existent-token").Return(nil, sessions.ErrSessionNotFound)

				config := Config{
					JWTSecret:          "test-secret",
					AccessTokenExpiry:  15 * time.Minute,
					RefreshTokenExpiry: 7 * 24 * time.Hour,
					BCryptCost:         10,
				}

				service := NewService(mockAccountsRepo, mockSessionsRepo, mockUsersRepo, cache, config, logger.NewTest())

				return service
			},
			contextToken:  "non-existent-token",
			expectedError: false, // Returns 401 response
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			service := tt.setupMocks(t)
			handler := NewHandler(service, logger.NewTest())

			// Execute with token in context
			ctx := context.WithValue(context.Background(), SessionTokenContextKey, tt.contextToken)

			response, err := handler.Logout(ctx, tt.request)

			// Assert
			if tt.expectedError && err == nil {
				t.Errorf("Expected error, got nil")
			} else if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check response type
			if err == nil && response != nil {
				switch response.(type) {
				case Logout204Response:
					// Success case
				case Logout401ApplicationProblemPlusJSONResponse:
					// Unauthorized case
				default:
					t.Errorf("Unexpected response type: %T", response)
				}
			}
		})
	}
}
