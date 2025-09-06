package auth

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	// testPassword is a common test password used across tests
	testPassword = "SecurePass123!"
)

// MockRepository implements ExtendedRepository for testing
type MockRepository struct {
	users    map[uuid.UUID]*User
	sessions map[uuid.UUID]*Session
	accounts map[uuid.UUID]*Account
	err      error // Used to simulate errors
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		users:    make(map[uuid.UUID]*User),
		sessions: make(map[uuid.UUID]*Session),
		accounts: make(map[uuid.UUID]*Account),
	}
}

// User operations
func (m *MockRepository) CreateUser(_ context.Context, entity *User) (*User, error) {
	if m.err != nil {
		return nil, m.err
	}

	// Check for duplicate email
	for _, u := range m.users {
		if u.Email == entity.Email {
			return nil, ErrUserExists
		}
	}

	if entity.Id == uuid.Nil {
		entity.Id = uuid.New()
	}
	entity.CreatedAt = time.Now()
	entity.UpdatedAt = time.Now()
	m.users[entity.Id] = entity
	return entity, nil
}

func (m *MockRepository) GetUser(_ context.Context, id uuid.UUID) (*User, error) {
	if m.err != nil {
		return nil, m.err
	}
	user, exists := m.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (m *MockRepository) UpdateUser(_ context.Context, id uuid.UUID, entity *User) (*User, error) {
	if m.err != nil {
		return nil, m.err
	}
	if _, exists := m.users[id]; !exists {
		return nil, ErrUserNotFound
	}
	entity.Id = id
	entity.UpdatedAt = time.Now()
	m.users[id] = entity
	return entity, nil
}

func (m *MockRepository) DeleteUser(_ context.Context, id uuid.UUID) error {
	if m.err != nil {
		return m.err
	}
	if _, exists := m.users[id]; !exists {
		return ErrUserNotFound
	}
	delete(m.users, id)
	return nil
}

func (m *MockRepository) ListUsers(_ context.Context, params ListUsersParams) ([]*User, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}

	users := make([]*User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}

	// Apply pagination
	start := params.Offset
	end := start + params.Limit
	if start > len(users) {
		return []*User{}, int64(len(users)), nil
	}
	if end > len(users) {
		end = len(users)
	}

	return users[start:end], int64(len(users)), nil
}

