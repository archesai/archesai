package labels

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Service provides label management business logic
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService creates a new labels service
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// Create creates a new label for an organization
func (s *Service) Create(ctx context.Context, req *CreateLabelJSONRequestBody, orgID string) (*Label, error) {
	s.logger.Debug("creating label",
		slog.String("name", req.Name),
		slog.String("org", orgID))

	// Validate inputs
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if orgID == "" {
		return nil, errors.New("organization ID is required")
	}

	// Check if label with same name already exists
	// TODO: Implement GetLabelByName when repository supports it
	// existing, err := s.repo.GetLabelByName(ctx, req.Name, orgID)
	// if err == nil && existing != nil {
	// 	return nil, ErrLabelExists
	// }

	label := &Label{
		Id:             uuid.New(),
		Name:           s.sanitizeName(req.Name),
		OrganizationId: orgID,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	// TODO: Add color and description when fields are added to schema

	createdLabel, err := s.repo.Create(ctx, label)
	if err != nil {
		s.logger.Error("failed to create label",
			slog.String("error", err.Error()),
			slog.String("name", req.Name))
		return nil, fmt.Errorf("failed to create label: %w", err)
	}

	s.logger.Info("label created successfully",
		slog.String("id", createdLabel.Id.String()),
		slog.String("name", createdLabel.Name),
		slog.String("org", orgID))

	return createdLabel, nil
}

// Get retrieves a label by ID
func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Label, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid label ID")
	}

	label, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrLabelNotFound
		}
		s.logger.Error("failed to get label",
			slog.String("id", id.String()),
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get label: %w", err)
	}

	return label, nil
}

// Update updates a label
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateLabelJSONRequestBody) (*Label, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid label ID")
	}

	// Get existing label
	label, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// Track if any changes were made
	hasChanges := false

	// Update name if provided
	if req.Name != "" && req.Name != label.Name {
		// TODO: Check for duplicate names
		label.Name = s.sanitizeName(req.Name)
		hasChanges = true
	}

	// TODO: Update color and description when fields are added to schema

	// Only update if changes were made
	if !hasChanges {
		return label, nil
	}

	label.UpdatedAt = time.Now().UTC()

	updatedLabel, err := s.repo.Update(ctx, id, label)
	if err != nil {
		s.logger.Error("failed to update label",
			slog.String("id", id.String()),
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to update label: %w", err)
	}

	s.logger.Info("label updated successfully",
		slog.String("id", id.String()))

	return updatedLabel, nil
}

// Delete deletes a label
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid label ID")
	}

	// Check if label exists
	_, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	// TODO: Check if label is in use by any artifacts
	// This would require querying artifact_labels junction table

	err = s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("failed to delete label",
			slog.String("id", id.String()),
			slog.String("error", err.Error()))
		return fmt.Errorf("failed to delete label: %w", err)
	}

	s.logger.Info("label deleted successfully",
		slog.String("id", id.String()))

	return nil
}

// List retrieves labels for an organization with pagination
func (s *Service) List(ctx context.Context, orgID string, limit, offset int) ([]*Label, int, error) {
	// Validate pagination parameters
	if limit <= 0 {
		limit = 50 // Default limit
	}
	if limit > 500 {
		limit = 500 // Max limit for labels
	}
	if offset < 0 {
		offset = 0
	}

	params := ListLabelsParams{
		Page: PageQuery{
			Number: offset/limit + 1,
			Size:   limit,
		},
		// TODO: Add organization filtering when repository supports it
	}

	labels, total, err := s.repo.List(ctx, params)
	if err != nil {
		s.logger.Error("failed to list labels",
			slog.String("org", orgID),
			slog.String("error", err.Error()))
		return nil, 0, fmt.Errorf("failed to list labels: %w", err)
	}

	// Filter by organization if not handled by repository
	// This is temporary until repository layer supports filtering
	filtered := make([]*Label, 0, len(labels))
	for _, label := range labels {
		if label.OrganizationId == orgID {
			filtered = append(filtered, label)
		}
	}

	return filtered, int(total), nil
}

// GetLabelByName retrieves a label by name within an organization
func (s *Service) GetLabelByName(ctx context.Context, name, orgID string) (*Label, error) {
	if name == "" {
		return nil, errors.New("label name is required")
	}
	if orgID == "" {
		return nil, errors.New("organization ID is required")
	}

	// TODO: Implement when repository supports GetLabelByName
	// For now, list all and filter
	labels, _, err := s.List(ctx, orgID, 1000, 0)
	if err != nil {
		return nil, err
	}

	normalizedName := strings.ToLower(strings.TrimSpace(name))
	for _, label := range labels {
		if strings.ToLower(label.Name) == normalizedName {
			return label, nil
		}
	}

	return nil, ErrLabelNotFound
}

