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
	"time"
)

// Event represents a domain event with common fields
type Event struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Domain    string            `json:"domain"` // e.g., "auth", "organizations"
	Timestamp time.Time         `json:"timestamp"`
	Source    string            `json:"source"`
	Data      interface{}       `json:"data"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// Publisher is the shared event publisher interface used by all domains
type Publisher interface {
	// Publish sends an event to the event system
	Publish(ctx context.Context, event Event) error

	// PublishRaw publishes an event with domain context
	PublishRaw(ctx context.Context, domain string, eventType string, data interface{}) error
}

// Subscriber handles event subscriptions
type Subscriber interface {
	// Subscribe to all events from a specific domain
	Subscribe(ctx context.Context, domain string, handler func(event Event) error) error

	// SubscribeToType subscribes to specific event types
	SubscribeToType(ctx context.Context, domain string, eventType string, handler func(event Event) error) error
}

// PublisherSubscriber combines publishing and subscribing capabilities
type PublisherSubscriber interface {
	Publisher
	Subscriber
}
