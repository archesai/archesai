package runs

import (
	"context"

	"github.com/archesai/archesai/internal/core/valueobjects"
)

// FailRunCommand represents a command to fail a run.
type FailRunCommand struct {
	RunID  valueobjects.RunID
	Reason string
}

// NewFailRunCommand creates a new fail run command.
func NewFailRunCommand(runID valueobjects.RunID, reason string) *FailRunCommand {
	return &FailRunCommand{
		RunID:  runID,
		Reason: reason,
	}
}

// FailRunCommandHandler handles the fail run command.
type FailRunCommandHandler interface {
	Handle(ctx context.Context, command *FailRunCommand) error
}
