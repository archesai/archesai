// Package codegen provides unified code generation from OpenAPI schemas with x-codegen extensions.
package codegen

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/archesai/archesai/internal/logger"
	"github.com/archesai/archesai/internal/parsers"
	"github.com/archesai/archesai/internal/templates"
)

// stringPtr returns a pointer to the given string.
func stringPtr(s string) *string {
	return &s
}

// GetDefaultConfig returns a new Config with default values.
func GetDefaultConfig() *CodegenConfig {
	output := "internal"
	return &CodegenConfig{
		Generators: &struct {
			Cache *struct {
				Interface *string "json:\"interface,omitempty\" yaml:\"interface,omitempty\""
				Memory    *string "json:\"memory,omitempty\" yaml:\"memory,omitempty\""
				Redis     *string "json:\"redis,omitempty\" yaml:\"redis,omitempty\""
			} "json:\"cache,omitempty\" yaml:\"cache,omitempty\""
			EchoServer *string "json:\"echo_server,omitempty\" yaml:\"echo_server,omitempty\""
			Events     *struct {
				Interface *string "json:\"interface,omitempty\" yaml:\"interface,omitempty\""
				Nats      *string "json:\"nats,omitempty\" yaml:\"nats,omitempty\""
				Redis     *string "json:\"redis,omitempty\" yaml:\"redis,omitempty\""
			} "json:\"events,omitempty\" yaml:\"events,omitempty\""
			Repository *struct {
				Interface *string "json:\"interface,omitempty\" yaml:\"interface,omitempty\""
				Postgres  *string "json:\"postgres,omitempty\" yaml:\"postgres,omitempty\""
				Sqlite    *string "json:\"sqlite,omitempty\" yaml:\"sqlite,omitempty\""
			} "json:\"repository,omitempty\" yaml:\"repository,omitempty\""
			SQL *struct {
				Dialect   *string "json:\"dialect,omitempty\" yaml:\"dialect,omitempty\""
				QueryDir  *string "json:\"query_dir,omitempty\" yaml:\"query_dir,omitempty\""
				SchemaDir *string "json:\"schema_dir,omitempty\" yaml:\"schema_dir,omitempty\""
			} "json:\"sql,omitempty\" yaml:\"sql,omitempty\""
			Service *string "json:\"service,omitempty\" yaml:\"service,omitempty\""
			Types   *string "json:\"types,omitempty\" yaml:\"types,omitempty\""
		}{
			Repository: &struct {
				Interface *string "json:\"interface,omitempty\" yaml:\"interface,omitempty\""
				Postgres  *string "json:\"postgres,omitempty\" yaml:\"postgres,omitempty\""
				Sqlite    *string "json:\"sqlite,omitempty\" yaml:\"sqlite,omitempty\""
			}{
				Interface: stringPtr("repository.gen.go"),
				Postgres:  stringPtr("repository_postgres.gen.go"),
				Sqlite:    stringPtr("repository_sqlite.gen.go"),
			},
			EchoServer: stringPtr("handler.gen.go"),
			Service:    stringPtr("service.gen.go"),
			Types:      stringPtr("types.gen.go"),
		},
		Output:  &output,
		Openapi: "api/openapi.yaml",
	}

}

// Run executes the unified code generator with the given configuration.
func Run(config *CodegenConfig) error {
	// Create logger with error level by default (only show errors)
	// Set ARCHESAI_LOG_LEVEL=debug to see debug logs
	logLevel := os.Getenv("ARCHESAI_LOG_LEVEL")
	if logLevel == "" {
		logLevel = "error"
	}
	log := logger.New(logger.Config{Level: logLevel, Pretty: true})

	// Create parser and file writer
	parser := parsers.NewParser()
	fileWriter := templates.NewFileWriter()

	// Parse the OpenAPI specification
	openAPISchema, warnings, err := parser.Parse(config.Openapi)
	if err != nil {
		return fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}
	parser.OpenAPI = openAPISchema

	// Configure file writer
	fileWriter.WithOverwrite(true) // Always overwrite generated files
	fileWriter.WithHeader(templates.DefaultHeader())

	// Load templates from templates package
	templateMap, err := templates.LoadTemplates()
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	// Extract schemas from OpenAPI
	schemas := parser.OpenAPI.ExtractSchemas()

	// Add parsing warnings to parser
	if len(warnings) > 0 {
		for _, w := range warnings {
			log.Warn("OpenAPI parsing warning", "message", w)
		}
	}

	// Common context for all generators
	ctx := &GeneratorContext{
		Config:     config,
		Parser:     parser,
		Schemas:    schemas,
		Templates:  templateMap,
		FileWriter: fileWriter,
	}

	// Create all generators using the new unified system
	allGenerators := CreateGenerators(parser, log)

	// Define generator execution order
	// Types must be generated first as other generators may depend on them
	generatorOrder := []string{
		"types",
		"sql",
		"repository",
		"cache",
		"events",
		"service",
		"echo_server",
	}

	// Run each generator in order
	for _, name := range generatorOrder {
		gen, exists := allGenerators[name]
		if !exists {
			log.Debug("Generator not found", slog.String("name", name))
			continue
		}

		log.Debug("Running generator", slog.String("name", name))
		// Each generator checks if it's enabled internally
		if err := gen.Generate(ctx); err != nil {
			return fmt.Errorf("%s generation failed: %w", name, err)
		}
	}

	// Report warnings
	if warnings := parser.GetWarnings(); len(warnings) > 0 {
		log.Warn("Parsing warnings found")
		for _, w := range warnings {
			log.Debug("Warning", slog.String("message", w))
		}
	}

	log.Debug("Code generation completed successfully")
	return nil
}
