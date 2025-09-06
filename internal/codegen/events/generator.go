// Package events generates event publisher interfaces and implementations from OpenAPI specs.
package events

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"gopkg.in/yaml.v3"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

// Generator handles generation of event publisher code.
type Generator struct{}

// NewGenerator creates a new events generator.
func NewGenerator() *Generator {
	return &Generator{}
}

// Config represents events generation configuration.
type Config struct {
	Domains map[string]DomainConfig `yaml:"domains"`
}

// DomainConfig represents configuration for a single domain.
type DomainConfig struct {
	OpenAPI string       `yaml:"openapi"`
	Tags    []string     `yaml:"tags"`
	Events  EventsConfig `yaml:"events"`
}

// EventsConfig represents event publisher configuration.
type EventsConfig struct {
	Redis RedisConfig `yaml:"redis"`
	NATS  NATSConfig  `yaml:"nats"`
}

// RedisConfig represents Redis event publisher configuration.
type RedisConfig struct {
	Enabled bool   `yaml:"enabled"`
	Channel string `yaml:"channel"`
}

// NATSConfig represents NATS event publisher configuration.
type NATSConfig struct {
	Enabled bool   `yaml:"enabled"`
	Subject string `yaml:"subject"`
}

// EntityInfo represents parsed entity information.
type EntityInfo struct {
	Name       string
	Type       string
	PrimaryKey string
	Events     []EventType
}

// EventType represents an event type.
type EventType struct {
	Name        string
	Type        string
	Description string
}

// TemplateData represents data passed to templates.
type TemplateData struct {
	Domain   string
	Package  string
	Entities []EntityInfo
	Imports  []string
	Config   DomainConfig
}

// Generate generates event publisher code for all configured domains.
func (g *Generator) Generate(configPath string) error {
	// Read configuration
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// Generate for each domain
	for domain, domainConfig := range config.Domains {
		if err := g.generateDomainEvents(domain, domainConfig); err != nil {
			return fmt.Errorf("failed to generate events for %s: %w", domain, err)
		}
	}

	return nil
}

// generateDomainEvents generates event publisher code for a single domain.
func (g *Generator) generateDomainEvents(domain string, config DomainConfig) error {
	// Parse entities from configuration
	entities := g.parseEntities(config)

	if len(entities) == 0 {
		return nil // No entities to generate
	}

	// Template data
	data := TemplateData{
		Domain:   domain,
		Package:  domain,
		Entities: entities,
		Config:   config,
		Imports: []string{
			"context",
			"time",
			"github.com/google/uuid",
		},
	}

	// Generate events interface
	if err := g.generateFile("events.go.tmpl", filepath.Join("internal", domain, "events.go"), data); err != nil {
		return fmt.Errorf("failed to generate events interface: %w", err)
	}

	// Generate Redis implementation if enabled
	if config.Events.Redis.Enabled {
		data.Imports = append(data.Imports,
			"encoding/json",
			"fmt",
			"github.com/redis/go-redis/v9",
		)
		if err := g.generateFile("events_redis.go.tmpl", filepath.Join("internal", domain, "events_redis.gen.go"), data); err != nil {
			return fmt.Errorf("failed to generate redis events: %w", err)
		}
	}

	// Generate NATS implementation if enabled
	if config.Events.NATS.Enabled {
		data.Imports = append(data.Imports,
			"encoding/json",
			"fmt",
			"github.com/nats-io/nats.go",
		)
		if err := g.generateFile("events_nats.go.tmpl", filepath.Join("internal", domain, "events_nats.gen.go"), data); err != nil {
			return fmt.Errorf("failed to generate nats events: %w", err)
		}
	}

	return nil
}

// parseEntities parses entities from configuration.
func (g *Generator) parseEntities(config DomainConfig) []EntityInfo {
	var entities []EntityInfo

	// Map common entities based on tags
	for _, tag := range config.Tags {
		switch tag {
		case "Users":
			entities = append(entities, EntityInfo{
				Name:       "User",
				Type:       "UserEntity",
				PrimaryKey: "Id",
				Events: []EventType{
					{Name: "Created", Type: "UserCreated", Description: "User was created"},
					{Name: "Updated", Type: "UserUpdated", Description: "User was updated"},
					{Name: "Deleted", Type: "UserDeleted", Description: "User was deleted"},
					{Name: "EmailVerified", Type: "UserEmailVerified", Description: "User email was verified"},
					{Name: "PasswordChanged", Type: "UserPasswordChanged", Description: "User password was changed"},
				},
			})
		case "Sessions":
			entities = append(entities, EntityInfo{
				Name:       "Session",
				Type:       "SessionEntity",
				PrimaryKey: "Id",
				Events: []EventType{
					{Name: "Created", Type: "SessionCreated", Description: "Session was created"},
					{Name: "Deleted", Type: "SessionDeleted", Description: "Session was deleted"},
					{Name: "Expired", Type: "SessionExpired", Description: "Session expired"},
				},
			})
		case "Accounts":
			entities = append(entities, EntityInfo{
				Name:       "Account",
				Type:       "AccountEntity",
				PrimaryKey: "Id",
				Events: []EventType{
					{Name: "Linked", Type: "AccountLinked", Description: "Account was linked"},
					{Name: "Unlinked", Type: "AccountUnlinked", Description: "Account was unlinked"},
					{Name: "TokensRefreshed", Type: "AccountTokensRefreshed", Description: "Account tokens were refreshed"},
				},
			})
		case "Organizations":
			entities = append(entities, EntityInfo{
				Name:       "Organization",
				Type:       "OrganizationEntity",
				PrimaryKey: "Id",
				Events: []EventType{
					{Name: "Created", Type: "OrganizationCreated", Description: "Organization was created"},
					{Name: "Updated", Type: "OrganizationUpdated", Description: "Organization was updated"},
					{Name: "Deleted", Type: "OrganizationDeleted", Description: "Organization was deleted"},
					{Name: "MemberAdded", Type: "OrganizationMemberAdded", Description: "Member was added to organization"},
					{Name: "MemberRemoved", Type: "OrganizationMemberRemoved", Description: "Member was removed from organization"},
				},
			})
		}
	}

	return entities
}

// generateFile generates a file from template.
func (g *Generator) generateFile(templateName, outputPath string, data TemplateData) error {
	// Read template
	tmplContent, err := templatesFS.ReadFile(filepath.Join("templates", templateName))
	if err != nil {
		return fmt.Errorf("failed to read template %s: %w", templateName, err)
	}

	// Create template with helper functions
	tmpl, err := template.New(templateName).Funcs(template.FuncMap{
		"title": func(s string) string {
			if s == "" {
				return s
			}
			return string(unicode.ToUpper(rune(s[0]))) + s[1:]
		},
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"join":  strings.Join,
		"camelCase": func(s string) string {
			if s == "" {
				return s
			}
			return strings.ToLower(s[:1]) + s[1:]
		},
		"kebabCase": func(s string) string {
			var result []rune
			for i, r := range s {
				if i > 0 && unicode.IsUpper(r) {
					result = append(result, '-')
				}
				result = append(result, unicode.ToLower(r))
			}
			return string(result)
		},
	}).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templateName, err)
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	// Format Go code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		// If formatting fails, write unformatted code for debugging
		formatted = buf.Bytes()
	}

	// Ensure output directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write file
	if err := os.WriteFile(outputPath, formatted, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", outputPath, err)
	}

	return nil
}
