package notifications

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/archesai/archesai/apis/studio/generated/core/models"
)

// OTPDeliverer displays OTP codes (for development).
type OTPDeliverer struct {
}

// NewOTPDeliverer creates a new OTP deliverer.
func NewOTPDeliverer() *OTPDeliverer {
	return &OTPDeliverer{}
}

// Deliver displays the OTP code.
func (d *OTPDeliverer) Deliver(
	_ context.Context,
	token *models.MagicLinkToken,
	_ string,
) error {
	slog.Info("🔐 OTP Code Generated",
		"identifier", token.Identifier,
		"code", token.Token,
		"expires_in", time.Until(token.ExpiresAt).Round(time.Second),
	)

	slog.Info(fmt.Sprintln("\n" + strings.Repeat("=", 80)))
	slog.Info(fmt.Sprintln("🔐 ONE-TIME PASSWORD"))
	slog.Info(fmt.Sprintln(strings.Repeat("-", 80)))
	slog.Info(fmt.Sprintf("For: %s\n", token.Identifier))
	slog.Info(fmt.Sprintf("Code: %s\n", *token.Token))
	slog.Info(fmt.Sprintf("Expires in: %v\n", time.Until(token.ExpiresAt).Round(time.Second)))
	slog.Info(fmt.Sprintln(strings.Repeat("=", 80) + "\n"))

	return nil
}
