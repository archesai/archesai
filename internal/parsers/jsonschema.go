// Package parsers provides OpenAPI and JSON Schema parsing functionality
package parsers

import (
	"fmt"
	"os"
	"strings"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

const (
	tagOmitEmpty = ",omitempty"
)

// JSONSchemaParser handles parsing JSON Schema files
type JSONSchemaParser struct {
	openAPIDoc *v3.Document
}

// NewJSONSchemaParser creates a new JSONSchemaParser instance
func NewJSONSchemaParser() *JSONSchemaParser {
	return &JSONSchemaParser{}
}

// WithOpenAPIDoc sets the OpenAPI document for reference resolution
func (p *JSONSchemaParser) WithOpenAPIDoc(doc *v3.Document) *JSONSchemaParser {
	p.openAPIDoc = doc
	return p
}

// Parse loads a schema from a YAML file
func (p *JSONSchemaParser) Parse(path string) (*base.Schema, error) {
	// Read the file
	specBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file %s: %w", path, err)
	}

	config := &datamodel.DocumentConfiguration{
		AllowFileReferences:   true,
		AllowRemoteReferences: true,
	}

	doc, err := libopenapi.NewDocumentWithConfiguration(specBytes, config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}

	v3Model, err := doc.BuildV3Model()
	if err != nil {
		return nil, fmt.Errorf("failed to build v3 model: %w", err)
	}

	// For standalone schema files, we'd need to extract the schema differently
	// This is a placeholder - may need adjustment based on actual use case
	if v3Model == nil || v3Model.Model.Components == nil ||
		v3Model.Model.Components.Schemas == nil {
		return nil, fmt.Errorf("no schemas found in file")
	}

	// Store the parsed document
	p.openAPIDoc = &v3Model.Model

	// Return the first schema found
	for pair := v3Model.Model.Components.Schemas.First(); pair != nil; pair = pair.Next() {
		return pair.Value().Schema(), nil
	}

	return nil, fmt.Errorf("no schema found")
}

// ExtractSchema processes a single schema with full context
func (p *JSONSchemaParser) ExtractSchema(
	schema *base.Schema,
	overrideName *string,
	currentSchemaType string,
) (*SchemaDef, error) {
	return p.extractSchemaRecursive(schema, overrideName, currentSchemaType, "")
}

