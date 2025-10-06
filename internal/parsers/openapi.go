package parsers

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/speakeasy-api/openapi/openapi"
)

// OpenAPIParser wraps an OpenAPI document and provides parsing utilities
type OpenAPIParser struct {
	openAPIDoc *openapi.OpenAPI
}

// NewOpenAPIParser creates a new OpenAPIParser instance
func NewOpenAPIParser() *OpenAPIParser {
	return &OpenAPIParser{}
}

// Parse parses an OpenAPI specification file and returns the document
func (p *OpenAPIParser) Parse(specPath string) (*openapi.OpenAPI, error) {
	// Create a new context
	ctx := context.Background()

	// Open the spec file
	f, err := os.Open(specPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			err = fmt.Errorf("failed to close file: %w", cerr)
		}
	}()

	// Parse the OpenAPI document
	doc, validationErrs, err := openapi.Unmarshal(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal OpenAPI document: %w", err)
	}

	// Resolve all references in the document
	additionalErrors, err := doc.ResolveAllReferences(ctx, openapi.ResolveAllOptions{
		OpenAPILocation: specPath,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to resolve references: %w", err)
	}
	validationErrs = append(validationErrs, additionalErrors...)

	// If there are validation errors, return them
	if len(validationErrs) > 0 {
		var msgs []string
		for _, ve := range validationErrs {
			msgs = append(msgs, ve.Error())
		}
		return nil, fmt.Errorf("OpenAPI validation errors:\n%s", strings.Join(msgs, "\n"))
	}

	// Store the parsed document
	p.openAPIDoc = doc
	return doc, nil
}

// ExtractOperations extracts all operations from the OpenAPI spec
func ExtractOperations(doc *openapi.OpenAPI) ([]OperationDef, error) {
	if doc == nil {
		return nil, nil
	}

	var operations []OperationDef
	ctx := context.Background()

	// Walk through the OpenAPI document
	for item := range openapi.Walk(ctx, doc) {
		err := item.Match(openapi.Matcher{
			Operation: func(op *openapi.Operation) error {
				if op == nil {
					return nil
				}

				if len(op.Tags) != 1 {
					return fmt.Errorf("operation must have exactly one tag")
				}

				if op.GetOperationID() == "" {
					return fmt.Errorf("operation must have an operationId")
				}

				if op.GetSummary() == "" {
					return fmt.Errorf("operation must have a summary")
				}

				// Extract HTTP method and path from the walker location
				method, path := openapi.ExtractMethodAndPath(item.Location)

				requestBody, err := ExtractRequestBody(op, doc)
				if err != nil {
					return fmt.Errorf(
						"failed to extract request body for %s %s: %w",
						method,
						path,
						err,
					)
				}

				responses, err := ExtractResponses(op, doc)
				if err != nil {
					return fmt.Errorf(
						"failed to extract responses for %s %s: %w",
						method,
						path,
						err,
					)
				}

				operationDef := OperationDef{
					Method:                strings.ToUpper(method),
					Path:                  path,
					ID:                    op.GetOperationID(),
					Description:           op.GetSummary(),
					Tag:                   op.Tags[0],
					Parameters:            ExtractParameters(op),
					Security:              ExtractSecurityRequirements(op, doc),
					Responses:             responses,
					XCodegenCustomHandler: ExtractXCodegenCustomHandler(op), // default to false
					XCodegenRepository:    ExtractXCodegenRepository(op),    // default to ""
					RequestBody:           requestBody,
				}

				operations = append(operations, operationDef)

				return nil
			},
		})
		if err != nil {
			return nil, err
		}
	}

	// Sort operations by ID for consistent ordering in generated code
	sort.Slice(operations, func(i, j int) bool {
		return operations[i].ID < operations[j].ID
	})

	return operations, nil
}

// ExtractXCodegenCustomHandler checks if the operation has the x-codegen-custom-handler extension set to true
func ExtractXCodegenCustomHandler(op *openapi.Operation) bool {
	if op.Extensions != nil {
		if customHandlerExt := op.Extensions.GetOrZero("x-codegen-custom-handler"); customHandlerExt != nil {
			return (customHandlerExt.Value == "true")
		}
	}
	return false
}

