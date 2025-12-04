// Package codegen provides unified code generation from OpenAPI schemas with x-codegen extensions.
package codegen

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/archesai/archesai/internal/parsers"
	"github.com/archesai/archesai/pkg/storage"
)

// Codegen handles code generation from OpenAPI specifications.
type Codegen struct {
	renderer         *Renderer
	parser           *parsers.OpenAPIParser
	storage          storage.Storage
	onlyFilter       string
	progressCallback ProgressCallback
}

// NewCodegen creates a new code generator instance.
func NewCodegen(outputDir string) *Codegen {
	return &Codegen{
		renderer: nil,
		parser:   parsers.NewOpenAPIParser(),
		storage:  storage.NewDiskStorage(outputDir),
	}
}

// WithStorage sets a custom storage implementation (useful for testing).
func (c *Codegen) WithStorage(s storage.Storage) *Codegen {
	c.storage = s
	return c
}

// WithLinting enables strict OpenAPI linting that blocks on any violations.
func (c *Codegen) WithLinting() *Codegen {
	c.parser.WithLinting()
	return c
}

// WithOnly sets which generators to run (comma-separated names).
func (c *Codegen) WithOnly(only string) *Codegen {
	c.onlyFilter = only
	return c
}

// WithProgress sets a callback for progress updates during generation.
func (c *Codegen) WithProgress(callback ProgressCallback) *Codegen {
	c.progressCallback = callback
	return c
}

// GetStorage returns the current storage implementation.
func (c *Codegen) GetStorage() storage.Storage {
	return c.storage
}

// Initialize sets up the generator with templates and renderer.
func (c *Codegen) Initialize() error {
	templates, err := LoadTemplates()
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	c.renderer = NewRenderer(templates)
	return nil
}

// GenerateAPI is the main generation function that orchestrates all code generation.
func (c *Codegen) GenerateAPI(specPath string) error {
	totalStart := time.Now()

	absSpecPath, err := filepath.Abs(specPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path of spec: %w", err)
	}

	slog.Debug("Parsing OpenAPI", slog.String("path", absSpecPath))

	_, err = c.parser.ParseFile(absSpecPath)
	if err != nil {
		return fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	specDef, err := c.parser.Extract()
	if err != nil {
		return fmt.Errorf("failed to extract definitions from openapi schema: %w", err)
	}

	slog.Debug("Initialized state",
		slog.Int("schemas", len(specDef.Schemas)),
		slog.Int("operations", len(specDef.Operations)),
		slog.String("project", specDef.ProjectName),
		"duration", time.Since(totalStart))

	ctx := &GeneratorContext{
		SpecDef:     specDef,
		SpecPath:    absSpecPath,
		Renderer:    c.renderer,
		Storage:     c.storage,
		ProjectName: specDef.ProjectName,
	}

	orchestrator := NewOrchestrator(DefaultGenerators()...)
	if c.onlyFilter != "" {
		orchestrator.WithOnly(c.onlyFilter)
	}
	if c.progressCallback != nil {
		orchestrator.WithProgress(c.progressCallback)
	}

	if err := orchestrator.Run(ctx); err != nil {
		return err
	}

	return nil
}
