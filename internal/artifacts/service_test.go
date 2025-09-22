package artifacts

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/archesai/archesai/internal/logger"
)

// Test helper functions.
func createTestArtifactsService(t *testing.T) (*Service, *MockRepository) {
	t.Helper()

	mockArtifactRepo := NewMockRepository(t)
	logger := logger.NewTest()

	service := NewService(mockArtifactRepo, nil, logger)
	return service, mockArtifactRepo
}

// TestArtifactsService_Create tests creating an artifact.
func TestArtifactsService_Create(t *testing.T) {
	orgID := uuid.New()
	producerID := uuid.New()

	t.Run("successful creation", func(t *testing.T) {
		service, mockRepo := createTestArtifactsService(t)

		req := &CreateArtifactJSONRequestBody{
			Name: "Test Artifact",
			Text: "Test content",
		}

		expectedArtifact := &Artifact{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Name:           "Test Artifact",
			Text:           "Test content",
			ProducerID:     producerID,
			Credits:        0.012,
			MimeType:       "",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		mockRepo.EXPECT().
			Create(mock.Anything, mock.AnythingOfType("*artifacts.Artifact")).
			Return(expectedArtifact, nil)

		request := CreateArtifactRequestObject{
			Body: req,
		}
		result, err := service.CreateArtifact(context.Background(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		if successResp, ok := result.(CreateArtifact201JSONResponse); ok {
			assert.Equal(t, "Test Artifact", successResp.Data.Name)
		} else {
			t.Fatal("Expected CreateArtifact201JSONResponse")
		}
	})

	t.Run("repository error", func(t *testing.T) {
		service, mockRepo := createTestArtifactsService(t)

		req := &CreateArtifactJSONRequestBody{
			Name: "Test Artifact",
			Text: "Test content",
		}

		// Mock repository to return an error
		mockRepo.EXPECT().
			Create(mock.Anything, mock.AnythingOfType("*artifacts.Artifact")).
			Return(nil, assert.AnError)

		request := CreateArtifactRequestObject{
			Body: req,
		}
		result, err := service.CreateArtifact(context.Background(), request)

		assert.NoError(t, err) // Service never returns Go errors
		assert.NotNil(t, result)

		// Check for error response type
		_, isErrorResp := result.(CreateArtifact400ApplicationProblemPlusJSONResponse)
		assert.True(t, isErrorResp, "Expected error response type")
	})
}

// TestArtifactsService_Get tests getting an artifact.
func TestArtifactsService_Get(t *testing.T) {
	artifactID := uuid.New()
	artifact := &Artifact{
		ID:             artifactID,
		OrganizationID: uuid.New(),
		Name:           "Test Artifact",
		Text:           "Test content",
		ProducerID:     uuid.New(),
		Credits:        10.0,
		MimeType:       "text/plain",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	t.Run("successful get", func(t *testing.T) {
		service, mockRepo := createTestArtifactsService(t)

		mockRepo.EXPECT().Get(mock.Anything, artifactID).Return(artifact, nil)

		request := GetArtifactRequestObject{
			ID: artifactID,
		}
		result, err := service.GetArtifact(context.Background(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		if successResp, ok := result.(GetArtifact200JSONResponse); ok {
			assert.Equal(t, artifactID, successResp.Data.ID)
		} else {
			t.Fatal("Expected GetArtifact200JSONResponse")
		}
	})

	t.Run("not found", func(t *testing.T) {
		service, mockRepo := createTestArtifactsService(t)

		mockRepo.EXPECT().Get(mock.Anything, artifactID).Return(nil, ErrArtifactNotFound)

		request := GetArtifactRequestObject{
			ID: artifactID,
		}
		result, err := service.GetArtifact(context.Background(), request)

		assert.NoError(t, err) // Service never returns Go errors
		assert.NotNil(t, result)

		// Check for error response type
		_, isErrorResp := result.(GetArtifact404ApplicationProblemPlusJSONResponse)
		assert.True(t, isErrorResp, "Expected error response type")
	})
}

// TestArtifactsService_Update tests updating an artifact.
func TestArtifactsService_Update(t *testing.T) {
	artifactID := uuid.New()

	t.Run("successful update", func(t *testing.T) {
		service, mockRepo := createTestArtifactsService(t)

		req := &UpdateArtifactJSONRequestBody{
			Name: "Updated Artifact",
			Text: "Updated content",
		}

		expectedArtifact := &Artifact{
			ID:             artifactID,
			OrganizationID: uuid.New(),
			Name:           "Updated Artifact",
			Text:           "Updated content",
			MimeType:       "",
			Credits:        0.016,
			UpdatedAt:      time.Now(),
		}

		// First get the artifact, then update it
		existingArtifact := &Artifact{
			ID:             artifactID,
			OrganizationID: uuid.New(),
			Name:           "Old Artifact",
			Text:           "Old content",
			MimeType:       "",
		}
		mockRepo.EXPECT().Get(mock.Anything, artifactID).Return(existingArtifact, nil)
		mockRepo.EXPECT().
			Update(mock.Anything, artifactID, mock.Anything).
			Return(expectedArtifact, nil)

		request := UpdateArtifactRequestObject{
			ID:   artifactID,
			Body: req,
		}
		result, err := service.UpdateArtifact(context.Background(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		if successResp, ok := result.(UpdateArtifact200JSONResponse); ok {
			assert.Equal(t, "Updated Artifact", successResp.Data.Name)
		} else {
			t.Fatal("Expected UpdateArtifact200JSONResponse")
		}
	})

	t.Run("not found", func(t *testing.T) {
		service, mockRepo := createTestArtifactsService(t)

		req := &UpdateArtifactJSONRequestBody{
			Name: "Updated Artifact",
		}

		mockRepo.EXPECT().Get(mock.Anything, artifactID).Return(nil, ErrArtifactNotFound)

		request := UpdateArtifactRequestObject{
			ID:   artifactID,
			Body: req,
		}
		result, err := service.UpdateArtifact(context.Background(), request)

		assert.NoError(t, err) // Service never returns Go errors
		assert.NotNil(t, result)

		// Check for error response type
		_, isErrorResp := result.(UpdateArtifact404ApplicationProblemPlusJSONResponse)
		assert.True(t, isErrorResp, "Expected error response type")
	})
}

// TestArtifactsService_Delete tests deleting an artifact.
func TestArtifactsService_Delete(t *testing.T) {
	artifactID := uuid.New()

	t.Run("successful deletion", func(t *testing.T) {
		service, mockRepo := createTestArtifactsService(t)

		existingArtifact := &Artifact{
			ID:             artifactID,
			OrganizationID: uuid.New(),
			Name:           "Test Artifact",
			Text:           "Test content",
		}

		mockRepo.EXPECT().Get(mock.Anything, artifactID).Return(existingArtifact, nil)
		mockRepo.EXPECT().Delete(mock.Anything, artifactID).Return(nil)

		request := DeleteArtifactRequestObject{
			ID: artifactID,
		}
		_, err := service.DeleteArtifact(context.Background(), request)

		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		service, mockRepo := createTestArtifactsService(t)

		mockRepo.EXPECT().Get(mock.Anything, artifactID).Return(nil, ErrArtifactNotFound)

		request := DeleteArtifactRequestObject{
			ID: artifactID,
		}
		result, err := service.DeleteArtifact(context.Background(), request)

		assert.NoError(t, err) // Service never returns Go errors
		assert.NotNil(t, result)

		// Check for error response type
		_, isErrorResp := result.(DeleteArtifact404ApplicationProblemPlusJSONResponse)
		assert.True(t, isErrorResp, "Expected error response type")
	})
}

// TestArtifactsService_List tests listing artifacts.
func TestArtifactsService_List(t *testing.T) {

	t.Run("successful list", func(t *testing.T) {
		service, mockRepo := createTestArtifactsService(t)

		orgID := uuid.New()
		limit := 10
		offset := 0

		artifacts := []*Artifact{
			{
				ID:             uuid.New(),
				OrganizationID: orgID,
				Name:           "Artifact 1",
				Text:           "Content 1",
				MimeType:       "text/plain",
			},
			{
				ID:             uuid.New(),
				OrganizationID: orgID,
				Name:           "Artifact 2",
				Text:           "Content 2",
				MimeType:       "text/plain",
			},
		}

		// Mock expects the internal repository call
		params := ListArtifactsParams{
			Page: PageQuery{
				Number: offset/limit + 1,
				Size:   limit,
			},
		}
		mockRepo.EXPECT().List(mock.Anything, params).Return(artifacts, int64(2), nil)

		request := ListArtifactsRequestObject{
			Params: ListArtifactsParams{
				Page: PageQuery{
					Number: offset/limit + 1,
					Size:   limit,
				},
			},
		}
		result, err := service.ListArtifacts(context.Background(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		if successResp, ok := result.(ListArtifacts200JSONResponse); ok {
			assert.Len(t, successResp.Data, 2)
			assert.Equal(t, float32(2), successResp.Meta.Total)
		} else {
			t.Fatal("Expected ListArtifacts200JSONResponse")
		}
	})
}
