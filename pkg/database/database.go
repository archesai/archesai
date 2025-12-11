// Package database provides database type detection and sqlc query generation.
//
// The package includes:
// - Database type constants for PostgreSQL and SQLite
// - Type detection from connection strings
// - Type-safe query generation using sqlc
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
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

// StartPostgreSQL starts a PostgreSQL testcontainer
func StartPostgreSQL(ctx context.Context) (*Database, *testpostgres.PostgresContainer, error) {

	// Create PostgreSQL container with pgvector extension
	postgresContainer, err := testpostgres.Run(ctx,
		"docker.io/pgvector/pgvector:pg15",
		testpostgres.WithDatabase("archesai-migrations"),
		testpostgres.WithUsername("postgres"),
		testpostgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	// Get connection string
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	// Open database connection using pgx driver
	sqlDB, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Create database.Database instance
	database := NewDatabase(sqlDB, nil, TypePostgreSQL)

	slog.Debug("PostgreSQL testcontainer started", slog.String("url", connStr))

	return database, postgresContainer, nil
}

// StartSQLite starts an in-memory SQLite database
func StartSQLite() (*Database, error) {

	sqlDB, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %w", err)
	}

	db := NewDatabase(sqlDB, nil, TypeSQLite)

	return db, nil
}

// Open opens a database connection based on the provided type and URL.
// For SQLite, the URL should be a file path or ":memory:".
// For PostgreSQL, the URL should be a connection string.
func Open(dbType, url string) (*Database, error) {
	t := ParseTypeFromString(dbType)

	switch t {
	case TypeSQLite:
		sqlDB, err := sql.Open("sqlite", url)
		if err != nil {
			return nil, fmt.Errorf("failed to open SQLite database: %w", err)
		}
		return NewDatabase(sqlDB, nil, TypeSQLite), nil

	case TypePostgreSQL:
		// Create pgxpool for PostgreSQL
		pool, err := pgxpool.New(context.Background(), url)
		if err != nil {
			return nil, fmt.Errorf("failed to create PostgreSQL pool: %w", err)
		}
		// Also create standard sql.DB for compatibility
		sqlDB, err := sql.Open("pgx", url)
		if err != nil {
			pool.Close()
			return nil, fmt.Errorf("failed to open PostgreSQL database: %w", err)
		}
		return NewDatabase(sqlDB, pool, TypePostgreSQL), nil

	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}
}
