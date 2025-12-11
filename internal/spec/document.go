package spec

import (
	"fmt"
	"io/fs"

	"go.yaml.in/yaml/v4"

	"github.com/archesai/archesai/internal/ref"
)

// OpenAPIDocument wraps the raw YAML representation of an OpenAPI spec.
// This is the single source of truth - all operations query this directly.
type OpenAPIDocument struct {
	root     *yaml.Node        // Root document node (for bundling)
	doc      *RawDocument      // Structured document (for config access)
	fsys     fs.FS             // Filesystem for resolving $refs
	baseDir  string            // Base directory for relative refs
	resolver *ref.FileResolver // Shared resolver for file refs
}

// NewOpenAPIDocumentFromFS loads an OpenAPI specification from a filesystem.
// The filesystem should already include any necessary includes (via CompositeFS).
func NewOpenAPIDocumentFromFS(fsys fs.FS, filename string) (*OpenAPIDocument, error) {
	// Load the root document
	data, err := fs.ReadFile(fsys, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read OpenAPI document: %w", err)
	}

	// Parse into yaml.Node (for bundling)
	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI document: %w", err)
	}

	// Parse into RawDocument struct (for config access)
	var parsedDoc RawDocument
	if err := yaml.Unmarshal(data, &parsedDoc); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI document: %w", err)
	}

	return &OpenAPIDocument{
		root:     &root,
		doc:      &parsedDoc,
		baseDir:  ".",
		fsys:     fsys,
		resolver: ref.NewFileResolver(fsys, "."),
	}, nil
}

// Raw returns the parsed RawDocument for external access.
func (d *OpenAPIDocument) Raw() *RawDocument {
	return d.doc
}
