package tui

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/archesai/archesai/internal/dev"
)

// Tab represents a log view tab
type Tab int

// Tab constants for different log views
const (
	TabAPI      Tab = iota // API service logs only
	TabPlatform            // Platform service logs only
)

func (t Tab) String() string {
	switch t {
	case TabAPI:
		return "API"
	case TabPlatform:
		return "Platform"
	default:
		return "Unknown"
	}
}

// LogEntry represents a simple log entry
type LogEntry struct {
	Time    time.Time
	Service string
	Message string
	Level   string
}

// DevModel represents the TUI model for dev command
type DevModel struct {
	manager   *dev.Manager
	viewport  viewport.Model
	spinner   spinner.Model
	activeTab Tab
	width     int
	height    int
	logs      []LogEntry
	logMutex  sync.RWMutex
	ready     bool
	quitting  bool
	styles    DevStyles
}

// DevStyles contains styles for the dev TUI
type DevStyles struct {
	TabBar       lipgloss.Style
	ActiveTab    lipgloss.Style
	InactiveTab  lipgloss.Style
	StatusBar    lipgloss.Style
	LogContainer lipgloss.Style
	ErrorLog     lipgloss.Style
	WarnLog      lipgloss.Style
	InfoLog      lipgloss.Style
	DebugLog     lipgloss.Style
	ServiceTag   lipgloss.Style
	Help         lipgloss.Style
}

// NewDevStyles creates default styles for the dev TUI
func NewDevStyles() DevStyles {
	return DevStyles{
		TabBar: lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(lipgloss.Color("240")),
		ActiveTab: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			Background(lipgloss.Color("235")).
			Padding(0, 2),
		InactiveTab: lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			Padding(0, 2),
		StatusBar: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Background(lipgloss.Color("235")).
			Padding(0, 1),
		LogContainer: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1),
		ErrorLog: lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")),
		WarnLog: lipgloss.NewStyle().
			Foreground(lipgloss.Color("214")),
		InfoLog: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")),
		DebugLog: lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")),
		ServiceTag: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("141")),
		Help: lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")),
	}
}

// NewDevModel creates a new TUI model for the dev command
func NewDevModel(manager *dev.Manager) *DevModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	model := &DevModel{
		manager:   manager,
		spinner:   s,
		activeTab: TabAPI,
		styles:    NewDevStyles(),
		logs:      []LogEntry{},
	}

	// Start goroutines to capture process output
	go model.captureProcessOutput()

	return model
}

// captureProcessOutput captures output from all processes
func (m *DevModel) captureProcessOutput() {
	processes := m.manager.GetProcesses()

	for name, process := range processes {
		// Capture output from each process
		go func(serviceName string, p *dev.Process) {
			outputCh := p.GetOutput()
			for line := range outputCh {
				m.addLog(serviceName, line)
			}
		}(name, process)

		// Capture errors from each process
		go func(serviceName string, p *dev.Process) {
			errorCh := p.GetErrors()
			for err := range errorCh {
				m.addLog(serviceName, fmt.Sprintf("ERROR: %v", err))
			}
		}(name, process)
	}
}

// addLog adds a log entry
func (m *DevModel) addLog(service, message string) {
	m.logMutex.Lock()
	defer m.logMutex.Unlock()

	// Detect log level from message
	level := "INFO"
	lowerMsg := strings.ToLower(message)
	if strings.Contains(lowerMsg, "error") {
		level = "ERROR"
	} else if strings.Contains(lowerMsg, "warn") {
		level = "WARN"
	} else if strings.Contains(lowerMsg, "debug") {
		level = "DEBUG"
	}

	entry := LogEntry{
		Time:    time.Now(),
		Service: service,
		Message: message,
		Level:   level,
	}

	m.logs = append(m.logs, entry)

	// Keep only last 1000 logs
	if len(m.logs) > 1000 {
		m.logs = m.logs[len(m.logs)-1000:]
	}
}

// Init initializes the model
func (m *DevModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.tickCmd(),
	)
}

// tickCmd returns a command that ticks periodically to update the view
func (m *DevModel) tickCmd() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type tickMsg time.Time

// Update handles messages
func (m *DevModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "left", "h", "1":
			m.activeTab = TabAPI
			m.updateViewportContent()

		case "right", "l", "2":
			m.activeTab = TabPlatform
			m.updateViewportContent()

		case "tab":
			// Toggle between tabs
			if m.activeTab == TabAPI {
				m.activeTab = TabPlatform
			} else {
				m.activeTab = TabAPI
			}
			m.updateViewportContent()

		case "c":
			// Clear logs
			m.logMutex.Lock()
			m.logs = []LogEntry{}
			m.logMutex.Unlock()
			m.updateViewportContent()

		case "r":
			// Restart current service
			switch m.activeTab {
			case TabAPI:
				go func() {
					_ = m.manager.RestartProcess("api")
				}()
			case TabPlatform:
				go func() {
					_ = m.manager.RestartProcess("platform")
				}()
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			// Initialize viewport
			m.viewport = viewport.New(msg.Width, msg.Height-5) // Leave room for tab bar and status
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - 5
		}
		m.updateViewportContent()

	case spinner.TickMsg:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case tickMsg:
		// Update viewport content periodically
		m.updateViewportContent()
		cmds = append(cmds, m.tickCmd())
	}

	// Update viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the UI
