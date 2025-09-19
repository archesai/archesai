// Package database provides database type detection and sqlc query generation.
//
// The package includes:
// - Database type constants for PostgreSQL and SQLite
// - Type detection from connection strings
// - Type-safe query generation using sqlc
package database

// Generate database queries from SQL files
//go:generate go tool sqlc generate

import (
	"database/sql"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Type represents the database type.
type Type string

// Database type constants.
const (
	TypePostgreSQL Type = "postgresql" // PostgreSQL database
	TypeSQLite     Type = "sqlite"     // SQLite database
)

// String converts Type to string.
func (t Type) String() string {
	switch t {
	case TypePostgreSQL:
		return "postgresql"
	case TypeSQLite:
		return "sqlite"
	default:
		return "unknown"
	}
}

// ParseTypeFromString converts a string to database Type.
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

// DetectTypeFromURL auto-detects database type from connection URL.
func DetectTypeFromURL(url string) Type {
	if strings.HasPrefix(url, "postgresql://") || strings.HasPrefix(url, "postgres://") {
		return TypePostgreSQL
	}
	if strings.HasPrefix(url, "sqlite://") || strings.Contains(url, ".db") || url == ":memory:" {
		return TypeSQLite
	}
	return TypePostgreSQL // Default to PostgreSQL
}

// Database wraps database connections for both PostgreSQL and SQLite.
type Database struct {
	sqlDB   *sql.DB       // Standard SQL connection (used by both)
	pgxPool *pgxpool.Pool // PostgreSQL specific pool (nil for SQLite)
	dbType  Type          // Database type
}

// NewDatabase creates a new database wrapper.
func NewDatabase(sqlDB *sql.DB, pgxPool *pgxpool.Pool, dbType Type) *Database {
	return &Database{
		sqlDB:   sqlDB,
		pgxPool: pgxPool,
		dbType:  dbType,
	}
}

// SQLDB returns the standard SQL database connection.
func (d *Database) SQLDB() *sql.DB {
	return d.sqlDB
}

// PgxPool returns the PostgreSQL connection pool (nil for SQLite).
func (d *Database) PgxPool() *pgxpool.Pool {
	return d.pgxPool
}

// Type returns the database type.
func (d *Database) Type() Type {
	return d.dbType
}

// TypeString returns the database type as a string.
func (d *Database) TypeString() string {
	return d.dbType.String()
}

// IsPostgreSQL returns true if this is a PostgreSQL database.
func (d *Database) IsPostgreSQL() bool {
	return d.dbType == TypePostgreSQL
}

// IsSQLite returns true if this is a SQLite database.
func (d *Database) IsSQLite() bool {
	return d.dbType == TypeSQLite
}
