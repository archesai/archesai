package pipeline

import (
	"context"
	"fmt"

	queries "github.com/archesai/archesai/apis/studio/generated/application/queries/pipeline"
	"github.com/archesai/archesai/apis/studio/generated/core/models"
	"github.com/archesai/archesai/apis/studio/generated/core/repositories"
)

// GetPipelineStepsQueryHandler handles the get pipeline steps query.
type GetPipelineStepsQueryHandler struct {
	pipelineRepo repositories.PipelineRepository
}

// NewGetPipelineStepsQueryHandler creates a new get pipeline steps query handler.
func NewGetPipelineStepsQueryHandler(
	pipelineRepo repositories.PipelineRepository,
) *GetPipelineStepsQueryHandler {
	return &GetPipelineStepsQueryHandler{
		pipelineRepo: pipelineRepo,
	}
}

// Handle executes the get pipeline steps query.
func (h *GetPipelineStepsQueryHandler) Handle(
	ctx context.Context,
	query *queries.GetPipelineStepsQuery,
) ([]*models.PipelineStep, error) {
	// Verify pipeline exists
	_, err := h.pipelineRepo.Get(ctx, query.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline: %w", err)
	}

	// TODO: Implement PipelineStepRepository to fetch steps
	// For now, return empty slice
	// steps, _, err := h.pipelineStepRepo.ListByPipeline(ctx, query.ID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to list pipeline steps: %w", err)
	// }
	steps := []*models.PipelineStep{}

	return steps, nil
}
