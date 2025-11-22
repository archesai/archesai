package cli

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
)

var (
	// These will be set by ldflags during build.
	version   = "dev"
	commit    = "none"
	buildDate = "unknown"
)

// versionCmd represents the version command.
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Display version, commit hash, and build date of the Arches server.`,
	Run: func(_ *cobra.Command, _ []string) {
		slog.Info(fmt.Sprintf("version %s", version))
		slog.Info(fmt.Sprintf("commit: %s", commit))
		slog.Info(fmt.Sprintf("built: %s", buildDate))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
