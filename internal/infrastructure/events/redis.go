package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"

	coreEvents "github.com/archesai/archesai/internal/core/events"
)

// RedisPublisher implements PublisherSubscriber using Redis pub/sub.
type RedisPublisher struct {
	client  *redis.Client
	channel string
}

// NewRedisPublisher creates a new Redis event publisher with default channel.
func NewRedisPublisher(client *redis.Client) coreEvents.Publisher {
	return &RedisPublisher{
		client:  client,
		channel: "events", // Global channel for all events
	}
}

// Publish sends a domain event to Redis.
func (p *RedisPublisher) Publish(ctx context.Context, event coreEvents.DomainEvent) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	// Convert domain event to infrastructure event
	infraEvent := Event{
		ID:        event.EventID(),
		Type:      event.EventType(),
		Domain:    event.AggregateType(),
		Timestamp: event.OccurredAt(),
		Source:    event.AggregateType(),
		Data: map[string]any{
			"aggregate_id":   event.AggregateID(),
			"aggregate_type": event.AggregateType(),
			"event":          event,
		},
	}

	// Marshal event to JSON
	data, err := json.Marshal(infraEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Publish to main channel
	if err := p.client.Publish(ctx, p.channel, data).Err(); err != nil {
		return fmt.Errorf("failed to publish event to Redis: %w", err)
	}

	// Also publish to domain-specific channel for selective subscription
	domainChannel := fmt.Sprintf("%s:%s", p.channel, infraEvent.Domain)
	if err := p.client.Publish(ctx, domainChannel, data).Err(); err != nil {
		// Log error but don't fail - domain-specific channel is optional
		_ = err
	}

	// Also publish to type-specific channel
	typeChannel := fmt.Sprintf("%s:%s:%s", p.channel, infraEvent.Domain, infraEvent.Type)
	if err := p.client.Publish(ctx, typeChannel, data).Err(); err != nil {
		// Log error but don't fail - type-specific channel is optional
		_ = err
	}

	return nil
}

// PublishMultiple publishes multiple domain events in order.
func (p *RedisPublisher) PublishMultiple(
	ctx context.Context,
	events []coreEvents.DomainEvent,
) error {
	for _, event := range events {
		if err := p.Publish(ctx, event); err != nil {
			return fmt.Errorf("failed to publish event %s: %w", event.EventID(), err)
		}
	}
	return nil
}

// Ensure RedisPublisher implements the core Publisher interface.
var _ coreEvents.Publisher = (*RedisPublisher)(nil)
