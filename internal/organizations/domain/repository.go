package domain

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for organization data persistence
type Repository interface {
	// Organization operations
	CreateOrganization(ctx context.Context, org *Organization) (*Organization, error)
	GetOrganization(ctx context.Context, id uuid.UUID) (*Organization, error)
	UpdateOrganization(ctx context.Context, org *Organization) (*Organization, error)
	DeleteOrganization(ctx context.Context, id uuid.UUID) error
	ListOrganizations(ctx context.Context, limit, offset int) ([]*Organization, int, error)

	// Member operations
	CreateMember(ctx context.Context, member *Member) (*Member, error)
	GetMember(ctx context.Context, id uuid.UUID) (*Member, error)
	GetMemberByUserAndOrg(ctx context.Context, userID, orgID string) (*Member, error)
	UpdateMember(ctx context.Context, member *Member) (*Member, error)
	DeleteMember(ctx context.Context, id uuid.UUID) error
	ListMembers(ctx context.Context, orgID string, limit, offset int) ([]*Member, int, error)

	// Invitation operations
	CreateInvitation(ctx context.Context, invitation *Invitation) (*Invitation, error)
	GetInvitation(ctx context.Context, id uuid.UUID) (*Invitation, error)
	UpdateInvitation(ctx context.Context, invitation *Invitation) (*Invitation, error)
	DeleteInvitation(ctx context.Context, id uuid.UUID) error
	ListInvitations(ctx context.Context, orgID string, limit, offset int) ([]*Invitation, int, error)
}
