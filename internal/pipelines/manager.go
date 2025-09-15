package pipelines

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// PipelineManager handles pipeline step operations with DAG support
type PipelineManager struct {
	pipelineRepository     Repository
	pipelineStepRepository PipelineStepRepository
	logger                 *slog.Logger
}

// NewPipelineManager creates a new pipeline manager
func NewPipelineManager(pipelineRepository Repository, pipelineStepRepository PipelineStepRepository, logger *slog.Logger) *PipelineManager {
	return &PipelineManager{
		pipelineRepository:     pipelineRepository,
		pipelineStepRepository: pipelineStepRepository,
		logger:                 logger,
	}
}

// CreatePipelineStep adds a new step to a pipeline
func (pm *PipelineManager) CreatePipelineStep(ctx context.Context, pipelineID, toolID uuid.UUID, config map[string]interface{}) (*PipelineStep, error) {
	step := &PipelineStep{
		Id:         uuid.New(),
		PipelineId: pipelineID,
		ToolId:     toolID,
		Config:     config,
	}

	// Store the step using the step repository
	createdStep, err := pm.pipelineStepRepository.CreatePipelineStep(ctx, step)
	if err != nil {
		return nil, fmt.Errorf("failed to create pipeline step: %w", err)
	}

	return createdStep, nil
}

// AddStepDependency creates a dependency between two steps
func (pm *PipelineManager) AddStepDependency(ctx context.Context, stepID, dependsOnID uuid.UUID) error {
	// Create the dependency
	err := pm.pipelineStepRepository.CreateStepDependency(ctx, stepID, dependsOnID)
	if err != nil {
		return fmt.Errorf("failed to create dependency: %w", err)
	}

	// Validate no cycles were created by getting all steps for validation
	// For now, we'll skip cycle validation until we have pipeline ID context
	// TODO: Add proper cycle validation when we have pipeline context
	pm.logger.Info("Created step dependency", "stepId", stepID, "dependsOnId", dependsOnID)

	return nil
}

// GetPipelineDAG retrieves all steps and dependencies for a pipeline
func (pm *PipelineManager) GetPipelineDAG(ctx context.Context, pipelineID uuid.UUID) ([]PipelineStep, map[uuid.UUID][]uuid.UUID, error) {
	// Get all steps for the pipeline
	steps, err := pm.pipelineStepRepository.GetPipelineSteps(ctx, pipelineID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get pipeline steps: %w", err)
	}

	// Get all dependencies for the pipeline
	dependencies, err := pm.pipelineStepRepository.GetPipelineDependencies(ctx, pipelineID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get dependencies: %w", err)
	}

	return steps, dependencies, nil
}

// ValidatePipeline checks if a pipeline is valid for execution
func (pm *PipelineManager) ValidatePipeline(ctx context.Context, pipelineID uuid.UUID) error {
	steps, deps, err := pm.GetPipelineDAG(ctx, pipelineID)
	if err != nil {
		return err
	}

	if len(steps) == 0 {
		return fmt.Errorf("pipeline has no steps")
	}

	// Create DAG and validate
	dag, err := NewDAG(steps, deps)
	if err != nil {
		return fmt.Errorf("invalid pipeline: %w", err)
	}

	// Check for unreachable nodes
	sorted, err := dag.TopologicalSort()
	if err != nil {
		return fmt.Errorf("pipeline validation failed: %w", err)
	}

	if len(sorted) != len(steps) {
		return fmt.Errorf("pipeline contains unreachable steps")
	}

	return nil
}

// GetExecutionPlan returns the execution plan for a pipeline
func (pm *PipelineManager) GetExecutionPlan(ctx context.Context, pipelineID uuid.UUID) ([][]uuid.UUID, error) {
	steps, deps, err := pm.GetPipelineDAG(ctx, pipelineID)
	if err != nil {
		return nil, err
	}

	dag, err := NewDAG(steps, deps)
	if err != nil {
		return nil, fmt.Errorf("failed to create DAG: %w", err)
	}

	return dag.GetExecutionPlan()
}

// PipelineStepResponse represents the API response for a pipeline step
type PipelineStepResponse struct {
	ID           openapi_types.UUID     `json:"id"`
	PipelineID   openapi_types.UUID     `json:"pipelineId"`
	ToolID       openapi_types.UUID     `json:"toolId"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	Config       map[string]interface{} `json:"config"`
	Position     int                    `json:"position"`
	Dependencies []openapi_types.UUID   `json:"dependencies"`
}

// ConvertStepToResponse converts a domain step to API response
func ConvertStepToResponse(step *PipelineStep, dependencies []uuid.UUID) PipelineStepResponse {
	deps := make([]openapi_types.UUID, len(dependencies))
	copy(deps, dependencies)

	return PipelineStepResponse{
		ID:           step.Id,
		PipelineID:   step.PipelineId,
		ToolID:       step.ToolId,
		Name:         step.Name,
		Description:  step.Description,
		Config:       step.Config,
		Position:     step.Position,
		Dependencies: deps,
	}
}

// ExecutionPlanResponse represents the execution plan for a pipeline
type ExecutionPlanResponse struct {
	PipelineID openapi_types.UUID `json:"pipelineId"`
	Levels     []ExecutionLevel   `json:"levels"`
	TotalSteps int                `json:"totalSteps"`
	IsValid    bool               `json:"isValid"`
}

// ExecutionLevel represents a level of parallel execution
type ExecutionLevel struct {
	Level int                  `json:"level"`
	Steps []openapi_types.UUID `json:"steps"`
}

// ConvertExecutionPlan converts internal execution plan to API response
func ConvertExecutionPlan(pipelineID uuid.UUID, plan [][]uuid.UUID) ExecutionPlanResponse {
	levels := make([]ExecutionLevel, len(plan))
	totalSteps := 0

	for i, level := range plan {
		steps := make([]openapi_types.UUID, len(level))
		copy(steps, level)
		levels[i] = ExecutionLevel{
			Level: i,
			Steps: steps,
		}
		totalSteps += len(level)
	}

	return ExecutionPlanResponse{
		PipelineID: pipelineID,
		Levels:     levels,
		TotalSteps: totalSteps,
		IsValid:    true,
	}
}
