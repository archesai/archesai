package organizations

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/archesai/archesai/internal/logger"
	"github.com/google/uuid"
)

// MockRepository implements ExtendedRepository for testing
type MockRepository struct {
	organizations map[uuid.UUID]*Organization
	members       map[uuid.UUID]*Member
	invitations   map[uuid.UUID]*Invitation
	err           error
}

// Compile-time check
var _ ExtendedRepository = (*MockRepository)(nil)

func NewMockRepository() *MockRepository {
	return &MockRepository{
		organizations: make(map[uuid.UUID]*Organization),
		members:       make(map[uuid.UUID]*Member),
		invitations:   make(map[uuid.UUID]*Invitation),
	}
}

// Organization methods
func (m *MockRepository) CreateOrganization(_ context.Context, org *Organization) (*Organization, error) {
	if m.err != nil {
		return nil, m.err
	}
	if org.Id == uuid.Nil {
		org.Id = uuid.New()
	}
	org.CreatedAt = time.Now()
	org.UpdatedAt = time.Now()
	m.organizations[org.Id] = org
	return org, nil
}

func (m *MockRepository) GetOrganization(_ context.Context, id uuid.UUID) (*Organization, error) {
	if m.err != nil {
		return nil, m.err
	}
	org, exists := m.organizations[id]
	if !exists {
		return nil, ErrOrganizationNotFound
	}
	return org, nil
}

func (m *MockRepository) UpdateOrganization(_ context.Context, id uuid.UUID, org *Organization) (*Organization, error) {
	if m.err != nil {
		return nil, m.err
	}
	if _, exists := m.organizations[id]; !exists {
		return nil, ErrOrganizationNotFound
	}
	org.UpdatedAt = time.Now()
	m.organizations[id] = org
	return org, nil
}

func (m *MockRepository) DeleteOrganization(_ context.Context, id uuid.UUID) error {
	if m.err != nil {
		return m.err
	}
	if _, exists := m.organizations[id]; !exists {
		return ErrOrganizationNotFound
	}
	delete(m.organizations, id)
	return nil
}

func (m *MockRepository) ListOrganizations(_ context.Context, _ ListOrganizationsParams) ([]*Organization, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	orgs := make([]*Organization, 0, len(m.organizations))
	for _, org := range m.organizations {
		orgs = append(orgs, org)
	}
	return orgs, int64(len(orgs)), nil
}

// Member methods
func (m *MockRepository) CreateMember(_ context.Context, member *Member) (*Member, error) {
	if m.err != nil {
		return nil, m.err
	}
	if member.Id == uuid.Nil {
		member.Id = uuid.New()
	}
	member.CreatedAt = time.Now()
	member.UpdatedAt = time.Now()
	m.members[member.Id] = member
	return member, nil
}

func (m *MockRepository) GetMember(_ context.Context, id uuid.UUID) (*Member, error) {
	if m.err != nil {
		return nil, m.err
	}
	member, exists := m.members[id]
	if !exists {
		return nil, ErrMemberNotFound
	}
	return member, nil
}

func (m *MockRepository) GetMemberByUserAndOrg(_ context.Context, userID, orgID string) (*Member, error) {
	if m.err != nil {
		return nil, m.err
	}
	// Simple implementation - find first member that matches
	for _, member := range m.members {
		if member.UserId == userID && member.OrganizationId == orgID {
			return member, nil
		}
	}
	return nil, ErrMemberNotFound
}

func (m *MockRepository) UpdateMember(_ context.Context, id uuid.UUID, member *Member) (*Member, error) {
	if m.err != nil {
		return nil, m.err
	}
	if _, exists := m.members[id]; !exists {
		return nil, ErrMemberNotFound
	}
	member.UpdatedAt = time.Now()
	m.members[id] = member
	return member, nil
}

func (m *MockRepository) DeleteMember(_ context.Context, id uuid.UUID) error {
	if m.err != nil {
		return m.err
	}
	if _, exists := m.members[id]; !exists {
		return ErrMemberNotFound
	}
	delete(m.members, id)
	return nil
}

func (m *MockRepository) ListMembers(_ context.Context, _ ListMembersParams) ([]*Member, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	members := make([]*Member, 0, len(m.members))
	for _, member := range m.members {
		members = append(members, member)
	}
	return members, int64(len(members)), nil
}

