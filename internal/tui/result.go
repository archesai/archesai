package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ResultModel displays operation results with status indicators.
type ResultModel struct {
	title   string
	items   []ResultItem
	footer  string
	styles  Styles
	width   int
	success bool
}

// ResultItem represents a single result item.
type ResultItem struct {
	Label   string
	Value   string
	Status  string // "success", "warning", "error", "info", ""
	Details string
}

// NewResult creates a new result display.
func NewResult(title string) *ResultModel {
	return &ResultModel{
		title:   title,
		items:   []ResultItem{},
		styles:  DefaultStyles(),
		success: true,
	}
}

// AddItem adds a result item.
func (m *ResultModel) AddItem(label, value, status string) *ResultModel {
	m.items = append(m.items, ResultItem{
		Label:  label,
		Value:  value,
		Status: status,
	})
	if status == StatusError {
		m.success = false
	}
	return m
}

// AddItemWithDetails adds a result item with additional details.
func (m *ResultModel) AddItemWithDetails(label, value, status, details string) *ResultModel {
	m.items = append(m.items, ResultItem{
		Label:   label,
		Value:   value,
		Status:  status,
		Details: details,
	})
	if status == StatusError {
		m.success = false
	}
	return m
}

// SetFooter sets the footer message.
func (m *ResultModel) SetFooter(footer string) *ResultModel {
	m.footer = footer
	return m
}

// WithStyles sets custom styles.
func (m ResultModel) WithStyles(styles Styles) ResultModel {
	m.styles = styles
	return m
}

// Init initializes the model.
func (m ResultModel) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (m ResultModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", KeyQuit, KeyEnter, " ":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
	}

	return m, nil
}

// View renders the result display.
func (m ResultModel) View() string {
	var sb strings.Builder

	// Title
	titleStyle := m.styles.Title
	if !m.success {
		titleStyle = titleStyle.Foreground(m.styles.Colors.Error)
	}
	sb.WriteString(titleStyle.Render(m.title))
	sb.WriteString("\n\n")

	// Items
	for _, item := range m.items {
		icon := m.styles.RenderIcon(item.Status)
		label := m.styles.Label.Render(item.Label + ":")
		value := m.styles.Value.Render(item.Value)

		sb.WriteString(fmt.Sprintf("  %s %s %s\n", icon, label, value))

		if item.Details != "" {
			detailStyle := m.styles.Muted
			if item.Status == StatusError {
				detailStyle = m.styles.Error
			}
			sb.WriteString(fmt.Sprintf("       %s\n", detailStyle.Render(item.Details)))
		}
	}

	// Footer
	if m.footer != "" {
		sb.WriteString("\n")
		sb.WriteString(m.styles.Muted.Render(m.footer))
		sb.WriteString("\n")
	}

	return sb.String()
}

// Render returns the result as a string without running the TUI.
func (m ResultModel) Render() string {
	return m.View()
}

// SummaryModel displays a compact summary with counts.
type SummaryModel struct {
	title    string
	counts   []SummaryCount
	messages []SummaryMessage
	styles   Styles
	width    int
}

// SummaryCount represents a count in the summary.
type SummaryCount struct {
	Label  string
	Count  int
	Status string
}

// SummaryMessage represents a message in the summary.
type SummaryMessage struct {
	Text   string
	Status string
}

// NewSummary creates a new summary display.
func NewSummary(title string) *SummaryModel {
	return &SummaryModel{
		title:  title,
		styles: DefaultStyles(),
	}
}

// AddCount adds a count to the summary.
func (m *SummaryModel) AddCount(label string, count int, status string) *SummaryModel {
	m.counts = append(m.counts, SummaryCount{
		Label:  label,
		Count:  count,
		Status: status,
	})
	return m
}

// AddMessage adds a message to the summary.
func (m *SummaryModel) AddMessage(text, status string) *SummaryModel {
	m.messages = append(m.messages, SummaryMessage{
		Text:   text,
		Status: status,
	})
	return m
}