// Extended User methods
func (m *MockRepository) GetUserByEmail(_ context.Context, email string) (*User, error) {
	if m.err != nil {
		return nil, m.err
	}
	emailType := openapi_types.Email(email)
	for _, user := range m.users {
		if user.Email == emailType {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

func (m *MockRepository) GetUserByUsername(_ context.Context, username string) (*User, error) {
	if m.err != nil {
		return nil, m.err
	}
	// For testing, use name as username
	for _, user := range m.users {
		if user.Name == username {
			return user, nil
		}
	}
	return nil, ErrUserNotFound
}

// Session operations
func (m *MockRepository) CreateSession(_ context.Context, entity *Session) (*Session, error) {
	if m.err != nil {
		return nil, m.err
	}

	if entity.Id == uuid.Nil {
		entity.Id = uuid.New()
	}
	entity.CreatedAt = time.Now()
	entity.UpdatedAt = time.Now()
	m.sessions[entity.Id] = entity
	return entity, nil
}

func (m *MockRepository) GetSession(_ context.Context, id uuid.UUID) (*Session, error) {
	if m.err != nil {
		return nil, m.err
	}
	session, exists := m.sessions[id]
	if !exists {
		return nil, ErrSessionNotFound
	}
	return session, nil
}

func (m *MockRepository) UpdateSession(_ context.Context, id uuid.UUID, entity *Session) (*Session, error) {
	if m.err != nil {
		return nil, m.err
	}
	if _, exists := m.sessions[id]; !exists {
		return nil, ErrSessionNotFound
	}
	entity.Id = id
	entity.UpdatedAt = time.Now()
	m.sessions[id] = entity
	return entity, nil
}

func (m *MockRepository) DeleteSession(_ context.Context, id uuid.UUID) error {
	if m.err != nil {
		return m.err
	}
	if _, exists := m.sessions[id]; !exists {
		return ErrSessionNotFound
	}
	delete(m.sessions, id)
	return nil
}

func (m *MockRepository) ListSessions(_ context.Context, params ListSessionsParams) ([]*Session, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}

	sessions := make([]*Session, 0, len(m.sessions))
	for _, session := range m.sessions {
		sessions = append(sessions, session)
	}

	// Apply pagination
	start := params.Offset
	end := start + params.Limit
	if start > len(sessions) {
		return []*Session{}, int64(len(sessions)), nil
	}
	if end > len(sessions) {
		end = len(sessions)
	}

	return sessions[start:end], int64(len(sessions)), nil
}

// Extended Session methods
func (m *MockRepository) GetSessionByToken(_ context.Context, token string) (*Session, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, session := range m.sessions {
		if session.Token == token {
			return session, nil
		}
	}
	return nil, ErrSessionNotFound
}

func (m *MockRepository) DeleteExpiredSessions(_ context.Context) error {
	if m.err != nil {
		return m.err
	}
	now := time.Now()
	for id, session := range m.sessions {
		// Parse ExpiresAt string to time
		expiresAt, err := time.Parse(time.RFC3339, session.ExpiresAt)
		if err == nil && expiresAt.Before(now) {
			delete(m.sessions, id)
		}
	}
	return nil
}

func (m *MockRepository) DeleteUserSessions(_ context.Context, userID uuid.UUID) error {
	if m.err != nil {
		return m.err
	}
	userIDStr := userID.String()
	for id, session := range m.sessions {
		if session.UserId == userIDStr {
			delete(m.sessions, id)
		}
	}
	return nil
}

// Account operations
func (m *MockRepository) CreateAccount(_ context.Context, entity *Account) (*Account, error) {
	if m.err != nil {
		return nil, m.err
	}

	if entity.Id == uuid.Nil {
		entity.Id = uuid.New()
	}
	entity.CreatedAt = time.Now()
	entity.UpdatedAt = time.Now()
	m.accounts[entity.Id] = entity
	return entity, nil
}

func (m *MockRepository) GetAccount(_ context.Context, id uuid.UUID) (*Account, error) {
	if m.err != nil {
		return nil, m.err
	}
	account, exists := m.accounts[id]
	if !exists {
		return nil, ErrAccountNotFound
	}
	return account, nil
}

func (m *MockRepository) UpdateAccount(_ context.Context, id uuid.UUID, entity *Account) (*Account, error) {
	if m.err != nil {
		return nil, m.err
	}
	if _, exists := m.accounts[id]; !exists {
		return nil, ErrAccountNotFound
	}
	entity.Id = id
	entity.UpdatedAt = time.Now()
	m.accounts[id] = entity
	return entity, nil
}

func (m *MockRepository) DeleteAccount(_ context.Context, id uuid.UUID) error {
	if m.err != nil {
		return m.err
	}
	if _, exists := m.accounts[id]; !exists {
		return ErrAccountNotFound
	}
	delete(m.accounts, id)
	return nil
}

func (m *MockRepository) ListAccounts(_ context.Context, params ListAccountsParams) ([]*Account, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}

	accounts := make([]*Account, 0, len(m.accounts))
	for _, account := range m.accounts {
		// Filter by UserID if provided
		if params.UserID != nil && account.UserId != *params.UserID {
			continue
		}
		accounts = append(accounts, account)
	}

	// Apply pagination
	start := params.Offset
	end := start + params.Limit
	if start > len(accounts) {
		return []*Account{}, int64(len(accounts)), nil
	}
	if end > len(accounts) {
		end = len(accounts)
	}

	return accounts[start:end], int64(len(accounts)), nil
}

// Extended Account methods
func (m *MockRepository) GetAccountByProviderAndProviderID(_ context.Context, provider, providerID string) (*Account, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, account := range m.accounts {
		if string(account.ProviderId) == provider && account.AccountId == providerID {
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
}

func (m *MockRepository) GetAccountsByUserID(_ context.Context, userID uuid.UUID) ([]*Account, error) {
	if m.err != nil {
		return nil, m.err
	}

	accounts := make([]*Account, 0)
	for _, account := range m.accounts {
		if account.UserId == userID {
			accounts = append(accounts, account)
		}
	}
	return accounts, nil
}

// Test helper functions
func createTestService(t *testing.T) (*Service, *MockRepository) {
	t.Helper()

	mockRepo := NewMockRepository()
	config := Config{
		JWTSecret:          "test-secret-key-for-testing-only",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
		SessionTokenExpiry: 30 * 24 * time.Hour,
		BCryptCost:         4, // Lower cost for faster tests
	}
	logger := slog.Default()

	service := NewService(mockRepo, config, logger)
	return service, mockRepo
}

// TestNewService tests the service constructor
func TestNewService(t *testing.T) {
	service, _ := createTestService(t)

	if service == nil {
		t.Fatal("Expected service to be created")
	}

	if service.repo == nil {
		t.Error("Expected repository to be set")
	}

	if service.logger == nil {
		t.Error("Expected logger to be set")
	}
}

// TestRegister tests user registration
func TestRegister(t *testing.T) {
	tests := []struct {
		name    string
		req     *RegisterRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid registration",
			req: &RegisterRequest{
				Email:    openapi_types.Email("test@example.com"),
				Password: testPassword,
				Name:     "Test User",
			},
			wantErr: false,
		},
		{
			name: "Duplicate email",
			req: &RegisterRequest{
				Email:    openapi_types.Email("duplicate@example.com"),
				Password: testPassword,
				Name:     "Test User",
			},
			wantErr: true,
			errMsg:  "user already exists",
		},
		{
			name: "Empty password",
			req: &RegisterRequest{
				Email:    openapi_types.Email("test2@example.com"),
				Password: "",
				Name:     "Test User",
			},
			wantErr: false, // Will succeed but with empty password hash
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, _ := createTestService(t)

			// Setup duplicate user for duplicate email test
			if tt.name == "Duplicate email" {
				_, _, err := service.Register(context.Background(), tt.req)
				if err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			user, tokens, err := service.Register(context.Background(), tt.req)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				} else if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Error message = %v, want to contain %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("Register() unexpected error: %v", err)
				return
			}

			if user == nil {
				t.Error("Expected user to be returned")
				return
			}

			if user.Email != tt.req.Email {
				t.Errorf("User email = %v, want %v", user.Email, tt.req.Email)
			}

			if tokens == nil || tokens.AccessToken == "" {
				t.Error("Expected tokens to be returned")
			}
		})
	}
}

