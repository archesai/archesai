package codegen

import (
	"bytes"
	"fmt"
)

// MainTemplateData holds the data for rendering the main.go template.
type MainTemplateData struct {
	ProjectName string
}

// MainGenerator generates the main.go entry point for the application.
type MainGenerator struct{}

// Name returns the generator name.
func (g *MainGenerator) Name() string { return "main" }

// Priority returns the generator priority.
func (g *MainGenerator) Priority() int { return PriorityNormal }

// Generate creates the main.go file for composition apps.
func (g *MainGenerator) Generate(ctx *GeneratorContext) error {
	// Only generate for composition apps (apps that compose internal packages)
	composedPkgs := ctx.ComposedPackages()
	if len(composedPkgs) == 0 {
		return nil
	}

	data := &MainTemplateData{
		ProjectName: ctx.ProjectName,
	}

	var buf bytes.Buffer
	if err := ctx.Renderer.Render(&buf, "main.go.tmpl", data); err != nil {
		return fmt.Errorf("failed to render main.go.tmpl: %w", err)
	}

	return ctx.Storage.WriteFile("main.gen.go", buf.Bytes(), 0644)
}
