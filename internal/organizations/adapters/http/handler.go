// Package http provides HTTP handlers for organization operations
package http

import (
	"context"
	"log/slog"

	"github.com/archesai/archesai/internal/organizations"
)

const (
	// Placeholder constants for development
	userPlaceholder = "user-placeholder"
)

// OrganizationHandler handles HTTP requests for organization operations
type OrganizationHandler struct {
	service *organizations.OrganizationService
	logger  *slog.Logger
}

// Ensure OrganizationHandler implements StrictServerInterface
var _ StrictServerInterface = (*OrganizationHandler)(nil)

// NewOrganizationHandler creates a new organization handler
func NewOrganizationHandler(service *organizations.OrganizationService, logger *slog.Logger) *OrganizationHandler {
	return &OrganizationHandler{
		service: service,
		logger:  logger,
	}
}

// NewOrganizationStrictHandler creates a StrictHandler with middleware
func NewOrganizationStrictHandler(handler StrictServerInterface) ServerInterface {
	return NewStrictHandler(handler, nil)
}

// Organization handlers

// FindManyOrganizations retrieves organizations (implements StrictServerInterface)
func (h *OrganizationHandler) FindManyOrganizations(ctx context.Context, req FindManyOrganizationsRequestObject) (FindManyOrganizationsResponseObject, error) {
	limit := 50
	offset := 0

	// Handle page-based pagination if provided
	if req.Params.Page.Number > 0 && req.Params.Page.Size > 0 {
		limit = req.Params.Page.Size
		offset = (req.Params.Page.Number - 1) * req.Params.Page.Size
	}

	orgs, total, err := h.service.ListOrganizations(ctx, limit, offset)
	if err != nil {
		h.logger.Error("failed to list organizations", "error", err)
		return nil, err
	}

	// Convert to API entities
	data := make([]organizations.OrganizationEntity, len(orgs))
	for i, org := range orgs {
		data[i] = org.OrganizationEntity
	}

	totalFloat32 := float32(total)
	return FindManyOrganizations200JSONResponse{
		Data: data,
		Meta: struct {
			Total float32 `json:"total"`
		}{
			Total: totalFloat32,
		},
	}, nil
}

// CreateOrganization creates a new organization (implements StrictServerInterface)
func (h *OrganizationHandler) CreateOrganization(ctx context.Context, req CreateOrganizationRequestObject) (CreateOrganizationResponseObject, error) {
	// TODO: Get user ID from auth context
	userID := userPlaceholder // This should come from JWT claims

	createReq := &organizations.CreateOrganizationRequest{
		BillingEmail:   req.Body.BillingEmail,
		OrganizationId: req.Body.OrganizationId,
	}

	org, err := h.service.CreateOrganization(ctx, createReq, userID)
	if err != nil {
		h.logger.Error("failed to create organization", "error", err)
		return nil, err
	}

	return CreateOrganization201JSONResponse{
		Data: org.OrganizationEntity,
	}, nil
}

