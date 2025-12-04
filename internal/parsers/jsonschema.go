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

const (
	omitEmptyTag = ",omitempty"
)

// JSONSchemaParser handles parsing JSON Schema files
type JSONSchemaParser struct {
	openAPIDoc *v3.Document
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

	return p.ParseBase(schema)
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

	// Extract schema type from extension
	schemaType := XCodegenSchemaTypeValueobject
	if schema.Extensions != nil {
		if ext, ok := schema.Extensions.Get("x-codegen-schema-type"); ok {
			var schemaTypeStr string
			if err := ext.Decode(&schemaTypeStr); err == nil {
				schemaType = XCodegenSchemaType(schemaTypeStr)
			}
		}
	}

	// Track the schema name from title
	schemaName := schema.Title

	processed, err := p.extractSchemaRecursive(schema, "", schemaName, schemaType)
	if err != nil {
		return nil, fmt.Errorf("failed to extract schema definition: %w", err)
	}

	// Set the schema type on the processed schema
	processed.XCodegenSchemaType = schemaType

	return processed, nil
}

// resolveReferencedSchema resolves a schema reference by name
func (p *JSONSchemaParser) resolveReferencedSchema(schemaName string) (*base.Schema, error) {
	if p.openAPIDoc == nil || p.openAPIDoc.Components == nil ||
		p.openAPIDoc.Components.Schemas == nil {
		return nil, fmt.Errorf("openAPI document or component schemas not available")
	}

	referencedSchemaProxy, ok := p.openAPIDoc.Components.Schemas.Get(schemaName)
	if !ok {
		return nil, fmt.Errorf("schema %s not found", schemaName)
	}

	resolvedSchema := referencedSchemaProxy.Schema()
	if resolvedSchema == nil {
		return nil, fmt.Errorf("resolved schema %s is nil", schemaName)
	}

	return resolvedSchema, nil
}

// extractSchemaRecursive builds a recursive SchemaDef structure
func (p *JSONSchemaParser) extractSchemaRecursive(
	schema *base.Schema,
	parentName string,
	schemaName string,
	schemaType XCodegenSchemaType,
) (*SchemaDef, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema is nil")
	}

	result := p.buildBasicSchemaDef(schema, schemaType)
	p.extractSchemaType(schema, result)

	// Handle different types
	switch result.Type {
	case schemaTypeObject:
		p.processObjectSchema(schema, result, parentName, schemaName, schemaType)
	case schemaTypeArray:
		p.processArraySchema(schema, result, parentName, schemaName, schemaType)
	default:
		p.processSimpleSchema(schema, result)
	}

	// Compute GoType
	result.GoType = p.computeGoType(result, parentName, schemaName, schemaType)

	return result, nil
}

// buildBasicSchemaDef creates the basic SchemaDef structure
func (p *JSONSchemaParser) buildBasicSchemaDef(
	schema *base.Schema,
	schemaType XCodegenSchemaType,
) *SchemaDef {
	result := &SchemaDef{
		Name:               schema.Title,
		Description:        schema.Description,
		Schema:             schema,
		Format:             schema.Format,
		DefaultValue:       schema.Default,
		XCodegenSchemaType: schemaType,
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
		// Extract x-internal extension - marks schema as imported from another package
		if ext, ok := schema.Extensions.Get("x-internal"); ok {
			var xinternal string
			if err := ext.Decode(&xinternal); err == nil {
				result.XInternal = xinternal
			}
		}
	}

	return result
}

// extractSchemaType determines the type from schema
func (p *JSONSchemaParser) extractSchemaType(schema *base.Schema, result *SchemaDef) {
	if len(schema.Type) > 0 {
		for _, t := range schema.Type {
			if t == "null" {
				result.Nullable = true
			} else {
				result.Type = t
			}
		}
	} else if len(schema.AllOf) > 0 {
		result.Type = schemaTypeObject
	}

	if schema.Nullable != nil {
		result.Nullable = *schema.Nullable
	}
}

// processObjectSchema processes object type schemas
func (p *JSONSchemaParser) processObjectSchema(
	schema *base.Schema,
	result *SchemaDef,
	parentName string,
	schemaName string,
	schemaType XCodegenSchemaType,
) {
	result.Properties = make(map[string]*SchemaDef)
	result.Required = schema.Required

	// Process allOf schemas
	p.processAllOfSchemas(schema, result, parentName, schemaName, schemaType)

	// Process direct properties if present
	if schema.Properties != nil {
		p.extractPropertiesInto(schema, result, parentName, schemaName, schemaType)
	}
}

