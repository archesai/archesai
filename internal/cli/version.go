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
	RunE:  runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(_ *cobra.Command, _ []string) error {
	slog.Info(fmt.Sprintf("version %s", version))
	slog.Info(fmt.Sprintf("commit: %s", commit))
	slog.Info(fmt.Sprintf("built: %s", buildDate))
	return nil
}
