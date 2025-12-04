package parsers

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/bundler"
	"github.com/pb33f/libopenapi/datamodel"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

const (
	boolTrueString = "true"
)

// OpenAPIParser wraps an OpenAPI document and provides parsing utilities
type OpenAPIParser struct {
	doc             *v3.Document
	lintEnabled     bool
	specBytes       []byte
	basePath        string
	enabledIncludes []string       // Names of enabled x-include-* extensions
	merger          *IncludeMerger // Optional merger for processing x-include-* extensions
}

// NewOpenAPIParser creates a new OpenAPIParser instance
func NewOpenAPIParser() *OpenAPIParser {
	return &OpenAPIParser{}
}

// WithLinting enables strict linting that will block parsing on any violations
func (p *OpenAPIParser) WithLinting() *OpenAPIParser {
	p.lintEnabled = true
	return p
}

// WithIncludeMerger sets the IncludeMerger to use for processing x-include-* extensions.
// If not set, a default empty merger will be used (no includes registered).
func (p *OpenAPIParser) WithIncludeMerger(m *IncludeMerger) *OpenAPIParser {
	p.merger = m
	return p
}

// Parse parses an OpenAPI specification from bytes and returns the document
func (p *OpenAPIParser) Parse(data []byte) (*v3.Document, error) {
	// Store spec bytes for potential linting
	p.specBytes = data

	// Run linting if enabled (blocking on any violations)
	if p.lintEnabled {
		// Use basePath if available for better reference resolution
		if p.basePath != "" {
			if err := p.LintWithBasePath(data, p.basePath); err != nil {
				return nil, err
			}
		} else {
			if err := p.Lint(data); err != nil {
				return nil, err
			}
		}
	}

	config := &datamodel.DocumentConfiguration{
		AllowFileReferences:   true,
		AllowRemoteReferences: true,
		BasePath:              p.basePath,
	}

	doc, err := libopenapi.NewDocumentWithConfiguration(data, config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}

	v3Model, err := doc.BuildV3Model()
	if err != nil {
		return nil, fmt.Errorf("failed to build v3 model: %w", err)
	}

	p.doc = &v3Model.Model
	return &v3Model.Model, nil
}

