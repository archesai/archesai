package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// APIKey represents an API key entity
type APIKey struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	KeyHash        string    `json:"key_hash"` // Store hashed version
	Prefix         string    `json:"prefix"`   // First 8 chars for identification
	Scopes         []string  `json:"scopes"`
	RateLimit      int       `json:"rate_limit"` // Requests per minute
	LastUsedAt     time.Time `json:"last_used_at"`
	ExpiresAt      time.Time `json:"expires_at"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// APIKeyCreate represents the response when creating a new API key
type APIKeyCreate struct {
	*APIKey
	PlainKey string `json:"key"` // Only returned once on creation
}

// GenerateAPIKey generates a new API key with prefix
func GenerateAPIKey() (string, string, error) {
	// Generate 32 bytes of random data
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", fmt.Errorf("generate random bytes: %w", err)
	}

	// Convert to hex string
	key := hex.EncodeToString(bytes)

	// Format as: sk_live_<64 hex chars> or sk_test_<64 hex chars>
	fullKey := fmt.Sprintf("sk_live_%s", key)
	prefix := fullKey[:8] // First 8 chars for identification

	return fullKey, prefix, nil
}

// ParseAPIKey extracts the key from various header formats
func ParseAPIKey(authHeader string) string {
	authHeader = strings.TrimSpace(authHeader)

	// Check for "ApiKey" scheme
	if strings.HasPrefix(strings.ToLower(authHeader), "apikey ") {
		return strings.TrimSpace(authHeader[7:])
	}

	// Check for "Bearer" scheme with API key format
	if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		key := strings.TrimSpace(authHeader[7:])
		if strings.HasPrefix(key, "sk_") {
			return key
		}
	}

	// Direct API key (from X-API-Key header)
	if strings.HasPrefix(authHeader, "sk_") {
		return authHeader
	}

	return ""
}

// ValidateAPIKeyFormat checks if the API key has valid format
func ValidateAPIKeyFormat(key string) bool {
	// Expected format: sk_live_<64 hex chars> or sk_test_<64 hex chars>
	if !strings.HasPrefix(key, "sk_") {
		return false
	}

	parts := strings.SplitN(key, "_", 3)
	if len(parts) != 3 {
		return false
	}

	// Check environment (live or test)
	if parts[1] != "live" && parts[1] != "test" {
		return false
	}

	// Check hex string length (32 bytes = 64 hex chars)
	if len(parts[2]) != 64 {
		return false
	}

	// Verify it's valid hex
	_, err := hex.DecodeString(parts[2])
	return err == nil
}

// HashAPIKey creates a hash of the API key for storage
func HashAPIKey(key string) string {
	// In production, use a proper hashing algorithm like bcrypt or argon2
	// For now, using a simple SHA256 (should be replaced)
	return fmt.Sprintf("hashed_%s", key) // Placeholder - implement proper hashing
}

// APIKeyRepository defines the interface for API key storage operations
type APIKeyRepository interface {
	CreateAPIKey(ctx context.Context, apiKey *APIKey) (*APIKey, error)
	GetAPIKeyByPrefix(ctx context.Context, prefix string) (*APIKey, error)
	GetAPIKeyByID(ctx context.Context, id uuid.UUID) (*APIKey, error)
	ListUserAPIKeys(ctx context.Context, userID uuid.UUID) ([]*APIKey, error)
	UpdateAPIKeyLastUsed(ctx context.Context, id uuid.UUID) error
	DeleteAPIKey(ctx context.Context, id uuid.UUID) error
	ValidateAPIKeyHash(ctx context.Context, prefix, keyHash string) (*APIKey, error)
}

// APIKeyCache provides caching for API keys
type APIKeyCache interface {
	GetAPIKey(ctx context.Context, prefix string) (*APIKey, error)
	SetAPIKey(ctx context.Context, prefix string, apiKey *APIKey, ttl time.Duration) error
	DeleteAPIKey(ctx context.Context, prefix string) error
}

// APIKeyService handles API key operations
type APIKeyService struct {
	repo  APIKeyRepository
	cache APIKeyCache
}

// NewAPIKeyService creates a new API key service
func NewAPIKeyService(repo APIKeyRepository, cache APIKeyCache) *APIKeyService {
	return &APIKeyService{
		repo:  repo,
		cache: cache,
	}
}

// CreateAPIKey creates a new API key for a user
func (s *APIKeyService) CreateAPIKey(ctx context.Context, userID, orgID uuid.UUID, name string, scopes []string, expiresIn time.Duration) (*APIKeyCreate, error) {
	// Generate the key
	plainKey, prefix, err := GenerateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("generate api key: %w", err)
	}

	// Create the API key entity
	apiKey := &APIKey{
		ID:             uuid.New(),
		UserID:         userID,
		OrganizationID: orgID,
		Name:           name,
		KeyHash:        HashAPIKey(plainKey),
		Prefix:         prefix,
		Scopes:         scopes,
		RateLimit:      100, // Default rate limit
		ExpiresAt:      time.Now().Add(expiresIn),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Store in database
	created, err := s.repo.CreateAPIKey(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("store api key: %w", err)
	}

	// Cache the key data
	if s.cache != nil {
		_ = s.cache.SetAPIKey(ctx, prefix, created, expiresIn)
	}

	return &APIKeyCreate{
		APIKey:   created,
		PlainKey: plainKey,
	}, nil
}

// ValidateAPIKey validates an API key and returns the associated data
func (s *APIKeyService) ValidateAPIKey(ctx context.Context, key string) (*APIKey, error) {
	// Validate format
	if !ValidateAPIKeyFormat(key) {
		return nil, fmt.Errorf("invalid api key format")
	}

	// Extract prefix
	prefix := key[:8]

	// Check cache first
	if s.cache != nil {
		cached, err := s.cache.GetAPIKey(ctx, prefix)
		if err == nil && cached != nil {
			// Validate the full key hash
			if cached.KeyHash == HashAPIKey(key) {
				// Check expiration
				if time.Now().Before(cached.ExpiresAt) {
					// Update last used timestamp asynchronously
					go func() {
						_ = s.repo.UpdateAPIKeyLastUsed(context.Background(), cached.ID)
					}()
					return cached, nil
				}
			}
		}
	}

	// Validate against database
	apiKey, err := s.repo.ValidateAPIKeyHash(ctx, prefix, HashAPIKey(key))
	if err != nil {
		return nil, fmt.Errorf("validate api key: %w", err)
	}

	// Check expiration
	if time.Now().After(apiKey.ExpiresAt) {
		return nil, fmt.Errorf("api key expired")
	}

	// Update last used timestamp asynchronously
	go func() {
		_ = s.repo.UpdateAPIKeyLastUsed(context.Background(), apiKey.ID)
	}()

	// Cache for future requests
	if s.cache != nil {
		ttl := time.Until(apiKey.ExpiresAt)
		_ = s.cache.SetAPIKey(ctx, apiKey.Prefix, apiKey, ttl)
	}

	return apiKey, nil
}

// RevokeAPIKey revokes an API key
func (s *APIKeyService) RevokeAPIKey(ctx context.Context, keyID uuid.UUID) error {
	// Get the key to find its prefix
	apiKey, err := s.repo.GetAPIKeyByID(ctx, keyID)
	if err != nil {
		return fmt.Errorf("get api key: %w", err)
	}

	// Delete from database
	if err := s.repo.DeleteAPIKey(ctx, keyID); err != nil {
		return fmt.Errorf("delete api key: %w", err)
	}

	// Remove from cache
	if s.cache != nil {
		_ = s.cache.DeleteAPIKey(ctx, apiKey.Prefix)
	}

	return nil
}
