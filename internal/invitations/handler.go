package invitations

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Handler handles HTTP requests for invitations
type Handler struct {
	service InvitationService
	logger  *slog.Logger
}

// NewHandler creates a new invitation handler
func NewHandler(service InvitationService, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// ListInvitations handles GET /organizations/{id}/invitations
func (h *Handler) ListInvitations(c echo.Context, organizationID openapi_types.UUID, params ListInvitationsParams) error {
	ctx := c.Request().Context()

	// Convert pagination params
	limit := 10
	offset := 0

	if params.Page.Size > 0 {
		limit = params.Page.Size
	}
	if params.Page.Number > 0 {
		offset = (params.Page.Number - 1) * limit
	}

	// Get invitations for the organization
	invitations, err := h.service.ListByOrganization(ctx, organizationID.String())
	if err != nil {
		h.logger.Error("failed to list invitations", "error", err, "organizationId", organizationID)
		return c.JSON(http.StatusInternalServerError, Problem{
			Title:  "Internal Server Error",
			Detail: "Failed to retrieve invitations",
			Status: http.StatusInternalServerError,
			Type:   "https://problems.archesai.com/internal-error",
		})
	}

	// Apply pagination
	total := len(invitations)
	start := offset
	end := start + limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	paginatedInvitations := invitations[start:end]

	// Convert to proper type
	responseData := make([]Invitation, 0, len(paginatedInvitations))
	for _, inv := range paginatedInvitations {
		if inv != nil {
			responseData = append(responseData, *inv)
		}
	}

	return c.JSON(http.StatusOK, ListInvitations200JSONResponse{
		Data: responseData,
		Meta: struct {
			Total float32 `json:"total"`
		}{
			Total: float32(total),
		},
	})
}

// CreateInvitation handles POST /organizations/{id}/invitations
func (h *Handler) CreateInvitation(c echo.Context, organizationID openapi_types.UUID) error {
	ctx := c.Request().Context()

	var req CreateInvitationJSONRequestBody
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, Problem{
			Title:  "Bad Request",
			Detail: "Invalid request body",
			Status: http.StatusBadRequest,
			Type:   "https://problems.archesai.com/bad-request",
		})
	}

	// Create invitation
	invitation := &Invitation{
		Id:             uuid.New(),
		Email:          req.Email,
		Role:           InvitationRole(req.Role),
		OrganizationId: organizationID.String(),
		InviterId:      "", // This should come from auth context
	}

	created, err := h.service.Create(ctx, invitation)
	if err != nil {
		if errors.Is(err, ErrInvitationAlreadyExists) {
			return c.JSON(http.StatusConflict, Problem{
				Title:  "Conflict",
				Detail: "Invitation already exists for this email",
				Status: http.StatusConflict,
				Type:   "https://problems.archesai.com/conflict",
			})
		}
		h.logger.Error("failed to create invitation", "error", err)
		return c.JSON(http.StatusInternalServerError, Problem{
			Title:  "Internal Server Error",
			Detail: "Failed to create invitation",
			Status: http.StatusInternalServerError,
			Type:   "https://problems.archesai.com/internal-error",
		})
	}

	return c.JSON(http.StatusCreated, CreateInvitation201JSONResponse{
		Data: *created,
	})
}

// GetInvitation handles GET /organizations/{id}/invitations/{invitationId}
func (h *Handler) GetInvitation(c echo.Context, organizationID openapi_types.UUID, invitationID openapi_types.UUID) error {
	ctx := c.Request().Context()

	invitation, err := h.service.Get(ctx, invitationID)
	if err != nil {
		if errors.Is(err, ErrInvitationNotFound) {
			return c.JSON(http.StatusNotFound, Problem{
				Title:  "Not Found",
				Detail: "Invitation not found",
				Status: http.StatusNotFound,
				Type:   "https://problems.archesai.com/not-found",
			})
		}
		h.logger.Error("failed to get invitation", "error", err, "id", invitationID)
		return c.JSON(http.StatusInternalServerError, Problem{
			Title:  "Internal Server Error",
			Detail: "Failed to retrieve invitation",
			Status: http.StatusInternalServerError,
			Type:   "https://problems.archesai.com/internal-error",
		})
	}

	// Verify invitation belongs to the organization
	if invitation.OrganizationId != organizationID.String() {
		return c.JSON(http.StatusNotFound, Problem{
			Title:  "Not Found",
			Detail: "Invitation not found",
			Status: http.StatusNotFound,
			Type:   "https://problems.archesai.com/not-found",
		})
	}

	return c.JSON(http.StatusOK, GetInvitation200JSONResponse{
		Data: *invitation,
	})
}

