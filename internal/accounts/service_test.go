package accounts

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Create(t *testing.T) {
	tests := []struct {
		name      string
		account   *Account
		setupMock func(*MockRepository, *MockCache, *MockEventPublisher)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "successful create with new UUID",
			account: &Account{
				AccountId:  "test123",
				ProviderId: Google,
				UserId:     uuid.New(),
			},
			setupMock: func(repo *MockRepository, cache *MockCache, events *MockEventPublisher) {
				repo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*accounts.Account")).
					Return(&Account{
						Id:         uuid.New(),
						AccountId:  "test123",
						ProviderId: Google,
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
					}, nil)
				cache.EXPECT().Set(mock.Anything, mock.AnythingOfType("*accounts.Account"), mock.Anything).Return(nil)
				events.EXPECT().PublishAccountCreated(mock.Anything, mock.AnythingOfType("*accounts.Account")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "successful create with existing UUID",
			account: &Account{
				Id:         uuid.New(),
				AccountId:  "test456",
				ProviderId: Github,
				UserId:     uuid.New(),
			},
			setupMock: func(repo *MockRepository, cache *MockCache, events *MockEventPublisher) {
				repo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*accounts.Account")).
					Return(&Account{
						Id:         uuid.New(),
						AccountId:  "test456",
						ProviderId: Github,
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
					}, nil)
				cache.EXPECT().Set(mock.Anything, mock.AnythingOfType("*accounts.Account"), mock.Anything).Return(nil)
				events.EXPECT().PublishAccountCreated(mock.Anything, mock.AnythingOfType("*accounts.Account")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "repository error",
			account: &Account{
				AccountId:  "test789",
				ProviderId: Microsoft,
				UserId:     uuid.New(),
			},
			setupMock: func(repo *MockRepository, _ *MockCache, _ *MockEventPublisher) {
				repo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*accounts.Account")).
					Return(nil, errors.New("database error"))
			},
			wantErr: true,
			errMsg:  "failed to create account",
		},
		{
			name: "cache error ignored",
			account: &Account{
				AccountId:  "test999",
				ProviderId: Local,
				UserId:     uuid.New(),
			},
			setupMock: func(repo *MockRepository, cache *MockCache, events *MockEventPublisher) {
				repo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*accounts.Account")).
					Return(&Account{
						Id:         uuid.New(),
						AccountId:  "test999",
						ProviderId: Local,
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
					}, nil)
				cache.EXPECT().Set(mock.Anything, mock.AnythingOfType("*accounts.Account"), mock.Anything).
					Return(errors.New("cache error"))
				events.EXPECT().PublishAccountCreated(mock.Anything, mock.AnythingOfType("*accounts.Account")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "event publish error ignored",
			account: &Account{
				AccountId:  "test111",
				ProviderId: Apple,
				UserId:     uuid.New(),
			},
			setupMock: func(repo *MockRepository, cache *MockCache, events *MockEventPublisher) {
				repo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*accounts.Account")).
					Return(&Account{
						Id:         uuid.New(),
						AccountId:  "test111",
						ProviderId: Apple,
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
					}, nil)
				cache.EXPECT().Set(mock.Anything, mock.AnythingOfType("*accounts.Account"), mock.Anything).Return(nil)
				events.EXPECT().PublishAccountCreated(mock.Anything, mock.AnythingOfType("*accounts.Account")).
					Return(errors.New("event error"))
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository(t)
			cache := NewMockCache(t)
			events := NewMockEventPublisher(t)

			tt.setupMock(repo, cache, events)

			service := NewService(repo, cache, events, slog.Default())
			result, err := service.Create(context.Background(), tt.account)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestService_Get(t *testing.T) {
	accountID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name      string
		id        uuid.UUID
		setupMock func(*MockRepository, *MockCache)
		want      *Account
		wantErr   error
	}{
		{
			name: "successful get from cache",
			id:   accountID,
			setupMock: func(_ *MockRepository, cache *MockCache) {
				cache.EXPECT().Get(mock.Anything, accountID).Return(&Account{
					Id:         accountID,
					AccountId:  "cached123",
					ProviderId: Google,
					UserId:     userID,
				}, nil)
			},
			want: &Account{
				Id:         accountID,
				AccountId:  "cached123",
				ProviderId: Google,
				UserId:     userID,
			},
		},
		{
			name: "cache miss, successful get from repo",
			id:   accountID,
			setupMock: func(repo *MockRepository, cache *MockCache) {
				cache.EXPECT().Get(mock.Anything, accountID).Return(nil, errors.New("cache miss"))
				repo.EXPECT().Get(mock.Anything, accountID).Return(&Account{
					Id:         accountID,
					AccountId:  "repo123",
					ProviderId: Github,
					UserId:     userID,
				}, nil)
				cache.EXPECT().Set(mock.Anything, mock.AnythingOfType("*accounts.Account"), mock.Anything).Return(nil)
			},
			want: &Account{
				Id:         accountID,
				AccountId:  "repo123",
				ProviderId: Github,
				UserId:     userID,
			},
		},
		{
			name: "account not found",
			id:   accountID,
			setupMock: func(repo *MockRepository, cache *MockCache) {
				cache.EXPECT().Get(mock.Anything, accountID).Return(nil, errors.New("cache miss"))
				repo.EXPECT().Get(mock.Anything, accountID).Return(nil, errors.New("not found"))
			},
			wantErr: ErrAccountNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository(t)
			cache := NewMockCache(t)
			events := NewMockEventPublisher(t)

			tt.setupMock(repo, cache)

			service := NewService(repo, cache, events, slog.Default())
			result, err := service.Get(context.Background(), tt.id)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}
		})
	}
}

func TestService_GetByProviderID(t *testing.T) {
	accountID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name              string
		provider          string
		providerAccountID string
		setupMock         func(*MockRepository, *MockCache)
		want              *Account
		wantErr           error
	}{
		{
			name:              "successful get from cache",
			provider:          "google",
			providerAccountID: "google123",
			setupMock: func(_ *MockRepository, cache *MockCache) {
				cache.EXPECT().GetByProviderId(mock.Anything, "google", "google123").Return(&Account{
					Id:         accountID,
					AccountId:  "google123",
					ProviderId: Google,
					UserId:     userID,
				}, nil)
			},
			want: &Account{
				Id:         accountID,
				AccountId:  "google123",
				ProviderId: Google,
				UserId:     userID,
			},
		},
		{
			name:              "cache miss, successful get from repo",
			provider:          "github",
			providerAccountID: "github456",
			setupMock: func(repo *MockRepository, cache *MockCache) {
				cache.EXPECT().GetByProviderId(mock.Anything, "github", "github456").
					Return(nil, errors.New("cache miss"))
				repo.EXPECT().GetByProviderId(mock.Anything, "github", "github456").Return(&Account{
					Id:         accountID,
					AccountId:  "github456",
					ProviderId: Github,
					UserId:     userID,
				}, nil)
				cache.EXPECT().Set(mock.Anything, mock.AnythingOfType("*accounts.Account"), mock.Anything).Return(nil)
			},
			want: &Account{
				Id:         accountID,
				AccountId:  "github456",
				ProviderId: Github,
				UserId:     userID,
			},
		},
		{
			name:              "account not found",
			provider:          "microsoft",
			providerAccountID: "ms789",
			setupMock: func(repo *MockRepository, cache *MockCache) {
				cache.EXPECT().GetByProviderId(mock.Anything, "microsoft", "ms789").
					Return(nil, errors.New("cache miss"))
				repo.EXPECT().GetByProviderId(mock.Anything, "microsoft", "ms789").
					Return(nil, errors.New("not found"))
			},
			wantErr: ErrAccountNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository(t)
			cache := NewMockCache(t)
			events := NewMockEventPublisher(t)

			tt.setupMock(repo, cache)

			service := NewService(repo, cache, events, slog.Default())
			result, err := service.GetByProviderID(context.Background(), tt.provider, tt.providerAccountID)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}
		})
	}
}

func TestService_List(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name      string
		params    ListAccountsParams
		setupMock func(*MockRepository)
		want      []*Account
		wantTotal int64
		wantErr   bool
	}{
		{
			name: "successful list",
			params: ListAccountsParams{
				Page: PageQuery{
					Number: 1,
					Size:   10,
				},
			},
			setupMock: func(repo *MockRepository) {
				accounts := []*Account{
					{Id: uuid.New(), AccountId: "acc1", ProviderId: Google, UserId: userID},
					{Id: uuid.New(), AccountId: "acc2", ProviderId: Github, UserId: userID},
				}
				repo.EXPECT().List(mock.Anything, mock.AnythingOfType("accounts.ListAccountsParams")).
					Return(accounts, int64(2), nil)
			},
			wantTotal: 2,
			wantErr:   false,
		},
		{
			name: "repository error",
			params: ListAccountsParams{
				Page: PageQuery{
					Number: 1,
					Size:   10,
				},
			},
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().List(mock.Anything, mock.AnythingOfType("accounts.ListAccountsParams")).
					Return(nil, int64(0), errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository(t)
			cache := NewMockCache(t)
			events := NewMockEventPublisher(t)

			tt.setupMock(repo)

			service := NewService(repo, cache, events, slog.Default())
			result, total, err := service.List(context.Background(), tt.params)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantTotal, total)
				if tt.want != nil {
					assert.Len(t, result, len(tt.want))
				}
			}
		})
	}
}

