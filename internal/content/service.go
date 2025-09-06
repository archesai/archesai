package content

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Service provides content business logic
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService creates a new content service
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// CreateArtifact creates a new artifact
func (s *Service) CreateArtifact(ctx context.Context, req *CreateArtifactRequest, orgID, producerID string) (*Artifact, error) {
	s.logger.Debug("creating artifact", "org", orgID, "producer", producerID)

	// Validate artifact size
	if len([]byte(req.Text)) > MaxArtifactSize {
		return nil, ErrArtifactTooLarge
	}

	var name string
	if req.Name != "" {
		name = req.Name
	}

	artifact := &Artifact{
		ArtifactEntity: ArtifactEntity{
			Id:             uuid.UUID{}, // Will be set by repository
			Name:           name,
			Text:           req.Text,
			OrganizationId: orgID,
			ProducerId:     producerID,
			Credits:        0.0, // Default credits
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	createdArtifact, err := s.repo.CreateArtifact(ctx, artifact)
	if err != nil {
		s.logger.Error("failed to create artifact", "error", err)
		return nil, fmt.Errorf("failed to create artifact: %w", err)
	}

	// Handle labels if provided
	/* Labels not in generated types
	if len(req.Labels) > 0 {
		for _, labelName := range req.Labels {
			label, err := s.repo.GetLabelByName(ctx, orgID, labelName)
			if err != nil {
				// Create label if it doesn't exist
				newLabel := &Label{
					LabelEntity: LabelEntity{
						Id:             uuid.UUID{}, // Will be set by repository
						Name:           labelName,
						OrganizationId: orgID,
						CreatedAt:      time.Now(),
						UpdatedAt:      time.Now(),
					},
				}
				label, err = s.repo.CreateLabel(ctx, newLabel)
				if err != nil {
					s.logger.Warn("failed to create label", "label", labelName, "error", err)
					continue
				}
			}

			// Add label to artifact
			if err := s.repo.AddLabelToArtifact(ctx, createdArtifact.Id, label.Id); err != nil {
				s.logger.Warn("failed to add label to artifact", "label", labelName, "error", err)
			}
		}
	}
	*/

	s.logger.Info("artifact created successfully", "id", createdArtifact.Id)
	return createdArtifact, nil
}

// GetArtifact retrieves an artifact by ID
func (s *Service) GetArtifact(ctx context.Context, id uuid.UUID) (*Artifact, error) {
	artifact, err := s.repo.GetArtifact(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get artifact: %w", err)
	}
	return artifact, nil
}

// UpdateArtifact updates an artifact
func (s *Service) UpdateArtifact(ctx context.Context, id uuid.UUID, req *UpdateArtifactRequest) (*Artifact, error) {
	artifact, err := s.repo.GetArtifact(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get artifact: %w", err)
	}

	// Update fields that were provided
	if req.Name != "" {
		artifact.Name = req.Name
	}
	if req.Text != "" {
		artifact.Text = req.Text
		// Validate size after update
		if len([]byte(artifact.Text)) > MaxArtifactSize {
			return nil, ErrArtifactTooLarge
		}
	}
	artifact.UpdatedAt = time.Now()

	updatedArtifact, err := s.repo.UpdateArtifact(ctx, artifact)
	if err != nil {
		return nil, fmt.Errorf("failed to update artifact: %w", err)
	}

	// Handle label updates if provided
	// TODO: Implement label updates - would need to compare current vs new labels
	// and add/remove as needed
	// _ = req.Labels - Labels not in generated types

	return updatedArtifact, nil
}

// DeleteArtifact deletes an artifact
func (s *Service) DeleteArtifact(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteArtifact(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete artifact: %w", err)
	}
	return nil
}

// ListArtifacts retrieves artifacts for an organization
func (s *Service) ListArtifacts(ctx context.Context, orgID string, limit, offset int) ([]*Artifact, int, error) {
	artifacts, total, err := s.repo.ListArtifacts(ctx, orgID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list artifacts: %w", err)
	}
	return artifacts, total, nil
}

// SearchArtifacts searches artifacts by text content
func (s *Service) SearchArtifacts(ctx context.Context, orgID, query string, limit, offset int) ([]*Artifact, int, error) {
	artifacts, total, err := s.repo.SearchArtifacts(ctx, orgID, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search artifacts: %w", err)
	}
	return artifacts, total, nil
}

// CreateLabel creates a new label
func (s *Service) CreateLabel(ctx context.Context, req *CreateLabelRequest, orgID string) (*Label, error) {
	s.logger.Debug("creating label", "name", req.Name, "org", orgID)

	// Check if label already exists
	existing, err := s.repo.GetLabelByName(ctx, orgID, req.Name)
	if err == nil && existing != nil {
		return nil, ErrLabelExists
	}

	label := &Label{
		LabelEntity: LabelEntity{
			Id:             uuid.UUID{}, // Will be set by repository
			Name:           req.Name,
			OrganizationId: orgID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	// Validate color if provided
	// Color validation would go here
	if false {
		return nil, fmt.Errorf("invalid color format")
	}

	createdLabel, err := s.repo.CreateLabel(ctx, label)
	if err != nil {
		s.logger.Error("failed to create label", "error", err)
		return nil, fmt.Errorf("failed to create label: %w", err)
	}

	s.logger.Info("label created successfully", "id", createdLabel.Id, "name", createdLabel.Name)
	return createdLabel, nil
}

// GetLabel retrieves a label by ID
func (s *Service) GetLabel(ctx context.Context, id uuid.UUID) (*Label, error) {
	label, err := s.repo.GetLabel(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get label: %w", err)
	}
	return label, nil
}

// UpdateLabel updates a label
func (s *Service) UpdateLabel(ctx context.Context, id uuid.UUID, req *UpdateLabelRequest) (*Label, error) {
	label, err := s.repo.GetLabel(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get label: %w", err)
	}

	// Update fields that were provided
	if req.Name != "" {
		label.Name = req.Name
	}
	label.UpdatedAt = time.Now()

	updatedLabel, err := s.repo.UpdateLabel(ctx, label)
	if err != nil {
		return nil, fmt.Errorf("failed to update label: %w", err)
	}

	return updatedLabel, nil
}

// DeleteLabel deletes a label
func (s *Service) DeleteLabel(ctx context.Context, id uuid.UUID) error {
	// TODO: Consider whether to prevent deletion if label is in use
	err := s.repo.DeleteLabel(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete label: %w", err)
	}
	return nil
}

// ListLabels retrieves labels for an organization
func (s *Service) ListLabels(ctx context.Context, orgID string, limit, offset int) ([]*Label, int, error) {
	labels, total, err := s.repo.ListLabels(ctx, orgID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list labels: %w", err)
	}
	return labels, total, nil
}

// GetArtifactsByLabel retrieves artifacts that have a specific label
func (s *Service) GetArtifactsByLabel(ctx context.Context, labelID uuid.UUID, limit, offset int) ([]*Artifact, int, error) {
	artifacts, total, err := s.repo.GetArtifactsByLabel(ctx, labelID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get artifacts by label: %w", err)
	}
	return artifacts, total, nil
}

// GetLabelsByArtifact retrieves labels for a specific artifact
func (s *Service) GetLabelsByArtifact(ctx context.Context, artifactID uuid.UUID) ([]*Label, error) {
	labels, err := s.repo.GetLabelsByArtifact(ctx, artifactID)
	if err != nil {
		return nil, fmt.Errorf("failed to get labels by artifact: %w", err)
	}
	return labels, nil
}

// AddLabelToArtifact adds a label to an artifact
func (s *Service) AddLabelToArtifact(ctx context.Context, artifactID, labelID uuid.UUID) error {
	err := s.repo.AddLabelToArtifact(ctx, artifactID, labelID)
	if err != nil {
		return fmt.Errorf("failed to add label to artifact: %w", err)
	}
	return nil
}

// RemoveLabelFromArtifact removes a label from an artifact
func (s *Service) RemoveLabelFromArtifact(ctx context.Context, artifactID, labelID uuid.UUID) error {
	err := s.repo.RemoveLabelFromArtifact(ctx, artifactID, labelID)
	if err != nil {
		return fmt.Errorf("failed to remove label from artifact: %w", err)
	}
	return nil
}
