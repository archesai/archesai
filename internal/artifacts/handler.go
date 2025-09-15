// Package artifacts provides HTTP handlers for artifact operations
package artifacts

import (
	"context"
	"log/slog"
	"net/http"
)

// Handler handles HTTP requests for artifacts
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new artifacts handler
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// FindManyArtifacts handles GET /content/artifacts
func (h *Handler) FindManyArtifacts(ctx context.Context, request FindManyArtifactsRequestObject) (FindManyArtifactsResponseObject, error) {
	// Get artifacts with pagination
	limit := 50
	offset := 0

	if request.Params.Page.Size > 0 {
		limit = request.Params.Page.Size
	}
	if request.Params.Page.Number > 0 {
		offset = (request.Params.Page.Number - 1) * limit
	}

	artifacts, total, err := h.service.List(ctx, limit, offset)
	if err != nil {
		h.logger.Error("failed to list artifacts", "error", err)
		return FindManyArtifacts400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: BadRequestApplicationProblemPlusJSONResponse{
				Type:   "list_failed",
				Title:  "Failed to list artifacts",
				Detail: "Unable to retrieve artifacts",
				Status: http.StatusInternalServerError,
			},
		}, nil
	}

	// Convert []*Artifact to []Artifact
	artifactList := make([]Artifact, len(artifacts))
	for i, a := range artifacts {
		artifactList[i] = *a
	}

	response := FindManyArtifacts200JSONResponse{
		Data: artifactList,
	}
	response.Meta.Total = float32(total)

	return response, nil
}

// GetOneArtifact handles GET /content/artifacts/{id}
func (h *Handler) GetOneArtifact(ctx context.Context, request GetOneArtifactRequestObject) (GetOneArtifactResponseObject, error) {
	artifactID := request.Id

	artifact, err := h.service.Get(ctx, artifactID)
	if err != nil {
		h.logger.Error("failed to get artifact", "error", err, "artifact_id", artifactID)
		return GetOneArtifact404ApplicationProblemPlusJSONResponse{
			NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
				Type:   "not_found",
				Title:  "Artifact not found",
				Detail: "The requested artifact was not found",
				Status: http.StatusNotFound,
			},
		}, nil
	}

	return GetOneArtifact200JSONResponse{
		Data: *artifact,
	}, nil
}

// CreateArtifact handles POST /content/artifacts
func (h *Handler) CreateArtifact(ctx context.Context, request CreateArtifactRequestObject) (CreateArtifactResponseObject, error) {
	// For now, use placeholder organization and producer IDs
	// TODO: Get these from auth context
	orgID := "default-org"
	producerID := "default-producer"

	createdArtifact, err := h.service.Create(ctx, request.Body, orgID, producerID)
	if err != nil {
		h.logger.Error("failed to create artifact", "error", err)
		return CreateArtifact400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: BadRequestApplicationProblemPlusJSONResponse{
				Type:   "create_failed",
				Title:  "Failed to create artifact",
				Detail: "Unable to create artifact",
				Status: http.StatusInternalServerError,
			},
		}, nil
	}

	return CreateArtifact201JSONResponse{
		Data: *createdArtifact,
	}, nil
}

// UpdateArtifact handles PATCH /content/artifacts/{id}
func (h *Handler) UpdateArtifact(ctx context.Context, request UpdateArtifactRequestObject) (UpdateArtifactResponseObject, error) {
	artifactID := request.Id

	updatedArtifact, err := h.service.Update(ctx, artifactID, request.Body)
	if err != nil {
		h.logger.Error("failed to update artifact", "error", err, "artifact_id", artifactID)
		return UpdateArtifact404ApplicationProblemPlusJSONResponse{
			NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
				Type:   "update_failed",
				Title:  "Failed to update artifact",
				Detail: "Unable to update artifact",
				Status: http.StatusInternalServerError,
			},
		}, nil
	}

	return UpdateArtifact200JSONResponse{
		Data: *updatedArtifact,
	}, nil
}

// DeleteArtifact handles DELETE /content/artifacts/{id}
func (h *Handler) DeleteArtifact(ctx context.Context, request DeleteArtifactRequestObject) (DeleteArtifactResponseObject, error) {
	artifactID := request.Id

	if err := h.service.Delete(ctx, artifactID); err != nil {
		h.logger.Error("failed to delete artifact", "error", err, "artifact_id", artifactID)
		return DeleteArtifact404ApplicationProblemPlusJSONResponse{
			NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
				Type:   "delete_failed",
				Title:  "Failed to delete artifact",
				Detail: "Unable to delete artifact",
				Status: http.StatusNotFound,
			},
		}, nil
	}

	return DeleteArtifact200JSONResponse{}, nil
}

// Ensure Handler implements StrictServerInterface
var _ StrictServerInterface = (*Handler)(nil)
