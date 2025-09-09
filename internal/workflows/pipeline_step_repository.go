package workflows

import (
	"context"

	"github.com/google/uuid"
)

// PipelineStepRepository provides operations for pipeline steps and their dependencies
type PipelineStepRepository interface {
	// Pipeline step operations
	GetPipelineSteps(ctx context.Context, pipelineID uuid.UUID) ([]PipelineStep, error)
	CreatePipelineStep(ctx context.Context, step *PipelineStep) (*PipelineStep, error)
	UpdatePipelineStep(ctx context.Context, stepID uuid.UUID, step *PipelineStep) (*PipelineStep, error)
	DeletePipelineStep(ctx context.Context, stepID uuid.UUID) error

	// Dependency operations
	GetStepDependencies(ctx context.Context, stepID uuid.UUID) ([]uuid.UUID, error)
	GetPipelineDependencies(ctx context.Context, pipelineID uuid.UUID) (map[uuid.UUID][]uuid.UUID, error)
	CreateStepDependency(ctx context.Context, stepID, dependsOnID uuid.UUID) error
	DeleteStepDependency(ctx context.Context, stepID, dependsOnID uuid.UUID) error
}
