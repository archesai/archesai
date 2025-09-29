package tools

import (
	"context"

	"github.com/archesai/archesai/internal/core/aggregates"
	"github.com/archesai/archesai/internal/core/valueobjects"
)

// CreateToolCommand represents a command to create a new tool.
type CreateToolCommand struct {
	OrganizationID valueobjects.OrganizationID
	Name           string
	ToolType       string
	Description    string
	Configuration  map[string]interface{}
	Capabilities   []string
	MaxTokens      int
	Temperature    float64
	TopP           float64
}

// NewCreateToolCommand creates a new create tool command.
func NewCreateToolCommand(
	organizationID valueobjects.OrganizationID,
	name string,
	toolType string,
) *CreateToolCommand {
	return &CreateToolCommand{
		OrganizationID: organizationID,
		Name:           name,
		ToolType:       toolType,
		Configuration:  make(map[string]interface{}),
		Capabilities:   []string{},
		MaxTokens:      4096,
		Temperature:    0.7,
		TopP:           1.0,
	}
}

// WithDescription sets the tool description.
func (c *CreateToolCommand) WithDescription(description string) *CreateToolCommand {
	c.Description = description
	return c
}

// WithConfiguration sets the tool configuration.
func (c *CreateToolCommand) WithConfiguration(config map[string]interface{}) *CreateToolCommand {
	c.Configuration = config
	return c
}

// WithCapabilities sets the tool capabilities.
func (c *CreateToolCommand) WithCapabilities(capabilities []string) *CreateToolCommand {
	c.Capabilities = capabilities
	return c
}

// WithParameters sets the LLM parameters.
func (c *CreateToolCommand) WithParameters(
	maxTokens int,
	temperature, topP float64,
) *CreateToolCommand {
	c.MaxTokens = maxTokens
	c.Temperature = temperature
	c.TopP = topP
	return c
}

// CreateToolCommandHandler handles the create tool command.
type CreateToolCommandHandler interface {
	Handle(ctx context.Context, command *CreateToolCommand) (*aggregates.Tool, error)
}
