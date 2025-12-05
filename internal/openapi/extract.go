package openapi

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"

	"github.com/archesai/archesai/internal/spec"
	"github.com/archesai/archesai/internal/strutil"
	"github.com/archesai/archesai/internal/typeconv"
)

// ExtractSpec extracts operations and component schemas from the OpenAPI spec
func (p *Parser) ExtractSpec() (*spec.Spec, error) {
	if p.doc == nil {
		return nil, fmt.Errorf("document not parsed; call Parse() first")
	}

	operations, err := extractOperations(p.doc)
	if err != nil {
		return nil, fmt.Errorf("failed to extract operations: %w", err)
	}

	schemas, err := extractSchemas(p.doc)
	if err != nil {
		return nil, fmt.Errorf("failed to extract component schemas: %w", err)
	}

	return &spec.Spec{
		Operations:  operations,
		Schemas:     schemas,
		Document:    p.doc,
		ProjectName: extractProjectName(p.doc),
	}, nil
}

// extractProjectName extracts the project name from the x-project-name extension
func extractProjectName(doc *v3.Document) string {
	if doc == nil || doc.Extensions == nil {
		return ""
	}

	if ext, ok := doc.Extensions.Get("x-project-name"); ok {
		var projectName string
		if err := ext.Decode(&projectName); err == nil {
			return projectName
		}
	}

	return ""
}

// extractOperations extracts all operations from the OpenAPI spec
func extractOperations(doc *v3.Document) ([]spec.Operation, error) {
	if doc == nil {
		return nil, fmt.Errorf("document not set")
	}

	var operations []spec.Operation

	// Iterate through all paths
	for pathPair := doc.Paths.PathItems.First(); pathPair != nil; pathPair = pathPair.Next() {
		path := pathPair.Key()
		pathItem := pathPair.Value()

		// Check each HTTP method
		methodOps := map[string]*v3.Operation{
			"GET":     pathItem.Get,
			"POST":    pathItem.Post,
			"PUT":     pathItem.Put,
			"PATCH":   pathItem.Patch,
			"DELETE":  pathItem.Delete,
			"HEAD":    pathItem.Head,
			"OPTIONS": pathItem.Options,
		}

		for method, op := range methodOps {
			if op == nil {
				continue
			}

			if len(op.Tags) != 1 {
				return nil, fmt.Errorf("operation %s %s must have exactly one tag", method, path)
			}

			if op.OperationId == "" {
				return nil, fmt.Errorf("operation %s %s must have an operationId", method, path)
			}

			if op.Summary == "" {
				return nil, fmt.Errorf("operation %s %s must have a summary", method, path)
			}

			requestBody, err := extractRequestBody(doc, op)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to extract request body for %s %s: %w",
					method,
					path,
					err,
				)
			}

			responses, err := extractResponses(doc, op)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to extract responses for %s %s: %w",
					method,
					path,
					err,
				)
			}

			operationDef := spec.Operation{
				Method:                strings.ToUpper(method),
				Path:                  path,
				ID:                    op.OperationId,
				Description:           op.Summary,
				Tag:                   op.Tags[0],
				Parameters:            extractParameters(op, pathItem),
				Security:              extractSecurityRequirements(doc, op),
				Responses:             responses,
				XCodegenCustomHandler: extractXCodegenCustomHandler(op),
				XCodegenRepository:    extractXCodegenRepository(op),
				XInternal:             extractXInternal(op),
				RequestBody:           requestBody,
			}

			operations = append(operations, operationDef)
		}
	}

	// Sort operations by ID for consistent ordering in generated code
	sort.Slice(operations, func(i, j int) bool {
		return operations[i].ID < operations[j].ID
	})

	return operations, nil
}

