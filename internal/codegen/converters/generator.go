// Package converters generates type converters between database and API models.
package converters

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

//go:embed templates/adapters.go.tmpl
var adaptersTemplate string

// Generator handles generation of type converters.
type Generator struct{}

// NewGenerator creates a new converters generator.
func NewGenerator() *Generator {
	return &Generator{}
}

// ConverterSpec represents a single converter function.
type ConverterSpec struct {
	Name     string
	From     string
	FromType string
	To       string
	ToType   string
	ToPrefix string
	ToVar    string
	Fields   []FieldMapping
}

// FieldMapping represents a field conversion.
type FieldMapping struct {
	ToField    string
	Conversion string
}

// Config represents the converter configuration.
type Config struct {
	Converters []ConverterConfig `yaml:"converters"`
}

// ConverterConfig represents a single converter configuration.
type ConverterConfig struct {
	Name      string                 `yaml:"name"`
	From      string                 `yaml:"from"`
	To        string                 `yaml:"to"`
	Automap   bool                   `yaml:"automap"`
	Mappings  map[string]string      `yaml:"mappings"`
	Overrides map[string]interface{} `yaml:"overrides"`
}

// TemplateData represents data passed to the template.
type TemplateData struct {
	Domain                 string
	Converters             []ConverterSpec
	Imports                []string
	NeedsHelpers           bool
	NeedsNullableString    bool
	NeedsNullableTime      bool
	NeedsNullableTimestamp bool
	NeedsNullableJSON      bool
	NeedsNullableMetadata  bool
}

// Generate generates converter code for all domains.
func (g *Generator) Generate() error {
	// Read converter configuration
	configPath := "internal/domains/adapters.yaml"
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// Generate converters for each domain
	domains := []string{"auth", "organizations", "workflows", "content"}
	for _, domain := range domains {
		if err := g.generateDomainConverters(domain, config); err != nil {
			log.Printf("Failed to generate converters for %s: %v", domain, err)
		}
	}

	return nil
}

func (g *Generator) generateDomainConverters(domain string, config Config) error {
	// Filter converters for this domain
	var domainConverters []ConverterSpec
	for _, conv := range config.Converters {
		if isDomainConverter(conv.Name, domain) {
			spec := g.buildConverterSpec(conv, domain)
			domainConverters = append(domainConverters, spec)
		}
	}

	if len(domainConverters) == 0 {
		return nil
	}

	// Analyze all conversions to determine what we need
	templateData := g.analyzeConverters(domain, domainConverters)

	// Parse and execute template
	tmpl, err := template.New("adapters").Parse(adaptersTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, templateData); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Format the generated code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to format code: %w", err)
	}

	// Write to file
	outputPath := fmt.Sprintf("internal/domains/%s/adapters/adapters.gen.go", domain)
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(outputPath, formatted, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	log.Printf("Generated converters for %s domain", domain)
	return nil
}

// isDomainConverter checks if a converter belongs to a specific domain based on naming
func isDomainConverter(name, domain string) bool {
	nameLower := strings.ToLower(name)

	// Check for domain-specific prefixes or entity names
	switch domain {
	case "auth":
		return strings.Contains(nameLower, "auth") ||
			strings.Contains(nameLower, "user") ||
			strings.Contains(nameLower, "session") ||
			strings.Contains(nameLower, "account")
	case "organizations":
		return strings.Contains(nameLower, "organization") ||
			strings.Contains(nameLower, "member") ||
			strings.Contains(nameLower, "invitation")
	case "workflows":
		return strings.Contains(nameLower, "pipeline") ||
			strings.Contains(nameLower, "run") ||
			strings.Contains(nameLower, "tool")
	case "content":
		return strings.Contains(nameLower, "artifact") ||
			strings.Contains(nameLower, "label")
	}
	return false
}

