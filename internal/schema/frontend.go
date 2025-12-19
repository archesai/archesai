package schema

import (
	"strings"

	"github.com/archesai/archesai/internal/strutil"
)

// Filter variant constants for DataTable columns.
const (
	filterVariantText        = "text"
	filterVariantMultiSelect = "multiSelect"
	filterVariantDate        = "date"
	filterVariantBoolean     = "boolean"
)

// Format constants for non-standard OpenAPI formats.
const formatURL = "url"

// Form field type constants.
const (
	fieldTypeEmail    = "email"
	fieldTypeURL      = "url"
	fieldTypePassword = "password"
	fieldTypeCheckbox = "checkbox"
	fieldTypeTextarea = "textarea"
	fieldTypeSelect   = "select"
)

// Link column field names.
const (
	linkFieldName  = "name"
	linkFieldTitle = "title"
)

// Default value constants.
const defaultValueFalse = "false"

// DataTableColumn represents a column in a generated DataTable component.
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

// FormField represents a field in a generated form component.
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

// ToDataTableColumn converts this schema property to a DataTable column definition.
// Returns nil for special fields like ID that shouldn't be shown as columns.
func (s *Schema) ToDataTableColumn(parent *Schema) *DataTableColumn {
	if s == nil {
		return nil
	}

	// Skip ID for now (usually handled specially as a link)
	if s.Title == "ID" {
		return nil
	}

	// Use JSONTag if available, otherwise fallback to CamelCase
	accessorKey := s.JSONTag
	if accessorKey == "" {
		accessorKey = strutil.CamelCase(s.Title)
	}
	// Strip ",omitempty" suffix if present
	if idx := strings.Index(accessorKey, ","); idx >= 0 {
		accessorKey = accessorKey[:idx]
	}

	col := &DataTableColumn{
		AccessorKey:  accessorKey,
		Label:        strutil.HumanReadable(s.Title),
		EnableFilter: true,
		EnableSort:   true,
	}

	// Determine filter variant and icon based on type
	switch {
	case len(s.Enum) > 0:
		col.FilterVariant = filterVariantMultiSelect
		col.Icon = "TextIcon"
		for _, v := range s.Enum {
			col.Options = append(col.Options, FilterOption{
				Label: strutil.HumanReadable(v),
				Value: v,
			})
		}

	case s.Format == FormatDateTime || s.Format == FormatDate:
		col.FilterVariant = filterVariantDate
		col.Icon = "CalendarIcon"

	case s.GoType == GoTypeBool:
		col.FilterVariant = filterVariantBoolean
		col.Icon = "CheckIcon"

	default:
		col.FilterVariant = filterVariantText
		col.Icon = "TextIcon"
	}

	// Check if this is a name column that should link to detail view
	if strings.ToLower(s.Title) == linkFieldName || strings.ToLower(s.Title) == linkFieldTitle {
		col.IsLink = true
		if parent != nil {
			col.LinkParam = strutil.CamelCase(parent.Title) + "ID"
		}
	}

	return col
}

// ToDataTableColumns converts this schema's properties to DataTable column definitions.
func (s *Schema) ToDataTableColumns() []DataTableColumn {
	if s == nil {
		return nil
	}

	var columns []DataTableColumn
	for _, prop := range s.GetSortedProperties() {
		col := prop.ToDataTableColumn(s)
		if col != nil {
			columns = append(columns, *col)
		}
	}
	return columns
}

// ToFormField converts this schema property to a form field definition.
// Returns nil for special fields like ID, CreatedAt, UpdatedAt that shouldn't be in forms.
func (s *Schema) ToFormField() *FormField {
	if s == nil {
		return nil
	}

	// Skip special fields that shouldn't be in forms
	switch s.Title {
	case "ID", "CreatedAt", "UpdatedAt":
		return nil
	}

	// Skip array fields for now (need special handling)
	if s.Type.PrimaryType() == TypeArray {
		return nil
	}

	// Use JSONTag if available, otherwise fallback to CamelCase
	fieldName := s.JSONTag
	if fieldName == "" {
		fieldName = strutil.CamelCase(s.Title)
	}
	// Strip ",omitempty" suffix if present
	if idx := strings.Index(fieldName, ","); idx >= 0 {
		fieldName = fieldName[:idx]
	}

	field := &FormField{
		Name:        fieldName,
		Label:       strutil.HumanReadable(s.Title),
		Description: s.Description,
		Required:    !s.Nullable,
	}

	// Determine field type based on schema
	switch {
	case len(s.Enum) > 0:
		field.Type = fieldTypeSelect
		for _, v := range s.Enum {
			field.Options = append(field.Options, FormFieldOption{
				Label: strutil.HumanReadable(v),
				Value: v,
			})
		}

	case s.Format == FormatDateTime || s.Format == FormatDate:
		field.Type = filterVariantDate

	case s.GoType == GoTypeBool:
		field.Type = fieldTypeCheckbox
		field.DefaultValue = defaultValueFalse

	case s.MaxLength != nil && *s.MaxLength > 100:
		field.Type = fieldTypeTextarea

	case s.Format == FormatEmail:
		field.Type = fieldTypeEmail

	case s.Format == FormatURI || s.Format == formatURL:
		field.Type = fieldTypeURL

	case s.Format == FormatPassword:
		field.Type = fieldTypePassword

	default:
		field.Type = filterVariantText
	}

	return field
}

// ToFormFields converts this schema's properties to form field definitions.
func (s *Schema) ToFormFields() []FormField {
	if s == nil {
		return nil
	}

	var fields []FormField
	for _, prop := range s.GetSortedProperties() {
		field := prop.ToFormField()
		if field != nil {
			fields = append(fields, *field)
		}
	}
	return fields
}
