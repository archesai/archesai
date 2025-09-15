// Package organizations provides organization management functionality including
// organization CRUD operations, member management, and invitation handling.
package organizations

import (
	"github.com/archesai/archesai/internal/invitations"
	"github.com/archesai/archesai/internal/members"
)

// Domain types
type (
	// OrganizationAlias is an alias to avoid conflicts with generated type
	OrganizationAlias = Organization

	// MemberAlias is an alias to avoid conflicts with generated type
	MemberAlias = members.Member

	// InvitationAlias is an alias to avoid conflicts with generated type
	InvitationAlias = invitations.Invitation

	// CreateOrganizationRequest represents a request to create an organization
	CreateOrganizationRequest = CreateOrganizationJSONBody

	// UpdateOrganizationRequest represents a request to update an organization
	UpdateOrganizationRequest = UpdateOrganizationJSONBody

	// CreateMemberRequest represents a request to create a member
	CreateMemberRequest = members.CreateMemberJSONBody

	// UpdateMemberRequest represents a request to update a member
	UpdateMemberRequest = members.UpdateMemberJSONBody

	// CreateInvitationRequest represents a request to create an invitation
	CreateInvitationRequest = invitations.CreateInvitationJSONBody
)

// Domain constants
const (
	// DefaultPlan is the default organization plan for new organizations
	DefaultPlan = "free"

	// MaxMembersPerOrganization defines the maximum number of members per organization
	MaxMembersPerOrganization = 100

	// InvitationExpiryDays defines how long invitations remain valid
	InvitationExpiryDays = 7
)
