package codegen

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"golang.org/x/sync/errgroup"
)

// ProgressCallback is called when a generator starts or completes.
type ProgressCallback func(event ProgressEvent)

// ProgressEvent represents a progress update from the orchestrator.
type ProgressEvent struct {
	Type          ProgressEventType
	GeneratorName string
	TotalCount    int
	CurrentIndex  int
	Error         error
}

// ProgressEventType indicates the type of progress event.
type ProgressEventType int

// Progress event types for orchestrator callbacks.
const (
	ProgressEventStart          ProgressEventType = iota // Generation started
	ProgressEventGeneratorStart                          // Individual generator started
	ProgressEventGeneratorDone                           // Individual generator completed
	ProgressEventDone                                    // All generation completed
	ProgressEventError                                   // Error occurred
)

// Orchestrator manages and executes code generators.
type Orchestrator struct {
	generators       []Generator
	onlyFilter       map[string]bool
	progressCallback ProgressCallback
}

// NewOrchestrator creates a new orchestrator with the given generators.
func NewOrchestrator(generators ...Generator) *Orchestrator {
	return &Orchestrator{
		generators: generators,
	}
}

// WithOnly sets which generators to run (comma-separated names).
func (o *Orchestrator) WithOnly(only string) *Orchestrator {
	o.onlyFilter = make(map[string]bool)
	for _, name := range splitByComma(only) {
		if name == "" {
			continue
		}
		o.onlyFilter[name] = true
	}
	return o
}

// WithProgress sets a callback for progress updates.
func (o *Orchestrator) WithProgress(callback ProgressCallback) *Orchestrator {
	o.progressCallback = callback
	return o
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

// Run executes all generators respecting their priorities.
// Generators with the same priority run in parallel.
func (o *Orchestrator) Run(ctx *GeneratorContext) error {
	// Filter generators based on --only flag
	var filtered []Generator
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
	generators []Generator
}

// groupByPriority groups generators by their priority level.
func (o *Orchestrator) groupByPriority(generators []Generator) []priorityGroup {
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
func (o *Orchestrator) runGroup(
	ctx *GeneratorContext,
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
