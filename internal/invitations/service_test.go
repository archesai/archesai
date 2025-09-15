package invitations

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
	ctx := context.Background()
	orgID := uuid.New().String()
	inviterID := uuid.New().String()

	tests := []struct {
		name       string
		invitation *Invitation
		setupMock  func(*MockRepository)
		want       *Invitation
		wantErr    error
	}{
		{
			name: "successful creation",
			invitation: &Invitation{
				Email:          "test@example.com",
				Role:           InvitationRoleAdmin,
				OrganizationId: orgID,
				InviterId:      inviterID,
			},
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().GetByEmail(mock.Anything, "test@example.com", orgID).
					Return(nil, errors.New("not found"))
				repo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*invitations.Invitation")).
					Return(&Invitation{
						Id:             uuid.New(),
						Email:          "test@example.com",
						Role:           InvitationRoleAdmin,
						Status:         StatusPending,
						OrganizationId: orgID,
						InviterId:      inviterID,
						CreatedAt:      time.Now(),
						UpdatedAt:      time.Now(),
					}, nil)
			},
			wantErr: nil,
		},
		{
			name: "invitation already exists",
			invitation: &Invitation{
				Email:          "existing@example.com",
				Role:           InvitationRoleMember,
				OrganizationId: orgID,
				InviterId:      inviterID,
			},
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().GetByEmail(mock.Anything, "existing@example.com", orgID).
					Return(&Invitation{
						Id:     uuid.New(),
						Email:  "existing@example.com",
						Status: StatusPending,
					}, nil)
			},
			wantErr: ErrInvitationAlreadyExists,
		},
		{
			name: "create error",
			invitation: &Invitation{
				Email:          "error@example.com",
				Role:           InvitationRoleOwner,
				OrganizationId: orgID,
				InviterId:      inviterID,
			},
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().GetByEmail(mock.Anything, "error@example.com", orgID).
					Return(nil, errors.New("not found"))
				repo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*invitations.Invitation")).
					Return(nil, errors.New("database error"))
			},
			wantErr: errors.New("failed to create invitation"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockRepository(t)
			tt.setupMock(mockRepo)

			service := NewService(mockRepo, slog.Default())
			result, err := service.Create(ctx, tt.invitation)

			if tt.wantErr != nil {
				assert.Error(t, err)
				if errors.Is(tt.wantErr, ErrInvitationAlreadyExists) {
					assert.ErrorIs(t, err, ErrInvitationAlreadyExists)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, StatusPending, result.Status)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Get(t *testing.T) {
	ctx := context.Background()
	invitationID := uuid.New()

	tests := []struct {
		name      string
		id        uuid.UUID
		setupMock func(*MockRepository)
		want      *Invitation
		wantErr   error
	}{
		{
			name: "successful get",
			id:   invitationID,
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, invitationID).
					Return(&Invitation{
						Id:    invitationID,
						Email: "test@example.com",
					}, nil)
			},
			wantErr: nil,
		},
		{
			name: "invitation not found",
			id:   invitationID,
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, invitationID).
					Return(nil, errors.New("not found"))
			},
			wantErr: ErrInvitationNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockRepository(t)
			tt.setupMock(mockRepo)

			service := NewService(mockRepo, slog.Default())
			result, err := service.Get(ctx, tt.id)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, invitationID, result.Id)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Accept(t *testing.T) {
	ctx := context.Background()
	invitationID := uuid.New()

	tests := []struct {
		name      string
		id        uuid.UUID
		setupMock func(*MockRepository)
		wantErr   error
	}{
		{
			name: "successful accept",
			id:   invitationID,
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, invitationID).
					Return(&Invitation{
						Id:        invitationID,
						Status:    StatusPending,
						ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
					}, nil)
				repo.EXPECT().Update(mock.Anything, invitationID, mock.AnythingOfType("*invitations.Invitation")).
					Return(&Invitation{
						Id:     invitationID,
						Status: StatusAccepted,
					}, nil)
			},
			wantErr: nil,
		},
		{
			name: "already accepted",
			id:   invitationID,
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, invitationID).
					Return(&Invitation{
						Id:     invitationID,
						Status: StatusAccepted,
					}, nil)
			},
			wantErr: ErrInvitationAlreadyAccepted,
		},
		{
			name: "already declined",
			id:   invitationID,
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, invitationID).
					Return(&Invitation{
						Id:     invitationID,
						Status: StatusDeclined,
					}, nil)
			},
			wantErr: ErrInvitationAlreadyDeclined,
		},
		{
			name: "expired invitation",
			id:   invitationID,
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, invitationID).
					Return(&Invitation{
						Id:        invitationID,
						Status:    StatusPending,
						ExpiresAt: time.Now().Add(-time.Hour).Format(time.RFC3339),
					}, nil)
			},
			wantErr: ErrInvitationExpired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockRepository(t)
			tt.setupMock(mockRepo)

			service := NewService(mockRepo, slog.Default())
			result, err := service.Accept(ctx, tt.id)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, StatusAccepted, result.Status)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Decline(t *testing.T) {
	ctx := context.Background()
	invitationID := uuid.New()

	tests := []struct {
		name      string
		id        uuid.UUID
		setupMock func(*MockRepository)
		wantErr   error
	}{
		{
			name: "successful decline",
			id:   invitationID,
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, invitationID).
					Return(&Invitation{
						Id:     invitationID,
						Status: StatusPending,
					}, nil)
				repo.EXPECT().Update(mock.Anything, invitationID, mock.AnythingOfType("*invitations.Invitation")).
					Return(&Invitation{
						Id:     invitationID,
						Status: StatusDeclined,
					}, nil)
			},
			wantErr: nil,
		},
		{
			name: "already accepted",
			id:   invitationID,
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, invitationID).
					Return(&Invitation{
						Id:     invitationID,
						Status: StatusAccepted,
					}, nil)
			},
			wantErr: ErrInvitationAlreadyAccepted,
		},
		{
			name: "already declined",
			id:   invitationID,
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, invitationID).
					Return(&Invitation{
						Id:     invitationID,
						Status: StatusDeclined,
					}, nil)
			},
			wantErr: ErrInvitationAlreadyDeclined,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockRepository(t)
			tt.setupMock(mockRepo)

			service := NewService(mockRepo, slog.Default())
			result, err := service.Decline(ctx, tt.id)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, StatusDeclined, result.Status)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Delete(t *testing.T) {
	ctx := context.Background()
	invitationID := uuid.New()

	tests := []struct {
		name      string
		id        uuid.UUID
		setupMock func(*MockRepository)
		wantErr   bool
	}{
		{
			name: "successful delete",
			id:   invitationID,
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, invitationID).
					Return(&Invitation{
						Id:     invitationID,
						Status: StatusPending,
					}, nil)
				repo.EXPECT().Delete(mock.Anything, invitationID).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "cannot delete accepted",
			id:   invitationID,
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, invitationID).
					Return(&Invitation{
						Id:     invitationID,
						Status: StatusAccepted,
					}, nil)
			},
			wantErr: true,
		},
		{
			name: "invitation not found",
			id:   invitationID,
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, invitationID).
					Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockRepository(t)
			tt.setupMock(mockRepo)

			service := NewService(mockRepo, slog.Default())
			err := service.Delete(ctx, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Update(t *testing.T) {
	ctx := context.Background()
	invitationID := uuid.New()

	tests := []struct {
		name       string
		id         uuid.UUID
		invitation *Invitation
		setupMock  func(*MockRepository)
		wantErr    error
	}{
		{
			name: "successful update",
			id:   invitationID,
			invitation: &Invitation{
				Role: InvitationRoleOwner,
			},
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, invitationID).
					Return(&Invitation{
						Id:        invitationID,
						Status:    StatusPending,
						ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
					}, nil)
				repo.EXPECT().Update(mock.Anything, invitationID, mock.AnythingOfType("*invitations.Invitation")).
					Return(&Invitation{
						Id:   invitationID,
						Role: InvitationRoleOwner,
					}, nil)
			},
			wantErr: nil,
		},
		{
			name: "cannot update accepted",
			id:   invitationID,
			invitation: &Invitation{
				Role: InvitationRoleOwner,
			},
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, invitationID).
					Return(&Invitation{
						Id:     invitationID,
						Status: StatusAccepted,
					}, nil)
			},
			wantErr: ErrInvitationAlreadyAccepted,
		},
		{
			name: "expired invitation",
			id:   invitationID,
			invitation: &Invitation{
				Role: InvitationRoleOwner,
			},
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, invitationID).
					Return(&Invitation{
						Id:        invitationID,
						Status:    StatusPending,
						ExpiresAt: time.Now().Add(-time.Hour).Format(time.RFC3339),
					}, nil)
			},
			wantErr: ErrInvitationExpired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockRepository(t)
			tt.setupMock(mockRepo)

			service := NewService(mockRepo, slog.Default())
			result, err := service.Update(ctx, tt.id, tt.invitation)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_ListByOrganization(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New().String()

	tests := []struct {
		name      string
		orgID     string
		setupMock func(*MockRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name:  "successful list",
			orgID: orgID,
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().ListByOrganization(mock.Anything, orgID).
					Return([]*Invitation{
						{Id: uuid.New(), Email: "user1@example.com"},
						{Id: uuid.New(), Email: "user2@example.com"},
					}, nil)
			},
			wantCount: 2,
			wantErr:   false,
		},
		{
			name:  "empty list",
			orgID: orgID,
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().ListByOrganization(mock.Anything, orgID).
					Return([]*Invitation{}, nil)
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:  "repository error",
			orgID: orgID,
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().ListByOrganization(mock.Anything, orgID).
					Return(nil, errors.New("database error"))
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockRepository(t)
			tt.setupMock(mockRepo)

			service := NewService(mockRepo, slog.Default())
			result, err := service.ListByOrganization(ctx, tt.orgID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.wantCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_Resend(t *testing.T) {
	ctx := context.Background()
	invitationID := uuid.New()

	tests := []struct {
		name      string
		id        uuid.UUID
		setupMock func(*MockRepository)
		wantErr   bool
		errMsg    string
	}{
		{
			name: "successful resend",
			id:   invitationID,
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, invitationID).
					Return(&Invitation{
						Id:     invitationID,
						Status: StatusPending,
					}, nil)
				repo.EXPECT().Update(mock.Anything, invitationID, mock.AnythingOfType("*invitations.Invitation")).
					Return(&Invitation{
						Id:        invitationID,
						Status:    StatusPending,
						ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Format(time.RFC3339),
					}, nil)
			},
			wantErr: false,
		},
		{
			name: "cannot resend accepted",
			id:   invitationID,
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, invitationID).
					Return(&Invitation{
						Id:     invitationID,
						Status: StatusAccepted,
					}, nil)
			},
			wantErr: true,
			errMsg:  "cannot resend invitation with status accepted",
		},
		{
			name: "invitation not found",
			id:   invitationID,
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().Get(mock.Anything, invitationID).
					Return(nil, errors.New("not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockRepository(t)
			tt.setupMock(mockRepo)

			service := NewService(mockRepo, slog.Default())
			result, err := service.Resend(ctx, tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_GetByEmail(t *testing.T) {
	ctx := context.Background()
	orgID := uuid.New().String()

	tests := []struct {
		name      string
		email     string
		orgID     string
		setupMock func(*MockRepository)
		wantErr   error
	}{
		{
			name:  "successful get by email",
			email: "test@example.com",
			orgID: orgID,
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().GetByEmail(mock.Anything, "test@example.com", orgID).
					Return(&Invitation{
						Email:          "test@example.com",
						OrganizationId: orgID,
					}, nil)
			},
			wantErr: nil,
		},
		{
			name:  "not found",
			email: "notfound@example.com",
			orgID: orgID,
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().GetByEmail(mock.Anything, "notfound@example.com", orgID).
					Return(nil, errors.New("not found"))
			},
			wantErr: ErrInvitationNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockRepository(t)
			tt.setupMock(mockRepo)

			service := NewService(mockRepo, slog.Default())
			result, err := service.GetByEmail(ctx, tt.email, tt.orgID)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.email, result.Email)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_ListByInviter(t *testing.T) {
	ctx := context.Background()
	inviterID := uuid.New().String()

	tests := []struct {
		name      string
		inviterID string
		setupMock func(*MockRepository)
		wantCount int
		wantErr   bool
	}{
		{
			name:      "successful list",
			inviterID: inviterID,
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().ListByInviter(mock.Anything, inviterID).
					Return([]*Invitation{
						{Id: uuid.New(), InviterId: inviterID},
						{Id: uuid.New(), InviterId: inviterID},
						{Id: uuid.New(), InviterId: inviterID},
					}, nil)
			},
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:      "empty list",
			inviterID: inviterID,
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().ListByInviter(mock.Anything, inviterID).
					Return([]*Invitation{}, nil)
			},
			wantCount: 0,
			wantErr:   false,
		},
		{
			name:      "repository error",
			inviterID: inviterID,
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().ListByInviter(mock.Anything, inviterID).
					Return(nil, errors.New("database error"))
			},
			wantCount: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockRepository(t)
			tt.setupMock(mockRepo)

			service := NewService(mockRepo, slog.Default())
			result, err := service.ListByInviter(ctx, tt.inviterID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.wantCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestService_List(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		params    ListInvitationsParams
		setupMock func(*MockRepository)
		wantCount int
		wantTotal int64
		wantErr   bool
	}{
		{
			name: "successful list with pagination",
			params: ListInvitationsParams{
				Page: PageQuery{Number: 1, Size: 10},
			},
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().List(mock.Anything, ListInvitationsParams{
					Page: PageQuery{Number: 1, Size: 10},
				}).Return([]*Invitation{
					{Id: uuid.New()},
					{Id: uuid.New()},
				}, int64(2), nil)
			},
			wantCount: 2,
			wantTotal: 2,
			wantErr:   false,
		},
		{
			name: "default limit applied",
			params: ListInvitationsParams{
				Page: PageQuery{Number: 1, Size: 0},
			},
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().List(mock.Anything, ListInvitationsParams{
					Page: PageQuery{Number: 1, Size: 10},
				}).Return([]*Invitation{}, int64(0), nil)
			},
			wantCount: 0,
			wantTotal: 0,
			wantErr:   false,
		},
		{
			name: "max limit enforced",
			params: ListInvitationsParams{
				Page: PageQuery{Number: 1, Size: 200},
			},
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().List(mock.Anything, ListInvitationsParams{
					Page: PageQuery{Number: 1, Size: 100},
				}).Return([]*Invitation{}, int64(0), nil)
			},
			wantCount: 0,
			wantTotal: 0,
			wantErr:   false,
		},
		{
			name: "repository error",
			params: ListInvitationsParams{
				Page: PageQuery{Number: 1, Size: 10},
			},
			setupMock: func(repo *MockRepository) {
				repo.EXPECT().List(mock.Anything, mock.Anything).
					Return(nil, int64(0), errors.New("database error"))
			},
			wantCount: 0,
			wantTotal: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockRepository(t)
			tt.setupMock(mockRepo)

			service := NewService(mockRepo, slog.Default())
			result, total, err := service.List(ctx, tt.params)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.wantCount)
				assert.Equal(t, tt.wantTotal, total)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestNewService(t *testing.T) {
	t.Run("with logger", func(t *testing.T) {
		mockRepo := NewMockRepository(t)
		logger := slog.Default()
		service := NewService(mockRepo, logger)

		assert.NotNil(t, service)
		assert.Equal(t, mockRepo, service.repo)
		assert.Equal(t, logger, service.logger)
	})

	t.Run("without logger uses default", func(t *testing.T) {
		mockRepo := NewMockRepository(t)
		service := NewService(mockRepo, nil)

		assert.NotNil(t, service)
		assert.Equal(t, mockRepo, service.repo)
		assert.NotNil(t, service.logger)
	})
}
