// Package redis provides Redis-based event infrastructure implementations.
package redis

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/archesai/archesai/internal/core/events"
	coreevents "github.com/archesai/archesai/internal/core/events"
	infraevents "github.com/archesai/archesai/internal/infrastructure/events"
)

// Publisher implements the Publisher interface using Redis.
type Publisher struct {
	client    *redis.Client
	publisher infraevents.Publisher
}

// NewPublisher creates a new Redis-based event publisher.
func NewPublisher(client *redis.Client) events.Publisher {
	return &Publisher{
		client:    client,
		publisher: infraevents.NewRedisPublisher(client),
	}
}

// Publish publishes a single domain event.
func (p *Publisher) Publish(ctx context.Context, event coreevents.DomainEvent) error {
	// Convert domain event to infrastructure event
	infraEvent := infraevents.Event{
		ID:        uuid.New().String(),
		Type:      event.EventType(),
		Domain:    event.AggregateType(),
		Timestamp: event.OccurredAt(),
		Source:    fmt.Sprintf("%s:%s", event.AggregateType(), event.AggregateID()),
	}

	return p.publisher.Publish(ctx, infraEvent)
}

// PublishMultiple publishes multiple domain events in order.
func (p *Publisher) PublishMultiple(
	ctx context.Context,
	events []coreevents.DomainEvent,
) error {
	for _, event := range events {
		if err := p.Publish(ctx, event); err != nil {
			return fmt.Errorf("failed to publish event %s: %w", event.EventType(), err)
		}
	}
	return nil
}

// Subscriber implements the Subscriber interface using Redis.
type Subscriber struct {
	client     *redis.Client
	subscriber infraevents.Subscriber
}

// NewSubscriber creates a new Redis-based event subscriber.
func NewSubscriber(client *redis.Client) events.Subscriber {
	redisPublisher := infraevents.NewRedisPublisher(client)
	subscriber, ok := redisPublisher.(infraevents.Subscriber)
	if !ok {
		panic("redis publisher does not implement subscriber interface")
	}
	return &Subscriber{
		client:     client,
		subscriber: subscriber,
	}
}

// Subscribe subscribes to events of a specific aggregate type.
func (s *Subscriber) Subscribe(
	ctx context.Context,
	aggregateType string,
	handler events.EventHandler,
) error {
	return s.subscriber.Subscribe(ctx, aggregateType, func(event infraevents.Event) error {
		// Convert infrastructure event to domain event
		domainEvent := s.toDomainEvent(event)
		return handler(ctx, domainEvent)
	})
}

// SubscribeToEventType subscribes to a specific event type.
func (s *Subscriber) SubscribeToEventType(
	ctx context.Context,
	eventType string,
	handler events.EventHandler,
) error {
	// Extract domain from event type (e.g., "ToolCreated" -> "tools")
	domain := s.extractDomain(eventType)

	return s.subscriber.SubscribeToType(
		ctx,
		domain,
		eventType,
		func(event infraevents.Event) error {
			domainEvent := s.toDomainEvent(event)
			return handler(ctx, domainEvent)
		},
	)
}

// Unsubscribe removes a subscription.
func (s *Subscriber) Unsubscribe(ctx context.Context, subscriptionID string) error {
	// Redis pub/sub doesn't support explicit unsubscribe by ID
	// This would need to be implemented with a subscription manager
	return nil
}

// toDomainEvent converts an infrastructure event to a domain event.
func (s *Subscriber) toDomainEvent(event infraevents.Event) coreevents.DomainEvent {
	aggregateID := ""
	aggregateType := ""

	if event.Metadata != nil {
		aggregateID = event.Metadata["aggregate_id"]
		aggregateType = event.Metadata["aggregate_type"]
	}

	return &coreevents.BaseEvent{
		ID:               event.ID,
		Type:             event.Type,
		AggregateIDVal:   aggregateID,
		AggregateTypeVal: aggregateType,
		Timestamp:        event.Timestamp,
	}
}

// extractDomain extracts the domain from an event type.
func (s *Subscriber) extractDomain(eventType string) string {
	// Map event types to domains
	// This is a simple implementation - could be enhanced
	switch {
	case contains(eventType, "Tool"):
		return "tools"
	case contains(eventType, "Run"):
		return "runs"
	case contains(eventType, "Label"):
		return "labels"
	case contains(eventType, "User"):
		return "users"
	case contains(eventType, "Organization"):
		return "organizations"
	default:
		return "unknown"
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}

// EventBus combines publisher and subscriber for Redis.
type EventBus struct {
	events.Publisher
	events.Subscriber
}

// NewEventBus creates a new Redis-based event bus.
func NewEventBus(client *redis.Client) events.EventBus {
	return &EventBus{
		Publisher:  NewPublisher(client),
		Subscriber: NewSubscriber(client),
	}
}
