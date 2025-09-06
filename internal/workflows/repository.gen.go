package workflows

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for workflow data persistence
type Repository interface {
	// Pipeline operations
	CreatePipeline(ctx context.Context, pipeline *Pipeline) (*Pipeline, error)
	GetPipeline(ctx context.Context, id uuid.UUID) (*Pipeline, error)
	UpdatePipeline(ctx context.Context, pipeline *Pipeline) (*Pipeline, error)
	DeletePipeline(ctx context.Context, id uuid.UUID) error
	ListPipelines(ctx context.Context, orgID string, limit, offset int) ([]*Pipeline, int, error)

	// Run operations
	CreateRun(ctx context.Context, run *Run) (*Run, error)
	GetRun(ctx context.Context, id uuid.UUID) (*Run, error)
	UpdateRun(ctx context.Context, run *Run) (*Run, error)
	DeleteRun(ctx context.Context, id uuid.UUID) error
	ListRuns(ctx context.Context, orgID string, limit, offset int) ([]*Run, int, error)
	ListRunsByPipeline(ctx context.Context, pipelineID string, limit, offset int) ([]*Run, int, error)

	// Tool operations
	CreateTool(ctx context.Context, tool *Tool) (*Tool, error)
	GetTool(ctx context.Context, id uuid.UUID) (*Tool, error)
	UpdateTool(ctx context.Context, tool *Tool) (*Tool, error)
	DeleteTool(ctx context.Context, id uuid.UUID) error
	ListTools(ctx context.Context, orgID string, limit, offset int) ([]*Tool, int, error)
}
