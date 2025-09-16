// Package redis provides Redis infrastructure for caching, queuing, sessions, and pub/sub.
//
// The package wraps the go-redis client with domain-specific functionality including:
// - Cache: Key-value caching with TTL support
// - Queue: Task queue implementation
// - Session: Session management for authentication
// - PubSub: Publish/subscribe messaging
package redis

import (
	"time"
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
	// QueueKeyPrefix is the prefix for queue keys
	QueueKeyPrefix = "queue:"

	// PubSubKeyPrefix is the prefix for pubsub channel names
	PubSubKeyPrefix = "pubsub:"
)
