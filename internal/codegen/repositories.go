package codegen

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/archesai/archesai/internal/parsers"
)

// GenerateRepositories generates all repository interfaces and implementations
func (g *Generator) GenerateRepositories(schemas map[string]*parsers.ProcessedSchema) error {
	for name, processed := range schemas {
		// Only generate repositories for schemas with x-codegen.repository config
		if processed.XCodegen != nil && processed.XCodegen.Repository != nil {
			if err := g.generateRepositoryForSchema(processed); err != nil {
				return fmt.Errorf("failed to generate repository for %s: %w", name, err)
			}
		}
	}
	return nil
}

// generateRepositoryForSchema generates repository interface and implementations for a schema
func (g *Generator) generateRepositoryForSchema(
	schema *parsers.ProcessedSchema,
) error {
	// Extract additional methods from x-codegen configuration
	var additionalMethods []map[string]interface{}
	if schema.XCodegen != nil && schema.XCodegen.Repository != nil {
		for _, method := range schema.XCodegen.Repository.AdditionalMethods {
			var params []map[string]string
			for _, param := range method.Params {
				params = append(params, map[string]string{
					"Name": param,
					"Type": "string", // Default type, should be improved
				})
			}

			returns := []string{"error"}
			switch method.Returns {
			case "single":
				returns = []string{"*" + schema.Name, "error"}
			case "multiple":
				returns = []string{"[]*" + schema.Name, "error"}
			}

			additionalMethods = append(additionalMethods, map[string]interface{}{
				"Name":       method.Name,
				"Parameters": params,
				"Returns":    returns,
			})
		}
	}

	// Extract exclusion lists from schema
	excludeFromCreate := []string{}
	excludeFromUpdate := []string{}
	if schema.XCodegen != nil && schema.XCodegen.Repository != nil {
		excludeFromCreate = schema.XCodegen.Repository.ExcludeFromCreate
		excludeFromUpdate = schema.XCodegen.Repository.ExcludeFromUpdate
	}

	// Generate repository interface
	data := map[string]interface{}{
		"Package": "repositories",
		"Entities": []map[string]interface{}{
			{
				"Name":              schema.Name,
				"Type":              schema.Name,
				"Fields":            schema.Fields,
				"AdditionalMethods": additionalMethods,
				"ExcludeFromCreate": excludeFromCreate,
				"ExcludeFromUpdate": excludeFromUpdate,
			},
		},
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
	implData := map[string]interface{}{
		"Package": "repositories", // Implementation package
		"Entities": []map[string]interface{}{
			{
				"Name":              schema.Name,
				"Type":              schema.Name,
				"Fields":            schema.Fields,
				"AdditionalMethods": additionalMethods,
				"ExcludeFromCreate": excludeFromCreate,
				"ExcludeFromUpdate": excludeFromUpdate,
			},
		},
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
