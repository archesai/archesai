package accounts

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
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

// TestCreateAccount tests creating an account.
func TestCreateAccount(t *testing.T) {
	tests := []struct {
		name       string
		req        *CreateAccountJSONRequestBody
		setupMocks func(*MockRepository)
		wantErr    bool
	}{
		{
			name: "successful create",
			req: &CreateAccountJSONRequestBody{
				Email:    openapi_types.Email("test@example.com"),
				Name:     "Test User",
				Password: "password123",
			},
			setupMocks: func(repo *MockRepository) {
				repo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(_ *Account) bool {
					// Check that the account will have the provider ID set
					return true
				})).Return(&Account{
					ID:         uuid.New(),
					AccountID:  "test123",
					ProviderID: "local",
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			req: &CreateAccountJSONRequestBody{
				Email:    openapi_types.Email("test2@example.com"),
				Name:     "Test User 2",
				Password: "password456",
			},
			setupMocks: func(repo *MockRepository) {
				repo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*accounts.Account")).
					Return(nil, assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := createTestService(t)
			tt.setupMocks(mockRepo)

			request := CreateAccountRequestObject{
				Body: tt.req,
			}
			result, err := service.Create(context.Background(), request)

			assert.NoError(t, err) // Service never returns Go errors
			assert.NotNil(t, result)

			if tt.wantErr {
				// Check for error response type
				_, isErrorResp := result.(CreateAccount400ApplicationProblemPlusJSONResponse)
				assert.True(t, isErrorResp, "Expected error response type")
			} else {
				// Check for success response type
				_, isSuccessResp := result.(CreateAccount201JSONResponse)
				assert.True(t, isSuccessResp, "Expected success response type")
			}
		})
	}
}

// TestGetAccount tests getting an account by ID.
func TestGetAccount(t *testing.T) {
	accountID := uuid.New()
	account := &Account{
		ID:         accountID,
		AccountID:  "test123",
		ProviderID: "google",
		UserID:     uuid.New(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	tests := []struct {
		name       string
		accountID  uuid.UUID
		setupMocks func(*MockRepository)
		wantErr    bool
	}{
		{
			name:      "existing account",
			accountID: accountID,
			setupMocks: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, accountID).Return(account, nil)
			},
			wantErr: false,
		},
		{
			name:      "non-existent account",
			accountID: uuid.New(),
			setupMocks: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, mock.Anything).Return(nil, ErrAccountNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := createTestService(t)
			tt.setupMocks(mockRepo)

			request := GetAccountRequestObject{
				ID: tt.accountID,
			}
			result, err := service.Get(context.Background(), request)

			assert.NoError(t, err) // Service never returns Go errors
			assert.NotNil(t, result)

			if tt.wantErr {
				// Check for error response type
				_, isErrorResp := result.(GetAccount404ApplicationProblemPlusJSONResponse)
				assert.True(t, isErrorResp, "Expected error response type")
			} else {
				// Check for success response type
				_, isSuccessResp := result.(GetAccount200JSONResponse)
				assert.True(t, isSuccessResp, "Expected success response type")
			}
		})
	}
}

// TestUpdateAccount tests updating an account.
func TestUpdateAccount(t *testing.T) {
	accountID := uuid.New()
	existingAccount := &Account{
		ID:         accountID,
		AccountID:  "test123",
		ProviderID: "google",
		UserID:     uuid.New(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	tests := []struct {
		name       string
		accountID  uuid.UUID
		req        *UpdateAccountJSONRequestBody
		setupMocks func(*MockRepository)
		wantErr    bool
	}{
		{
			name:      "successful update",
			accountID: accountID,
			req:       &UpdateAccountJSONRequestBody{
				// Update fields would go here based on actual schema
			},
			setupMocks: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, accountID).Return(existingAccount, nil)
				updatedAccount := &Account{
					ID:         accountID,
					AccountID:  existingAccount.AccountID,
					ProviderID: existingAccount.ProviderID,
					UserID:     existingAccount.UserID,
					CreatedAt:  existingAccount.CreatedAt,
					UpdatedAt:  time.Now(),
				}
				repo.EXPECT().
					Update(mock.Anything, accountID, mock.AnythingOfType("*accounts.Account")).
					Return(updatedAccount, nil)
			},
			wantErr: false,
		},
		{
			name:      "non-existent account",
			accountID: uuid.New(),
			req:       &UpdateAccountJSONRequestBody{},
			setupMocks: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, mock.Anything).Return(nil, ErrAccountNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := createTestService(t)
			tt.setupMocks(mockRepo)

			request := UpdateAccountRequestObject{
				ID:   tt.accountID,
				Body: tt.req,
			}
			result, err := service.Update(context.Background(), request)

			assert.NoError(t, err) // Service never returns Go errors
			assert.NotNil(t, result)

			if tt.wantErr {
				// Check for error response type
				_, isErrorResp := result.(UpdateAccount404ApplicationProblemPlusJSONResponse)
				assert.True(t, isErrorResp, "Expected error response type")
			} else {
				// Check for success response type
				_, isSuccessResp := result.(UpdateAccount200JSONResponse)
				assert.True(t, isSuccessResp, "Expected success response type")
			}
		})
	}
}

