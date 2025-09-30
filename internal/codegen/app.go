package codegen

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/archesai/archesai/internal/parsers"
)

const (
	schemaNameSession  = "session"
	httpMethodGet      = "GET"
	securityBearerAuth = "bearerAuth"
	securitySession    = "sessionCookie"
)

// DomainData represents data for a single domain in the app template
type DomainData struct {
	Name                  string // e.g., "Account"
	NamePlural            string // e.g., "Accounts"
	NameLower             string // e.g., "account"
	NamePluralLower       string // e.g., "accounts"
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

// GenerateAppWiring generates the main application wiring file
func (g *Generator) GenerateAppWiring(
	schemas map[string]*parsers.ProcessedSchema,
	operations []parsers.OperationDef,
) error {
	// Collect domain data from schemas with x-codegen
	domains := make([]DomainData, 0)
	domainMap := make(map[string]*DomainData)

	// First, collect all tags that have actual operations
	tagsWithOperations := make(map[string]bool)
	for _, op := range operations {
		for _, tag := range op.Tags {
			tagsWithOperations[strings.ToLower(tag)] = true
		}
	}

	// Process each schema to build domain information
	for name, schema := range schemas {
		if schema.XCodegen == nil {
			continue
		}

		// Use the schema name to determine the domain
		// The plural form of the schema name becomes the domain
		namePlural := Pluralize(name)
		nameSingular := name

		// Check if this schema should be included
		hasAPIOperations := tagsWithOperations[strings.ToLower(namePlural)]
		hasRepository := schema.XCodegen.Repository != nil

		// Include schemas that either:
		// 1. Have API operations tagged with their plural name
		// 2. Have a repository defined (needed for cross-cutting concerns like auth)
		if !hasAPIOperations && !hasRepository {
			continue
		}

		// Get or create domain data
		domain, exists := domainMap[namePlural]
		if !exists {
			domain = &DomainData{
				Name:            nameSingular,
				NamePlural:      namePlural,
				NameLower:       strings.ToLower(nameSingular),
				NamePluralLower: strings.ToLower(namePlural),
				Commands:        []string{},
				Queries:         []string{},
				Operations:      []parsers.OperationDef{},
			}
			domainMap[namePlural] = domain
		}

		// Update domain based on x-codegen settings
		// Controllers are only created for schemas with operations
		if hasAPIOperations {
			domain.HasController = true
		}

		// Determine route protection level based on domain name
		// This is a simplified logic - you may need to adjust based on your needs
		lowerName := strings.ToLower(namePlural)
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
		if schema.XCodegen.SchemaType == schemaTypeEntity {
			domain.HasCommands = true
			domain.HasQueries = true
		}
	}

	// Now populate Commands, Queries, and Operations from operations
	for _, op := range operations {
		if len(op.Tags) == 0 {
			continue
		}

		// Get the domain from the first tag (plural form)
		// Need to find the domain by case-insensitive comparison
		tagLower := strings.ToLower(op.Tags[0])
		var domain *DomainData
		for domainName, d := range domainMap {
			if strings.ToLower(domainName) == tagLower {
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
		operationID := op.OperationID
		if op.Method == httpMethodGet {
			// It's a query
			if !contains(domain.Queries, operationID) {
				domain.Queries = append(domain.Queries, operationID)
			}
		} else {
			// It's a command (POST, PUT, PATCH, DELETE)
			if !contains(domain.Commands, operationID) {
				domain.Commands = append(domain.Commands, operationID)
			}
		}
	}

	// Convert map to slice and sort for consistent order
	for _, domain := range domainMap {
		// Sort commands and queries for consistent ordering
		domain.Commands = sortCommands(domain.Commands)
		domain.Queries = sortQueries(domain.Queries)
		domains = append(domains, *domain)
	}

	// Sort domains by name for consistent generation order
	sort.Slice(domains, func(i, j int) bool {
		return domains[i].Name < domains[j].Name
	})

	// Generate app.go
	appData := map[string]interface{}{
		"Domains": domains,
	}

	appPath := filepath.Join("internal/application/app", "app.gen.go")
	appTmpl, ok := g.templates["app.tmpl"]
	if !ok {
		return fmt.Errorf("app template not found")
	}

	if err := g.filewriter.WriteTemplate(appPath, appTmpl, appData); err != nil {
		return fmt.Errorf("failed to generate app.go: %w", err)
	}

	// Generate infrastructure.go
	infraPath := filepath.Join("internal/application/app", "infrastructure.gen.go")
	infraTmpl, ok := g.templates["infrastructure.tmpl"]
	if !ok {
		return fmt.Errorf("infrastructure template not found")
	}

	if err := g.filewriter.WriteTemplate(infraPath, infraTmpl, appData); err != nil {
		return fmt.Errorf("failed to generate infrastructure.go: %w", err)
	}

	return nil
}

// contains checks if a string is in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// sortCommands sorts command names alphabetically
func sortCommands(commands []string) []string {
	sorted := make([]string, len(commands))
	copy(sorted, commands)
	sort.Strings(sorted)
	return sorted
}

// sortQueries sorts query names alphabetically
func sortQueries(queries []string) []string {
	sorted := make([]string, len(queries))
	copy(sorted, queries)
	sort.Strings(sorted)
	return sorted
}
