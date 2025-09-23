package codegen

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/archesai/archesai/internal/parsers"
	"github.com/archesai/archesai/internal/templates"
)

// Generator is the interface that all code generators must implement.
type Generator interface {
	// Generate performs the code generation.
	Generate(ctx *GeneratorContext) error
}

// GeneratorContext holds common dependencies for all generators.
type GeneratorContext struct {
	Config     *CodegenConfig
	Parser     *parsers.Parser
	Schemas    []*parsers.JSONSchema
	Templates  map[string]*template.Template
	FileWriter *templates.FileWriter
}

// WriteTemplateForDomain writes a template for a specific domain.
func (ctx *GeneratorContext) WriteTemplateForDomain(
	domain string,
	templateName string,
	outputFile string,
	data interface{},
) error {
	outputPath := filepath.Join(*ctx.Config.Output, domain, outputFile)

	tmpl, ok := ctx.Templates[templateName]
	if !ok {
		return fmt.Errorf("template %s not found", templateName)
	}

	if err := ctx.FileWriter.WriteTemplate(outputPath, tmpl, data); err != nil {
		return fmt.Errorf("failed to write %s for %s: %w", outputFile, domain, err)
	}

	return nil
}

// ExtractCodegenExtension extracts and parses the x-codegen extension from a schema.
func ExtractCodegenExtension(schema *parsers.JSONSchema) (*CodegenExtension, error) {
	if schema == nil {
		return nil, fmt.Errorf("schema is nil")
	}

	// Try to get the x-codegen extension
	xCodegen, ok := schema.GetExtension("x-codegen")
	if !ok {
		// Return a default config if no extension found
		return &CodegenExtension{
			Domain: &schema.Name,
		}, nil
	}

	// Convert the raw extension to JSON then unmarshal to CodegenExtension
	jsonData, err := json.Marshal(xCodegen)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal x-codegen: %w", err)
	}

	var ext CodegenExtension
	if err := json.Unmarshal(jsonData, &ext); err != nil {
		return nil, fmt.Errorf("failed to unmarshal x-codegen: %w", err)
	}

	// Set domain if not specified
	if ext.Domain == nil {
		ext.Domain = &schema.Name
	}

	return &ext, nil
}

// =============================================================================
// Generic Generator System
// =============================================================================

// DataPreparer transforms parsed data into template data.
type DataPreparer[T any] func(domain string, schemas []*parsers.JSONSchema, operations []templates.OperationData, allSchemas map[string]*parsers.JSONSchema) T

// SchemaFilter determines which schemas to include.
type SchemaFilter func(*parsers.JSONSchema) bool

// GenericGenerator handles all template-based code generation.
type GenericGenerator[T any] struct {
	Name           string                                              // Generator name for logging
	TemplateName   string                                              // Template file name
	OutputFile     string                                              // Output file name pattern
	PrepareFunc    DataPreparer[T]                                     // Function to prepare template data
	Filter         SchemaFilter                                        // Which schemas to process
	Enabled        func(*CodegenConfig) bool                           // Check if generator is enabled
	UseOperations  bool                                                // Use operations instead of schemas
	CustomGenerate func(*GeneratorContext, *GenericGenerator[T]) error // Custom generation logic
	Parser         *parsers.Parser                                     // Parser instance
	logger         *slog.Logger                                        // Logger instance
}

// Generate implements the Generator interface.
func (g *GenericGenerator[T]) Generate(ctx *GeneratorContext) error {
	if !g.Enabled(ctx.Config) {
		g.logger.Debug(fmt.Sprintf("%s generator disabled", g.Name))
		return nil
	}

	g.logger.Debug(fmt.Sprintf("Running %s generator", g.Name))

	// Use custom generate if provided
	if g.CustomGenerate != nil {
		return g.CustomGenerate(ctx, g)
	}

	// Handle operation-based generators differently
	if g.UseOperations {
		return g.generateWithOperations(ctx)
	}

	// Some generators need operations even if they're primarily schema-based
	if g.Name == "types" {
		return g.generateWithSchemasAndOperations(ctx)
	}

	// Standard schema-based generation
	return g.generateWithSchemas(ctx)
}

func (g *GenericGenerator[T]) generateWithSchemas(ctx *GeneratorContext) error {
	// Group schemas by domain with filter
	domainSchemas := g.Parser.OpenAPI.GroupSchemasByDomain()

	g.logger.Debug(fmt.Sprintf("%s: Found %d domains", g.Name, len(domainSchemas)))

	// Get sorted domains for consistent output
	domains := g.Parser.OpenAPI.GetSortedDomains(domainSchemas)

	for _, domain := range domains {
		schemas := domainSchemas[domain]
		if len(schemas) == 0 {
			continue
		}

		// Skip schemas that are not actual domain entities
		if domain == "Base" || domain == "FilterNode" || domain == "Page" || domain == "Problem" {
			g.logger.Debug(fmt.Sprintf("Skipping %s for non-domain schema %s", g.Name, domain))
			continue
		}

		g.logger.Debug(fmt.Sprintf("Generating %s for domain", g.Name),
			slog.String("domain", domain),
			slog.Int("schemas", len(schemas)))

		// Prepare data using the specific preparer
		// Convert schemas slice to map for lookup
		allSchemas := make(map[string]*parsers.JSONSchema)
		for _, s := range ctx.Schemas {
			allSchemas[s.Name] = s
		}
		data := g.PrepareFunc(domain, schemas, nil, allSchemas)

		// Determine output file from config
		outputFile := g.getOutputFile(ctx.Config)
		if outputFile == "" {
			return fmt.Errorf("%s generator: output file not configured", g.Name)
		}

		// Write using template
		if err := ctx.WriteTemplateForDomain(domain, g.TemplateName, outputFile, data); err != nil {
			return fmt.Errorf("%s generation failed for %s: %w", g.Name, domain, err)
		}
	}

	return nil
}

