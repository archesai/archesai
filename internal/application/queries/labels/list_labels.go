package labels

import (
	"context"

	"github.com/archesai/archesai/internal/core/entities"
	"github.com/archesai/archesai/internal/core/valueobjects"
)

// ListLabelsQuery represents a query to list labels.
type ListLabelsQuery struct {
	OrganizationID valueobjects.OrganizationID
	SearchQuery    *string
	Limit          int
	Offset         int
}

// NewListLabelsQuery creates a new list labels query.
func NewListLabelsQuery(organizationID valueobjects.OrganizationID) *ListLabelsQuery {
	return &ListLabelsQuery{
		OrganizationID: organizationID,
		Limit:          50,
		Offset:         0,
	}
}

// WithSearch sets the search query filter.
func (q *ListLabelsQuery) WithSearch(search string) *ListLabelsQuery {
	q.SearchQuery = &search
	return q
}

// WithPagination sets pagination parameters.
func (q *ListLabelsQuery) WithPagination(limit, offset int) *ListLabelsQuery {
	q.Limit = limit
	q.Offset = offset
	return q
}

// ListLabelsQueryHandler handles the list labels query.
type ListLabelsQueryHandler interface {
	Handle(ctx context.Context, query *ListLabelsQuery) ([]*entities.Label, error)
}
