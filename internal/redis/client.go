// Package redis provides Redis client and utilities for caching, queuing, and session storage.
package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	ctx    = context.Background()
)

// Config represents Redis configuration
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

// Initialize creates and configures the Redis client
func Initialize(cfg Config) error {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	// Test connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return nil
}

// GetClient returns the Redis client instance
func GetClient() *redis.Client {
	return client
}

// Close closes the Redis connection
func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
