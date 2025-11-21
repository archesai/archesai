// Package notifications provides notification delivery implementations
package notifications

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/archesai/archesai/apps/studio/generated/core/models"
)

// ConsoleDeliverer prints magic links to console (for development).
type ConsoleDeliverer struct {
}

// NewConsoleDeliverer creates a new console deliverer.
func NewConsoleDeliverer() *ConsoleDeliverer {
	return &ConsoleDeliverer{}
}

// Deliver prints the magic link to console.
func (d *ConsoleDeliverer) Deliver(
	_ context.Context,
	token *models.MagicLinkToken,
	baseURL string,
) error {
	magicLink := fmt.Sprintf("%s/auth/magic-link/verify?token=%s", baseURL, *token.Token)

	slog.Info("🔮 Magic Link Generated",
		"identifier", token.Identifier,
		"link", magicLink,
		"expires_in", time.Until(token.ExpiresAt).Round(time.Second),
	)

	slog.Info(fmt.Sprintln("\n" + strings.Repeat("=", 80)))
	slog.Info(fmt.Sprintln("🔮 MAGIC LINK AUTHENTICATION"))
	slog.Info(fmt.Sprintln(strings.Repeat("-", 80)))
	slog.Info(fmt.Sprintf("For: %s\n", token.Identifier))
	slog.Info(fmt.Sprintf("Link: %s\n", magicLink))
	slog.Info(fmt.Sprintf("Expires in: %v\n", time.Until(token.ExpiresAt).Round(time.Second)))
	slog.Info(fmt.Sprintln(strings.Repeat("=", 80) + "\n"))

	return nil
}
