// Package database provides database migration functionality
package database

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log/slog"
	"strings"

	"ariga.io/atlas/sql/migrate"
	"ariga.io/atlas/sql/postgres"
	"ariga.io/atlas/sql/sqlite"

	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver
	_ "modernc.org/sqlite"             // sqlite driver
)

// MigrationRunner handles database migrations.
type MigrationRunner struct {
	db           *Database
	migrationsFS fs.FS // Required migrations filesystem
}

// NewMigrationRunner creates a new migration runner.
func NewMigrationRunner(db *Database, migrationsFS fs.FS) *MigrationRunner {
	return &MigrationRunner{
		db:           db,
		migrationsFS: migrationsFS,
	}
}

// Up applies all pending migrations.
func (m *MigrationRunner) Up() error {
	ctx := context.Background()

	// Create Atlas driver based on database type
	var driver migrate.Driver
	var err error

	switch m.db.dbType {
	case TypePostgreSQL:
		driver, err = postgres.Open(m.db.sqlDB)
		if err != nil {
			return fmt.Errorf("failed to create postgres driver: %w", err)
		}

	case TypeSQLite:
		driver, err = sqlite.Open(m.db.sqlDB)
		if err != nil {
			return fmt.Errorf("failed to create sqlite driver: %w", err)
		}

	default:
		return fmt.Errorf("unsupported database type: %s", m.db.dbType)
	}

	// Create migration directory from provided FS
	// The migrations are expected to be in a "migrations" subdirectory
	migrationsFS, err := fs.Sub(m.migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create sub filesystem: %w", err)
	}

	// Create memory-based migration directory
	dir := &migrate.MemDir{}

	// Read all migration files from embed.FS and add to MemDir
	entries, err := fs.ReadDir(migrationsFS, ".")
	if err != nil {
		return fmt.Errorf("failed to read migrations: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		content, err := fs.ReadFile(migrationsFS, entry.Name())
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", entry.Name(), err)
		}
		if err := dir.WriteFile(entry.Name(), content); err != nil {
			return fmt.Errorf("failed to write migration %s: %w", entry.Name(), err)
		}
	}

	// Generate checksum file (atlas.sum) for the migration directory
	// This is required for Atlas to validate the migration files
	files, err := dir.Files()
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// If no migration files exist, nothing to do
	if len(files) == 0 {
		slog.Info("No migrations to apply")
		return nil
	}

	hashFile, err := migrate.NewHashFile(files)
	if err != nil {
		return fmt.Errorf("failed to create hash file: %w", err)
	}
	if err := migrate.WriteSumFile(dir, hashFile); err != nil {
		return fmt.Errorf("failed to write sum file: %w", err)
	}

	// Create executor with NopRevisionReadWriter for one-time migration replay
	// This is suitable for applying migrations without tracking revision history
	executor, err := migrate.NewExecutor(driver, dir, migrate.NopRevisionReadWriter{})
	if err != nil {
		return fmt.Errorf("failed to create executor: %w", err)
	}

	// Execute all pending migrations
	if err := executor.ExecuteN(ctx, 0); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

// Down is not supported in Atlas native migrations.
// Atlas uses a "roll-forward" approach instead of rollback.
func (m *MigrationRunner) Down() error {
	return fmt.Errorf(
		"down migrations not supported with Atlas - use roll-forward migrations instead",
	)
}

// Version returns the current migration version from Atlas's revision table.
func (m *MigrationRunner) Version() (int64, error) {
	var version string
	err := m.db.sqlDB.QueryRow(
		"SELECT version FROM atlas_schema_revisions ORDER BY executed_at DESC LIMIT 1",
	).Scan(&version)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // No migrations applied yet
		}
		return 0, fmt.Errorf("failed to get version: %w", err)
	}

	// Atlas versions are strings (timestamps), we'll return 1 if any migrations exist
	if version != "" {
		return 1, nil
	}
	return 0, nil
}
