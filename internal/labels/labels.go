// Package labels provides label management functionality.
package labels

//go:generate go tool oapi-codegen --config=../../.codegen.types.yaml --package labels --include-tags Labels ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=../../.codegen.server.yaml --package labels --include-tags Labels ../../api/openapi.bundled.yaml

import "errors"

// Domain constants
const (

	// MaxLabelsPerOrganization defines the maximum number of labels per organization
	MaxLabelsPerOrganization = 1000
)

// Domain errors
var (
	// ErrArtifactNotFound is returned when an artifact is not found
	ErrArtifactNotFound = errors.New("artifact not found")

	// ErrArtifactTooLarge is returned when an artifact exceeds the maximum size
	ErrArtifactTooLarge = errors.New("artifact too large")

	// ErrLabelNotFound is returned when a label is not found
	ErrLabelNotFound = errors.New("label not found")

	// ErrLabelExists is returned when a label already exists
	ErrLabelExists = errors.New("label already exists")
)
