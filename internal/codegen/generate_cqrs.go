package codegen

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/archesai/archesai/internal/parsers"
)

// CQRSHandlerTemplateData defines a template data structure
type CQRSHandlerTemplateData struct {
	Operation  *parsers.OperationDef
	OutputPath string // Import path for generated code
}

// GenerateCommandQueryHandlers generates command and query handlers
func (g *Generator) GenerateCommandQueryHandlers(
	operations []parsers.OperationDef,
) error {

	grouped := make(map[string][]parsers.OperationDef)
	for _, op := range operations {
		grouped[op.Tag] = append(grouped[op.Tag], op)
	}

	// For each domain, generate command handlers for write operations
	for tag, ops := range grouped {

		for _, op := range ops {
			// Determine template and output directory based on method
			var tmplName string
			var output string
			if op.Method == "GET" {
				tmplName = "query_handler.go.tmpl"
				output = "queries"
			} else {
				tmplName = "command_handler.go.tmpl"
				output = "commands"
			}

			// Get the output path
			outputPath := filepath.Join(
				"generated", "application",
				output,
				strings.ToLower(tag),
				fmt.Sprintf("%s.gen.go", parsers.SnakeCase(op.ID)),
			)

			importPath := "github.com/archesai/archesai" + strings.TrimPrefix(
				g.storage.BaseDir(),
				".",
			)

			// Create minimal template data
			data := &CQRSHandlerTemplateData{
				Operation:  &op,
				OutputPath: importPath,
			}

			// Render to buffer
			var buf bytes.Buffer
			if err := g.renderer.Render(&buf, tmplName, data); err != nil {
				return fmt.Errorf("failed to render handler for %s: %w", op.ID, err)
			}

			// Write the handler file
			if err := g.storage.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
				return fmt.Errorf("failed to write handler for %s: %w", op.ID, err)
			}
		}
	}

	return nil
}
