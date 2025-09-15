package pipelines

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/archesai/archesai/internal/runs"
	"github.com/archesai/archesai/internal/tools"
	"github.com/google/uuid"
)

// QueueService defines the interface for workflow queueing
type QueueService interface {
	EnqueueRun(ctx context.Context, runID uuid.UUID) error
	DequeueRun(ctx context.Context, status string) (*uuid.UUID, error)
}

// ToolRunner defines the interface for executing tools
type ToolRunner interface {
	// Run executes a tool with the given configuration and input
	Run(ctx context.Context, tool *tools.Tool, config map[string]interface{}, input interface{}) (interface{}, error)
}

// ToolRunnerFunc is an adapter to allow functions to be used as ToolRunner
type ToolRunnerFunc func(ctx context.Context, tool *tools.Tool, config map[string]interface{}, input interface{}) (interface{}, error)

// Run implements ToolRunner
func (f ToolRunnerFunc) Run(ctx context.Context, tool *tools.Tool, config map[string]interface{}, input interface{}) (interface{}, error) {
	return f(ctx, tool, config, input)
}

// WorkflowExecutor handles the execution of workflow pipelines
type WorkflowExecutor struct {
	pipelineRepo    Repository
	runRepo         runs.Repository
	toolRepo        tools.Repository
	pipelineManager *PipelineManager
	runner          ToolRunner
	queue           QueueService
	logger          *slog.Logger
	maxParallel     int
	mu              sync.RWMutex
	executions      map[uuid.UUID]*ExecutionContext
}

// ExecutionContext tracks the state of a running pipeline
type ExecutionContext struct {
	RunID      uuid.UUID
	PipelineID uuid.UUID
	DAG        *DAG
	Executor   *DAGExecutor
	StartTime  time.Time
	EndTime    *time.Time
	Status     runs.RunStatus
	Error      error
	Results    map[uuid.UUID]interface{}
	mu         sync.RWMutex
}

// NewWorkflowExecutor creates a new workflow executor
func NewWorkflowExecutor(pipelineRepo Repository, runRepo runs.Repository, toolRepo tools.Repository, pipelineManager *PipelineManager, runner ToolRunner, queue QueueService, logger *slog.Logger, maxParallel int) *WorkflowExecutor {
	if maxParallel <= 0 {
		maxParallel = 4 // Default parallelism
	}
	return &WorkflowExecutor{
		pipelineRepo:    pipelineRepo,
		runRepo:         runRepo,
		toolRepo:        toolRepo,
		pipelineManager: pipelineManager,
		runner:          runner,
		queue:           queue,
		logger:          logger,
		maxParallel:     maxParallel,
		executions:      make(map[uuid.UUID]*ExecutionContext),
	}
}

// ExecutePipeline starts the execution of a pipeline
func (we *WorkflowExecutor) ExecutePipeline(ctx context.Context, pipelineID uuid.UUID, input map[string]interface{}) (*runs.Run, error) {
	// Get pipeline
	pipeline, err := we.pipelineRepo.Get(ctx, pipelineID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline: %w", err)
	}

	// Create run record
	run := &runs.Run{
		Id:         uuid.New(),
		PipelineId: pipeline.Id,
		Status:     runs.QUEUED,
		Progress:   0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Store input as metadata in error field for now
	// TODO: Add proper metadata field or use a separate table
	if input != nil {
		inputJSON, err := json.Marshal(input)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal input: %w", err)
		}
		// Temporarily store in error field - this should be refactored
		run.Error = string(inputJSON)
	}

	// Save run to database
	createdRun, err := we.runRepo.Create(ctx, run)
	if err != nil {
		return nil, fmt.Errorf("failed to create run: %w", err)
	}

	// Queue the execution
	err = we.queue.EnqueueRun(ctx, createdRun.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to queue run: %w", err)
	}

	// Start execution in background
	go we.executeRun(context.Background(), createdRun.Id)

	return createdRun, nil
}

