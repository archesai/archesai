// Package organizations provides organization management functionality including
// organization CRUD operations, member management, and invitation handling.
package organizations

//go:generate go tool oapi-codegen --config=../../.types.codegen.yaml --package organizations --include-tags Organizations ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.server.codegen.yaml --package organizations --include-tags Organizations ../../api/openapi.bundled.yaml

import (
	"errors"

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

	// ErrInvitationNotPending is returned when an invitation is not in pending status
	ErrInvitationNotPending = errors.New("invitation is not pending")

	// ErrInvitationExpired is returned when an invitation has expired
	ErrInvitationExpired = errors.New("invitation expired")
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
