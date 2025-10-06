// Package parsers provides OpenAPI and JSON Schema parsing functionality
package parsers

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/speakeasy-api/openapi/jsonschema/oas3"
	"github.com/speakeasy-api/openapi/marshaller"
	"github.com/speakeasy-api/openapi/openapi"
)

const (
	tagOmitEmpty = ",omitempty"
)

// JSONSchemaParser handles parsing JSON Schema files
type JSONSchemaParser struct {
	openAPIDoc *openapi.OpenAPI
}

// NewJSONSchemaParser creates a new JSONSchemaParser instance
func NewJSONSchemaParser() *JSONSchemaParser {
	return &JSONSchemaParser{}
}

// WithOpenAPIDoc sets the OpenAPI document for reference resolution
func (p *JSONSchemaParser) WithOpenAPIDoc(doc *openapi.OpenAPI) *JSONSchemaParser {
	p.openAPIDoc = doc
	return p
}

// Parse loads a schema from a YAML file
func (p *JSONSchemaParser) Parse(filePath string) (*oas3.Schema, error) {
	// Create a new context
	ctx := context.Background()

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file %s: %w", filePath, err)
	}

	// Unmarshal directly to a JSONSchema
	var schema oas3.JSONSchema[oas3.Concrete]
	validationErrs, err := marshaller.Unmarshal(ctx, bytes.NewReader(data), &schema)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema: %w", err)
	}

	// Validate the schema
	additionalErrs := schema.Validate(ctx)
	validationErrs = append(validationErrs, additionalErrs...)

	// If there are validation errors, return them
	if len(validationErrs) > 0 {
		var msgs []string
		for _, ve := range validationErrs {
			msgs = append(msgs, ve.Error())
		}
		return nil, fmt.Errorf("OpenAPI validation errors:\n%s", strings.Join(msgs, "\n"))
	}

	return schema.Left, nil
}

// ExtractSchema processes a single schema with full context
func (p *JSONSchemaParser) ExtractSchema(
	schema *oas3.Schema,
	overrideName *string,
	currentSchemaType string,
) (*SchemaDef, error) {
	return p.extractSchemaRecursive(schema, overrideName, currentSchemaType, "")
}

// extractSchemaRecursive builds a recursive SchemaDef structure
// processAllOfSchemas handles allOf schema composition
func (p *JSONSchemaParser) processAllOfSchemas(
	allOfItem *oas3.JSONSchemaReferenceable,
	result *SchemaDef,
	parentName, currentSchemaType string,
) {
	if allOfItem.GetResolvedObject() == nil || !allOfItem.GetResolvedObject().IsLeft() {
		return
	}

	allOfSchema := allOfItem.GetResolvedObject().GetLeft()

	// Handle references in allOf
	if allOfSchema.Ref != nil && allOfSchema.Ref.String() != "" {
		p.processAllOfReference(allOfSchema, result, parentName, currentSchemaType)
		return
	}

	// Handle direct properties in allOf
	if allOfSchema.Properties != nil {
		result.Required = append(result.Required, allOfSchema.Required...)
		p.extractPropertiesInto(allOfSchema, result, parentName, currentSchemaType)
	}
}

// processAllOfReference resolves and processes a reference in an allOf schema
func (p *JSONSchemaParser) processAllOfReference(
	allOfSchema *oas3.Schema,
	result *SchemaDef,
	parentName, currentSchemaType string,
) {
	refString := allOfSchema.Ref.String()
	if !strings.HasPrefix(refString, "#/components/schemas/") {
		return
	}

	schemaName := strings.TrimPrefix(refString, "#/components/schemas/")
	schemaName = strings.TrimSuffix(schemaName, ".yaml")

	// Check if we can resolve the reference
	if p.openAPIDoc == nil || p.openAPIDoc.Components == nil ||
		p.openAPIDoc.Components.Schemas == nil {
		return
	}

	referencedSchema := p.openAPIDoc.Components.Schemas.GetOrZero(schemaName)
	if referencedSchema == nil {
		return
	}

	resolved := referencedSchema.GetResolvedObject()
	if resolved == nil || !resolved.IsLeft() {
		return
	}

	resolvedSchema := resolved.GetLeft()
	// Merge required fields FIRST, before extracting properties
	result.Required = append(result.Required, resolvedSchema.Required...)
	// Then extract properties from the resolved schema
	p.extractPropertiesInto(resolvedSchema, result, parentName, currentSchemaType)
}

