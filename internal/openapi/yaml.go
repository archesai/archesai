package openapi

import (
	"bytes"

	"go.yaml.in/yaml/v4"
)

// openapiKeyOrder defines the standard order for OpenAPI document keys.
// This ensures consistent output across all bundled files.
// Keys are assigned priorities that work across all OpenAPI contexts.
var openapiKeyOrder = map[string]int{
	// JSON Schema / Reference keys (highest priority - always first)
	"$ref":           -100,
	"$schema":        -99,
	"$id":            -98,
	"$anchor":        -97,
	"$dynamicRef":    -96,
	"$dynamicAnchor": -95,
	"$vocabulary":    -94,
	"$comment":       -93,
	"$defs":          -92,

	// OpenAPI root level keys
	"openapi":           0,
	"x-project-name":    1,
	"info":              2,
	"jsonSchemaDialect": 3,

	// Identification keys (work in multiple contexts)
	"operationId": 5,
	"name":        6,
	"title":       7,
	"summary":     8,
	"description": 9,
	"version":     10,

	// Info-specific keys
	"termsOfService": 11,
	"contact":        12,
	"license":        13,

	// Type/Schema definition keys
	"type":     20,
	"enum":     21,
	"const":    22,
	"default":  23,
	"format":   24,
	"nullable": 25,

	// Composition keywords
	"allOf":            30,
	"anyOf":            31,
	"oneOf":            32,
	"not":              33,
	"if":               34,
	"then":             35,
	"else":             36,
	"dependentSchemas": 37,

	// Array keywords
	"prefixItems": 40,
	"items":       41,
	"contains":    42,

	// Object keywords
	"properties":            45,
	"patternProperties":     46,
	"additionalProperties":  47,
	"propertyNames":         48,
	"unevaluatedItems":      49,
	"unevaluatedProperties": 50,

	// Validation keywords - strings
	"minLength": 55,
	"maxLength": 56,
	"pattern":   57,

	// Validation keywords - numbers
	"minimum":          60,
	"maximum":          61,
	"exclusiveMinimum": 62,
	"exclusiveMaximum": 63,
	"multipleOf":       64,

	// Validation keywords - arrays
	"minItems":    70,
	"maxItems":    71,
	"uniqueItems": 72,
	"minContains": 73,
	"maxContains": 74,

	// Validation keywords - objects
	"minProperties":     80,
	"maxProperties":     81,
	"required":          82,
	"dependentRequired": 83,

	// Content keywords
	"contentEncoding":  90,
	"contentMediaType": 91,
	"contentSchema":    92,

	// OpenAPI structural keys
	"servers":      100,
	"security":     101,
	"tags":         102,
	"externalDocs": 103,
	"paths":        104,
	"webhooks":     105,
	"components":   106,

	// Components section keys
	"schemas":         110,
	"responses":       111,
	"parameters":      112,
	"examples":        113,
	"requestBodies":   114,
	"headers":         115,
	"securitySchemes": 116,
	"links":           117,
	"callbacks":       118,
	"pathItems":       119,

	// Operation keys
	"deprecated":  130,
	"requestBody": 131,

	// HTTP methods (for path items)
	"get":     140,
	"put":     141,
	"post":    142,
	"delete":  143,
	"options": 144,
	"head":    145,
	"patch":   146,
	"trace":   147,

	// Response/content keys
	"content": 150,
	"schema":  151,

	// Example/documentation keys
	"example":       160,
	"readOnly":      161,
	"writeOnly":     162,
	"xml":           163,
	"discriminator": 164,

	// Parameter/header keys
	"in":              170,
	"style":           171,
	"explode":         172,
	"allowReserved":   173,
	"allowEmptyValue": 174,

	// Security keys
	"scheme":       180,
	"bearerFormat": 181,
	"flows":        182,
	"openIdConnectUrl": 183,
}

// marshalYAML marshals a yaml.Node with 2-space indentation and consistent key ordering.
func marshalYAML(node *yaml.Node) ([]byte, error) {
	// Sort keys for consistent output
	sortYAMLNode(node)

	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(node); err != nil {
		return nil, err
	}
	if err := encoder.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// sortYAMLNode recursively sorts mapping node keys according to OpenAPI conventions.
func sortYAMLNode(node *yaml.Node) {
	if node == nil {
		return
	}

	switch node.Kind {
	case yaml.DocumentNode:
		for _, child := range node.Content {
			sortYAMLNode(child)
		}
	case yaml.MappingNode:
		sortMappingNode(node)
		// Recursively sort children
		for i := 1; i < len(node.Content); i += 2 {
			sortYAMLNode(node.Content[i])
		}
	case yaml.SequenceNode:
		for _, child := range node.Content {
			sortYAMLNode(child)
		}
	}
}

// sortMappingNode sorts the key-value pairs in a mapping node.
func sortMappingNode(node *yaml.Node) {
	if node.Kind != yaml.MappingNode || len(node.Content) < 4 {
		return // Need at least 2 pairs to sort
	}

	// Extract pairs
	pairs := make([]struct {
		key   *yaml.Node
		value *yaml.Node
	}, len(node.Content)/2)

	for i := 0; i < len(node.Content); i += 2 {
		pairs[i/2].key = node.Content[i]
		pairs[i/2].value = node.Content[i+1]
	}

	// Sort pairs using stable sort to maintain relative order of keys not in openapiKeyOrder
	stableSort(pairs, func(a, b struct {
		key   *yaml.Node
		value *yaml.Node
	}) bool {
		aKey := a.key.Value
		bKey := b.key.Value

		aOrder, aHasOrder := openapiKeyOrder[aKey]
		bOrder, bHasOrder := openapiKeyOrder[bKey]

		// Both have defined order - use it
		if aHasOrder && bHasOrder {
			return aOrder < bOrder
		}

		// Only one has defined order - it comes first
		if aHasOrder {
			return true
		}
		if bHasOrder {
			return false
		}

		// Neither has defined order - x- prefixed keys go last, then alphabetical
		aIsExtension := len(aKey) > 2 && aKey[:2] == "x-"
		bIsExtension := len(bKey) > 2 && bKey[:2] == "x-"

		if aIsExtension && !bIsExtension {
			return false
		}
		if !aIsExtension && bIsExtension {
			return true
		}

		// Both are extensions or neither - alphabetical
		return aKey < bKey
	})

	// Rebuild content
	for i, pair := range pairs {
		node.Content[i*2] = pair.key
		node.Content[i*2+1] = pair.value
	}
}

// stableSort performs a stable insertion sort on the pairs slice.
func stableSort[T any](items []T, less func(a, b T) bool) {
	for i := 1; i < len(items); i++ {
		for j := i; j > 0 && less(items[j], items[j-1]); j-- {
			items[j], items[j-1] = items[j-1], items[j]
		}
	}
}
