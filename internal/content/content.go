// Package content provides content management functionality including
// artifact storage and organization, and label management for content categorization.
package content

// Domain constants
const (
	// MaxArtifactSize is the maximum size for an artifact in bytes (10MB)
	MaxArtifactSize = 10 * 1024 * 1024

	// MaxLabelsPerOrganization defines the maximum number of labels per organization
	MaxLabelsPerOrganization = 1000

	// MaxArtifactsPerOrganization defines the maximum number of artifacts per organization
	MaxArtifactsPerOrganization = 100000
)
