package openapi

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

const (
	boolTrueString = "true"
)

// Parser wraps an OpenAPI document and provides parsing utilities
type Parser struct {
	doc      *v3.Document
	basePath string
}

// NewParser creates a new Parser instance
func NewParser() *Parser {
	return &Parser{}
}

// Parse reads and parses an OpenAPI specification from a file path.
// It automatically processes any x-include-* extensions to merge in referenced specs.
func (p *Parser) Parse(path string) (*v3.Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	p.basePath = filepath.Dir(path)
	return p.ParseBytes(data)
}

// ParseBytes parses an OpenAPI specification from bytes and returns the document
func (p *Parser) ParseBytes(data []byte) (*v3.Document, error) {
	config := &datamodel.DocumentConfiguration{
		AllowFileReferences:   true,
		AllowRemoteReferences: true,
		BasePath:              p.basePath,
	}

	doc, err := libopenapi.NewDocumentWithConfiguration(data, config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}

	v3Model, err := doc.BuildV3Model()
	if err != nil {
		return nil, fmt.Errorf("failed to build v3 model: %w", err)
	}

	p.doc = &v3Model.Model
	return &v3Model.Model, nil
}
