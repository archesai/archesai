// Package workflows provides workflow management functionality including
// pipeline definitions, run executions, and tool management.
package workflows

//go:generate go tool oapi-codegen --config=../../types.codegen.yaml --package workflows --include-tags Workflows,Pipelines,Runs,Tools ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../server.codegen.yaml --package workflows --include-tags Workflows,Pipelines,Runs,Tools ../../api/openapi.bundled.yaml

import "errors"

// Common errors
var (
	ErrNotFound = errors.New("not found")
)

// Domain types
type (

	// CreatePipelineRequest represents a request to create a pipeline
	CreatePipelineRequest = CreatePipelineJSONBody

	// UpdatePipelineRequest represents a request to update a pipeline
	UpdatePipelineRequest = UpdatePipelineJSONBody

	// CreateRunRequest represents a request to create a run
	CreateRunRequest = CreateRunJSONBody

	// UpdateRunRequest represents a request to update a run
	UpdateRunRequest = UpdateRunJSONBody

	// CreateToolRequest represents a request to create a tool
	CreateToolRequest = CreateToolJSONBody

	// UpdateToolRequest represents a request to update a tool
	UpdateToolRequest = UpdateToolJSONBody
)

// Domain errors
var (
	// ErrPipelineNotFound is returned when a pipeline is not found
	ErrPipelineNotFound = errors.New("pipeline not found")

	// ErrRunNotFound is returned when a run is not found
	ErrRunNotFound = errors.New("run not found")

	// ErrToolNotFound is returned when a tool is not found
	ErrToolNotFound = errors.New("tool not found")

	// ErrToolExists is returned when a tool already exists
	ErrToolExists = errors.New("tool already exists")

	// ErrInvalidTransition is returned when an invalid state transition is attempted
	ErrInvalidTransition = errors.New("invalid state transition")
)

// CanStart checks if the run can be started
func (r *Run) CanStart() bool {
	return r.Status == QUEUED
}

// IsRunning checks if the run is currently running
func (r *Run) IsRunning() bool {
	return r.Status == PROCESSING
}

// CanCancel checks if the run can be cancelled
func (r *Run) CanCancel() bool {
	return r.Status == PROCESSING || r.Status == QUEUED
}

// UpdateProgress updates the run's progress
func (r *Run) UpdateProgress(progress float32) {
	r.Progress = progress
}

// Domain constants
const (
	// DefaultTimeout is the default timeout for workflow runs
	DefaultTimeout = 3600 // 1 hour in seconds

	// MaxPipelinesPerOrganization defines the maximum number of pipelines per organization
	MaxPipelinesPerOrganization = 1000

	// MaxRunsToKeep defines how many completed runs to keep per pipeline
	MaxRunsToKeep = 100
)
