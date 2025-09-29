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
	operationsByDomain, err := g.groupOperationsByDomain(operations)
	if err != nil {
		return err
	}

	// Track schema types to determine imports
	schemaTypeMap := g.buildSchemaTypeMap(schemas)

	// Generate a handler file for each domain
	for domain, operations := range operationsByDomain {
		if domain == "" {
			continue
		}

		if err := g.generateDomainController(domain, operations, schemaTypeMap); err != nil {
			return err
		}
	}

	return nil
}

// groupOperationsByDomain groups operations by their first tag
func (g *Generator) groupOperationsByDomain(
	operations []parsers.OperationDef,
) (map[string][]parsers.OperationDef, error) {
	operationsByDomain := make(map[string][]parsers.OperationDef)
	for _, op := range operations {
		if len(op.Tags) > 0 {
			operationsByDomain[op.Tags[0]] = append(operationsByDomain[op.Tags[0]], op)
		} else {
			return nil, fmt.Errorf("operation %s has no tags defined", op.OperationID)
		}
	}
	return operationsByDomain, nil
}

// buildSchemaTypeMap creates a map of schema names to their package types
func (g *Generator) buildSchemaTypeMap(
	schemas map[string]*parsers.ProcessedSchema,
) map[string]string {
	schemaTypeMap := make(map[string]string)
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
			schemaTypeMap[name] = "dto"
		}
	}
	return schemaTypeMap
}

// generateDomainController generates a controller for a specific domain
func (g *Generator) generateDomainController(
	domain string,
	operations []parsers.OperationDef,
	schemaTypeMap map[string]string,
) error {
	domainTitle := strings.ToUpper(domain[:1]) + domain[1:]
	domainSingular := Singularize(domainTitle)

	// Track which imports are needed for this handler
	importsNeeded := make(map[string]bool)

	// Process operations to add template data
	processedOps := g.processOperations(operations, schemaTypeMap, domainTitle, importsNeeded)

	// Determine which handlers are needed
	crudFlags := g.determineCRUDHandlers(operations)

	data := map[string]interface{}{
		"Package":             "controllers",
		"Domain":              domainTitle,
		"DomainSingular":      domainSingular,
		"DomainLower":         strings.ToLower(domain),
		"DomainSingularLower": strings.ToLower(domainSingular),
		"Operations":          processedOps,
		"Imports":             []map[string]string{}, // no longer needed as we have standard imports in template
		"HasCreate":           crudFlags["create"],
		"HasUpdate":           crudFlags["update"],
		"HasDelete":           crudFlags["delete"],
		"HasGet":              crudFlags["get"],
		"HasList":             crudFlags["list"],
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

	return nil
}

// processOperations processes operations for template data
func (g *Generator) processOperations(
	operations []parsers.OperationDef,
	schemaTypeMap map[string]string,
	domainTitle string,
	importsNeeded map[string]bool,
) []map[string]interface{} {
	processedOps := make([]map[string]interface{}, 0, len(operations))

	for _, op := range operations {
		processedOp := g.processOperation(op, schemaTypeMap, domainTitle, importsNeeded)
		processedOps = append(processedOps, processedOp)
	}

	return processedOps
}

// processOperation processes a single operation for template data
func (g *Generator) processOperation(
	op parsers.OperationDef,
	schemaTypeMap map[string]string,
	domainTitle string,
	importsNeeded map[string]bool,
) map[string]interface{} {
	responseType, responsePackage, successResponse := g.extractResponseInfo(
		op,
		schemaTypeMap,
		importsNeeded,
	)

	// If no response type found, try to determine from domain
	if responseType == "" && domainTitle != "" {
		domainSingular := strings.TrimSuffix(domainTitle, "s")
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

	return map[string]interface{}{
		"Name":                op.Name,
		"GoName":              op.GoName,
		"Method":              op.Method,
		"Path":                op.Path,
		"Description":         op.Description,
		"OperationID":         op.OperationID,
		"Tags":                op.Tags,
		"PathParams":          op.PathParams,
		"QueryParams":         op.QueryParams,
		"HeaderParams":        op.HeaderParams,
		"RequestBodyRequired": op.RequestBodyRequired,
		"RequestBodySchema":   requestBodySchema,
		"HasRequestBody":      op.RequestBodySchema != nil,
		"Responses":           op.Responses,
		"Security":            op.Security,
		"ResponseType":        responseType,
		"ResponsePackage":     responsePackage,
		"SuccessResponse":     successResponse,
	}
}

// extractResponseInfo extracts response type and package information
func (g *Generator) extractResponseInfo(
	op parsers.OperationDef,
	schemaTypeMap map[string]string,
	importsNeeded map[string]bool,
) (string, string, map[string]interface{}) {
	responseType := ""
	responsePackage := ""
	var successResponse map[string]interface{}

	for _, resp := range op.Responses {
		if resp.IsSuccess && resp.Schema != "" {
			responseType = resp.Schema
			if pkg, ok := schemaTypeMap[resp.Schema]; ok {
				importsNeeded[pkg] = true
				responsePackage = pkg
			}
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

	return responseType, responsePackage, successResponse
}

// determineCRUDHandlers determines which CRUD handlers are needed based on operation names
func (g *Generator) determineCRUDHandlers(operations []parsers.OperationDef) map[string]bool {
	flags := map[string]bool{
		"create": false,
		"update": false,
		"delete": false,
		"get":    false,
		"list":   false,
	}

	for _, op := range operations {
		opNameLower := strings.ToLower(op.Name)

		if strings.HasPrefix(opNameLower, "create") {
			flags["create"] = true
		} else if strings.HasPrefix(opNameLower, "update") || strings.HasPrefix(opNameLower, "patch") {
			flags["update"] = true
		} else if strings.HasPrefix(opNameLower, "delete") || strings.HasPrefix(opNameLower, "remove") {
			flags["delete"] = true
		} else if strings.HasPrefix(opNameLower, "list") || strings.HasPrefix(opNameLower, "getall") {
			flags["list"] = true
		} else if strings.HasPrefix(opNameLower, "get") {
			flags["get"] = true
		}
	}

	return flags
}
