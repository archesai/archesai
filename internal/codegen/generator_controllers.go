package codegen

import (
	"fmt"
	"path/filepath"

	"github.com/archesai/archesai/internal/parsers"
)

// ControllerTemplateData holds the data for rendering controller templates.
type ControllerTemplateData struct {
	Operation         *parsers.OperationDef
	ProjectName       string
	NeedsServerModels bool
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
		data := &ControllerTemplateData{
			Operation:         &op,
			ProjectName:       ctx.ProjectName,
			NeedsServerModels: op.NeedsServerModels(),
		}

		if err := ctx.RenderToFile("controller.go.tmpl", outputPath, data); err != nil {
			return fmt.Errorf("failed to generate controller for %s: %w", op.ID, err)
		}
	}
	return nil
}
