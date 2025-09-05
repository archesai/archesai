// Package database provides data persistence infrastructure including
// database connections, query generation, and migrations.
//
// The package includes:
// - Database abstraction layer supporting PostgreSQL and SQLite
// - Type-safe query generation using sqlc
// - Database migrations using golang-migrate
// - Connection pooling and health checks
// - Transaction management
package database

// Generate database queries from SQL files
//go:generate sqlc generate

import (
	"context"
	"database/sql"
	"strings"
)

// Type represents the database type
type Type string

// Database type constants
const (
	TypePostgreSQL Type = "postgresql" // PostgreSQL database
	TypeSQLite     Type = "sqlite"     // SQLite database
)

// Database defines the common interface for all database implementations
type Database interface {
	// Core operations
	Query(ctx context.Context, query string, args ...interface{}) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) Row
	Exec(ctx context.Context, query string, args ...interface{}) (Result, error)

	// Transaction support
	Begin(ctx context.Context) (Transaction, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Transaction, error)

	// Connection management
	Ping(ctx context.Context) error
	Close() error
	Stats() Stats

	// Database type identification
	Type() Type

	// Get underlying connection for driver-specific operations
	// Returns *sql.DB for SQLite, *pgxpool.Pool for PostgreSQL
	Underlying() interface{}
}

// Transaction defines the interface for database transactions
type Transaction interface {
	Query(ctx context.Context, query string, args ...interface{}) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) Row
	Exec(ctx context.Context, query string, args ...interface{}) (Result, error)
	Commit() error
	Rollback() error
}

// Rows defines the interface for query result rows
type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
	Err() error
}

// Row defines the interface for a single query result row
type Row interface {
	Scan(dest ...interface{}) error
}

// Result defines the interface for exec results
type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

// Stats represents database connection pool statistics
type Stats struct {
	OpenConnections   int
	InUse             int
	Idle              int
	WaitCount         int64
	WaitDuration      string
	MaxIdleClosed     int64
	MaxLifetimeClosed int64
}

// ParseTypeFromString converts a string to database Type
func ParseTypeFromString(s string) Type {
	switch s {
	case "postgresql", "postgres", "pg":
		return TypePostgreSQL
	case "sqlite", "sqlite3":
		return TypeSQLite
	default:
		return TypePostgreSQL // Default to PostgreSQL
	}
}

// DetectTypeFromURL auto-detects database type from connection URL
func DetectTypeFromURL(url string) Type {
	if strings.HasPrefix(url, "postgresql://") || strings.HasPrefix(url, "postgres://") {
		return TypePostgreSQL
	}
	if strings.HasPrefix(url, "sqlite://") || strings.Contains(url, ".db") || url == ":memory:" {
		return TypeSQLite
	}
	return TypePostgreSQL // Default to PostgreSQL
}
