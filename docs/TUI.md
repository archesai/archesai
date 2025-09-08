# ArchesAI Terminal User Interface (TUI)

## Overview

The ArchesAI TUI provides an interactive terminal interface for configuration management and AI agent interactions.

## Features

- **Configuration Viewer** (default mode) - Browse all settings without needing API keys
- **AI Chat Interface** (--chat mode) - Interactive conversations with AI agents
- **Beautiful Design** - Centered layout with color-coded status indicators
- **Multi-Provider Support** - OpenAI, Claude, Gemini, Ollama, DeepSeek

## Installation

The TUI is built into the main ArchesAI binary:

```bash
go build -o archesai cmd/archesai/main.go
```

## Usage

### Configuration Viewer (Default)

Launch the configuration viewer without any API keys:

```bash
archesai tui
```

This displays:

- 🗄️ Database Configuration
- 🌐 Server Configuration
- 🔐 Authentication Settings
- 📦 Redis Configuration
- 💾 Storage Settings
- 🤖 AI Agents & LLM Providers
- 📊 System Status
- 🔧 Environment Variables

### AI Chat Mode

Launch the interactive AI chat interface:

```bash
# With OpenAI (default)
export OPENAI_API_KEY=your-key
archesai tui --chat

# With specific provider
archesai tui --chat --provider=claude --model=claude-3-opus

# With API key flag
archesai tui --chat --api-key=your-key --provider=openai

# With Ollama (local)
archesai tui --chat --provider=ollama --model=llama2
```

## Keyboard Controls

| Key      | Action            | Context           |
| -------- | ----------------- | ----------------- |
| `↑/↓`    | Navigate items    | Menu/Lists        |
| `Enter`  | Select item       | Menu/Send message |
| `ESC`    | Go back           | Sub-menus         |
| `Tab`    | Toggle agent list | Chat mode         |
| `Ctrl+C` | Quit              | Any time          |

## Configuration Display

### Status Indicators

- **●** Green dot - Active/Configured
- **●** Orange dot - Inactive/Optional
- **○** Red circle - Not configured
- **✓** Check mark - Enabled
- **✗** X mark - Disabled

### Value Display

- **Regular values**: Shown in light gray
- **Sensitive data**: Displayed as `••••••••`
- **Empty values**: Shown as "_not set_" in italics

## Examples

### Basic Usage

The TUI provides interactive prompts and feedback for all operations.

### Basic Usage Example

```go
package main

import (
    "log"
    "os"

    "github.com/archesai/archesai/internal/llm"
    "github.com/archesai/archesai/internal/tui"
    tea "github.com/charmbracelet/bubbletea"
)

func main() {
    // Initialize OpenAI client
    apiKey := os.Getenv("OPENAI_API_KEY")
    llmClient := llm.NewOpenAILLM(apiKey)
    chatClient := llm.NewChatClient(llmClient)

    // Use built-in default personas
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
```

### Using Default Personas

```go
// Use built-in default personas
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
```

## Architecture

The TUI is built with:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [ArchesAI LLM Package](../internal/llm/) - Direct LLM client interfaces

## Design Principles

1. **No Setup Required** - Config viewer works without any API keys
2. **Centered Layout** - Fixed 60-character width for clean appearance
3. **Visual Feedback** - Color-coded status indicators
4. **Responsive** - Adapts to terminal size
5. **Consistent** - Same formatting across all sections
6. **Simple Architecture** - Direct LLM usage without complex abstractions

## Chat Interface Features

### Multi-Agent Support

- Switch between different AI personas using `Tab`
- Each persona has its own conversation history
- Configurable system prompts and temperature settings

### Session Management

- Automatic conversation history management
- Context preservation across messages
- Clean session state for each persona

### Provider Support

- **OpenAI**: Full support with all models (gpt-4, gpt-3.5-turbo, etc.)
- **Ollama**: Local model support (llama2, mistral, etc.)
- **Claude/Gemini/DeepSeek**: Ready for implementation

## Troubleshooting

### Terminal Issues

- Ensure your terminal supports 256 colors
- Use a modern terminal (iTerm2, Windows Terminal, Alacritty)
- Try resizing if display appears broken

### API Key Issues

- Set environment variables: `OPENAI_API_KEY`, `ANTHROPIC_API_KEY`, etc.
- Or use `--api-key` flag
- Verify key has proper permissions

### Provider Notes

- **OpenAI**: Fully functional with streaming support
- **Ollama**: Requires local Ollama server running (http://localhost:11434)
- **Claude/Gemini**: Implementation ready, just need client setup

## Development

To modify the TUI:

1. Edit `internal/tui/config_tui.go` for config viewer
2. Edit `internal/tui/tui.go` for chat interface
3. Edit `internal/cli/tui.go` for CLI integration
4. Edit `internal/llm/chat.go` for chat client functionality

Run tests:

```bash
go test ./internal/tui/...
go test ./internal/llm/...
```

## Future Enhancements

- [ ] Search/filter in configuration
- [ ] Edit configuration values directly
- [ ] Export/import configuration
- [ ] Connection testing for services
- [ ] Theme customization
- [ ] Streaming response support in TUI
- [ ] File upload support for chat
- [ ] Tool/function calling integration