// processAllOfSchemas processes allOf schemas
func (p *JSONSchemaParser) processAllOfSchemas(
	schema *base.Schema,
	result *SchemaDef,
	parentName string,
	schemaName string,
	schemaType XCodegenSchemaType,
) {
	for _, allOfItem := range schema.AllOf {
		allOfSchema := allOfItem.Schema()
		if allOfSchema == nil {
			continue
		}

		if refString := allOfItem.GetReference(); refString != "" {
			p.processAllOfReference(refString, result, parentName, schemaName, schemaType)
		} else if allOfSchema.Properties != nil {
			result.Required = append(result.Required, allOfSchema.Required...)
			p.extractPropertiesInto(allOfSchema, result, parentName, schemaName, schemaType)
		}
	}
}

// processAllOfReference processes a reference in allOf
func (p *JSONSchemaParser) processAllOfReference(
	refString string,
	result *SchemaDef,
	parentName string,
	schemaName string,
	schemaType XCodegenSchemaType,
) {
	if !strings.HasPrefix(refString, "#/components/schemas/") {
		return
	}

	refSchemaName := strings.TrimPrefix(refString, "#/components/schemas/")
	refSchemaName = strings.TrimSuffix(refSchemaName, ".yaml")

	resolvedSchema, err := p.resolveReferencedSchema(refSchemaName)
	if err == nil {
		result.Required = append(result.Required, resolvedSchema.Required...)
		p.extractPropertiesInto(resolvedSchema, result, parentName, schemaName, schemaType)
	}
}

// processArraySchema processes array type schemas
func (p *JSONSchemaParser) processArraySchema(
	schema *base.Schema,
	result *SchemaDef,
	parentName string,
	schemaName string,
	schemaType XCodegenSchemaType,
) {
	if schema.Items == nil || schema.Items.A == nil {
		return
	}

	itemProxy := schema.Items.A
	itemSchema := itemProxy.Schema()
	if itemSchema == nil {
		return
	}

	if refString := itemProxy.GetReference(); refString != "" {
		p.processArrayItemReference(refString, result, schemaName, schemaType)
	} else {
		p.processArrayItemInline(itemSchema, result, parentName, schemaName, schemaType)
	}
}

// processArrayItemReference processes a reference in array items
func (p *JSONSchemaParser) processArrayItemReference(
	refString string,
	result *SchemaDef,
	schemaName string,
	schemaType XCodegenSchemaType,
) {
	if !strings.HasPrefix(refString, "#/components/schemas/") {
		return
	}

	refSchemaName := strings.TrimPrefix(refString, "#/components/schemas/")
	refSchemaName = strings.TrimSuffix(refSchemaName, ".yaml")

	goType, itemSchemaType, format := p.resolvePropertyType(refSchemaName, schemaName, schemaType)
	result.Items = &SchemaDef{
		Name:   refSchemaName,
		GoType: goType,
		Type:   itemSchemaType,
		Format: format,
	}
}

// processArrayItemInline processes inline array item schemas
func (p *JSONSchemaParser) processArrayItemInline(
	itemSchema *base.Schema,
	result *SchemaDef,
	parentName string,
	schemaName string,
	schemaType XCodegenSchemaType,
) {
	arrayName := result.Name
	if arrayName == "" && parentName != "" {
		arrayName = parentName
	}

	itemName := ""
	if arrayName != "" {
		itemName = arrayName + "Item"
	}

	result.Items, _ = p.extractSchemaRecursive(itemSchema, itemName, schemaName, schemaType)

	if result.Items != nil && result.Items.Type == schemaTypeObject &&
		len(result.Items.Properties) > 0 && itemName != "" {
		result.Items.Name = itemName
		result.Items.GoType = itemName
	}
}

// processSimpleSchema processes simple type schemas
func (p *JSONSchemaParser) processSimpleSchema(schema *base.Schema, result *SchemaDef) {
	if schema.Enum == nil {
		return
	}

	for _, enumVal := range schema.Enum {
		if enumVal == nil {
			continue
		}
		var strVal string
		if err := enumVal.Decode(&strVal); err == nil {
			result.Enum = append(result.Enum, strVal)
		}
	}
}

