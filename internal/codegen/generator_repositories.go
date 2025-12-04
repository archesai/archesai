package codegen

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/archesai/archesai/internal/parsers"
)

// RepositoriesTemplateData holds the data for rendering repository templates.
type RepositoriesTemplateData struct {
	Entity              *parsers.SchemaDef
	ProjectName         string
	ModelImportPath     string
	RepositoryInterface string
}

// RepositoriesGenerator generates repository interface code for entities.
type RepositoriesGenerator struct{}

// Name returns the generator name.
func (g *RepositoriesGenerator) Name() string { return "repositories" }

// Priority returns the generator priority.
func (g *RepositoriesGenerator) Priority() int { return PriorityNormal }

// Generate creates repository interface code for each entity schema.
func (g *RepositoriesGenerator) Generate(ctx *GeneratorContext) error {
	for _, schema := range ctx.OwnEntitySchemas() {
		modelImportPath, repositoryInterface := getRepositoryImportPaths(ctx, schema)
		data := &RepositoriesTemplateData{
			Entity:              schema,
			ProjectName:         ctx.ProjectName,
			ModelImportPath:     modelImportPath,
			RepositoryInterface: repositoryInterface,
		}

		outputPath := filepath.Join("repositories", strings.ToLower(schema.Name)+".gen.go")
		if err := ctx.RenderToFile("repository.go.tmpl", outputPath, data); err != nil {
			return fmt.Errorf(
				"failed to generate repository interface for %s: %w",
				schema.Name,
				err,
			)
		}
	}
	return nil
}

func getRepositoryImportPaths(
	ctx *GeneratorContext,
	schema *parsers.SchemaDef,
) (modelImportPath, repositoryInterface string) {
	internalContext := ctx.InternalContext()
	if schema.IsInternal(internalContext) && schema.XInternal != "" {
		modelImportPath = InternalPackageModelsPath(schema.XInternal)
		repositoryInterface = InternalPackageRepositoriesPath(schema.XInternal)
	} else {
		modelImportPath = ctx.ProjectName + "/models"
		repositoryInterface = ctx.ProjectName + "/repositories"
	}
	return
}
