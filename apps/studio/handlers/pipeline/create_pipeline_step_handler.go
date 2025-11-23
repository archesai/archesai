// Package pipeline provides command and query handlers for pipeline operations.
package pipeline

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	commands "github.com/archesai/archesai/apps/studio/generated/application/commands/pipeline"
	domainevents "github.com/archesai/archesai/apps/studio/generated/core/events"
	"github.com/archesai/archesai/apps/studio/generated/core/models"
	"github.com/archesai/archesai/apps/studio/generated/core/repositories"
	"github.com/archesai/archesai/pkg/events"
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
	cmd *commands.CreatePipelineStepCommand,
) (*models.PipelineStep, error) {
	// Verify pipeline exists
	pipeline, err := h.pipelineRepo.Get(ctx, cmd.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline: %w", err)
	}

	// Create the pipeline step entity
	now := time.Now().UTC()
	step := &models.PipelineStep{
		ID:         uuid.New(),
		PipelineID: pipeline.ID,
		ToolID:     cmd.ToolID,
		// Name:         cmd.Name,
		// Description:  cmd.Description,
		// Config:       cmd.Config,
		// Position:     cmd.Position,
		// Dependencies: cmd.Dependencies,
		// Status:       models.PipelineStepStatusPending,
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
	event := domainevents.NewPipelineStepCreatedEvent(step.ID)
	if err := h.publisher.Publish(ctx, event); err != nil {
		slog.Error("failed to publish event", "error", err)
	}

	return step, nil
}
