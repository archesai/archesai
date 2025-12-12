package generators

import (
	"path/filepath"
	"strings"
)

// GenerateGoMod generates go.mod file.
func GenerateGoMod(ctx *GeneratorContext) error {
	path := "go.mod"
	if ctx.FileExists(path) {
		return nil
	}

	data := map[string]string{"ProjectName": ctx.ProjectName}
	return ctx.RenderToFile("go.mod.tmpl", path, data)
}

// GenerateMain generates main.go for composition apps.
func GenerateMain(ctx *GeneratorContext) error {
	// Only generate for composition apps
	if len(ctx.ComposedPackages()) == 0 {
		return nil
	}

	data := map[string]string{"ProjectName": ctx.ProjectName}
	return ctx.RenderToFile("main.go.tmpl", "main.gen.go", data)
}

// GenerateSchemas generates model files for each schema.
func GenerateSchemas(ctx *GeneratorContext) error {
	internalContext := ctx.InternalContext()

	for _, schema := range ctx.Spec.Schemas {
		if schema.IsInternal(internalContext) {
			continue
		}

		// Skip request/response schemas
		if isRequestResponseSchema(schema.Name) {
			continue
		}

		path := filepath.Join("models", strings.ToLower(schema.Name)+".gen.go")
		data := &SchemasTemplateData{
			Package: "models",
			Schema:  schema,
		}

		if err := ctx.RenderToFile("schema.go.tmpl", path, data); err != nil {
			return err
		}
	}

	return nil
}

func isRequestResponseSchema(name string) bool {
	suffixes := []string{"Request", "Response", "Input", "Output"}
	for _, suffix := range suffixes {
		if len(name) > len(suffix) && name[len(name)-len(suffix):] == suffix {
			return true
		}
	}
	return false
}
