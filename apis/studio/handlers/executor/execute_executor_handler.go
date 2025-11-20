// Package executor provides handlers for code execution operations.
package executor

import (
	"context"
	"fmt"

	commands "github.com/archesai/archesai/apis/studio/generated/application/commands/executor"
	"github.com/archesai/archesai/apis/studio/generated/core/repositories"
	"github.com/archesai/archesai/pkg/events"
	"github.com/archesai/archesai/pkg/executor"
)

// ExecutionResult represents the result of executing an executor
type ExecutionResult struct {
	Output          map[string]any `json:"output"`
	ExecutionTimeMs int            `json:"executionTimeMs,omitempty"`
	Logs            string         `json:"logs,omitempty"`
}

// ExecuteExecutorCommandHandler handles the execute executor command.
type ExecuteExecutorCommandHandler struct {
	repo            repositories.ExecutorRepository
	executorService executor.ExecutorService[map[string]any, map[string]any]
	publisher       events.Publisher
}

// NewExecuteExecutorCommandHandler creates a new execute executor command handler.
func NewExecuteExecutorCommandHandler(
	repo repositories.ExecutorRepository,
	executorService executor.ExecutorService[map[string]any, map[string]any],
	publisher events.Publisher,
) *ExecuteExecutorCommandHandler {
	return &ExecuteExecutorCommandHandler{
		repo:            repo,
		executorService: executorService,
		publisher:       publisher,
	}
}

// Handle executes the executor command.
func (h *ExecuteExecutorCommandHandler) Handle(
	ctx context.Context,
	cmd *commands.ExecuteExecutorCommand,
) (*ExecutionResult, error) {
	// Execute using the executor service
	result, err := h.executorService.Execute(ctx, cmd.ID.String(), cmd.Input)
	if err != nil {
		return nil, fmt.Errorf("failed to execute executor: %w", err)
	}

	return &ExecutionResult{
		Output:          result.Output,
		ExecutionTimeMs: int(result.ExecutionTimeMs),
		Logs:            result.Logs,
	}, nil
}
