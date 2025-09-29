package labels

import (
	"context"

	"github.com/archesai/archesai/internal/core/valueobjects"
)

// DeleteLabelCommand represents a command to delete a label.
type DeleteLabelCommand struct {
	LabelID valueobjects.LabelID
}

// NewDeleteLabelCommand creates a new delete label command.
func NewDeleteLabelCommand(labelID valueobjects.LabelID) *DeleteLabelCommand {
	return &DeleteLabelCommand{
		LabelID: labelID,
	}
}

// DeleteLabelCommandHandler handles the delete label command.
type DeleteLabelCommandHandler interface {
	Handle(ctx context.Context, command *DeleteLabelCommand) error
}
