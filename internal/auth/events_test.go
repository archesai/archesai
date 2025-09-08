package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/archesai/archesai/internal/auth"
	"github.com/google/uuid"
)

func TestEventSystem(t *testing.T) {
	// Test NoOpEventPublisher implementation
	publisher := auth.NewNoOpEventPublisher()

	ctx := context.Background()
	account := &auth.Account{
		Id:         uuid.New(),
		AccountId:  "test-account-123",
		ProviderId: auth.Local,
		UserId:     uuid.New(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Test all event publishing methods
	tests := []struct {
		name string
		fn   func() error
	}{
		{
			name: "PublishAccountCreated",
			fn: func() error {
				return publisher.PublishAccountCreated(ctx, account)
			},
		},
		{
			name: "PublishAccountUpdated",
			fn: func() error {
				return publisher.PublishAccountUpdated(ctx, account)
			},
		},
		{
			name: "PublishAccountDeleted",
			fn: func() error {
				return publisher.PublishAccountDeleted(ctx, account)
			},
		},
		{
			name: "PublishAccountLinked",
			fn: func() error {
				return publisher.PublishAccountLinked(ctx, account)
			},
		},
		{
			name: "PublishAccountUnlinked",
			fn: func() error {
				return publisher.PublishAccountUnlinked(ctx, account)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fn(); err != nil {
				t.Errorf("%s() returned error: %v", tt.name, err)
			}
		})
	}

	// Test session events
	session := &auth.Session{
		Id:                   uuid.New(),
		Token:                "test-token",
		ActiveOrganizationId: uuid.New(),
		ExpiresAt:            time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		IpAddress:            "127.0.0.1",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	sessionTests := []struct {
		name string
		fn   func() error
	}{
		{
			name: "PublishSessionCreated",
			fn: func() error {
				return publisher.PublishSessionCreated(ctx, session)
			},
		},
		{
			name: "PublishSessionRefreshed",
			fn: func() error {
				return publisher.PublishSessionRefreshed(ctx, session)
			},
		},
		{
			name: "PublishSessionExpired",
			fn: func() error {
				return publisher.PublishSessionExpired(ctx, session)
			},
		},
		{
			name: "PublishSessionDeleted",
			fn: func() error {
				return publisher.PublishSessionDeleted(ctx, session)
			},
		},
	}

	for _, tt := range sessionTests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fn(); err != nil {
				t.Errorf("%s() returned error: %v", tt.name, err)
			}
		})
	}
}

func TestEventTypes(t *testing.T) {
	// Verify event type constants are defined correctly
	expectedTypes := map[string]auth.EventType{
		"account.created":   auth.EventAccountCreated,
		"account.updated":   auth.EventAccountUpdated,
		"account.deleted":   auth.EventAccountDeleted,
		"account.linked":    auth.EventAccountLinked,
		"account.unlinked":  auth.EventAccountUnlinked,
		"session.created":   auth.EventSessionCreated,
		"session.refreshed": auth.EventSessionRefreshed,
		"session.expired":   auth.EventSessionExpired,
		"session.deleted":   auth.EventSessionDeleted,
	}

	for expected, actual := range expectedTypes {
		if string(actual) != expected {
			t.Errorf("EventType mismatch: expected %s, got %s", expected, actual)
		}
	}
}

func TestEventStructures(t *testing.T) {
	account := &auth.Account{
		Id:         uuid.New(),
		AccountId:  "test-account-789",
		ProviderId: auth.Google,
		UserId:     uuid.New(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Test event creation functions
	createdEvent := auth.NewAccountCreatedEvent(account)
	if createdEvent.Entity != account {
		t.Error("Event entity does not match original")
	}
	if createdEvent.Type != auth.EventAccountCreated {
		t.Errorf("Event type mismatch: expected %s, got %s", auth.EventAccountCreated, createdEvent.Type)
	}
	if createdEvent.Source != "auth" {
		t.Errorf("Event source mismatch: expected auth, got %s", createdEvent.Source)
	}

	// Verify event has required fields
	if createdEvent.ID == "" {
		t.Error("Event ID is empty")
	}
	if createdEvent.Timestamp.IsZero() {
		t.Error("Event timestamp is zero")
	}
	if createdEvent.Data != account {
		t.Error("Event data does not match entity")
	}
}
