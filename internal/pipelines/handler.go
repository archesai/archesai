// Package pipelines provides HTTP handlers for pipeline operations
package pipelines

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

const (
	// Placeholder constants for development
	orgPlaceholder = "org-placeholder"
)

// Handler handles HTTP requests for workflow operations
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// Ensure Handler implements StrictServerInterface
var _ StrictServerInterface = (*Handler)(nil)

// NewHandler creates a new workflow handler
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// NewWorkflowStrictHandler creates a StrictHandler with middleware
func NewWorkflowStrictHandler(handler StrictServerInterface) ServerInterface {
	return NewStrictHandler(handler, nil)
}

// Pipeline handlers

// ListPipelines retrieves pipelines (implements StrictServerInterface)
func (h *Handler) ListPipelines(ctx context.Context, req ListPipelinesRequestObject) (ListPipelinesResponseObject, error) {
	limit := 50
	offset := 0

	// Handle page-based pagination if provided
	if req.Params.Page.Number > 0 && req.Params.Page.Size > 0 {
		limit = req.Params.Page.Size
		offset = (req.Params.Page.Number - 1) * req.Params.Page.Size
	}

	// TODO: Get organization ID from context
	orgID := orgPlaceholder

	pipelines, total, err := h.service.List(ctx, orgID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list pipelines", "error", err)
		return nil, err
	}

	// Convert to API entities
	data := make([]Pipeline, len(pipelines))
	for i, pipeline := range pipelines {
		data[i] = *pipeline
	}

	totalFloat32 := float32(total)
	return ListPipelines200JSONResponse{
		Data: data,
		Meta: struct {
			Total float32 `json:"total"`
		}{
			Total: totalFloat32,
		},
	}, nil
}

// CreatePipeline creates a new pipeline (implements StrictServerInterface)
func (h *Handler) CreatePipeline(ctx context.Context, req CreatePipelineRequestObject) (CreatePipelineResponseObject, error) {
	// TODO: Get organization ID from context
	orgID := orgPlaceholder

	createReq := &CreatePipelineRequest{
		Name:        req.Body.Name,
		Description: req.Body.Description,
	}

	pipeline, err := h.service.Create(ctx, createReq, orgID)
	if err != nil {
		h.logger.Error("failed to create pipeline", "error", err)
		return nil, err
	}

	return CreatePipeline201JSONResponse{
		Data: *pipeline,
	}, nil
}

// GetPipeline retrieves a pipeline by ID (implements StrictServerInterface)
func (h *Handler) GetPipeline(ctx context.Context, req GetPipelineRequestObject) (GetPipelineResponseObject, error) {
	pipeline, err := h.service.Get(ctx, req.Id)
	if err != nil {
		if err == ErrPipelineNotFound {
			return GetPipeline404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Detail: "Pipeline not found",
					Status: 404,
					Title:  "Pipeline not found",
				},
			}, nil
		}
		h.logger.Error("failed to get pipeline", "error", err)
		return nil, err
	}

	return GetPipeline200JSONResponse{
		Data: *pipeline,
	}, nil
}

// UpdatePipeline updates a pipeline (implements StrictServerInterface)
func (h *Handler) UpdatePipeline(ctx context.Context, req UpdatePipelineRequestObject) (UpdatePipelineResponseObject, error) {
	updateReq := &UpdatePipelineRequest{
		Name:        req.Body.Name,
		Description: req.Body.Description,
	}

	pipeline, err := h.service.Update(ctx, req.Id, updateReq)
	if err != nil {
		if err == ErrPipelineNotFound {
			return UpdatePipeline404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Detail: "Pipeline not found",
					Status: 404,
					Title:  "Pipeline not found",
				},
			}, nil
		}
		h.logger.Error("failed to update pipeline", "error", err)
		return nil, err
	}

	return UpdatePipeline200JSONResponse{
		Data: *pipeline,
	}, nil
}

// DeletePipeline deletes a pipeline (implements StrictServerInterface)
func (h *Handler) DeletePipeline(ctx context.Context, req DeletePipelineRequestObject) (DeletePipelineResponseObject, error) {
	err := h.service.Delete(ctx, req.Id)
	if err != nil {
		if err == ErrPipelineNotFound {
			return DeletePipeline404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Detail: "Pipeline not found",
					Status: 404,
					Title:  "Pipeline not found",
				},
			}, nil
		}
		h.logger.Error("failed to delete pipeline", "error", err)
		return nil, err
	}

	return DeletePipeline200JSONResponse{}, nil
}

// GetPipelineSteps retrieves all steps for a pipeline
func (h *Handler) GetPipelineSteps(_ context.Context, req GetPipelineStepsRequestObject) (GetPipelineStepsResponseObject, error) {
	// TODO: Implement GetPipelineSteps in service
	_ = req
	steps := []PipelineStep{}

	return GetPipelineSteps200JSONResponse{
		Data: steps,
	}, nil
}

// CreatePipelineStep adds a step to a pipeline
func (h *Handler) CreatePipelineStep(_ context.Context, req CreatePipelineStepRequestObject) (CreatePipelineStepResponseObject, error) {
	// TODO: Implement CreatePipelineStep in service
	step := PipelineStep{
		Id:           uuid.New(),
		PipelineId:   req.Id,
		ToolId:       req.Body.ToolId,
		Name:         req.Body.Name,
		Description:  req.Body.Description,
		Config:       req.Body.Config,
		Position:     req.Body.Position,
		Dependencies: req.Body.Dependencies,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	return CreatePipelineStep201JSONResponse{
		Data: step,
	}, nil
}

// GetPipelineExecutionPlan gets the execution plan for a pipeline
func (h *Handler) GetPipelineExecutionPlan(_ context.Context, req GetPipelineExecutionPlanRequestObject) (GetPipelineExecutionPlanResponseObject, error) {
	// TODO: Implement GetPipelineExecutionPlan in service
	_ = req

	return GetPipelineExecutionPlan200JSONResponse{
		Data: struct {
			EstimatedDuration *int `json:"estimatedDuration,omitempty"`
			IsValid           bool `json:"isValid"`
			Levels            []struct {
				Level int                  `json:"level"`
				Steps []openapi_types.UUID `json:"steps"`
			} `json:"levels"`
			PipelineId openapi_types.UUID `json:"pipelineId"` //nolint:revive // matches generated code
			TotalSteps int                `json:"totalSteps"`
		}{
			PipelineId: req.Id,
			IsValid:    true,
			TotalSteps: 0,
			Levels: []struct {
				Level int                  `json:"level"`
				Steps []openapi_types.UUID `json:"steps"`
			}{},
		},
	}, nil
}

// ValidatePipelineExecutionPlan validates a pipeline configuration
func (h *Handler) ValidatePipelineExecutionPlan(_ context.Context, req ValidatePipelineExecutionPlanRequestObject) (ValidatePipelineExecutionPlanResponseObject, error) {
	// TODO: Implement ValidatePipelineExecutionPlan in service
	_ = req

	return ValidatePipelineExecutionPlan200JSONResponse{
		Data: struct {
			Issues *[]string `json:"issues,omitempty"`
			Valid  bool      `json:"valid"`
		}{
			Valid:  true,
			Issues: &[]string{},
		},
	}, nil
}