func (g *GenericGenerator[T]) generateWithSchemasAndOperations(ctx *GeneratorContext) error {
	// Group schemas by domain with filter
	domainSchemas := g.Parser.OpenAPI.GroupSchemasByDomain()
	g.logger.Debug(fmt.Sprintf("%s (with ops) found %d domains", g.Name, len(domainSchemas)))

	// Get sorted domains for consistent output
	domains := g.Parser.OpenAPI.GetSortedDomains(domainSchemas)
	g.logger.Debug(fmt.Sprintf("Domains: %v", domains))

	// Extract operations from OpenAPI spec
	operations := g.Parser.OpenAPI.ExtractOperations()
	g.logger.Debug(fmt.Sprintf("Found %d operations", len(operations)))

	for _, domain := range domains {
		schemas := domainSchemas[domain]
		if len(schemas) == 0 {
			continue
		}

		// Skip schemas that are not actual domain entities
		if domain == "Base" || domain == "FilterNode" || domain == "Page" || domain == "Problem" {
			g.logger.Debug(
				fmt.Sprintf("Skipping %s generation for schemas without domain tags", g.Name),
				slog.Int("schemas", len(schemas)),
			)
			continue
		}

		// Filter operations for this domain
		domainOps := []templates.OperationData{}
		for _, op := range operations {
			// Check if operation belongs to this domain based on tags or path
			for _, tag := range op.Tags {
				if strings.ToLower(tag) == domain || strings.ToLower(tag) == domain+"s" {
					domainOps = append(domainOps, op)
					break
				}
			}
		}

		g.logger.Debug(fmt.Sprintf("Generating %s for domain", g.Name),
			slog.String("domain", domain),
			slog.Int("schemas", len(schemas)),
			slog.Int("operations", len(domainOps)))

		// Prepare data using the specific preparer with both schemas and operations
		// Convert schemas slice to map for lookup
		allSchemas := make(map[string]*parsers.JSONSchema)
		for _, s := range ctx.Schemas {
			allSchemas[s.Name] = s
		}
		data := g.PrepareFunc(domain, schemas, domainOps, allSchemas)

		// Determine output file from config
		outputFile := g.getOutputFile(ctx.Config)
		if outputFile == "" {
			return fmt.Errorf("%s generator: output file not configured", g.Name)
		}

		// Write using template
		if err := ctx.WriteTemplateForDomain(domain, g.TemplateName, outputFile, data); err != nil {
			return fmt.Errorf("%s generation failed for %s: %w", g.Name, domain, err)
		}
	}

	return nil
}

func (g *GenericGenerator[T]) generateWithOperations(ctx *GeneratorContext) error {
	// Extract operations from OpenAPI spec
	operations := g.Parser.OpenAPI.ExtractOperations()
	if len(operations) == 0 {
		return nil
	}

	// Group operations by domain
	domainOps := g.Parser.OpenAPI.GroupOperationsByDomain()

	for domain, ops := range domainOps {
		if len(ops) == 0 {
			continue
		}

		// Skip "default" domain as it's not a valid Go package name
		if domain == "default" {
			g.logger.Debug(
				fmt.Sprintf("Skipping %s generation for operations without domain tags", g.Name),
				slog.Int("operations", len(ops)),
			)
			continue
		}

		g.logger.Debug(fmt.Sprintf("Generating %s for domain", g.Name),
			slog.String("domain", domain),
			slog.Int("operations", len(ops)))

		// Prepare data using the specific preparer
		// Convert schemas slice to map for lookup
		allSchemas := make(map[string]*parsers.JSONSchema)
		for _, s := range ctx.Schemas {
			allSchemas[s.Name] = s
		}
		data := g.PrepareFunc(domain, nil, ops, allSchemas)

		// Determine output file from config
		outputFile := g.getOutputFile(ctx.Config)
		if outputFile == "" {
			return fmt.Errorf("%s generator: output file not configured", g.Name)
		}

		// Write using template
		outputPath := filepath.Join(*ctx.Config.Output, domain, outputFile)
		if err := ctx.FileWriter.WriteTemplate(outputPath, ctx.Templates[g.TemplateName], data); err != nil {
			return fmt.Errorf("%s generation failed for %s: %w", g.Name, domain, err)
		}
	}

	return nil
}

// getOutputFile gets the output file name from config based on generator type.
func (g *GenericGenerator[T]) getOutputFile(config *CodegenConfig) string {
	switch g.Name {
	case "types":
		return *config.Generators.Types
	case "service":
		if config.Generators.Service != nil {
			return *config.Generators.Service
		}
	case "echo_server":
		return *config.Generators.EchoServer
	case "cache":
		if config.Generators.Cache != nil {
			return *config.Generators.Cache.Interface
		}
	}
	return ""
}

// =============================================================================
// Generator Registry
// =============================================================================

// CreateGenerators creates all configured generators.
func CreateGenerators(parser *parsers.Parser, logger *slog.Logger) map[string]Generator {
	generators := map[string]Generator{
		"types": &GenericGenerator[*TypesTemplateInput]{
			Name:         "types",
			TemplateName: "types.tmpl",
			OutputFile:   "", // Will be set from config
			PrepareFunc:  prepareTypesData,
			Filter:       func(s *parsers.JSONSchema) bool { return true }, // All schemas
			Enabled:      func(c *CodegenConfig) bool { return c.Generators.Types != nil },
			Parser:       parser,
			logger:       logger,
		},
		"service": &GenericGenerator[*ServiceTemplateInput]{
			Name:         "service",
			TemplateName: "service.tmpl",
			OutputFile:   "", // Will be set from config
			PrepareFunc:  prepareServiceData,
			Filter:       func(s *parsers.JSONSchema) bool { return true }, // Service check done later
			Enabled: func(c *CodegenConfig) bool {
				return c.Generators.Service != nil
			},
			CustomGenerate: generateServiceWithManualCheck,
			Parser:         parser,
			logger:         logger,
		},
		"repository": &GenericGenerator[*RepositoryTemplateInput]{
			Name:        "repository",
			PrepareFunc: prepareRepositoryData,
			Filter:      func(s *parsers.JSONSchema) bool { return true }, // Repository check done later
			Enabled: func(c *CodegenConfig) bool {
				return c.Generators.Repository != nil &&
					(c.Generators.Repository.Interface != nil ||
						c.Generators.Repository.Postgres != nil ||
						c.Generators.Repository.Sqlite != nil)
			},
			CustomGenerate: generateRepositoryMultipleFiles,
			Parser:         parser,
			logger:         logger,
		},
		"echo_server": &GenericGenerator[*HandlerTemplateInput]{
			Name:          "echo_server",
			TemplateName:  "echo_server.tmpl",
			OutputFile:    "", // Will be set from config
			PrepareFunc:   prepareHandlerData,
			UseOperations: true,
			Enabled:       func(c *CodegenConfig) bool { return c.Generators.EchoServer != nil },
			Parser:        parser,
			logger:        logger,
		},
		"events": &GenericGenerator[*EventsTemplateInput]{
			Name:         "events",
			TemplateName: "events.tmpl",
			OutputFile:   "", // Will be set from config
			PrepareFunc:  prepareEventsData,
			Filter:       func(s *parsers.JSONSchema) bool { return true }, // Events check done later
			Enabled: func(c *CodegenConfig) bool {
				return c.Generators.Events != nil &&
					(c.Generators.Events.Nats != nil || c.Generators.Events.Redis != nil)
			},
			CustomGenerate: generateEventsMultipleFiles,
			Parser:         parser,
			logger:         logger,
		},
		"cache": &GenericGenerator[*CacheTemplateInput]{
			Name:         "cache",
			TemplateName: "cache.tmpl",
			OutputFile:   "", // Will be set from config
			PrepareFunc:  prepareCacheData,
			Filter:       func(s *parsers.JSONSchema) bool { return true }, // Cache check done later
			Enabled:      func(c *CodegenConfig) bool { return c.Generators.Cache != nil && c.Generators.Cache.Interface != nil },
			Parser:       parser,
			logger:       logger,
		},
	}

	// Add SQL generator separately as it has different logic
	generators["sql"] = NewSQLGenerator(parser, logger)

	return generators
}

