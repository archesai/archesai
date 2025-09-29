package runs

import (
	"context"

	"github.com/archesai/archesai/internal/core/valueobjects"
)

// CancelRunCommand represents a command to cancel a run.
type CancelRunCommand struct {
	RunID valueobjects.RunID
}

// NewCancelRunCommand creates a new cancel run command.
func NewCancelRunCommand(runID valueobjects.RunID) *CancelRunCommand {
	return &CancelRunCommand{
		RunID: runID,
	}
}

// CancelRunCommandHandler handles the cancel run command.
type CancelRunCommandHandler interface {
	Handle(ctx context.Context, command *CancelRunCommand) error
}
