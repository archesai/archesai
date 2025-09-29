package labels

import (
	"context"

	"github.com/archesai/archesai/internal/core/entities"
	"github.com/archesai/archesai/internal/core/valueobjects"
)

// SearchLabelsQuery represents a query to search labels.
type SearchLabelsQuery struct {
	OrganizationID valueobjects.OrganizationID
	Query          string
	Limit          int
}

// NewSearchLabelsQuery creates a new search labels query.
func NewSearchLabelsQuery(
	organizationID valueobjects.OrganizationID,
	query string,
) *SearchLabelsQuery {
	return &SearchLabelsQuery{
		OrganizationID: organizationID,
		Query:          query,
		Limit:          10,
	}
}

// WithLimit sets the result limit.
func (q *SearchLabelsQuery) WithLimit(limit int) *SearchLabelsQuery {
	q.Limit = limit
	return q
}

// SearchLabelsQueryHandler handles the search labels query.
type SearchLabelsQueryHandler interface {
	Handle(ctx context.Context, query *SearchLabelsQuery) ([]*entities.Label, error)
}
