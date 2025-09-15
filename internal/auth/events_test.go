package auth

import (
	"testing"
	"time"

	"github.com/archesai/archesai/internal/accounts"
	"github.com/archesai/archesai/internal/sessions"
	"github.com/google/uuid"
)

func TestEventSystem(t *testing.T) {
	t.Skip("Skipping test - event publisher not implemented")
	// Test NoOpEventPublisher implementation
	// publisher := auth.NewNoOpEventPublisher()
	/*

		ctx := context.Background()
		account := &accounts.Account{
			Id:         uuid.New(),
			AccountId:  "test-account-123",
			ProviderId: accounts.Local,
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
		session := &sessions.Session{
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
	*/
}

func TestEventTypes(t *testing.T) {
	// Verify event type constants are defined correctly
	expectedTypes := map[string]string{
		"account.created":   accounts.EventAccountCreated,
		"account.updated":   accounts.EventAccountUpdated,
		"account.deleted":   accounts.EventAccountDeleted,
		"account.linked":    accounts.EventAccountLinked,
		"account.unlinked":  accounts.EventAccountUnlinked,
		"session.created":   sessions.EventSessionCreated,
		"session.refreshed": sessions.EventSessionRefreshed,
		"session.expired":   sessions.EventSessionExpired,
		"session.deleted":   sessions.EventSessionDeleted,
	}

	for expected, actual := range expectedTypes {
		if actual != expected {
			t.Errorf("EventType mismatch: expected %s, got %s", expected, actual)
		}
	}
}

func TestEventStructures(t *testing.T) {
	account := &accounts.Account{
		Id:         uuid.New(),
		AccountId:  "test-account-789",
		ProviderId: accounts.Google,
		UserId:     uuid.New(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Test event creation functions
	createdEvent := accounts.NewAccountCreatedEvent(account)
	if createdEvent.Account != account {
		t.Error("Event entity does not match original")
	}
	if createdEvent.EventType() != accounts.EventAccountCreated {
		t.Errorf("Event type mismatch: expected %s, got %s", accounts.EventAccountCreated, createdEvent.EventType())
	}
	if createdEvent.EventDomain() != "accounts" {
		t.Errorf("Event domain mismatch: expected accounts, got %s", createdEvent.EventDomain())
	}

	// Verify event has required fields
	if createdEvent.ID == "" {
		t.Error("Event ID is empty")
	}
	if createdEvent.Timestamp.IsZero() {
		t.Error("Event timestamp is zero")
	}
	if createdEvent.EventData() != account {
		t.Error("Event data does not match entity")
	}
}
