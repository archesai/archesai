// Package handlers provides HTTP request handlers for the auth domain.
package handlers

import (
	"github.com/archesai/archesai/internal/domains/auth/entities"
	"github.com/archesai/archesai/internal/generated/api"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// convertToGeneratedUser converts a domain User to the generated UserEntity type
func convertToGeneratedUser(u *entities.User) api.UserEntity {
	user := api.UserEntity{
		Id:            u.ID,
		Email:         openapi_types.Email(u.Email),
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
func convertToGeneratedUsers(domainUsers []*entities.User) []api.UserEntity {
	result := make([]api.UserEntity, len(domainUsers))
	for i, u := range domainUsers {
		result[i] = convertToGeneratedUser(u)
	}
	return result
}

// convertPagination converts generated pagination params to domain options
func convertPagination(page api.Page) (limit, offset int32) {
	limit = 50 // default
	offset = 0 // default

	if page.Size > 0 {
		limit = int32(page.Size)
		if limit > 100 {
			limit = 100 // max limit
		}
	}
	if page.Number > 0 {
		offset = int32(page.Number-1) * limit
	}

	return limit, offset
}
