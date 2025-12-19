// Package dev provides process management for running development services
package dev

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// ProcessConfig defines the configuration for a process
type ProcessConfig struct {
	Name       string   // Name of the service (e.g., "api", "platform")
	Command    string   // Command to execute
	Args       []string // Arguments for the command
	Dir        string   // Working directory
	Env        []string // Environment variables
	HotReload  bool     // Enable hot reload for this process
	BuildCmd   string   // Command to build before running (for hot reload)
	BuildArgs  []string // Arguments for build command
	WatchPaths []string // Paths to watch for hot reload
	WatchExts  []string // File extensions to watch (e.g., .go, .mod)
}

// ProcessState represents the current state of a process
type ProcessState int

// Process state constants
const (
	ProcessStateStopped ProcessState = iota
	ProcessStateStarting
	ProcessStateRunning
	ProcessStateStopping
	ProcessStateError
)

func (s ProcessState) String() string {
	switch s {
	case ProcessStateStopped:
		return "stopped"
	case ProcessStateStarting:
		return "starting"
	case ProcessStateRunning:
		return "running"
	case ProcessStateStopping:
		return "stopping"
	case ProcessStateError:
		return "error"
	default:
		return "unknown"
	}
}

// serviceColors maps service names to lipgloss colors for log output
var serviceColors = map[string]lipgloss.Color{
	"api":      lipgloss.Color("39"),  // cyan
	"frontend": lipgloss.Color("213"), // pink
}

// Process represents a managed process
type Process struct {
	config      ProcessConfig
	cmd         *exec.Cmd
	state       ProcessState
	logger      *slog.Logger
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
	outputCh    chan string
	errorCh     chan error
	doneCh      chan struct{}
	prefixStyle lipgloss.Style
}

// NewProcess creates a new process instance
func NewProcess(config ProcessConfig, logger *slog.Logger) *Process {
	ctx, cancel := context.WithCancel(context.Background())

	// Set up colored prefix for service
	color, ok := serviceColors[config.Name]
	if !ok {
		color = lipgloss.Color("245") // default gray
	}
	prefixStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(color)

	return &Process{
		config:      config,
		state:       ProcessStateStopped,
		logger:      logger.With(slog.String("service", config.Name)),
		ctx:         ctx,
		cancel:      cancel,
		outputCh:    make(chan string, 100),
		errorCh:     make(chan error, 10),
		doneCh:      make(chan struct{}),
		prefixStyle: prefixStyle,
	}
}

