package main

type EntityConfig struct {
	Name         string
	Package      string
	ParentModule string
	Table        string
	Fields       []Field
	CreateFields []string
	UpdateFields []string
}

type Field struct {
	Name     string
	Type     string
	DBName   string
	JsonTag  string
	Validate string
	Required bool
	Nullable bool
}

type TemplateData struct {
	Package        string
	Entity         string
	LowerEntity    string
	EntityFields   []TemplateField
	CreateFields   []TemplateField
	UpdateFields   []TemplateField
	RequiredFields []TemplateField
	LogFields      []TemplateField
	HasEmailField  bool
	HasSlugField   bool
	HasSearchField bool
}

type TemplateField struct {
	Name      string
	Type      string
	JsonTag   string
	LogName   string
	Validate  string
	OmitEmpty bool
	Required  bool
}

func main() {
	// 	var (
	// 		specDir = flag.String("spec", "api/specifications", "OpenAPI specifications directory")
	// 		entity  = flag.String("entity", "", "Entity name to generate (optional, generates all if not specified)")
	// 		output  = flag.String("output", "", "Output directory (required if entity specified)")
	// 		genType = flag.String("type", "all", "Type to generate: entity, repository, service, handler, or all")
	// 		parent  = flag.String("parent", "", "Parent module (optional, e.g., 'auth' to generate under internal/auth/package)")
	// 	)
	// 	flag.Parse()

	// 	if *entity != "" && *output == "" {
	// 		log.Fatal("output directory is required when entity is specified")
	// 	}

	// 	// Create a merged document with all specifications
	// 	mergedSpec := mergeSpecifications(*specDir)

	// 	// Load with kin-openapi
	// 	loader := openapi3.NewLoader()
	// 	loader.IsExternalRefsAllowed = true
	// 	loader.ReadFromURIFunc = openapi3.ReadFromFile

	// 	// Save merged spec to temp file for loading
	// 	tempFile, err := os.CreateTemp("", "merged-spec-*.yaml")
	// 	if err != nil {
	// 		log.Fatalf("Failed to create temp file: %v", err)
	// 	}
	// 	defer os.Remove(tempFile.Name())

	// 	mergedYAML, err := yaml.Marshal(mergedSpec)
	// 	if err != nil {
	// 		log.Fatalf("Failed to marshal merged spec: %v", err)
	// 	}

	// 	if err := os.WriteFile(tempFile.Name(), mergedYAML, 0644); err != nil {
	// 		log.Fatalf("Failed to write temp file: %v", err)
	// 	}

	// 	// Load without validation to handle circular refs
	// 	docBytes, err := os.ReadFile(tempFile.Name())
	// 	if err != nil {
	// 		log.Fatalf("Failed to read temp file: %v", err)
	// 	}

	// 	doc, err := loader.LoadFromData(docBytes)
	// 	if err != nil {
	// 		log.Fatalf("Failed to load merged spec: %v", err)
	// 	}

	// 	entities := make(map[string]EntityConfig)

	// 	// Process each entity type
	// 	specFiles, err := filepath.Glob(filepath.Join(*specDir, "*.yaml"))
	// 	if err != nil {
	// 		log.Fatalf("Failed to find spec files: %v", err)
	// 	}

	// 	for _, specFile := range specFiles {
	// 		baseName := strings.TrimSuffix(filepath.Base(specFile), ".yaml")
	// 		if baseName == "openapi" || baseName == "config" || baseName == "errors" || baseName == "filters" || baseName == "health" {
	// 			continue // Skip non-entity files
	// 		}

	// 		// Convert plural to singular and capitalize
	// 		entityName := pluralToSingular(baseName)
	// 		entityName = cases.Title(language.AmericanEnglish).String(entityName)

	// 		// Find entity schema in merged doc
	// 		var entitySchema *openapi3.SchemaRef
	// 		if doc.Components != nil && doc.Components.Schemas != nil {
	// 			for schemaName, schema := range doc.Components.Schemas {
	// 				if strings.HasSuffix(schemaName, entityName+"Entity") {
	// 					entitySchema = schema
	// 					break
	// 				}
	// 			}
	// 		}

	// 		if entitySchema == nil {
	// 			log.Printf("No entity schema found for %s", entityName)
	// 			continue
	// 		}

	// 		// Extract fields from schema
	// 		fields := extractFields(entitySchema)

	// 		// Extract create/update fields from paths
	// 		createFields, updateFields := extractRequestFields(doc, baseName)

	// 		// Extract parent module from first tag
	// 		parentModule := extractParentModuleFromTags(doc, baseName)
	// 		if *parent != "" {
	// 			parentModule = *parent // Override with command line flag if provided
	// 		}

	// 		// Create entity config
	// 		config := EntityConfig{
	// 			Name:         entityName,
	// 			Package:      strings.ToLower(baseName),
	// 			ParentModule: parentModule,
	// 			Table:        baseName,
	// 			Fields:       fields,
	// 			CreateFields: createFields,
	// 			UpdateFields: updateFields,
	// 		}

	// 		entities[entityName] = config
	// 		log.Printf("Found entity %s with %d fields", entityName, len(fields))
	// 	}

	// 	// Generate for specific entity or all entities
	// 	if *entity != "" {
	// 		config, exists := entities[*entity]
	// 		if !exists {
	// 			log.Fatalf("Entity %s not found in specifications", *entity)
	// 		}

	// 		if err := generateEntity(config, *output, *genType); err != nil {
	// 			log.Fatalf("Failed to generate entity %s: %v", *entity, err)
	// 		}
	// 		fmt.Printf("Generated %s in %s\n", *entity, *output)
	// 	} else {
	// 		// Generate all entities
	// 		for _, config := range entities {
	// 			var outputDir string
	// 			if config.ParentModule != "" {
	// 				outputDir = filepath.Join("internal", config.ParentModule, config.Package)
	// 			} else {
	// 				outputDir = filepath.Join("internal", config.Package)
	// 			}
	// 			if err := generateEntity(config, outputDir, *genType); err != nil {
	// 				log.Printf("Failed to generate entity %s: %v", config.Name, err)
	// 				continue
	// 			}
	// 			fmt.Printf("Generated %s in %s\n", config.Name, outputDir)
	// 		}
	// 	}
	// }

	// func mergeSpecifications(specDir string) map[string]interface{} {
	// 	merged := map[string]interface{}{
	// 		"openapi": "3.1.1",
	// 		"info": map[string]interface{}{
	// 			"title":   "Merged API",
	// 			"version": "1.0.0",
	// 		},
	// 		"paths": map[string]interface{}{},
	// 		"components": map[string]interface{}{
	// 			"schemas":         map[string]interface{}{},
	// 			"securitySchemes": map[string]interface{}{},
	// 		},
	// 	}

	// 	// Load main OpenAPI file if exists
	// 	mainPath := filepath.Join("api", "openapi.yaml")
	// 	if data, err := os.ReadFile(mainPath); err == nil {
	// 		var mainSpec map[string]interface{}
	// 		if err := yaml.Unmarshal(data, &mainSpec); err == nil {
	// 			if components, ok := mainSpec["components"].(map[string]interface{}); ok {
	// 				if secSchemes, ok := components["securitySchemes"]; ok {
	// 					merged["components"].(map[string]interface{})["securitySchemes"] = secSchemes
	// 				}
	// 			}
	// 		}
	// 	}

	// 	// Get all spec files
	// 	specFiles, err := filepath.Glob(filepath.Join(specDir, "*.yaml"))
	// 	if err != nil {
	// 		log.Printf("Failed to find spec files: %v", err)
	// 		return merged
	// 	}

	// 	// Load and merge all specifications
	// 	for _, specFile := range specFiles {
	// 		data, err := os.ReadFile(specFile)
	// 		if err != nil {
	// 			log.Printf("Failed to read %s: %v", specFile, err)
	// 			continue
	// 		}

	// 		var spec map[string]interface{}
	// 		if err := yaml.Unmarshal(data, &spec); err != nil {
	// 			log.Printf("Failed to parse %s: %v", specFile, err)
	// 			continue
	// 		}

	// 		// Merge paths
	// 		if paths, ok := spec["paths"].(map[string]interface{}); ok {
	// 			for path, pathItem := range paths {
	// 				merged["paths"].(map[string]interface{})[path] = pathItem
	// 			}
	// 		}

	// 		// Merge components
	// 		if components, ok := spec["components"].(map[string]interface{}); ok {
	// 			if schemas, ok := components["schemas"].(map[string]interface{}); ok {
	// 				mergedSchemas := merged["components"].(map[string]interface{})["schemas"].(map[string]interface{})
	// 				for name, schema := range schemas {
	// 					mergedSchemas[name] = schema
	// 				}
	// 			}
	// 		}
	// 	}

	// return merged
}