// UpdateInvitation handles PATCH /organizations/{id}/invitations/{invitationId}
func (h *Handler) UpdateInvitation(c echo.Context, organizationID openapi_types.UUID, invitationID openapi_types.UUID) error {
	ctx := c.Request().Context()

	var req UpdateInvitationJSONRequestBody
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, Problem{
			Title:  "Bad Request",
			Detail: "Invalid request body",
			Status: http.StatusBadRequest,
			Type:   "https://problems.archesai.com/bad-request",
		})
	}

	// Get existing invitation first to verify organization
	existing, err := h.service.Get(ctx, invitationID)
	if err != nil {
		if errors.Is(err, ErrInvitationNotFound) {
			return c.JSON(http.StatusNotFound, Problem{
				Title:  "Not Found",
				Detail: "Invitation not found",
				Status: http.StatusNotFound,
				Type:   "https://problems.archesai.com/not-found",
			})
		}
		h.logger.Error("failed to get invitation", "error", err, "id", invitationID)
		return c.JSON(http.StatusInternalServerError, Problem{
			Title:  "Internal Server Error",
			Detail: "Failed to retrieve invitation",
			Status: http.StatusInternalServerError,
			Type:   "https://problems.archesai.com/internal-error",
		})
	}

	// Verify invitation belongs to the organization
	if existing.OrganizationId != organizationID.String() {
		return c.JSON(http.StatusNotFound, Problem{
			Title:  "Not Found",
			Detail: "Invitation not found",
			Status: http.StatusNotFound,
			Type:   "https://problems.archesai.com/not-found",
		})
	}

	// Update invitation fields
	if req.Email != "" {
		existing.Email = req.Email
	}
	if req.Role != "" {
		existing.Role = InvitationRole(req.Role)
	}

	updated, err := h.service.Update(ctx, invitationID, existing)
	if err != nil {
		if errors.Is(err, ErrInvitationExpired) {
			return c.JSON(http.StatusBadRequest, Problem{
				Title:  "Bad Request",
				Detail: "Invitation has expired",
				Status: http.StatusBadRequest,
				Type:   "https://problems.archesai.com/invitation-expired",
			})
		}
		if errors.Is(err, ErrInvitationAlreadyAccepted) {
			return c.JSON(http.StatusBadRequest, Problem{
				Title:  "Bad Request",
				Detail: "Invitation has already been accepted",
				Status: http.StatusBadRequest,
				Type:   "https://problems.archesai.com/invitation-already-accepted",
			})
		}
		h.logger.Error("failed to update invitation", "error", err, "id", invitationID)
		return c.JSON(http.StatusInternalServerError, Problem{
			Title:  "Internal Server Error",
			Detail: "Failed to update invitation",
			Status: http.StatusInternalServerError,
			Type:   "https://problems.archesai.com/internal-error",
		})
	}

	return c.JSON(http.StatusOK, UpdateInvitation200JSONResponse{
		Data: *updated,
	})
}

// DeleteInvitation handles DELETE /organizations/{id}/invitations/{invitationId}
func (h *Handler) DeleteInvitation(c echo.Context, organizationID openapi_types.UUID, invitationID openapi_types.UUID) error {
	ctx := c.Request().Context()

	// Get invitation first to verify organization
	invitation, err := h.service.Get(ctx, invitationID)
	if err != nil {
		if errors.Is(err, ErrInvitationNotFound) {
			return c.JSON(http.StatusNotFound, Problem{
				Title:  "Not Found",
				Detail: "Invitation not found",
				Status: http.StatusNotFound,
				Type:   "https://problems.archesai.com/not-found",
			})
		}
		h.logger.Error("failed to get invitation", "error", err, "id", invitationID)
		return c.JSON(http.StatusInternalServerError, Problem{
			Title:  "Internal Server Error",
			Detail: "Failed to retrieve invitation",
			Status: http.StatusInternalServerError,
			Type:   "https://problems.archesai.com/internal-error",
		})
	}

	// Verify invitation belongs to the organization
	if invitation.OrganizationId != organizationID.String() {
		return c.JSON(http.StatusNotFound, Problem{
			Title:  "Not Found",
			Detail: "Invitation not found",
			Status: http.StatusNotFound,
			Type:   "https://problems.archesai.com/not-found",
		})
	}

	err = h.service.Delete(ctx, invitationID)
	if err != nil {
		if errors.Is(err, ErrInvitationNotFound) {
			return c.JSON(http.StatusNotFound, Problem{
				Title:  "Not Found",
				Detail: "Invitation not found",
				Status: http.StatusNotFound,
				Type:   "https://problems.archesai.com/not-found",
			})
		}
		h.logger.Error("failed to delete invitation", "error", err, "id", invitationID)
		return c.JSON(http.StatusInternalServerError, Problem{
			Title:  "Internal Server Error",
			Detail: "Failed to delete invitation",
			Status: http.StatusInternalServerError,
			Type:   "https://problems.archesai.com/internal-error",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// Ensure Handler implements ServerInterface
var _ ServerInterface = (*Handler)(nil)
