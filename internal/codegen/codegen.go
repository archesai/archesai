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

// runGenerator is a generic function to run any generator type
func runGenerator(generatorType string, config *Config, schemas map[string]*ParsedSchema, templates map[string]*template.Template, fileWriter *FileWriter, filterFunc func(*ParsedSchema) bool) error {
	log.Printf("▶️  Running %s generator...", generatorType)

	// Group schemas by domain
	domainSchemas := make(map[string][]*ParsedSchema)
	for _, s := range schemas {
		if filterFunc(s) {
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

		log.Printf("  Generating %s for domain '%s' with %d entities", generatorType, domain, len(schemas))

		// Prepare template data based on generator type
		var templateData interface{}
		var outputFiles []struct{ path, template string }

		switch generatorType {
		case "repository":
			templateData = prepareRepositoryData(domain, schemas)
			outputFiles = []struct{ path, template string }{
				{filepath.Join(config.OutputDir, domain, "repository.gen.go"), "repository.go.tmpl"},
			}
		case "cache":
			templateData = prepareCacheData(domain, schemas)
			outputFiles = []struct{ path, template string }{
				{filepath.Join(config.OutputDir, domain, "cache.gen.go"), "cache.go.tmpl"},
				{filepath.Join(config.OutputDir, domain, "cache_memory.gen.go"), "cache_memory.go.tmpl"},
				{filepath.Join(config.OutputDir, domain, "cache_redis.gen.go"), "cache_redis.go.tmpl"},
			}
		case "events":
			templateData = prepareEventsData(domain, schemas)
			outputFiles = []struct{ path, template string }{
				{filepath.Join(config.OutputDir, domain, "events.gen.go"), "events.go.tmpl"},
				{filepath.Join(config.OutputDir, domain, "events_redis.gen.go"), "events_redis.go.tmpl"},
				{filepath.Join(config.OutputDir, domain, "events_nats.gen.go"), "events_nats.go.tmpl"},
			}
		case "adapters":
			templateData = prepareAdaptersData(domain, schemas)
			outputFiles = []struct{ path, template string }{
				{filepath.Join(config.OutputDir, domain, "mappers.gen.go"), "adapters.go.tmpl"},
			}
		}

		// Write all template files
		for _, file := range outputFiles {
			if err := fileWriter.WriteTemplate(file.path, templates[file.template], templateData); err != nil {
				return fmt.Errorf("failed to write %s for %s: %w", generatorType, domain, err)
			}
		}
	}

	return nil
}

// Helper functions to prepare data for each generator type
func prepareRepositoryData(domain string, schemas []*ParsedSchema) interface{} {
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

	return struct {
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
}

func prepareCacheData(domain string, schemas []*ParsedSchema) interface{} {
	return struct {
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
}

func prepareEventsData(domain string, schemas []*ParsedSchema) interface{} {
	return struct {
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
}

func prepareAdaptersData(domain string, schemas []*ParsedSchema) interface{} {
	return struct {
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
}

// generateRepository generates repository interfaces and implementations.
func generateRepository(config *Config, schemas map[string]*ParsedSchema, templates map[string]*template.Template, fileWriter *FileWriter) error {
	return runGenerator("repository", config, schemas, templates, fileWriter, NeedsRepository)
}

// generateCache generates cache interfaces and implementations.
func generateCache(config *Config, schemas map[string]*ParsedSchema, templates map[string]*template.Template, fileWriter *FileWriter) error {
	return runGenerator("cache", config, schemas, templates, fileWriter, NeedsCache)
}

// generateEvents generates event interfaces and implementations.
func generateEvents(config *Config, schemas map[string]*ParsedSchema, templates map[string]*template.Template, fileWriter *FileWriter) error {
	return runGenerator("events", config, schemas, templates, fileWriter, NeedsEvents)
}

// generateHandlers generates HTTP handler stubs.
func generateHandlers(_ *Config, _ map[string]*ParsedSchema, _ map[string]*template.Template, _ *FileWriter) error {
	log.Printf("▶️  Running handlers generator...")
	// Stub implementation - handlers generation is complex and may not be needed initially
	return nil
}

// generateAdapters generates type adapters/mappers.
func generateAdapters(config *Config, schemas map[string]*ParsedSchema, templates map[string]*template.Template, fileWriter *FileWriter) error {
	// Check if there are any schemas that need adapters first
	hasAdapters := false
	for _, s := range schemas {
		if NeedsAdapter(s) {
			hasAdapters = true
			break
		}
	}

	if !hasAdapters {
		log.Printf("▶️  Running adapters generator...")
		log.Printf("  No schemas configured for adapter generation")
		return nil
	}

	return runGenerator("adapters", config, schemas, templates, fileWriter, NeedsAdapter)
}

// generateDefaults generates configuration defaults.
func generateDefaults(config *Config, _ map[string]*ParsedSchema, _ map[string]*template.Template, fileWriter *FileWriter) error {
	log.Printf("▶️  Running defaults generator...")

	// Create parser
	parser := NewParser(filepath.Dir(config.OpenAPI))

	// Parse OpenAPI spec
	_, err := parser.ParseOpenAPISpec(config.OpenAPI)
	if err != nil {
		return fmt.Errorf("failed to parse spec: %w", err)
	}

	// Get complete defaults
	defaults, err := parser.GetCompleteConfigDefaults()
	if err != nil {
		return fmt.Errorf("failed to get defaults: %w", err)
	}

	// Generate Go code
	code := GenerateDefaultsCode(defaults)

	// Write to file
	outputPath := "./internal/config/defaults.gen.go"
	if err := fileWriter.WriteFile(outputPath, []byte(code)); err != nil {
		return fmt.Errorf("failed to write defaults file: %w", err)
	}

	log.Printf("  ✅ Generated %s", outputPath)
	log.Printf("  Total defaults generated: %d", len(parser.FlattenConfigDefaults(defaults)))

	return nil
}
