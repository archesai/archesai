package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// Step represents a single step in a multi-step operation.
type Step struct {
	ID          string
	Title       string
	Description string
	Status      StepStatus
	StartTime   time.Time
	EndTime     time.Time
	Error       error
	Output      string
}

// StepStatus represents the status of a step.
type StepStatus int

const (
	StepPending StepStatus = iota
	StepRunning
	StepSuccess
	StepWarning
	StepError
	StepSkipped
)

func (s StepStatus) String() string {
	switch s {
	case StepPending:
		return "pending"
	case StepRunning:
		return "running"
	case StepSuccess:
		return "success"
	case StepWarning:
		return "warning"
	case StepError:
		return "error"
	case StepSkipped:
		return "skipped"
	default:
		return "unknown"
	}
}

// StepsModel displays a list of steps with their status.
type StepsModel struct {
	title   string
	steps   []Step
	spinner spinner.Model
	styles  Styles
	width   int
	done    bool
	current int
}

// StepUpdateMsg updates a step's status.
type StepUpdateMsg struct {
	ID     string
	Status StepStatus
	Output string
	Error  error
}

// StepStartMsg starts a step.
type StepStartMsg struct {
	ID string
}

// StepsDoneMsg signals all steps are complete.
type StepsDoneMsg struct{}

// NewSteps creates a new steps model.
func NewSteps(title string, steps []Step) StepsModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = DefaultStyles().Info

	return StepsModel{
		title:   title,
		steps:   steps,
		spinner: s,
		styles:  DefaultStyles(),
		current: -1,
	}
}

// NewStepsFromStrings creates steps from simple string titles.
func NewStepsFromStrings(title string, stepTitles ...string) StepsModel {
	steps := make([]Step, len(stepTitles))
	for i, t := range stepTitles {
		steps[i] = Step{
			ID:     fmt.Sprintf("step-%d", i),
			Title:  t,
			Status: StepPending,
		}
	}
	return NewSteps(title, steps)
}

// WithStyles sets custom styles.
func (m StepsModel) WithStyles(styles Styles) StepsModel {
	m.styles = styles
	m.spinner.Style = styles.Info
	return m
}

// Init initializes the model.
func (m StepsModel) Init() tea.Cmd {
	return m.spinner.Tick
}

