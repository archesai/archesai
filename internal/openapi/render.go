package openapi

import (
	"fmt"
)

type RenderFormat string

const (
	RenderFormatYAML RenderFormat = "yaml"
	RenderFormatJSON RenderFormat = "json"
)

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

// ParseAndRender parses an OpenAPI specification from a file and renders it in the specified format.
// This is a convenience method that combines Parse and RenderDocument.
func (p *Parser) ParseAndRender(path string, format RenderFormat) ([]byte, error) {
	if _, err := p.Parse(path); err != nil {
		return nil, fmt.Errorf("failed to parse spec: %w", err)
	}
	return p.RenderDocument(format)
}
