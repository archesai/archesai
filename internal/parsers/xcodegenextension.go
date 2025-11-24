package parsers

// XCodegenExtension represents Configuration for code generation from OpenAPI schemas
type XCodegenExtension struct {

	// Repository Repository generation configuration
	Repository *XCodegenExtensionRepository `json:"repository,omitempty" yaml:"repository,omitempty"`
}

// XCodegenExtensionRepository represents Repository generation configuration
type XCodegenExtensionRepository struct {

	// AdditionalMethods Additional repository methods to generate
	AdditionalMethods []XCodegenExtensionRepositoryAdditionalMethodsItem `json:"additionalMethods,omitempty" yaml:"additionalMethods,omitempty"`

	// ExcludeFromCreate Fields to exclude from Create operations (e.g., auto-generated or DB-set fields)
	ExcludeFromCreate []string `json:"excludeFromCreate,omitempty" yaml:"excludeFromCreate,omitempty"`

	// ExcludeFromUpdate Fields to exclude from Update operations (e.g., immutable fields)
	ExcludeFromUpdate []string `json:"excludeFromUpdate,omitempty" yaml:"excludeFromUpdate,omitempty"`

	// Indices Database indices to create
	Indices []string `json:"indices,omitempty" yaml:"indices,omitempty"`

	// Relations Foreign key relationships to other entities
	Relations []XCodegenExtensionRepositoryRelationsItem `json:"relations,omitempty" yaml:"relations,omitempty"`
}

// XCodegenExtensionRepositoryAdditionalMethodsItem represents a nested type for XCodegenExtension
type XCodegenExtensionRepositoryAdditionalMethodsItem struct {

	// Name Method name
	Name string `json:"name" yaml:"name"`

	// Params Method parameters
	Params []XCodegenExtensionRepositoryAdditionalMethodsItemParamsItem `json:"params,omitempty" yaml:"params,omitempty"`

	// Returns Return type (single or multiple)
	Returns string `json:"returns" yaml:"returns"`
}

// XCodegenExtensionRepositoryAdditionalMethodsItemParamsItem represents a nested type for XCodegenExtension
type XCodegenExtensionRepositoryAdditionalMethodsItemParamsItem struct {

	// Format Parameter format (e.g., uuid, email) - optional
	Format *string `json:"format,omitempty" yaml:"format,omitempty"`

	// Name Parameter name
	Name string `json:"name" yaml:"name"`

	// Type Parameter type (e.g., string, int)
	Type string `json:"type" yaml:"type"`
}

// XCodegenExtensionRepositoryRelationsItem represents a nested type for XCodegenExtension
type XCodegenExtensionRepositoryRelationsItem struct {

	// Field The field name in this entity that references another entity
	Field string `json:"field" yaml:"field"`

	// OnDelete Foreign key ON DELETE action
	OnDelete *string `json:"onDelete,omitempty" yaml:"onDelete,omitempty"`

	// OnUpdate Foreign key ON UPDATE action
	OnUpdate *string `json:"onUpdate,omitempty" yaml:"onUpdate,omitempty"`

	// References The table name being referenced (snake_case)
	References string `json:"references" yaml:"references"`

	// ReferencesField The field in the referenced table (defaults to 'id')
	ReferencesField *string `json:"referencesField,omitempty" yaml:"referencesField,omitempty"`
}
