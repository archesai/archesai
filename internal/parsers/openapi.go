package parsers

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/bundler"
	"github.com/pb33f/libopenapi/datamodel"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

// OpenAPIParser wraps an OpenAPI document and provides parsing utilities
type OpenAPIParser struct {
	openAPIDoc *v3.Document
}

// NewOpenAPIParser creates a new OpenAPIParser instance
func NewOpenAPIParser() *OpenAPIParser {
	return &OpenAPIParser{}
}

// Parse parses an OpenAPI specification file and returns the document
func (p *OpenAPIParser) Parse(path string) (*v3.Document, error) {
	// Read the spec file
	specBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	config := &datamodel.DocumentConfiguration{
		AllowFileReferences:   true,
		AllowRemoteReferences: true,
	}

	doc, err := libopenapi.NewDocumentWithConfiguration(specBytes, config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}

	v3Model, err := doc.BuildV3Model()
	if err != nil {
		return nil, fmt.Errorf("failed to build v3 model: %w", err)
	}

	// Store the parsed document
	p.openAPIDoc = &v3Model.Model
	return &v3Model.Model, nil
}

// Bundle bundles an OpenAPI specification with external references into a single document
func (p *OpenAPIParser) Bundle(specPath, outputPath string, orvalFix bool) error {
	// Read the spec file
	specBytes, err := os.ReadFile(specPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Calculate base path from input spec directory
	basePath := filepath.Dir(specPath)

	// Configure bundler
	config := &datamodel.DocumentConfiguration{
		BasePath:                basePath,
		ExtractRefsSequentially: true,
		AllowRemoteReferences:   true,
		AllowFileReferences:     true,
	}

	// Bundle using composed mode
	bundled, err := bundler.BundleBytesComposed(specBytes, config, &bundler.BundleCompositionConfig{
		Delimiter: "__",
	})
	if err != nil {
		return fmt.Errorf("failed to bundle spec: %w", err)
	}

	// Write bundled output
	err = os.WriteFile(outputPath, bundled, 0644)
	if err != nil {
		return fmt.Errorf("failed to write bundled file: %w", err)
	}

	// Resolve pathItems if orval fix is enabled
	if orvalFix {
		if err := resolvePathItems(outputPath); err != nil {
			return fmt.Errorf("failed to resolve pathItems: %w", err)
		}
	}

	return nil
}

// ExtractOperations extracts all operations from the OpenAPI spec
func ExtractOperations(doc *v3.Document) ([]OperationDef, error) {
	if doc == nil {
		return nil, nil
	}

	var operations []OperationDef

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

			requestBody, err := ExtractRequestBody(op, doc)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to extract request body for %s %s: %w",
					method,
					path,
					err,
				)
			}

			responses, err := ExtractResponses(op, doc)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to extract responses for %s %s: %w",
					method,
					path,
					err,
				)
			}

			operationDef := OperationDef{
				Method:                strings.ToUpper(method),
				Path:                  path,
				ID:                    op.OperationId,
				Description:           op.Summary,
				Tag:                   op.Tags[0],
				Parameters:            ExtractParameters(op, pathItem),
				Security:              ExtractSecurityRequirements(op, doc),
				Responses:             responses,
				XCodegenCustomHandler: ExtractXCodegenCustomHandler(op),
				XCodegenRepository:    ExtractXCodegenRepository(op),
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

// ExtractXCodegenCustomHandler checks if the operation has the x-codegen-custom-handler extension set to true
func ExtractXCodegenCustomHandler(op *v3.Operation) bool {
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
			return strVal == "true"
		}
	}
	return false
}

