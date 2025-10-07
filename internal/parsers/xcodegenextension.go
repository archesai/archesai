package parsers

import (
	"fmt"

	"go.yaml.in/yaml/v4"
)

// XCodegenParser handles parsing of x-codegen extensions
type XCodegenParser struct {
}

// NewXCodegenParser creates a new x-codegen extension parser
func NewXCodegenParser() *XCodegenParser {
	return &XCodegenParser{}
}

// ParseExtension parses an x-codegen extension from any type (typically from Extensions map)
func (p *XCodegenParser) ParseExtension(
	ext *yaml.Node,
	schemaName string,
) (*XCodegenExtension, error) {
	if ext == nil {
		return nil, nil
	}

	var xcodegen XCodegenExtension
	if err := ext.Decode(&xcodegen); err != nil {
		return nil, fmt.Errorf(
			"failed to decode x-codegen extension for schema '%s': %w",
			schemaName,
			err,
		)
	}
	return &xcodegen, nil
}
