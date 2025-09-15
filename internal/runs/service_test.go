package runs

import (
	"context"
	"errors"
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

// TestService_CreateRun tests creating a run
func TestService_CreateRun(t *testing.T) {
	pipelineID := uuid.New().String()
	orgID := uuid.New().String()

	t.Run("successful creation", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		req := &CreateRunJSONRequestBody{
			PipelineId: pipelineID,
		}

		expectedRun := &Run{
			Id:             uuid.New(),
			PipelineId:     pipelineID,
			OrganizationId: orgID,
			Status:         QUEUED,
			Progress:       0,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		mockRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(r *Run) bool {
			return r.PipelineId == pipelineID && r.Status == QUEUED
		})).Return(expectedRun, nil)

		result, err := service.Create(context.Background(), req, orgID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, pipelineID, result.PipelineId)
		assert.Equal(t, QUEUED, result.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("creation with invalid input", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		req := &CreateRunJSONRequestBody{
			PipelineId: "", // Invalid empty pipeline ID
		}

		// Mock the repository to return an error for invalid input
		mockRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*runs.Run")).Return(nil, errors.New("invalid pipeline ID"))

		result, err := service.Create(context.Background(), req, orgID)

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

// TestService_GetRun tests getting a run
func TestService_GetRun(t *testing.T) {
	runID := uuid.New()

	t.Run("successful get", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		expectedRun := &Run{
			Id:             runID,
			PipelineId:     uuid.New().String(),
			OrganizationId: uuid.New().String(),
			Status:         PROCESSING,
			Progress:       50,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		mockRepo.EXPECT().Get(mock.Anything, runID).Return(expectedRun, nil)

		result, err := service.Get(context.Background(), runID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, runID, result.Id)
		assert.Equal(t, PROCESSING, result.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("run not found", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().Get(mock.Anything, runID).Return(nil, ErrRunNotFound)

		result, err := service.Get(context.Background(), runID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrRunNotFound)
		mockRepo.AssertExpectations(t)
	})
}

// TestService_ListRuns tests listing runs
func TestService_ListRuns(t *testing.T) {
	orgID := uuid.New().String()

	t.Run("successful list", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		runs := []*Run{
			{
				Id:             uuid.New(),
				PipelineId:     uuid.New().String(),
				OrganizationId: orgID,
				Status:         COMPLETED,
				Progress:       100,
				CreatedAt:      time.Now().Add(-2 * time.Hour),
				UpdatedAt:      time.Now().Add(-1 * time.Hour),
			},
			{
				Id:             uuid.New(),
				PipelineId:     uuid.New().String(),
				OrganizationId: orgID,
				Status:         PROCESSING,
				Progress:       50,
				CreatedAt:      time.Now().Add(-30 * time.Minute),
				UpdatedAt:      time.Now(),
			},
		}

		params := ListRunsParams{
			Limit:  10,
			Offset: 0,
		}

		mockRepo.EXPECT().List(mock.Anything, params).Return(runs, int64(2), nil)

		result, total, err := service.List(context.Background(), orgID, 10, 0)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, int64(2), total)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		params := ListRunsParams{
			Limit:  10,
			Offset: 0,
		}

		mockRepo.EXPECT().List(mock.Anything, params).Return([]*Run{}, int64(0), nil)

		result, total, err := service.List(context.Background(), orgID, 10, 0)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 0)
		assert.Equal(t, int64(0), total)
		mockRepo.AssertExpectations(t)
	})
}

// TestService_DeleteRun tests deleting a run
func TestService_DeleteRun(t *testing.T) {
	runID := uuid.New()

	t.Run("successful delete", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().Delete(mock.Anything, runID).Return(nil)

		err := service.Delete(context.Background(), runID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("delete non-existent run", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().Delete(mock.Anything, runID).Return(ErrRunNotFound)

		err := service.Delete(context.Background(), runID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrRunNotFound)
		mockRepo.AssertExpectations(t)
	})
}