// extractSchemaRecursive builds a recursive SchemaDef structure
// processAllOfSchemas handles allOf schema composition
func (p *JSONSchemaParser) processAllOfSchemas(
	allOfItem *base.SchemaProxy,
	result *SchemaDef,
	parentName, currentSchemaType string,
) {
	allOfSchema := allOfItem.Schema()
	if allOfSchema == nil {
		return
	}

	// Handle references in allOf
	if allOfItem.GetReference() != "" {
		p.processAllOfReference(
			result,
			parentName,
			currentSchemaType,
			allOfItem.GetReference(),
		)
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
	result *SchemaDef,
	parentName, currentSchemaType string,
	refString string,
) {
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

	referencedSchemaProxy, ok := p.openAPIDoc.Components.Schemas.Get(schemaName)
	if !ok {
		return
	}

	resolvedSchema := referencedSchemaProxy.Schema()
	if resolvedSchema == nil {
		return
	}

	// Merge required fields FIRST, before extracting properties
	result.Required = append(result.Required, resolvedSchema.Required...)
	// Then extract properties from the resolved schema
	p.extractPropertiesInto(resolvedSchema, result, parentName, currentSchemaType)
}

// extractSchemaBasicInfo extracts basic schema information (name, type, format, etc.)
func (p *JSONSchemaParser) extractSchemaBasicInfo(
	schema *base.Schema,
	overrideName *string,
) *SchemaDef {
	// Determine name
	name := ""
	if overrideName != nil && *overrideName != "" {
		name = *overrideName
	} else if schema.Title != "" {
		name = schema.Title
	}

	result := &SchemaDef{
		Name:        name,
		Description: schema.Description,
		Schema:      schema,
	}

	// Extract x-codegen extension if at top level
	if name != "" && schema.Extensions != nil {
		if ext, ok := schema.Extensions.Get("x-codegen"); ok {
			parser := &XCodegenParser{}
			if xcodegen, err := parser.ParseExtension(ext, name); err == nil {
				result.XCodegen = xcodegen
			}
		}
	}

	// Determine type from schema
	if len(schema.Type) > 0 {
		result.Type = schema.Type[0]
	} else if len(schema.AllOf) > 0 {
		// If there's an allOf with no explicit type, assume object
		result.Type = schemaTypeObject
	}

	// Extract format
	if schema.Format != "" {
		result.Format = schema.Format
	}

	// Extract common properties
	if schema.Nullable != nil {
		result.Nullable = *schema.Nullable
	}
	if schema.Default != nil {
		result.DefaultValue = schema.Default
	}

	return result
}

// processObjectType handles object type schemas
func (p *JSONSchemaParser) processObjectType(
	schema *base.Schema,
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
	schema *base.Schema,
	result *SchemaDef,
	name, currentSchemaType string,
) {
	if schema.Items == nil || schema.Items.A == nil {
		return
	}

	itemSchema := schema.Items.A.Schema()
	if itemSchema == nil {
		return
	}

	itemName := ""
	if name != "" {
		itemName = name + "Item"
	}
	// For array items, use the itemName as the parentName to avoid duplication
	result.Items, _ = p.extractSchemaRecursive(itemSchema, &itemName, currentSchemaType, itemName)
}

// processEnumType handles enum value extraction for simple types
func processEnumType(schema *base.Schema, result *SchemaDef) {
	if schema.Enum == nil {
		return
	}

	for _, enumVal := range schema.Enum {
		if enumVal == nil {
			continue
		}
		// Decode yaml.Node to string
		var strVal string
		if err := enumVal.Decode(&strVal); err == nil {
			result.Enum = append(result.Enum, strVal)
		}
	}
}

func (p *JSONSchemaParser) extractSchemaRecursive(
	schema *base.Schema,
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

	referencedSchemaProxy, ok := p.openAPIDoc.Components.Schemas.Get(schemaName)
	if !ok {
		return schemaName, "", ""
	}

	resolvedSchema := referencedSchemaProxy.Schema()
	if resolvedSchema == nil {
		return schemaName, "", ""
	}

	// Get the actual type name from the schema's title
	actualTypeName := resolvedSchema.Title
	if actualTypeName == "" {
		actualTypeName = schemaName
	}

	// Determine package qualification
	targetPackage := p.extractPropertyPackage(resolvedSchema, actualTypeName)

	// Set the GoType with proper qualification
	qualifiedType := p.qualifyPropertyType(actualTypeName, targetPackage, currentSchemaType)

	// Copy type info from resolved schema
	if len(resolvedSchema.Type) > 0 {
		schemaType = resolvedSchema.Type[0]
	}
	if resolvedSchema.Format != "" {
		format = resolvedSchema.Format
	}

	return qualifiedType, schemaType, format
}

// extractPropertyPackage extracts the package type for a property schema
func (p *JSONSchemaParser) extractPropertyPackage(schema *base.Schema, typeName string) string {
	if schema.Extensions == nil {
		return string(XCodegenExtensionSchemaTypeEntity)
	}

	xExt, ok := schema.Extensions.Get("x-codegen")
	if !ok {
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
	propSchema *base.Schema,
	propName, fieldName string,
	parent *SchemaDef,
	currentSchemaType string,
	refString string,
) {
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
		if propSchema != nil && propSchema.Title != "" {
			refDef.GoType = propSchema.Title
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
	propSchema *base.Schema,
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
	schema *base.Schema,
	parent *SchemaDef,
	parentName string,
	currentSchemaType string,
) {
	if schema.Properties == nil {
		return
	}

	for propPair := schema.Properties.First(); propPair != nil; propPair = propPair.Next() {
		propName := propPair.Key()
		propProxy := propPair.Value()

		if propProxy == nil {
			continue
		}

		propSchema := propProxy.Schema()
		if propSchema == nil {
			continue
		}

		fieldName := PascalCase(propName)

		// Check if it's a reference
		if propProxy.GetReference() != "" {
			p.processPropertyReference(
				propSchema,
				propName,
				fieldName,
				parent,
				currentSchemaType,
				propProxy.GetReference(),
			)
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
