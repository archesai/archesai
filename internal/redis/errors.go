package redis

import "errors"

// Package errors.
var (
	// ErrNoRedisConfig is returned when no Redis configuration is provided.
	ErrNoRedisConfig = errors.New("no Redis configuration provided")

	// ErrNotInitialized is returned when Redis client is not initialized.
	ErrNotInitialized = errors.New("redis client not initialized")

	// ErrKeyNotFound is returned when a key doesn't exist.
	ErrKeyNotFound = errors.New("key not found")

	// ErrInvalidValue is returned when a value cannot be decoded.
	ErrInvalidValue = errors.New("invalid value")
)
