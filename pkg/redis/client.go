package redis

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

// Client wraps the Redis client with domain-specific functionality.
type Client struct {
	redis  *redis.Client
	Queue  *Queue
	PubSub *PubSub
	config *Config
}

// NewClient creates a new Redis client with all features.
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid Redis config: %w", err)
	}

	// Create Redis options
	opts := &redis.Options{
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		Password:     config.Password,
	}

	// Use URL if provided, otherwise use host:port
	if config.URL != "" {
		parsedOpts, err := redis.ParseURL(config.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
		}
		opts = parsedOpts
		// Override with explicitly set options
		if config.PoolSize > 0 {
			opts.PoolSize = config.PoolSize
		}
		if config.MinIdleConns > 0 {
			opts.MinIdleConns = config.MinIdleConns
		}
	} else {
		opts.Addr = fmt.Sprintf("%s:%d", config.Host, config.Port)
	}

	// Create Redis client
	redisClient := redis.NewClient(opts)

	// Test connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	slog.Info("Connected to Redis",
		"addr", opts.Addr,
		"db", opts.DB,
		"pool_size", opts.PoolSize,
	)

	// Create the client with all features
	client := &Client{
		redis:  redisClient,
		config: config,
	}

	// Initialize features based on config
	if config.EnableQueue {
		client.Queue = NewQueue(redisClient)
		slog.Info("Redis queue enabled")
	}

	if config.EnablePubSub {
		client.PubSub = NewPubSub(redisClient)
		slog.Info("Redis pub/sub enabled")
	}

	return client, nil
}

// Close closes the Redis connection.
func (c *Client) Close() error {
	if c.redis != nil {
		return c.redis.Close()
	}
	return nil
}

// Ping checks if Redis is reachable.
func (c *Client) Ping(ctx context.Context) error {
	return c.redis.Ping(ctx).Err()
}

// GetRedisClient returns the underlying Redis client for advanced usage.
func (c *Client) GetRedisClient() *redis.Client {
	return c.redis
}

// FlushDB flushes the current database (use with caution!)
func (c *Client) FlushDB(ctx context.Context) error {
	return c.redis.FlushDB(ctx).Err()
}

// FlushAll flushes all databases (use with extreme caution!)
func (c *Client) FlushAll(ctx context.Context) error {
	return c.redis.FlushAll(ctx).Err()
}

// Info returns Redis server information.
func (c *Client) Info(ctx context.Context, sections ...string) (string, error) {
	var cmd *redis.StringCmd
	if len(sections) > 0 {
		cmd = c.redis.Info(ctx, sections...)
	} else {
		cmd = c.redis.Info(ctx)
	}
	return cmd.Result()
}

// DBSize returns the number of keys in the current database.
func (c *Client) DBSize(ctx context.Context) (int64, error) {
	return c.redis.DBSize(ctx).Result()
}
