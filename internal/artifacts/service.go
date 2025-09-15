package artifacts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/archesai/archesai/internal/labels"
	"github.com/google/uuid"
)

// Service provides artifact business logic
type Service struct {
	repo      Repository
	labelRepo labels.Repository // For label operations
	logger    *slog.Logger
	maxSize   int // Max artifact size in bytes
}

// ServiceConfig contains configuration for the artifacts service
type ServiceConfig struct {
	MaxArtifactSize int
	CacheEnabled    bool
}

// NewArtifactsService creates a new artifacts service
func NewArtifactsService(repo Repository, labelRepo labels.Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:      repo,
		labelRepo: labelRepo,
		logger:    logger,
		maxSize:   10 * 1024 * 1024,
	}
}

// Create creates a new artifact with validation and processing
func (s *Service) Create(ctx context.Context, req *CreateArtifactJSONRequestBody, orgID, producerID UUID) (*Artifact, error) {
	s.logger.Debug("creating artifact",
		slog.String("org", orgID.String()),
		slog.String("producer", producerID.String()),
		slog.String("name", req.Name))

	// Validate inputs
	if err := s.validateCreateRequest(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Prepare artifact
	artifact := &Artifact{
		Id:             uuid.New(),
		Name:           s.sanitizeName(req.Name),
		Text:           req.Text,
		OrganizationId: orgID,
		ProducerId:     producerID,
		Credits:        float32(s.calculateCredits(req.Text)),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	// Create artifact in repository
	createdArtifact, err := s.repo.Create(ctx, artifact)
	if err != nil {
		s.logger.Error("failed to create artifact",
			slog.String("error", err.Error()),
			slog.String("org", orgID.String()))
		return nil, fmt.Errorf("failed to create artifact: %w", err)
	}

	// Handle labels if provided (future implementation)
	// This would involve a many-to-many relationship table

	s.logger.Info("artifact created successfully",
		slog.String("id", createdArtifact.Id.String()),
		slog.String("name", createdArtifact.Name),
		slog.String("org", orgID.String()))

	return createdArtifact, nil
}

// Get retrieves an artifact by ID with proper error handling
func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Artifact, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid artifact ID")
	}

	artifact, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrArtifactNotFound
		}
		s.logger.Error("failed to get artifact",
			slog.String("id", id.String()),
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to get artifact: %w", err)
	}

	return artifact, nil
}

// Update updates an artifact with validation
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateArtifactJSONRequestBody) (*Artifact, error) {
	if id == uuid.Nil {
		return nil, errors.New("invalid artifact ID")
	}

	// Get existing artifact
	artifact, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// Track if any changes were made
	hasChanges := false

	// Update name if provided
	if req.Name != "" && req.Name != artifact.Name {
		artifact.Name = s.sanitizeName(req.Name)
		hasChanges = true
	}

	// Update text if provided
	if req.Text != "" && req.Text != artifact.Text {
		if err := s.validateTextSize(req.Text); err != nil {
			return nil, err
		}
		artifact.Text = req.Text
		artifact.Credits = float32(s.calculateCredits(req.Text))
		hasChanges = true
	}

	// Only update if changes were made
	if !hasChanges {
		return artifact, nil
	}

	artifact.UpdatedAt = time.Now().UTC()

	// Update in repository
	updatedArtifact, err := s.repo.Update(ctx, id, artifact)
	if err != nil {
		s.logger.Error("failed to update artifact",
			slog.String("id", id.String()),
			slog.String("error", err.Error()))
		return nil, fmt.Errorf("failed to update artifact: %w", err)
	}

	s.logger.Info("artifact updated successfully",
		slog.String("id", id.String()))

	return updatedArtifact, nil
}

// Delete deletes an artifact with proper cleanup
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid artifact ID")
	}

	// Check if artifact exists
	_, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	// Delete artifact (cascade delete should handle relationships)
	err = s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("failed to delete artifact",
			slog.String("id", id.String()),
			slog.String("error", err.Error()))
		return fmt.Errorf("failed to delete artifact: %w", err)
	}

	s.logger.Info("artifact deleted successfully",
		slog.String("id", id.String()))

	return nil
}

