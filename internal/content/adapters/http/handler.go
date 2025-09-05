// Package http provides HTTP handlers for content operations
package http

import (
	"context"
	"log/slog"

	"github.com/archesai/archesai/internal/content/domain"
)

const (
	// Placeholder constants for development
	orgPlaceholder = "org-placeholder"
)

// ContentHandler handles HTTP requests for content operations
type ContentHandler struct {
	service *domain.ContentService
	logger  *slog.Logger
}

// Ensure ContentHandler implements StrictServerInterface
var _ StrictServerInterface = (*ContentHandler)(nil)

// NewContentHandler creates a new content handler
func NewContentHandler(service *domain.ContentService, logger *slog.Logger) *ContentHandler {
	return &ContentHandler{
		service: service,
		logger:  logger,
	}
}

// NewContentStrictHandler creates a StrictHandler with middleware
func NewContentStrictHandler(handler StrictServerInterface) ServerInterface {
	return NewStrictHandler(handler, nil)
}

// Artifact handlers

// FindManyArtifacts retrieves artifacts (implements StrictServerInterface)
func (h *ContentHandler) FindManyArtifacts(ctx context.Context, req FindManyArtifactsRequestObject) (FindManyArtifactsResponseObject, error) {
	limit := 50
	offset := 0

	// Handle page-based pagination if provided
	if req.Params.Page.Number > 0 && req.Params.Page.Size > 0 {
		limit = req.Params.Page.Size
		offset = (req.Params.Page.Number - 1) * req.Params.Page.Size
	}

	// TODO: Get organization ID from context
	orgID := orgPlaceholder

	artifacts, total, err := h.service.ListArtifacts(ctx, orgID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list artifacts", "error", err)
		return nil, err
	}

	// Convert to API entities
	data := make([]domain.ArtifactEntity, len(artifacts))
	for i, artifact := range artifacts {
		data[i] = artifact.ArtifactEntity
	}

	totalFloat32 := float32(total)
	return FindManyArtifacts200JSONResponse{
		Data: data,
		Meta: struct {
			Total float32 `json:"total"`
		}{
			Total: totalFloat32,
		},
	}, nil
}

// CreateArtifact creates a new artifact (implements StrictServerInterface)
func (h *ContentHandler) CreateArtifact(ctx context.Context, req CreateArtifactRequestObject) (CreateArtifactResponseObject, error) {
	// TODO: Get organization ID from context
	orgID := orgPlaceholder
	// TODO: Get producer ID from context
	producerID := "producer-placeholder"

	name := &req.Body.Name
	createReq := &domain.CreateArtifactRequest{
		Name: name,
		Text: req.Body.Text,
	}

	artifact, err := h.service.CreateArtifact(ctx, createReq, orgID, producerID)
	if err != nil {
		h.logger.Error("failed to create artifact", "error", err)
		return nil, err
	}

	return CreateArtifact201JSONResponse{
		Data: artifact.ArtifactEntity,
	}, nil
}

// GetOneArtifact retrieves an artifact by ID (implements StrictServerInterface)
func (h *ContentHandler) GetOneArtifact(ctx context.Context, req GetOneArtifactRequestObject) (GetOneArtifactResponseObject, error) {
	artifact, err := h.service.GetArtifact(ctx, req.Id)
	if err != nil {
		if err == domain.ErrArtifactNotFound {
			return GetOneArtifact404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Detail: "Artifact not found",
					Status: 404,
					Title:  "Artifact not found",
				},
			}, nil
		}
		h.logger.Error("failed to get artifact", "error", err)
		return nil, err
	}

	return GetOneArtifact200JSONResponse{
		Data: artifact.ArtifactEntity,
	}, nil
}

// UpdateArtifact updates an artifact (implements StrictServerInterface)
func (h *ContentHandler) UpdateArtifact(ctx context.Context, req UpdateArtifactRequestObject) (UpdateArtifactResponseObject, error) {
	name := &req.Body.Name
	text := &req.Body.Text
	updateReq := &domain.UpdateArtifactRequest{
		Name: name,
		Text: text,
	}

	artifact, err := h.service.UpdateArtifact(ctx, req.Id, updateReq)
	if err != nil {
		if err == domain.ErrArtifactNotFound {
			return UpdateArtifact404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Detail: "Artifact not found",
					Status: 404,
					Title:  "Artifact not found",
				},
			}, nil
		}
		h.logger.Error("failed to update artifact", "error", err)
		return nil, err
	}

	return UpdateArtifact200JSONResponse{
		Data: artifact.ArtifactEntity,
	}, nil
}

