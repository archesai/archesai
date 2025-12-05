// Package spec provides data structures and constants for OpenAPI specifications.
package spec

// Schema type constants
const (
	SchemaTypeString  = "string"
	SchemaTypeInteger = "integer"
	SchemaTypeNumber  = "number"
	SchemaTypeBoolean = "boolean"
	SchemaTypeArray   = "array"
	SchemaTypeObject  = "object"
)

// Format constants
const (
	FormatDateTime = "date-time"
	FormatDate     = "date"
	FormatUUID     = "uuid"
	FormatEmail    = "email"
	FormatURI      = "uri"
	FormatHostname = "hostname"
	FormatInt32    = "int32"
	FormatInt64    = "int64"
	FormatFloat    = "float"
	FormatDouble   = "double"
)

// Go type constants
const (
	GoTypeInterface = "any"
	GoTypeString    = "string"
	GoTypeInt       = "int"
	GoTypeInt32     = "int32"
	GoTypeInt64     = "int64"
	GoTypeFloat32   = "float32"
	GoTypeFloat64   = "float64"
	GoTypeBool      = "bool"
	GoTypeTime      = "time.Time"
	GoTypeUUID      = "uuid.UUID"
	GoTypeMapString = "map[string]any"
	GoTypeSliceAny  = "[]any"
)
