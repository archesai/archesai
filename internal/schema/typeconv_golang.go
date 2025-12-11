package schema

// Go type constants.
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

// Default validation constraints for OpenAPI schemas.
const (
	// String defaults
	DefaultStringMinLength = 0
	DefaultStringMaxLength = 10000

	// UUID defaults (exactly 36 chars: 8-4-4-4-12 with hyphens)
	DefaultUUIDLength  = 36
	DefaultUUIDExample = "550e8400-e29b-41d4-a716-446655440000"

	// Date-time defaults
	DefaultDateTimeMinLength = 1
	DefaultDateTimeMaxLength = 255
	DefaultDateTimeExample   = "2024-01-15T09:30:00Z"
	DefaultDateExample       = "2024-01-15"

	// Integer defaults
	DefaultIntegerMinimum = 0
	DefaultIntegerMaximum = 1000000

	// Int32 defaults
	DefaultInt32Minimum = 0
	DefaultInt32Maximum = 2147483647

	// Int64 defaults
	DefaultInt64Minimum = 0
	DefaultInt64Maximum = 9223372036854775807

	// Number/float defaults
	DefaultNumberMinimum = 0.0
	DefaultNumberMaximum = 1000000.0

	// Array defaults
	DefaultArrayMaxItems = 10000
)

// ToGoType converts a schema type and format to a Go type.
func ToGoType(schemaType, format string) string {
	switch schemaType {
	case TypeString:
		return stringFormatToGoType(format)
	case TypeInteger:
		return integerFormatToGoType(format)
	case TypeNumber:
		return numberFormatToGoType(format)
	case TypeBoolean:
		return GoTypeBool
	case TypeArray:
		return GoTypeSliceAny
	case TypeObject:
		return GoTypeMapString
	default:
		return GoTypeInterface
	}
}

// stringFormatToGoType converts a string format to a Go type.
func stringFormatToGoType(format string) string {
	switch format {
	case FormatDateTime, FormatDate:
		return GoTypeTime
	case FormatUUID:
		return GoTypeUUID
	case FormatEmail, FormatURI, FormatHostname, FormatPassword:
		return GoTypeString
	default:
		return GoTypeString
	}
}

// integerFormatToGoType converts an integer format to a Go type.
func integerFormatToGoType(format string) string {
	switch format {
	case FormatInt32:
		return GoTypeInt32
	case FormatInt64:
		return GoTypeInt64
	default:
		return GoTypeInt
	}
}

// numberFormatToGoType converts a number format to a Go type.
func numberFormatToGoType(format string) string {
	switch format {
	case FormatFloat:
		return GoTypeFloat32
	case FormatDouble:
		return GoTypeFloat64
	default:
		return GoTypeFloat64
	}
}