// func extractFields(schema *openapi3.SchemaRef) []Field {
// 	var fields []Field

// 	if schema.Value == nil || schema.Value.Properties == nil {
// 		return fields
// 	}

// 	// Define field ordering to ensure consistent output
// 	fieldOrder := []string{"id", "createdAt", "updatedAt"}
// 	orderedProps := make([]string, 0, len(schema.Value.Properties))

// 	// Add ordered fields first
// 	for _, fieldName := range fieldOrder {
// 		if _, exists := schema.Value.Properties[fieldName]; exists {
// 			orderedProps = append(orderedProps, fieldName)
// 		}
// 	}

// 	// Add remaining fields in alphabetical order
// 	for propName := range schema.Value.Properties {
// 		if !contains(fieldOrder, propName) {
// 			orderedProps = append(orderedProps, propName)
// 		}
// 	}

// 	for _, propName := range orderedProps {
// 		propSchema := schema.Value.Properties[propName]
// 		field := Field{
// 			Name:     toPascalCase(propName),
// 			JsonTag:  propName,
// 			DBName:   camelToSnake(propName),
// 			Required: contains(schema.Value.Required, propName),
// 		}

// 		// Determine Go type
// 		field.Type = schemaToGoType(propSchema)

// 		// Check nullable
// 		if propSchema.Value != nil && propSchema.Value.Nullable {
// 			field.Nullable = true
// 			if !strings.HasPrefix(field.Type, "*") && field.Type != "interface{}" {
// 				field.Type = "*" + field.Type
// 			}
// 		}

