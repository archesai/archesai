// Package tui provides terminal user interface components for Arches.
package tui

import "github.com/charmbracelet/lipgloss"

// Status string constants used throughout TUI components.
const (
	StatusSuccess    = "success"
	StatusDone       = "done"
	StatusCompleted  = "completed"
	StatusWarning    = "warning"
	StatusWarn       = "warn"
	StatusError      = "error"
	StatusFailed     = "failed"
	StatusInfo       = "info"
	StatusRunning    = "running"
	StatusInProgress = "in_progress"
	StatusPending    = "pending"
)

// Key string constants for keyboard handling.
const (
	KeyQuit  = "ctrl+c"
	KeyEnter = "enter"
)

// Colors defines the color palette for TUI components.
type Colors struct {
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Success   lipgloss.Color
	Warning   lipgloss.Color
	Error     lipgloss.Color
	Info      lipgloss.Color
	Muted     lipgloss.Color
	Text      lipgloss.Color
	TextDim   lipgloss.Color
	Border    lipgloss.Color
	BgDark    lipgloss.Color
	BgLight   lipgloss.Color
}

// DefaultColors returns the default color palette.
func DefaultColors() Colors {
	return Colors{
		Primary:   lipgloss.Color("212"), // Pink/Magenta
		Secondary: lipgloss.Color("141"), // Purple
		Success:   lipgloss.Color("82"),  // Green
		Warning:   lipgloss.Color("214"), // Orange
		Error:     lipgloss.Color("196"), // Red
		Info:      lipgloss.Color("75"),  // Blue
		Muted:     lipgloss.Color("241"), // Gray
		Text:      lipgloss.Color("252"), // Light gray
		TextDim:   lipgloss.Color("245"), // Dimmed text
		Border:    lipgloss.Color("240"), // Border gray
		BgDark:    lipgloss.Color("235"), // Dark background
		BgLight:   lipgloss.Color("238"), // Light background
	}
}

// Styles holds common styles used across TUI components.
type Styles struct {
	Colors Colors

	// Text styles
	Title    lipgloss.Style
	Subtitle lipgloss.Style
	Label    lipgloss.Style
	Value    lipgloss.Style
	Muted    lipgloss.Style
	Bold     lipgloss.Style
	Italic   lipgloss.Style

	// Status styles
	Success lipgloss.Style
	Warning lipgloss.Style
	Error   lipgloss.Style
	Info    lipgloss.Style

	// Container styles
	Box         lipgloss.Style
	BorderedBox lipgloss.Style
	StatusBar   lipgloss.Style

	// Interactive styles
	Selected   lipgloss.Style
	Unselected lipgloss.Style
	Focused    lipgloss.Style

	// Icons
	IconSuccess string
	IconWarning string
	IconError   string
	IconInfo    string
	IconPending string
	IconRunning string
	IconDone    string
}

// DefaultStyles returns the default styles.
func DefaultStyles() Styles {
	colors := DefaultColors()

	return Styles{
		Colors: colors,

		// Text styles
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(colors.Primary),

		Subtitle: lipgloss.NewStyle().
			Foreground(colors.Secondary),

		Label: lipgloss.NewStyle().
			Foreground(colors.TextDim).
			Width(20),

		Value: lipgloss.NewStyle().
			Foreground(colors.Text),

		Muted: lipgloss.NewStyle().
			Foreground(colors.Muted).
			Italic(true),

		Bold: lipgloss.NewStyle().
			Bold(true).
			Foreground(colors.Text),

		Italic: lipgloss.NewStyle().
			Italic(true).
			Foreground(colors.Text),

		// Status styles
		Success: lipgloss.NewStyle().
			Foreground(colors.Success).
			Bold(true),

		Warning: lipgloss.NewStyle().
			Foreground(colors.Warning).
			Bold(true),

		Error: lipgloss.NewStyle().
			Foreground(colors.Error).
			Bold(true),

		Info: lipgloss.NewStyle().
			Foreground(colors.Info),

		// Container styles
		Box: lipgloss.NewStyle().
			Padding(1, 2),

		BorderedBox: lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colors.Border).
			Padding(1, 2),

		StatusBar: lipgloss.NewStyle().
			Foreground(colors.TextDim).
			Background(colors.BgDark).
			Padding(0, 1),

		// Interactive styles
		Selected: lipgloss.NewStyle().
			Foreground(colors.Primary).
			Background(colors.BgDark).
			Bold(true).
			Padding(0, 1),

		Unselected: lipgloss.NewStyle().
			Foreground(colors.TextDim).
			Padding(0, 1),

		Focused: lipgloss.NewStyle().
			Foreground(colors.Primary).
			Bold(true),

		// Icons
		IconSuccess: "✓",
		IconWarning: "⚠",
		IconError:   "✗",
		IconInfo:    "ℹ",
		IconPending: "○",
		IconRunning: "●",
		IconDone:    "◉",
	}
}

// WithColors returns a new Styles with custom colors.
func (s Styles) WithColors(colors Colors) Styles {
	s.Colors = colors
	// Re-apply colors to all styles
	s.Title = s.Title.Foreground(colors.Primary)
	s.Subtitle = s.Subtitle.Foreground(colors.Secondary)
	s.Label = s.Label.Foreground(colors.TextDim)
	s.Value = s.Value.Foreground(colors.Text)
	s.Muted = s.Muted.Foreground(colors.Muted)
	s.Bold = s.Bold.Foreground(colors.Text)
	s.Italic = s.Italic.Foreground(colors.Text)
	s.Success = s.Success.Foreground(colors.Success)
	s.Warning = s.Warning.Foreground(colors.Warning)
	s.Error = s.Error.Foreground(colors.Error)
	s.Info = s.Info.Foreground(colors.Info)
	s.BorderedBox = s.BorderedBox.BorderForeground(colors.Border)
	s.StatusBar = s.StatusBar.Foreground(colors.TextDim).Background(colors.BgDark)
	s.Selected = s.Selected.Foreground(colors.Primary).Background(colors.BgDark)
	s.Unselected = s.Unselected.Foreground(colors.TextDim)
	s.Focused = s.Focused.Foreground(colors.Primary)
	return s
}

// RenderIcon returns a styled icon based on status.
func (s Styles) RenderIcon(status string) string {
	switch status {
	case StatusSuccess, StatusDone, StatusCompleted:
		return s.Success.Render(s.IconSuccess)
	case StatusWarning, StatusWarn:
		return s.Warning.Render(s.IconWarning)
	case StatusError, StatusFailed:
		return s.Error.Render(s.IconError)
	case StatusInfo:
		return s.Info.Render(s.IconInfo)
	case StatusRunning, StatusInProgress:
		return s.Info.Render(s.IconRunning)
	default:
		return s.Muted.Render(s.IconPending)
	}
}

// RenderStatus returns styled status text with icon.
func (s Styles) RenderStatus(status, text string) string {
	icon := s.RenderIcon(status)
	var style lipgloss.Style
	switch status {
	case StatusSuccess, StatusDone, StatusCompleted:
		style = s.Success
	case StatusWarning, StatusWarn:
		style = s.Warning
	case StatusError, StatusFailed:
		style = s.Error
	case StatusInfo, StatusRunning, StatusInProgress:
		style = s.Info
	default:
		style = s.Muted
	}
	return icon + " " + style.Render(text)
}

// RenderKeyValue renders a key-value pair with consistent formatting.
func (s Styles) RenderKeyValue(key, value string) string {
	return s.Label.Render(key+":") + " " + s.Value.Render(value)
}
