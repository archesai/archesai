// Package workflows provides HTTP handlers for workflow operations
package workflows

import (
	"context"
	"log/slog"
)

const (
	// Placeholder constants for development
	orgPlaceholder = "org-placeholder"
)

// WorkflowHandler handles HTTP requests for workflow operations
type WorkflowHandler struct {
	service *WorkflowService
	logger  *slog.Logger
}

// Ensure WorkflowHandler implements StrictServerInterface
var _ StrictServerInterface = (*WorkflowHandler)(nil)

// NewWorkflowHandler creates a new workflow handler
func NewWorkflowHandler(service *WorkflowService, logger *slog.Logger) *WorkflowHandler {
	return &WorkflowHandler{
		service: service,
		logger:  logger,
	}
}

// NewWorkflowStrictHandler creates a StrictHandler with middleware
func NewWorkflowStrictHandler(handler StrictServerInterface) ServerInterface {
	return NewStrictHandler(handler, nil)
}

// Pipeline handlers

// FindManyPipelines retrieves pipelines (implements StrictServerInterface)
func (h *WorkflowHandler) FindManyPipelines(ctx context.Context, req FindManyPipelinesRequestObject) (FindManyPipelinesResponseObject, error) {
	limit := 50
	offset := 0

	// Handle page-based pagination if provided
	if req.Params.Page.Number > 0 && req.Params.Page.Size > 0 {
		limit = req.Params.Page.Size
		offset = (req.Params.Page.Number - 1) * req.Params.Page.Size
	}

	// TODO: Get organization ID from context
	orgID := orgPlaceholder

	pipelines, total, err := h.service.ListPipelines(ctx, orgID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list pipelines", "error", err)
		return nil, err
	}

	// Convert to API entities
	data := make([]PipelineEntity, len(pipelines))
	for i, pipeline := range pipelines {
		data[i] = pipeline.PipelineEntity
	}

	totalFloat32 := float32(total)
	return FindManyPipelines200JSONResponse{
		Data: data,
		Meta: struct {
			Total float32 `json:"total"`
		}{
			Total: totalFloat32,
		},
	}, nil
}

// CreatePipeline creates a new pipeline (implements StrictServerInterface)
func (h *WorkflowHandler) CreatePipeline(ctx context.Context, req CreatePipelineRequestObject) (CreatePipelineResponseObject, error) {
	// TODO: Get organization ID from context
	orgID := orgPlaceholder

	createReq := &CreatePipelineRequest{
		Name:        req.Body.Name,
		Description: req.Body.Description,
	}

	pipeline, err := h.service.CreatePipeline(ctx, createReq, orgID)
	if err != nil {
		h.logger.Error("failed to create pipeline", "error", err)
		return nil, err
	}

	return CreatePipeline201JSONResponse{
		Data: pipeline.PipelineEntity,
	}, nil
}

