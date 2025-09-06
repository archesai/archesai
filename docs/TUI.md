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

See [examples/tui_demo.go](examples/tui.go) for complete working examples.

### Creating Custom Agents

```go
// Create an agent with custom instructions
agent := swarm.NewAgent("Assistant", "gpt-4", llm.OpenAI).
    WithInstructions("You are a helpful AI assistant.")

// Create agent with functions
weatherAgent := swarm.NewAgent("WeatherBot", "gpt-4", llm.OpenAI).
    WithInstructions("You provide weather information.").
    WithFunctions([]swarm.AgentFunction{
        {
            Name:        "get_weather",
            Description: "Get weather for a location",
            Parameters: map[string]interface{}{
                "type": "object",
                "properties": map[string]interface{}{
                    "location": map[string]interface{}{
                        "type": "string",
                    },
                },
            },
            Function: func(args, ctx map[string]interface{}) swarm.Result {
                location := args["location"].(string)
                return swarm.Result{
                    Data: fmt.Sprintf("Weather in %s: Sunny, 72°F", location),
                }
            },
        },
    })
```

## Architecture

The TUI is built with:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [SwarmGo](../internal/swarm/) - Multi-agent orchestration

## Design Principles

1. **No Setup Required** - Config viewer works without any API keys
2. **Centered Layout** - Fixed 60-character width for clean appearance
3. **Visual Feedback** - Color-coded status indicators
4. **Responsive** - Adapts to terminal size
5. **Consistent** - Same formatting across all sections

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

- **OpenAI**: Fully functional
- **Claude/Gemini**: Currently using stub implementations
- **Ollama**: Requires local Ollama server running

## Development

To modify the TUI:

1. Edit `internal/tui/config_tui.go` for config viewer
2. Edit `internal/tui/tui.go` for chat interface
3. Edit `internal/cli/tui.go` for CLI integration

Run tests:

```bash
go test ./internal/tui/...
```

## Future Enhancements

- [ ] Search/filter in configuration
- [ ] Edit configuration values directly
- [ ] Export/import configuration
- [ ] Connection testing for services
- [ ] Theme customization
