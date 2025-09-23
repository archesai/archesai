package parsers

import (
	"context"
	"sort"
	"strings"

	"github.com/speakeasy-api/openapi/openapi"

	"github.com/archesai/archesai/internal/templates"
)

// OpenAPISchema handles OpenAPI document operations and extraction.
type OpenAPISchema struct {
	*openapi.OpenAPI
	FilePath string
}

// NewOpenAPISchema creates a new OpenAPISchema instance.
func NewOpenAPISchema(doc *openapi.OpenAPI, filepath string) *OpenAPISchema {
	return &OpenAPISchema{
		OpenAPI:  doc,
		FilePath: filepath,
	}
}

// ExtractOperations extracts all operations from the OpenAPI spec.
func (o *OpenAPISchema) ExtractOperations() []templates.OperationData {
	if o == nil || o.OpenAPI == nil {
		return nil
	}

	var operations []templates.OperationData
	ctx := context.Background()

	// Walk through the OpenAPI document
	for item := range openapi.Walk(ctx, o.OpenAPI) {
		_ = item.Match(openapi.Matcher{
			// Process operations
			Operation: func(op *openapi.Operation) error {
				if op == nil {
					return nil
				}

				// Extract HTTP method and path from the walker location
				method, path := openapi.ExtractMethodAndPath(item.Location)

				operation := templates.OperationData{
					Method: strings.ToUpper(method), // Convert to uppercase (GET, POST, etc.)
					Path:   path,
				}

				// Extract operation ID
				if op.OperationID != nil {
					operation.OperationID = *op.OperationID
				}

				// Extract description
				if op.Summary != nil {
					operation.Description = *op.Summary
				}

				// Extract tags
				operation.Tags = op.Tags

				// Extract parameters
				if op.Parameters != nil {
					for _, paramRef := range op.Parameters {
						param := paramRef.GetResolvedObject()
						if param != nil {
							opParam := templates.ParamData{
								Name: param.Name,
								In:   string(param.In),
							}

							// Handle optional fields
							if param.Required != nil {
								opParam.Required = *param.Required
							}
							if param.Description != nil {
								opParam.Description = *param.Description
							}

							// Extract parameter type from schema
							if param.Schema != nil {
								schema := param.Schema.GetResolvedObject()
								if schema != nil && schema.Left != nil {
									jsonSchema := NewJSONSchema(schema.Left)
									opParam.Type = jsonSchema.SchemaToGoType(schema.Left)
									if schema.Left.Format != nil {
										opParam.Format = *schema.Left.Format
									}
								}
							}

							// Categorize by location
							switch param.In {
							case "path":
								operation.PathParams = append(operation.PathParams, opParam)
							case "query":
								operation.QueryParams = append(operation.QueryParams, opParam)
							case "header":
								operation.HeaderParams = append(operation.HeaderParams, opParam)
							}
						}
					}
				}

				// Check for request body
				if op.RequestBody != nil {
					reqBody := op.RequestBody.GetResolvedObject()
					if reqBody != nil {
						operation.HasRequestBody = true
						if reqBody.Required != nil {
							operation.RequestBodyRequired = *reqBody.Required
						}

						// Extract request body schema name if available
						if reqBody.Content != nil {
							if jsonContent := reqBody.Content.GetOrZero("application/json"); jsonContent != nil {
								if jsonContent.Schema != nil {
									// Try to get reference first
									ref := jsonContent.Schema.GetReference()
									if ref.String() != "" {
										// Extract schema name from ref like "#/components/schemas/User"
										parts := strings.Split(ref.String(), "/")
										if len(parts) > 0 {
											operation.RequestBodySchema = parts[len(parts)-1]
										}
									} else {
										// For inline schemas, generate a type name based on operation ID
										// e.g., updateUser -> UpdateUserRequestBody
										if operation.OperationID != "" {
											// Convert operationID to PascalCase and add RequestBody suffix
											baseOperationName := templates.PascalCase(operation.OperationID)
											operation.RequestBodySchema = baseOperationName + "RequestBody"
										}
									}
								}
							}
						}
					}
				}

				// Extract security requirements
				// First check if operation has its own security requirements
				securityRequirements := op.Security

				// If no operation-level security, use global security
				if securityRequirements == nil && o.Security != nil {
					securityRequirements = o.Security
				}

				for _, secReq := range securityRequirements {
					if secReq != nil {
						for secSchemeName, scopes := range secReq.All() {
							// Look up the security scheme in the global security definitions
							if o.Components != nil && o.Components.SecuritySchemes != nil {
								secSchemeRef := o.Components.SecuritySchemes.GetOrZero(
									secSchemeName,
								)
								if secSchemeRef != nil {
									secScheme := secSchemeRef.GetResolvedObject()
									if secScheme != nil {
										secReqData := templates.SecurityRequirement{
											Name:   secSchemeName,
											Type:   string(secScheme.Type),
											Scopes: scopes,
										}

										if secScheme.Type == "http" &&
											secScheme.Scheme != nil &&
											*secScheme.Scheme == "bearer" {
											operation.HasBearerAuth = true
											secReqData.Scheme = *secScheme.Scheme
										}

										if secScheme.Type == "apiKey" &&
											secScheme.In != nil &&
											*secScheme.In == "cookie" {
											operation.HasCookieAuth = true
											// secReqData.Scheme = *secScheme.Scheme
										}

										operation.Security = append(
											operation.Security,
											secReqData,
										)

									}
								}
							}
						}
					}
				}

				// Extract responses
				if op.Responses != nil {
					for statusCode := range op.Responses.Keys() {
						responseRef := op.Responses.GetOrZero(statusCode)
						if responseRef != nil {
							response := responseRef.GetResolvedObject()
							if response != nil {
								opResp := templates.OperationResponse{
									StatusCode:  statusCode,
									Description: response.Description,
								}

								// Extract response schema if available
								if response.Content != nil {
									if jsonContent := response.Content.GetOrZero("application/json"); jsonContent != nil {
										if jsonContent.Schema != nil {
											// Try to get reference
											ref := jsonContent.Schema.GetReference()
											// Extract schema name from ref like "#/components/schemas/User"
											parts := strings.Split(ref.String(), "/")
											if len(parts) > 0 {
												opResp.Schema = parts[len(parts)-1]
											}
										}
									}
								}

								// Check if it's a success response (2xx)
								if strings.HasPrefix(statusCode, "2") {
									opResp.IsSuccess = true
									if operation.SuccessResponse == nil {
										operation.SuccessResponse = &opResp
										// Set ResponseType from first success response with a schema
										if opResp.Schema != "" && operation.ResponseType == "" {
											operation.ResponseType = opResp.Schema
										}
									}
								} else {
									operation.ErrorResponses = append(operation.ErrorResponses, opResp)
								}

								operation.Responses = append(operation.Responses, opResp)
							}
						}
					}
				}

				operation.Name = operation.OperationID
				operations = append(operations, operation)

				return nil
			},
		})
	}

	return operations
}

