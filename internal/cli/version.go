package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// These will be set by ldflags during build
	version   = "dev"
	commit    = "none"
	buildDate = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Display version, commit hash, and build date of the ArchesAI server.`,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("ArchesAI version %s\n", version)
		fmt.Printf("  commit: %s\n", commit)
		fmt.Printf("  built:  %s\n", buildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