// =============================================================================
// Custom Generation Functions
// =============================================================================

func generateServiceWithManualCheck(
	ctx *GeneratorContext,
	g *GenericGenerator[*ServiceTemplateInput],
) error {
	// Group schemas by domain, filtering only those that need services
	domainSchemas := g.Parser.OpenAPI.GroupSchemasByDomain()

	for domain, schemas := range domainSchemas {
		if len(schemas) == 0 {
			continue
		}

		// Skip "default" domain as it's not a valid Go package name
		if domain == "default" {
			g.logger.Debug("Skipping service generation for schemas without domain tags",
				slog.Int("schemas", len(schemas)))
			continue
		}

		// Filter schemas that have repository configuration
		serviceSchemas := []*parsers.JSONSchema{}
		for _, s := range schemas {
			ext, _ := ExtractCodegenExtension(s)
			if ext != nil && ext.Repository != nil && len(ext.Repository.Operations) > 0 {
				serviceSchemas = append(serviceSchemas, s)
			}
		}

		// Skip if no schemas need service generation
		if len(serviceSchemas) == 0 {
			g.logger.Debug("Skipping service generation - no schemas with repository config",
				slog.String("domain", domain))
			continue
		}

		// Check if a manual service.go file exists
		manualServicePath := filepath.Join(*ctx.Config.Output, domain, "service.go")
		if _, err := os.Stat(manualServicePath); err == nil {
			g.logger.Debug("Skipping service generation - manual service.go exists",
				slog.String("domain", domain),
				slog.String("path", manualServicePath))
			continue
		}

		g.logger.Debug("Generating service for domain",
			slog.String("domain", domain),
			slog.Int("schemas", len(serviceSchemas)))

		// Prepare template data
		// Convert schemas slice to map for lookup
		allSchemas := make(map[string]*parsers.JSONSchema)
		for _, s := range ctx.Schemas {
			allSchemas[s.Name] = s
		}
		data := g.PrepareFunc(domain, serviceSchemas, nil, allSchemas)

		// Generate service file
		if err := ctx.WriteTemplateForDomain(
			domain,
			"service.tmpl",
			*ctx.Config.Generators.Service,
			data,
		); err != nil {
			return fmt.Errorf("failed to write service for %s: %w", domain, err)
		}
	}

	return nil
}

func generateRepositoryMultipleFiles(
	ctx *GeneratorContext,
	g *GenericGenerator[*RepositoryTemplateInput],
) error {
	// Group schemas by domain, filtering only those that need repositories
	domainSchemas := g.Parser.OpenAPI.GroupSchemasByDomain()

	for domain, schemas := range domainSchemas {
		if len(schemas) == 0 {
			continue
		}

		// Skip "default" domain as it's not a valid Go package name
		if domain == "default" {
			g.logger.Debug("Skipping repository generation for schemas without domain tags",
				slog.Int("schemas", len(schemas)))
			continue
		}

		// Filter schemas that have repository configuration
		repoSchemas := []*parsers.JSONSchema{}
		for _, s := range schemas {
			ext, _ := ExtractCodegenExtension(s)
			if ext != nil && ext.Repository != nil && len(ext.Repository.Operations) > 0 {
				repoSchemas = append(repoSchemas, s)
			}
		}

		// Skip if no schemas need repository generation
		if len(repoSchemas) == 0 {
			g.logger.Debug("Skipping repository generation - no schemas with repository config",
				slog.String("domain", domain))
			continue
		}

		g.logger.Debug("Generating repository for domain",
			slog.String("domain", domain),
			slog.Int("schemas", len(repoSchemas)))

		// Prepare template data
		// Convert schemas slice to map for lookup
		allSchemas := make(map[string]*parsers.JSONSchema)
		for _, s := range ctx.Schemas {
			allSchemas[s.Name] = s
		}
		data := g.PrepareFunc(domain, repoSchemas, nil, allSchemas)

		// Generate repository interface if configured
		if ctx.Config.Generators.Repository.Interface != nil {
			if err := ctx.WriteTemplateForDomain(
				domain,
				"repository.tmpl",
				*ctx.Config.Generators.Repository.Interface,
				data,
			); err != nil {
				return fmt.Errorf("failed to write repository interface for %s: %w", domain, err)
			}
		}

		// Generate PostgreSQL implementation if configured
		if ctx.Config.Generators.Repository.Postgres != nil {
			if err := ctx.WriteTemplateForDomain(
				domain,
				"repository_postgres.tmpl",
				*ctx.Config.Generators.Repository.Postgres,
				data,
			); err != nil {
				return fmt.Errorf("failed to write PostgreSQL repository for %s: %w", domain, err)
			}
		}

		// Generate SQLite implementation if configured
		if ctx.Config.Generators.Repository.Sqlite != nil {
			if err := ctx.WriteTemplateForDomain(
				domain,
				"repository_sqlite.tmpl",
				*ctx.Config.Generators.Repository.Sqlite,
				data,
			); err != nil {
				return fmt.Errorf("failed to write SQLite repository for %s: %w", domain, err)
			}
		}
	}

	return nil
}

