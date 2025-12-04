package codegen

import (
	"fmt"
	"path/filepath"
)

// AppTemplateData holds the data for rendering the app bootstrap template.
type AppTemplateData struct {
	ProjectName string
}

// AppGenerator generates the app bootstrap code.
type AppGenerator struct{}

// Name returns the generator name.
func (g *AppGenerator) Name() string { return "app" }

// Priority returns the generator priority.
func (g *AppGenerator) Priority() int { return PriorityNormal }

// Generate creates the app bootstrap file for composition apps.
func (g *AppGenerator) Generate(ctx *GeneratorContext) error {
	// Only generate for composition apps (apps that compose internal packages)
	composedPkgs := ctx.ComposedPackages()
	if len(composedPkgs) == 0 {
		return nil
	}

	data := &AppTemplateData{
		ProjectName: ctx.ProjectName,
	}

	outputPath := filepath.Join("bootstrap", "app.gen.go")
	if err := ctx.RenderToFile("app.go.tmpl", outputPath, data); err != nil {
		return fmt.Errorf("failed to generate app bootstrap: %w", err)
	}
	return nil
}
