package organizations

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Service provides organization business logic
type Service struct {
	repo   OrganizationRepository
	logger *slog.Logger
}

// NewService creates a new organization service
func NewService(repo OrganizationRepository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// CreateOrganization creates a new organization
func (s *Service) CreateOrganization(ctx context.Context, req *CreateOrganizationRequest, creatorUserID string) (*Organization, error) {
	s.logger.Debug("creating organization", "id", req.OrganizationId, "creator", creatorUserID)

	// Set default plan
	plan := OrganizationEntityPlan(DefaultPlan)

	org := &Organization{
		OrganizationEntity: OrganizationEntity{
			Id:           req.OrganizationId,
			Name:         "", // Name should be set from somewhere else
			BillingEmail: openapi_types.Email(req.BillingEmail),
			Plan:         plan,
			Credits:      0.0, // Start with 0 credits
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	createdOrg, err := s.repo.CreateOrganization(ctx, org)
	if err != nil {
		s.logger.Error("failed to create organization", "error", err)
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	// Create initial owner member
	member := &Member{
		MemberEntity: MemberEntity{
			Id: uuid.UUID{}, // Will be set by repository
			// Note: MemberEntity doesn't have UserId field in API
			OrganizationId: createdOrg.Id.String(),
			Role:           MemberEntityRoleOwner,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	_, err = s.repo.CreateMember(ctx, member)
	if err != nil {
		s.logger.Error("failed to create initial member", "error", err)
		// Try to clean up the created organization
		if delErr := s.repo.DeleteOrganization(ctx, createdOrg.Id); delErr != nil {
			s.logger.Warn("failed to cleanup organization after member creation failure", "error", delErr)
		}
		return nil, fmt.Errorf("failed to create initial member: %w", err)
	}

	s.logger.Info("organization created successfully", "id", createdOrg.Id, "name", createdOrg.Name)
	return createdOrg, nil
}

// GetOrganization retrieves an organization by ID
func (s *Service) GetOrganization(ctx context.Context, id uuid.UUID) (*Organization, error) {
	org, err := s.repo.GetOrganization(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	return org, nil
}

// UpdateOrganization updates an organization
func (s *Service) UpdateOrganization(ctx context.Context, id uuid.UUID, req *UpdateOrganizationRequest) (*Organization, error) {
	org, err := s.repo.GetOrganization(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}

	// Update fields that were provided
	if req.BillingEmail != "" {
		org.BillingEmail = openapi_types.Email(req.BillingEmail)
	}
	org.UpdatedAt = time.Now()

	updatedOrg, err := s.repo.UpdateOrganization(ctx, org)
	if err != nil {
		return nil, fmt.Errorf("failed to update organization: %w", err)
	}

	return updatedOrg, nil
}

// DeleteOrganization deletes an organization
func (s *Service) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
	// TODO: Add additional checks (e.g., organization has no active resources)
	err := s.repo.DeleteOrganization(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}
	return nil
}

// ListOrganizations retrieves a list of organizations
func (s *Service) ListOrganizations(ctx context.Context, limit, offset int) ([]*Organization, int, error) {
	orgs, total, err := s.repo.ListOrganizations(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list organizations: %w", err)
	}
	return orgs, total, nil
}

// CreateMember adds a member to an organization
func (s *Service) CreateMember(ctx context.Context, req *CreateMemberRequest, orgID string) (*Member, error) {
	// Check if member already exists
	// CreateMemberRequest doesn't have UserID field in the generated types
	// We need to use a different approach
	existing, err := s.repo.GetMemberByUserAndOrg(ctx, "", orgID)
	if err == nil && existing != nil {
		return nil, ErrMemberExists
	}

	member := &Member{
		MemberEntity: MemberEntity{
			Id: uuid.UUID{}, // Will be set by repository
			// Note: MemberEntity doesn't have UserId field in API
			OrganizationId: orgID,
			Role:           MemberEntityRole(req.Role),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	createdMember, err := s.repo.CreateMember(ctx, member)
	if err != nil {
		return nil, fmt.Errorf("failed to create member: %w", err)
	}

	return createdMember, nil
}

// GetMember retrieves a member by ID
func (s *Service) GetMember(ctx context.Context, id uuid.UUID) (*Member, error) {
	member, err := s.repo.GetMember(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get member: %w", err)
	}
	return member, nil
}

// UpdateMember updates a member's role
func (s *Service) UpdateMember(ctx context.Context, id uuid.UUID, req *UpdateMemberRequest) (*Member, error) {
	member, err := s.repo.GetMember(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get member: %w", err)
	}

	// Prevent removing the last owner
	// TODO: Check if this is the last owner in the organization
	// For now, we'll allow it but this should be implemented
	_ = member.Role
	_ = req.Role

	if req.Role != "" {
		member.Role = MemberEntityRole(req.Role)
	}
	member.UpdatedAt = time.Now()

	updatedMember, err := s.repo.UpdateMember(ctx, member)
	if err != nil {
		return nil, fmt.Errorf("failed to update member: %w", err)
	}

	return updatedMember, nil
}

// DeleteMember removes a member from an organization
func (s *Service) DeleteMember(ctx context.Context, id uuid.UUID) error {
	member, err := s.repo.GetMember(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get member: %w", err)
	}

	// Prevent removing the last owner
	// TODO: Check if this is the last owner in the organization
	// For now, we'll allow it but this should be implemented
	_ = member.Role

	err = s.repo.DeleteMember(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete member: %w", err)
	}

	return nil
}

// ListMembers retrieves members of an organization
func (s *Service) ListMembers(ctx context.Context, orgID string, limit, offset int) ([]*Member, int, error) {
	members, total, err := s.repo.ListMembers(ctx, orgID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list members: %w", err)
	}
	return members, total, nil
}

// CreateInvitation creates a new invitation
func (s *Service) CreateInvitation(ctx context.Context, req *CreateInvitationRequest, orgID, inviterID string) (*Invitation, error) {
	expiresAt := time.Now().AddDate(0, 0, 7).Format(time.RFC3339) // 7 days expiry

	invitation := &Invitation{
		InvitationEntity: InvitationEntity{
			Id:             uuid.UUID{}, // Will be set by repository
			OrganizationId: orgID,
			InviterId:      inviterID,
			Email:          req.Email,
			Role:           InvitationEntityRole(req.Role),
			Status:         "pending", // Use string instead of enum
			ExpiresAt:      expiresAt, // String field, not pointer
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	createdInvitation, err := s.repo.CreateInvitation(ctx, invitation)
	if err != nil {
		return nil, fmt.Errorf("failed to create invitation: %w", err)
	}

	// TODO: Send invitation email

	return createdInvitation, nil
}

// AcceptInvitation accepts an invitation and creates a member
func (s *Service) AcceptInvitation(ctx context.Context, id uuid.UUID, _ string) (*Member, error) {
	invitation, err := s.repo.GetInvitation(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get invitation: %w", err)
	}

	if invitation.Status != "pending" {
		return nil, fmt.Errorf("invitation is not pending")
	}

	// Check if invitation is expired
	expiresAt, err := time.Parse(time.RFC3339, invitation.ExpiresAt)
	if err == nil && time.Now().After(expiresAt) {
		return nil, fmt.Errorf("invitation expired")
	}

	// Create member
	member := &Member{
		MemberEntity: MemberEntity{
			Id: uuid.UUID{}, // Will be set by repository
			// Note: MemberEntity doesn't have UserId field
			OrganizationId: invitation.OrganizationId,
			Role:           MemberEntityRole(invitation.Role),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	createdMember, err := s.repo.CreateMember(ctx, member)
	if err != nil {
		return nil, fmt.Errorf("failed to create member: %w", err)
	}

	// Update invitation status
	invitation.Status = "accepted"
	invitation.UpdatedAt = time.Now()
	_, err = s.repo.UpdateInvitation(ctx, invitation)
	if err != nil {
		s.logger.Warn("failed to update invitation status", "error", err)
		// Don't fail the whole operation for this
	}

	return createdMember, nil
}

// GetInvitation retrieves a single invitation by ID
func (s *Service) GetInvitation(ctx context.Context, id uuid.UUID) (*Invitation, error) {
	invitation, err := s.repo.GetInvitation(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get invitation: %w", err)
	}
	return invitation, nil
}

// DeleteInvitation removes an invitation
func (s *Service) DeleteInvitation(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteInvitation(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete invitation: %w", err)
	}
	return nil
}

// ListInvitations retrieves invitations for an organization
func (s *Service) ListInvitations(ctx context.Context, orgID string, limit, offset int) ([]*Invitation, int, error) {
	invitations, total, err := s.repo.ListInvitations(ctx, orgID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list invitations: %w", err)
	}
	return invitations, total, nil
}
