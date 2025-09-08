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

	// Example 2: TUI with Ollama (local)
	// runOllamaExample()

	// Example 3: TUI with custom personas
	// runCustomPersonasExample()
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

	// Create personas using the default ones
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

// runOllamaExample shows how to run the TUI with local Ollama
func RunOllamaExample() {
	// Initialize Ollama LLM client (no API key needed)
	llmClient, err := llm.NewOllamaLLM()
	if err != nil {
		log.Fatalf("Failed to create Ollama client: %v", err)
	}
	chatClient := llm.NewChatClient(llmClient)

	// Create personas for Ollama
	personas := []*llm.ChatPersona{
		{
			Name:         "Local Assistant",
			SystemPrompt: "You are a helpful AI assistant running locally.",
			Model:        "llama2",
			Temperature:  0.7,
		},
		{
			Name:         "Local Coder",
			SystemPrompt: "You are a coding expert running locally. Help with programming questions.",
			Model:        "llama2",
			Temperature:  0.3,
		},
	}

	// Create and run TUI
	model := tui.New(chatClient, personas)
	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		log.Fatalf("Error running TUI: %v", err)
	}
}

// runCustomPersonasExample shows how to create custom personas
func RunCustomPersonasExample() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set OPENAI_API_KEY environment variable")
	}

	// Initialize OpenAI LLM client
	llmClient := llm.NewOpenAILLM(apiKey)
	chatClient := llm.NewChatClient(llmClient)

	// Create custom personas
	personas := []*llm.ChatPersona{
		{
			Name: "Shakespeare",
			SystemPrompt: `You are William Shakespeare, the famous English playwright and poet. 
			Respond in the style of Shakespearean English, with flowery language and poetic flair. 
			Use 'thee', 'thou', 'thy' and other archaic terms appropriately.`,
			Model:       "gpt-4",
			Temperature: 0.9, // Higher temperature for more creative responses
		},
		{
			Name: "Pirate Captain",
			SystemPrompt: `You are a swashbuckling pirate captain from the golden age of piracy. 
			Speak like a stereotypical pirate with 'arr', 'matey', 'ye' and nautical terminology. 
			Tell tales of treasure hunts and adventures on the high seas.`,
			Model:       "gpt-4",
			Temperature: 0.8,
		},
		{
			Name: "Zen Master",
			SystemPrompt: `You are a wise Zen master who speaks in philosophical riddles and koans. 
			Give thoughtful, meditative responses that encourage deep reflection. 
			Use nature metaphors and speak about the path to enlightenment.`,
			Model:       "gpt-4",
			Temperature: 0.6,
		},
	}

	// Create and run TUI
	model := tui.New(chatClient, personas)
	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		log.Fatalf("Error running TUI: %v", err)
	}
}
