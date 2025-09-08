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

	flag.StringVar(&configPath, "config", "codegen.yaml", "Config file path")
	flag.StringVar(&configPath, "c", "codegen.yaml", "Config file path (shorthand)")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")
	flag.BoolVar(&verbose, "v", false, "Verbose output (shorthand)")
	flag.Parse()

	if verbose {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	} else {
		log.SetFlags(0)
	}

	if err := codegen.Run(configPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
