package generators

import "fmt"

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

	if err := ctx.RenderToFile("main.go.tmpl", "main.gen.go", data); err != nil {
		return fmt.Errorf("failed to generate main.go: %w", err)
	}
	return nil
}
