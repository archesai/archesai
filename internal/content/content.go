// Package content provides content management functionality including
// artifact storage and organization, and label management for content categorization.
package content

//go:generate go tool oapi-codegen --config=models.cfg.yaml ../../api/openapi.bundled.yaml
//go:generate go tool oapi-codegen --config=server.cfg.yaml ../../api/openapi.bundled.yaml

import "errors"

// Domain types
type (
	// Artifact represents a content artifact with its entity
	Artifact struct {
		ArtifactEntity
	}

	// Label represents a content label with its entity
	Label struct {
		LabelEntity
	}

	// CreateArtifactRequest represents a request to create an artifact
	CreateArtifactRequest = CreateArtifactJSONBody

	// UpdateArtifactRequest represents a request to update an artifact
	UpdateArtifactRequest = UpdateArtifactJSONBody

	// CreateLabelRequest represents a request to create a label
	CreateLabelRequest = CreateLabelJSONBody

	// UpdateLabelRequest represents a request to update a label
	UpdateLabelRequest = UpdateLabelJSONBody
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

// Domain constants
const (
	// MaxArtifactSize is the maximum size for an artifact in bytes (10MB)
	MaxArtifactSize = 10 * 1024 * 1024

	// MaxLabelsPerOrganization defines the maximum number of labels per organization
	MaxLabelsPerOrganization = 1000

	// MaxArtifactsPerOrganization defines the maximum number of artifacts per organization
	MaxArtifactsPerOrganization = 100000
)