// extractSchemaBasicInfo extracts basic schema information (name, type, format, etc.)
func (p *JSONSchemaParser) extractSchemaBasicInfo(
	schema *oas3.Schema,
	overrideName *string,
) *SchemaDef {
	// Determine name
	name := ""
	if overrideName != nil && *overrideName != "" {
		name = *overrideName
	} else if schema.GetTitle() != "" {
		name = schema.GetTitle()
	}

	result := &SchemaDef{
		Name:        name,
		Description: schema.GetDescription(),
		Schema:      schema,
	}

	// Extract x-codegen extension if at top level
	if name != "" && schema.Extensions != nil {
		result.XCodegen = NewXCodegenParser().Parse(*schema.Extensions)
	}

	// Determine type from schema
	types := schema.GetType()
	if len(types) > 0 {
		result.Type = string(types[0])
	} else if len(schema.AllOf) > 0 {
		// If there's an allOf with no explicit type, assume object
		result.Type = schemaTypeObject
	}

	// Extract format
	if schema.Format != nil {
		result.Format = *schema.Format
	}

	// Extract common properties
	if schema.Nullable != nil {
		result.Nullable = *schema.Nullable
	}
	if schema.Default != nil {
		var defaultVal any
		if err := schema.Default.Decode(&defaultVal); err == nil {
			result.DefaultValue = defaultVal
		}
	}

	return result
}

// processObjectType handles object type schemas
func (p *JSONSchemaParser) processObjectType(
	schema *oas3.Schema,
	result *SchemaDef,
	parentName, currentSchemaType string,
) {
	// Initialize properties map for objects
	result.Properties = make(map[string]*SchemaDef)
	result.Required = schema.Required

	// Process allOf first (even if there are no direct properties)
	for _, allOfItem := range schema.AllOf {
		p.processAllOfSchemas(allOfItem, result, parentName, currentSchemaType)
	}

	// Process direct properties if present
	if schema.Properties != nil {
		p.extractPropertiesInto(schema, result, parentName, currentSchemaType)
	}
}

// processArrayType handles array type schemas
func (p *JSONSchemaParser) processArrayType(
	schema *oas3.Schema,
	result *SchemaDef,
	name, currentSchemaType string,
) {
	if schema.Items == nil || !schema.Items.IsLeft() {
		return
	}

	itemSchema := schema.Items.GetLeft()
	itemName := ""
	if name != "" {
		itemName = name + "Item"
	}
	// For array items, use the itemName as the parentName to avoid duplication
	result.Items, _ = p.extractSchemaRecursive(itemSchema, &itemName, currentSchemaType, itemName)
}

// processEnumType handles enum value extraction for simple types
func processEnumType(schema *oas3.Schema, result *SchemaDef) {
	if schema.Enum == nil {
		return
	}

	for _, enumVal := range schema.Enum {
		if enumVal == nil {
			continue
		}
		var val string
		if err := enumVal.Decode(&val); err == nil {
			result.Enum = append(result.Enum, val)
		}
	}
}

func (p *JSONSchemaParser) extractSchemaRecursive(
	schema *oas3.Schema,
	overrideName *string,
	currentSchemaType string,
	parentName string,
) (*SchemaDef, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema is nil")
	}

	result := p.extractSchemaBasicInfo(schema, overrideName)

	// Handle different types
	switch result.Type {
	case schemaTypeObject:
		p.processObjectType(schema, result, parentName, currentSchemaType)
	case schemaTypeArray:
		p.processArrayType(schema, result, result.Name, currentSchemaType)
	default:
		// Simple types - extract enum values
		processEnumType(schema, result)
	}

	// Compute GoType
	result.GoType = p.computeGoType(result, parentName, currentSchemaType)

	return result, nil
}