// executeRun performs the actual pipeline execution
func (we *WorkflowExecutor) executeRun(ctx context.Context, runID uuid.UUID) {
	we.logger.Info("Starting pipeline execution", "runId", runID)

	// Get run details
	run, err := we.runRepo.Get(ctx, runID)
	if err != nil {
		we.logger.Error("Failed to get run", "error", err, "runId", runID)
		return
	}

	// Update run status to processing
	run.Status = runs.PROCESSING
	run.StartedAt = time.Now()
	_, err = we.runRepo.Update(ctx, run.Id, run)
	if err != nil {
		we.logger.Error("Failed to update run status", "error", err, "runId", runID)
		return
	}

	// Build DAG from pipeline steps
	dag, err := we.buildDAG(ctx, run.PipelineId)
	if err != nil {
		we.handleRunFailure(ctx, run, err)
		return
	}

	// Create execution context
	execCtx := &ExecutionContext{
		RunID:      runID,
		PipelineID: run.PipelineId,
		DAG:        dag,
		StartTime:  time.Now(),
		Status:     runs.PROCESSING,
		Results:    make(map[uuid.UUID]interface{}),
	}

	// Store execution context
	we.mu.Lock()
	we.executions[runID] = execCtx
	we.mu.Unlock()

	// Create DAG executor with tool runner adapter
	toolExecutor := &dagToolExecutor{
		executor: we,
		runID:    runID,
	}
	execCtx.Executor = NewDAGExecutor(dag, toolExecutor, we.maxParallel)

	// Execute the DAG
	err = execCtx.Executor.Execute(ctx)
	if err != nil {
		we.handleRunFailure(ctx, run, err)
		return
	}

	// Mark run as completed
	we.handleRunSuccess(ctx, run)
}

// buildDAG creates a DAG from pipeline steps
func (we *WorkflowExecutor) buildDAG(ctx context.Context, pipelineID uuid.UUID) (*DAG, error) {
	we.logger.Info("Building DAG for pipeline", "pipelineId", pipelineID)

	// Use the pipeline manager to get DAG structure
	steps, dependencies, err := we.pipelineManager.GetPipelineDAG(ctx, pipelineID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline DAG: %w", err)
	}

	return NewDAG(steps, dependencies)
}

// handleRunFailure updates the run status to failed
func (we *WorkflowExecutor) handleRunFailure(ctx context.Context, run *runs.Run, err error) {
	we.logger.Error("Pipeline execution failed", "error", err, "runId", run.Id)

	run.Status = runs.FAILED
	run.CompletedAt = time.Now()

	// Store error message
	run.Error = err.Error()

	_, updateErr := we.runRepo.Update(ctx, run.Id, run)
	if updateErr != nil {
		we.logger.Error("Failed to update failed run", "error", updateErr, "runId", run.Id)
	}
}

// handleRunSuccess updates the run status to completed
func (we *WorkflowExecutor) handleRunSuccess(ctx context.Context, run *runs.Run) {
	we.logger.Info("Pipeline execution completed", "runId", run.Id)

	run.Status = runs.COMPLETED
	run.CompletedAt = time.Now()
	run.Progress = 100

	_, err := we.runRepo.Update(ctx, run.Id, run)
	if err != nil {
		we.logger.Error("Failed to update completed run", "error", err, "runId", run.Id)
	}
}

// GetExecutionStatus returns the current status of a run
func (we *WorkflowExecutor) GetExecutionStatus(runID uuid.UUID) (*ExecutionContext, bool) {
	we.mu.RLock()
	defer we.mu.RUnlock()

	exec, exists := we.executions[runID]
	return exec, exists
}

// dagToolExecutor adapts WorkflowExecutor to work with DAGExecutor
type dagToolExecutor struct {
	executor *WorkflowExecutor
	runID    uuid.UUID
}