// GetOneOrganization retrieves an organization by ID (implements StrictServerInterface)
func (h *OrganizationHandler) GetOneOrganization(ctx context.Context, req GetOneOrganizationRequestObject) (GetOneOrganizationResponseObject, error) {
	org, err := h.service.GetOrganization(ctx, req.Id)
	if err != nil {
		if err == organizations.ErrOrganizationNotFound {
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
		Data: org.OrganizationEntity,
	}, nil
}

// UpdateOrganization updates an organization (implements StrictServerInterface)
func (h *OrganizationHandler) UpdateOrganization(ctx context.Context, req UpdateOrganizationRequestObject) (UpdateOrganizationResponseObject, error) {
	updateReq := &organizations.UpdateOrganizationRequest{
		BillingEmail:   req.Body.BillingEmail,
		OrganizationId: req.Body.OrganizationId,
	}

	org, err := h.service.UpdateOrganization(ctx, req.Id, updateReq)
	if err != nil {
		if err == organizations.ErrOrganizationNotFound {
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
		Data: org.OrganizationEntity,
	}, nil
}

// DeleteOrganization deletes an organization (implements StrictServerInterface)
func (h *OrganizationHandler) DeleteOrganization(ctx context.Context, req DeleteOrganizationRequestObject) (DeleteOrganizationResponseObject, error) {
	err := h.service.DeleteOrganization(ctx, req.Id)
	if err != nil {
		if err == organizations.ErrOrganizationNotFound {
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

// Member handlers

// FindManyMembers retrieves members of an organization (implements StrictServerInterface)
func (h *OrganizationHandler) FindManyMembers(ctx context.Context, req FindManyMembersRequestObject) (FindManyMembersResponseObject, error) {
	limit := 50
	offset := 0

	// Handle page-based pagination if provided
	if req.Params.Page.Number > 0 && req.Params.Page.Size > 0 {
		limit = req.Params.Page.Size
		offset = (req.Params.Page.Number - 1) * req.Params.Page.Size
	}

	members, total, err := h.service.ListMembers(ctx, req.Id, limit, offset)
	if err != nil {
		h.logger.Error("failed to list members", "error", err)
		return nil, err
	}

	// Convert to API entities
	data := make([]organizations.MemberEntity, len(members))
	for i, member := range members {
		data[i] = member.MemberEntity
	}

	totalFloat32 := float32(total)
	return FindManyMembers200JSONResponse{
		Data: data,
		Meta: struct {
			Total float32 `json:"total"`
		}{
			Total: totalFloat32,
		},
	}, nil
}

// CreateMember adds a member to an organization (implements StrictServerInterface)
func (h *OrganizationHandler) CreateMember(ctx context.Context, req CreateMemberRequestObject) (CreateMemberResponseObject, error) {
	createReq := &organizations.CreateMemberRequest{
		Role: req.Body.Role,
	}

	member, err := h.service.CreateMember(ctx, createReq, req.Id)
	if err != nil {
		if err == organizations.ErrMemberExists {
			// Return 400 bad request since there's no 409 response defined
			return CreateMember400ApplicationProblemPlusJSONResponse{
				BadRequestApplicationProblemPlusJSONResponse: BadRequestApplicationProblemPlusJSONResponse{
					Detail: "Member already exists",
					Status: 400,
					Title:  "Member already exists",
				},
			}, nil
		}
		h.logger.Error("failed to create member", "error", err)
		return nil, err
	}

	return CreateMember201JSONResponse{
		Data: member.MemberEntity,
	}, nil
}

// GetOneMember retrieves a member by ID (implements StrictServerInterface)
func (h *OrganizationHandler) GetOneMember(ctx context.Context, req GetOneMemberRequestObject) (GetOneMemberResponseObject, error) {
	member, err := h.service.GetMember(ctx, req.MemberId)
	if err != nil {
		if err == organizations.ErrMemberNotFound {
			return GetOneMember404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Detail: "Member not found",
					Status: 404,
					Title:  "Member not found",
				},
			}, nil
		}
		h.logger.Error("failed to get member", "error", err)
		return nil, err
	}

	return GetOneMember200JSONResponse{
		Data: member.MemberEntity,
	}, nil
}

// UpdateMember updates a member's role (implements StrictServerInterface)
func (h *OrganizationHandler) UpdateMember(ctx context.Context, req UpdateMemberRequestObject) (UpdateMemberResponseObject, error) {
	updateReq := &organizations.UpdateMemberRequest{
		Role: req.Body.Role,
	}

	member, err := h.service.UpdateMember(ctx, req.MemberId, updateReq)
	if err != nil {
		if err == organizations.ErrMemberNotFound {
			return UpdateMember404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Detail: "Member not found",
					Status: 404,
					Title:  "Member not found",
				},
			}, nil
		}
		h.logger.Error("failed to update member", "error", err)
		return nil, err
	}

	return UpdateMember200JSONResponse{
		Data: member.MemberEntity,
	}, nil
}

// DeleteMember removes a member from an organization (implements StrictServerInterface)
func (h *OrganizationHandler) DeleteMember(ctx context.Context, req DeleteMemberRequestObject) (DeleteMemberResponseObject, error) {
	err := h.service.DeleteMember(ctx, req.MemberId)
	if err != nil {
		if err == organizations.ErrMemberNotFound {
			return DeleteMember404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Detail: "Member not found",
					Status: 404,
					Title:  "Member not found",
				},
			}, nil
		}
		h.logger.Error("failed to delete member", "error", err)
		return nil, err
	}

	return DeleteMember200JSONResponse{}, nil
}

