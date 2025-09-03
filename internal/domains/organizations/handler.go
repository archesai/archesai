// Package organizations provides organization management functionality.
package organizations

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/archesai/archesai/internal/generated/api"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
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

// NewHandler creates a new organization handler
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// RegisterRoutes registers organization routes
func (h *Handler) RegisterRoutes(g *echo.Group) {
	// Organization routes
	g.POST("", h.CreateOrganization)
	g.GET("", h.FindManyOrganizations)
	g.GET("/:id", h.FindOrganizationByID)
	g.PUT("/:id", h.UpdateOrganization)
	g.DELETE("/:id", h.DeleteOrganization)

	// Member routes
	g.POST("/:id/members", h.CreateMember)
	g.GET("/:id/members", h.FindManyMembers)
	g.GET("/:id/members/:memberId", h.FindMemberByID)
	g.PUT("/:id/members/:memberId", h.UpdateMember)
	g.DELETE("/:id/members/:memberId", h.DeleteMember)

	// Invitation routes
	g.POST("/:id/invitations", h.CreateInvitation)
	g.GET("/:id/invitations", h.FindManyInvitations)
	g.GET("/:id/invitations/:invitationId", h.FindInvitationByID)
	g.POST("/:id/invitations/:invitationId/accept", h.AcceptInvitation)
	g.DELETE("/:id/invitations/:invitationId", h.DeleteInvitation)
}

// Organization handlers

// CreateOrganization creates a new organization
func (h *Handler) CreateOrganization(c echo.Context) error {
	var req CreateOrganizationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	// TODO: Get user ID from auth context
	userID := userPlaceholder // This should come from JWT claims

	org, err := h.service.CreateOrganization(c.Request().Context(), &req, userID)
	if err != nil {
		h.logger.Error("failed to create organization", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to create organization",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": org.OrganizationEntity,
	})
}

// FindManyOrganizations retrieves organizations
func (h *Handler) FindManyOrganizations(c echo.Context) error {
	limit := 50
	offset := 0

	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if o := c.QueryParam("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	orgs, total, err := h.service.ListOrganizations(c.Request().Context(), limit, offset)
	if err != nil {
		h.logger.Error("failed to list organizations", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to retrieve organizations",
		})
	}

	// Convert to API entities
	data := make([]api.OrganizationEntity, len(orgs))
	for i, org := range orgs {
		data[i] = org.OrganizationEntity
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"meta": map[string]interface{}{
			"total": total,
		},
	})
}

// FindOrganizationByID retrieves an organization by ID
func (h *Handler) FindOrganizationByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid organization ID",
		})
	}

	org, err := h.service.GetOrganization(c.Request().Context(), id)
	if err != nil {
		if err == ErrOrganizationNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Organization not found",
			})
		}
		h.logger.Error("failed to get organization", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to retrieve organization",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": org.OrganizationEntity,
	})
}

// UpdateOrganization updates an organization
func (h *Handler) UpdateOrganization(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid organization ID",
		})
	}

	var req UpdateOrganizationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	org, err := h.service.UpdateOrganization(c.Request().Context(), id, &req)
	if err != nil {
		if err == ErrOrganizationNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Organization not found",
			})
		}
		h.logger.Error("failed to update organization", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to update organization",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": org.OrganizationEntity,
	})
}

// DeleteOrganization deletes an organization
func (h *Handler) DeleteOrganization(c echo.Context) error {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid organization ID",
		})
	}

	err = h.service.DeleteOrganization(c.Request().Context(), id)
	if err != nil {
		if err == ErrOrganizationNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Organization not found",
			})
		}
		h.logger.Error("failed to delete organization", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to delete organization",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// Member handlers

// CreateMember adds a member to an organization
func (h *Handler) CreateMember(c echo.Context) error {
	orgID := c.Param("id")

	var req CreateMemberRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	member, err := h.service.CreateMember(c.Request().Context(), &req, orgID)
	if err != nil {
		if err == ErrMemberExists {
			return c.JSON(http.StatusConflict, map[string]interface{}{
				"error": "Member already exists",
			})
		}
		h.logger.Error("failed to create member", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to create member",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": member.MemberEntity,
	})
}

// FindManyMembers retrieves members of an organization
func (h *Handler) FindManyMembers(c echo.Context) error {
	orgID := c.Param("id")

	limit := 50
	offset := 0

	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if o := c.QueryParam("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	members, total, err := h.service.ListMembers(c.Request().Context(), orgID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list members", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to retrieve members",
		})
	}

	// Convert to API entities
	data := make([]api.MemberEntity, len(members))
	for i, member := range members {
		data[i] = member.MemberEntity
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"meta": map[string]interface{}{
			"total": total,
		},
	})
}

