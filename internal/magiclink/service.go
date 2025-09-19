// Package magiclink provides magic link authentication functionality
package magiclink

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrTokenNotFound is returned when a magic link token cannot be found
	ErrTokenNotFound = errors.New("magic link token not found")
	// ErrTokenExpired is returned when a magic link token has expired
	ErrTokenExpired = errors.New("magic link token has expired")
	// ErrTokenAlreadyUsed is returned when a magic link token has already been used
	ErrTokenAlreadyUsed = errors.New("magic link token has already been used")
	// ErrInvalidOTP is returned when an invalid OTP code is provided
	ErrInvalidOTP = errors.New("invalid OTP code")
	// ErrRateLimitExceeded is returned when the rate limit is exceeded
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
)

// DeliveryMethod represents how the magic link should be delivered
type DeliveryMethod string

const (
	// DeliveryEmail sends magic links via email
	DeliveryEmail DeliveryMethod = "email"
	// DeliveryConsole prints magic links to console/stdout
	DeliveryConsole DeliveryMethod = "console"
	// DeliveryOTP sends one-time passwords
	DeliveryOTP DeliveryMethod = "otp"
	// DeliveryWebhook sends magic links to a webhook
	DeliveryWebhook DeliveryMethod = "webhook"
)

// Token represents a magic link token
type Token struct {
	ID             uuid.UUID
	UserID         *uuid.UUID
	Token          string
	TokenHash      string
	Code           string
	Identifier     string
	DeliveryMethod DeliveryMethod
	ExpiresAt      time.Time
	UsedAt         *time.Time
	IPAddress      string
	UserAgent      string
	CreatedAt      time.Time
}

// Repository defines the interface for magic link storage
type Repository interface {
	Create(ctx context.Context, token *Token) error
	GetByToken(ctx context.Context, tokenHash string) (*Token, error)
	GetByCode(ctx context.Context, identifier string, code string) (*Token, error)
	MarkUsed(ctx context.Context, tokenID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
	CountRecentByIdentifier(ctx context.Context, identifier string, since time.Time) (int, error)
}

// Deliverer defines the interface for sending magic links
type Deliverer interface {
	Deliver(ctx context.Context, token *Token, baseURL string) error
}

// Service handles magic link operations
type Service struct {
	repo       Repository
	logger     *slog.Logger
	deliverers map[DeliveryMethod]Deliverer
	baseURL    string
	ttl        time.Duration
	rateLimit  int
}

// NewService creates a new magic link service
func NewService(repo Repository, logger *slog.Logger, baseURL string) *Service {
	s := &Service{
		repo:       repo,
		logger:     logger,
		baseURL:    baseURL,
		ttl:        15 * time.Minute,
		rateLimit:  5, // 5 requests per hour
		deliverers: make(map[DeliveryMethod]Deliverer),
	}

	// Register default deliverers
	s.deliverers[DeliveryConsole] = &ConsoleDeliverer{logger: logger}
	s.deliverers[DeliveryOTP] = &OTPDeliverer{logger: logger}

	return s
}

// RegisterDeliverer adds a custom deliverer for a delivery method
func (s *Service) RegisterDeliverer(method DeliveryMethod, deliverer Deliverer) {
	s.deliverers[method] = deliverer
}

// RequestMagicLink generates and sends a magic link
func (s *Service) RequestMagicLink(
	ctx context.Context,
	identifier string,
	method DeliveryMethod,
	userID *uuid.UUID,
	ipAddress string,
	userAgent string,
) (*Token, error) {
	// Check rate limit
	since := time.Now().Add(-1 * time.Hour)
	count, err := s.repo.CountRecentByIdentifier(ctx, identifier, since)
	if err != nil {
		return nil, fmt.Errorf("checking rate limit: %w", err)
	}
	if count >= s.rateLimit {
		return nil, ErrRateLimitExceeded
	}

	// Generate secure token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("generating token: %w", err)
	}
	tokenString := hex.EncodeToString(tokenBytes)

	// Hash token for storage
	hash := sha256.Sum256([]byte(tokenString))
	tokenHash := hex.EncodeToString(hash[:])

	// Generate OTP code if needed
	var code string
	if method == DeliveryOTP {
		code = s.generateOTP()
	}

	// Create token record
	token := &Token{
		ID:             uuid.New(),
		UserID:         userID,
		Token:          tokenString,
		TokenHash:      tokenHash,
		Code:           code,
		Identifier:     identifier,
		DeliveryMethod: method,
		ExpiresAt:      time.Now().Add(s.ttl),
		IPAddress:      ipAddress,
		UserAgent:      userAgent,
		CreatedAt:      time.Now(),
	}

	// Store in database
	if err := s.repo.Create(ctx, token); err != nil {
		return nil, fmt.Errorf("storing token: %w", err)
	}

	// Deliver the token
	deliverer, ok := s.deliverers[method]
	if !ok {
		return nil, fmt.Errorf("unsupported delivery method: %s", method)
	}

	if err := deliverer.Deliver(ctx, token, s.baseURL); err != nil {
		return nil, fmt.Errorf("delivering token: %w", err)
	}

	// Don't return the actual token for security
	token.Token = ""
	return token, nil
}

