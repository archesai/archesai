// Package deliverers provides delivery implementations for authentication services
package notifications

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/archesai/archesai/internal/core/valueobjects"
)

// ConsoleDeliverer prints magic links to console (for development).
type ConsoleDeliverer struct {
	logger *slog.Logger
}

// NewConsoleDeliverer creates a new console deliverer.
func NewConsoleDeliverer(logger *slog.Logger) *ConsoleDeliverer {
	return &ConsoleDeliverer{
		logger: logger,
	}
}

// Deliver prints the magic link to console.
func (d *ConsoleDeliverer) Deliver(
	_ context.Context,
	token *valueobjects.MagicLinkToken,
	baseURL string,
) error {
	magicLink := fmt.Sprintf("%s/auth/magic-link/verify?token=%s", baseURL, *token.Token)

	d.logger.Info("🔮 Magic Link Generated",
		"identifier", token.Identifier,
		"link", magicLink,
		"expires_in", time.Until(token.ExpiresAt).Round(time.Second),
	)

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔮 MAGIC LINK AUTHENTICATION")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("For: %s\n", token.Identifier)
	fmt.Printf("Link: %s\n", magicLink)
	fmt.Printf("Expires in: %v\n", time.Until(token.ExpiresAt).Round(time.Second))
	fmt.Println(strings.Repeat("=", 80) + "\n")

	return nil
}
