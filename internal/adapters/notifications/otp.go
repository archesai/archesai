package notifications

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/archesai/archesai/internal/core/valueobjects"
)

// OTPDeliverer displays OTP codes (for development).
type OTPDeliverer struct {
	logger *slog.Logger
}

// NewOTPDeliverer creates a new OTP deliverer.
func NewOTPDeliverer(logger *slog.Logger) *OTPDeliverer {
	return &OTPDeliverer{
		logger: logger,
	}
}

// Deliver displays the OTP code.
func (d *OTPDeliverer) Deliver(
	_ context.Context,
	token *valueobjects.MagicLinkToken,
	_ string,
) error {
	d.logger.Info("🔐 OTP Code Generated",
		"identifier", token.Identifier,
		"code", token.Code,
		"expires_in", time.Until(token.ExpiresAt).Round(time.Second),
	)

	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🔐 ONE-TIME PASSWORD")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("For: %s\n", token.Identifier)
	fmt.Printf("Code: %s\n", token.Code)
	fmt.Printf("Expires in: %v\n", time.Until(token.ExpiresAt).Round(time.Second))
	fmt.Println(strings.Repeat("=", 80) + "\n")

	return nil
}