// ExtractXCodegenRepository extracts the repository name from x-codegen-repository extension
func ExtractXCodegenRepository(op *openapi.Operation) string {
	if op.Extensions != nil {
		if repoExt := op.Extensions.GetOrZero("x-codegen-repository"); repoExt != nil {
			if repoExt.Value != "" {
				return repoExt.Value
			}
		}
	}
	return ""
}

// ExtractParameters extracts all parameters from an operation
func ExtractParameters(op *openapi.Operation) []ParamDef {
	var params []ParamDef

	if op.GetParameters() == nil {
		return params
	}

	for _, paramRef := range op.GetParameters() {
		param := paramRef.GetResolvedObject()
		if param == nil {
			continue
		}

		// Create embedded SchemaDef for the parameter
		schemaDef := &SchemaDef{
			Name:        PascalCase(param.Name),
			Description: param.GetDescription(),
			Required:    []string{}, // Parameters handle required differently
		}

		// Extract type information from schema
		if param.Schema != nil {
			schema := param.Schema.GetResolvedObject()
			if schema != nil && schema.Left != nil {
				// Get type info
				types := schema.Left.GetType()
				if len(types) > 0 {
					schemaDef.Type = string(types[0])
				}
				schemaDef.Format = schema.Left.GetFormat()
				schemaDef.Schema = schema.Left

				// Parameters are used in controllers, so pass empty currentPackage
				schemaDef.GoType = SchemaToGoType(schema.Left, nil, "")

				// Make optional parameters pointers (unless they're already slices/maps)
				if !param.GetRequired() &&
					!strings.HasPrefix(schemaDef.GoType, "[]") &&
					!strings.HasPrefix(schemaDef.GoType, "map") &&
					!strings.HasPrefix(schemaDef.GoType, "*") {
					schemaDef.GoType = "*" + schemaDef.GoType
				}
			}
		}

		// Create ParamDef with embedded SchemaDef
		paramDef := ParamDef{
			SchemaDef: schemaDef,
			In:        string(param.In),
			Style:     string(param.GetStyle()),
			Explode:   param.GetExplode(),
		}

		// Override Required at the ParamDef level since it's different for parameters
		if param.GetRequired() {
			paramDef.Required = []string{schemaDef.Name}
		}

		params = append(params, paramDef)
	}

	return params
}

// ExtractSecurityRequirements extracts security requirements from an operation
func ExtractSecurityRequirements(op *openapi.Operation, doc *openapi.OpenAPI) []SecurityDef {
	var securityDefs []SecurityDef

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
		if secReq != nil {
			for secSchemeName, scopes := range secReq.All() {
				// Look up the security scheme in the global security definitions
				if doc != nil && doc.Components != nil && doc.Components.SecuritySchemes != nil {
					secSchemeRef := doc.Components.SecuritySchemes.GetOrZero(secSchemeName)
					if secSchemeRef != nil {
						secScheme := secSchemeRef.GetResolvedObject()
						if secScheme != nil {
							secReqData := SecurityDef{
								Name:   secSchemeName,
								Type:   string(secScheme.Type),
								Scopes: scopes,
							}

							if secScheme.Type == "http" &&
								secScheme.Scheme != nil &&
								*secScheme.Scheme == "bearer" {
								secReqData.Scheme = *secScheme.Scheme
							}

							if secScheme.Type == "apiKey" &&
								secScheme.In != nil &&
								*secScheme.In == "cookie" {
								// Store in Scheme for consistency
								secReqData.Scheme = "cookie"
							}

							securityDefs = append(securityDefs, secReqData)
						}
					}
				}
			}
		}
	}

	return securityDefs
}

