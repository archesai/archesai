package codegen

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/archesai/archesai/internal/parsers"
)

// GenerateControllers generates all HTTP controllers grouped by domain
func (g *Generator) GenerateControllers(
	operations []parsers.OperationDef,
	schemas map[string]*parsers.ProcessedSchema,
) error {

	// Group operations by their first tag (domain)
	operationsByDomain := make(map[string][]parsers.OperationDef)
	for _, op := range operations {
		if len(op.Tags) > 0 {
			operationsByDomain[op.Tags[0]] = append(operationsByDomain[op.Tags[0]], op)
		} else {
			return fmt.Errorf("operation %s has no tags defined", op.OperationID)
		}
	}

	// Track schema types to determine imports
	schemaTypeMap := make(map[string]string) // schema name -> package
	for name, processed := range schemas {
		if processed.XCodegen != nil {
			switch processed.XCodegen.SchemaType {
			case schemaTypeEntity:
				schemaTypeMap[name] = "entities"
			case schemaTypeValueObject:
				schemaTypeMap[name] = "valueobjects"
			default:
				schemaTypeMap[name] = "dto"
			}
		} else {
			// No x-codegen means it's a DTO
			schemaTypeMap[name] = "dto"
		}
	}

	// Build a map of schemas that have controller configuration
	schemasWithControllers := make(map[string]bool)
	for name, processed := range schemas {
		if processed.XCodegen != nil && processed.XCodegen.Controller != nil {
			schemasWithControllers[name] = true
		}
	}

	// Generate a handler file for each domain
	for domain, operations := range operationsByDomain {
		// Skip empty domains
		if domain == "" {
			continue
		}

		// Check if any schema has a controller configuration for operations in this domain
		// We generate controllers for all tags that have operations, regardless of x-codegen
		// This allows flexibility in tagging (e.g., Auth tag can handle Session operations)

		// Capitalize first letter of domain for title case
		domainTitle := strings.ToUpper(domain[:1]) + domain[1:]

		// Get singular form for entity name (e.g., "Labels" -> "Label")
		domainSingular := Singularize(domainTitle)

		// Track which imports are needed for this handler
		importsNeeded := make(map[string]bool)

		// Process operations to add template data
		processedOps := make([]map[string]interface{}, 0, len(operations))
		for _, op := range operations {
			// Use the already-split parameters from OperationDef
			pathParams := op.PathParams
			queryParams := op.QueryParams
			headerParams := op.HeaderParams

			// Determine response type and track imports
			responseType := ""
			responsePackage := ""
			var successResponse map[string]interface{}
			for _, resp := range op.Responses {
				if resp.IsSuccess && resp.Schema != "" {
					responseType = resp.Schema
					// Track which package this schema comes from
					if pkg, ok := schemaTypeMap[resp.Schema]; ok {
						importsNeeded[pkg] = true
						responsePackage = pkg
					}
					// Create success response data for template
					successResponse = map[string]interface{}{
						"StatusCode":  resp.StatusCode,
						"Schema":      resp.Schema,
						"Package":     responsePackage,
						"IsArray":     resp.IsArray,
						"Description": resp.Description,
					}
					break
				}
			}

			// If no response type found, try to determine from domain
			if responseType == "" && domain != "" {
				// Try to find the schema for the domain (singular form)
				domainSingular := strings.TrimSuffix(domainTitle, "s") // Simple singularize
				if pkg, ok := schemaTypeMap[domainSingular]; ok {
					responsePackage = pkg
					importsNeeded[pkg] = true
				}
			}

			// Process request body schema if present
			var requestBodySchema map[string]interface{}
			if op.RequestBodySchema != nil {
				requestBodySchema = map[string]interface{}{
					"Name":           op.RequestBodySchema.Name,
					"Fields":         op.RequestBodySchema.Fields,
					"RequiredFields": op.RequestBodySchema.RequiredFields,
				}
			}

			processedOps = append(processedOps, map[string]interface{}{
				"Name":                op.Name,
				"GoName":              op.GoName,
				"Method":              op.Method,
				"Path":                op.Path,
				"Description":         op.Description,
				"OperationID":         op.OperationID,
				"Tags":                op.Tags,
				"PathParams":          pathParams,
				"QueryParams":         queryParams,
				"HeaderParams":        headerParams,
				"RequestBodyRequired": op.RequestBodyRequired,
				"RequestBodySchema":   requestBodySchema,
				"HasRequestBody":      op.RequestBodySchema != nil,
				"Responses":           op.Responses,
				"Security":            op.Security,
				"ResponseType":        responseType,
				"ResponsePackage":     responsePackage,
				"SuccessResponse":     successResponse,
			})
		}

		// Build imports list - no longer needed as we have standard imports in template
		var imports []map[string]string

		// Determine which handlers are needed based on operation names, not HTTP methods
		// This ensures we only generate handlers for actual CRUD operations, not custom ones
		hasCreate := false
		hasUpdate := false
		hasDelete := false
		hasGet := false
		hasList := false

		for _, op := range operations {
			// Use operation name/ID as the primary indicator
			// The template expects operations that start with these prefixes
			opNameLower := strings.ToLower(op.Name)

			// Check for standard CRUD operation patterns in the operation name
			if strings.HasPrefix(opNameLower, "create") {
				hasCreate = true
			} else if strings.HasPrefix(opNameLower, "update") || strings.HasPrefix(opNameLower, "patch") {
				hasUpdate = true
			} else if strings.HasPrefix(opNameLower, "delete") || strings.HasPrefix(opNameLower, "remove") {
				hasDelete = true
			} else if strings.HasPrefix(opNameLower, "list") || strings.HasPrefix(opNameLower, "getall") {
				hasList = true
			} else if strings.HasPrefix(opNameLower, "get") {
				// Only consider it a "get" operation if it's explicitly named that way
				// This avoids treating OAuth callbacks and other GET endpoints as "get" operations
				hasGet = true
			}
			// Custom operations (like oauthAuthorize, requestMagicLink) won't match any of these
			// and will be handled as custom operations in the template
		}

		data := map[string]interface{}{
			"Package":             "controllers",
			"Domain":              domainTitle,                     // e.g., "Labels" (as it comes from tags)
			"DomainSingular":      domainSingular,                  // e.g., "Label" (singular form)
			"DomainLower":         strings.ToLower(domain),         // e.g., "labels"
			"DomainSingularLower": strings.ToLower(domainSingular), // e.g., "label"
			"Operations":          processedOps,
			"Imports":             imports,
			"HasCreate":           hasCreate,
			"HasUpdate":           hasUpdate,
			"HasDelete":           hasDelete,
			"HasGet":              hasGet,
			"HasList":             hasList,
		}

		outputPath := filepath.Join(
			"internal/adapters/http/controllers",
			strings.ToLower(domain)+".gen.go",
		)

		tmpl, ok := g.templates["controller.tmpl"]
		if !ok {
			return fmt.Errorf("controller template not found")
		}

		if err := g.filewriter.WriteTemplate(outputPath, tmpl, data); err != nil {
			return fmt.Errorf("failed to generate handler for %s: %w", domain, err)
		}
	}

	return nil
}
