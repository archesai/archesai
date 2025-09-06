// Package swarm provides agent orchestration and workflow management for LLM-based systems
package swarm

import (
	"github.com/archesai/archesai/internal/llm"
)

// Agent represents an entity with specific attributes and behaviors.
type Agent struct {
	Name              string                                               // The name of the agent.
	Model             string                                               // The model identifier.
	Provider          llm.Provider                                         // The LLM provider to use.
	Config            *ClientConfig                                        // Provider-specific configuration.
	Instructions      string                                               // Static instructions for the agent.
	InstructionsFunc  func(contextVariables map[string]interface{}) string // Function to generate dynamic instructions based on context.
	Functions         []AgentFunction                                      // A list of functions the agent can perform.
	Memory            *MemoryStore                                         // Memory store for the agent.
	ParallelToolCalls bool                                                 // Whether to allow parallel tool calls.
}

// AgentFunction represents a function that can be performed by an agent
type AgentFunction struct {
	Name        string                                                                            // The name of the function.
	Description string                                                                            // Description of what the function does.
	Parameters  map[string]interface{}                                                            // Parameters for the function.
	Function    func(args map[string]interface{}, contextVariables map[string]interface{}) Result // The actual function implementation.
}

// FunctionToDefinition converts an AgentFunction to a llm.Function
func FunctionToDefinition(af AgentFunction) llm.Function {
	return llm.Function{
		Name:        af.Name,
		Description: af.Description,
		Parameters:  af.Parameters,
	}
}

// NewAgent creates a new agent with initialized memory store
func NewAgent(name, model string, provider llm.Provider) *Agent {
	return &Agent{
		Name:     name,
		Model:    model,
		Provider: provider,
		Memory:   NewMemoryStore(100), // Default to 100 short-term memories
	}
}

// WithConfig sets the configuration for the agent
func (a *Agent) WithConfig(config *ClientConfig) *Agent {
	a.Config = config
	return a
}

// WithInstructions sets the static instructions for the agent
func (a *Agent) WithInstructions(instructions string) *Agent {
	a.Instructions = instructions
	return a
}

// WithInstructionsFunc sets the dynamic instructions function for the agent
func (a *Agent) WithInstructionsFunc(f func(map[string]interface{}) string) *Agent {
	a.InstructionsFunc = f
	return a
}

// WithFunctions sets the functions available to the agent
func (a *Agent) WithFunctions(functions []AgentFunction) *Agent {
	a.Functions = functions
	return a
}

// WithParallelToolCalls enables or disables parallel tool calls
func (a *Agent) WithParallelToolCalls(enabled bool) *Agent {
	a.ParallelToolCalls = enabled
	return a
}
