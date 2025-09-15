// Package tools provides HTTP handlers for tool operations
package tools

import (
	"context"
	"log/slog"
)

const (
	// Placeholder constants for development
	orgPlaceholder = "org-placeholder"
)

// Handler handles HTTP requests for tool operations
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new handler for tool operations
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// ListTools retrieves tools (implements StrictServerInterface)
func (h *Handler) ListTools(ctx context.Context, req ListToolsRequestObject) (ListToolsResponseObject, error) {
	limit := 50
	offset := 0

	// Handle page-based pagination if provided
	if req.Params.Page.Number > 0 && req.Params.Page.Size > 0 {
		limit = req.Params.Page.Size
		offset = (req.Params.Page.Number - 1) * req.Params.Page.Size
	}

	// TODO: Get organization ID from context
	orgID := orgPlaceholder

	tools, total, err := h.service.List(ctx, orgID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list tools", "error", err)
		return nil, err
	}

	// Convert to API entities
	data := make([]Tool, len(tools))
	for i, tool := range tools {
		data[i] = *tool
	}

	totalFloat32 := float32(total)
	return ListTools200JSONResponse{
		Data: data,
		Meta: struct {
			Total float32 `json:"total"`
		}{
			Total: totalFloat32,
		},
	}, nil
}

// CreateTool creates a new tool (implements StrictServerInterface)
func (h *Handler) CreateTool(ctx context.Context, req CreateToolRequestObject) (CreateToolResponseObject, error) {
	// TODO: Get organization ID from context
	orgID := orgPlaceholder

	tool, err := h.service.Create(ctx, req.Body, orgID)
	if err != nil {
		h.logger.Error("failed to create tool", "error", err)
		return CreateTool400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: BadRequestApplicationProblemPlusJSONResponse{
				Detail: "Failed to create tool",
				Status: 400,
				Title:  "Bad Request",
			},
		}, nil
	}

	return CreateTool201JSONResponse{
		Data: *tool,
	}, nil
}

// GetTool retrieves a tool by ID (implements StrictServerInterface)
func (h *Handler) GetTool(ctx context.Context, req GetToolRequestObject) (GetToolResponseObject, error) {
	tool, err := h.service.Get(ctx, req.Id)
	if err != nil {
		if err == ErrToolNotFound {
			return GetTool404ApplicationProblemPlusJSONResponse{
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

	return GetTool200JSONResponse{
		Data: *tool,
	}, nil
}

// UpdateTool updates a tool (implements StrictServerInterface)
func (h *Handler) UpdateTool(ctx context.Context, req UpdateToolRequestObject) (UpdateToolResponseObject, error) {
	tool, err := h.service.Update(ctx, req.Id, req.Body)
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
		Data: *tool,
	}, nil
}

// DeleteTool deletes a tool (implements StrictServerInterface)
func (h *Handler) DeleteTool(ctx context.Context, req DeleteToolRequestObject) (DeleteToolResponseObject, error) {
	err := h.service.Delete(ctx, req.Id)
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
