package tui

import (
	"fmt"
	"io"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Runner provides easy integration of TUI components into CLI commands.
type Runner struct {
	styles    Styles
	altScreen bool
	output    io.Writer
}

// NewRunner creates a new TUI runner.
func NewRunner() *Runner {
	return &Runner{
		styles:    DefaultStyles(),
		altScreen: false,
		output:    os.Stdout,
	}
}

// WithStyles sets custom styles.
func (r *Runner) WithStyles(styles Styles) *Runner {
	r.styles = styles
	return r
}

// WithAltScreen enables alternate screen mode.
func (r *Runner) WithAltScreen() *Runner {
	r.altScreen = true
	return r
}

// WithOutput sets the output writer.
func (r *Runner) WithOutput(w io.Writer) *Runner {
	r.output = w
	return r
}

// Run executes a tea.Model.
func (r *Runner) Run(model tea.Model) error {
	opts := []tea.ProgramOption{
		tea.WithOutput(r.output),
	}
	if r.altScreen {
		opts = append(opts, tea.WithAltScreen())
	}

	p := tea.NewProgram(model, opts...)
	_, err := p.Run()
	return err
}

// Spinner runs a spinner while executing a task.
func (r *Runner) Spinner(message string, task func() (string, error)) error {
	model := NewSpinner(message).WithStyles(r.styles)

	opts := []tea.ProgramOption{
		tea.WithOutput(r.output),
	}
	if r.altScreen {
		opts = append(opts, tea.WithAltScreen())
	}

	p := tea.NewProgram(model, opts...)

	go func() {
		result, err := task()
		p.Send(SpinnerDoneMsg{Result: result, Err: err})
	}()

	_, err := p.Run()
	return err
}

// Progress runs a progress bar while executing a task.
func (r *Runner) Progress(
	total int,
	message string,
	task func(update func(current int, msg string)) error,
) error {
	model := NewProgress(total, message).WithStyles(r.styles)

	opts := []tea.ProgramOption{
		tea.WithOutput(r.output),
	}
	if r.altScreen {
		opts = append(opts, tea.WithAltScreen())
	}

	p := tea.NewProgram(model, opts...)

	go func() {
		err := task(func(current int, msg string) {
			p.Send(ProgressUpdateMsg{Current: current, Message: msg})
		})
		p.Send(ProgressDoneMsg{Err: err})
	}()

	_, err := p.Run()
	return err
}

// Steps runs a multi-step operation.
func (r *Runner) Steps(
	title string,
	steps []StepDef,
	task func(stepID string) (string, error),
) error {
	stepList := make([]Step, len(steps))
	for i, s := range steps {
		stepList[i] = Step{
			ID:          s.ID,
			Title:       s.Title,
			Description: s.Description,
			Status:      StepPending,
		}
	}

	model := NewSteps(title, stepList).WithStyles(r.styles)

	opts := []tea.ProgramOption{
		tea.WithOutput(r.output),
	}
	if r.altScreen {
		opts = append(opts, tea.WithAltScreen())
	}

	p := tea.NewProgram(model, opts...)

	go func() {
		for _, step := range steps {
			p.Send(StepStartMsg{ID: step.ID})
			output, err := task(step.ID)
			if err != nil {
				p.Send(StepUpdateMsg{
					ID:     step.ID,
					Status: StepError,
					Error:  err,
				})
				break
			}
			p.Send(StepUpdateMsg{
				ID:     step.ID,
				Status: StepSuccess,
				Output: output,
			})
		}
		p.Send(StepsDoneMsg{})
	}()

	_, err := p.Run()
	return err
}

// Print outputs styled text without running a full TUI.
func (r *Runner) Print(text string) {
	_, _ = fmt.Fprintln(r.output, text)
}

// PrintSuccess outputs a success message.
func (r *Runner) PrintSuccess(text string) {
	_, _ = fmt.Fprintln(r.output, r.styles.RenderStatus("success", text))
}

// PrintError outputs an error message.
func (r *Runner) PrintError(text string) {
	_, _ = fmt.Fprintln(r.output, r.styles.RenderStatus("error", text))
}

// PrintWarning outputs a warning message.
func (r *Runner) PrintWarning(text string) {
	_, _ = fmt.Fprintln(r.output, r.styles.RenderStatus("warning", text))
}

// PrintInfo outputs an info message.
func (r *Runner) PrintInfo(text string) {
	_, _ = fmt.Fprintln(r.output, r.styles.RenderStatus("info", text))
}

// PrintTitle outputs a styled title.
func (r *Runner) PrintTitle(text string) {
	_, _ = fmt.Fprintln(r.output, r.styles.Title.Render(text))
}

// PrintMuted outputs muted text.
func (r *Runner) PrintMuted(text string) {
	_, _ = fmt.Fprintln(r.output, r.styles.Muted.Render(text))
}

// PrintKeyValue outputs a key-value pair.
func (r *Runner) PrintKeyValue(key, value string) {
	_, _ = fmt.Fprintln(r.output, r.styles.RenderKeyValue(key, value))
}

// PrintResult outputs a result display.
func (r *Runner) PrintResult(result *ResultModel) {
	_, _ = fmt.Fprint(r.output, result.Render())
}

// PrintSummary outputs a summary display.
func (r *Runner) PrintSummary(summary *SummaryModel) {
	_, _ = fmt.Fprint(r.output, summary.Render())
}

// PrintTable outputs a table display.
func (r *Runner) PrintTable(table *TableModel) {
	_, _ = fmt.Fprint(r.output, table.Render())
}

// PrintDivider outputs a horizontal divider.
func (r *Runner) PrintDivider() {
	_, _ = fmt.Fprintln(r.output, r.styles.Muted.Render(strings.Repeat("─", 40)))
}

// PrintNewline outputs a blank line.
func (r *Runner) PrintNewline() {
	_, _ = fmt.Fprintln(r.output)
}

// Styles returns the runner's styles for custom rendering.
func (r *Runner) Styles() Styles {
	return r.styles
}
