package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var _ Publisher = (*RedisPublisher)(nil)

// RedisPublisher implements Publisher using Redis pub/sub.
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

	// Marshal event to JSON
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Publish to main channel
	if err := p.client.Publish(ctx, p.channel, data).Err(); err != nil {
		return fmt.Errorf("failed to publish event to Redis: %w", err)
	}

	// Also publish to type-specific channel
	typeChannel := fmt.Sprintf("%s:%s", p.channel, event.EventType())
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
			return fmt.Errorf("failed to publish event %s: %w", event.EventData(), err)
		}
	}
	return nil
}
