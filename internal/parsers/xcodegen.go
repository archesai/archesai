package parsers

import (
	"fmt"
	"slices"
	"strings"

	"github.com/speakeasy-api/openapi/extensions"
	"github.com/speakeasy-api/openapi/jsonschema/oas3"
)

// XCodegenExtension represents the x-codegen extension structure
type XCodegenExtension struct {
	Type       string            `yaml:"type"` // Required: entity, aggregate, valueobject
	Repository *RepositoryConfig `yaml:"repository,omitempty"`
	Service    *ServiceConfig    `yaml:"service,omitempty"`
	Handler    *HandlerConfig    `yaml:"handler,omitempty"`
}

// MethodParameter defines a method parameter
type MethodParameter struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

// RepositoryConfig defines repository generation configuration
type RepositoryConfig struct {
	Name              string             `yaml:"name,omitempty"`
	Operations        []string           `yaml:"operations,omitempty"`
	Indices           []string           `yaml:"indices,omitempty"`
	AdditionalMethods []AdditionalMethod `yaml:"additional_methods,omitempty"`
	SoftDelete        bool               `yaml:"soft_delete,omitempty"`
	Timestamps        bool               `yaml:"timestamps,omitempty"`
	TableName         string             `yaml:"table_name,omitempty"`
}

// AdditionalMethod defines an additional repository method
type AdditionalMethod struct {
	Name    string   `yaml:"name"`
	Params  []string `yaml:"params,omitempty"`
	Returns string   `yaml:"returns"` // single, multiple
}

// ServiceConfig defines service layer generation configuration
type ServiceConfig struct {
	Generate           bool              `yaml:"generate,omitempty"`
	TransactionSupport bool              `yaml:"transaction_support,omitempty"`
	ErrorHandling      string            `yaml:"error_handling,omitempty"` // error_return, panic, custom
	Logging            LoggingConfig     `yaml:"logging,omitempty"`
	CustomOperations   []CustomOperation `yaml:"custom_operations,omitempty"`
}

// LoggingConfig defines logging configuration
type LoggingConfig struct {
	Level      string `yaml:"level,omitempty"` // debug, info, warn, error
	Structured bool   `yaml:"structured,omitempty"`
}

// CustomOperation defines a custom service operation
type CustomOperation struct {
	Name        string `yaml:"name"`
	OperationID string `yaml:"operation_id"`
	Method      string `yaml:"method"` // GET, POST, PUT, PATCH, DELETE
}

// HandlerConfig defines HTTP handler generation configuration
type HandlerConfig struct {
	Generate   bool             `yaml:"generate,omitempty"`
	PathPrefix string           `yaml:"path_prefix,omitempty"`
	Pagination PaginationConfig `yaml:"pagination,omitempty"`
}

// PaginationConfig defines pagination configuration
type PaginationConfig struct {
	DefaultLimit int `yaml:"default_limit,omitempty"`
	MaxLimit     int `yaml:"max_limit,omitempty"`
}

// XCodegenParser handles parsing of x-codegen extensions
type XCodegenParser struct {
}

// NewXCodegenParser creates a new x-codegen extension parser
func NewXCodegenParser() *XCodegenParser {
	return &XCodegenParser{}
}

// ParseExtension parses an x-codegen extension from a schema
func (p *XCodegenParser) ParseExtension(
	ext extensions.Extension,
	schemaName string,
) (*XCodegenExtension, error) {
	if ext == nil {
		return nil, nil
	}

	var xcodegen XCodegenExtension
	// ext is already *yaml.Node, so we can decode directly
	if err := ext.Decode(&xcodegen); err != nil {
		return nil, fmt.Errorf(
			"failed to decode x-codegen extension for schema '%s': %w",
			schemaName,
			err,
		)
	}

	// Validate the extension
	if err := p.validateExtension(&xcodegen, schemaName); err != nil {
		return nil, err
	}

	return &xcodegen, nil
}

// ParseSchemaExtensions parses x-codegen extensions from a schema
func (p *XCodegenParser) ParseSchemaExtensions(
	schema *oas3.Schema,
	schemaName string,
) (*XCodegenExtension, error) {
	if schema == nil || schema.Extensions == nil {
		return nil, nil
	}

	ext := schema.Extensions.GetOrZero("x-codegen")
	if ext == nil {
		return nil, nil
	}

	return p.ParseExtension(ext, schemaName)
}

// validateExtension validates an x-codegen extension
func (p *XCodegenParser) validateExtension(ext *XCodegenExtension, schemaName string) error {
	if ext == nil {
		return nil
	}

	// Validate type (required)
	if ext.Type == "" {
		return fmt.Errorf("x-codegen.type is required for schema '%s'", schemaName)
	}

	validTypes := []string{"entity", "aggregate", "valueobject"}
	if !contains(validTypes, ext.Type) {
		return fmt.Errorf(
			"invalid type '%s' in schema '%s'. Valid types: %s",
			ext.Type,
			schemaName,
			strings.Join(validTypes, ", "),
		)
	}

	// Validate repository configuration
	if ext.Repository != nil {
		if err := p.validateRepositoryConfig(ext.Repository, schemaName); err != nil {
			return err
		}
	}

	// Validate service configuration
	if ext.Service != nil {
		if err := p.validateServiceConfig(ext.Service, schemaName); err != nil {
			return err
		}
	}

	// Validate handler configuration
	if ext.Handler != nil {
		if err := p.validateHandlerConfig(ext.Handler, schemaName); err != nil {
			return err
		}
	}

	return nil
}

