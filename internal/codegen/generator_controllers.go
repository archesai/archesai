package codegen

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/archesai/archesai/internal/parsers"
)

// HandlerControllerTemplateData holds the data for rendering controller templates.
type HandlerControllerTemplateData struct {
	Operation   *parsers.OperationDef
	ProjectName string
}

// ControllersGenerator generates controller code for API operations.
type ControllersGenerator struct{}

// Name returns the generator name.
func (g *ControllersGenerator) Name() string { return "controllers" }

// Priority returns the generator priority.
func (g *ControllersGenerator) Priority() int { return PriorityNormal }

// Generate creates controller code for each API operation.
func (g *ControllersGenerator) Generate(ctx *GeneratorContext) error {
	internalContext := ctx.InternalContext()
	for _, op := range ctx.SpecDef.Operations {
		if op.IsInternal(internalContext) {
			continue
		}

		outputPath := filepath.Join(
			"controllers",
			fmt.Sprintf("%s.gen.go", parsers.SnakeCase(op.ID)),
		)
		data := &HandlerControllerTemplateData{
			Operation:   &op,
			ProjectName: ctx.ProjectName,
		}

		var buf bytes.Buffer
		if err := ctx.Renderer.Render(&buf, "controller.go.tmpl", data); err != nil {
			return fmt.Errorf("failed to render handler controller for %s: %w", op.ID, err)
		}

		if err := ctx.Storage.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write handler controller for %s: %w", op.ID, err)
		}
	}
	return nil
}
