// Package main demonstrates how to use the ArchesAI TUI with chat interfaces
package main

import (
	"log"
	"os"

	"github.com/archesai/archesai/internal/llm"
	"github.com/archesai/archesai/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Example 1: Basic TUI with OpenAI
	RunOpenAIExample()
}

// runOpenAIExample shows how to run the TUI with OpenAI
func RunOpenAIExample() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set OPENAI_API_KEY environment variable")
	}

	// Initialize OpenAI LLM client
	llmClient := llm.NewOpenAILLM(apiKey)
	chatClient := llm.NewChatClient(llmClient)

	// Use the built-in default personas
	personas := []*llm.ChatPersona{
		{
			Name:         llm.DefaultPersonas.Assistant.Name,
			SystemPrompt: llm.DefaultPersonas.Assistant.SystemPrompt,
			Model:        "gpt-4",
			Temperature:  llm.DefaultPersonas.Assistant.Temperature,
		},
		{
			Name:         llm.DefaultPersonas.CodeHelper.Name,
			SystemPrompt: llm.DefaultPersonas.CodeHelper.SystemPrompt,
			Model:        "gpt-4",
			Temperature:  llm.DefaultPersonas.CodeHelper.Temperature,
		},
	}

	// Create and run TUI
	model := tui.New(chatClient, personas)
	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		log.Fatalf("Error running TUI: %v", err)
	}
}
