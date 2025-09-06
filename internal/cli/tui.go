package cli

import (
	"fmt"
	"os"

	"github.com/archesai/archesai/internal/llm"
	"github.com/archesai/archesai/internal/swarm"
	"github.com/archesai/archesai/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// tuiCmd represents the tui command
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
	tuiCmd.Flags().BoolVar(&tuiChatMode, "chat", false, "Launch AI chat interface instead of config viewer")
	tuiCmd.Flags().StringVar(&tuiProvider, "provider", "openai", "LLM provider (openai, claude, gemini, ollama)")
	tuiCmd.Flags().StringVar(&tuiModel, "model", "gpt-4", "Model to use")
	tuiCmd.Flags().StringVar(&tuiAPIKey, "api-key", "", "API key for the provider (or use environment variable)")
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
			return fmt.Errorf("API key required for chat mode. Set --api-key flag or appropriate environment variable")
		}
	}

	// Convert provider string to LLMProvider
	var provider llm.Provider
	switch tuiProvider {
	case "openai":
		provider = llm.OpenAI
	case "claude":
		provider = llm.Claude
	case "gemini":
		provider = llm.Gemini
	case providerOllama:
		provider = llm.Ollama
	case "deepseek":
		provider = llm.DeepSeek
	default:
		return fmt.Errorf("unsupported provider: %s", tuiProvider)
	}

	// Initialize Swarm client
	swarmClient := swarm.NewSwarm(tuiAPIKey, provider)
	if swarmClient == nil {
		return fmt.Errorf("failed to initialize Swarm client")
	}

	// Create sample agents
	agents := createSampleAgents(provider)

	// Create and run the TUI
	model := tui.New(swarmClient, agents)
	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}

	return nil
}

func createSampleAgents(provider llm.Provider) []*swarm.Agent {
	agents := []*swarm.Agent{
		// General Assistant
		swarm.NewAgent("Assistant", "gpt-4", provider).
			WithInstructions("You are a helpful AI assistant. Be concise and informative."),

		// Code Helper
		swarm.NewAgent("CodeHelper", "gpt-4", provider).
			WithInstructions("You are a coding assistant. Help with programming questions, debugging, and best practices."),

		// Creative Writer
		swarm.NewAgent("CreativeWriter", "gpt-4", provider).
			WithInstructions("You are a creative writing assistant. Help with stories, poems, and creative content."),

		// Data Analyst
		swarm.NewAgent("DataAnalyst", "gpt-4", provider).
			WithInstructions("You are a data analysis expert. Help with data interpretation, statistics, and insights."),

		// Research Assistant
		swarm.NewAgent("Researcher", "gpt-4", provider).
			WithInstructions("You are a research assistant. Help find information, summarize topics, and provide citations when possible."),
	}

	// Adjust model based on provider
	modelName := getModelForProvider(provider)
	for _, agent := range agents {
		agent.Model = modelName
	}

	return agents
}

func getModelForProvider(provider llm.Provider) string {
	switch provider {
	case llm.Claude:
		return "claude-3-opus-20240229"
	case llm.Gemini:
		return "gemini-pro"
	case llm.Ollama:
		return "llama2"
	case llm.DeepSeek:
		return "deepseek-chat"
	default:
		return "gpt-4"
	}
}
