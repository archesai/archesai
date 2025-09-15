// Package organizations provides HTTP handlers for organization operations
package organizations

import (
	"context"
	"log/slog"
)

const (
	// Placeholder constants for development
	userPlaceholder = "user-placeholder"
)

// Handler handles HTTP requests for organization operations
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// Ensure Handler implements StrictServerInterface
var _ StrictServerInterface = (*Handler)(nil)

// NewHandler creates a new organizations handler
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// ListOrganizations retrieves organizations (implements StrictServerInterface)
func (h *Handler) ListOrganizations(ctx context.Context, req ListOrganizationsRequestObject) (ListOrganizationsResponseObject, error) {
	// Default pagination
	limit := 50
	offset := 0

	// Handle page-based pagination if provided
	if req.Params.Page.Number > 0 && req.Params.Page.Size > 0 {
		limit = req.Params.Page.Size
		offset = (req.Params.Page.Number - 1) * req.Params.Page.Size
	}

	organizations, total, err := h.service.List(ctx, limit, offset)
	if err != nil {
		h.logger.Error("failed to list organizations", "error", err)
		return nil, err
	}

	// Convert to API entities
	data := make([]Organization, len(organizations))
	for i, org := range organizations {
		data[i] = *org
	}

	totalFloat32 := float32(total)
	return ListOrganizations200JSONResponse{
		Data: data,
		Meta: struct {
			Total float32 `json:"total"`
		}{
			Total: totalFloat32,
		},
	}, nil
}

// CreateOrganization creates a new organization (implements StrictServerInterface)
func (h *Handler) CreateOrganization(ctx context.Context, req CreateOrganizationRequestObject) (CreateOrganizationResponseObject, error) {
	// TODO: Get creator user ID from auth context
	creatorUserID := userPlaceholder

	createReq := &CreateOrganizationRequest{
		OrganizationId: req.Body.OrganizationId,
		BillingEmail:   req.Body.BillingEmail,
	}

	organization, err := h.service.Create(ctx, createReq, creatorUserID)
	if err != nil {
		h.logger.Error("failed to create organization", "error", err)
		return nil, err
	}

	return CreateOrganization201JSONResponse{
		Data: *organization,
	}, nil
}

// GetOneOrganization retrieves an organization by ID (implements StrictServerInterface)
func (h *Handler) GetOneOrganization(ctx context.Context, req GetOneOrganizationRequestObject) (GetOneOrganizationResponseObject, error) {
	organization, err := h.service.Get(ctx, req.Id)
	if err != nil {
		if err == ErrOrganizationNotFound {
			return GetOneOrganization404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Detail: "Organization not found",
					Status: 404,
					Title:  "Organization not found",
				},
			}, nil
		}
		h.logger.Error("failed to get organization", "error", err)
		return nil, err
	}

	return GetOneOrganization200JSONResponse{
		Data: *organization,
	}, nil
}

// UpdateOrganization updates an organization (implements StrictServerInterface)
func (h *Handler) UpdateOrganization(ctx context.Context, req UpdateOrganizationRequestObject) (UpdateOrganizationResponseObject, error) {
	updateReq := &UpdateOrganizationRequest{
		BillingEmail: req.Body.BillingEmail,
	}

	organization, err := h.service.Update(ctx, req.Id, updateReq)
	if err != nil {
		if err == ErrOrganizationNotFound {
			return UpdateOrganization404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Detail: "Organization not found",
					Status: 404,
					Title:  "Organization not found",
				},
			}, nil
		}
		h.logger.Error("failed to update organization", "error", err)
		return nil, err
	}

	return UpdateOrganization200JSONResponse{
		Data: *organization,
	}, nil
}

// DeleteOrganization deletes an organization (implements StrictServerInterface)
func (h *Handler) DeleteOrganization(ctx context.Context, req DeleteOrganizationRequestObject) (DeleteOrganizationResponseObject, error) {
	err := h.service.Delete(ctx, req.Id)
	if err != nil {
		if err == ErrOrganizationNotFound {
			return DeleteOrganization404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Detail: "Organization not found",
					Status: 404,
					Title:  "Organization not found",
				},
			}, nil
		}
		h.logger.Error("failed to delete organization", "error", err)
		return nil, err
	}

	return DeleteOrganization200JSONResponse{}, nil
}
