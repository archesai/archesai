package openapi

import (
	"fmt"
	"os"

	"github.com/pb33f/libopenapi/bundler"
)

// Bundle bundles an OpenAPI specification with external references into a single document.
func (p *Parser) Bundle(outputPath string, orvalFix bool) (string, error) {
	bundled, err := bundler.BundleDocumentComposed(p.doc, &bundler.BundleCompositionConfig{
		Delimiter: "__",
	})
	if err != nil {
		return "", fmt.Errorf("failed to bundle spec: %w", err)
	}

	// Remove duplicate component entries created by composed bundling.
	// The bundler creates both "Foo: $ref: #/components/.../Foo__type" and "Foo__type: {...}".
	// We keep only the actual definitions (with suffix) and rename them to remove the suffix.
	bundled, err = cleanupComposedBundle(bundled)
	if err != nil {
		return "", fmt.Errorf("failed to cleanup bundled spec: %w", err)
	}

	// Write bundled output
	err = os.WriteFile(outputPath, bundled, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write bundled file: %w", err)
	}

	// Resolve pathItems if orval fix is enabled
	if orvalFix {
		if err := resolvePathItems(outputPath); err != nil {
			return "", fmt.Errorf("failed to resolve pathItems: %w", err)
		}
	}

	return outputPath, nil
}
