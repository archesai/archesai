package runs

import (
	"context"

	"github.com/archesai/archesai/internal/core/valueobjects"
)

// UpdateRunProgressCommand represents a command to update run progress.
type UpdateRunProgressCommand struct {
	RunID    valueobjects.RunID
	Progress int
	Message  string
}

// NewUpdateRunProgressCommand creates a new update run progress command.
func NewUpdateRunProgressCommand(runID valueobjects.RunID, progress int) *UpdateRunProgressCommand {
	return &UpdateRunProgressCommand{
		RunID:    runID,
		Progress: progress,
	}
}

// WithMessage sets the progress message.
func (c *UpdateRunProgressCommand) WithMessage(message string) *UpdateRunProgressCommand {
	c.Message = message
	return c
}

// UpdateRunProgressCommandHandler handles the update run progress command.
type UpdateRunProgressCommandHandler interface {
	Handle(ctx context.Context, command *UpdateRunProgressCommand) error
}
