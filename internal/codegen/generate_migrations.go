// Package codegen provides code generation utilities including database migration generation using Atlas.
package codegen

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"ariga.io/atlas/sql/migrate"
	"ariga.io/atlas/sql/postgres"
	"ariga.io/atlas/sql/schema"
	"ariga.io/atlas/sql/sqlite"

	// pgx driver is required for postgres connections
	_ "github.com/jackc/pgx/v5/stdlib"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	// modernc.org/sqlite driver is required for SQLite connections
	_ "modernc.org/sqlite"

	database "github.com/archesai/archesai/internal/infrastructure/persistence"
)

const (
	// PostgresMigrationDir is the default directory for migrations
	PostgresMigrationDir = "internal/infrastructure/persistence/postgres/migrations"

	// SQLiteMigrationDir is the default directory for SQLite migrations
	SQLiteMigrationDir = "internal/infrastructure/persistence/sqlite/migrations"

	// PostgresHCLSchemaFile is the generated HCL schema file for PostgreSQL
	PostgresHCLSchemaFile = "internal/infrastructure/persistence/postgres/schema.gen.hcl"

	// SQLiteHCLSchemaFile is the generated HCL schema file for SQLite
	SQLiteHCLSchemaFile = "internal/infrastructure/persistence/sqlite/schema.gen.hcl"

	// PostgreSQL schema name
	postgresSchemaName = "public"

	// SQLite schema name
	sqliteSchemaName = "main"
)

// GetHCLSchemaFile returns the path to the HCL schema file for a given database type
func GetHCLSchemaFile(dbType database.Type) string {
	switch dbType {
	case database.TypePostgreSQL:
		return PostgresHCLSchemaFile
	case database.TypeSQLite:
		return SQLiteHCLSchemaFile
	default:
		return ""
	}
}

// GenerateMigrations generates a migration for a specific database type
func (g *Generator) GenerateMigrations(ctx context.Context, dbType database.Type) error {
	m := &MigrationGenerator{}

	// Start database
	if err := m.Start(ctx, dbType); err != nil {
		return err
	}
	defer func() {
		if err := m.Stop(ctx); err != nil {
			slog.Error(
				"Failed to stop database",
				slog.String("error", err.Error()),
				slog.String("type", dbType.String()),
			)
		}
	}()

	// Generate migration
	return m.GenerateMigration(ctx)
}

// MigrationGenerator handles database migration generation using Atlas
type MigrationGenerator struct {
	container *testpostgres.PostgresContainer
	database  *database.Database
}

// Start spins up a database for migration generation (testcontainer for PostgreSQL, in-memory for SQLite)
func (m *MigrationGenerator) Start(ctx context.Context, dbType database.Type) error {
	switch dbType {
	case database.TypePostgreSQL:
		database, container, err := database.StartPostgreSQL(ctx)
		if err != nil {
			return fmt.Errorf("failed to start PostgreSQL container: %w", err)
		}
		m.database = database
		m.container = container
	case database.TypeSQLite:
		database, err := database.StartSQLite()
		if err != nil {
			return fmt.Errorf("failed to start SQLite database: %w", err)
		}
		m.database = database
	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}
	return nil
}

// Stop terminates the PostgreSQL testcontainer
func (m *MigrationGenerator) Stop(ctx context.Context) error {
	if m.database != nil && m.database.SQLDB() != nil {
		if err := m.database.SQLDB().Close(); err != nil {
			slog.Warn("Failed to close database connection", slog.String("error", err.Error()))
		}
	}
	if m.container != nil {
		slog.Debug("Stopping PostgreSQL testcontainer")
		return m.container.Terminate(ctx)
	}
	return nil
}

