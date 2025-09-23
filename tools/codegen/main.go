// Package main provides the entry point for the codegen tool
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/archesai/archesai/internal/codegen"
)

func main() {
	var configPath string
	var verbose bool
	var openapiPath string
	var outputPath string

	flag.StringVar(&configPath, "config", ".archesai.codegen.yaml", "Config file path")
	flag.StringVar(
		&openapiPath,
		"openapi",
		"api/openapi.bundled.yaml",
		"OpenAPI spec file path",
	)
	flag.StringVar(&outputPath, "output", "internal", "Output directory for generated code")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")
	flag.BoolVar(&verbose, "v", false, "Verbose output (shorthand)")
	flag.Parse()

	if verbose {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	} else {
		log.SetFlags(0)
	}

	cfg := codegen.GetDefaultConfig()
	cfg.Openapi = openapiPath
	cfg.Output = &outputPath

	// Always use config file
	if err := codegen.Run(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
