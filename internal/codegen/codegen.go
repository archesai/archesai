// Package codegen provides unified code generation from OpenAPI schemas with x-codegen extensions.
//
// This package generates repository interfaces, cache implementations, event publishers,
// and other boilerplate code from OpenAPI specifications annotated with x-codegen metadata.
//
// The generator reads codegen.yaml configuration and processes OpenAPI schemas to produce:
//   - Repository interfaces and database implementations
//   - Cache interfaces and memory/redis implementations
//   - Event publisher interfaces and NATS/Redis implementations
//   - HTTP handlers and adapters
//   - Configuration defaults
//
// All generated files follow the pattern *.gen.go and should not be edited manually.
//
//go:generate go tool oapi-codegen --config=../../types.codegen.yaml --package codegen --generate skip-prune,models ../../api/components/schemas/XCodegenWrapper.yaml
package codegen

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Config represents the unified codegen configuration.
type Config struct {
	// OpenAPI spec file
	OpenAPI string `yaml:"openapi"`

	// Output directory
	OutputDir string `yaml:"output"`

	// Domains to generate
	Domains map[string]DomainConfig `yaml:"domains"`

	// Generators to enable
	Generators GeneratorFlags `yaml:"generators"`

	// Global settings
	Settings GlobalSettings `yaml:"settings"`
}

// DomainConfig configures generation for a specific domain.
type DomainConfig struct {
	Tags    []string `yaml:"tags"`
	Schemas []string `yaml:"schemas"`
}

// GeneratorFlags controls which generators are enabled.
type GeneratorFlags struct {
	Repository bool `yaml:"repository"`
	Cache      bool `yaml:"cache"`
	Events     bool `yaml:"events"`
	Handlers   bool `yaml:"handlers"`
	Adapters   bool `yaml:"adapters"`
	Defaults   bool `yaml:"defaults"`
}

// GlobalSettings contains global generation settings.
type GlobalSettings struct {
	OverwriteExisting bool   `yaml:"overwrite"`
	GenerateTests     bool   `yaml:"generate_tests"`
	FileHeader        string `yaml:"header"`
}

// Run executes the unified code generator with the given configuration.
func Run(configPath string) error {
	// Load configuration
	config, err := loadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Determine base directory from config path
	baseDir := filepath.Dir(configPath)
	if baseDir == "." {
		baseDir = ""
	}

	// Create parser and file writer
	parser := NewParser(baseDir)
	fileWriter := NewFileWriter()

	// Configure file writer
	if config.Settings.FileHeader != "" {
		fileWriter.WithHeader(config.Settings.FileHeader)
	} else {
		fileWriter.WithHeader(DefaultHeader())
	}
	fileWriter.WithOverwrite(config.Settings.OverwriteExisting)

	// Load templates
	templates, err := loadTemplates()
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	log.Println("🚀 Starting unified code generation...")

	// Parse OpenAPI spec
	schemas, err := parser.ParseOpenAPISpec(config.OpenAPI)
	if err != nil {
		return fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	log.Printf("📋 Parsed %d schemas from %s", len(schemas), config.OpenAPI)

	// Filter schemas based on domain configuration
	filteredSchemas := filterSchemas(config, schemas)

	log.Printf("🎯 Filtered to %d schemas with x-codegen annotations", len(filteredSchemas))

	// Run each enabled generator
	flags := config.Generators
	if flags.Repository {
		if err := generateRepository(config, filteredSchemas, templates, fileWriter); err != nil {
			return fmt.Errorf("repository generator failed: %w", err)
		}
		log.Printf("✅ repository generator completed")
	}

	if flags.Cache {
		if err := generateCache(config, filteredSchemas, templates, fileWriter); err != nil {
			return fmt.Errorf("cache generator failed: %w", err)
		}
		log.Printf("✅ cache generator completed")
	}

	if flags.Events {
		if err := generateEvents(config, filteredSchemas, templates, fileWriter); err != nil {
			return fmt.Errorf("events generator failed: %w", err)
		}
		log.Printf("✅ events generator completed")
	}

	if flags.Handlers {
		if err := generateHandlers(config, filteredSchemas, templates, fileWriter); err != nil {
			return fmt.Errorf("handlers generator failed: %w", err)
		}
		log.Printf("✅ handlers generator completed")
	}

	if flags.Adapters {
		if err := generateAdapters(config, filteredSchemas, templates, fileWriter); err != nil {
			return fmt.Errorf("adapters generator failed: %w", err)
		}
		log.Printf("✅ adapters generator completed")
	}

	if flags.Defaults {
		if err := generateDefaults(config, filteredSchemas, templates, fileWriter); err != nil {
			return fmt.Errorf("defaults generator failed: %w", err)
		}
		log.Printf("✅ defaults generator completed")
	}

	// Report warnings
	if warnings := parser.GetWarnings(); len(warnings) > 0 {
		log.Println("⚠️  Warnings:")
		for _, w := range warnings {
			log.Printf("  - %s", w)
		}
	}

	log.Println("✨ Code generation completed successfully!")
	return nil
}

// loadTemplates loads all template files.
func loadTemplates() (map[string]*template.Template, error) {
	templates := make(map[string]*template.Template)
	templateFiles := []string{
		"repository.go.tmpl",
		"repository_postgres.go.tmpl",
		"repository_sqlite.go.tmpl",
		"cache.go.tmpl",
		"cache_memory.go.tmpl",
		"cache_redis.go.tmpl",
		"events.go.tmpl",
		"events_nats.go.tmpl",
		"events_redis.go.tmpl",
		"adapters.go.tmpl",
		"config.go.tmpl",
	}

	for _, file := range templateFiles {
		content, err := GetTemplate(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read template %s: %w", file, err)
		}

		tmpl, err := template.New(file).Funcs(TemplateFuncs()).Parse(content)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", file, err)
		}

		templates[file] = tmpl
	}

	return templates, nil
}