// TestLogin tests user login
func TestLogin(t *testing.T) {
	tests := []struct {
		name      string
		email     openapi_types.Email
		password  string
		setupUser bool
		wantErr   bool
	}{
		{
			name:      "Valid login",
			email:     openapi_types.Email("test@example.com"),
			password:  testPassword,
			setupUser: true,
			wantErr:   false,
		},
		{
			name:      "Invalid email",
			email:     openapi_types.Email("nonexistent@example.com"),
			password:  testPassword,
			setupUser: false,
			wantErr:   true,
		},
		{
			name:      "Invalid password",
			email:     openapi_types.Email("test@example.com"),
			password:  "WrongPassword",
			setupUser: true,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, _ := createTestService(t)

			// Setup user if needed
			if tt.setupUser {
				_, _, err := service.Register(context.Background(), &RegisterRequest{
					Email:    openapi_types.Email("test@example.com"),
					Password: testPassword,
					Name:     "Test User",
				})
				if err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			req := &LoginRequest{
				Email:    tt.email,
				Password: tt.password,
			}

			user, tokens, err := service.Login(context.Background(), req, "127.0.0.1", "TestAgent")

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Login() unexpected error: %v", err)
				return
			}

			if user == nil {
				t.Error("Expected user to be returned")
				return
			}

			if user.Email != tt.email {
				t.Errorf("User email = %v, want %v", user.Email, tt.email)
			}

			if tokens == nil || tokens.AccessToken == "" {
				t.Error("Expected tokens to be returned")
			}
		})
	}
}

// TestValidateToken tests JWT token validation
func TestValidateToken(t *testing.T) {
	service, _ := createTestService(t)

	// Create a user and get tokens
	user, tokens, err := service.Register(context.Background(), &RegisterRequest{
		Email:    openapi_types.Email("test@example.com"),
		Password: testPassword,
		Name:     "Test User",
	})
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "Valid token",
			token:   tokens.AccessToken,
			wantErr: false,
		},
		{
			name:    "Invalid token",
			token:   "invalid.token.here",
			wantErr: true,
		},
		{
			name:    "Empty token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := service.ValidateToken(tt.token)

			if tt.wantErr {
				if err == nil {
					t.Error("ValidateToken() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateToken() unexpected error: %v", err)
				return
			}

			if claims == nil {
				t.Error("Expected claims to be returned")
				return
			}

			if claims.UserID != user.Id {
				t.Errorf("Claims UserID = %v, want %v", claims.UserID, user.Id)
			}

			if claims.Email != string(user.Email) {
				t.Errorf("Claims Email = %v, want %v", claims.Email, user.Email)
			}
		})
	}
}

