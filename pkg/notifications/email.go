package notifications

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/archesai/archesai/pkg/auth/schemas"
)

// EmailSender interface for sending emails.
type EmailSender interface {
	Send(to, subject, body string) error
}

// EmailDeliverer sends magic links via email.
type EmailDeliverer struct {
	sender EmailSender
}

// NewEmailDeliverer creates an email deliverer.
func NewEmailDeliverer(sender EmailSender) *EmailDeliverer {
	return &EmailDeliverer{
		sender: sender,
	}
}

// Deliver sends the magic link via email.
func (d *EmailDeliverer) Deliver(
	_ context.Context,
	token *schemas.MagicLinkToken,
	baseURL string,
) error {
	magicLink := fmt.Sprintf("%s/auth/magic-link/verify?token=%s", baseURL, *token.Token)

	subject := "Sign in to Arches"
	body := fmt.Sprintf(`
Hi there,

Click the link below to sign in to Arches:

%s

This link will expire in %v.

If you didn't request this, you can safely ignore this email.

Best,
The Arches Team
`, magicLink, time.Until(token.ExpiresAt).Round(time.Minute))

	if err := d.sender.Send(token.Identifier, subject, body); err != nil {
		return fmt.Errorf("sending email: %w", err)
	}

	slog.Info("Magic link sent via email",
		"identifier", token.Identifier,
		"expires_in", time.Until(token.ExpiresAt).Round(time.Second),
	)

	return nil
}
