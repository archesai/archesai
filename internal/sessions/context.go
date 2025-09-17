package sessions

import (
	"context"

	"github.com/google/uuid"
)

// GetAuthContextFromGoContext retrieves authentication data from standard Go context
// This is useful for service and repository layers that don't have access to Echo context.
func GetAuthContextFromGoContext(ctx context.Context) (*EnhancedClaims, uuid.UUID, bool) {
	// Try to get enhanced claims first
	if claims, ok := ctx.Value("auth.claims").(*EnhancedClaims); ok {
		return claims, claims.UserID, true
	}

	// Fallback to getting user ID directly
	if userID, ok := ctx.Value("auth.user_id").(uuid.UUID); ok {
		return nil, userID, true
	}

	return nil, uuid.Nil, false
}
