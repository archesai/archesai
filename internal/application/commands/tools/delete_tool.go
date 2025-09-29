package tools

import (
	"context"

	"github.com/archesai/archesai/internal/core/valueobjects"
)

// DeleteToolCommand represents a command to delete a tool.
type DeleteToolCommand struct {
	ToolID valueobjects.ToolID
}

// NewDeleteToolCommand creates a new delete tool command.
func NewDeleteToolCommand(toolID valueobjects.ToolID) *DeleteToolCommand {
	return &DeleteToolCommand{
		ToolID: toolID,
	}
}

// DeleteToolCommandHandler handles the delete tool command.
type DeleteToolCommandHandler interface {
	Handle(ctx context.Context, command *DeleteToolCommand) error
}
