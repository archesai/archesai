package members

import "errors"

// Member errors
var (
	ErrMemberNotFound = errors.New("member not found")
	// ErrDuplicateMember is returned when a member already exists
	ErrDuplicateMember = errors.New("member already exists")
)
