package codegen

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/archesai/archesai/internal/parsers"
)

// SchemasTemplateData defines a template data structure
type SchemasTemplateData struct {
	Package string
	Schema  *parsers.SchemaDef
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
		outputDir = "internal/core/models"
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

	data := SchemasTemplateData{
		Package: packageName,
		Schema:  schema,
	}

	outputPath := filepath.Join(outputDir, strings.ToLower(schema.Name)+".gen.go")
	tmpl, ok := g.templates["schema.tmpl"]
	if !ok {
		return fmt.Errorf("schema template not found")
	}

	return g.filewriter.WriteTemplate(outputPath, tmpl, data)
}
