package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SpinnerStyle defines spinner animation styles.
type SpinnerStyle int

// Spinner animation style constants.
const (
	SpinnerDot     SpinnerStyle = iota // Dot animation (default)
	SpinnerLine                        // Line animation
	SpinnerMiniDot                     // Mini dot animation
	SpinnerJump                        // Jump animation
	SpinnerPulse                       // Pulse animation
	SpinnerPoints                      // Points animation
	SpinnerGlobe                       // Globe animation
	SpinnerMoon                        // Moon phases animation
	SpinnerMonkey                      // Monkey animation
)

// SpinnerModel is a simple spinner with a message.
type SpinnerModel struct {
	spinner spinner.Model
	message string
	done    bool
	err     error
	result  string
	styles  Styles
	width   int
}

// SpinnerDoneMsg signals the spinner operation is complete.
type SpinnerDoneMsg struct {
	Result string
	Err    error
}

// SpinnerUpdateMsg updates the spinner message.
type SpinnerUpdateMsg struct {
	Message string
}

// NewSpinner creates a new spinner model.
func NewSpinner(message string) SpinnerModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	return SpinnerModel{
		spinner: s,
		message: message,
		styles:  DefaultStyles(),
	}
}

// WithStyle sets the spinner animation style.
func (m SpinnerModel) WithStyle(style SpinnerStyle) SpinnerModel {
	switch style {
	case SpinnerDot:
		m.spinner.Spinner = spinner.Dot
	case SpinnerLine:
		m.spinner.Spinner = spinner.Line
	case SpinnerMiniDot:
		m.spinner.Spinner = spinner.MiniDot
	case SpinnerJump:
		m.spinner.Spinner = spinner.Jump
	case SpinnerPulse:
		m.spinner.Spinner = spinner.Pulse
	case SpinnerPoints:
		m.spinner.Spinner = spinner.Points
	case SpinnerGlobe:
		m.spinner.Spinner = spinner.Globe
	case SpinnerMoon:
		m.spinner.Spinner = spinner.Moon
	case SpinnerMonkey:
		m.spinner.Spinner = spinner.Monkey
	}
	return m
}

// WithColor sets the spinner color.
func (m SpinnerModel) WithColor(color lipgloss.Color) SpinnerModel {
	m.spinner.Style = lipgloss.NewStyle().Foreground(color)
	return m
}

// WithStyles sets custom styles.
func (m SpinnerModel) WithStyles(styles Styles) SpinnerModel {
	m.styles = styles
	m.spinner.Style = lipgloss.NewStyle().Foreground(styles.Colors.Primary)
	return m
}

// Init initializes the spinner.
func (m SpinnerModel) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update handles messages.
func (m SpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", KeyQuit:
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width

	case SpinnerUpdateMsg:
		m.message = msg.Message

	case SpinnerDoneMsg:
		m.done = true
		m.err = msg.Err
		m.result = msg.Result
		return m, tea.Quit

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View renders the spinner.
func (m SpinnerModel) View() string {
	if m.done {
		if m.err != nil {
			return m.styles.RenderStatus("error", fmt.Sprintf("Failed: %v", m.err)) + "\n"
		}
		if m.result != "" {
			return m.styles.RenderStatus("success", m.result) + "\n"
		}
		return m.styles.RenderStatus("success", "Done") + "\n"
	}

	return fmt.Sprintf("%s %s\n", m.spinner.View(), m.message)
}

// MultiSpinnerModel manages multiple concurrent spinners.
type MultiSpinnerModel struct {
	spinners []spinnerItem
	styles   Styles
	width    int
	allDone  bool
}

type spinnerItem struct {
	id      string
	spinner spinner.Model
	message string
	status  string // "pending", "running", "success", "error"
	result  string
}

// MultiSpinnerItemMsg updates a specific spinner item.
type MultiSpinnerItemMsg struct {
	ID      string
	Status  string
	Message string
	Result  string
}

// MultiSpinnerDoneMsg signals all spinners are complete.
type MultiSpinnerDoneMsg struct{}

// NewMultiSpinner creates a multi-spinner model.
func NewMultiSpinner(items map[string]string) MultiSpinnerModel {
	styles := DefaultStyles()
	spinners := make([]spinnerItem, 0, len(items))

	for id, message := range items {
		s := spinner.New()
		s.Spinner = spinner.Dot
		s.Style = lipgloss.NewStyle().Foreground(styles.Colors.Primary)

		spinners = append(spinners, spinnerItem{
			id:      id,
			spinner: s,
			message: message,
			status:  "pending",
		})
	}

	return MultiSpinnerModel{
		spinners: spinners,
		styles:   styles,
	}
}

// Init initializes all spinners.
func (m MultiSpinnerModel) Init() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.spinners))
	for i := range m.spinners {
		cmds[i] = m.spinners[i].spinner.Tick
	}
	return tea.Batch(cmds...)
}

// Update handles messages.
func (m MultiSpinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", KeyQuit:
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width

	case MultiSpinnerItemMsg:
		for i := range m.spinners {
			if m.spinners[i].id == msg.ID {
				m.spinners[i].status = msg.Status
				if msg.Message != "" {
					m.spinners[i].message = msg.Message
				}
				if msg.Result != "" {
					m.spinners[i].result = msg.Result
				}
				break
			}
		}
		// Check if all done
		allDone := true
		for _, item := range m.spinners {
			if item.status != StatusSuccess && item.status != StatusError {
				allDone = false
				break
			}
		}
		if allDone {
			m.allDone = true
			return m, tea.Quit
		}

	case MultiSpinnerDoneMsg:
		m.allDone = true
		return m, tea.Quit

	case spinner.TickMsg:
		cmds := make([]tea.Cmd, len(m.spinners))
		for i := range m.spinners {
			if m.spinners[i].status == StatusRunning {
				m.spinners[i].spinner, cmds[i] = m.spinners[i].spinner.Update(msg)
			}
		}
		return m, tea.Batch(cmds...)
	}

	return m, nil
}

// View renders all spinners.
func (m MultiSpinnerModel) View() string {
	var sb strings.Builder

	for _, item := range m.spinners {
		switch item.status {
		case StatusSuccess:
			sb.WriteString(m.styles.RenderStatus(StatusSuccess, item.message))
			if item.result != "" {
				sb.WriteString(" " + m.styles.Muted.Render(item.result))
			}
		case StatusError:
			sb.WriteString(m.styles.RenderStatus(StatusError, item.message))
			if item.result != "" {
				sb.WriteString(" " + m.styles.Error.Render(item.result))
			}
		case StatusRunning:
			sb.WriteString(fmt.Sprintf("%s %s", item.spinner.View(), item.message))
		default: // pending
			sb.WriteString(m.styles.RenderStatus(StatusPending, item.message))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// RunSpinner runs a simple spinner with a task function.
func RunSpinner(message string, task func() (string, error)) error {
	model := NewSpinner(message)

	p := tea.NewProgram(model)

	go func() {
		result, err := task()
		p.Send(SpinnerDoneMsg{Result: result, Err: err})
	}()

	_, err := p.Run()
	return err
}