// TestHashPassword tests password hashing
func TestHashPassword(t *testing.T) {
	service, _ := createTestService(t)

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid password",
			password: testPassword,
			wantErr:  false,
		},
		{
			name:     "Empty password",
			password: "",
			wantErr:  false, // bcrypt can handle empty strings
		},
		{
			name:     "Very long password",
			password: strings.Repeat("a", 73), // bcrypt max is 72
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := service.hashPassword(tt.password)

			if tt.wantErr {
				if err == nil {
					t.Error("hashPassword() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("hashPassword() unexpected error: %v", err)
				return
			}

			if hash == "" {
				t.Error("Expected hash to be returned")
				return
			}

			// Verify the hash
			err = service.verifyPassword(tt.password, hash)
			if err != nil {
				t.Errorf("Failed to verify hashed password: %v", err)
			}
		})
	}
}

// TestGetUser tests getting a user by ID
func TestGetUser(t *testing.T) {
	service, mockRepo := createTestService(t)

	// Create a test user
	user, _, err := service.Register(context.Background(), &RegisterRequest{
		Email:    openapi_types.Email("test@example.com"),
		Password: testPassword,
		Name:     "Test User",
	})
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	tests := []struct {
		name    string
		userID  uuid.UUID
		wantErr error
	}{
		{
			name:    "Existing user",
			userID:  user.Id,
			wantErr: nil,
		},
		{
			name:    "Non-existent user",
			userID:  uuid.New(),
			wantErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, err := service.GetUser(context.Background(), tt.userID)

			if tt.wantErr != nil {
				if err == nil || !errors.Is(err, tt.wantErr) {
					t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("GetUser() unexpected error: %v", err)
				return
			}

			if gotUser == nil {
				t.Error("Expected user to be returned")
				return
			}

			if gotUser.Id != tt.userID {
				t.Errorf("User ID = %v, want %v", gotUser.Id, tt.userID)
			}
		})
	}

	// Test with repository error
	t.Run("Repository error", func(t *testing.T) {
		mockRepo.err = errors.New("database error")
		defer func() { mockRepo.err = nil }()

		_, err := service.GetUser(context.Background(), user.Id)
		if err == nil {
			t.Error("Expected error when repository fails")
		}
	})
}

// TestUpdateUser tests updating a user
func TestUpdateUser(t *testing.T) {
	service, _ := createTestService(t)

	// Create a test user
	user, _, err := service.Register(context.Background(), &RegisterRequest{
		Email:    openapi_types.Email("test@example.com"),
		Password: testPassword,
		Name:     "Test User",
	})
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	tests := []struct {
		name    string
		userID  uuid.UUID
		req     *UpdateUserRequest
		wantErr bool
	}{
		{
			name:   "Update email",
			userID: user.Id,
			req: &UpdateUserRequest{
				Email: "newemail@example.com",
			},
			wantErr: false,
		},
		{
			name:   "Update image",
			userID: user.Id,
			req: &UpdateUserRequest{
				Image: "https://example.com/avatar.png",
			},
			wantErr: false,
		},
		{
			name:   "Non-existent user",
			userID: uuid.New(),
			req: &UpdateUserRequest{
				Email: "test@example.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedUser, err := service.UpdateUser(context.Background(), tt.userID, tt.req)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateUser() unexpected error: %v", err)
				return
			}

			if updatedUser == nil {
				t.Error("Expected user to be returned")
				return
			}

			if tt.req.Email != "" && string(updatedUser.Email) != tt.req.Email {
				t.Errorf("User email = %v, want %v", updatedUser.Email, tt.req.Email)
			}

			if tt.req.Image != "" && updatedUser.Image != tt.req.Image {
				t.Errorf("User image = %v, want %v", updatedUser.Image, tt.req.Image)
			}
		})
	}
}