// Invitation methods
func (m *MockRepository) CreateInvitation(_ context.Context, inv *Invitation) (*Invitation, error) {
	if m.err != nil {
		return nil, m.err
	}
	if inv.Id == uuid.Nil {
		inv.Id = uuid.New()
	}
	inv.CreatedAt = time.Now()
	inv.UpdatedAt = time.Now()
	m.invitations[inv.Id] = inv
	return inv, nil
}

func (m *MockRepository) GetInvitation(_ context.Context, id uuid.UUID) (*Invitation, error) {
	if m.err != nil {
		return nil, m.err
	}
	inv, exists := m.invitations[id]
	if !exists {
		return nil, ErrInvitationNotFound
	}
	return inv, nil
}

func (m *MockRepository) UpdateInvitation(_ context.Context, _ uuid.UUID, inv *Invitation) (*Invitation, error) {
	if m.err != nil {
		return nil, m.err
	}
	if _, exists := m.invitations[inv.Id]; !exists {
		return nil, ErrInvitationNotFound
	}
	inv.UpdatedAt = time.Now()
	m.invitations[inv.Id] = inv
	return inv, nil
}

func (m *MockRepository) DeleteInvitation(_ context.Context, id uuid.UUID) error {
	if m.err != nil {
		return m.err
	}
	if _, exists := m.invitations[id]; !exists {
		return ErrInvitationNotFound
	}
	delete(m.invitations, id)
	return nil
}

func (m *MockRepository) ListInvitations(_ context.Context, _ ListInvitationsParams) ([]*Invitation, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	invs := make([]*Invitation, 0, len(m.invitations))
	for _, inv := range m.invitations {
		invs = append(invs, inv)
	}
	return invs, int64(len(invs)), nil
}

// Test cases
func TestService_CreateOrganization(t *testing.T) {
	tests := []struct {
		name      string
		req       *CreateOrganizationRequest
		creatorID string
		repoErr   error
		wantErr   bool
	}{
		{
			name: "successful creation",
			req: &CreateOrganizationRequest{
				OrganizationId: uuid.New(),
				BillingEmail:   "billing@example.com",
			},
			creatorID: "user-123",
			wantErr:   false,
		},
		{
			name: "repository error",
			req: &CreateOrganizationRequest{
				OrganizationId: uuid.New(),
				BillingEmail:   "billing@example.com",
			},
			creatorID: "user-123",
			repoErr:   errors.New("database error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			repo.err = tt.repoErr
			service := NewService(repo, logger.NewTest())

			org, err := service.CreateOrganization(context.Background(), tt.req, tt.creatorID)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateOrganization() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && org == nil {
				t.Error("CreateOrganization() returned nil organization")
			}
		})
	}
}

func TestService_GetOrganization(t *testing.T) {
	tests := []struct {
		name    string
		orgID   uuid.UUID
		setup   func(*MockRepository)
		wantErr bool
	}{
		{
			name:  "existing organization",
			orgID: uuid.New(),
			setup: func(r *MockRepository) {
				org := &Organization{
					Id:           uuid.New(),
					Name:         "Test Org",
					BillingEmail: "billing@example.com",
					Plan:         OrganizationPlan(DefaultPlan),
				}
				r.organizations[org.Id] = org
			},
			wantErr: false,
		},
		{
			name:    "non-existent organization",
			orgID:   uuid.New(),
			setup:   func(_ *MockRepository) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, logger.NewTest())

			// Get the correct ID from the setup
			var testID uuid.UUID
			if tt.name == "existing organization" {
				for id := range repo.organizations {
					testID = id
					break
				}
			} else {
				testID = tt.orgID
			}

			org, err := service.GetOrganization(context.Background(), testID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrganization() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && org == nil {
				t.Error("GetOrganization() returned nil organization")
			}
		})
	}
}