func generateEventsMultipleFiles(
	ctx *GeneratorContext,
	g *GenericGenerator[*EventsTemplateInput],
) error {
	// Group schemas by domain
	domainSchemas := g.Parser.OpenAPI.GroupSchemasByDomain()

	for domain, schemas := range domainSchemas {
		if len(schemas) == 0 {
			continue
		}

		g.logger.Debug("Generating events for domain",
			slog.String("domain", domain),
			slog.Int("schemas", len(schemas)))

		// Prepare template data
		// Convert schemas slice to map for lookup
		allSchemas := make(map[string]*parsers.JSONSchema)
		for _, s := range ctx.Schemas {
			allSchemas[s.Name] = s
		}
		data := g.PrepareFunc(domain, schemas, nil, allSchemas)

		// Generate base events file if configured
		if ctx.Config.Generators.Events.Interface != nil {
			if err := ctx.WriteTemplateForDomain(
				domain,
				"events.tmpl",
				*ctx.Config.Generators.Events.Interface,
				data,
			); err != nil {
				return fmt.Errorf("failed to write events interface for %s: %w", domain, err)
			}
		}

		// Generate NATS implementation if configured
		if ctx.Config.Generators.Events.Nats != nil {
			if err := ctx.WriteTemplateForDomain(
				domain,
				"events_nats.tmpl",
				*ctx.Config.Generators.Events.Nats,
				data,
			); err != nil {
				return fmt.Errorf("failed to write NATS events for %s: %w", domain, err)
			}
		}

		// Generate Redis implementation if configured
		if ctx.Config.Generators.Events.Redis != nil {
			if err := ctx.WriteTemplateForDomain(
				domain,
				"events_redis.tmpl",
				*ctx.Config.Generators.Events.Redis,
				data,
			); err != nil {
				return fmt.Errorf("failed to write Redis events for %s: %w", domain, err)
			}
		}
	}

	return nil
}

// =============================================================================
// Data Preparation Functions
// =============================================================================