func TestService_ListByUserID(t *testing.T) {
	userID := uuid.New()

	tests := []struct {
		name      string
		userID    uuid.UUID
		setupMock func(*MockRepository, *MockCache)
		want      []*Account
		wantErr   bool
	}{
		{
			name:   "successful list from cache",
			userID: userID,
			setupMock: func(_ *MockRepository, cache *MockCache) {
				accounts := []*Account{
					{Id: uuid.New(), AccountId: "acc1", ProviderId: Google, UserId: userID},
					{Id: uuid.New(), AccountId: "acc2", ProviderId: Github, UserId: userID},
				}
				cache.EXPECT().ListByUserId(mock.Anything, userID).Return(accounts, nil)
			},
			wantErr: false,
		},
		{
			name:   "cache miss, successful list from repo",
			userID: userID,
			setupMock: func(repo *MockRepository, cache *MockCache) {
				accounts := []*Account{
					{Id: uuid.New(), AccountId: "acc1", ProviderId: Google, UserId: userID},
				}
				cache.EXPECT().ListByUserId(mock.Anything, userID).Return(nil, errors.New("cache miss"))
				repo.EXPECT().ListByUserId(mock.Anything, userID).Return(accounts, nil)
			},
			wantErr: false,
		},
		{
			name:   "repository error",
			userID: userID,
			setupMock: func(repo *MockRepository, cache *MockCache) {
				cache.EXPECT().ListByUserId(mock.Anything, userID).Return(nil, errors.New("cache miss"))
				repo.EXPECT().ListByUserId(mock.Anything, userID).Return(nil, errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository(t)
			cache := NewMockCache(t)
			events := NewMockEventPublisher(t)

			tt.setupMock(repo, cache)

			service := NewService(repo, cache, events, slog.Default())
			result, err := service.ListByUserID(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	accountID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name      string
		id        uuid.UUID
		account   *Account
		setupMock func(*MockRepository, *MockCache, *MockEventPublisher)
		wantErr   error
	}{
		{
			name: "successful update",
			id:   accountID,
			account: &Account{
				AccountId:  "updated123",
				ProviderId: Google,
				UserId:     userID,
			},
			setupMock: func(repo *MockRepository, cache *MockCache, events *MockEventPublisher) {
				existing := &Account{
					Id:         accountID,
					AccountId:  "old123",
					ProviderId: Google,
					UserId:     userID,
					CreatedAt:  time.Now().Add(-time.Hour),
					UpdatedAt:  time.Now().Add(-time.Hour),
				}
				repo.EXPECT().Get(mock.Anything, accountID).Return(existing, nil)
				repo.EXPECT().Update(mock.Anything, accountID, mock.AnythingOfType("*accounts.Account")).
					Return(&Account{
						Id:         accountID,
						AccountId:  "updated123",
						ProviderId: Google,
						UserId:     userID,
						CreatedAt:  existing.CreatedAt,
						UpdatedAt:  time.Now(),
					}, nil)
				cache.EXPECT().Delete(mock.Anything, accountID).Return(nil)
				cache.EXPECT().Set(mock.Anything, mock.AnythingOfType("*accounts.Account"), mock.Anything).Return(nil)
				events.EXPECT().PublishAccountUpdated(mock.Anything, mock.AnythingOfType("*accounts.Account")).Return(nil)
			},
		},
		{
			name: "account not found",
			id:   accountID,
			account: &Account{
				AccountId:  "notfound",
				ProviderId: Github,
			},
			setupMock: func(repo *MockRepository, _ *MockCache, _ *MockEventPublisher) {
				repo.EXPECT().Get(mock.Anything, accountID).Return(nil, errors.New("not found"))
			},
			wantErr: ErrAccountNotFound,
		},
		{
			name: "update error",
			id:   accountID,
			account: &Account{
				AccountId:  "error123",
				ProviderId: Microsoft,
			},
			setupMock: func(repo *MockRepository, _ *MockCache, _ *MockEventPublisher) {
				existing := &Account{
					Id:        accountID,
					CreatedAt: time.Now().Add(-time.Hour),
				}
				repo.EXPECT().Get(mock.Anything, accountID).Return(existing, nil)
				repo.EXPECT().Update(mock.Anything, accountID, mock.AnythingOfType("*accounts.Account")).
					Return(nil, errors.New("update failed"))
			},
			wantErr: errors.New("failed to update account"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository(t)
			cache := NewMockCache(t)
			events := NewMockEventPublisher(t)

			tt.setupMock(repo, cache, events)

			service := NewService(repo, cache, events, slog.Default())
			result, err := service.Update(context.Background(), tt.id, tt.account)

			if tt.wantErr != nil {
				if tt.wantErr == ErrAccountNotFound {
					assert.Equal(t, tt.wantErr, err)
				} else {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tt.wantErr.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	accountID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name      string
		id        uuid.UUID
		setupMock func(*MockRepository, *MockCache, *MockEventPublisher)
		wantErr   error
	}{
		{
			name: "successful delete",
			id:   accountID,
			setupMock: func(repo *MockRepository, cache *MockCache, events *MockEventPublisher) {
				account := &Account{
					Id:         accountID,
					AccountId:  "delete123",
					ProviderId: Google,
					UserId:     userID,
				}
				repo.EXPECT().Get(mock.Anything, accountID).Return(account, nil)
				repo.EXPECT().Delete(mock.Anything, accountID).Return(nil)
				cache.EXPECT().Delete(mock.Anything, accountID).Return(nil)
				events.EXPECT().PublishAccountDeleted(mock.Anything, account).Return(nil)
			},
		},
		{
			name: "account not found",
			id:   accountID,
			setupMock: func(repo *MockRepository, _ *MockCache, _ *MockEventPublisher) {
				repo.EXPECT().Get(mock.Anything, accountID).Return(nil, errors.New("not found"))
			},
			wantErr: ErrAccountNotFound,
		},
		{
			name: "delete error",
			id:   accountID,
			setupMock: func(repo *MockRepository, _ *MockCache, _ *MockEventPublisher) {
				account := &Account{
					Id:         accountID,
					AccountId:  "error123",
					ProviderId: Github,
					UserId:     userID,
				}
				repo.EXPECT().Get(mock.Anything, accountID).Return(account, nil)
				repo.EXPECT().Delete(mock.Anything, accountID).Return(errors.New("delete failed"))
			},
			wantErr: errors.New("failed to delete account"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository(t)
			cache := NewMockCache(t)
			events := NewMockEventPublisher(t)

			tt.setupMock(repo, cache, events)

			service := NewService(repo, cache, events, slog.Default())
			err := service.Delete(context.Background(), tt.id)

			if tt.wantErr != nil {
				if tt.wantErr == ErrAccountNotFound {
					assert.Equal(t, tt.wantErr, err)
				} else {
					assert.Error(t, err)
					assert.Contains(t, err.Error(), tt.wantErr.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_LinkAccount(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()

	tests := []struct {
		name      string
		userID    uuid.UUID
		account   *Account
		setupMock func(*MockRepository, *MockCache, *MockEventPublisher)
		wantErr   error
	}{
		{
			name:   "link new account",
			userID: userID,
			account: &Account{
				AccountId:  "new123",
				ProviderId: Google,
			},
			setupMock: func(repo *MockRepository, cache *MockCache, events *MockEventPublisher) {
				repo.EXPECT().GetByProviderId(mock.Anything, "google", "new123").
					Return(nil, errors.New("not found"))
				repo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*accounts.Account")).
					Return(&Account{
						Id:         uuid.New(),
						AccountId:  "new123",
						ProviderId: Google,
						UserId:     userID,
					}, nil)
				cache.EXPECT().Set(mock.Anything, mock.AnythingOfType("*accounts.Account"), mock.Anything).Return(nil)
				events.EXPECT().PublishAccountCreated(mock.Anything, mock.AnythingOfType("*accounts.Account")).Return(nil)
				events.EXPECT().PublishAccountLinked(mock.Anything, mock.AnythingOfType("*accounts.Account")).Return(nil)
			},
		},
		{
			name:   "update existing account for same user",
			userID: userID,
			account: &Account{
				AccountId:  "existing123",
				ProviderId: Github,
			},
			setupMock: func(repo *MockRepository, cache *MockCache, events *MockEventPublisher) {
				existingID := uuid.New()
				existing := &Account{
					Id:         existingID,
					AccountId:  "existing123",
					ProviderId: Github,
					UserId:     userID,
					CreatedAt:  time.Now().Add(-time.Hour),
				}
				repo.EXPECT().GetByProviderId(mock.Anything, "github", "existing123").Return(existing, nil)
				repo.EXPECT().Get(mock.Anything, existingID).Return(existing, nil)
				repo.EXPECT().Update(mock.Anything, existingID, mock.AnythingOfType("*accounts.Account")).
					Return(&Account{
						Id:         existingID,
						AccountId:  "existing123",
						ProviderId: Github,
						UserId:     userID,
					}, nil)
				cache.EXPECT().Delete(mock.Anything, existingID).Return(nil)
				cache.EXPECT().Set(mock.Anything, mock.AnythingOfType("*accounts.Account"), mock.Anything).Return(nil)
				events.EXPECT().PublishAccountUpdated(mock.Anything, mock.AnythingOfType("*accounts.Account")).Return(nil)
			},
		},
		{
			name:   "account already linked to different user",
			userID: userID,
			account: &Account{
				AccountId:  "conflict123",
				ProviderId: Microsoft,
			},
			setupMock: func(repo *MockRepository, _ *MockCache, _ *MockEventPublisher) {
				existing := &Account{
					Id:         uuid.New(),
					AccountId:  "conflict123",
					ProviderId: Microsoft,
					UserId:     otherUserID,
				}
				repo.EXPECT().GetByProviderId(mock.Anything, "microsoft", "conflict123").Return(existing, nil)
			},
			wantErr: ErrDuplicateAccount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository(t)
			cache := NewMockCache(t)
			events := NewMockEventPublisher(t)

			tt.setupMock(repo, cache, events)

			service := NewService(repo, cache, events, slog.Default())
			result, err := service.LinkAccount(context.Background(), tt.userID, tt.account)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestService_UnlinkAccount(t *testing.T) {
	accountID := uuid.New()
	userID := uuid.New()
	otherUserID := uuid.New()

	tests := []struct {
		name      string
		userID    uuid.UUID
		accountID uuid.UUID
		setupMock func(*MockRepository, *MockCache, *MockEventPublisher)
		wantErr   error
	}{
		{
			name:      "successful unlink",
			userID:    userID,
			accountID: accountID,
			setupMock: func(repo *MockRepository, cache *MockCache, events *MockEventPublisher) {
				account := &Account{
					Id:         accountID,
					AccountId:  "unlink123",
					ProviderId: Google,
					UserId:     userID,
				}
				cache.EXPECT().Get(mock.Anything, accountID).Return(account, nil)
				repo.EXPECT().Get(mock.Anything, accountID).Return(account, nil)
				repo.EXPECT().Delete(mock.Anything, accountID).Return(nil)
				cache.EXPECT().Delete(mock.Anything, accountID).Return(nil)
				events.EXPECT().PublishAccountDeleted(mock.Anything, account).Return(nil)
				events.EXPECT().PublishAccountUnlinked(mock.Anything, account).Return(nil)
			},
		},
		{
			name:      "account not found",
			userID:    userID,
			accountID: accountID,
			setupMock: func(repo *MockRepository, cache *MockCache, _ *MockEventPublisher) {
				cache.EXPECT().Get(mock.Anything, accountID).Return(nil, errors.New("cache miss"))
				repo.EXPECT().Get(mock.Anything, accountID).Return(nil, errors.New("not found"))
			},
			wantErr: ErrAccountNotFound,
		},
		{
			name:      "account belongs to different user",
			userID:    userID,
			accountID: accountID,
			setupMock: func(_ *MockRepository, cache *MockCache, _ *MockEventPublisher) {
				account := &Account{
					Id:         accountID,
					AccountId:  "other123",
					ProviderId: Github,
					UserId:     otherUserID,
				}
				cache.EXPECT().Get(mock.Anything, accountID).Return(account, nil)
			},
			wantErr: ErrAccountNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository(t)
			cache := NewMockCache(t)
			events := NewMockEventPublisher(t)

			tt.setupMock(repo, cache, events)

			service := NewService(repo, cache, events, slog.Default())
			err := service.UnlinkAccount(context.Background(), tt.userID, tt.accountID)

			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_WithNilDependencies(t *testing.T) {
	t.Run("service works without cache", func(t *testing.T) {
		repo := NewMockRepository(t)
		events := NewMockEventPublisher(t)

		accountID := uuid.New()
		account := &Account{
			Id:         accountID,
			AccountId:  "nocache123",
			ProviderId: Google,
		}

		repo.EXPECT().Get(mock.Anything, accountID).Return(account, nil)

		service := NewService(repo, nil, events, slog.Default())
		result, err := service.Get(context.Background(), accountID)

		assert.NoError(t, err)
		assert.Equal(t, account, result)
	})

	t.Run("service works without events", func(t *testing.T) {
		repo := NewMockRepository(t)
		cache := NewMockCache(t)

		account := &Account{
			AccountId:  "noevents123",
			ProviderId: Github,
			UserId:     uuid.New(),
		}

		repo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*accounts.Account")).
			Return(&Account{
				Id:         uuid.New(),
				AccountId:  "noevents123",
				ProviderId: Github,
			}, nil)
		cache.EXPECT().Set(mock.Anything, mock.AnythingOfType("*accounts.Account"), mock.Anything).Return(nil)

		service := NewService(repo, cache, nil, slog.Default())
		result, err := service.Create(context.Background(), account)

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}
