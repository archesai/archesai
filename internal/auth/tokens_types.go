package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

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
