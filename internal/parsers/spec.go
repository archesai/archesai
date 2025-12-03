package parsers

import (
	"slices"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// SpecDef represents the entire specification definition.
type SpecDef struct {
	Operations      []OperationDef // All operations in the spec
	Schemas         []*SchemaDef   // All schemas defined in the spec
	Document        *v3.Document   // The underlying OpenAPI document
	ProjectName     string         // Project name from x-project-name extension
	EnabledIncludes []string       // Names of enabled x-include-* extensions
}

// HasInclude returns true if the named include is enabled (e.g., "auth", "executors")
func (s *SpecDef) HasInclude(name string) bool {
	return slices.Contains(s.EnabledIncludes, name)
}