// prepareTypesData transforms schemas to TypesTemplateInput for type generation.
func prepareTypesData(
	domain string,
	schemas []*parsers.JSONSchema,
	operations []templates.OperationData,
	allSchemas map[string]*parsers.JSONSchema,
) *TypesTemplateInput {
	data := &TypesTemplateInput{
		TemplateData: templates.TemplateData{
			Package: domain,
			Domain:  domain,
			Imports: []string{},
		},
		Schemas:   []templates.SchemaData{},
		Constants: []templates.ConstantDef{},
	}

	// Track which imports are needed
	needsTime := false
	needsUUID := false

	// Process each schema
	for _, schema := range schemas {
		if schema.Schema == nil {
			continue
		}

		// Handle enum constants
		if len(schema.Schema.Enum) > 0 {
			var enumValues []string
			for _, val := range schema.Schema.Enum {
				if val != nil {
					var strVal string
					if err := val.Decode(&strVal); err == nil {
						enumValues = append(enumValues, strVal)
					}
				}
			}
			if len(enumValues) > 0 {
				data.Constants = append(data.Constants, templates.ConstantDef{
					Name:   schema.Name,
					Values: enumValues,
				})
			}
			continue // Skip enum schemas from regular schema processing
		}

		// Extract fields using parser
		entityFields := schema.ExtractSchemaFields()

		var fields []templates.FieldData
		for _, ef := range entityFields {
			field := templates.FieldData{
				Name:          ef.Name,
				FieldName:     ef.FieldName,
				GoType:        ef.GoType,
				SQLCType:      ef.SQLCType,
				SQLCFieldName: ef.SQLCFieldName,
				JSONTag:       ef.JSONTag,
				YAMLTag:       ef.YAMLTag,
				Format:        ef.Format,
				Enum:          ef.Enum,
				Description:   ef.Description,
				Required:      ef.Required,
				Nullable:      ef.Nullable,
				DefaultValue:  ef.DefaultValue,
				IsEnumType:    ef.IsEnumType,
			}

			// If the field has enum values, update the GoType to use the enum type
			if len(ef.Enum) > 0 {
				enumTypeName := fmt.Sprintf("%s%s", schema.Name, ef.FieldName)
				field.GoType = enumTypeName
				field.IsEnumType = true
			}

			fields = append(fields, field)

			// Check for import requirements
			if strings.Contains(ef.GoType, "time.Time") {
				needsTime = true
			}
			if strings.Contains(ef.GoType, "uuid.UUID") {
				needsUUID = true
			}
		}

		description := ""
		if schema.Description != nil {
			description = *schema.Description
		}

		data.Schemas = append(data.Schemas, templates.SchemaData{
			Name:        schema.Name,
			Description: description,
			Fields:      fields,
		})

		// Process field-level enums to generate constants
		for _, field := range fields {
			if len(field.Enum) > 0 {
				// Create a type name for the enum (e.g., AccountProviderID)
				enumTypeName := fmt.Sprintf("%s%s", schema.Name, field.FieldName)

				// Add the enum type and constants
				data.Constants = append(data.Constants, templates.ConstantDef{
					Name:   enumTypeName,
					Values: field.Enum,
				})

				// TODO: Update the field's GoType to use the enum type instead of string
				// This would require updating the field after it's been added to the schema
			}
		}
	}

	// Process operations to generate request/response types
	for _, op := range operations {
		// Generate request body types
		if op.HasRequestBody && op.RequestBodySchema != "" {
			// Generate CreateXRequestBody, UpdateXRequestBody, etc.
			requestTypeName := fmt.Sprintf("%sRequestBody", strings.Title(op.OperationID))

			// Check if this request type already exists
			typeExists := false
			for _, s := range data.Schemas {
				if s.Name == requestTypeName {
					typeExists = true
					break
				}
			}

			if !typeExists {
				requestType := templates.SchemaData{
					Name:        requestTypeName,
					Description: fmt.Sprintf("Request body for %s", op.OperationID),
					Fields:      []templates.FieldData{}, // Fields would be extracted from OpenAPI
				}

				// For create/update operations, add common fields
				if strings.HasPrefix(op.OperationID, "create") ||
					strings.HasPrefix(op.OperationID, "update") {
					requestType.Fields = []templates.FieldData{
						{
							FieldName:   "Description",
							GoType:      "*string",
							JSONTag:     "description,omitempty",
							YAMLTag:     "description,omitempty",
							Description: fmt.Sprintf("The %s description", strings.ToLower(domain)),
						},
						{
							FieldName:   "Name",
							GoType:      "*string",
							JSONTag:     "name,omitempty",
							YAMLTag:     "name,omitempty",
							Description: fmt.Sprintf("The %s name", strings.ToLower(domain)),
						},
					}
				}

				data.Schemas = append(data.Schemas, requestType)
			}
		}
	}

	// Generate types for main entity schemas
	for _, schema := range schemas {
		if schema.Schema == nil || len(schema.Schema.Enum) > 0 {
			continue
		}

		// Only for main entities with x-codegen configuration
		ext, _ := ExtractCodegenExtension(schema)
		if ext != nil && ext.Repository != nil &&
			len(ext.Repository.Operations) > 0 {
			// Generate response type aliases
			// CreateXResponse = X
			createResponseAlias := templates.SchemaData{
				Name:        fmt.Sprintf("Create%sResponse", schema.Name),
				IsTypeAlias: true,
				TypeAlias:   fmt.Sprintf("= %s", schema.Name),
			}
			data.Schemas = append(data.Schemas, createResponseAlias)

			// UpdateXResponse = X
			updateResponseAlias := templates.SchemaData{
				Name:        fmt.Sprintf("Update%sResponse", schema.Name),
				IsTypeAlias: true,
				TypeAlias:   fmt.Sprintf("= %s", schema.Name),
			}
			data.Schemas = append(data.Schemas, updateResponseAlias)

			// GetXResponse = X
			getResponseAlias := templates.SchemaData{
				Name:        fmt.Sprintf("Get%sResponse", schema.Name),
				IsTypeAlias: true,
				TypeAlias:   fmt.Sprintf("= %s", schema.Name),
			}
			data.Schemas = append(data.Schemas, getResponseAlias)
		}
	}

	// Track which list params types will be generated by entity list operations
	entityListOps := make(map[string]bool)
	for _, schema := range schemas {
		ext, _ := ExtractCodegenExtension(schema)
		if ext != nil && ext.Repository != nil {
			for _, op := range ext.Repository.Operations {
				if op == "list" {
					// Will generate ListSchemaNameParams
					pluralName := schema.Name + "s"
					entityListOps[fmt.Sprintf("List%sParams", pluralName)] = true
					break
				}
			}
		}
	}

	// Generate operation parameter types (e.g., OauthAuthorizeParams)
	for _, op := range operations {
		// Capitalize the first letter of operation ID for the type name
		typeName := op.OperationID
		if len(typeName) > 0 {
			typeName = strings.ToUpper(typeName[:1]) + typeName[1:]
		}
		paramTypeName := fmt.Sprintf("%sParams", typeName)

		// Skip if this would duplicate an entity list params type
		if entityListOps[paramTypeName] {
			continue
		}

		// Check if operation has any parameters (query, path, or header)
		hasParams := false
		var paramFields []templates.FieldData

		// Process query parameters
		for _, param := range op.QueryParams {
			hasParams = true
			goType := getGoType(param.Type, param.Format)
			if !param.Required {
				goType = "*" + goType
			}
			// Clean up description - remove newlines and excessive whitespace
			description := strings.ReplaceAll(param.Description, "\n", " ")
			description = strings.Join(strings.Fields(description), " ")

			paramFields = append(paramFields, templates.FieldData{
				FieldName:   toCamelCase(param.Name),
				GoType:      goType,
				JSONTag:     param.Name + ",omitempty",
				YAMLTag:     param.Name + ",omitempty",
				Description: description,
			})
		}

		// Process header parameters
		for _, param := range op.HeaderParams {
			hasParams = true
			goType := getGoType(param.Type, param.Format)
			if !param.Required {
				goType = "*" + goType
			}
			// Clean up description - remove newlines and excessive whitespace
			description := strings.ReplaceAll(param.Description, "\n", " ")
			description = strings.Join(strings.Fields(description), " ")

			paramFields = append(paramFields, templates.FieldData{
				FieldName:   toCamelCase(param.Name),
				GoType:      goType,
				JSONTag:     param.Name + ",omitempty",
				YAMLTag:     param.Name + ",omitempty",
				Description: description,
			})
		}

		// Generate params type if operation has parameters
		if hasParams {
			paramsType := templates.SchemaData{
				Name:        paramTypeName,
				Description: fmt.Sprintf("Parameters for %s operation", op.OperationID),
				Fields:      paramFields,
			}
			data.Schemas = append(data.Schemas, paramsType)
		}
	}

	// Generate list-related types for entities with list operations
	for _, schema := range schemas {
		// Skip enums and non-entity schemas
		if schema.Schema == nil || len(schema.Schema.Enum) > 0 {
			continue
		}

		// Check if this entity has a list operation
		hasListOp := false
		ext, _ := ExtractCodegenExtension(schema)
		if ext != nil && ext.Repository != nil {
			for _, op := range ext.Repository.Operations {
				if op == "list" {
					hasListOp = true
					break
				}
			}
		}

		if !hasListOp {
			continue
		}

		// Use plural name for list operations
		pluralName := schema.Name + "s"

		// Generate ListXParams type
		listParamsType := templates.SchemaData{
			Name:        fmt.Sprintf("List%sParams", pluralName),
			Description: fmt.Sprintf("Parameters for listing %s", strings.ToLower(pluralName)),
			Fields: []templates.FieldData{
				{
					FieldName:   "Filter",
					GoType:      fmt.Sprintf("*List%sParamsFilter", pluralName),
					JSONTag:     "filter,omitempty",
					YAMLTag:     "filter,omitempty",
					Description: "Filter parameters",
				},
				{
					FieldName:   "Page",
					GoType:      fmt.Sprintf("*List%sParamsPage", pluralName),
					JSONTag:     "page,omitempty",
					YAMLTag:     "page,omitempty",
					Description: "Pagination parameters",
				},
				{
					FieldName:   "Sort",
					GoType:      fmt.Sprintf("*List%sParamsSort", pluralName),
					JSONTag:     "sort,omitempty",
					YAMLTag:     "sort,omitempty",
					Description: "Sort parameters",
				},
			},
		}
		data.Schemas = append(data.Schemas, listParamsType)

		// Generate ListXParamsFilter type
		filterType := templates.SchemaData{
			Name:        fmt.Sprintf("List%sParamsFilter", pluralName),
			Description: fmt.Sprintf("Filter %s by field values", strings.ToLower(pluralName)),
			Fields: []templates.FieldData{
				{
					FieldName:   "Filter",
					GoType:      "interface{}",
					JSONTag:     "filter",
					YAMLTag:     "filter",
					Description: "A recursive filter node that can be a condition or group",
				},
			},
		}
		data.Schemas = append(data.Schemas, filterType)

		// Generate ListXParamsPage type
		pageType := templates.SchemaData{
			Name: fmt.Sprintf("List%sParamsPage", pluralName),
			Description: fmt.Sprintf(
				"Pagination parameters for listing %s",
				strings.ToLower(pluralName),
			),
			Fields: []templates.FieldData{
				{
					FieldName:   "Number",
					GoType:      "int",
					JSONTag:     "number",
					YAMLTag:     "number",
					Description: "Page number (1-indexed)",
				},
				{
					FieldName:   "Size",
					GoType:      "int",
					JSONTag:     "size",
					YAMLTag:     "size",
					Description: "Page size (items per page)",
				},
			},
		}
		data.Schemas = append(data.Schemas, pageType)

		// Generate ListXParamsSort type
		sortType := templates.SchemaData{
			Name:        fmt.Sprintf("List%sParamsSort", pluralName),
			Description: fmt.Sprintf("Sort %s by field and order", strings.ToLower(pluralName)),
			Fields: []templates.FieldData{
				{
					FieldName:   "Field",
					GoType:      fmt.Sprintf("List%sParamsSortField", pluralName),
					JSONTag:     "field",
					YAMLTag:     "field",
					Description: "Field to sort by",
				},
				{
					FieldName:   "Order",
					GoType:      fmt.Sprintf("List%sParamsSortOrder", pluralName),
					JSONTag:     "order",
					YAMLTag:     "order",
					Description: "Sort order (asc or desc)",
				},
			},
		}
		data.Schemas = append(data.Schemas, sortType)

		// Note: Sort field and order types are generated as enums with constants below
		// No need for type aliases since the enum generation handles the type declaration

		// Generate sort field enum constants
		schemaFields := schema.ExtractSchemaFields()
		var sortFieldValues []string
		for _, field := range schemaFields {
			// Use the JSON tag name (camelCase) for the sort field value
			fieldName := field.JSONTag
			// Remove any tag modifiers like ,omitempty
			if idx := strings.Index(fieldName, ","); idx > 0 {
				fieldName = fieldName[:idx]
			}
			sortFieldValues = append(sortFieldValues, fieldName)
		}

		sortFieldEnum := templates.ConstantDef{
			Name:           fmt.Sprintf("List%sParamsSortField", pluralName),
			Values:         sortFieldValues,
			ConstantPrefix: pluralName, // Use plural name as prefix (e.g., "Sessions", "APIKeys")
		}
		data.Constants = append(data.Constants, sortFieldEnum)

		// Generate sort order enum constants
		sortOrderEnum := templates.ConstantDef{
			Name:           fmt.Sprintf("List%sParamsSortOrder", pluralName),
			Values:         []string{"asc", "desc"},
			ConstantPrefix: pluralName, // Use plural name as prefix
		}
		data.Constants = append(data.Constants, sortOrderEnum)

		// Generate DeleteXResponse (empty struct)
		deleteResponseType := templates.SchemaData{
			Name:        fmt.Sprintf("Delete%sResponse", schema.Name),
			Description: fmt.Sprintf("Response for deleting a %s", strings.ToLower(schema.Name)),
			Fields:      []templates.FieldData{},
		}
		data.Schemas = append(data.Schemas, deleteResponseType)

		// Generate ListXResponse type
		listResponseType := templates.SchemaData{
			Name:        fmt.Sprintf("List%sResponse", pluralName),
			Description: fmt.Sprintf("Response for listing %s", strings.ToLower(pluralName)),
			Fields: []templates.FieldData{
				{
					FieldName:   "Data",
					GoType:      fmt.Sprintf("[]*%s", schema.Name),
					JSONTag:     "data",
					YAMLTag:     "data",
					Description: fmt.Sprintf("List of %s", strings.ToLower(pluralName)),
				},
				{
					FieldName:   "Meta",
					GoType:      "ListMetadata",
					JSONTag:     "meta",
					YAMLTag:     "meta",
					Description: "Metadata about the list response",
				},
			},
		}
		data.Schemas = append(data.Schemas, listResponseType)
	}

	// Generate common ListMetadata type (only once per domain)
	needsListMetadata := false
	for _, schema := range schemas {
		ext, _ := ExtractCodegenExtension(schema)
		if ext != nil && ext.Repository != nil {
			for _, op := range ext.Repository.Operations {
				if op == "list" {
					needsListMetadata = true
					break
				}
			}
		}
		if needsListMetadata {
			break
		}
	}

	if needsListMetadata {
		listMetadataType := templates.SchemaData{
			Name:        "ListMetadata",
			Description: "Metadata for list responses",
			Fields: []templates.FieldData{
				{
					FieldName:   "Total",
					GoType:      "int64",
					JSONTag:     "total",
					YAMLTag:     "total",
					Description: "Total number of items",
				},
			},
		}
		data.Schemas = append(data.Schemas, listMetadataType)
	}

	// Add required imports
	if needsTime {
		data.Imports = append(data.Imports, "time")
	}
	if needsUUID {
		data.Imports = append(data.Imports, "github.com/google/uuid")
	}

	return data
}

