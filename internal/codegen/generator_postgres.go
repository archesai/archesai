package codegen

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/archesai/archesai/internal/parsers"
)

// PostgresGenerator generates PostgreSQL repository implementations.
type PostgresGenerator struct{}

// Name returns the generator name.
func (g *PostgresGenerator) Name() string { return "postgres" }

// Priority returns the generator priority.
func (g *PostgresGenerator) Priority() int { return PriorityNormal }

// Generate creates PostgreSQL repository code for each entity schema.
func (g *PostgresGenerator) Generate(ctx *GeneratorContext) error {
	internalContext := ctx.InternalContext()
	for _, schema := range ctx.SpecDef.Schemas {
		if schema.XCodegenSchemaType != parsers.XCodegenSchemaTypeEntity {
			continue
		}

		if err := generateSQLQueriesForSchema(ctx, schema, "postgres"); err != nil {
			return fmt.Errorf("failed to generate PostgreSQL queries for %s: %w", schema.Name, err)
		}

		modelImportPath, repositoryInterface := getRepositoryImportPaths(
			ctx.ProjectName,
			internalContext,
			schema,
		)
		data := &RepositoriesTemplateData{
			Entity:              schema,
			ProjectName:         ctx.ProjectName,
			ModelImportPath:     modelImportPath,
			RepositoryInterface: repositoryInterface,
		}

		var buf bytes.Buffer
		if err := ctx.Renderer.Render(&buf, "postgres.go.tmpl", data); err != nil {
			return fmt.Errorf("failed to render PostgreSQL repository for %s: %w", schema.Name, err)
		}

		outputPath := filepath.Join(
			"infrastructure",
			"postgres",
			"repositories",
			strings.ToLower(schema.Name)+"_repository.gen.go",
		)
		if err := ctx.Storage.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write PostgreSQL repository for %s: %w", schema.Name, err)
		}
	}
	return nil
}

func generateSQLQueriesForSchema(
	ctx *GeneratorContext,
	schema *parsers.SchemaDef,
	dbType string,
) error {
	internalContext := ctx.InternalContext()
	modelImportPath, repositoryInterface := getRepositoryImportPaths(
		ctx.ProjectName,
		internalContext,
		schema,
	)

	data := &RepositoriesTemplateData{
		Entity:              schema,
		ProjectName:         ctx.ProjectName,
		ModelImportPath:     modelImportPath,
		RepositoryInterface: repositoryInterface,
	}

	var buf bytes.Buffer
	if err := ctx.Renderer.Render(&buf, "sql_queries.sql.tmpl", data); err != nil {
		return fmt.Errorf("failed to render %s queries: %w", dbType, err)
	}

	outputPath := filepath.Join(
		"infrastructure",
		dbType,
		"queries",
		strings.ToLower(schema.Name)+"s.gen.sql",
	)
	return ctx.Storage.WriteFile(outputPath, buf.Bytes(), 0644)
}
