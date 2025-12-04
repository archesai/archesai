package codegen

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/archesai/archesai/internal/parsers"
)

// RoutesTemplateData holds the data for rendering the routes template.
type RoutesTemplateData struct {
	Operations       []parsers.OperationDef
	ProjectName      string
	InternalPackages []InternalPackage
}

// RoutesGenerator generates HTTP route registration code.
type RoutesGenerator struct{}

// Name returns the generator name.
func (g *RoutesGenerator) Name() string { return "routes" }

// Priority returns the generator priority.
func (g *RoutesGenerator) Priority() int { return PriorityNormal }

// Generate creates route registration code for all API operations.
func (g *RoutesGenerator) Generate(ctx *GeneratorContext) error {
	data := g.buildTemplateData(ctx)
	if data == nil {
		return nil
	}

	outputPath := filepath.Join("bootstrap", "routes.gen.go")
	if err := ctx.RenderToFile("routes.go.tmpl", outputPath, data); err != nil {
		return fmt.Errorf("failed to generate routes: %w", err)
	}
	return nil
}

func (g *RoutesGenerator) buildTemplateData(ctx *GeneratorContext) *RoutesTemplateData {
	// Check for internal package (has own operations)
	operations := ctx.OwnOperations()
	if len(operations) > 0 {
		sort.Slice(operations, func(i, j int) bool {
			return operations[i].ID < operations[j].ID
		})
		return &RoutesTemplateData{
			Operations:  operations,
			ProjectName: ctx.ProjectName,
		}
	}

	// Check for composition app (composes other packages)
	composedPkgs := ctx.ComposedPackages()
	if len(composedPkgs) > 0 {
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
		return &RoutesTemplateData{
			ProjectName:      ctx.ProjectName,
			InternalPackages: internalPackages,
		}
	}

	return nil
}
