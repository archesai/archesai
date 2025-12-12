package generators

import (
	"path/filepath"

	"github.com/archesai/archesai/internal/spec"
)

// GenerateAppBootstrap generates the app bootstrap file.
func GenerateAppBootstrap(ctx *GeneratorContext) error {
	actx := BuildAppContext(ctx)
	if actx.ShouldSkip() {
		return nil
	}

	// Bundle the OpenAPI spec
	var openAPISpec string
	var apiTitle string
	if ctx.SpecPath != "" {
		doc, err := spec.NewOpenAPIDocumentWithIncludes(ctx.SpecPath, ctx.Spec.EnabledIncludes)
		if err == nil {
			bundler := spec.NewBundler(doc)
			specBytes, err := bundler.BundleToYAML()
			if err == nil {
				openAPISpec = string(specBytes)
			}
			title, _, _ := doc.Info()
			apiTitle = title
		}
	}

	path := filepath.Join("app", "bootstrap.gen.go")
	data := &AppBootstrapTemplateData{
		ProjectName:      ctx.ProjectName,
		InternalPackages: actx.InternalPackages,
		OpenAPISpec:      openAPISpec,
		APITitle:         apiTitle,
	}

	return ctx.RenderToFile("app_bootstrap.go.tmpl", path, data)
}

// GenerateAppContainer generates the app container file.
func GenerateAppContainer(ctx *GeneratorContext) error {
	actx := BuildAppContext(ctx)
	if actx.ShouldSkip() {
		return nil
	}

	path := filepath.Join("app", "container.gen.go")
	data := &AppContainerTemplateData{
		Entities:         actx.Entities,
		ProjectName:      ctx.ProjectName,
		InternalPackages: actx.InternalPackagesWithEntities,
	}

	return ctx.RenderToFile("app_container.go.tmpl", path, data)
}

// GenerateAppHandlers generates the app handlers file.
func GenerateAppHandlers(ctx *GeneratorContext) error {
	actx := BuildAppContext(ctx)
	if actx.ShouldSkip() {
		return nil
	}

	path := filepath.Join("app", "handlers.gen.go")
	data := &AppHandlersTemplateData{
		Operations:        actx.Operations,
		Repositories:      actx.Repositories,
		ProjectName:       ctx.ProjectName,
		NeedsPublisher:    actx.NeedsPublisher,
		HasCustomHandlers: HasCustomHandlers(actx.Operations),
		InternalPackages:  actx.InternalPackages,
	}

	return ctx.RenderToFile("handlers.go.tmpl", path, data)
}

// GenerateAppInfrastructure generates the app infrastructure file.
func GenerateAppInfrastructure(ctx *GeneratorContext) error {
	actx := BuildAppContext(ctx)
	if actx.ShouldSkip() {
		return nil
	}

	path := filepath.Join("app", "infrastructure.gen.go")
	data := &AppHandlersTemplateData{
		Operations:       actx.Operations,
		Repositories:     actx.Repositories,
		ProjectName:      ctx.ProjectName,
		NeedsPublisher:   actx.NeedsPublisher,
		InternalPackages: actx.InternalPackages,
	}

	return ctx.RenderToFile("infrastructure.go.tmpl", path, data)
}

// GenerateAppRoutes generates the app routes file.
func GenerateAppRoutes(ctx *GeneratorContext) error {
	actx := BuildAppContext(ctx)
	if actx.ShouldSkip() {
		return nil
	}

	path := filepath.Join("app", "routes.gen.go")
	data := &AppRoutesTemplateData{
		Operations:       actx.Operations,
		ProjectName:      ctx.ProjectName,
		InternalPackages: actx.InternalPackages,
	}

	return ctx.RenderToFile("routes.go.tmpl", path, data)
}
