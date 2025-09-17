package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // PostgreSQL driver for database/sql compatibility

	"github.com/archesai/archesai/internal/config"
)

// PGDatabase implements the Database interface for PostgreSQL.
type PGDatabase struct {
	pool   *pgxpool.Pool
	sqlDB  *sql.DB // For stdlib compatibility when needed
	logger *slog.Logger
}

// NewPostgreSQL creates a new PostgreSQL database connection.
func NewPostgreSQL(cfg *config.DatabaseConfig, logger *slog.Logger) (Database, error) {
	ctx := context.Background()

	// Parse the connection string
	poolConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Configure connection pool with defaults
	if cfg.MaxConns > 0 {
		poolConfig.MaxConns = int32(cfg.MaxConns)
	} else {
		poolConfig.MaxConns = 25
	}

	if cfg.MinConns > 0 {
		poolConfig.MinConns = int32(cfg.MinConns)
	} else {
		poolConfig.MinConns = 5
	}

	if cfg.ConnMaxLifetime != "" {
		duration, err := time.ParseDuration(cfg.ConnMaxLifetime)
		if err == nil {
			poolConfig.MaxConnLifetime = duration
		}
	}

	if cfg.ConnMaxIdleTime != "" {
		duration, err := time.ParseDuration(cfg.ConnMaxIdleTime)
		if err == nil {
			poolConfig.MaxConnIdleTime = duration
		}
	}

	if cfg.HealthCheckPeriod != "" {
		duration, err := time.ParseDuration(cfg.HealthCheckPeriod)
		if err == nil {
			poolConfig.HealthCheckPeriod = duration
		}
	}

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create stdlib connection for compatibility
	sqlDB, err := sql.Open("pgx", cfg.URL)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to create stdlib connection: %w", err)
	}

	logger.Info("PostgreSQL connection established",
		"max_connections", poolConfig.MaxConns,
		"min_connections", poolConfig.MinConns,
	)

	return &PGDatabase{
		pool:   pool,
		sqlDB:  sqlDB,
		logger: logger,
	}, nil
}

// Query executes a query that returns rows.
func (db *PGDatabase) Query(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	rows, err := db.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &pgRows{rows: rows}, nil
}

// QueryRow executes a query that returns at most one row.
func (db *PGDatabase) QueryRow(ctx context.Context, query string, args ...interface{}) Row {
	return &pgRow{row: db.pool.QueryRow(ctx, query, args...)}
}

// Exec executes a query without returning any rows.
func (db *PGDatabase) Exec(ctx context.Context, query string, args ...interface{}) (Result, error) {
	result, err := db.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &pgResult{result: result}, nil
}

// Begin starts a transaction.
func (db *PGDatabase) Begin(ctx context.Context) (Transaction, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &pgTransaction{tx: tx}, nil
}

// BeginTx starts a transaction with options.
func (db *PGDatabase) BeginTx(ctx context.Context, opts *sql.TxOptions) (Transaction, error) {
	// Map SQL isolation levels to pgx string values
	var isoLevel pgx.TxIsoLevel
	switch opts.Isolation {
	case sql.LevelDefault:
		isoLevel = pgx.ReadCommitted // PostgreSQL default
	case sql.LevelReadUncommitted:
		isoLevel = pgx.ReadUncommitted
	case sql.LevelReadCommitted:
		isoLevel = pgx.ReadCommitted
	case sql.LevelRepeatableRead:
		isoLevel = pgx.RepeatableRead
	case sql.LevelSerializable:
		isoLevel = pgx.Serializable
	default:
		isoLevel = pgx.ReadCommitted
	}

	tx, err := db.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:       isoLevel,
		AccessMode:     pgx.ReadWrite,
		DeferrableMode: pgx.NotDeferrable,
	})
	if err != nil {
		return nil, err
	}
	return &pgTransaction{tx: tx}, nil
}

// Ping verifies the connection to the database.
func (db *PGDatabase) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}

// Close closes the database connection.
func (db *PGDatabase) Close() error {
	db.logger.Info("Closing PostgreSQL connection")
	db.pool.Close()
	if db.sqlDB != nil {
		return db.sqlDB.Close()
	}
	return nil
}

// Stats returns database statistics.
func (db *PGDatabase) Stats() Stats {
	stats := db.pool.Stat()
	return Stats{
		OpenConnections:   int(stats.TotalConns()),
		InUse:             int(stats.AcquiredConns()),
		Idle:              int(stats.IdleConns()),
		WaitCount:         stats.EmptyAcquireCount(),
		WaitDuration:      stats.AcquireDuration().String(),
		MaxIdleClosed:     stats.CanceledAcquireCount(),
		MaxLifetimeClosed: 0, // Not directly available in pgxpool
	}
}

// Type returns the database type.
func (db *PGDatabase) Type() Type {
	return TypePostgreSQL
}

// Underlying returns the underlying connection pool.
func (db *PGDatabase) Underlying() interface{} {
	return db.pool
}

// GetSQLDB returns the stdlib SQL database connection for migrations.
func (db *PGDatabase) GetSQLDB() *sql.DB {
	return db.sqlDB
}

// Adapter types for pgx compatibility

type pgRows struct {
	rows pgx.Rows
}

func (r *pgRows) Next() bool {
	return r.rows.Next()
}

func (r *pgRows) Scan(dest ...interface{}) error {
	return r.rows.Scan(dest...)
}

func (r *pgRows) Close() error {
	r.rows.Close()
	return nil
}

func (r *pgRows) Err() error {
	return r.rows.Err()
}

type pgRow struct {
	row pgx.Row
}

func (r *pgRow) Scan(dest ...interface{}) error {
	return r.row.Scan(dest...)
}

type pgResult struct {
	result pgconn.CommandTag
}

func (r *pgResult) LastInsertId() (int64, error) {
	return 0, fmt.Errorf("LastInsertId not supported in PostgreSQL")
}

func (r *pgResult) RowsAffected() (int64, error) {
	return r.result.RowsAffected(), nil
}

type pgTransaction struct {
	tx pgx.Tx
}

func (t *pgTransaction) Query(
	ctx context.Context,
	query string,
	args ...interface{},
) (Rows, error) {
	rows, err := t.tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &pgRows{rows: rows}, nil
}

func (t *pgTransaction) QueryRow(ctx context.Context, query string, args ...interface{}) Row {
	return &pgRow{row: t.tx.QueryRow(ctx, query, args...)}
}

func (t *pgTransaction) Exec(
	ctx context.Context,
	query string,
	args ...interface{},
) (Result, error) {
	result, err := t.tx.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &pgResult{result: result}, nil
}

func (t *pgTransaction) Commit() error {
	return t.tx.Commit(context.Background())
}

func (t *pgTransaction) Rollback() error {
	return t.tx.Rollback(context.Background())
}
