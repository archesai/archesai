package codegen

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/archesai/archesai/internal/parsers"
)

// ControllersTemplateData defines a template data structure
type ControllersTemplateData struct {
	Tag        string
	Operations []parsers.OperationDef
}

// GenerateControllers generates all HTTP controllers grouped by tag
func (g *Generator) GenerateControllers(
	operations []parsers.OperationDef,
) error {

	// Group operations by their domain tag (first tag)
	grouped := make(map[string][]parsers.OperationDef)
	for _, op := range operations {
		grouped[op.Tag] = append(grouped[op.Tag], op)
	}

	// Generate a handler file for each tag
	for tag, ops := range grouped {
		if err := g.generateControllerForTag(tag, ops); err != nil {
			return err
		}
	}

	return nil
}

// generateControllerForTag generates a controller for a specific domain
func (g *Generator) generateControllerForTag(
	tag string,
	operations []parsers.OperationDef,
) error {

	// Sort operations in standard REST order: POST, GET (singular), GET (list), PATCH/PUT, DELETE
	sort.SliceStable(operations, func(i, j int) bool {
		return getOperationOrder(operations[i]) < getOperationOrder(operations[j])
	})

	data := &ControllersTemplateData{
		Tag:        tag,
		Operations: operations,
	}

	outputPath := filepath.Join(
		"internal/adapters/http/controllers",
		strings.ToLower(tag)+".gen.go",
	)

	tmpl, ok := g.templates["controller.tmpl"]
	if !ok {
		return fmt.Errorf("controller template not found")
	}

	if err := g.filewriter.WriteTemplate(outputPath, tmpl, data); err != nil {
		return fmt.Errorf("failed to generate handler for %s: %w", tag, err)
	}

	return nil
}

// getOperationOrder returns the sort order for an operation
// Standard REST order: POST, GET (singular), GET (list), PATCH/PUT, DELETE
func getOperationOrder(op parsers.OperationDef) int {
	switch op.Method {
	case "POST":
		return 1
	case "GET":
		if strings.HasPrefix(op.ID, "List") {
			return 3
		}
		return 2
	case "PATCH", "PUT":
		return 4
	case "DELETE":
		return 5
	default:
		return 6
	}
}
