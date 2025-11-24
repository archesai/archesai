package codegen

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/archesai/archesai/internal/parsers"
)

// RepositoriesTemplateData defines a template data structure
type RepositoriesTemplateData struct {
	Entity     *parsers.SchemaDef
	OutputPath string // Import path for generated code
}

// GenerateRepositories generates all repository interfaces and implementations
func (g *Generator) GenerateRepositories(schemas []*parsers.SchemaDef) error {
	for _, schema := range schemas {
		if schema.XCodegenSchemaType == parsers.XCodegenSchemaTypeEntity {
			if err := g.generateRepositoryForSchema(schema); err != nil {
				return fmt.Errorf("failed to generate repository for %s: %w", schema.Name, err)
			}
		}
	}
	return nil
}

// generateRepositoryForSchema generates repository interface and implementations for a schema
func (g *Generator) generateRepositoryForSchema(
	schema *parsers.SchemaDef,
) error {

	importPath := "github.com/archesai/archesai" + strings.TrimPrefix(g.storage.BaseDir(), ".")

	// First, generate SQL queries for this schema
	if err := g.generateSQLQueriesForSchema(schema, importPath); err != nil {
		return fmt.Errorf("failed to generate SQL queries: %w", err)
	}

	// Generate repository interface
	data := &RepositoriesTemplateData{
		Entity:     schema,
		OutputPath: importPath,
	}

	// Generate interface in repositories folder
	var buf bytes.Buffer
	if err := g.renderer.Render(&buf, "repository.go.tmpl", data); err != nil {
		return fmt.Errorf("failed to render repository interface: %w", err)
	}

	outputPath := filepath.Join(
		"generated", "core", "repositories",
		strings.ToLower(schema.Name)+".gen.go",
	)
	if err := g.storage.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return err
	}

	// Generate concrete implementations with different package
	implData := &RepositoriesTemplateData{
		Entity:     schema,
		OutputPath: importPath,
	}

	// PostgreSQL
	buf.Reset()
	if err := g.renderer.Render(&buf, "repository_postgres.go.tmpl", implData); err != nil {
		return fmt.Errorf("failed to render PostgreSQL repository: %w", err)
	}

	outputPath = filepath.Join(
		"generated",
		"infrastructure",
		"persistence",
		"postgres",
		"repositories",
		strings.ToLower(schema.Name)+"_repository.gen.go",
	)
	if err := g.storage.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write PostgreSQL repository: %w", err)
	}

	// SQLite FIXME: SQLite repository generation is currently disabled
	// buf.Reset()
	// if err := g.renderer.Render(&buf, "repository_sqlite.go.tmpl", implData); err != nil {
	// 	return fmt.Errorf("failed to render SQLite repository: %w", err)
	// }

	// outputPath = filepath.Join(
	// 	"generated", "infrastructure", "persistence", "sqlite", "repositories",
	// 	strings.ToLower(schema.Name)+"_repository.gen.go",
	// )
	// if err := g.storage.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
	// 	return fmt.Errorf("failed to write SQLite repository: %w", err)
	// }

	return nil
}

// generateSQLQueriesForSchema generates SQL query files for SQLC
func (g *Generator) generateSQLQueriesForSchema(
	schema *parsers.SchemaDef,
	importPath string,
) error {
	data := &RepositoriesTemplateData{
		Entity:     schema,
		OutputPath: importPath,
	}

	// Generate PostgreSQL queries
	var buf bytes.Buffer
	if err := g.renderer.Render(&buf, "sql_queries.sql.tmpl", data); err != nil {
		return fmt.Errorf("failed to render PostgreSQL queries: %w", err)
	}

	postgresPath := filepath.Join(
		"generated",
		"infrastructure",
		"persistence",
		"postgres",
		"queries",
		strings.ToLower(schema.Name)+"s.sql",
	)
	if err := g.storage.WriteFile(postgresPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write PostgreSQL queries: %w", err)
	}

	// Generate SQLite queries (using same template for now)
	sqlitePath := filepath.Join(
		"generated",
		"infrastructure",
		"persistence",
		"sqlite",
		"queries",
		strings.ToLower(schema.Name)+"s.sql",
	)
	if err := g.storage.WriteFile(sqlitePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write SQLite queries: %w", err)
	}

	return nil
}
