// Package artifacts provides artifact management functionality including
// artifact storage and organization.
package artifacts

// Domain constants
const (
	// MaxArtifactSize is the maximum size for an artifact in bytes (10MB)
	MaxArtifactSize = 10 * 1024 * 1024

	// MaxLabelsPerOrganization defines the maximum number of labels per organization
	MaxLabelsPerOrganization = 1000

	// MaxArtifactsPerOrganization defines the maximum number of artifacts per organization
	MaxArtifactsPerOrganization = 100000
)