func (o *OpenAPISchema) ExtractSchemas() []*JSONSchema {
	if o == nil || o.Components == nil || o.Components.Schemas == nil {
		return nil
	}

	var schemas []*JSONSchema
	for name, schemaRef := range o.Components.Schemas.All() {
		resolvedSchema := schemaRef.GetResolvedObject()
		if resolvedSchema != nil && resolvedSchema.Left != nil {
			schema := resolvedSchema.Left
			jsonSchema := NewJSONSchema(schema)
			jsonSchema.Name = name

			// Extract extensions if available
			if schema.Extensions != nil {
				// Extract tags
				if ext := schema.Extensions.GetOrZero("x-tags"); ext != nil {
					var tags []interface{}
					if err := ext.Decode(&tags); err == nil {
						for _, tag := range tags {
							if tagStr, ok := tag.(string); ok {
								jsonSchema.Tags = append(jsonSchema.Tags, tagStr)
							}
						}
					}
				}

				// Store x-codegen extension as raw data
				if codegenExt := schema.Extensions.GetOrZero("x-codegen"); codegenExt != nil {
					var rawCodegen any
					if err := codegenExt.Decode(&rawCodegen); err == nil {
						jsonSchema.Extensions["x-codegen"] = rawCodegen
					}
				}
			}

			schemas = append(schemas, jsonSchema)
		}
	}

	return schemas
}

// GroupOperationsByDomain groups operations by their domain based on tags.
func (o *OpenAPISchema) GroupOperationsByDomain() map[string][]templates.OperationData {
	operations := o.ExtractOperations()

	domainOps := make(map[string][]templates.OperationData)
	for _, op := range operations {
		// Use the first tag as domain
		if len(op.Tags) > 0 {
			domain := strings.ToLower(op.Tags[0])
			domainOps[domain] = append(domainOps[domain], op)
		}
	}

	// Sort operations within each domain for consistent output
	for domain, ops := range domainOps {
		sort.Slice(ops, func(i, j int) bool {
			if ops[i].Path != ops[j].Path {
				return ops[i].Path < ops[j].Path
			}
			return ops[i].Method < ops[j].Method
		})
		domainOps[domain] = ops
	}

	return domainOps
}

// FilterOperationsForDomain filters operations for a specific domain.
func (o *OpenAPISchema) FilterOperationsForDomain(
	operations []templates.OperationData,
	domain string,
) []templates.OperationData {
	var filtered []templates.OperationData

	for _, op := range operations {
		// Check if any tag matches the domain
		for _, tag := range op.Tags {
			normalized := strings.ToLower(tag)
			// Handle both singular and plural forms
			if normalized == domain || normalized == domain+"s" ||
				strings.TrimSuffix(normalized, "s") == domain {
				filtered = append(filtered, op)
				break
			}
		}
	}

	return filtered
}

// GroupSchemasByDomain groups schemas by their x-codegen domain.
func (o *OpenAPISchema) GroupSchemasByDomain() map[string][]*JSONSchema {

	schemas := o.ExtractSchemas()
	domainSchemas := make(map[string][]*JSONSchema)
	for _, schema := range schemas {
		// Try to get domain from x-codegen extension
		domain := ""
		if xCodegen, ok := schema.GetExtension("x-codegen"); ok {
			// Try to extract domain from the raw extension data
			if codegenMap, ok := xCodegen.(map[string]any); ok {
				if d, ok := codegenMap["domain"].(string); ok && d != "" {
					domain = d
				}
			}
		}

		// Use schema name as fallback domain
		if domain == "" {
			domain = schema.Name
		}

		// Skip if still no domain
		if domain == "" {
			continue
		}

		domainSchemas[domain] = append(domainSchemas[domain], schema)
	}

	// Sort schemas within each domain for consistent output
	for domain, schs := range domainSchemas {
		sort.Slice(schs, func(i, j int) bool {
			return schs[i].Name < schs[j].Name
		})
		domainSchemas[domain] = schs
	}

	return domainSchemas
}

// GetSortedDomains returns a sorted list of domain names from the grouped schemas.
func (o *OpenAPISchema) GetSortedDomains(domainSchemas map[string][]*JSONSchema) []string {
	var domains []string
	for domain := range domainSchemas {
		domains = append(domains, domain)
	}
	sort.Strings(domains)
	return domains
}