// Invitation handlers

// FindManyInvitations retrieves invitations for an organization (implements StrictServerInterface)
func (h *OrganizationHandler) FindManyInvitations(ctx context.Context, req FindManyInvitationsRequestObject) (FindManyInvitationsResponseObject, error) {
	limit := 50
	offset := 0

	// Handle page-based pagination if provided
	if req.Params.Page.Number > 0 && req.Params.Page.Size > 0 {
		limit = req.Params.Page.Size
		offset = (req.Params.Page.Number - 1) * req.Params.Page.Size
	}

	invitations, total, err := h.service.ListInvitations(ctx, req.Id.String(), limit, offset)
	if err != nil {
		h.logger.Error("failed to list invitations", "error", err)
		return nil, err
	}

	// Convert to API entities
	data := make([]organizations.InvitationEntity, len(invitations))
	for i, invitation := range invitations {
		data[i] = invitation.InvitationEntity
	}

	totalFloat32 := float32(total)
	return FindManyInvitations200JSONResponse{
		Data: data,
		Meta: struct {
			Total float32 `json:"total"`
		}{
			Total: totalFloat32,
		},
	}, nil
}

// CreateInvitation creates a new invitation (implements StrictServerInterface)
func (h *OrganizationHandler) CreateInvitation(ctx context.Context, req CreateInvitationRequestObject) (CreateInvitationResponseObject, error) {
	// TODO: Get inviter ID from auth context
	inviterID := userPlaceholder

	createReq := &organizations.CreateInvitationRequest{
		Email: req.Body.Email,
		Role:  req.Body.Role,
	}

	invitation, err := h.service.CreateInvitation(ctx, createReq, req.Id.String(), inviterID)
	if err != nil {
		h.logger.Error("failed to create invitation", "error", err)
		return nil, err
	}

	return CreateInvitation201JSONResponse{
		Data: invitation.InvitationEntity,
	}, nil
}

// GetOneInvitation retrieves an invitation by ID (implements StrictServerInterface)
func (h *OrganizationHandler) GetOneInvitation(ctx context.Context, req GetOneInvitationRequestObject) (GetOneInvitationResponseObject, error) {
	invitation, err := h.service.GetInvitation(ctx, req.InvitationId)
	if err != nil {
		if err == organizations.ErrInvitationNotFound {
			return GetOneInvitation404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Detail: "Invitation not found",
					Status: 404,
					Title:  "Invitation not found",
				},
			}, nil
		}
		h.logger.Error("failed to get invitation", "error", err)
		return nil, err
	}

	return GetOneInvitation200JSONResponse{
		Data: invitation.InvitationEntity,
	}, nil
}

// UpdateInvitation updates an invitation (implements StrictServerInterface)
func (h *OrganizationHandler) UpdateInvitation(_ context.Context, _ UpdateInvitationRequestObject) (UpdateInvitationResponseObject, error) {
	// TODO: Implement invitation updates when needed
	// Return 404 since we don't support updates yet
	return UpdateInvitation404ApplicationProblemPlusJSONResponse{
		NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
			Detail: "Invitation updates not implemented",
			Status: 404,
			Title:  "Not Implemented",
		},
	}, nil
}

// DeleteInvitation deletes an invitation (implements StrictServerInterface)
func (h *OrganizationHandler) DeleteInvitation(ctx context.Context, req DeleteInvitationRequestObject) (DeleteInvitationResponseObject, error) {
	err := h.service.DeleteInvitation(ctx, req.InvitationId)
	if err != nil {
		if err == organizations.ErrInvitationNotFound {
			return DeleteInvitation404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Detail: "Invitation not found",
					Status: 404,
					Title:  "Invitation not found",
				},
			}, nil
		}
		h.logger.Error("failed to delete invitation", "error", err)
		return nil, err
	}

	return DeleteInvitation200JSONResponse{}, nil
}
