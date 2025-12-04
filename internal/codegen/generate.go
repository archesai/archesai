// Package codegen provides unified code generation from OpenAPI schemas with x-codegen extensions.
package codegen

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"sort"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/archesai/archesai/internal/parsers"
	"github.com/archesai/archesai/pkg/storage"
)

// ProgressCallback is called when a generator starts or completes.
type ProgressCallback func(event ProgressEvent)

// ProgressEvent represents a progress update from the codegen.
type ProgressEvent struct {
	Type          ProgressEventType
	GeneratorName string
	TotalCount    int
	CurrentIndex  int
	Error         error
}

// ProgressEventType indicates the type of progress event.
type ProgressEventType int

// Progress event types for codegen callbacks.
const (
	ProgressEventStart          ProgressEventType = iota // Generation started
	ProgressEventGeneratorStart                          // Individual generator started
	ProgressEventGeneratorDone                           // Individual generator completed
	ProgressEventDone                                    // All generation completed
	ProgressEventError                                   // Error occurred
)

// Codegen handles code generation from OpenAPI specifications.
type Codegen struct {
	renderer         *Renderer
	parser           *parsers.OpenAPIParser
	storage          storage.Storage
	generators       []Generator
	onlyFilter       map[string]bool
	progressCallback ProgressCallback
}

// NewCodegen creates a new code generator instance.
// By default, it uses NewDefaultIncludeMerger to register all standard includes.
func NewCodegen(outputDir string) *Codegen {
	merger := NewDefaultIncludeMerger()
	return &Codegen{
		renderer:   nil,
		parser:     parsers.NewOpenAPIParser().WithIncludeMerger(merger),
		storage:    storage.NewDiskStorage(outputDir),
		generators: DefaultGenerators(),
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
	c.onlyFilter = make(map[string]bool)
	for _, name := range splitByComma(only) {
		if name == "" {
			continue
		}
		c.onlyFilter[name] = true
	}
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

	return c.runGenerators(ctx)
}

// emitProgress sends a progress event if a callback is registered.
func (c *Codegen) emitProgress(event ProgressEvent) {
	if c.progressCallback != nil {
		c.progressCallback(event)
	}
}

// shouldRun returns true if the generator should run based on --only filter.
func (c *Codegen) shouldRun(name string) bool {
	if len(c.onlyFilter) == 0 {
		return true
	}
	return c.onlyFilter[name]
}

// runGenerators executes all generators respecting their priorities.
// Generators with the same priority run in parallel.
func (c *Codegen) runGenerators(ctx *GeneratorContext) error {
	// Filter generators based on --only flag
	var filtered []Generator
	for _, g := range c.generators {
		if c.shouldRun(g.Name()) {
			filtered = append(filtered, g)
		}
	}

	if len(filtered) == 0 {
		return nil
	}

	c.emitProgress(ProgressEvent{
		Type:       ProgressEventStart,
		TotalCount: len(filtered),
	})

	// Group generators by priority
	groups := c.groupByPriority(filtered)

	// Execute each priority group
	completedCount := 0
	for _, group := range groups {
		if err := c.runGroup(ctx, group, &completedCount, len(filtered)); err != nil {
			c.emitProgress(ProgressEvent{
				Type:  ProgressEventError,
				Error: err,
			})
			return err
		}
	}

	c.emitProgress(ProgressEvent{
		Type:       ProgressEventDone,
		TotalCount: len(filtered),
	})

	return nil
}

// priorityGroup holds generators with the same priority.
type priorityGroup struct {
	priority   int
	generators []Generator
}

// groupByPriority groups generators by their priority level.
func (c *Codegen) groupByPriority(generators []Generator) []priorityGroup {
	// Map priority to generators
	groupMap := make(map[int][]Generator)
	for _, g := range generators {
		p := g.Priority()
		groupMap[p] = append(groupMap[p], g)
	}

	// Convert to slice and sort by priority
	var groups []priorityGroup
	for priority, gens := range groupMap {
		groups = append(groups, priorityGroup{
			priority:   priority,
			generators: gens,
		})
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].priority < groups[j].priority
	})

	return groups
}

// runGroup executes a group of generators in parallel.
func (c *Codegen) runGroup(
	ctx *GeneratorContext,
	group priorityGroup,
	completedCount *int,
	totalCount int,
) error {
	eg, _ := errgroup.WithContext(context.Background())

	for _, g := range group.generators {
		g := g // capture for goroutine
		c.emitProgress(ProgressEvent{
			Type:          ProgressEventGeneratorStart,
			GeneratorName: g.Name(),
			CurrentIndex:  *completedCount,
			TotalCount:    totalCount,
		})

		eg.Go(func() error {
			start := time.Now()
			if err := g.Generate(ctx); err != nil {
				return fmt.Errorf("%s: %w", g.Name(), err)
			}
			slog.Debug("Generator completed",
				slog.String("name", g.Name()),
				slog.Duration("duration", time.Since(start)))

			c.emitProgress(ProgressEvent{
				Type:          ProgressEventGeneratorDone,
				GeneratorName: g.Name(),
				CurrentIndex:  *completedCount,
				TotalCount:    totalCount,
			})
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	*completedCount += len(group.generators)
	return nil
}

// splitByComma splits a string by comma and trims whitespace.
func splitByComma(s string) []string {
	var result []string
	current := ""
	for _, c := range s {
		if c == ',' {
			trimmed := trimSpace(current)
			if trimmed != "" {
				result = append(result, trimmed)
			}
			current = ""
		} else {
			current += string(c)
		}
	}
	trimmed := trimSpace(current)
	if trimmed != "" {
		result = append(result, trimmed)
	}
	return result
}

// trimSpace trims leading and trailing whitespace.
func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

// DefaultGenerators returns all standard generators.
func DefaultGenerators() []Generator {
	return []Generator{
		&GoModGenerator{},
		&SchemasGenerator{},
		&RepositoriesGenerator{},
		&PostgresGenerator{},
		&SQLiteGenerator{},
		&HandlersGenerator{},
		&ControllersGenerator{},
		&HCLGenerator{},
		&SQLCGenerator{},
		&ClientGenerator{},
		&MainGenerator{},
		&AppGenerator{},
		&RoutesGenerator{},
		&BootstrapHandlersGenerator{},
		&ContainerGenerator{},
	}
}
