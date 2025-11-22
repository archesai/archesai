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
	"time"
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

// Process represents a managed process
type Process struct {
	config   ProcessConfig
	cmd      *exec.Cmd
	state    ProcessState
	logger   *slog.Logger
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
	outputCh chan string
	errorCh  chan error
	doneCh   chan struct{}
}

// NewProcess creates a new process instance
func NewProcess(config ProcessConfig, logger *slog.Logger) *Process {
	ctx, cancel := context.WithCancel(context.Background())
	return &Process{
		config:   config,
		state:    ProcessStateStopped,
		logger:   logger.With(slog.String("service", config.Name)),
		ctx:      ctx,
		cancel:   cancel,
		outputCh: make(chan string, 100),
		errorCh:  make(chan error, 10),
		doneCh:   make(chan struct{}),
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

	// Create the command
	p.cmd = exec.CommandContext(p.ctx, p.config.Command, p.config.Args...)

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
	defer p.mu.Unlock()

	if p.state != ProcessStateRunning {
		return fmt.Errorf("process %s is not running", p.config.Name)
	}

	p.state = ProcessStateStopping
	p.logger.Info("Stopping process")

	// Cancel the context to signal shutdown
	p.cancel()

	// Give the process time to shutdown gracefully
	done := make(chan error, 1)
	go func() {
		done <- p.cmd.Wait()
	}()

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
	case err := <-done:
		if err != nil && err.Error() != "signal: terminated" {
			p.logger.Debug("Process exited with error", "error", err)
		}
	}

	p.state = ProcessStateStopped
	close(p.doneCh)
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
func (p *Process) readOutput(pipe io.Reader, source string) {
	scanner := bufio.NewScanner(pipe)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024) // 1MB max line size

	for scanner.Scan() {
		line := scanner.Text()

		// Log to stdout/stderr normally
		p.logger.Info(line, "source", source)

		// Also send to output channel for buffer consumption
		select {
		case p.outputCh <- line:
		case <-p.ctx.Done():
			return
		default:
			// Drop message if channel is full
		}
	}

	if err := scanner.Err(); err != nil {
		p.errorCh <- fmt.Errorf("error reading %s: %w", source, err)
	}
}

// monitor monitors the process and updates state
func (p *Process) monitor() {
	err := p.cmd.Wait()

	p.mu.Lock()
	defer p.mu.Unlock()

	if err != nil {
		if p.state != ProcessStateStopping {
			p.state = ProcessStateError
			p.logger.Error("Process exited with error", "error", err)
			p.errorCh <- fmt.Errorf("process exited: %w", err)
		}
	} else {
		if p.state != ProcessStateStopping {
			p.state = ProcessStateStopped
			p.logger.Info("Process exited normally")
		}
	}

	if p.state == ProcessStateStopping {
		p.state = ProcessStateStopped
	}
}