// extractParameters extracts all parameters from an operation and path item
func extractParameters(op *v3.Operation, pathItem *v3.PathItem) []spec.Param {
	var params []spec.Param

	// Collect parameters from both operation and path item
	allParams := make([]*v3.Parameter, 0)
	if pathItem.Parameters != nil {
		allParams = append(allParams, pathItem.Parameters...)
	}
	if op.Parameters != nil {
		allParams = append(allParams, op.Parameters...)
	}

	for _, param := range allParams {
		if param == nil {
			continue
		}

		// Create embedded Schema for the parameter
		schemaDef := &spec.Schema{
			Name:        strutil.PascalCase(param.Name),
			Description: param.Description,
			Required:    []string{}, // Parameters handle required differently
		}

		// Extract type information from schema
		if param.Schema != nil {
			schema := param.Schema.Schema()
			if schema != nil {
				// Get type info
				if len(schema.Type) > 0 {
					schemaDef.Type = schema.Type[0]
				}
				if schema.Format != "" {
					schemaDef.Format = schema.Format
				}
				schemaDef.Schema = schema

				// Parameters are used in controllers, so pass empty currentPackage
				schemaDef.GoType = typeconv.SchemaToGoType(schema, nil, "")

				// Make optional parameters pointers (unless they're already slices/maps)
				if param.Required == nil || !*param.Required {
					if !strings.HasPrefix(schemaDef.GoType, "[]") &&
						!strings.HasPrefix(schemaDef.GoType, "map") &&
						!strings.HasPrefix(schemaDef.GoType, "*") {
						schemaDef.GoType = "*" + schemaDef.GoType
					}
				}
			}
		}

		// Create Param with embedded Schema
		paramDef := spec.Param{
			Schema: schemaDef,
			In:     param.In,
		}
		if param.Style != "" {
			paramDef.Style = param.Style
		}
		if param.Explode != nil {
			paramDef.Explode = *param.Explode
		}

		// Override Required at the Param level since it's different for parameters
		if param.Required != nil && *param.Required {
			paramDef.Required = []string{schemaDef.Name}
		}

		params = append(params, paramDef)
	}

	return params
}

// extractSecurityRequirements extracts security requirements from an operation
func extractSecurityRequirements(doc *v3.Document, op *v3.Operation) []spec.Security {
	var securityDefs []spec.Security

	securityRequirements := op.Security

	// Important: If operation explicitly sets empty security array, it means no auth
	// Only use global security if the operation has no security field at all (nil)
	if securityRequirements == nil && doc != nil && doc.Security != nil {
		securityRequirements = doc.Security
	}

	// If security is explicitly empty, return empty slice (no auth required)
	if securityRequirements != nil && len(securityRequirements) == 0 {
		return securityDefs
	}

	for _, secReq := range securityRequirements {
		if secReq == nil || secReq.Requirements == nil {
			continue
		}

		for schemePair := secReq.Requirements.First(); schemePair != nil; schemePair = schemePair.Next() {
			secSchemeName := schemePair.Key()
			scopes := schemePair.Value()

			// Look up the security scheme in the global security definitions
			if doc != nil && doc.Components != nil && doc.Components.SecuritySchemes != nil {
				if secScheme, ok := doc.Components.SecuritySchemes.Get(secSchemeName); ok {
					if secScheme != nil {
						secReqData := spec.Security{
							Name:   secSchemeName,
							Type:   secScheme.Type,
							Scopes: scopes,
						}

						if secScheme.Type == "http" && secScheme.Scheme == "bearer" {
							secReqData.Scheme = secScheme.Scheme
						}

						if secScheme.Type == "apiKey" && secScheme.In == "cookie" {
							// Store in Scheme for consistency
							secReqData.Scheme = "cookie"
						}

						securityDefs = append(securityDefs, secReqData)
					}
				}
			}
		}
	}

	return securityDefs
}

// extractResponses extracts all responses from an operation
func extractResponses(doc *v3.Document, op *v3.Operation) ([]spec.ResponseDef, error) {
	var responses []spec.ResponseDef
	if op.Responses == nil || op.Responses.Codes == nil {
		return responses, nil
	}

	for statusCode, response := range op.Responses.Codes.FromNewest() {
		if response != nil {
			// Initialize response definition with basic info
			responseDef := spec.ResponseDef{
				StatusCode: statusCode,
				Headers:    make(map[string]*spec.Schema),
			}

			// Extract content-type and schema from response content
			if response.Content != nil {
				// Try application/json first
				for contentType, content := range response.Content.FromNewest() {
					responseDef.ContentType = contentType
					if content.Schema != nil {
						schema := content.Schema.Schema()
						if schema != nil {
							schema.Title = fmt.Sprintf("%s%sResponse", op.OperationId, statusCode)
							jsonParser := NewJSONSchemaParser(doc)
							processed, err := jsonParser.ParseBase(schema)
							if err != nil {
								return nil, fmt.Errorf(
									"failed to process response schema for status code %s: %w",
									statusCode,
									err,
								)
							}
							// Add the response description to the schema
							if processed.Description == "" {
								processed.Description = response.Description
							}
							responseDef.Schema = processed
						}
					}
				}
			}

			// Extract headers from the response
			if response.Headers != nil {
				for headerName, header := range response.Headers.FromNewest() {
					if header != nil && header.Schema != nil {
						schema := header.Schema.Schema()
						if schema != nil {
							// Create a Schema for the header
							headerDef := &spec.Schema{
								Name:        headerName,
								Description: header.Description,
							}

							// Get type info
							if len(schema.Type) > 0 {
								headerDef.Type = schema.Type[0]
							}
							if schema.Format != "" {
								headerDef.Format = schema.Format
							}
							headerDef.Schema = schema
							headerDef.GoType = typeconv.SchemaToGoType(schema, nil, "")

							responseDef.Headers[headerName] = headerDef
						}
					}
				}
			}

			// For responses without a schema, create a minimal Schema
			if responseDef.Schema == nil {
				responseDef.Schema = &spec.Schema{
					Description: response.Description,
				}
			}

			// Add the response even if it doesn't have a schema (e.g., error responses using $ref)
			responses = append(responses, responseDef)
		}
	}

	// Sort responses by status code (success responses first, then errors in numerical order)
	sort.Slice(responses, func(i, j int) bool {
		iCode, _ := strconv.Atoi(responses[i].StatusCode)
		jCode, _ := strconv.Atoi(responses[j].StatusCode)
		return iCode < jCode
	})

	return responses, nil
}

