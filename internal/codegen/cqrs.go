package codegen

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/archesai/archesai/internal/parsers"
)

// COMMANDS AND QUERIES

// GenerateCommandQueryHandlers generates command and query handlers
func (g *Generator) GenerateCommandQueryHandlers(
	operations []parsers.OperationDef,
	schemas map[string]*parsers.ProcessedSchema,
) error {
	// Generate command handlers (includes types)
	if err := g.generateCommandHandlers(operations, schemas); err != nil {
		return fmt.Errorf("failed to generate command handlers: %w", err)
	}

	// Generate query handlers (includes types)
	if err := g.generateQueryHandlers(operations, schemas); err != nil {
		return fmt.Errorf("failed to generate query handlers: %w", err)
	}

	return nil
}

// generateCommandHandlers generates individual command handler files for each operation
func (g *Generator) generateCommandHandlers(
	operations []parsers.OperationDef,
	schemas map[string]*parsers.ProcessedSchema,
) error {
	tmpl, ok := g.templates["single_command_handler.tmpl"]
	if !ok {
		// Template doesn't exist yet, skip
		return nil
	}

	// Build a map of domains that have commands configuration
	domainsWithCommands := make(map[string]bool)
	for name, processed := range schemas {
		if processed.XCodegen != nil && processed.XCodegen.Commands != nil {
			// Mark this domain as having commands
			// Use the plural form of the name as the domain (e.g., "Tool" -> "tools")
			domainsWithCommands[strings.ToLower(name+"s")] = true
		}
	}

	// Group operations by their domain/tag
	operationsByDomain := make(map[string][]parsers.OperationDef)
	for _, op := range operations {
		if len(op.Tags) > 0 {
			domain := strings.ToLower(op.Tags[0])
			operationsByDomain[domain] = append(operationsByDomain[domain], op)
		}
	}

	// For each domain, generate command handlers for write operations
	for domain, operations := range operationsByDomain {
		// Only generate commands for domains that have x-codegen configuration
		if !domainsWithCommands[domain] {
			continue
		}
		for _, op := range operations {
			// Only generate command handlers for write operations
			if op.Method == "POST" || op.Method == "PUT" || op.Method == "PATCH" ||
				op.Method == "DELETE" {
				// Determine command type based on operation
				var commandType string
				switch {
				case strings.HasPrefix(op.OperationID, "create"):
					commandType = "Create"
				case strings.HasPrefix(op.OperationID, "update"):
					commandType = "Update"
				case strings.HasPrefix(op.OperationID, "delete"):
					commandType = "Delete"
				default:
					// Custom command, use operation ID with title case
					commandType = strings.Title(op.OperationID)
				}

				// Get the entity name from the schema if available
				entityName := Singularize(strings.Title(domain))
				entityNameLower := strings.ToLower(entityName)

				// Create template data
				data := map[string]interface{}{
					"Package":         domain,
					"CommandType":     commandType,
					"EntityName":      entityName,
					"EntityNameLower": entityNameLower,
				}

				// Generate the command handler file
				outputPath := filepath.Join(
					"internal/application/commands",
					domain,
					fmt.Sprintf("%s_%s.gen.go", strings.ToLower(commandType), entityNameLower),
				)

				// Write the handler file
				if err := g.filewriter.WriteTemplate(outputPath, tmpl, data); err != nil {
					return fmt.Errorf(
						"failed to generate command handler for %s: %w",
						op.OperationID,
						err,
					)
				}
			}
		}
	}

	return nil
}

// generateQueryHandlers generates individual query handler files for each operation
func (g *Generator) generateQueryHandlers(
	operations []parsers.OperationDef,
	schemas map[string]*parsers.ProcessedSchema,
) error {
	tmpl, ok := g.templates["single_query_handler.tmpl"]
	if !ok {
		// Template doesn't exist yet, skip
		return nil
	}

	// Build a map of domains that have queries configuration
	domainsWithQueries := make(map[string]bool)
	for name, processed := range schemas {
		if processed.XCodegen != nil && processed.XCodegen.Queries != nil {
			// Mark this domain as having queries
			// Use the plural form of the name as the domain (e.g., "Tool" -> "tools")
			domainsWithQueries[strings.ToLower(name+"s")] = true
		}
	}

	// Group operations by their domain/tag
	operationsByDomain := make(map[string][]parsers.OperationDef)
	for _, op := range operations {
		if len(op.Tags) > 0 {
			domain := strings.ToLower(op.Tags[0])
			operationsByDomain[domain] = append(operationsByDomain[domain], op)
		}
	}

	// For each domain, generate query handlers for read operations
	for domain, operations := range operationsByDomain {
		// Only generate queries for domains that have x-codegen configuration
		if !domainsWithQueries[domain] {
			continue
		}
		for _, op := range operations {
			// Only generate query handlers for read operations
			if op.Method == "GET" {
				// Determine query type based on operation
				var queryType string
				var usesPluralName bool
				switch {
				case strings.HasPrefix(op.OperationID, "list"):
					queryType = "List"
					usesPluralName = true
				case strings.HasPrefix(op.OperationID, "get"):
					queryType = "Get"
				case strings.HasPrefix(op.OperationID, "search"):
					queryType = "Search"
					usesPluralName = true
				default:
					// Custom query, use operation ID with title case
					queryType = strings.Title(op.OperationID)
				}

				// Get the entity name from the schema if available
				entityName := Singularize(strings.Title(domain))
				entityNameLower := strings.ToLower(entityName)
				entityNamePlural := Pluralize(entityName)

				// Create template data
				data := map[string]interface{}{
					"Package":          domain,
					"QueryType":        queryType,
					"EntityName":       entityName,
					"EntityNameLower":  entityNameLower,
					"EntityNamePlural": entityNamePlural,
				}

				// Generate the query handler file
				var fileName string
				if usesPluralName {
					fileName = fmt.Sprintf(
						"%s_%s.gen.go",
						strings.ToLower(queryType),
						strings.ToLower(entityNamePlural),
					)
				} else {
					fileName = fmt.Sprintf("%s_%s.gen.go", strings.ToLower(queryType), entityNameLower)
				}

				outputPath := filepath.Join(
					"internal/application/queries",
					domain,
					fileName,
				)

				// Write the handler file
				if err := g.filewriter.WriteTemplate(outputPath, tmpl, data); err != nil {
					return fmt.Errorf(
						"failed to generate query handler for %s: %w",
						op.OperationID,
						err,
					)
				}
			}
		}
	}

	return nil
}