// WithStyles sets custom styles.
func (m SummaryModel) WithStyles(styles Styles) SummaryModel {
	m.styles = styles
	return m
}

// Init initializes the model.
func (m SummaryModel) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (m SummaryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", KeyQuit, KeyEnter, " ":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
	}

	return m, nil
}

// View renders the summary.
func (m SummaryModel) View() string {
	var sb strings.Builder

	// Title
	sb.WriteString(m.styles.Title.Render(m.title))
	sb.WriteString("\n\n")

	// Counts in a horizontal layout
	if len(m.counts) > 0 {
		countStrs := make([]string, len(m.counts))
		for i, c := range m.counts {
			var style lipgloss.Style
			switch c.Status {
			case StatusSuccess:
				style = m.styles.Success
			case StatusWarning:
				style = m.styles.Warning
			case StatusError:
				style = m.styles.Error
			default:
				style = m.styles.Info
			}
			countStrs[i] = style.Render(fmt.Sprintf("%s: %d", c.Label, c.Count))
		}
		sb.WriteString("  " + strings.Join(countStrs, "  │  "))
		sb.WriteString("\n")
	}

	// Messages
	if len(m.messages) > 0 {
		sb.WriteString("\n")
		for _, msg := range m.messages {
			sb.WriteString("  " + m.styles.RenderStatus(msg.Status, msg.Text) + "\n")
		}
	}

	return sb.String()
}

// Render returns the summary as a string without running the TUI.
func (m SummaryModel) Render() string {
	return m.View()
}

// TableModel displays data in a table format.
type TableModel struct {
	title   string
	headers []string
	rows    [][]string
	styles  Styles
	width   int
}

// NewTable creates a new table display.
func NewTable(title string, headers ...string) *TableModel {
	return &TableModel{
		title:   title,
		headers: headers,
		rows:    [][]string{},
		styles:  DefaultStyles(),
	}
}

// AddRow adds a row to the table.
func (m *TableModel) AddRow(values ...string) *TableModel {
	m.rows = append(m.rows, values)
	return m
}

// WithStyles sets custom styles.
func (m TableModel) WithStyles(styles Styles) TableModel {
	m.styles = styles
	return m
}

// Init initializes the model.
func (m TableModel) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", KeyQuit, KeyEnter, " ":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
	}

	return m, nil
}

// View renders the table.
func (m TableModel) View() string {
	var sb strings.Builder

	// Title
	if m.title != "" {
		sb.WriteString(m.styles.Title.Render(m.title))
		sb.WriteString("\n\n")
	}

	if len(m.headers) == 0 {
		return sb.String()
	}

	// Calculate column widths
	colWidths := make([]int, len(m.headers))
	for i, h := range m.headers {
		colWidths[i] = len(h)
	}
	for _, row := range m.rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// Header row
	headerCells := make([]string, len(m.headers))
	for i, h := range m.headers {
		headerCells[i] = m.styles.Bold.Render(padRight(h, colWidths[i]))
	}
	sb.WriteString("  " + strings.Join(headerCells, "  ") + "\n")

	// Separator
	sepParts := make([]string, len(colWidths))
	for i, w := range colWidths {
		sepParts[i] = strings.Repeat("─", w)
	}
	sb.WriteString("  " + m.styles.Muted.Render(strings.Join(sepParts, "──")) + "\n")

	// Data rows
	for _, row := range m.rows {
		cells := make([]string, len(m.headers))
		for i := range m.headers {
			val := ""
			if i < len(row) {
				val = row[i]
			}
			cells[i] = m.styles.Value.Render(padRight(val, colWidths[i]))
		}
		sb.WriteString("  " + strings.Join(cells, "  ") + "\n")
	}

	return sb.String()
}

// Render returns the table as a string without running the TUI.
func (m TableModel) Render() string {
	return m.View()
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}
