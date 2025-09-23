// Package parsers provides utilities for parsing and manipulating OpenAPI documents
package parsers

import (
	"context"
	"fmt"
	"os"

	"github.com/speakeasy-api/openapi/openapi"
)

// Parser handles parsing of OpenAPI documents
type Parser struct {
	OpenAPI  *OpenAPISchema // OpenAPI operations handler
	warnings []string
}

// NewParser creates a new parser
func NewParser() *Parser {
	return &Parser{
		warnings: []string{},
	}
}

// Parse parses an OpenAPI specification file
func (p *Parser) Parse(specPath string) (*OpenAPISchema, []string, error) {
	ctx := context.Background()

	f, err := os.Open(specPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	doc, validationErrs, err := openapi.Unmarshal(ctx, f)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal OpenAPI document: %w", err)
	}

	if len(validationErrs) > 0 {
		for _, vErr := range validationErrs {
			p.warnings = append(p.warnings, vErr.Error())
		}
	}

	// Resolve all references in the document
	resolveValidationErrs, resolveErrs := doc.ResolveAllReferences(ctx, openapi.ResolveAllOptions{
		OpenAPILocation: specPath,
	})
	if resolveErrs != nil {
		return nil, nil, fmt.Errorf("failed to resolve references: %w", resolveErrs)
	}
	if len(resolveValidationErrs) > 0 {
		for _, vErr := range resolveValidationErrs {
			p.warnings = append(p.warnings, vErr.Error())
		}
	}

	return &OpenAPISchema{
		OpenAPI:  doc,
		FilePath: specPath,
	}, nil, nil

}

// GetWarnings returns any warnings accumulated during parsing
func (p *Parser) GetWarnings() []string {
	return p.warnings
}
