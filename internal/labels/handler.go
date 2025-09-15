// Package labels provides HTTP handlers for label operations
package labels

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// Handler handles HTTP requests for labels
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new labels handler
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// ListLabels handles GET /content/labels
func (h *Handler) ListLabels(ctx context.Context, request ListLabelsRequestObject) (ListLabelsResponseObject, error) {
	// Get labels with pagination
	limit := 50
	offset := 0

	if request.Params.Page.Size > 0 {
		limit = request.Params.Page.Size
	}
	if request.Params.Page.Number > 0 {
		offset = (request.Params.Page.Number - 1) * limit
	}

	// For now, use a placeholder organization ID
	// TODO: Get from auth context
	orgID := uuid.New()

	labels, total, err := h.service.List(ctx, orgID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list labels", "error", err)
		return ListLabels400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: BadRequestApplicationProblemPlusJSONResponse{
				Type:   "list_failed",
				Title:  "Failed to list labels",
				Detail: "Unable to retrieve labels",
				Status: http.StatusInternalServerError,
			},
		}, nil
	}

	// Convert []*Label to []Label
	labelList := make([]Label, len(labels))
	for i, l := range labels {
		labelList[i] = *l
	}

	response := ListLabels200JSONResponse{
		Data: labelList,
	}
	response.Meta.Total = float32(total)

	return response, nil
}

// GetLabel handles GET /content/labels/{id}
func (h *Handler) GetLabel(ctx context.Context, request GetLabelRequestObject) (GetLabelResponseObject, error) {
	labelID := request.Id

	label, err := h.service.Get(ctx, labelID)
	if err != nil {
		h.logger.Error("failed to get label", "error", err, "label_id", labelID)
		return GetLabel404ApplicationProblemPlusJSONResponse{
			NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
				Type:   "not_found",
				Title:  "Label not found",
				Detail: "The requested label was not found",
				Status: http.StatusNotFound,
			},
		}, nil
	}

	return GetLabel200JSONResponse{
		Data: *label,
	}, nil
}

// CreateLabel handles POST /content/labels
func (h *Handler) CreateLabel(ctx context.Context, request CreateLabelRequestObject) (CreateLabelResponseObject, error) {
	// For now, use a placeholder organization ID
	// TODO: Get from auth context
	orgID := uuid.New()

	createdLabel, err := h.service.Create(ctx, request.Body, orgID)
	if err != nil {
		h.logger.Error("failed to create label", "error", err)
		return CreateLabel400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: BadRequestApplicationProblemPlusJSONResponse{
				Type:   "create_failed",
				Title:  "Failed to create label",
				Detail: "Unable to create label",
				Status: http.StatusInternalServerError,
			},
		}, nil
	}

	return CreateLabel201JSONResponse{
		Data: *createdLabel,
	}, nil
}

// UpdateLabel handles PATCH /content/labels/{id}
func (h *Handler) UpdateLabel(ctx context.Context, request UpdateLabelRequestObject) (UpdateLabelResponseObject, error) {
	labelID := request.Id

	// Get existing label
	label, err := h.service.Get(ctx, labelID)
	if err != nil {
		h.logger.Error("failed to get label for update", "error", err, "label_id", labelID)
		return UpdateLabel404ApplicationProblemPlusJSONResponse{
			NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
				Type:   "not_found",
				Title:  "Label not found",
				Detail: "The requested label was not found",
				Status: http.StatusNotFound,
			},
		}, nil
	}

	// Update fields if provided
	if request.Body.Name != "" {
		label.Name = request.Body.Name
	}

	updatedLabel, err := h.service.Update(ctx, labelID, request.Body)
	if err != nil {
		h.logger.Error("failed to update label", "error", err, "label_id", labelID)
		return UpdateLabel404ApplicationProblemPlusJSONResponse{
			NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
				Type:   "update_failed",
				Title:  "Failed to update label",
				Detail: "Unable to update label",
				Status: http.StatusInternalServerError,
			},
		}, nil
	}

	return UpdateLabel200JSONResponse{
		Data: *updatedLabel,
	}, nil
}

// DeleteLabel handles DELETE /content/labels/{id}
func (h *Handler) DeleteLabel(ctx context.Context, request DeleteLabelRequestObject) (DeleteLabelResponseObject, error) {
	labelID := request.Id

	if err := h.service.Delete(ctx, labelID); err != nil {
		h.logger.Error("failed to delete label", "error", err, "label_id", labelID)
		return DeleteLabel404ApplicationProblemPlusJSONResponse{
			NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
				Type:   "delete_failed",
				Title:  "Failed to delete label",
				Detail: "Unable to delete label",
				Status: http.StatusNotFound,
			},
		}, nil
	}

	return DeleteLabel200JSONResponse{}, nil
}

// Ensure Handler implements StrictServerInterface
var _ StrictServerInterface = (*Handler)(nil)
