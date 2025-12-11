// Package generators provides code generation from OpenAPI specifications.
package generators

// GenerateFunc is the signature for all generators.
// Each generator receives the context and handles its own iteration logic.
type GenerateFunc func(ctx *GeneratorContext) error

// Generator defines a code generator.
type Generator struct {
	// Name is the unique identifier for this generator.
	Name string

	// Priority determines execution order. Lower values run first.
	// Generators with the same priority run in parallel.
	Priority int

	// Generate is the function that performs the generation.
	Generate GenerateFunc
}
