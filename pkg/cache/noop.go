package cache

import (
	"context"
	"time"
)

var _ Cache[any] = (*NoOpCache[any])(nil)

// NoOpCache implements Cache but does nothing (always returns cache miss).
type NoOpCache[T any] struct{}

// NewNoOpCache creates a new no-op cache.
func NewNoOpCache[T any]() *NoOpCache[T] {
	return &NoOpCache[T]{}
}

// Get always returns cache miss.
func (c *NoOpCache[T]) Get(_ context.Context, _ string) (*T, error) {
	return nil, ErrCacheMiss
}

// GetMany always returns empty map.
func (c *NoOpCache[T]) GetMany(_ context.Context, _ []string) (map[string]*T, error) {
	return make(map[string]*T), nil
}

// Set does nothing.
func (c *NoOpCache[T]) Set(_ context.Context, _ string, _ *T, _ time.Duration) error {
	return nil
}

// SetMany does nothing.
func (c *NoOpCache[T]) SetMany(_ context.Context, _ map[string]*T, _ time.Duration) error {
	return nil
}

// Delete does nothing.
func (c *NoOpCache[T]) Delete(_ context.Context, _ string) error {
	return nil
}

// DeleteMany does nothing.
func (c *NoOpCache[T]) DeleteMany(_ context.Context, _ []string) error {
	return nil
}

// Exists always returns false.
func (c *NoOpCache[T]) Exists(_ context.Context, _ string) (bool, error) {
	return false, nil
}

// Clear does nothing.
func (c *NoOpCache[T]) Clear(_ context.Context) error {
	return nil
}
