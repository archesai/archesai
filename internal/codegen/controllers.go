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

	// List of domains that have command/query handlers
	cqrsDomains := map[string]bool{
		"accounts":      true,
		"apikeys":       true,
		"artifacts":     true,
		"invitations":   true,
		"labels":        true,
		"members":       true,
		"organizations": true,
		"pipelines":     true,
		"runs":          true,
		"tools":         true,
		"users":         true,
	}

	// Generate a handler file for each domain
	for domain, operations := range operationsByDomain {
		// Skip empty domains
		if domain == "" {
			continue
		}

		// Skip domains that don't have CQRS handlers for now
		// TODO: Generate simpler controllers for non-CQRS domains
		if !cqrsDomains[strings.ToLower(domain)] {
			continue
		}

		// Capitalize first letter of domain for title case
		domainTitle := strings.ToUpper(domain[:1]) + domain[1:]

		// Get singular form for entity name (e.g., "Labels" -> "Label")
		domainSingular := Singularize(domainTitle)

		// Track which imports are needed for this handler
		importsNeeded := make(map[string]bool)

		// Process operations to split parameters by type
		processedOps := make([]map[string]interface{}, 0, len(operations))
		for _, op := range operations {
			var pathParams []parsers.ParamDef
			var queryParams []parsers.ParamDef
			var headerParams []parsers.ParamDef

			for _, param := range op.Parameters {
				switch param.In {
				case "path":
					pathParams = append(pathParams, param)
				case "query":
					queryParams = append(queryParams, param)
				case "header":
					headerParams = append(headerParams, param)
				}
			}

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

			// Note: Request body handling removed - operations no longer have inline request bodies
			// They should be defined as schemas in components instead

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
				"Responses":           op.Responses,
				"Security":            op.Security,
				"ResponseType":        responseType,
				"ResponsePackage":     responsePackage,
				"SuccessResponse":     successResponse,
			})
		}

		// Build imports list - no longer needed as we have standard imports in template
		var imports []map[string]string

		data := map[string]interface{}{
			"Package":             "controllers",
			"Domain":              domainTitle,                     // e.g., "Labels" (as it comes from tags)
			"DomainSingular":      domainSingular,                  // e.g., "Label" (singular form)
			"DomainLower":         strings.ToLower(domain),         // e.g., "labels"
			"DomainSingularLower": strings.ToLower(domainSingular), // e.g., "label"
			"Operations":          processedOps,
			"Imports":             imports,
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
