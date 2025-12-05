package typeconv

import (
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"

	"github.com/archesai/archesai/internal/spec"
)

// SchemaToGoType converts a JSON Schema to a Go type with proper package qualification
func SchemaToGoType(schema *base.Schema, doc *v3.Document, currentPackage string) string {
	if schema == nil {
		return spec.GoTypeInterface
	}

	// Get the types array from the schema
	if len(schema.Type) == 0 {
		return spec.GoTypeInterface
	}

	// Use the first type (most schemas have only one type)
	schemaType := schema.Type[0]

	// Delegate to type-specific handlers
	switch schemaType {
	case spec.SchemaTypeString:
		return stringToGoType(schema)
	case spec.SchemaTypeInteger:
		return integerToGoType(schema)
	case spec.SchemaTypeNumber:
		return numberToGoType(schema)
	case spec.SchemaTypeBoolean:
		return spec.GoTypeBool
	case spec.SchemaTypeArray:
		return arrayToGoType(schema, doc, currentPackage)
	case spec.SchemaTypeObject:
		return objectToGoType(schema, doc, currentPackage)
	default:
		return spec.GoTypeInterface
	}
}

// stringToGoType converts a string schema to a Go type
func stringToGoType(schema *base.Schema) string {
	if schema.Format == "" {
		return spec.GoTypeString
	}

	switch schema.Format {
	case spec.FormatDateTime, spec.FormatDate:
		return spec.GoTypeTime
	case spec.FormatUUID:
		return spec.GoTypeUUID
	case spec.FormatEmail, spec.FormatURI, spec.FormatHostname:
		return spec.GoTypeString
	default:
		return spec.GoTypeString
	}
}

// integerToGoType converts an integer schema to a Go type
func integerToGoType(schema *base.Schema) string {
	if schema.Format == "" {
		return spec.GoTypeInt
	}

	switch schema.Format {
	case spec.FormatInt32:
		return spec.GoTypeInt32
	case spec.FormatInt64:
		return spec.GoTypeInt64
	default:
		return spec.GoTypeInt
	}
}

// numberToGoType converts a number schema to a Go type
func numberToGoType(schema *base.Schema) string {
	if schema.Format == "" {
		return spec.GoTypeFloat64
	}

	switch schema.Format {
	case spec.FormatFloat:
		return spec.GoTypeFloat32
	case spec.FormatDouble:
		return spec.GoTypeFloat64
	default:
		return spec.GoTypeFloat64
	}
}

// arrayToGoType converts an array schema to a Go type
func arrayToGoType(schema *base.Schema, doc *v3.Document, currentPackage string) string {
	if schema.Items == nil || schema.Items.A == nil {
		return spec.GoTypeSliceAny
	}

	itemSchema := schema.Items.A.Schema()
	if itemSchema == nil {
		return spec.GoTypeSliceAny
	}

	itemTypes := itemSchema.Type
	// Check if the array item is an object with properties
	if len(itemTypes) > 0 && itemTypes[0] == spec.SchemaTypeObject &&
		itemSchema.Properties != nil && itemSchema.Properties.Len() > 0 {
		// For arrays of objects with properties, we'll need special handling
		// This will be handled by the caller (processSchemaProperties)
		return "[]map[string]any"
	}

	itemType := SchemaToGoType(itemSchema, doc, currentPackage)
	return "[]" + itemType
}

// objectToGoType converts an object schema to a Go type
func objectToGoType(schema *base.Schema, doc *v3.Document, currentPackage string) string {
	if schema.AdditionalProperties != nil && schema.AdditionalProperties.IsA() {
		valueSchema := schema.AdditionalProperties.A.Schema()
		if valueSchema != nil {
			valueType := SchemaToGoType(valueSchema, doc, currentPackage)
			return "map[string]" + valueType
		}
	}
	// For objects with properties, just use map[string]any
	return spec.GoTypeMapString
}
