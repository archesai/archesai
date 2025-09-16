package codegen

// Common type name constants used throughout code generation
const (
	// Format types
	formatEmail = "email"
	formatInt64 = "int64"
	formatUUID  = "uuid"

	// SQL types
	sqlTypeUUID = "UUID"

	// Go types
	goTypeString    = "string"
	goTypeBool      = "bool"
	goTypeInt       = "int"
	goTypeInt32     = "int32"
	goTypeInt64     = "int64"
	goTypeFloat32   = "float32"
	goTypeFloat64   = "float64"
	goTypeTimeTime  = "time.Time"
	goTypeUUIDType  = "uuid.UUID"
	goTypeEmail     = "Email"
	goTypeEmailFull = "openapi_types.Email"
	goTypeMapString = "map[string]interface{}"

	// Go pointer types
	goTypePtrString   = "*string"
	goTypePtrUUID     = "*uuid.UUID"
	goTypePtrTime     = "*time.Time"
	goTypePtrInt32    = "*int32"
	goTypePtrFloat64  = "*float64"
	goTypePtrBool     = "*bool"
	goTypeUUIDLiteral = "UUID"
)
