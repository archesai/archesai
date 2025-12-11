package generators

import (
	"fmt"
	"path/filepath"

	"github.com/archesai/archesai/internal/spec"
	"github.com/archesai/archesai/internal/strutil"
)

// ControllerTemplateData holds the data for rendering controller templates.
type ControllerTemplateData struct {
	Operation   *spec.Operation
	ProjectName string
}

// ControllersGenerator generates controller code for API operations.
type ControllersGenerator struct{}

// Name returns the generator name.
func (g *ControllersGenerator) Name() string { return "routes" }

// Priority returns the generator priority.
func (g *ControllersGenerator) Priority() int { return PriorityNormal }

// Generate creates controller code for each API operation.
func (g *ControllersGenerator) Generate(ctx *GeneratorContext) error {
	internalContext := ctx.InternalContext()
	for _, op := range ctx.Spec.Operations {
		if op.IsInternal(internalContext) {
			continue
		}

		outputPath := filepath.Join(
			"routes",
			fmt.Sprintf("%s.gen.go", strutil.SnakeCase(op.ID)),
		)
		data := &ControllerTemplateData{
			Operation:   &op,
			ProjectName: ctx.ProjectName,
		}

		if err := ctx.RenderToFile("controller.go.tmpl", outputPath, data); err != nil {
			return fmt.Errorf("failed to generate controller for %s: %w", op.ID, err)
		}
	}
	return nil
}
