package organizations

import (
	"context"
)

// ExtendedRepository adds custom methods to the generated Repository interface.
// This interface is what the service actually uses.
type ExtendedRepository interface {
	Repository

	// Custom member lookup method
	GetMemberByUserAndOrg(ctx context.Context, userID, orgID string) (*Member, error)
}
