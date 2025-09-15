// Package pipelines provides pipeline management functionality
package pipelines

// Domain types
type (
	// CreatePipelineRequest represents a request to create a pipeline
	CreatePipelineRequest = CreatePipelineJSONBody

	// UpdatePipelineRequest represents a request to update a pipeline
	UpdatePipelineRequest = UpdatePipelineJSONBody
)

// Domain constants
const (
	// DefaultTimeout is the default timeout for pipelines
	DefaultTimeout = 3600 // 1 hour in seconds

	// MaxPipelinesPerOrganization defines the maximum number of pipelines per organization
	MaxPipelinesPerOrganization = 1000
)
