// Package tui provides terminal user interface components for Arches.
package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
)

const (
	viewMenu = "menu"
)

// ConfigModel represents the configuration TUI
type ConfigModel struct {
	sections     []string
	selectedItem int
	currentView  string // viewMenu, "database", "server", "auth", "redis", "storage", "agents"
	width        int
	height       int
	styles       *ConfigStyles
}

// ConfigStyles holds styling for the config viewer
type ConfigStyles struct {
	Title    lipgloss.Style
	Menu     lipgloss.Style
	Selected lipgloss.Style
	Key      lipgloss.Style
	Value    lipgloss.Style
	Section  lipgloss.Style
	Help     lipgloss.Style
	Success  lipgloss.Style
	Warning  lipgloss.Style
	Error    lipgloss.Style
}

// NewConfigStyles creates the config viewer styles
func NewConfigStyles() *ConfigStyles {
	return &ConfigStyles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15")).
			Background(lipgloss.AdaptiveColor{Light: "63", Dark: "57"}).
			Padding(0, 3).
			MarginTop(1).
			MarginBottom(1).
			Align(lipgloss.Center),

		Menu: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1),

		Selected: lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Background(lipgloss.Color("235")).
			Bold(true).
			PaddingLeft(1).
			PaddingRight(2),

		Key: lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")).
			Bold(true).
			Width(20),

		Value: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")),

		Section: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("213")).
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(lipgloss.Color("238")).
			MarginBottom(1).
			PaddingBottom(1),

		Help: lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Italic(true).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238")).
			Padding(0, 1).
			MarginTop(1),

		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color("82")).
			Bold(true),

		Warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")).
			Bold(true),

		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true),
	}
}

// NewConfigModel creates a new configuration viewer
func NewConfigModel() ConfigModel {
	// Initialize viper to load config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.archesai")
	viper.SetEnvPrefix("ARCHESAI")
	viper.AutomaticEnv()

	// Try to read config file (ignore errors if not found)
	_ = viper.ReadInConfig()

	return ConfigModel{
		sections: []string{
			"🗄️  Database Configuration",
			"🌐 Server Configuration",
			"🔐 Authentication Settings",
			"📦 Redis Configuration",
			"💾 Storage Settings",
			"🤖 AI Agents & LLM Providers",
			"📊 System Status",
			"🔧 Environment Variables",
		},
		currentView: viewMenu,
		styles:      NewConfigStyles(),
	}
}

// Init initializes the config model
func (m ConfigModel) Init() tea.Cmd {
	return tea.EnterAltScreen
}

// Update handles messages
func (m ConfigModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			if m.currentView != viewMenu {
				m.currentView = viewMenu
				return m, nil
			}
			return m, tea.Quit

		case tea.KeyUp:
			if m.currentView == viewMenu && m.selectedItem > 0 {
				m.selectedItem--
			}

		case tea.KeyDown:
			if m.currentView == viewMenu && m.selectedItem < len(m.sections)-1 {
				m.selectedItem++
			}

		case tea.KeyEnter:
			if m.currentView == viewMenu {
				switch m.selectedItem {
				case 0:
					m.currentView = "database"
				case 1:
					m.currentView = "server"
				case 2:
					m.currentView = "auth"
				case 3:
					m.currentView = "redis"
				case 4:
					m.currentView = "storage"
				case 5:
					m.currentView = "agents"
				case 6:
					m.currentView = "status"
				case 7:
					m.currentView = "env"
				}
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View renders the TUI
func (m ConfigModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	var sections []string

	// Title - make it smaller and centered
	title := m.styles.Title.Render("⚙️  Arches Configuration Viewer")
	titleLine := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, title)
	sections = append(sections, titleLine)

	// Main content - center it
	var mainContent string
	switch m.currentView {
	case viewMenu:
		mainContent = m.renderMenu()
	case "database":
		mainContent = m.renderDatabase()
	case "server":
		mainContent = m.renderServer()
	case "auth":
		mainContent = m.renderAuth()
	case "redis":
		mainContent = m.renderRedis()
	case "storage":
		mainContent = m.renderStorage()
	case "agents":
		mainContent = m.renderAgents()
	case "status":
		mainContent = m.renderStatus()
	case "env":
		mainContent = m.renderEnv()
	}

	// Center the main content horizontally
	centeredContent := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, mainContent)
	sections = append(sections, centeredContent)

	// Help text at bottom
	helpText := "↑/↓: Navigate │ Enter: Select │ ESC: Back │ Ctrl+C: Quit"
	help := m.styles.Help.Render(helpText)

	// Calculate remaining height and add spacer
	usedHeight := lipgloss.Height(strings.Join(sections, "\n")) + lipgloss.Height(help) + 2
	if remainingHeight := m.height - usedHeight; remainingHeight > 0 {
		sections = append(sections, strings.Repeat("\n", remainingHeight))
	}

	sections = append(sections, lipgloss.PlaceHorizontal(m.width, lipgloss.Center, help))

	return lipgloss.JoinVertical(lipgloss.Center, sections...)
}

