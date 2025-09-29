package events

import (
	"time"

	"github.com/google/uuid"
)

// DomainEvent represents a domain event that has occurred in the system.
type DomainEvent interface {
	EventID() string
	EventType() string
	AggregateID() string
	AggregateType() string
	OccurredAt() time.Time
}

// BaseEvent provides common fields for all domain events.
type BaseEvent struct {
	ID               string    `json:"id"`
	Type             string    `json:"type"`
	AggregateIDVal   string    `json:"aggregate_id"`
	AggregateTypeVal string    `json:"aggregate_type"`
	Timestamp        time.Time `json:"timestamp"`
}

// NewBaseEvent creates a new base event with the given aggregate type and event type.
func NewBaseEvent(aggregateType, eventType string) BaseEvent {
	return BaseEvent{
		ID:               uuid.New().String(),
		Type:             eventType,
		AggregateTypeVal: aggregateType,
		Timestamp:        time.Now().UTC(),
	}
}

// EventID returns the event's unique identifier.
func (e BaseEvent) EventID() string {
	return e.ID
}

// EventType returns the type of the event.
func (e BaseEvent) EventType() string {
	return e.Type
}

// AggregateID returns the aggregate ID.
func (e BaseEvent) AggregateID() string {
	return e.AggregateIDVal
}

// AggregateType returns the aggregate type.
func (e BaseEvent) AggregateType() string {
	return e.AggregateTypeVal
}

// OccurredAt returns when the event occurred.
func (e BaseEvent) OccurredAt() time.Time {
	return e.Timestamp
}