// resolvePropertyType resolves a schema reference and returns type info
func (p *JSONSchemaParser) resolvePropertyType(
	refSchemaName string,
	currentSchemaName string,
	currentSchemaType XCodegenSchemaType,
) (goType string, schemaType string, format string) {
	resolvedSchema, err := p.resolveReferencedSchema(refSchemaName)
	if err != nil {
		return refSchemaName, "", ""
	}

	// Get the actual type name from the schema's title
	actualTypeName := resolvedSchema.Title
	if actualTypeName == "" {
		actualTypeName = refSchemaName
	}

	// Determine the target package type from the referenced schema
	targetPackage := string(XCodegenSchemaTypeEntity) // Default to entity
	xInternal := ""
	if resolvedSchema.Extensions != nil {
		if ext, ok := resolvedSchema.Extensions.Get("x-codegen-schema-type"); ok {
			var schemaTypeStr string
			if err := ext.Decode(&schemaTypeStr); err == nil {
				targetPackage = schemaTypeStr
			}
		}
		// Check for x-internal extension - marks schema as imported from another package
		if ext, ok := resolvedSchema.Extensions.Get("x-internal"); ok {
			var xinternalStr string
			if err := ext.Decode(&xinternalStr); err == nil {
				xInternal = xinternalStr
			}
		}
	}

	// Qualify the type if needed (inline qualifyPropertyType)
	// If schema is internal to another package (e.g., server), use that package's models
	if xInternal == "server" {
		goType = "servermodels." + actualTypeName
	} else {
		// Response schemas always qualify model types
		isResponseSchema := currentSchemaName != "" && strings.HasSuffix(currentSchemaName, "Response")
		if isResponseSchema {
			goType = "models." + actualTypeName
		} else if string(currentSchemaType) == targetPackage {
			// Within the same package, no qualification needed
			goType = actualTypeName
		} else {
			// Different package, always qualify
			goType = "models." + actualTypeName
		}
	}

	// Copy type info from resolved schema
	if len(resolvedSchema.Type) > 0 {
		schemaType = resolvedSchema.Type[0]
	}
	if resolvedSchema.Format != "" {
		format = resolvedSchema.Format
	}

	return goType, schemaType, format
}

// extractPropertiesInto extracts properties into a parent SchemaDef
func (p *JSONSchemaParser) extractPropertiesInto(
	schema *base.Schema,
	parent *SchemaDef,
	parentName string,
	currentSchemaName string,
	currentSchemaType XCodegenSchemaType,
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
		if refString := propProxy.GetReference(); refString != "" {
			// Handle reference property (inline processPropertyReference)
			if strings.HasPrefix(refString, "#/components/schemas/") {
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
				goType, schemaType, format := p.resolvePropertyType(
					schemaName,
					currentSchemaName,
					currentSchemaType,
				)
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

				// Add omitempty if not required (inline addOmitEmptyTags)
				if !parent.IsPropertyRequired(propName) {
					if refDef.JSONTag != "" && !strings.Contains(refDef.JSONTag, omitEmptyTag) {
						refDef.JSONTag += omitEmptyTag
					}
					if refDef.YAMLTag != "" && !strings.Contains(refDef.YAMLTag, omitEmptyTag) {
						refDef.YAMLTag += omitEmptyTag
					}
				}

				parent.Properties[fieldName] = refDef
			}
		} else {
			// Handle direct property (inline processDirectProperty)
			// For nested objects, generate consistent type names
			effectiveParentName := parentName
			if effectiveParentName == "" {
				effectiveParentName = parent.Name
			}

			// Inline generateNestedTypeName
			typeName := ""
			if effectiveParentName != "" && fieldName != "" {
				typeName = effectiveParentName + fieldName
			} else if fieldName != "" {
				typeName = fieldName
			} else {
				typeName = effectiveParentName
			}

			propDef, _ := p.extractSchemaRecursive(propSchema, typeName, currentSchemaName, currentSchemaType)
			if propDef == nil {
				continue
			}

			// Use the original property name for JSON/YAML tags
			propDef.JSONTag = propName
			propDef.YAMLTag = propName

			// Add omitempty if not required (inline addOmitEmptyTags)
			if !parent.IsPropertyRequired(propName) {
				if propDef.JSONTag != "" && !strings.Contains(propDef.JSONTag, omitEmptyTag) {
					propDef.JSONTag += omitEmptyTag
				}
				if propDef.YAMLTag != "" && !strings.Contains(propDef.YAMLTag, omitEmptyTag) {
					propDef.YAMLTag += omitEmptyTag
				}
			}

			parent.Properties[fieldName] = propDef

			// If this is a nested object with properties, set its GoType to the prefixed name
			// Objects with only additionalProperties should keep their map type
			// Don't change Name which should remain as the field name
			if propDef.Type == schemaTypeObject && typeName != "" && len(propDef.Properties) > 0 {
				propDef.GoType = typeName // Type name
			}
		}
	}
}

// computeGoType determines the Go type for a schema
func (p *JSONSchemaParser) computeGoType(
	schema *SchemaDef,
	parentName string,
	_ string,
	currentSchemaType XCodegenSchemaType,
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
		return "[]" + p.computeGoType(
			schema.Items,
			parentName,
			"",
			currentSchemaType,
		)
	}

	// For simple types, use type conversion
	return SchemaToGoType(schema.Schema, p.openAPIDoc, string(currentSchemaType))
}