func (g *Generator) buildConverterSpec(config ConverterConfig, domain string) ConverterSpec {
	spec := ConverterSpec{
		Name:     config.Name,
		From:     config.From,
		FromType: resolveType(config.From),
		To:       config.To,
		ToType:   resolveType(config.To),
		ToPrefix: getTypePrefix(config.To),
		ToVar:    getTypeName(config.To),
		Fields:   []FieldMapping{},
	}

	// Handle automap if enabled
	if config.Automap {
		// Get source and target struct fields
		fromStructName := getTypeName(config.From)
		toStructName := getTypeName(config.To)

		// Determine file paths
		var fromPath, toPath string
		if strings.Contains(config.From, "postgresql") {
			fromPath = "internal/generated/database/postgresql/models.go"
		}
		if strings.Contains(config.To, "api") {
			toPath = fmt.Sprintf("internal/domains/%s/generated/api/types.gen.go", domain)
		}

		// Get fields from both structs
		var fromFields, toFields []string
		if fromPath != "" {
			fromFields, _ = getStructFields(fromPath, fromStructName)
		}
		if toPath != "" {
			toFields, _ = getStructFields(toPath, toStructName)
		}

		// Create a set of target fields
		toFieldSet := make(map[string]bool)
		for _, f := range toFields {
			toFieldSet[f] = true
		}

		// Create a set of overridden fields (to skip during automap)
		overrideSet := make(map[string]bool)
		for field := range config.Overrides {
			overrideSet[field] = true
		}

		// Auto-map matching fields (except overridden ones)
		for _, field := range fromFields {
			if toFieldSet[field] && !overrideSet[field] {
				// Determine smart conversion based on field name and types
				conversion := determineAutoConversion(field, config.From, config.To)
				spec.Fields = append(spec.Fields, FieldMapping{
					ToField:    field,
					Conversion: conversion,
				})
			}
		}
	}

	// Add explicit mappings
	for fromField, toField := range config.Mappings {
		// Remove any existing mapping for this field
		newFields := []FieldMapping{}
		for _, fm := range spec.Fields {
			if fm.ToField != toField {
				newFields = append(newFields, fm)
			}
		}
		spec.Fields = newFields

		mapping := FieldMapping{
			ToField:    toField,
			Conversion: fmt.Sprintf("from.%s", fromField),
		}
		spec.Fields = append(spec.Fields, mapping)
	}

	// Add overrides (these take precedence)
	for field, value := range config.Overrides {
		// Remove any existing mapping for this field
		newFields := []FieldMapping{}
		for _, fm := range spec.Fields {
			if fm.ToField != field {
				newFields = append(newFields, fm)
			}
		}
		spec.Fields = newFields

		// Add the override
		spec.Fields = append(spec.Fields, FieldMapping{
			ToField:    field,
			Conversion: fmt.Sprintf("%v", value),
		})
	}

	// Sort fields alphabetically for consistent output
	sort.Slice(spec.Fields, func(i, j int) bool {
		return spec.Fields[i].ToField < spec.Fields[j].ToField
	})

	return spec
}

// determineAutoConversion determines the conversion for auto-mapped fields
func determineAutoConversion(fieldName, fromType, toType string) string {
	// Check if it's a database to API conversion
	isDBToAPI := strings.Contains(fromType, "postgresql") && strings.Contains(toType, "api")

	if isDBToAPI {
		// Handle common nullable string fields
		nullableStringFields := []string{
			"Image", "Logo", "StripeCustomerId", "AccessToken", "RefreshToken",
			"IdToken", "Scope", "ActiveOrganizationId", "IpAddress", "UserAgent",
			"Description", "Color", "StorageProvider", "StorageKey", "ContentType",
			"Checksum", "Error", "UpdatedBy", "Name", "BillingEmail", "PipelineId",
			"ProducerId", "PreviewImage", "Text",
		}
		for _, f := range nullableStringFields {
			if fieldName == f {
				return fmt.Sprintf("handleNullableString(from.%s)", fieldName)
			}
		}

		// Handle timestamp conversions
		if strings.Contains(fieldName, "At") && fieldName != "CreatedAt" && fieldName != "UpdatedAt" {
			if fieldName == "StartedAt" || fieldName == "CompletedAt" {
				return fmt.Sprintf("handleNullableTimestamp(from.%s)", fieldName)
			}
			// ExpiresAt for invitations needs time formatting
			if fieldName == "ExpiresAt" && strings.Contains(toType, "InvitationEntity") {
				return fmt.Sprintf("from.%s.Format(time.RFC3339)", fieldName)
			}
		}
	}

	// Default: direct mapping
	return fmt.Sprintf("from.%s", fieldName)
}

