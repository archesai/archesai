package database

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines basic Create, Read, Update, Delete operations for an entity of type T.
type Repository[T any] interface {
	Create(ctx context.Context, entity *T) (*T, error)
	Get(ctx context.Context, id uuid.UUID) (*T, error)
	Update(ctx context.Context, id uuid.UUID, entity *T) (*T, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int32) ([]*T, int64, error)
}
