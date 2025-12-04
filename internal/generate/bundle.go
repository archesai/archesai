package generate

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/archesai/archesai/internal/codegen"
	"github.com/archesai/archesai/internal/parsers"
)

// BundleSpec bundles the OpenAPI specification if needed.
// Returns the path to the bundled spec (or original if bundling was skipped).
func BundleSpec(opts *Options) error {
	shouldBundle := true
	if opts.Only != "" {
		shouldBundle = false
		for component := range strings.SplitSeq(opts.Only, ",") {
			if strings.TrimSpace(strings.ToLower(component)) == "bundle" {
				shouldBundle = true
				break
			}
		}
	}

	if shouldBundle {
		dir := filepath.Dir(opts.SpecPath)
		bundledPath := filepath.Join(dir, "openapi.bundled.yaml")

		merger := codegen.NewDefaultIncludeMerger()
		parser := parsers.NewOpenAPIParser().WithIncludeMerger(merger)
		if err := parser.Bundle(opts.SpecPath, bundledPath, opts.OrvalFix); err != nil {
			return fmt.Errorf("bundling failed: %w", err)
		}

		slog.Debug("Bundled OpenAPI specification", slog.String("output", bundledPath))
		opts.SpecPath = bundledPath
	}

	return nil
}
