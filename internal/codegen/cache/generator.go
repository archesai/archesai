// Package cache generates cache interfaces and implementations from OpenAPI specs.
package cache

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

// Generator handles generation of cache code.
type Generator struct{}

// NewGenerator creates a new cache generator.
func NewGenerator() *Generator {
	return &Generator{}
}

// Config represents cache generation configuration.
type Config struct {
	Domains map[string]DomainConfig `yaml:"domains"`
}

// DomainConfig represents configuration for a single domain.
type DomainConfig struct {
	OpenAPI string      `yaml:"openapi"`
	Tags    []string    `yaml:"tags"`
	Cache   CacheConfig `yaml:"cache"`
}

// CacheConfig represents cache adapter configuration.
type CacheConfig struct {
	Redis  RedisConfig  `yaml:"redis"`
	Memory MemoryConfig `yaml:"memory"`
}

// RedisConfig represents Redis cache configuration.
type RedisConfig struct {
	Enabled    bool   `yaml:"enabled"`
	Prefix     string `yaml:"prefix"`
	DefaultTTL int    `yaml:"default_ttl"` // in seconds
}

// MemoryConfig represents in-memory cache configuration.
type MemoryConfig struct {
	Enabled    bool `yaml:"enabled"`
	MaxItems   int  `yaml:"max_items"`
	DefaultTTL int  `yaml:"default_ttl"` // in seconds
}

// EntityInfo represents parsed entity information.
type EntityInfo struct {
	Name       string
	Type       string
	PrimaryKey string
	Cacheable  bool
	TTL        int // TTL in seconds
}

// TemplateData represents data passed to templates.
type TemplateData struct {
	Domain   string
	Package  string
	Entities []EntityInfo
	Imports  []string
	Config   DomainConfig
}

// Generate generates cache code for all configured domains.
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
		if err := g.generateDomainCache(domain, domainConfig); err != nil {
			return fmt.Errorf("failed to generate cache for %s: %w", domain, err)
		}
	}

	return nil
}

// generateDomainCache generates cache code for a single domain.
func (g *Generator) generateDomainCache(domain string, config DomainConfig) error {
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

	// Generate cache interface
	if err := g.generateFile("cache.go.tmpl", filepath.Join("internal", domain, "cache.go"), data); err != nil {
		return fmt.Errorf("failed to generate cache interface: %w", err)
	}

	// Generate Redis implementation if enabled
	if config.Cache.Redis.Enabled {
		data.Imports = append(data.Imports,
			"encoding/json",
			"fmt",
			"github.com/redis/go-redis/v9",
		)
		if err := g.generateFile("cache_redis.go.tmpl", filepath.Join("internal", domain, "cache_redis.gen.go"), data); err != nil {
			return fmt.Errorf("failed to generate redis cache: %w", err)
		}
	}

	// Generate in-memory implementation if enabled
	if config.Cache.Memory.Enabled {
		data.Imports = append(data.Imports,
			"sync",
			"fmt",
		)
		if err := g.generateFile("cache_memory.go.tmpl", filepath.Join("internal", domain, "cache_memory.gen.go"), data); err != nil {
			return fmt.Errorf("failed to generate memory cache: %w", err)
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
			ttl := 300 // 5 minutes default
			if config.Cache.Redis.DefaultTTL > 0 {
				ttl = config.Cache.Redis.DefaultTTL
			}
			entities = append(entities, EntityInfo{
				Name:       "User",
				Type:       "UserEntity",
				PrimaryKey: "Id",
				Cacheable:  true,
				TTL:        ttl,
			})
		case "Sessions":
			ttl := 900 // 15 minutes default for sessions
			if config.Cache.Redis.DefaultTTL > 0 {
				ttl = config.Cache.Redis.DefaultTTL
			}
			entities = append(entities, EntityInfo{
				Name:       "Session",
				Type:       "SessionEntity",
				PrimaryKey: "Id",
				Cacheable:  true,
				TTL:        ttl,
			})
		case "Accounts":
			ttl := 600 // 10 minutes default
			if config.Cache.Redis.DefaultTTL > 0 {
				ttl = config.Cache.Redis.DefaultTTL
			}
			entities = append(entities, EntityInfo{
				Name:       "Account",
				Type:       "AccountEntity",
				PrimaryKey: "Id",
				Cacheable:  true,
				TTL:        ttl,
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
