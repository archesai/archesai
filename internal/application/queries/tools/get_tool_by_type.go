package tools

import (
	"context"

	"github.com/archesai/archesai/internal/core/aggregates"
	"github.com/archesai/archesai/internal/core/valueobjects"
)

// GetToolByTypeQuery represents a query to get a tool by type.
type GetToolByTypeQuery struct {
	OrganizationID valueobjects.OrganizationID
	ToolType       string
}

// NewGetToolByTypeQuery creates a new get tool by type query.
func NewGetToolByTypeQuery(
	organizationID valueobjects.OrganizationID,
	toolType string,
) *GetToolByTypeQuery {
	return &GetToolByTypeQuery{
		OrganizationID: organizationID,
		ToolType:       toolType,
	}
}

// GetToolByTypeQueryHandler handles the get tool by type query.
type GetToolByTypeQueryHandler interface {
	Handle(ctx context.Context, query *GetToolByTypeQuery) (*aggregates.Tool, error)
}
