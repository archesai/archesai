package parsers

import (
	"strings"

	"github.com/speakeasy-api/openapi/jsonschema/oas3"
)

// Constants for various type mappings and validations
const (
	// SQL Dialects
	SQLDialectPostgres = "postgresql"
	SQLDialectSQLite   = "sqlite"

	// Default SQL Types
	SQLTypeText      = "TEXT"
	SQLTypeInteger   = "INTEGER"
	SQLTypeBigInt    = "BIGINT"
	SQLTypeBoolean   = "BOOLEAN"
	SQLTypeTimestamp = "TIMESTAMPTZ"
	SQLTypeDateTime  = "DATETIME"
	SQLTypeDate      = "DATE"
	SQLTypeUUID      = "UUID"
	SQLTypeJSONB     = "JSONB"
	SQLTypeNumeric   = "NUMERIC"
	SQLTypeReal      = "REAL"
	SQLTypeDouble    = "DOUBLE PRECISION"
)

// ValidRepositoryOperations lists all valid repository operations
var ValidRepositoryOperations = []string{"create", "read", "update", "delete", "list"}

// ValidHTTPMethods lists all valid HTTP methods
var ValidHTTPMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

// ValidErrorHandling lists all valid error handling types
var ValidErrorHandling = []string{"error_return", "panic", "custom"}

// ValidLogLevels lists all valid logging levels
var ValidLogLevels = []string{"debug", "info", "warn", "error"}

// ValidDomainTypes lists all valid domain types
var ValidDomainTypes = []string{"entity", "aggregate", "valueobject", "dto"}

// OperationDef represents an API operation
type OperationDef struct {
	Name                string                      // Original operation ID from OpenAPI
	GoName              string                      // Go-friendly name (PascalCase)
	Method              string                      // HTTP method (GET, POST, etc.)
	Path                string                      // URL path
	Description         string                      // Operation description
	OperationID         string                      // Operation ID from OpenAPI
	Tags                []string                    // Operation tags
	Parameters          []ParamDef                  // All parameters (backward compat)
	PathParams          []ParamDef                  // Path parameters
	QueryParams         []ParamDef                  // Query parameters
	HeaderParams        []ParamDef                  // Header parameters
	Responses           []ResponseDef               // All responses
	Security            []SecurityDef               // Security requirements
	RequestBodyRequired bool                        // Whether request body is required
	RequestBodySchema   *ProcessedSchema            // Processed request body schema
	ResponseSchemas     map[string]*ProcessedSchema // Processed response schemas by status code
}

// GetSuccessResponse returns the first successful response (2xx status code)
func (o *OperationDef) GetSuccessResponse() *ResponseDef {
	for _, resp := range o.Responses {
		if resp.IsSuccess {
			return &resp
		}
	}
	return nil
}

// GetErrorResponses returns all error responses (non-2xx status codes)
func (o *OperationDef) GetErrorResponses() []ResponseDef {
	var errors []ResponseDef
	for _, resp := range o.Responses {
		if !resp.IsSuccess {
			errors = append(errors, resp)
		}
	}
	return errors
}

// HasBearerAuth checks if the operation requires bearer token authentication
func (o *OperationDef) HasBearerAuth() bool {
	for _, sec := range o.Security {
		if sec.Type == "http" && strings.EqualFold(sec.Scheme, "bearer") {
			return true
		}
	}
	return false
}

// HasCookieAuth checks if the operation requires cookie-based authentication
func (o *OperationDef) HasCookieAuth() bool {
	for _, sec := range o.Security {
		if sec.Type == "apiKey" && strings.EqualFold(sec.Scheme, "cookie") {
			return true
		}
	}
	return false
}

// ResponseDef represents a response in an operation
type ResponseDef struct {
	StatusCode  string // HTTP status code
	Description string // Response description
	Schema      string // Name of the response schema
	IsSuccess   bool   // Whether this is a success response (2xx)
	IsArray     bool   // Whether the response data is an array
}

// ParamDef represents a parameter in an operation
type ParamDef struct {
	Name        string // Parameter name
	In          string // Location (path, query, header)
	Type        string // Go type
	GoType      string // Generated Go type
	Format      string // Format hint (uuid, email, date-time)
	Required    bool   // Whether parameter is required
	Description string // Parameter description
	Schema      string // Referenced schema name if any
	Style       string // Parameter style (form, simple, etc.)
	Explode     bool   // Whether to explode array/object parameters
}

// SecurityDef represents a security requirement
type SecurityDef struct {
	Name   string   // Security scheme name
	Type   string   // Security type (http, apiKey, oauth2)
	Scheme string   // Security scheme (bearer for http, cookie for apiKey)
	Scopes []string // Required scopes
}

// FieldDef represents a schema field/property
type FieldDef struct {
	Name         string   // Original field name from schema
	FieldName    string   // Go field name (PascalCase)
	GoType       string   // Go type
	JSONTag      string   // JSON tag value
	YAMLTag      string   // YAML tag value
	Format       string   // Format hint (uuid, email, date-time)
	Enum         []string // Enum values if applicable
	Description  string   // Field description
	Required     bool     // Whether field is required
	Nullable     bool     // Whether field is nullable
	DefaultValue string   // Default value if any
	IsEnumType   bool     // Whether this field is an enum type
}

// ProcessedSchema contains all extracted information from a schema
type ProcessedSchema struct {
	// Basic info
	Name        string
	Title       string
	Description string
	Package     string

	// Fields
	Fields         []FieldDef
	RequiredFields []FieldDef
	OptionalFields []FieldDef

	// Type information
	IsEntity        bool
	IsAggregate     bool
	IsValueObject   bool
	IsDTO           bool
	IsEnum          bool
	HasDomainEvents bool

	// X-codegen extension
	XCodegen *XCodegenExtension

	// Enum values if applicable
	EnumValues []string

	// Original schema reference
	Schema *oas3.Schema
}
