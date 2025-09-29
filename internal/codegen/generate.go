// Package codegen provides unified code generation from OpenAPI schemas with x-codegen extensions.
package codegen

import (
	"bytes"
	"fmt"
	"go/format"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/speakeasy-api/openapi/jsonschema/oas3"

	"github.com/archesai/archesai/internal/parsers"
	"github.com/archesai/archesai/internal/shared/logger"
)

// Schema type constants
const (
	schemaTypeEntity      = "entity"
	schemaTypeAggregate   = "aggregate"
	schemaTypeValueObject = "valueobject"
)

// Generate is the main generation function that orchestrates all code generation
func Generate(specPath string, opts Configuration) (string, error) {
	// 1. Setup logging
	logLevel := os.Getenv("ARCHESAI_LOGGING_LEVEL")
	if logLevel == "" {
		logLevel = "error"
	}
	log := logger.New(logger.Config{Level: logLevel, Pretty: true})

	// 2. Initialize global state
	globalState = NewGlobalState()
	globalState.Options = opts

	// 3. Create file writer
	fileWriter := NewFileWriter()
	fileWriter.WithOverwrite(true)
	fileWriter.WithHeader(DefaultHeader())
	globalState.FileWriter = fileWriter

	// 4. Load templates
	templates, err := LoadTemplates()
	if err != nil {
		return "", fmt.Errorf("failed to load templates: %w", err)
	}
	globalState.Templates = templates

	// 5. Parse the OpenAPI specification
	log.Info("Parsing OpenAPI specification", slog.String("path", specPath))
	openAPISchema, warnings, err := parsers.ParseOpenAPI(specPath)
	if err != nil {
		return "", fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	// Log any warnings
	for _, warning := range warnings {
		log.Warn("OpenAPI warning", slog.String("warning", warning))
	}

	// 6. Initialize global state with parsed data
	if err := globalState.Initialize(openAPISchema, templates, fileWriter, opts); err != nil {
		return "", fmt.Errorf("failed to initialize state: %w", err)
	}

	log.Info("Initialized state",
		slog.Int("schemas", len(globalState.ProcessedSchemas)),
		slog.Int("operations", len(globalState.Operations)))

	// Buffer to collect all output
	var output bytes.Buffer

	log.Info("Generating models (DTOs, entities, value objects)")
	if err := GenerateModels(globalState); err != nil {
		return "", fmt.Errorf("failed to generate models: %w", err)
	}

	log.Info("Generating repositories")
	if err := GenerateRepositories(globalState); err != nil {
		return "", fmt.Errorf("failed to generate repositories: %w", err)
	}

	log.Info("Generating command and query handlers")
	if err := GenerateCommandQueryHandlers(globalState); err != nil {
		return "", fmt.Errorf("failed to generate handlers: %w", err)
	}

	log.Info("Generating handlers")
	if err := GenerateControllers(globalState); err != nil {
		return "", fmt.Errorf("failed to generate handlers: %w", err)
	}
	log.Info("Generating events")
	if err := GenerateEvents(globalState); err != nil {
		return "", fmt.Errorf("failed to generate events: %w", err)
	}

	// 10. Format output if needed
	outputStr := output.String()
	if outputStr != "" {
		formatted, err := format.Source([]byte(outputStr))
		if err != nil {
			log.Warn("Failed to format output", slog.String("error", err.Error()))
		} else {
			outputStr = string(formatted)
		}
	}

	log.Info("Code generation completed successfully")
	return outputStr, nil
}

// Configuration types and defaults
type Configuration struct {
	// Path to the OpenAPI spec
	SpecPath string
}

// GenerateRepositories generates all repository interfaces and implementations
func GenerateRepositories(state *GlobalState) error {
	for name, processed := range state.ProcessedSchemas {
		// Only generate repositories for schemas with x-codegen.repository config
		if processed.XCodegen != nil && processed.XCodegen.Repository != nil {
			if err := generateRepositoryForSchema(state, processed.Schema, name); err != nil {
				return fmt.Errorf("failed to generate repository for %s: %w", name, err)
			}
		}
	}
	return nil
}

// generateRepositoryForSchema generates repository interface and implementations for a schema
func generateRepositoryForSchema(state *GlobalState, schema *oas3.Schema, name string) error {
	// Determine entity package based on schema type
	entityPackage := "entities"
	title := name
	xcodegen := state.XCodegenMap[title]

	if xcodegen != nil {
		switch xcodegen.Type {
		case "aggregate":
			entityPackage = "aggregates"
		case "valueobject":
			// Value objects don't get repositories
			return nil
		}
	}

	// Extract additional methods from x-codegen configuration
	var additionalMethods []map[string]interface{}
	if xcodegen != nil && xcodegen.Repository != nil {
		for _, method := range xcodegen.Repository.AdditionalMethods {
			var params []map[string]string
			for _, param := range method.Params {
				params = append(params, map[string]string{
					"Name": param,
					"Type": "string", // Default type, should be improved
				})
			}

			returns := []string{"error"}
			switch method.Returns {
			case "single":
				returns = []string{"*" + title, "error"}
			case "multiple":
				returns = []string{"[]*" + title, "error"}
			}

			additionalMethods = append(additionalMethods, map[string]interface{}{
				"Name":       method.Name,
				"Parameters": params,
				"Returns":    returns,
			})
		}
	}

	// Generate repository interface
	data := map[string]interface{}{
		"Package": "repositories",
		"Entities": []map[string]interface{}{
			{
				"Name":              title,
				"Type":              title,
				"EntityPackage":     entityPackage,
				"AdditionalMethods": additionalMethods,
			},
		},
	}

	// Generate interface in repositories folder
	outputPath := filepath.Join(
		"internal/core/repositories",
		strings.ToLower(title)+".gen.go",
	)

	tmpl, ok := state.GetTemplate("repository.tmpl")
	if !ok {
		return fmt.Errorf("repository template not found")
	}

	if err := state.FileWriter.WriteTemplate(outputPath, tmpl, data); err != nil {
		return err
	}

	// Generate concrete implementations with different package
	implData := map[string]interface{}{
		"Package": "repositories", // Implementation package
		"Entities": []map[string]interface{}{
			{
				"Name":              title,
				"Type":              title,
				"EntityPackage":     entityPackage, // Import from entities/aggregates package
				"AdditionalMethods": additionalMethods,
			},
		},
	}

	// PostgreSQL
	if tmpl, ok := state.GetTemplate("repository_postgres.tmpl"); ok {
		outputPath := filepath.Join(
			"internal/infrastructure/persistence/postgres/repositories",
			strings.ToLower(title)+"_repository.gen.go",
		)
		if err := state.FileWriter.WriteTemplate(outputPath, tmpl, implData); err != nil {
			return fmt.Errorf("failed to generate PostgreSQL repository: %w", err)
		}
	}

	// SQLite
	if tmpl, ok := state.GetTemplate("repository_sqlite.tmpl"); ok {
		outputPath := filepath.Join(
			"internal/infrastructure/persistence/sqlite/repositories",
			strings.ToLower(title)+"_repository.gen.go",
		)
		if err := state.FileWriter.WriteTemplate(outputPath, tmpl, implData); err != nil {
			return fmt.Errorf("failed to generate SQLite repository: %w", err)
		}
	}

	return nil
}

// GenerateCommandQueryHandlers generates command and query handlers
func GenerateCommandQueryHandlers(state *GlobalState) error {
	// Generate command handlers (includes types)
	if err := generateCommandHandlers(state); err != nil {
		return fmt.Errorf("failed to generate command handlers: %w", err)
	}

	// Generate query handlers (includes types)
	if err := generateQueryHandlers(state); err != nil {
		return fmt.Errorf("failed to generate query handlers: %w", err)
	}

	return nil
}

// generateCommandHandlers generates individual command handler files for each operation
func generateCommandHandlers(state *GlobalState) error {
	tmpl, ok := state.GetTemplate("single_command_handler.tmpl")
	if !ok {
		// Template doesn't exist yet, skip
		return nil
	}

	// Group operations by their domain/tag
	operationsByDomain := make(map[string][]parsers.OperationDef)
	for _, op := range state.Operations {
		if len(op.Tags) > 0 {
			domain := strings.ToLower(op.Tags[0])
			operationsByDomain[domain] = append(operationsByDomain[domain], op)
		}
	}

	// For each domain, generate command handlers for write operations
	for domain, operations := range operationsByDomain {
		for _, op := range operations {
			// Only generate command handlers for write operations
			if op.Method == "POST" || op.Method == "PUT" || op.Method == "PATCH" ||
				op.Method == "DELETE" {
				// Determine command type based on operation
				var commandType string
				switch {
				case strings.HasPrefix(op.OperationID, "create"):
					commandType = "Create"
				case strings.HasPrefix(op.OperationID, "update"):
					commandType = "Update"
				case strings.HasPrefix(op.OperationID, "delete"):
					commandType = "Delete"
				default:
					// Custom command, use operation ID with title case
					commandType = strings.Title(op.OperationID)
				}

				// Get the entity name from the schema if available
				entityName := Singularize(strings.Title(domain))
				entityNameLower := strings.ToLower(entityName)

				// Create template data
				data := map[string]interface{}{
					"Package":         domain,
					"CommandType":     commandType,
					"EntityName":      entityName,
					"EntityNameLower": entityNameLower,
				}

				// Generate the command handler file
				outputPath := filepath.Join(
					"internal/application/commands",
					domain,
					fmt.Sprintf("%s_%s.gen.go", strings.ToLower(commandType), entityNameLower),
				)

				// Write the handler file
				if err := state.FileWriter.WriteTemplate(outputPath, tmpl, data); err != nil {
					return fmt.Errorf(
						"failed to generate command handler for %s: %w",
						op.OperationID,
						err,
					)
				}
			}
		}
	}

	return nil
}

// generateQueryHandlers generates individual query handler files for each operation
func generateQueryHandlers(state *GlobalState) error {
	tmpl, ok := state.GetTemplate("single_query_handler.tmpl")
	if !ok {
		// Template doesn't exist yet, skip
		return nil
	}

	// Group operations by their domain/tag
	operationsByDomain := make(map[string][]parsers.OperationDef)
	for _, op := range state.Operations {
		if len(op.Tags) > 0 {
			domain := strings.ToLower(op.Tags[0])
			operationsByDomain[domain] = append(operationsByDomain[domain], op)
		}
	}

	// For each domain, generate query handlers for read operations
	for domain, operations := range operationsByDomain {
		for _, op := range operations {
			// Only generate query handlers for read operations
			if op.Method == "GET" {
				// Determine query type based on operation
				var queryType string
				var usesPluralName bool
				switch {
				case strings.HasPrefix(op.OperationID, "list"):
					queryType = "List"
					usesPluralName = true
				case strings.HasPrefix(op.OperationID, "get"):
					queryType = "Get"
				case strings.HasPrefix(op.OperationID, "search"):
					queryType = "Search"
					usesPluralName = true
				default:
					// Custom query, use operation ID with title case
					queryType = strings.Title(op.OperationID)
				}

				// Get the entity name from the schema if available
				entityName := Singularize(strings.Title(domain))
				entityNameLower := strings.ToLower(entityName)
				entityNamePlural := strings.Title(domain)

				// Create template data
				data := map[string]interface{}{
					"Package":          domain,
					"QueryType":        queryType,
					"EntityName":       entityName,
					"EntityNameLower":  entityNameLower,
					"EntityNamePlural": entityNamePlural,
				}

				// Generate the query handler file
				var fileName string
				if usesPluralName {
					fileName = fmt.Sprintf(
						"%s_%s.gen.go",
						strings.ToLower(queryType),
						strings.ToLower(entityNamePlural),
					)
				} else {
					fileName = fmt.Sprintf("%s_%s.gen.go", strings.ToLower(queryType), entityNameLower)
				}

				outputPath := filepath.Join(
					"internal/application/queries",
					domain,
					fileName,
				)

				// Write the handler file
				if err := state.FileWriter.WriteTemplate(outputPath, tmpl, data); err != nil {
					return fmt.Errorf(
						"failed to generate query handler for %s: %w",
						op.OperationID,
						err,
					)
				}
			}
		}
	}

	return nil
}

// GenerateControllers generates all HTTP controllers grouped by domain
func GenerateControllers(state *GlobalState) error {
	// Group operations by their first tag (domain)
	operationsByDomain := make(map[string][]parsers.OperationDef)

	for _, op := range state.Operations {
		domain := ""
		if len(op.Tags) > 0 {
			domain = op.Tags[0]
		} else {
			// Try to extract domain from path
			parts := strings.Split(strings.Trim(op.Path, "/"), "/")
			if len(parts) > 1 {
				domain = parts[1]
			} else if len(parts) > 0 {
				domain = parts[0]
			}
		}

		if domain != "" {
			operationsByDomain[domain] = append(operationsByDomain[domain], op)
		}
	}

	// Track schema types to determine imports
	schemaTypeMap := make(map[string]string) // schema name -> package
	for name, processed := range state.ProcessedSchemas {
		if processed.XCodegen != nil {
			switch processed.XCodegen.Type {
			case schemaTypeEntity:
				schemaTypeMap[name] = "entities"
			case schemaTypeAggregate:
				schemaTypeMap[name] = "aggregates"
			case schemaTypeValueObject:
				schemaTypeMap[name] = "valueobjects"
			default:
				schemaTypeMap[name] = "dto"
			}
		} else {
			// No x-codegen means it's a DTO
			schemaTypeMap[name] = "dto"
		}
	}

	// List of domains that have command/query handlers
	cqrsDomains := map[string]bool{
		"accounts":      true,
		"apikeys":       true,
		"artifacts":     true,
		"invitations":   true,
		"labels":        true,
		"members":       true,
		"organizations": true,
		"pipelines":     true,
		"runs":          true,
		"tools":         true,
		"users":         true,
	}

	// Generate a handler file for each domain
	for domain, operations := range operationsByDomain {
		// Skip empty domains
		if domain == "" {
			continue
		}

		// Skip domains that don't have CQRS handlers for now
		// TODO: Generate simpler controllers for non-CQRS domains
		if !cqrsDomains[strings.ToLower(domain)] {
			continue
		}

		// Capitalize first letter of domain for title case
		domainTitle := strings.ToUpper(domain[:1]) + domain[1:]

		// Get singular form for entity name (e.g., "Labels" -> "Label")
		domainSingular := Singularize(domainTitle)

		// Track which imports are needed for this handler
		importsNeeded := make(map[string]bool)

		// Process operations to split parameters by type
		processedOps := make([]map[string]interface{}, 0, len(operations))
		for _, op := range operations {
			var pathParams []parsers.ParamDef
			var queryParams []parsers.ParamDef
			var headerParams []parsers.ParamDef

			for _, param := range op.Parameters {
				switch param.In {
				case "path":
					pathParams = append(pathParams, param)
				case "query":
					queryParams = append(queryParams, param)
				case "header":
					headerParams = append(headerParams, param)
				}
			}

			// Determine response type and track imports
			responseType := ""
			responsePackage := ""
			var successResponse map[string]interface{}
			for _, resp := range op.Responses {
				if resp.IsSuccess && resp.Schema != "" {
					responseType = resp.Schema
					// Track which package this schema comes from
					if pkg, ok := schemaTypeMap[resp.Schema]; ok {
						importsNeeded[pkg] = true
						responsePackage = pkg
					}
					// Create success response data for template
					successResponse = map[string]interface{}{
						"StatusCode":  resp.StatusCode,
						"Schema":      resp.Schema,
						"Package":     responsePackage,
						"IsArray":     resp.IsArray,
						"Description": resp.Description,
					}
					break
				}
			}

			// If no response type found, try to determine from domain
			if responseType == "" && domain != "" {
				// Try to find the schema for the domain (singular form)
				domainSingular := strings.TrimSuffix(domainTitle, "s") // Simple singularize
				if pkg, ok := schemaTypeMap[domainSingular]; ok {
					responsePackage = pkg
					importsNeeded[pkg] = true
				}
			}

			// Note: Request body handling removed - operations no longer have inline request bodies
			// They should be defined as schemas in components instead

			processedOps = append(processedOps, map[string]interface{}{
				"Name":                op.Name,
				"GoName":              op.GoName,
				"Method":              op.Method,
				"Path":                op.Path,
				"Description":         op.Description,
				"OperationID":         op.OperationID,
				"Tags":                op.Tags,
				"PathParams":          pathParams,
				"QueryParams":         queryParams,
				"HeaderParams":        headerParams,
				"RequestBodyRequired": op.RequestBodyRequired,
				"Responses":           op.Responses,
				"Security":            op.Security,
				"ResponseType":        responseType,
				"ResponsePackage":     responsePackage,
				"SuccessResponse":     successResponse,
			})
		}

		// Build imports list - no longer needed as we have standard imports in template
		var imports []map[string]string

		data := map[string]interface{}{
			"Package":             "controllers",
			"Domain":              domainTitle,                     // e.g., "Labels" (as it comes from tags)
			"DomainSingular":      domainSingular,                  // e.g., "Label" (singular form)
			"DomainLower":         strings.ToLower(domain),         // e.g., "labels"
			"DomainSingularLower": strings.ToLower(domainSingular), // e.g., "label"
			"Operations":          processedOps,
			"Imports":             imports,
		}

		outputPath := filepath.Join(
			"internal/adapters/http/controllers",
			strings.ToLower(domain)+".gen.go",
		)

		tmpl, ok := state.GetTemplate("controller.tmpl")
		if !ok {
			return fmt.Errorf("controller template not found")
		}

		if err := state.FileWriter.WriteTemplate(outputPath, tmpl, data); err != nil {
			return fmt.Errorf("failed to generate handler for %s: %w", domain, err)
		}
	}

	return nil
}

// GenerateEvents generates domain events for entities and aggregates
func GenerateEvents(state *GlobalState) error {
	for name, processed := range state.ProcessedSchemas {
		// Only generate events for entities and aggregates, not value objects
		if processed.XCodegen != nil && processed.XCodegen.Type != "valueobject" {
			// Check if entity has an ID field
			hasIDField := false
			fields := processed.Fields

			for _, field := range fields {
				if field.FieldName == "ID" {
					hasIDField = true
					break
				}
			}

			// Only generate events for entities with ID fields
			if hasIDField && processed.Schema != nil {
				if err := generateEventsForSchema(state, processed.Schema, name); err != nil {
					return fmt.Errorf(
						"failed to generate events for %s: %w",
						name,
						err,
					)
				}
			}
		}
	}
	return nil
}

// generateEventsForSchema generates domain events for a schema
func generateEventsForSchema(state *GlobalState, schema *oas3.Schema, name string) error {
	title := name

	data := map[string]interface{}{
		"Package": "events",
		"Domain":  title,
		"Entities": []map[string]interface{}{
			{
				"Name":            title,
				"NameLower":       strings.ToLower(title),
				"NamePlural":      Pluralize(title),
				"NamePluralLower": strings.ToLower(Pluralize(title)),
			},
		},
	}

	outputPath := filepath.Join(
		"internal/core/events",
		strings.ToLower(title)+"_events.gen.go",
	)

	tmpl, ok := state.GetTemplate("events.tmpl")
	if !ok {
		return fmt.Errorf("events template not found")
	}

	return state.FileWriter.WriteTemplate(outputPath, tmpl, data)
}

// GenerateModels generates all model types (DTOs, entities, value objects)
func GenerateModels(state *GlobalState) error {
	for name, processed := range state.ProcessedSchemas {
		if processed.Schema == nil {
			continue
		}

		// Determine what type of model this is based on x-codegen
		modelType := "dto" // default
		if processed.XCodegen != nil {
			modelType = processed.XCodegen.Type
			if modelType == "" {
				modelType = "dto"
			}
			log.Info(
				"Generating model",
				slog.String("name", name),
				slog.String("type", modelType),
			)
		}

		if err := generateModel(state, processed.Schema, name, modelType); err != nil {
			return fmt.Errorf("failed to generate %s %s: %w", modelType, name, err)
		}
	}

	// Note: Inline schemas from operations are handled by the OpenAPI components
	// No need to generate them separately

	// Generate Params types for list operations
	if err := generateListParamTypes(state); err != nil {
		return fmt.Errorf("failed to generate list param types: %w", err)
	}

	return nil
}

// generateModel is the unified function for generating all model types (DTOs, entities, aggregates, value objects)
func generateModel(state *GlobalState, schema *oas3.Schema, name string, modelType string) error {
	// Ensure schema has a title set
	if schema.Title == nil || *schema.Title == "" {
		schema.Title = &name
	}

	// Determine package and output path based on model type
	var packageName, outputDir string
	var isEntity, isAggregate, isValueObject bool

	switch modelType {
	case schemaTypeEntity:
		packageName = "entities"
		outputDir = "internal/core/entities"
		isEntity = true
	case schemaTypeAggregate:
		packageName = "aggregates"
		outputDir = "internal/core/aggregates"
		isAggregate = true
	case schemaTypeValueObject:
		packageName = "valueobjects"
		outputDir = "internal/core/valueobjects"
		isValueObject = true

		// Handle batched value objects for referenced schemas
		if referencedSchemas := findReferencedSchemas(state, schema); len(referencedSchemas) > 0 {
			allSchemas := append([]*oas3.Schema{schema}, referencedSchemas...)
			return generateBatchedModels(
				state,
				allSchemas,
				strings.ToLower(name),
				packageName,
				outputDir,
			)
		}
	default: // "dto" or any other type
		packageName = "dto"
		outputDir = "internal/application/dto"

		// Handle special DTO subtypes based on naming convention
		if strings.HasSuffix(name, "Request") || strings.HasSuffix(name, "RequestBody") {
			packageName = "requests"
			outputDir = filepath.Join(outputDir, "requests")
		} else if strings.HasSuffix(name, "Response") {
			packageName = "responses"
			outputDir = filepath.Join(outputDir, "responses")
		} else if strings.HasSuffix(name, "Params") || strings.HasSuffix(name, "Query") {
			packageName = "params"
			outputDir = filepath.Join(outputDir, "params")
		}
	}

	// Extract fields
	fields := parsers.ExtractFields(schema)
	if len(fields) == 0 {
		return nil // Skip if no fields
	}

	// Sort fields alphabetically
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Name < fields[j].Name
	})

	// Prepare template data based on model type
	if modelType == "dto" {
		// DTO uses a simpler structure with Types array
		var templateFields []map[string]interface{}
		for _, field := range fields {
			jsonName := field.JSONTag
			if jsonName == "" {
				jsonName = field.Name
			}
			// Remove ,omitempty if present as template will add it
			if strings.Contains(jsonName, ",omitempty") {
				jsonName = strings.Split(jsonName, ",")[0]
			}

			templateFields = append(templateFields, map[string]interface{}{
				"FieldName":    field.FieldName,
				"GoType":       field.GoType,
				"JSONName":     jsonName,
				"YAMLName":     field.YAMLTag,
				"Required":     field.Required,
				"DefaultValue": field.DefaultValue,
			})
		}

		data := map[string]interface{}{
			"Package": packageName,
			"Types": []map[string]interface{}{
				{
					"Name":   name,
					"Fields": templateFields,
				},
			},
		}

		outputPath := filepath.Join(outputDir, strings.ToLower(name)+".gen.go")
		tmpl, ok := state.GetTemplate("schema.tmpl")
		if !ok {
			return fmt.Errorf("schema template not found")
		}
		return state.FileWriter.WriteTemplate(outputPath, tmpl, data)
	}

	// For domain types (Entity, Aggregate, ValueObject)
	var requiredFields, optionalFields []parsers.FieldDef
	hasIDField, hasTimeField, hasUUIDField := false, false, false

	for _, field := range fields {
		// Categorize fields
		if field.Required {
			requiredFields = append(requiredFields, field)
		} else {
			optionalFields = append(optionalFields, field)
		}

		// Check field types
		if field.FieldName == "ID" {
			hasIDField = true
		}
		if strings.Contains(field.GoType, "time.Time") {
			hasTimeField = true
		}
		if strings.Contains(field.GoType, "uuid.UUID") {
			hasUUIDField = true
		}
	}

	// Only entities/aggregates with ID fields should have domain events
	hasDomainEvents := hasIDField && !isValueObject

	desc := schema.GetDescription()
	data := map[string]interface{}{
		"Package":         packageName,
		"Name":            name,
		"NameLower":       strings.ToLower(name),
		"NamePlural":      Pluralize(name),
		"Description":     getDescription(&desc),
		"Type":            modelType,
		"Fields":          fields,
		"RequiredFields":  requiredFields,
		"OptionalFields":  optionalFields,
		"HasTimeFields":   hasTimeField,
		"HasUUIDFields":   hasUUIDField,
		"HasDomainEvents": hasDomainEvents,
		"IsEntity":        isEntity,
		"IsAggregate":     isAggregate,
		"IsValueObject":   isValueObject,
	}

	outputPath := filepath.Join(outputDir, strings.ToLower(name)+".gen.go")
	tmpl, ok := state.GetTemplate("schema.tmpl")
	if !ok {
		return fmt.Errorf("schema template not found")
	}

	return state.FileWriter.WriteTemplate(outputPath, tmpl, data)
}

