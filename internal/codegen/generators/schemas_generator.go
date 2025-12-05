package generators

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/archesai/archesai/internal/spec"
)

// SchemasTemplateData holds the data for rendering schema/model templates.
type SchemasTemplateData struct {
	Package string
	Schema  *spec.Schema
}

// SchemasGenerator generates model code from OpenAPI schemas.
type SchemasGenerator struct{}

// Name returns the generator name.
func (g *SchemasGenerator) Name() string { return "models" }

// Priority returns the generator priority.
func (g *SchemasGenerator) Priority() int { return PriorityNormal }

// Generate creates model code for each schema in the OpenAPI spec.
func (g *SchemasGenerator) Generate(ctx *GeneratorContext) error {
	internalContext := ctx.InternalContext()
	for _, schema := range ctx.Spec.Schemas {
		if schema.IsInternal(internalContext) {
			continue
		}

		data := SchemasTemplateData{
			Package: "models",
			Schema:  schema,
		}

		outputPath := filepath.Join("models", strings.ToLower(schema.Name)+".gen.go")
		if err := ctx.RenderToFile("schema.go.tmpl", outputPath, data); err != nil {
			return fmt.Errorf(
				"failed to generate %s %s: %w",
				schema.XCodegenSchemaType,
				schema.Name,
				err,
			)
		}
	}
	return nil
}
