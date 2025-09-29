package labels

import (
	"context"

	"github.com/archesai/archesai/internal/core/entities"
	"github.com/archesai/archesai/internal/core/valueobjects"
)

// GetLabelQuery represents a query to get a label by ID.
type GetLabelQuery struct {
	LabelID valueobjects.LabelID
}

// NewGetLabelQuery creates a new get label query.
func NewGetLabelQuery(labelID valueobjects.LabelID) *GetLabelQuery {
	return &GetLabelQuery{
		LabelID: labelID,
	}
}

// GetLabelQueryHandler handles the get label query.
type GetLabelQueryHandler interface {
	Handle(ctx context.Context, query *GetLabelQuery) (*entities.Label, error)
}
