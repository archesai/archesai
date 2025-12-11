package spec

import (
	"strings"

	"github.com/archesai/archesai/internal/schema"
	"github.com/archesai/archesai/internal/strutil"
)

// EntityOperations tracks which CRUD operations exist for an entity.
type EntityOperations struct {
	// HasList is true if a list operation exists.
	HasList bool

	// HasGet is true if a get operation exists.
	HasGet bool

	// HasCreate is true if a create operation exists.
	HasCreate bool

	// HasUpdate is true if an update operation exists.
	HasUpdate bool

	// HasDelete is true if a delete operation exists.
	HasDelete bool

	// ListHasFilterParam is true if the list operation has a filter parameter.
	ListHasFilterParam bool

	// ListHasPageParam is true if the list operation has a page parameter.
	ListHasPageParam bool

	// ListHasSortParam is true if the list operation has a sort parameter.
	ListHasSortParam bool

	// CreateRequiresParentID is true if create needs a parent ID (nested resources).
	CreateRequiresParentID bool

	// UpdateRequiresParentID is true if update needs a parent ID (nested resources).
	UpdateRequiresParentID bool

	// DeleteRequiresParentID is true if delete needs a parent ID (nested resources).
	DeleteRequiresParentID bool

	// ParentIDParamName is the parent ID parameter name for nested resources.
	ParentIDParamName string

	// ListOp is the list operation if it exists.
	ListOp *Operation

	// GetOp is the get operation if it exists.
	GetOp *Operation

	// CreateOp is the create operation if it exists.
	CreateOp *Operation

	// UpdateOp is the update operation if it exists.
	UpdateOp *Operation

	// DeleteOp is the delete operation if it exists.
	DeleteOp *Operation
}

// IsNested returns true if this entity is a nested resource (requires parent ID).
// This is determined by checking if the list operation has a parent path parameter.
func (ops *EntityOperations) IsNested() bool {
	if ops == nil || ops.ListOp == nil {
		return false
	}
	// Count path parameters in the list operation
	pathParamCount := 0
	for _, param := range ops.ListOp.Parameters {
		if param.In == paramLocationPath {
			pathParamCount++
		}
	}
	// Nested resources have at least one path param (the parent ID)
	return pathParamCount > 0
}

// FindEntityOperations finds CRUD operations for an entity by name.
func (s *Spec) FindEntityOperations(entityName string) EntityOperations {
	ops := EntityOperations{}
	entityLower := strings.ToLower(entityName)

	for i := range s.Operations {
		op := &s.Operations[i]
		opLower := strings.ToLower(op.ID)

		// Match common CRUD patterns
		switch {
		case strings.Contains(opLower, "list"+entityLower) ||
			strings.Contains(opLower, "get"+entityLower+"s"):
			ops.HasList = true
			ops.ListOp = op
			// Check for filter/page/sort parameters
			ops.ListHasFilterParam, ops.ListHasPageParam, ops.ListHasSortParam = detectListParams(
				op,
			)

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
func detectListParams(op *Operation) (hasFilter, hasPage, hasSort bool) {
	for _, param := range op.Parameters {
		if param.In != "query" {
			continue
		}
		// Skip auto-generated parameters (they have servermodels.* GoType)
		if param.Schema != nil && strings.Contains(param.GoType, "servermodels.") {
			continue
		}
		switch strings.ToLower(param.Title) {
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

// detectParentIDParam checks if an operation has a parent ID path parameter.
func detectParentIDParam(op *Operation) string {
	var pathParams []string
	for _, param := range op.Parameters {
		if param.In == paramLocationPath {
			pathParams = append(pathParams, param.Title)
		}
	}

	// For operations with 2+ path params, the first one is typically the parent
	if len(pathParams) >= 2 {
		return pathParams[0]
	}

	return ""
}

// EntityFrontendData provides all data needed for frontend templates.
type EntityFrontendData struct {
	// Entity is the schema being rendered.
	Entity *schema.Schema

	// EntityName is the PascalCase entity name.
	EntityName string

	// EntityNameLower is the camelCase entity name.
	EntityNameLower string

	// EntityNameKebab is the kebab-case entity name.
	EntityNameKebab string

	// EntityKey is the snake_case plural entity key (e.g., "users", "pipeline_steps").
	EntityKey string

	// Columns are the DataTable column definitions.
	Columns []schema.DataTableColumn

	// CreateFormFields are the form field definitions for create operations.
	CreateFormFields []schema.FormField

	// UpdateFormFields are the form field definitions for update operations.
	UpdateFormFields []schema.FormField

	// Operations contains CRUD operation info for this entity.
	Operations EntityOperations

	// ProjectName is the project module path.
	ProjectName string

	// ImportPath is the generated client import path.
	ImportPath string
}

// BuildEntityFrontendData creates the frontend template data for an entity.
func (s *Spec) BuildEntityFrontendData(entity *schema.Schema) *EntityFrontendData {
	data := &EntityFrontendData{
		Entity:          entity,
		EntityName:      entity.Title,
		EntityNameLower: strutil.CamelCase(entity.Title),
		EntityNameKebab: strutil.KebabCase(entity.Title),
		EntityKey:       strutil.SnakeCase(entity.Title) + "s",
		ProjectName:     s.ProjectName,
		ImportPath: "#lib/client/" + strings.ToLower(
			entity.Title,
		) + "/" + strings.ToLower(
			entity.Title,
		),
	}

	// Build columns from schema properties
	data.Columns = entity.ToDataTableColumns()

	// Map operations for this entity (do this first to access request bodies)
	data.Operations = s.FindEntityOperations(entity.Title)

	// Build form fields from request bodies of create/update operations
	if data.Operations.CreateOp != nil && data.Operations.CreateOp.RequestBody != nil {
		data.CreateFormFields = data.Operations.CreateOp.RequestBody.ToFormFields()
	}
	if data.Operations.UpdateOp != nil && data.Operations.UpdateOp.RequestBody != nil {
		data.UpdateFormFields = data.Operations.UpdateOp.RequestBody.ToFormFields()
	}

	return data
}