// List retrieves artifacts with pagination and filtering
func (s *Service) List(ctx context.Context, limit, offset int) ([]*Artifact, int64, error) {
	// For now, list without organization filtering
	return s.ListByOrganization(ctx, uuid.Nil, limit, offset)
}

// ListByOrganization retrieves artifacts for a specific organization
func (s *Service) ListByOrganization(ctx context.Context, orgID UUID, limit, offset int) ([]*Artifact, int64, error) {
	// Validate pagination parameters
	if limit <= 0 {
		limit = 50 // Default limit
	}
	if limit > 1000 {
		limit = 1000 // Max limit
	}
	if offset < 0 {
		offset = 0
	}

	params := ListArtifactsParams{
		Page: PageQuery{
			Number: offset/limit + 1,
			Size:   limit,
		},
		// TODO: Add organization filtering when repository supports it
	}

	artifacts, total, err := s.repo.List(ctx, params)
	if err != nil {
		s.logger.Error("failed to list artifacts",
			slog.String("org", orgID.String()),
			slog.String("error", err.Error()))
		return nil, 0, fmt.Errorf("failed to list artifacts: %w", err)
	}

	// Filter by organization if not handled by repository
	// This is temporary until repository layer supports filtering
	filtered := make([]*Artifact, 0, len(artifacts))
	for _, artifact := range artifacts {
		if artifact.OrganizationId == orgID {
			filtered = append(filtered, artifact)
		}
	}

	return filtered, total, nil
}

// Search performs full-text search on artifacts
func (s *Service) Search(ctx context.Context, orgID UUID, query string, limit, offset int) ([]*Artifact, int, error) {
	if query == "" {
		artifacts, total, err := s.ListByOrganization(ctx, orgID, limit, offset)
		return artifacts, int(total), err
	}

	// Validate pagination
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100 // Lower limit for search operations
	}

	// TODO: Implement actual search functionality
	// This would typically use PostgreSQL full-text search or a dedicated search engine

	// For now, do a simple in-memory filter (not efficient for production)
	allArtifacts, _, err := s.ListByOrganization(ctx, orgID, 1000, 0)
	if err != nil {
		return nil, 0, err
	}

	queryLower := strings.ToLower(query)
	var results []*Artifact
	for _, artifact := range allArtifacts {
		if strings.Contains(strings.ToLower(artifact.Name), queryLower) ||
			strings.Contains(strings.ToLower(artifact.Text), queryLower) {
			results = append(results, artifact)
		}
	}

	// Apply pagination to results
	start := offset
	if start > len(results) {
		start = len(results)
	}
	end := start + limit
	if end > len(results) {
		end = len(results)
	}

	return results[start:end], len(results), nil
}

// GetArtifactsByLabel retrieves artifacts that have a specific label
func (s *Service) GetArtifactsByLabel(ctx context.Context, labelID uuid.UUID, _, _ int) ([]*Artifact, int, error) {
	if labelID == uuid.Nil {
		return nil, 0, errors.New("invalid label ID")
	}

	// Validate label exists
	_, err := s.labelRepo.Get(ctx, labelID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, ErrLabelNotFound
		}
		return nil, 0, err
	}

	// TODO: Implement when many-to-many relationship is added to database
	// This would query the artifact_labels junction table

	return []*Artifact{}, 0, nil
}

// AddLabelToArtifact associates a label with an artifact
func (s *Service) AddLabelToArtifact(ctx context.Context, artifactID, labelID uuid.UUID) error {
	if artifactID == uuid.Nil || labelID == uuid.Nil {
		return errors.New("invalid artifact or label ID")
	}

	// Validate artifact exists
	_, err := s.Get(ctx, artifactID)
	if err != nil {
		return err
	}

	// Validate label exists
	_, err = s.labelRepo.Get(ctx, labelID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrLabelNotFound
		}
		return err
	}

	// TODO: Implement when many-to-many relationship is added to database
	// This would insert into artifact_labels junction table

	s.logger.Info("label added to artifact",
		slog.String("artifact_id", artifactID.String()),
		slog.String("label_id", labelID.String()))

	return nil
}

