package codegen

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/archesai/archesai/internal/parsers"
)

// GenerateEvents generates domain events for entities
func (g *Generator) GenerateEvents(schemas map[string]*parsers.ProcessedSchema) error {
	for name, processed := range schemas {
		// Only generate events for entities, not value objects
		if processed.XCodegen != nil && processed.XCodegen.SchemaType != "valueobject" {
			// Check if entity has an ID field
			hasIDField := false
			fields := processed.Fields

			for _, field := range fields {
				if field.FieldName == "ID" {
					hasIDField = true
					break
				}
			}

			// Only generate events for entities with ID fields
			if hasIDField && processed.Schema != nil {
				if err := g.generateEventsForSchema(name); err != nil {
					return fmt.Errorf(
						"failed to generate events for %s: %w",
						name,
						err,
					)
				}
			}
		}
	}
	return nil
}

// generateEventsForSchema generates domain events for a schema
func (g *Generator) generateEventsForSchema(name string) error {
	title := name

	data := map[string]interface{}{
		"Package": "events",
		"Domain":  title,
		"Entities": []map[string]interface{}{
			{
				"Name":            title,
				"NameLower":       strings.ToLower(title),
				"NamePlural":      Pluralize(title),
				"NamePluralLower": strings.ToLower(Pluralize(title)),
			},
		},
	}

	outputPath := filepath.Join(
		"internal/core/events",
		strings.ToLower(title)+"_events.gen.go",
	)

	tmpl, ok := g.templates["events.tmpl"]
	if !ok {
		return fmt.Errorf("events template not found")
	}

	return g.filewriter.WriteTemplate(outputPath, tmpl, data)
}
