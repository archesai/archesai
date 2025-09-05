// Package workflows provides workflow management functionality including
// pipeline definitions, run executions, and tool management.
package workflows

// Domain constants
const (
	// DefaultTimeout is the default timeout for workflow runs
	DefaultTimeout = 3600 // 1 hour in seconds

	// MaxPipelinesPerOrganization defines the maximum number of pipelines per organization
	MaxPipelinesPerOrganization = 1000

	// MaxRunsToKeep defines how many completed runs to keep per pipeline
	MaxRunsToKeep = 100
)
