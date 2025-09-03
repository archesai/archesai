// Package database provides database abstraction and connection management.
package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/postgresql/*.sql
var postgresqlMigrations embed.FS

//go:embed migrations/sqlite/*.sql
var sqliteMigrations embed.FS

// MigrationRunner handles database migrations
type MigrationRunner struct {
	db     Database
	logger *slog.Logger
}

// NewMigrationRunner creates a new migration runner
func NewMigrationRunner(db Database, logger *slog.Logger) *MigrationRunner {
	return &MigrationRunner{
		db:     db,
		logger: logger,
	}
}

// Up applies all pending migrations
func (m *MigrationRunner) Up() error {
	migration, err := m.createMigration()
	if err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}
	defer func() {
		if _, err := migration.Close(); err != nil {
			m.logger.Debug("failed to close migration", "error", err)
		}
	}()

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		m.logger.Info("No migrations to apply")
	} else {
		m.logger.Info("Migrations applied successfully")
	}

	return nil
}

// Down rolls back the last migration
func (m *MigrationRunner) Down() error {
	migration, err := m.createMigration()
	if err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}
	defer func() {
		if _, err := migration.Close(); err != nil {
			m.logger.Debug("failed to close migration", "error", err)
		}
	}()

	if err := migration.Steps(-1); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	if err == migrate.ErrNoChange {
		m.logger.Info("No migrations to rollback")
	} else {
		m.logger.Info("Migration rolled back successfully")
	}

	return nil
}

// Version returns the current migration version
func (m *MigrationRunner) Version() (uint, bool, error) {
	migration, err := m.createMigration()
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migration: %w", err)
	}
	defer func() {
		if _, err := migration.Close(); err != nil {
			m.logger.Debug("failed to close migration", "error", err)
		}
	}()

	return migration.Version()
}

// Force sets the migration version without running migrations
func (m *MigrationRunner) Force(version int) error {
	migration, err := m.createMigration()
	if err != nil {
		return fmt.Errorf("failed to create migration: %w", err)
	}
	defer func() {
		if _, err := migration.Close(); err != nil {
			m.logger.Debug("failed to close migration", "error", err)
		}
	}()

	if err := migration.Force(version); err != nil {
		return fmt.Errorf("failed to force migration version: %w", err)
	}

	m.logger.Info("Migration version forced", "version", version)
	return nil
}

// createMigration creates a migrate instance based on the database type
func (m *MigrationRunner) createMigration() (*migrate.Migrate, error) {
	var driver database.Driver
	var sourceDriver source.Driver
	var err error

	switch m.db.Type() {
	case TypePostgreSQL:
		// Get the underlying connection for PostgreSQL
		// We need to create a database/sql connection for golang-migrate
		sqlDB, err := m.getPostgreSQLConnection()
		if err != nil {
			return nil, fmt.Errorf("failed to get PostgreSQL connection: %w", err)
		}

		driver, err = postgres.WithInstance(sqlDB, &postgres.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to create PostgreSQL driver: %w", err)
		}

		sourceDriver, err = iofs.New(postgresqlMigrations, "migrations/postgresql")
		if err != nil {
			return nil, fmt.Errorf("failed to create PostgreSQL source driver: %w", err)
		}

	case TypeSQLite:
		// Get the underlying connection for SQLite
		sqlDB, ok := m.db.Underlying().(*sql.DB)
		if !ok {
			return nil, fmt.Errorf("SQLite database does not have *sql.DB underlying connection")
		}

		driver, err = sqlite.WithInstance(sqlDB, &sqlite.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to create SQLite driver: %w", err)
		}

		sourceDriver, err = iofs.New(sqliteMigrations, "migrations/sqlite")
		if err != nil {
			return nil, fmt.Errorf("failed to create SQLite source driver: %w", err)
		}

	default:
		return nil, fmt.Errorf("unsupported database type: %s", m.db.Type())
	}

	migration, err := migrate.NewWithInstance("iofs", sourceDriver, m.db.Type().String(), driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}

	return migration, nil
}

// getPostgreSQLConnection gets a database/sql connection for PostgreSQL
func (m *MigrationRunner) getPostgreSQLConnection() (*sql.DB, error) {
	// For PostgreSQL, we need to extract the config and create a sql.DB connection
	// since golang-migrate doesn't support pgxpool directly
	cfg, ok := m.db.(*PostgreSQL)
	if !ok {
		return nil, fmt.Errorf("database is not PostgreSQL")
	}

	// Return the stdlib connection
	return cfg.sqlDB, nil
}

// String converts Type to string
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

// RunMigrations is a convenience function to run migrations on a database
func RunMigrations(db Database, logger *slog.Logger) error {
	runner := NewMigrationRunner(db, logger)
	return runner.Up()
}

// MigrationConfig holds configuration for migrations
type MigrationConfig struct {
	MigrationsPath string
	Direction      string // "up" or "down"
	Steps          int    // Number of steps (0 means all)
	Force          int    // Force to a specific version (-1 means don't force)
}

// RunMigrationsWithConfig runs migrations with specific configuration
func RunMigrationsWithConfig(db Database, cfg MigrationConfig, logger *slog.Logger) error {
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

// GetMigrationFiles returns the list of migration files for the database type
func GetMigrationFiles(dbType Type) ([]string, error) {
	var files []string
	var fsys embed.FS
	var basePath string

	switch dbType {
	case TypePostgreSQL:
		fsys = postgresqlMigrations
		basePath = "migrations/postgresql"
	case TypeSQLite:
		fsys = sqliteMigrations
		basePath = "migrations/sqlite"
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	entries, err := fsys.ReadDir(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read migration directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".sql" {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}
