// Command schema-to-struct generates Go struct definitions from OpenAPI schemas.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/archesai/archesai/internal/schema"
	"github.com/archesai/archesai/internal/templates"
)

var (
	specPath    = flag.String("spec", "", "Path to OpenAPI schema YAML file (required)")
	outputPath  = flag.String("output", "", "Output Go file path or directory (required)")
	packageName = flag.String("package", "", "Go package name (required)")
	rootType    = flag.String("root", "", "Root type name (defaults to filename without extension)")
)

// SchemasTemplateData matches the template data structure used by the generator.
type SchemasTemplateData struct {
	Package string
	Schema  *schema.Schema   // For multi-file mode (one schema per file)
	Schemas []*schema.Schema // For single-file mode (all schemas in one file)
}

// IsSingleMode returns true if generating all schemas in a single file.
func (d *SchemasTemplateData) IsSingleMode() bool {
	return len(d.Schemas) > 0
}

// GetSchemas returns the schemas to render.
func (d *SchemasTemplateData) GetSchemas() []*schema.Schema {
	if d.IsSingleMode() {
		return d.Schemas
	}
	if d.Schema != nil {
		return []*schema.Schema{d.Schema}
	}
	return nil
}

func main() {
	flag.Parse()

	if *specPath == "" || *outputPath == "" || *packageName == "" {
		fmt.Fprintln(
			os.Stderr,
			"Usage: schema-to-struct -spec <schema.yaml> -output <output.go|dir/> -package <pkg>",
		)
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Determine root type name
	rootName := *rootType
	if rootName == "" {
		base := filepath.Base(*specPath)
		rootName = base[:len(base)-len(filepath.Ext(base))]
	}

	baseDir := filepath.Dir(*specPath)

	// Use SchemaResolver to parse and resolve the schema
	resolver := schema.NewLoader(baseDir)
	rootSchema, err := resolver.LoadSchemaFile(*specPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing schema: %v\n", err)
		os.Exit(1)
	}

	// Override title with root name if specified
	if *rootType != "" {
		rootSchema.Title = rootName
		rootSchema.GoType = rootName
	}

	// Load templates
	tmpl, err := templates.LoadTemplates()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading templates: %v\n", err)
		os.Exit(1)
	}

	renderer := templates.NewRenderer(tmpl)

	// Get all loaded schemas (includes referenced schemas)
	allSchemas := resolver.Schemas()

	// Determine if output is a directory (ends with / or is a directory path without .go extension)
	isDir := strings.HasSuffix(*outputPath, "/") || !strings.HasSuffix(*outputPath, ".go")

	if isDir {
		// Generate separate files for each schema
		outputDir := strings.TrimSuffix(*outputPath, "/")
		if err := os.MkdirAll(outputDir, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "error creating output dir: %v\n", err)
			os.Exit(1)
		}

		for _, s := range allSchemas {
			var buf bytes.Buffer
			data := &SchemasTemplateData{
				Package: *packageName,
				Schema:  s,
			}

			if err := renderer.Render(&buf, "schema.go.tmpl", data); err != nil {
				fmt.Fprintf(os.Stderr, "error rendering template for %s: %v\n", s.Title, err)
				os.Exit(1)
			}

			// Generate filename from schema title (e.g., ConfigAPI -> configapi.gen.go)
			filename := strings.ToLower(s.Title) + ".gen.go"
			outputFile := filepath.Join(outputDir, filename)

			if err := os.WriteFile(outputFile, buf.Bytes(), 0o644); err != nil {
				fmt.Fprintf(os.Stderr, "error writing %s: %v\n", outputFile, err)
				os.Exit(1)
			}
		}

	} else {
		// Single file output - generate all schemas in one file
		var schemaSlice []*schema.Schema
		for _, s := range allSchemas {
			schemaSlice = append(schemaSlice, s)
		}

		var buf bytes.Buffer
		data := &SchemasTemplateData{
			Package: *packageName,
			Schemas: schemaSlice,
		}

		if err := renderer.Render(&buf, "schema.go.tmpl", data); err != nil {
			fmt.Fprintf(os.Stderr, "error rendering template: %v\n", err)
			os.Exit(1)
		}

		// Ensure output directory exists
		if err := os.MkdirAll(filepath.Dir(*outputPath), 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "error creating output dir: %v\n", err)
			os.Exit(1)
		}

		if err := os.WriteFile(*outputPath, buf.Bytes(), 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "error writing output: %v\n", err)
			os.Exit(1)
		}

	}
}
