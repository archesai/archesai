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
package codegen

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"text/template"

	"github.com/archesai/archesai/internal/logger"
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

	// Generators configuration
	Generators *GeneratorConfig `yaml:"generators,omitempty"`

	// Global settings
	Settings GlobalSettings `yaml:"settings"`
}

// DomainConfig configures generation for a specific domain.
type DomainConfig struct {
	Tags    []string `yaml:"tags"`
	Schemas []string `yaml:"schemas"`
}

// GeneratorConfig controls generator output paths.
// If a path is specified, the generator is enabled for that output.
// Paths are relative to the domain directory.
type GeneratorConfig struct {
	Repository *RepositoryGeneratorConfig `yaml:"repository,omitempty"`
	Cache      *CacheGeneratorConfig      `yaml:"cache,omitempty"`
	Events     *EventsGeneratorConfig     `yaml:"events,omitempty"`
	Handlers   string                     `yaml:"handlers,omitempty"` // e.g., "handlers.gen.go"
	Adapters   string                     `yaml:"adapters,omitempty"` // e.g., "adapters.gen.go"
	Tests      *TestsGeneratorConfig      `yaml:"tests,omitempty"`
	Defaults   string                     `yaml:"defaults,omitempty"` // Global path for defaults
}

// RepositoryGeneratorConfig specifies repository generator outputs.
type RepositoryGeneratorConfig struct {
	Interface string `yaml:"interface,omitempty"` // e.g., "ports.gen.go"
	Postgres  string `yaml:"postgres,omitempty"`  // e.g., "postgres_repository.gen.go"
	SQLite    string `yaml:"sqlite,omitempty"`    // e.g., "sqlite_repository.gen.go"
}

// CacheGeneratorConfig specifies cache generator outputs.
type CacheGeneratorConfig struct {
	Interface string `yaml:"interface,omitempty"` // e.g., "cache.gen.go"
	Memory    string `yaml:"memory,omitempty"`    // e.g., "cache_memory.gen.go"
	Redis     string `yaml:"redis,omitempty"`     // e.g., "cache_redis.gen.go"
}

// EventsGeneratorConfig specifies event generator outputs.
type EventsGeneratorConfig struct {
	Interface string `yaml:"interface,omitempty"` // e.g., "events.gen.go"
	Redis     string `yaml:"redis,omitempty"`     // e.g., "events_redis.gen.go"
	NATS      string `yaml:"nats,omitempty"`      // e.g., "events_nats.gen.go"
}

// TestsGeneratorConfig specifies test generator outputs.
type TestsGeneratorConfig struct {
	Service    string `yaml:"service,omitempty"`    // e.g., "service_test.gen.go"
	Repository string `yaml:"repository,omitempty"` // e.g., "repository_test.gen.go"
	Handler    string `yaml:"handler,omitempty"`    // e.g., "handler_test.gen.go"
}

// GlobalSettings contains global generation settings.
type GlobalSettings struct {
	OverwriteExisting bool   `yaml:"overwrite"`
	GenerateTests     bool   `yaml:"generate_tests"`
	FileHeader        string `yaml:"header"`
}

