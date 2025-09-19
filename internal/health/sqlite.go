package health

import (
	"context"
	"database/sql"
	"fmt"
)

// SQLiteRepository implements Repository for SQLite
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository creates a new SQLite health repository
func NewSQLiteRepository(db *sql.DB) Repository {
	return &SQLiteRepository{db: db}
}

// CheckDatabase performs a health check on the SQLite database connection
func (r *SQLiteRepository) CheckDatabase(ctx context.Context) error {
	if r.db == nil {
		return ErrDatabaseUnavailable
	}

	// Ping the database to check connectivity
	if err := r.db.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	// Run a simple query to ensure the database is responsive
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
func (r *SQLiteRepository) CheckRedis(_ context.Context) error {
	// TODO: Implement Redis health check when Redis is configured
	return nil
}

// CheckEmail performs a health check on email service
func (r *SQLiteRepository) CheckEmail(_ context.Context) error {
	// TODO: Implement email service health check when email is configured
	return nil
}
