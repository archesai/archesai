package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// RedisPublisher implements PublisherSubscriber using Redis pub/sub
type RedisPublisher struct {
	client  *redis.Client
	channel string
}

// NewRedisPublisher creates a new Redis event publisher with default channel
func NewRedisPublisher(client *redis.Client) PublisherSubscriber {
	return &RedisPublisher{
		client:  client,
		channel: "events", // Global channel for all events
	}
}

// Publish sends an event to Redis
func (p *RedisPublisher) Publish(ctx context.Context, event Event) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	// Marshal event to JSON
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Publish to main channel
	if err := p.client.Publish(ctx, p.channel, data).Err(); err != nil {
		return fmt.Errorf("failed to publish event to Redis: %w", err)
	}

	// Also publish to domain-specific channel for selective subscription
	domainChannel := fmt.Sprintf("%s:%s", p.channel, event.Domain)
	if err := p.client.Publish(ctx, domainChannel, data).Err(); err != nil {
		// Log error but don't fail - domain-specific channel is optional
		_ = err
	}

	// Also publish to type-specific channel
	typeChannel := fmt.Sprintf("%s:%s:%s", p.channel, event.Domain, event.Type)
	if err := p.client.Publish(ctx, typeChannel, data).Err(); err != nil {
		// Log error but don't fail - type-specific channel is optional
		_ = err
	}

	return nil
}

// PublishRaw publishes an event with domain context
func (p *RedisPublisher) PublishRaw(ctx context.Context, domain string, eventType string, data interface{}) error {
	event := Event{
		ID:        uuid.New().String(),
		Type:      eventType,
		Domain:    domain,
		Timestamp: time.Now().UTC(),
		Source:    domain,
		Data:      data,
	}

	// If data is already an event-like structure with ID, preserve it
	if dataMap, ok := data.(map[string]interface{}); ok {
		if id, exists := dataMap["id"]; exists {
			if idStr, ok := id.(string); ok {
				event.Metadata = map[string]string{
					"entity_id": idStr,
				}
			}
		}
	}

	return p.Publish(ctx, event)
}

// Subscribe subscribes to all events from a specific domain
func (p *RedisPublisher) Subscribe(ctx context.Context, domain string, handler func(event Event) error) error {
	// Subscribe to domain-specific channel
	domainChannel := fmt.Sprintf("%s:%s", p.channel, domain)

	pubsub := p.client.Subscribe(ctx, domainChannel)
	defer func() {
		_ = pubsub.Close()
	}()

	// Wait for subscription confirmation
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return fmt.Errorf("failed to subscribe to domain channel: %w", err)
	}

	// Get channel for messages
	ch := pubsub.Channel()

	// Process messages
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-ch:
			if msg == nil {
				continue
			}

			// Unmarshal event
			var event Event
			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				// Log error but continue processing
				continue
			}

			// Handle event
			if err := handler(event); err != nil {
				// Log error but continue processing
				continue
			}
		}
	}
}

// SubscribeToType subscribes to specific event types
func (p *RedisPublisher) SubscribeToType(ctx context.Context, domain string, eventType string, handler func(event Event) error) error {
	// Subscribe to type-specific channel
	typeChannel := fmt.Sprintf("%s:%s:%s", p.channel, domain, eventType)

	pubsub := p.client.Subscribe(ctx, typeChannel)
	defer func() {
		_ = pubsub.Close()
	}()

	// Wait for subscription confirmation
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return fmt.Errorf("failed to subscribe to type channel: %w", err)
	}

	// Get channel for messages
	ch := pubsub.Channel()

	// Process messages
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-ch:
			if msg == nil {
				continue
			}

			// Unmarshal event
			var event Event
			if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
				// Log error but continue processing
				continue
			}

			// Handle event
			if err := handler(event); err != nil {
				// Log error but continue processing
				continue
			}
		}
	}
}

// Ensure RedisPublisher implements the interfaces
var _ Publisher = (*RedisPublisher)(nil)
var _ PublisherSubscriber = (*RedisPublisher)(nil)
