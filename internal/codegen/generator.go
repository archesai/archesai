package codegen

import (
	"strings"

	"github.com/archesai/archesai/internal/parsers"
	"github.com/archesai/archesai/pkg/storage"
)

// Priority levels for generators.
const (
	PriorityFirst  = 0
	PriorityNormal = 100
	PriorityLast   = 200
	PriorityFinal  = 300 // Runs after PriorityLast (e.g., sqlc needs migrations from hcl)
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
