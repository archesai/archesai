// Package notifications provides notification delivery implementations
package notifications

import (
	"context"

	"github.com/archesai/archesai/pkg/auth/schemas"
)

// Deliverer handles notification delivery for magic links and OTPs.
type Deliverer interface {
	Deliver(ctx context.Context, token *schemas.MagicLinkToken, baseURL string) error
}
