// Package generators provides code generation from OpenAPI specifications.
package generators

import (
	"github.com/archesai/archesai/internal/spec"
)

// EntityFrontendData provides all data needed for frontend templates.
type EntityFrontendData struct {
	// Entity is the schema being rendered.
	Entity *spec.Schema

	// EntityName is the PascalCase entity name.
	EntityName string

	// EntityNameLower is the camelCase entity name.
	EntityNameLower string

	// EntityNameKebab is the kebab-case entity name.
	EntityNameKebab string

	// EntityKey is the snake_case plural entity key (e.g., "users", "pipeline_steps").
	EntityKey string

	// Columns are the DataTable column definitions.
	Columns []DataTableColumn

	// CreateFormFields are the form field definitions for create operations.
	CreateFormFields []FormField

	// UpdateFormFields are the form field definitions for update operations.
	UpdateFormFields []FormField

	// Operations contains CRUD operation info for this entity.
	Operations EntityOperations

	// ProjectName is the project module path.
	ProjectName string

	// ImportPath is the generated client import path.
	ImportPath string
}

// DataTableColumn represents a column in the generated DataTable.
type DataTableColumn struct {
	// AccessorKey is the property name to access.
	AccessorKey string

	// Label is the human-readable column header.
	Label string

	// FilterVariant is the filter type: text, multiSelect, date, boolean.
	FilterVariant string

	// Icon is the icon component name for this column.
	Icon string

	// EnableFilter enables filtering on this column.
	EnableFilter bool

	// EnableSort enables sorting on this column.
	EnableSort bool

	// IsLink indicates if this column should render as a link.
	IsLink bool

	// LinkParam is the route parameter name for the link.
	LinkParam string

	// Options are the available options for multiSelect filters.
	Options []FilterOption
}

// FilterOption represents an option in a multiSelect filter.
type FilterOption struct {
	Label string
	Value string
}

// FormField represents a field in the generated form.
type FormField struct {
	// Name is the field property name.
	Name string

	// Label is the human-readable field label.
	Label string

	// Description is optional help text.
	Description string

	// Type is the input type: text, textarea, select, checkbox, date.
	Type string

	// Required indicates if the field is required.
	Required bool

	// Options are the available options for select fields.
	Options []FormFieldOption

	// DefaultValue is the default value expression.
	DefaultValue string
}

// FormFieldOption represents an option in a select field.
type FormFieldOption struct {
	Label string
	Value string
}

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
	ListOp *spec.Operation

	// GetOp is the get operation if it exists.
	GetOp *spec.Operation

	// CreateOp is the create operation if it exists.
	CreateOp *spec.Operation

	// UpdateOp is the update operation if it exists.
	UpdateOp *spec.Operation

	// DeleteOp is the delete operation if it exists.
	DeleteOp *spec.Operation
}

// RouteTemplateData holds data for generating route files.
type FrontendRouteTemplateData struct {
	// Entity is the entity this route is for.
	Entity *EntityFrontendData

	// RouteType is the type of route: list, detail, create.
	RouteType string

	// RoutePath is the file-based route path.
	RoutePath string

	// ProjectName is the project module path.
	ProjectName string
}

// FrontendConfigTemplateData holds data for generating config files.
type FrontendConfigTemplateData struct {
	// ProjectName is the project module path.
	ProjectName string

	// PackageName is the package.json name.
	PackageName string

	// Entities are all entities in the project.
	Entities []*EntityFrontendData

	// Tags are the OpenAPI tags used for client code generation.
	Tags []string

	// APIHost is the API host URL.
	APIHost string

	// PlatformURL is the platform URL.
	PlatformURL string

	// Port is the development server port.
	Port int

	// SiteName is the display name for the site.
	SiteName string

	// SiteDescription is the site description for SEO.
	SiteDescription string

	// SiteURL is the base URL for the site.
	SiteURL string

	// SpecPath is the relative path to the bundled OpenAPI spec from web/.
	SpecPath string
}
