package codegen

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/archesai/archesai/internal/parsers"
)

// SchemasTemplateData defines a template data structure
type SchemasTemplateData struct {
	Package    string
	Schema     *parsers.SchemaDef
	OutputPath string // Import path for generated code
}

// GenerateSchemas generates all model types (DTOs, entities, value objects)
func (g *Generator) GenerateSchemas(schemas []*parsers.SchemaDef) error {
	for _, processed := range schemas {
		if err := g.generateSchema(processed, nil); err != nil {
			return fmt.Errorf(
				"failed to generate %s %s: %w",
				processed.XCodegenSchemaType,
				processed.Name,
				err,
			)
		}
	}

	return nil
}

// generateSchema generates a single model file (simplified - no more batching)
func (g *Generator) generateSchema(
	schema *parsers.SchemaDef,
	customOutputDir *string,
) error {

	// Use a single package for all schema types
	var packageName, outputDir string

	switch schema.XCodegenSchemaType {
	case schemaTypeEntity, schemaTypeValueObject:
		packageName = "models"
		outputDir = filepath.Join(g.outputDir, "generated", "core", "models")
		if customOutputDir != nil {
			outputDir = *customOutputDir
			packageName = filepath.Base(*customOutputDir)
		}
	default:
		return fmt.Errorf(
			"unsupported model type: %s for schema %s",
			schema.XCodegenSchemaType,
			schema.Name,
		)
	}

	importPath := "github.com/archesai/archesai" + strings.TrimPrefix(g.outputDir, ".")

	data := SchemasTemplateData{
		Package:    packageName,
		Schema:     schema,
		OutputPath: importPath,
	}

	// Render to buffer
	var buf bytes.Buffer
	if err := g.renderer.Render(&buf, "schema.go.tmpl", data); err != nil {
		return fmt.Errorf("failed to render schema: %w", err)
	}

	// Write using storage interface
	outputPath := filepath.Join(outputDir, strings.ToLower(schema.Name)+".gen.go")
	return g.storage.WriteFile(outputPath, buf.Bytes(), 0644)
}