// Start starts the process
func (p *Process) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state == ProcessStateRunning || p.state == ProcessStateStarting {
		return fmt.Errorf("process %s is already running or starting", p.config.Name)
	}

	p.state = ProcessStateStarting
	p.logger.Info("Starting process", "command", p.config.Command, "args", p.config.Args)

	// Create the command (don't use CommandContext as it sends SIGKILL on cancel)
	p.cmd = exec.Command(p.config.Command, p.config.Args...)

	if p.config.Dir != "" {
		p.cmd.Dir = p.config.Dir
	}

	if len(p.config.Env) > 0 {
		p.cmd.Env = append(os.Environ(), p.config.Env...)
	}

	// Create pipes for stdout and stderr
	stdout, err := p.cmd.StdoutPipe()
	if err != nil {
		p.state = ProcessStateError
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := p.cmd.StderrPipe()
	if err != nil {
		p.state = ProcessStateError
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := p.cmd.Start(); err != nil {
		p.state = ProcessStateError
		return fmt.Errorf("failed to start process: %w", err)
	}

	p.state = ProcessStateRunning

	// Start goroutines to read output
	go p.readOutput(stdout, "stdout")
	go p.readOutput(stderr, "stderr")

	// Monitor process in a goroutine
	go p.monitor()

	return nil
}

// Stop stops the process gracefully
func (p *Process) Stop() error {
	p.mu.Lock()
	if p.state != ProcessStateRunning {
		p.mu.Unlock()
		return fmt.Errorf("process %s is not running", p.config.Name)
	}

	p.state = ProcessStateStopping
	p.logger.Info("Stopping process")
	p.mu.Unlock()

	// Send SIGINT first to allow graceful shutdown
	if err := p.cmd.Process.Signal(syscall.SIGINT); err != nil {
		// Process may have already exited
		if !strings.Contains(err.Error(), "process already finished") {
			p.logger.Debug("Failed to send SIGINT", "error", err)
		}
	}

	// Wait for the process to exit (monitor goroutine will close doneCh)
	select {
	case <-time.After(5 * time.Second):
		// Force kill if graceful shutdown doesn't work
		p.logger.Warn("Process did not stop gracefully, forcing kill")
		if err := p.cmd.Process.Kill(); err != nil {
			// Ignore "already finished" errors
			if !strings.Contains(err.Error(), "process already finished") {
				return fmt.Errorf("failed to kill process: %w", err)
			}
		}
		// Wait for monitor to finish
		<-p.doneCh
	case <-p.doneCh:
		// Process exited normally
	}

	// Cancel context to clean up any remaining goroutines
	p.cancel()

	return nil
}

// Restart restarts the process
func (p *Process) Restart() error {
	p.logger.Info("Restarting process")

	// Stop if running
	if p.GetState() == ProcessStateRunning {
		if err := p.Stop(); err != nil {
			return fmt.Errorf("failed to stop process for restart: %w", err)
		}
		// Wait a bit for cleanup
		time.Sleep(500 * time.Millisecond)
	}

	// Reset context
	p.ctx, p.cancel = context.WithCancel(context.Background())
	p.doneCh = make(chan struct{})

	// Start again
	return p.Start()
}

// GetState returns the current state of the process
func (p *Process) GetState() ProcessState {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state
}

// GetOutput returns the output channel
func (p *Process) GetOutput() <-chan string {
	return p.outputCh
}

// GetErrors returns the error channel
func (p *Process) GetErrors() <-chan error {
	return p.errorCh
}

// readOutput reads output from a pipe and sends it to the output channel
func (p *Process) readOutput(pipe io.Reader, _ string) {
	scanner := bufio.NewScanner(pipe)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024) // 1MB max line size

	prefix := p.prefixStyle.Render(fmt.Sprintf("[%s]", p.config.Name))

	for scanner.Scan() {
		line := scanner.Text()

		// Skip noisy shutdown messages from pnpm/npm
		if p.GetState() == ProcessStateStopping && p.isShutdownNoise(line) {
			continue
		}

		// Print directly to stdout with colored prefix
		fmt.Printf("%s %s\n", prefix, line)

		// Also send to output channel for any consumers
		select {
		case p.outputCh <- line:
		case <-p.ctx.Done():
			return
		default:
			// Drop message if channel is full
		}
	}

	if err := scanner.Err(); err != nil {
		p.errorCh <- fmt.Errorf("error reading output: %w", err)
	}
}

// isShutdownNoise returns true if the line is expected shutdown noise that can be ignored
func (p *Process) isShutdownNoise(line string) bool {
	noisePatterns := []string{
		"ERR_PNPM_RECURSIVE_RUN_FIRST_FAIL",
		"Command failed with signal \"SIGINT\"",
		"Command failed with signal \"SIGTERM\"",
		"SIGINT",
		"SIGTERM",
	}
	for _, pattern := range noisePatterns {
		if strings.Contains(line, pattern) {
			return true
		}
	}
	return false
}

// monitor monitors the process and updates state
func (p *Process) monitor() {
	err := p.cmd.Wait()

	p.mu.Lock()
	defer p.mu.Unlock()

	if err != nil {
		errStr := err.Error()
		// Ignore expected exit signals during shutdown
		isExpectedSignal := errStr == "signal: interrupt" ||
			errStr == "signal: terminated" ||
			errStr == "signal: killed"
		if p.state == ProcessStateStopping && isExpectedSignal {
			p.state = ProcessStateStopped
		} else if p.state != ProcessStateStopping {
			p.state = ProcessStateError
			p.logger.Error("Process exited with error", "error", err)
			select {
			case p.errorCh <- fmt.Errorf("process exited: %w", err):
			default:
			}
		} else {
			p.state = ProcessStateStopped
		}
	} else {
		if p.state != ProcessStateStopping {
			p.logger.Info("Process exited normally")
		}
		p.state = ProcessStateStopped
	}

	// Signal that monitor is done - use select to avoid panic if already closed
	select {
	case <-p.doneCh:
		// Already closed
	default:
		close(p.doneCh)
	}
}
