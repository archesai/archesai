// Package organizations provides organization management functionality.
package organizations

import (
	"errors"
	"time"

	"github.com/archesai/archesai/internal/generated/api"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Domain errors
var (
	ErrOrganizationNotFound    = errors.New("organization not found")
	ErrOrganizationExists      = errors.New("organization already exists")
	ErrMemberNotFound          = errors.New("member not found")
	ErrMemberExists            = errors.New("member already exists")
	ErrInvitationNotFound      = errors.New("invitation not found")
	ErrInvitationExists        = errors.New("invitation already exists")
	ErrInvitationExpired       = errors.New("invitation expired")
	ErrInvalidRole             = errors.New("invalid role")
	ErrCannotRemoveOwner       = errors.New("cannot remove organization owner")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
)

// Organization extends the generated API OrganizationEntity with domain-specific fields
type Organization struct {
	api.OrganizationEntity
	// Add any domain-specific fields that aren't in the API
}

// Member extends the generated API MemberEntity with domain-specific fields
type Member struct {
	api.MemberEntity
	// Add any domain-specific fields that aren't in the API
}

// Invitation extends the generated API InvitationEntity with domain-specific fields
type Invitation struct {
	api.InvitationEntity
	// Add any domain-specific fields that aren't in the API
}

// CreateOrganizationRequest represents a request to create an organization
type CreateOrganizationRequest struct {
	Name         string                     `json:"name" validate:"required,min=1,max=100"`
	BillingEmail openapi_types.Email        `json:"billing_email" validate:"required,email"`
	Plan         api.OrganizationEntityPlan `json:"plan,omitempty"`
}

// UpdateOrganizationRequest represents a request to update an organization
type UpdateOrganizationRequest struct {
	Name         *string                     `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	BillingEmail *openapi_types.Email        `json:"billing_email,omitempty" validate:"omitempty,email"`
	Plan         *api.OrganizationEntityPlan `json:"plan,omitempty"`
}

// CreateMemberRequest represents a request to add a member
type CreateMemberRequest struct {
	UserID openapi_types.UUID   `json:"user_id" validate:"required"`
	Role   api.MemberEntityRole `json:"role" validate:"required"`
}

// UpdateMemberRequest represents a request to update a member
type UpdateMemberRequest struct {
	Role *api.MemberEntityRole `json:"role,omitempty" validate:"omitempty"`
}

// CreateInvitationRequest represents a request to create an invitation
type CreateInvitationRequest struct {
	Email string                   `json:"email" validate:"required,email"`
	Role  api.InvitationEntityRole `json:"role" validate:"required"`
}

// NewOrganization creates a new organization from the API entity
func NewOrganization(entity api.OrganizationEntity) *Organization {
	return &Organization{OrganizationEntity: entity}
}

// NewMember creates a new member from the API entity
func NewMember(entity api.MemberEntity) *Member {
	return &Member{MemberEntity: entity}
}

// NewInvitation creates a new invitation from the API entity
func NewInvitation(entity api.InvitationEntity) *Invitation {
	return &Invitation{InvitationEntity: entity}
}

// IsExpired checks if an invitation has expired
func (i *Invitation) IsExpired() bool {
	if i.ExpiresAt == "" {
		return false
	}
	expiresAt, err := time.Parse(time.RFC3339, i.ExpiresAt)
	if err != nil {
		return true
	}
	return time.Now().After(expiresAt)
}

// HasPermission checks if a member has sufficient permissions for an action
func (m *Member) HasPermission(requiredRole api.MemberEntityRole) bool {
	// Role hierarchy: owner > admin > member
	memberWeight := getRoleWeight(m.Role)
	requiredWeight := getRoleWeight(requiredRole)
	return memberWeight >= requiredWeight
}

// getRoleWeight returns a numeric weight for role comparison
func getRoleWeight(role api.MemberEntityRole) int {
	switch role {
	case api.MemberEntityRoleOwner:
		return 3
	case api.MemberEntityRoleAdmin:
		return 2
	case api.MemberEntityRoleMember:
		return 1
	default:
		return 0
	}
}
