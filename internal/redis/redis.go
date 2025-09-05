// Package redis provides Redis infrastructure for caching, queuing, sessions, and pub/sub.
//
// The package wraps the go-redis client with domain-specific functionality including:
// - Cache: Key-value caching with TTL support
// - Queue: Task queue implementation
// - Session: Session management for authentication
// - PubSub: Publish/subscribe messaging
package redis

import (
	"errors"
	"time"
)

// Package errors
var (
	// ErrNoRedisConfig is returned when no Redis configuration is provided
	ErrNoRedisConfig = errors.New("no Redis configuration provided")

	// ErrNotInitialized is returned when Redis client is not initialized
	ErrNotInitialized = errors.New("redis client not initialized")

	// ErrKeyNotFound is returned when a key doesn't exist
	ErrKeyNotFound = errors.New("key not found")

	// ErrInvalidValue is returned when a value cannot be decoded
	ErrInvalidValue = errors.New("invalid value")
)

// Default configuration values
const (
	// DefaultMaxRetries is the default number of retries
	DefaultMaxRetries = 3

	// DefaultDialTimeout is the default dial timeout
	DefaultDialTimeout = 5 * time.Second

	// DefaultReadTimeout is the default read timeout
	DefaultReadTimeout = 3 * time.Second

	// DefaultWriteTimeout is the default write timeout
	DefaultWriteTimeout = 3 * time.Second

	// DefaultPoolSize is the default connection pool size
	DefaultPoolSize = 10

	// DefaultMinIdleConns is the default minimum idle connections
	DefaultMinIdleConns = 5

	// DefaultDB is the default Redis database
	DefaultDB = 0
)

// Key prefixes for different subsystems
const (
	// CacheKeyPrefix is the prefix for cache keys
	CacheKeyPrefix = "cache:"

	// SessionKeyPrefix is the prefix for session keys
	SessionKeyPrefix = "session:"

	// QueueKeyPrefix is the prefix for queue keys
	QueueKeyPrefix = "queue:"

	// PubSubKeyPrefix is the prefix for pubsub channel names
	PubSubKeyPrefix = "pubsub:"
)
