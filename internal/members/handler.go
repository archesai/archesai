// Package members provides HTTP handlers for member operations
package members

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/archesai/archesai/internal/auth"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for members
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new members handler
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// FindManyMembers handles GET /organizations/{id}/members
func (h *Handler) FindManyMembers(ctx context.Context, request FindManyMembersRequestObject) (FindManyMembersResponseObject, error) {
	// Parse organization ID
	organizationID, err := uuid.Parse(request.Id)
	if err != nil {
		return FindManyMembers400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: BadRequestApplicationProblemPlusJSONResponse{
				Type:   "invalid_id",
				Title:  "Invalid organization ID",
				Detail: "The provided organization ID is not valid",
				Status: http.StatusBadRequest,
			},
		}, nil
	}

	members, err := h.service.ListOrganizationMembers(ctx, organizationID)
	if err != nil {
		h.logger.Error("failed to list organization members", "error", err, "organization_id", organizationID)
		return FindManyMembers400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: BadRequestApplicationProblemPlusJSONResponse{
				Type:   "list_failed",
				Title:  "Failed to list members",
				Detail: "Unable to retrieve organization members",
				Status: http.StatusInternalServerError,
			},
		}, nil
	}

	// Convert []*Member to []Member
	memberList := make([]Member, len(members))
	for i, m := range members {
		memberList[i] = *m
	}

	// Calculate total (for now, just the count of members)
	total := float32(len(members))

	response := FindManyMembers200JSONResponse{
		Data: memberList,
	}
	response.Meta.Total = total

	return response, nil
}

// GetOneMember handles GET /organizations/{id}/members/{memberId}
func (h *Handler) GetOneMember(ctx context.Context, request GetOneMemberRequestObject) (GetOneMemberResponseObject, error) {
	// The request.Id is already a UUID type from openapi_types
	organizationID := request.Id
	memberID := request.MemberId

	member, err := h.service.GetMember(ctx, memberID)
	if err != nil {
		h.logger.Error("failed to get member", "error", err, "member_id", memberID)
		return GetOneMember404ApplicationProblemPlusJSONResponse{
			NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
				Type:   "not_found",
				Title:  "Member not found",
				Detail: "The requested member was not found",
				Status: http.StatusNotFound,
			},
		}, nil
	}

	// Parse organization ID from member to verify
	memberOrgID, err := uuid.Parse(member.OrganizationId)
	if err != nil || memberOrgID != organizationID {
		return GetOneMember404ApplicationProblemPlusJSONResponse{
			NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
				Type:   "not_found",
				Title:  "Member not found",
				Detail: "The requested member was not found in this organization",
				Status: http.StatusNotFound,
			},
		}, nil
	}

	return GetOneMember200JSONResponse{
		Data: *member,
	}, nil
}

// CreateMember handles POST /organizations/{id}/members
func (h *Handler) CreateMember(ctx context.Context, request CreateMemberRequestObject) (CreateMemberResponseObject, error) {
	// Parse organization ID
	organizationID, err := uuid.Parse(request.Id)
	if err != nil {
		return CreateMember400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: BadRequestApplicationProblemPlusJSONResponse{
				Type:   "invalid_id",
				Title:  "Invalid organization ID",
				Detail: "The provided organization ID is not valid",
				Status: http.StatusBadRequest,
			},
		}, nil
	}

	// Get authenticated user ID from context
	_, userID, ok := auth.GetAuthContextFromGoContext(ctx)
	if !ok {
		return CreateMember400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: BadRequestApplicationProblemPlusJSONResponse{
				Type:   "auth_required",
				Title:  "Authentication required",
				Detail: "User must be authenticated to create a member",
				Status: http.StatusUnauthorized,
			},
		}, nil
	}

	// Convert role to domain type
	var role MemberRole
	switch request.Body.Role {
	case CreateMemberJSONBodyRoleAdmin:
		role = MemberRoleAdmin
	case CreateMemberJSONBodyRoleMember:
		role = MemberRoleMember
	case CreateMemberJSONBodyRoleOwner:
		role = MemberRoleOwner
	default:
		role = MemberRoleMember
	}

	member, err := h.service.CreateMember(ctx, organizationID, userID, role)
	if err != nil {
		h.logger.Error("failed to create member", "error", err, "organization_id", organizationID, "user_id", userID)
		return CreateMember400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: BadRequestApplicationProblemPlusJSONResponse{
				Type:   "create_failed",
				Title:  "Failed to create member",
				Detail: "Unable to add member to organization",
				Status: http.StatusInternalServerError,
			},
		}, nil
	}

	return CreateMember201JSONResponse{
		Data: *member,
	}, nil
}