// VerifyToken verifies a magic link token
func (s *Service) VerifyToken(ctx context.Context, token string) (*Token, error) {
	// Hash the provided token
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	// Look up the token
	t, err := s.repo.GetByToken(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, ErrTokenNotFound) {
			return nil, ErrTokenNotFound
		}
		return nil, fmt.Errorf("looking up token: %w", err)
	}

	// Check if already used
	if t.UsedAt != nil {
		return nil, ErrTokenAlreadyUsed
	}

	// Check if expired
	if time.Now().After(t.ExpiresAt) {
		return nil, ErrTokenExpired
	}

	// Mark as used
	if err := s.repo.MarkUsed(ctx, t.ID); err != nil {
		return nil, fmt.Errorf("marking token as used: %w", err)
	}

	return t, nil
}

// VerifyOTP verifies an OTP code
func (s *Service) VerifyOTP(ctx context.Context, identifier string, code string) (*Token, error) {
	// Look up the token by identifier and code
	t, err := s.repo.GetByCode(ctx, identifier, code)
	if err != nil {
		if errors.Is(err, ErrTokenNotFound) {
			return nil, ErrInvalidOTP
		}
		return nil, fmt.Errorf("looking up OTP: %w", err)
	}

	// Check if already used
	if t.UsedAt != nil {
		return nil, ErrTokenAlreadyUsed
	}

	// Check if expired
	if time.Now().After(t.ExpiresAt) {
		return nil, ErrTokenExpired
	}

	// Mark as used
	if err := s.repo.MarkUsed(ctx, t.ID); err != nil {
		return nil, fmt.Errorf("marking OTP as used: %w", err)
	}

	return t, nil
}

// CleanupExpired removes expired tokens
func (s *Service) CleanupExpired(ctx context.Context) error {
	return s.repo.DeleteExpired(ctx)
}

// generateOTP generates a 6-digit OTP code
func (s *Service) generateOTP() string {
	maxNum := big.NewInt(999999)
	n, _ := rand.Int(rand.Reader, maxNum)
	return fmt.Sprintf("%06d", n.Int64())
}

// ConsoleDeliverer prints the magic link to console (for development)
type ConsoleDeliverer struct {
	logger *slog.Logger
}

// Deliver prints the magic link to console
func (d *ConsoleDeliverer) Deliver(_ context.Context, token *Token, baseURL string) error {
	magicLink := fmt.Sprintf("%s/auth/magic-link/verify?token=%s", baseURL, token.Token)

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

// OTPDeliverer displays the OTP code
type OTPDeliverer struct {
	logger *slog.Logger
}

// Deliver displays the OTP code
func (d *OTPDeliverer) Deliver(_ context.Context, token *Token, _ string) error {
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

// EmailDeliverer sends magic links via email
type EmailDeliverer struct {
	sender EmailSender
	logger *slog.Logger
}

// EmailSender interface for sending emails
type EmailSender interface {
	Send(to, subject, body string) error
}

// NewEmailDeliverer creates an email deliverer
func NewEmailDeliverer(sender EmailSender, logger *slog.Logger) *EmailDeliverer {
	return &EmailDeliverer{
		sender: sender,
		logger: logger,
	}
}

// Deliver sends the magic link via email
func (d *EmailDeliverer) Deliver(_ context.Context, token *Token, baseURL string) error {
	magicLink := fmt.Sprintf("%s/auth/magic-link/verify?token=%s", baseURL, token.Token)

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

	d.logger.Info("Magic link sent via email",
		"identifier", token.Identifier,
		"expires_in", time.Until(token.ExpiresAt).Round(time.Second),
	)

	return nil
}
