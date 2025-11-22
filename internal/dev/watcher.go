package dev

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// WatcherConfig defines configuration for the file watcher
type WatcherConfig struct {
	Paths      []string      // Paths to watch
	Extensions []string      // File extensions to trigger on (e.g., .go, .mod)
	Ignore     []string      // Paths to ignore
	Debounce   time.Duration // Time to wait after last change before triggering
}

// Watcher watches files and triggers callbacks on changes
type Watcher struct {
	config   WatcherConfig
	watcher  *fsnotify.Watcher
	logger   *slog.Logger
	callback func()
	ctx      context.Context
	cancel   context.CancelFunc
	mu       sync.Mutex
	timer    *time.Timer
}

// NewWatcher creates a new file watcher
func NewWatcher(config WatcherConfig, logger *slog.Logger, callback func()) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Default debounce if not set
	if config.Debounce == 0 {
		config.Debounce = 100 * time.Millisecond
	}

	w := &Watcher{
		config:   config,
		watcher:  fsWatcher,
		logger:   logger,
		callback: callback,
		ctx:      ctx,
		cancel:   cancel,
	}

	// Add paths to watch
	for _, path := range config.Paths {
		if err := w.addPath(path); err != nil {
			logger.Warn("Failed to watch path", "path", path, "error", err)
		}
	}

	return w, nil
}

// addPath recursively adds a path to watch
func (w *Watcher) addPath(path string) error {
	// Check if path should be ignored
	for _, ignore := range w.config.Ignore {
		if matched, _ := filepath.Match(ignore, path); matched {
			return nil
		}
		// Also check if path contains the ignore pattern
		if strings.Contains(path, ignore) {
			return nil
		}
	}

	// Add the path
	if err := w.watcher.Add(path); err != nil {
		return err
	}

	w.logger.Debug("Watching path", "path", path)

	// If it's a directory, recursively add subdirectories
	return filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		// Skip if should be ignored
		for _, ignore := range w.config.Ignore {
			if strings.Contains(p, ignore) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// Add directories
		if info.IsDir() && p != path {
			if err := w.watcher.Add(p); err == nil {
				w.logger.Debug("Watching subdirectory", "path", p)
			}
		}

		return nil
	})
}

// shouldTrigger checks if a file change should trigger the callback
func (w *Watcher) shouldTrigger(path string) bool {
	// Check if file has the right extension
	if len(w.config.Extensions) > 0 {
		ext := filepath.Ext(path)
		found := false
		for _, e := range w.config.Extensions {
			if ext == e {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check if path should be ignored
	for _, ignore := range w.config.Ignore {
		if strings.Contains(path, ignore) {
			return false
		}
	}

	return true
}

// Start starts watching for file changes
func (w *Watcher) Start() {
	go func() {
		for {
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					return
				}

				// Check if we should trigger on this file
				if !w.shouldTrigger(event.Name) {
					continue
				}

				// Handle the event
				switch {
				case event.Op&fsnotify.Write == fsnotify.Write,
					event.Op&fsnotify.Create == fsnotify.Create:
					w.logger.Debug("File changed", "file", event.Name, "op", event.Op.String())
					w.triggerCallback()

				case event.Op&fsnotify.Remove == fsnotify.Remove:
					// File was removed, might be part of atomic save
					w.logger.Debug("File removed", "file", event.Name)
					w.triggerCallback()
				}

			case err, ok := <-w.watcher.Errors:
				if !ok {
					return
				}
				w.logger.Error("Watcher error", "error", err)

			case <-w.ctx.Done():
				return
			}
		}
	}()
}

// triggerCallback triggers the callback with debouncing
func (w *Watcher) triggerCallback() {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Cancel existing timer if any
	if w.timer != nil {
		w.timer.Stop()
	}

	// Set new timer
	w.timer = time.AfterFunc(w.config.Debounce, func() {
		w.logger.Info("File changes detected, triggering reload")
		w.callback()
	})
}

// Stop stops the watcher
func (w *Watcher) Stop() error {
	w.cancel()
	return w.watcher.Close()
}
