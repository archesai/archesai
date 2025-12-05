package openapi

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/index"

	"github.com/archesai/archesai/pkg/logger"
)

const (
	boolTrueString = "true"
)

// Parser wraps an OpenAPI document and provides parsing utilities
type Parser struct {
	doc      *v3.Document
	basePath string
	localFS  fs.FS // Optional custom filesystem for resolving references
}

// NewParser creates a new Parser instance
func NewParser() *Parser {
	return &Parser{}
}

// SetLocalFS sets a custom filesystem for resolving local file references.
// This is useful when combining multiple filesystems (e.g., local + embedded).
func (p *Parser) SetLocalFS(fsys fs.FS) {
	p.localFS = fsys
}

// Parse reads and parses an OpenAPI specification from a file path.
// It automatically processes any x-include-* extensions to merge in referenced specs,
// then bundles all external references into a single document.
func (p *Parser) Parse(path string) (*v3.Document, error) {
	p.basePath = filepath.Dir(path)

	// Build merger and composite filesystem for reference resolution
	merger := NewDefaultIncludeMerger()

	// Merge spec content (paths, tags, components, security from includes)
	data, _, err := merger.MergeSpec(path)
	if err != nil {
		return nil, fmt.Errorf("failed to merge spec: %w", err)
	}

	// Build composite filesystem so $ref can resolve from embedded includes
	if p.localFS == nil {
		compositeFS, _, err := merger.BuildCompositeFS(path)
		if err != nil {
			return nil, fmt.Errorf("failed to build composite filesystem: %w", err)
		}
		p.localFS = compositeFS
	}

	// First parse to resolve all references
	if _, err := p.ParseBytes(data); err != nil {
		return nil, err
	}

	// Bundle to inline all external references
	bundled, err := p.Bundle()
	if err != nil {
		return nil, fmt.Errorf("failed to bundle spec: %w", err)
	}

	// Re-parse the bundled spec for a clean document
	p.localFS = nil
	return p.ParseBytes(bundled)
}

// ParseBytes parses an OpenAPI specification from bytes and returns the document
func (p *Parser) ParseBytes(data []byte) (*v3.Document, error) {
	config := &datamodel.DocumentConfiguration{
		AllowFileReferences:     true,
		AllowRemoteReferences:   true,
		BasePath:                p.basePath,
		ExtractRefsSequentially: true,
		Logger:                  logger.NewDiscard(),
	}

	// Use custom filesystem if provided - wrap it in a LocalFS for libopenapi
	if p.localFS != nil {
		localFSConfig := &index.LocalFSConfig{
			BaseDirectory: p.basePath,
			DirFS:         p.localFS,
		}
		localFS, err := index.NewLocalFSWithConfig(localFSConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create local filesystem: %w", err)
		}
		config.LocalFS = localFS
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
