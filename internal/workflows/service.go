package workflows

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

// WorkflowService provides workflow business logic
type WorkflowService struct {
	repo   WorkflowRepository
	logger *slog.Logger
}

// NewWorkflowService creates a new workflow service
func NewWorkflowService(repo WorkflowRepository, logger *slog.Logger) *WorkflowService {
	return &WorkflowService{
		repo:   repo,
		logger: logger,
	}
}

// CreatePipeline creates a new pipeline
func (s *WorkflowService) CreatePipeline(ctx context.Context, req *CreatePipelineRequest, orgID string) (*Pipeline, error) {
	s.logger.Debug("creating pipeline", "name", req.Name, "org", orgID)

	pipeline := &Pipeline{
		PipelineEntity: PipelineEntity{
			Id:             UUID{}, // Will be set by repository
			Name:           req.Name,
			Description:    req.Description,
			OrganizationId: UUID{}, // Convert orgID to UUID
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	createdPipeline, err := s.repo.CreatePipeline(ctx, pipeline)
	if err != nil {
		s.logger.Error("failed to create pipeline", "error", err)
		return nil, fmt.Errorf("failed to create pipeline: %w", err)
	}

	s.logger.Info("pipeline created successfully", "id", createdPipeline.Id, "name", createdPipeline.Name)
	return createdPipeline, nil
}

// GetPipeline retrieves a pipeline by ID
func (s *WorkflowService) GetPipeline(ctx context.Context, id uuid.UUID) (*Pipeline, error) {
	pipeline, err := s.repo.GetPipeline(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline: %w", err)
	}
	return pipeline, nil
}

// UpdatePipeline updates a pipeline
func (s *WorkflowService) UpdatePipeline(ctx context.Context, id uuid.UUID, req *UpdatePipelineRequest) (*Pipeline, error) {
	pipeline, err := s.repo.GetPipeline(ctx, id)
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

	updatedPipeline, err := s.repo.UpdatePipeline(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to update pipeline: %w", err)
	}

	return updatedPipeline, nil
}

// DeletePipeline deletes a pipeline
func (s *WorkflowService) DeletePipeline(ctx context.Context, id uuid.UUID) error {
	// TODO: Add checks for active runs
	err := s.repo.DeletePipeline(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete pipeline: %w", err)
	}
	return nil
}

// ListPipelines retrieves pipelines for an organization
func (s *WorkflowService) ListPipelines(ctx context.Context, orgID string, limit, offset int) ([]*Pipeline, int, error) {
	pipelines, total, err := s.repo.ListPipelines(ctx, orgID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list pipelines: %w", err)
	}
	return pipelines, total, nil
}

// CreateRun creates a new run
func (s *WorkflowService) CreateRun(ctx context.Context, req *CreateRunRequest, orgID string) (*Run, error) {
	s.logger.Debug("creating run", "pipeline_id", req.PipelineId, "org", orgID)

	// Validate pipeline exists
	pipelineUUID, err := uuid.Parse(req.PipelineId)
	if err != nil {
		return nil, fmt.Errorf("invalid pipeline ID: %w", err)
	}

	pipeline, err := s.repo.GetPipeline(ctx, pipelineUUID)
	if err != nil {
		return nil, fmt.Errorf("pipeline not found: %w", err)
	}

	run := &Run{
		RunEntity: RunEntity{
			Id:             uuid.UUID{}, // Will be set by repository
			PipelineId:     req.PipelineId,
			OrganizationId: orgID,
			ToolId:         "",     // No tool ID in request
			Status:         QUEUED, // Use QUEUED instead of Pending
			Progress:       0.0,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	createdRun, err := s.repo.CreateRun(ctx, run)
	if err != nil {
		s.logger.Error("failed to create run", "error", err)
		return nil, fmt.Errorf("failed to create run: %w", err)
	}

	s.logger.Info("run created successfully", "id", createdRun.Id, "pipeline", pipeline.Name)
	return createdRun, nil
}

// GetRun retrieves a run by ID
func (s *WorkflowService) GetRun(ctx context.Context, id uuid.UUID) (*Run, error) {
	run, err := s.repo.GetRun(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get run: %w", err)
	}
	return run, nil
}

// StartRun starts a pending run
func (s *WorkflowService) StartRun(ctx context.Context, id uuid.UUID) (*Run, error) {
	run, err := s.repo.GetRun(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get run: %w", err)
	}

	if !run.CanStart() {
		return nil, fmt.Errorf("run cannot be started in current status: %s", run.Status)
	}

	run.Status = PROCESSING
	run.StartedAt = time.Now()
	run.UpdatedAt = time.Now()

	updatedRun, err := s.repo.UpdateRun(ctx, run)
	if err != nil {
		return nil, fmt.Errorf("failed to start run: %w", err)
	}

	// TODO: Trigger actual workflow execution

	return updatedRun, nil
}

// UpdateRunProgress updates a run's progress
func (s *WorkflowService) UpdateRunProgress(ctx context.Context, id uuid.UUID, progress float32) (*Run, error) {
	run, err := s.repo.GetRun(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get run: %w", err)
	}

	run.UpdateProgress(progress)

	updatedRun, err := s.repo.UpdateRun(ctx, run)
	if err != nil {
		return nil, fmt.Errorf("failed to update run progress: %w", err)
	}

	return updatedRun, nil
}

// CompleteRun marks a run as completed
func (s *WorkflowService) CompleteRun(ctx context.Context, id uuid.UUID, _ map[string]interface{}) (*Run, error) {
	run, err := s.repo.GetRun(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get run: %w", err)
	}

	if !run.IsRunning() {
		return nil, fmt.Errorf("run is not running")
	}

	run.Status = COMPLETED
	run.Progress = 100.0
	// Note: RunEntity doesn't have Output field in API
	run.CompletedAt = time.Now()
	run.UpdatedAt = time.Now()

	updatedRun, err := s.repo.UpdateRun(ctx, run)
	if err != nil {
		return nil, fmt.Errorf("failed to complete run: %w", err)
	}

	s.logger.Info("run completed successfully", "id", run.Id)
	return updatedRun, nil
}

// CancelRun cancels a run
func (s *WorkflowService) CancelRun(ctx context.Context, id uuid.UUID) (*Run, error) {
	run, err := s.repo.GetRun(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get run: %w", err)
	}

	if !run.CanCancel() {
		return nil, fmt.Errorf("run cannot be cancelled in current status: %s", run.Status)
	}

	run.Status = FAILED // Use FAILED for cancelled runs
	run.CompletedAt = time.Now()
	run.UpdatedAt = time.Now()

	updatedRun, err := s.repo.UpdateRun(ctx, run)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel run: %w", err)
	}

	s.logger.Info("run cancelled", "id", run.Id)
	return updatedRun, nil
}

// ListRuns retrieves runs for an organization
func (s *WorkflowService) ListRuns(ctx context.Context, orgID string, limit, offset int) ([]*Run, int, error) {
	runs, total, err := s.repo.ListRuns(ctx, orgID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list runs: %w", err)
	}
	return runs, total, nil
}

// ListRunsByPipeline retrieves runs for a specific pipeline
func (s *WorkflowService) ListRunsByPipeline(ctx context.Context, pipelineID string, limit, offset int) ([]*Run, int, error) {
	runs, total, err := s.repo.ListRunsByPipeline(ctx, pipelineID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list runs by pipeline: %w", err)
	}
	return runs, total, nil
}

// CreateTool creates a new tool
func (s *WorkflowService) CreateTool(ctx context.Context, req *CreateToolRequest, orgID string) (*Tool, error) {
	s.logger.Debug("creating tool", "name", req.Name, "org", orgID)

	tool := &Tool{
		ToolEntity: ToolEntity{
			Id:             uuid.UUID{}, // Will be set by repository
			Name:           req.Name,
			Description:    req.Description,
			OrganizationId: orgID,
			// Note: Config field not available in ToolEntity API
			InputMimeType:  "application/json", // Default mime type
			OutputMimeType: "application/json", // Default mime type
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		},
	}

	createdTool, err := s.repo.CreateTool(ctx, tool)
	if err != nil {
		s.logger.Error("failed to create tool", "error", err)
		return nil, fmt.Errorf("failed to create tool: %w", err)
	}

	s.logger.Info("tool created successfully", "id", createdTool.Id, "name", createdTool.Name)
	return createdTool, nil
}

// GetTool retrieves a tool by ID
func (s *WorkflowService) GetTool(ctx context.Context, id uuid.UUID) (*Tool, error) {
	tool, err := s.repo.GetTool(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tool: %w", err)
	}
	return tool, nil
}

// UpdateTool updates a tool
func (s *WorkflowService) UpdateTool(ctx context.Context, id uuid.UUID, req *UpdateToolRequest) (*Tool, error) {
	tool, err := s.repo.GetTool(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tool: %w", err)
	}

	if req.Name != "" {
		tool.Name = req.Name
	}
	if req.Description != "" {
		tool.Description = req.Description
	}
	// Note: Config field not available in ToolEntity API
	tool.UpdatedAt = time.Now()

	updatedTool, err := s.repo.UpdateTool(ctx, tool)
	if err != nil {
		return nil, fmt.Errorf("failed to update tool: %w", err)
	}

	return updatedTool, nil
}

// DeleteRun deletes a run
func (s *WorkflowService) DeleteRun(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteRun(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete run: %w", err)
	}
	return nil
}

// DeleteTool deletes a tool
func (s *WorkflowService) DeleteTool(ctx context.Context, id uuid.UUID) error {
	// TODO: Add checks for active runs using this tool
	err := s.repo.DeleteTool(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete tool: %w", err)
	}
	return nil
}

// ListTools retrieves tools for an organization
func (s *WorkflowService) ListTools(ctx context.Context, orgID string, limit, offset int) ([]*Tool, int, error) {
	tools, total, err := s.repo.ListTools(ctx, orgID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list tools: %w", err)
	}
	return tools, total, nil
}
