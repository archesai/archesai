package codegen

import (
	"context"
	"fmt"
	"os/exec"
	"sort"

	"golang.org/x/sync/errgroup"

	database "github.com/archesai/archesai/internal/infrastructure/persistence"
	"github.com/archesai/archesai/internal/parsers"
)

// HCLTemplateData defines the template data for HCL generation
type HCLTemplateData struct {
	Schemas      []*parsers.SchemaDef
	DatabaseType database.Type
}

// GenerateHCL generates HCL database schema from OpenAPI schemas
func (g *Generator) GenerateHCL(schemas []*parsers.SchemaDef) error {
	var entities []*parsers.SchemaDef
	for _, schema := range schemas {
		if schema.GetSchemaType() == "entity" {
			entities = append(entities, schema)
		}
	}

	// Sort entities by their table name (snake_case version of schema name)
	sort.Slice(entities, func(i, j int) bool {
		tableNameI := parsers.SnakeCase(entities[i].Name)
		tableNameJ := parsers.SnakeCase(entities[j].Name)
		return tableNameI < tableNameJ
	})

	tmpl, ok := g.templates["schema_hcl.tmpl"]
	if !ok {
		return fmt.Errorf("HCL template not found")
	}

	// Phase 1: Write HCL schema files in parallel
	eg := &errgroup.Group{}

	// Generate PostgreSQL schema
	postgresData := HCLTemplateData{
		Schemas:      entities,
		DatabaseType: database.TypePostgreSQL,
	}
	eg.Go(func() error {
		if err := g.filewriter.WriteTemplate(PostgresHCLSchemaFile, tmpl, postgresData); err != nil {
			return fmt.Errorf("failed to write PostgreSQL HCL file: %w", err)
		}
		return nil
	})

	// Generate SQLite schema
	sqliteData := HCLTemplateData{
		Schemas:      entities,
		DatabaseType: database.TypeSQLite,
	}
	eg.Go(func() error {
		if err := g.filewriter.WriteTemplate(SQLiteHCLSchemaFile, tmpl, sqliteData); err != nil {
			return fmt.Errorf("failed to write SQLite HCL file: %w", err)
		}
		return nil
	})

	// Wait for HCL file generation to complete
	if err := eg.Wait(); err != nil {
		return err
	}

	// Check Docker availability before migrations
	if err := exec.Command("docker", "version").Run(); err != nil {
		return fmt.Errorf("docker not available, skipping HCL formatting: %w", err)
	}

	// Phase 2: Generate migrations in parallel
	ctx := context.Background()
	eg2 := &errgroup.Group{}

	eg2.Go(func() error {
		if err := g.GenerateMigrations(ctx, database.TypePostgreSQL); err != nil {
			return fmt.Errorf("failed to generate PostgreSQL migration: %w", err)
		}
		return nil
	})

	eg2.Go(func() error {
		if err := g.GenerateMigrations(ctx, database.TypeSQLite); err != nil {
			return fmt.Errorf("failed to generate SQLite migration: %w", err)
		}
		return nil
	})

	// Wait for migrations to complete
	if err := eg2.Wait(); err != nil {
		return err
	}

	return nil
}
