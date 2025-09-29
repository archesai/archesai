package labels

import (
	"context"

	"github.com/archesai/archesai/internal/core/entities"
	"github.com/archesai/archesai/internal/core/valueobjects"
)

// UpdateLabelCommand represents a command to update a label.
type UpdateLabelCommand struct {
	LabelID valueobjects.LabelID
	Name    string
}

// NewUpdateLabelCommand creates a new update label command.
func NewUpdateLabelCommand(labelID valueobjects.LabelID, name string) *UpdateLabelCommand {
	return &UpdateLabelCommand{
		LabelID: labelID,
		Name:    name,
	}
}

// UpdateLabelCommandHandler handles the update label command.
type UpdateLabelCommandHandler interface {
	Handle(ctx context.Context, command *UpdateLabelCommand) (*entities.Label, error)
}
