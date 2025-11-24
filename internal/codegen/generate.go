// Package codegen provides unified code generation from OpenAPI schemas with x-codegen extensions.
package codegen

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/archesai/archesai/internal/parsers"
	"github.com/archesai/archesai/pkg/storage"
)

// Generator handles code generation from OpenAPI specifications
type Generator struct {
	renderer *Renderer
	parser   *parsers.OpenAPIParser
	storage  storage.Storage
}

// NewGenerator creates a new code generator instance
func NewGenerator(outputDir string) *Generator {
	return &Generator{
		renderer: nil,
		parser:   parsers.NewOpenAPIParser(),
		storage:  storage.NewDiskStorage(outputDir),
	}
}

// WithStorage sets a custom storage implementation (useful for testing)
func (g *Generator) WithStorage(s storage.Storage) *Generator {
	g.storage = s
	return g
}

// GetStorage returns the current storage implementation
func (g *Generator) GetStorage() storage.Storage {
	return g.storage
}

// Initialize sets up the generator with templates and renderer
func (g *Generator) Initialize() error {
	templates, err := LoadTemplates()
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	// Wire templates directly to the Renderer
	g.renderer = NewRenderer(templates)

	return nil
}

// GenerateAPI is the main generation function that orchestrates all code generation
func (g *Generator) GenerateAPI(specPath string) error {
	totalStart := time.Now()

	slog.Info(
		"Parsing OpenAPI",
		slog.String("path", specPath),
	)
	_, err := g.parser.ParseFile(specPath)
	if err != nil {
		return fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	specDef, err := g.parser.Extract()
	if err != nil {
		return fmt.Errorf("failed to extraact definitions from openapi schema: %w", err)
	}

	slog.Info("Initialized state",
		slog.Int("schemas", len(specDef.Schemas)),
		slog.Int("operations", len(specDef.Operations)),
		"duration", time.Since(totalStart))

	eg, _ := errgroup.WithContext(context.Background())

	// Track timing for each generator
	type generatorTiming struct {
		name     string
		duration time.Duration
	}
	timings := make(chan generatorTiming, 10)

	eg.Go(func() error {
		start := time.Now()
		if err := g.GenerateSchemas(specDef.Schemas); err != nil {
			return fmt.Errorf("failed to generate models: %w", err)
		}
		timings <- generatorTiming{"GenerateSchemas", time.Since(start)}
		return nil
	})

	eg.Go(func() error {
		start := time.Now()
		if err := g.GenerateRepositories(specDef.Schemas); err != nil {
			return fmt.Errorf("failed to generate repositories: %w", err)
		}
		timings <- generatorTiming{"GenerateRepositories", time.Since(start)}
		return nil
	})

	eg.Go(func() error {
		start := time.Now()
		if err := g.GenerateEvents(specDef.Schemas); err != nil {
			return fmt.Errorf("failed to generate events: %w", err)
		}
		timings <- generatorTiming{"GenerateEvents", time.Since(start)}
		return nil
	})

	eg.Go(func() error {
		start := time.Now()
		if err := g.GenerateCommandQueryHandlers(specDef.Operations); err != nil {
			return fmt.Errorf("failed to generate handlers: %w", err)
		}
		timings <- generatorTiming{"GenerateCommandQueryHandlers", time.Since(start)}
		return nil
	})

	eg.Go(func() error {
		start := time.Now()
		if err := g.GenerateControllers(specDef.Operations); err != nil {
			return fmt.Errorf("failed to generate handlers: %w", err)
		}
		timings <- generatorTiming{"GenerateControllers", time.Since(start)}
		return nil
	})

	eg.Go(func() error {
		start := time.Now()
		if err := g.GenerateHCL(specDef.Schemas); err != nil {
			return fmt.Errorf("failed to generate HCL schema: %w", err)
		}
		hclDuration := time.Since(start)
		timings <- generatorTiming{"GenerateHCL", hclDuration}

		sqlcStart := time.Now()
		if err := g.GenerateSQLC(); err != nil {
			return fmt.Errorf("failed to generate SQLC files: %w", err)
		}
		timings <- generatorTiming{"GenerateSQLC", time.Since(sqlcStart)}
		return nil
	})

	eg.Go(func() error {
		start := time.Now()
		if err := g.GenerateJSClient(specPath); err != nil {
			return fmt.Errorf("failed to generate JavaScript client: %w", err)
		}
		timings <- generatorTiming{"GenerateJSClient", time.Since(start)}
		return nil
	})

	eg.Go(func() error {
		start := time.Now()
		if err := g.GenerateBootstrap(specDef.Schemas, specDef.Operations); err != nil {
			return fmt.Errorf("failed to generate bootstrap files: %w", err)
		}
		timings <- generatorTiming{"GenerateBootstrap", time.Since(start)}
		return nil
	})

	// Wait for all parallel generators to complete
	if err := eg.Wait(); err != nil {
		return err
	}

	// Close timings channel and collect results
	close(timings)
	for timing := range timings {
		slog.Info("Generator completed",
			slog.String("name", timing.name),
			slog.Duration("duration", timing.duration))
	}

	totalDuration := time.Since(totalStart)
	slog.Info("Code generation completed successfully",
		slog.Duration("total_duration", totalDuration))

	return nil
}
