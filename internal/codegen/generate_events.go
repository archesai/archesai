package codegen

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/archesai/archesai/internal/parsers"
)

// EventsTemplateData defines a template data structure
type EventsTemplateData struct {
	Entity     *parsers.SchemaDef
	OutputPath string // Import path for generated code
}

// GenerateEvents generates domain events for all entities
func (g *Generator) GenerateEvents(schemas []*parsers.SchemaDef) error {

	outputPath := "github.com/archesai/archesai" + strings.TrimPrefix(g.outputDir, ".")

	for _, schema := range schemas {

		if schema.XCodegenSchemaType != parsers.XCodegenSchemaTypeEntity {
			continue
		}

		data := &EventsTemplateData{
			Entity:     schema,
			OutputPath: outputPath,
		}

		eventFilePath := filepath.Join(
			g.outputDir, "generated", "core", "events",
			strings.ToLower(schema.Name)+"_events.gen.go",
		)

		tmpl, ok := g.templates["events.tmpl"]
		if !ok {
			return fmt.Errorf("events template not found")
		}

		if err := g.filewriter.WriteTemplate(eventFilePath, tmpl, data); err != nil {
			return fmt.Errorf("failed to generate events for %s: %w", schema.Name, err)
		}

	}
	return nil
}
