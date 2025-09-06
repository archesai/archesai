package workflows

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
)

// MockRepository implements Repository for testing
type MockRepository struct {
	pipelines map[uuid.UUID]*Pipeline
	runs      map[uuid.UUID]*Run
	err       error
}

// Compile-time check
var _ Repository = (*MockRepository)(nil)

func NewMockRepository() *MockRepository {
	return &MockRepository{
		pipelines: make(map[uuid.UUID]*Pipeline),
		runs:      make(map[uuid.UUID]*Run),
	}
}

// Pipeline methods
func (m *MockRepository) CreatePipeline(_ context.Context, pipeline *Pipeline) (*Pipeline, error) {
	if m.err != nil {
		return nil, m.err
	}
	if pipeline.Id == (UUID{}) {
		pipeline.Id = uuid.New()
	}
	pipeline.CreatedAt = time.Now()
	pipeline.UpdatedAt = time.Now()
	m.pipelines[pipeline.Id] = pipeline
	return pipeline, nil
}

func (m *MockRepository) GetPipeline(_ context.Context, id uuid.UUID) (*Pipeline, error) {
	if m.err != nil {
		return nil, m.err
	}
	pipeline, exists := m.pipelines[id]
	if !exists {
		return nil, ErrPipelineNotFound
	}
	return pipeline, nil
}

func (m *MockRepository) UpdatePipeline(_ context.Context, id uuid.UUID, pipeline *Pipeline) (*Pipeline, error) {
	if m.err != nil {
		return nil, m.err
	}
	if _, exists := m.pipelines[id]; !exists {
		return nil, ErrPipelineNotFound
	}
	pipeline.UpdatedAt = time.Now()
	m.pipelines[id] = pipeline
	return pipeline, nil
}

func (m *MockRepository) DeletePipeline(_ context.Context, id uuid.UUID) error {
	if m.err != nil {
		return m.err
	}
	if _, exists := m.pipelines[id]; !exists {
		return ErrPipelineNotFound
	}
	delete(m.pipelines, id)
	return nil
}

func (m *MockRepository) ListPipelines(_ context.Context, _ ListPipelinesParams) ([]*Pipeline, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	pipelines := make([]*Pipeline, 0, len(m.pipelines))
	for _, pipeline := range m.pipelines {
		pipelines = append(pipelines, pipeline)
	}
	return pipelines, int64(len(pipelines)), nil
}

// Run methods
func (m *MockRepository) CreateRun(_ context.Context, run *Run) (*Run, error) {
	if m.err != nil {
		return nil, m.err
	}
	if run.Id == (UUID{}) {
		run.Id = uuid.New()
	}
	run.CreatedAt = time.Now()
	run.UpdatedAt = time.Now()
	m.runs[run.Id] = run
	return run, nil
}

func (m *MockRepository) GetRun(_ context.Context, id uuid.UUID) (*Run, error) {
	if m.err != nil {
		return nil, m.err
	}
	run, exists := m.runs[id]
	if !exists {
		return nil, ErrRunNotFound
	}
	return run, nil
}

func (m *MockRepository) UpdateRun(_ context.Context, id uuid.UUID, run *Run) (*Run, error) {
	if m.err != nil {
		return nil, m.err
	}
	if _, exists := m.runs[id]; !exists {
		return nil, ErrRunNotFound
	}
	run.UpdatedAt = time.Now()
	m.runs[id] = run
	return run, nil
}

func (m *MockRepository) DeleteRun(_ context.Context, id uuid.UUID) error {
	if m.err != nil {
		return m.err
	}
	if _, exists := m.runs[id]; !exists {
		return ErrRunNotFound
	}
	delete(m.runs, id)
	return nil
}

func (m *MockRepository) ListRuns(_ context.Context, _ ListRunsParams) ([]*Run, int64, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	runs := make([]*Run, 0, len(m.runs))
	for _, run := range m.runs {
		runs = append(runs, run)
	}
	return runs, int64(len(runs)), nil
}

