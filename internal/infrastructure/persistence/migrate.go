// Package migrations provides database migration functionality
package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/pressly/goose/v3"
)

//go:embed postgres/migrations/*.sql
var migrations embed.FS

//go:embed sqlite/migrations/*.sql
var _ embed.FS

const (
	dbTypePostgreSQL = "postgresql"
	dbTypePostgres   = "postgres"
	dbTypeSQLite     = "sqlite"
	dbTypeSQLite3    = "sqlite3"
)

// MigrationRunner handles database migrations.
type MigrationRunner struct {
	db     *sql.DB
	dbType string
	logger *slog.Logger
}

// NewMigrationRunner creates a new migration runner.
func NewMigrationRunner(db *sql.DB, dbType string, logger *slog.Logger) *MigrationRunner {
	return &MigrationRunner{
		db:     db,
		dbType: dbType,
		logger: logger,
	}
}

// Up applies all pending migrations.
func (m *MigrationRunner) Up() error {
	if err := m.setEnvironmentVariables(); err != nil {
		return fmt.Errorf("failed to set environment variables: %w", err)
	}

	goose.SetBaseFS(migrations)

	// Set the dialect based on database type
	dialect := m.getDialect()
	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(m.db, "migrations"); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	m.logger.Info("Migrations applied successfully")
	return nil
}

// Down rolls back the last migration.
func (m *MigrationRunner) Down() error {
	if err := m.setEnvironmentVariables(); err != nil {
		return fmt.Errorf("failed to set environment variables: %w", err)
	}

	goose.SetBaseFS(migrations)

	dialect := m.getDialect()
	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Down(m.db, "migrations"); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	m.logger.Info("Migration rolled back successfully")
	return nil
}

// Version returns the current migration version.
func (m *MigrationRunner) Version() (int64, error) {
	goose.SetBaseFS(migrations)

	dialect := m.getDialect()
	if err := goose.SetDialect(dialect); err != nil {
		return 0, fmt.Errorf("failed to set dialect: %w", err)
	}

	version, err := goose.GetDBVersion(m.db)
	if err != nil {
		return 0, fmt.Errorf("failed to get version: %w", err)
	}

	return version, nil
}

// Force sets the migration version without running migrations.
func (m *MigrationRunner) Force(version int64) error {
	if err := m.setEnvironmentVariables(); err != nil {
		return fmt.Errorf("failed to set environment variables: %w", err)
	}

	goose.SetBaseFS(migrations)

	dialect := m.getDialect()
	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	// Force the version by resetting and then applying up to the target version
	if err := goose.Reset(m.db, "migrations"); err != nil {
		return fmt.Errorf("failed to reset migrations: %w", err)
	}

	if version > 0 {
		if err := goose.UpTo(m.db, "migrations", version); err != nil {
			return fmt.Errorf("failed to migrate to version %d: %w", version, err)
		}
	}

	m.logger.Info("Migration version forced", "version", version)
	return nil
}

// setEnvironmentVariables sets database-specific environment variables for migrations.
func (m *MigrationRunner) setEnvironmentVariables() error {
	switch m.dbType {
	case dbTypePostgreSQL, dbTypePostgres:
		_ = os.Setenv("TIMESTAMP_TYPE", "TIMESTAMPTZ")
		_ = os.Setenv("TIMESTAMP_DEFAULT", "CURRENT_TIMESTAMP")
		_ = os.Setenv("REAL_TYPE", "DOUBLE PRECISION")
	case dbTypeSQLite, dbTypeSQLite3:
		_ = os.Setenv("TIMESTAMP_TYPE", "TEXT")
		_ = os.Setenv("TIMESTAMP_DEFAULT", "(strftime('%Y-%m-%d %H:%M:%f', 'now'))")
		_ = os.Setenv("REAL_TYPE", "REAL")
	default:
		return fmt.Errorf("unsupported database type: %s", m.dbType)
	}
	return nil
}

// getDialect returns the goose dialect string for the database type.
func (m *MigrationRunner) getDialect() string {
	switch m.dbType {
	case dbTypePostgreSQL, dbTypePostgres:
		return dbTypePostgres
	case dbTypeSQLite, dbTypeSQLite3:
		return dbTypeSQLite3
	default:
		return ""
	}
}

// RunMigrations is a convenience function to run migrations on a database.
// It attempts to detect the database type from the driver name.
func RunMigrations(db *sql.DB, dbType string, logger *slog.Logger) error {
	// If dbType is not provided, try to detect it
	if dbType == "" {
		dbType = detectDatabaseType(db)
	}
	runner := NewMigrationRunner(db, dbType, logger)
	return runner.Up()
}

// detectDatabaseType attempts to detect the database type from the connection.
func detectDatabaseType(db *sql.DB) string {
	// Try to get driver name from the database stats
	driverName := db.Driver()
	if driverName != nil {
		driverType := fmt.Sprintf("%T", driverName)
		if strings.Contains(strings.ToLower(driverType), "postgres") ||
			strings.Contains(strings.ToLower(driverType), "pgx") {
			return dbTypePostgreSQL
		}
		if strings.Contains(strings.ToLower(driverType), dbTypeSQLite) {
			return dbTypeSQLite
		}
	}

	// Try a simple query to detect the database type
	var version string

	// Try PostgreSQL version query
	err := db.QueryRow("SELECT version()").Scan(&version)
	if err == nil && strings.Contains(strings.ToLower(version), dbTypePostgreSQL) {
		return dbTypePostgreSQL
	}

	// Try SQLite version query
	err = db.QueryRow("SELECT sqlite_version()").Scan(&version)
	if err == nil {
		return dbTypeSQLite
	}

	// Default to PostgreSQL if we can't detect
	return dbTypePostgreSQL
}

// MigrationConfig holds configuration for migrations.
type MigrationConfig struct {
	Direction string // "up" or "down"
	Steps     int    // Number of steps (0 means all)
	Force     int64  // Force to a specific version (-1 means don't force)
}

// RunMigrationsWithConfig runs migrations with specific configuration.
func RunMigrationsWithConfig(
	db *sql.DB,
	dbType string,
	cfg MigrationConfig,
	logger *slog.Logger,
) error {
	runner := NewMigrationRunner(db, dbType, logger)

	if cfg.Force >= 0 {
		return runner.Force(cfg.Force)
	}

	switch cfg.Direction {
	case "up":
		return runner.Up()
	case "down":
		return runner.Down()
	default:
		return fmt.Errorf("invalid migration direction: %s", cfg.Direction)
	}
}
