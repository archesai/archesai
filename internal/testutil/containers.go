// Package testutil provides testing utilities for integration tests
package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver for database/sql
	"github.com/pressly/goose/v3"
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
	db, err := sql.Open("pgx", pc.DSN)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer func() { _ = db.Close() }()

	// Set environment variables for PostgreSQL
	_ = os.Setenv("TIMESTAMP_TYPE", "TIMESTAMPTZ")
	_ = os.Setenv("TIMESTAMP_DEFAULT", "CURRENT_TIMESTAMP")
	_ = os.Setenv("REAL_TYPE", "DOUBLE PRECISION")

	// Set goose dialect
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	// Run migrations
	if err := goose.Up(db, migrationsPath); err != nil {
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