// getStructFields gets field names from a struct type
func getStructFields(filePath, structName string) ([]string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return nil, err
	}

	var fields []string
	ast.Inspect(file, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok || typeSpec.Name.Name != structName {
			return true
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		for _, field := range structType.Fields.List {
			if len(field.Names) > 0 {
				fields = append(fields, field.Names[0].Name)
			}
		}
		return false
	})
	return fields, nil
}

func resolveType(typeName string) string {
	// Keep the type as-is for proper imports
	return typeName
}

func getTypeName(fullType string) string {
	parts := strings.Split(fullType, ".")
	return parts[len(parts)-1]
}

func getTypePrefix(fullType string) string {
	if strings.Contains(fullType, "api.") {
		return "api."
	} else if strings.Contains(fullType, "postgresql.") {
		return "postgresql."
	}
	return ""
}

func (g *Generator) analyzeConverters(domain string, converters []ConverterSpec) TemplateData {
	// Collect all conversion strings
	var allConversions []string
	for _, conv := range converters {
		for _, field := range conv.Fields {
			allConversions = append(allConversions, field.Conversion)
		}
	}
	conversionsStr := strings.Join(allConversions, " ")

	// Determine required imports
	imports := []string{
		fmt.Sprintf(`"github.com/archesai/archesai/internal/domains/%s/generated/api"`, domain),
		`"github.com/archesai/archesai/internal/infrastructure/database/generated/postgresql"`,
	}

	// Add conditional imports
	if containsPattern(conversionsStr, `uuid\.MustParse|openapi_types\.UUID`) {
		imports = append(imports, `"github.com/google/uuid"`)
	}

	needsOpenapiTypes := containsPattern(conversionsStr, `openapi_types\.`)
	if needsOpenapiTypes {
		imports = append(imports, `openapi_types "github.com/oapi-codegen/runtime/types"`)
	}

	needsTime := containsPattern(conversionsStr, `time\.`) || containsPattern(conversionsStr, `handleNullableTime|Format\(time\.RFC3339\)`)
	if needsTime {
		imports = append(imports, `"time"`)
	}

	needsJSON := containsPattern(conversionsStr, `handleNullableJSON|json\.RawMessage|handleNullableMetadata`)
	if needsJSON {
		imports = append(imports, `"encoding/json"`)
	}

	needsPgtype := containsPattern(conversionsStr, `handleNullableTimestamp|pgtype\.`)
	if needsPgtype {
		imports = append(imports, `"github.com/jackc/pgx/v5/pgtype"`)
	}

	// Check what helpers we need
	needsNullableString := containsPattern(conversionsStr, `handleNullableString`)
	needsNullableTime := containsPattern(conversionsStr, `handleNullableTime`)
	needsNullableTimestamp := containsPattern(conversionsStr, `handleNullableTimestamp`)
	needsNullableJSON := containsPattern(conversionsStr, `handleNullableJSON`)
	needsNullableMetadata := containsPattern(conversionsStr, `handleNullableMetadata`)

	return TemplateData{
		Domain:                 domain,
		Converters:             converters,
		Imports:                imports,
		NeedsHelpers:           needsNullableString || needsNullableTime || needsNullableTimestamp || needsNullableJSON || needsNullableMetadata,
		NeedsNullableString:    needsNullableString,
		NeedsNullableTime:      needsNullableTime,
		NeedsNullableTimestamp: needsNullableTimestamp,
		NeedsNullableJSON:      needsNullableJSON,
		NeedsNullableMetadata:  needsNullableMetadata,
	}
}

func containsPattern(text, pattern string) bool {
	matched, _ := regexp.MatchString(pattern, text)
	return matched
}
