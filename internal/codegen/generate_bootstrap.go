package codegen

import (
	"fmt"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/archesai/archesai/internal/parsers"
)

const (
	httpMethodGet = "GET"
)

// TagDef represents data for a single domain in the app template
type TagDef struct {
	Name                  string
	HasController         bool
	HasRepository         bool
	HasCommands           bool
	HasQueries            bool
	Commands              []string // e.g., ["Create", "Update", "Delete"]
	Queries               []string // e.g., ["Get", "List"]
	IsPublic              bool
	RequiresOrgMembership bool
	Operations            []parsers.OperationDef // Operations for this domain
}

// GenerateBootstrap generates the main application wiring file
func (g *Generator) GenerateBootstrap(
	schemas []*parsers.SchemaDef,
	operations []parsers.OperationDef,
) error {
	// Collect domain data from schemas with x-codegen
	tags := make([]TagDef, 0)
	tagMap := make(map[string]*TagDef)

	// First, collect all tags that have actual operations
	tagsWithOperations := make(map[string]bool)
	for _, op := range operations {
		tagsWithOperations[op.Tag] = true
	}

	// Process each schema to build domain information
	for _, schema := range schemas {

		// Check if this schema should be included
		// Tags are now singular, so check for the singular form
		hasAPIOperations := tagsWithOperations[schema.Name]
		hasRepository := schema.XCodegenSchemaType == parsers.XCodegenSchemaTypeEntity

		// Include schemas that either:
		// 1. Have API operations tagged with their singular name
		// 2. Have a repository defined (needed for cross-cutting concerns like auth)
		if !hasAPIOperations && !hasRepository {
			continue
		}

		// Get or create domain data
		domain, exists := tagMap[schema.Name]
		if !exists {
			domain = &TagDef{
				Name:       schema.Name,
				Commands:   []string{},
				Queries:    []string{},
				Operations: []parsers.OperationDef{},
			}
			tagMap[schema.Name] = domain
		}

		// Update domain based on x-codegen settings
		// Controllers are only created for schemas with operations
		if hasAPIOperations {
			domain.HasController = true
		}

		// Determine route protection level based on domain schema.Name
		// This is a simplified logic - you may need to adjust based on your needs
		lowerName := strings.ToLower(schema.Name)
		if strings.Contains(lowerName, "auth") ||
			strings.Contains(lowerName, "health") ||
			strings.Contains(lowerName, "config") {
			domain.IsPublic = true
		} else if strings.Contains(lowerName, "member") ||
			strings.Contains(lowerName, "organization") {
			domain.RequiresOrgMembership = true
		}

		if hasRepository {
			domain.HasRepository = true
		}

		// Commands and Queries are now auto-determined from operations
		// For entities, we assume they have both commands and queries
		if schema.XCodegenSchemaType == parsers.XCodegenSchemaTypeEntity {
			domain.HasCommands = true
			domain.HasQueries = true
		}
	}

	// Now populate Commands, Queries, and Operations from operations
	for _, op := range operations {

		// Get the domain from the first tag (now singular)
		// Need to find the domain by matching singular tag to domain's singular name
		tagLower := strings.ToLower(op.Tag)
		var domain *TagDef
		for _, d := range tagMap {
			// Match the singular tag to the domain's singular name
			if strings.ToLower(d.Name) == tagLower {
				domain = d
				break
			}
		}
		if domain == nil {
			continue
		}

		// Add operation to the domain's operations list
		domain.Operations = append(domain.Operations, op)

		// Determine if it's a command or query based on HTTP method
		// Use raw operationID to preserve original casing
		operationID := op.ID
		if op.Method == httpMethodGet {
			// It's a query
			if !slices.Contains(domain.Queries, operationID) {
				domain.Queries = append(domain.Queries, operationID)
			}
		} else {
			// It's a command (POST, PUT, PATCH, DELETE)
			if !slices.Contains(domain.Commands, operationID) {
				domain.Commands = append(domain.Commands, operationID)
			}
		}
	}

	// Convert map to slice and sort for consistent order
	for _, domain := range tagMap {
		// Sort operations first (this determines handler order)
		sortOperationsByRESTOrder(domain.Operations)

		// Now extract commands and queries in the order they appear in sorted operations
		domain.Commands = extractCommandsInOrder(domain.Operations)
		domain.Queries = extractQueriesInOrder(domain.Operations)

		tags = append(tags, *domain)
	}

	// Sort tags by name for consistent generation order
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].Name < tags[j].Name
	})

	// Generate app.go
	appData := map[string]any{
		"Domains": tags,
	}

	appPath := filepath.Join("internal/infrastructure/bootstrap", "app.gen.go")
	appTmpl, ok := g.templates["bootstrap.tmpl"]
	if !ok {
		return fmt.Errorf("app template not found")
	}

	if err := g.filewriter.WriteTemplate(appPath, appTmpl, appData); err != nil {
		return fmt.Errorf("failed to generate app.go: %w", err)
	}

	// Generate infrastructure.go
	infraPath := filepath.Join("internal/infrastructure/bootstrap", "infrastructure.gen.go")
	infraTmpl, ok := g.templates["infrastructure.tmpl"]
	if !ok {
		return fmt.Errorf("infrastructure template not found")
	}

	if err := g.filewriter.WriteTemplate(infraPath, infraTmpl, appData); err != nil {
		return fmt.Errorf("failed to generate infrastructure.go: %w", err)
	}

	return nil
}

// sortOperationsByRESTOrder sorts operations using the same order as controllers
func sortOperationsByRESTOrder(ops []parsers.OperationDef) {
	sort.SliceStable(ops, func(i, j int) bool {
		return getOperationOrder(ops[i]) < getOperationOrder(ops[j])
	})
}

// extractCommandsInOrder extracts command operation IDs in the order they appear
func extractCommandsInOrder(ops []parsers.OperationDef) []string {
	commands := []string{}
	for _, op := range ops {
		if op.Method != httpMethodGet {
			commands = append(commands, op.ID)
		}
	}
	return commands
}

// extractQueriesInOrder extracts query operation IDs in the order they appear
func extractQueriesInOrder(ops []parsers.OperationDef) []string {
	queries := []string{}
	for _, op := range ops {
		if op.Method == httpMethodGet {
			queries = append(queries, op.ID)
		}
	}
	return queries
}