// FindMemberByID retrieves a member by ID
// FindMemberByID retrieves a member by ID
func (h *Handler) FindMemberByID(c echo.Context) error {
	memberIDParam := c.Param("memberId")
	memberID, err := uuid.Parse(memberIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid member ID",
		})
	}

	member, err := h.service.GetMember(c.Request().Context(), memberID)
	if err != nil {
		if err == ErrMemberNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Member not found",
			})
		}
		h.logger.Error("failed to get member", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to retrieve member",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": member.MemberEntity,
	})
}

// UpdateMember updates a member's role
func (h *Handler) UpdateMember(c echo.Context) error {
	memberIDParam := c.Param("memberId")
	memberID, err := uuid.Parse(memberIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid member ID",
		})
	}

	var req UpdateMemberRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	member, err := h.service.UpdateMember(c.Request().Context(), memberID, &req)
	if err != nil {
		if err == ErrMemberNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Member not found",
			})
		}
		h.logger.Error("failed to update member", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to update member",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": member.MemberEntity,
	})
}

// DeleteMember removes a member from an organization
func (h *Handler) DeleteMember(c echo.Context) error {
	memberIDParam := c.Param("memberId")
	memberID, err := uuid.Parse(memberIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid member ID",
		})
	}

	err = h.service.DeleteMember(c.Request().Context(), memberID)
	if err != nil {
		if err == ErrMemberNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Member not found",
			})
		}
		h.logger.Error("failed to delete member", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to delete member",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// Invitation handlers

// CreateInvitation creates a new invitation
func (h *Handler) CreateInvitation(c echo.Context) error {
	orgID := c.Param("id")

	var req CreateInvitationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	// TODO: Get inviter ID from auth context
	inviterID := userPlaceholder

	invitation, err := h.service.CreateInvitation(c.Request().Context(), &req, orgID, inviterID)
	if err != nil {
		h.logger.Error("failed to create invitation", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to create invitation",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data": invitation.InvitationEntity,
	})
}

// FindManyInvitations retrieves invitations for an organization
func (h *Handler) FindManyInvitations(c echo.Context) error {
	orgID := c.Param("id")

	limit := 50
	offset := 0

	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if o := c.QueryParam("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	invitations, total, err := h.service.ListInvitations(c.Request().Context(), orgID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list invitations", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to retrieve invitations",
		})
	}

	// Convert to API entities
	data := make([]api.InvitationEntity, len(invitations))
	for i, invitation := range invitations {
		data[i] = invitation.InvitationEntity
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": data,
		"meta": map[string]interface{}{
			"total": total,
		},
	})
}

// FindInvitationByID retrieves an invitation by ID
// FindInvitationByID retrieves an invitation by ID
func (h *Handler) FindInvitationByID(c echo.Context) error {
	invitationIDParam := c.Param("invitationId")
	invitationID, err := uuid.Parse(invitationIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid invitation ID",
		})
	}

	invitation, err := h.service.GetInvitation(c.Request().Context(), invitationID)
	if err != nil {
		if err == ErrInvitationNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Invitation not found",
			})
		}
		h.logger.Error("failed to get invitation", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to retrieve invitation",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": invitation.InvitationEntity,
	})
}

// AcceptInvitation accepts an invitation and creates a member
func (h *Handler) AcceptInvitation(c echo.Context) error {
	invitationIDParam := c.Param("invitationId")
	invitationID, err := uuid.Parse(invitationIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid invitation ID",
		})
	}

	// TODO: Get user ID from auth context
	userID := userPlaceholder

	member, err := h.service.AcceptInvitation(c.Request().Context(), invitationID, userID)
	if err != nil {
		if err == ErrInvitationNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Invitation not found",
			})
		}
		if err == ErrInvitationExpired {
			return c.JSON(http.StatusGone, map[string]interface{}{
				"error": "Invitation expired",
			})
		}
		h.logger.Error("failed to accept invitation", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to accept invitation",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": member.MemberEntity,
	})
}

// DeleteInvitation deletes an invitation
func (h *Handler) DeleteInvitation(c echo.Context) error {
	invitationIDParam := c.Param("invitationId")
	invitationID, err := uuid.Parse(invitationIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid invitation ID",
		})
	}

	err = h.service.DeleteInvitation(c.Request().Context(), invitationID)
	if err != nil {
		if err == ErrInvitationNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "Invitation not found",
			})
		}
		h.logger.Error("failed to delete invitation", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to delete invitation",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetInvitation stub for repository interface
func (s *Service) GetInvitation(ctx context.Context, id openapi_types.UUID) (*Invitation, error) {
	invitation, err := s.repo.GetInvitation(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get invitation: %w", err)
	}
	return invitation, nil
}

// DeleteInvitation stub for service interface
func (s *Service) DeleteInvitation(ctx context.Context, id openapi_types.UUID) error {
	err := s.repo.DeleteInvitation(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete invitation: %w", err)
	}
	return nil
}
