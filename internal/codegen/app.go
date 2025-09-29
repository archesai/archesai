package codegen

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/archesai/archesai/internal/parsers"
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
}

// GenerateAppWiring generates the main application wiring file
func (g *Generator) GenerateAppWiring(schemas map[string]*parsers.ProcessedSchema) error {
	// Collect domain data from schemas with x-codegen
	domains := make([]DomainData, 0)
	domainMap := make(map[string]*DomainData)

	// Process each schema to build domain information
	for name, schema := range schemas {
		if schema.XCodegen == nil {
			continue
		}

		// Use the schema name to determine the domain
		// The plural form of the schema name becomes the domain
		namePlural := Pluralize(name)
		nameSingular := name

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
			}
			domainMap[namePlural] = domain
		}

		// Update domain based on x-codegen settings
		if schema.XCodegen.Controller != nil {
			domain.HasController = true

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
		}

		if schema.XCodegen.Repository != nil {
			domain.HasRepository = true
		}

		if schema.XCodegen.Commands != nil {
			domain.HasCommands = true
			for _, op := range schema.XCodegen.Commands.Operations {
				commandName := strings.Title(strings.ToLower(op))
				if !contains(domain.Commands, commandName) {
					domain.Commands = append(domain.Commands, commandName)
				}
			}
		}

		if schema.XCodegen.Queries != nil {
			domain.HasQueries = true
			for _, op := range schema.XCodegen.Queries.Operations {
				queryName := strings.Title(strings.ToLower(op))
				if !contains(domain.Queries, queryName) {
					domain.Queries = append(domain.Queries, queryName)
				}
			}
		}
	}

	// Convert map to slice and sort for consistent order
	for _, domain := range domainMap {
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
