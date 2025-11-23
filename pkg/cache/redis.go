package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var _ Cache[any] = (*RedisCache[any])(nil)

// RedisCache implements Cache using Redis as the backend.
type RedisCache[T any] struct {
	client    *redis.Client
	keyPrefix string
}

// NewRedisCache creates a new Redis-backed cache with an optional key prefix.
func NewRedisCache[T any](client *redis.Client, keyPrefix string) *RedisCache[T] {
	return &RedisCache[T]{
		client:    client,
		keyPrefix: keyPrefix,
	}
}

// formatKey adds the prefix to a cache key.
func (c *RedisCache[T]) formatKey(key string) string {
	if c.keyPrefix != "" {
		return fmt.Sprintf("%s:%s", c.keyPrefix, key)
	}
	return key
}

// Get retrieves an item from cache by key.
func (c *RedisCache[T]) Get(ctx context.Context, key string) (*T, error) {
	data, err := c.client.Get(ctx, c.formatKey(key)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil // Key doesn't exist
		}
		return nil, fmt.Errorf("redis get: %w", err)
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return &value, nil
}

// Set stores an item in cache with TTL.
func (c *RedisCache[T]) Set(ctx context.Context, key string, value *T, ttl time.Duration) error {
	if value == nil {
		return errors.New("cannot cache nil value")
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	if err := c.client.Set(ctx, c.formatKey(key), data, ttl).Err(); err != nil {
		return fmt.Errorf("redis set: %w", err)
	}

	return nil
}

// Delete removes an item from cache.
func (c *RedisCache[T]) Delete(ctx context.Context, key string) error {
	if err := c.client.Del(ctx, c.formatKey(key)).Err(); err != nil {
		return fmt.Errorf("redis del: %w", err)
	}
	return nil
}

// Exists checks if a key exists in cache.
func (c *RedisCache[T]) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, c.formatKey(key)).Result()
	if err != nil {
		return false, fmt.Errorf("redis exists: %w", err)
	}
	return count > 0, nil
}

// Clear removes all items from cache (use with caution).
func (c *RedisCache[T]) Clear(ctx context.Context) error {
	if c.keyPrefix == "" {
		return errors.New("cannot clear cache without key prefix")
	}

	pattern := fmt.Sprintf("%s:*", c.keyPrefix)
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()

	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("redis scan: %w", err)
	}

	if len(keys) > 0 {
		if err := c.client.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("redis del: %w", err)
		}
	}

	return nil
}

// GetMany retrieves multiple items from cache.
func (c *RedisCache[T]) GetMany(ctx context.Context, keys []string) (map[string]*T, error) {
	if len(keys) == 0 {
		return make(map[string]*T), nil
	}

	// Format keys with prefix
	formattedKeys := make([]string, len(keys))
	for i, key := range keys {
		formattedKeys[i] = c.formatKey(key)
	}

	// Get values from Redis
	values, err := c.client.MGet(ctx, formattedKeys...).Result()
	if err != nil {
		return nil, fmt.Errorf("redis mget: %w", err)
	}

	// Parse results
	result := make(map[string]*T)
	for i, val := range values {
		if val == nil {
			continue // Key doesn't exist
		}

		data, ok := val.(string)
		if !ok {
			continue
		}

		var value T
		if err := json.Unmarshal([]byte(data), &value); err != nil {
			continue // Skip invalid entries
		}

		result[keys[i]] = &value
	}

	return result, nil
}

// SetMany stores multiple items in cache.
func (c *RedisCache[T]) SetMany(ctx context.Context, items map[string]*T, ttl time.Duration) error {
	if len(items) == 0 {
		return nil
	}

	pipe := c.client.Pipeline()

	for key, value := range items {
		if value == nil {
			continue
		}

		data, err := json.Marshal(value)
		if err != nil {
			continue // Skip items that can't be marshaled
		}

		pipe.Set(ctx, c.formatKey(key), data, ttl)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("redis pipeline exec: %w", err)
	}

	return nil
}

// DeleteMany removes multiple items from cache.
func (c *RedisCache[T]) DeleteMany(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	formattedKeys := make([]string, len(keys))
	for i, key := range keys {
		formattedKeys[i] = c.formatKey(key)
	}

	if err := c.client.Del(ctx, formattedKeys...).Err(); err != nil {
		return fmt.Errorf("redis del: %w", err)
	}

	return nil
}
