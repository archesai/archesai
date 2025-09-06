// Package workflows provides SQLite-based repository implementation for workflows domain
package workflows

import (
	"context"

	"github.com/archesai/archesai/internal/database/sqlite"
	"github.com/google/uuid"
)

// SQLiteRepository handles workflows data persistence using SQLite
type SQLiteRepository struct {
	q *sqlite.Queries
}

// CreatePipeline implements Repository.
func (s *SQLiteRepository) CreatePipeline(_ context.Context, _ *Pipeline) (*Pipeline, error) {
	panic("unimplemented")
}

// CreateRun implements Repository.
func (s *SQLiteRepository) CreateRun(_ context.Context, _ *Run) (*Run, error) {
	panic("unimplemented")
}

// DeletePipeline implements Repository.
func (s *SQLiteRepository) DeletePipeline(_ context.Context, _ uuid.UUID) error {
	panic("unimplemented")
}

// DeleteRun implements Repository.
func (s *SQLiteRepository) DeleteRun(_ context.Context, _ uuid.UUID) error {
	panic("unimplemented")
}

// GetPipeline implements Repository.
func (s *SQLiteRepository) GetPipeline(_ context.Context, _ uuid.UUID) (*Pipeline, error) {
	panic("unimplemented")
}

// GetRun implements Repository.
func (s *SQLiteRepository) GetRun(_ context.Context, _ uuid.UUID) (*Run, error) {
	panic("unimplemented")
}

// ListPipelines implements Repository.
func (s *SQLiteRepository) ListPipelines(_ context.Context, _ ListPipelinesParams) ([]*Pipeline, int64, error) {
	panic("unimplemented")
}

// ListRuns implements Repository.
func (s *SQLiteRepository) ListRuns(_ context.Context, _ ListRunsParams) ([]*Run, int64, error) {
	panic("unimplemented")
}

// UpdatePipeline implements Repository.
func (s *SQLiteRepository) UpdatePipeline(_ context.Context, _ uuid.UUID, _ *Pipeline) (*Pipeline, error) {
	panic("unimplemented")
}

// UpdateRun implements Repository.
func (s *SQLiteRepository) UpdateRun(_ context.Context, _ uuid.UUID, _ *Run) (*Run, error) {
	panic("unimplemented")
}

// NewSQLiteRepository creates a new SQLite repository for workflows
func NewSQLiteRepository(q *sqlite.Queries) *SQLiteRepository {
	return &SQLiteRepository{
		q: q,
	}
}

// Ensure SQLiteRepository implements Repository
var _ Repository = (*SQLiteRepository)(nil)
