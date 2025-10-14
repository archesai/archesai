package codegen

import (
	"fmt"
	"sort"

	"github.com/archesai/archesai/internal/parsers"
)

// HCLTemplateData defines the template data for HCL generation
type HCLTemplateData struct {
	Schemas []*parsers.SchemaDef
}

// GenerateHCL generates HCL database schema from OpenAPI schemas
func (g *Generator) GenerateHCL(schemas []*parsers.SchemaDef) error {
	// Filter for entities only (not value objects)
	var entities []*parsers.SchemaDef
	for _, schema := range schemas {
		if schema.GetSchemaType() == "entity" {
			entities = append(entities, schema)
		}
	}

	// Sort entities by their table name (snake_case version of schema name)
	sort.Slice(entities, func(i, j int) bool {
		tableNameI := parsers.SnakeCase(entities[i].Name)
		tableNameJ := parsers.SnakeCase(entities[j].Name)
		return tableNameI < tableNameJ
	})

	data := HCLTemplateData{
		Schemas: entities,
	}

	tmpl, ok := g.templates["schema_hcl.tmpl"]
	if !ok {
		return fmt.Errorf("HCL template not found")
	}

	outputPath := "schema.gen.hcl"
	return g.filewriter.WriteTemplate(outputPath, tmpl, data)
}