func TestService_UpdateOrganization(t *testing.T) {
	existingID := uuid.New()

	tests := []struct {
		name    string
		orgID   uuid.UUID
		req     *UpdateOrganizationRequest
		setup   func(*MockRepository)
		wantErr bool
	}{
		{
			name:  "successful update",
			orgID: existingID,
			req: &UpdateOrganizationRequest{
				BillingEmail: "new@example.com",
			},
			setup: func(r *MockRepository) {
				org := &Organization{
					Id:           existingID,
					Name:         "Test Org",
					BillingEmail: "old@example.com",
					Plan:         OrganizationPlan(DefaultPlan),
				}
				r.organizations[existingID] = org
			},
			wantErr: false,
		},
		{
			name:  "non-existent organization",
			orgID: uuid.New(),
			req: &UpdateOrganizationRequest{
				BillingEmail: "new@example.com",
			},
			setup:   func(_ *MockRepository) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, logger.NewTest())

			org, err := service.UpdateOrganization(context.Background(), tt.orgID, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateOrganization() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if org == nil {
					t.Error("UpdateOrganization() returned nil organization")
				} else if tt.req.BillingEmail != "" && string(org.BillingEmail) != tt.req.BillingEmail {
					t.Errorf("UpdateOrganization() billing email = %v, want %v", org.BillingEmail, tt.req.BillingEmail)
				}
			}
		})
	}
}

func TestService_DeleteOrganization(t *testing.T) {
	existingID := uuid.New()

	tests := []struct {
		name    string
		orgID   uuid.UUID
		setup   func(*MockRepository)
		wantErr bool
	}{
		{
			name:  "successful deletion",
			orgID: existingID,
			setup: func(r *MockRepository) {
				org := &Organization{
					Id:   existingID,
					Name: "Test Org",
				}
				r.organizations[existingID] = org
			},
			wantErr: false,
		},
		{
			name:    "non-existent organization",
			orgID:   uuid.New(),
			setup:   func(_ *MockRepository) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, logger.NewTest())

			err := service.DeleteOrganization(context.Background(), tt.orgID)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteOrganization() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if _, exists := repo.organizations[tt.orgID]; exists {
					t.Error("DeleteOrganization() organization still exists")
				}
			}
		})
	}
}

func TestService_CreateMember(t *testing.T) {
	tests := []struct {
		name    string
		req     *CreateMemberRequest
		orgID   string
		wantErr bool
	}{
		{
			name: "successful creation",
			req: &CreateMemberRequest{
				Role: CreateMemberJSONBodyRole("member"),
			},
			orgID:   uuid.New().String(),
			wantErr: false,
		},
		{
			name: "admin role",
			req: &CreateMemberRequest{
				Role: CreateMemberJSONBodyRole("admin"),
			},
			orgID:   uuid.New().String(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			service := NewService(repo, logger.NewTest())

			member, err := service.CreateMember(context.Background(), tt.req, tt.orgID)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateMember() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && member == nil {
				t.Error("CreateMember() returned nil member")
			}
		})
	}
}

func TestService_UpdateMember(t *testing.T) {
	existingID := uuid.New()

	tests := []struct {
		name     string
		memberID uuid.UUID
		req      *UpdateMemberRequest
		setup    func(*MockRepository)
		wantErr  bool
	}{
		{
			name:     "successful update",
			memberID: existingID,
			req: &UpdateMemberRequest{
				Role: UpdateMemberJSONBodyRoleAdmin,
			},
			setup: func(r *MockRepository) {
				member := &Member{
					Id:             existingID,
					OrganizationId: uuid.New().String(),
					Role:           MemberRoleMember,
				}
				r.members[existingID] = member
			},
			wantErr: false,
		},
		{
			name:     "non-existent member",
			memberID: uuid.New(),
			req: &UpdateMemberRequest{
				Role: UpdateMemberJSONBodyRoleAdmin,
			},
			setup:   func(_ *MockRepository) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, logger.NewTest())

			member, err := service.UpdateMember(context.Background(), tt.memberID, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateMember() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && member == nil {
				t.Error("UpdateMember() returned nil member")
			}
		})
	}
}