// RemoveLabelFromArtifact removes a label association from an artifact
func (s *Service) RemoveLabelFromArtifact(_ context.Context, artifactID, labelID uuid.UUID) error {
	if artifactID == uuid.Nil || labelID == uuid.Nil {
		return errors.New("invalid artifact or label ID")
	}

	// TODO: Implement when many-to-many relationship is added to database
	// This would delete from artifact_labels junction table

	s.logger.Info("label removed from artifact",
		slog.String("artifact_id", artifactID.String()),
		slog.String("label_id", labelID.String()))

	return nil
}

// GetLabelsByArtifact retrieves all labels for a specific artifact
func (s *Service) GetLabelsByArtifact(ctx context.Context, artifactID uuid.UUID) ([]*labels.Label, error) {
	if artifactID == uuid.Nil {
		return nil, errors.New("invalid artifact ID")
	}

	// Validate artifact exists
	_, err := s.Get(ctx, artifactID)
	if err != nil {
		return nil, err
	}

	// TODO: Implement when many-to-many relationship is added to database
	// This would query the artifact_labels junction table

	return []*labels.Label{}, nil
}

// Helper methods

// validateCreateRequest validates the create artifact request
func (s *Service) validateCreateRequest(req *CreateArtifactJSONRequestBody) error {
	if req == nil {
		return errors.New("request is required")
	}

	if req.Text == "" {
		return errors.New("artifact text is required")
	}

	return s.validateTextSize(req.Text)
}

// validateTextSize checks if the text size is within limits
func (s *Service) validateTextSize(text string) error {
	size := len([]byte(text))
	if size > s.maxSize {
		return fmt.Errorf("%w: size %d exceeds maximum %d", ErrArtifactTooLarge, size, s.maxSize)
	}
	return nil
}

// sanitizeName cleans and validates the artifact name
func (s *Service) sanitizeName(name string) string {
	// Trim whitespace
	name = strings.TrimSpace(name)

	// Limit length
	if len(name) > 255 {
		name = name[:255]
	}

	return name
}

// calculateCredits calculates the credit cost for an artifact based on its size
func (s *Service) calculateCredits(text string) float64 {
	// Simple calculation: 1 credit per 1000 characters
	// This can be adjusted based on business requirements
	return float64(len(text)) / 1000.0
}

// BulkCreateArtifacts creates multiple artifacts in a single operation
func (s *Service) BulkCreateArtifacts(ctx context.Context, artifacts []*CreateArtifactJSONRequestBody, orgID, producerID UUID) ([]*Artifact, error) {
	if len(artifacts) == 0 {
		return nil, errors.New("no artifacts to create")
	}

	if len(artifacts) > 100 {
		return nil, errors.New("bulk create limited to 100 artifacts")
	}

	created := make([]*Artifact, 0, len(artifacts))
	for _, req := range artifacts {
		artifact, err := s.Create(ctx, req, orgID, producerID)
		if err != nil {
			// Log error but continue with other artifacts
			s.logger.Error("failed to create artifact in bulk operation",
				slog.String("name", req.Name),
				slog.String("error", err.Error()))
			continue
		}
		created = append(created, artifact)
	}

	if len(created) == 0 {
		return nil, errors.New("all artifacts failed to create")
	}

	return created, nil
}

// GetArtifactStats returns statistics for artifacts in an organization
func (s *Service) GetArtifactStats(ctx context.Context, orgID UUID) (map[string]interface{}, error) {
	// Get total count
	_, total, err := s.ListByOrganization(ctx, orgID, 1, 0)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_artifacts": total,
		"organization_id": orgID,
		"timestamp":       time.Now().UTC(),
	}

	return stats, nil
}
