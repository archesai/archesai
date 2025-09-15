// Package organizations provides organization management functionality including
// organization CRUD operations, member management, and invitation handling.
package organizations

// Domain types
type (
	// OrganizationAlias is an alias to avoid conflicts with generated type
	OrganizationAlias = Organization

	// MemberAlias is an alias to avoid conflicts with generated type
	MemberAlias = Member

	// InvitationAlias is an alias to avoid conflicts with generated type
	InvitationAlias = Invitation

	// CreateOrganizationRequest represents a request to create an organization
	CreateOrganizationRequest = CreateOrganizationJSONBody

	// UpdateOrganizationRequest represents a request to update an organization
	UpdateOrganizationRequest = UpdateOrganizationJSONBody

	// CreateMemberRequest represents a request to create a member
	CreateMemberRequest = CreateMemberJSONBody

	// UpdateMemberRequest represents a request to update a member
	UpdateMemberRequest = UpdateMemberJSONBody

	// CreateInvitationRequest represents a request to create an invitation
	CreateInvitationRequest = CreateInvitationJSONBody
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
