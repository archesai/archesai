package spec

import (
	"slices"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// Spec represents the entire specification definition.
type Spec struct {
	Operations      []Operation  // All operations in the spec
	Schemas         []*Schema    // All schemas defined in the spec
	ProjectName     string       // Project name from x-project-name extension
	EnabledIncludes []string     // Names of enabled x-include-* extensions
	Document        *v3.Document // The underlying OpenAPI document

}

// HasInclude returns true if the named include is enabled (e.g., "auth", "executors")
func (s *Spec) HasInclude(name string) bool {
	return slices.Contains(s.EnabledIncludes, name)
}