// resolvePropertyReference resolves a schema reference and returns the resolved schema and type info
func (p *JSONSchemaParser) resolvePropertyReference(
	schemaName string,
	currentSchemaType string,
) (goType string, schemaType string, format string) {
	if p.openAPIDoc == nil || p.openAPIDoc.Components == nil ||
		p.openAPIDoc.Components.Schemas == nil {
		return schemaName, "", ""
	}

	referencedSchema := p.openAPIDoc.Components.Schemas.GetOrZero(schemaName)
	if referencedSchema == nil {
		return schemaName, "", ""
	}

	resolved := referencedSchema.GetResolvedObject()
	if resolved == nil || !resolved.IsLeft() {
		return schemaName, "", ""
	}

	resolvedSchema := resolved.GetLeft()

	// Get the actual type name from the schema's title
	actualTypeName := resolvedSchema.GetTitle()
	if actualTypeName == "" {
		actualTypeName = schemaName
	}

	// Determine package qualification
	targetPackage := p.extractPropertyPackage(resolvedSchema, actualTypeName)

	// Set the GoType with proper qualification
	qualifiedType := p.qualifyPropertyType(actualTypeName, targetPackage, currentSchemaType)

	// Copy type info from resolved schema
	types := resolvedSchema.GetType()
	if len(types) > 0 {
		schemaType = string(types[0])
	}
	if resolvedSchema.Format != nil {
		format = *resolvedSchema.Format
	}

	return qualifiedType, schemaType, format
}

// extractPropertyPackage extracts the package type for a property schema
func (p *JSONSchemaParser) extractPropertyPackage(schema *oas3.Schema, typeName string) string {
	if schema.Extensions == nil {
		return string(XCodegenExtensionSchemaTypeEntity)
	}

	xExt := schema.Extensions.GetOrZero("x-codegen")
	if xExt == nil {
		return string(XCodegenExtensionSchemaTypeEntity)
	}

	parser := &XCodegenParser{}
	xcodegen, err := parser.ParseExtension(xExt, typeName)
	if err != nil || xcodegen == nil {
		return string(XCodegenExtensionSchemaTypeEntity)
	}

	return string(xcodegen.GetSchemaType())
}

// qualifyPropertyType qualifies a property type name based on context
func (p *JSONSchemaParser) qualifyPropertyType(
	typeName, targetPackage, currentSchemaType string,
) string {
	switch currentSchemaType {
	case "":
		// Controller context - always qualify
		return p.addPropertyPackagePrefix(typeName, targetPackage)
	case targetPackage:
		// Same package
		return typeName
	default:
		// Different packages
		return p.addPropertyPackagePrefix(typeName, targetPackage)
	}
}

// addPropertyPackagePrefix adds package prefix to property type
func (p *JSONSchemaParser) addPropertyPackagePrefix(typeName, targetPackage string) string {
	if targetPackage == string(XCodegenExtensionSchemaTypeValueobject) {
		return "valueobjects." + typeName
	}
	return "entities." + typeName
}

// addOmitEmptyTags adds omitempty to tags if property is not required
func addOmitEmptyTags(def *SchemaDef, isRequired bool) {
	if isRequired {
		return
	}

	if def.JSONTag != "" && !strings.Contains(def.JSONTag, tagOmitEmpty) {
		def.JSONTag += tagOmitEmpty
	}
	if def.YAMLTag != "" && !strings.Contains(def.YAMLTag, tagOmitEmpty) {
		def.YAMLTag += tagOmitEmpty
	}
}

