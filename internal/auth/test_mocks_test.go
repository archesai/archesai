package auth

import (
	"context"
	"time"

	"github.com/archesai/archesai/internal/users"
	"github.com/google/uuid"
)

// MockUsersRepository implements users.Repository for testing
// This manual mock is needed because Go doesn't allow importing test types from other packages
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
	if user.Id == uuid.Nil {
		user.Id = uuid.New()
	}
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

// SetError sets an error to be returned by repository methods
func (m *MockUsersRepository) SetError(err error) {
	m.err = err
}

// AddUser adds a user to the mock repository
func (m *MockUsersRepository) AddUser(user *users.User) {
	if m.users == nil {
		m.users = make(map[uuid.UUID]*users.User)
	}
	m.users[user.Id] = user
}

// Verify interface compliance
var _ users.Repository = (*MockUsersRepository)(nil)

// MockOAuthProvider implements OAuthProvider for testing
type MockOAuthProvider struct {
	providerID     string
	authURL        string
	exchangeTokens *OAuthTokens
	exchangeErr    error
	userInfo       *OAuthUserInfo
	userInfoErr    error
	refreshTokens  *OAuthTokens
	refreshErr     error
}

func (m *MockOAuthProvider) GetProviderID() string {
	return m.providerID
}

func (m *MockOAuthProvider) GetAuthURL(_ string, _ string) string {
	return m.authURL
}

func (m *MockOAuthProvider) ExchangeCode(_ context.Context, _ string, _ string) (*OAuthTokens, error) {
	return m.exchangeTokens, m.exchangeErr
}

func (m *MockOAuthProvider) GetUserInfo(_ context.Context, _ string) (*OAuthUserInfo, error) {
	return m.userInfo, m.userInfoErr
}

func (m *MockOAuthProvider) RefreshToken(_ context.Context, _ string) (*OAuthTokens, error) {
	return m.refreshTokens, m.refreshErr
}

// Test constants
const (
	testPassword   = "SecurePass123!"
	testValidToken = "valid-test-token"
)