// Execute implements ToolExecutor interface for DAGExecutor
func (dte *dagToolExecutor) Execute(ctx context.Context, tool *tools.Tool, input interface{}) (interface{}, error) {
	// Get the step configuration from the execution context
	dte.executor.mu.RLock()
	execCtx, exists := dte.executor.executions[dte.runID]
	dte.executor.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("execution context not found")
	}

	// Execute the tool using the runner
	config := make(map[string]interface{}) // Get from step config
	result, err := dte.executor.runner.Run(ctx, tool, config, input)
	if err != nil {
		return nil, fmt.Errorf("tool execution failed: %w", err)
	}

	// Store result in execution context
	execCtx.mu.Lock()
	// Store by tool ID for now - you might want to use step ID instead
	execCtx.Results[tool.Id] = result
	execCtx.mu.Unlock()

	return result, nil
}

// ContainerToolRunner executes tools in containers
type ContainerToolRunner struct {
	logger *slog.Logger
}

// NewContainerToolRunner creates a new container-based tool runner
func NewContainerToolRunner(logger *slog.Logger) *ContainerToolRunner {
	return &ContainerToolRunner{
		logger: logger,
	}
}

// Run executes a tool in a container
func (ctr *ContainerToolRunner) Run(_ context.Context, tool *tools.Tool, config map[string]interface{}, input interface{}) (interface{}, error) {
	_ = input // will be used in actual implementation
	ctr.logger.Info("Executing tool in container",
		"toolId", tool.Id,
		"toolName", tool.Name,
		"config", config,
	)

	// TODO: Implement actual container execution
	// This would typically:
	// 1. Pull the container image specified in tool configuration
	// 2. Mount input data as volume or pass as environment
	// 3. Run the container with resource limits
	// 4. Capture output and logs
	// 5. Clean up resources

	// For now, return a mock result
	return map[string]interface{}{
		"toolId":     tool.Id,
		"toolName":   tool.Name,
		"executedAt": time.Now(),
		"status":     "success",
		"output":     "Mock execution result",
	}, nil
}

// HTTPToolRunner executes tools via HTTP API calls
type HTTPToolRunner struct {
	logger *slog.Logger
}

// NewHTTPToolRunner creates a new HTTP-based tool runner
func NewHTTPToolRunner(logger *slog.Logger) *HTTPToolRunner {
	return &HTTPToolRunner{
		logger: logger,
	}
}

// Run executes a tool via HTTP API
func (htr *HTTPToolRunner) Run(_ context.Context, tool *tools.Tool, config map[string]interface{}, input interface{}) (interface{}, error) {
	_ = input // will be used in actual implementation
	htr.logger.Info("Executing tool via HTTP",
		"toolId", tool.Id,
		"toolName", tool.Name,
		"config", config,
	)

	// TODO: Implement actual HTTP execution
	// This would typically:
	// 1. Build HTTP request from tool configuration
	// 2. Add authentication if required
	// 3. Send request with input data
	// 4. Handle retries and timeouts
	// 5. Parse response

	return map[string]interface{}{
		"toolId":     tool.Id,
		"toolName":   tool.Name,
		"executedAt": time.Now(),
		"status":     "success",
		"response":   "Mock HTTP response",
	}, nil
}

// CompositeToolRunner routes to different runners based on tool type
type CompositeToolRunner struct {
	runners map[string]ToolRunner
	logger  *slog.Logger
}

// NewCompositeToolRunner creates a runner that delegates to specific runners
func NewCompositeToolRunner(logger *slog.Logger) *CompositeToolRunner {
	return &CompositeToolRunner{
		runners: make(map[string]ToolRunner),
		logger:  logger,
	}
}

// RegisterRunner registers a runner for a specific tool type
func (ctr *CompositeToolRunner) RegisterRunner(toolType string, runner ToolRunner) {
	ctr.runners[toolType] = runner
}

// Run delegates to the appropriate runner based on tool type
func (ctr *CompositeToolRunner) Run(ctx context.Context, tool *tools.Tool, config map[string]interface{}, input interface{}) (interface{}, error) {
	// Determine tool type from configuration or metadata
	toolType := "container" // Default type
	if typeVal, ok := config["type"].(string); ok {
		toolType = typeVal
	}

	runner, exists := ctr.runners[toolType]
	if !exists {
		return nil, fmt.Errorf("no runner registered for tool type: %s", toolType)
	}

	return runner.Run(ctx, tool, config, input)
}