// 		// Add validation based on schema
// 		var validations []string
// 		if field.Required {
// 			validations = append(validations, "required")
// 		}
// 		if propSchema.Value != nil {
// 			if propSchema.Value.MinLength > 0 {
// 				validations = append(validations, fmt.Sprintf("min=%d", propSchema.Value.MinLength))
// 			}
// 			if propSchema.Value.MaxLength != nil && *propSchema.Value.MaxLength > 0 {
// 				validations = append(validations, fmt.Sprintf("max=%d", *propSchema.Value.MaxLength))
// 			}
// 			if propSchema.Value.Format == "email" {
// 				validations = append(validations, "email")
// 			}
// 			if propSchema.Value.Format == "uuid" {
// 				validations = append(validations, "uuid")
// 			}
// 		}
// 		field.Validate = strings.Join(validations, ",")

// 		fields = append(fields, field)
// 	}

// 	return fields
// }

// func extractParentModuleFromTags(doc *openapi3.T, baseName string) string {
// 	if doc.Paths == nil {
// 		return ""
// 	}

// 	// Look for the first operation with this baseName and extract the first tag
// 	for path, pathItem := range doc.Paths.Map() {
// 		if !strings.HasPrefix(path, "/"+baseName) {
// 			continue
// 		}

// 		// Check all HTTP methods for tags
// 		operations := []*openapi3.Operation{
// 			pathItem.Get,
// 			pathItem.Post,
// 			pathItem.Put,
// 			pathItem.Patch,
// 			pathItem.Delete,
// 		}

// 		for _, op := range operations {
// 			if op != nil && len(op.Tags) > 0 {
// 				firstTag := op.Tags[0]
// 				// Convert tag to lowercase module name
// 				// e.g., "Authentication" -> "auth", "Members" -> "members"
// 				return strings.ToLower(tagToModuleName(firstTag))
// 			}
// 		}
// 	}

