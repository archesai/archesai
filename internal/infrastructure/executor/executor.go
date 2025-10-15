// Package executor provides a generic interface for executing code in isolated containers
package executor

import (
	"context"
	"time"
)

// Executor is a generic interface for executing functions that transform input A to output B
type Executor[A any, B any] interface {
	// Execute runs the executor with the given input and returns the output or an error
	Execute(ctx context.Context, input A) (B, error)
}

// Config holds configuration for an executor
type Config struct {
	Timeout time.Duration // Execution timeout

	// Schema validation
	SchemaIn  []byte // JSON Schema for input validation
	SchemaOut []byte // JSON Schema for output validation
}