// getGoType converts OpenAPI types to Go types
func getGoType(apiType, format string) string {
	switch apiType {
	case "string":
		switch format {
		case "uuid":
			return "uuid.UUID"
		case "date-time":
			return "time.Time"
		case "email":
			return "string"
		default:
			return "string"
		}
	case "integer":
		switch format {
		case "int32":
			return "int32"
		case "int64":
			return "int64"
		default:
			return "int"
		}
	case "number":
		switch format {
		case "float":
			return "float32"
		case "double":
			return "float64"
		default:
			return "float64"
		}
	case "boolean":
		return "bool"
	case "array":
		return "[]interface{}"
	case "object":
		return "interface{}"
	default:
		return "interface{}"
	}
}

// toCamelCase converts snake_case or kebab-case to CamelCase
func toCamelCase(s string) string {
	// Handle empty string
	if s == "" {
		return s
	}

	// Split by underscore or hyphen
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-'
	})

	// Capitalize each part
	for i := range parts {
		if parts[i] != "" {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}

	return strings.Join(parts, "")
}

// prepareServiceData transforms schemas to ServiceTemplateInput for service generation.
func prepareServiceData(
	domain string,
	schemas []*parsers.JSONSchema,
	operations []templates.OperationData,
	allSchemas map[string]*parsers.JSONSchema,
) *ServiceTemplateInput {
	// For services, we generate for each individual schema
	// Most domains have only one entity, but we handle multiple entities
	if len(schemas) == 0 {
		return nil
	}

	// For now, take the first schema (most domains have only one entity)
	// TODO: Handle multiple entities in a single domain
	schema := schemas[0]

	data := &ServiceTemplateInput{
		TemplateData: templates.TemplateData{
			Package: domain,
			Domain:  domain,
		},
	}

	// Create entity data from schema
	entity := templates.EntityData{
		Domain:          domain,
		Package:         domain,
		Name:            schema.Name,
		NameLower:       strings.ToLower(schema.Name),
		NamePlural:      schema.Name + "s",
		NamePluralLower: strings.ToLower(schema.Name + "s"),
		Type:            schema.Name,
		CodegenExtension: func() *CodegenExtension {
			ext, _ := ExtractCodegenExtension(schema)
			return ext
		}(), // Extract the x-codegen config
	}

	// Extract fields
	fields := schema.ExtractSchemaFields()
	entity.Fields = fields

	data.Entities = []templates.EntityData{entity}

	return data
}

