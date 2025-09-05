package redis

import (
	"encoding/json"
	"fmt"
	"time"
)

// Cache provides caching operations
type Cache struct{}

// NewCache creates a new cache instance
func NewCache() *Cache {
	return &Cache{}
}

// Set stores a value in cache with TTL
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return client.Set(ctx, key, data, ttl).Err()
}

// Get retrieves a value from cache
func (c *Cache) Get(key string, dest interface{}) error {
	data, err := client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// Delete removes a value from cache
func (c *Cache) Delete(key string) error {
	return client.Del(ctx, key).Err()
}

// Exists checks if a key exists
func (c *Cache) Exists(key string) (bool, error) {
	n, err := client.Exists(ctx, key).Result()
	return n > 0, err
}

// SetNX sets a value only if it doesn't exist
func (c *Cache) SetNX(key string, value interface{}, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	return client.SetNX(ctx, key, data, ttl).Result()
}

// Increment increments a counter
func (c *Cache) Increment(key string) (int64, error) {
	return client.Incr(ctx, key).Result()
}

// Decrement decrements a counter
func (c *Cache) Decrement(key string) (int64, error) {
	return client.Decr(ctx, key).Result()
}
