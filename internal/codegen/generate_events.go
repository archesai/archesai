package codegen

import (
	"bytes"
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

		// Render to buffer
		var buf bytes.Buffer
		if err := g.renderer.Render(&buf, "events.go.tmpl", data); err != nil {
			return fmt.Errorf("failed to render events for %s: %w", schema.Name, err)
		}

		// Write using storage
		eventFilePath := filepath.Join(
			g.outputDir, "generated", "core", "events",
			strings.ToLower(schema.Name)+"_events.gen.go",
		)
		if err := g.storage.WriteFile(eventFilePath, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write events for %s: %w", schema.Name, err)
		}

	}
	return nil
}
