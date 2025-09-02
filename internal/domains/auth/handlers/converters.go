package handlers

import (
	"github.com/archesai/archesai/internal/domains/auth/entities"
	"github.com/archesai/archesai/internal/generated/api/auth/users"
	"github.com/archesai/archesai/internal/generated/api/common"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// convertToGeneratedUser converts a domain User to the generated UserEntity type
func convertToGeneratedUser(u *entities.User) users.UserEntity {
	user := users.UserEntity{
		Id:            openapi_types.UUID(u.ID),
		Email:         u.Email,
		Name:          u.Name,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
		EmailVerified: u.EmailVerified,
	}

	// Add optional fields if they exist
	if u.Image != nil {
		user.Image = u.Image
	}

	return user
}

// convertToGeneratedUsers converts a slice of domain Users to generated UserEntity
func convertToGeneratedUsers(domainUsers []*entities.User) []users.UserEntity {
	result := make([]users.UserEntity, len(domainUsers))
	for i, u := range domainUsers {
		result[i] = convertToGeneratedUser(u)
	}
	return result
}

// convertUpdateRequest converts generated update request to domain input
func convertUpdateRequest(req *users.UpdateUserJSONRequestBody) map[string]interface{} {
	updates := make(map[string]interface{})

	if req.Email != nil {
		updates["email"] = *req.Email
	}

	// The UpdateUserJSONBody only has Email and Image fields based on the generated code
	if req.Image != nil {
		updates["image"] = *req.Image
	}

	return updates
}

// convertPagination converts generated pagination params to domain options
func convertPagination(page *common.Page) (limit, offset int32) {
	limit = 50 // default
	offset = 0 // default

	if page != nil {
		if page.Size != nil {
			limit = int32(*page.Size)
			if limit > 100 {
				limit = 100 // max limit
			}
		}
		if page.Number != nil && *page.Number > 0 {
			offset = int32(*page.Number-1) * limit
		}
	}

	return limit, offset
}

// convertFilter converts generated filter to domain filter
func convertFilter(filter *users.UsersFilter) map[string]interface{} {
	if filter == nil {
		return nil
	}

	// For now, return empty filter - implement actual filter logic based on your needs
	// The UsersFilterNode is a union type that needs special handling
	result := make(map[string]interface{})

	// You would need to implement the actual filter parsing based on the union type
	// filter.AsUsersFilterNode0() or filter.AsUsersFilterNode1()

	return result
}

// convertSort converts generated sort to domain sort
func convertSort(sortList *users.UsersSort) (string, string) {
	if sortList == nil || len(*sortList) == 0 {
		return "created_at", "desc"
	}

	// Take the first sort parameter
	sort := (*sortList)[0]

	orderBy := "created_at"
	orderDir := "desc"

	// Map the field
	switch sort.Field {
	case "email":
		orderBy = "email"
	case "name":
		orderBy = "name"
	case "createdAt":
		orderBy = "created_at"
	case "updatedAt":
		orderBy = "updated_at"
	}

	// Map the order
	if sort.Order == "asc" {
		orderDir = "asc"
	} else {
		orderDir = "desc"
	}

	return orderBy, orderDir
}

// Helper to convert UUID types
func toUUID(id openapi_types.UUID) uuid.UUID {
	return uuid.UUID(id)
}

func toOpenAPIUUID(id uuid.UUID) openapi_types.UUID {
	return openapi_types.UUID(id)
}