// 	return ""
// }

// func tagToModuleName(tag string) string {
// 	// Convert common tag patterns to module names
// 	tag = strings.ToLower(tag)
// 	switch tag {
// 	case "authentication":
// 		return "auth"
// 	case "email verification", "email change", "password reset":
// 		return "auth"
// 	case "members":
// 		return "members"
// 	case "organizations":
// 		return "organizations"
// 	case "users":
// 		return "users"
// 	default:
// 		// Remove spaces and common suffixes
// 		tag = strings.ReplaceAll(tag, " ", "")
// 		tag = strings.TrimSuffix(tag, "s")
// 		return tag
// 	}
// }

// func extractRequestFields(doc *openapi3.T, baseName string) ([]string, []string) {
// 	var createFields, updateFields []string

// 	if doc.Paths == nil {
// 		return createFields, updateFields
// 	}

// 	for path, pathItem := range doc.Paths.Map() {
// 		// Look for POST endpoint (create)
// 		if pathItem.Post != nil && strings.HasPrefix(path, "/"+baseName) && !strings.Contains(path, "{") {
// 			if pathItem.Post.RequestBody != nil && pathItem.Post.RequestBody.Value != nil && pathItem.Post.RequestBody.Value.Content != nil {
// 				if jsonContent, ok := pathItem.Post.RequestBody.Value.Content["application/json"]; ok && jsonContent.Schema != nil {
// 					createFields = extractFieldNames(jsonContent.Schema)
// 				}
// 			}
// 		}

// 		// Look for PATCH endpoint (update)
// 		if pathItem.Patch != nil && strings.Contains(path, "{id}") {
// 			if pathItem.Patch.RequestBody != nil && pathItem.Patch.RequestBody.Value != nil && pathItem.Patch.RequestBody.Value.Content != nil {
// 				if jsonContent, ok := pathItem.Patch.RequestBody.Value.Content["application/json"]; ok && jsonContent.Schema != nil {
// 					updateFields = extractFieldNames(jsonContent.Schema)
// 				}
// 			}
// 		}
// 	}

// 	return createFields, updateFields
// }

// func extractFieldNames(schema *openapi3.SchemaRef) []string {
// 	var names []string
// 	if schema.Value != nil && schema.Value.Properties != nil {
// 		for propName := range schema.Value.Properties {
// 			names = append(names, toPascalCase(propName))
// 		}
// 	}
// 	return names
// }

// func schemaToGoType(schema *openapi3.SchemaRef) string {
// 	if schema.Value == nil {
// 		return "interface{}"
// 	}

// 	// Handle nullable types first
// 	if schema.Value.Nullable {
// 		baseType := schemaToGoTypeNonNullable(schema)
// 		if !strings.HasPrefix(baseType, "*") && baseType != "interface{}" {
// 			return "*" + baseType
// 		}
// 		return baseType
// 	}

// 	return schemaToGoTypeNonNullable(schema)
// }

// func schemaToGoTypeNonNullable(schema *openapi3.SchemaRef) string {
// 	if schema.Value == nil {
// 		return "interface{}"
// 	}

// 	// Check if type is set
// 	if schema.Value.Type == nil || len(*schema.Value.Type) == 0 {
// 		// Try to infer from format or other properties
// 		if schema.Value.Format != "" {
// 			switch schema.Value.Format {
// 			case "date-time":
// 				return "time.Time"
// 			case "uuid":
// 				return "string"
// 			case "int64":
// 				return "int64"
// 			case "int32":
// 				return "int32"
// 			}
// 		}
// 		return "interface{}"
// 	}

// 	// Get first type if multiple types
// 	typeStr := (*schema.Value.Type)[0]

