package runs

import (
	"context"

	"github.com/archesai/archesai/internal/core/valueobjects"
)

// StartRunCommand represents a command to start a run.
type StartRunCommand struct {
	RunID valueobjects.RunID
}

// NewStartRunCommand creates a new start run command.
func NewStartRunCommand(runID valueobjects.RunID) *StartRunCommand {
	return &StartRunCommand{
		RunID: runID,
	}
}

// StartRunCommandHandler handles the start run command.
type StartRunCommandHandler interface {
	Handle(ctx context.Context, command *StartRunCommand) error
}