// validateRepositoryConfig validates repository configuration
func (p *XCodegenParser) validateRepositoryConfig(cfg *RepositoryConfig, schemaName string) error {
	if cfg == nil {
		return nil
	}

	// Validate operations
	validOps := []string{"create", "read", "update", "delete", "list"}
	for _, op := range cfg.Operations {
		if !contains(validOps, op) {
			// Just log warning, don't fail
			fmt.Printf(
				"Warning: invalid repository operation '%s' in schema '%s'. Valid operations: %s\n",
				op,
				schemaName,
				strings.Join(validOps, ", "),
			)
		}
	}

	// Validate additional methods
	for i, method := range cfg.AdditionalMethods {
		if method.Name == "" {
			return fmt.Errorf(
				"additional method at index %d missing name in schema '%s'",
				i,
				schemaName,
			)
		}
		if method.Returns != "single" && method.Returns != "multiple" {
			return fmt.Errorf(
				"invalid returns value '%s' for method '%s' in schema '%s'. Use 'single' or 'multiple'",
				method.Returns,
				method.Name,
				schemaName,
			)
		}
	}

	return nil
}

// validateServiceConfig validates service configuration
func (p *XCodegenParser) validateServiceConfig(cfg *ServiceConfig, schemaName string) error {
	if cfg == nil {
		return nil
	}

	// Validate error handling
	validErrorHandling := []string{"error_return", "panic", "custom"}
	if cfg.ErrorHandling != "" && !contains(validErrorHandling, cfg.ErrorHandling) {
		fmt.Printf("Warning: invalid error handling '%s' in schema '%s'. Valid options: %s\n",
			cfg.ErrorHandling, schemaName, strings.Join(validErrorHandling, ", "))
	}

	// Validate logging level
	if cfg.Logging.Level != "" {
		validLevels := []string{"debug", "info", "warn", "error"}
		if !contains(validLevels, cfg.Logging.Level) {
			fmt.Printf("Warning: invalid logging level '%s' in schema '%s'. Valid levels: %s\n",
				cfg.Logging.Level, schemaName, strings.Join(validLevels, ", "))
		}
	}

	// Validate custom operations
	for i, op := range cfg.CustomOperations {
		if op.Name == "" || op.OperationID == "" || op.Method == "" {
			return fmt.Errorf(
				"custom operation at index %d missing required fields in schema '%s'",
				i,
				schemaName,
			)
		}

		validMethods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
		if !contains(validMethods, op.Method) {
			return fmt.Errorf(
				"invalid HTTP method '%s' for operation '%s' in schema '%s'. Valid methods: %s",
				op.Method,
				op.Name,
				schemaName,
				strings.Join(validMethods, ", "),
			)
		}
	}

	return nil
}

// validateHandlerConfig validates handler configuration
func (p *XCodegenParser) validateHandlerConfig(cfg *HandlerConfig, schemaName string) error {
	if cfg == nil {
		return nil
	}

	// Validate path prefix
	if cfg.PathPrefix != "" && !strings.HasPrefix(cfg.PathPrefix, "/") {
		fmt.Printf(
			"Warning: path prefix '%s' should start with '/' in schema '%s'. Use '%s' instead\n",
			cfg.PathPrefix,
			schemaName,
			"/"+cfg.PathPrefix,
		)
	}

	// Validate pagination limits
	if cfg.Pagination.DefaultLimit > 0 {
		if cfg.Pagination.DefaultLimit > 1000 {
			fmt.Printf(
				"Warning: default_limit %d is very high in schema '%s'. Consider using a value <= 100 for better performance\n",
				cfg.Pagination.DefaultLimit,
				schemaName,
			)
		}
	}

	if cfg.Pagination.MaxLimit > 0 {
		if cfg.Pagination.MaxLimit > 1000 {
			fmt.Printf(
				"Warning: max_limit %d is very high in schema '%s'. Consider using a value <= 1000 for better performance\n",
				cfg.Pagination.MaxLimit,
				schemaName,
			)
		}

		if cfg.Pagination.DefaultLimit > cfg.Pagination.MaxLimit {
			return fmt.Errorf("default_limit %d exceeds max_limit %d in schema '%s'",
				cfg.Pagination.DefaultLimit, cfg.Pagination.MaxLimit, schemaName)
		}
	}

	return nil
}

// Parse is a simpler helper that just parses extensions
func (p *XCodegenParser) Parse(extensions extensions.Extensions) *XCodegenExtension {
	ext := extensions.GetOrZero("x-codegen")
	if ext == nil {
		return nil
	}

	var xcodegen XCodegenExtension
	if err := ext.Decode(&xcodegen); err != nil {
		// Just return nil if we can't decode
		return nil
	}

	return &xcodegen
}

// contains checks if a string is in a slice
func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}