// 	switch typeStr {
// 	case "string":
// 		switch schema.Value.Format {
// 		case "date-time":
// 			return "time.Time"
// 		case "uuid":
// 			return "string"
// 		case "date":
// 			return "string" // Could be time.Time if date-only is needed
// 		case "email":
// 			return "string"
// 		case "uri", "url":
// 			return "string"
// 		default:
// 			return "string"
// 		}
// 	case "integer":
// 		switch schema.Value.Format {
// 		case "int64":
// 			return "int64"
// 		case "int32":
// 			return "int32"
// 		default:
// 			return "int"
// 		}
// 	case "number":
// 		switch schema.Value.Format {
// 		case "float":
// 			return "float32"
// 		case "double":
// 			return "float64"
// 		default:
// 			return "float64"
// 		}
// 	case "boolean":
// 		return "bool"
// 	case "array":
// 		if schema.Value.Items != nil {
// 			itemType := schemaToGoType(schema.Value.Items)
// 			return "[]" + itemType
// 		}
// 		return "[]interface{}"
// 	case "object":
// 		// Check if it has properties defined
// 		if len(schema.Value.Properties) > 0 {
// 			// This would be a struct, but we'll keep it simple for now
// 			return "map[string]interface{}"
// 		}
// 		// Check for additionalProperties
// 		if schema.Value.AdditionalProperties.Has != nil && *schema.Value.AdditionalProperties.Has {
// 			if schema.Value.AdditionalProperties.Schema != nil {
// 				valueType := schemaToGoType(schema.Value.AdditionalProperties.Schema)
// 				return "map[string]" + valueType
// 			}
// 		}
// 		return "map[string]interface{}"
// 	default:
// 		return "interface{}"
// 	}
// }

// func generateEntity(config EntityConfig, outputDir, genType string) error {
// 	// Create output directory
// 	if err := os.MkdirAll(outputDir, 0755); err != nil {
// 		return fmt.Errorf("failed to create output directory: %w", err)
// 	}

// 	// Prepare template data
// 	templateData := prepareTemplateData(config)

// 	// Generate files based on type
// 	if genType == "all" || genType == "entity" {
// 		if err := generateFromTemplate("entity.go.tmpl", filepath.Join(outputDir, "entity.gen.go"), templateData); err != nil {
// 			return fmt.Errorf("failed to generate entity: %w", err)
// 		}
// 	}

// 	if genType == "all" || genType == "repository" {
// 		if err := generateFromTemplate("repository.go.tmpl", filepath.Join(outputDir, "repository.gen.go"), templateData); err != nil {
// 			return fmt.Errorf("failed to generate repository: %w", err)
// 		}
// 	}

// 	if genType == "all" || genType == "service" {
// 		if err := generateFromTemplate("service.go.tmpl", filepath.Join(outputDir, "service.gen.go"), templateData); err != nil {
// 			return fmt.Errorf("failed to generate service: %w", err)
// 		}
// 	}

// 	if genType == "all" || genType == "handler" {
// 		if err := generateFromTemplate("handler.go.tmpl", filepath.Join(outputDir, "handler.gen.go"), templateData); err != nil {
// 			return fmt.Errorf("failed to generate handler: %w", err)
// 		}
// 	}

// 	return nil
// }

// func prepareTemplateData(config EntityConfig) TemplateData {
// 	data := TemplateData{
// 		Package:     config.Package,
// 		Entity:      config.Name,
// 		LowerEntity: strings.ToLower(config.Name),
// 	}

// 	// Convert fields
// 	fieldMap := make(map[string]Field)
// 	for _, field := range config.Fields {
// 		fieldMap[field.Name] = field

// 		templateField := TemplateField{
// 			Name:      field.Name,
// 			Type:      field.Type,
// 			JsonTag:   field.JsonTag,
// 			LogName:   strings.ToLower(field.Name),
// 			Validate:  field.Validate,
// 			OmitEmpty: field.Nullable || strings.HasPrefix(field.Type, "*"),
// 			Required:  field.Required,
// 		}

// 		data.EntityFields = append(data.EntityFields, templateField)

