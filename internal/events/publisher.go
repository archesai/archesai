package events

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DomainEvent is the interface that all domain events must implement.
type DomainEvent interface {
	// EventType returns the event type string (e.g., "user.created")
	EventType() string
	// EventDomain returns the domain this event belongs to (e.g., "users")
	EventDomain() string
	// EventData returns the actual event data
	EventData() interface{}
}

// BaseEvent provides common fields for all events.
type BaseEvent struct {
	ID        string    `json:"id"`
	Domain    string    `json:"domain"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version,omitempty"`
}

// NewBaseEvent creates a new base event.
func NewBaseEvent(domain, eventType string) BaseEvent {
	return BaseEvent{
		ID:        uuid.New().String(),
		Domain:    domain,
		Type:      eventType,
		Timestamp: time.Now().UTC(),
		Version:   "1.0",
	}
}

// PublishDomainEvent is a helper function to publish domain events.
func PublishDomainEvent(ctx context.Context, publisher Publisher, event DomainEvent) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	// Create the full event structure
	fullEvent := Event{
		ID:        uuid.New().String(),
		Type:      event.EventType(),
		Domain:    event.EventDomain(),
		Timestamp: time.Now().UTC(),
		Source:    event.EventDomain(),
		Data:      event.EventData(),
	}

	return publisher.Publish(ctx, fullEvent)
}
