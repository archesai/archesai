// Package redis provides Redis infrastructure for caching, queuing, sessions, and pub/sub
package redis

import (
	"time"
)

// Config holds Redis-specific configuration
type Config struct {
	// Connection settings
	URL      string `mapstructure:"url" yaml:"url" env:"REDIS_URL"`
	Host     string `mapstructure:"host" yaml:"host" env:"REDIS_HOST"`
	Port     int    `mapstructure:"port" yaml:"port" env:"REDIS_PORT"`
	Password string `mapstructure:"password" yaml:"password" env:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"db" yaml:"db" env:"REDIS_DB"`

	// Connection pool settings
	PoolSize     int           `mapstructure:"pool_size" yaml:"pool_size" env:"REDIS_POOL_SIZE"`
	MinIdleConns int           `mapstructure:"min_idle_conns" yaml:"min_idle_conns" env:"REDIS_MIN_IDLE_CONNS"`
	MaxRetries   int           `mapstructure:"max_retries" yaml:"max_retries" env:"REDIS_MAX_RETRIES"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout" yaml:"dial_timeout" env:"REDIS_DIAL_TIMEOUT"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" yaml:"read_timeout" env:"REDIS_READ_TIMEOUT"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" yaml:"write_timeout" env:"REDIS_WRITE_TIMEOUT"`

	// Feature flags
	EnableCache   bool `mapstructure:"enable_cache" yaml:"enable_cache" env:"REDIS_ENABLE_CACHE"`
	EnableQueue   bool `mapstructure:"enable_queue" yaml:"enable_queue" env:"REDIS_ENABLE_QUEUE"`
	EnableSession bool `mapstructure:"enable_session" yaml:"enable_session" env:"REDIS_ENABLE_SESSION"`
	EnablePubSub  bool `mapstructure:"enable_pubsub" yaml:"enable_pubsub" env:"REDIS_ENABLE_PUBSUB"`

	// Cache settings
	DefaultCacheTTL time.Duration `mapstructure:"default_cache_ttl" yaml:"default_cache_ttl" env:"REDIS_DEFAULT_CACHE_TTL"`

	// Session settings
	SessionTTL time.Duration `mapstructure:"session_ttl" yaml:"session_ttl" env:"REDIS_SESSION_TTL"`

	// Queue settings
	QueueBlockTimeout time.Duration `mapstructure:"queue_block_timeout" yaml:"queue_block_timeout" env:"REDIS_QUEUE_BLOCK_TIMEOUT"`
}

// DefaultConfig returns a Config with default values
func DefaultConfig() *Config {
	return &Config{
		Host:         "localhost",
		Port:         6379,
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 2,
		MaxRetries:   3,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,

		EnableCache:   true,
		EnableQueue:   true,
		EnableSession: true,
		EnablePubSub:  true,

		DefaultCacheTTL:   1 * time.Hour,
		SessionTTL:        24 * time.Hour,
		QueueBlockTimeout: 30 * time.Second,
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.URL == "" && c.Host == "" {
		return ErrNoRedisConfig
	}

	if c.PoolSize < 1 {
		c.PoolSize = 10
	}

	if c.MinIdleConns < 0 {
		c.MinIdleConns = 0
	}

	if c.MaxRetries < 0 {
		c.MaxRetries = 3
	}

	return nil
}
