package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// PubSub provides pub/sub functionality using Redis
type PubSub struct {
	client *redis.Client
}

// NewPubSub creates a new PubSub instance
func NewPubSub(client *redis.Client) *PubSub {
	return &PubSub{
		client: client,
	}
}

// Publish publishes a message to a channel
func (p *PubSub) Publish(channel string, message interface{}) error {
	ctx := context.Background()

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return p.client.Publish(ctx, channel, data).Err()
}

// Subscribe subscribes to a channel and returns a subscription
func (p *PubSub) Subscribe(channels ...string) *Subscription {
	pubsub := p.client.Subscribe(context.Background(), channels...)
	return &Subscription{
		pubsub: pubsub,
	}
}

// Subscription represents an active subscription to Redis channels
type Subscription struct {
	pubsub *redis.PubSub
}

// Channel returns the channel for receiving messages
func (s *Subscription) Channel() <-chan *redis.Message {
	return s.pubsub.Channel()
}

// ReceiveMessage blocks and waits for a message
func (s *Subscription) ReceiveMessage(ctx context.Context) (*redis.Message, error) {
	return s.pubsub.ReceiveMessage(ctx)
}

// Close closes the subscription
func (s *Subscription) Close() error {
	return s.pubsub.Close()
}

// PublishJSON publishes a JSON message to a channel
func (p *PubSub) PublishJSON(channel string, v interface{}) error {
	return p.Publish(channel, v)
}

// PSubscribe subscribes to patterns and returns a subscription
func (p *PubSub) PSubscribe(patterns ...string) *Subscription {
	pubsub := p.client.PSubscribe(context.Background(), patterns...)
	return &Subscription{
		pubsub: pubsub,
	}
}
