package runs

import (
	"context"

	"github.com/archesai/archesai/internal/core/valueobjects"
)

// GetRunStatusQuery represents a query to get the status of a run.
type GetRunStatusQuery struct {
	RunID valueobjects.RunID
}

// NewGetRunStatusQuery creates a new get run status query.
func NewGetRunStatusQuery(runID valueobjects.RunID) *GetRunStatusQuery {
	return &GetRunStatusQuery{
		RunID: runID,
	}
}

// RunStatusResponse represents the response to a run status query.
type RunStatusResponse struct {
	RunID      valueobjects.RunID
	Status     valueobjects.RunStatus
	Progress   int
	StartedAt  *string
	FinishedAt *string
	Error      *string
}

// GetRunStatusQueryHandler handles the get run status query.
type GetRunStatusQueryHandler interface {
	Handle(ctx context.Context, query *GetRunStatusQuery) (*RunStatusResponse, error)
}
