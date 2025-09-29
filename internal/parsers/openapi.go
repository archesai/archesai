package parsers

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/speakeasy-api/openapi/jsonschema/oas3"
	"github.com/speakeasy-api/openapi/openapi"
)

// ParseOpenAPI parses an OpenAPI specification file and returns the document
func ParseOpenAPI(specPath string) (*openapi.OpenAPI, []string, error) {
	ctx := context.Background()

	// Open the spec file
	f, err := os.Open(specPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			err = fmt.Errorf("failed to close file: %w", cerr)
		}
	}()

	// Parse the OpenAPI document
	doc, validationErrs, err := openapi.Unmarshal(ctx, f)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal OpenAPI document: %w", err)
	}

	// Collect validation warnings
	var warnings []string
	for _, vErr := range validationErrs {
		warnings = append(warnings, vErr.Error())
	}

	// Resolve all references in the document
	resolveValidationErrs, resolveErrs := doc.ResolveAllReferences(ctx, openapi.ResolveAllOptions{
		OpenAPILocation: specPath,
	})
	if resolveErrs != nil {
		return nil, warnings, fmt.Errorf("failed to resolve references: %w", resolveErrs)
	}

	// Add resolve validation warnings
	for _, vErr := range resolveValidationErrs {
		warnings = append(warnings, vErr.Error())
	}

	return doc, warnings, nil
}

// ExtractOperations extracts all operations from the OpenAPI spec
func ExtractOperations(doc *openapi.OpenAPI) []OperationDef {
	if doc == nil {
		return nil
	}

	var operations []OperationDef
	ctx := context.Background()

	// Walk through the OpenAPI document
	for item := range openapi.Walk(ctx, doc) {
		_ = item.Match(openapi.Matcher{
			Operation: func(op *openapi.Operation) error {
				if op == nil {
					return nil
				}

				// Extract HTTP method and path from the walker location
				method, path := openapi.ExtractMethodAndPath(item.Location)

				// Capitalize first letter of operation ID for GoName
				opID := op.GetOperationID()
				goName := opID
				if len(opID) > 0 {
					goName = strings.ToUpper(opID[:1]) + opID[1:]
				}

				operationDef := OperationDef{
					Method:              strings.ToUpper(method),
					Path:                path,
					OperationID:         opID,
					Name:                opID,
					GoName:              goName,
					Description:         op.GetSummary(),
					Tags:                op.Tags,
					RequestBodyRequired: hasRequiredRequestBody(op),
				}

				// Extract parameters and split by type
				operationDef.Parameters = extractParameters(op)
				for _, param := range operationDef.Parameters {
					switch param.In {
					case "path":
						operationDef.PathParams = append(operationDef.PathParams, param)
					case "query":
						operationDef.QueryParams = append(operationDef.QueryParams, param)
					case "header":
						operationDef.HeaderParams = append(operationDef.HeaderParams, param)
					}
				}

				// Extract and process request body schema
				if required, schema := ExtractRequestBody(op); schema != nil {
					operationDef.RequestBodyRequired = required
					// Process the schema to get field definitions
					processed, err := ProcessSchema(schema, "RequestBody")
					if err == nil {
						operationDef.RequestBodySchema = processed
					}
				}

				// Extract security requirements
				operationDef.Security = extractSecurityRequirements(op, doc)

				// Extract responses
				operationDef.Responses = extractResponses(op)

				// Extract and process response schemas
				responseSchemas := ExtractResponseSchemas(op)
				if len(responseSchemas) > 0 {
					operationDef.ResponseSchemas = make(map[string]*ProcessedSchema)
					for statusCode, schema := range responseSchemas {
						processed, err := ProcessSchema(
							schema,
							fmt.Sprintf("Response%s", statusCode),
						)
						if err == nil {
							operationDef.ResponseSchemas[statusCode] = processed
						}
					}
				}

				operations = append(operations, operationDef)

				return nil
			},
		})
	}

	return operations
}

