package workflows

import (
	"context"
	"testing"
	"time"

	"github.com/archesai/archesai/internal/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test helper functions
func createTestService(t *testing.T) (*Service, *MockRepository) {
	t.Helper()

	mockRepo := NewMockRepository(t)
	logger := logger.NewTest()

	service := NewService(mockRepo, logger)
	return service, mockRepo
}

// TestService_CreatePipeline tests creating a pipeline
func TestService_CreatePipeline(t *testing.T) {
	orgID := uuid.New().String()

	t.Run("successful creation", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		req := &CreatePipelineRequest{
			Name:        "Test Pipeline",
			Description: "Test description",
		}

		expectedPipeline := &Pipeline{
			Id:             uuid.New(),
			OrganizationId: uuid.New(),
			Name:           "Test Pipeline",
			Description:    "Test description",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		mockRepo.EXPECT().CreatePipeline(mock.Anything, mock.MatchedBy(func(p *Pipeline) bool {
			return p.Name == "Test Pipeline"
		})).Return(expectedPipeline, nil)

		result, err := service.CreatePipeline(context.Background(), req, orgID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Test Pipeline", result.Name)
	})
}

// TestService_GetPipeline tests getting a pipeline
func TestService_GetPipeline(t *testing.T) {
	pipelineID := uuid.New()
	pipeline := &Pipeline{
		Id:             pipelineID,
		OrganizationId: uuid.New(),
		Name:           "Test Pipeline",
		Description:    "Test description",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	t.Run("successful get", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().GetPipeline(mock.Anything, pipelineID).Return(pipeline, nil)

		result, err := service.GetPipeline(context.Background(), pipelineID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, pipelineID, result.Id)
	})

	t.Run("not found", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().GetPipeline(mock.Anything, pipelineID).Return(nil, ErrPipelineNotFound)

		result, err := service.GetPipeline(context.Background(), pipelineID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrPipelineNotFound)
	})
}

// TestService_CreateRun tests creating a run
func TestService_CreateRun(t *testing.T) {
	pipelineID := uuid.New().String()

	t.Run("successful creation", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		pipeline := &Pipeline{
			Id:             uuid.New(),
			OrganizationId: uuid.New(),
			Name:           "Test Pipeline",
			Description:    "Test description",
		}

		req := &CreateRunRequest{
			PipelineId: pipelineID,
		}

		expectedRun := &Run{
			Id:             uuid.New(),
			PipelineId:     pipelineID,
			OrganizationId: uuid.New().String(),
			Status:         QUEUED,
			Progress:       0,
			StartedAt:      time.Now(),
			CreatedAt:      time.Now(),
		}

		// First get the pipeline to validate it exists
		mockRepo.EXPECT().GetPipeline(mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(pipeline, nil)
		mockRepo.EXPECT().CreateRun(mock.Anything, mock.MatchedBy(func(r *Run) bool {
			return r.PipelineId == pipelineID
		})).Return(expectedRun, nil)

		orgID := uuid.New().String()
		result, err := service.CreateRun(context.Background(), req, orgID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, pipelineID, result.PipelineId)
	})

	t.Run("pipeline not found", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		req := &CreateRunRequest{
			PipelineId: pipelineID,
		}

		mockRepo.EXPECT().GetPipeline(mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, ErrPipelineNotFound)

		orgID := uuid.New().String()
		result, err := service.CreateRun(context.Background(), req, orgID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrPipelineNotFound)
	})
}

// TestService_UpdateRunProgress tests updating a run's progress
func TestService_UpdateRunProgress(t *testing.T) {
	runID := uuid.New()

	t.Run("successful update", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		// UpdateRunProgress takes progress directly
		progress := float32(100)

		expectedRun := &Run{
			Id:         runID,
			PipelineId: uuid.New().String(),
			Status:     COMPLETED,
			Progress:   100,
			StartedAt:  time.Now().Add(-1 * time.Hour),
			UpdatedAt:  time.Now(),
		}

		// First get the run to update
		existingRun := &Run{
			Id:         runID,
			PipelineId: uuid.New().String(),
			Status:     PROCESSING,
			Progress:   50,
		}
		mockRepo.EXPECT().GetRun(mock.Anything, runID).Return(existingRun, nil)
		mockRepo.EXPECT().UpdateRun(mock.Anything, runID, mock.MatchedBy(func(r *Run) bool {
			return r.Progress == 100
		})).Return(expectedRun, nil)

		result, err := service.UpdateRunProgress(context.Background(), runID, progress)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, COMPLETED, result.Status)
		assert.Equal(t, float32(100), result.Progress)
	})
}
