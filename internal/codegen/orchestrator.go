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

	"github.com/archesai/archesai/internal/codegen/generators"
	"github.com/archesai/archesai/internal/openapi"
	"github.com/archesai/archesai/internal/strutil"
	"github.com/archesai/archesai/internal/templates"
	"github.com/archesai/archesai/pkg/storage"
)

// Orchestrator handles code generation from OpenAPI specifications.
type Orchestrator struct {
	renderer         *templates.Renderer
	storage          storage.Storage
	generators       []generators.Generator
	onlyFilter       map[string]bool
	progressCallback ProgressCallback
}

// NewOrchestrator creates a new code generator instance.
// By default, it uses NewDefaultIncludeMerger to register all standard includes.
func NewOrchestrator(outputDir string) *Orchestrator {
	return &Orchestrator{
		renderer:   nil,
		storage:    storage.NewDiskStorage(outputDir),
		generators: generators.DefaultGenerators(),
	}
}

// WithStorage sets a custom storage implementation (useful for testing).
func (o *Orchestrator) WithStorage(s storage.Storage) *Orchestrator {
	o.storage = s
	return o
}

// WithOnly sets which generators to run (comma-separated names).
func (o *Orchestrator) WithOnly(only string) *Orchestrator {
	o.onlyFilter = make(map[string]bool)
	split := strutil.SplitByComma(only)
	for _, name := range split {
		if name == "" {
			continue
		}
		o.onlyFilter[name] = true
	}
	return o
}

// WithProgress sets a callback for progress updates during generation.
func (o *Orchestrator) WithProgress(callback ProgressCallback) *Orchestrator {
	o.progressCallback = callback
	return o
}

// GetStorage returns the current storage implementation.
func (o *Orchestrator) GetStorage() storage.Storage {
	return o.storage
}

// Initialize sets up the generator with templates and renderer.
func (o *Orchestrator) Initialize() error {
	template, err := templates.LoadTemplates()
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	o.renderer = templates.NewRenderer(template)
	return nil
}

// Generate is the main generation function that orchestrates all code generation.
func (o *Orchestrator) Generate(specPath string) error {
	totalStart := time.Now()

	absSpecPath, err := filepath.Abs(specPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path of spec: %w", err)
	}

	slog.Debug("Parsing OpenAPI", slog.String("path", absSpecPath))

	parser := openapi.NewParser()
	_, err = parser.Parse(absSpecPath)
	if err != nil {
		return fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	spec, err := parser.ExtractSpec()
	if err != nil {
		return fmt.Errorf("failed to extract definitions from openapi schema: %w", err)
	}

	slog.Debug("Initialized state",
		slog.Int("schemas", len(spec.Schemas)),
		slog.Int("operations", len(spec.Operations)),
		slog.String("project", spec.ProjectName),
		"duration", time.Since(totalStart))

	ctx := &generators.GeneratorContext{
		Spec:        spec,
		SpecPath:    absSpecPath,
		Renderer:    o.renderer,
		Storage:     o.storage,
		ProjectName: spec.ProjectName,
	}

	return o.runGenerators(ctx)
}

// emitProgress sends a progress event if a callback is registered.
func (o *Orchestrator) emitProgress(event ProgressEvent) {
	if o.progressCallback != nil {
		o.progressCallback(event)
	}
}

// shouldRun returns true if the generator should run based on --only filter.
func (o *Orchestrator) shouldRun(name string) bool {
	if len(o.onlyFilter) == 0 {
		return true
	}
	return o.onlyFilter[name]
}

// runGenerators executes all generators respecting their priorities.
// Generators with the same priority run in parallel.
func (o *Orchestrator) runGenerators(ctx *generators.GeneratorContext) error {
	// Filter generators based on --only flag
	var filtered []generators.Generator
	for _, g := range o.generators {
		if o.shouldRun(g.Name()) {
			filtered = append(filtered, g)
		}
	}

	if len(filtered) == 0 {
		return nil
	}

	o.emitProgress(ProgressEvent{
		Type:       ProgressEventStart,
		TotalCount: len(filtered),
	})

	// Group generators by priority
	groups := o.groupByPriority(filtered)

	// Execute each priority group
	completedCount := 0
	for _, group := range groups {
		if err := o.runGroup(ctx, group, &completedCount, len(filtered)); err != nil {
			o.emitProgress(ProgressEvent{
				Type:  ProgressEventError,
				Error: err,
			})
			return err
		}
	}

	o.emitProgress(ProgressEvent{
		Type:       ProgressEventDone,
		TotalCount: len(filtered),
	})

	return nil
}

// priorityGroup holds generators with the same priority.
type priorityGroup struct {
	priority   int
	generators []generators.Generator
}

// groupByPriority groups generators by their priority level.
func (o *Orchestrator) groupByPriority(gens []generators.Generator) []priorityGroup {
	// Map priority to generators
	groupMap := make(map[int][]generators.Generator)
	for _, g := range gens {
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
func (o *Orchestrator) runGroup(
	ctx *generators.GeneratorContext,
	group priorityGroup,
	completedCount *int,
	totalCount int,
) error {
	eg, _ := errgroup.WithContext(context.Background())

	for _, g := range group.generators {
		g := g // capture for goroutine
		o.emitProgress(ProgressEvent{
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

			o.emitProgress(ProgressEvent{
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