// Test cases
func TestService_CreatePipeline(t *testing.T) {
	tests := []struct {
		name    string
		req     *CreatePipelineRequest
		orgID   string
		repoErr error
		wantErr bool
	}{
		{
			name: "successful creation",
			req: &CreatePipelineRequest{
				Name:        "Test Pipeline",
				Description: "Test Description",
			},
			orgID:   uuid.New().String(),
			wantErr: false,
		},
		{
			name: "repository error",
			req: &CreatePipelineRequest{
				Name:        "Test Pipeline",
				Description: "Test Description",
			},
			orgID:   uuid.New().String(),
			repoErr: errors.New("database error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			repo.err = tt.repoErr
			service := NewService(repo, slog.Default())

			pipeline, err := service.CreatePipeline(context.Background(), tt.req, tt.orgID)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreatePipeline() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && pipeline == nil {
				t.Error("CreatePipeline() returned nil pipeline")
			}
		})
	}
}

func TestService_GetPipeline(t *testing.T) {
	tests := []struct {
		name       string
		pipelineID uuid.UUID
		setup      func(*MockRepository)
		wantErr    bool
	}{
		{
			name:       "existing pipeline",
			pipelineID: uuid.New(),
			setup: func(r *MockRepository) {
				pipeline := &Pipeline{
					Id:          uuid.New(),
					Name:        "Test Pipeline",
					Description: "Test Description",
				}
				r.pipelines[pipeline.Id] = pipeline
			},
			wantErr: false,
		},
		{
			name:       "non-existent pipeline",
			pipelineID: uuid.New(),
			setup:      func(_ *MockRepository) {},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, slog.Default())

			// Get the correct ID from the setup
			var testID uuid.UUID
			if tt.name == "existing pipeline" {
				for id := range repo.pipelines {
					testID = id
					break
				}
			} else {
				testID = tt.pipelineID
			}

			pipeline, err := service.GetPipeline(context.Background(), testID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetPipeline() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && pipeline == nil {
				t.Error("GetPipeline() returned nil pipeline")
			}
		})
	}
}

func TestService_UpdatePipeline(t *testing.T) {
	existingID := uuid.New()

	tests := []struct {
		name       string
		pipelineID uuid.UUID
		req        *UpdatePipelineRequest
		setup      func(*MockRepository)
		wantErr    bool
	}{
		{
			name:       "successful update",
			pipelineID: existingID,
			req: &UpdatePipelineRequest{
				Name:        "Updated Pipeline",
				Description: "Updated Description",
			},
			setup: func(r *MockRepository) {
				pipeline := &Pipeline{
					Id:          existingID,
					Name:        "Original Pipeline",
					Description: "Original Description",
				}
				r.pipelines[existingID] = pipeline
			},
			wantErr: false,
		},
		{
			name:       "non-existent pipeline",
			pipelineID: uuid.New(),
			req: &UpdatePipelineRequest{
				Name: "Updated Pipeline",
			},
			setup:   func(_ *MockRepository) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, slog.Default())

			pipeline, err := service.UpdatePipeline(context.Background(), tt.pipelineID, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdatePipeline() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if pipeline == nil {
					t.Error("UpdatePipeline() returned nil pipeline")
				} else {
					if tt.req.Name != "" && pipeline.Name != tt.req.Name {
						t.Errorf("UpdatePipeline() name = %v, want %v", pipeline.Name, tt.req.Name)
					}
					if tt.req.Description != "" && pipeline.Description != tt.req.Description {
						t.Errorf("UpdatePipeline() description = %v, want %v", pipeline.Description, tt.req.Description)
					}
				}
			}
		})
	}
}

func TestService_DeletePipeline(t *testing.T) {
	existingID := uuid.New()

	tests := []struct {
		name       string
		pipelineID uuid.UUID
		setup      func(*MockRepository)
		wantErr    bool
	}{
		{
			name:       "successful deletion",
			pipelineID: existingID,
			setup: func(r *MockRepository) {
				pipeline := &Pipeline{
					Id:   existingID,
					Name: "Test Pipeline",
				}
				r.pipelines[existingID] = pipeline
			},
			wantErr: false,
		},
		{
			name:       "non-existent pipeline",
			pipelineID: uuid.New(),
			setup:      func(_ *MockRepository) {},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, slog.Default())

			err := service.DeletePipeline(context.Background(), tt.pipelineID)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeletePipeline() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if _, exists := repo.pipelines[tt.pipelineID]; exists {
					t.Error("DeletePipeline() pipeline still exists")
				}
			}
		})
	}
}

