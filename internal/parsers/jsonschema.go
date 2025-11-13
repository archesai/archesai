// Package parsers provides OpenAPI and JSON Schema parsing utilities.
package parsers

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/datamodel/low"
	lowbase "github.com/pb33f/libopenapi/datamodel/low/base"
	"github.com/pb33f/libopenapi/index"
	"go.yaml.in/yaml/v4"
)

// JSONSchemaParser handles parsing JSON Schema files
type JSONSchemaParser struct {
	openAPIDoc *v3.Document
	schemaType XCodegenSchemaType
	schemaName string // Track the current schema name being parsed
}

// NewJSONSchemaParser creates a new JSONSchemaParser instance
func NewJSONSchemaParser(doc *v3.Document) *JSONSchemaParser {
	return &JSONSchemaParser{
		openAPIDoc: doc,
	}
}

// Parse parses a JSON Schema from bytes
func (p *JSONSchemaParser) Parse(data []byte) (*SchemaDef, error) {
	// Parse as YAML node first
	var rootNode yaml.Node
	if err := yaml.Unmarshal(data, &rootNode); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema: %w", err)
	}

	// Create an index for the schema
	idx := index.NewSpecIndex(&rootNode)

	// Build low-level schema proxy
	lowProxy := &lowbase.SchemaProxy{}
	ctx := context.Background()

	// The rootNode is a Document node, we need to get to the actual schema content
	var schemaNode *yaml.Node
	if rootNode.Kind == yaml.DocumentNode && len(rootNode.Content) > 0 {
		schemaNode = rootNode.Content[0]
	} else {
		schemaNode = &rootNode
	}

	if err := lowProxy.Build(ctx, nil, schemaNode, idx); err != nil {
		return nil, fmt.Errorf("failed to build schema proxy: %w", err)
	}

	// Create a NodeReference wrapper for the low-level schema proxy
	nodeRef := &low.NodeReference[*lowbase.SchemaProxy]{
		Value:     lowProxy,
		ValueNode: schemaNode,
		Context:   ctx,
	}

	// Convert to high-level schema
	highProxy := base.NewSchemaProxy(nodeRef)
	schema := highProxy.Schema()
	if schema == nil {
		return nil, fmt.Errorf("failed to build high-level schema")
	}

	// Track the schema name from title
	p.schemaName = schema.Title

	// Determine schema type from x-codegen extension
	p.schemaType = XCodegenSchemaTypeValueobject // Default
	if schema.Extensions != nil {
		if ext, ok := schema.Extensions.Get("x-codegen-schema-type"); ok {
			var schemaType string
			if err := ext.Decode(&schemaType); err == nil {
				p.schemaType = XCodegenSchemaType(schemaType)
			}
		}
	}

	processed, err := p.extractSchemaRecursive(schema, "")
	if err != nil {
		return nil, fmt.Errorf("failed to extract schema definition: %w", err)
	}

	// Set the schema type on the processed schema
	processed.XCodegenSchemaType = p.schemaType

	return processed, nil
}

// ParseFile reads and parses a JSON Schema from a file
func (p *JSONSchemaParser) ParseFile(path string) (*SchemaDef, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file %s: %w", path, err)
	}
	return p.Parse(data)
}

// ParseBase parses a schema from base.Schema
func (p *JSONSchemaParser) ParseBase(schema *base.Schema) (*SchemaDef, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema is nil")
	}

	// Track the schema name from title
	p.schemaName = schema.Title

	// Determine schema type from x-codegen extension
	p.schemaType = XCodegenSchemaTypeValueobject // Default
	if schema.Extensions != nil {
		if ext, ok := schema.Extensions.Get("x-codegen-schema-type"); ok {
			var schemaType string
			if err := ext.Decode(&schemaType); err == nil {
				p.schemaType = XCodegenSchemaType(schemaType)
			}
		}
	}

	processed, err := p.extractSchemaRecursive(schema, "")
	if err != nil {
		return nil, fmt.Errorf("failed to extract schema definition: %w", err)
	}

	// Set the schema type on the processed schema
	processed.XCodegenSchemaType = p.schemaType

	return processed, nil
}

// extractSchemaRecursive builds a recursive SchemaDef structure
// processAllOfSchemas handles allOf schema composition
func (p *JSONSchemaParser) processAllOfSchemas(
	allOfItem *base.SchemaProxy,
	result *SchemaDef,
	parentName string,
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
			allOfItem.GetReference(),
		)
		return
	}

	// Handle direct properties in allOf
	if allOfSchema.Properties != nil {
		result.Required = append(result.Required, allOfSchema.Required...)
		p.extractPropertiesInto(allOfSchema, result, parentName)
	}
}

