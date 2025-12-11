// Package generators provides code generation from OpenAPI specifications.
package generators

import (
	"regexp"
	"strings"

	"github.com/archesai/archesai/internal/spec"
	"github.com/archesai/archesai/internal/strutil"
)

// BuildEntityFrontendData creates the frontend template data for an entity.
func BuildEntityFrontendData(ctx *GeneratorContext, entity *spec.Schema) *EntityFrontendData {
	data := &EntityFrontendData{
		Entity:          entity,
		EntityName:      entity.Name,
		EntityNameLower: strutil.CamelCase(entity.Name),
		EntityNameKebab: strutil.KebabCase(entity.Name),
		EntityKey:       strutil.SnakeCase(entity.Name) + "s",
		ProjectName:     ctx.ProjectName,
		ImportPath: "#lib/client/" + strings.ToLower(
			entity.Name,
		) + "/" + strings.ToLower(
			entity.Name,
		),
	}

	// Build columns from schema properties
	data.Columns = buildColumnsFromSchema(entity)

	// Map operations for this entity (do this first to access request bodies)
	data.Operations = findEntityOperations(ctx, entity.Name)

	// Build form fields from request bodies of create/update operations
	if data.Operations.CreateOp != nil && data.Operations.CreateOp.RequestBody != nil {
		data.CreateFormFields = buildFormFieldsFromSchema(data.Operations.CreateOp.RequestBody.Schema)
	}
	if data.Operations.UpdateOp != nil && data.Operations.UpdateOp.RequestBody != nil {
		data.UpdateFormFields = buildFormFieldsFromSchema(data.Operations.UpdateOp.RequestBody.Schema)
	}

	return data
}

// IsNestedResource checks if an entity is a nested resource (requires parent ID).
// This is determined by checking if the list operation has a parent path parameter.
func IsNestedResource(ops EntityOperations) bool {
	if ops.ListOp == nil {
		return false
	}
	// Count path parameters in the list operation
	pathParamCount := 0
	for _, param := range ops.ListOp.Parameters {
		if param.In == "path" {
			pathParamCount++
		}
	}
	// Nested resources have at least one path param (the parent ID)
	return pathParamCount > 0
}

// buildColumnsFromSchema creates DataTable columns from schema properties.
func buildColumnsFromSchema(schema *spec.Schema) []DataTableColumn {
	var columns []DataTableColumn

	for _, prop := range schema.GetSortedProperties() {
		col := mapPropertyToColumn(prop, schema)
		if col != nil {
			columns = append(columns, *col)
		}
	}

	return columns
}

// mapPropertyToColumn maps a schema property to a DataTable column.
func mapPropertyToColumn(prop *spec.Schema, parent *spec.Schema) *DataTableColumn {
	// Skip ID for now (usually handled specially as a link)
	if prop.Name == "ID" {
		return nil
	}

	// Use JSONTag if available, otherwise fallback to CamelCase
	accessorKey := prop.JSONTag
	if accessorKey == "" {
		accessorKey = strutil.CamelCase(prop.Name)
	}
	// Strip ",omitempty" suffix if present
	if idx := strings.Index(accessorKey, ","); idx >= 0 {
		accessorKey = accessorKey[:idx]
	}

	col := &DataTableColumn{
		AccessorKey:  accessorKey,
		Label:        toHumanReadable(prop.Name),
		EnableFilter: true,
		EnableSort:   true,
	}

	// Determine filter variant and icon based on type
	switch {
	case len(prop.Enum) > 0:
		col.FilterVariant = "multiSelect"
		col.Icon = "TextIcon"
		for _, v := range prop.Enum {
			col.Options = append(col.Options, FilterOption{
				Label: toHumanReadable(v),
				Value: v,
			})
		}

	case prop.Format == "date-time" || prop.Format == "date":
		col.FilterVariant = "date"
		col.Icon = "CalendarIcon"

	case prop.GoType == "bool":
		col.FilterVariant = "boolean"
		col.Icon = "CheckIcon"

	default:
		col.FilterVariant = "text"
		col.Icon = "TextIcon"
	}

	// Check if this is a name column that should link to detail view
	if strings.ToLower(prop.Name) == "name" || strings.ToLower(prop.Name) == "title" {
		col.IsLink = true
		col.LinkParam = strutil.CamelCase(parent.Name) + "ID"
	}

	return col
}

// buildFormFieldsFromSchema creates form fields from schema properties.
func buildFormFieldsFromSchema(schema *spec.Schema) []FormField {
	var fields []FormField

	for _, prop := range schema.GetSortedProperties() {
		field := mapPropertyToFormField(prop)
		if field != nil {
			fields = append(fields, *field)
		}
	}

	return fields
}

