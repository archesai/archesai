// Package domain contains the content domain business logic and entities
package domain

import (
	"errors"

	"github.com/archesai/archesai/internal/content/generated/api"
)

// Domain constants
const (
	MaxArtifactSize = 10 * 1024 * 1024 // 10MB maximum artifact size
)

// Domain errors
var (
	ErrArtifactNotFound = errors.New("artifact not found")
	ErrArtifactExists   = errors.New("artifact already exists")
	ErrArtifactTooLarge = errors.New("artifact exceeds maximum size")
	ErrLabelNotFound    = errors.New("label not found")
	ErrLabelExists      = errors.New("label already exists")
	ErrInvalidArtifact  = errors.New("invalid artifact")
	ErrInvalidLabel     = errors.New("invalid label")
)

// Artifact extends the generated API ArtifactEntity with domain-specific fields
type Artifact struct {
	api.ArtifactEntity
	// Add any domain-specific fields that aren't in the API
}

// Label extends the generated API LabelEntity with domain-specific fields
type Label struct {
	api.LabelEntity
	// Add any domain-specific fields that aren't in the API
}

// CreateArtifactRequest represents a request to create an artifact
type CreateArtifactRequest struct {
	Name   *string  `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Text   string   `json:"text" validate:"required,min=1"`
	Labels []string `json:"labels,omitempty"`
}

// UpdateArtifactRequest represents a request to update an artifact
type UpdateArtifactRequest struct {
	Name   *string  `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Text   *string  `json:"text,omitempty" validate:"omitempty,min=1"`
	Labels []string `json:"labels,omitempty"`
}

// CreateLabelRequest represents a request to create a label
type CreateLabelRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
	Color       *string `json:"color,omitempty" validate:"omitempty,hexcolor"`
}

// UpdateLabelRequest represents a request to update a label
type UpdateLabelRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=500"`
	Color       *string `json:"color,omitempty" validate:"omitempty,hexcolor"`
}

// NewArtifact creates a new artifact from the API entity
func NewArtifact(entity api.ArtifactEntity) *Artifact {
	return &Artifact{ArtifactEntity: entity}
}

// NewLabel creates a new label from the API entity
func NewLabel(entity api.LabelEntity) *Label {
	return &Label{LabelEntity: entity}
}

// GetSize returns the size of the artifact text in bytes
func (a *Artifact) GetSize() int {
	return len([]byte(a.Text))
}

// IsTooBig checks if the artifact exceeds the maximum size
func (a *Artifact) IsTooBig() bool {
	return a.GetSize() > MaxArtifactSize
}

// HasName checks if the artifact has a name
func (a *Artifact) HasName() bool {
	return a.Name != ""
}

// GetDisplayName returns the artifact name or a default based on ID
func (a *Artifact) GetDisplayName() string {
	if a.HasName() {
		return a.Name
	}
	return "Artifact " + a.Id.String()
}

// AddLabel adds a label to the artifact (if labels were stored directly)
// This would be used if labels were embedded in the artifact
func (a *Artifact) AddLabel(_ string) {
	// This is a conceptual method - actual implementation would depend on
	// whether labels are stored as relationships or embedded data
}

// RemoveLabel removes a label from the artifact
func (a *Artifact) RemoveLabel(_ string) {
	// This is a conceptual method - actual implementation would depend on
	// whether labels are stored as relationships or embedded data
}

// IsValidColor always returns true since LabelEntity doesn't have Color field in API
func (l *Label) IsValidColor() bool {
	return true // Color field not present in API
}

// GetDisplayColor returns a default color since LabelEntity doesn't have Color field
func (l *Label) GetDisplayColor() string {
	return "#6B7280" // Default gray color
}
