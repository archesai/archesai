package codegen

import (
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
	for _, schema := range ctx.SpecDef.Schemas {
		if schema.XCodegenSchemaType != parsers.XCodegenSchemaTypeEntity {
			continue
		}

		if err := generateSQLQueriesForSchema(ctx, schema, "postgres"); err != nil {
			return fmt.Errorf("failed to generate PostgreSQL queries for %s: %w", schema.Name, err)
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
			"postgres",
			"repositories",
			strings.ToLower(schema.Name)+"_repository.gen.go",
		)
		if err := ctx.RenderToFile("postgres.go.tmpl", outputPath, data); err != nil {
			return fmt.Errorf(
				"failed to generate PostgreSQL repository for %s: %w",
				schema.Name,
				err,
			)
		}
	}
	return nil
}

func generateSQLQueriesForSchema(
	ctx *GeneratorContext,
	schema *parsers.SchemaDef,
	dbType string,
) error {
	modelImportPath, repositoryInterface := getRepositoryImportPaths(ctx, schema)
	data := &RepositoriesTemplateData{
		Entity:              schema,
		ProjectName:         ctx.ProjectName,
		ModelImportPath:     modelImportPath,
		RepositoryInterface: repositoryInterface,
	}

	outputPath := filepath.Join(
		"infrastructure",
		dbType,
		"queries",
		strings.ToLower(schema.Name)+"s.gen.sql",
	)
	return ctx.RenderToFile("sql_queries.sql.tmpl", outputPath, data)
}
