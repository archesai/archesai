package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/archesai/archesai/internal/adapters/tui"
	"github.com/archesai/archesai/internal/core/models"
	"github.com/archesai/archesai/internal/infrastructure/llm"
)

// tuiCmd represents the tui command.
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the interactive TUI for configuration and AI agents",
	Long: `Launch an interactive terminal user interface (TUI) for viewing
configuration and optionally interacting with AI agents.

The TUI provides:
- Configuration viewer for all settings
- Database, server, and auth configuration display
- AI provider status and configuration
- Optional chat interface with AI agents`,
	Example: `  archesai tui                  # Launch config viewer
  archesai tui --chat            # Launch AI chat interface
  archesai tui --chat --provider=openai --model=gpt-4`,
	RunE: runTUI,
}

const (
	providerOllama = "ollama"
)

var (
	tuiProvider string
	tuiModel    string
	tuiAPIKey   string
	tuiChatMode bool
)

func init() {
	rootCmd.AddCommand(tuiCmd)

	// TUI specific flags
	tuiCmd.Flags().
		BoolVar(&tuiChatMode, "chat", false, "Launch AI chat interface instead of config viewer")
	tuiCmd.Flags().
		StringVar(&tuiProvider, "provider", "openai", "LLM provider (openai, claude, gemini, ollama)")
	tuiCmd.Flags().StringVar(&tuiModel, "model", "gpt-4", "Model to use")
	tuiCmd.Flags().
		StringVar(&tuiAPIKey, "api-key", "", "API key for the provider (or use environment variable)")
}

func runTUI(_ *cobra.Command, _ []string) error {
	// If not in chat mode, launch the config viewer
	if !tuiChatMode {
		model := tui.NewConfigModel()
		program := tea.NewProgram(model, tea.WithAltScreen())

		if _, err := program.Run(); err != nil {
			return fmt.Errorf("error running config TUI: %w", err)
		}
		return nil
	}

	// Chat mode - get API key from environment if not provided
	if tuiAPIKey == "" {
		switch tuiProvider {
		case "openai":
			tuiAPIKey = os.Getenv("OPENAI_API_KEY")
		case "claude":
			tuiAPIKey = os.Getenv("ANTHROPIC_API_KEY")
		case "gemini":
			tuiAPIKey = os.Getenv("GEMINI_API_KEY")
		case providerOllama:
			// Ollama doesn't need an API key
			tuiAPIKey = "not-required"
		}

		if tuiAPIKey == "" && tuiProvider != providerOllama {
			return fmt.Errorf(
				"API key required for chat mode. Set --api-key flag or appropriate environment variable",
			)
		}
	}

	// Convert provider string to LLMProvider
	var provider models.Provider
	switch tuiProvider {
	case "openai":
		provider = models.ProviderOpenAI
	case "claude":
		provider = models.ProviderClaude
	case "gemini":
		provider = models.ProviderGemini
	case providerOllama:
		provider = models.ProviderOllama
	case "deepseek":
		provider = models.ProviderDeepSeek
	default:
		return fmt.Errorf("unsupported provider: %s", tuiProvider)
	}

	// Initialize LLM client
	var llmService llm.LLM
	switch provider {
	case models.ProviderOpenAI:
		llmService = llm.NewOpenAILLM(tuiAPIKey)
	case models.ProviderOllama:
		var err error
		llmService, err = llm.NewOllamaLLM()
		if err != nil {
			return fmt.Errorf("failed to create Ollama client: %w", err)
		}
	default:
		return fmt.Errorf("provider %s not yet implemented in new chat interface", tuiProvider)
	}

	if llmService == nil {
		return fmt.Errorf("failed to initialize LLM client")
	}

	// Create chat client
	chatClient := llm.NewChatClient(llmService)

	// Create sample personas with the specified model
	modelName := getModelForProvider(provider)
	personas := createSamplePersonas(modelName)

	// Create and run the TUI
	model := tui.New(chatClient, personas)
	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}

	return nil
}

func createSamplePersonas(model string) []*llm.ChatPersona {
	// Use the default personas but update the model
	personas := []*llm.ChatPersona{
		{
			Name:         llm.DefaultPersonas.Assistant.Name,
			SystemPrompt: llm.DefaultPersonas.Assistant.SystemPrompt,
			Model:        model,
			Temperature:  llm.DefaultPersonas.Assistant.Temperature,
		},
		{
			Name:         llm.DefaultPersonas.CodeHelper.Name,
			SystemPrompt: llm.DefaultPersonas.CodeHelper.SystemPrompt,
			Model:        model,
			Temperature:  llm.DefaultPersonas.CodeHelper.Temperature,
		},
		{
			Name:         llm.DefaultPersonas.CreativeWriter.Name,
			SystemPrompt: llm.DefaultPersonas.CreativeWriter.SystemPrompt,
			Model:        model,
			Temperature:  llm.DefaultPersonas.CreativeWriter.Temperature,
		},
		{
			Name:         llm.DefaultPersonas.DataAnalyst.Name,
			SystemPrompt: llm.DefaultPersonas.DataAnalyst.SystemPrompt,
			Model:        model,
			Temperature:  llm.DefaultPersonas.DataAnalyst.Temperature,
		},
		{
			Name:         llm.DefaultPersonas.Researcher.Name,
			SystemPrompt: llm.DefaultPersonas.Researcher.SystemPrompt,
			Model:        model,
			Temperature:  llm.DefaultPersonas.Researcher.Temperature,
		},
	}

	return personas
}

func getModelForProvider(provider models.Provider) string {
	switch provider {
	case models.ProviderClaude:
		return "claude-3-opus-20240229"
	case models.ProviderGemini:
		return "gemini-pro"
	case models.ProviderOllama:
		return "llama2"
	case models.ProviderDeepSeek:
		return "deepseek-chat"
	default:
		return "gpt-4"
	}
}
