// Package labels provides label management functionality.
package labels

//go:generate go tool oapi-codegen --config=../../.codegen.types.yaml --package labels --include-tags Labels ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.codegen.server.yaml --package labels --include-tags Labels ../../api/openapi.bundled.yaml

// Domain constants
const (

	// MaxLabelsPerOrganization defines the maximum number of labels per organization
	MaxLabelsPerOrganization = 1000
)
