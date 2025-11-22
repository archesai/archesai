package codegen

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"sort"

	"golang.org/x/sync/errgroup"

	"github.com/archesai/archesai/internal/parsers"
	"github.com/archesai/archesai/pkg/database"
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
		if schema.XCodegenSchemaType == "entity" {
			entities = append(entities, schema)
		}
	}

	// Sort entities by their table name (snake_case version of schema name)
	sort.Slice(entities, func(i, j int) bool {
		tableNameI := parsers.SnakeCase(entities[i].Name)
		tableNameJ := parsers.SnakeCase(entities[j].Name)
		return tableNameI < tableNameJ
	})

	// Phase 1: Write HCL schema files in parallel
	eg := &errgroup.Group{}

	// Generate PostgreSQL schema
	postgresData := HCLTemplateData{
		Schemas:      entities,
		DatabaseType: database.TypePostgreSQL,
	}
	postgresHCLPath := filepath.Join(
		g.outputDir,
		"generated",
		"infrastructure",
		"persistence",
		"postgres",
		"schema.gen.hcl",
	)
	eg.Go(func() error {
		var buf bytes.Buffer
		if err := g.renderer.Render(&buf, "db.hcl.tmpl", postgresData); err != nil {
			return fmt.Errorf("failed to render PostgreSQL HCL: %w", err)
		}
		if err := g.storage.WriteFile(postgresHCLPath, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write PostgreSQL HCL file: %w", err)
		}
		return nil
	})

	// Generate SQLite schema
	sqliteData := HCLTemplateData{
		Schemas:      entities,
		DatabaseType: database.TypeSQLite,
	}
	sqliteHCLPath := filepath.Join(
		g.outputDir,
		"generated",
		"infrastructure",
		"persistence",
		"sqlite",
		"schema.gen.hcl",
	)
	eg.Go(func() error {
		var buf bytes.Buffer
		if err := g.renderer.Render(&buf, "db.hcl.tmpl", sqliteData); err != nil {
			return fmt.Errorf("failed to render SQLite HCL: %w", err)
		}
		if err := g.storage.WriteFile(sqliteHCLPath, buf.Bytes(), 0644); err != nil {
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