// GetLabelsByArtifact retrieves all labels for a specific artifact
func (s *Service) GetLabelsByArtifact(_ context.Context, artifactID uuid.UUID) ([]*Label, error) {
	if artifactID == uuid.Nil {
		return nil, errors.New("invalid artifact ID")
	}

	// TODO: Implement when many-to-many relationship is added to database
	// This would query the artifact_labels junction table

	return []*Label{}, nil
}

// GetLabelStats returns statistics for labels in an organization
func (s *Service) GetLabelStats(ctx context.Context, orgID string) (map[string]interface{}, error) {
	// Get total count
	_, total, err := s.List(ctx, orgID, 1, 0)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_labels":     total,
		"organization_id":  orgID,
		"max_labels":       MaxLabelsPerOrganization,
		"labels_used":      total,
		"labels_remaining": MaxLabelsPerOrganization - total,
		"timestamp":        time.Now().UTC(),
	}

	return stats, nil
}

// Helper methods

// validateCreateRequest validates the create label request
func (s *Service) validateCreateRequest(req *CreateLabelJSONRequestBody) error {
	if req == nil {
		return errors.New("request is required")
	}

	if req.Name == "" {
		return errors.New("label name is required")
	}

	if len(req.Name) > 100 {
		return errors.New("label name must be 100 characters or less")
	}

	// TODO: Validate description when field is added to schema

	return nil
}

// sanitizeName cleans and validates the label name
func (s *Service) sanitizeName(name string) string {
	// Trim whitespace
	name = strings.TrimSpace(name)

	// Limit length
	if len(name) > 100 {
		name = name[:100]
	}

	return name
}

// isValidColor validates a hex color string
// nolint:unused // kept for future use when color field is added
func (s *Service) isValidColor(color string) bool {
	// Accept colors with or without #
	color = strings.TrimPrefix(color, "#")

	// Check if it's a valid 6-character hex color
	if len(color) != 6 {
		return false
	}

	// Check if all characters are valid hex
	for _, c := range color {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
			return false
		}
	}

	return true
}

// BulkCreateLabels creates multiple labels in a single operation
func (s *Service) BulkCreateLabels(ctx context.Context, labels []*CreateLabelJSONRequestBody, orgID string) ([]*Label, error) {
	if len(labels) == 0 {
		return nil, errors.New("no labels to create")
	}

	if len(labels) > 50 {
		return nil, errors.New("bulk create limited to 50 labels")
	}

	created := make([]*Label, 0, len(labels))
	for _, req := range labels {
		label, err := s.Create(ctx, req, orgID)
		if err != nil {
			// Log error but continue with other labels
			s.logger.Error("failed to create label in bulk operation",
				slog.String("name", req.Name),
				slog.String("error", err.Error()))
			continue
		}
		created = append(created, label)
	}

	if len(created) == 0 {
		return nil, errors.New("all labels failed to create")
	}

	return created, nil
}

// MergeLabels merges one label into another
func (s *Service) MergeLabels(ctx context.Context, sourceID, targetID uuid.UUID) error {
	if sourceID == uuid.Nil || targetID == uuid.Nil {
		return errors.New("invalid label IDs")
	}

	if sourceID == targetID {
		return errors.New("cannot merge label with itself")
	}

	// Verify both labels exist
	source, err := s.Get(ctx, sourceID)
	if err != nil {
		return fmt.Errorf("source label not found: %w", err)
	}

	target, err := s.Get(ctx, targetID)
	if err != nil {
		return fmt.Errorf("target label not found: %w", err)
	}

	// Verify they're in the same organization
	if source.OrganizationId != target.OrganizationId {
		return errors.New("labels must be in the same organization")
	}

	// TODO: Update all artifacts with source label to use target label
	// This would require updating artifact_labels junction table

	// Delete the source label
	if err := s.Delete(ctx, sourceID); err != nil {
		return fmt.Errorf("failed to delete source label: %w", err)
	}

	s.logger.Info("labels merged successfully",
		slog.String("source_id", sourceID.String()),
		slog.String("target_id", targetID.String()))

	return nil
}
