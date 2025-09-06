// Package organizations provides organization management functionality including
// organization CRUD operations, member management, and invitation handling.
package organizations

//go:generate go tool oapi-codegen --config=models.cfg.yaml ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=server.cfg.yaml ../../api/openapi.bundled.yaml

import "errors"

// Domain types
type (
	// Organization represents an organization with its entity
	Organization struct {
		OrganizationEntity
	}

	// Member represents a member with its entity
	Member struct {
		MemberEntity
	}

	// Invitation represents an invitation with its entity
	Invitation struct {
		InvitationEntity
	}

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

// Domain errors
var (
	// ErrOrganizationNotFound is returned when an organization is not found
	ErrOrganizationNotFound = errors.New("organization not found")

	// ErrMemberExists is returned when a member already exists
	ErrMemberExists = errors.New("member already exists")

	// ErrMemberNotFound is returned when a member is not found
	ErrMemberNotFound = errors.New("member not found")

	// ErrInvitationNotFound is returned when an invitation is not found
	ErrInvitationNotFound = errors.New("invitation not found")
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