// renderMenu renders the main menu
func (m ConfigModel) renderMenu() string {
	var items []string

	// Add a header
	header := lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")).
		Italic(true).
		Render("Select a configuration section to view:")
	items = append(items, header)
	items = append(items, "") // Empty line for spacing

	for i, section := range m.sections {
		var item string
		if i == m.selectedItem {
			// Selected item with arrow and highlighting
			item = m.styles.Selected.Render(fmt.Sprintf("  ▸ %s", section))
		} else {
			// Normal item with proper indentation
			item = fmt.Sprintf("    %s", section)
		}
		items = append(items, item)

		// Add spacing between items for better readability
		if i < len(m.sections)-1 {
			items = append(items, "")
		}
	}

	menu := strings.Join(items, "\n")

	// Use fixed width for better centering
	menuWidth := 60
	if m.width < 65 {
		menuWidth = m.width - 5
	}

	return m.styles.Menu.
		Width(menuWidth).
		Render(menu)
}

// renderDatabase renders database configuration
func (m ConfigModel) renderDatabase() string {
	var content strings.Builder

	content.WriteString(m.styles.Section.Render("📊 Database Configuration") + "\n\n")

	// PostgreSQL settings
	content.WriteString(m.styles.Key.Render("PostgreSQL") + "\n")
	content.WriteString(m.renderConfigItem("Host", viper.GetString("database.postgres.host"), false))
	content.WriteString(m.renderConfigItem("Port", fmt.Sprintf("%d", viper.GetInt("database.postgres.port")), false))
	content.WriteString(m.renderConfigItem("Database", viper.GetString("database.postgres.database"), false))
	content.WriteString(m.renderConfigItem("User", viper.GetString("database.postgres.user"), false))
	content.WriteString(m.renderConfigItem("SSL Mode", viper.GetString("database.postgres.sslmode"), false))

	// Connection pool settings
	content.WriteString("\n" + m.styles.Key.Render("Connection Pool") + "\n")
	content.WriteString(m.renderConfigItem("Max Open", fmt.Sprintf("%d", viper.GetInt("database.postgres.max_open_conns")), false))
	content.WriteString(m.renderConfigItem("Max Idle", fmt.Sprintf("%d", viper.GetInt("database.postgres.max_idle_conns")), false))

	// SQLite settings
	content.WriteString("\n" + m.styles.Key.Render("SQLite") + "\n")
	sqlitePath := viper.GetString("database.sqlite.path")
	if sqlitePath == "" {
		sqlitePath = "Not configured"
	}
	content.WriteString(m.renderConfigItem("Path", sqlitePath, false))

	// Current driver
	content.WriteString("\n" + m.styles.Key.Render("Active Driver") + "\n")
	driver := viper.GetString("database.driver")
	switch driver {
	case "postgres":
		content.WriteString("  " + m.styles.Success.Render("● PostgreSQL") + "\n")
	case "sqlite":
		content.WriteString("  " + m.styles.Warning.Render("● SQLite") + "\n")
	default:
		content.WriteString("  " + m.styles.Error.Render("○ Not configured") + "\n")
	}

	// Use fixed width for better centering
	contentWidth := 60
	if m.width < 65 {
		contentWidth = m.width - 5
	}

	return m.styles.Menu.Width(contentWidth).Render(content.String())
}

// renderConfigItem renders a configuration key-value pair with consistent formatting
func (m ConfigModel) renderConfigItem(key, value string, _ bool) string {
	if value == "" || value == "0" {
		value = lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Italic(true).Render("not set")
	} else {
		value = m.styles.Value.Render(value)
	}

	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")).
		Width(20).
		Render("  " + key + ":")

	return fmt.Sprintf("%s %s\n", keyStyle, value)
}

