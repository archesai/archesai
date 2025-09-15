package members

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

// Service handles member business logic
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService creates a new members service
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// ListOrganizationMembers returns all members for an organization
func (s *Service) ListOrganizationMembers(ctx context.Context, organizationID uuid.UUID) ([]*Member, error) {
	// Get all members and filter by organization
	// TODO: Add organization filter to ListMembersParams when available
	params := ListMembersParams{
		Page: PageQuery{
			Number: 1,
			Size:   1000, // Get a reasonable number of members
		},
	}
	allMembers, _, err := s.repo.List(ctx, params)
	if err != nil {
		s.logger.Error("failed to list organization members", "error", err, "organization_id", organizationID)
		return nil, fmt.Errorf("failed to list organization members: %w", err)
	}

	// Filter members by organization ID
	var members []*Member
	orgIDStr := organizationID.String()
	for _, member := range allMembers {
		if member.OrganizationId == orgIDStr {
			members = append(members, member)
		}
	}

	return members, nil
}

// GetMember returns a specific member
func (s *Service) GetMember(ctx context.Context, memberID uuid.UUID) (*Member, error) {
	member, err := s.repo.Get(ctx, memberID)
	if err != nil {
		s.logger.Error("failed to get member", "error", err, "member_id", memberID)
		return nil, fmt.Errorf("failed to get member: %w", err)
	}

	return member, nil
}

// CreateMember creates a new organization member
func (s *Service) CreateMember(ctx context.Context, organizationID uuid.UUID, userID uuid.UUID, role MemberRole) (*Member, error) {
	member := &Member{
		Id:             uuid.New(),
		OrganizationId: organizationID.String(),
		UserId:         userID.String(),
		Role:           role,
	}

	createdMember, err := s.repo.Create(ctx, member)
	if err != nil {
		s.logger.Error("failed to create member", "error", err, "organization_id", organizationID, "user_id", userID)
		return nil, fmt.Errorf("failed to create member: %w", err)
	}

	return createdMember, nil
}

// UpdateMemberRole updates a member's role
func (s *Service) UpdateMemberRole(ctx context.Context, memberID uuid.UUID, role MemberRole) (*Member, error) {
	member, err := s.repo.Get(ctx, memberID)
	if err != nil {
		s.logger.Error("failed to get member for update", "error", err, "member_id", memberID)
		return nil, fmt.Errorf("failed to get member: %w", err)
	}

	member.Role = role
	updatedMember, err := s.repo.Update(ctx, memberID, member)
	if err != nil {
		s.logger.Error("failed to update member role", "error", err, "member_id", memberID)
		return nil, fmt.Errorf("failed to update member role: %w", err)
	}

	return updatedMember, nil
}

// RemoveMember removes a member from an organization
func (s *Service) RemoveMember(ctx context.Context, memberID uuid.UUID) error {
	if err := s.repo.Delete(ctx, memberID); err != nil {
		s.logger.Error("failed to remove member", "error", err, "member_id", memberID)
		return fmt.Errorf("failed to remove member: %w", err)
	}

	return nil
}

// GetMemberByUserAndOrganization returns a member by user and organization IDs
func (s *Service) GetMemberByUserAndOrganization(ctx context.Context, userID uuid.UUID, organizationID uuid.UUID) (*Member, error) {
	// List all members and find the matching one
	// TODO: Add a specific repository method for this query
	params := ListMembersParams{
		Page: PageQuery{
			Number: 1,
			Size:   1000,
		},
	}
	allMembers, _, err := s.repo.List(ctx, params)
	if err != nil {
		s.logger.Error("failed to list members", "error", err)
		return nil, fmt.Errorf("failed to get member: %w", err)
	}

	userIDStr := userID.String()
	orgIDStr := organizationID.String()
	for _, member := range allMembers {
		if member.UserId == userIDStr && member.OrganizationId == orgIDStr {
			return member, nil
		}
	}

	return nil, fmt.Errorf("member not found for user %s in organization %s", userID, organizationID)
}
