// Package domain contains the workflows domain business logic and entities
package domain

import (
	"errors"
	"time"

	"github.com/archesai/archesai/internal/workflows/generated/api"
)

// Domain errors
var (
	ErrPipelineNotFound  = errors.New("pipeline not found")
	ErrPipelineExists    = errors.New("pipeline already exists")
	ErrRunNotFound       = errors.New("run not found")
	ErrRunAlreadyStarted = errors.New("run already started")
	ErrRunNotStarted     = errors.New("run not started")
	ErrToolNotFound      = errors.New("tool not found")
	ErrToolExists        = errors.New("tool already exists")
	ErrInvalidStatus     = errors.New("invalid status")
	ErrInvalidProgress   = errors.New("invalid progress value")
)

// Pipeline extends the generated API PipelineEntity with domain-specific fields
type Pipeline struct {
	api.PipelineEntity
	// Add any domain-specific fields that aren't in the API
}

// Run extends the generated API RunEntity with domain-specific fields
type Run struct {
	api.RunEntity
	// Add any domain-specific fields that aren't in the API
}

// Tool extends the generated API ToolEntity with domain-specific fields
type Tool struct {
	api.ToolEntity
	// Add any domain-specific fields that aren't in the API
}

// CreatePipelineRequest represents a request to create a pipeline
type CreatePipelineRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"required,min=1,max=500"`
}

// UpdatePipelineRequest represents a request to update a pipeline
type UpdatePipelineRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,min=1,max=500"`
}

// CreateRunRequest represents a request to create a run
type CreateRunRequest struct {
	PipelineID *string                `json:"pipeline_id,omitempty"`
	ToolID     string                 `json:"tool_id" validate:"required"`
	Input      map[string]interface{} `json:"input,omitempty"`
}

// UpdateRunRequest represents a request to update a run
type UpdateRunRequest struct {
	Status   *api.RunEntityStatus   `json:"status,omitempty"`
	Progress *float32               `json:"progress,omitempty" validate:"omitempty,min=0,max=100"`
	Output   map[string]interface{} `json:"output,omitempty"`
}

// CreateToolRequest represents a request to create a tool
type CreateToolRequest struct {
	Name        string                 `json:"name" validate:"required,min=1,max=100"`
	Description string                 `json:"description" validate:"required,min=1,max=500"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// UpdateToolRequest represents a request to update a tool
type UpdateToolRequest struct {
	Name        *string                `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string                `json:"description,omitempty" validate:"omitempty,min=1,max=500"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// NewPipeline creates a new pipeline from the API entity
func NewPipeline(entity api.PipelineEntity) *Pipeline {
	return &Pipeline{PipelineEntity: entity}
}

// NewRun creates a new run from the API entity
func NewRun(entity api.RunEntity) *Run {
	return &Run{RunEntity: entity}
}

// NewTool creates a new tool from the API entity
func NewTool(entity api.ToolEntity) *Tool {
	return &Tool{ToolEntity: entity}
}

// IsCompleted checks if a run is in a completed state
func (r *Run) IsCompleted() bool {
	return r.Status == api.COMPLETED || r.Status == api.FAILED
}

// IsRunning checks if a run is currently executing
func (r *Run) IsRunning() bool {
	return r.Status == api.PROCESSING
}

// CanStart checks if a run can be started
func (r *Run) CanStart() bool {
	return r.Status == api.QUEUED
}

// CanCancel checks if a run can be cancelled
func (r *Run) CanCancel() bool {
	return r.Status == api.PROCESSING || r.Status == api.QUEUED
}

// UpdateProgress updates the run's progress and validates the value
func (r *Run) UpdateProgress(progress float32) error {
	if progress < 0 || progress > 100 {
		return ErrInvalidProgress
	}
	r.Progress = progress
	r.UpdatedAt = time.Now()
	return nil
}

// GetDuration calculates the duration of a run if it's completed
func (r *Run) GetDuration() (time.Duration, error) {
	if !r.IsCompleted() || r.StartedAt.IsZero() || r.CompletedAt.IsZero() {
		return 0, errors.New("run not completed or missing timestamps")
	}

	return r.CompletedAt.Sub(r.StartedAt), nil
}

// EstimatedCompletion estimates when a run will complete based on current progress
func (r *Run) EstimatedCompletion() (*time.Time, error) {
	if r.Progress <= 0 || r.StartedAt.IsZero() || r.IsCompleted() {
		return nil, errors.New("cannot estimate completion")
	}

	elapsed := time.Since(r.StartedAt)
	totalEstimated := time.Duration(float64(elapsed) / float64(r.Progress) * 100)
	estimated := r.StartedAt.Add(totalEstimated)

	return &estimated, nil
}
