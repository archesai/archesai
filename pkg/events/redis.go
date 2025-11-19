package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// RedisPublisher implements PublisherSubscriber using Redis pub/sub.
type RedisPublisher struct {
	client  *redis.Client
	channel string
}

// NewRedisPublisher creates a new Redis event publisher with default channel.
func NewRedisPublisher(client *redis.Client) Publisher {
	return &RedisPublisher{
		client:  client,
		channel: "events", // Global channel for all events
	}
}

// Publish sends a domain event to Redis.
func (p *RedisPublisher) Publish(ctx context.Context, event Event) error {

	// Convert domain event to infrastructure event
	infraEvent := Event{
		ID:        event.ID,
		Type:      event.Type,
		Domain:    event.Domain,
		Timestamp: event.Timestamp,
		Source:    event.Source,
		Data: map[string]any{
			"event": event,
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
	events []Event,
) error {
	for _, event := range events {
		if err := p.Publish(ctx, event); err != nil {
			return fmt.Errorf("failed to publish event %s: %w", event.ID, err)
		}
	}
	return nil
}

// Ensure RedisPublisher implements the core Publisher interface.
var _ Publisher = (*RedisPublisher)(nil)
