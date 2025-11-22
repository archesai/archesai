package dev

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"os/exec"
	"sync"
	"time"
)

// Manager orchestrates multiple processes for development
type Manager struct {
	processes map[string]*Process
	watchers  map[string]*Watcher
	logger    *slog.Logger
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewManager creates a new process manager
func NewManager(logger *slog.Logger) *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	return &Manager{
		processes: make(map[string]*Process),
		watchers:  make(map[string]*Watcher),
		logger:    logger,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// AddProcess adds a process to the manager
func (m *Manager) AddProcess(config ProcessConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.processes[config.Name]; exists {
		return fmt.Errorf("process %s already exists", config.Name)
	}

	process := NewProcess(config, m.logger)
	m.processes[config.Name] = process

	// Setup file watcher for hot reload if enabled
	if config.HotReload {
		watcherConfig := WatcherConfig{
			Paths:      config.WatchPaths,
			Extensions: config.WatchExts,
			Ignore:     []string{".git", "node_modules", ".air.toml", "bin", "tmp"},
			Debounce:   500 * time.Millisecond,
		}

		callback := func() {
			m.logger.Info("Rebuilding and restarting process", "name", config.Name)

			// Run build command if specified
			if config.BuildCmd != "" {
				if err := m.buildProcess(config); err != nil {
					m.logger.Error("Build failed", "name", config.Name, "error", err)
					return
				}
			}

			// Restart the process
			if err := process.Restart(); err != nil {
				m.logger.Error("Failed to restart process", "name", config.Name, "error", err)
			}
		}

		watcher, err := NewWatcher(watcherConfig, m.logger, callback)
		if err != nil {
			m.logger.Warn("Failed to create file watcher", "name", config.Name, "error", err)
		} else {
			m.watchers[config.Name] = watcher
		}
	}

	return nil
}

// StartAll starts all processes
func (m *Manager) StartAll() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var wg sync.WaitGroup
	errCh := make(chan error, len(m.processes))

	// Build processes that need building first
	for name, process := range m.processes {
		if process.config.HotReload && process.config.BuildCmd != "" {
			m.logger.Info("Building process", "name", name)
			if err := m.buildProcess(process.config); err != nil {
				return fmt.Errorf("failed to build %s: %w", name, err)
			}
		}
	}

	// Start all processes
	for name, process := range m.processes {
		wg.Add(1)
		go func(n string, p *Process) {
			defer wg.Done()
			if err := p.Start(); err != nil {
				errCh <- fmt.Errorf("failed to start %s: %w", n, err)
			}
		}(name, process)
	}

	wg.Wait()
	close(errCh)

	// Check for errors
	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to start processes: %v", errs)
	}

	// Start file watchers
	for name, watcher := range m.watchers {
		m.logger.Info("Starting file watcher", "name", name)
		watcher.Start()
	}

	m.logger.Info("All processes started successfully")
	return nil
}

// StopAll stops all processes gracefully
func (m *Manager) StopAll() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.logger.Info("Stopping all processes")

	var wg sync.WaitGroup
	errCh := make(chan error, len(m.processes))

	for name, process := range m.processes {
		if process.GetState() != ProcessStateRunning {
			continue
		}

		wg.Add(1)
		go func(n string, p *Process) {
			defer wg.Done()
			if err := p.Stop(); err != nil {
				errCh <- fmt.Errorf("failed to stop %s: %w", n, err)
			}
		}(name, process)
	}

	wg.Wait()
	close(errCh)

	// Check for errors
	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to stop processes: %v", errs)
	}

	m.logger.Info("All processes stopped")
	return nil
}

// RestartProcess restarts a specific process
func (m *Manager) RestartProcess(name string) error {
	m.mu.RLock()
	process, exists := m.processes[name]
	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("process %s not found", name)
	}

	return process.Restart()
}

// GetProcess gets a specific process
func (m *Manager) GetProcess(name string) (*Process, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	process, exists := m.processes[name]
	if !exists {
		return nil, fmt.Errorf("process %s not found", name)
	}

	return process, nil
}

// GetProcesses returns all processes
func (m *Manager) GetProcesses() map[string]*Process {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy to avoid race conditions
	processes := make(map[string]*Process)
	maps.Copy(processes, m.processes)
	return processes
}

// Shutdown gracefully shuts down the manager and all processes
func (m *Manager) Shutdown() error {
	m.logger.Info("Shutting down process manager")

	// Stop all watchers first
	for name, watcher := range m.watchers {
		m.logger.Info("Stopping file watcher", "name", name)
		if err := watcher.Stop(); err != nil {
			m.logger.Warn("Failed to stop watcher", "name", name, "error", err)
		}
	}

	// Cancel context to stop monitoring
	m.cancel()

	// Stop all processes
	return m.StopAll()
}

// buildProcess runs the build command for a process
func (m *Manager) buildProcess(config ProcessConfig) error {
	if config.BuildCmd == "" {
		return nil
	}

	cmd := exec.Command(config.BuildCmd, config.BuildArgs...)
	cmd.Dir = config.Dir
	cmd.Env = append(os.Environ(), config.Env...)

	// Capture output for logging
	output, err := cmd.CombinedOutput()
	if err != nil {
		m.logger.Error("Build failed", "name", config.Name, "output", string(output))
		return fmt.Errorf("build failed: %w\n%s", err, output)
	}

	m.logger.Info("Build successful", "name", config.Name)
	return nil
}
