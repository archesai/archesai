// Package events provides shared event infrastructure for all domains.
//
// This package implements a centralized event system that all domains use
// for publishing domain events. It provides:
//   - A common Publisher interface
//   - Redis-based event publishing and subscription
//   - NoOp implementation for testing
//
// Each domain generates its own event types but uses this shared infrastructure
// for actually publishing and subscribing to events.
package events

import (
	"context"
)

// Event is the interface that all domain events must implement.
type Event interface {
	// EventType returns the event type string (e.g., "user.created")
	EventType() string
	// EventData returns the actual event data
	EventData() any
}

// Publisher is the shared event publisher interface used by all domains.
type Publisher interface {
	// Publish sends an event to the event system
	Publish(ctx context.Context, event Event) error

	// PublishMultiple sends multiple events to the event system
	PublishMultiple(ctx context.Context, events []Event) error
}

// Subscriber handles event subscriptions.
type Subscriber interface {
	// Subscribe to all events from a specific domain
	Subscribe(ctx context.Context, domain string, handler func(event Event) error) error

	// SubscribeToType subscribes to specific event types
	SubscribeToType(
		ctx context.Context,
		domain string,
		eventType string,
		handler func(event Event) error,
	) error
}