// mapPropertyToFormField maps a schema property to a form field.
func mapPropertyToFormField(prop *spec.Schema) *FormField {
	// Skip special fields that shouldn't be in forms
	switch prop.Name {
	case "ID", "CreatedAt", "UpdatedAt":
		return nil
	}

	// Skip array fields for now (need special handling)
	if prop.Type.PrimaryType() == spec.SchemaTypeArray {
		return nil
	}

	// Use JSONTag if available, otherwise fallback to CamelCase
	fieldName := prop.JSONTag
	if fieldName == "" {
		fieldName = strutil.CamelCase(prop.Name)
	}
	// Strip ",omitempty" suffix if present
	if idx := strings.Index(fieldName, ","); idx >= 0 {
		fieldName = fieldName[:idx]
	}

	field := &FormField{
		Name:        fieldName,
		Label:       toHumanReadable(prop.Name),
		Description: prop.Description,
		Required:    !prop.Nullable,
	}

	// Determine field type based on schema
	switch {
	case len(prop.Enum) > 0:
		field.Type = "select"
		for _, v := range prop.Enum {
			field.Options = append(field.Options, FormFieldOption{
				Label: toHumanReadable(v),
				Value: v,
			})
		}

	case prop.Format == "date-time" || prop.Format == "date":
		field.Type = "date"

	case prop.GoType == "bool":
		field.Type = "checkbox"
		field.DefaultValue = "false"

	case prop.MaxLength != nil && *prop.MaxLength > 100:
		field.Type = "textarea"

	case prop.Format == "email":
		field.Type = "email"

	case prop.Format == "uri" || prop.Format == "url":
		field.Type = "url"

	case prop.Format == "password":
		field.Type = "password"

	default:
		field.Type = "text"
	}

	return field
}

// findEntityOperations finds CRUD operations for an entity.
func findEntityOperations(ctx *GeneratorContext, entityName string) EntityOperations {
	ops := EntityOperations{}
	entityLower := strings.ToLower(entityName)

	for i := range ctx.Spec.Operations {
		op := &ctx.Spec.Operations[i]
		opLower := strings.ToLower(op.ID)

		// Match common CRUD patterns
		switch {
		case strings.Contains(opLower, "list"+entityLower) ||
			strings.Contains(opLower, "get"+entityLower+"s"):
			ops.HasList = true
			ops.ListOp = op
			// Check for filter/page/sort parameters
			ops.ListHasFilterParam, ops.ListHasPageParam, ops.ListHasSortParam = detectListParams(op)

		case strings.HasPrefix(opLower, "get"+entityLower) &&
			!strings.Contains(opLower, "list"):
			ops.HasGet = true
			ops.GetOp = op

		case strings.Contains(opLower, "create"+entityLower) ||
			strings.Contains(opLower, "add"+entityLower):
			ops.HasCreate = true
			ops.CreateOp = op
			// Check for parent ID requirement
			parentID := detectParentIDParam(op)
			if parentID != "" {
				ops.CreateRequiresParentID = true
				ops.ParentIDParamName = parentID
			}

		case strings.Contains(opLower, "update"+entityLower) ||
			strings.Contains(opLower, "edit"+entityLower):
			ops.HasUpdate = true
			ops.UpdateOp = op
			// Check for parent ID requirement
			parentID := detectParentIDParam(op)
			if parentID != "" {
				ops.UpdateRequiresParentID = true
				if ops.ParentIDParamName == "" {
					ops.ParentIDParamName = parentID
				}
			}

		case strings.Contains(opLower, "delete"+entityLower) ||
			strings.Contains(opLower, "remove"+entityLower):
			ops.HasDelete = true
			ops.DeleteOp = op
			// Check for parent ID requirement
			parentID := detectParentIDParam(op)
			if parentID != "" {
				ops.DeleteRequiresParentID = true
				if ops.ParentIDParamName == "" {
					ops.ParentIDParamName = parentID
				}
			}
		}
	}

	return ops
}

// detectListParams checks if an operation has filter, page, and sort parameters
// that were originally defined in the OpenAPI spec (not auto-generated).
// Auto-generated params have GoType like "servermodels.Page" or "*servermodels.FilterNode".
func detectListParams(op *spec.Operation) (hasFilter, hasPage, hasSort bool) {
	for _, param := range op.Parameters {
		if param.In != "query" {
			continue
		}
		// Skip auto-generated parameters (they have servermodels.* GoType)
		if param.Schema != nil && strings.Contains(param.Schema.GoType, "servermodels.") {
			continue
		}
		switch strings.ToLower(param.Name) {
		case "filter":
			hasFilter = true
		case "page":
			hasPage = true
		case "sort":
			hasSort = true
		}
	}
	return
}

// detectParentIDParam checks if an operation has a parent ID path parameter
// (e.g., organizationID for nested resources like /organizations/{organizationID}/members).
func detectParentIDParam(op *spec.Operation) string {
	// Count path parameters - if there's more than one, the first is likely a parent ID
	var pathParams []string
	for _, param := range op.Parameters {
		if param.In == "path" {
			pathParams = append(pathParams, param.Name)
		}
	}

	// For operations with 2+ path params, the first one is typically the parent
	if len(pathParams) >= 2 {
		return pathParams[0]
	}

	return ""
}

// toHumanReadable converts camelCase/PascalCase to human-readable format.
func toHumanReadable(s string) string {
	// Insert spaces before uppercase letters
	re := regexp.MustCompile(`([a-z])([A-Z])`)
	result := re.ReplaceAllString(s, "${1} ${2}")

	// Also handle sequences of uppercase (e.g., "ID" -> "ID")
	re2 := regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
	result = re2.ReplaceAllString(result, "${1} ${2}")

	// Capitalize first letter
	if len(result) > 0 {
		result = strings.ToUpper(result[:1]) + result[1:]
	}

	return result
}