func TestService_CreateInvitation(t *testing.T) {
	tests := []struct {
		name      string
		req       *CreateInvitationRequest
		orgID     string
		inviterID string
		wantErr   bool
	}{
		{
			name: "successful creation",
			req: &CreateInvitationRequest{
				Email: "invite@example.com",
				Role:  CreateInvitationJSONBodyRoleMember,
			},
			orgID:     uuid.New().String(),
			inviterID: "user-123",
			wantErr:   false,
		},
		{
			name: "admin invitation",
			req: &CreateInvitationRequest{
				Email: "admin@example.com",
				Role:  CreateInvitationJSONBodyRoleAdmin,
			},
			orgID:     uuid.New().String(),
			inviterID: "user-456",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			service := NewService(repo, logger.NewTest())

			invitation, err := service.CreateInvitation(context.Background(), tt.req, tt.orgID, tt.inviterID)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateInvitation() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if invitation == nil {
					t.Error("CreateInvitation() returned nil invitation")
				} else {
					if invitation.Status != "pending" {
						t.Errorf("CreateInvitation() status = %v, want pending", invitation.Status)
					}
					if invitation.Email != tt.req.Email {
						t.Errorf("CreateInvitation() email = %v, want %v", invitation.Email, tt.req.Email)
					}
				}
			}
		})
	}
}

func TestService_AcceptInvitation(t *testing.T) {
	validID := uuid.New()
	expiredID := uuid.New()

	tests := []struct {
		name    string
		invID   uuid.UUID
		userID  string
		setup   func(*MockRepository)
		wantErr bool
	}{
		{
			name:   "successful acceptance",
			invID:  validID,
			userID: "user-123",
			setup: func(r *MockRepository) {
				inv := &Invitation{
					Id:             validID,
					OrganizationId: uuid.New().String(),
					Email:          "test@example.com",
					Role:           InvitationRoleMember,
					Status:         "pending",
					ExpiresAt:      time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				}
				r.invitations[validID] = inv
			},
			wantErr: false,
		},
		{
			name:   "expired invitation",
			invID:  expiredID,
			userID: "user-456",
			setup: func(r *MockRepository) {
				inv := &Invitation{
					Id:             expiredID,
					OrganizationId: uuid.New().String(),
					Email:          "expired@example.com",
					Role:           InvitationRoleMember,
					Status:         "pending",
					ExpiresAt:      time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
				}
				r.invitations[expiredID] = inv
			},
			wantErr: true,
		},
		{
			name:    "non-existent invitation",
			invID:   uuid.New(),
			userID:  "user-789",
			setup:   func(_ *MockRepository) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, logger.NewTest())

			member, err := service.AcceptInvitation(context.Background(), tt.invID, tt.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("AcceptInvitation() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && member == nil {
				t.Error("AcceptInvitation() returned nil member")
			}
		})
	}
}

func TestService_ListOrganizations(t *testing.T) {
	tests := []struct {
		name    string
		limit   int
		offset  int
		setup   func(*MockRepository)
		wantErr bool
	}{
		{
			name:   "list with organizations",
			limit:  10,
			offset: 0,
			setup: func(r *MockRepository) {
				for i := 0; i < 5; i++ {
					org := &Organization{
						Id:           uuid.New(),
						Name:         "Test Org",
						BillingEmail: "test@example.com",
						Plan:         OrganizationPlan(DefaultPlan),
					}
					r.organizations[org.Id] = org
				}
			},
			wantErr: false,
		},
		{
			name:    "empty list",
			limit:   10,
			offset:  0,
			setup:   func(_ *MockRepository) {},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, logger.NewTest())

			orgs, total, err := service.ListOrganizations(context.Background(), tt.limit, tt.offset)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListOrganizations() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if orgs == nil {
					t.Error("ListOrganizations() returned nil organizations")
				}
				if total != len(orgs) {
					t.Errorf("ListOrganizations() total = %v, want %v", total, len(orgs))
				}
			}
		})
	}
}

// TestService_GetMember tests the GetMember method
func TestService_GetMember(t *testing.T) {
	existingID := uuid.New()

	tests := []struct {
		name     string
		memberID uuid.UUID
		setup    func(*MockRepository)
		wantErr  bool
	}{
		{
			name:     "existing member",
			memberID: existingID,
			setup: func(r *MockRepository) {
				member := &Member{
					Id:             existingID,
					OrganizationId: uuid.New().String(),
					UserId:         "user-123",
					Role:           MemberRole("member"),
				}
				r.members[existingID] = member
			},
			wantErr: false,
		},
		{
			name:     "non-existent member",
			memberID: uuid.New(),
			setup:    func(_ *MockRepository) {},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, logger.NewTest())

			member, err := service.GetMember(context.Background(), tt.memberID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetMember() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && member == nil {
				t.Error("GetMember() returned nil member")
			}
		})
	}
}