func TestService_CreateRun(t *testing.T) {
	pipelineID := uuid.New()

	tests := []struct {
		name    string
		req     *CreateRunRequest
		orgID   string
		setup   func(*MockRepository)
		repoErr error
		wantErr bool
	}{
		{
			name: "successful creation",
			req: &CreateRunRequest{
				PipelineId: pipelineID.String(),
			},
			orgID: uuid.New().String(),
			setup: func(r *MockRepository) {
				pipeline := &Pipeline{
					Id:          pipelineID,
					Name:        "Test Pipeline",
					Description: "Test Description",
				}
				r.pipelines[pipelineID] = pipeline
			},
			wantErr: false,
		},
		{
			name: "pipeline not found",
			req: &CreateRunRequest{
				PipelineId: uuid.New().String(),
			},
			orgID:   uuid.New().String(),
			setup:   func(_ *MockRepository) {},
			wantErr: true,
		},
		{
			name: "invalid pipeline ID",
			req: &CreateRunRequest{
				PipelineId: "invalid-uuid",
			},
			orgID:   uuid.New().String(),
			setup:   func(_ *MockRepository) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			repo.err = tt.repoErr
			service := NewService(repo, slog.Default())

			run, err := service.CreateRun(context.Background(), tt.req, tt.orgID)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRun() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && run == nil {
				t.Error("CreateRun() returned nil run")
			}
		})
	}
}

func TestService_GetRun(t *testing.T) {
	tests := []struct {
		name    string
		runID   uuid.UUID
		setup   func(*MockRepository)
		wantErr bool
	}{
		{
			name:  "existing run",
			runID: uuid.New(),
			setup: func(r *MockRepository) {
				run := &Run{
					Id:             uuid.New(),
					PipelineId:     uuid.New().String(),
					OrganizationId: uuid.New().String(),
					Status:         QUEUED,
				}
				r.runs[run.Id] = run
			},
			wantErr: false,
		},
		{
			name:    "non-existent run",
			runID:   uuid.New(),
			setup:   func(_ *MockRepository) {},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, slog.Default())

			// Get the correct ID from the setup
			var testID uuid.UUID
			if tt.name == "existing run" {
				for id := range repo.runs {
					testID = id
					break
				}
			} else {
				testID = tt.runID
			}

			run, err := service.GetRun(context.Background(), testID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetRun() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && run == nil {
				t.Error("GetRun() returned nil run")
			}
		})
	}
}

func TestService_ListPipelines(t *testing.T) {
	tests := []struct {
		name    string
		orgID   string
		limit   int
		offset  int
		setup   func(*MockRepository)
		wantErr bool
	}{
		{
			name:   "list with pipelines",
			orgID:  uuid.New().String(),
			limit:  10,
			offset: 0,
			setup: func(r *MockRepository) {
				for i := 0; i < 5; i++ {
					pipeline := &Pipeline{
						Id:          uuid.New(),
						Name:        "Test Pipeline",
						Description: "Test Description",
					}
					r.pipelines[pipeline.Id] = pipeline
				}
			},
			wantErr: false,
		},
		{
			name:    "empty list",
			orgID:   uuid.New().String(),
			limit:   10,
			offset:  0,
			setup:   func(_ *MockRepository) {},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, slog.Default())

			pipelines, total, err := service.ListPipelines(context.Background(), tt.orgID, tt.limit, tt.offset)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListPipelines() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if pipelines == nil {
					t.Error("ListPipelines() returned nil pipelines")
				}
				if total != len(pipelines) {
					t.Errorf("ListPipelines() total = %v, want %v", total, len(pipelines))
				}
			}
		})
	}
}

