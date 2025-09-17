package events

import (
	"context"
)

// NoOpPublisher is a no-op event publisher implementation for testing.
type NoOpPublisher struct{}

// NewNoOpPublisher creates a new no-op event publisher.
func NewNoOpPublisher() PublisherSubscriber {
	return &NoOpPublisher{}
}

// Publish does nothing in no-op implementation.
func (p *NoOpPublisher) Publish(_ context.Context, _ Event) error {
	return nil
}

// PublishRaw does nothing in no-op implementation.
func (p *NoOpPublisher) PublishRaw(_ context.Context, _ string, _ string, _ interface{}) error {
	return nil
}

// Subscribe does nothing in no-op implementation.
func (p *NoOpPublisher) Subscribe(_ context.Context, _ string, _ func(event Event) error) error {
	return nil
}

// SubscribeToType does nothing in no-op implementation.
func (p *NoOpPublisher) SubscribeToType(
	_ context.Context,
	_ string,
	_ string,
	_ func(event Event) error,
) error {
	return nil
}

// Ensure NoOpPublisher implements the interfaces.
var _ Publisher = (*NoOpPublisher)(nil)
var _ PublisherSubscriber = (*NoOpPublisher)(nil)
