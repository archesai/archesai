package codegen

import (
	"path/filepath"
	"strings"

	"github.com/archesai/archesai/internal/spec"
	"github.com/archesai/archesai/internal/strutil"
)

// GroupOperations is the generator group for operation types (use cases).
const GroupOperations = "operations"

// GroupHTTP is the generator group for HTTP handlers.
const GroupHTTP = "http"

const (
	genOperations = "operations"
	genHTTP       = "http"
)

// OperationsView wraps a spec for rendering with package metadata.
type OperationsView struct {
	*spec.Spec
	Package          string
	OperationsPrefix string // "operations." or "" (for single mode)
	SchemasPrefix    string // "schemas." or "" (for single mode)
	operation        *spec.Operation
}

// OwnOperations returns operations for this view.
func (v *OperationsView) OwnOperations() []spec.Operation {
	if v.operation != nil {
		return []spec.Operation{*v.operation}
	}
	return v.Spec.OwnOperations()
}

// IsSingleMode returns true if generating in single file mode.
func (v *OperationsView) IsSingleMode() bool {
	return v.OperationsPrefix == ""
}

// ModelType converts a Go type to use the correct schemas prefix for this view.
// In single mode (SchemasPrefix=""), it strips existing "schemas." prefix.
// In multi mode (SchemasPrefix="schemas."), it ensures the prefix is present.
func (v *OperationsView) ModelType(goType string) string {
	return modelType(goType, v.SchemasPrefix)
}

// modelType converts a Go type to use the given schemas prefix.
func modelType(goType string, prefix string) string {
	// Handle slice types
	if strings.HasPrefix(goType, "[]") {
		inner := strings.TrimPrefix(goType, "[]")
		return "[]" + modelType(inner, prefix)
	}

	// Handle pointer types
	if strings.HasPrefix(goType, "*") {
		inner := strings.TrimPrefix(goType, "*")
		return "*" + modelType(inner, prefix)
	}

	// Strip existing schemas. prefix if present
	goType = strings.TrimPrefix(goType, "schemas.")

	// Keep external package qualifications (serverschemas., etc.)
	if strings.Contains(goType, ".") {
		return goType
	}

	// Primitive types don't need prefix
	switch goType {
	case "string", "int", "int32", "int64", "float32", "float64", "bool", "any", "time.Time":
		return goType
	}

	// uuid.UUID is a qualified type
	if goType == "uuid.UUID" {
		return goType
	}

	// Schema reference - prefix with schemas prefix (may be empty in single mode)
	return prefix + goType
}

// generateOperations generates operation type files (interfaces + input/output) for each operation.
func (c *Codegen) generateOperations(s *spec.Spec) error {
	if c.isSingleStyle() {
		pkg := filepath.Base(c.storage.BaseDir())
		view := &OperationsView{Spec: s, Package: pkg, SchemasPrefix: ""}
		return c.RenderToFile(genOperations+".go.tmpl", "operations.gen.go", view)
	}

	// Multi file mode: one file per operation
	for _, op := range s.OwnOperations() {
		path := filepath.Join("operations", strutil.SnakeCase(op.ID)+".gen.go")
		view := &OperationsView{
			Spec:          s,
			Package:       "operations",
			SchemasPrefix: "schemas.",
			operation:     &op,
		}
		if err := c.RenderToFile(genOperations+".go.tmpl", path, view); err != nil {
			return err
		}
	}

	return nil
}

// generateHTTP generates HTTP handler files for each operation.
func (c *Codegen) generateHTTP(s *spec.Spec) error {
	if c.isSingleStyle() {
		pkg := filepath.Base(c.storage.BaseDir())
		view := &OperationsView{Spec: s, Package: pkg, OperationsPrefix: "", SchemasPrefix: ""}
		return c.RenderToFile(genHTTP+".go.tmpl", "http.gen.go", view)
	}

	// Multi file mode: one file per handler
	for _, op := range s.OwnOperations() {
		path := filepath.Join("http", strutil.SnakeCase(op.ID)+".gen.go")
		view := &OperationsView{
			Spec:             s,
			Package:          "http",
			OperationsPrefix: "operations.",
			SchemasPrefix:    "schemas.",
			operation:        &op,
		}
		if err := c.RenderToFile(genHTTP+".go.tmpl", path, view); err != nil {
			return err
		}
	}

	return nil
}
