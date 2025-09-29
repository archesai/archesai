package runs

import (
	"context"

	"github.com/archesai/archesai/internal/core/aggregates"
	"github.com/archesai/archesai/internal/core/valueobjects"
)

// CreateRunCommand represents a command to create a new run.
type CreateRunCommand struct {
	PipelineID     valueobjects.PipelineID
	ToolID         valueobjects.ToolID
	OrganizationID valueobjects.OrganizationID
	CreatedBy      valueobjects.UserID
	InputData      map[string]interface{}
	Configuration  map[string]interface{}
}

// NewCreateRunCommand creates a new create run command.
func NewCreateRunCommand(
	pipelineID valueobjects.PipelineID,
	toolID valueobjects.ToolID,
	organizationID valueobjects.OrganizationID,
	createdBy valueobjects.UserID,
) *CreateRunCommand {
	return &CreateRunCommand{
		PipelineID:     pipelineID,
		ToolID:         toolID,
		OrganizationID: organizationID,
		CreatedBy:      createdBy,
		InputData:      make(map[string]interface{}),
		Configuration:  make(map[string]interface{}),
	}
}

// WithInputData sets the input data for the run.
func (c *CreateRunCommand) WithInputData(data map[string]interface{}) *CreateRunCommand {
	c.InputData = data
	return c
}

// WithConfiguration sets the configuration for the run.
func (c *CreateRunCommand) WithConfiguration(config map[string]interface{}) *CreateRunCommand {
	c.Configuration = config
	return c
}

// CreateRunCommandHandler handles the create run command.
type CreateRunCommandHandler interface {
	Handle(ctx context.Context, command *CreateRunCommand) (*aggregates.Run, error)
}
