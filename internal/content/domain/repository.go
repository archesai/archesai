package domain

import (
	"context"

	"github.com/google/uuid"
)

// ContentRepository defines the interface for content data persistence
type ContentRepository interface {
	// Artifact operations
	CreateArtifact(ctx context.Context, artifact *Artifact) (*Artifact, error)
	GetArtifact(ctx context.Context, id uuid.UUID) (*Artifact, error)
	UpdateArtifact(ctx context.Context, artifact *Artifact) (*Artifact, error)
	DeleteArtifact(ctx context.Context, id uuid.UUID) error
	ListArtifacts(ctx context.Context, orgID string, limit, offset int) ([]*Artifact, int, error)
	SearchArtifacts(ctx context.Context, orgID, query string, limit, offset int) ([]*Artifact, int, error)

	// Label operations
	CreateLabel(ctx context.Context, label *Label) (*Label, error)
	GetLabel(ctx context.Context, id uuid.UUID) (*Label, error)
	GetLabelByName(ctx context.Context, orgID, name string) (*Label, error)
	UpdateLabel(ctx context.Context, label *Label) (*Label, error)
	DeleteLabel(ctx context.Context, id uuid.UUID) error
	ListLabels(ctx context.Context, orgID string, limit, offset int) ([]*Label, int, error)

	// Artifact-Label relationships
	AddLabelToArtifact(ctx context.Context, artifactID, labelID uuid.UUID) error
	RemoveLabelFromArtifact(ctx context.Context, artifactID, labelID uuid.UUID) error
	GetArtifactsByLabel(ctx context.Context, labelID uuid.UUID, limit, offset int) ([]*Artifact, int, error)
	GetLabelsByArtifact(ctx context.Context, artifactID uuid.UUID) ([]*Label, error)
}