// renderServer renders server configuration
func (m ConfigModel) renderServer() string {
	var content strings.Builder

	content.WriteString(m.styles.Section.Render("🌐 Server Configuration") + "\n\n")

	content.WriteString(m.styles.Key.Render("API Server") + "\n")
	content.WriteString(m.renderConfigItem("Host", viper.GetString("server.host"), false))
	content.WriteString(m.renderConfigItem("Port", fmt.Sprintf("%d", viper.GetInt("server.port")), false))
	content.WriteString(m.renderConfigItem("Mode", viper.GetString("server.mode"), false))

	content.WriteString("\n" + m.styles.Key.Render("CORS Settings") + "\n")
	origins := viper.GetStringSlice("server.cors.allowed_origins")
	if len(origins) > 0 {
		for i, origin := range origins {
			if i == 0 {
				content.WriteString(m.renderConfigItem("Allowed Origins", origin, false))
			} else {
				content.WriteString(m.renderConfigItem("", origin, false))
			}
		}
	} else {
		content.WriteString(m.renderConfigItem("Allowed Origins", "none configured", false))
	}

	content.WriteString("\n" + m.styles.Key.Render("Rate Limiting") + "\n")
	enabled := viper.GetBool("server.rate_limit.enabled")
	if enabled {
		content.WriteString(m.renderConfigItem("Status", m.styles.Success.Render("Enabled"), false))
	} else {
		content.WriteString(m.renderConfigItem("Status", m.styles.Warning.Render("Disabled"), false))
	}
	content.WriteString(m.renderConfigItem("Requests/Min", fmt.Sprintf("%d", viper.GetInt("server.rate_limit.requests_per_minute")), false))

	// Use fixed width for better centering
	contentWidth := 60
	if m.width < 65 {
		contentWidth = m.width - 5
	}

	return m.styles.Menu.Width(contentWidth).Render(content.String())
}

// renderAuth renders authentication configuration
func (m ConfigModel) renderAuth() string {
	var content strings.Builder

	content.WriteString(m.styles.Section.Render("Authentication Settings") + "\n\n")

	content.WriteString(m.styles.Key.Render("JWT Configuration:") + "\n")
	secretSet := viper.GetString("auth.jwt.secret") != ""
	if secretSet {
		content.WriteString("  Secret: " + m.styles.Success.Render("✓ Configured") + "\n")
	} else {
		content.WriteString("  Secret: " + m.styles.Error.Render("✗ Not Set") + "\n")
	}
	content.WriteString(fmt.Sprintf("  Expiry: %s\n", m.styles.Value.Render(viper.GetString("auth.jwt.expiry"))))

	content.WriteString("\n" + m.styles.Key.Render("OAuth Providers:") + "\n")

	// Google OAuth
	googleEnabled := viper.GetString("auth.oauth.google.client_id") != ""
	if googleEnabled {
		content.WriteString("  Google: " + m.styles.Success.Render("✓ Configured") + "\n")
	} else {
		content.WriteString("  Google: " + m.styles.Warning.Render("○ Not Configured") + "\n")
	}

	// GitHub OAuth
	githubEnabled := viper.GetString("auth.oauth.github.client_id") != ""
	if githubEnabled {
		content.WriteString("  GitHub: " + m.styles.Success.Render("✓ Configured") + "\n")
	} else {
		content.WriteString("  GitHub: " + m.styles.Warning.Render("○ Not Configured") + "\n")
	}

	// Use fixed width for better centering
	contentWidth := 60
	if m.width < 65 {
		contentWidth = m.width - 5
	}

	return m.styles.Menu.Width(contentWidth).Render(content.String())
}

// renderRedis renders Redis configuration
func (m ConfigModel) renderRedis() string {
	var content strings.Builder

	content.WriteString(m.styles.Section.Render("Redis Configuration") + "\n\n")

	content.WriteString(m.styles.Key.Render("Connection:") + "\n")
	content.WriteString(fmt.Sprintf("  Host: %s\n", m.styles.Value.Render(viper.GetString("redis.host"))))
	content.WriteString(fmt.Sprintf("  Port: %s\n", m.styles.Value.Render(fmt.Sprintf("%d", viper.GetInt("redis.port")))))
	content.WriteString(fmt.Sprintf("  Database: %s\n", m.styles.Value.Render(fmt.Sprintf("%d", viper.GetInt("redis.db")))))

	passwordSet := viper.GetString("redis.password") != ""
	if passwordSet {
		content.WriteString("  Password: " + m.styles.Success.Render("✓ Set") + "\n")
	} else {
		content.WriteString("  Password: " + m.styles.Warning.Render("○ Not Set") + "\n")
	}

	content.WriteString("\n" + m.styles.Key.Render("Pool Settings:") + "\n")
	content.WriteString(fmt.Sprintf("  Max Retries: %s\n", m.styles.Value.Render(fmt.Sprintf("%d", viper.GetInt("redis.max_retries")))))
	content.WriteString(fmt.Sprintf("  Pool Size: %s\n", m.styles.Value.Render(fmt.Sprintf("%d", viper.GetInt("redis.pool_size")))))

	// Use fixed width for better centering
	contentWidth := 60
	if m.width < 65 {
		contentWidth = m.width - 5
	}

	return m.styles.Menu.Width(contentWidth).Render(content.String())
}