// filterSchemas filters schemas based on domain configuration and x-codegen presence.
func filterSchemas(config *Config, allSchemas map[string]*ParsedSchema) map[string]*ParsedSchema {
	filtered := make(map[string]*ParsedSchema)

	for name, s := range allSchemas {
		// Check if schema has any x-codegen configuration
		if !HasXCodegen(s) {
			continue
		}

		// Check if schema matches domain filters
		if matchesDomainFilter(config, s) {
			filtered[name] = s
		}
	}

	return filtered
}

// matchesDomainFilter checks if a schema matches the configured domain filters.
func matchesDomainFilter(config *Config, s *ParsedSchema) bool {
	// If no domain config, include all
	if len(config.Domains) == 0 {
		return true
	}

	// Check if schema's domain is configured
	domainConfig, ok := config.Domains[s.Domain]
	if !ok {
		return false
	}

	// If specific schemas are listed, check if this schema is included
	if len(domainConfig.Schemas) > 0 {
		for _, schemaName := range domainConfig.Schemas {
			if s.Name == schemaName {
				return true
			}
		}
		return false
	}

	// Include all schemas in the domain
	return true
}

// loadConfig loads the codegen configuration from a YAML file.
func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		// Try default locations
		defaultPaths := []string{
			"codegen.yaml",
			"codegen.yml",
			".codegen.yaml",
			".codegen.yml",
		}

		for _, defaultPath := range defaultPaths {
			if data, err = os.ReadFile(defaultPath); err == nil {
				break
			}
		}

		if err != nil {
			return nil, fmt.Errorf("config file not found: %w", err)
		}
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Set defaults
	if config.OutputDir == "" {
		config.OutputDir = "internal"
	}

	// Enable all generators by default if none specified
	if !config.Generators.Repository && !config.Generators.Cache &&
		!config.Generators.Events && !config.Generators.Handlers &&
		!config.Generators.Adapters && !config.Generators.Defaults {
		config.Generators.Repository = true
		config.Generators.Cache = true
		config.Generators.Events = true
		config.Generators.Handlers = true
		config.Generators.Adapters = true
		config.Generators.Defaults = true
	}

	return &config, nil
}