// TestService_DeleteMember tests the DeleteMember method
func TestService_DeleteMember(t *testing.T) {
	existingID := uuid.New()

	tests := []struct {
		name     string
		memberID uuid.UUID
		setup    func(*MockRepository)
		wantErr  bool
	}{
		{
			name:     "delete existing member",
			memberID: existingID,
			setup: func(r *MockRepository) {
				member := &Member{
					Id:             existingID,
					OrganizationId: uuid.New().String(),
					UserId:         "user-123",
					Role:           MemberRole("member"),
				}
				r.members[existingID] = member
			},
			wantErr: false,
		},
		{
			name:     "delete non-existent member",
			memberID: uuid.New(),
			setup:    func(_ *MockRepository) {},
			wantErr:  true,
		},
		{
			name:     "repository error",
			memberID: existingID,
			setup: func(r *MockRepository) {
				r.err = errors.New("database error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, logger.NewTest())

			err := service.DeleteMember(context.Background(), tt.memberID)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteMember() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestService_ListMembers tests the ListMembers method
func TestService_ListMembers(t *testing.T) {
	orgID := uuid.New()

	tests := []struct {
		name    string
		orgID   uuid.UUID
		limit   int
		offset  int
		setup   func(*MockRepository)
		wantErr bool
		wantLen int
	}{
		{
			name:   "list members for organization",
			orgID:  orgID,
			limit:  10,
			offset: 0,
			setup: func(r *MockRepository) {
				for i := 0; i < 3; i++ {
					member := &Member{
						Id:             uuid.New(),
						OrganizationId: orgID.String(),
						UserId:         "user-" + string(rune('0'+i)),
						Role:           MemberRole("member"),
					}
					r.members[member.Id] = member
				}
			},
			wantErr: false,
			wantLen: 3,
		},
		{
			name:    "empty member list",
			orgID:   uuid.New(),
			limit:   10,
			offset:  0,
			setup:   func(_ *MockRepository) {},
			wantErr: false,
			wantLen: 0,
		},
		{
			name:   "with pagination",
			orgID:  orgID,
			limit:  2,
			offset: 1,
			setup: func(r *MockRepository) {
				for i := 0; i < 5; i++ {
					member := &Member{
						Id:             uuid.New(),
						OrganizationId: orgID.String(),
						UserId:         "user-" + string(rune('0'+i)),
						Role:           MemberRole("member"),
					}
					r.members[member.Id] = member
				}
			},
			wantErr: false,
			wantLen: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, logger.NewTest())

			members, total, err := service.ListMembers(context.Background(), tt.orgID.String(), tt.limit, tt.offset)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListMembers() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if members == nil {
					t.Error("ListMembers() returned nil members")
				}
				if total != tt.wantLen {
					t.Errorf("ListMembers() total = %v, want %v", total, tt.wantLen)
				}
			}
		})
	}
}

// TestService_GetInvitation tests the GetInvitation method
func TestService_GetInvitation(t *testing.T) {
	existingID := uuid.New()

	tests := []struct {
		name    string
		invID   uuid.UUID
		setup   func(*MockRepository)
		wantErr bool
	}{
		{
			name:  "existing invitation",
			invID: existingID,
			setup: func(r *MockRepository) {
				inv := &Invitation{
					Id:             existingID,
					OrganizationId: uuid.New().String(),
					Email:          "test@example.com",
					Role:           InvitationRoleMember,
					Status:         "pending",
					ExpiresAt:      time.Now().Add(24 * time.Hour).Format(time.RFC3339),
				}
				r.invitations[existingID] = inv
			},
			wantErr: false,
		},
		{
			name:    "non-existent invitation",
			invID:   uuid.New(),
			setup:   func(_ *MockRepository) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, logger.NewTest())

			inv, err := service.GetInvitation(context.Background(), tt.invID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetInvitation() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && inv == nil {
				t.Error("GetInvitation() returned nil invitation")
			}
		})
	}
}

// TestService_DeleteInvitation tests the DeleteInvitation method
func TestService_DeleteInvitation(t *testing.T) {
	existingID := uuid.New()

	tests := []struct {
		name    string
		invID   uuid.UUID
		setup   func(*MockRepository)
		wantErr bool
	}{
		{
			name:  "delete existing invitation",
			invID: existingID,
			setup: func(r *MockRepository) {
				inv := &Invitation{
					Id:             existingID,
					OrganizationId: uuid.New().String(),
					Email:          "test@example.com",
					Role:           InvitationRoleMember,
					Status:         "pending",
				}
				r.invitations[existingID] = inv
			},
			wantErr: false,
		},
		{
			name:    "delete non-existent invitation",
			invID:   uuid.New(),
			setup:   func(_ *MockRepository) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, logger.NewTest())

			err := service.DeleteInvitation(context.Background(), tt.invID)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteInvitation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestService_ListInvitations tests the ListInvitations method
func TestService_ListInvitations(t *testing.T) {
	orgID := uuid.New()

	tests := []struct {
		name    string
		orgID   string
		limit   int
		offset  int
		setup   func(*MockRepository)
		wantErr bool
	}{
		{
			name:   "list invitations",
			orgID:  orgID.String(),
			limit:  10,
			offset: 0,
			setup: func(r *MockRepository) {
				for i := 0; i < 3; i++ {
					inv := &Invitation{
						Id:             uuid.New(),
						OrganizationId: orgID.String(),
						Email:          "test" + string(rune('0'+i)) + "@example.com",
						Role:           InvitationRoleMember,
						Status:         "pending",
					}
					r.invitations[inv.Id] = inv
				}
			},
			wantErr: false,
		},
		{
			name:    "empty invitation list",
			orgID:   uuid.New().String(),
			limit:   10,
			offset:  0,
			setup:   func(_ *MockRepository) {},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, logger.NewTest())

			invs, total, err := service.ListInvitations(context.Background(), tt.orgID, tt.limit, tt.offset)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListInvitations() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if invs == nil {
					t.Error("ListInvitations() returned nil invitations")
				}
				if total < 0 {
					t.Errorf("ListInvitations() total = %v, want >= 0", total)
				}
			}
		})
	}
}