// GenerateMigration generates a migration by comparing the current database state
// with the desired state defined in the HCL schema
func (m *MigrationGenerator) GenerateMigration(ctx context.Context) error {
	if m.database == nil {
		return fmt.Errorf("database not initialized, call Start() first")
	}

	// Determine migration directory and schema name based on database type
	var migrationDir, schemaName string
	switch m.database.Type() {
	case database.TypePostgreSQL:
		migrationDir = PostgresMigrationDir
		schemaName = postgresSchemaName
	case database.TypeSQLite:
		migrationDir = SQLiteMigrationDir
		schemaName = sqliteSchemaName
	default:
		return fmt.Errorf("unsupported database type: %s", m.database.Type())
	}

	// Ensure migration directory exists
	if err := os.MkdirAll(migrationDir, 0755); err != nil {
		return fmt.Errorf("failed to create migration directory: %w", err)
	}

	// Step 1: Apply existing migrations to get current state
	if err := database.RunMigrations(m.database); err != nil {
		return fmt.Errorf("failed to apply existing migrations: %w", err)
	}
	slog.Debug("Existing migrations applied successfully to dev database")

	// Create Atlas driver based on database type
	var driver migrate.Driver
	var err error
	switch m.database.Type() {
	case database.TypePostgreSQL:
		driver, err = postgres.Open(m.database.SQLDB())
	case database.TypeSQLite:
		driver, err = sqlite.Open(m.database.SQLDB())
	default:
		return fmt.Errorf("unsupported database type: %s", m.database.Type())
	}
	if err != nil {
		return fmt.Errorf("failed to create atlas driver: %w", err)
	}

	// Step 2: Inspect current database state
	currentSchema, err := driver.InspectSchema(ctx, schemaName, &schema.InspectOptions{})
	if err != nil {
		return fmt.Errorf("failed to inspect current schema: %w", err)
	}

	// Step 3: Load desired schema from HCL
	desiredSchema, err := m.loadHCLSchema()
	if err != nil {
		return fmt.Errorf("failed to load HCL schema: %w", err)
	}

	// Step 4: Compute diff between current and desired state
	changes, err := driver.SchemaDiff(currentSchema, desiredSchema)
	if err != nil {
		return fmt.Errorf("failed to compute schema diff: %w", err)
	}

	if len(changes) == 0 {
		slog.Debug("No schema changes detected")
		return nil
	}

	// Step 5: Generate migration plan
	plan, err := driver.PlanChanges(ctx, "", changes)
	if err != nil {
		return fmt.Errorf("failed to plan changes: %w", err)
	}

	// Step 6: Generate migration SQL
	migrationSQL := m.formatMigrationSQL(plan)

	// Step 7: Write migration file
	timestamp := time.Now().Format("20060102150405")
	migrationFile := fmt.Sprintf("%s.sql", timestamp)
	migrationPath := filepath.Join(migrationDir, migrationFile)

	migrationContent := fmt.Sprintf(`-- Generated at: %s
%s`, time.Now().Format(time.RFC3339), migrationSQL)

	if err := os.WriteFile(migrationPath, []byte(migrationContent), 0644); err != nil {
		return fmt.Errorf("failed to write migration file: %w", err)
	}

	slog.Info("Migration file created", slog.String("path", migrationPath))

	return nil
}

// loadHCLSchema loads the desired schema from HCL file using Atlas
func (m *MigrationGenerator) loadHCLSchema() (*schema.Schema, error) {
	hclSchemaFile := GetHCLSchemaFile(m.database.Type())
	if hclSchemaFile == "" {
		return nil, fmt.Errorf("unsupported database type for HCL schema: %s", m.database.Type())
	}

	hclData, err := os.ReadFile(hclSchemaFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read HCL file: %w", err)
	}

	// Determine which evaluator and schema name to use based on database type
	var s schema.Schema
	var defaultSchemaName string

	switch m.database.Type() {
	case database.TypePostgreSQL:
		defaultSchemaName = postgresSchemaName
		if err := postgres.EvalHCLBytes(hclData, &s, nil); err != nil {
			return nil, fmt.Errorf("failed to parse HCL schema: %w", err)
		}

	case database.TypeSQLite:
		defaultSchemaName = sqliteSchemaName
		if err := sqlite.EvalHCLBytes(hclData, &s, nil); err != nil {
			return nil, fmt.Errorf("failed to parse HCL schema: %w", err)
		}

	default:
		return nil, fmt.Errorf("unsupported database type: %s", m.database.Type())
	}

	// Always set schema name to match the database type
	// This ensures PostgreSQL uses "public" and SQLite uses "main"
	s.Name = defaultSchemaName

	return &s, nil
}

// formatMigrationSQL formats the migration plan into SQL statements
func (m *MigrationGenerator) formatMigrationSQL(plan *migrate.Plan) string {
	if len(plan.Changes) == 0 {
		return ""
	}

	var sql string
	for _, change := range plan.Changes {
		if change.Comment != "" {
			sql += "-- " + change.Comment + "\n"
		}
		sql += change.Cmd
		if change.Cmd != "" && change.Cmd[len(change.Cmd)-1] != ';' {
			sql += ";"
		}
		sql += "\n"
	}
	return sql
}