// generateRepository generates repository interfaces and implementations.
func generateRepository(config *Config, schemas map[string]*ParsedSchema, templates map[string]*template.Template, fileWriter *FileWriter) error {
	log.Printf("▶️  Running repository generator...")

	// Group schemas by domain
	domainSchemas := make(map[string][]*ParsedSchema)
	for _, s := range schemas {
		if NeedsRepository(s) {
			domainSchemas[s.Domain] = append(domainSchemas[s.Domain], s)
		}
	}

	// Sort domains for consistent output
	domains := make([]string, 0, len(domainSchemas))
	for domain := range domainSchemas {
		domains = append(domains, domain)
	}
	sort.Strings(domains)

	// Generate for each domain
	for _, domain := range domains {
		schemas := domainSchemas[domain]
		log.Printf("  Generating repository for domain '%s' with %d entities", domain, len(schemas))

		// Sort schemas by name for consistent output
		sort.Slice(schemas, func(i, j int) bool {
			return schemas[i].Name < schemas[j].Name
		})

		// Convert schemas to entities for template
		var entities []struct {
			Name              string
			Type              string
			Operations        []string
			AdditionalMethods []interface{}
		}
		for _, schema := range schemas {
			ops := []string{}
			if schema.XCodegen.Repository.Operations != nil {
				for _, op := range schema.XCodegen.Repository.Operations {
					ops = append(ops, string(op))
				}
			}
			entities = append(entities, struct {
				Name              string
				Type              string
				Operations        []string
				AdditionalMethods []interface{}
			}{
				Name:              schema.Name,
				Type:              schema.Name,
				Operations:        ops,
				AdditionalMethods: []interface{}{},
			})
		}

		// Generate repository interface
		repoData := struct {
			Domain   string
			Package  string
			Entities []struct {
				Name              string
				Type              string
				Operations        []string
				AdditionalMethods []interface{}
			}
			Imports []string
		}{
			Domain:   domain,
			Package:  domain,
			Entities: entities,
			Imports:  []string{"github.com/google/uuid"},
		}

		outputPath := filepath.Join(config.OutputDir, domain, "repository.gen.go")
		if err := fileWriter.WriteTemplate(outputPath, templates["repository.go.tmpl"], repoData); err != nil {
			return fmt.Errorf("failed to write repository for %s: %w", domain, err)
		}

		// Generate database implementations
		// TODO: Fix postgres and sqlite repository templates to handle new schema structure
		// The templates have hard-coded field mappings that break with schema changes
		// Commenting out for now until templates can be made more dynamic
		/*
			for _, schema := range schemas {
				ops := []string{}
				if schema.XCodegen.Repository.Operations != nil {
					for _, op := range schema.XCodegen.Repository.Operations {
						ops = append(ops, string(op))
					}
				}
				entity := struct {
					Name              string
					Type              string
					Operations        []string
					AdditionalMethods []interface{}
				}{
					Name:              schema.Name,
					Type:              schema.Name,
					Operations:        ops,
					AdditionalMethods: []interface{}{},
				}

				postgresData := struct {
					Domain   string
					Package  string
					Schema   *ParsedSchema
					Entities []struct {
						Name              string
						Type              string
						Operations        []string
						AdditionalMethods []interface{}
					}
					Imports []string
				}{
					Domain:  domain,
					Package: domain,
					Schema:  schema,
					Entities: []struct {
						Name              string
						Type              string
						Operations        []string
						AdditionalMethods []interface{}
					}{entity},
					Imports: []string{
						"github.com/google/uuid",
						fmt.Sprintf("github.com/archesai/archesai/internal/%s", domain),
					},
				}
				outputPath := filepath.Join(config.OutputDir, domain, "repository_postgres.gen.go")
				if err := fileWriter.WriteTemplate(outputPath, templates["repository_postgres.go.tmpl"], postgresData); err != nil {
					return fmt.Errorf("failed to write postgres repository for %s: %w", schema.Name, err)
				}

				ops = []string{}
				if schema.XCodegen.Repository.Operations != nil {
					for _, op := range schema.XCodegen.Repository.Operations {
						ops = append(ops, string(op))
					}
				}
				entity = struct {
					Name              string
					Type              string
					Operations        []string
					AdditionalMethods []interface{}
				}{
					Name:              schema.Name,
					Type:              schema.Name,
					Operations:        ops,
					AdditionalMethods: []interface{}{},
				}

				sqliteData := struct {
					Domain   string
					Package  string
					Schema   *ParsedSchema
					Entities []struct {
						Name              string
						Type              string
						Operations        []string
						AdditionalMethods []interface{}
					}
					Imports []string
				}{
					Domain:  domain,
					Package: domain,
					Schema:  schema,
					Entities: []struct {
						Name              string
						Type              string
						Operations        []string
						AdditionalMethods []interface{}
					}{entity},
					Imports: []string{
						"github.com/google/uuid",
						fmt.Sprintf("github.com/archesai/archesai/internal/%s", domain),
					},
				}
				outputPath = filepath.Join(config.OutputDir, domain, "repository_sqlite.gen.go")
				if err := fileWriter.WriteTemplate(outputPath, templates["repository_sqlite.go.tmpl"], sqliteData); err != nil {
					return fmt.Errorf("failed to write sqlite repository for %s: %w", schema.Name, err)
				}
				}
		*/
	}

	return nil
}

