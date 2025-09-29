package runs

import (
	"context"

	"github.com/archesai/archesai/internal/core/valueobjects"
)

// CompleteRunCommand represents a command to complete a run.
type CompleteRunCommand struct {
	RunID      valueobjects.RunID
	OutputData map[string]interface{}
}

// NewCompleteRunCommand creates a new complete run command.
func NewCompleteRunCommand(runID valueobjects.RunID) *CompleteRunCommand {
	return &CompleteRunCommand{
		RunID:      runID,
		OutputData: make(map[string]interface{}),
	}
}

// WithOutputData sets the output data for the run.
func (c *CompleteRunCommand) WithOutputData(data map[string]interface{}) *CompleteRunCommand {
	c.OutputData = data
	return c
}

// CompleteRunCommandHandler handles the complete run command.
type CompleteRunCommandHandler interface {
	Handle(ctx context.Context, command *CompleteRunCommand) error
}
