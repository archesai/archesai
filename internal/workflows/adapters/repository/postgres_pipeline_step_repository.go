// Package repository provides concrete implementations of repository interfaces
package repository

import (
	"context"
	"fmt"

	"github.com/archesai/archesai/internal/database/postgresql"
	"github.com/archesai/archesai/internal/workflows"
	"github.com/google/uuid"
)

// PostgresPipelineStepRepository implements PipelineStepRepository using PostgreSQL
type PostgresPipelineStepRepository struct {
	queries *postgresql.Queries
}

// NewPostgresPipelineStepRepository creates a new PostgreSQL pipeline step repository
func NewPostgresPipelineStepRepository(queries *postgresql.Queries) workflows.PipelineStepRepository {
	return &PostgresPipelineStepRepository{
		queries: queries,
	}
}

// GetPipelineSteps retrieves all steps for a pipeline
func (r *PostgresPipelineStepRepository) GetPipelineSteps(ctx context.Context, pipelineID uuid.UUID) ([]workflows.PipelineStep, error) {
	dbSteps, err := r.queries.ListPipelineStepsByPipeline(ctx, postgresql.ListPipelineStepsByPipelineParams{
		PipelineId: pipelineID,
		Limit:      1000, // TODO: Add proper pagination
		Offset:     0,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline steps: %w", err)
	}

	steps := make([]workflows.PipelineStep, len(dbSteps))
	for i, dbStep := range dbSteps {
		steps[i] = workflows.PipelineStep{
			Id:         dbStep.Id,
			PipelineId: dbStep.PipelineId,
			ToolId:     dbStep.ToolId,
			// TODO: Add other fields like Config, Name, Description when available
		}
	}

	return steps, nil
}

// CreatePipelineStep creates a new pipeline step
func (r *PostgresPipelineStepRepository) CreatePipelineStep(ctx context.Context, step *workflows.PipelineStep) (*workflows.PipelineStep, error) {
	dbStep, err := r.queries.CreatePipelineStep(ctx, postgresql.CreatePipelineStepParams{
		Id:         step.Id,
		PipelineId: step.PipelineId,
		ToolId:     step.ToolId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create pipeline step: %w", err)
	}

	return &workflows.PipelineStep{
		Id:         dbStep.Id,
		PipelineId: dbStep.PipelineId,
		ToolId:     dbStep.ToolId,
		// TODO: Map additional fields
	}, nil
}

// UpdatePipelineStep updates an existing pipeline step
func (r *PostgresPipelineStepRepository) UpdatePipelineStep(ctx context.Context, stepID uuid.UUID, step *workflows.PipelineStep) (*workflows.PipelineStep, error) {
	dbStep, err := r.queries.UpdatePipelineStep(ctx, postgresql.UpdatePipelineStepParams{
		Id:     stepID,
		ToolId: &step.ToolId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update pipeline step: %w", err)
	}

	return &workflows.PipelineStep{
		Id:         dbStep.Id,
		PipelineId: dbStep.PipelineId,
		ToolId:     dbStep.ToolId,
	}, nil
}

// DeletePipelineStep deletes a pipeline step
func (r *PostgresPipelineStepRepository) DeletePipelineStep(ctx context.Context, stepID uuid.UUID) error {
	err := r.queries.DeletePipelineStep(ctx, stepID)
	if err != nil {
		return fmt.Errorf("failed to delete pipeline step: %w", err)
	}
	return nil
}

// GetStepDependencies retrieves dependencies for a specific step
func (r *PostgresPipelineStepRepository) GetStepDependencies(ctx context.Context, stepID uuid.UUID) ([]uuid.UUID, error) {
	deps, err := r.queries.GetStepDependencies(ctx, stepID)
	if err != nil {
		return nil, fmt.Errorf("failed to get step dependencies: %w", err)
	}
	return deps, nil
}

// GetPipelineDependencies retrieves all dependencies for a pipeline
func (r *PostgresPipelineStepRepository) GetPipelineDependencies(ctx context.Context, pipelineID uuid.UUID) (map[uuid.UUID][]uuid.UUID, error) {
	deps, err := r.queries.GetPipelineStepDependencies(ctx, pipelineID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline dependencies: %w", err)
	}

	dependencies := make(map[uuid.UUID][]uuid.UUID)
	for _, dep := range deps {
		stepID := dep.PipelineStepId
		prerequisiteID := dep.PrerequisiteId

		if _, exists := dependencies[stepID]; !exists {
			dependencies[stepID] = []uuid.UUID{}
		}
		dependencies[stepID] = append(dependencies[stepID], prerequisiteID)
	}

	return dependencies, nil
}

// CreateStepDependency creates a dependency between two steps
func (r *PostgresPipelineStepRepository) CreateStepDependency(ctx context.Context, stepID, dependsOnID uuid.UUID) error {
	err := r.queries.CreatePipelineStepDependency(ctx, postgresql.CreatePipelineStepDependencyParams{
		PipelineStepId: stepID,
		PrerequisiteId: dependsOnID,
	})
	if err != nil {
		return fmt.Errorf("failed to create step dependency: %w", err)
	}
	return nil
}

// DeleteStepDependency deletes a dependency between two steps
func (r *PostgresPipelineStepRepository) DeleteStepDependency(ctx context.Context, stepID, dependsOnID uuid.UUID) error {
	err := r.queries.DeletePipelineStepDependency(ctx, postgresql.DeletePipelineStepDependencyParams{
		PipelineStepId: stepID,
		PrerequisiteId: dependsOnID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete step dependency: %w", err)
	}
	return nil
}
