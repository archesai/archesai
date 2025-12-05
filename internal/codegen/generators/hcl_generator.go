package generators

import (
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"path/filepath"
	"sort"

	"golang.org/x/sync/errgroup"

	"github.com/archesai/archesai/internal/spec"
	"github.com/archesai/archesai/internal/strutil"
	"github.com/archesai/archesai/pkg/database"
	"github.com/archesai/archesai/pkg/storage"
)

// HCLTemplateData holds the data for rendering HCL schema templates.
type HCLTemplateData struct {
	Schemas      []*spec.Schema
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
	var entities []*spec.Schema
	for _, schema := range ctx.Spec.Schemas {
		if schema.XCodegenSchemaType == spec.XCodegenSchemaTypeEntity {
			entities = append(entities, schema)
		}
	}

	sort.Slice(entities, func(i, j int) bool {
		return strutil.SnakeCase(entities[i].Name) < strutil.SnakeCase(entities[j].Name)
	})

	// Generate HCL files in parallel
	eg := &errgroup.Group{}

	// PostgreSQL
	eg.Go(func() error {
		data := HCLTemplateData{Schemas: entities, DatabaseType: database.TypePostgreSQL}
		path := filepath.Join("infrastructure", "postgres", "schema.gen.hcl")
		return ctx.RenderToFile("hcl.tmpl", path, data)
	})

	// SQLite
	eg.Go(func() error {
		data := HCLTemplateData{Schemas: entities, DatabaseType: database.TypeSQLite}
		path := filepath.Join("infrastructure", "sqlite", "schema.gen.hcl")
		return ctx.RenderToFile("hcl.tmpl", path, data)
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
