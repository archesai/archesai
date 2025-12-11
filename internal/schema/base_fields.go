package schema

import "github.com/archesai/archesai/internal/ref"

// BaseFieldDefs returns the standard base field definitions for entity schemas.
func BaseFieldDefs() map[string]*Schema {
	return map[string]*Schema{
		"ID": {
			Title:       "ID",
			Description: "Unique identifier for the resource",
			Type:        PropertyType{Types: []string{TypeString}},
			Format:      FormatUUID,
			GoType:      GoTypeUUID,
			JSONTag:     "id",
			YAMLTag:     "id",
		},
		"CreatedAt": {
			Title:       "CreatedAt",
			Description: "The date and time when the resource was created",
			Type:        PropertyType{Types: []string{TypeString}},
			Format:      FormatDateTime,
			GoType:      GoTypeTime,
			JSONTag:     "createdAt",
			YAMLTag:     "createdAt",
		},
		"UpdatedAt": {
			Title:       "UpdatedAt",
			Description: "The date and time when the resource was last updated",
			Type:        PropertyType{Types: []string{TypeString}},
			Format:      FormatDateTime,
			GoType:      GoTypeTime,
			JSONTag:     "updatedAt",
			YAMLTag:     "updatedAt",
		},
	}
}

// AddBaseFields adds id, createdAt, updatedAt fields to a schema if they don't exist.
// Also adds them to the required list.
func AddBaseFields(schema *Schema) {
	if schema.Properties == nil {
		schema.Properties = make(map[string]*ref.Ref[Schema])
	}

	for fieldName, defaults := range BaseFieldDefs() {
		if existingRef, ok := schema.Properties[fieldName]; ok {
			// Field exists - add description if missing
			existing := existingRef.GetOrNil()
			if existing != nil && existing.Description == "" {
				existing.Description = defaults.Description
			}
		} else {
			// Field doesn't exist - add it
			schema.Properties[fieldName] = ref.NewInline(defaults)
		}
	}

	// Add base fields to required list if not already present
	requiredFields := []string{"id", "createdAt", "updatedAt"}
	requiredMap := make(map[string]bool, len(schema.Required))
	for _, r := range schema.Required {
		requiredMap[r] = true
	}
	for _, bf := range requiredFields {
		if !requiredMap[bf] {
			schema.Required = append(schema.Required, bf)
		}
	}
}

// GetSchemasAlias returns the Go schemas package alias for an include name.
func GetSchemasAlias(include string) string {
	switch include {
	case "server":
		return "serverschemas"
	case "auth":
		return "authschemas"
	case "config":
		return "configschemas"
	case "pipelines":
		return "pipelinesschemas"
	case "executor":
		return "executorschemas"
	case "storage":
		return "storageschemas"
	default:
		return ""
	}
}