// TestService_StartRun tests the StartRun method
func TestService_StartRun(t *testing.T) {
	existingID := uuid.New()

	tests := []struct {
		name    string
		runID   uuid.UUID
		setup   func(*MockRepository)
		wantErr bool
	}{
		{
			name:  "successful start",
			runID: existingID,
			setup: func(r *MockRepository) {
				run := &Run{
					Id:             existingID,
					PipelineId:     uuid.New().String(),
					OrganizationId: uuid.New().String(),
					Status:         QUEUED,
				}
				r.runs[existingID] = run
			},
			wantErr: false,
		},
		{
			name:    "non-existent run",
			runID:   uuid.New(),
			setup:   func(_ *MockRepository) {},
			wantErr: true,
		},
		{
			name:  "repository error",
			runID: existingID,
			setup: func(r *MockRepository) {
				r.err = errors.New("database error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, slog.Default())

			run, err := service.StartRun(context.Background(), tt.runID)

			if (err != nil) != tt.wantErr {
				t.Errorf("StartRun() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && run == nil {
				t.Error("StartRun() returned nil run")
			}
		})
	}
}

// TestService_DeleteRun tests the DeleteRun method
func TestService_DeleteRun(t *testing.T) {
	existingID := uuid.New()

	tests := []struct {
		name    string
		runID   uuid.UUID
		setup   func(*MockRepository)
		wantErr bool
	}{
		{
			name:  "delete existing run",
			runID: existingID,
			setup: func(r *MockRepository) {
				run := &Run{
					Id:             existingID,
					PipelineId:     uuid.New().String(),
					OrganizationId: uuid.New().String(),
					Status:         COMPLETED,
				}
				r.runs[existingID] = run
			},
			wantErr: false,
		},
		{
			name:    "delete non-existent run",
			runID:   uuid.New(),
			setup:   func(_ *MockRepository) {},
			wantErr: true,
		},
		{
			name:  "repository error",
			runID: existingID,
			setup: func(r *MockRepository) {
				r.err = errors.New("database error")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, slog.Default())

			err := service.DeleteRun(context.Background(), tt.runID)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteRun() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestService_ListRuns tests the ListRuns method
func TestService_ListRuns(t *testing.T) {
	pipelineID := uuid.New().String()

	tests := []struct {
		name       string
		pipelineID string
		limit      int
		offset     int
		setup      func(*MockRepository)
		wantErr    bool
		wantLen    int
	}{
		{
			name:       "list runs for pipeline",
			pipelineID: pipelineID,
			limit:      10,
			offset:     0,
			setup: func(r *MockRepository) {
				for i := 0; i < 3; i++ {
					run := &Run{
						Id:             uuid.New(),
						PipelineId:     pipelineID,
						OrganizationId: uuid.New().String(),
						Status:         COMPLETED,
					}
					r.runs[run.Id] = run
				}
			},
			wantErr: false,
			wantLen: 3,
		},
		{
			name:       "empty run list",
			pipelineID: uuid.New().String(),
			limit:      10,
			offset:     0,
			setup:      func(_ *MockRepository) {},
			wantErr:    false,
			wantLen:    0,
		},
		{
			name:       "with pagination",
			pipelineID: pipelineID,
			limit:      2,
			offset:     1,
			setup: func(r *MockRepository) {
				for i := 0; i < 5; i++ {
					run := &Run{
						Id:             uuid.New(),
						PipelineId:     pipelineID,
						OrganizationId: uuid.New().String(),
						Status:         PROCESSING,
					}
					r.runs[run.Id] = run
				}
			},
			wantErr: false,
			wantLen: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockRepository()
			tt.setup(repo)
			service := NewService(repo, slog.Default())

			runs, total, err := service.ListRuns(context.Background(), tt.pipelineID, tt.limit, tt.offset)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListRuns() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if runs == nil {
					t.Error("ListRuns() returned nil runs")
				}
				if total != tt.wantLen {
					t.Errorf("ListRuns() total = %v, want %v", total, tt.wantLen)
				}
			}
		})
	}
}

// TestNewService tests the service constructor
func TestNewService(t *testing.T) {
	repo := NewMockRepository()
	logger := slog.Default()

	service := NewService(repo, logger)

	if service == nil {
		t.Fatal("NewService() returned nil")
	}

	if service.repo == nil {
		t.Error("NewService() service.repo is nil")
	}

	if service.logger == nil {
		t.Error("NewService() service.logger is nil")
	}
}

// TestMockRepository_EdgeCases tests edge cases in the mock repository
func TestMockRepository_EdgeCases(t *testing.T) {
	t.Run("ListPipelines with error", func(t *testing.T) {
		repo := NewMockRepository()
		repo.err = errors.New("database error")

		pipelines, total, err := repo.ListPipelines(context.Background(), ListPipelinesParams{})

		if err == nil {
			t.Error("ListPipelines() expected error but got none")
		}

		if pipelines != nil {
			t.Error("ListPipelines() returned non-nil slice on error")
		}

		if total != 0 {
			t.Errorf("ListPipelines() total = %v, want 0", total)
		}
	})

	t.Run("UpdateRun on non-existent", func(t *testing.T) {
		repo := NewMockRepository()
		run := &Run{
			Status: FAILED,
		}

		_, err := repo.UpdateRun(context.Background(), uuid.New(), run)

		if !errors.Is(err, ErrRunNotFound) {
			t.Errorf("UpdateRun() error = %v, want %v", err, ErrRunNotFound)
		}
	})

	t.Run("GetRun on non-existent", func(t *testing.T) {
		repo := NewMockRepository()

		_, err := repo.GetRun(context.Background(), uuid.New())

		if !errors.Is(err, ErrRunNotFound) {
			t.Errorf("GetRun() error = %v, want %v", err, ErrRunNotFound)
		}
	})

	t.Run("DeletePipeline with runs", func(t *testing.T) {
		repo := NewMockRepository()
		pipelineID := uuid.New()

		// Create pipeline
		pipeline := &Pipeline{
			Id:   pipelineID,
			Name: "Test Pipeline",
		}
		repo.pipelines[pipelineID] = pipeline

		// Create run for this pipeline
		run := &Run{
			Id:         uuid.New(),
			PipelineId: pipelineID.String(),
		}
		repo.runs[run.Id] = run

		// Delete pipeline
		err := repo.DeletePipeline(context.Background(), pipelineID)

		if err != nil {
			t.Errorf("DeletePipeline() unexpected error: %v", err)
		}

		// Verify pipeline was deleted
		if _, exists := repo.pipelines[pipelineID]; exists {
			t.Error("Pipeline was not deleted")
		}

		// Note: runs are not automatically deleted in mock
		if _, exists := repo.runs[run.Id]; !exists {
			t.Error("Run was unexpectedly deleted")
		}
	})
}

// TestRunStatusTransitions tests run status transitions
func TestRunStatusTransitions(t *testing.T) {
	tests := []struct {
		name          string
		currentStatus RunStatus
		newStatus     RunStatus
		valid         bool
	}{
		{
			name:          "queued to running",
			currentStatus: QUEUED,
			newStatus:     PROCESSING,
			valid:         true,
		},
		{
			name:          "running to completed",
			currentStatus: PROCESSING,
			newStatus:     COMPLETED,
			valid:         true,
		},
		{
			name:          "running to failed",
			currentStatus: PROCESSING,
			newStatus:     FAILED,
			valid:         true,
		},
		{
			name:          "completed to running",
			currentStatus: COMPLETED,
			newStatus:     PROCESSING,
			valid:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This is a placeholder for status transition logic
			// In real implementation, you would validate transitions
			if tt.currentStatus == COMPLETED && tt.newStatus == PROCESSING {
				if tt.valid {
					t.Error("Invalid status transition marked as valid")
				}
			}
		})
	}
}

// BenchmarkCreatePipeline benchmarks pipeline creation
func BenchmarkCreatePipeline(b *testing.B) {
	repo := NewMockRepository()
	service := NewService(repo, slog.Default())

	req := &CreatePipelineRequest{
		Name:        "Benchmark Pipeline",
		Description: "Benchmark Description",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.CreatePipeline(context.Background(), req, uuid.New().String())
	}
}

// BenchmarkListRuns benchmarks listing runs
func BenchmarkListRuns(b *testing.B) {
	repo := NewMockRepository()
	service := NewService(repo, slog.Default())

	// Setup some runs
	pipelineID := uuid.New().String()
	for i := 0; i < 100; i++ {
		run := &Run{
			Id:         uuid.New(),
			PipelineId: pipelineID,
			Status:     COMPLETED,
		}
		repo.runs[run.Id] = run
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = service.ListRuns(context.Background(), pipelineID, 10, 0)
	}
}
