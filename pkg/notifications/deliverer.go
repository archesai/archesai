// Package notifications provides notification delivery implementations
package notifications

import (
	"context"

	"github.com/archesai/archesai/pkg/auth"
)

// Deliverer handles notification delivery for magic links and OTPs.
type Deliverer interface {
	Deliver(ctx context.Context, token *auth.MagicLinkToken, baseURL string) error
}
