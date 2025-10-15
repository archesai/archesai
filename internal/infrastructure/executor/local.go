package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// LocalExecutor runs code locally without container isolation
type LocalExecutor[A any, B any] struct {
	config      LocalConfig
	validator   *SchemaValidator
	executeFunc func(context.Context, A) (B, error)
}

// LocalConfig holds configuration for a local executor
type LocalConfig struct {
	Config
}

// NewLocalExecutor creates a new local executor with the given execution function
func NewLocalExecutor[A any, B any](
	config LocalConfig,
	executeFunc func(context.Context, A) (B, error),
) (*LocalExecutor[A, B], error) {
	if executeFunc == nil {
		return nil, fmt.Errorf("execute function is required")
	}

	// Set defaults
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	// Create schema validator if schemas are provided
	var validator *SchemaValidator
	if len(config.SchemaIn) > 0 && len(config.SchemaOut) > 0 {
		var err error
		validator, err = NewSchemaValidator(config.SchemaIn, config.SchemaOut)
		if err != nil {
			return nil, fmt.Errorf("create schema validator: %w", err)
		}
	}

	return &LocalExecutor[A, B]{
		config:      config,
		executeFunc: executeFunc,
		validator:   validator,
	}, nil
}

// Execute runs the local executor with the given input and returns the output
func (e *LocalExecutor[A, B]) Execute(ctx context.Context, input A) (B, error) {
	var zero B

	// Marshal input to JSON
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return zero, fmt.Errorf("marshal input: %w", err)
	}

	// Validate input against schema if validator is configured
	if e.validator != nil {
		var inputAny any
		if err := json.Unmarshal(inputBytes, &inputAny); err != nil {
			return zero, fmt.Errorf("unmarshal input for validation: %w", err)
		}
		if err := e.validator.ValidateInput(inputAny); err != nil {
			return zero, fmt.Errorf("input validation failed: %w", err)
		}
	}

	// Create context with timeout
	execCtx, cancel := context.WithTimeout(ctx, e.config.Timeout)
	defer cancel()

	// Execute the function
	output, err := e.executeFunc(execCtx, input)
	if err != nil {
		// Check if it was a timeout
		if execCtx.Err() == context.DeadlineExceeded {
			return zero, fmt.Errorf("execution timed out after %s", e.config.Timeout)
		}
		return zero, fmt.Errorf("execution failed: %w", err)
	}

	// Validate output against schema if validator is configured
	if e.validator != nil {
		outputBytes, err := json.Marshal(output)
		if err != nil {
			return zero, fmt.Errorf("marshal output: %w", err)
		}

		var outputAny any
		if err := json.Unmarshal(outputBytes, &outputAny); err != nil {
			return zero, fmt.Errorf("unmarshal output for validation: %w", err)
		}

		if err := e.validator.ValidateOutput(outputAny); err != nil {
			return zero, fmt.Errorf("output validation failed: %w", err)
		}
	}

	return output, nil
}
