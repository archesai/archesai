package codegen

// GenerateOptions configures the code generation process.
type Options struct {
	OutputPath string
	SpecPath   string
	OrvalFix   bool
	DryRun     bool
	Lint       bool
	Only       string
	TUI        bool
}
