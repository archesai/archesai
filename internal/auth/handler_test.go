package auth

import (
	"context"
	"testing"
	"time"

	"github.com/archesai/archesai/internal/logger"
	"github.com/archesai/archesai/internal/users"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUsersRepository implements users.Repository for testing
type MockUsersRepository struct {
	users map[uuid.UUID]*users.User
	err   error
}

func NewMockUsersRepository() *MockUsersRepository {
	return &MockUsersRepository{
		users: make(map[uuid.UUID]*users.User),
	}
}

func (m *MockUsersRepository) CreateUser(_ context.Context, user *users.User) (*users.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	user.Id = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.users[user.Id] = user
	return user, nil
}

func (m *MockUsersRepository) GetUser(_ context.Context, id uuid.UUID) (*users.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	user, exists := m.users[id]
	if !exists {
		return nil, users.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUsersRepository) GetUserByID(_ context.Context, id uuid.UUID) (*users.User, error) {
	return m.GetUser(context.Background(), id)
}

func (m *MockUsersRepository) GetUserByEmail(_ context.Context, email string) (*users.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, user := range m.users {
		if string(user.Email) == email {
			return user, nil
		}
	}
	return nil, users.ErrUserNotFound
}

func (m *MockUsersRepository) UpdateUser(_ context.Context, id uuid.UUID, user *users.User) (*users.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	if _, exists := m.users[id]; !exists {
		return nil, users.ErrUserNotFound
	}
	user.Id = id
	user.UpdatedAt = time.Now()
	m.users[id] = user
	return user, nil
}

func (m *MockUsersRepository) DeleteUser(_ context.Context, id uuid.UUID) error {
	if m.err != nil {
		return m.err
	}
	if _, exists := m.users[id]; !exists {
		return users.ErrUserNotFound
	}
	delete(m.users, id)
	return nil
}

func (m *MockUsersRepository) ListUsers(_ context.Context, _ users.ListUsersParams) ([]*users.User, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	userList := make([]*users.User, 0, len(m.users))
	for _, user := range m.users {
		userList = append(userList, user)
	}
	return userList, int64(len(userList)), nil
}

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
					Password: "Password123!",
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
			mockRepo := NewMockRepository(t)
			mockUsersRepo := NewMockUsersRepository()

			// Setup mock service
			mockService := &Service{
				repo:      mockRepo,
				usersRepo: mockUsersRepo,
				logger:    logger.NewTest(),
				config: Config{
					JWTSecret:          "test-secret",
					AccessTokenExpiry:  15 * time.Minute,
					RefreshTokenExpiry: 7 * 24 * time.Hour,
				},
			}

			// Setup users repository mock and expectations
			switch tt.name {
			case "user already exists":
				// Pre-populate with existing user
				existingUser := &users.User{
					Id:    uuid.New(),
					Email: "existing@example.com",
					Name:  "Existing User",
				}
				mockUsersRepo.users[existingUser.Id] = existingUser
			case "successful registration":
				// Setup expectations for successful registration
				mockRepo.EXPECT().CreateAccount(mock.Anything, mock.AnythingOfType("*auth.Account")).Return(&Account{
					Id:         uuid.New(),
					UserId:     uuid.New(),
					ProviderId: Local,
					AccountId:  "test@example.com",
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}, nil)

				mockRepo.EXPECT().CreateSession(mock.Anything, mock.AnythingOfType("*auth.Session")).Return(&Session{
					Id:        uuid.New(),
					Token:     "test-token",
					UserId:    uuid.New(),
					ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
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
		setupMocks     func(*testing.T) (*Service, *MockUsersRepository)
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
			setupMocks: func(t *testing.T) (*Service, *MockUsersRepository) {
				mockRepo := NewMockRepository(t)
				mockUsersRepo := NewMockUsersRepository()

				userID := uuid.New()
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)

				// Setup account retrieval by provider and provider ID
				mockRepo.EXPECT().GetAccountByProviderAndProviderID(mock.Anything, string(Local), "test@example.com").Return(&Account{
					Id:         uuid.New(),
					UserId:     userID,
					ProviderId: Local,
					AccountId:  "test@example.com",
					Password:   string(hashedPassword),
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}, nil)

				// Setup session creation
				mockRepo.EXPECT().CreateSession(mock.Anything, mock.Anything).Return(&Session{
					Id:        uuid.New(),
					UserId:    userID,
					Token:     "test-token",
					ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)

				// Setup user in users repo
				user := &users.User{
					Id:    userID,
					Email: "test@example.com",
					Name:  "Test User",
				}
				mockUsersRepo.users[userID] = user

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

				return service, mockUsersRepo
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
			setupMocks: func(t *testing.T) (*Service, *MockUsersRepository) {
				mockRepo := NewMockRepository(t)
				mockUsersRepo := NewMockUsersRepository()

				// The user doesn't exist in the users repository, so GetAccountByProviderAndProviderID should never be called
				// No expectation needed for GetAccountByProviderAndProviderID since the service returns early

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

				return service, mockUsersRepo
			},
			expectedError:  false, // Returns 401 response, not error
			expectedStatus: 401,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			service, _ := tt.setupMocks(t)
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
				mockRepo := NewMockRepository(t)
				mockUsersRepo := NewMockUsersRepository()

				// Setup get session and delete session
				testSession := &Session{
					Id:        uuid.New(),
					UserId:    uuid.New(),
					Token:     "test-token",
					ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
				}
				mockRepo.EXPECT().GetSessionByToken(mock.Anything, "test-token").Return(testSession, nil)
				mockRepo.EXPECT().DeleteSession(mock.Anything, testSession.Id).Return(nil)

				service := &Service{
					repo:      mockRepo,
					usersRepo: mockUsersRepo,
					logger:    logger.NewTest(),
					config: Config{
						JWTSecret: "test-secret",
					},
				}

				return service
			},
			contextToken:  "test-token",
			expectedError: false,
		},
		{
			name:    "session not found",
			request: LogoutRequestObject{},
			setupMocks: func(t *testing.T) *Service {
				mockRepo := NewMockRepository(t)
				mockUsersRepo := NewMockUsersRepository()

				// Setup get session returns not found
				mockRepo.EXPECT().GetSessionByToken(mock.Anything, "non-existent-token").Return(nil, ErrSessionNotFound)

				service := &Service{
					repo:      mockRepo,
					usersRepo: mockUsersRepo,
					logger:    logger.NewTest(),
					config: Config{
						JWTSecret: "test-secret",
					},
				}

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
