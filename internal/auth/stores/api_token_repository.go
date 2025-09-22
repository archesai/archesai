// Package stores provides storage implementations for authentication data
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
	"github.com/archesai/archesai/internal/database/postgresql"
)

// APITokenRepository implements auth.APITokenStore using PostgreSQL
type APITokenRepository struct {
	db *postgresql.Queries
}

// NewAPITokenRepository creates a new API token repository
func NewAPITokenRepository(db *postgresql.Queries) *APITokenRepository {
	return &APITokenRepository{
		db: db,
	}
}

// CreateToken creates a new API token
const tokenPrefix = "sk_live_"

// CreateToken creates a new API token
func (r *APITokenRepository) CreateToken(
	ctx context.Context,
	userID, organizationID uuid.UUID,
	name string,
	scopes []string,
	expiresIn time.Duration,
) (*auth.APITokenResponse, error) {
	// Generate the raw token
	rawToken, err := generateSecureToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Create the formatted token with prefix
	fullToken := tokenPrefix + rawToken
	prefix := fullToken[:8] // First 8 chars for identification

	// Hash the token for storage
	hash := hashToken(fullToken)

	// Calculate expiration time
	var expiresAt *time.Time
	if expiresIn > 0 {
		exp := time.Now().Add(expiresIn)
		expiresAt = &exp
	}

	// Create the token in database
	tokenID := uuid.New()
	params := postgresql.CreateApiTokenParams{
		ID:             tokenID,
		UserID:         userID,
		OrganizationID: organizationID,
		Name:           &name,
		KeyHash:        hash,
		Prefix:         &prefix,
		Scopes:         scopes,
		RateLimit:      60, // Default rate limit
		ExpiresAt:      expiresAt,
	}

	dbToken, err := r.db.CreateApiToken(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	// Return response
	var expiresAtTime time.Time
	if dbToken.ExpiresAt != nil {
		expiresAtTime = *dbToken.ExpiresAt
	}

	return &auth.APITokenResponse{
		ID:        dbToken.ID,
		Name:      name,
		Key:       fullToken, // Use Key field, not Token
		Prefix:    prefix,
		Scopes:    dbToken.Scopes,
		RateLimit: int(dbToken.RateLimit),
		ExpiresAt: expiresAtTime,
		CreatedAt: dbToken.CreatedAt,
	}, nil
}

// ValidateToken validates an API token and returns its details
func (r *APITokenRepository) ValidateToken(
	ctx context.Context,
	tokenString string,
) (*auth.APIToken, error) {
	// Hash the provided token
	hash := hashToken(tokenString)

	// Look up the token by hash
	dbToken, err := r.db.GetApiTokenByKeyHash(ctx, hash)
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	// Check expiration
	if dbToken.ExpiresAt != nil && time.Now().After(*dbToken.ExpiresAt) {
		return nil, fmt.Errorf("token expired")
	}

	// Update last used timestamp
	_ = r.db.UpdateApiTokenLastUsed(ctx, dbToken.ID)

	// Convert to auth.APIToken
	name := ""
	if dbToken.Name != nil {
		name = *dbToken.Name
	}

	prefix := tokenPrefix
	if dbToken.Prefix != nil {
		prefix = *dbToken.Prefix
	}

	return &auth.APIToken{
		ID:             dbToken.ID,
		UserID:         dbToken.UserID,
		OrganizationID: dbToken.OrganizationID,
		Name:           name,
		Prefix:         prefix,
		Scopes:         dbToken.Scopes,
		ExpiresAt:      dbToken.ExpiresAt,
		CreatedAt:      dbToken.CreatedAt,
	}, nil
}

// GetToken retrieves a token by ID
func (r *APITokenRepository) GetToken(
	ctx context.Context,
	tokenID uuid.UUID,
) (*auth.APIToken, error) {
	dbToken, err := r.db.GetApiToken(ctx, tokenID)
	if err != nil {
		return nil, fmt.Errorf("token not found")
	}

	name := ""
	if dbToken.Name != nil {
		name = *dbToken.Name
	}

	prefix := tokenPrefix
	if dbToken.Prefix != nil {
		prefix = *dbToken.Prefix
	}

	return &auth.APIToken{
		ID:             dbToken.ID,
		UserID:         dbToken.UserID,
		OrganizationID: dbToken.OrganizationID,
		Name:           name,
		Prefix:         prefix,
		Scopes:         dbToken.Scopes,
		ExpiresAt:      dbToken.ExpiresAt,
		CreatedAt:      dbToken.CreatedAt,
	}, nil
}

// ListTokensByUser lists all tokens for a user
func (r *APITokenRepository) ListTokensByUser(
	ctx context.Context,
	userID uuid.UUID,
) ([]*auth.APIToken, error) {
	// Default pagination
	dbTokens, err := r.db.ListApiTokensByUser(ctx, postgresql.ListApiTokensByUserParams{
		UserID: userID,
		Limit:  100,
		Offset: 0,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list tokens: %w", err)
	}

	tokens := make([]*auth.APIToken, len(dbTokens))
	for i, dbToken := range dbTokens {
		name := ""
		if dbToken.Name != nil {
			name = *dbToken.Name
		}

		prefix := tokenPrefix
		if dbToken.Prefix != nil {
			prefix = *dbToken.Prefix
		}

		tokens[i] = &auth.APIToken{
			ID:             dbToken.ID,
			UserID:         dbToken.UserID,
			OrganizationID: dbToken.OrganizationID,
			Name:           name,
			Prefix:         prefix,
			Scopes:         dbToken.Scopes,
			ExpiresAt:      dbToken.ExpiresAt,
			CreatedAt:      dbToken.CreatedAt,
		}
	}

	return tokens, nil
}

// RevokeToken revokes an API token
func (r *APITokenRepository) RevokeToken(
	ctx context.Context,
	tokenID uuid.UUID,
) error {
	err := r.db.DeleteApiToken(ctx, tokenID)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}
	return nil
}

// RevokeUserTokens revokes all tokens for a user
func (r *APITokenRepository) RevokeUserTokens(
	ctx context.Context,
	userID uuid.UUID,
) error {
	err := r.db.DeleteApiTokensByUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke user tokens: %w", err)
	}
	return nil
}

// CleanupExpired removes expired tokens
func (r *APITokenRepository) CleanupExpired(ctx context.Context) error {
	return r.db.DeleteExpiredApiTokens(ctx)
}

// generateSecureToken generates a cryptographically secure random token
func generateSecureToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// hashToken creates a SHA256 hash of the token
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// ParseAPIKey extracts the key from various header formats
func (r *APITokenRepository) ParseAPIKey(authHeader string) string {
	authHeader = strings.TrimSpace(authHeader)

	// Check for Bearer token format
	if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		token := strings.TrimSpace(authHeader[7:])
		if strings.HasPrefix(token, "sk_") {
			return token
		}
	}

	// Check for direct API key
	if strings.HasPrefix(authHeader, "sk_") {
		return authHeader
	}

	return ""
}
