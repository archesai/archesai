package health

import (
	"context"
	"log/slog"
	"testing"
	"time"
)

func TestNewService(t *testing.T) {
	logger := slog.Default()
	service := NewService(logger)

	if service == nil {
		t.Fatal("Expected non-nil service")
	}

	if service.logger != logger {
		t.Error("Expected logger to be set correctly")
	}

	if time.Since(service.startTime) > time.Second {
		t.Error("Expected start time to be recent")
	}
}

func TestService_CheckHealth(t *testing.T) {
	tests := []struct {
		name           string
		setupTime      time.Time
		expectedStatus ServiceStatus
		checkUptime    bool
	}{
		{
			name:      "all services healthy",
			setupTime: time.Now(),
			expectedStatus: ServiceStatus{
				Database: StatusHealthy,
				Email:    StatusHealthy,
				Redis:    StatusHealthy,
			},
			checkUptime: true,
		},
		{
			name:      "check after running for a while",
			setupTime: time.Now().Add(-5 * time.Second),
			expectedStatus: ServiceStatus{
				Database: StatusHealthy,
				Email:    StatusHealthy,
				Redis:    StatusHealthy,
			},
			checkUptime: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &Service{
				startTime: tt.setupTime,
				logger:    slog.Default(),
			}

			ctx := context.Background()
			status := service.CheckHealth(ctx)

			if status.Database != tt.expectedStatus.Database {
				t.Errorf("Expected Database status %s, got %s", tt.expectedStatus.Database, status.Database)
			}

			if status.Email != tt.expectedStatus.Email {
				t.Errorf("Expected Email status %s, got %s", tt.expectedStatus.Email, status.Email)
			}

			if status.Redis != tt.expectedStatus.Redis {
				t.Errorf("Expected Redis status %s, got %s", tt.expectedStatus.Redis, status.Redis)
			}

			if tt.checkUptime {
				expectedMinUptime := time.Since(tt.setupTime).Seconds()
				// Allow small variance for timing
				if status.Uptime < expectedMinUptime-0.1 || status.Uptime > expectedMinUptime+1 {
					t.Errorf("Expected uptime around %f, got %f", expectedMinUptime, status.Uptime)
				}
			}
		})
	}
}

func TestHealthConstants(t *testing.T) {
	// Test status constants
	if StatusHealthy != "healthy" {
		t.Errorf("Expected StatusHealthy to be 'healthy', got %s", StatusHealthy)
	}
	if StatusUnhealthy != "unhealthy" {
		t.Errorf("Expected StatusUnhealthy to be 'unhealthy', got %s", StatusUnhealthy)
	}
	if StatusDegraded != "degraded" {
		t.Errorf("Expected StatusDegraded to be 'degraded', got %s", StatusDegraded)
	}

	// Test timeout constants
	if DefaultTimeout != 5*time.Second {
		t.Errorf("Expected DefaultTimeout to be 5s, got %v", DefaultTimeout)
	}
	if DefaultInterval != 30*time.Second {
		t.Errorf("Expected DefaultInterval to be 30s, got %v", DefaultInterval)
	}
	if DefaultRetries != 3 {
		t.Errorf("Expected DefaultRetries to be 3, got %d", DefaultRetries)
	}
	if DefaultRetryDelay != 1*time.Second {
		t.Errorf("Expected DefaultRetryDelay to be 1s, got %v", DefaultRetryDelay)
	}

	// Test endpoint paths
	if LivenessPath != "/health/live" {
		t.Errorf("Expected LivenessPath to be '/health/live', got %s", LivenessPath)
	}
	if ReadinessPath != "/health/ready" {
		t.Errorf("Expected ReadinessPath to be '/health/ready', got %s", ReadinessPath)
	}
	if HealthPath != "/health" {
		t.Errorf("Expected HealthPath to be '/health', got %s", HealthPath)
	}

	// Test component names
	if ComponentAPI != "api" {
		t.Errorf("Expected ComponentAPI to be 'api', got %s", ComponentAPI)
	}
	if ComponentDatabase != "database" {
		t.Errorf("Expected ComponentDatabase to be 'database', got %s", ComponentDatabase)
	}
	if ComponentRedis != "redis" {
		t.Errorf("Expected ComponentRedis to be 'redis', got %s", ComponentRedis)
	}
	if ComponentStorage != "storage" {
		t.Errorf("Expected ComponentStorage to be 'storage', got %s", ComponentStorage)
	}
	if ComponentWorker != "worker" {
		t.Errorf("Expected ComponentWorker to be 'worker', got %s", ComponentWorker)
	}
}

func TestServiceStatus(t *testing.T) {
	status := ServiceStatus{
		Database: StatusHealthy,
		Email:    StatusDegraded,
		Redis:    StatusUnhealthy,
		Uptime:   123.45,
	}

	if status.Database != StatusHealthy {
		t.Errorf("Expected Database status to be %s, got %s", StatusHealthy, status.Database)
	}
	if status.Email != StatusDegraded {
		t.Errorf("Expected Email status to be %s, got %s", StatusDegraded, status.Email)
	}
	if status.Redis != StatusUnhealthy {
		t.Errorf("Expected Redis status to be %s, got %s", StatusUnhealthy, status.Redis)
	}
	if status.Uptime != 123.45 {
		t.Errorf("Expected Uptime to be 123.45, got %f", status.Uptime)
	}
}

func BenchmarkCheckHealth(b *testing.B) {
	service := NewService(slog.Default())
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.CheckHealth(ctx)
	}
}

func TestConcurrentHealthChecks(t *testing.T) {
	service := NewService(slog.Default())
	ctx := context.Background()

	// Run multiple health checks concurrently
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			status := service.CheckHealth(ctx)
			if status.Database != StatusHealthy {
				t.Errorf("Expected healthy database status in concurrent check")
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func TestUptimeCalculation(t *testing.T) {
	// Test uptime calculation with different start times
	testCases := []struct {
		name              string
		startTimeOffset   time.Duration
		expectedMinUptime float64
		expectedMaxUptime float64
	}{
		{
			name:              "just started",
			startTimeOffset:   0,
			expectedMinUptime: 0,
			expectedMaxUptime: 1,
		},
		{
			name:              "running for 10 seconds",
			startTimeOffset:   -10 * time.Second,
			expectedMinUptime: 10,
			expectedMaxUptime: 11,
		},
		{
			name:              "running for 1 minute",
			startTimeOffset:   -1 * time.Minute,
			expectedMinUptime: 60,
			expectedMaxUptime: 61,
		},
		{
			name:              "running for 1 hour",
			startTimeOffset:   -1 * time.Hour,
			expectedMinUptime: 3600,
			expectedMaxUptime: 3601,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := &Service{
				startTime: time.Now().Add(tc.startTimeOffset),
				logger:    slog.Default(),
			}

			status := service.CheckHealth(context.Background())

			if status.Uptime < tc.expectedMinUptime || status.Uptime > tc.expectedMaxUptime {
				t.Errorf("Expected uptime between %f and %f, got %f",
					tc.expectedMinUptime, tc.expectedMaxUptime, status.Uptime)
			}
		})
	}
}