// DeleteArtifact deletes an artifact (implements StrictServerInterface)
func (h *ContentHandler) DeleteArtifact(ctx context.Context, req DeleteArtifactRequestObject) (DeleteArtifactResponseObject, error) {
	err := h.service.DeleteArtifact(ctx, req.Id)
	if err != nil {
		if err == domain.ErrArtifactNotFound {
			return DeleteArtifact404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Detail: "Artifact not found",
					Status: 404,
					Title:  "Artifact not found",
				},
			}, nil
		}
		h.logger.Error("failed to delete artifact", "error", err)
		return nil, err
	}

	return DeleteArtifact200JSONResponse{}, nil
}

// Label handlers

// FindManyLabels retrieves labels (implements StrictServerInterface)
func (h *ContentHandler) FindManyLabels(ctx context.Context, req FindManyLabelsRequestObject) (FindManyLabelsResponseObject, error) {
	limit := 50
	offset := 0

	// Handle page-based pagination if provided
	if req.Params.Page.Number > 0 && req.Params.Page.Size > 0 {
		limit = req.Params.Page.Size
		offset = (req.Params.Page.Number - 1) * req.Params.Page.Size
	}

	// TODO: Get organization ID from context
	orgID := orgPlaceholder

	labels, total, err := h.service.ListLabels(ctx, orgID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list labels", "error", err)
		return nil, err
	}

	// Convert to API entities
	data := make([]domain.LabelEntity, len(labels))
	for i, label := range labels {
		data[i] = label.LabelEntity
	}

	totalFloat32 := float32(total)
	return FindManyLabels200JSONResponse{
		Data: data,
		Meta: struct {
			Total float32 `json:"total"`
		}{
			Total: totalFloat32,
		},
	}, nil
}

// CreateLabel creates a new label (implements StrictServerInterface)
func (h *ContentHandler) CreateLabel(ctx context.Context, req CreateLabelRequestObject) (CreateLabelResponseObject, error) {
	// TODO: Get organization ID from context
	orgID := orgPlaceholder

	createReq := &domain.CreateLabelRequest{
		Name: req.Body.Name,
	}

	label, err := h.service.CreateLabel(ctx, createReq, orgID)
	if err != nil {
		if err == domain.ErrLabelExists {
			// Return 400 bad request since there's no 409 response defined
			return CreateLabel400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: BadRequestApplicationProblemPlusJSONResponse{
					Detail: "Label already exists",
					Status: 400,
					Title:  "Label already exists",
				},
			}, nil
		}
		h.logger.Error("failed to create label", "error", err)
		return nil, err
	}

	return CreateLabel201JSONResponse{
		Data: label.LabelEntity,
	}, nil
}

// GetOneLabel retrieves a label by ID (implements StrictServerInterface)
func (h *ContentHandler) GetOneLabel(ctx context.Context, req GetOneLabelRequestObject) (GetOneLabelResponseObject, error) {
	label, err := h.service.GetLabel(ctx, req.Id)
	if err != nil {
		if err == domain.ErrLabelNotFound {
			return GetOneLabel404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Detail: "Label not found",
					Status: 404,
					Title:  "Label not found",
				},
			}, nil
		}
		h.logger.Error("failed to get label", "error", err)
		return nil, err
	}

	return GetOneLabel200JSONResponse{
		Data: label.LabelEntity,
	}, nil
}

// UpdateLabel updates a label (implements StrictServerInterface)
func (h *ContentHandler) UpdateLabel(ctx context.Context, req UpdateLabelRequestObject) (UpdateLabelResponseObject, error) {
	name := &req.Body.Name
	updateReq := &domain.UpdateLabelRequest{
		Name: name,
	}

	label, err := h.service.UpdateLabel(ctx, req.Id, updateReq)
	if err != nil {
		if err == domain.ErrLabelNotFound {
			return UpdateLabel404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Detail: "Label not found",
					Status: 404,
					Title:  "Label not found",
				},
			}, nil
		}
		h.logger.Error("failed to update label", "error", err)
		return nil, err
	}

	return UpdateLabel200JSONResponse{
		Data: label.LabelEntity,
	}, nil
}

// DeleteLabel deletes a label (implements StrictServerInterface)
func (h *ContentHandler) DeleteLabel(ctx context.Context, req DeleteLabelRequestObject) (DeleteLabelResponseObject, error) {
	err := h.service.DeleteLabel(ctx, req.Id)
	if err != nil {
		if err == domain.ErrLabelNotFound {
			return DeleteLabel404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Detail: "Label not found",
					Status: 404,
					Title:  "Label not found",
				},
			}, nil
		}
		h.logger.Error("failed to delete label", "error", err)
		return nil, err
	}

	return DeleteLabel200JSONResponse{}, nil
}
