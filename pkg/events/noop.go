package events

import (
	"context"
)

// NoOpPublisher is a no-op event publisher implementation for testing.
type NoOpPublisher struct{}

// NewNoOpPublisher creates a new no-op event publisher.
func NewNoOpPublisher() Publisher {
	return &NoOpPublisher{}
}

// Publish does nothing in no-op implementation.
func (p *NoOpPublisher) Publish(_ context.Context, _ Event) error {
	return nil
}

// PublishMultiple does nothing in no-op implementation.
func (p *NoOpPublisher) PublishMultiple(_ context.Context, _ []Event) error {
	return nil
}

// Ensure NoOpPublisher implements the interface.
var _ Publisher = (*NoOpPublisher)(nil)
