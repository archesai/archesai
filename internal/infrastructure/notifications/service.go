// Package notifications provides notification delivery implementations
package notifications

import (
	"context"

	"github.com/archesai/archesai/internal/core/valueobjects"
)

// Deliverer handles notification delivery for magic links and OTPs.
type Deliverer interface {
	Deliver(ctx context.Context, token *valueobjects.MagicLinkToken, baseURL string) error
}