// renderStorage renders storage configuration
func (m ConfigModel) renderStorage() string {
	var content strings.Builder

	content.WriteString(m.styles.Section.Render("Storage Settings") + "\n\n")

	content.WriteString(m.styles.Key.Render("Local Storage:") + "\n")
	content.WriteString(fmt.Sprintf("  Upload Dir: %s\n", m.styles.Value.Render(viper.GetString("storage.local.upload_dir"))))
	content.WriteString(fmt.Sprintf("  Max File Size: %s MB\n", m.styles.Value.Render(fmt.Sprintf("%d", viper.GetInt("storage.local.max_file_size_mb")))))

	content.WriteString("\n" + m.styles.Key.Render("S3 Storage:") + "\n")
	s3Enabled := viper.GetString("storage.s3.bucket") != ""
	if s3Enabled {
		content.WriteString("  Status: " + m.styles.Success.Render("✓ Configured") + "\n")
		content.WriteString(fmt.Sprintf("  Bucket: %s\n", m.styles.Value.Render(viper.GetString("storage.s3.bucket"))))
		content.WriteString(fmt.Sprintf("  Region: %s\n", m.styles.Value.Render(viper.GetString("storage.s3.region"))))
	} else {
		content.WriteString("  Status: " + m.styles.Warning.Render("○ Not Configured") + "\n")
	}

	// Use fixed width for better centering
	contentWidth := 60
	if m.width < 65 {
		contentWidth = m.width - 5
	}

	return m.styles.Menu.Width(contentWidth).Render(content.String())
}

// renderAgents renders AI agents configuration
func (m ConfigModel) renderAgents() string {
	var content strings.Builder

	content.WriteString(m.styles.Section.Render("AI Agents & LLM Providers") + "\n\n")

	// Check for API keys
	content.WriteString(m.styles.Key.Render("OpenAI:") + "\n")
	openaiKey := viper.GetString("llm.openai.api_key")
	if openaiKey != "" {
		content.WriteString("  API Key: " + m.styles.Success.Render("✓ Configured") + "\n")
		content.WriteString(fmt.Sprintf("  Model: %s\n", m.styles.Value.Render(viper.GetString("llm.openai.model"))))
	} else {
		content.WriteString("  API Key: " + m.styles.Warning.Render("○ Not Set") + "\n")
	}

	content.WriteString("\n" + m.styles.Key.Render("Anthropic (Claude):") + "\n")
	claudeKey := viper.GetString("llm.anthropic.api_key")
	if claudeKey != "" {
		content.WriteString("  API Key: " + m.styles.Success.Render("✓ Configured") + "\n")
		content.WriteString(fmt.Sprintf("  Model: %s\n", m.styles.Value.Render(viper.GetString("llm.anthropic.model"))))
	} else {
		content.WriteString("  API Key: " + m.styles.Warning.Render("○ Not Set") + "\n")
	}

	content.WriteString("\n" + m.styles.Key.Render("Google (Gemini):") + "\n")
	geminiKey := viper.GetString("llm.gemini.api_key")
	if geminiKey != "" {
		content.WriteString("  API Key: " + m.styles.Success.Render("✓ Configured") + "\n")
	} else {
		content.WriteString("  API Key: " + m.styles.Warning.Render("○ Not Set") + "\n")
	}

	content.WriteString("\n" + m.styles.Key.Render("Ollama (Local):") + "\n")
	ollamaHost := viper.GetString("llm.ollama.host")
	if ollamaHost != "" {
		content.WriteString(fmt.Sprintf("  Host: %s\n", m.styles.Value.Render(ollamaHost)))
		content.WriteString("  Status: " + m.styles.Success.Render("✓ Available") + "\n")
	} else {
		content.WriteString("  Status: " + m.styles.Warning.Render("○ Not Configured") + "\n")
	}

	// SwarmGo agents status
	content.WriteString("\n" + m.styles.Key.Render("SwarmGo Agents:") + "\n")
	content.WriteString("  Status: " + m.styles.Success.Render("✓ Integrated") + "\n")
	content.WriteString("  Available Agents:\n")
	content.WriteString("    - Assistant (General Purpose)\n")
	content.WriteString("    - CodeHelper (Programming)\n")
	content.WriteString("    - CreativeWriter (Content)\n")
	content.WriteString("    - DataAnalyst (Analysis)\n")
	content.WriteString("    - Researcher (Information)\n")

	// Use fixed width for better centering
	contentWidth := 60
	if m.width < 65 {
		contentWidth = m.width - 5
	}

	return m.styles.Menu.Width(contentWidth).Render(content.String())
}

