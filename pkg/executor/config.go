package executor

import "time"

// Config holds configuration for an executor
type Config struct {
	Timeout time.Duration // Execution timeout

	// Schema validation
	SchemaIn  []byte // JSON Schema for input validation
	SchemaOut []byte // JSON Schema for output validation
}