// ExtractXCodegenRepository extracts the repository name from x-codegen-repository extension
func ExtractXCodegenRepository(op *v3.Operation) string {
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

// ExtractParameters extracts all parameters from an operation and path item
func ExtractParameters(op *v3.Operation, pathItem *v3.PathItem) []ParamDef {
	var params []ParamDef

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

		// Create embedded SchemaDef for the parameter
		schemaDef := &SchemaDef{
			Name:        PascalCase(param.Name),
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
				schemaDef.GoType = SchemaToGoType(schema, nil, "")

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

		// Create ParamDef with embedded SchemaDef
		paramDef := ParamDef{
			SchemaDef: schemaDef,
			In:        param.In,
		}
		if param.Style != "" {
			paramDef.Style = param.Style
		}
		if param.Explode != nil {
			paramDef.Explode = *param.Explode
		}

		// Override Required at the ParamDef level since it's different for parameters
		if param.Required != nil && *param.Required {
			paramDef.Required = []string{schemaDef.Name}
		}

		params = append(params, paramDef)
	}

	return params
}

// ExtractSecurityRequirements extracts security requirements from an operation
func ExtractSecurityRequirements(op *v3.Operation, doc *v3.Document) []SecurityDef {
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
						secReqData := SecurityDef{
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

// ExtractResponses extracts all responses from an operation
func ExtractResponses(op *v3.Operation, doc *v3.Document) ([]ResponseDef, error) {
	var responses []ResponseDef
	if op.Responses == nil || op.Responses.Codes == nil {
		return responses, nil
	}

	for responsePair := op.Responses.Codes.First(); responsePair != nil; responsePair = responsePair.Next() {
		statusCode := responsePair.Key()
		response := responsePair.Value()

		if response != nil {
			// Initialize response definition with basic info
			responseDef := ResponseDef{
				StatusCode: statusCode,
			}

			// Try to extract schema if present
			if response.Content != nil {
				if mediaType, ok := response.Content.Get("application/json"); ok {
					if mediaType.Schema != nil {
						schema := mediaType.Schema.Schema()
						if schema != nil {
							responseName := fmt.Sprintf("%s%sResponse", op.OperationId, statusCode)
							jsonParser := NewJSONSchemaParser().WithOpenAPIDoc(doc)
							processed, err := jsonParser.ExtractSchema(schema, &responseName, "")
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
							responseDef.SchemaDef = processed
						}
					}
				}
			}

			// For responses without a schema, create a minimal SchemaDef
			if responseDef.SchemaDef == nil {
				responseDef.SchemaDef = &SchemaDef{
					Description: response.Description,
				}
			}

			// Add the response even if it doesn't have a schema (e.g., error responses using $ref)
			responses = append(responses, responseDef)
		}
	}

	return responses, nil
}

// ExtractRequestBody checks if an operation has a required request body and extracts its schema
func ExtractRequestBody(op *v3.Operation, doc *v3.Document) (*RequestBodyDef, error) {
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
					requestName := fmt.Sprintf("%sRequestBody", op.OperationId)
					jsonParser := NewJSONSchemaParser().WithOpenAPIDoc(doc)
					processed, err := jsonParser.ExtractSchema(schema, &requestName, "")
					if err != nil {
						return nil, fmt.Errorf("failed to process request body schema: %w", err)
					}

					required := false
					if rb.Required != nil {
						required = *rb.Required
					}

					return &RequestBodyDef{
						SchemaDef: processed,
						Required:  required,
					}, nil
				}
			}
		}
	}

	return nil, nil
}

// ExtractComponentSchemas processes all schemas from the OpenAPI document
func ExtractComponentSchemas(doc *v3.Document) ([]*SchemaDef, error) {
	results := []*SchemaDef{}

	if doc == nil {
		return results, nil
	}

	if doc.Components == nil || doc.Components.Schemas == nil {
		return results, nil
	}

	// Process each schema
	for schemaPair := doc.Components.Schemas.First(); schemaPair != nil; schemaPair = schemaPair.Next() {
		schemaName := schemaPair.Key()
		schemaRef := schemaPair.Value()

		if schemaRef == nil {
			continue
		}

		// Get the resolved schema object
		schema := schemaRef.Schema()
		if schema == nil {
			continue
		}

		// Determine schema type from x-codegen extension
		schemaType := ""
		if schema.Extensions != nil {
			if ext, ok := schema.Extensions.Get("x-codegen"); ok {
				var xcodegen XCodegenExtension
				err := ext.Decode(&xcodegen)
				if err == nil {
					schemaType = string(xcodegen.GetSchemaType())
				}
			}
		}

		jsonParser := NewJSONSchemaParser().WithOpenAPIDoc(doc)
		processed, err := jsonParser.ExtractSchema(schema, &schemaName, schemaType)
		if err != nil {
			return nil, fmt.Errorf("failed to process schema %s: %w", schemaName, err)
		}

		results = append(results, processed)
	}

	return results, nil
}
