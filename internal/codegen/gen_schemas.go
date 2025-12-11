package codegen

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/archesai/archesai/internal/schema"
	"github.com/archesai/archesai/internal/spec"
)

// GroupSchemas is the generator group for model schemas.
const GroupSchemas = "schemas"

const (
	genSchemas = "schemas"
	tmplSchema = "schema.go.tmpl"
)

// SchemasTemplateData holds data for schema template rendering.
type SchemasTemplateData struct {
	Package string
	Schema  *schema.Schema   // Single schema (multi-file mode)
	Schemas []*schema.Schema // All schemas (single-file mode)
}

// IsSingleMode returns true if generating in single file mode.
func (d *SchemasTemplateData) IsSingleMode() bool {
	return len(d.Schemas) > 0
}

// GetSchemas returns the schemas to render.
// In single mode, returns all schemas. In multi mode, returns just the one schema.
func (d *SchemasTemplateData) GetSchemas() []*schema.Schema {
	if d.IsSingleMode() {
		return d.Schemas
	}
	if d.Schema != nil {
		return []*schema.Schema{d.Schema}
	}
	return nil
}

// generateSchemas generates model files for each schema.
func (c *Codegen) generateSchemas(s *spec.Spec) error {
	internalContext := s.InternalContext()

	// Collect non-internal, non-request/response schemas
	var schemas []*schema.Schema
	for _, sch := range s.Schemas {
		if sch.IsInternal(internalContext) {
			continue
		}
		if isRequestResponseSchema(sch.Title) {
			continue
		}
		schemas = append(schemas, sch)
	}

	// Sort schemas by title for consistent output
	sort.Slice(schemas, func(i, j int) bool {
		return schemas[i].Title < schemas[j].Title
	})

	if c.isSingleStyle() {
		// Single file mode: all schemas in one file
		pkg := filepath.Base(c.storage.BaseDir())
		data := &SchemasTemplateData{
			Package: pkg,
			Schemas: schemas,
		}
		return c.RenderToFile(tmplSchema, "schemas.gen.go", data)
	}

	// Multi file mode: one file per schema
	for _, sch := range schemas {
		path := filepath.Join("schemas", strings.ToLower(sch.Title)+".gen.go")
		data := &SchemasTemplateData{
			Package: "schemas",
			Schema:  sch,
		}

		if err := c.RenderToFile(tmplSchema, path, data); err != nil {
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
