package codegen

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/archesai/archesai/internal/parsers"
)

// RepositoriesTemplateData defines a template data structure
type RepositoriesTemplateData struct {
	Entity *parsers.SchemaDef
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

	// Generate repository interface
	data := &RepositoriesTemplateData{
		Entity: schema,
	}

	// Generate interface in repositories folder
	outputPath := filepath.Join(
		"internal/core/repositories",
		strings.ToLower(schema.Name)+".gen.go",
	)

	tmpl, ok := g.templates["repository.tmpl"]
	if !ok {
		return fmt.Errorf("repository template not found")
	}

	if err := g.filewriter.WriteTemplate(outputPath, tmpl, data); err != nil {
		return err
	}

	// Generate concrete implementations with different package
	implData := &RepositoriesTemplateData{
		Entity: schema,
	}

	// PostgreSQL
	if tmpl, ok := g.templates["repository_postgres.tmpl"]; ok {
		outputPath := filepath.Join(
			"internal/infrastructure/persistence/postgres/repositories",
			strings.ToLower(schema.Name)+"_repository.gen.go",
		)
		if err := g.filewriter.WriteTemplate(outputPath, tmpl, implData); err != nil {
			return fmt.Errorf("failed to generate PostgreSQL repository: %w", err)
		}
	}

	// SQLite
	if tmpl, ok := g.templates["repository_sqlite.tmpl"]; ok {
		outputPath := filepath.Join(
			"internal/infrastructure/persistence/sqlite/repositories",
			strings.ToLower(schema.Name)+"_repository.gen.go",
		)
		if err := g.filewriter.WriteTemplate(outputPath, tmpl, implData); err != nil {
			return fmt.Errorf("failed to generate SQLite repository: %w", err)
		}
	}

	return nil
}
