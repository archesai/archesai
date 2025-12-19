package schema

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.yaml.in/yaml/v4"

	"github.com/archesai/archesai/internal/ref"
	"github.com/archesai/archesai/internal/strutil"
)

// Loader loads and resolves schemas from YAML files.
// It handles both file-based $refs and inline schemas.
type Loader struct {
	baseDir   string
	schemas   map[string]*Schema
	refToName map[string]string
}

// NewLoader creates a Loader for standalone schema files.
func NewLoader(baseDir string) *Loader {
	return &Loader{
		baseDir:   baseDir,
		schemas:   make(map[string]*Schema),
		refToName: make(map[string]string),
	}
}

// Load implements the ref.Loader interface.
// Parses raw bytes into a Schema with computed fields (GoType, JSONTag, etc.).
func (l *Loader) Load(data []byte, name string) (*Schema, error) {
	return l.LoadSchemaFromBytes(data, name)
}

// LoadSchemaFile loads and fully resolves a schema from a YAML file path.
// This includes: parsing YAML, resolving $refs, computing GoType/JSONTag/YAMLTag,
// processing x-codegen extensions, and adding base fields for entities.
func (l *Loader) LoadSchemaFile(path string) (*Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	// Update baseDir to the schema file's directory for resolving relative refs
	oldBaseDir := l.baseDir
	l.baseDir = filepath.Dir(path)
	defer func() { l.baseDir = oldBaseDir }()

	return l.LoadSchemaFromBytes(data, filepath.Base(path))
}

// LoadSchemaFromBytes parses and resolves a schema from YAML bytes.
// defaultName is used if the schema has no title (typically the filename without extension).
func (l *Loader) LoadSchemaFromBytes(data []byte, defaultName string) (*Schema, error) {
	var sf Schema
	if err := yaml.Unmarshal(data, &sf); err != nil {
		return nil, fmt.Errorf("parsing schema: %w", err)
	}

	// Derive name from title or default (strip extension if present)
	name := sf.Title
	if name == "" {
		name = strings.TrimSuffix(defaultName, filepath.Ext(defaultName))
	}

	schema := &Schema{
		Title:              name,
		Description:        sf.Description,
		Type:               PropertyType{Types: []string{TypeObject}},
		Properties:         make(map[string]*ref.Ref[Schema]),
		Required:           sf.Required,
		XCodegen:           sf.XCodegen,
		XCodegenSchemaType: sf.XCodegenSchemaType,
		XInternal:          sf.XInternal,
	}

	// Handle allOf composition (common pattern: allOf: [$ref: Base.yaml, {type: object, properties: ...}])
	if len(sf.AllOf) > 0 {
		for _, itemRef := range sf.AllOf {
			if itemRef.IsRef() {
				// This is a reference to a base schema (e.g., Base.yaml)
				// The base fields (id, createdAt, updatedAt) will be added later
				continue
			}
			// This is the inline object definition with properties
			item := itemRef.GetOrNil()
			if item == nil {
				continue
			}
			for propName, propRef := range item.Properties {
				fieldName := strutil.PascalCase(propName)
				propSchema := propRef.GetOrNil()
				if propSchema != nil {
					l.ProcessProperty(propSchema, fieldName, propName, name)
				}
				schema.Properties[fieldName] = propRef
			}
			// Merge required fields
			schema.Required = append(schema.Required, item.Required...)
		}
	} else {
		// Regular schema without allOf
		for propName, propRef := range sf.Properties {
			fieldName := strutil.PascalCase(propName)

			// Handle ref properties by loading them
			if propRef.IsRef() {
				refPath := propRef.RefPath
				refName := ExtractSchemaNameFromRef(refPath)

				// Try to find already loaded schema
				refSchema := l.resolveSchemaByRef(refName)
				if refSchema == nil {
					// Load from file
					loadedSchema, err := l.LoadRef(refPath)
					if err == nil {
						refSchema = loadedSchema
					}
				}

				if refSchema != nil {
					// Create inline schema that references the resolved schema
					propSchema := &Schema{
						Title:   fieldName,
						GoType:  refSchema.Title,
						Type:    refSchema.Type,
						JSONTag: propName,
						YAMLTag: propName,
					}
					schema.Properties[fieldName] = ref.NewInline(propSchema)
					continue
				}
			}

			// Handle inline properties
			propSchema := propRef.GetOrNil()
			if propSchema != nil {
				l.ProcessProperty(propSchema, fieldName, propName, name)
			}
			schema.Properties[fieldName] = propRef
		}
	}

	// Build required map and set omitempty on optional fields
	requiredMap := make(map[string]bool)
	for _, req := range schema.Required {
		requiredMap[req] = true
	}
	for propName, propRef := range schema.Properties {
		propSchema := propRef.GetOrNil()
		if propSchema == nil {
			continue
		}
		jsonName := propSchema.JSONTag
		if jsonName == "" {
			jsonName = strutil.CamelCase(propName)
		}
		if !requiredMap[jsonName] && !requiredMap[propName] {
			propSchema.JSONTag = jsonName + ",omitempty"
			propSchema.YAMLTag = jsonName + ",omitempty"
		}
	}

	// Add base fields for entity schemas
	if sf.XCodegenSchemaType == TypeEntity {
		AddBaseFields(schema)
	}

	// Set GoType with models alias prefix if this schema comes from an include
	if schema.XInternal != "" {
		alias := GetSchemasAlias(schema.XInternal)
		if alias != "" {
			schema.GoType = alias + "." + schema.Title
		} else {
			schema.GoType = schema.Title
		}
	} else {
		schema.GoType = schema.Title
	}

	// Cache the schema
	l.schemas[schema.Title] = schema

	return schema, nil
}

