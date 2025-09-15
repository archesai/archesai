// Package runs provides HTTP handlers for run operations
package runs

import (
	"context"
	"log/slog"
)

const (
	// Placeholder constants for development
	orgPlaceholder = "org-placeholder"
)

// Handler handles HTTP requests for run operations
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new handler for run operations
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// CreateRun creates a new run (implements StrictServerInterface)
func (h *Handler) CreateRun(ctx context.Context, req CreateRunRequestObject) (CreateRunResponseObject, error) {
	// TODO: Get organization ID from context
	orgID := orgPlaceholder

	// Create the run
	run, err := h.service.Create(ctx, req.Body, orgID)
	if err != nil {
		h.logger.Error("failed to create run", "error", err)
		return nil, err
	}

	// Convert to API response
	return CreateRun201JSONResponse{
		Data: Run{
			Id:             run.Id,
			PipelineId:     run.PipelineId,
			OrganizationId: run.OrganizationId,
			Status:         run.Status,
			Progress:       run.Progress,
			StartedAt:      run.StartedAt,
			CompletedAt:    run.CompletedAt,
			CreatedAt:      run.CreatedAt,
			UpdatedAt:      run.UpdatedAt,
		},
	}, nil
}

// DeleteRun deletes a run (implements StrictServerInterface)
func (h *Handler) DeleteRun(ctx context.Context, req DeleteRunRequestObject) (DeleteRunResponseObject, error) {
	err := h.service.Delete(ctx, req.Id)
	if err != nil {
		if err == ErrRunNotFound {
			return DeleteRun404ApplicationProblemPlusJSONResponse{}, nil
		}
		h.logger.Error("failed to delete run", "error", err, "id", req.Id)
		return nil, err
	}

	// For delete, we can return an empty response or a success message
	// Since the API expects a Data field, let's return nil for now
	return DeleteRun200JSONResponse{}, nil
}

// GetRun retrieves a single run (implements StrictServerInterface)
func (h *Handler) GetRun(ctx context.Context, req GetRunRequestObject) (GetRunResponseObject, error) {
	run, err := h.service.Get(ctx, req.Id)
	if err != nil {
		if err == ErrRunNotFound {
			return GetRun404ApplicationProblemPlusJSONResponse{}, nil
		}
		h.logger.Error("failed to get run", "error", err, "id", req.Id)
		return nil, err
	}

	// Convert to API response
	return GetRun200JSONResponse{
		Data: Run{
			Id:             run.Id,
			PipelineId:     run.PipelineId,
			OrganizationId: run.OrganizationId,
			Status:         run.Status,
			Progress:       run.Progress,
			StartedAt:      run.StartedAt,
			CompletedAt:    run.CompletedAt,
			CreatedAt:      run.CreatedAt,
			UpdatedAt:      run.UpdatedAt,
		},
	}, nil
}

// UpdateRun updates a run (implements StrictServerInterface)
func (h *Handler) UpdateRun(ctx context.Context, req UpdateRunRequestObject) (UpdateRunResponseObject, error) {
	run, err := h.service.Update(ctx, req.Id, req.Body)
	if err != nil {
		if err == ErrRunNotFound {
			return UpdateRun404ApplicationProblemPlusJSONResponse{}, nil
		}
		h.logger.Error("failed to update run", "error", err, "id", req.Id)
		return nil, err
	}

	// Convert to API response
	return UpdateRun200JSONResponse{
		Data: Run{
			Id:             run.Id,
			PipelineId:     run.PipelineId,
			OrganizationId: run.OrganizationId,
			Status:         run.Status,
			Progress:       run.Progress,
			StartedAt:      run.StartedAt,
			CompletedAt:    run.CompletedAt,
			CreatedAt:      run.CreatedAt,
			UpdatedAt:      run.UpdatedAt,
		},
	}, nil
}

// ListRuns retrieves runs (implements StrictServerInterface)
func (h *Handler) ListRuns(ctx context.Context, req ListRunsRequestObject) (ListRunsResponseObject, error) {
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
	runs, total, err := h.service.List(ctx, orgID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list runs", "error", err)
		return nil, err
	}

	// Convert to API entities
	data := make([]Run, len(runs))
	for i, run := range runs {
		data[i] = *run
	}

	totalFloat32 := float32(total)
	return ListRuns200JSONResponse{
		Data: data,
		Meta: struct {
			Total float32 `json:"total"`
		}{
			Total: totalFloat32,
		},
	}, nil
}

// Create creates a new run (implements StrictServerInterface)
func (h *Handler) Create(ctx context.Context, req CreateRunRequestObject) (CreateRunResponseObject, error) {
	// TODO: Get organization ID from context
	orgID := orgPlaceholder

	run, err := h.service.Create(ctx, req.Body, orgID)
	if err != nil {
		h.logger.Error("failed to create run", "error", err)
		return nil, err
	}

	return CreateRun201JSONResponse{
		Data: *run,
	}, nil
}

// Get retrieves a run by ID (implements StrictServerInterface)
func (h *Handler) Get(ctx context.Context, req GetRunRequestObject) (GetRunResponseObject, error) {
	run, err := h.service.Get(ctx, req.Id)
	if err != nil {
		if err == ErrRunNotFound {
			return GetRun404ApplicationProblemPlusJSONResponse{
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

	return GetRun200JSONResponse{
		Data: *run,
	}, nil
}

// Update updates a run (implements StrictServerInterface)
func (h *Handler) Update(_ context.Context, _ UpdateRunRequestObject) (UpdateRunResponseObject, error) {
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

// Delete deletes a run (implements StrictServerInterface)
func (h *Handler) Delete(ctx context.Context, req DeleteRunRequestObject) (DeleteRunResponseObject, error) {
	err := h.service.Delete(ctx, req.Id)
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
