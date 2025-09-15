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

	service := NewService(mockRepo, nil, logger)
	return service, mockRepo
}

// TestService_CreatePipeline tests creating a pipeline
func TestService_CreatePipeline(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		req := &CreatePipelineJSONRequestBody{
			Name:        "Test Pipeline",
			Description: "Test description",
		}

		expectedPipeline := &Pipeline{
			ID:             uuid.New(),
			OrganizationID: uuid.New(),
			Name:           "Test Pipeline",
			Description:    "Test description",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		mockRepo.EXPECT().Create(mock.Anything, mock.AnythingOfType("*pipelines.Pipeline")).Return(expectedPipeline, nil)

		request := CreatePipelineRequestObject{
			Body: req,
		}
		result, err := service.Create(context.Background(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		if successResp, ok := result.(CreatePipeline201JSONResponse); ok {
			assert.Equal(t, "Test Pipeline", successResp.Data.Name)
		} else {
			t.Fatalf("expected CreatePipeline201JSONResponse, got %T", result)
		}
	})
}

// TestService_GetPipeline tests getting a pipeline
func TestService_GetPipeline(t *testing.T) {
	pipelineID := uuid.New()
	pipeline := &Pipeline{
		ID:             pipelineID,
		OrganizationID: uuid.New(),
		Name:           "Test Pipeline",
		Description:    "Test description",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	t.Run("successful get", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().Get(mock.Anything, pipelineID).Return(pipeline, nil)

		request := GetPipelineRequestObject{
			ID: pipelineID,
		}
		result, err := service.Get(context.Background(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		if successResp, ok := result.(GetPipeline200JSONResponse); ok {
			assert.Equal(t, pipelineID, successResp.Data.ID)
		} else {
			t.Fatalf("expected GetPipeline200JSONResponse, got %T", result)
		}
	})

	t.Run("not found", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().Get(mock.Anything, pipelineID).Return(nil, ErrPipelineNotFound)

		request := GetPipelineRequestObject{
			ID: pipelineID,
		}
		result, err := service.Get(context.Background(), request)

		assert.NoError(t, err) // Service never returns Go errors
		assert.NotNil(t, result)

		// Check for error response type
		_, isErrorResp := result.(GetPipeline404ApplicationProblemPlusJSONResponse)
		assert.True(t, isErrorResp, "Expected error response type")
	})
}

// TestService_UpdatePipeline tests updating a pipeline
func TestService_UpdatePipeline(t *testing.T) {
	pipelineID := uuid.New()

	t.Run("successful update", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		existingPipeline := &Pipeline{
			ID:             pipelineID,
			OrganizationID: uuid.New(),
			Name:           "Old Name",
			Description:    "Old description",
			CreatedAt:      time.Now().Add(-1 * time.Hour),
			UpdatedAt:      time.Now().Add(-1 * time.Hour),
		}

		updatedPipeline := &Pipeline{
			ID:             pipelineID,
			OrganizationID: existingPipeline.OrganizationID,
			Name:           "New Name",
			Description:    "New description",
			CreatedAt:      existingPipeline.CreatedAt,
			UpdatedAt:      time.Now(),
		}

		req := &UpdatePipelineJSONRequestBody{
			Name:        "New Name",
			Description: "New description",
		}

		mockRepo.EXPECT().Get(mock.Anything, pipelineID).Return(existingPipeline, nil)
		mockRepo.EXPECT().Update(mock.Anything, pipelineID, mock.AnythingOfType("*pipelines.Pipeline")).Return(updatedPipeline, nil)

		request := UpdatePipelineRequestObject{
			ID:   pipelineID,
			Body: req,
		}
		result, err := service.Update(context.Background(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		if successResp, ok := result.(UpdatePipeline200JSONResponse); ok {
			assert.Equal(t, "New Name", successResp.Data.Name)
			assert.Equal(t, "New description", successResp.Data.Description)
		} else {
			t.Fatalf("expected UpdatePipeline200JSONResponse, got %T", result)
		}
	})
}

// TestService_DeletePipeline tests deleting a pipeline
func TestService_DeletePipeline(t *testing.T) {
	pipelineID := uuid.New()

	t.Run("successful delete", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		pipeline := &Pipeline{
			ID:          pipelineID,
			Name:        "Test Pipeline",
			Description: "Test description",
		}

		mockRepo.EXPECT().Get(mock.Anything, pipelineID).Return(pipeline, nil)
		mockRepo.EXPECT().Delete(mock.Anything, pipelineID).Return(nil)

		request := DeletePipelineRequestObject{
			ID: pipelineID,
		}
		res, err := service.Delete(context.Background(), request)
		assert.NotNil(t, res)

		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().Get(mock.Anything, pipelineID).Return(nil, ErrPipelineNotFound)

		request := DeletePipelineRequestObject{
			ID: pipelineID,
		}
		result, err := service.Delete(context.Background(), request)

		assert.NoError(t, err) // Service never returns Go errors
		assert.NotNil(t, result)

		// Check for error response type
		_, isErrorResp := result.(DeletePipeline404ApplicationProblemPlusJSONResponse)
		assert.True(t, isErrorResp, "Expected error response type")
	})
}

// TestService_ListPipelines tests listing pipelines
func TestService_ListPipelines(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		pipelines := []*Pipeline{
			{
				ID:             uuid.New(),
				OrganizationID: uuid.New(),
				Name:           "Pipeline 1",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			{
				ID:             uuid.New(),
				OrganizationID: uuid.New(),
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

		request := ListPipelinesRequestObject{
			Params: ListPipelinesParams{
				Page: PageQuery{
					Number: 1,
					Size:   10,
				},
			},
		}
		result, err := service.List(context.Background(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		if successResp, ok := result.(ListPipelines200JSONResponse); ok {
			assert.Len(t, successResp.Data, 2)
			assert.Equal(t, float32(2), successResp.Meta.Total)
		} else {
			t.Fatalf("expected ListPipelines200JSONResponse, got %T", result)
		}
	})
}
