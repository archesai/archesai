package codegen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/archesai/archesai/internal/parsers"
	"github.com/archesai/archesai/pkg/storage"
)

// Priority levels for generators.
//
// Generator Priority Guidelines:
//
//   - PriorityFirst (0): Generators that create foundational files.
//     GoModGenerator creates go.mod and must run first.
//
//   - PriorityNormal (100): Independent generators that can run in parallel.
//     Includes: schemas, handlers, controllers, routes, postgres, sqlite,
//     repositories, client, main, app, container, bootstrap handlers.
//
//   - PriorityLast (200): Generators that depend on PriorityNormal outputs.
//     HCLGenerator needs all entity schemas to be defined.
//
//   - PriorityFinal (300): Generators that depend on PriorityLast outputs.
//     SQLCGenerator needs HCL migrations to generate type-safe queries.
const (
	PriorityFirst  = 0
	PriorityNormal = 100
	PriorityLast   = 200
	PriorityFinal  = 300
)

// Generator defines the interface for code generators.
type Generator interface {
	Name() string
	Priority() int
	Generate(ctx *GeneratorContext) error
}

// GeneratorContext provides shared context and dependencies for generators.
type GeneratorContext struct {
	SpecDef     *parsers.SpecDef
	SpecPath    string
	Renderer    *Renderer
	Storage     storage.Storage
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
func (ctx *GeneratorContext) OwnOperations() []parsers.OperationDef {
	internalContext := ctx.InternalContext()
	var operations []parsers.OperationDef
	for _, op := range ctx.SpecDef.Operations {
		if op.XInternal == "" || op.XInternal == internalContext {
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
	for _, op := range ctx.SpecDef.Operations {
		if op.XInternal != "" && op.XInternal != internalContext {
			pkgMap[op.XInternal] = true
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
// Filters by XCodegenSchemaType == "entity" and excludes internal schemas.
func (ctx *GeneratorContext) OwnEntitySchemas() []*parsers.SchemaDef {
	internalContext := ctx.InternalContext()
	var entities []*parsers.SchemaDef
	for _, schema := range ctx.SpecDef.Schemas {
		if schema.XCodegenSchemaType != parsers.XCodegenSchemaTypeEntity {
			continue
		}
		if schema.IsInternal(internalContext) {
			continue
		}
		entities = append(entities, schema)
	}
	return entities
}