// ProcessProperty sets computed fields on a property schema.
// Sets: Title, JSONTag, YAMLTag, GoType, Nullable
// Handles: oneOf nullable patterns, nested objects, arrays
func (l *Loader) ProcessProperty(propSchema *Schema, fieldName, jsonName, parentName string) {
	propSchema.Title = fieldName
	propSchema.JSONTag = jsonName
	propSchema.YAMLTag = jsonName
	propSchema.Nullable = propSchema.Type.Nullable

	// Handle oneOf for nullable types with constraints
	// e.g., oneOf: [{type: string, minLength: 1}, {type: 'null'}]
	if len(propSchema.OneOf) == 2 {
		var nonNullSchema *Schema
		hasNull := false
		for _, r := range propSchema.OneOf {
			s := r.GetOrNil()
			if s == nil {
				continue
			}
			pt := s.Type.PrimaryType()
			if pt == TypeNull || pt == "" && s.Type.Nullable {
				hasNull = true
			} else if nonNullSchema == nil {
				nonNullSchema = s
			}
		}
		// If it's a nullable union type, use the non-null schema's type
		if hasNull && nonNullSchema != nil {
			propSchema.Nullable = true
			propSchema.Type = nonNullSchema.Type
			propSchema.Format = nonNullSchema.Format
			propSchema.MinLength = nonNullSchema.MinLength
			propSchema.MaxLength = nonNullSchema.MaxLength
			propSchema.Pattern = nonNullSchema.Pattern
			propSchema.Minimum = nonNullSchema.Minimum
			propSchema.Maximum = nonNullSchema.Maximum
			propSchema.Enum = nonNullSchema.Enum
		}
	}

	// Determine Go type based on the schema type
	primaryType := propSchema.Type.PrimaryType()
	switch primaryType {
	case TypeArray:
		if propSchema.Items == nil {
			propSchema.GoType = "[]any"
		} else if propSchema.Items.IsRef() {
			// Array items are a $ref - use the referenced schema name
			itemName := ExtractSchemaNameFromRef(propSchema.Items.RefPath)
			propSchema.GoType = "[]" + itemName
		} else {
			// Array items are inline
			itemSchema := propSchema.Items.GetOrNil()
			if itemSchema == nil {
				propSchema.GoType = "[]any"
			} else if itemSchema.Type.PrimaryType() == TypeObject && len(itemSchema.Properties) > 0 {
				// Inline object items - generate a type name
				itemTypeName := parentName + fieldName + "Item"
				itemSchema.GoType = itemTypeName
				// Process nested properties
				processedProps := make(map[string]*ref.Ref[Schema])
				for nestedPropName, nestedPropRef := range itemSchema.Properties {
					nestedFieldName := strutil.PascalCase(nestedPropName)
					nestedPropSchema := nestedPropRef.GetOrNil()
					if nestedPropSchema != nil {
						l.ProcessProperty(nestedPropSchema, nestedFieldName, nestedPropName, itemTypeName)
					}
					processedProps[nestedFieldName] = nestedPropRef
				}
				itemSchema.Properties = processedProps
				propSchema.GoType = "[]" + itemTypeName
			} else {
				// Primitive array items
				itemGoType := ToGoType(itemSchema.Type.PrimaryType(), itemSchema.Format)
				propSchema.GoType = "[]" + itemGoType
			}
		}
	case TypeObject:
		// Objects with no properties are generic maps
		if len(propSchema.Properties) == 0 {
			propSchema.GoType = GoTypeMapString
		} else {
			// Nested object type name includes parent name to avoid conflicts
			nestedTypeName := parentName + fieldName
			propSchema.GoType = nestedTypeName
			// Recursively process nested object properties
			processedProps := make(map[string]*ref.Ref[Schema])
			for nestedPropName, nestedPropRef := range propSchema.Properties {
				nestedFieldName := strutil.PascalCase(nestedPropName)
				nestedPropSchema := nestedPropRef.GetOrNil()
				if nestedPropSchema != nil {
					l.ProcessProperty(nestedPropSchema, nestedFieldName, nestedPropName, nestedTypeName)
				}
				processedProps[nestedFieldName] = nestedPropRef
			}
			propSchema.Properties = processedProps
		}
	default:
		propSchema.GoType = ToGoType(primaryType, propSchema.Format)
	}
}

