package codegen

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/archesai/archesai/internal/parsers"
)

// ControllersTemplateData defines a template data structure
type ControllersTemplateData struct {
	Tag               string
	Operations        []parsers.OperationDef
	OutputPath        string // Import path for generated code
	HasCustomHandlers bool   // Whether this controller has any custom handlers
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

	importPath := "github.com/archesai/archesai" + strings.TrimPrefix(g.storage.BaseDir(), ".")

	// Check if any operations have custom handlers
	hasCustom := false
	for _, op := range operations {
		if op.XCodegenCustomHandler {
			hasCustom = true
			break
		}
	}

	data := &ControllersTemplateData{
		Tag:               tag,
		Operations:        operations,
		OutputPath:        importPath,
		HasCustomHandlers: hasCustom,
	}

	// Render to buffer
	var buf bytes.Buffer
	if err := g.renderer.Render(&buf, "controller.go.tmpl", data); err != nil {
		return fmt.Errorf("failed to render controller for %s: %w", tag, err)
	}

	// Write using storage
	outputPath := filepath.Join(
		"generated", "adapters", "http", "controllers",
		strings.ToLower(tag)+".gen.go",
	)
	return g.storage.WriteFile(outputPath, buf.Bytes(), 0644)
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