// prepareRepositoryData transforms schemas to RepositoryTemplateInput for repository generation.
func prepareRepositoryData(
	domain string,
	schemas []*parsers.JSONSchema,
	operations []templates.OperationData,
	allSchemas map[string]*parsers.JSONSchema,
) *RepositoryTemplateInput {

	data := &RepositoryTemplateInput{
		TemplateData: templates.TemplateData{
			Package: domain,
			Domain:  domain,
		},
	}

	// Create entities from schemas
	for _, schema := range schemas {
		// Extract x-codegen configuration
		ext, _ := ExtractCodegenExtension(schema)

		// Determine repository name
		repositoryName := "Repository" // default
		if ext != nil && ext.Repository != nil && ext.Repository.Name != nil {
			repositoryName = *ext.Repository.Name
		}

		entity := templates.EntityData{
			Domain:            domain,
			Package:           domain,
			Name:              schema.Name,
			NameLower:         strings.ToLower(schema.Name),
			NamePlural:        schema.Name + "s",
			NamePluralLower:   strings.ToLower(schema.Name + "s"),
			Type:              schema.Name,
			RepositoryName:    repositoryName,
			CodegenExtension:  ext,
			AdditionalMethods: []templates.MethodData{}, // Initialize empty
		}

		// Extract fields
		fields := schema.ExtractSchemaFields()

		// Fix up enum types - if a field has enum values, update its GoType and set IsEnumType
		for i := range fields {
			if len(fields[i].Enum) > 0 {
				// Generate the enum type name (SchemaName + FieldName)
				enumTypeName := fmt.Sprintf("%s%s", schema.Name, fields[i].FieldName)
				fields[i].GoType = enumTypeName
				fields[i].IsEnumType = true
			}
		}

		entity.Fields = fields

		// Extract additional methods from x-codegen
		// ext already extracted above
		if ext != nil && ext.Repository != nil {
			if len(ext.Repository.AdditionalMethods) > 0 {
				for _, additionalMethod := range ext.Repository.AdditionalMethods {
					methodData := templates.MethodData{
						Name: additionalMethod.Name,
						Description: fmt.Sprintf(
							"%s performs %s operation",
							additionalMethod.Name,
							additionalMethod.Name,
						),
						Parameters: []templates.ParamData{},
						Returns:    []string{},
					}

					// Convert parameters
					for _, param := range additionalMethod.Params {
						paramType := ""

						// Look up the field type from the schema
						for _, field := range fields {
							fieldLower := strings.ToLower(field.Name)
							paramLower := strings.ToLower(param)

							// Exact match (case-insensitive)
							if fieldLower == paramLower {
								paramType = field.GoType
								break
							}

							// Try adding "ID" suffix to param to match field
							// e.g., "provider" matches "providerID"
							if fieldLower == paramLower+"id" {
								paramType = field.GoType
								break
							}

							// For compound params like "providerAccountID", extract and match the suffix
							// e.g., "providerAccountID" -> try to match "accountID"
							if strings.Contains(paramLower, "account") &&
								strings.HasSuffix(paramLower, "id") {
								// Try to find accountID field
								if fieldLower == "accountid" {
									paramType = field.GoType
									break
								}
							}
						}

						// If not found, it's likely a reference to another entity's ID
						if paramType == "" {
							if strings.HasSuffix(param, "ID") {
								paramType = "uuid.UUID"
							} else {
								paramType = "interface{}" // Fallback
							}
						}

						methodData.Parameters = append(methodData.Parameters, templates.ParamData{
							Name: param,
							Type: paramType,
						})
					}

					// Convert returns
					if additionalMethod.Returns == "multiple" {
						methodData.Returns = []string{
							fmt.Sprintf("[]*%s", schema.Name),
							"error",
						}
					} else if additionalMethod.Returns == "single" {
						methodData.Returns = []string{fmt.Sprintf("*%s", schema.Name), "error"}
					} else {
						methodData.Returns = []string{"error"}
					}

					entity.AdditionalMethods = append(entity.AdditionalMethods, methodData)
				}
			}
		}

		data.Entities = append(data.Entities, entity)
	}

	return data
}