// extractParameters extracts all parameters from an operation
func extractParameters(op *openapi.Operation) []ParamDef {
	var params []ParamDef

	if op.GetParameters() == nil {
		return params
	}

	for _, paramRef := range op.GetParameters() {
		param := paramRef.GetResolvedObject()
		if param == nil {
			continue
		}

		paramDef := ParamDef{
			Name:        param.Name,
			In:          string(param.In),
			Required:    param.GetRequired(),
			Description: param.GetDescription(),
			Style:       string(param.GetStyle()),
			Explode:     param.GetExplode(),
		}

		if param.Schema != nil {
			// Check if it's a reference to another schema
			schemaRef := param.Schema.GetReference()
			if schemaRef != "" {
				// Extract schema name from reference
				// e.g., "../schemas/Page.yaml" -> "Page"
				parts := strings.Split(schemaRef.String(), "/")
				if len(parts) > 0 {
					schemaName := parts[len(parts)-1]
					// Remove .yaml extension if present
					schemaName = strings.TrimSuffix(schemaName, ".yaml")
					paramDef.Schema = schemaName
				}
			}

			// Also extract the type
			schema := param.Schema.GetResolvedObject()
			if schema != nil && schema.Left != nil {
				paramDef.Type = SchemaToGoType(schema.Left)
				paramDef.GoType = paramDef.Type
				paramDef.Format = schema.Left.GetFormat()
			}
		}

		params = append(params, paramDef)
	}

	return params
}

// extractSecurityRequirements extracts security requirements from an operation
func extractSecurityRequirements(op *openapi.Operation, doc *openapi.OpenAPI) []SecurityDef {
	var securityDefs []SecurityDef

	securityRequirements := op.Security
	if securityRequirements == nil && doc != nil && doc.Security != nil {
		securityRequirements = doc.Security
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

// hasRequiredRequestBody checks if an operation has a required request body
func hasRequiredRequestBody(op *openapi.Operation) bool {
	if op.RequestBody == nil {
		return false
	}

	rb := op.RequestBody.GetResolvedObject()
	if rb == nil {
		return false
	}

	return rb.GetRequired()
}

// extractResponses extracts all responses from an operation
func extractResponses(op *openapi.Operation) []ResponseDef {
	var responses []ResponseDef

	if op.GetResponses() != nil {
		for statusCode, responseRef := range op.GetResponses().All() {
			response := responseRef.GetResolvedObject()
			if response != nil {
				responseDef := ResponseDef{
					StatusCode:  statusCode,
					Description: response.GetDescription(),
				}

				// Check if it's a success response
				if code, err := strconv.Atoi(statusCode); err == nil && code >= 200 && code < 300 {
					responseDef.IsSuccess = true
				}

				responses = append(responses, responseDef)
			}
		}
	}

	return responses
}

// ExtractRequestBody checks if an operation has a required request body and extracts its schema
func ExtractRequestBody(op *openapi.Operation) (bool, *oas3.Schema) {
	if op.RequestBody == nil {
		return false, nil
	}

	rb := op.RequestBody.GetResolvedObject()
	if rb == nil {
		return false, nil
	}

	// Check if required
	required := rb.GetRequired()

	// Extract schema from request body content
	if rb.Content != nil {
		if jsonContent := rb.Content.GetOrZero("application/json"); jsonContent != nil {
			if jsonContent.Schema != nil {
				schema := jsonContent.Schema.GetResolvedObject()
				if schema != nil && schema.Left != nil {
					return required, schema.Left
				}
			}
		}
	}

	return required, nil
}

// ExtractResponseSchemas extracts schemas from operation responses
func ExtractResponseSchemas(op *openapi.Operation) map[string]*oas3.Schema {
	schemas := make(map[string]*oas3.Schema)

	if op.GetResponses() == nil {
		return schemas
	}

	for statusCode, responseRef := range op.GetResponses().All() {
		response := responseRef.GetResolvedObject()
		if response == nil {
			continue
		}

		if response.Content != nil {
			if jsonContent := response.Content.GetOrZero("application/json"); jsonContent != nil {
				if jsonContent.Schema != nil {
					schema := jsonContent.Schema.GetResolvedObject()
					if schema != nil && schema.Left != nil {
						schemas[statusCode] = schema.Left
					}
				}
			}
		}
	}

	return schemas
}

// ProcessAllSchemas processes all schemas from the OpenAPI document
func ProcessAllSchemas(doc *openapi.OpenAPI) (map[string]*ProcessedSchema, error) {
	results := make(map[string]*ProcessedSchema)

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

		processed, err := ProcessSchema(schemaObj.Left, schemaName)
		if err != nil {
			return nil, fmt.Errorf("failed to process schema %s: %w", schemaName, err)
		}

		results[schemaName] = processed
	}

	return results, nil
}
