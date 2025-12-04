package codegen

import (
	"fmt"
	"path/filepath"

	"github.com/archesai/archesai/internal/parsers"
)

// ApplicationTemplateData holds the data for rendering application handler templates.
type ApplicationTemplateData struct {
	Operation         *parsers.OperationDef
	ProjectName       string
	NeedsUUID         bool
	NeedsServerModels bool
}

// ApplicationStubTemplateData holds the data for rendering custom handler stub templates.
type ApplicationStubTemplateData struct {
	Operation   *parsers.OperationDef
	ProjectName string
	NeedsUUID   bool
}

// HandlersGenerator generates handler code for API operations.
type HandlersGenerator struct{}

// Name returns the generator name.
func (g *HandlersGenerator) Name() string { return "application" }

// Priority returns the generator priority.
func (g *HandlersGenerator) Priority() int { return PriorityNormal }

// Generate creates handler code for each API operation.
func (g *HandlersGenerator) Generate(ctx *GeneratorContext) error {
	internalContext := ctx.InternalContext()
	for _, op := range ctx.SpecDef.Operations {
		if op.IsInternal(internalContext) {
			continue
		}

		if err := generateHandlerGenFile(ctx, op); err != nil {
			return fmt.Errorf("failed to generate handler for %s: %w", op.ID, err)
		}

		if op.XCodegenCustomHandler {
			if err := generateCustomHandlerStub(ctx, op); err != nil {
				return fmt.Errorf("failed to generate handler stub for %s: %w", op.ID, err)
			}
		}
	}
	return nil
}

func generateHandlerGenFile(ctx *GeneratorContext, op parsers.OperationDef) error {
	fileName := parsers.SnakeCase(op.ID) + ".gen.go"
	outputPath := filepath.Join("application", fileName)

	data := &ApplicationTemplateData{
		Operation:         &op,
		ProjectName:       ctx.ProjectName,
		NeedsUUID:         op.NeedsUUID(),
		NeedsServerModels: op.NeedsServerModels(),
	}

	if err := ctx.RenderToFile("handler.gen.go.tmpl", outputPath, data); err != nil {
		return fmt.Errorf("failed to generate handler for %s: %w", op.ID, err)
	}
	return nil
}

func generateCustomHandlerStub(ctx *GeneratorContext, op parsers.OperationDef) error {
	fileName := parsers.SnakeCase(op.ID) + ".impl.go"
	outputPath := filepath.Join("application", fileName)

	data := &ApplicationStubTemplateData{
		Operation:   &op,
		ProjectName: ctx.ProjectName,
		NeedsUUID:   op.NeedsUUID(),
	}

	if err := ctx.RenderToFileIfNotExists("handler_stub.go.tmpl", outputPath, data); err != nil {
		return fmt.Errorf("failed to generate handler stub for %s: %w", op.ID, err)
	}
	return nil
}
