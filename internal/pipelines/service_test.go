package pipelines

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

		mockRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(p *Pipeline) bool {
			return p.Name == "Test Pipeline"
		})).Return(expectedPipeline, nil)

		result, err := service.Create(context.Background(), req, orgID)

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

		mockRepo.EXPECT().Get(mock.Anything, pipelineID).Return(pipeline, nil)

		result, err := service.Get(context.Background(), pipelineID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, pipelineID, result.Id)
	})

	t.Run("not found", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().Get(mock.Anything, pipelineID).Return(nil, ErrPipelineNotFound)

		result, err := service.Get(context.Background(), pipelineID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrPipelineNotFound)
	})
}

// TestService_UpdatePipeline tests updating a pipeline
func TestService_UpdatePipeline(t *testing.T) {
	pipelineID := uuid.New()

	t.Run("successful update", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		existingPipeline := &Pipeline{
			Id:             pipelineID,
			OrganizationId: uuid.New(),
			Name:           "Old Name",
			Description:    "Old description",
			CreatedAt:      time.Now().Add(-1 * time.Hour),
			UpdatedAt:      time.Now().Add(-1 * time.Hour),
		}

		updatedPipeline := &Pipeline{
			Id:             pipelineID,
			OrganizationId: existingPipeline.OrganizationId,
			Name:           "New Name",
			Description:    "New description",
			CreatedAt:      existingPipeline.CreatedAt,
			UpdatedAt:      time.Now(),
		}

		req := &UpdatePipelineRequest{
			Name:        "New Name",
			Description: "New description",
		}

		mockRepo.EXPECT().Get(mock.Anything, pipelineID).Return(existingPipeline, nil)
		mockRepo.EXPECT().Update(mock.Anything, pipelineID, mock.MatchedBy(func(p *Pipeline) bool {
			return p.Name == "New Name" && p.Description == "New description"
		})).Return(updatedPipeline, nil)

		result, err := service.Update(context.Background(), pipelineID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "New Name", result.Name)
		assert.Equal(t, "New description", result.Description)
	})
}

// TestService_DeletePipeline tests deleting a pipeline
func TestService_DeletePipeline(t *testing.T) {
	pipelineID := uuid.New()

	t.Run("successful delete", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().Delete(mock.Anything, pipelineID).Return(nil)

		err := service.Delete(context.Background(), pipelineID)

		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().Delete(mock.Anything, pipelineID).Return(ErrPipelineNotFound)

		err := service.Delete(context.Background(), pipelineID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrPipelineNotFound)
	})
}

// TestService_ListPipelines tests listing pipelines
func TestService_ListPipelines(t *testing.T) {
	orgID := uuid.New().String()

	t.Run("successful list", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		pipelines := []*Pipeline{
			{
				Id:             uuid.New(),
				OrganizationId: uuid.New(),
				Name:           "Pipeline 1",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			{
				Id:             uuid.New(),
				OrganizationId: uuid.New(),
				Name:           "Pipeline 2",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
		}

		params := ListPipelinesParams{
			Page: PageQuery{
				Number: 1,
				Size:   10,
			},
		}

		mockRepo.EXPECT().List(mock.Anything, params).Return(pipelines, int64(2), nil)

		result, total, err := service.List(context.Background(), orgID, 10, 0)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, 2, total)
	})
}
