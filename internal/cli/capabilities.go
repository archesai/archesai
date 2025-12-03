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
	"github.com/archesai/archesai/internal/tui"
)

var (
	capabilitiesJSON bool
	capabilitiesAll  bool
	capabilitiesTUI  bool
)

var capabilitiesCmd = &cobra.Command{
	Use:   "capabilities",
	Short: "Check system capabilities and dependencies",
	Long: `Check for required and optional system dependencies.

This command verifies that required tools are installed and displays
their versions. Use this to diagnose environment issues.

Examples:
  archesai capabilities           # Check all capabilities
  archesai capabilities --json    # Output as JSON for scripting
  archesai capabilities --all     # Show all, including optional
  archesai capabilities --tui     # Display with TUI styling`,
	RunE: func(_ *cobra.Command, _ []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		detector := capabilities.DefaultDetector(slog.Default())
		result := detector.Check(ctx)

		if capabilitiesJSON {
			return outputJSON(result)
		}

		if capabilitiesTUI {
			return outputTUI(result)
		}

		return outputTable(result)
	},
}

func init() {
	rootCmd.AddCommand(capabilitiesCmd)
	capabilitiesCmd.Flags().BoolVar(&capabilitiesJSON, "json", false, "Output as JSON")
	capabilitiesCmd.Flags().BoolVar(
		&capabilitiesAll, "all", false, "Show all capabilities including optional",
	)
	capabilitiesCmd.Flags().
		BoolVarP(&capabilitiesTUI, "tui", "t", false, "Display with TUI styling")
}

func outputJSON(result capabilities.CheckResult) error {
	type jsonCapability struct {
		Name        string `json:"name"`
		Command     string `json:"command"`
		Version     string `json:"version,omitempty"`
		Required    bool   `json:"required"`
		Found       bool   `json:"found"`
		Error       string `json:"error,omitempty"`
		Description string `json:"description"`
	}

	type jsonResult struct {
		Capabilities []jsonCapability `json:"capabilities"`
		AllRequired  bool             `json:"allRequired"`
		AllFound     bool             `json:"allFound"`
	}

	out := jsonResult{
		Capabilities: make([]jsonCapability, len(result.Capabilities)),
		AllRequired:  result.AllRequired,
		AllFound:     result.AllFound,
	}

	for i, cap := range result.Capabilities {
		jc := jsonCapability{
			Name:        cap.Name,
			Command:     cap.Command,
			Version:     cap.Version,
			Required:    cap.Required,
			Found:       cap.Found,
			Description: cap.Description,
		}
		if cap.Error != nil {
			jc.Error = cap.Error.Error()
		}
		out.Capabilities[i] = jc
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
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
			version = "installed"
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

	if hasOptional && capabilitiesAll {
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
				version = "installed"
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

func outputTable(result capabilities.CheckResult) error {
	fmt.Println("System Capabilities")
	fmt.Println("===================")
	fmt.Println()

	// Required capabilities
	fmt.Println("Required:")
	for _, cap := range result.Capabilities {
		if !cap.Required && !capabilitiesAll {
			continue
		}
		if cap.Required {
			printCapability(cap)
		}
	}

	// Optional capabilities
	hasOptional := false
	for _, capability := range result.Capabilities {
		if !capability.Required {
			hasOptional = true
			break
		}
	}

	if hasOptional {
		fmt.Println()
		fmt.Println("Optional:")
		for _, capability := range result.Capabilities {
			if !capability.Required {
				printCapability(capability)
			}
		}
	}

	fmt.Println()

	// Summary
	if !result.AllRequired {
		missing := result.RequiredMissing()
		fmt.Printf("Missing required dependencies (%d):\n", len(missing))
		for _, cap := range missing {
			fmt.Printf("  - %s (%s)\n", cap.Name, cap.Command)
		}
		return fmt.Errorf("missing required dependencies")
	}

	fmt.Println("All required dependencies are installed.")
	return nil
}

func printCapability(c capabilities.Capability) {
	status := "OK"
	if !c.Found {
		status = "MISSING"
	}

	version := c.Version
	if version == "" && c.Found {
		version = "installed"
	} else if version == "" {
		version = "-"
	}

	fmt.Printf("  [%s] %-12s %-20s %s\n", status, c.Name, version, c.Description)
}