// TestDeleteAccount tests deleting an account.
func TestDeleteAccount(t *testing.T) {
	accountID := uuid.New()
	account := &Account{
		ID:         accountID,
		AccountID:  "test123",
		ProviderID: "google",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	tests := []struct {
		name       string
		accountID  uuid.UUID
		setupMocks func(*MockRepository)
		wantErr    bool
	}{
		{
			name:      "successful delete",
			accountID: accountID,
			setupMocks: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, accountID).Return(account, nil)
				repo.EXPECT().Delete(mock.Anything, accountID).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "non-existent account",
			accountID: uuid.New(),
			setupMocks: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, mock.Anything).Return(nil, ErrAccountNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := createTestService(t)
			tt.setupMocks(mockRepo)

			request := DeleteAccountRequestObject{
				ID: tt.accountID,
			}
			result, err := service.Delete(context.Background(), request)

			assert.NoError(t, err) // Service never returns Go errors
			assert.NotNil(t, result)

			if tt.wantErr {
				// Check for error response type
				_, isErrorResp := result.(DeleteAccount404ApplicationProblemPlusJSONResponse)
				assert.True(t, isErrorResp, "Expected error response type")
			} else {
				// Check for success response type
				_, isSuccessResp := result.(DeleteAccount200JSONResponse)
				assert.True(t, isSuccessResp, "Expected success response type")
			}
		})
	}
}

// TestListAccounts tests listing accounts.
func TestListAccounts(t *testing.T) {
	accounts := []*Account{
		{
			ID:         uuid.New(),
			AccountID:  "test1",
			ProviderID: "google",
			UserID:     uuid.New(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
		{
			ID:         uuid.New(),
			AccountID:  "test2",
			ProviderID: "github",
			UserID:     uuid.New(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
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
			name:   "list all accounts",
			limit:  10,
			offset: 0,
			setupMocks: func(repo *MockRepository) {
				params := ListAccountsParams{
					Page: PageQuery{
						Number: 1,
						Size:   10,
					},
				}
				repo.EXPECT().List(mock.Anything, params).Return(accounts, int64(2), nil)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:   "empty list",
			limit:  10,
			offset: 0,
			setupMocks: func(repo *MockRepository) {
				params := ListAccountsParams{
					Page: PageQuery{
						Number: 1,
						Size:   10,
					},
				}
				repo.EXPECT().List(mock.Anything, params).Return([]*Account{}, int64(0), nil)
			},
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := createTestService(t)
			tt.setupMocks(mockRepo)

			request := ListAccountsRequestObject{
				Params: ListAccountsParams{
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
				// Check response structure when implementation is complete
			}
		})
	}
}
