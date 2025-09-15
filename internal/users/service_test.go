package users

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test helper functions
func createTestService(t *testing.T) (*Service, *MockRepository, *MockCache, *MockEventPublisher) {
	t.Helper()

	mockRepo := NewMockRepository(t)
	mockCache := NewMockCache(t)
	mockEvents := NewMockEventPublisher(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	service := NewService(mockRepo, mockCache, mockEvents, logger)
	return service, mockRepo, mockCache, mockEvents
}

// TestNewService tests the service constructor
func TestNewService(t *testing.T) {
	service, _, _, _ := createTestService(t)

	assert.NotNil(t, service)
	assert.NotNil(t, service.repo)
	assert.NotNil(t, service.logger)
}

// TestGetUser tests getting a user by ID
func TestGetUser(t *testing.T) {
	// Create a test user
	user := &User{
		Id:            uuid.New(),
		Email:         openapi_types.Email("test@example.com"),
		Name:          "Test User",
		EmailVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	tests := []struct {
		name       string
		userID     uuid.UUID
		setupMocks func(*MockRepository, *MockCache)
		wantErr    error
	}{
		{
			name:   "Existing user - cache hit",
			userID: user.Id,
			setupMocks: func(_ *MockRepository, cache *MockCache) {
				cache.EXPECT().Get(mock.Anything, user.Id).Return(user, nil)
			},
			wantErr: nil,
		},
		{
			name:   "Existing user - cache miss",
			userID: user.Id,
			setupMocks: func(repo *MockRepository, cache *MockCache) {
				cache.EXPECT().Get(mock.Anything, user.Id).Return(nil, ErrCacheMiss)
				repo.EXPECT().Get(mock.Anything, user.Id).Return(user, nil)
				cache.EXPECT().Set(mock.Anything, user, mock.Anything).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:   "Non-existent user",
			userID: uuid.New(),
			setupMocks: func(repo *MockRepository, cache *MockCache) {
				cache.EXPECT().Get(mock.Anything, mock.Anything).Return(nil, ErrCacheMiss)
				repo.EXPECT().Get(mock.Anything, mock.Anything).Return(nil, ErrUserNotFound)
			},
			wantErr: ErrUserNotFound,
		},
		{
			name:   "Repository error",
			userID: user.Id,
			setupMocks: func(repo *MockRepository, cache *MockCache) {
				cache.EXPECT().Get(mock.Anything, user.Id).Return(nil, ErrCacheMiss)
				repo.EXPECT().Get(mock.Anything, user.Id).Return(nil, errors.New("database error"))
			},
			wantErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo, mockCache, _ := createTestService(t)
			tt.setupMocks(mockRepo, mockCache)

			gotUser, err := service.Get(context.Background(), tt.userID)

			if tt.wantErr != nil {
				assert.Error(t, err)
				if tt.wantErr.Error() != "database error" {
					assert.ErrorIs(t, err, tt.wantErr)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, gotUser)
				assert.Equal(t, tt.userID, gotUser.Id)
			}
		})
	}
}

// TestUpdateUser tests updating a user
func TestUpdateUser(t *testing.T) {
	// Create a test user
	user := &User{
		Id:            uuid.New(),
		Email:         openapi_types.Email("test@example.com"),
		Name:          "Test User",
		EmailVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	tests := []struct {
		name       string
		userID     uuid.UUID
		req        *UpdateUserJSONBody
		setupMocks func(*MockRepository, *MockCache, *MockEventPublisher)
		wantErr    bool
	}{
		{
			name:   "Update email",
			userID: user.Id,
			req: &UpdateUserJSONBody{
				Email: "newemail@example.com",
			},
			setupMocks: func(repo *MockRepository, cache *MockCache, events *MockEventPublisher) {
				repo.EXPECT().Get(mock.Anything, user.Id).Return(user, nil)
				updatedUser := &User{
					Id:            user.Id,
					Email:         openapi_types.Email("newemail@example.com"),
					Name:          user.Name,
					EmailVerified: user.EmailVerified,
					CreatedAt:     user.CreatedAt,
					UpdatedAt:     time.Now(),
				}
				repo.EXPECT().Update(mock.Anything, user.Id, mock.MatchedBy(func(u *User) bool {
					return string(u.Email) == "newemail@example.com"
				})).Return(updatedUser, nil)
				cache.EXPECT().Set(mock.Anything, mock.Anything, mock.Anything).Return(nil)
				events.EXPECT().PublishUserUpdated(mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "Non-existent user",
			userID: uuid.New(),
			req: &UpdateUserJSONBody{
				Email: "test@example.com",
			},
			setupMocks: func(repo *MockRepository, _ *MockCache, _ *MockEventPublisher) {
				repo.EXPECT().Get(mock.Anything, mock.Anything).Return(nil, ErrUserNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo, mockCache, mockEvents := createTestService(t)
			tt.setupMocks(mockRepo, mockCache, mockEvents)

			updatedUser, err := service.Update(context.Background(), tt.userID, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, updatedUser)
				if tt.req.Email != "" {
					assert.Equal(t, tt.req.Email, string(updatedUser.Email))
				}
			}
		})
	}
}

// TestDeleteUser tests deleting a user
func TestDeleteUser(t *testing.T) {
	// Create a test user
	user := &User{
		Id:            uuid.New(),
		Email:         openapi_types.Email("test@example.com"),
		Name:          "Test User",
		EmailVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	tests := []struct {
		name       string
		userID     uuid.UUID
		setupMocks func(*MockRepository, *MockCache, *MockEventPublisher)
		wantErr    error
	}{
		{
			name:   "Delete existing user",
			userID: user.Id,
			setupMocks: func(repo *MockRepository, cache *MockCache, events *MockEventPublisher) {
				repo.EXPECT().Get(mock.Anything, user.Id).Return(user, nil)
				repo.EXPECT().Delete(mock.Anything, user.Id).Return(nil)
				cache.EXPECT().Delete(mock.Anything, user.Id).Return(nil)
				events.EXPECT().PublishUserDeleted(mock.Anything, user).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:   "Delete non-existent user",
			userID: uuid.New(),
			setupMocks: func(repo *MockRepository, _ *MockCache, _ *MockEventPublisher) {
				repo.EXPECT().Get(mock.Anything, mock.Anything).Return(nil, ErrUserNotFound)
			},
			wantErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo, mockCache, mockEvents := createTestService(t)
			tt.setupMocks(mockRepo, mockCache, mockEvents)

			err := service.Delete(context.Background(), tt.userID)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestListUsers tests listing users
func TestListUsers(t *testing.T) {
	// Create test users
	users := []*User{
		{
			Id:            uuid.New(),
			Email:         openapi_types.Email("test1@example.com"),
			Name:          "Test User 1",
			EmailVerified: false,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			Id:            uuid.New(),
			Email:         openapi_types.Email("test2@example.com"),
			Name:          "Test User 2",
			EmailVerified: false,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	tests := []struct {
		name       string
		limit      int32
		offset     int32
		setupMocks func(*MockRepository)
		wantCount  int
	}{
		{
			name:   "Get all users",
			limit:  10,
			offset: 0,
			setupMocks: func(repo *MockRepository) {
				repo.EXPECT().List(mock.Anything, ListUsersParams{
					Page: PageQuery{
						Number: 1,
						Size:   10,
					},
				}).Return(users, int64(2), nil)
			},
			wantCount: 2,
		},
		{
			name:   "Get first user",
			limit:  1,
			offset: 0,
			setupMocks: func(repo *MockRepository) {
				repo.EXPECT().List(mock.Anything, ListUsersParams{
					Page: PageQuery{
						Number: 1,
						Size:   1,
					},
				}).Return(users[:1], int64(2), nil)
			},
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo, _, _ := createTestService(t)
			tt.setupMocks(mockRepo)

			gotUsers, err := service.List(context.Background(), tt.limit, tt.offset)

			assert.NoError(t, err)
			assert.Len(t, gotUsers, tt.wantCount)
		})
	}
}

// TestGetUserByEmail tests fetching user by email
func TestGetUserByEmail(t *testing.T) {
	// Create a test user
	user := &User{
		Id:    uuid.New(),
		Email: openapi_types.Email("test@example.com"),
		Name:  "Test User",
	}

	tests := []struct {
		name       string
		email      string
		setupMocks func(*MockRepository, *MockCache)
		wantErr    error
	}{
		{
			name:  "Existing user - cache miss",
			email: "test@example.com",
			setupMocks: func(repo *MockRepository, cache *MockCache) {
				cache.EXPECT().GetByEmail(mock.Anything, "test@example.com").Return(nil, ErrCacheMiss)
				repo.EXPECT().GetByEmail(mock.Anything, "test@example.com").Return(user, nil)
				cache.EXPECT().Set(mock.Anything, user, mock.Anything).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:  "Non-existent user",
			email: "nonexistent@example.com",
			setupMocks: func(repo *MockRepository, cache *MockCache) {
				cache.EXPECT().GetByEmail(mock.Anything, "nonexistent@example.com").Return(nil, ErrCacheMiss)
				repo.EXPECT().GetByEmail(mock.Anything, "nonexistent@example.com").Return(nil, ErrUserNotFound)
			},
			wantErr: ErrUserNotFound,
		},
		{
			name:  "Repository error",
			email: "test@example.com",
			setupMocks: func(repo *MockRepository, cache *MockCache) {
				cache.EXPECT().GetByEmail(mock.Anything, "test@example.com").Return(nil, ErrCacheMiss)
				repo.EXPECT().GetByEmail(mock.Anything, "test@example.com").Return(nil, errors.New("database error"))
			},
			wantErr: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo, mockCache, _ := createTestService(t)
			tt.setupMocks(mockRepo, mockCache)

			gotUser, err := service.GetByEmail(context.Background(), tt.email)

			if tt.wantErr != nil {
				assert.Error(t, err)
				if tt.wantErr.Error() != "database error" {
					assert.ErrorIs(t, err, tt.wantErr)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, gotUser)
				assert.Equal(t, tt.email, string(gotUser.Email))
				assert.Equal(t, user.Id, gotUser.Id)
			}
		})
	}
}
