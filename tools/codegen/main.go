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

	flag.StringVar(&configPath, "config", ".archesai.codegen.yaml", "Config file path")
	flag.StringVar(&configPath, "c", ".archesai.codegen.yaml", "Config file path (shorthand)")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")
	flag.BoolVar(&verbose, "v", false, "Verbose output (shorthand)")
	flag.Parse()

	if verbose {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	} else {
		log.SetFlags(0)
	}

	// Always use config file
	if err := codegen.Run(configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