// UpdateMember handles PATCH /organizations/{id}/members/{memberId}
func (h *Handler) UpdateMember(ctx context.Context, request UpdateMemberRequestObject) (UpdateMemberResponseObject, error) {
	// The request.Id is already a UUID type from openapi_types
	organizationID := request.Id
	memberID := request.MemberId

	// Get the member first to verify it belongs to the organization
	member, err := h.service.GetMember(ctx, memberID)
	if err != nil {
		h.logger.Error("failed to get member for update", "error", err, "member_id", memberID)
		return UpdateMember404ApplicationProblemPlusJSONResponse{
			NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
				Type:   "not_found",
				Title:  "Member not found",
				Detail: "The requested member was not found",
				Status: http.StatusNotFound,
			},
		}, nil
	}

	// Parse organization ID from member to verify
	memberOrgID, err := uuid.Parse(member.OrganizationId)
	if err != nil || memberOrgID != organizationID {
		return UpdateMember404ApplicationProblemPlusJSONResponse{
			NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
				Type:   "not_found",
				Title:  "Member not found",
				Detail: "The requested member was not found in this organization",
				Status: http.StatusNotFound,
			},
		}, nil
	}

	// Convert role to domain type
	var role MemberRole
	switch request.Body.Role {
	case UpdateMemberJSONBodyRoleAdmin:
		role = MemberRoleAdmin
	case UpdateMemberJSONBodyRoleMember:
		role = MemberRoleMember
	case UpdateMemberJSONBodyRoleOwner:
		role = MemberRoleOwner
	default:
		role = MemberRoleMember
	}

	updatedMember, err := h.service.UpdateMemberRole(ctx, memberID, role)
	if err != nil {
		h.logger.Error("failed to update member", "error", err, "member_id", memberID)
		return UpdateMember404ApplicationProblemPlusJSONResponse{
			NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
				Type:   "update_failed",
				Title:  "Failed to update member",
				Detail: "Unable to update member role",
				Status: http.StatusInternalServerError,
			},
		}, nil
	}

	return UpdateMember200JSONResponse{
		Data: *updatedMember,
	}, nil
}

// DeleteMember handles DELETE /organizations/{id}/members/{memberId}
func (h *Handler) DeleteMember(ctx context.Context, request DeleteMemberRequestObject) (DeleteMemberResponseObject, error) {
	// The request.Id is already a UUID type from openapi_types
	organizationID := request.Id
	memberID := request.MemberId

	// Get the member first to verify it belongs to the organization
	member, err := h.service.GetMember(ctx, memberID)
	if err != nil {
		h.logger.Error("failed to get member for deletion", "error", err, "member_id", memberID)
		return DeleteMember404ApplicationProblemPlusJSONResponse{
			NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
				Type:   "not_found",
				Title:  "Member not found",
				Detail: "The requested member was not found",
				Status: http.StatusNotFound,
			},
		}, nil
	}

	// Parse organization ID from member to verify
	memberOrgID, err := uuid.Parse(member.OrganizationId)
	if err != nil || memberOrgID != organizationID {
		return DeleteMember404ApplicationProblemPlusJSONResponse{
			NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
				Type:   "not_found",
				Title:  "Member not found",
				Detail: "The requested member was not found in this organization",
				Status: http.StatusNotFound,
			},
		}, nil
	}

	if err := h.service.RemoveMember(ctx, memberID); err != nil {
		h.logger.Error("failed to delete member", "error", err, "member_id", memberID)
		return DeleteMember404ApplicationProblemPlusJSONResponse{
			NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
				Type:   "delete_failed",
				Title:  "Failed to delete member",
				Detail: "Unable to remove member from organization",
				Status: http.StatusInternalServerError,
			},
		}, nil
	}

	return DeleteMember200JSONResponse{}, nil
}

// Ensure Handler implements StrictServerInterface
var _ StrictServerInterface = (*Handler)(nil)
