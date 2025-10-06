// Package pipeline provides command and query handlers for pipeline operations.
package pipeline

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/archesai/archesai/internal/core/entities"
	"github.com/archesai/archesai/internal/core/events"
	"github.com/archesai/archesai/internal/core/repositories"
)

// CreatePipelineStepCommandHandler handles the create pipeline step command.
type CreatePipelineStepCommandHandler struct {
	pipelineRepo repositories.PipelineRepository
	publisher    events.Publisher
}

// NewCreatePipelineStepCommandHandler creates a new create pipeline step command handler.
func NewCreatePipelineStepCommandHandler(
	pipelineRepo repositories.PipelineRepository,
	publisher events.Publisher,
) *CreatePipelineStepCommandHandler {
	return &CreatePipelineStepCommandHandler{
		pipelineRepo: pipelineRepo,
		publisher:    publisher,
	}
}

// Handle executes the create pipeline step command.
func (h *CreatePipelineStepCommandHandler) Handle(
	ctx context.Context,
	cmd *CreatePipelineStepCommand,
) (*entities.PipelineStep, error) {
	// Verify pipeline exists
	pipeline, err := h.pipelineRepo.Get(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline: %w", err)
	}

	// Create the pipeline step entity
	now := time.Now().UTC()
	step := &entities.PipelineStep{
		ID:         uuid.New(),
		PipelineID: pipeline.ID,
		ToolID:     cmd.ToolID,
		// Name:         cmd.Name,
		// Description:  cmd.Description,
		// Config:       cmd.Config,
		// Position:     cmd.Position,
		// Dependencies: cmd.Dependencies,
		// Status:       entities.PipelineStepStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// TODO: Implement PipelineStepRepository and use it here
	// For now, return the created step without persisting
	// created, err := h.pipelineStepRepo.Create(ctx, step)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create pipeline step: %w", err)
	// }
	_ = step // Suppress unused variable warning

	// Publish domain event
	event := events.NewPipelineStepCreatedEvent(step.ID)
	if err := h.publisher.Publish(ctx, event); err != nil {
		// Log error but don't fail the operation
		_ = err
	}

	return step, nil
}
