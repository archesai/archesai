// Package cache provides a generic caching interface and implementations
package cache

import (
	"context"
	"time"
)

// Cache provides a generic caching interface for any type T
type Cache[T any] interface {
	// Get retrieves an item from cache by key
	Get(ctx context.Context, key string) (*T, error)

	// Set stores an item in cache with TTL
	Set(ctx context.Context, key string, value *T, ttl time.Duration) error

	// Delete removes an item from cache
	Delete(ctx context.Context, key string) error

	// Exists checks if a key exists in cache
	Exists(ctx context.Context, key string) (bool, error)

	// Clear removes all items from cache (use with caution)
	Clear(ctx context.Context) error
}

// MultiCache provides batch operations for cache
type MultiCache[T any] interface {
	Cache[T]

	// GetMany retrieves multiple items from cache
	GetMany(ctx context.Context, keys []string) (map[string]*T, error)

	// SetMany stores multiple items in cache
	SetMany(ctx context.Context, items map[string]*T, ttl time.Duration) error

	// DeleteMany removes multiple items from cache
	DeleteMany(ctx context.Context, keys []string) error
}
