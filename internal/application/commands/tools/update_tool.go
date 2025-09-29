package tools

import (
	"context"

	"github.com/archesai/archesai/internal/core/aggregates"
	"github.com/archesai/archesai/internal/core/valueobjects"
)

// UpdateToolCommand represents a command to update a tool.
type UpdateToolCommand struct {
	ToolID        valueobjects.ToolID
	Name          *string
	Description   *string
	Configuration map[string]interface{}
	Capabilities  []string
	MaxTokens     *int
	Temperature   *float64
	TopP          *float64
	Active        *bool
}

// NewUpdateToolCommand creates a new update tool command.
func NewUpdateToolCommand(toolID valueobjects.ToolID) *UpdateToolCommand {
	return &UpdateToolCommand{
		ToolID: toolID,
	}
}

// WithName sets the tool name.
func (c *UpdateToolCommand) WithName(name string) *UpdateToolCommand {
	c.Name = &name
	return c
}

// WithDescription sets the tool description.
func (c *UpdateToolCommand) WithDescription(description string) *UpdateToolCommand {
	c.Description = &description
	return c
}

// WithActive sets the tool active status.
func (c *UpdateToolCommand) WithActive(active bool) *UpdateToolCommand {
	c.Active = &active
	return c
}

// UpdateToolCommandHandler handles the update tool command.
type UpdateToolCommandHandler interface {
	Handle(ctx context.Context, command *UpdateToolCommand) (*aggregates.Tool, error)
}