// generateCache generates cache interfaces and implementations.
func generateCache(config *Config, schemas map[string]*ParsedSchema, templates map[string]*template.Template, fileWriter *FileWriter) error {
	log.Printf("▶️  Running cache generator...")

	// Group schemas by domain
	domainSchemas := make(map[string][]*ParsedSchema)
	for _, s := range schemas {
		if NeedsCache(s) {
			domainSchemas[s.Domain] = append(domainSchemas[s.Domain], s)
		}
	}

	// Sort domains for consistent output
	domains := make([]string, 0, len(domainSchemas))
	for domain := range domainSchemas {
		domains = append(domains, domain)
	}
	sort.Strings(domains)

	// Generate for each domain
	for _, domain := range domains {
		schemas := domainSchemas[domain]

		// Sort schemas by name for consistent output
		sort.Slice(schemas, func(i, j int) bool {
			return schemas[i].Name < schemas[j].Name
		})

		log.Printf("  Generating cache for domain '%s' with %d entities", domain, len(schemas))

		cacheData := struct {
			Domain  string
			Package string
			Schemas []*ParsedSchema
			Imports []string
		}{
			Domain:  domain,
			Package: domain,
			Schemas: schemas,
			Imports: []string{"github.com/google/uuid"},
		}

		// Generate cache interface
		outputPath := filepath.Join(config.OutputDir, domain, "cache.gen.go")
		if err := fileWriter.WriteTemplate(outputPath, templates["cache.go.tmpl"], cacheData); err != nil {
			return fmt.Errorf("failed to write cache for %s: %w", domain, err)
		}

		// Generate cache implementations
		memoryPath := filepath.Join(config.OutputDir, domain, "cache_memory.gen.go")
		if err := fileWriter.WriteTemplate(memoryPath, templates["cache_memory.go.tmpl"], cacheData); err != nil {
			return fmt.Errorf("failed to write memory cache for %s: %w", domain, err)
		}

		redisPath := filepath.Join(config.OutputDir, domain, "cache_redis.gen.go")
		if err := fileWriter.WriteTemplate(redisPath, templates["cache_redis.go.tmpl"], cacheData); err != nil {
			return fmt.Errorf("failed to write redis cache for %s: %w", domain, err)
		}
	}

	return nil
}

