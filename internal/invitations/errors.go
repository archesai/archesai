package invitations

import "errors"

var (
	// ErrInvitationNotFound is returned when an invitation is not found
	ErrInvitationNotFound = errors.New("invitation not found")

	// ErrInvitationExpired is returned when an invitation has expired
	ErrInvitationExpired = errors.New("invitation expired")

	// ErrInvitationAlreadyAccepted is returned when trying to accept an already accepted invitation
	ErrInvitationAlreadyAccepted = errors.New("invitation already accepted")

	// ErrInvitationAlreadyDeclined is returned when trying to accept a declined invitation
	ErrInvitationAlreadyDeclined = errors.New("invitation already declined")

	// ErrInvitationAlreadyExists is returned when an invitation already exists for the email
	ErrInvitationAlreadyExists = errors.New("invitation already exists for this email")

	// ErrInvalidInvitationStatus is returned when an invalid status is provided
	ErrInvalidInvitationStatus = errors.New("invalid invitation status")
)