// RunWithConfig executes the unified code generator with the given configuration.
func RunWithConfig(config *Config) error {
	// Create logger with error level by default (only show errors)
	// Set ARCHESAI_LOG_LEVEL=debug to see debug logs
	logLevel := os.Getenv("ARCHESAI_LOG_LEVEL")
	if logLevel == "" {
		logLevel = "error"
	}
	log := logger.New(logger.Config{Level: logLevel, Pretty: true})

	// Create parser and file writer
	parser := NewParser("")
	fileWriter := NewFileWriter()

	// Configure file writer
	fileWriter.WithOverwrite(true) // Always overwrite generated files
	if config.Settings.FileHeader != "" {
		fileWriter.WithHeader(config.Settings.FileHeader)
	} else {
		fileWriter.WithHeader(DefaultHeader())
	}

	// Load templates
	templates, err := loadTemplates()
	if err != nil {
		return fmt.Errorf("failed to load templates: %w", err)
	}

	log.Debug("Starting unified code generation...")

	// Parse OpenAPI spec
	schemas, err := parser.ParseOpenAPISpec(config.OpenAPI)
	if err != nil {
		return fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	log.Debug("Parsed schemas", slog.Int("count", len(schemas)), slog.String("spec", config.OpenAPI))

	// Auto-detect domains if not configured
	if len(config.Domains) == 0 {
		config.Domains = autoDetectDomains(schemas)
		log.Debug("Auto-detected domains", slog.Int("count", len(config.Domains)))
	}

	// Filter schemas based on domain configuration
	filteredSchemas := filterSchemas(config, schemas)

	log.Debug("Filtered schemas with x-codegen", slog.Int("count", len(filteredSchemas)))

	// Run each enabled generator based on path configuration
	if config.Generators != nil && config.Generators.Repository != nil {
		if err := generateRepository(config, filteredSchemas, templates, fileWriter, log); err != nil {
			return fmt.Errorf("repository generator failed: %w", err)
		}
		log.Debug("repository generator completed")
	}

	if config.Generators != nil && config.Generators.Cache != nil {
		if err := generateCache(config, filteredSchemas, templates, fileWriter, log); err != nil {
			return fmt.Errorf("cache generator failed: %w", err)
		}
		log.Debug("cache generator completed")
	}

	if config.Generators != nil && config.Generators.Events != nil {
		if err := generateEvents(config, filteredSchemas, templates, fileWriter, log); err != nil {
			return fmt.Errorf("events generator failed: %w", err)
		}
		log.Debug("events generator completed")
	}

	if config.Generators != nil && config.Generators.Handlers != "" {
		if err := generateHandlers(config, filteredSchemas, templates, fileWriter, log); err != nil {
			return fmt.Errorf("handlers generator failed: %w", err)
		}
		log.Debug("handlers generator completed")
	}

	if config.Generators != nil && config.Generators.Adapters != "" {
		if err := generateAdapters(config, filteredSchemas, templates, fileWriter, log); err != nil {
			return fmt.Errorf("adapters generator failed: %w", err)
		}
		log.Debug("adapters generator completed")
	}

	if config.Generators != nil && config.Generators.Defaults != "" {
		if err := generateDefaults(config, filteredSchemas, templates, fileWriter, log); err != nil {
			return fmt.Errorf("defaults generator failed: %w", err)
		}
		log.Debug("defaults generator completed")
	}

	// Generate tests if enabled
	if config.Settings.GenerateTests {
		generateTests(config, filteredSchemas, templates, fileWriter, log)
		log.Debug("test generator completed")
	}

	// Report warnings
	if warnings := parser.GetWarnings(); len(warnings) > 0 {
		log.Warn("Parsing warnings found")
		for _, w := range warnings {
			log.Debug("Warning", slog.String("message", w))
		}
	}

	log.Debug("Code generation completed successfully")
	return nil
}

// Run executes the unified code generator with the given configuration file.
func Run(configPath string) error {
	// Config path is required
	if configPath == "" {
		return fmt.Errorf("config file path is required")
	}

	// Load configuration
	config, err := loadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	return RunWithConfig(config)
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
			"archesai.codegen.yaml",
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

	// If no generators configured, use defaults
	if config.Generators == nil {
		config.Generators = &GeneratorConfig{
			Repository: &RepositoryGeneratorConfig{
				Interface: "ports.gen.go",
				Postgres:  "postgres_repository.gen.go",
				SQLite:    "sqlite_repository.gen.go",
			},
			Cache: &CacheGeneratorConfig{
				Interface: "cache.gen.go",
				Memory:    "cache_memory.gen.go",
				Redis:     "cache_redis.gen.go",
			},
			Events: &EventsGeneratorConfig{
				Interface: "events.gen.go",
				Redis:     "events_redis.gen.go",
			},
			Handlers: "handlers.gen.go",
			Adapters: "adapters.gen.go",
			Defaults: "./internal/config/defaults.gen.go",
		}
	}

	return &config, nil
}

// runGeneratorWithPaths runs generators using path-based configuration
func runGeneratorWithPaths(generatorType string, config *Config, schemas map[string]*ParsedSchema, templates map[string]*template.Template, fileWriter *FileWriter, filterFunc func(*ParsedSchema) bool, log *slog.Logger) error {
	log.Debug("Running generator", slog.String("type", generatorType))

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

		log.Debug("Generating for domain",
			slog.String("generator", generatorType),
			slog.String("domain", domain),
			slog.Int("entities", len(schemas)))

		// Prepare template data and output files based on generator type and config
		var templateData interface{}
		var outputFiles []struct{ path, template string }

		switch generatorType {
		case "repository":
			if config.Generators.Repository == nil {
				continue
			}
			templateData = prepareRepositoryData(domain, schemas)
			outputFiles = []struct{ path, template string }{}
			if config.Generators.Repository.Interface != "" {
				outputFiles = append(outputFiles, struct{ path, template string }{
					filepath.Join(config.OutputDir, domain, config.Generators.Repository.Interface), "repository.go.tmpl",
				})
			}
			if config.Generators.Repository.Postgres != "" {
				outputFiles = append(outputFiles, struct{ path, template string }{
					filepath.Join(config.OutputDir, domain, config.Generators.Repository.Postgres), "repository_postgres.go.tmpl",
				})
			}
			if config.Generators.Repository.SQLite != "" {
				outputFiles = append(outputFiles, struct{ path, template string }{
					filepath.Join(config.OutputDir, domain, config.Generators.Repository.SQLite), "repository_sqlite.go.tmpl",
				})
			}
		case "cache":
			if config.Generators.Cache == nil {
				continue
			}
			templateData = prepareCacheData(domain, schemas)
			outputFiles = []struct{ path, template string }{}
			if config.Generators.Cache.Interface != "" {
				outputFiles = append(outputFiles, struct{ path, template string }{
					filepath.Join(config.OutputDir, domain, config.Generators.Cache.Interface), "cache.go.tmpl",
				})
			}
			if config.Generators.Cache.Memory != "" {
				outputFiles = append(outputFiles, struct{ path, template string }{
					filepath.Join(config.OutputDir, domain, config.Generators.Cache.Memory), "cache_memory.go.tmpl",
				})
			}
			if config.Generators.Cache.Redis != "" {
				outputFiles = append(outputFiles, struct{ path, template string }{
					filepath.Join(config.OutputDir, domain, config.Generators.Cache.Redis), "cache_redis.go.tmpl",
				})
			}
		case "events":
			if config.Generators.Events == nil {
				continue
			}
			templateData = prepareEventsData(domain, schemas)
			outputFiles = []struct{ path, template string }{}
			if config.Generators.Events.Interface != "" {
				outputFiles = append(outputFiles, struct{ path, template string }{
					filepath.Join(config.OutputDir, domain, config.Generators.Events.Interface), "events.go.tmpl",
				})
			}
			if config.Generators.Events.Redis != "" {
				outputFiles = append(outputFiles, struct{ path, template string }{
					filepath.Join(config.OutputDir, domain, config.Generators.Events.Redis), "events_redis.go.tmpl",
				})
			}
			if config.Generators.Events.NATS != "" {
				outputFiles = append(outputFiles, struct{ path, template string }{
					filepath.Join(config.OutputDir, domain, config.Generators.Events.NATS), "events_nats.go.tmpl",
				})
			}
		case "adapters":
			if config.Generators.Adapters == "" {
				continue
			}
			templateData = prepareAdaptersData(domain, schemas)
			outputFiles = []struct{ path, template string }{
				{filepath.Join(config.OutputDir, domain, config.Generators.Adapters), "adapters.go.tmpl"},
			}
		case "handlers":
			if config.Generators.Handlers == "" {
				continue
			}
			templateData = prepareHandlersData(domain, schemas)
			outputFiles = []struct{ path, template string }{
				{filepath.Join(config.OutputDir, domain, config.Generators.Handlers), "handlers.go.tmpl"},
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

		// Extract additional methods from x-codegen
		var additionalMethods []interface{}
		if schema.XCodegen != nil && schema.XCodegen.Repository.AdditionalMethods != nil {
			for _, method := range schema.XCodegen.Repository.AdditionalMethods {
				additionalMethods = append(additionalMethods, struct {
					Name string
				}{
					Name: method.Name,
				})
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
			AdditionalMethods: additionalMethods,
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

func prepareHandlersData(domain string, schemas []*ParsedSchema) interface{} {
	return struct {
		Domain   string
		Package  string
		Entities []*ParsedSchema
		Imports  []string
	}{
		Domain:   domain,
		Package:  domain,
		Entities: schemas,
		Imports:  []string{"github.com/google/uuid", "github.com/labstack/echo/v4"},
	}
}

// generateRepository generates repository interfaces and implementations.
func generateRepository(config *Config, schemas map[string]*ParsedSchema, templates map[string]*template.Template, fileWriter *FileWriter, log *slog.Logger) error {
	return runGeneratorWithPaths("repository", config, schemas, templates, fileWriter, NeedsRepository, log)
}

// generateCache generates cache interfaces and implementations.
func generateCache(config *Config, schemas map[string]*ParsedSchema, templates map[string]*template.Template, fileWriter *FileWriter, log *slog.Logger) error {
	return runGeneratorWithPaths("cache", config, schemas, templates, fileWriter, NeedsCache, log)
}

// generateEvents generates event interfaces and implementations.
func generateEvents(config *Config, schemas map[string]*ParsedSchema, templates map[string]*template.Template, fileWriter *FileWriter, log *slog.Logger) error {
	return runGeneratorWithPaths("events", config, schemas, templates, fileWriter, NeedsEvents, log)
}

// generateHandlers generates HTTP handler stubs.
func generateHandlers(_ *Config, _ map[string]*ParsedSchema, _ map[string]*template.Template, _ *FileWriter, log *slog.Logger) error {
	// For now, return without generating since handlers are complex
	// and typically generated by oapi-codegen
	log.Debug("Running handlers generator")
	log.Debug("Handler generation delegated to oapi-codegen")
	return nil
}

// generateAdapters generates type adapters/mappers.
func generateAdapters(config *Config, schemas map[string]*ParsedSchema, templates map[string]*template.Template, fileWriter *FileWriter, log *slog.Logger) error {
	// Check if there are any schemas that need adapters first
	hasAdapters := false
	for _, s := range schemas {
		if NeedsAdapter(s) {
			hasAdapters = true
			break
		}
	}

	if !hasAdapters {
		log.Debug("Running adapters generator")
		log.Debug("No schemas configured for adapter generation")
		return nil
	}

	return runGeneratorWithPaths("adapters", config, schemas, templates, fileWriter, NeedsAdapter, log)
}

// generateDefaults generates configuration defaults.
func generateDefaults(config *Config, _ map[string]*ParsedSchema, _ map[string]*template.Template, fileWriter *FileWriter, log *slog.Logger) error {
	log.Debug("Running defaults generator")

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

	log.Debug("Generated defaults",
		slog.String("path", outputPath),
		slog.Int("count", len(parser.FlattenConfigDefaults(defaults))))

	return nil
}

// generateTests generates test boilerplate for all domains
func generateTests(_ *Config, schemas map[string]*ParsedSchema, _ map[string]*template.Template, _ *FileWriter, log *slog.Logger) {
	log.Debug("Running test generator")

	// Group schemas by domain
	domainSchemas := make(map[string][]*ParsedSchema)
	for _, s := range schemas {
		domainSchemas[s.Domain] = append(domainSchemas[s.Domain], s)
	}

	// Sort domains for consistent output
	domains := make([]string, 0, len(domainSchemas))
	for domain := range domainSchemas {
		domains = append(domains, domain)
	}
	sort.Strings(domains)

	// Generate tests for each domain
	for _, domain := range domains {
		schemas := domainSchemas[domain]

		// Sort schemas by name for consistent output
		sort.Slice(schemas, func(i, j int) bool {
			return schemas[i].Name < schemas[j].Name
		})

		log.Debug("Generating tests for domain",
			slog.String("domain", domain),
			slog.Int("entities", len(schemas)))

		// TODO: Implement test generation with templates
		// For now, we'll skip the actual generation
		log.Debug("Test generation not yet implemented", slog.String("domain", domain))
	}
}

// autoDetectDomains analyzes schemas and auto-detects domain configuration
func autoDetectDomains(schemas map[string]*ParsedSchema) map[string]DomainConfig {
	domains := make(map[string]DomainConfig)

	// Collect unique domains from schemas
	domainSchemas := make(map[string][]string)
	for name, schema := range schemas {
		if schema.Domain == "" {
			continue
		}

		domainSchemas[schema.Domain] = append(domainSchemas[schema.Domain], name)
	}

	// Convert to DomainConfig
	for domain, schemaNames := range domainSchemas {
		// For now, we don't have tags in ParsedSchema, so we'll use domain names
		// In the future, this could be enhanced to extract tags from the OpenAPI spec
		domains[domain] = DomainConfig{
			Schemas: schemaNames,
		}
	}

	// If no domains detected, use default mapping
	if len(domains) == 0 {
		domains = map[string]DomainConfig{
			"auth":          {Tags: []string{"Auth", "Users", "Sessions", "Accounts"}},
			"organizations": {Tags: []string{"Organizations", "Members", "Invitations"}},
			"workflows":     {Tags: []string{"Workflows", "Pipelines", "Runs", "Tools"}},
			"content":       {Tags: []string{"Content", "Artifacts", "Labels"}},
		}
	}

	return domains
}