// LoadRef loads a schema from a $ref path relative to baseDir.
func (l *Loader) LoadRef(refPath string) (*Schema, error) {
	// Handle relative paths
	fullPath := filepath.Join(l.baseDir, refPath)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("reading ref %s: %w", fullPath, err)
	}

	// Update baseDir for nested refs
	oldBaseDir := l.baseDir
	l.baseDir = filepath.Dir(fullPath)
	defer func() { l.baseDir = oldBaseDir }()

	schema, err := l.LoadSchemaFromBytes(data, filepath.Base(refPath))
	if err != nil {
		return nil, err
	}

	return schema, nil
}

// resolveSchemaByRef looks up a schema by its ref name.
func (l *Loader) resolveSchemaByRef(refName string) *Schema {
	// First try direct lookup
	if schema, ok := l.schemas[refName]; ok {
		return schema
	}
	// Then try mapped name
	if actualName, ok := l.refToName[refName]; ok {
		if schema, ok := l.schemas[actualName]; ok {
			return schema
		}
	}
	return nil
}

// Schemas returns all loaded schemas.
func (l *Loader) Schemas() map[string]*Schema {
	return l.schemas
}

// SetSchema adds or updates a schema in the cache.
func (l *Loader) SetSchema(name string, schema *Schema) {
	l.schemas[name] = schema
}

// SetRefToName maps a ref name to an actual schema name.
func (l *Loader) SetRefToName(refName, actualName string) {
	l.refToName[refName] = actualName
}

// ExtractSchemaNameFromRef extracts the schema name from a $ref path.
// e.g., "./User.yaml" -> "User", "#/components/schemas/User" -> "User"
func ExtractSchemaNameFromRef(refPath string) string {
	// Handle internal refs (#/components/schemas/Name)
	if strings.HasPrefix(refPath, "#/") {
		parts := strings.Split(refPath, "/")
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
		return ""
	}

	// Handle file refs (./Name.yaml, ../schemas/Name.yaml)
	base := filepath.Base(refPath)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}
