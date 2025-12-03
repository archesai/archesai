package codegen

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/archesai/archesai/internal/parsers"
)

// BootstrapHandlersTemplateData holds the data for rendering the handlers bootstrap template.
type BootstrapHandlersTemplateData struct {
	Operations     []parsers.OperationDef
	Repositories   []string
	ProjectName    string
	NeedsPublisher bool
}

// BootstrapHandlersGenerator generates handler initialization code.
type BootstrapHandlersGenerator struct{}

// Name returns the generator name.
func (g *BootstrapHandlersGenerator) Name() string { return "bootstrap_handlers" }

// Priority returns the generator priority.
func (g *BootstrapHandlersGenerator) Priority() int { return PriorityNormal }

// Generate creates handler bootstrap code for internal packages.
func (g *BootstrapHandlersGenerator) Generate(ctx *GeneratorContext) error {
	// Only generate for internal packages (packages with their own operations)
	operations := ctx.OwnOperations()
	if len(operations) == 0 {
		return nil
	}

	sort.Slice(operations, func(i, j int) bool {
		return operations[i].ID < operations[j].ID
	})

	// Collect unique repositories (skip custom handlers which don't need repos)
	repoMap := make(map[string]bool)
	for _, op := range operations {
		if op.XCodegenCustomHandler {
			continue
		}
		repoName := op.XCodegenRepository
		if repoName == "" {
			repoName = op.Tag
		}
		if repoName != "" {
			repoMap[repoName] = true
		}
	}

	var repositories []string
	for repo := range repoMap {
		repositories = append(repositories, repo)
	}
	sort.Strings(repositories)

	// Check if we need a publisher (any non-GET, non-custom operations)
	needsPublisher := false
	for _, op := range operations {
		if !op.XCodegenCustomHandler && op.Method != "GET" {
			needsPublisher = true
			break
		}
	}

	data := &BootstrapHandlersTemplateData{
		Operations:     operations,
		Repositories:   repositories,
		ProjectName:    ctx.ProjectName,
		NeedsPublisher: needsPublisher,
	}

	var buf bytes.Buffer
	if err := ctx.Renderer.Render(&buf, "handlers.go.tmpl", data); err != nil {
		return fmt.Errorf("failed to render handlers.go.tmpl: %w", err)
	}

	return ctx.Storage.WriteFile(filepath.Join("bootstrap", "handlers.gen.go"), buf.Bytes(), 0644)
}
