package events

import (
	"context"

	coreEvents "github.com/archesai/archesai/internal/core/events"
)

// NoOpPublisher is a no-op event publisher implementation for testing.
type NoOpPublisher struct{}

// NewNoOpPublisher creates a new no-op event publisher.
func NewNoOpPublisher() coreEvents.Publisher {
	return &NoOpPublisher{}
}

// Publish does nothing in no-op implementation.
func (p *NoOpPublisher) Publish(_ context.Context, _ coreEvents.DomainEvent) error {
	return nil
}

// PublishMultiple does nothing in no-op implementation.
func (p *NoOpPublisher) PublishMultiple(_ context.Context, _ []coreEvents.DomainEvent) error {
	return nil
}

// Ensure NoOpPublisher implements the interface.
var _ coreEvents.Publisher = (*NoOpPublisher)(nil)
