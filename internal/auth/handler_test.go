package auth

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		request        RegisterRequestObject
		mockUser       *User
		mockTokens     *TokenResponse
		mockError      error
		expectedStatus int
	}{
		{
			name: "successful registration",
			request: RegisterRequestObject{
				Body: &RegisterJSONRequestBody{
					Email:    "test@example.com",
					Password: "password123",
					Name:     "Test User",
				},
			},
			mockUser: &User{
				Id:            uuid.New(),
				Email:         "test@example.com",
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
				ExpiresAt:    time.Now().Add(time.Hour),
			},
			mockError:      nil,
			expectedStatus: 201,
		},
		{
			name: "user already exists",
			request: RegisterRequestObject{
				Body: &RegisterJSONRequestBody{
					Email:    "existing@example.com",
					Password: "password123",
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
			// Setup mock service
			mockService := &Service{
				repo:   &MockRepository{},
				logger: slog.Default(),
				config: Config{
					JWTSecret:          "test-secret",
					AccessTokenExpiry:  15 * time.Minute,
					RefreshTokenExpiry: 7 * 24 * time.Hour,
				},
			}

			// Override the Register method behavior
			mockRepo := mockService.repo.(*MockRepository)
			mockRepo.users = make(map[uuid.UUID]*User)
			mockRepo.accounts = make(map[uuid.UUID]*Account)
			mockRepo.sessions = make(map[uuid.UUID]*Session)

			if tt.mockError == ErrUserExists {
				// For user exists case, add a user with same email
				existingUser := &User{
					Id:    uuid.New(),
					Email: tt.request.Body.Email,
					Name:  "Existing User",
				}
				mockRepo.users[existingUser.Id] = existingUser
			}

			handler := NewHandler(mockService, slog.Default())

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
		setupMock      func(*MockRepository)
		expectedError  bool
		expectedStatus int
	}{
		{
			name: "successful login",
			request: LoginRequestObject{
				Body: &LoginJSONRequestBody{
					Email:    "test@example.com",
					Password: "password123",
				},
			},
			setupMock: func(m *MockRepository) {
				userID := uuid.New()
				m.users = map[uuid.UUID]*User{
					userID: {
						Id:            userID,
						Email:         "test@example.com",
						Name:          "Test User",
						EmailVerified: true,
						CreatedAt:     time.Now(),
						UpdatedAt:     time.Now(),
					},
				}
				// Use empty password to skip verification in mock
				m.accounts = map[uuid.UUID]*Account{
					uuid.New(): {
						Id:         uuid.New(),
						UserId:     userID,
						ProviderId: Local,
						AccountId:  "test@example.com",
						Password:   "", // Empty password skips verification
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
					},
				}
				m.sessions = make(map[uuid.UUID]*Session)
			},
			expectedError:  false,
			expectedStatus: 200,
		},
		{
			name: "user not found",
			request: LoginRequestObject{
				Body: &LoginJSONRequestBody{
					Email:    "nonexistent@example.com",
					Password: "password123",
				},
			},
			setupMock: func(m *MockRepository) {
				m.users = make(map[uuid.UUID]*User)
				m.accounts = make(map[uuid.UUID]*Account)
				m.sessions = make(map[uuid.UUID]*Session)
			},
			expectedError:  false, // Returns 401 response, not error
			expectedStatus: 401,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock service
			mockRepo := &MockRepository{}
			tt.setupMock(mockRepo)

			mockService := &Service{
				repo:   mockRepo,
				logger: slog.Default(),
				config: Config{
					JWTSecret:          "test-secret",
					AccessTokenExpiry:  15 * time.Minute,
					RefreshTokenExpiry: 7 * 24 * time.Hour,
				},
			}

			handler := NewHandler(mockService, slog.Default())

			// Execute with context containing IP and user agent
			ctx := context.WithValue(context.Background(), ipAddressKey, "192.168.1.1")
			ctx = context.WithValue(ctx, userAgentKey, "TestAgent/1.0")

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
		setupMock     func(*MockRepository)
		contextToken  string
		expectedError bool
	}{
		{
			name:    "successful logout",
			request: LogoutRequestObject{},
			setupMock: func(m *MockRepository) {
				sessionID := uuid.New()
				m.sessions = map[uuid.UUID]*Session{
					sessionID: {
						Id:        sessionID,
						UserId:    uuid.New(),
						Token:     "test-token",
						ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}
			},
			contextToken:  "test-token",
			expectedError: false,
		},
		{
			name:    "session not found",
			request: LogoutRequestObject{},
			setupMock: func(m *MockRepository) {
				m.sessions = make(map[uuid.UUID]*Session)
			},
			contextToken:  "non-existent-token",
			expectedError: false, // Returns 401 response
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock repository
			mockRepo := &MockRepository{}
			tt.setupMock(mockRepo)

			mockService := &Service{
				repo:   mockRepo,
				logger: slog.Default(),
				config: Config{
					JWTSecret: "test-secret",
				},
			}

			handler := NewHandler(mockService, slog.Default())

			// Execute with token in context
			ctx := context.WithValue(context.Background(), authTokenKey, tt.contextToken)

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

func TestHandler_GetOneUser(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name          string
		request       GetOneUserRequestObject
		setupMock     func(*MockRepository)
		contextUser   *User
		expectedError bool
	}{
		{
			name: "get existing user",
			request: GetOneUserRequestObject{
				Id: userID,
			},
			setupMock: func(m *MockRepository) {
				m.users = map[uuid.UUID]*User{
					userID: {
						Id:            userID,
						Email:         "test@example.com",
						Name:          "Test User",
						EmailVerified: true,
						CreatedAt:     time.Now(),
						UpdatedAt:     time.Now(),
					},
				}
			},
			contextUser: &User{
				Id:    userID,
				Email: "test@example.com",
			},
			expectedError: false,
		},
		{
			name: "user not found",
			request: GetOneUserRequestObject{
				Id: uuid.New(),
			},
			setupMock: func(m *MockRepository) {
				m.users = make(map[uuid.UUID]*User)
			},
			contextUser: &User{
				Id:    userID,
				Email: "test@example.com",
			},
			expectedError: false, // Returns 404 response
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock repository
			mockRepo := &MockRepository{}
			tt.setupMock(mockRepo)

			mockService := &Service{
				repo:   mockRepo,
				logger: slog.Default(),
				config: Config{
					JWTSecret: "test-secret",
				},
			}

			handler := NewHandler(mockService, slog.Default())

			// Execute with user in context
			ctx := context.Background()
			if tt.contextUser != nil {
				ctx = context.WithValue(ctx, UserContextKey, tt.contextUser)
			}

			response, err := handler.GetOneUser(ctx, tt.request)

			// Assert
			if tt.expectedError && err == nil {
				t.Errorf("Expected error, got nil")
			} else if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check response type
			if err == nil && response != nil {
				switch resp := response.(type) {
				case GetOneUser200JSONResponse:
					// Success case
					if resp.Data.Id != userID {
						t.Errorf("Expected user ID %v, got %v", userID, resp.Data.Id)
					}
				case GetOneUser404ApplicationProblemPlusJSONResponse:
					// Not found case - expected
				default:
					t.Errorf("Unexpected response type: %T", response)
				}
			}
		})
	}
}

func TestHandler_UpdateUser(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name          string
		request       UpdateUserRequestObject
		setupMock     func(*MockRepository)
		contextUser   *User
		expectedError bool
	}{
		{
			name: "successful update",
			request: UpdateUserRequestObject{
				Id: userID,
				Body: &UpdateUserJSONRequestBody{
					Email: "newemail@example.com",
					Image: "https://example.com/avatar.jpg",
				},
			},
			setupMock: func(m *MockRepository) {
				m.users = map[uuid.UUID]*User{
					userID: {
						Id:            userID,
						Email:         "test@example.com",
						Name:          "Test User",
						EmailVerified: true,
						CreatedAt:     time.Now(),
						UpdatedAt:     time.Now(),
					},
				}
			},
			contextUser: &User{
				Id:    userID,
				Email: "test@example.com",
			},
			expectedError: false,
		},
		{
			name: "user not found",
			request: UpdateUserRequestObject{
				Id: uuid.New(),
				Body: &UpdateUserJSONRequestBody{
					Email: "newemail@example.com",
				},
			},
			setupMock: func(m *MockRepository) {
				m.users = make(map[uuid.UUID]*User)
			},
			contextUser: &User{
				Id:    userID,
				Email: "test@example.com",
			},
			expectedError: false, // Returns 404 response
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock repository
			mockRepo := &MockRepository{}
			tt.setupMock(mockRepo)

			mockService := &Service{
				repo:   mockRepo,
				logger: slog.Default(),
				config: Config{
					JWTSecret: "test-secret",
				},
			}

			handler := NewHandler(mockService, slog.Default())

			// Execute with user in context
			ctx := context.Background()
			if tt.contextUser != nil {
				ctx = context.WithValue(ctx, UserContextKey, tt.contextUser)
			}

			response, err := handler.UpdateUser(ctx, tt.request)

			// Assert
			if tt.expectedError && err == nil {
				t.Errorf("Expected error, got nil")
			} else if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check response type
			if err == nil && response != nil {
				switch resp := response.(type) {
				case UpdateUser200JSONResponse:
					// Success case
					if resp.Data.Id != userID {
						t.Errorf("Expected user ID %v, got %v", userID, resp.Data.Id)
					}
				case UpdateUser404ApplicationProblemPlusJSONResponse:
					// Not found case - expected
				default:
					t.Errorf("Unexpected response type: %T", response)
				}
			}
		})
	}
}

func TestHandler_FindManyUsers(t *testing.T) {
	tests := []struct {
		name          string
		request       FindManyUsersRequestObject
		setupMock     func(*MockRepository)
		expectedCount int
		expectedError bool
	}{
		{
			name: "list users with pagination",
			request: FindManyUsersRequestObject{
				Params: FindManyUsersParams{},
			},
			setupMock: func(m *MockRepository) {
				// Add some test users
				for i := 0; i < 5; i++ {
					userID := uuid.New()
					m.users[userID] = &User{
						Id:            userID,
						Email:         Email(fmt.Sprintf("user%d@example.com", i)),
						Name:          fmt.Sprintf("User %d", i),
						EmailVerified: i%2 == 0,
						CreatedAt:     time.Now(),
						UpdatedAt:     time.Now(),
					}
				}
			},
			expectedCount: 5,
			expectedError: false,
		},
		{
			name: "empty list",
			request: FindManyUsersRequestObject{
				Params: FindManyUsersParams{},
			},
			setupMock: func(m *MockRepository) {
				m.users = make(map[uuid.UUID]*User)
			},
			expectedCount: 0,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock repository
			mockRepo := &MockRepository{
				users:    make(map[uuid.UUID]*User),
				accounts: make(map[uuid.UUID]*Account),
				sessions: make(map[uuid.UUID]*Session),
			}
			tt.setupMock(mockRepo)

			mockService := &Service{
				repo:   mockRepo,
				logger: slog.Default(),
				config: Config{
					JWTSecret: "test-secret",
				},
			}

			handler := NewHandler(mockService, slog.Default())

			// Execute
			ctx := context.Background()
			response, err := handler.FindManyUsers(ctx, tt.request)

			// Assert
			if tt.expectedError && err == nil {
				t.Errorf("Expected error, got nil")
			} else if !tt.expectedError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Check response
			if resp, ok := response.(FindManyUsers200JSONResponse); ok {
				if len(resp.Data) != tt.expectedCount {
					t.Errorf("Expected %d users, got %d", tt.expectedCount, len(resp.Data))
				}
			} else if err == nil {
				t.Errorf("Expected FindManyUsers200JSONResponse, got %T", response)
			}
		})
	}
}
