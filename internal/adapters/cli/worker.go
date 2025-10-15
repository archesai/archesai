package cli

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// workerCmd represents the worker command.
var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Run background job worker",
	Long: `Start the Arches background job worker which processes
asynchronous tasks like email sending, webhook delivery, and
long-running data processing jobs.`,
	Example: `  archesai worker
  archesai worker --queues=default,email,webhook
  archesai worker --concurrency=10`,
	RunE: runWorker,
}

var (
	workerQueues      []string
	workerConcurrency int
)

func init() {
	rootCmd.AddCommand(workerCmd)

	// Local flags for worker
	workerCmd.Flags().
		StringSliceVar(&workerQueues, "queues", []string{"default"}, "Queues to process")
	workerCmd.Flags().IntVar(&workerConcurrency, "concurrency", 10, "Number of concurrent workers")

	// Bind to viper
	if err := viper.BindPFlag("worker.queues", workerCmd.Flags().Lookup("queues")); err != nil {
		slog.Error("Failed to bind queues flag", "err", err)
	}
	if err := viper.BindPFlag("worker.concurrency", workerCmd.Flags().Lookup("concurrency")); err != nil {
		slog.Error("Failed to bind concurrency flag", "err", err)
	}
}

func runWorker(_ *cobra.Command, _ []string) error {
	queues := viper.GetStringSlice("worker.queues")
	concurrency := viper.GetInt("worker.concurrency")

	slog.Info(fmt.Sprintf("⚙️  Worker would start processing queues: %v", queues))
	slog.Info(fmt.Sprintf("   Concurrency: %d workers", concurrency))

	// TODO: Implement worker
	// This would:
	// 1. Connect to Redis or other job queue
	// 2. Process jobs from specified queues
	// 3. Handle graceful shutdown

	return fmt.Errorf("worker not yet implemented")
}
