package codegen

import (
	"bytes"
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
	internalContext := ctx.InternalContext()
	for _, schema := range ctx.SpecDef.Schemas {
		if schema.XCodegenSchemaType != parsers.XCodegenSchemaTypeEntity {
			continue
		}
		if schema.IsInternal(internalContext) {
			continue
		}

		modelImportPath, repositoryInterface := getRepositoryImportPaths(
			ctx.ProjectName,
			internalContext,
			schema,
		)
		data := &RepositoriesTemplateData{
			Entity:              schema,
			ProjectName:         ctx.ProjectName,
			ModelImportPath:     modelImportPath,
			RepositoryInterface: repositoryInterface,
		}

		var buf bytes.Buffer
		if err := ctx.Renderer.Render(&buf, "repository.go.tmpl", data); err != nil {
			return fmt.Errorf("failed to render repository interface for %s: %w", schema.Name, err)
		}

		outputPath := filepath.Join("repositories", strings.ToLower(schema.Name)+".gen.go")
		if err := ctx.Storage.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write repository interface for %s: %w", schema.Name, err)
		}
	}
	return nil
}

func getRepositoryImportPaths(
	projectName, internalContext string,
	schema *parsers.SchemaDef,
) (modelImportPath, repositoryInterface string) {
	if schema.IsInternal(internalContext) && schema.XInternal != "" {
		modelImportPath = "github.com/archesai/archesai/pkg/" + schema.XInternal + "/models"
		repositoryInterface = "github.com/archesai/archesai/pkg/" + schema.XInternal + "/repositories"
	} else {
		modelImportPath = projectName + "/models"
		repositoryInterface = projectName + "/repositories"
	}
	return
}
