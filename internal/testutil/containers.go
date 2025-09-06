// Package testutil provides testing utilities for integration tests
package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Required for file-based migrations
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq" // PostgreSQL driver for database/sql
	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

// PostgresContainer wraps a PostgreSQL test container
type PostgresContainer struct {
	container testcontainers.Container
	DSN       string
	Pool      *pgxpool.Pool
}

// StartPostgresContainer starts a PostgreSQL container for testing
func StartPostgresContainer(ctx context.Context, t *testing.T) *PostgresContainer {
	t.Helper()

	container, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("testuser"),
		tcpostgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Failed to get container port: %v", err)
	}

	dsn := fmt.Sprintf("postgresql://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port.Port())

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	pc := &PostgresContainer{
		container: container,
		DSN:       dsn,
		Pool:      pool,
	}

	// Register cleanup
	t.Cleanup(func() {
		_ = pc.Stop(context.Background())
	})

	return pc
}

// RunMigrations runs database migrations on the test database
func (pc *PostgresContainer) RunMigrations(migrationsPath string) error {
	db, err := sql.Open("postgres", pc.DSN)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer func() { _ = db.Close() }()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+absPath,
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// Stop stops the PostgreSQL container
func (pc *PostgresContainer) Stop(ctx context.Context) error {
	if pc.Pool != nil {
		pc.Pool.Close()
	}
	if pc.container != nil {
		return pc.container.Terminate(ctx)
	}
	return nil
}

// RedisContainer wraps a Redis test container
type RedisContainer struct {
	container testcontainers.Container
	Client    *redis.Client
	Address   string
}

// StartRedisContainer starts a Redis container for testing
func StartRedisContainer(ctx context.Context, t *testing.T) *RedisContainer {
	t.Helper()

	container, err := tcredis.Run(ctx,
		"redis:7-alpine",
		tcredis.WithConfigFile(filepath.Join("testdata", "redis.conf")),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		// Try without config file if it doesn't exist
		container, err = tcredis.Run(ctx,
			"redis:7-alpine",
			testcontainers.WithWaitStrategy(
				wait.ForLog("Ready to accept connections").
					WithStartupTimeout(30*time.Second),
			),
		)
		if err != nil {
			t.Fatalf("Failed to start Redis container: %v", err)
		}
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "6379")
	if err != nil {
		t.Fatalf("Failed to get container port: %v", err)
	}

	address := fmt.Sprintf("%s:%s", host, port.Port())

	client := redis.NewClient(&redis.Options{
		Addr: address,
	})

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatalf("Failed to connect to Redis: %v", err)
	}

	rc := &RedisContainer{
		container: container,
		Client:    client,
		Address:   address,
	}

	// Register cleanup
	t.Cleanup(func() {
		_ = rc.Stop(context.Background())
	})

	return rc
}

// Stop stops the Redis container
func (rc *RedisContainer) Stop(ctx context.Context) error {
	if rc.Client != nil {
		if err := rc.Client.Close(); err != nil {
			return err
		}
	}
	if rc.container != nil {
		return rc.container.Terminate(ctx)
	}
	return nil
}
