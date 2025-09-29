package tools

import (
	"context"

	"github.com/archesai/archesai/internal/core/aggregates"
	"github.com/google/uuid"
)

// GetToolQuery represents a query to get a tool by ID.
type GetToolQuery struct {
	ToolID uuid.UUID
}

// NewGetToolQuery creates a new get tool query.
func NewGetToolQuery(toolID uuid.UUID) *GetToolQuery {
	return &GetToolQuery{
		ToolID: toolID,
	}
}

// GetToolQueryHandler handles the get tool query.
type GetToolQueryHandler interface {
	Handle(ctx context.Context, query *GetToolQuery) (*aggregates.Tool, error)
}
