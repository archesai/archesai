package capabilities

import (
	"context"
	"log/slog"
	"sync"
)

// Checker defines the interface for checking a single capability.
type Checker interface {
	Check(ctx context.Context) Capability
}

// Detector orchestrates capability detection across multiple checkers.
type Detector struct {
	logger   *slog.Logger
	checkers []Checker
}

// NewDetector creates a new Detector with the given logger.
func NewDetector(logger *slog.Logger) *Detector {
	if logger == nil {
		logger = slog.Default()
	}
	return &Detector{
		logger:   logger,
		checkers: []Checker{},
	}
}

// WithChecker adds a checker to the detector.
func (d *Detector) WithChecker(c Checker) *Detector {
	d.checkers = append(d.checkers, c)
	return d
}

// WithCheckers adds multiple checkers to the detector.
func (d *Detector) WithCheckers(checkers ...Checker) *Detector {
	d.checkers = append(d.checkers, checkers...)
	return d
}

// Check runs all registered checkers and returns the aggregated result.
// Checkers are run concurrently for better performance.
func (d *Detector) Check(ctx context.Context) CheckResult {
	result := CheckResult{
		Capabilities: make([]Capability, len(d.checkers)),
		AllRequired:  true,
		AllFound:     true,
	}

	if len(d.checkers) == 0 {
		return result
	}

	var wg sync.WaitGroup
	wg.Add(len(d.checkers))

	for i, checker := range d.checkers {
		go func(idx int, c Checker) {
			defer wg.Done()
			capability := c.Check(ctx)
			result.Capabilities[idx] = capability

			d.logger.Debug("capability check complete",
				slog.String("name", capability.Name),
				slog.Bool("found", capability.Found),
				slog.String("version", capability.Version),
			)
		}(i, checker)
	}

	wg.Wait()

	// Calculate aggregate status
	for _, cap := range result.Capabilities {
		if !cap.Found {
			result.AllFound = false
			if cap.Required {
				result.AllRequired = false
			}
		}
	}

	return result
}

// DefaultDetector returns a detector configured with all standard checkers.
func DefaultDetector(logger *slog.Logger) *Detector {
	return NewDetector(logger).WithCheckers(
		// Required capabilities
		NewGoChecker(true),
		NewDockerChecker(true),

		// Optional capabilities
		NewNodeChecker(true),
		NewPnpmChecker(true),
		NewPythonChecker(true),
		NewGitChecker(false),
		NewKubectlChecker(false),
		NewHelmChecker(false),
	)
}
