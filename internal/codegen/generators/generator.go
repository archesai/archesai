package generators

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/archesai/archesai/internal/spec"
	"github.com/archesai/archesai/internal/templates"
)

// GeneratorContext provides shared context and dependencies for generators.
type GeneratorContext struct {
	Spec        *spec.Spec
	SpecPath    string
	Renderer    *templates.Renderer
	Storage     Storage
	ProjectName string
}

// InternalContext returns the last segment of the project name.
func (ctx *GeneratorContext) InternalContext() string {
	if ctx.ProjectName == "" {
		return ""
	}
	parts := strings.Split(ctx.ProjectName, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

// OwnOperations returns operations that belong to this package.
// An operation belongs to this package if x-internal is empty or matches InternalContext.
func (ctx *GeneratorContext) OwnOperations() []spec.Operation {
	internalContext := ctx.InternalContext()
	var operations []spec.Operation
	for _, op := range ctx.Spec.Operations {
		if op.Internal == "" || op.Internal == internalContext {
			operations = append(operations, op)
		}
	}
	return operations
}

// ComposedPackages returns the unique x-internal package names from operations
// that belong to OTHER packages (not this one).
func (ctx *GeneratorContext) ComposedPackages() []string {
	internalContext := ctx.InternalContext()
	pkgMap := make(map[string]bool)
	for _, op := range ctx.Spec.Operations {
		if op.Internal != "" && op.Internal != internalContext {
			pkgMap[op.Internal] = true
		}
	}

	var packages []string
	for pkg := range pkgMap {
		packages = append(packages, pkg)
	}
	return packages
}

// InternalPackage represents an internal package for composition.
type InternalPackage struct {
	Name       string
	Alias      string
	ImportPath string
}

// RenderToFile renders a template and writes it to the specified path.
// This replaces the common pattern of creating a buffer, rendering, and writing.
func (ctx *GeneratorContext) RenderToFile(templateName, outputPath string, data any) error {
	var buf bytes.Buffer
	if err := ctx.Renderer.Render(&buf, templateName, data); err != nil {
		return fmt.Errorf("failed to render %s: %w", templateName, err)
	}
	if err := ctx.Storage.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", outputPath, err)
	}
	return nil
}

// RenderToFileIfNotExists renders a template only if the file doesn't exist.
// Useful for stub files that should not be overwritten.
func (ctx *GeneratorContext) RenderToFileIfNotExists(
	templateName, outputPath string,
	data any,
) error {
	fullPath := filepath.Join(ctx.Storage.BaseDir(), outputPath)
	if _, err := os.Stat(fullPath); err == nil {
		return nil // File exists, skip
	}
	return ctx.RenderToFile(templateName, outputPath, data)
}

// OwnEntitySchemas returns entity schemas that belong to this package.
// Filters by SchemaType == "entity" and excludes internal schemas.
func (ctx *GeneratorContext) OwnEntitySchemas() []*spec.Schema {
	internalContext := ctx.InternalContext()
	var entities []*spec.Schema
	for _, schema := range ctx.Spec.Schemas {
		if schema.XCodegenSchemaType != spec.SchemaTypeEntity {
			continue
		}
		if schema.IsInternal(internalContext) {
			continue
		}
		entities = append(entities, schema)
	}
	return entities
}

// AllEntitySchemas returns all entity schemas regardless of x-internal.
// Used by database generators that need to generate repositories
// for all entities including those from included packages.
func (ctx *GeneratorContext) AllEntitySchemas() []*spec.Schema {
	var entities []*spec.Schema
	for _, schema := range ctx.Spec.Schemas {
		if schema.XCodegenSchemaType != spec.SchemaTypeEntity {
			continue
		}
		entities = append(entities, schema)
	}
	return entities
}

// FileExists checks if a file already exists at the given path.
func (ctx *GeneratorContext) FileExists(path string) bool {
	exists, _ := ctx.Storage.Exists(path)
	return exists
}

// RenderTSXToFile renders a TSX template to the specified path.
// TSX templates use [[ ]] delimiters and are looked up with tsx/ prefix.
func (ctx *GeneratorContext) RenderTSXToFile(templateName, outputPath string, data any) error {
	var buf bytes.Buffer
	// Prefix with tsx/ to select TSX template collection in renderer
	if err := ctx.Renderer.Render(&buf, "tsx/"+templateName, data); err != nil {
		return fmt.Errorf("failed to render %s: %w", templateName, err)
	}
	if err := ctx.Storage.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", outputPath, err)
	}
	return nil
}
