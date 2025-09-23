// Package templates provides template management and data structures for code generation.
package templates

// =============================================================================
// Base Template Data Structures
// =============================================================================

// TemplateData is the foundation all templates build on.
// Every generated file needs these basics.
type TemplateData struct {
	Package string   // The Go package name (e.g., "users", "accounts")
	Domain  string   // The logical domain (might differ from package)
	Imports []string // Required imports for generated code
}

// Operation represents an API operation.
type OperationData struct {
	Name                string // Original operation ID from OpenAPI
	GoName              string // Go-friendly name (PascalCase)
	Method              string
	Path                string
	Description         string
	OperationID         string
	Tags                []string
	Parameters          []ParamData
	PathParams          []ParamData
	QueryParams         []ParamData
	HeaderParams        []ParamData
	HasRequestBody      bool
	RequestBodyRequired bool
	RequestBodySchema   string
	ResponseType        string // Type for the success response
	Responses           []OperationResponse
	SuccessResponse     *OperationResponse
	ErrorResponses      []OperationResponse
	Security            []SecurityRequirement // Security requirements for this operation
	HasBearerAuth       bool                  // Whether this operation requires bearer authentication
	HasCookieAuth       bool                  // Whether this operation requires cookie authentication
}

// OperationResponse represents a response in an operation.
type OperationResponse struct {
	StatusCode  string
	Description string
	Schema      string
	IsSuccess   bool
}

// ParamData represents a parameter in an operation.
type ParamData struct {
	Name        string // Parameter name
	In          string // Location (path, query, header)
	Type        string // Go type (string, int, uuid.UUID, etc.)
	Format      string // Format hint (uuid, email, date-time)
	Required    bool   // Whether parameter is required
	Description string // Parameter description
}

// SecurityRequirement represents a security requirement for an operation.
type SecurityRequirement struct {
	Name   string   // Security scheme name (e.g., "bearerAuth", "sessionCookie")
	Type   string   // Security type (e.g., "http", "apiKey", "oauth2")
	Scheme string   // Security scheme (e.g., "bearer" for http type)
	Scopes []string // Required scopes for this security requirement
}

// EntityData represents a domain entity/model.
// Multiple templates need entity information in consistent format.
type EntityData struct {
	Domain            string       // Domain name
	Package           string       // Package name
	Name              string       // Entity name (singular, PascalCase)
	NameLower         string       // Entity name (singular, lowercase)
	NamePlural        string       // Entity name (plural, PascalCase)
	NamePluralLower   string       // Entity name (plural, lowercase)
	Type              string       // Entity type
	RepositoryName    string       // Custom repository interface name (e.g., SessionRepository)
	Fields            []FieldData  // Entity fields
	Operations        []string     // CRUD operations
	AdditionalMethods []MethodData // Additional custom methods
	UpdateExclude     []string     // Fields to exclude from updates
	HasEmailField     bool         // Whether entity has email fields
	HasTimeField      bool         // Whether entity has time fields
	HasUUIDField      bool         // Whether entity has UUID fields
	CodegenExtension  interface{}  // X-codegen configuration from schema
}

// FieldData represents a single field/property.
// Fields need consistent representation across types, repositories, and services.
type FieldData struct {
	Name          string   // Field name in schema
	FieldName     string   // Go field name (PascalCase)
	GoType        string   // Go type (string, int, uuid.UUID, etc.)
	SQLCType      string   // SQLC type representation
	SQLCFieldName string   // SQLC field name
	JSONTag       string   // JSON tag value
	YAMLTag       string   // YAML tag value
	Format        string   // Format hint (uuid, email, date-time)
	Enum          []string // Enum values if applicable
	Description   string   // Field description
	Required      bool     // Whether field is required
	Nullable      bool     // Whether field is nullable
	DefaultValue  string   // Default value if any
	IsEnumType    bool     // Whether this field is an enum type
}

// MethodData represents a custom method in a service or repository.
type MethodData struct {
	Name        string      // Method name
	Description string      // Method description
	Parameters  []ParamData // Method parameters
	Returns     []string    // Return types
}

// ConstantDef represents a constant definition (usually from enums).
type ConstantDef struct {
	Name           string   // Constant name
	Values         []string // Constant values
	ConstantPrefix string   // Optional prefix for constant names to avoid conflicts
}

// =============================================================================
// Schema and Type-Related Data Structures
// =============================================================================

// SchemaData represents a single schema for type generation.
type SchemaData struct {
	Name        string      // Schema name
	Description string      // Schema description
	Fields      []FieldData // Schema fields
	IsTypeAlias bool        // Whether this is a type alias
	TypeAlias   string      // The underlying type for type aliases
}

// TypeAliasData represents a type alias definition.
type TypeAliasData struct {
	Name           string // Alias name
	UnderlyingType string // Underlying type
}

// EventData represents a single event.
type EventData struct {
	Name        string // Event name
	Type        string // Event type
	PayloadType string // Payload type name
}

// TableData represents a database table.
type TableData struct {
	Name        string           // Table name
	Columns     []ColumnData     // Table columns
	PrimaryKey  []string         // Primary key columns
	ForeignKeys []ForeignKeyData // Foreign key constraints
	Indexes     []IndexData      // Table indexes
}

// ColumnData represents a database column.
type ColumnData struct {
	Name         string // Column name
	Type         string // SQL type
	Nullable     bool   // Whether column is nullable
	DefaultValue string // Default value
	IsPrimaryKey bool   // Whether part of primary key
}

// ForeignKeyData represents a foreign key constraint.
type ForeignKeyData struct {
	Name              string   // Constraint name
	Columns           []string // Local columns
	ReferencedTable   string   // Referenced table
	ReferencedColumns []string // Referenced columns
	OnDelete          string   // ON DELETE action
	OnUpdate          string   // ON UPDATE action
}

// IndexData represents a database index.
type IndexData struct {
	Name    string   // Index name
	Columns []string // Indexed columns
	Unique  bool     // Whether index is unique
}
