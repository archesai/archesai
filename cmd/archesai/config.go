package main

import (
	"github.com/spf13/cobra"
)

// configCmd represents the config parent command.
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Arches configuration",
	Long:  `Commands for managing Arches configuration.`,
}

func init() {
	rootCmd.AddCommand(configCmd)
}
