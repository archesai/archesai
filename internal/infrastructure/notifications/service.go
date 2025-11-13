// Package notifications provides notification delivery implementations
package notifications

import (
	"context"

	"github.com/archesai/archesai/internal/core/models"
)

// Deliverer handles notification delivery for magic links and OTPs.
type Deliverer interface {
	Deliver(ctx context.Context, token *models.MagicLinkToken, baseURL string) error
}
