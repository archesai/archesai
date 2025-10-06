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
	// Pass the OpenAPI document to the extraction pipeline for proper reference resolution
	for _, processed := range schemas {
		if processed.Schema == nil || processed.XCodegen == nil {
			continue
		}

		// Determine the schema type for this schema
		schemaType := "entity"
		if processed.XCodegen != nil && processed.XCodegen.SchemaType != "" {
			schemaType = string(processed.XCodegen.GetSchemaType())
		}

		// Re-extract with the OpenAPI document context for proper reference resolution
		reExtracted, err := g.jsonSchemaParser.ExtractSchema(
			processed.Schema,
			nil,
			schemaType,
		)
		if err != nil {
			return fmt.Errorf(
				"failed to re-extract schema %s with document context: %w",
				processed.Name,
				err,
			)
		}

		if err := g.generateModel(reExtracted, nil); err != nil {
			return fmt.Errorf(
				"failed to generate %s %s: %w",
				reExtracted.XCodegen.GetSchemaType(),
				reExtracted.Name,
				err,
			)
		}
	}

	return nil
}

// generateModel generates a single model file (simplified - no more batching)
func (g *Generator) generateModel(
	schema *parsers.SchemaDef,
	customOutputDir *string,
) error {

	// Determine package and output path based on model type
	var packageName, outputDir string

	schemaType := schema.GetSchemaType()
	switch schemaType {
	case schemaTypeEntity:
		packageName = "entities"
		outputDir = "internal/core/entities"
		if customOutputDir != nil {
			outputDir = *customOutputDir
			packageName = filepath.Base(*customOutputDir)
		}
	case schemaTypeValueObject:
		packageName = "valueobjects"
		outputDir = "internal/core/valueobjects"
		if customOutputDir != nil {
			outputDir = *customOutputDir
			packageName = filepath.Base(*customOutputDir)
		}
	default:
		return fmt.Errorf("unsupported model type: %s", schemaType)
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
