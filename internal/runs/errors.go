package runs

import "errors"

// Domain errors
var (
	// ErrNotFound is returned when a resource is not found
	ErrNotFound = errors.New("not found")

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
