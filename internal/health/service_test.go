package health

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/archesai/archesai/internal/database"
)

func TestNewService(t *testing.T) {
	logger := slog.Default()
	// Database is an interface, use nil for testing
	var db *database.Database
	service := NewService(db, logger)

	if service == nil {
		t.Fatal("Expected non-nil service")
	}

	if service.logger != logger {
		t.Error("Expected logger to be set correctly")
	}

	// Can't compare interface values directly
	if service.db != nil {
		t.Error("Expected db to be nil for this test")
	}

	if time.Since(service.start) > time.Second {
		t.Error("Expected start time to be recent")
	}
}

func TestService_CheckHealth(t *testing.T) {
	tests := []struct {
		name        string
		setupTime   time.Time
		withDB      bool
		expectError bool
	}{
		{
			name:        "health check with nil database",
			setupTime:   time.Now(),
			withDB:      true,
			expectError: false,
		},
		{
			name:        "health check without database",
			setupTime:   time.Now(),
			withDB:      false,
			expectError: false, // We don't error, just report unhealthy
		},
		{
			name:        "health check after running for a while",
			setupTime:   time.Now().Add(-5 * time.Second),
			withDB:      true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.Default()
			// Database is an interface, can't instantiate directly - using nil for unit tests
			var db *database.Database

			service := &Service{
				db:     db,
				logger: logger,
				start:  tt.setupTime,
			}

			ctx := context.Background()
			status, err := service.CheckHealth(ctx)

			if tt.expectError && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if status != nil {
				// Check database status - always unhealthy in unit tests since we use nil
				expectedDBStatus := StatusUnhealthy

				if status.Services.Database != expectedDBStatus {
					t.Errorf("Database status = %v, want %v", status.Services.Database, expectedDBStatus)
				}

				// Email and Redis should always be healthy (TODO in implementation)
				if status.Services.Email != StatusHealthy {
					t.Errorf("Email status = %v, want %v", status.Services.Email, StatusHealthy)
				}
				if status.Services.Redis != StatusHealthy {
					t.Errorf("Redis status = %v, want %v", status.Services.Redis, StatusHealthy)
				}

				// Check timestamp format
				if status.Timestamp == "" {
					t.Error("Expected timestamp to be set")
				}

				// Check uptime
				expectedUptime := time.Since(tt.setupTime).Seconds()
				if status.Uptime < 0 || status.Uptime > float32(expectedUptime+1) {
					t.Errorf("Uptime = %v, expected around %v", status.Uptime, expectedUptime)
				}
			}
		})
	}
}

func TestService_checkDatabase(t *testing.T) {
	tests := []struct {
		name        string
		db          *database.Database
		expectError bool
	}{
		{
			name:        "nil database connection returns error",
			db:          nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.Default()
			service := &Service{
				db:     tt.db,
				logger: logger,
				start:  time.Now(),
			}

			ctx := context.Background()
			err := service.checkDatabase(ctx)

			if tt.expectError && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
