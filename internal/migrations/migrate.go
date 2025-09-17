// Package migrations provides database migration functionality
package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"os"

	"github.com/pressly/goose/v3"

	"github.com/archesai/archesai/internal/database"
)

//go:embed postgresql/*.sql
var migrations embed.FS

// MigrationRunner handles database migrations.
type MigrationRunner struct {
	db     database.Database
	logger *slog.Logger
}

// NewMigrationRunner creates a new migration runner.
func NewMigrationRunner(db database.Database, logger *slog.Logger) *MigrationRunner {
	return &MigrationRunner{
		db:     db,
		logger: logger,
	}
}

// Up applies all pending migrations.
func (m *MigrationRunner) Up() error {
	if err := m.setEnvironmentVariables(); err != nil {
		return fmt.Errorf("failed to set environment variables: %w", err)
	}

	sqlDB, err := m.getSQLDB()
	if err != nil {
		return fmt.Errorf("failed to get SQL database: %w", err)
	}

	goose.SetBaseFS(migrations)

	// Set the dialect based on database type
	dialect := m.getDialect()
	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(sqlDB, "migrations"); err != nil {
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

	sqlDB, err := m.getSQLDB()
	if err != nil {
		return fmt.Errorf("failed to get SQL database: %w", err)
	}

	goose.SetBaseFS(migrations)

	dialect := m.getDialect()
	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Down(sqlDB, "migrations"); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	m.logger.Info("Migration rolled back successfully")
	return nil
}

// Version returns the current migration version.
func (m *MigrationRunner) Version() (int64, error) {
	sqlDB, err := m.getSQLDB()
	if err != nil {
		return 0, fmt.Errorf("failed to get SQL database: %w", err)
	}

	goose.SetBaseFS(migrations)

	dialect := m.getDialect()
	if err := goose.SetDialect(dialect); err != nil {
		return 0, fmt.Errorf("failed to set dialect: %w", err)
	}

	version, err := goose.GetDBVersion(sqlDB)
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

	sqlDB, err := m.getSQLDB()
	if err != nil {
		return fmt.Errorf("failed to get SQL database: %w", err)
	}

	goose.SetBaseFS(migrations)

	dialect := m.getDialect()
	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	// Force the version by resetting and then applying up to the target version
	if err := goose.Reset(sqlDB, "migrations"); err != nil {
		return fmt.Errorf("failed to reset migrations: %w", err)
	}

	if version > 0 {
		if err := goose.UpTo(sqlDB, "migrations", version); err != nil {
			return fmt.Errorf("failed to migrate to version %d: %w", version, err)
		}
	}

	m.logger.Info("Migration version forced", "version", version)
	return nil
}

// setEnvironmentVariables sets database-specific environment variables for migrations.
func (m *MigrationRunner) setEnvironmentVariables() error {
	switch m.db.Type() {
	case database.TypePostgreSQL:
		_ = os.Setenv("TIMESTAMP_TYPE", "TIMESTAMPTZ")
		_ = os.Setenv("TIMESTAMP_DEFAULT", "CURRENT_TIMESTAMP")
		_ = os.Setenv("REAL_TYPE", "DOUBLE PRECISION")
	case database.TypeSQLite:
		_ = os.Setenv("TIMESTAMP_TYPE", "TEXT")
		_ = os.Setenv("TIMESTAMP_DEFAULT", "(strftime('%Y-%m-%d %H:%M:%f', 'now'))")
		_ = os.Setenv("REAL_TYPE", "REAL")
	default:
		return fmt.Errorf("unsupported database type: %s", m.db.Type())
	}
	return nil
}

// getDialect returns the goose dialect string for the database type.
func (m *MigrationRunner) getDialect() string {
	switch m.db.Type() {
	case database.TypePostgreSQL:
		return "postgres"
	case database.TypeSQLite:
		return "sqlite3"
	default:
		return ""
	}
}

// getSQLDB gets a database/sql connection.
func (m *MigrationRunner) getSQLDB() (*sql.DB, error) {
	switch m.db.Type() {
	case database.TypePostgreSQL:
		cfg, ok := m.db.(*database.PGDatabase)
		if !ok {
			return nil, fmt.Errorf("database is not PostgreSQL")
		}
		return cfg.GetSQLDB(), nil
	case database.TypeSQLite:
		sqlDB, ok := m.db.Underlying().(*sql.DB)
		if !ok {
			return nil, fmt.Errorf("SQLite database does not have *sql.DB underlying connection")
		}
		return sqlDB, nil
	default:
		return nil, fmt.Errorf("unsupported database type: %s", m.db.Type())
	}
}

// RunMigrations is a convenience function to run migrations on a database.
func RunMigrations(db database.Database, logger *slog.Logger) error {
	runner := NewMigrationRunner(db, logger)
	return runner.Up()
}

// MigrationConfig holds configuration for migrations.
type MigrationConfig struct {
	Direction string // "up" or "down"
	Steps     int    // Number of steps (0 means all)
	Force     int64  // Force to a specific version (-1 means don't force)
}

// RunMigrationsWithConfig runs migrations with specific configuration.
func RunMigrationsWithConfig(db database.Database, cfg MigrationConfig, logger *slog.Logger) error {
	runner := NewMigrationRunner(db, logger)

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
