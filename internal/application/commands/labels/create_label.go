package labels

import (
	"context"

	"github.com/archesai/archesai/internal/core/entities"
	"github.com/archesai/archesai/internal/core/valueobjects"
)

// CreateLabelCommand represents a command to create a new label.
type CreateLabelCommand struct {
	Name           string
	OrganizationID valueobjects.OrganizationID
	CreatedBy      valueobjects.UserID
}

// NewCreateLabelCommand creates a new create label command.
func NewCreateLabelCommand(
	name string,
	organizationID valueobjects.OrganizationID,
	createdBy valueobjects.UserID,
) *CreateLabelCommand {
	return &CreateLabelCommand{
		Name:           name,
		OrganizationID: organizationID,
		CreatedBy:      createdBy,
	}
}

// CreateLabelCommandHandler handles the create label command.
type CreateLabelCommandHandler interface {
	Handle(ctx context.Context, command *CreateLabelCommand) (*entities.Label, error)
}
