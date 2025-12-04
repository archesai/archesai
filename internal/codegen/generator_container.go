package codegen

import (
	"fmt"
	"path/filepath"
	"sort"
)

// ContainerTemplateData holds the data for rendering the container template.
type ContainerTemplateData struct {
	InternalPackages []InternalPackage
	ProjectName      string
}

// ContainerGenerator generates dependency injection container code.
type ContainerGenerator struct{}

// Name returns the generator name.
func (g *ContainerGenerator) Name() string { return "container" }

// Priority returns the generator priority.
func (g *ContainerGenerator) Priority() int { return PriorityNormal }

// Generate creates the container code for the application.
func (g *ContainerGenerator) Generate(ctx *GeneratorContext) error {
	// Only generate for composition apps (apps that compose internal packages)
	composedPkgs := ctx.ComposedPackages()
	if len(composedPkgs) == 0 {
		return nil
	}

	var internalPackages []InternalPackage
	for _, pkgName := range composedPkgs {
		internalPackages = append(internalPackages, InternalPackage{
			Name:       pkgName,
			Alias:      pkgName,
			ImportPath: InternalPackageImportPath(pkgName),
		})
	}
	sort.Slice(internalPackages, func(i, j int) bool {
		return internalPackages[i].Name < internalPackages[j].Name
	})

	data := &ContainerTemplateData{
		InternalPackages: internalPackages,
		ProjectName:      ctx.ProjectName,
	}

	outputPath := filepath.Join("bootstrap", "container.gen.go")
	if err := ctx.RenderToFile("container.go.tmpl", outputPath, data); err != nil {
		return fmt.Errorf("failed to generate container: %w", err)
	}
	return nil
}
