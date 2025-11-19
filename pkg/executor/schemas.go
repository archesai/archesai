package executor

import (
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// SchemaValidator handles JSON Schema validation
type SchemaValidator struct {
	compiler       *jsonschema.Compiler
	schemaIn       *jsonschema.Schema
	schemaOut      *jsonschema.Schema
	schemaInBytes  []byte
	schemaOutBytes []byte
}

// NewSchemaValidator creates a new schema validator with the given input and output schemas
func NewSchemaValidator(schemaInBytes, schemaOutBytes []byte) (*SchemaValidator, error) {
	compiler := jsonschema.NewCompiler()

	// Unmarshal schemas to any first
	var schemaInData any
	if err := json.Unmarshal(schemaInBytes, &schemaInData); err != nil {
		return nil, fmt.Errorf("unmarshal input schema: %w", err)
	}

	var schemaOutData any
	if err := json.Unmarshal(schemaOutBytes, &schemaOutData); err != nil {
		return nil, fmt.Errorf("unmarshal output schema: %w", err)
	}

	// Add schemas as resources
	schemaInURL := "https://example.com/schema_in.json"
	if err := compiler.AddResource(schemaInURL, schemaInData); err != nil {
		return nil, fmt.Errorf("add input schema: %w", err)
	}

	schemaOutURL := "https://example.com/schema_out.json"
	if err := compiler.AddResource(schemaOutURL, schemaOutData); err != nil {
		return nil, fmt.Errorf("add output schema: %w", err)
	}

	// Compile schemas
	schemaIn, err := compiler.Compile(schemaInURL)
	if err != nil {
		return nil, fmt.Errorf("compile input schema: %w", err)
	}

	schemaOut, err := compiler.Compile(schemaOutURL)
	if err != nil {
		return nil, fmt.Errorf("compile output schema: %w", err)
	}

	return &SchemaValidator{
		compiler:       compiler,
		schemaIn:       schemaIn,
		schemaOut:      schemaOut,
		schemaInBytes:  schemaInBytes,
		schemaOutBytes: schemaOutBytes,
	}, nil
}

// ValidateInput validates the input data against the input schema
func (v *SchemaValidator) ValidateInput(data any) error {
	if err := v.schemaIn.Validate(data); err != nil {
		return fmt.Errorf("input validation failed: %w", err)
	}
	return nil
}

// ValidateOutput validates the output data against the output schema
func (v *SchemaValidator) ValidateOutput(data any) error {
	if err := v.schemaOut.Validate(data); err != nil {
		return fmt.Errorf("output validation failed: %w", err)
	}
	return nil
}

// GetInputSchema returns the raw input schema bytes
func (v *SchemaValidator) GetInputSchema() []byte {
	return v.schemaInBytes
}

// GetOutputSchema returns the raw output schema bytes
func (v *SchemaValidator) GetOutputSchema() []byte {
	return v.schemaOutBytes
}
