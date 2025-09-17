package users

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test helper functions.
func createTestService(t *testing.T) (*Service, *MockRepository) {
	t.Helper()

	mockRepo := NewMockRepository(t)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))

	service := NewService(mockRepo, nil, logger)
	return service, mockRepo
}

// TestNewService tests the service constructor.
func TestNewService(t *testing.T) {
	service, _ := createTestService(t)

	assert.NotNil(t, service)
	assert.NotNil(t, service.repo)
	assert.NotNil(t, service.logger)
}

// TestGetUser tests getting a user by ID.
func TestGetUser(t *testing.T) {
	// Create a test user
	user := &User{
		ID:            uuid.New(),
		Email:         "test@example.com",
		Name:          "Test User",
		EmailVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	tests := []struct {
		name       string
		userID     uuid.UUID
		setupMocks func(*MockRepository)
		wantErr    error
	}{
		{
			name:   "Existing user",
			userID: user.ID,
			setupMocks: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, user.ID).Return(user, nil)
			},
			wantErr: nil,
		},
		{
			name:   "Non-existent user",
			userID: uuid.New(),
			setupMocks: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, mock.Anything).Return(nil, ErrUserNotFound)
			},
			wantErr: ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := createTestService(t)
			tt.setupMocks(mockRepo)

			request := GetUserRequestObject{
				ID: tt.userID,
			}
			result, err := service.Get(context.Background(), request)

			assert.NoError(t, err) // Service never returns Go errors
			assert.NotNil(t, result)

			if tt.wantErr != nil {
				// Check for error response type
				_, isErrorResp := result.(GetUser404ApplicationProblemPlusJSONResponse)
				assert.True(t, isErrorResp, "Expected error response type")
			} else {
				// Check for success response type
				if resp, ok := result.(GetUser200JSONResponse); ok {
					assert.Equal(t, tt.userID, resp.Data.ID)
				} else {
					t.Fatal("Expected GetUser200JSONResponse")
				}
			}
		})
	}
}

// TestUpdateUser tests updating a user.
func TestUpdateUser(t *testing.T) {
	// Create a test user
	user := &User{
		ID:            uuid.New(),
		Email:         "test@example.com",
		Name:          "Test User",
		EmailVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	tests := []struct {
		name       string
		userID     uuid.UUID
		req        *UpdateUserJSONRequestBody
		setupMocks func(*MockRepository)
		wantErr    bool
	}{
		{
			name:   "Update email",
			userID: user.ID,
			req: &UpdateUserJSONRequestBody{
				Email: "newemail@example.com",
			},
			setupMocks: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, user.ID).Return(user, nil)
				updatedUser := &User{
					ID:            user.ID,
					Email:         "newemail@example.com",
					Name:          user.Name,
					EmailVerified: user.EmailVerified,
					CreatedAt:     user.CreatedAt,
					UpdatedAt:     time.Now(),
				}
				repo.EXPECT().
					Update(mock.Anything, user.ID, mock.AnythingOfType("*users.User")).
					Return(updatedUser, nil)
			},
			wantErr: false,
		},
		{
			name:   "Non-existent user",
			userID: uuid.New(),
			req: &UpdateUserJSONRequestBody{
				Email: "test@example.com",
			},
			setupMocks: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, mock.Anything).Return(nil, ErrUserNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := createTestService(t)
			tt.setupMocks(mockRepo)

			request := UpdateUserRequestObject{
				ID:   tt.userID,
				Body: tt.req,
			}
			result, err := service.Update(context.Background(), request)

			assert.NoError(t, err) // Service never returns Go errors
			assert.NotNil(t, result)

			if tt.wantErr {
				// Check for error response type
				_, isErrorResp := result.(UpdateUser404ApplicationProblemPlusJSONResponse)
				assert.True(t, isErrorResp, "Expected error response type")
			} else {
				// Check for success response type
				if resp, ok := result.(UpdateUser200JSONResponse); ok {
					if tt.req.Email != "" {
						assert.Equal(t, tt.req.Email, string(resp.Data.Email))
					}
				} else {
					t.Fatal("Expected UpdateUser200JSONResponse")
				}
			}
		})
	}
}

// TestDeleteUser tests deleting a user.
func TestDeleteUser(t *testing.T) {
	// Create a test user
	user := &User{
		ID:            uuid.New(),
		Email:         "test@example.com",
		Name:          "Test User",
		EmailVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	tests := []struct {
		name       string
		userID     uuid.UUID
		setupMocks func(*MockRepository)
		wantErr    bool
	}{
		{
			name:   "Delete existing user",
			userID: user.ID,
			setupMocks: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, user.ID).Return(user, nil)
				repo.EXPECT().Delete(mock.Anything, user.ID).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "Delete non-existent user",
			userID: uuid.New(),
			setupMocks: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, mock.Anything).Return(nil, ErrUserNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := createTestService(t)
			tt.setupMocks(mockRepo)

			request := DeleteUserRequestObject{
				ID: tt.userID,
			}
			result, err := service.Delete(context.Background(), request)

			assert.NoError(t, err) // Service never returns Go errors
			assert.NotNil(t, result)

			if tt.wantErr {
				// Check for error response type
				_, isErrorResp := result.(DeleteUser404ApplicationProblemPlusJSONResponse)
				assert.True(t, isErrorResp, "Expected error response type")
			} else {
				// Check for success response type
				_, isSuccessResp := result.(DeleteUser200JSONResponse)
				assert.True(t, isSuccessResp, "Expected success response type")
			}
		})
	}
}

// TestListUsers tests listing users.
func TestListUsers(t *testing.T) {
	// Create test users
	users := []*User{
		{
			ID:            uuid.New(),
			Email:         "user1@example.com",
			Name:          "User 1",
			EmailVerified: true,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
		{
			ID:            uuid.New(),
			Email:         "user2@example.com",
			Name:          "User 2",
			EmailVerified: false,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		},
	}

	tests := []struct {
		name       string
		limit      int
		offset     int
		setupMocks func(*MockRepository)
		wantCount  int
		wantErr    bool
	}{
		{
			name:   "List all users",
			limit:  10,
			offset: 0,
			setupMocks: func(repo *MockRepository) {
				params := ListUsersParams{
					Page: PageQuery{
						Number: 1,
						Size:   10,
					},
				}
				repo.EXPECT().List(mock.Anything, params).Return(users, int64(2), nil)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:   "Empty list",
			limit:  10,
			offset: 0,
			setupMocks: func(repo *MockRepository) {
				params := ListUsersParams{
					Page: PageQuery{
						Number: 1,
						Size:   10,
					},
				}
				repo.EXPECT().List(mock.Anything, params).Return([]*User{}, int64(0), nil)
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := createTestService(t)
			tt.setupMocks(mockRepo)

			request := ListUsersRequestObject{
				Params: ListUsersParams{
					Page: PageQuery{
						Number: tt.offset/tt.limit + 1,
						Size:   tt.limit,
					},
				},
			}
			result, err := service.List(context.Background(), request)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if resp, ok := result.(ListUsers200JSONResponse); ok {
					assert.Len(t, resp.Data, tt.wantCount)
					assert.Equal(t, float32(tt.wantCount), resp.Meta.Total)
				} else {
					t.Fatal("Expected ListUsers200JSONResponse")
				}
			}
		})
	}
}
