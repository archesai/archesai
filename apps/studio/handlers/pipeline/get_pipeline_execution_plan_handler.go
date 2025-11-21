// Package pipeline provides command and query handlers for pipeline operations.
package pipeline

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	queries "github.com/archesai/archesai/apps/studio/generated/application/queries/pipeline"
	"github.com/archesai/archesai/apps/studio/generated/core/models"
	"github.com/archesai/archesai/apps/studio/generated/core/repositories"
)

// ExecutionLevel represents a level in the execution plan.
type ExecutionLevel struct {
	Level int         `json:"level"`
	Steps []uuid.UUID `json:"steps"`
}

// ExecutionPlan represents the complete execution plan for a pipeline.
type ExecutionPlan struct {
	PipelineID        uuid.UUID        `json:"pipelineID"`
	Levels            []ExecutionLevel `json:"levels"`
	TotalSteps        int              `json:"totalSteps"`
	IsValid           bool             `json:"isValid"`
	EstimatedDuration *int             `json:"estimatedDuration,omitempty"`
}

// GetPipelineExecutionPlanQueryHandler handles the get pipeline execution plan query.
type GetPipelineExecutionPlanQueryHandler struct {
	pipelineRepo repositories.PipelineRepository
}

// NewGetPipelineExecutionPlanQueryHandler creates a new get pipeline execution plan query handler.
func NewGetPipelineExecutionPlanQueryHandler(
	pipelineRepo repositories.PipelineRepository,
) *GetPipelineExecutionPlanQueryHandler {
	return &GetPipelineExecutionPlanQueryHandler{
		pipelineRepo: pipelineRepo,
	}
}

// Handle executes the get pipeline execution plan query.
func (h *GetPipelineExecutionPlanQueryHandler) Handle(
	ctx context.Context,
	query *queries.GetPipelineExecutionPlanQuery,
) (*ExecutionPlan, error) {
	// Verify pipeline exists
	pipeline, err := h.pipelineRepo.Get(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline: %w", err)
	}

	// TODO: Implement PipelineStepRepository to fetch steps
	// For now, use empty slice
	// steps, _, err := h.pipelineStepRepo.ListByPipeline(ctx, query.ID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to list pipeline steps: %w", err)
	// }
	steps := []*models.PipelineStep{}

	// Build execution plan based on dependencies
	levels := h.buildExecutionLevels(steps)

	// // Calculate estimated duration
	// estimatedDuration := h.calculateEstimatedDuration(steps, levels)

	// Check if plan is valid (no cycles)
	isValid := !h.hasCycles(steps)

	plan := &ExecutionPlan{
		PipelineID: pipeline.ID,
		Levels:     levels,
		TotalSteps: len(steps),
		IsValid:    isValid,
		// EstimatedDuration: &estimatedDuration,
	}

	return plan, nil
}

// buildExecutionLevels creates execution levels based on step dependencies.
func (h *GetPipelineExecutionPlanQueryHandler) buildExecutionLevels(
	steps []*models.PipelineStep,
) []ExecutionLevel {
	if len(steps) == 0 {
		return []ExecutionLevel{}
	}

	// 	// Build dependency map
	// 	dependsOn := make(map[uuid.UUID][]uuid.UUID)
	// 	for _, step := range steps {
	// 		dependsOn[step.ID] = step.Dependencies
	// 	}

	// 	// Calculate levels using topological sort
	// 	levels := []ExecutionLevel{}
	// 	processed := make(map[uuid.UUID]bool)
	// 	level := 0

	// 	for len(processed) < len(steps) {
	// 		currentLevel := ExecutionLevel{
	// 			Level: level,
	// 			Steps: []uuid.UUID{},
	// 		}

	// 		// Find steps that can run at this level
	// 		for _, step := range steps {
	// 			if processed[step.ID] {
	// 				continue
	// 			}

	// 			// Check if all dependencies are processed
	// 			canRun := true
	// 			for _, dep := range dependsOn[step.ID] {
	// 				if !processed[dep] {
	// 					canRun = false
	// 					break
	// 				}
	// 			}

	// 			if canRun {
	// 				currentLevel.Steps = append(currentLevel.Steps, step.ID)
	// 				processed[step.ID] = true
	// 			}
	// 		}

	// 		if len(currentLevel.Steps) == 0 {
	// 			// Cycle detected or no more steps can be processed
	// 			break
	// 		}

	// 		levels = append(levels, currentLevel)
	// 		level++
	// 	}

	// 	return levels
	// }

	// // calculateEstimatedDuration estimates total execution time.
	// func (h *GetPipelineExecutionPlanQueryHandler) calculateEstimatedDuration(
	// 	steps []*models.PipelineStep,
	// 	levels []ExecutionLevel,
	// ) int {
	// 	if len(steps) == 0 {
	// 		return 0
	// 	}

	// 	// Build step timeout map
	// 	stepTimeouts := make(map[uuid.UUID]int)
	// 	for _, step := range steps {
	// 		timeout := 3600 // Default timeout
	// 		if step.Timeout != nil {
	// 			timeout = *step.Timeout
	// 		}
	// 		stepTimeouts[step.ID] = timeout
	// 	}

	// 	// Calculate duration as sum of max timeout per level
	// 	totalDuration := 0
	// 	for _, level := range levels {
	// 		maxTimeout := 0
	// 		for _, stepID := range level.Steps {
	// 			if timeout, ok := stepTimeouts[stepID]; ok && timeout > maxTimeout {
	// 				maxTimeout = timeout
	// 			}
	// 		}
	// 		totalDuration += maxTimeout
	// 	}

	return []ExecutionLevel{}
}

// hasCycles checks for circular dependencies using DFS.
func (h *GetPipelineExecutionPlanQueryHandler) hasCycles(_ []*models.PipelineStep) bool {
	// Build adjacency list
	// graph := make(map[string][]string)
	// for _, step := range steps {
	// 	stepID := step.ID.String()
	// 	if _, exists := graph[stepID]; !exists {
	// 		graph[stepID] = []string{}
	// 	}
	// 	for _, dep := range step.Dependencies {
	// 		graph[stepID] = append(graph[stepID], dep.String())
	// 	}
	// }

	// // Track visited nodes and recursion stack
	// visited := make(map[string]bool)
	// recStack := make(map[string]bool)

	// // DFS to detect cycles
	// var dfs func(string) bool
	// dfs = func(node string) bool {
	// 	visited[node] = true
	// 	recStack[node] = true

	// 	for _, neighbor := range graph[node] {
	// 		if !visited[neighbor] {
	// 			if dfs(neighbor) {
	// 				return true
	// 			}
	// 		} else if recStack[neighbor] {
	// 			return true // Cycle detected
	// 		}
	// 	}

	// 	recStack[node] = false
	// 	return false
	// }

	// // Check all nodes
	// for stepID := range graph {
	// 	if !visited[stepID] {
	// 		if dfs(stepID) {
	// 			return true
	// 		}
	// 	}
	// }

	return false
}
