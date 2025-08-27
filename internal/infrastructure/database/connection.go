package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// DB represents a database connection pool
type DB struct {
	*pgxpool.Pool
	logger *zap.Logger
}

// Config holds database configuration
type Config struct {
	URL               string
	MaxConns          int32
	MinConns          int32
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
}

// NewConnection creates a new database connection pool
func NewConnection(cfg Config, logger *zap.Logger) (*DB, error) {
	ctx := context.Background()

	// Parse the connection string
	config, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Configure connection pool
	config.MaxConns = cfg.MaxConns
	if config.MaxConns == 0 {
		config.MaxConns = 25
	}

	config.MinConns = cfg.MinConns
	if config.MinConns == 0 {
		config.MinConns = 5
	}

	config.MaxConnLifetime = cfg.MaxConnLifetime
	if config.MaxConnLifetime == 0 {
		config.MaxConnLifetime = 1 * time.Hour
	}

	config.MaxConnIdleTime = cfg.MaxConnIdleTime
	if config.MaxConnIdleTime == 0 {
		config.MaxConnIdleTime = 30 * time.Minute
	}

	config.HealthCheckPeriod = cfg.HealthCheckPeriod
	if config.HealthCheckPeriod == 0 {
		config.HealthCheckPeriod = 1 * time.Minute
	}

	// Set up connection hooks for logging
	config.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		logger.Debug("acquiring database connection")
		return true
	}

	config.AfterRelease = func(conn *pgx.Conn) bool {
		logger.Debug("releasing database connection")
		return true
	}

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("database connection established",
		zap.Int32("max_connections", config.MaxConns),
		zap.Int32("min_connections", config.MinConns),
	)

	return &DB{
		Pool:   pool,
		logger: logger,
	}, nil
}

// Close closes the database connection pool
func (db *DB) Close() {
	db.logger.Info("closing database connection pool")
	db.Pool.Close()
}

// Health checks the health of the database connection
func (db *DB) Health(ctx context.Context) error {
	return db.Pool.Ping(ctx)
}

// Stats returns database pool statistics
func (db *DB) Stats() *pgxpool.Stat {
	return db.Pool.Stat()
}

// Transaction executes a function within a database transaction
func (db *DB) Transaction(ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			// Rollback on panic
			tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			// Rollback on error
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				db.logger.Error("failed to rollback transaction", zap.Error(rbErr))
			}
		} else {
			// Commit on success
			err = tx.Commit(ctx)
			if err != nil {
				db.logger.Error("failed to commit transaction", zap.Error(err))
			}
		}
	}()

	err = fn(tx)
	return err
}

// Migrate runs database migrations
func (db *DB) Migrate(ctx context.Context, migrationsPath string) error {
	// This is a placeholder for migration logic
	// In a real implementation, you would use a migration tool like golang-migrate
	db.logger.Info("running database migrations", zap.String("path", migrationsPath))

	// Example: Check if tables exist
	var exists bool
	err := db.Pool.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'user'
		)
	`).Scan(&exists)

	if err != nil {
		return fmt.Errorf("failed to check database schema: %w", err)
	}

	if !exists {
		db.logger.Info("database schema not found, migrations needed")
		// Run migrations here
	} else {
		db.logger.Info("database schema already exists")
	}

	return nil
}

// QueryBuilder provides a fluent interface for building queries
type QueryBuilder struct {
	db     *DB
	query  string
	args   []interface{}
	logger *zap.Logger
}

// NewQueryBuilder creates a new query builder
func (db *DB) NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		db:     db,
		args:   make([]interface{}, 0),
		logger: db.logger,
	}
}

// Select sets the SELECT clause
func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
	qb.query = "SELECT " + joinStrings(columns, ", ")
	return qb
}

// From sets the FROM clause
func (qb *QueryBuilder) From(table string) *QueryBuilder {
	qb.query += " FROM " + table
	return qb
}

// Where adds a WHERE clause
func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
	qb.query += " WHERE " + condition
	qb.args = append(qb.args, args...)
	return qb
}

// OrderBy adds an ORDER BY clause
func (qb *QueryBuilder) OrderBy(column string, direction string) *QueryBuilder {
	qb.query += fmt.Sprintf(" ORDER BY %s %s", column, direction)
	return qb
}

// Limit adds a LIMIT clause
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.query += fmt.Sprintf(" LIMIT %d", limit)
	return qb
}

// Offset adds an OFFSET clause
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.query += fmt.Sprintf(" OFFSET %d", offset)
	return qb
}

// Build returns the built query and arguments
func (qb *QueryBuilder) Build() (string, []interface{}) {
	return qb.query, qb.args
}

// Helper function to join strings
func joinStrings(strings []string, separator string) string {
	result := ""
	for i, s := range strings {
		if i > 0 {
			result += separator
		}
		result += s
	}
	return result
}
