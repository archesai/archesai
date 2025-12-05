package openapi

import (
	"fmt"

	"github.com/pb33f/libopenapi/bundler"
)

// Bundle bundles an OpenAPI specification and returns the bundled bytes.
func (p *Parser) Bundle() ([]byte, error) {
	bundled, err := bundler.BundleDocumentComposed(p.doc, &bundler.BundleCompositionConfig{
		Delimiter: "__",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to bundle spec: %w", err)
	}

	// Remove duplicate component entries created by composed bundling.
	// The bundler creates both "Foo: $ref: #/components/.../Foo__type" and "Foo__type: {...}".
	// We keep only the actual definitions (with suffix) and rename them to remove the suffix.
	bundled, err = cleanupComposedBundle(bundled)
	if err != nil {
		return nil, fmt.Errorf("failed to cleanup bundled spec: %w", err)
	}

	return bundled, nil
}
