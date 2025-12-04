// Package generate provides the generate command orchestration logic.
package generate

// Options configures the code generation process.
type Options struct {
	OutputPath string
	SpecPath   string
	OrvalFix   bool
	DryRun     bool
	Lint       bool
	Only       string
}
