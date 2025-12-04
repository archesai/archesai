package codegen

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/archesai/archesai/internal/parsers"
)

// SQLiteGenerator generates SQLite repository implementations.
type SQLiteGenerator struct{}

// Name returns the generator name.
func (g *SQLiteGenerator) Name() string { return "sqlite" }

// Priority returns the generator priority.
func (g *SQLiteGenerator) Priority() int { return PriorityNormal }

// Generate creates SQLite repository code for each entity schema.
func (g *SQLiteGenerator) Generate(ctx *GeneratorContext) error {
	// Generate db.gen.go
	dbPath := filepath.Join("infrastructure", "sqlite", "repositories", "db.gen.go")
	if err := ctx.RenderToFile("sqlite_db.go.tmpl", dbPath, nil); err != nil {
		return fmt.Errorf("failed to generate SQLite db file: %w", err)
	}

	for _, schema := range ctx.SpecDef.Schemas {
		if schema.XCodegenSchemaType != parsers.XCodegenSchemaTypeEntity {
			continue
		}

		if err := generateSQLQueriesForSchema(ctx, schema, "sqlite"); err != nil {
			return fmt.Errorf("failed to generate SQLite queries for %s: %w", schema.Name, err)
		}

		modelImportPath, repositoryInterface := getRepositoryImportPaths(ctx, schema)
		data := &RepositoriesTemplateData{
			Entity:              schema,
			ProjectName:         ctx.ProjectName,
			ModelImportPath:     modelImportPath,
			RepositoryInterface: repositoryInterface,
		}

		outputPath := filepath.Join(
			"infrastructure",
			"sqlite",
			"repositories",
			strings.ToLower(schema.Name)+"_repository.gen.go",
		)
		if err := ctx.RenderToFile("sqlite.go.tmpl", outputPath, data); err != nil {
			return fmt.Errorf("failed to generate SQLite repository for %s: %w", schema.Name, err)
		}
	}
	return nil
}