// generateEvents generates event interfaces and implementations.
func generateEvents(config *Config, schemas map[string]*ParsedSchema, templates map[string]*template.Template, fileWriter *FileWriter) error {
	log.Printf("▶️  Running events generator...")

	// Group schemas by domain
	domainSchemas := make(map[string][]*ParsedSchema)
	for _, s := range schemas {
		if NeedsEvents(s) {
			domainSchemas[s.Domain] = append(domainSchemas[s.Domain], s)
		}
	}

	// Sort domains for consistent output
	domains := make([]string, 0, len(domainSchemas))
	for domain := range domainSchemas {
		domains = append(domains, domain)
	}
	sort.Strings(domains)

	// Generate for each domain
	for _, domain := range domains {
		schemas := domainSchemas[domain]

		// Sort schemas by name for consistent output
		sort.Slice(schemas, func(i, j int) bool {
			return schemas[i].Name < schemas[j].Name
		})

		log.Printf("  Generating events for domain '%s' with %d entities", domain, len(schemas))

		eventsData := struct {
			Domain   string
			Package  string
			Schemas  []*ParsedSchema
			Entities []*ParsedSchema // Alias for compatibility with templates
			Imports  []string
		}{
			Domain:   domain,
			Package:  domain,
			Schemas:  schemas,
			Entities: schemas, // Same data, different field name for template compatibility
			Imports:  []string{"github.com/google/uuid"},
		}

		// Generate events interface
		outputPath := filepath.Join(config.OutputDir, domain, "events.gen.go")
		if err := fileWriter.WriteTemplate(outputPath, templates["events.go.tmpl"], eventsData); err != nil {
			return fmt.Errorf("failed to write events for %s: %w", domain, err)
		}

		// Generate event implementations
		redisPath := filepath.Join(config.OutputDir, domain, "events_redis.gen.go")
		if err := fileWriter.WriteTemplate(redisPath, templates["events_redis.go.tmpl"], eventsData); err != nil {
			return fmt.Errorf("failed to write redis events for %s: %w", domain, err)
		}

		natsPath := filepath.Join(config.OutputDir, domain, "events_nats.gen.go")
		if err := fileWriter.WriteTemplate(natsPath, templates["events_nats.go.tmpl"], eventsData); err != nil {
			return fmt.Errorf("failed to write nats events for %s: %w", domain, err)
		}
	}

	return nil
}

// generateHandlers generates HTTP handler stubs.
func generateHandlers(_ *Config, _ map[string]*ParsedSchema, _ map[string]*template.Template, _ *FileWriter) error {
	log.Printf("▶️  Running handlers generator...")
	// Stub implementation - handlers generation is complex and may not be needed initially
	return nil
}

// generateAdapters generates type adapters/mappers.
func generateAdapters(config *Config, schemas map[string]*ParsedSchema, templates map[string]*template.Template, fileWriter *FileWriter) error {
	log.Printf("▶️  Running adapters generator...")

	// Group schemas by domain
	domainSchemas := make(map[string][]*ParsedSchema)
	for _, s := range schemas {
		if NeedsAdapter(s) {
			domainSchemas[s.Domain] = append(domainSchemas[s.Domain], s)
		}
	}

	// Generate for each domain
	if len(domainSchemas) == 0 {
		log.Printf("  No schemas configured for adapter generation")
		return nil
	}

	// Sort domains for consistent output
	domains := make([]string, 0, len(domainSchemas))
	for domain := range domainSchemas {
		domains = append(domains, domain)
	}
	sort.Strings(domains)

	for _, domain := range domains {
		schemas := domainSchemas[domain]

		// Sort schemas by name for consistent output
		sort.Slice(schemas, func(i, j int) bool {
			return schemas[i].Name < schemas[j].Name
		})

		log.Printf("  Generating adapters for domain '%s' with %d entities", domain, len(schemas))

		adaptersData := struct {
			Domain  string
			Package string
			Schemas []*ParsedSchema
			Imports []string
		}{
			Domain:  domain,
			Package: domain,
			Schemas: schemas,
			Imports: []string{"github.com/google/uuid"},
		}

		outputPath := filepath.Join(config.OutputDir, domain, "mappers.gen.go")
		if err := fileWriter.WriteTemplate(outputPath, templates["adapters.go.tmpl"], adaptersData); err != nil {
			return fmt.Errorf("failed to write adapters for %s: %w", domain, err)
		}
	}

	return nil
}

// generateDefaults generates configuration defaults.
func generateDefaults(_ *Config, _ map[string]*ParsedSchema, _ map[string]*template.Template, _ *FileWriter) error {
	log.Printf("▶️  Running defaults generator...")

	// For now, skip defaults generation as it's causing issues
	// TODO: Implement proper config schema detection and field mapping
	log.Printf("  Skipping defaults generation (not implemented)")
	return nil
}