func (m *DevModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	var s strings.Builder

	// Tab bar
	s.WriteString(m.renderTabBar())
	s.WriteString("\n")

	// Main content (viewport)
	s.WriteString(m.viewport.View())
	s.WriteString("\n")

	// Status bar
	s.WriteString(m.renderStatusBar())

	return s.String()
}

// renderTabBar renders the tab navigation bar
func (m *DevModel) renderTabBar() string {
	tabs := []string{}
	for i := 0; i < 2; i++ {
		tab := Tab(i)
		label := fmt.Sprintf("[%d] %s", i+1, tab.String())
		if tab == m.activeTab {
			tabs = append(tabs, m.styles.ActiveTab.Render(label))
		} else {
			tabs = append(tabs, m.styles.InactiveTab.Render(label))
		}
	}

	tabBar := strings.Join(tabs, " ")
	return m.styles.TabBar.Width(m.width).Render(tabBar)
}

// renderStatusBar renders the bottom status bar
func (m *DevModel) renderStatusBar() string {
	// Get process states
	processes := m.manager.GetProcesses()
	apiState := "stopped"
	platformState := "stopped"

	if proc, ok := processes["api"]; ok {
		apiState = proc.GetState().String()
	}
	if proc, ok := processes["platform"]; ok {
		platformState = proc.GetState().String()
	}

	m.logMutex.RLock()
	logCount := len(m.logs)
	m.logMutex.RUnlock()

	left := fmt.Sprintf(" API: %s | Platform: %s | Logs: %d",
		m.getStateIcon(apiState)+" "+apiState,
		m.getStateIcon(platformState)+" "+platformState,
		logCount,
	)

	right := " ←/→: Switch Tab | C: Clear | R: Restart | Q: Quit "

	width := m.width - lipgloss.Width(left) - lipgloss.Width(right)
	if width < 0 {
		width = 0
	}

	status := left + strings.Repeat(" ", width) + right
	return m.styles.StatusBar.Width(m.width).Render(status)
}

// getStateIcon returns an icon for the process state
func (m *DevModel) getStateIcon(state string) string {
	switch state {
	case "running":
		return "●"
	case "starting":
		return m.spinner.View()
	case "stopping":
		return "◌"
	case "error":
		return "✗"
	default:
		return "○"
	}
}

// updateViewportContent updates the viewport with filtered logs
func (m *DevModel) updateViewportContent() {
	m.logMutex.RLock()
	defer m.logMutex.RUnlock()

	var logs []LogEntry

	switch m.activeTab {
	case TabAPI:
		logs = m.filterLogsByService("api")
	case TabPlatform:
		logs = m.filterLogsByService("platform")
	}

	// Format logs for display
	var content strings.Builder
	for _, entry := range logs {
		content.WriteString(m.formatLogEntry(entry))
		content.WriteString("\n")
	}

	m.viewport.SetContent(content.String())
}

// filterLogsByService filters logs by service name
func (m *DevModel) filterLogsByService(service string) []LogEntry {
	filtered := []LogEntry{}
	for _, log := range m.logs {
		if log.Service == service {
			filtered = append(filtered, log)
		}
	}
	return filtered
}

// formatLogEntry formats a single log entry for display
func (m *DevModel) formatLogEntry(entry LogEntry) string {
	// Format timestamp
	timestamp := entry.Time.Format("15:04:05.000")

	// Format service tag
	serviceTag := m.styles.ServiceTag.Render(fmt.Sprintf("[%s]", entry.Service))

	// Apply color based on log level
	var logStyle lipgloss.Style
	switch entry.Level {
	case "ERROR":
		logStyle = m.styles.ErrorLog
	case "WARN":
		logStyle = m.styles.WarnLog
	case "DEBUG":
		logStyle = m.styles.DebugLog
	default:
		logStyle = m.styles.InfoLog
	}

	// Format the complete line
	return fmt.Sprintf("%s %s %s",
		timestamp,
		serviceTag,
		logStyle.Render(entry.Message),
	)
}

// RunDevTUI runs the dev TUI interface
func RunDevTUI(manager *dev.Manager) error {
	// Create and run the TUI model
	model := NewDevModel(manager)
	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err := p.Run()
	return err
}
