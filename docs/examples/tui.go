// Package main demonstrates how to use the ArchesAI TUI with SwarmGo agents
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/archesai/archesai/internal/llm"
	"github.com/archesai/archesai/internal/swarm"
	"github.com/archesai/archesai/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Example 1: Basic TUI with OpenAI
	RunOpenAIExample()

	// Example 2: TUI with Claude
	// runClaudeExample()

	// Example 3: TUI with custom agents
	// runCustomAgentsExample()
}

// runOpenAIExample shows how to run the TUI with OpenAI
func RunOpenAIExample() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set OPENAI_API_KEY environment variable")
	}

	// Initialize Swarm with OpenAI
	swarmClient := swarm.NewSwarm(apiKey, llm.OpenAI)

	// Create agents
	agents := []*swarm.Agent{
		swarm.NewAgent("Assistant", "gpt-4", llm.OpenAI).
			WithInstructions("You are a helpful AI assistant."),

		swarm.NewAgent("Coder", "gpt-4", llm.OpenAI).
			WithInstructions("You are a coding expert. Help with programming questions."),
	}

	// Create and run TUI
	model := tui.New(swarmClient, agents)
	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		log.Fatalf("Error running TUI: %v", err)
	}
}

// runClaudeExample shows how to run the TUI with Claude
func RunClaudeExample() {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set ANTHROPIC_API_KEY environment variable")
	}

	// Initialize Swarm with Claude
	swarmClient := swarm.NewSwarm(apiKey, llm.Claude)

	// Create agents
	agents := []*swarm.Agent{
		swarm.NewAgent("Claude", "claude-3-opus-20240229", llm.Claude).
			WithInstructions("You are Claude, an AI assistant created by Anthropic."),
	}

	// Create and run TUI
	model := tui.New(swarmClient, agents)
	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		log.Fatalf("Error running TUI: %v", err)
	}
}

// runCustomAgentsExample shows how to create custom agents with functions
func RunCustomAgentsExample() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set OPENAI_API_KEY environment variable")
	}

	// Initialize Swarm
	swarmClient := swarm.NewSwarm(apiKey, llm.OpenAI)

	// Create a weather agent with a function
	weatherAgent := swarm.NewAgent("WeatherBot", "gpt-4", llm.OpenAI).
		WithInstructions("You are a weather assistant. Use the get_weather function to provide weather information.").
		WithFunctions([]swarm.AgentFunction{
			{
				Name:        "get_weather",
				Description: "Get the weather for a location",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"location": map[string]interface{}{
							"type":        "string",
							"description": "The city and state, e.g., San Francisco, CA",
						},
					},
					"required": []string{"location"},
				},
				Function: func(args map[string]interface{}, _ map[string]interface{}) swarm.Result {
					location := args["location"].(string)
					// Mock weather data
					weather := fmt.Sprintf("The weather in %s is sunny and 72°F", location)
					return swarm.Result{
						Data: weather,
					}
				},
			},
		})

	// Create a calculator agent
	calcAgent := swarm.NewAgent("Calculator", "gpt-4", llm.OpenAI).
		WithInstructions("You are a calculator assistant. Help with math problems.").
		WithFunctions([]swarm.AgentFunction{
			{
				Name:        "calculate",
				Description: "Perform a calculation",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"expression": map[string]interface{}{
							"type":        "string",
							"description": "The mathematical expression to evaluate",
						},
					},
					"required": []string{"expression"},
				},
				Function: func(args map[string]interface{}, _ map[string]interface{}) swarm.Result {
					expr := args["expression"].(string)
					// In a real implementation, you'd evaluate the expression
					return swarm.Result{
						Data: fmt.Sprintf("Result of %s = 42", expr),
					}
				},
			},
		})

	agents := []*swarm.Agent{weatherAgent, calcAgent}

	// Create and run TUI
	model := tui.New(swarmClient, agents)
	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		log.Fatalf("Error running TUI: %v", err)
	}
}
