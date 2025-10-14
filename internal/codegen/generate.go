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

// Generator handles code generation from OpenAPI specifications
type Generator struct {
	templates        map[string]*template.Template
	filewriter       *FileWriter
	openAPIParser    *parsers.OpenAPIParser
	jsonSchemaParser *parsers.JSONSchemaParser
}

// NewGenerator creates a new code generator instance
func NewGenerator() *Generator {
	return &Generator{
		templates:        nil,
		filewriter:       nil,
		openAPIParser:    parsers.NewOpenAPIParser(),
		jsonSchemaParser: parsers.NewJSONSchemaParser(),
	}
}

// Initialize sets up the generator with templates and file writer
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

// GenerateAPI is the main generation function that orchestrates all code generation
func (g *Generator) GenerateAPI(specPath string) (string, error) {
	logLevel := os.Getenv("ARCHESAI_LOGGING_LEVEL")
	if logLevel == "" {
		logLevel = "error"
	}
	log := logger.New(logger.Config{Level: logLevel, Pretty: true})

	log.Info("Parsing OpenAPI specification", slog.String("path", specPath))
	openAPISchema, err := g.openAPIParser.Parse(specPath)
	if err != nil {
		return "", fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	// Set up the JSON schema parser with the OpenAPI document context
	g.jsonSchemaParser.WithOpenAPIDoc(openAPISchema)

	operations, err := parsers.ExtractOperations(openAPISchema)
	if err != nil {
		return "", fmt.Errorf("failed to process operations: %w", err)
	}

	schemas, err := parsers.ExtractComponentSchemas(openAPISchema)
	if err != nil {
		return "", fmt.Errorf("failed to process schemas: %w", err)
	}

	log.Info("Initialized state",
		slog.Int("schemas", len(schemas)),
		slog.Int("operations", len(operations)))

	// Buffer to collect all output
	var output bytes.Buffer

	log.Info("Generating models (DTOs, entities, value objects)")
	if err := g.GenerateSchemas(schemas); err != nil {
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
	if err := g.GenerateControllers(operations); err != nil {
		return "", fmt.Errorf("failed to generate handlers: %w", err)
	}
	log.Info("Generating events")
	if err := g.GenerateEvents(schemas); err != nil {
		return "", fmt.Errorf("failed to generate events: %w", err)
	}

	log.Info("Generating bootstrap files")
	if err := g.GenerateBootstrap(schemas, operations); err != nil {
		return "", fmt.Errorf("failed to generate bootstrap files: %w", err)
	}

	log.Info("Generating HCL database schema")
	if err := g.GenerateHCL(schemas); err != nil {
		return "", fmt.Errorf("failed to generate HCL schema: %w", err)
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

// GenerateJSONSchema generates Go structs from a JSON Schema file
func (g *Generator) GenerateJSONSchema(specPath string, outputDir string) (string, error) {
	logLevel := os.Getenv("ARCHESAI_LOGGING_LEVEL")
	if logLevel == "" {
		logLevel = "error"
	}
	log := logger.New(logger.Config{Level: logLevel, Pretty: true})

	log.Info("Parsing JSONSchema specification", slog.String("path", specPath))
	jsonSchema, err := g.jsonSchemaParser.Parse(specPath)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSONSchema spec: %w", err)
	}

	schema, err := g.jsonSchemaParser.ExtractSchema(jsonSchema, nil, "")
	if err != nil {
		return "", fmt.Errorf("failed to extract schema definition: %w", err)
	}

	// Buffer to collect all output
	var output bytes.Buffer

	log.Info("Generating models (DTOs, entities, value objects)")
	if err := g.generateModel(schema, &outputDir); err != nil {
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
