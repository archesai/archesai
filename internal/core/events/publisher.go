// Package events defines event-related port interfaces.
package events

import (
	"context"
)

// Publisher defines the interface for publishing domain events.
type Publisher interface {
	// Publish publishes a single domain event.
	Publish(ctx context.Context, event DomainEvent) error

	// PublishMultiple publishes multiple domain events in order.
	PublishMultiple(ctx context.Context, events []DomainEvent) error
}

// Subscriber defines the interface for subscribing to domain events.
type Subscriber interface {
	// Subscribe subscribes to events of a specific aggregate type.
	Subscribe(ctx context.Context, aggregateType string, handler EventHandler) error

	// SubscribeToEventType subscribes to a specific event type.
	SubscribeToEventType(ctx context.Context, eventType string, handler EventHandler) error

	// Unsubscribe removes a subscription.
	Unsubscribe(ctx context.Context, subscriptionID string) error
}

// EventHandler is a function that handles domain events.
type EventHandler func(ctx context.Context, event DomainEvent) error

// EventStore defines the interface for persisting and retrieving domain events.
type EventStore interface {
	// Save persists a domain event.
	Save(ctx context.Context, event DomainEvent) error

	// SaveMultiple persists multiple domain events in a transaction.
	SaveMultiple(ctx context.Context, events []DomainEvent) error

	// GetEvents retrieves all events for an aggregate.
	GetEvents(ctx context.Context, aggregateID string) ([]DomainEvent, error)

	// GetEventsAfterVersion retrieves events for an aggregate after a specific version.
	GetEventsAfterVersion(
		ctx context.Context,
		aggregateID string,
		version int64,
	) ([]DomainEvent, error)

	// GetEventsByType retrieves events of a specific type.
	GetEventsByType(
		ctx context.Context,
		eventType string,
		limit int,
		offset int,
	) ([]DomainEvent, error)
}

// EventBus combines publishing and subscribing capabilities.
type EventBus interface {
	Publisher
	Subscriber
}
