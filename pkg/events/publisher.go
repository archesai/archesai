package events

import (
	"time"

	"github.com/google/uuid"
)

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
