// Package database provides database abstraction and connection management.
package database

import (
	"context"
	"database/sql"
	"fmt"
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

// Config holds database configuration
type Config struct {
	URL                 string
	Type                Type
	PgMaxConns          int32
	PgMinConns          int32
	ConnMaxLifetime     string
	ConnMaxIdleTime     string
	PgHealthCheckPeriod string
	RunMigrations       bool
}

// DefaultConfig returns default configuration for the specified database type
func DefaultConfig(dbType Type) *Config {
	switch dbType {
	case TypePostgreSQL:
		return &Config{
			Type:                TypePostgreSQL,
			PgMaxConns:          25,
			PgMinConns:          5,
			ConnMaxLifetime:     "1h",
			ConnMaxIdleTime:     "30m",
			PgHealthCheckPeriod: "1m",
		}
	case TypeSQLite:
		return &Config{
			Type:            TypeSQLite,
			ConnMaxLifetime: "0",
			ConnMaxIdleTime: "0",
		}
	default:
		return nil
	}
}

// ParseType parses a string into a database Type
func ParseType(s string) (Type, error) {
	switch s {
	case "postgresql", "postgres", "pg":
		return TypePostgreSQL, nil
	case "sqlite", "sqlite3":
		return TypeSQLite, nil
	default:
		return "", fmt.Errorf("unknown database type: %s", s)
	}
}

// ParseTypeString converts a string to Type (returns TypePostgreSQL as default)
func ParseTypeString(s string) Type {
	t, err := ParseType(s)
	if err != nil {
		return TypePostgreSQL
	}
	return t
}
