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
)

// DatabaseHCLTemplateData holds the data for rendering HCL schema templates.
type DatabaseHCLTemplateData struct {
	Schemas      []*spec.Schema
	DatabaseType database.Type
}

// GenerateHCL creates HCL schema files for database migrations.
func GenerateHCL(ctx *GeneratorContext) error {
	var entities []*spec.Schema
	for _, schema := range ctx.Spec.Schemas {
		if schema.XCodegenSchemaType == spec.SchemaTypeEntity {
			entities = append(entities, schema)
		}
	}

	sort.Slice(entities, func(i, j int) bool {
		return strutil.SnakeCase(entities[i].Name) < strutil.SnakeCase(entities[j].Name)
	})

	slog.Debug("DatabaseHCL generator starting", slog.Int("entities", len(entities)))

	// Generate HCL files in parallel
	eg := &errgroup.Group{}

	// PostgreSQL
	eg.Go(func() error {
		data := DatabaseHCLTemplateData{
			Schemas:      entities,
			DatabaseType: database.TypePostgreSQL,
		}
		path := filepath.Join("database", "postgres", "schema.gen.hcl")
		return ctx.RenderToFile("hcl.tmpl", path, data)
	})

	// SQLite
	eg.Go(func() error {
		data := DatabaseHCLTemplateData{Schemas: entities, DatabaseType: database.TypeSQLite}
		path := filepath.Join("database", "sqlite", "schema.gen.hcl")
		return ctx.RenderToFile("hcl.tmpl", path, data)
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	// Skip migrations for memory storage
	// FIXME

	// Check Docker availability
	if err := exec.Command("docker", "version").Run(); err != nil {
		slog.Debug("DatabaseHCL generator skipping migrations (docker not available)")
		return fmt.Errorf("docker not available, skipping migrations: %w", err)
	}

	slog.Debug("DatabaseHCL generator running migrations")

	// Generate migrations in parallel
	bgCtx := context.Background()
	eg2 := &errgroup.Group{}

	eg2.Go(func() error {
		return runMigrationGenerator(bgCtx, ctx.Storage.BaseDir(), database.TypePostgreSQL)
	})

	eg2.Go(func() error {
		return runMigrationGenerator(bgCtx, ctx.Storage.BaseDir(), database.TypeSQLite)
	})

	if err := eg2.Wait(); err != nil {
		return err
	}

	slog.Debug("DatabaseHCL generator completed", slog.Int("entities", len(entities)))

	return nil
}

func runMigrationGenerator(ctx context.Context, outputDir string, dbType database.Type) error {
	m := &DatabaseMigrationsGenerator{outputDir: outputDir}
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