// Update handles messages.
func (m StepsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width

	case StepStartMsg:
		for i := range m.steps {
			if m.steps[i].ID == msg.ID {
				m.steps[i].Status = StepRunning
				m.steps[i].StartTime = time.Now()
				m.current = i
				break
			}
		}

	case StepUpdateMsg:
		for i := range m.steps {
			if m.steps[i].ID == msg.ID {
				m.steps[i].Status = msg.Status
				m.steps[i].Output = msg.Output
				m.steps[i].Error = msg.Error
				if msg.Status != StepRunning && msg.Status != StepPending {
					m.steps[i].EndTime = time.Now()
				}
				break
			}
		}
		// Check if all done
		allDone := true
		for _, step := range m.steps {
			if step.Status == StepPending || step.Status == StepRunning {
				allDone = false
				break
			}
		}
		if allDone {
			m.done = true
			return m, tea.Quit
		}

	case StepsDoneMsg:
		m.done = true
		return m, tea.Quit

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View renders the steps.
func (m StepsModel) View() string {
	var sb strings.Builder

	// Title
	sb.WriteString(m.styles.Title.Render(m.title))
	sb.WriteString("\n\n")

	// Steps
	for i, step := range m.steps {
		// Status icon
		var icon string
		switch step.Status {
		case StepRunning:
			icon = m.spinner.View()
		case StepSuccess:
			icon = m.styles.Success.Render(m.styles.IconSuccess)
		case StepWarning:
			icon = m.styles.Warning.Render(m.styles.IconWarning)
		case StepError:
			icon = m.styles.Error.Render(m.styles.IconError)
		case StepSkipped:
			icon = m.styles.Muted.Render("○")
		default:
			icon = m.styles.Muted.Render(m.styles.IconPending)
		}

		// Step number and title
		stepNum := m.styles.Muted.Render(fmt.Sprintf("%d.", i+1))
		title := step.Title
		switch step.Status {
		case StepRunning:
			title = m.styles.Info.Render(title)
		case StepSuccess:
			title = m.styles.Value.Render(title)
		case StepError:
			title = m.styles.Error.Render(title)
		default:
			title = m.styles.Muted.Render(title)
		}

		sb.WriteString(fmt.Sprintf("  %s %s %s", icon, stepNum, title))

		// Duration for completed steps
		if !step.EndTime.IsZero() && !step.StartTime.IsZero() {
			duration := step.EndTime.Sub(step.StartTime).Round(time.Millisecond)
			sb.WriteString(" " + m.styles.Muted.Render(fmt.Sprintf("(%s)", duration)))
		}

		sb.WriteString("\n")

		// Description or output
		if step.Description != "" && step.Status == StepRunning {
			sb.WriteString(fmt.Sprintf("       %s\n", m.styles.Muted.Render(step.Description)))
		}
		if step.Output != "" {
			sb.WriteString(fmt.Sprintf("       %s\n", m.styles.Muted.Render(step.Output)))
		}
		if step.Error != nil {
			sb.WriteString(fmt.Sprintf("       %s\n", m.styles.Error.Render(step.Error.Error())))
		}
	}

	// Summary when done
	if m.done {
		sb.WriteString("\n")
		successCount := 0
		errorCount := 0
		for _, step := range m.steps {
			switch step.Status {
			case StepSuccess:
				successCount++
			case StepError:
				errorCount++
			}
		}
		if errorCount > 0 {
			sb.WriteString(
				m.styles.Error.Render(fmt.Sprintf("Completed with %d error(s)", errorCount)),
			)
		} else {
			sb.WriteString(m.styles.Success.Render(fmt.Sprintf("All %d steps completed successfully", successCount)))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// RunSteps runs a series of steps with a task function for each.
func RunSteps(title string, steps []StepDef, task func(stepID string) (string, error)) error {
	stepList := make([]Step, len(steps))
	for i, s := range steps {
		stepList[i] = Step{
			ID:          s.ID,
			Title:       s.Title,
			Description: s.Description,
			Status:      StepPending,
		}
	}

	model := NewSteps(title, stepList)
	p := tea.NewProgram(model)

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

// StepDef defines a step for RunSteps.
type StepDef struct {
	ID          string
	Title       string
	Description string
}

// StepRunner helps run steps programmatically.
type StepRunner struct {
	program *tea.Program
	model   StepsModel
}

// NewStepRunner creates a step runner for programmatic control.
func NewStepRunner(title string, steps []Step) *StepRunner {
	model := NewSteps(title, steps)
	return &StepRunner{
		model: model,
	}
}

// Start begins the TUI.
func (r *StepRunner) Start() {
	r.program = tea.NewProgram(r.model)
	go func() {
		_, _ = r.program.Run()
	}()
}

// StartStep marks a step as running.
func (r *StepRunner) StartStep(id string) {
	if r.program != nil {
		r.program.Send(StepStartMsg{ID: id})
	}
}

// CompleteStep marks a step as complete.
func (r *StepRunner) CompleteStep(id string, output string) {
	if r.program != nil {
		r.program.Send(StepUpdateMsg{
			ID:     id,
			Status: StepSuccess,
			Output: output,
		})
	}
}

// FailStep marks a step as failed.
func (r *StepRunner) FailStep(id string, err error) {
	if r.program != nil {
		r.program.Send(StepUpdateMsg{
			ID:     id,
			Status: StepError,
			Error:  err,
		})
	}
}

// SkipStep marks a step as skipped.
func (r *StepRunner) SkipStep(id string) {
	if r.program != nil {
		r.program.Send(StepUpdateMsg{
			ID:     id,
			Status: StepSkipped,
		})
	}
}

// Done signals completion.
func (r *StepRunner) Done() {
	if r.program != nil {
		r.program.Send(StepsDoneMsg{})
	}
}
