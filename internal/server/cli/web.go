package cli

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Run the web UI server",
	Long: `Start the ArchesAI web UI server which serves the frontend
application for the platform.

This command will serve the built frontend assets and proxy API
requests to the configured API server.`,
	Example: `  archesai web --port=3000
  archesai web --api-url=http://localhost:8080`,
	RunE: runWeb,
}

var (
	webPort   int
	webAPIURL string
)

func init() {
	rootCmd.AddCommand(webCmd)

	// Local flags for web server
	webCmd.Flags().IntVar(&webPort, "port", 3000, "Port to bind the web server to")
	webCmd.Flags().StringVar(&webAPIURL, "api-url", "http://localhost:8080", "URL of the API server to proxy requests to")

	// Bind to viper
	if err := viper.BindPFlag("web.port", webCmd.Flags().Lookup("port")); err != nil {
		log.Fatalf("Failed to bind port flag: %v", err)
	}
	if err := viper.BindPFlag("web.api_url", webCmd.Flags().Lookup("api-url")); err != nil {
		log.Fatalf("Failed to bind api-url flag: %v", err)
	}
}

func runWeb(_ *cobra.Command, _ []string) error {
	port := viper.GetInt("web.port")
	apiURL := viper.GetString("web.api_url")

	log.Printf("🌐 Web server would start on port %d", port)
	log.Printf("   Proxying API requests to %s", apiURL)

	// TODO: Implement web server
	// This would:
	// 1. Serve static files from web/platform/dist
	// 2. Proxy /api/* requests to the API server
	// 3. Handle client-side routing for SPA

	return fmt.Errorf("web server not yet implemented")
}
