// Package http provides HTTP handlers for organization operations
package organizations

import (
	"context"
	"log/slog"
)

const (
	// Placeholder constants for development
	userPlaceholder = "user-placeholder"
)

// OrganizationHandler handles HTTP requests for organization operations
type OrganizationHandler struct {
	service *Service
	logger  *slog.Logger
}

// Ensure OrganizationHandler implements StrictServerInterface
var _ StrictServerInterface = (*OrganizationHandler)(nil)

// NewOrganizationHandler creates a new organization handler
func NewOrganizationHandler(service *Service, logger *slog.Logger) *OrganizationHandler {
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
	data := make([]Organization, len(orgs))
	for i, org := range orgs {
		data[i] = *org
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

	createReq := &CreateOrganizationRequest{
		BillingEmail:   req.Body.BillingEmail,
		OrganizationId: req.Body.OrganizationId,
	}

	org, err := h.service.CreateOrganization(ctx, createReq, userID)
	if err != nil {
		h.logger.Error("failed to create organization", "error", err)
		return nil, err
	}

	return CreateOrganization201JSONResponse{
		Data: *org,
	}, nil
}

// GetOneOrganization retrieves an organization by ID (implements StrictServerInterface)
func (h *OrganizationHandler) GetOneOrganization(ctx context.Context, req GetOneOrganizationRequestObject) (GetOneOrganizationResponseObject, error) {
	org, err := h.service.repo.GetOrganizationByID(ctx, req.Id)
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
		Data: *org,
	}, nil
}

// UpdateOrganization updates an organization (implements StrictServerInterface)
func (h *OrganizationHandler) UpdateOrganization(ctx context.Context, req UpdateOrganizationRequestObject) (UpdateOrganizationResponseObject, error) {
	updateReq := &UpdateOrganizationRequest{
		BillingEmail:   req.Body.BillingEmail,
		OrganizationId: req.Body.OrganizationId,
	}

	org, err := h.service.UpdateOrganization(ctx, req.Id, updateReq)
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
		Data: *org,
	}, nil
}

// DeleteOrganization deletes an organization (implements StrictServerInterface)
func (h *OrganizationHandler) DeleteOrganization(ctx context.Context, req DeleteOrganizationRequestObject) (DeleteOrganizationResponseObject, error) {
	err := h.service.DeleteOrganization(ctx, req.Id)
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
	data := make([]Member, len(members))
	for i, member := range members {
		data[i] = *member
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
	createReq := &CreateMemberRequest{
		Role: req.Body.Role,
	}

	member, err := h.service.CreateMember(ctx, createReq, req.Id)
	if err != nil {
		if err == ErrMemberExists {
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
		Data: *member,
	}, nil
}

// GetOneMember retrieves a member by ID (implements StrictServerInterface)
func (h *OrganizationHandler) GetOneMember(ctx context.Context, req GetOneMemberRequestObject) (GetOneMemberResponseObject, error) {
	member, err := h.service.GetMember(ctx, req.MemberId)
	if err != nil {
		if err == ErrMemberNotFound {
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
		Data: *member,
	}, nil
}

// UpdateMember updates a member's role (implements StrictServerInterface)
func (h *OrganizationHandler) UpdateMember(ctx context.Context, req UpdateMemberRequestObject) (UpdateMemberResponseObject, error) {
	updateReq := &UpdateMemberRequest{
		Role: req.Body.Role,
	}

	member, err := h.service.UpdateMember(ctx, req.MemberId, updateReq)
	if err != nil {
		if err == ErrMemberNotFound {
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
		Data: *member,
	}, nil
}

// DeleteMember removes a member from an organization (implements StrictServerInterface)
func (h *OrganizationHandler) DeleteMember(ctx context.Context, req DeleteMemberRequestObject) (DeleteMemberResponseObject, error) {
	err := h.service.DeleteMember(ctx, req.MemberId)
	if err != nil {
		if err == ErrMemberNotFound {
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
	data := make([]Invitation, len(invitations))
	for i, invitation := range invitations {
		data[i] = *invitation
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

	createReq := &CreateInvitationRequest{
		Email: req.Body.Email,
		Role:  req.Body.Role,
	}

	invitation, err := h.service.CreateInvitation(ctx, createReq, req.Id.String(), inviterID)
	if err != nil {
		h.logger.Error("failed to create invitation", "error", err)
		return nil, err
	}

	return CreateInvitation201JSONResponse{
		Data: *invitation,
	}, nil
}

// GetOneInvitation retrieves an invitation by ID (implements StrictServerInterface)
func (h *OrganizationHandler) GetOneInvitation(ctx context.Context, req GetOneInvitationRequestObject) (GetOneInvitationResponseObject, error) {
	invitation, err := h.service.GetInvitation(ctx, req.InvitationId)
	if err != nil {
		if err == ErrInvitationNotFound {
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
		Data: *invitation,
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
		if err == ErrInvitationNotFound {
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