// ParseFile reads and parses an OpenAPI specification from a file.
// It automatically processes any x-include-* extensions to merge in referenced specs.
func (p *OpenAPIParser) ParseFile(path string) (*v3.Document, error) {
	// Process x-include-* extensions first
	merger := p.merger
	if merger == nil {
		merger = NewIncludeMerger() // empty, no includes registered
	}
	mergedSpecPath, cleanup, enabledIncludes, err := merger.ProcessIncludes(path)
	if err != nil {
		return nil, fmt.Errorf("failed to process includes: %w", err)
	}
	defer cleanup()

	// Store enabled includes for later use
	p.enabledIncludes = enabledIncludes

	data, err := os.ReadFile(mergedSpecPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	// Store base path for reference resolution during linting
	p.basePath = filepath.Dir(mergedSpecPath)
	return p.Parse(data)
}

// Bundle bundles an OpenAPI specification with external references into a single document.
// It first processes any x-include-* extensions to merge in referenced specs, then bundles
// all external $ref references into a single document.
func (p *OpenAPIParser) Bundle(specPath, outputPath string, orvalFix bool) error {
	// Process x-include-* extensions first
	merger := p.merger
	if merger == nil {
		merger = NewIncludeMerger() // empty, no includes registered
	}
	mergedSpecPath, cleanup, _, err := merger.ProcessIncludes(specPath)
	if err != nil {
		return fmt.Errorf("failed to process includes: %w", err)
	}
	defer cleanup()

	// Read the (potentially merged) spec file
	specBytes, err := os.ReadFile(mergedSpecPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Calculate base path from the spec's directory
	basePath := filepath.Dir(mergedSpecPath)

	// Run linting if enabled (blocking on any violations)
	if p.lintEnabled {
		if err := p.LintWithBasePath(specBytes, basePath); err != nil {
			return err
		}
	}

	// Configure bundler
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))
	config := &datamodel.DocumentConfiguration{
		BasePath:                basePath,
		ExtractRefsSequentially: true,
		AllowRemoteReferences:   true,
		Logger:                  logger,
		AllowFileReferences:     true,
	}

	// Bundle using composed mode
	bundled, err := bundler.BundleBytesComposed(specBytes, config, &bundler.BundleCompositionConfig{
		Delimiter: "__",
	})
	if err != nil {
		return fmt.Errorf("failed to bundle spec: %w", err)
	}

	// Remove duplicate component entries created by composed bundling.
	// The bundler creates both "Foo: $ref: #/components/.../Foo__type" and "Foo__type: {...}".
	// We keep only the actual definitions (with suffix) and rename them to remove the suffix.
	bundled, err = cleanupComposedBundle(bundled)
	if err != nil {
		return fmt.Errorf("failed to cleanup bundled spec: %w", err)
	}

	// Inject default responses (e.g., 500 for all operations)
	bundled, err = injectDefaults(bundled)
	if err != nil {
		return fmt.Errorf("failed to inject defaults: %w", err)
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

// Extract extracts operations and component schemas from the OpenAPI spec
func (p *OpenAPIParser) Extract() (*SpecDef, error) {
	operations, err := p.extractOperations()
	if err != nil {
		return nil, fmt.Errorf("failed to extract operations: %w", err)
	}

	schemas, err := p.extractComponentSchemas()
	if err != nil {
		return nil, fmt.Errorf("failed to extract component schemas: %w", err)
	}

	// Extract project name from x-project-name extension
	projectName := p.extractProjectName()

	return &SpecDef{
		Operations:  operations,
		Schemas:     schemas,
		Document:    p.doc,
		ProjectName: projectName,
	}, nil
}

// extractProjectName extracts the project name from the x-project-name extension
func (p *OpenAPIParser) extractProjectName() string {
	if p.doc == nil || p.doc.Extensions == nil {
		return ""
	}

	if ext, ok := p.doc.Extensions.Get("x-project-name"); ok {
		var projectName string
		if err := ext.Decode(&projectName); err == nil {
			return projectName
		}
	}

	return ""
}

// extractOperations extracts all operations from the OpenAPI spec
func (p *OpenAPIParser) extractOperations() ([]OperationDef, error) {
	if p.doc == nil {
		return nil, fmt.Errorf("document not set")
	}

	var operations []OperationDef

	// Iterate through all paths
	for pathPair := p.doc.Paths.PathItems.First(); pathPair != nil; pathPair = pathPair.Next() {
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

			requestBody, err := p.extractRequestBody(op)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to extract request body for %s %s: %w",
					method,
					path,
					err,
				)
			}

			responses, err := p.extractResponses(op)
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
				Parameters:            p.extractParameters(op, pathItem),
				Security:              p.extractSecurityRequirements(op),
				Responses:             responses,
				XCodegenCustomHandler: p.extractXCodegenCustomHandler(op),
				XCodegenRepository:    p.extractXCodegenRepository(op),
				XInternal:             p.extractXInternal(op),
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

// extractXCodegenCustomHandler checks if the operation has the x-codegen-custom-handler extension set to true
func (p *OpenAPIParser) extractXCodegenCustomHandler(op *v3.Operation) bool {
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
func (p *OpenAPIParser) extractXCodegenRepository(op *v3.Operation) string {
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
func (p *OpenAPIParser) extractXInternal(op *v3.Operation) string {
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

// extractParameters extracts all parameters from an operation and path item
func (p *OpenAPIParser) extractParameters(op *v3.Operation, pathItem *v3.PathItem) []ParamDef {
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

// extractSecurityRequirements extracts security requirements from an operation
func (p *OpenAPIParser) extractSecurityRequirements(op *v3.Operation) []SecurityDef {
	var securityDefs []SecurityDef

	securityRequirements := op.Security

	// Important: If operation explicitly sets empty security array, it means no auth
	// Only use global security if the operation has no security field at all (nil)
	if securityRequirements == nil && p.doc != nil && p.doc.Security != nil {
		securityRequirements = p.doc.Security
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
			if p.doc != nil && p.doc.Components != nil && p.doc.Components.SecuritySchemes != nil {
				if secScheme, ok := p.doc.Components.SecuritySchemes.Get(secSchemeName); ok {
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

// extractResponses extracts all responses from an operation
func (p *OpenAPIParser) extractResponses(op *v3.Operation) ([]ResponseDef, error) {
	var responses []ResponseDef
	if op.Responses == nil || op.Responses.Codes == nil {
		return responses, nil
	}

	for statusCode, response := range op.Responses.Codes.FromNewest() {
		if response != nil {
			// Initialize response definition with basic info
			responseDef := ResponseDef{
				StatusCode: statusCode,
				Headers:    make(map[string]*SchemaDef),
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
							jsonParser := NewJSONSchemaParser(p.doc)
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
							responseDef.SchemaDef = processed
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
							// Create a SchemaDef for the header
							headerDef := &SchemaDef{
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
							headerDef.GoType = SchemaToGoType(schema, nil, "")

							responseDef.Headers[headerName] = headerDef
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

	// Sort responses by status code (success responses first, then errors in numerical order)
	sort.Slice(responses, func(i, j int) bool {
		iCode, _ := strconv.Atoi(responses[i].StatusCode)
		jCode, _ := strconv.Atoi(responses[j].StatusCode)
		return iCode < jCode
	})

	return responses, nil
}

// extractRequestBody checks if an operation has a required request body and extracts its schema
func (p *OpenAPIParser) extractRequestBody(op *v3.Operation) (*RequestBodyDef, error) {
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
					jsonParser := NewJSONSchemaParser(p.doc)
					processed, err := jsonParser.ParseBase(schema)
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
func (p *OpenAPIParser) extractComponentSchemas() ([]*SchemaDef, error) {
	if p.doc == nil {
		return nil, fmt.Errorf("document not set")
	}

	if p.doc.Components == nil || p.doc.Components.Schemas == nil {
		return nil, nil
	}

	// Initialize results slice
	results := []*SchemaDef{}

	// Track processed schemas by title to avoid duplicates from ref resolution
	processedTitles := make(map[string]bool)

	// Process each schema
	for schemaPair := p.doc.Components.Schemas.First(); schemaPair != nil; schemaPair = schemaPair.Next() {
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

		jsonParser := NewJSONSchemaParser(p.doc)
		processed, err := jsonParser.ParseBase(schema)
		if err != nil {
			return nil, fmt.Errorf("failed to process schema %s: %w", schemaName, err)
		}

		results = append(results, processed)
	}

	return results, nil
}