// ExtractResponses extracts all responses from an operation
func ExtractResponses(op *openapi.Operation, doc *openapi.OpenAPI) ([]ResponseDef, error) {
	var responses []ResponseDef
	if op.GetResponses() == nil {
		return responses, nil
	}

	for statusCode, responseRef := range op.GetResponses().All() {
		response := responseRef.GetResolvedObject()
		if response != nil {
			// Initialize response definition with basic info
			responseDef := ResponseDef{
				StatusCode: statusCode,
			}

			// Try to extract schema if present
			mediaType := response.Content.GetOrZero("application/json")
			if mediaType != nil && mediaType.Schema != nil {
				schemaObj := mediaType.Schema.GetResolvedObject()
				if schemaObj != nil && schemaObj.Left != nil {
					responseName := fmt.Sprintf("%s%sResponse", op.GetOperationID(), statusCode)
					jsonParser := NewJSONSchemaParser().WithOpenAPIDoc(doc)
					processed, err := jsonParser.ExtractSchema(schemaObj.Left, &responseName, "")
					if err != nil {
						return nil, fmt.Errorf(
							"failed to process response schema for status code %s: %w",
							statusCode,
							err,
						)
					}
					// Add the response description to the schema
					if processed.Description == "" {
						processed.Description = response.GetDescription()
					}
					responseDef.SchemaDef = processed
				}
			} else {
				// For responses without a schema, create a minimal SchemaDef
				responseDef.SchemaDef = &SchemaDef{
					Description: response.GetDescription(),
				}
			}

			// Add the response even if it doesn't have a schema (e.g., error responses using $ref)
			responses = append(responses, responseDef)
		}
	}

	return responses, nil
}

// ExtractRequestBody checks if an operation has a required request body and extracts its schema
func ExtractRequestBody(op *openapi.Operation, doc *openapi.OpenAPI) (*RequestBodyDef, error) {
	if op.RequestBody == nil {
		return nil, nil
	}

	rb := op.RequestBody.GetResolvedObject()
	if rb == nil {
		return nil, fmt.Errorf("failed to resolve request body")
	}

	// Extract schema from request body content
	if rb.Content != nil {
		if jsonContent := rb.Content.GetOrZero("application/json"); jsonContent != nil {
			if jsonContent.Schema != nil {
				schema := jsonContent.Schema.GetResolvedObject()
				if schema != nil && schema.Left != nil {
					requestName := fmt.Sprintf("%sRequestBody", op.GetOperationID())
					jsonParser := NewJSONSchemaParser().WithOpenAPIDoc(doc)
					processed, err := jsonParser.ExtractSchema(schema.Left, &requestName, "")
					if err != nil {
						return nil, fmt.Errorf("failed to process request body schema: %w", err)
					}

					return &RequestBodyDef{
						SchemaDef: processed,
						Required:  rb.GetRequired(),
					}, nil
				}
			}
		}
	}

	return nil, nil
}

// ExtractComponentSchemas processes all schemas from the OpenAPI document
func ExtractComponentSchemas(doc *openapi.OpenAPI) ([]*SchemaDef, error) {
	results := []*SchemaDef{}

	if doc == nil {
		return results, nil
	}

	if doc.Components == nil || doc.Components.Schemas == nil {
		return results, nil
	}

	// Process each schema
	for schemaName, schemaRef := range doc.Components.Schemas.All() {
		if schemaRef == nil {
			continue
		}

		// Get the resolved schema object
		schemaObj := schemaRef.GetResolvedObject()
		if schemaObj == nil || schemaObj.Left == nil {
			continue
		}

		// Determine schema type from x-codegen extension
		schemaType := ""
		if schemaObj.Left.Extensions != nil {
			if xExt := schemaObj.Left.Extensions.GetOrZero("x-codegen"); xExt != nil {
				parser := &XCodegenParser{}
				if xcodegen, err := parser.ParseExtension(xExt, schemaName); err == nil &&
					xcodegen != nil {
					schemaType = string(xcodegen.GetSchemaType())
				}
			}
		}

		jsonParser := NewJSONSchemaParser().WithOpenAPIDoc(doc)
		processed, err := jsonParser.ExtractSchema(schemaObj.Left, &schemaName, schemaType)
		if err != nil {
			return nil, fmt.Errorf("failed to process schema %s: %w", schemaName, err)
		}

		results = append(results, processed)
	}

	return results, nil
}
