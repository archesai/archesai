package organizations

import (
	"errors"
)

// Domain errors.
var (
	// ErrOrganizationNotFound is returned when an organization is not found.
	ErrOrganizationNotFound = errors.New("organization not found")

	// ErrMemberExists is returned when a member already exists.
	ErrMemberExists = errors.New("member already exists")

	// ErrMemberNotFound is returned when a member is not found.
	ErrMemberNotFound = errors.New("member not found")

	// ErrInvitationNotFound is returned when an invitation is not found.
	ErrInvitationNotFound = errors.New("invitation not found")

	// ErrInvitationNotPending is returned when an invitation is not in pending status.
	ErrInvitationNotPending = errors.New("invitation is not pending")

	// ErrInvitationExpired is returned when an invitation has expired.
	ErrInvitationExpired = errors.New("invitation expired")
)
