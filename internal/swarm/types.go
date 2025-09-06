package swarm

import (
	"github.com/archesai/archesai/internal/llm"
)

// Response represents the response from an agent
type Response struct {
	Messages         []llm.Message
	Agent            *Agent
	ContextVariables map[string]interface{}
	ToolResults      []ToolResult // Results from tool calls
}

// ToolResult represents the result of a tool call
type ToolResult struct {
	ToolName string      // Name of the tool that was called
	Args     interface{} // Arguments passed to the tool
	Result   Result      // Result returned by the tool
}

// Result represents the result of a function execution
type Result struct {
	Success bool        // Whether the function execution was successful
	Data    interface{} // Any data returned by the function
	Error   error       // Any error that occurred during execution
	Agent   *Agent      // Active agent
}
