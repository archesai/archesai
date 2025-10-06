package codegen

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/archesai/archesai/internal/parsers"
)

// EventsTemplateData defines a template data structure
type EventsTemplateData struct {
	Entity *parsers.SchemaDef
}

// GenerateEvents generates domain events for all entities
func (g *Generator) GenerateEvents(schemas []*parsers.SchemaDef) error {
	for _, schema := range schemas {
		if schema.XCodegen != nil && schema.XCodegen.SchemaType != "valueobject" &&
			schema.Schema != nil {
			data := &EventsTemplateData{
				Entity: schema,
			}

			outputPath := filepath.Join(
				"internal/core/events",
				strings.ToLower(schema.Name)+"_events.gen.go",
			)

			tmpl, ok := g.templates["events.tmpl"]
			if !ok {
				return fmt.Errorf("events template not found")
			}

			if err := g.filewriter.WriteTemplate(outputPath, tmpl, data); err != nil {
				return fmt.Errorf("failed to generate events for %s: %w", schema.Name, err)
			}
		}
	}
	return nil
}