// 		// Check for special fields
// 		if strings.ToLower(field.Name) == "email" {
// 			data.HasEmailField = true
// 		}
// 		if strings.ToLower(field.Name) == "slug" {
// 			data.HasSlugField = true
// 		}
// 		if strings.Contains(strings.ToLower(field.Name), "name") || strings.Contains(strings.ToLower(field.Name), "title") {
// 			data.HasSearchField = true
// 		}
// 	}

// 	// Create fields
// 	for _, fieldName := range config.CreateFields {
// 		if field, exists := fieldMap[fieldName]; exists {
// 			templateField := TemplateField{
// 				Name:     field.Name,
// 				Type:     field.Type,
// 				JsonTag:  field.JsonTag,
// 				Validate: field.Validate,
// 				Required: field.Required,
// 			}
// 			data.CreateFields = append(data.CreateFields, templateField)

// 			if field.Required {
// 				data.RequiredFields = append(data.RequiredFields, templateField)
// 			}

// 			// Add to log fields if it's a string field
// 			if field.Type == "string" {
// 				data.LogFields = append(data.LogFields, templateField)
// 			}
// 		}
// 	}

// 	// Update fields (make them pointers for optional updates)
// 	for _, fieldName := range config.UpdateFields {
// 		if field, exists := fieldMap[fieldName]; exists {
// 			fieldType := field.Type
// 			if !strings.HasPrefix(fieldType, "*") && fieldType != "interface{}" {
// 				fieldType = "*" + fieldType
// 			}

// 			templateField := TemplateField{
// 				Name:      field.Name,
// 				Type:      fieldType,
// 				JsonTag:   field.JsonTag,
// 				Validate:  "omitempty," + field.Validate,
// 				OmitEmpty: true,
// 			}
// 			data.UpdateFields = append(data.UpdateFields, templateField)
// 		}
// 	}

// 	return data
// }

// func generateFromTemplate(templateFile, outputFile string, data TemplateData) error {
// 	tmplPath := filepath.Join("internal/generator/templates", templateFile)
// 	tmpl, err := template.ParseFiles(tmplPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to parse template: %w", err)
// 	}

// 	var buf bytes.Buffer
// 	if err := tmpl.Execute(&buf, data); err != nil {
// 		return fmt.Errorf("failed to execute template: %w", err)
// 	}

// 	// Format the generated Go code
// 	formatted, err := format.Source(buf.Bytes())
// 	if err != nil {
// 		// If formatting fails, write unformatted code for debugging
// 		if err := os.WriteFile(outputFile, buf.Bytes(), 0644); err != nil {
// 			return fmt.Errorf("failed to write file: %w", err)
// 		}
// 		return fmt.Errorf("generated code has formatting errors: %w", err)
// 	}

// 	// Write the formatted code
// 	if err := os.WriteFile(outputFile, formatted, 0644); err != nil {
// 		return fmt.Errorf("failed to write file: %w", err)
// 	}

// 	return nil
// }

// // Helper functions
// func contains(s []string, item string) bool {
// 	return slices.Contains(s, item)
// }

// func toPascalCase(s string) string {
// 	if s == "" {
// 		return s
// 	}

// 	// Handle common camelCase to PascalCase conversion
// 	runes := []rune(s)
// 	runes[0] = rune(strings.ToUpper(string(runes[0]))[0])
// 	return string(runes)
// }

// func camelToSnake(s string) string {
// 	var result strings.Builder
// 	for i, r := range s {
// 		if i > 0 && r >= 'A' && r <= 'Z' {
// 			result.WriteRune('_')
// 		}
// 		result.WriteRune(r)
// 	}
// 	return strings.ToLower(result.String())
// }

// func pluralToSingular(s string) string {
// 	// Simple pluralization rules
// 	if strings.HasSuffix(s, "ies") {
// 		return s[:len(s)-3] + "y"
// 	}
// 	if strings.HasSuffix(s, "ses") || strings.HasSuffix(s, "xes") || strings.HasSuffix(s, "zes") {
// 		return s[:len(s)-2]
// 	}
// 	if strings.HasSuffix(s, "s") && !strings.HasSuffix(s, "ss") {
// 		return s[:len(s)-1]
// 	}
// 	return s
// }