// TestNewService tests the service constructor
func TestNewService(t *testing.T) {
	repo := NewMockRepository()
	logger := logger.NewTest()

	service := NewService(repo, logger)

	if service == nil {
		t.Fatal("NewService() returned nil")
	}

	if service.repo == nil {
		t.Error("NewService() service.repo is nil")
	}

	if service.logger == nil {
		t.Error("NewService() service.logger is nil")
	}
}

// TestMockRepository_EdgeCases tests edge cases in the mock repository
func TestMockRepository_EdgeCases(t *testing.T) {
	t.Run("ListOrganizations with nil params", func(t *testing.T) {
		repo := NewMockRepository()
		orgs, total, err := repo.ListOrganizations(context.Background(), ListOrganizationsParams{})

		if err != nil {
			t.Errorf("ListOrganizations() unexpected error: %v", err)
		}

		if orgs == nil {
			t.Error("ListOrganizations() returned nil slice")
		}

		if total != 0 {
			t.Errorf("ListOrganizations() total = %v, want 0", total)
		}
	})

	t.Run("UpdateOrganization on non-existent", func(t *testing.T) {
		repo := NewMockRepository()
		org := &Organization{
			Name: "Updated",
		}

		_, err := repo.UpdateOrganization(context.Background(), uuid.New(), org)

		if !errors.Is(err, ErrOrganizationNotFound) {
			t.Errorf("UpdateOrganization() error = %v, want %v", err, ErrOrganizationNotFound)
		}
	})

	t.Run("GetMemberByID on non-existent", func(t *testing.T) {
		repo := NewMockRepository()

		_, err := repo.GetMember(context.Background(), uuid.New())

		if !errors.Is(err, ErrMemberNotFound) {
			t.Errorf("GetMemberByID() error = %v, want %v", err, ErrMemberNotFound)
		}
	})
}
