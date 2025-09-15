package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/archesai/archesai/internal/config"
	_ "modernc.org/sqlite" // Pure Go SQLite driver
)

// SQLite implements the Database interface for SQLite
type SQLite struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewSQLite creates a new SQLite database connection
func NewSQLite(cfg *config.DatabaseConfig, logger *slog.Logger) (Database, error) {
	ctx := context.Background()

	// Parse connection URL
	connStr := cfg.URL
	if connStr == "" {
		connStr = ":memory:" // Default to in-memory database
	}

	// Open database connection
	db, err := sql.Open("sqlite", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	// Configure connection pool
	// SQLite performs best with a single connection
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if cfg.ConnMaxLifetime != "" && cfg.ConnMaxLifetime != "0" {
		duration, err := time.ParseDuration(cfg.ConnMaxLifetime)
		if err == nil {
			db.SetConnMaxLifetime(duration)
		}
	}

	if cfg.ConnMaxIdleTime != "" && cfg.ConnMaxIdleTime != "0" {
		duration, err := time.ParseDuration(cfg.ConnMaxIdleTime)
		if err == nil {
			db.SetConnMaxIdleTime(duration)
		}
	}

	// Test connection
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping SQLite database: %w", err)
	}

	// Configure SQLite pragmas
	if err := configureSQLitePragmas(db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to configure SQLite pragmas: %w", err)
	}

	logger.Info("SQLite connection established",
		"database", connStr,
	)

	return &SQLite{
		db:     db,
		logger: logger,
	}, nil
}

// configureSQLitePragmas sets SQLite-specific pragmas for performance and reliability
func configureSQLitePragmas(db *sql.DB) error {
	pragmas := []string{
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous = NORMAL",
		"PRAGMA cache_size = -2000", // 2MB cache
		"PRAGMA foreign_keys = ON",
		"PRAGMA busy_timeout = 5000",
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			return fmt.Errorf("failed to execute %s: %w", pragma, err)
		}
	}

	return nil
}

// Query executes a query that returns rows
func (db *SQLite) Query(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	rows, err := db.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &sqlRows{rows: rows}, nil
}

// QueryRow executes a query that returns at most one row
func (db *SQLite) QueryRow(ctx context.Context, query string, args ...interface{}) Row {
	return &sqlRow{row: db.db.QueryRowContext(ctx, query, args...)}
}

// Exec executes a query without returning any rows
func (db *SQLite) Exec(ctx context.Context, query string, args ...interface{}) (Result, error) {
	return db.db.ExecContext(ctx, query, args...)
}

// Begin starts a transaction
func (db *SQLite) Begin(ctx context.Context) (Transaction, error) {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &sqlTransaction{tx: tx}, nil
}

// BeginTx starts a transaction with options
func (db *SQLite) BeginTx(ctx context.Context, opts *sql.TxOptions) (Transaction, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &sqlTransaction{tx: tx}, nil
}

// Ping verifies the connection to the database
func (db *SQLite) Ping(ctx context.Context) error {
	return db.db.PingContext(ctx)
}

// Close closes the database connection
func (db *SQLite) Close() error {
	db.logger.Info("Closing SQLite connection")
	return db.db.Close()
}

// Stats returns database statistics
func (db *SQLite) Stats() Stats {
	stats := db.db.Stats()
	return Stats{
		OpenConnections:   stats.OpenConnections,
		InUse:             stats.InUse,
		Idle:              stats.Idle,
		WaitCount:         stats.WaitCount,
		WaitDuration:      stats.WaitDuration.String(),
		MaxIdleClosed:     stats.MaxIdleClosed,
		MaxLifetimeClosed: stats.MaxLifetimeClosed,
	}
}

// Type returns the database type
func (db *SQLite) Type() Type {
	return TypeSQLite
}

// Underlying returns the underlying database connection
func (db *SQLite) Underlying() interface{} {
	return db.db
}

// Adapter types for database/sql compatibility

type sqlRows struct {
	rows *sql.Rows
}

func (r *sqlRows) Next() bool {
	return r.rows.Next()
}

func (r *sqlRows) Scan(dest ...interface{}) error {
	return r.rows.Scan(dest...)
}

func (r *sqlRows) Close() error {
	return r.rows.Close()
}

func (r *sqlRows) Err() error {
	return r.rows.Err()
}

type sqlRow struct {
	row *sql.Row
}

func (r *sqlRow) Scan(dest ...interface{}) error {
	return r.row.Scan(dest...)
}

type sqlTransaction struct {
	tx *sql.Tx
}

func (t *sqlTransaction) Query(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	rows, err := t.tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &sqlRows{rows: rows}, nil
}

func (t *sqlTransaction) QueryRow(ctx context.Context, query string, args ...interface{}) Row {
	return &sqlRow{row: t.tx.QueryRowContext(ctx, query, args...)}
}

func (t *sqlTransaction) Exec(ctx context.Context, query string, args ...interface{}) (Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

func (t *sqlTransaction) Commit() error {
	return t.tx.Commit()
}

func (t *sqlTransaction) Rollback() error {
	return t.tx.Rollback()
}
