package main

import (
	"github.com/spf13/cobra"
)

// specCmd represents the spec parent command
var specCmd = &cobra.Command{
	Use:   "spec",
	Short: "OpenAPI specification utilities",
	Long:  `Commands for working with OpenAPI specifications.`,
}

func init() {
	rootCmd.AddCommand(specCmd)
}
