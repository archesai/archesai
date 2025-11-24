package codegen

import (
	"bytes"
	"fmt"
	"io"
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
		var buf bytes.Buffer
		if err := g.generateSchema(processed, &buf); err != nil {
			return fmt.Errorf(
				"failed to generate %s %s: %w",
				processed.XCodegenSchemaType,
				processed.Name,
				err,
			)
		}
		outputPath := filepath.Join(
			"generated",
			"core",
			strings.ToLower(processed.Name)+".gen.go",
		)
		err := g.storage.WriteFile(outputPath, buf.Bytes(), 0644)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", outputPath, err)
		}
	}

	return nil
}

// generateSchema generates a single model file
func (g *Generator) generateSchema(
	schema *parsers.SchemaDef,
	out io.Writer,
) error {
	data := SchemasTemplateData{
		Package: "core",
		Schema:  schema,
	}
	if err := g.renderer.Render(out, "schema.go.tmpl", data); err != nil {
		return fmt.Errorf("failed to render schema: %w", err)
	}
	return nil

}
