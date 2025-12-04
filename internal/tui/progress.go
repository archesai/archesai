package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

// ProgressModel displays a progress bar with optional message.
type ProgressModel struct {
	progress progress.Model
	message  string
	current  int
	total    int
	done     bool
	err      error
	result   string
	styles   Styles
	width    int
}

// ProgressUpdateMsg updates progress.
type ProgressUpdateMsg struct {
	Current int
	Total   int
	Message string
}

// ProgressDoneMsg signals completion.
type ProgressDoneMsg struct {
	Result string
	Err    error
}

// ProgressIncrementMsg increments progress by 1.
type ProgressIncrementMsg struct {
	Message string
}

// NewProgress creates a new progress bar model.
func NewProgress(total int, message string) ProgressModel {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
	)

	return ProgressModel{
		progress: p,
		message:  message,
		total:    total,
		styles:   DefaultStyles(),
	}
}

// WithWidth sets the progress bar width.
func (m ProgressModel) WithWidth(width int) ProgressModel {
	m.progress.Width = width
	return m
}

// WithGradient sets custom gradient colors.
func (m ProgressModel) WithGradient(colorA, colorB string) ProgressModel {
	m.progress = progress.New(
		progress.WithGradient(colorA, colorB),
		progress.WithWidth(m.progress.Width),
	)
	return m
}

// WithSolidFill sets a solid fill color.
func (m ProgressModel) WithSolidFill(color string) ProgressModel {
	m.progress = progress.New(
		progress.WithSolidFill(color),
		progress.WithWidth(m.progress.Width),
	)
	return m
}

// WithStyles sets custom styles.
func (m ProgressModel) WithStyles(styles Styles) ProgressModel {
	m.styles = styles
	return m
}

// Init initializes the progress bar.
func (m ProgressModel) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", KeyQuit:
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.progress.Width = msg.Width - 20
		if m.progress.Width < 20 {
			m.progress.Width = 20
		}
		if m.progress.Width > 60 {
			m.progress.Width = 60
		}

	case ProgressUpdateMsg:
		m.current = msg.Current
		if msg.Total > 0 {
			m.total = msg.Total
		}
		if msg.Message != "" {
			m.message = msg.Message
		}

	case ProgressIncrementMsg:
		m.current++
		if msg.Message != "" {
			m.message = msg.Message
		}
		if m.current >= m.total {
			m.done = true
			return m, tea.Quit
		}

	case ProgressDoneMsg:
		m.done = true
		m.err = msg.Err
		m.result = msg.Result
		return m, tea.Quit

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	return m, nil
}

// View renders the progress bar.
func (m ProgressModel) View() string {
	if m.done {
		if m.err != nil {
			return m.styles.RenderStatus("error", fmt.Sprintf("Failed: %v", m.err)) + "\n"
		}
		if m.result != "" {
			return m.styles.RenderStatus("success", m.result) + "\n"
		}
		return m.styles.RenderStatus("success", fmt.Sprintf("Completed %d items", m.total)) + "\n"
	}

	percent := 0.0
	if m.total > 0 {
		percent = float64(m.current) / float64(m.total)
	}

	var sb strings.Builder
	sb.WriteString(m.message)
	sb.WriteString("\n")
	sb.WriteString(m.progress.ViewAs(percent))
	sb.WriteString(fmt.Sprintf(" %d/%d", m.current, m.total))
	sb.WriteString("\n")

	return sb.String()
}

// FileProgressModel tracks file generation/processing progress.
type FileProgressModel struct {
	progress progress.Model
	files    []FileItem
	current  int
	done     bool
	err      error
	styles   Styles
	width    int
	showList bool
	maxShow  int
}

// FileItem represents a file being processed.
type FileItem struct {
	Path   string
	Status string // "pending", "processing", "done", "error"
}

// FileProgressUpdateMsg adds or updates a file.
type FileProgressUpdateMsg struct {
	Path   string
	Status string
}

// FileProgressDoneMsg signals completion.
type FileProgressDoneMsg struct {
	Err error
}

// NewFileProgress creates a new file progress model.
func NewFileProgress(showList bool) FileProgressModel {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(50),
	)

	return FileProgressModel{
		progress: p,
		files:    []FileItem{},
		styles:   DefaultStyles(),
		showList: showList,
		maxShow:  10,
	}
}

// WithMaxShow sets max files to show in list.
func (m FileProgressModel) WithMaxShow(maxFiles int) FileProgressModel {
	m.maxShow = maxFiles
	return m
}

// WithStyles sets custom styles.
func (m FileProgressModel) WithStyles(styles Styles) FileProgressModel {
	m.styles = styles
	return m
}

// Init initializes the model.
func (m FileProgressModel) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (m FileProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", KeyQuit:
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.progress.Width = msg.Width - 20
		if m.progress.Width < 20 {
			m.progress.Width = 20
		}
		if m.progress.Width > 60 {
			m.progress.Width = 60
		}

	case FileProgressUpdateMsg:
		found := false
		for i := range m.files {
			if m.files[i].Path == msg.Path {
				m.files[i].Status = msg.Status
				found = true
				break
			}
		}
		if !found {
			m.files = append(m.files, FileItem(msg))
		}
		// Count completed
		m.current = 0
		for _, f := range m.files {
			if f.Status == StatusDone || f.Status == StatusError {
				m.current++
			}
		}

	case FileProgressDoneMsg:
		m.done = true
		m.err = msg.Err
		return m, tea.Quit

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	return m, nil
}

// View renders the file progress.
func (m FileProgressModel) View() string {
	var sb strings.Builder

	total := len(m.files)
	if total == 0 {
		return m.styles.Muted.Render("Waiting for files...") + "\n"
	}

	if m.done {
		if m.err != nil {
			sb.WriteString(m.styles.RenderStatus("error", fmt.Sprintf("Failed: %v", m.err)))
		} else {
			sb.WriteString(m.styles.RenderStatus("success", fmt.Sprintf("Generated %d files", total)))
		}
		sb.WriteString("\n")
		return sb.String()
	}

	// Progress bar
	percent := 0.0
	if total > 0 {
		percent = float64(m.current) / float64(total)
	}
	sb.WriteString(m.progress.ViewAs(percent))
	sb.WriteString(fmt.Sprintf(" %d/%d files\n", m.current, total))

	// File list (if enabled)
	if m.showList {
		sb.WriteString("\n")
		start := 0
		if len(m.files) > m.maxShow {
			start = len(m.files) - m.maxShow
		}
		for i := start; i < len(m.files); i++ {
			f := m.files[i]
			icon := m.styles.RenderIcon(f.Status)
			path := f.Path
			// Truncate long paths
			if len(path) > m.width-10 && m.width > 20 {
				path = "..." + path[len(path)-(m.width-13):]
			}
			sb.WriteString(fmt.Sprintf("  %s %s\n", icon, path))
		}
		if start > 0 {
			sb.WriteString(m.styles.Muted.Render(fmt.Sprintf("  ... and %d more\n", start)))
		}
	}

	return sb.String()
}

// RunProgress runs a progress bar with incremental updates.
func RunProgress(
	total int,
	message string,
	task func(update func(current int, msg string)) error,
) error {
	model := NewProgress(total, message)

	p := tea.NewProgram(model)

	go func() {
		err := task(func(current int, msg string) {
			p.Send(ProgressUpdateMsg{Current: current, Message: msg})
		})
		p.Send(ProgressDoneMsg{Err: err})
	}()

	_, err := p.Run()
	return err
}
