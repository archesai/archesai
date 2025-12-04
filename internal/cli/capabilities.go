package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/archesai/archesai/internal/capabilities"
	"github.com/archesai/archesai/internal/cli/flags"
	"github.com/archesai/archesai/internal/tui"
)

const statusInstalled = "installed"

var capabilitiesCmd = &cobra.Command{
	Use:   "capabilities",
	Short: "Check system capabilities and dependencies",
	Long: `Check for required and optional system dependencies.

This command verifies that required tools are installed and displays
their versions. Use this to diagnose environment issues.

Examples:
  archesai capabilities           # Check all capabilities
  archesai capabilities --json    # Output as JSON for scripting
  archesai capabilities --all     # Show all, including optional`,
	RunE: runCapabilities,
}

func init() {
	rootCmd.AddCommand(capabilitiesCmd)
	flags.SetCapabilitiesFlags(capabilitiesCmd)
}

func outputJSON(result capabilities.CheckResult) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}

func outputTUI(result capabilities.CheckResult) error {
	runner := tui.NewRunner()

	// Title
	runner.PrintTitle("System Capabilities")
	runner.PrintNewline()

	// Required capabilities table
	requiredTable := tui.NewTable(
		"Required Dependencies",
		"Status",
		"Name",
		"Version",
		"Description",
	)
	for _, cap := range result.Capabilities {
		if !cap.Required {
			continue
		}
		status := "✓"
		if !cap.Found {
			status = "✗"
		}
		version := cap.Version
		if version == "" && cap.Found {
			version = statusInstalled
		} else if version == "" {
			version = "-"
		}
		requiredTable.AddRow(status, cap.Name, version, cap.Description)
	}
	runner.PrintTable(requiredTable)

	// Optional capabilities
	hasOptional := false
	for _, capability := range result.Capabilities {
		if !capability.Required {
			hasOptional = true
			break
		}
	}

	if hasOptional && flags.Capabilities.All {
		runner.PrintNewline()
		optionalTable := tui.NewTable(
			"Optional Dependencies",
			"Status",
			"Name",
			"Version",
			"Description",
		)
		for _, cap := range result.Capabilities {
			if cap.Required {
				continue
			}
			status := "✓"
			if !cap.Found {
				status = "○"
			}
			version := cap.Version
			if version == "" && cap.Found {
				version = statusInstalled
			} else if version == "" {
				version = "-"
			}
			optionalTable.AddRow(status, cap.Name, version, cap.Description)
		}
		runner.PrintTable(optionalTable)
	}

	runner.PrintNewline()

	// Summary
	requiredCount := 0
	foundCount := 0
	missingCount := 0

	for _, cap := range result.Capabilities {
		if cap.Required {
			requiredCount++
			if cap.Found {
				foundCount++
			} else {
				missingCount++
			}
		}
	}

	summary := tui.NewSummary("Summary")
	summary.AddCount("Required", requiredCount, "info")
	summary.AddCount("Found", foundCount, "success")

	if missingCount > 0 {
		summary.AddCount("Missing", missingCount, "error")
	}

	if result.AllRequired {
		summary.AddMessage("All required dependencies are installed", "success")
	} else {
		for _, cap := range result.RequiredMissing() {
			summary.AddMessage(fmt.Sprintf("Missing: %s (%s)", cap.Name, cap.Command), "error")
		}
	}

	runner.PrintSummary(summary)

	if !result.AllRequired {
		return fmt.Errorf("missing required dependencies")
	}

	return nil
}

func runCapabilities(_ *cobra.Command, _ []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	detector := capabilities.DefaultDetector(slog.Default())
	result := detector.Check(ctx)

	if flags.Capabilities.JSON {
		return outputJSON(result)
	}

	return outputTUI(result)
}
