package runs

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

// Service implements the Service interface for run operations
type Service struct {
	repo   Repository
	logger *slog.Logger
}

// NewService creates a new run service
func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

// List retrieves runs for an organization
func (s *Service) List(ctx context.Context, _ string, limit, offset int) ([]*Run, int64, error) {
	params := ListRunsParams{
		Page: PageQuery{
			Number: offset/limit + 1,
			Size:   limit,
		},
	}

	runs, total, err := s.repo.List(ctx, params)
	if err != nil {
		s.logger.Error("failed to list runs", "error", err)
		return nil, 0, err
	}

	return runs, total, nil
}

// Create creates a new run
func (s *Service) Create(ctx context.Context, req *CreateRunJSONRequestBody, orgID string) (*Run, error) {
	orgUUID, err := uuid.Parse(orgID)
	if err != nil {
		s.logger.Error("invalid organization ID", "orgID", orgID, "error", err)
		return nil, err
	}

	run := &Run{
		Id:             uuid.New(),
		OrganizationId: orgUUID,
		PipelineId:     req.PipelineId,
		Status:         QUEUED,
		Progress:       0,
	}

	createdRun, err := s.repo.Create(ctx, run)
	if err != nil {
		s.logger.Error("failed to create run", "error", err)
		return nil, err
	}

	return createdRun, nil
}

// Get retrieves a run by ID
func (s *Service) Get(ctx context.Context, id uuid.UUID) (*Run, error) {
	run, err := s.repo.Get(ctx, id)
	if err != nil {
		s.logger.Error("failed to get run", "id", id, "error", err)
		return nil, err
	}

	return run, nil
}

// Update updates a run by ID
func (s *Service) Update(ctx context.Context, id uuid.UUID, req *UpdateRunJSONRequestBody) (*Run, error) {
	// Get existing run
	run, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.PipelineId != uuid.Nil {
		run.PipelineId = req.PipelineId
	}

	// Update the run in repository
	updated, err := s.repo.Update(ctx, id, run)
	if err != nil {
		s.logger.Error("failed to update run", "id", id, "error", err)
		return nil, err
	}

	return updated, nil
}

// Delete deletes a run by ID
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.logger.Error("failed to delete run", "id", id, "error", err)
		return err
	}

	return nil
}

// Domain constants
const (
	// MaxRunsToKeep defines how many completed runs to keep per pipeline
	MaxRunsToKeep = 100
)

// Domain type aliases
type (
	// CreateRunRequest represents a request to create a run
	CreateRunRequest = CreateRunJSONBody

	// UpdateRunRequest represents a request to update a run
	UpdateRunRequest = UpdateRunJSONBody
)

// CanStart checks if the run can be started
func (r *Run) CanStart() bool {
	return r.Status == QUEUED
}

// IsRunning checks if the run is currently running
func (r *Run) IsRunning() bool {
	return r.Status == PROCESSING
}

// CanCancel checks if the run can be cancelled
func (r *Run) CanCancel() bool {
	return r.Status == PROCESSING || r.Status == QUEUED
}

// UpdateProgress updates the run's progress
func (r *Run) UpdateProgress(progress float32) {
	r.Progress = progress
}
