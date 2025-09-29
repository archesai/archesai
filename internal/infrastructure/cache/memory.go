package cache

import (
	"context"
	"sync"
	"time"
)

// MemoryCache implements Cache using in-memory storage (for testing).
type MemoryCache[T any] struct {
	mu    sync.RWMutex
	items map[string]*memoryItem[T]
}

type memoryItem[T any] struct {
	value     *T
	expiresAt time.Time
}

// NewMemoryCache creates a new in-memory cache.
func NewMemoryCache[T any]() *MemoryCache[T] {
	return &MemoryCache[T]{
		items: make(map[string]*memoryItem[T]),
	}
}

// Get retrieves an item from cache by key.
func (c *MemoryCache[T]) Get(_ context.Context, key string) (*T, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return nil, nil
	}

	// Check expiration
	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		return nil, nil
	}

	return item.value, nil
}

// Set stores an item in cache with TTL.
func (c *MemoryCache[T]) Set(_ context.Context, key string, value *T, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	c.items[key] = &memoryItem[T]{
		value:     value,
		expiresAt: expiresAt,
	}

	return nil
}

// Delete removes an item from cache.
func (c *MemoryCache[T]) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
	return nil
}

// Exists checks if a key exists in cache.
func (c *MemoryCache[T]) Exists(_ context.Context, key string) (bool, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return false, nil
	}

	// Check expiration
	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		return false, nil
	}

	return true, nil
}

// Clear removes all items from cache.
func (c *MemoryCache[T]) Clear(_ context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*memoryItem[T])
	return nil
}

// GetMany retrieves multiple items from cache.
func (c *MemoryCache[T]) GetMany(_ context.Context, keys []string) (map[string]*T, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]*T)
	now := time.Now()

	for _, key := range keys {
		item, exists := c.items[key]
		if !exists {
			continue
		}

		// Check expiration
		if !item.expiresAt.IsZero() && now.After(item.expiresAt) {
			continue
		}

		result[key] = item.value
	}

	return result, nil
}

// SetMany stores multiple items in cache.
func (c *MemoryCache[T]) SetMany(_ context.Context, items map[string]*T, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiresAt time.Time
	if ttl > 0 {
		expiresAt = time.Now().Add(ttl)
	}

	for key, value := range items {
		c.items[key] = &memoryItem[T]{
			value:     value,
			expiresAt: expiresAt,
		}
	}

	return nil
}

// DeleteMany removes multiple items from cache.
func (c *MemoryCache[T]) DeleteMany(_ context.Context, keys []string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, key := range keys {
		delete(c.items, key)
	}

	return nil
}

// Ensure MemoryCache implements MultiCache.
var _ MultiCache[any] = (*MemoryCache[any])(nil)
