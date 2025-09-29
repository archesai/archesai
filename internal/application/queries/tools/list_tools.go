package tools

import (
	"context"

	"github.com/archesai/archesai/internal/core/aggregates"
	"github.com/archesai/archesai/internal/core/valueobjects"
)

// ListToolsQuery represents a query to list tools.
type ListToolsQuery struct {
	OrganizationID *valueobjects.OrganizationID
	ToolType       *string
	ActiveOnly     bool
	Limit          int
	Offset         int
}

// NewListToolsQuery creates a new list tools query.
func NewListToolsQuery() *ListToolsQuery {
	return &ListToolsQuery{
		ActiveOnly: true,
		Limit:      50,
		Offset:     0,
	}
}

// WithOrganizationID sets the organization ID filter.
func (q *ListToolsQuery) WithOrganizationID(id valueobjects.OrganizationID) *ListToolsQuery {
	q.OrganizationID = &id
	return q
}

// WithToolType sets the tool type filter.
func (q *ListToolsQuery) WithToolType(toolType string) *ListToolsQuery {
	q.ToolType = &toolType
	return q
}

// WithPagination sets pagination parameters.
func (q *ListToolsQuery) WithPagination(limit, offset int) *ListToolsQuery {
	q.Limit = limit
	q.Offset = offset
	return q
}

// ListToolsQueryHandler handles the list tools query.
type ListToolsQueryHandler interface {
	Handle(ctx context.Context, query *ListToolsQuery) ([]*aggregates.Tool, error)
}
