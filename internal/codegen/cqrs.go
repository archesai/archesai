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

	// Build a map of schemas that are entities (not value objects)
	entitySchemas := make(map[string]bool)
	for name, processed := range schemas {
		if processed.XCodegen != nil && processed.XCodegen.SchemaType == "entity" {
			// Use both singular and plural forms of the name as the domain
			// (e.g., "Tool" -> "tools", "Auth" -> "auth")
			entitySchemas[strings.ToLower(name)] = true     // Singular form
			entitySchemas[strings.ToLower(name+"s")] = true // Plural form

			// Special case: Session entity is used for Auth domain
			if strings.ToLower(name) == schemaNameSession {
				entitySchemas["auth"] = true
			}
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
		// Only generate commands for entity domains
		if !entitySchemas[domain] {
			continue
		}
		for _, op := range operations {
			// Only generate command handlers for write operations (POST, PUT, PATCH, DELETE)
			if op.Method == httpMethodGet {
				continue // Skip read operations
			}

			// Use the operationId directly for naming
			commandName := PascalCase(op.OperationID) + "Command"
			fileName := SnakeCase(op.OperationID)

			// Extract command type from operationId for template compatibility
			var commandType string
			entityName := Singularize(Title(domain))
			entityNameLower := strings.ToLower(entityName)

			// For naming the command types
			commandEntityName := entityName

			// Check if it's a standard CRUD operation
			standardCreate := op.OperationID == "create"+entityName
			standardUpdate := op.OperationID == "update"+entityName
			standardDelete := op.OperationID == "delete"+entityName

			isStandardCRUD := false
			if standardCreate {
				commandType = "Create"
				isStandardCRUD = true
			} else if standardUpdate {
				commandType = "Update"
				isStandardCRUD = true
			} else if standardDelete {
				commandType = "Delete"
				isStandardCRUD = true
			} else {
				// For non-standard operations, use the full operationId as type
				commandType = PascalCase(op.OperationID)
			}

			// Automatically set CustomHandler for non-standard operations if not explicitly set
			customHandler := op.CustomHandler
			if !isStandardCRUD && !customHandler {
				customHandler = true
			}

			// For non-standard operations, clear entity name to avoid redundant suffixes
			// Standard CRUD operations keep the entity name to avoid conflicts
			if !isStandardCRUD {
				commandEntityName = ""
			}

			// Check for authentication requirements
			requiresAuth := false
			for _, sec := range op.Security {
				if sec.Name == "bearerAuth" || sec.Name == "sessionCookie" {
					requiresAuth = true
					break
				}
			}

			// Create template data
			data := map[string]interface{}{
				"Package":           domain,
				"CommandName":       commandName,
				"CommandType":       commandType,
				"EntityName":        entityName,        // Keep original for repository and return types
				"EntityNameLower":   entityNameLower,   // Keep original for descriptions
				"CommandEntityName": commandEntityName, // Use for command type names
				"Operation":         op,
				"RequestBody":       op.RequestBodySchema,
				"PathParams":        op.PathParams,
				"QueryParams":       op.QueryParams,
				"RequiresAuth":      requiresAuth,
				"CustomHandler":     customHandler, // Pass custom handler flag
			}

			// Generate the command handler file
			outputPath := filepath.Join(
				"internal/application/commands",
				domain,
				fmt.Sprintf("%s.gen.go", fileName),
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

	// Build a map of schemas that are entities (not value objects)
	entitySchemas := make(map[string]bool)
	for name, processed := range schemas {
		if processed.XCodegen != nil && processed.XCodegen.SchemaType == "entity" {
			// Use both singular and plural forms of the name as the domain
			// (e.g., "Tool" -> "tools", "Auth" -> "auth")
			entitySchemas[strings.ToLower(name)] = true     // Singular form
			entitySchemas[strings.ToLower(name+"s")] = true // Plural form

			// Special case: Session entity is used for Auth domain
			if strings.ToLower(name) == schemaNameSession {
				entitySchemas["auth"] = true
			}
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
		// Only generate queries for entity domains
		if !entitySchemas[domain] {
			continue
		}
		for _, op := range operations {
			// Only generate query handlers for read operations (GET)
			if op.Method != "GET" {
				continue // Skip write operations
			}

			// Use the operationId directly for naming
			queryName := PascalCase(op.OperationID) + "Query"
			fileName := SnakeCase(op.OperationID)

			// Extract query type from operationId for template compatibility
			var queryType string
			entityName := Singularize(Title(domain))
			entityNameLower := strings.ToLower(entityName)
			entityNamePlural := Pluralize(entityName)

			// Special handling for auth domain entity names
			if domain == "auth" {
				// Determine the actual entity based on the operation
				if strings.Contains(strings.ToLower(op.OperationID), "session") {
					entityName = "Session"
					entityNameLower = "session"
					entityNamePlural = "Sessions"
				} else if strings.Contains(strings.ToLower(op.OperationID), "account") {
					entityName = "Account"
					entityNameLower = "account"
					entityNamePlural = "Accounts"
				} else {
					// Default to User for other auth operations
					entityName = "User"
					entityNameLower = "user"
					entityNamePlural = "Users"
				}
			}

			// For naming the query/command types
			queryEntityName := entityName
			queryEntityNamePlural := entityNamePlural

			// Check if it's a standard CRUD operation
			standardGet := op.OperationID == "get"+entityName
			standardList := op.OperationID == "list"+entityNamePlural

			// Special case for auth domain - only clear entity names for non-standard operations
			// Standard CRUD operations keep the entity name to avoid conflicts
			isStandardCRUD := false
			if standardGet {
				queryType = "Get"
				isStandardCRUD = true
			} else if standardList {
				queryType = "List"
				isStandardCRUD = true
			} else if strings.HasPrefix(op.OperationID, "search") {
				queryType = "Search"
			} else {
				// For non-standard operations, use the full operationId as type
				queryType = PascalCase(op.OperationID)
			}

			// Automatically set CustomHandler for non-standard operations if not explicitly set
			customHandler := op.CustomHandler
			if !isStandardCRUD && !customHandler {
				customHandler = true
			}

			// For non-standard operations, clear entity names to avoid redundant suffixes
			// Standard CRUD operations keep the entity name to avoid conflicts
			if !isStandardCRUD {
				queryEntityName = ""
				queryEntityNamePlural = ""
			}

			// Check for authentication requirements
			requiresAuth := false
			for _, sec := range op.Security {
				if sec.Name == "bearerAuth" || sec.Name == "sessionCookie" {
					requiresAuth = true
					break
				}
			}

			// Create template data
			data := map[string]interface{}{
				"Package":               domain,
				"QueryName":             queryName,
				"QueryType":             queryType,
				"EntityName":            entityName,            // Keep original for repository and return types
				"EntityNameLower":       entityNameLower,       // Keep original for descriptions
				"EntityNamePlural":      entityNamePlural,      // Keep original for return types
				"QueryEntityName":       queryEntityName,       // Use for query type names
				"QueryEntityNamePlural": queryEntityNamePlural, // Use for query type names
				"Operation":             op,
				"PathParams":            op.PathParams,
				"QueryParams":           op.QueryParams,
				"RequiresAuth":          requiresAuth,
				"CustomHandler":         customHandler, // Pass custom handler flag
			}

			// Generate the query handler file
			outputPath := filepath.Join(
				"internal/application/queries",
				domain,
				fmt.Sprintf("%s.gen.go", fileName),
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

	return nil
}