// TestDeleteUser tests deleting a user
func TestDeleteUser(t *testing.T) {
	service, mockRepo := createTestService(t)

	// Create a test user
	user, _, err := service.Register(context.Background(), &RegisterRequest{
		Email:    openapi_types.Email("test@example.com"),
		Password: testPassword,
		Name:     "Test User",
	})
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	tests := []struct {
		name    string
		userID  uuid.UUID
		wantErr error
	}{
		{
			name:    "Delete existing user",
			userID:  user.Id,
			wantErr: nil,
		},
		{
			name:    "Delete non-existent user",
			userID:  uuid.New(),
			wantErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteUser(context.Background(), tt.userID)

			if tt.wantErr != nil {
				if err == nil || !errors.Is(err, tt.wantErr) {
					t.Errorf("DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("DeleteUser() unexpected error: %v", err)
				return
			}

			// Verify user was deleted
			_, err = mockRepo.GetUser(context.Background(), tt.userID)
			if !errors.Is(err, ErrUserNotFound) {
				t.Error("User was not deleted from repository")
			}
		})
	}
}

// TestListUsers tests listing users
func TestListUsers(t *testing.T) {
	service, _ := createTestService(t)

	// Create test users
	for i := 0; i < 5; i++ {
		email := openapi_types.Email("test" + string(rune('0'+i)) + "@example.com")
		_, _, err := service.Register(context.Background(), &RegisterRequest{
			Email:    email,
			Password: testPassword,
			Name:     "Test User",
		})
		if err != nil {
			t.Fatalf("Setup failed: %v", err)
		}
	}

	tests := []struct {
		name      string
		limit     int32
		offset    int32
		wantCount int
	}{
		{
			name:      "Get all users",
			limit:     10,
			offset:    0,
			wantCount: 5,
		},
		{
			name:      "Get first 3 users",
			limit:     3,
			offset:    0,
			wantCount: 3,
		},
		{
			name:      "Get users with offset",
			limit:     3,
			offset:    2,
			wantCount: 3,
		},
		{
			name:      "Offset beyond users",
			limit:     10,
			offset:    10,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, err := service.ListUsers(context.Background(), tt.limit, tt.offset)

			if err != nil {
				t.Errorf("ListUsers() unexpected error: %v", err)
				return
			}

			if len(users) != tt.wantCount {
				t.Errorf("ListUsers() returned %d users, want %d", len(users), tt.wantCount)
			}
		})
	}
}

// TestLogout tests user logout
func TestLogout(t *testing.T) {
	service, _ := createTestService(t)

	// Create a user and login
	_, tokens, err := service.Register(context.Background(), &RegisterRequest{
		Email:    openapi_types.Email("test@example.com"),
		Password: testPassword,
		Name:     "Test User",
	})
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "Valid logout",
			token:   tokens.RefreshToken,
			wantErr: false,
		},
		{
			name:    "Invalid token",
			token:   "invalid-token",
			wantErr: true,
		},
		{
			name:    "Empty token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.Logout(context.Background(), tt.token)

			if tt.wantErr {
				if err == nil {
					t.Error("Logout() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Logout() unexpected error: %v", err)
				return
			}
		})
	}
}

// TestCleanupExpiredSessions tests cleaning up expired sessions
func TestCleanupExpiredSessions(t *testing.T) {
	service, mockRepo := createTestService(t)

	// Create a user
	user, _, err := service.Register(context.Background(), &RegisterRequest{
		Email:    openapi_types.Email("test@example.com"),
		Password: testPassword,
		Name:     "Test User",
	})
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// Create expired and valid sessions
	expiredTime := time.Now().Add(-1 * time.Hour)
	validTime := time.Now().Add(1 * time.Hour)

	expiredSession := &Session{
		Token:     "expired-token",
		ExpiresAt: expiredTime.Format(time.RFC3339),
		UserId:    user.Id.String(),
		IpAddress: "127.0.0.1",
		UserAgent: "TestAgent",
	}
	validSession := &Session{
		Token:     "valid-token",
		ExpiresAt: validTime.Format(time.RFC3339),
		UserId:    user.Id.String(),
		IpAddress: "127.0.0.1",
		UserAgent: "TestAgent",
	}

	_, _ = mockRepo.CreateSession(context.Background(), expiredSession)
	_, _ = mockRepo.CreateSession(context.Background(), validSession)

	// Run cleanup
	err = service.CleanupExpiredSessions(context.Background())
	if err != nil {
		t.Errorf("CleanupExpiredSessions() unexpected error: %v", err)
		return
	}

	// Verify expired session was removed
	_, err = mockRepo.GetSessionByToken(context.Background(), "expired-token")
	if !errors.Is(err, ErrSessionNotFound) {
		t.Error("Expired session was not cleaned up")
	}

	// Verify valid session still exists
	session, err := mockRepo.GetSessionByToken(context.Background(), "valid-token")
	if err != nil || session == nil {
		t.Error("Valid session was incorrectly removed")
	}
}

