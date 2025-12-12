package generators

import (
	"path/filepath"

	"github.com/archesai/archesai/internal/strutil"
)

// GenerateHandlers generates handler files for each operation.
func GenerateHandlers(ctx *GeneratorContext) error {
	internalContext := ctx.InternalContext()

	for _, op := range ctx.Spec.Operations {
		if op.IsInternal(internalContext) {
			continue
		}

		path := filepath.Join("handlers", strutil.SnakeCase(op.ID)+".gen.go")
		data := &ApplicationTemplateData{
			Operation:   &op,
			ProjectName: ctx.ProjectName,
		}

		if err := ctx.RenderToFile("application_handler.go.tmpl", path, data); err != nil {
			return err
		}
	}

	return nil
}

// GenerateHandlerStubs generates stub files for custom handlers in the implement package.
func GenerateHandlerStubs(ctx *GeneratorContext) error {
	internalContext := ctx.InternalContext()

	for _, op := range ctx.Spec.Operations {
		if op.IsInternal(internalContext) {
			continue
		}

		if !op.CustomHandler {
			continue
		}

		path := filepath.Join("implement", strutil.SnakeCase(op.ID)+".gen.go")

		data := &ApplicationStubTemplateData{
			Operation:   &op,
			ProjectName: ctx.ProjectName,
		}

		if err := ctx.RenderToFile("handler_stub.go.tmpl", path, data); err != nil {
			return err
		}
	}

	return nil
}

// GenerateRoutes generates route/controller files for each operation.
func GenerateRoutes(ctx *GeneratorContext) error {
	internalContext := ctx.InternalContext()

	for _, op := range ctx.Spec.Operations {
		if op.IsInternal(internalContext) {
			continue
		}

		path := filepath.Join("routes", strutil.SnakeCase(op.ID)+".gen.go")
		data := &RouteTemplateData{
			Operation:   &op,
			ProjectName: ctx.ProjectName,
		}

		if err := ctx.RenderToFile("controller.go.tmpl", path, data); err != nil {
			return err
		}
	}

	return nil
}
