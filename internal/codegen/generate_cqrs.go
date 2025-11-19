package codegen

import (
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
			var tmplPath string
			var output string
			if op.Method == "GET" {
				tmplPath = "query_handler.tmpl"
				output = "queries"
			} else {
				tmplPath = "command_handler.tmpl"
				output = "commands"
			}

			// Load the template
			tmpl, ok := g.templates[tmplPath]
			if !ok {
				return fmt.Errorf("command handler template not found")
			}

			// Get the output path
			outputPath := filepath.Join(
				g.outputDir, "generated", "application",
				output,
				strings.ToLower(tag),
				fmt.Sprintf("%s.gen.go", parsers.SnakeCase(op.ID)),
			)

			importPath := "github.com/archesai/archesai" + strings.TrimPrefix(g.outputDir, ".")

			// Create minimal template data
			data := &CQRSHandlerTemplateData{
				Operation:  &op,
				OutputPath: importPath,
			}

			// Write the handler file
			if err := g.filewriter.WriteTemplate(outputPath, tmpl, data); err != nil {
				return fmt.Errorf(
					"failed to generate handler for %s: %w",
					op.ID,
					err,
				)
			}
		}
	}

	return nil
}
