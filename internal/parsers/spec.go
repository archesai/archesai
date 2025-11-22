package parsers

// SpecDef represents the entire specification definition.
type SpecDef struct {
	Operations []OperationDef // All operations in the spec
	Schemas    []*SchemaDef   // All schemas defined in the spec
}