// generateBatchedModels generates multiple models in a single file (used for value objects with references)
func generateBatchedModels(
	state *GlobalState,
	schemas []*oas3.Schema,
	outputName string,
	packageName string,
	outputDir string,
) error {
	allTypes := []map[string]interface{}{}
	hasTimeFields := false
	hasUUIDFields := false

	for _, schema := range schemas {
		fields := parsers.ExtractFields(schema)
		if len(fields) == 0 {
			continue
		}

		// Sort fields alphabetically
		sort.Slice(fields, func(i, j int) bool {
			return fields[i].Name < fields[j].Name
		})

		// Check for time and UUID fields
		for _, field := range fields {
			if strings.Contains(field.GoType, "time.Time") {
				hasTimeFields = true
			}
			if strings.Contains(field.GoType, "uuid.UUID") {
				hasUUIDFields = true
			}
		}

		// Separate required and optional fields
		var requiredFields, optionalFields []parsers.FieldDef
		for _, field := range fields {
			if field.Required {
				requiredFields = append(requiredFields, field)
			} else {
				optionalFields = append(optionalFields, field)
			}
		}

		title := schema.GetTitle()
		desc := schema.GetDescription()
		typeData := map[string]interface{}{
			"Name":           title,
			"NameLower":      strings.ToLower(title),
			"NamePlural":     Pluralize(title),
			"Description":    getDescription(&desc),
			"Fields":         fields,
			"RequiredFields": requiredFields,
			"OptionalFields": optionalFields,
		}

		allTypes = append(allTypes, typeData)
	}

	// Generate the consolidated file
	data := map[string]interface{}{
		"Package":       packageName,
		"Types":         allTypes,
		"HasTimeFields": hasTimeFields,
		"HasUUIDFields": hasUUIDFields,
	}

	outputPath := filepath.Join(outputDir, outputName+".gen.go")
	tmpl, ok := state.GetTemplate("schema.tmpl")
	if !ok {
		return fmt.Errorf("schema template not found")
	}

	return state.FileWriter.WriteTemplate(outputPath, tmpl, data)
}