// renderStatus renders system status
func (m ConfigModel) renderStatus() string {
	var content strings.Builder

	content.WriteString(m.styles.Section.Render("System Status") + "\n\n")

	// Environment
	content.WriteString(m.styles.Key.Render("Environment:") + "\n")
	env := viper.GetString("environment")
	switch env {
	case "production":
		content.WriteString("  Mode: " + m.styles.Error.Render("Production") + "\n")
	case "development":
		content.WriteString("  Mode: " + m.styles.Warning.Render("Development") + "\n")
	default:
		content.WriteString(fmt.Sprintf("  Mode: %s\n", m.styles.Value.Render(env)))
	}

	// Debug mode
	debugMode := viper.GetBool("debug")
	if debugMode {
		content.WriteString("  Debug: " + m.styles.Warning.Render("✓ Enabled") + "\n")
	} else {
		content.WriteString("  Debug: " + m.styles.Value.Render("○ Disabled") + "\n")
	}

	// Config file
	content.WriteString("\n" + m.styles.Key.Render("Configuration:") + "\n")
	configFile := viper.ConfigFileUsed()
	if configFile != "" {
		content.WriteString(fmt.Sprintf("  Config File: %s\n", m.styles.Value.Render(configFile)))
	} else {
		content.WriteString("  Config File: " + m.styles.Warning.Render("Using defaults") + "\n")
	}

	// Feature flags
	content.WriteString("\n" + m.styles.Key.Render("Features:") + "\n")
	content.WriteString("  TUI: " + m.styles.Success.Render("✓ Enabled") + "\n")
	content.WriteString("  SwarmGo: " + m.styles.Success.Render("✓ Integrated") + "\n")
	content.WriteString("  Multi-Agent: " + m.styles.Success.Render("✓ Available") + "\n")

	// Use fixed width for better centering
	contentWidth := 60
	if m.width < 65 {
		contentWidth = m.width - 5
	}

	return m.styles.Menu.Width(contentWidth).Render(content.String())
}

// renderEnv renders environment variables
func (m ConfigModel) renderEnv() string {
	var content strings.Builder

	content.WriteString(m.styles.Section.Render("Environment Variables") + "\n\n")

	// Database
	content.WriteString(m.styles.Key.Render("Database:") + "\n")
	dbURL := viper.GetString("DATABASE_URL")
	if dbURL != "" {
		content.WriteString("  DATABASE_URL: " + m.styles.Success.Render("✓ Set") + "\n")
	} else {
		content.WriteString("  DATABASE_URL: " + m.styles.Warning.Render("○ Not Set") + "\n")
	}

	// API Keys
	content.WriteString("\n" + m.styles.Key.Render("API Keys:") + "\n")

	envVars := map[string]string{
		"OPENAI_API_KEY":        "OpenAI",
		"ANTHROPIC_API_KEY":     "Anthropic",
		"GEMINI_API_KEY":        "Gemini",
		"REDIS_URL":             "Redis URL",
		"JWT_SECRET":            "JWT Secret",
		"AWS_ACCESS_KEY_ID":     "AWS Access Key",
		"AWS_SECRET_ACCESS_KEY": "AWS Secret Key",
	}

	for env, name := range envVars {
		value := viper.GetString(env)
		if value != "" {
			content.WriteString(fmt.Sprintf("  %s: %s\n", name, m.styles.Success.Render("✓ Set")))
		} else {
			content.WriteString(fmt.Sprintf("  %s: %s\n", name, m.styles.Warning.Render("○ Not Set")))
		}
	}

	// Server
	content.WriteString("\n" + m.styles.Key.Render("Server:") + "\n")
	content.WriteString(fmt.Sprintf("  PORT: %s\n", m.styles.Value.Render(viper.GetString("PORT"))))
	content.WriteString(fmt.Sprintf("  HOST: %s\n", m.styles.Value.Render(viper.GetString("HOST"))))

	// Use fixed width for better centering
	contentWidth := 60
	if m.width < 65 {
		contentWidth = m.width - 5
	}

	return m.styles.Menu.Width(contentWidth).Render(content.String())
}