// prepareHandlerData transforms operations to HandlerTemplateInput for HTTP handler generation.
func prepareHandlerData(
	domain string,
	schemas []*parsers.JSONSchema,
	operations []templates.OperationData,
	allSchemas map[string]*parsers.JSONSchema,
) *HandlerTemplateInput {
	data := &HandlerTemplateInput{
		TemplateData: templates.TemplateData{
			Package: domain,
			Domain:  domain,
		},
		Operations: []templates.OperationData{},
	}

	// Operations are already templates.OperationData
	for _, op := range operations {
		// Create Go-friendly name (capitalize first letter)
		goName := op.OperationID
		if len(goName) > 0 {
			goName = strings.ToUpper(goName[:1]) + goName[1:]
		}

		operation := templates.OperationData{
			Name:              op.OperationID,
			GoName:            goName,
			Method:            op.Method,
			Path:              op.Path,
			Description:       op.Description,
			HasRequestBody:    op.HasRequestBody,
			RequestBodySchema: op.RequestBodySchema,
			ResponseType:      op.ResponseType,
			Tags:              op.Tags,
			HasBearerAuth:     op.HasBearerAuth,
			HasCookieAuth:     op.HasCookieAuth,
			Security:          op.Security,
		}

		// Convert parameters - use the Type field that was set during parsing
		for _, param := range op.PathParams {
			paramType := param.Type
			if paramType == "" {
				paramType = inferParamType(param) // Fallback to inference
			}
			operation.PathParams = append(operation.PathParams, templates.ParamData{
				Name:        param.Name,
				Type:        paramType,
				Required:    param.Required,
				Description: param.Description,
			})
		}

		for _, param := range op.QueryParams {
			paramType := param.Type
			if paramType == "" {
				paramType = inferParamType(param) // Fallback to inference
			}
			operation.QueryParams = append(operation.QueryParams, templates.ParamData{
				Name:        param.Name,
				Type:        paramType,
				Required:    param.Required,
				Description: param.Description,
			})
		}

		for _, param := range op.HeaderParams {
			paramType := param.Type
			if paramType == "" {
				paramType = inferParamType(param) // Fallback to inference
			}
			operation.HeaderParams = append(operation.HeaderParams, templates.ParamData{
				Name:        param.Name,
				Type:        paramType,
				Required:    param.Required,
				Description: param.Description,
			})
		}

		data.Operations = append(data.Operations, operation)
	}

	return data
}

// prepareEventsData transforms schemas to EventsTemplateInput for event generation.
func prepareEventsData(
	domain string,
	schemas []*parsers.JSONSchema,
	operations []templates.OperationData,
	allSchemas map[string]*parsers.JSONSchema,
) *EventsTemplateInput {

	data := &EventsTemplateInput{
		TemplateData: templates.TemplateData{
			Package: domain,
			Domain:  domain,
		},
		Events:     []templates.EventData{},
		EventTypes: []string{"created", "updated", "deleted"},
	}

	// Create events for each schema/entity
	for _, schema := range schemas {
		nameLower := strings.ToLower(schema.Name)
		for _, eventType := range data.EventTypes {
			data.Events = append(data.Events, templates.EventData{
				Name:        fmt.Sprintf("%s_%s", nameLower, eventType),
				Type:        eventType,
				PayloadType: schema.Name,
			})
		}
	}

	return data
}

// prepareCacheData transforms schemas to CacheTemplateInput for cache generation.
func prepareCacheData(
	domain string,
	schemas []*parsers.JSONSchema,
	operations []templates.OperationData,
	allSchemas map[string]*parsers.JSONSchema,
) *CacheTemplateInput {

	var entities []templates.EntityData
	for _, schema := range schemas {
		// Extract x-codegen configuration
		ext, _ := ExtractCodegenExtension(schema)

		// Determine repository name
		repositoryName := "Repository" // default
		if ext != nil && ext.Repository != nil && ext.Repository.Name != nil {
			repositoryName = *ext.Repository.Name
		}

		entity := templates.EntityData{
			Domain:            domain,
			Package:           domain,
			Name:              schema.Name,
			NameLower:         strings.ToLower(schema.Name),
			NamePlural:        schema.Name + "s",
			NamePluralLower:   strings.ToLower(schema.Name + "s"),
			Type:              schema.Name,
			RepositoryName:    repositoryName,
			CodegenExtension:  ext,
			AdditionalMethods: []templates.MethodData{}, // Initialize empty
		}

		// Extract fields
		fields := schema.ExtractSchemaFields()

		// Fix up enum types - if a field has enum values, update its GoType and set IsEnumType
		for i := range fields {
			if len(fields[i].Enum) > 0 {
				// Generate the enum type name (SchemaName + FieldName)
				enumTypeName := fmt.Sprintf("%s%s", schema.Name, fields[i].FieldName)
				fields[i].GoType = enumTypeName
				fields[i].IsEnumType = true
			}
		}

		entity.Fields = fields

		entities = append(entities, entity)
	}

	return &CacheTemplateInput{
		TemplateData: templates.TemplateData{
			Package: domain,
			Domain:  domain,
		},
		Entities:  entities,
		CacheType: "redis", // Default cache type
		TTL:       3600,    // Default TTL: 1 hour
	}
}

// =============================================================================
// Helper Functions
// =============================================================================

// inferParamType infers the Go type from an operation parameter.
func inferParamType(param templates.ParamData) string {
	// Check if Type field exists and is populated
	// This is a placeholder - the Type field should be populated by the parser
	// For now, we'll default to string
	// TODO: Update once templates.ParamData has Type field
	if param.Name == "id" {
		return "uuid.UUID"
	}
	// Default to string for now
	return "string"
}
