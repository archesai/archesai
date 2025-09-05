package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache provides caching operations
type Cache struct {
	client *redis.Client
}

// NewCache creates a new cache instance
func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client: client,
	}
}

// Set stores a value in cache with TTL
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) error {
	ctx := context.Background()
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.client.Set(ctx, key, data, ttl).Err()
}

// Get retrieves a value from cache
func (c *Cache) Get(key string, dest interface{}) error {
	ctx := context.Background()
	data, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// Delete removes a value from cache
func (c *Cache) Delete(key string) error {
	ctx := context.Background()
	return c.client.Del(ctx, key).Err()
}

// Exists checks if a key exists
func (c *Cache) Exists(key string) (bool, error) {
	ctx := context.Background()
	n, err := c.client.Exists(ctx, key).Result()
	return n > 0, err
}

// SetNX sets a value only if it doesn't exist
func (c *Cache) SetNX(key string, value interface{}, ttl time.Duration) (bool, error) {
	ctx := context.Background()
	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.client.SetNX(ctx, key, data, ttl).Result()
}

// Increment increments a counter
func (c *Cache) Increment(key string) (int64, error) {
	ctx := context.Background()
	return c.client.Incr(ctx, key).Result()
}

// Decrement decrements a counter
func (c *Cache) Decrement(key string) (int64, error) {
	ctx := context.Background()
	return c.client.Decr(ctx, key).Result()
}

// TTL gets the remaining TTL of a key
func (c *Cache) TTL(key string) (time.Duration, error) {
	ctx := context.Background()
	return c.client.TTL(ctx, key).Result()
}

// Expire sets/updates the TTL of a key
func (c *Cache) Expire(key string, ttl time.Duration) error {
	ctx := context.Background()
	return c.client.Expire(ctx, key, ttl).Err()
}
