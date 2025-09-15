package users

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

// Handler provides HTTP handlers for user operations
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// Ensure Handler implements StrictServerInterface
var _ StrictServerInterface = (*Handler)(nil)

// NewHandler creates a new user handler
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// NewUserStrictHandler creates a StrictHandler with middleware
func NewUserStrictHandler(handler StrictServerInterface) ServerInterface {
	return NewStrictHandler(handler, nil)
}

// GetOneUser handles getting a single user
func (h *Handler) GetOneUser(ctx context.Context, req GetOneUserRequestObject) (GetOneUserResponseObject, error) {
	userID, err := uuid.Parse(req.Id.String())
	if err != nil {
		return GetOneUser404ApplicationProblemPlusJSONResponse{
			NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
				Title:  "Invalid user ID",
				Status: 400,
				Type:   "invalid-user-id",
				Detail: "The provided user ID is not a valid UUID",
			},
		}, nil
	}

	user, err := h.service.Get(ctx, userID)
	if err != nil {
		switch err {
		case ErrUserNotFound:
			return GetOneUser404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Title:  "User not found",
					Status: 404,
					Type:   "user-not-found",
					Detail: "The requested user could not be found",
				},
			}, nil
		default:
			h.logger.Error("failed to get user", "error", err)
			return GetOneUser404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Title:  "Internal server error",
					Status: 500,
					Type:   "internal-error",
					Detail: "An internal error occurred while retrieving the user",
				},
			}, nil
		}
	}

	return GetOneUser200JSONResponse{
		Data: *user,
	}, nil
}

// UpdateUser handles updating a user
func (h *Handler) UpdateUser(ctx context.Context, req UpdateUserRequestObject) (UpdateUserResponseObject, error) {
	userID, err := uuid.Parse(req.Id.String())
	if err != nil {
		return UpdateUser404ApplicationProblemPlusJSONResponse{
			NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
				Title:  "Invalid user ID",
				Status: 400,
				Type:   "invalid-user-id",
				Detail: "The provided user ID is not a valid UUID",
			},
		}, nil
	}

	updateReq := &UpdateUserJSONBody{}
	if req.Body != nil {
		updateReq.Email = req.Body.Email
		updateReq.Image = req.Body.Image
	}

	user, err := h.service.Update(ctx, userID, updateReq)
	if err != nil {
		switch err {
		case ErrUserNotFound:
			return UpdateUser404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Title:  "User not found",
					Status: 404,
					Type:   "user-not-found",
					Detail: "The requested user could not be found",
				},
			}, nil
		default:
			h.logger.Error("failed to update user", "error", err)
			return UpdateUser404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Title:  "Internal server error",
					Status: 500,
					Type:   "internal-error",
					Detail: "An internal error occurred while updating the user",
				},
			}, nil
		}
	}

	return UpdateUser200JSONResponse{
		Data: *user,
	}, nil
}

// DeleteUser handles deleting a user
func (h *Handler) DeleteUser(ctx context.Context, req DeleteUserRequestObject) (DeleteUserResponseObject, error) {
	userID, err := uuid.Parse(req.Id.String())
	if err != nil {
		return DeleteUser404ApplicationProblemPlusJSONResponse{
			NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
				Title:  "Invalid user ID",
				Status: 400,
				Type:   "invalid-user-id",
				Detail: "The provided user ID is not a valid UUID",
			},
		}, nil
	}

	// Get the user first for the response
	user, err := h.service.Get(ctx, userID)
	if err != nil {
		switch err {
		case ErrUserNotFound:
			return DeleteUser404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Title:  "User not found",
					Status: 404,
					Type:   "user-not-found",
					Detail: "The requested user could not be found",
				},
			}, nil
		default:
			h.logger.Error("failed to get user for deletion", "error", err)
			return DeleteUser404ApplicationProblemPlusJSONResponse{
				NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
					Title:  "Internal server error",
					Status: 500,
					Type:   "internal-error",
					Detail: "An internal error occurred while retrieving the user",
				},
			}, nil
		}
	}

	err = h.service.Delete(ctx, userID)
	if err != nil {
		h.logger.Error("failed to delete user", "error", err)
		return DeleteUser404ApplicationProblemPlusJSONResponse{
			NotFoundApplicationProblemPlusJSONResponse: NotFoundApplicationProblemPlusJSONResponse{
				Title:  "Internal server error",
				Status: 500,
				Type:   "internal-error",
				Detail: "An internal error occurred while deleting the user",
			},
		}, nil
	}

	return DeleteUser200JSONResponse{
		Data: *user,
	}, nil
}

// FindManyUsers handles listing multiple users
func (h *Handler) FindManyUsers(ctx context.Context, req FindManyUsersRequestObject) (FindManyUsersResponseObject, error) {
	limit := int32(50) // Default limit
	offset := int32(0) // Default offset

	// Extract pagination from request if provided
	if req.Params.Page.Size > 0 {
		limit = int32(req.Params.Page.Size)
	}
	if req.Params.Page.Number > 0 {
		offset = int32((req.Params.Page.Number - 1) * int(limit))
	}

	users, err := h.service.List(ctx, limit, offset)
	if err != nil {
		h.logger.Error("failed to list users", "error", err)
		return FindManyUsers400ApplicationProblemPlusJSONResponse{
			BadRequestApplicationProblemPlusJSONResponse: BadRequestApplicationProblemPlusJSONResponse{
				Title:  "Internal server error",
				Status: 500,
				Type:   "internal-error",
				Detail: "An internal error occurred while retrieving users",
			},
		}, nil
	}

	// Convert to response format
	responseUsers := make([]User, len(users))
	for i, user := range users {
		responseUsers[i] = *user
	}

	response := FindManyUsers200JSONResponse{
		Data: responseUsers,
	}
	response.Meta.Total = float32(len(responseUsers)) // In a real implementation, this would be the total count

	return response, nil
}
