package codegen

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"path/filepath"
	"sort"

	"golang.org/x/sync/errgroup"

	"github.com/archesai/archesai/internal/parsers"
	"github.com/archesai/archesai/pkg/database"
	"github.com/archesai/archesai/pkg/storage"
)

// HCLTemplateData holds the data for rendering HCL schema templates.
type HCLTemplateData struct {
	Schemas      []*parsers.SchemaDef
	DatabaseType database.Type
}

// HCLGenerator generates Atlas HCL schema files for database migrations.
type HCLGenerator struct{}

// Name returns the generator name.
func (g *HCLGenerator) Name() string { return "hcl" }

// Priority returns the generator priority.
func (g *HCLGenerator) Priority() int { return PriorityLast }

// Generate creates HCL schema files for database migrations.
func (g *HCLGenerator) Generate(ctx *GeneratorContext) error {
	var entities []*parsers.SchemaDef
	for _, schema := range ctx.SpecDef.Schemas {
		if schema.XCodegenSchemaType == "entity" {
			entities = append(entities, schema)
		}
	}

	sort.Slice(entities, func(i, j int) bool {
		return parsers.SnakeCase(entities[i].Name) < parsers.SnakeCase(entities[j].Name)
	})

	// Generate HCL files in parallel
	eg := &errgroup.Group{}

	// PostgreSQL
	eg.Go(func() error {
		data := HCLTemplateData{Schemas: entities, DatabaseType: database.TypePostgreSQL}
		var buf bytes.Buffer
		if err := ctx.Renderer.Render(&buf, "hcl.tmpl", data); err != nil {
			return fmt.Errorf("failed to render PostgreSQL HCL: %w", err)
		}
		path := filepath.Join("infrastructure", "postgres", "schema.gen.hcl")
		return ctx.Storage.WriteFile(path, buf.Bytes(), 0644)
	})

	// SQLite
	eg.Go(func() error {
		data := HCLTemplateData{Schemas: entities, DatabaseType: database.TypeSQLite}
		var buf bytes.Buffer
		if err := ctx.Renderer.Render(&buf, "hcl.tmpl", data); err != nil {
			return fmt.Errorf("failed to render SQLite HCL: %w", err)
		}
		path := filepath.Join("infrastructure", "sqlite", "schema.gen.hcl")
		return ctx.Storage.WriteFile(path, buf.Bytes(), 0644)
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	// Skip migrations for memory storage
	if _, isMemory := ctx.Storage.(*storage.MemoryStorage); isMemory {
		return nil
	}

	// Check Docker availability
	if err := exec.Command("docker", "version").Run(); err != nil {
		return fmt.Errorf("docker not available, skipping migrations: %w", err)
	}

	// Generate migrations in parallel
	bgCtx := context.Background()
	eg2 := &errgroup.Group{}

	eg2.Go(func() error {
		return runMigrationGenerator(bgCtx, ctx.Storage.BaseDir(), database.TypePostgreSQL)
	})

	eg2.Go(func() error {
		return runMigrationGenerator(bgCtx, ctx.Storage.BaseDir(), database.TypeSQLite)
	})

	return eg2.Wait()
}

func runMigrationGenerator(ctx context.Context, outputDir string, dbType database.Type) error {
	m := &MigrationGenerator{outputDir: outputDir}
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
	return m.GenerateMigration(ctx)
}