// extractRequestBody checks if an operation has a required request body and extracts its schema
func extractRequestBody(doc *v3.Document, op *v3.Operation) (*spec.RequestBody, error) {
	if op.RequestBody == nil {
		return nil, nil
	}

	rb := op.RequestBody
	if rb == nil {
		return nil, nil
	}

	// Extract schema from request body content
	if rb.Content != nil {
		if jsonContent, ok := rb.Content.Get("application/json"); ok {
			if jsonContent.Schema != nil {
				schema := jsonContent.Schema.Schema()
				if schema != nil {
					schema.Title = fmt.Sprintf("%sRequestBody", op.OperationId)
					jsonParser := NewJSONSchemaParser(doc)
					processed, err := jsonParser.ParseBase(schema)
					if err != nil {
						return nil, fmt.Errorf("failed to process request body schema: %w", err)
					}

					required := false
					if rb.Required != nil {
						required = *rb.Required
					}

					return &spec.RequestBody{
						Schema:   processed,
						Required: required,
					}, nil
				}
			}
		}
	}

	return nil, nil
}

// ExtractComponentSchemas processes all schemas from the OpenAPI document
func extractSchemas(doc *v3.Document) ([]*spec.Schema, error) {
	if doc == nil {
		return nil, fmt.Errorf("document not set")
	}

	if doc.Components == nil || doc.Components.Schemas == nil {
		return nil, nil
	}

	// Initialize results slice
	results := []*spec.Schema{}

	// Track processed schemas by title to avoid duplicates from ref resolution
	processedTitles := make(map[string]bool)

	// Process each schema
	for schemaPair := doc.Components.Schemas.First(); schemaPair != nil; schemaPair = schemaPair.Next() {
		schemaName := schemaPair.Key()
		schemaRef := schemaPair.Value()

		if schemaRef == nil {
			return nil, fmt.Errorf("schema %s is nil", schemaName)
		}

		// Get the resolved schema object
		schema := schemaRef.Schema()
		if schema == nil {
			return nil, fmt.Errorf("schema %s is nil", schemaName)
		}

		// Override the schema name with the component key
		if schema.Title == "" {
			schema.Title = schemaName
		}

		// Skip duplicates (can occur when refs resolve to same schema)
		if processedTitles[schema.Title] {
			continue
		}
		processedTitles[schema.Title] = true

		jsonParser := NewJSONSchemaParser(doc)
		processed, err := jsonParser.ParseBase(schema)
		if err != nil {
			return nil, fmt.Errorf("failed to process schema %s: %w", schemaName, err)
		}

		results = append(results, processed)
	}

	return results, nil
}

// extractXCodegenCustomHandler checks if the operation has the x-codegen-custom-handler extension set to true
func extractXCodegenCustomHandler(op *v3.Operation) bool {
	if op.Extensions == nil {
		return false
	}
	if val, ok := op.Extensions.Get("x-codegen-custom-handler"); ok {
		var boolVal bool
		if err := val.Decode(&boolVal); err == nil {
			return boolVal
		}
		var strVal string
		if err := val.Decode(&strVal); err == nil {
			return strVal == boolTrueString
		}
	}
	return false
}

// extractXCodegenRepository extracts the repository name from x-codegen-repository extension
func extractXCodegenRepository(op *v3.Operation) string {
	if op.Extensions == nil {
		return ""
	}
	if val, ok := op.Extensions.Get("x-codegen-repository"); ok {
		var strVal string
		if err := val.Decode(&strVal); err == nil {
			return strVal
		}
	}
	return ""
}

// extractXInternal extracts the x-internal extension value (e.g., "server", "config")
// When set, the operation should be imported from another package instead of generated
func extractXInternal(op *v3.Operation) string {
	if op.Extensions == nil {
		return ""
	}
	if val, ok := op.Extensions.Get("x-internal"); ok {
		var strVal string
		if err := val.Decode(&strVal); err == nil {
			return strVal
		}
	}
	return ""
}