// generateListParamTypes generates Params types for all operations that have parameters
func generateListParamTypes(state *GlobalState) error {

	// Map to track unique param types we need to generate
	paramTypesMap := make(map[string]bool)

	// Extract all operations from the OpenAPI spec
	operations := state.Operations

	log.Debug("Processing operations for params", slog.Int("count", len(operations)))

	for _, op := range operations {
		// Skip operations without parameters
		if len(op.Parameters) == 0 {
			continue
		}

		// Generate a params type name based on the operation ID
		if op.OperationID == "" {
			continue
		}

		// Convert operationId like "listAccounts" to "ListAccountsParams"
		paramTypeName := strings.ToUpper(op.OperationID[:1]) + op.OperationID[1:] + "Params"

		// Skip if we've already processed this param type
		if paramTypesMap[paramTypeName] {
			continue
		}
		paramTypesMap[paramTypeName] = true

		// Collect all parameters and their schemas
		var fields []map[string]interface{}

		for _, param := range op.Parameters {
			// Skip path parameters - they're part of the URL
			if param.In == "path" {
				continue
			}

			// Check if the parameter references a schema
			if param.Schema != "" {
				// This parameter references a schema (like Page)
				// Use the schema type directly instead of flattening
				fieldName := param.Name
				if fieldName == "" {
					fieldName = param.Schema
				}

				// Clean up description to be single-line
				cleanDesc := cleanDescription(param.Description)

				// Capitalize the field name properly
				capitalizedFieldName := strings.ToUpper(fieldName[:1]) + fieldName[1:]

				field := map[string]interface{}{
					"FieldName":   capitalizedFieldName,
					"GoType":      param.Schema,
					"JSONName":    param.Name,
					"YAMLName":    param.Name,
					"Required":    param.Required,
					"Description": cleanDesc,
				}

				// Make optional fields pointers
				if !param.Required {
					field["GoType"] = "*" + param.Schema
				}

				fields = append(fields, field)
			} else {
				// Simple parameter without schema reference
				// Clean up description to be single-line
				cleanDesc := cleanDescription(param.Description)
				field := map[string]interface{}{
					"FieldName":   strings.ToUpper(param.Name[:1]) + param.Name[1:],
					"GoType":      param.GoType,
					"JSONName":    param.Name,
					"YAMLName":    param.Name,
					"Required":    param.Required,
					"Description": cleanDesc,
				}

				// Handle optional fields
				if !param.Required && param.GoType != "" && !strings.HasPrefix(param.GoType, "*") {
					field["GoType"] = "*" + param.GoType
				}

				fields = append(fields, field)
			}
		}

		// Only create param type if there are non-path parameters
		if len(fields) > 0 {
			// Generate individual file for each param type
			if err := generateIndividualParamFile(state, paramTypeName, fields); err != nil {
				return fmt.Errorf("failed to generate param file for %s: %w", paramTypeName, err)
			}
		}
	}

	return nil
}

