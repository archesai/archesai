package runs

import (
	"context"

	"github.com/archesai/archesai/internal/core/aggregates"
	"github.com/archesai/archesai/internal/core/valueobjects"
)

// GetRunQuery represents a query to get a run by ID.
type GetRunQuery struct {
	RunID valueobjects.RunID
}

// NewGetRunQuery creates a new get run query.
func NewGetRunQuery(runID valueobjects.RunID) *GetRunQuery {
	return &GetRunQuery{
		RunID: runID,
	}
}

// GetRunQueryHandler handles the get run query.
type GetRunQueryHandler interface {
	Handle(ctx context.Context, query *GetRunQuery) (*aggregates.Run, error)
}
