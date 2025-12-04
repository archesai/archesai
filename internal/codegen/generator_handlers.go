package codegen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/archesai/archesai/internal/parsers"
)

// HandlerTemplateData holds the data for rendering handler templates.
type HandlerTemplateData struct {
	Operation   *parsers.OperationDef
	ProjectName string
	NeedsUUID   bool
}

// HandlerStubTemplateData holds the data for rendering custom handler stub templates.
type HandlerStubTemplateData struct {
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

	needsUUID := op.HasBearerAuth() || op.HasCookieAuth()
	for _, p := range op.GetPathParams() {
		if p.GoType == "uuid.UUID" {
			needsUUID = true
		}
	}

	data := &HandlerTemplateData{
		Operation:   &op,
		ProjectName: ctx.ProjectName,
		NeedsUUID:   needsUUID,
	}

	var buf bytes.Buffer
	if err := ctx.Renderer.Render(&buf, "handler.gen.go.tmpl", data); err != nil {
		return fmt.Errorf("failed to render handler for %s: %w", op.ID, err)
	}

	return ctx.Storage.WriteFile(outputPath, buf.Bytes(), 0644)
}

func generateCustomHandlerStub(ctx *GeneratorContext, op parsers.OperationDef) error {
	fileName := parsers.SnakeCase(op.ID) + ".impl.go"
	outputPath := filepath.Join("application", fileName)

	fullPath := filepath.Join(ctx.Storage.BaseDir(), outputPath)
	if _, err := os.Stat(fullPath); err == nil {
		return nil
	}

	needsUUID := op.HasBearerAuth() || op.HasCookieAuth()
	for _, p := range op.GetPathParams() {
		if p.GoType == "uuid.UUID" {
			needsUUID = true
		}
	}

	data := &HandlerStubTemplateData{
		Operation:   &op,
		ProjectName: ctx.ProjectName,
		NeedsUUID:   needsUUID,
	}

	var buf bytes.Buffer
	if err := ctx.Renderer.Render(&buf, "handler_stub.go.tmpl", data); err != nil {
		return fmt.Errorf("failed to render handler stub for %s: %w", op.ID, err)
	}

	return ctx.Storage.WriteFile(outputPath, buf.Bytes(), 0644)
}