// processAllOfReference resolves and processes a reference in an allOf schema
func (p *JSONSchemaParser) processAllOfReference(
	result *SchemaDef,
	parentName,
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
	p.extractPropertiesInto(resolvedSchema, result, parentName)
}

// extractSchemaBasicInfo extracts basic schema information (name, type, format, etc.)
func (p *JSONSchemaParser) extractSchemaBasicInfo(
	schema *base.Schema,
) *SchemaDef {

	result := &SchemaDef{
		Name:               schema.Title,
		Description:        schema.Description,
		Schema:             schema,
		Format:             schema.Format,
		DefaultValue:       schema.Default,
		XCodegenSchemaType: p.schemaType,
	}

	// Extract x-codegen extension if at top level
	if schema.Extensions != nil {
		if ext, ok := schema.Extensions.Get("x-codegen"); ok {
			var xcodegen XCodegenExtension
			if err := ext.Decode(&xcodegen); err == nil {
				if xcodegen.Repository != nil && len(xcodegen.Repository.Indices) > 0 {
					sort.Strings(xcodegen.Repository.Indices)
				}
				result.XCodegen = &xcodegen
			}
		}
	}

	// Determine type from schema
	if len(schema.Type) > 0 {
		for _, t := range schema.Type {
			if t == "null" {
				result.Nullable = true
			} else {
				result.Type = t // Use the non-null type
			}
		}
	} else if len(schema.AllOf) > 0 {
		result.Type = schemaTypeObject
	}

	if schema.Nullable != nil {
		result.Nullable = *schema.Nullable
	}

	return result
}

// processObjectType handles object type schemas
func (p *JSONSchemaParser) processObjectType(
	schema *base.Schema,
	result *SchemaDef,
	parentName string,
) {
	// Initialize properties map for objects
	result.Properties = make(map[string]*SchemaDef)
	result.Required = schema.Required

	// Process allOf first (even if there are no direct properties)
	for _, allOfItem := range schema.AllOf {
		p.processAllOfSchemas(allOfItem, result, parentName)
	}

	// Process direct properties if present
	if schema.Properties != nil {
		p.extractPropertiesInto(schema, result, parentName)
	}
}

// processArrayItemReference resolves array item references and returns a SchemaDef
func (p *JSONSchemaParser) processArrayItemReference(
	refString string,
) *SchemaDef {
	if !strings.HasPrefix(refString, "#/components/schemas/") {
		return nil
	}

	schemaName := strings.TrimPrefix(refString, "#/components/schemas/")
	schemaName = strings.TrimSuffix(schemaName, ".yaml")

	// Resolve the reference to get type info
	goType, schemaType, format := p.resolvePropertyReference(schemaName)

	return &SchemaDef{
		Name:   schemaName,
		GoType: goType,
		Type:   schemaType,
		Format: format,
	}
}

// processArrayType handles array type schemas
func (p *JSONSchemaParser) processArrayType(
	schema *base.Schema,
	result *SchemaDef,
	name string,
) {
	if schema.Items == nil || schema.Items.A == nil {
		return
	}

	itemProxy := schema.Items.A
	itemSchema := itemProxy.Schema()
	if itemSchema == nil {
		return
	}

	// Check if the item is a reference to another schema
	if refString := itemProxy.GetReference(); refString != "" {
		// Handle reference - extract the referenced type name
		result.Items = p.processArrayItemReference(refString)
		return
	}

	// For non-reference items (inline schemas), create a nested type name
	itemName := ""
	if name != "" {
		itemName = name + "Item"
	}
	// For array items, use the itemName as the parentName to avoid duplication
	result.Items, _ = p.extractSchemaRecursive(itemSchema, itemName)

	// If the item is an object with properties, ensure it has a name for type generation
	if result.Items != nil && result.Items.Type == schemaTypeObject &&
		len(result.Items.Properties) > 0 &&
		itemName != "" {
		result.Items.Name = itemName
		result.Items.GoType = itemName
	}
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
	parentName string,
) (*SchemaDef, error) {

	if schema == nil {
		return nil, fmt.Errorf("schema is nil")
	}

	result := p.extractSchemaBasicInfo(schema)

	// Handle different types
	switch result.Type {
	case schemaTypeObject:
		p.processObjectType(schema, result, parentName)
	case schemaTypeArray:
		// Use parentName if Name is empty (for nested properties)
		arrayName := result.Name
		if arrayName == "" && parentName != "" {
			arrayName = parentName
		}
		p.processArrayType(schema, result, arrayName)
	default:
		// Simple types - extract enum values
		processEnumType(schema, result)
	}

	// Compute GoType
	result.GoType = p.computeGoType(result, parentName)

	return result, nil
}

