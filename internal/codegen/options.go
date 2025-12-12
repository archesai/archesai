package codegen

import (
	"github.com/archesai/archesai/pkg/config"
)

// Options configures the code generation process.
type Options struct {
	OutputPath string   `yaml:"output"`
	SpecPath   string   `yaml:"spec"`
	Lint       bool     `yaml:"lint"`
	Only       []string `yaml:"only"`
	Includes   []string `yaml:"includes"`
	ConfigPath string   `yaml:"-"` // Path to config file
	WorkDir    string   `yaml:"-"` // Working directory (config file's dir)
}

// generationConfig wraps Options for parsing from arches.yaml.
type generationConfig struct {
	Generation Options `yaml:"generation"`
}

// LoadOptionsFromConfig loads generation options from arches.yaml.
// If configFile is provided, it loads from that specific file.
// Otherwise, it searches for config files in standard locations.
func LoadOptionsFromConfig(configFile string) Options {
	parser := config.NewParser[generationConfig]()
	cfg, err := parser.LoadFrom(configFile)
	if err != nil {
		return Options{}
	}
	opts := cfg.Config.Generation
	opts.ConfigPath = cfg.ConfigPath
	opts.WorkDir = cfg.WorkDir()
	return opts
}
