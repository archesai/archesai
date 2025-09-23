package stores

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/archesai/archesai/internal/auth"
)

// APIKeyStore handles API token operations.
type APIKeyStore struct {
	repo  auth.APIKeyRepository
	cache auth.APIKeyCache
}

// NewAPIKeyStore creates a new API token store.
func NewAPIKeyStore(repo auth.APIKeyRepository, cache auth.APIKeyCache) *APIKeyStore {
	return &APIKeyStore{
		repo:  repo,
		cache: cache,
	}
}

// CreateToken creates a new API token.
func (s *APIKeyStore) CreateToken(
	ctx context.Context,
	userID, organizationID uuid.UUID,
	name string,
	scopes []string,
	expiresIn time.Duration,
) (*auth.APIKey, error) {
	// Generate the key
	plainKey, prefix, err := s.generateAPIKey()
	if err != nil {
		return nil, fmt.Errorf("generate api key: %w", err)
	}

	hash := s.hashAPIKey(plainKey)

	// Create the API key entity
	apiKey := &auth.APIKey{
		ID:             uuid.New(),
		UserID:         &userID,
		OrganizationID: &organizationID,
		Name:           name,
		KeyHash:        &hash,
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

	return &auth.APIKey{
		ID:        created.ID,
		Name:      created.Name,
		Prefix:    created.Prefix,
		Scopes:    created.Scopes,
		RateLimit: created.RateLimit,
		ExpiresAt: created.ExpiresAt,
		CreatedAt: created.CreatedAt,
	}, nil
}

// ValidateToken validates an API token and returns the associated data.
func (s *APIKeyStore) ValidateToken(ctx context.Context, key string) (*auth.APIKey, error) {
	// Validate format
	if !s.validateAPIKeyFormat(key) {
		return nil, auth.ErrInvalidAPIKeyFormat
	}

	// Extract prefix
	prefix := key[:8]

	// Check cache first
	if s.cache != nil {
		cached, err := s.cache.GetAPIKey(ctx, prefix)
		if err == nil && cached != nil {
			// Validate the full key hash
			if *cached.KeyHash == s.hashAPIKey(key) {
				// Check expiration
				if time.Now().Before(cached.ExpiresAt) {
					// Update last used timestamp asynchronously
					go func() {
						_ = s.repo.UpdateAPIKeyLastUsed(context.Background(), cached.ID)
					}()
					return s.convertToAuthAPIKey(cached), nil
				}
			}
		}
	}

	// Validate against database
	apiKey, err := s.repo.ValidateAPIKeyHash(ctx, prefix, s.hashAPIKey(key))
	if err != nil {
		return nil, auth.ErrInvalidAPIKey
	}

	// Check expiration
	if time.Now().After(apiKey.ExpiresAt) {
		return nil, auth.ErrAPIKeyExpired
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

	return s.convertToAuthAPIKey(apiKey), nil
}

// RevokeToken revokes an API token.
func (s *APIKeyStore) RevokeToken(ctx context.Context, keyID uuid.UUID) error {
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

// ListTokensByUser returns all API tokens for a user.
func (s *APIKeyStore) ListTokensByUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]*auth.APIKey, error) {
	apiKeys, err := s.repo.ListUserAPIKeys(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*auth.APIKey, len(apiKeys))
	for i, key := range apiKeys {
		result[i] = s.convertToAuthAPIKey(key)
	}

	return result, nil
}

// ParseAPIKey extracts the key from various header formats.
func (s *APIKeyStore) ParseAPIKey(authHeader string) string {
	authHeader = strings.TrimSpace(authHeader)

	// Check for "APIKey" scheme
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

// Helper methods

func (s *APIKeyStore) generateAPIKey() (string, string, error) {
	// Generate 32 bytes of random data
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", fmt.Errorf("generate random bytes: %w", err)
	}

	// Convert to hex string
	key := hex.EncodeToString(bytes)

	// Format as: sk_live_<64 hex chars> or sk_test_<64 hex chars>
	const keyPrefix = "sk_live_"
	fullKey := fmt.Sprintf("%s%s", keyPrefix, key)
	prefix := fullKey[:8] // First 8 chars for identification

	return fullKey, prefix, nil
}

func (s *APIKeyStore) hashAPIKey(key string) string {
	// Use SHA256 for consistent hashing
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

func (s *APIKeyStore) validateAPIKeyFormat(key string) bool {
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

func (s *APIKeyStore) convertToAuthAPIKey(key *auth.APIKey) *auth.APIKey {
	if key == nil {
		return nil
	}
	return &auth.APIKey{
		ID:             key.ID,
		UserID:         key.UserID,
		OrganizationID: key.OrganizationID,
		Name:           key.Name,
		Prefix:         key.Prefix,
		Scopes:         key.Scopes,
		RateLimit:      key.RateLimit,
		ExpiresAt:      key.ExpiresAt,
		LastUsedAt:     key.LastUsedAt,
		CreatedAt:      key.CreatedAt,
		UpdatedAt:      key.UpdatedAt,
	}
}