// resolvePropertyReference resolves a schema reference and returns the resolved schema and type info
func (p *JSONSchemaParser) resolvePropertyReference(
	schemaName string,
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

	// Set the GoType with proper qualification
	// Determine the target package type from the referenced schema
	targetPackage := string(XCodegenSchemaTypeEntity) // Default to entity
	if resolvedSchema.Extensions != nil {
		if ext, ok := resolvedSchema.Extensions.Get("x-codegen-schema-type"); ok {
			var schemaTypeStr string
			if err := ext.Decode(&schemaTypeStr); err == nil {
				targetPackage = schemaTypeStr
			}
		}
	}
	qualifiedType := p.qualifyPropertyType(actualTypeName, targetPackage)

	// Copy type info from resolved schema
	if len(resolvedSchema.Type) > 0 {
		schemaType = resolvedSchema.Type[0]
	}
	if resolvedSchema.Format != "" {
		format = resolvedSchema.Format
	}

	return qualifiedType, schemaType, format
}

// qualifyPropertyType qualifies a property type name based on context
func (p *JSONSchemaParser) qualifyPropertyType(
	typeName, targetPackage string,
) string {
	// For response schemas (which default to valueobject but aren't really valueobjects),
	// we need to check if we're in a different context.
	// If the current schema is a valueobject but the name suggests it's a response schema,
	// we should always qualify references to actual valueobjects and entities.
	if p.schemaName != "" && strings.HasSuffix(p.schemaName, "Response") {
		// This is a response schema, always qualify model types
		if targetPackage == string(XCodegenSchemaTypeValueobject) ||
			targetPackage == string(XCodegenSchemaTypeEntity) {
			return "models." + typeName
		}
		return "models." + typeName
	}

	// If we're not parsing an entity or valueobject (e.g., we're in a controller),
	// always qualify entity and valueobject types
	if p.schemaType != XCodegenSchemaTypeEntity && p.schemaType != XCodegenSchemaTypeValueobject {
		if targetPackage == string(XCodegenSchemaTypeValueobject) {
			return "models." + typeName
		}
		return "models." + typeName
	}

	// Within entity/valueobject packages, only qualify if different package
	switch string(p.schemaType) {
	case targetPackage:
		return typeName
	default:
		if targetPackage == string(XCodegenSchemaTypeValueobject) {
			return "models." + typeName
		}
		return "models." + typeName
	}
}

// addOmitEmptyTags adds omitempty to tags if property is not required
func addOmitEmptyTags(def *SchemaDef, isRequired bool) {

	tagOmitEmpty := ",omitempty"

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
	goType, schemaType, format := p.resolvePropertyReference(schemaName)
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
	parentName string,
) {
	// For nested objects, always prefix with parent name to avoid conflicts
	typeName := fieldName
	if parentName != "" {
		typeName = parentName + fieldName
	} else if parent.Name != "" {
		typeName = parent.Name + fieldName
	}

	propDef, _ := p.extractSchemaRecursive(propSchema, typeName)
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
				propProxy.GetReference(),
			)
		} else {
			p.processDirectProperty(propSchema, propName, fieldName, parent, parentName)
		}
	}
}

// computeGoType determines the Go type for a schema
func (p *JSONSchemaParser) computeGoType(
	schema *SchemaDef,
	parentName string,
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
		// If the item is an object with properties and has a name, use that name
		// This ensures we use the generated nested type (e.g., LevelsItem) instead of map[string]any
		if schema.Items.Type == "object" && schema.Items.Name != "" &&
			len(schema.Items.Properties) > 0 {
			return "[]" + schema.Items.Name
		}
		// Otherwise compute the item type recursively
		return "[]" + p.computeGoType(schema.Items, parentName)
	}

	// For simple types, use type conversion
	return SchemaToGoType(schema.Schema, p.openAPIDoc, string(p.schemaType))
}
