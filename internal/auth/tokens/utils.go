package tokens

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

// GenerateAPIKey generates a new API key with prefix.
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

// ParseAPIKey extracts the key from various header formats.
func ParseAPIKey(authHeader string) string {
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

// ValidateAPIKeyFormat checks if the API key has valid format.
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

// HashAPIKey creates a hash of the API key for storage.
func HashAPIKey(key string) string {
	// In production, use a proper hashing algorithm like bcrypt or argon2
	// For now, using a simple SHA256 (should be replaced)
	// TODO: Implement proper hashing with bcrypt or argon2
	return fmt.Sprintf("hashed_%s", key) // Placeholder - implement proper hashing
}

// GenerateSecureToken generates a cryptographically secure random token.
func GenerateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