// generateIndividualParamFile generates a single param file for a param type
func generateIndividualParamFile(
	state *GlobalState,
	typeName string,
	fields []map[string]interface{},
) error {
	tmpl, ok := state.GetTemplate("schema.tmpl")
	if !ok {
		return fmt.Errorf("schema template not found")
	}

	data := map[string]interface{}{
		"Package": "dto",
		"Types": []map[string]interface{}{
			{
				"Name":   typeName,
				"Fields": fields,
			},
		},
	}

	// Create filename from type name (e.g., ListAccountsParams -> listaccountsparams.gen.go)
	outputPath := fmt.Sprintf("internal/application/dto/%s.gen.go", strings.ToLower(typeName))
	return state.FileWriter.WriteTemplate(outputPath, tmpl, data)
}

// extractSchemaNameFromRef extracts the schema name from a reference string
func extractSchemaNameFromRef(ref string) string {
	// Handle both local and remote references
	// Local: #/components/schemas/Config
	// Remote: ./ConfigAPI.yaml

	if strings.HasPrefix(ref, "#/") {
		// Local reference
		parts := strings.Split(ref, "/")
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
	} else if strings.HasSuffix(ref, ".yaml") {
		// Remote file reference
		base := filepath.Base(ref)
		return strings.TrimSuffix(base, ".yaml")
	}

	// Default: return the last part after any slash
	parts := strings.Split(ref, "/")
	return parts[len(parts)-1]
}

