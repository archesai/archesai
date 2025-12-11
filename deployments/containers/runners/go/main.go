package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// ContainerRequest is the JSON structure sent to the container via stdin
type ContainerRequest struct {
	SchemaIn  json.RawMessage `json:"schema_in"`  // JSON Schema for input validation
	SchemaOut json.RawMessage `json:"schema_out"` // JSON Schema for output validation
	Input     json.RawMessage `json:"input"`      // The actual input data
}

// ContainerResponse is the JSON structure returned by the container via stdout
type ContainerResponse struct {
	OK     bool            `json:"ok"`               // Whether execution was successful
	Output json.RawMessage `json:"output,omitempty"` // The output data (if ok=true)
	Error  *ContainerError `json:"error,omitempty"`  // Error details (if ok=false)
}

// ContainerError represents an error from the container
type ContainerError struct {
	Message string `json:"message"`           // Error message
	Details string `json:"details,omitempty"` // Optional additional details
}

// executeFunction is defined in execute.go
// It can be replaced by mounting a custom execute.go or building a derived image

func main() {
	var response ContainerResponse

	// Ensure we always output a JSON response
	defer func() {
		enc := json.NewEncoder(os.Stdout)
		_ = enc.Encode(response)
	}()

	// Read input from stdin
	rawInput, err := io.ReadAll(os.Stdin)
	if err != nil {
		response = ContainerResponse{
			OK:    false,
			Error: &ContainerError{Message: fmt.Sprintf("failed to read input: %v", err)},
		}
		return
	}

	if len(rawInput) == 0 {
		response = ContainerResponse{
			OK:    false,
			Error: &ContainerError{Message: "no input provided"},
		}
		return
	}

	// Parse the request
	var request ContainerRequest
	if err := json.Unmarshal(rawInput, &request); err != nil {
		response = ContainerResponse{
			OK:    false,
			Error: &ContainerError{Message: fmt.Sprintf("failed to parse request: %v", err)},
		}
		return
	}

	// Validate that input is present
	if len(request.Input) == 0 {
		response = ContainerResponse{
			OK:    false,
			Error: &ContainerError{Message: "missing 'input' field in request"},
		}
		return
	}

	// Validate input against schema if provided
	if len(request.SchemaIn) > 0 {
		if err := validateSchema(request.Input, request.SchemaIn, "input"); err != nil {
			response = ContainerResponse{
				OK:    false,
				Error: &ContainerError{Message: err.Error()},
			}
			return
		}
	}

	// Execute the function
	output, err := executeFunction(request.Input)
	if err != nil {
		response = ContainerResponse{
			OK:    false,
			Error: &ContainerError{Message: fmt.Sprintf("execution failed: %v", err)},
		}
		return
	}

	// Validate output against schema if provided
	if len(request.SchemaOut) > 0 {
		if err := validateSchema(output, request.SchemaOut, "output"); err != nil {
			response = ContainerResponse{
				OK:    false,
				Error: &ContainerError{Message: err.Error()},
			}
			return
		}
	}

	// Return success response
	response = ContainerResponse{
		OK:     true,
		Output: output,
	}
}

// validateSchema validates data against a JSON schema
func validateSchema(data json.RawMessage, schemaBytes json.RawMessage, dataType string) error {
	compiler := jsonschema.NewCompiler()

	// Add schema resource
	if err := compiler.AddResource("schema.json", bytes.NewReader(schemaBytes)); err != nil {
		return fmt.Errorf("failed to add %s schema: %v", dataType, err)
	}

	// Compile schema
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return fmt.Errorf("failed to compile %s schema: %v", dataType, err)
	}

	// Parse data to any for validation
	var dataAny any
	if err := json.Unmarshal(data, &dataAny); err != nil {
		return fmt.Errorf("failed to unmarshal %s for validation: %v", dataType, err)
	}

	// Validate
	if err := schema.Validate(dataAny); err != nil {
		return fmt.Errorf("%s validation failed: %v", dataType, err)
	}

	return nil
}
