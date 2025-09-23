package health

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"
)

// MockRepository is a mock implementation of Repository for testing
type MockRepository struct {
	dbErr    error
	redisErr error
	emailErr error
}

func (m *MockRepository) CheckDatabase(_ context.Context) error {
	return m.dbErr
}

func (m *MockRepository) CheckRedis(_ context.Context) error {
	return m.redisErr
}

func (m *MockRepository) CheckEmail(_ context.Context) error {
	return m.emailErr
}

func TestNewService(t *testing.T) {
	logger := slog.Default()
	repo := &MockRepository{}
	service := NewService(repo, logger)

	if service == nil {
		t.Fatal("Expected non-nil service")
	}

	if service.logger != logger {
		t.Error("Expected logger to be set correctly")
	}

	if service.repo != repo {
		t.Error("Expected repository to be set correctly")
	}

	if time.Since(service.start) > time.Second {
		t.Error("Expected start time to be recent")
	}
}

func TestService_CheckHealth(t *testing.T) {
	tests := []struct {
		name        string
		setupRepo   func() *MockRepository
		expectError bool
		wantDB      string
		wantRedis   string
		wantEmail   string
	}{
		{
			name: "all services healthy",
			setupRepo: func() *MockRepository {
				return &MockRepository{}
			},
			expectError: false,
			wantDB:      StatusHealthy,
			wantRedis:   StatusHealthy,
			wantEmail:   StatusHealthy,
		},
		{
			name: "database unhealthy",
			setupRepo: func() *MockRepository {
				return &MockRepository{
					dbErr: errors.New("database connection failed"),
				}
			},
			expectError: false,
			wantDB:      StatusUnhealthy,
			wantRedis:   StatusHealthy,
			wantEmail:   StatusHealthy,
		},
		{
			name: "redis unhealthy",
			setupRepo: func() *MockRepository {
				return &MockRepository{
					redisErr: errors.New("redis connection failed"),
				}
			},
			expectError: false,
			wantDB:      StatusHealthy,
			wantRedis:   StatusUnhealthy,
			wantEmail:   StatusHealthy,
		},
		{
			name: "email unhealthy",
			setupRepo: func() *MockRepository {
				return &MockRepository{
					emailErr: errors.New("email service unavailable"),
				}
			},
			expectError: false,
			wantDB:      StatusHealthy,
			wantRedis:   StatusHealthy,
			wantEmail:   StatusUnhealthy,
		},
		{
			name: "all services unhealthy",
			setupRepo: func() *MockRepository {
				return &MockRepository{
					dbErr:    errors.New("database connection failed"),
					redisErr: errors.New("redis connection failed"),
					emailErr: errors.New("email service unavailable"),
				}
			},
			expectError: false,
			wantDB:      StatusUnhealthy,
			wantRedis:   StatusUnhealthy,
			wantEmail:   StatusUnhealthy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.setupRepo()
			service := &Service{
				repo:   repo,
				logger: slog.Default(),
				start:  time.Now(),
			}

			response, err := service.CheckHealth(context.Background())

			if (err != nil) != tt.expectError {
				t.Errorf("CheckHealth() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if response == nil {
				t.Fatal("Expected non-nil response")
			}

			if response.Services.Database != tt.wantDB {
				t.Errorf("Database status = %v, want %v", response.Services.Database, tt.wantDB)
			}

			if response.Services.Redis != tt.wantRedis {
				t.Errorf("Redis status = %v, want %v", response.Services.Redis, tt.wantRedis)
			}

			if response.Services.Email != tt.wantEmail {
				t.Errorf("Email status = %v, want %v", response.Services.Email, tt.wantEmail)
			}

			if response.Timestamp == "" {
				t.Error("Expected timestamp to be set")
			}

			if response.Uptime <= 0 {
				t.Error("Expected uptime to be greater than 0")
			}
		})
	}
}

func TestService_CheckHealthWithUptime(t *testing.T) {
	repo := &MockRepository{}
	startTime := time.Now().Add(-5 * time.Minute) // Service started 5 minutes ago
	service := &Service{
		repo:   repo,
		logger: slog.Default(),
		start:  startTime,
	}

	response, err := service.CheckHealth(context.Background())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Uptime should be approximately 300 seconds (5 minutes)
	expectedUptime := float64(300)
	tolerance := float64(2) // Allow 2 seconds tolerance

	if response.Uptime < expectedUptime-tolerance || response.Uptime > expectedUptime+tolerance {
		t.Errorf("Uptime = %v, want approximately %v", response.Uptime, expectedUptime)
	}
}
