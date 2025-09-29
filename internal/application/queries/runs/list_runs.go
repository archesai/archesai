package runs

import (
	"context"

	"github.com/archesai/archesai/internal/core/aggregates"
	"github.com/archesai/archesai/internal/core/valueobjects"
)

// ListRunsQuery represents a query to list runs.
type ListRunsQuery struct {
	PipelineID     *valueobjects.PipelineID
	ToolID         *valueobjects.ToolID
	OrganizationID *valueobjects.OrganizationID
	Status         *valueobjects.RunStatus
	CreatedBy      *valueobjects.UserID
	Limit          int
	Offset         int
}

// NewListRunsQuery creates a new list runs query.
func NewListRunsQuery() *ListRunsQuery {
	return &ListRunsQuery{
		Limit:  50,
		Offset: 0,
	}
}

// WithPipelineID sets the pipeline ID filter.
func (q *ListRunsQuery) WithPipelineID(id valueobjects.PipelineID) *ListRunsQuery {
	q.PipelineID = &id
	return q
}

// WithToolID sets the tool ID filter.
func (q *ListRunsQuery) WithToolID(id valueobjects.ToolID) *ListRunsQuery {
	q.ToolID = &id
	return q
}

// WithOrganizationID sets the organization ID filter.
func (q *ListRunsQuery) WithOrganizationID(id valueobjects.OrganizationID) *ListRunsQuery {
	q.OrganizationID = &id
	return q
}

// WithStatus sets the status filter.
func (q *ListRunsQuery) WithStatus(status valueobjects.RunStatus) *ListRunsQuery {
	q.Status = &status
	return q
}

// WithCreatedBy sets the created by filter.
func (q *ListRunsQuery) WithCreatedBy(userID valueobjects.UserID) *ListRunsQuery {
	q.CreatedBy = &userID
	return q
}

// WithPagination sets pagination parameters.
func (q *ListRunsQuery) WithPagination(limit, offset int) *ListRunsQuery {
	q.Limit = limit
	q.Offset = offset
	return q
}

// ListRunsQueryHandler handles the list runs query.
type ListRunsQueryHandler interface {
	Handle(ctx context.Context, query *ListRunsQuery) ([]*aggregates.Run, error)
}
