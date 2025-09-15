package pipelines

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// Service provides pipeline business logic
type Service struct {
	pipelineRepo Repository
	logger       *slog.Logger
}

// NewService creates a new pipeline service
func NewService(pipelineRepo Repository, logger *slog.Logger) *Service {
	return &Service{
		pipelineRepo: pipelineRepo,
		logger:       logger,
	}
}

// Create creates a new pipeline
func (s *Service) Create(ctx context.Context, req *CreatePipelineRequest, orgID string) (*Pipeline, error) {
	s.logger.Debug("creating pipeline", "name", req.Name, "org", orgID)

	pipeline := &Pipeline{
		Id:             UUID{}, // Will be set by repository
		Name:           req.Name,
		Description:    req.Description,
		OrganizationId: UUID{}, // Convert orgID to UUID
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	createdPipeline, err := s.pipelineRepo.Create(ctx, pipeline)
	if err != nil {
		s.logger.Error("failed to create pipeline", "error", err)
		return nil, fmt.Errorf("failed to create pipeline: %w", err)
	}

	s.logger.Info("pipeline created successfully", "id", createdPipeline.Id, "name", createdPipeline.Name)
	return createdPipeline, nil
}

// Get retrieves a pipeline by ID
func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Pipeline, error) {
	pipeline, err := s.pipelineRepo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline: %w", err)
	}
	return pipeline, nil
}

// Update updates a pipeline
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdatePipelineRequest) (*Pipeline, error) {
	pipeline, err := s.pipelineRepo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline: %w", err)
	}

	if req.Name != "" {
		pipeline.Name = req.Name
	}
	if req.Description != "" {
		pipeline.Description = req.Description
	}
	pipeline.UpdatedAt = time.Now()

	updatedPipeline, err := s.pipelineRepo.Update(ctx, id, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to update pipeline: %w", err)
	}

	return updatedPipeline, nil
}

// Delete deletes a pipeline
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	// TODO: Add checks for active runs
	err := s.pipelineRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete pipeline: %w", err)
	}
	return nil
}

// List retrieves pipelines for an organization
func (s *Service) List(ctx context.Context, orgID string, limit, offset int) ([]*Pipeline, int, error) {
	// TODO: Add organization filtering to repository when available
	_ = orgID
	params := ListPipelinesParams{
		Limit:  limit,
		Offset: offset,
	}
	pipelines, total, err := s.pipelineRepo.List(ctx, params)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list pipelines: %w", err)
	}
	return pipelines, int(total), nil
}
