package codegen

// Options configures the code generation process.
type Options struct {
	OutputPath string
	SpecPath   string
	Lint       bool
	Only       string
	TUI        bool
}