// GetOnePipeline retrieves a pipeline by ID (implements StrictServerInterface)
func (h *WorkflowHandler) GetOnePipeline(ctx context.Context, req GetOnePipelineRequestObject) (GetOnePipelineResponseObject, error) {
	pipeline, err := h.service.GetPipeline(ctx, req.Id)
	if err != nil {
		if err == ErrPipelineNotFound {
			return GetOnePipeline404ApplicationProblemPlusJSONResponse{
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

	return GetOnePipeline200JSONResponse{
		Data: pipeline.PipelineEntity,
	}, nil
}

// UpdatePipeline updates a pipeline (implements StrictServerInterface)
func (h *WorkflowHandler) UpdatePipeline(ctx context.Context, req UpdatePipelineRequestObject) (UpdatePipelineResponseObject, error) {
	updateReq := &UpdatePipelineRequest{
		Name:        req.Body.Name,
		Description: req.Body.Description,
	}

	pipeline, err := h.service.UpdatePipeline(ctx, req.Id, updateReq)
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
		Data: pipeline.PipelineEntity,
	}, nil
}

// DeletePipeline deletes a pipeline (implements StrictServerInterface)
func (h *WorkflowHandler) DeletePipeline(ctx context.Context, req DeletePipelineRequestObject) (DeletePipelineResponseObject, error) {
	err := h.service.DeletePipeline(ctx, req.Id)
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

// Run handlers

// FindManyRuns retrieves runs (implements StrictServerInterface)
func (h *WorkflowHandler) FindManyRuns(ctx context.Context, req FindManyRunsRequestObject) (FindManyRunsResponseObject, error) {
	limit := 50
	offset := 0

	// Handle page-based pagination if provided
	if req.Params.Page.Number > 0 && req.Params.Page.Size > 0 {
		limit = req.Params.Page.Size
		offset = (req.Params.Page.Number - 1) * req.Params.Page.Size
	}

	// TODO: Get organization ID from context
	orgID := orgPlaceholder

	// TODO: Add filter support when service method is updated
	runs, total, err := h.service.ListRuns(ctx, orgID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list runs", "error", err)
		return nil, err
	}

	// Convert to API entities
	data := make([]RunEntity, len(runs))
	for i, run := range runs {
		data[i] = run.RunEntity
	}

	totalFloat32 := float32(total)
	return FindManyRuns200JSONResponse{
		Data: data,
		Meta: struct {
			Total float32 `json:"total"`
		}{
			Total: totalFloat32,
		},
	}, nil
}

// CreateRun creates a new run (implements StrictServerInterface)
func (h *WorkflowHandler) CreateRun(ctx context.Context, req CreateRunRequestObject) (CreateRunResponseObject, error) {
	// TODO: Get organization ID from context
	orgID := orgPlaceholder

	createReq := &CreateRunRequest{
		PipelineId: req.Body.PipelineId,
	}

	run, err := h.service.CreateRun(ctx, createReq, orgID)
	if err != nil {
		h.logger.Error("failed to create run", "error", err)
		return nil, err
	}

	return CreateRun201JSONResponse{
		Data: run.RunEntity,
	}, nil
}

// GetOneRun retrieves a run by ID (implements StrictServerInterface)
func (h *WorkflowHandler) GetOneRun(ctx context.Context, req GetOneRunRequestObject) (GetOneRunResponseObject, error) {
	run, err := h.service.GetRun(ctx, req.Id)
	if err != nil {
		if err == ErrRunNotFound {
			return GetOneRun404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Detail: "Run not found",
					Status: 404,
					Title:  "Run not found",
				},
			}, nil
		}
		h.logger.Error("failed to get run", "error", err)
		return nil, err
	}

	return GetOneRun200JSONResponse{
		Data: run.RunEntity,
	}, nil
}

// UpdateRun updates a run (implements StrictServerInterface)
func (h *WorkflowHandler) UpdateRun(_ context.Context, _ UpdateRunRequestObject) (UpdateRunResponseObject, error) {
	// Runs are typically not directly updated - their status changes through state transitions
	// Return 404 since we don't support direct updates
	return UpdateRun404ApplicationProblemPlusJSONResponse{
		NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
			Detail: "Run updates not implemented - use state transition endpoints",
			Status: 404,
			Title:  "Not Implemented",
		},
	}, nil
}

// DeleteRun deletes a run (implements StrictServerInterface)
func (h *WorkflowHandler) DeleteRun(ctx context.Context, req DeleteRunRequestObject) (DeleteRunResponseObject, error) {
	err := h.service.DeleteRun(ctx, req.Id)
	if err != nil {
		if err == ErrRunNotFound {
			return DeleteRun404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Detail: "Run not found",
					Status: 404,
					Title:  "Run not found",
				},
			}, nil
		}
		h.logger.Error("failed to delete run", "error", err)
		return nil, err
	}

	return DeleteRun200JSONResponse{}, nil
}

// Tool handlers