// findReferencedSchemas recursively finds all schemas referenced by a given schema's properties
func findReferencedSchemas(state *GlobalState, schema *oas3.Schema) []*oas3.Schema {
	// Track which schemas we've already processed to avoid cycles
	seen := make(map[string]bool)
	var result []*oas3.Schema

	// Helper function to recursively find references
	var findRefs func(*oas3.Schema)
	findRefs = func(s *oas3.Schema) {
		if s == nil || s.Properties == nil {
			return
		}

		// Check each property for references
		for propName := range s.Properties.Keys() {
			propRef := s.Properties.GetOrZero(propName)
			if propRef != nil && propRef.IsLeft() {
				prop := propRef.GetLeft()
				if prop != nil && prop.Ref != nil && prop.Ref.String() != "" {
					// Extract the schema name from the reference
					refString := prop.Ref.String()
					schemaName := extractSchemaNameFromRef(refString)

					// Skip if we've already seen this schema
					if seen[schemaName] {
						continue
					}
					seen[schemaName] = true

					// Find the referenced schema in our state
					if processed, exists := state.ProcessedSchemas[schemaName]; exists {
						// Only include if it doesn't have its own x-codegen
						// (schemas with x-codegen will be generated separately)
						if processed.XCodegen == nil && processed.Schema != nil {
							// Set the title on the schema if it's not already set
							if processed.Schema.Title == nil || *processed.Schema.Title == "" {
								title := schemaName
								processed.Schema.Title = &title
							}
							result = append(result, processed.Schema)

							// Recursively find references in this schema
							findRefs(processed.Schema)
						}
					}
				}
			}
		}
	}

	// Start the recursive search
	findRefs(schema)
	return result
}

// getDescription safely gets the description from a pointer
func getDescription(desc *string) string {
	if desc != nil {
		return *desc
	}
	return ""
}

// cleanDescription cleans up multi-line descriptions to be single-line for Go comments
func cleanDescription(desc string) string {
	if desc == "" {
		return ""
	}
	// Replace newlines with spaces and clean up
	desc = strings.ReplaceAll(desc, "\n", " ")
	// Remove list markers that would break Go syntax
	desc = strings.ReplaceAll(desc, "- ", "")
	// Clean up extra spaces
	desc = strings.TrimSpace(desc)
	// Replace multiple spaces with single space
	for strings.Contains(desc, "  ") {
		desc = strings.ReplaceAll(desc, "  ", " ")
	}
	return desc
}
