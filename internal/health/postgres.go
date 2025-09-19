package health

import (
	"context"
	"database/sql"
	"fmt"
)

// PostgresRepository implements Repository for PostgreSQL
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL health repository
func NewPostgresRepository(db *sql.DB) Repository {
	return &PostgresRepository{db: db}
}

// CheckDatabase performs a health check on the PostgreSQL database connection
func (r *PostgresRepository) CheckDatabase(ctx context.Context) error {
	if r.db == nil {
		return ErrDatabaseUnavailable
	}

	// Ping the database to check connectivity
	if err := r.db.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	// Optionally run a simple query to ensure the database is responsive
	var result int
	err := r.db.QueryRowContext(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("database query failed: %w", err)
	}

	if result != 1 {
		return fmt.Errorf("unexpected query result: %d", result)
	}

	return nil
}

// CheckRedis performs a health check on Redis connection
func (r *PostgresRepository) CheckRedis(_ context.Context) error {
	// TODO: Implement Redis health check when Redis is configured
	return nil
}

// CheckEmail performs a health check on email service
func (r *PostgresRepository) CheckEmail(_ context.Context) error {
	// TODO: Implement email service health check when email is configured
	return nil
}
