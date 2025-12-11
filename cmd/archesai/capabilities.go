package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/archesai/archesai/cmd/archesai/flags"
	"github.com/archesai/archesai/internal/capabilities"
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

// Styles for capabilities output
var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	dimStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
)

func outputConsole(result capabilities.CheckResult) error {
	// Title
	fmt.Println(titleStyle.Render("System Capabilities"))
	fmt.Println()

	// Required capabilities
	fmt.Println(titleStyle.Render("Required Dependencies"))
	for _, cap := range result.Capabilities {
		if !cap.Required {
			continue
		}
		status := successStyle.Render("✓")
		if !cap.Found {
			status = errorStyle.Render("✗")
		}
		version := cap.Version
		if version == "" && cap.Found {
			version = statusInstalled
		} else if version == "" {
			version = "-"
		}
		fmt.Printf(
			"  %s %-12s %-15s %s\n",
			status,
			cap.Name,
			dimStyle.Render(version),
			cap.Description,
		)
	}

	// Optional capabilities
	hasOptional := false
	for _, capability := range result.Capabilities {
		if !capability.Required {
			hasOptional = true
			break
		}
	}

	if hasOptional && flags.Capabilities.All {
		fmt.Println()
		fmt.Println(titleStyle.Render("Optional Dependencies"))
		for _, cap := range result.Capabilities {
			if cap.Required {
				continue
			}
			status := successStyle.Render("✓")
			if !cap.Found {
				status = dimStyle.Render("○")
			}
			version := cap.Version
			if version == "" && cap.Found {
				version = statusInstalled
			} else if version == "" {
				version = "-"
			}
			fmt.Printf(
				"  %s %-12s %-15s %s\n",
				status,
				cap.Name,
				dimStyle.Render(version),
				cap.Description,
			)
		}
	}

	fmt.Println()

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

	fmt.Printf(
		"Required: %d  Found: %s",
		requiredCount,
		successStyle.Render(fmt.Sprintf("%d", foundCount)),
	)
	if missingCount > 0 {
		fmt.Printf("  Missing: %s", errorStyle.Render(fmt.Sprintf("%d", missingCount)))
	}
	fmt.Println()

	if result.AllRequired {
		fmt.Println(successStyle.Render("All required dependencies are installed"))
	} else {
		for _, cap := range result.RequiredMissing() {
			fmt.Println(errorStyle.Render(fmt.Sprintf("Missing: %s (%s)", cap.Name, cap.Command)))
		}
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

	return outputConsole(result)
}