// FindManyTools retrieves tools (implements StrictServerInterface)
func (h *WorkflowHandler) FindManyTools(ctx context.Context, req FindManyToolsRequestObject) (FindManyToolsResponseObject, error) {
	limit := 50
	offset := 0

	// Handle page-based pagination if provided
	if req.Params.Page.Number > 0 && req.Params.Page.Size > 0 {
		limit = req.Params.Page.Size
		offset = (req.Params.Page.Number - 1) * req.Params.Page.Size
	}

	// TODO: Get organization ID from context
	orgID := orgPlaceholder

	tools, total, err := h.service.ListTools(ctx, orgID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list tools", "error", err)
		return nil, err
	}

	// Convert to API entities
	data := make([]ToolEntity, len(tools))
	for i, tool := range tools {
		data[i] = tool.ToolEntity
	}

	totalFloat32 := float32(total)
	return FindManyTools200JSONResponse{
		Data: data,
		Meta: struct {
			Total float32 `json:"total"`
		}{
			Total: totalFloat32,
		},
	}, nil
}

// CreateTool creates a new tool (implements StrictServerInterface)
func (h *WorkflowHandler) CreateTool(ctx context.Context, req CreateToolRequestObject) (CreateToolResponseObject, error) {
	// TODO: Get organization ID from context
	orgID := orgPlaceholder

	createReq := &CreateToolRequest{
		Name:        req.Body.Name,
		Description: req.Body.Description,
	}

	tool, err := h.service.CreateTool(ctx, createReq, orgID)
	if err != nil {
		if err == ErrToolExists {
			// Return 400 bad request since there's no 409 response defined
			return CreateTool400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: BadRequestApplicationProblemPlusJSONResponse{
					Detail: "Tool already exists",
					Status: 400,
					Title:  "Tool already exists",
				},
			}, nil
		}
		h.logger.Error("failed to create tool", "error", err)
		return nil, err
	}

	return CreateTool201JSONResponse{
		Data: tool.ToolEntity,
	}, nil
}

// GetOneTool retrieves a tool by ID (implements StrictServerInterface)
func (h *WorkflowHandler) GetOneTool(ctx context.Context, req GetOneToolRequestObject) (GetOneToolResponseObject, error) {
	tool, err := h.service.GetTool(ctx, req.Id)
	if err != nil {
		if err == ErrToolNotFound {
			return GetOneTool404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Detail: "Tool not found",
					Status: 404,
					Title:  "Tool not found",
				},
			}, nil
		}
		h.logger.Error("failed to get tool", "error", err)
		return nil, err
	}

	return GetOneTool200JSONResponse{
		Data: tool.ToolEntity,
	}, nil
}

// UpdateTool updates a tool (implements StrictServerInterface)
func (h *WorkflowHandler) UpdateTool(ctx context.Context, req UpdateToolRequestObject) (UpdateToolResponseObject, error) {
	updateReq := &UpdateToolRequest{
		Name:        req.Body.Name,
		Description: req.Body.Description,
	}

	tool, err := h.service.UpdateTool(ctx, req.Id, updateReq)
	if err != nil {
		if err == ErrToolNotFound {
			return UpdateTool404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Detail: "Tool not found",
					Status: 404,
					Title:  "Tool not found",
				},
			}, nil
		}
		h.logger.Error("failed to update tool", "error", err)
		return nil, err
	}

	return UpdateTool200JSONResponse{
		Data: tool.ToolEntity,
	}, nil
}

// DeleteTool deletes a tool (implements StrictServerInterface)
func (h *WorkflowHandler) DeleteTool(ctx context.Context, req DeleteToolRequestObject) (DeleteToolResponseObject, error) {
	err := h.service.DeleteTool(ctx, req.Id)
	if err != nil {
		if err == ErrToolNotFound {
			return DeleteTool404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Detail: "Tool not found",
					Status: 404,
					Title:  "Tool not found",
				},
			}, nil
		}
		h.logger.Error("failed to delete tool", "error", err)
		return nil, err
	}

	return DeleteTool200JSONResponse{}, nil
}
