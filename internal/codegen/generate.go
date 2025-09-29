// Package codegen provides unified code generation from OpenAPI schemas with x-codegen extensions.
package codegen

import (
	"bytes"
	"fmt"
	"go/format"
	"log/slog"
	"os"
	"text/template"

	"github.com/archesai/archesai/internal/parsers"
	"github.com/archesai/archesai/internal/shared/logger"
)

// Schema type constants
const (
	schemaTypeEntity      = "entity"
	schemaTypeValueObject = "valueobject"
)

type Generator struct {
	templates  map[string]*template.Template
	filewriter *FileWriter
}

func NewGenerator() *Generator {
	return &Generator{
		templates:  nil,
		filewriter: nil,
	}
}

func (g *Generator) Initialize() error {
	filewriter := NewFileWriter()
	filewriter.WithOverwrite(true)
	filewriter.WithHeader(DefaultHeader())
	g.filewriter = filewriter

	templates, err := LoadTemplates()
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}
	g.templates = templates

	return nil
}

// Generate is the main generation function that orchestrates all code generation
func (g *Generator) GenerateAPI(specPath string) (string, error) {
	logLevel := os.Getenv("ARCHESAI_LOGGING_LEVEL")
	if logLevel == "" {
		logLevel = "error"
	}
	log := logger.New(logger.Config{Level: logLevel, Pretty: true})

	log.Info("Parsing OpenAPI specification", slog.String("path", specPath))
	openAPISchema, warnings, err := parsers.ParseOpenAPI(specPath)
	if err != nil {
		return "", fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	// Log any warnings
	for _, warning := range warnings {
		log.Warn("OpenAPI warning", slog.String("warning", warning))
	}

	operations := parsers.ExtractOperations(openAPISchema)
	schemas, err := parsers.ProcessAllSchemas(openAPISchema)
	if err != nil {
		return "", fmt.Errorf("failed to process schemas: %w", err)
	}

	log.Info("Initialized state",
		slog.Int("schemas", len(schemas)),
		slog.Int("operations", len(operations)))

	// Buffer to collect all output
	var output bytes.Buffer

	log.Info("Generating models (DTOs, entities, value objects)")
	if err := g.GenerateModels(schemas); err != nil {
		return "", fmt.Errorf("failed to generate models: %w", err)
	}

	log.Info("Generating repositories")
	if err := g.GenerateRepositories(schemas); err != nil {
		return "", fmt.Errorf("failed to generate repositories: %w", err)
	}

	log.Info("Generating command and query handlers")
	if err := g.GenerateCommandQueryHandlers(operations); err != nil {
		return "", fmt.Errorf("failed to generate handlers: %w", err)
	}

	log.Info("Generating handlers")
	if err := g.GenerateControllers(operations, schemas); err != nil {
		return "", fmt.Errorf("failed to generate handlers: %w", err)
	}
	log.Info("Generating events")
	if err := g.GenerateEvents(schemas); err != nil {
		return "", fmt.Errorf("failed to generate events: %w", err)
	}

	// 10. Format output if needed
	outputStr := output.String()
	if outputStr != "" {
		formatted, err := format.Source([]byte(outputStr))
		if err != nil {
			log.Warn("Failed to format output", slog.String("error", err.Error()))
		} else {
			outputStr = string(formatted)
		}
	}

	log.Info("Code generation completed successfully")
	return outputStr, nil
}

// Generate is the main generation function that orchestrates all code generation
func (g *Generator) GenerateSchema(specPath string, outputDir string) (string, error) {
	logLevel := os.Getenv("ARCHESAI_LOGGING_LEVEL")
	if logLevel == "" {
		logLevel = "error"
	}
	log := logger.New(logger.Config{Level: logLevel, Pretty: true})

	log.Info("Parsing OpenAPI specification", slog.String("path", specPath))
	jsonSchema, xcodegen, err := parsers.ParseJSONSchema(specPath)
	if err != nil {
		return "", fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	title := jsonSchema.GetTitle()
	if title == "" {
		return "", fmt.Errorf("JSON Schema must have a title")
	}
	schema, err := parsers.ProcessSchema(jsonSchema, title)
	if err != nil {
		return "", fmt.Errorf("failed to process schemas: %w", err)
	}
	log.Info("Initialized state",
		slog.String("schema", title),
	)

	// Buffer to collect all output
	var output bytes.Buffer

	log.Info("Generating models (DTOs, entities, value objects)")
	if err := g.generateModel(schema, schema.Title, xcodegen.SchemaType, &outputDir); err != nil {
		return "", fmt.Errorf("failed to generate models: %w", err)
	}

	// 10. Format output if needed
	log.Info("Formatting output")
	outputStr := output.String()
	if outputStr != "" {
		formatted, err := format.Source([]byte(outputStr))
		if err != nil {
			log.Warn("Failed to format output", slog.String("error", err.Error()))
		} else {
			outputStr = string(formatted)
		}
	}

	log.Info("Code generation completed successfully")
	return outputStr, nil
}