// TestRefreshToken tests token refresh functionality
func TestRefreshToken(t *testing.T) {
	service, _ := createTestService(t)

	// Create a user and get tokens
	_, tokens, err := service.Register(context.Background(), &RegisterRequest{
		Email:    openapi_types.Email("test@example.com"),
		Password: testPassword,
		Name:     "Test User",
	})
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "Valid refresh token",
			token:   tokens.RefreshToken,
			wantErr: false,
		},
		{
			name:    "Invalid refresh token",
			token:   "invalid-refresh-token",
			wantErr: true,
		},
		{
			name:    "Empty refresh token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newTokens, err := service.RefreshToken(context.Background(), tt.token)

			if tt.wantErr {
				if err == nil {
					t.Error("RefreshToken() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("RefreshToken() unexpected error: %v", err)
				return
			}

			if newTokens == nil || newTokens.AccessToken == "" {
				t.Error("Expected new tokens to be returned")
			}
		})
	}
}

// TestVerifyPassword tests password verification
func TestVerifyPassword(t *testing.T) {
	service, _ := createTestService(t)

	password := testPassword
	hash, _ := service.hashPassword(password)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password,
			hash:     hash,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "WrongPassword",
			hash:     hash,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password,
			hash:     "invalid-hash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.verifyPassword(tt.password, tt.hash)

			if tt.wantErr {
				if err == nil {
					t.Error("verifyPassword() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("verifyPassword() unexpected error: %v", err)
			}
		})
	}
}