// processPropertyReference handles reference properties
func (p *JSONSchemaParser) processPropertyReference(
	propSchema *oas3.Schema,
	propName, fieldName string,
	parent *SchemaDef,
	currentSchemaType string,
) {
	refString := propSchema.Ref.String()
	if !strings.HasPrefix(refString, "#/components/schemas/") {
		return
	}

	schemaName := strings.TrimPrefix(refString, "#/components/schemas/")
	schemaName = strings.TrimSuffix(schemaName, ".yaml")

	// Create a simple SchemaDef for the reference
	refDef := &SchemaDef{
		Name:    fieldName,
		JSONTag: propName,
		YAMLTag: propName,
		Schema:  propSchema,
	}

	// Resolve the reference to get type info
	goType, schemaType, format := p.resolvePropertyReference(schemaName, currentSchemaType)
	refDef.GoType = goType
	refDef.Type = schemaType
	refDef.Format = format

	// If we couldn't resolve, try to use the title from any available info
	if refDef.GoType == "" {
		if propSchema != nil && propSchema.GetTitle() != "" {
			refDef.GoType = propSchema.GetTitle()
		} else {
			refDef.GoType = schemaName
		}
	}

	// Add omitempty if not required
	addOmitEmptyTags(refDef, parent.IsPropertyRequired(propName))

	parent.Properties[fieldName] = refDef
}

// processDirectProperty handles non-reference properties
func (p *JSONSchemaParser) processDirectProperty(
	propSchema *oas3.Schema,
	propName, fieldName string,
	parent *SchemaDef,
	parentName, currentSchemaType string,
) {
	// For nested objects, always prefix with parent name to avoid conflicts
	typeName := fieldName
	if parentName != "" {
		typeName = parentName + fieldName
	} else if parent.Name != "" {
		typeName = parent.Name + fieldName
	}

	propDef, _ := p.extractSchemaRecursive(propSchema, &fieldName, currentSchemaType, typeName)
	if propDef == nil {
		return
	}

	// Use the original property name for JSON/YAML tags
	propDef.JSONTag = propName
	propDef.YAMLTag = propName

	// Add omitempty if not required
	addOmitEmptyTags(propDef, parent.IsPropertyRequired(propName))

	parent.Properties[fieldName] = propDef

	// If this is a nested object with properties, set its GoType to the prefixed name
	// Objects with only additionalProperties should keep their map type
	// Don't change Name which should remain as the field name
	if propDef.Type == schemaTypeObject && typeName != "" && len(propDef.Properties) > 0 {
		propDef.GoType = typeName // Type name
	}
}

// extractPropertiesInto extracts properties into a parent SchemaDef
func (p *JSONSchemaParser) extractPropertiesInto(
	schema *oas3.Schema,
	parent *SchemaDef,
	parentName string,
	currentSchemaType string,
) {
	if schema.Properties == nil {
		return
	}

	for propName := range schema.Properties.Keys() {
		propRef := schema.Properties.GetOrZero(propName)
		if propRef == nil || !propRef.IsLeft() {
			continue
		}

		propSchema := propRef.GetLeft()
		fieldName := PascalCase(propName)

		// Check if it's a reference
		if propSchema.Ref != nil && propSchema.Ref.String() != "" {
			p.processPropertyReference(propSchema, propName, fieldName, parent, currentSchemaType)
		} else {
			p.processDirectProperty(propSchema, propName, fieldName, parent, parentName, currentSchemaType)
		}
	}
}

// computeGoType determines the Go type for a schema
func (p *JSONSchemaParser) computeGoType(
	schema *SchemaDef,
	parentName string,
	currentSchemaType string,
) string {
	// If already set (e.g., for references), use it
	if schema.GoType != "" {
		return schema.GoType
	}

	// For objects with a name AND properties, use the name as the type
	// Objects with only additionalProperties should be handled as maps
	if schema.Type == "object" && schema.Name != "" && len(schema.Properties) > 0 {
		if parentName != "" && parentName != schema.Name {
			return parentName + schema.Name
		}
		return schema.Name
	}

	// For arrays, prepend [] to item type
	if schema.Type == "array" && schema.Items != nil {
		return "[]" + p.computeGoType(schema.Items, parentName, currentSchemaType)
	}

	// For simple types, use type conversion
	return SchemaToGoType(schema.Schema, p.openAPIDoc, currentSchemaType)
}
