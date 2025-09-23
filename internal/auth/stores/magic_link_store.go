package stores

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"

	"github.com/archesai/archesai/internal/auth"
)

// MagicLinkStore handles magic link token operations.
type MagicLinkStore struct {
	repo      MagicLinkRepository
	ttl       time.Duration
	rateLimit int
}

// NewMagicLinkStore creates a new magic link store.
func NewMagicLinkStore(
	repo MagicLinkRepository,
	ttl time.Duration,
	rateLimit int,
) *MagicLinkStore {
	if ttl == 0 {
		ttl = 15 * time.Minute
	}
	if rateLimit == 0 {
		rateLimit = 5 // 5 requests per hour
	}
	return &MagicLinkStore{
		repo:      repo,
		ttl:       ttl,
		rateLimit: rateLimit,
	}
}

// CreateToken generates and stores a magic link token.
func (s *MagicLinkStore) CreateToken(
	ctx context.Context,
	identifier string,
	deliveryMethod auth.DeliveryMethod,
	userID *uuid.UUID,
	IPAddress string,
	userAgent string,
) (*auth.MagicLinkToken, error) {
	// Check rate limit
	since := time.Now().Add(-1 * time.Hour)
	count, err := s.repo.CountRecentByIdentifier(ctx, identifier, since)
	if err != nil {
		return nil, fmt.Errorf("checking rate limit: %w", err)
	}
	if count >= s.rateLimit {
		return nil, auth.ErrRateLimitExceeded
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
	if deliveryMethod == auth.DeliveryOTP {
		code = s.generateOTP()
	}

	// Create token record
	token := &MagicLinkToken{
		ID:             uuid.New(),
		UserID:         userID,
		Token:          tokenString,
		TokenHash:      tokenHash,
		Code:           code,
		Identifier:     identifier,
		DeliveryMethod: string(deliveryMethod),
		ExpiresAt:      time.Now().Add(s.ttl),
		IPAddress:      IPAddress,
		UserAgent:      userAgent,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Store in database
	if err := s.repo.Create(ctx, token); err != nil {
		return nil, fmt.Errorf("storing token: %w", err)
	}

	// Convert to auth module type (without exposing actual token)
	return &auth.MagicLinkToken{
		ID:             token.ID,
		UserID:         token.UserID,
		Token:          "", // Don't expose the token
		Code:           token.Code,
		Identifier:     token.Identifier,
		DeliveryMethod: deliveryMethod,
		ExpiresAt:      token.ExpiresAt,
		IPAddress:      token.IPAddress,
		UserAgent:      token.UserAgent,
		CreatedAt:      token.CreatedAt,
	}, nil
}

// VerifyToken verifies a magic link token.
func (s *MagicLinkStore) VerifyToken(
	ctx context.Context,
	token string,
) (*auth.MagicLinkToken, error) {
	// Hash the provided token
	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	// Look up the token
	t, err := s.repo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("looking up token: %w", err)
	}

	// Check if already used
	if t.UsedAt != nil {
		return nil, auth.ErrTokenAlreadyUsed
	}

	// Check if expired
	if time.Now().After(t.ExpiresAt) {
		return nil, auth.ErrTokenExpired
	}

	// Mark as used
	if err := s.repo.MarkUsed(ctx, t.ID); err != nil {
		return nil, fmt.Errorf("marking token as used: %w", err)
	}

	// Convert to auth module type
	return &auth.MagicLinkToken{
		ID:             t.ID,
		UserID:         t.UserID,
		Token:          t.Token,
		Code:           t.Code,
		Identifier:     t.Identifier,
		DeliveryMethod: auth.DeliveryMethod(t.DeliveryMethod),
		ExpiresAt:      t.ExpiresAt,
		UsedAt:         t.UsedAt,
		IPAddress:      t.IPAddress,
		UserAgent:      t.UserAgent,
		CreatedAt:      t.CreatedAt,
	}, nil
}

// VerifyOTP verifies an OTP code.
func (s *MagicLinkStore) VerifyOTP(
	ctx context.Context,
	identifier string,
	code string,
) (*auth.MagicLinkToken, error) {
	// Look up the token by identifier and code
	t, err := s.repo.GetByCode(ctx, identifier, code)
	if err != nil {
		return nil, auth.ErrInvalidOTP
	}

	// Check if already used
	if t.UsedAt != nil {
		return nil, auth.ErrTokenAlreadyUsed
	}

	// Check if expired
	if time.Now().After(t.ExpiresAt) {
		return nil, auth.ErrTokenExpired
	}

	// Mark as used
	if err := s.repo.MarkUsed(ctx, t.ID); err != nil {
		return nil, fmt.Errorf("marking OTP as used: %w", err)
	}

	// Convert to auth module type
	return &auth.MagicLinkToken{
		ID:             t.ID,
		UserID:         t.UserID,
		Token:          t.Token,
		Code:           t.Code,
		Identifier:     t.Identifier,
		DeliveryMethod: auth.DeliveryMethod(t.DeliveryMethod),
		ExpiresAt:      t.ExpiresAt,
		UsedAt:         t.UsedAt,
		IPAddress:      t.IPAddress,
		UserAgent:      t.UserAgent,
		CreatedAt:      t.CreatedAt,
	}, nil
}

// CleanupExpired removes expired tokens.
func (s *MagicLinkStore) CleanupExpired(ctx context.Context) error {
	return s.repo.DeleteExpired(ctx)
}

// generateOTP generates a 6-digit OTP code.
func (s *MagicLinkStore) generateOTP() string {
	maxNum := big.NewInt(999999)
	n, _ := rand.Int(rand.Reader, maxNum)
	return fmt.Sprintf("%06d", n.Int64())
}
