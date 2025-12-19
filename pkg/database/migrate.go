package database

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"log/slog"
	"strings"
	"time"

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
		slog.Debug("No migrations to apply")
		return nil
	}

	hashFile, err := migrate.NewHashFile(files)
	if err != nil {
		return fmt.Errorf("failed to create hash file: %w", err)
	}
	if err := migrate.WriteSumFile(dir, hashFile); err != nil {
		return fmt.Errorf("failed to write sum file: %w", err)
	}

	// Create revision tracker for tracking applied migrations
	revisions := newDBRevisions(m.db.sqlDB, m.db.dbType)
	if err := revisions.init(ctx); err != nil {
		return fmt.Errorf("failed to initialize revision table: %w", err)
	}

	// Create executor with revision tracking
	executor, err := migrate.NewExecutor(driver, dir, revisions,
		migrate.WithAllowDirty(true),
	)
	if err != nil {
		return fmt.Errorf("failed to create executor: %w", err)
	}

	// Check for pending migrations
	pending, err := executor.Pending(ctx)
	if err != nil {
		// "no pending migration files" means all migrations are applied
		if strings.Contains(err.Error(), "no pending migration files") {
			slog.Debug("no pending migrations")
			return nil
		}
		return fmt.Errorf("failed to check pending migrations: %w", err)
	}

	if len(pending) == 0 {
		slog.Debug("no pending migrations")
		return nil
	}

	slog.Debug("applying migrations", "count", len(pending))

	// Execute only pending migrations
	if err := executor.ExecuteN(ctx, len(pending)); err != nil {
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

// Version returns the number of applied migrations.
func (m *MigrationRunner) Version() (int64, error) {
	ctx := context.Background()

	revisions := newDBRevisions(m.db.sqlDB, m.db.dbType)
	applied, err := revisions.ReadRevisions(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to read revisions: %w", err)
	}

	return int64(len(applied)), nil
}

// dbRevisions implements migrate.RevisionReadWriter for tracking migrations in a database table.
type dbRevisions struct {
	db     *sql.DB
	dbType Type
}

func newDBRevisions(db *sql.DB, dbType Type) *dbRevisions {
	return &dbRevisions{db: db, dbType: dbType}
}

// Ident returns the table identifier.
func (r *dbRevisions) Ident() *migrate.TableIdent {
	return &migrate.TableIdent{Name: "atlas_schema_revisions", Schema: ""}
}

// init creates the revisions table if it doesn't exist.
func (r *dbRevisions) init(ctx context.Context) error {
	var query string
	switch r.dbType {
	case TypePostgreSQL:
		query = `CREATE TABLE IF NOT EXISTS atlas_schema_revisions (
			version VARCHAR(255) PRIMARY KEY,
			description VARCHAR(255) NOT NULL,
			type INTEGER NOT NULL DEFAULT 2,
			applied INTEGER NOT NULL DEFAULT 0,
			total INTEGER NOT NULL DEFAULT 0,
			executed_at TIMESTAMP NOT NULL,
			execution_time BIGINT NOT NULL,
			error TEXT,
			error_stmt TEXT,
			hash VARCHAR(255) NOT NULL,
			partial_hashes TEXT,
			operator_version VARCHAR(255) NOT NULL
		)`
	case TypeSQLite:
		query = `CREATE TABLE IF NOT EXISTS atlas_schema_revisions (
			version TEXT PRIMARY KEY,
			description TEXT NOT NULL,
			type INTEGER NOT NULL DEFAULT 2,
			applied INTEGER NOT NULL DEFAULT 0,
			total INTEGER NOT NULL DEFAULT 0,
			executed_at TEXT NOT NULL,
			execution_time INTEGER NOT NULL,
			error TEXT,
			error_stmt TEXT,
			hash TEXT NOT NULL,
			partial_hashes TEXT,
			operator_version TEXT NOT NULL
		)`
	default:
		return fmt.Errorf("unsupported database type: %s", r.dbType)
	}
	_, err := r.db.ExecContext(ctx, query)
	return err
}

// ReadRevisions returns all revisions.
func (r *dbRevisions) ReadRevisions(ctx context.Context) ([]*migrate.Revision, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT version, description, type, applied, total, executed_at, execution_time, error, error_stmt, hash, partial_hashes, operator_version
		FROM atlas_schema_revisions ORDER BY version`,
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var revisions []*migrate.Revision
	for rows.Next() {
		var rev migrate.Revision
		var executedAt time.Time
		var errorStr, errorStmt, partialHashes sql.NullString
		if err := rows.Scan(
			&rev.Version, &rev.Description, &rev.Type, &rev.Applied, &rev.Total,
			&executedAt, &rev.ExecutionTime, &errorStr, &errorStmt, &rev.Hash,
			&partialHashes, &rev.OperatorVersion,
		); err != nil {
			return nil, err
		}
		rev.ExecutedAt = executedAt
		if errorStr.Valid {
			rev.Error = errorStr.String
		}
		if errorStmt.Valid {
			rev.ErrorStmt = errorStmt.String
		}
		if partialHashes.Valid {
			rev.PartialHashes = strings.Split(partialHashes.String, ",")
		}
		revisions = append(revisions, &rev)
	}
	return revisions, rows.Err()
}

// ReadRevision returns a revision by version.
func (r *dbRevisions) ReadRevision(ctx context.Context, version string) (*migrate.Revision, error) {
	var rev migrate.Revision
	var executedAt time.Time
	var errorStr, errorStmt, partialHashes sql.NullString
	err := r.db.QueryRowContext(
		ctx,
		`SELECT version, description, type, applied, total, executed_at, execution_time, error, error_stmt, hash, partial_hashes, operator_version
		FROM atlas_schema_revisions WHERE version = $1`,
		version,
	).Scan(
		&rev.Version, &rev.Description, &rev.Type, &rev.Applied, &rev.Total,
		&executedAt, &rev.ExecutionTime, &errorStr, &errorStmt, &rev.Hash,
		&partialHashes, &rev.OperatorVersion,
	)
	if err == sql.ErrNoRows {
		return nil, migrate.ErrRevisionNotExist
	}
	if err != nil {
		return nil, err
	}
	rev.ExecutedAt = executedAt
	if errorStr.Valid {
		rev.Error = errorStr.String
	}
	if errorStmt.Valid {
		rev.ErrorStmt = errorStmt.String
	}
	if partialHashes.Valid {
		rev.PartialHashes = strings.Split(partialHashes.String, ",")
	}
	return &rev, nil
}

// WriteRevision saves the revision to the storage.
func (r *dbRevisions) WriteRevision(ctx context.Context, rev *migrate.Revision) error {
	var partialHashes string
	if len(rev.PartialHashes) > 0 {
		partialHashes = strings.Join(rev.PartialHashes, ",")
	}

	var query string
	switch r.dbType {
	case TypePostgreSQL:
		query = `INSERT INTO atlas_schema_revisions
			(version, description, type, applied, total, executed_at, execution_time, error, error_stmt, hash, partial_hashes, operator_version)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
			ON CONFLICT (version) DO UPDATE SET
			description = EXCLUDED.description, type = EXCLUDED.type, applied = EXCLUDED.applied,
			total = EXCLUDED.total, executed_at = EXCLUDED.executed_at, execution_time = EXCLUDED.execution_time,
			error = EXCLUDED.error, error_stmt = EXCLUDED.error_stmt, hash = EXCLUDED.hash,
			partial_hashes = EXCLUDED.partial_hashes, operator_version = EXCLUDED.operator_version`
	case TypeSQLite:
		query = `INSERT OR REPLACE INTO atlas_schema_revisions
			(version, description, type, applied, total, executed_at, execution_time, error, error_stmt, hash, partial_hashes, operator_version)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	default:
		return fmt.Errorf("unsupported database type: %s", r.dbType)
	}

	_, err := r.db.ExecContext(ctx, query,
		rev.Version, rev.Description, rev.Type, rev.Applied, rev.Total,
		rev.ExecutedAt, rev.ExecutionTime, rev.Error, rev.ErrorStmt, rev.Hash,
		partialHashes, rev.OperatorVersion,
	)
	return err
}

// DeleteRevision deletes a revision by version from the storage.
func (r *dbRevisions) DeleteRevision(ctx context.Context, version string) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM atlas_schema_revisions WHERE version = $1`,
		version,
	)
	return err
}
