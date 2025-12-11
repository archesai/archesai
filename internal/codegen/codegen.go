// Package codegen provides unified code generation from OpenAPI schemas with x-codegen extensions.
package codegen

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/archesai/archesai/internal/located"
	"github.com/archesai/archesai/internal/spec"
	"github.com/archesai/archesai/internal/templates"
	"github.com/archesai/archesai/pkg/config/schemas"
	"github.com/archesai/archesai/pkg/storage"
)

// Priority constants for generator execution order.
// Lower values run first. Generators with the same priority run in parallel.
const (
	PriorityFirst  = 0   // Runs first (e.g., go.mod, package.json)
	PriorityNormal = 100 // Normal priority (most generators)
	PriorityLast   = 200 // Runs after normal generators (e.g., HCL schema)
	PriorityFinal  = 300 // Runs last (e.g., SQLC after migrations)
)

// Codegen handles code generation from OpenAPI specifications.
type Codegen struct {
	renderer *templates.Renderer
	storage  storage.Storage
	progress ProgressCallback
	cfg      schemas.ConfigGeneration
}

// New creates a new Codegen instance with disk storage at the given output directory.
func New(outputDir string) (*Codegen, error) {
	tmpl, err := templates.LoadTemplates()
	if err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	return &Codegen{
		renderer: templates.NewRenderer(tmpl),
		storage:  storage.NewTrackedStorage(storage.NewDiskStorage(outputDir)),
	}, nil
}

// WithStorage sets a custom storage implementation (useful for testing).
func (c *Codegen) WithStorage(s storage.Storage) *Codegen {
	c.storage = s
	return c
}

// WithProgress sets a callback for progress updates during generation.
func (c *Codegen) WithProgress(callback ProgressCallback) *Codegen {
	c.progress = callback
	return c
}

// GetStorage returns the current storage implementation.
func (c *Codegen) GetStorage() storage.Storage {
	return c.storage
}

// Generate runs all generators for the given spec and configuration.
func (c *Codegen) Generate(s *located.Located[spec.Spec], cfg schemas.ConfigGeneration) error {
	totalStart := time.Now()

	// Store config for generators to access
	c.cfg = cfg

	// Ensure spec path is absolute
	if s.Path != "" {
		absPath, err := filepath.Abs(s.Path)
		if err != nil {
			return fmt.Errorf("failed to get absolute path of spec: %w", err)
		}
		s.Path = absPath
	}

	slog.Debug("Starting code generation",
		slog.Int("schemas", len(s.Value.Schemas)),
		slog.Int("operations", len(s.Value.Operations)),
		slog.String("project", s.Value.ProjectName))

	if err := c.runGenerators(s, cfg); err != nil {
		return err
	}

	slog.Debug("Code generation complete", slog.Duration("duration", time.Since(totalStart)))
	return nil
}

// isSingleStyle returns true if the generation style is "single".
func (c *Codegen) isSingleStyle() bool {
	return c.cfg.Style != nil && *c.cfg.Style == schemas.ConfigGenerationStyleSingle
}

// RenderToFile renders a template and writes it to the specified path.
func (c *Codegen) RenderToFile(templateName, outputPath string, data any) error {
	var buf bytes.Buffer
	if err := c.renderer.Render(&buf, templateName, data); err != nil {
		return fmt.Errorf("failed to render %s: %w", templateName, err)
	}
	if err := c.storage.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", outputPath, err)
	}
	return nil
}

// RenderToFileIfNotExists renders a template only if the file doesn't exist.
// Useful for stub files that should not be overwritten.
func (c *Codegen) RenderToFileIfNotExists(templateName, outputPath string, data any) error {
	fullPath := filepath.Join(c.storage.BaseDir(), outputPath)
	if _, err := os.Stat(fullPath); err == nil {
		return nil // File exists, skip
	}
	return c.RenderToFile(templateName, outputPath, data)
}

// RenderTSXToFile renders a TSX template to the specified path.
// TSX templates use [[ ]] delimiters and are looked up with tsx/ prefix.
func (c *Codegen) RenderTSXToFile(templateName, outputPath string, data any) error {
	var buf bytes.Buffer
	if err := c.renderer.Render(&buf, "tsx/"+templateName, data); err != nil {
		return fmt.Errorf("failed to render %s: %w", templateName, err)
	}
	if err := c.storage.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", outputPath, err)
	}
	return nil
}

// generator defines an internal generator with name, group, priority, and method.
type generator struct {
	name     string
	group    string
	priority int
	run      func() error
}

// isGroupEnabled checks if a generator group is enabled in the config.
// If Groups is nil, all groups are enabled by default.
func isGroupEnabled(groups *schemas.ConfigGenerationGroups, group string) bool {
	if groups == nil {
		return true
	}
	switch group {
	case GroupModule:
		return groups.Module
	case GroupSchemas:
		return groups.Schemas
	case GroupOperations:
		return groups.Operations
	case GroupHTTP:
		return groups.HTTP
	case GroupApp:
		return groups.Wire
	case GroupPostgres:
		return groups.Postgres
	case GroupSQLite:
		return groups.Sqlite
	case GroupWeb:
		return groups.Web
	default:
		return true
	}
}

// runGenerators executes all generators respecting their priorities.
func (c *Codegen) runGenerators(s *located.Located[spec.Spec], cfg schemas.ConfigGeneration) error {
	// Build generator list
	generators := c.defaultGenerators(s)

	// Filter generators based on config.Groups
	// If Groups is nil, all groups are enabled by default
	var filtered []generator
	for _, gen := range generators {
		if isGroupEnabled(cfg.Groups, gen.group) {
			filtered = append(filtered, gen)
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
		if err := c.runGroup(group, &completedCount, len(filtered)); err != nil {
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

// emitProgress sends a progress event if a callback is registered.
func (c *Codegen) emitProgress(event ProgressEvent) {
	if c.progress != nil {
		c.progress(event)
	}
}

// priorityGroup holds generators with the same priority.
type priorityGroup struct {
	priority   int
	generators []generator
}

// groupByPriority groups generators by their priority level.
func (c *Codegen) groupByPriority(gens []generator) []priorityGroup {
	groupMap := make(map[int][]generator)
	for _, gen := range gens {
		groupMap[gen.priority] = append(groupMap[gen.priority], gen)
	}

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
func (c *Codegen) runGroup(group priorityGroup, completedCount *int, totalCount int) error {
	eg, _ := errgroup.WithContext(context.Background())

	for _, gen := range group.generators {
		gen := gen // capture for goroutine
		c.emitProgress(ProgressEvent{
			Type:          ProgressEventGeneratorStart,
			GeneratorName: gen.name,
			CurrentIndex:  *completedCount,
			TotalCount:    totalCount,
		})

		eg.Go(func() error {
			start := time.Now()
			if err := gen.run(); err != nil {
				return fmt.Errorf("%s: %w", gen.name, err)
			}
			slog.Debug("Generator completed",
				slog.String("name", gen.name),
				slog.Duration("duration", time.Since(start)))

			c.emitProgress(ProgressEvent{
				Type:          ProgressEventGeneratorDone,
				GeneratorName: gen.name,
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
