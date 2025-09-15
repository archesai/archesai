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

//go:generate go tool oapi-codegen --config=../../.types.codegen.yaml --package codegen --generate skip-prune,models ../../api/components/schemas/XCodegenWrapper.yaml

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/archesai/archesai/internal/logger"
	"gopkg.in/yaml.v3"
)

// Config is an alias to the generated CodegenConfig type for backward compatibility
type Config = CodegenConfig

// RunWithConfig executes the unified code generator with the given configuration.
func RunWithConfig(config *CodegenConfig) error {
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
	if config.Settings.Header != "" {
		fileWriter.WithHeader(config.Settings.Header)
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
	schemas, err := parser.ParseOpenAPISpec(config.Openapi)
	if err != nil {
		return fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	log.Debug("Parsed schemas", slog.Int("count", len(schemas)), slog.String("spec", config.Openapi))

	// Auto-detect domains if not configured
	if len(config.Domains) == 0 {
		config.Domains = autoDetectDomains(schemas)
		log.Debug("Auto-detected domains", slog.Int("count", len(config.Domains)))
	}

	// Filter schemas based on domain configuration
	filteredSchemas := filterSchemas(config, schemas)

	log.Debug("Filtered schemas with x-codegen", slog.Int("count", len(filteredSchemas)))

	// Run each enabled generator based on path configuration
	if config.Generators.SQL.SchemaDir != "" || config.Generators.SQL.QueryDir != "" {
		if err := generateSQL(config, filteredSchemas, fileWriter, log); err != nil {
			return fmt.Errorf("SQL generator failed: %w", err)
		}
		log.Debug("SQL generator completed")
	}

	if config.Generators.Repository.Interface != "" || config.Generators.Repository.Postgres != "" || config.Generators.Repository.Sqlite != "" {
		if err := generateRepository(config, filteredSchemas, templates, fileWriter, log); err != nil {
			return fmt.Errorf("repository generator failed: %w", err)
		}
		log.Debug("repository generator completed")
	}

	if config.Generators.Cache.Interface != "" || config.Generators.Cache.Memory != "" || config.Generators.Cache.Redis != "" {
		if err := generateCache(config, filteredSchemas, templates, fileWriter, log); err != nil {
			return fmt.Errorf("cache generator failed: %w", err)
		}
		log.Debug("cache generator completed")
	}

	if config.Generators.Events.Interface != "" || config.Generators.Events.Redis != "" || config.Generators.Events.Nats != "" {
		if err := generateEvents(config, filteredSchemas, templates, fileWriter, log); err != nil {
			return fmt.Errorf("events generator failed: %w", err)
		}
		log.Debug("events generator completed")
	}

	if config.Generators.Service.Interface != "" || config.Generators.Service.Implementation != "" {
		if err := generateService(config, filteredSchemas, templates, fileWriter, log); err != nil {
			return fmt.Errorf("service generator failed: %w", err)
		}
		log.Debug("service generator completed")
	}

	if config.Generators.Defaults != "" {
		if err := generateDefaults(config, filteredSchemas, templates, fileWriter, log); err != nil {
			return fmt.Errorf("defaults generator failed: %w", err)
		}
		log.Debug("defaults generator completed")
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
		"service.go.tmpl",
		"service_impl.go.tmpl",
		"config.go.tmpl",
		"handler.go.tmpl",
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
			".archesai.codegen.yaml",
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
	if config.Output == "" {
		config.Output = "internal"
	}

	// Generators are always embedded in the generated type, no defaults needed

	return &config, nil
}

// runGeneratorWithPaths runs generators using path-based configuration
func runGeneratorWithPaths(generatorType string, config *Config, schemas map[string]*ParsedSchema, templates map[string]*template.Template, fileWriter *FileWriter, filterFunc func(*ParsedSchema) bool, log *slog.Logger) error {
	log.Debug("Running generator", slog.String("type", generatorType))

	// Group schemas by domain
	domainSchemas := groupSchemasByDomain(schemas, filterFunc)

	// Sort domains for consistent output
	domains := getSortedDomains(domainSchemas)

	// Generate for each domain
	for _, domain := range domains {
		schemas := domainSchemas[domain]
		sortSchemasByName(schemas)

		log.Debug("Generating for domain",
			slog.String("generator", generatorType),
			slog.String("domain", domain),
			slog.Int("entities", len(schemas)))

		// Get template data and output files for the generator type
		templateData, outputFiles := getGeneratorConfig(generatorType, config, domain, schemas)
		if templateData == nil {
			continue
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

// groupSchemasByDomain groups schemas by their domain
func groupSchemasByDomain(schemas map[string]*ParsedSchema, filterFunc func(*ParsedSchema) bool) map[string][]*ParsedSchema {
	domainSchemas := make(map[string][]*ParsedSchema)
	for _, s := range schemas {
		if filterFunc(s) {
			domainSchemas[s.Domain] = append(domainSchemas[s.Domain], s)
		}
	}
	return domainSchemas
}

// getSortedDomains returns sorted domain names
func getSortedDomains(domainSchemas map[string][]*ParsedSchema) []string {
	domains := make([]string, 0, len(domainSchemas))
	for domain := range domainSchemas {
		domains = append(domains, domain)
	}
	sort.Strings(domains)
	return domains
}

// sortSchemasByName sorts schemas by name
func sortSchemasByName(schemas []*ParsedSchema) {
	sort.Slice(schemas, func(i, j int) bool {
		return schemas[i].Name < schemas[j].Name
	})
}

// getGeneratorConfig returns template data and output files for a generator type
func getGeneratorConfig(generatorType string, config *Config, domain string, schemas []*ParsedSchema) (interface{}, []struct{ path, template string }) {
	switch generatorType {
	case "repository":
		return getRepositoryConfig(config, domain, schemas)
	case "cache":
		return getCacheConfig(config, domain, schemas)
	case "events":
		return getEventsConfig(config, domain, schemas)
	case "service":
		return nil, nil
	case "handler":
		return getHandlerConfig(config, domain, schemas)
	default:
		return nil, nil
	}
}

// getRepositoryConfig returns repository generator configuration
func getRepositoryConfig(config *Config, domain string, schemas []*ParsedSchema) (interface{}, []struct{ path, template string }) {
	if config.Generators.Repository.Interface == "" && config.Generators.Repository.Postgres == "" && config.Generators.Repository.Sqlite == "" {
		return nil, nil
	}
	templateData := prepareRepositoryData(domain, schemas)
	var outputFiles []struct{ path, template string }
	if config.Generators.Repository.Interface != "" {
		outputFiles = append(outputFiles, struct{ path, template string }{
			filepath.Join(config.Output, domain, config.Generators.Repository.Interface), "repository.go.tmpl",
		})
	}
	if config.Generators.Repository.Postgres != "" {
		outputFiles = append(outputFiles, struct{ path, template string }{
			filepath.Join(config.Output, domain, config.Generators.Repository.Postgres), "repository_postgres.go.tmpl",
		})
	}
	if config.Generators.Repository.Sqlite != "" {
		outputFiles = append(outputFiles, struct{ path, template string }{
			filepath.Join(config.Output, domain, config.Generators.Repository.Sqlite), "repository_sqlite.go.tmpl",
		})
	}
	return templateData, outputFiles
}

// getCacheConfig returns cache generator configuration
func getCacheConfig(config *Config, domain string, schemas []*ParsedSchema) (interface{}, []struct{ path, template string }) {
	if config.Generators.Cache.Interface == "" && config.Generators.Cache.Memory == "" && config.Generators.Cache.Redis == "" {
		return nil, nil
	}
	templateData := prepareCacheData(domain, schemas)
	var outputFiles []struct{ path, template string }
	if config.Generators.Cache.Interface != "" {
		outputFiles = append(outputFiles, struct{ path, template string }{
			filepath.Join(config.Output, domain, config.Generators.Cache.Interface), "cache.go.tmpl",
		})
	}
	if config.Generators.Cache.Memory != "" {
		outputFiles = append(outputFiles, struct{ path, template string }{
			filepath.Join(config.Output, domain, config.Generators.Cache.Memory), "cache_memory.go.tmpl",
		})
	}
	if config.Generators.Cache.Redis != "" {
		outputFiles = append(outputFiles, struct{ path, template string }{
			filepath.Join(config.Output, domain, config.Generators.Cache.Redis), "cache_redis.go.tmpl",
		})
	}
	return templateData, outputFiles
}

// getEventsConfig returns events generator configuration
func getEventsConfig(config *Config, domain string, schemas []*ParsedSchema) (interface{}, []struct{ path, template string }) {
	if config.Generators.Events.Interface == "" && config.Generators.Events.Redis == "" && config.Generators.Events.Nats == "" {
		return nil, nil
	}
	templateData := prepareEventsData(domain, schemas)
	var outputFiles []struct{ path, template string }
	if config.Generators.Events.Interface != "" {
		outputFiles = append(outputFiles, struct{ path, template string }{
			filepath.Join(config.Output, domain, config.Generators.Events.Interface), "events.go.tmpl",
		})
	}
	if config.Generators.Events.Redis != "" {
		outputFiles = append(outputFiles, struct{ path, template string }{
			filepath.Join(config.Output, domain, config.Generators.Events.Redis), "events_redis.go.tmpl",
		})
	}
	if config.Generators.Events.Nats != "" {
		outputFiles = append(outputFiles, struct{ path, template string }{
			filepath.Join(config.Output, domain, config.Generators.Events.Nats), "events_nats.go.tmpl",
		})
	}
	return templateData, outputFiles
}

// getHandlerConfig returns handler generator configuration (new OpenAPI-style)
func getHandlerConfig(config *Config, domain string, schemas []*ParsedSchema) (interface{}, []struct{ path, template string }) {
	// Check if handler generation is configured
	if config.Generators.Handlers == "" {
		return nil, nil
	}

	// For handler generation, we generate one handler per domain with all schemas
	if len(schemas) == 0 {
		return nil, nil
	}

	// Generate a single handler file for all schemas in the domain
	templateData := struct {
		Package string
		Imports []string
		Schemas []*ParsedSchema
		Domain  string
	}{
		Package: domain,
		Imports: []string{
			"context",
			"log/slog",
		},
		Schemas: schemas,
		Domain:  domain,
	}

	// Use the configured handlers output path
	outputFiles := []struct{ path, template string }{
		{filepath.Join(config.Output, domain, "server.impl.gen.go"), "handler.go.tmpl"},
	}
	return templateData, outputFiles
}

// prepareIndividualServiceData prepares template data for a single schema
func prepareIndividualServiceData(schema *ParsedSchema) interface{} {
	return struct {
		Package         string
		Name            string
		NameLower       string
		NamePlural      string
		NamePluralLower string
		XCodegen        *XCodegen
	}{
		Package:         schema.Domain,
		Name:            schema.Name,
		NameLower:       strings.ToLower(schema.Name),
		NamePlural:      Pluralize(schema.Name),
		NamePluralLower: strings.ToLower(Pluralize(schema.Name)),
		XCodegen:        schema.XCodegen,
	}
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
					Name    string
					Params  []string
					Returns string
				}{
					Name:    method.Name,
					Params:  method.Params,
					Returns: string(method.Returns),
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
	var entities []interface{}
	for _, schema := range schemas {
		ops := []string{}
		if schema.XCodegen.Cache.Enabled {
			ops = append(ops, "get", "set", "delete")
		}

		// Extract additional methods from x-codegen repository config
		var additionalMethods []interface{}
		if schema.XCodegen != nil && schema.XCodegen.Repository.AdditionalMethods != nil {
			for _, method := range schema.XCodegen.Repository.AdditionalMethods {
				additionalMethods = append(additionalMethods, struct {
					Name    string
					Params  []string
					Returns string
				}{
					Name:    method.Name,
					Params:  method.Params,
					Returns: string(method.Returns),
				})
			}
		}

		entities = append(entities, struct {
			Name              string
			Type              string
			Operations        []string
			AdditionalMethods []interface{}
			XCodegen          *XCodegen
		}{
			Name:              schema.Name,
			Type:              schema.Name, // Type is same as Name for cache
			Operations:        ops,
			AdditionalMethods: additionalMethods,
			XCodegen:          schema.XCodegen,
		})
	}

	return struct {
		Domain   string
		Package  string
		Entities []interface{}
		Imports  []string
	}{
		Domain:   domain,
		Package:  domain,
		Entities: entities,
		Imports:  []string{"github.com/google/uuid"},
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

// generateSQL generates SQL schema and query files from OpenAPI schemas.
func generateSQL(config *Config, schemas map[string]*ParsedSchema, fileWriter *FileWriter, log *slog.Logger) error {
	if config.Generators.SQL.SchemaDir == "" && config.Generators.SQL.QueryDir == "" {
		return nil
	}

	sqlConfig := config.Generators.SQL
	dialect := string(sqlConfig.Dialect)
	if dialect == "" {
		dialect = "postgresql" // Default to PostgreSQL
	}

	// Create generators
	schemaGen := NewSQLSchemaGenerator(dialect)
	queryGen := NewSQLQueryGenerator(Dialect(dialect))

	// Process each schema that has database configuration
	for name, parsedSchema := range schemas {
		if parsedSchema.XCodegen == nil || parsedSchema.XCodegen.Database.Table == "" {
			continue // Skip schemas without database configuration
		}

		log.Debug("Generating SQL for schema", slog.String("schema", name))

		// Get table name for file naming
		tableName := parsedSchema.XCodegen.Database.Table
		if tableName == "" {
			tableName = ToSnakeCase(name) + "s" // Default to plural
		}

		// Generate schema SQL
		if sqlConfig.SchemaDir != "" {
			schemaSQL, err := schemaGen.GenerateCreateTable(parsedSchema)
			if err != nil {
				log.Warn("Failed to generate schema", slog.String("schema", name), slog.String("error", err.Error()))
				continue
			}

			// Generate indices
			indices, err := schemaGen.GenerateIndices(parsedSchema)
			if err != nil {
				log.Warn("Failed to generate indices", slog.String("schema", name), slog.String("error", err.Error()))
			}

			// Combine schema and indices
			fullSQL := schemaSQL
			if len(indices) > 0 {
				fullSQL += "\n" + strings.Join(indices, "\n")
			}

			// Write schema file - use table name for filename
			schemaFile := filepath.Join(sqlConfig.SchemaDir, tableName+".sql")
			if err := fileWriter.WriteFile(schemaFile, []byte(fullSQL)); err != nil {
				return fmt.Errorf("failed to write schema file %s: %w", schemaFile, err)
			}
		}

		// Generate query SQL
		if sqlConfig.QueryDir != "" {
			querySQL, err := queryGen.GenerateQueries(parsedSchema)
			if err != nil {
				log.Warn("Failed to generate queries", slog.String("schema", name), slog.String("error", err.Error()))
				continue
			}

			// Write query file - use table name for filename
			queryFile := filepath.Join(sqlConfig.QueryDir, tableName+".sql")
			if err := fileWriter.WriteFile(queryFile, []byte(querySQL)); err != nil {
				return fmt.Errorf("failed to write query file %s: %w", queryFile, err)
			}
		}
	}

	return nil
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

// generateService generates service interfaces and implementations.
func generateService(config *Config, schemas map[string]*ParsedSchema, templates map[string]*template.Template, fileWriter *FileWriter, log *slog.Logger) error {
	log.Debug("Running service generator")

	// Generate individual service files per schema
	for _, schema := range schemas {
		if !NeedsService(schema) {
			continue
		}

		// Check if a manual service.go file exists
		manualServicePath := filepath.Join(config.Output, schema.Domain, "service.go")
		if _, err := os.Stat(manualServicePath); err == nil {
			log.Debug("Skipping service generation - manual service.go exists",
				slog.String("domain", schema.Domain),
				slog.String("path", manualServicePath))
			continue
		}

		log.Debug("Generating service",
			slog.String("domain", schema.Domain),
			slog.String("entity", schema.Name))

		// Prepare template data for this specific schema
		templateData := prepareIndividualServiceData(schema)

		// Generate service.gen.go
		servicePath := filepath.Join(config.Output, schema.Domain, "service.gen.go")
		if err := fileWriter.WriteTemplate(servicePath, templates["service.go.tmpl"], templateData); err != nil {
			return fmt.Errorf("failed to write service for %s: %w", schema.Name, err)
		}

		// Check if a manual handler.go file exists
		manualHandlerPath := filepath.Join(config.Output, schema.Domain, "handler.go")
		if _, err := os.Stat(manualHandlerPath); err == nil {
			log.Debug("Skipping handler generation - manual handler.go exists",
				slog.String("domain", schema.Domain),
				slog.String("path", manualHandlerPath))
		} else {
			// Generate server.gen.go
			handlerPath := filepath.Join(config.Output, schema.Domain, "server.gen.go")
			if err := fileWriter.WriteTemplate(handlerPath, templates["handler.go.tmpl"], templateData); err != nil {
				return fmt.Errorf("failed to write handler for %s: %w", schema.Name, err)
			}
		}
	}

	return nil
}

// generateDefaults generates configuration defaults.
func generateDefaults(config *Config, _ map[string]*ParsedSchema, _ map[string]*template.Template, fileWriter *FileWriter, log *slog.Logger) error {
	log.Debug("Running defaults generator")

	// Create parser
	parser := NewParser(filepath.Dir(config.Openapi))

	// Parse OpenAPI spec
	_, err := parser.ParseOpenAPISpec(config.Openapi)
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

// autoDetectDomains analyzes schemas and auto-detects domain configuration
func autoDetectDomains(schemas map[string]*ParsedSchema) map[string]struct {
	Schemas []string `json:"schemas,omitempty,omitzero" yaml:"schemas,omitempty"`
	Tags    []string `json:"tags,omitempty,omitzero" yaml:"tags,omitempty"`
} {
	domains := make(map[string]struct {
		Schemas []string `json:"schemas,omitempty,omitzero" yaml:"schemas,omitempty"`
		Tags    []string `json:"tags,omitempty,omitzero" yaml:"tags,omitempty"`
	})

	// Collect unique domains from schemas
	domainSchemas := make(map[string][]string)
	for name, schema := range schemas {
		if schema.Domain == "" {
			continue
		}

		domainSchemas[schema.Domain] = append(domainSchemas[schema.Domain], name)
	}

	// Convert to domain config struct
	for domain, schemaNames := range domainSchemas {
		// For now, we don't have tags in ParsedSchema, so we'll use domain names
		// In the future, this could be enhanced to extract tags from the OpenAPI spec
		domains[domain] = struct {
			Schemas []string `json:"schemas,omitempty,omitzero" yaml:"schemas,omitempty"`
			Tags    []string `json:"tags,omitempty,omitzero" yaml:"tags,omitempty"`
		}{
			Schemas: schemaNames,
		}
	}

	return domains
}
