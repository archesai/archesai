package parsers

import (
	"fmt"
	"slices"

	"github.com/speakeasy-api/openapi/extensions"
	"github.com/speakeasy-api/openapi/jsonschema/oas3"
)

// type XCodegenExtension struct {
// 	SchemaType string `yaml:"schemaType,omitempty"`
// 	Repository *struct {
// 		AdditionalMethods []struct {
// 			Name    string   `yaml:"name"`
// 			Params  []string `yaml:"params"`
// 			Returns string   `yaml:"returns"` // "single", "multiple", "error"
// 		} `yaml:"additional_methods,omitempty"`
// 	} `yaml:"repository,omitempty"`
// }

// XCodegenParser handles parsing of x-codegen extensions
type XCodegenParser struct {
}

// NewXCodegenParser creates a new x-codegen extension parser
func NewXCodegenParser() *XCodegenParser {
	return &XCodegenParser{}
}

// ParseExtension parses an x-codegen extension from a schema
func (p *XCodegenParser) ParseExtension(
	ext extensions.Extension,
	schemaName string,
) (*XCodegenExtension, error) {
	if ext == nil {
		return nil, nil
	}

	var xcodegen XCodegenExtension
	// ext is already *yaml.Node, so we can decode directly
	if err := ext.Decode(&xcodegen); err != nil {
		return nil, fmt.Errorf(
			"failed to decode x-codegen extension for schema '%s': %w",
			schemaName,
			err,
		)
	}

	return &xcodegen, nil
}

// ParseSchemaExtensions parses x-codegen extensions from a schema
func (p *XCodegenParser) ParseSchemaExtensions(
	schema *oas3.Schema,
	schemaName string,
) (*XCodegenExtension, error) {
	if schema == nil || schema.Extensions == nil {
		return nil, nil
	}

	ext := schema.Extensions.GetOrZero("x-codegen")
	if ext == nil {
		return nil, nil
	}

	return p.ParseExtension(ext, schemaName)
}

// Parse is a simpler helper that just parses extensions
func (p *XCodegenParser) Parse(extensions extensions.Extensions) *XCodegenExtension {
	ext := extensions.GetOrZero("x-codegen")
	if ext == nil {
		return nil
	}

	var xcodegen XCodegenExtension
	if err := ext.Decode(&xcodegen); err != nil {
		// Just return nil if we can't decode
		return nil
	}

	return &xcodegen
}

// contains checks if a string is in a slice
func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}
