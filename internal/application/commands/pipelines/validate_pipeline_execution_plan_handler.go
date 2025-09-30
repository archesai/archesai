package pipelines

import (
	"context"
	"fmt"

	"github.com/archesai/archesai/internal/core/entities"
	"github.com/archesai/archesai/internal/core/events"
	"github.com/archesai/archesai/internal/core/repositories"
)

// ValidationResult represents the result of pipeline validation.
type ValidationResult struct {
	Valid  bool     `json:"valid"`
	Issues []string `json:"issues,omitempty"`
}

// ValidatePipelineExecutionPlanCommandHandler handles the validate pipeline execution plan command.
type ValidatePipelineExecutionPlanCommandHandler struct {
	pipelineRepo repositories.PipelineRepository
	publisher    events.Publisher
}

// NewValidatePipelineExecutionPlanCommandHandler creates a new validate pipeline execution plan command handler.
func NewValidatePipelineExecutionPlanCommandHandler(
	pipelineRepo repositories.PipelineRepository,
	publisher events.Publisher,
) *ValidatePipelineExecutionPlanCommandHandler {
	return &ValidatePipelineExecutionPlanCommandHandler{
		pipelineRepo: pipelineRepo,
		publisher:    publisher,
	}
}

// Handle executes the validate pipeline execution plan command.
func (h *ValidatePipelineExecutionPlanCommandHandler) Handle(
	ctx context.Context,
	cmd *ValidatePipelineExecutionPlanCommand,
) (*ValidationResult, error) {
	// Verify pipeline exists
	_, err := h.pipelineRepo.Get(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline: %w", err)
	}

	// TODO: Implement PipelineStepRepository to fetch steps
	// For now, use an empty slice
	// steps, _, err := h.pipelineStepRepo.List(ctx, 1000, 0)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to list pipeline steps: %w", err)
	// }
	steps := []*entities.PipelineStep{}

	result := &ValidationResult{
		Valid:  true,
		Issues: []string{},
	}

	// Basic validation checks
	if len(steps) == 0 {
		result.Issues = append(result.Issues, "Pipeline has no steps")
	}

	// Build dependency graph and check for cycles
	if hasCycles := h.detectCycles(steps); hasCycles {
		result.Valid = false
		result.Issues = append(result.Issues, "Pipeline has circular dependencies")
	}

	// Check for orphaned dependencies
	stepIDs := make(map[string]bool)
	for _, step := range steps {
		stepIDs[step.ID.String()] = true
	}

	for _, step := range steps {
		for _, depID := range step.Dependencies {
			if !stepIDs[depID.String()] {
				result.Valid = false
				result.Issues = append(
					result.Issues,
					fmt.Sprintf("Step %s depends on non-existent step %s", step.Name, depID),
				)
			}
		}
	}

	return result, nil
}

// detectCycles checks for circular dependencies in the pipeline steps using DFS.
func (h *ValidatePipelineExecutionPlanCommandHandler) detectCycles(
	steps []*entities.PipelineStep,
) bool {
	// Build adjacency list
	graph := make(map[string][]string)
	for _, step := range steps {
		stepID := step.ID.String()
		if _, exists := graph[stepID]; !exists {
			graph[stepID] = []string{}
		}
		for _, dep := range step.Dependencies {
			graph[stepID] = append(graph[stepID], dep.String())
		}
	}

	// Track visited nodes and recursion stack
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	// DFS to detect cycles
	var dfs func(string) bool
	dfs = func(node string) bool {
		visited[node] = true
		recStack[node] = true

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				if dfs(neighbor) {
					return true
				}
			} else if recStack[neighbor] {
				return true // Cycle detected
			}
		}

		recStack[node] = false
		return false
	}

	// Check all nodes
	for stepID := range graph {
		if !visited[stepID] {
			if dfs(stepID) {
				return true
			}
		}
	}

	return false
}
