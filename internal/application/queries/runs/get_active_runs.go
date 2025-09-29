package runs

import (
	"context"

	"github.com/archesai/archesai/internal/core/aggregates"
	"github.com/archesai/archesai/internal/core/valueobjects"
)

// GetActiveRunsQuery represents a query to get all active runs.
type GetActiveRunsQuery struct {
	OrganizationID *valueobjects.OrganizationID
}

// NewGetActiveRunsQuery creates a new get active runs query.
func NewGetActiveRunsQuery() *GetActiveRunsQuery {
	return &GetActiveRunsQuery{}
}

// WithOrganizationID sets the organization filter.
func (q *GetActiveRunsQuery) WithOrganizationID(
	id valueobjects.OrganizationID,
) *GetActiveRunsQuery {
	q.OrganizationID = &id
	return q
}

// GetActiveRunsQueryHandler handles the get active runs query.
type GetActiveRunsQueryHandler interface {
	Handle(ctx context.Context, query *GetActiveRunsQuery) ([]*aggregates.Run, error)
}
