package openapi

import (
	"fmt"
	"os"
)

// RenderFormat specifies the output format for rendering OpenAPI documents.
type RenderFormat string

// Render format constants.
const (
	RenderFormatYAML RenderFormat = "yaml"
	RenderFormatJSON RenderFormat = "json"
)

// RenderOptions configures document rendering.
type RenderOptions struct {
	Format RenderFormat
}

// RenderDocument renders the parsed document in the specified format.
// Parse must be called before RenderDocument to populate the internal document.
func (p *Parser) RenderDocument(format RenderFormat) ([]byte, error) {
	if p.doc == nil {
		return nil, fmt.Errorf("document not parsed; call Parse first")
	}
	switch format {
	case RenderFormatYAML:
		return p.doc.Render()
	case RenderFormatJSON:
		return p.doc.RenderJSON("  ")
	default:
		return nil, fmt.Errorf("unsupported render format: %s", format)
	}
}

// RenderToFile renders the parsed document and writes it to the specified path.
func (p *Parser) RenderToFile(outputPath string) (string, error) {
	output, err := p.RenderDocument(RenderFormatYAML)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(outputPath, output, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// Resolve pathItems for compatibility
	if err := resolvePathItems(outputPath); err != nil {
		return "", fmt.Errorf("failed to resolve pathItems: %w", err)
	}

	return outputPath, nil
}