// TestGenerateTokens tests token generation
func TestGenerateTokens(t *testing.T) {
	service, _ := createTestService(t)

	user := &User{
		Id:    uuid.New(),
		Email: openapi_types.Email("test@example.com"),
		Name:  "Test User",
	}

	tests := []struct {
		name      string
		user      *User
		wantErr   bool
		checkFunc func(*testing.T, *TokenResponse)
	}{
		{
			name:    "Valid user",
			user:    user,
			wantErr: false,
			checkFunc: func(t *testing.T, tokens *TokenResponse) {
				if tokens.AccessToken == "" {
					t.Error("Expected access token")
				}
				if tokens.RefreshToken == "" {
					t.Error("Expected refresh token")
				}
				if tokens.ExpiresIn <= 0 {
					t.Error("Expected positive expiry time")
				}
			},
		},
		{
			name: "User with empty ID",
			user: &User{
				Id:    uuid.Nil,
				Email: openapi_types.Email("test@example.com"),
				Name:  "Test User",
			},
			wantErr: false,
			checkFunc: func(t *testing.T, tokens *TokenResponse) {
				if tokens.AccessToken == "" {
					t.Error("Expected access token even with nil UUID")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := service.generateTokens(tt.user)

			if tt.wantErr {
				if err == nil {
					t.Error("generateTokens() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("generateTokens() unexpected error: %v", err)
				return
			}

			if tokens == nil {
				t.Error("Expected tokens to be returned")
				return
			}

			if tt.checkFunc != nil {
				tt.checkFunc(t, tokens)
			}
		})
	}
}

// TestGetUserByEmail tests fetching user by email through repository
func TestGetUserByEmail(t *testing.T) {
	_, mockRepo := createTestService(t)

	// Create a test user
	user := &User{
		Id:    uuid.New(),
		Email: openapi_types.Email("test@example.com"),
		Name:  "Test User",
	}
	_, _ = mockRepo.CreateUser(context.Background(), user)

	tests := []struct {
		name    string
		email   string
		wantErr error
		setErr  error
	}{
		{
			name:    "Existing user",
			email:   "test@example.com",
			wantErr: nil,
		},
		{
			name:    "Non-existent user",
			email:   "nonexistent@example.com",
			wantErr: ErrUserNotFound,
		},
		{
			name:    "Repository error",
			email:   "test@example.com",
			setErr:  errors.New("database error"),
			wantErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setErr != nil {
				mockRepo.err = tt.setErr
				defer func() { mockRepo.err = nil }()
			}

			gotUser, err := mockRepo.GetUserByEmail(context.Background(), tt.email)

			if tt.wantErr != nil {
				if err == nil {
					t.Error("GetUserByEmail() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("GetUserByEmail() unexpected error: %v", err)
				return
			}

			if gotUser == nil {
				t.Error("Expected user to be returned")
				return
			}

			if string(gotUser.Email) != tt.email {
				t.Errorf("User email = %v, want %v", gotUser.Email, tt.email)
			}

			if gotUser.Id != user.Id {
				t.Errorf("User ID = %v, want %v", gotUser.Id, user.Id)
			}
		})
	}
}

// TestCreateSession tests session creation through repository
func TestCreateSession(t *testing.T) {
	_, mockRepo := createTestService(t)

	userID := uuid.New()

	tests := []struct {
		name      string
		userID    uuid.UUID
		ipAddress string
		userAgent string
		wantErr   bool
		setErr    error
	}{
		{
			name:      "Valid session",
			userID:    userID,
			ipAddress: "192.168.1.1",
			userAgent: "Mozilla/5.0",
			wantErr:   false,
		},
		{
			name:      "Empty IP address",
			userID:    userID,
			ipAddress: "",
			userAgent: "Mozilla/5.0",
			wantErr:   false,
		},
		{
			name:      "Empty user agent",
			userID:    userID,
			ipAddress: "192.168.1.1",
			userAgent: "",
			wantErr:   false,
		},
		{
			name:      "Repository error",
			userID:    userID,
			ipAddress: "192.168.1.1",
			userAgent: "Mozilla/5.0",
			setErr:    errors.New("database error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setErr != nil {
				mockRepo.err = tt.setErr
				defer func() { mockRepo.err = nil }()
			}

			session := &Session{
				UserId:    tt.userID.String(),
				IpAddress: tt.ipAddress,
				UserAgent: tt.userAgent,
				Token:     "test-token-" + uuid.New().String(),
				ExpiresAt: time.Now().Add(1 * time.Hour).Format(time.RFC3339),
			}

			createdSession, err := mockRepo.CreateSession(context.Background(), session)

			if tt.wantErr {
				if err == nil {
					t.Error("CreateSession() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("CreateSession() unexpected error: %v", err)
				return
			}

			if createdSession == nil {
				t.Error("Expected session to be returned")
				return
			}

			if createdSession.UserId != tt.userID.String() {
				t.Errorf("Session UserID = %v, want %v", createdSession.UserId, tt.userID.String())
			}

			if createdSession.IpAddress != tt.ipAddress {
				t.Errorf("Session IpAddress = %v, want %v", createdSession.IpAddress, tt.ipAddress)
			}

			if createdSession.UserAgent != tt.userAgent {
				t.Errorf("Session UserAgent = %v, want %v", createdSession.UserAgent, tt.userAgent)
			}

			if createdSession.Token == "" {
				t.Error("Expected session token to be generated")
			}
		})
	}
}

// TestMockRepositoryEdgeCases tests edge cases in the mock repository
func TestMockRepositoryEdgeCases(t *testing.T) {
	mockRepo := NewMockRepository()

	t.Run("ListUsers with empty repository", func(t *testing.T) {
		users, total, err := mockRepo.ListUsers(context.Background(), ListUsersParams{
			Limit:  10,
			Offset: 0,
		})

		if err != nil {
			t.Errorf("ListUsers() unexpected error: %v", err)
		}

		if len(users) != 0 {
			t.Errorf("Expected 0 users, got %d", len(users))
		}

		if total != 0 {
			t.Errorf("Expected total of 0, got %d", total)
		}
	})

	t.Run("ListSessions with filter", func(t *testing.T) {
		// Create sessions for different users
		userID1 := uuid.New()
		userID2 := uuid.New()

		session1 := &Session{
			Id:        uuid.New(),
			UserId:    userID1.String(),
			Token:     "token1",
			ExpiresAt: time.Now().Add(1 * time.Hour).Format(time.RFC3339),
		}
		session2 := &Session{
			Id:        uuid.New(),
			UserId:    userID2.String(),
			Token:     "token2",
			ExpiresAt: time.Now().Add(1 * time.Hour).Format(time.RFC3339),
		}

		_, _ = mockRepo.CreateSession(context.Background(), session1)
		_, _ = mockRepo.CreateSession(context.Background(), session2)

		// List all sessions
		sessions, total, err := mockRepo.ListSessions(context.Background(), ListSessionsParams{
			Limit:  10,
			Offset: 0,
		})

		if err != nil {
			t.Errorf("ListSessions() unexpected error: %v", err)
		}

		if len(sessions) != 2 {
			t.Errorf("Expected 2 sessions, got %d", len(sessions))
		}

		if total != 2 {
			t.Errorf("Expected total of 2, got %d", total)
		}
	})

	t.Run("GetAccountsByUserID with no accounts", func(t *testing.T) {
		userID := uuid.New()
		accounts, err := mockRepo.GetAccountsByUserID(context.Background(), userID)

		if err != nil {
			t.Errorf("GetAccountsByUserID() unexpected error: %v", err)
		}

		if len(accounts) != 0 {
			t.Errorf("Expected 0 accounts, got %d", len(accounts))
		}
	})

	t.Run("DeleteUserSessions with no sessions", func(t *testing.T) {
		userID := uuid.New()
		err := mockRepo.DeleteUserSessions(context.Background(), userID)

		if err != nil {
			t.Errorf("DeleteUserSessions() unexpected error: %v", err)
		}
	})
}

// TestConcurrentOperations tests concurrent access to the service
func TestConcurrentOperations(t *testing.T) {
	service, _ := createTestService(t)

	// Create a base user
	user, _, err := service.Register(context.Background(), &RegisterRequest{
		Email:    openapi_types.Email("concurrent@example.com"),
		Password: testPassword,
		Name:     "Concurrent User",
	})
	if err != nil {
		t.Fatalf("Setup failed: %v", err)
	}

	// Run concurrent operations
	done := make(chan bool, 3)

	// Concurrent reads
	go func() {
		for i := 0; i < 10; i++ {
			_, _ = service.GetUser(context.Background(), user.Id)
		}
		done <- true
	}()

	// Concurrent updates
	go func() {
		for i := 0; i < 10; i++ {
			req := &UpdateUserRequest{
				Image: "https://example.com/avatar" + string(rune('0'+i)) + ".png",
			}
			_, _ = service.UpdateUser(context.Background(), user.Id, req)
		}
		done <- true
	}()

	// Concurrent token validations
	go func() {
		tokens, _ := service.generateTokens(user)
		for i := 0; i < 10; i++ {
			_, _ = service.ValidateToken(tokens.AccessToken)
		}
		done <- true
	}()

	// Wait for all goroutines
	for i := 0; i < 3; i++ {
		<-done
	}
}

// BenchmarkHashPassword benchmarks password hashing
func BenchmarkHashPassword(b *testing.B) {
	b.Helper()
	mockRepo := NewMockRepository()
	config := Config{
		JWTSecret:          "test-secret",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
		SessionTokenExpiry: 30 * 24 * time.Hour,
		BCryptCost:         4,
	}
	logger := slog.Default()
	service := NewService(mockRepo, config, logger)
	password := testPassword

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.hashPassword(password)
	}
}

// BenchmarkVerifyPassword benchmarks password verification
func BenchmarkVerifyPassword(b *testing.B) {
	b.Helper()
	mockRepo := NewMockRepository()
	config := Config{
		JWTSecret:          "test-secret",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
		SessionTokenExpiry: 30 * 24 * time.Hour,
		BCryptCost:         4,
	}
	logger := slog.Default()
	service := NewService(mockRepo, config, logger)
	password := testPassword
	hash, _ := service.hashPassword(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.verifyPassword(password, hash)
	}
}

// BenchmarkGenerateTokens benchmarks token generation
func BenchmarkGenerateTokens(b *testing.B) {
	b.Helper()
	mockRepo := NewMockRepository()
	config := Config{
		JWTSecret:          "test-secret",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
		SessionTokenExpiry: 30 * 24 * time.Hour,
		BCryptCost:         4,
	}
	logger := slog.Default()
	service := NewService(mockRepo, config, logger)
	user := &User{
		Id:    uuid.New(),
		Email: openapi_types.Email("bench@example.com"),
		Name:  "Benchmark User",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.generateTokens(user)
	}
}

// BenchmarkValidateToken benchmarks token validation
func BenchmarkValidateToken(b *testing.B) {
	b.Helper()
	mockRepo := NewMockRepository()
	config := Config{
		JWTSecret:          "test-secret",
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
		SessionTokenExpiry: 30 * 24 * time.Hour,
		BCryptCost:         4,
	}
	logger := slog.Default()
	service := NewService(mockRepo, config, logger)
	user := &User{
		Id:    uuid.New(),
		Email: openapi_types.Email("bench@example.com"),
		Name:  "Benchmark User",
	}
	tokens, _ := service.generateTokens(user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.ValidateToken(tokens.AccessToken)
	}
}
