package artifacts

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
func createTestArtifactsService(t *testing.T) (*Service, *MockRepository) {
	t.Helper()

	mockArtifactRepo := new(MockRepository)
	logger := logger.NewTest()

	service := NewArtifactsService(mockArtifactRepo, nil, logger)
	return service, mockArtifactRepo
}

// TestArtifactsService_Create tests creating an artifact
func TestArtifactsService_Create(t *testing.T) {
	orgID := uuid.New().String()
	producerID := uuid.New().String()

	t.Run("successful creation", func(t *testing.T) {
		service, mockRepo := createTestArtifactsService(t)

		req := &CreateArtifactJSONRequestBody{
			Name: "Test Artifact",
			Text: "Test content",
		}

		expectedArtifact := &Artifact{
			Id:             uuid.New(),
			OrganizationId: orgID,
			Name:           "Test Artifact",
			Text:           "Test content",
			ProducerId:     producerID,
			Credits:        0.012,
			MimeType:       "",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		mockRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(a *Artifact) bool {
			return a.Name == "Test Artifact" && a.Text == "Test content"
		})).Return(expectedArtifact, nil)

		result, err := service.Create(context.Background(), req, orgID, producerID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Test Artifact", result.Name)
	})

	t.Run("artifact too large", func(t *testing.T) {
		service, _ := createTestArtifactsService(t)

		// Create a request with text larger than MaxArtifactSize
		largeText := make([]byte, MaxArtifactSize+1)
		for i := range largeText {
			largeText[i] = 'a'
		}

		req := &CreateArtifactJSONRequestBody{
			Name: "Large Artifact",
			Text: string(largeText),
		}

		result, err := service.Create(context.Background(), req, orgID, producerID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrArtifactTooLarge)
	})
}

// TestArtifactsService_Get tests getting an artifact
func TestArtifactsService_Get(t *testing.T) {
	artifactID := uuid.New()
	artifact := &Artifact{
		Id:             artifactID,
		OrganizationId: uuid.New().String(),
		Name:           "Test Artifact",
		Text:           "Test content",
		ProducerId:     uuid.New().String(),
		Credits:        10.0,
		MimeType:       "text/plain",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	t.Run("successful get", func(t *testing.T) {
		service, mockRepo := createTestArtifactsService(t)

		mockRepo.EXPECT().Get(mock.Anything, artifactID).Return(artifact, nil)

		result, err := service.Get(context.Background(), artifactID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, artifactID, result.Id)
	})

	t.Run("not found", func(t *testing.T) {
		service, mockRepo := createTestArtifactsService(t)

		mockRepo.EXPECT().Get(mock.Anything, artifactID).Return(nil, ErrArtifactNotFound)

		result, err := service.Get(context.Background(), artifactID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrArtifactNotFound)
	})
}

// TestArtifactsService_Update tests updating an artifact
func TestArtifactsService_Update(t *testing.T) {
	artifactID := uuid.New()

	t.Run("successful update", func(t *testing.T) {
		service, mockRepo := createTestArtifactsService(t)

		req := &UpdateArtifactJSONRequestBody{
			Name: "Updated Artifact",
			Text: "Updated content",
		}

		expectedArtifact := &Artifact{
			Id:             artifactID,
			OrganizationId: uuid.New().String(),
			Name:           "Updated Artifact",
			Text:           "Updated content",
			MimeType:       "",
			Credits:        0.016,
			UpdatedAt:      time.Now(),
		}

		// First get the artifact, then update it
		existingArtifact := &Artifact{
			Id:             artifactID,
			OrganizationId: uuid.New().String(),
			Name:           "Old Artifact",
			Text:           "Old content",
			MimeType:       "",
		}
		mockRepo.EXPECT().Get(mock.Anything, artifactID).Return(existingArtifact, nil)
		mockRepo.EXPECT().Update(mock.Anything, artifactID, mock.Anything).Return(expectedArtifact, nil)

		result, err := service.Update(context.Background(), artifactID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Updated Artifact", result.Name)
	})

	t.Run("not found", func(t *testing.T) {
		service, mockRepo := createTestArtifactsService(t)

		req := &UpdateArtifactJSONRequestBody{
			Name: "Updated Artifact",
		}

		mockRepo.EXPECT().Get(mock.Anything, artifactID).Return(nil, ErrArtifactNotFound)

		result, err := service.Update(context.Background(), artifactID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrArtifactNotFound)
	})
}

// TestArtifactsService_Delete tests deleting an artifact
func TestArtifactsService_Delete(t *testing.T) {
	artifactID := uuid.New()

	t.Run("successful deletion", func(t *testing.T) {
		service, mockRepo := createTestArtifactsService(t)

		existingArtifact := &Artifact{
			Id:             artifactID,
			OrganizationId: uuid.New().String(),
			Name:           "Test Artifact",
			Text:           "Test content",
		}

		mockRepo.EXPECT().Get(mock.Anything, artifactID).Return(existingArtifact, nil)
		mockRepo.EXPECT().Delete(mock.Anything, artifactID).Return(nil)

		err := service.Delete(context.Background(), artifactID)

		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		service, mockRepo := createTestArtifactsService(t)

		mockRepo.EXPECT().Get(mock.Anything, artifactID).Return(nil, ErrArtifactNotFound)

		err := service.Delete(context.Background(), artifactID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrArtifactNotFound)
	})
}

// TestArtifactsService_List tests listing artifacts
func TestArtifactsService_List(t *testing.T) {

	t.Run("successful list", func(t *testing.T) {
		service, mockRepo := createTestArtifactsService(t)

		orgID := uuid.New().String()
		limit := 10
		offset := 0

		artifacts := []*Artifact{
			{
				Id:             uuid.New(),
				OrganizationId: orgID,
				Name:           "Artifact 1",
				Text:           "Content 1",
				MimeType:       "text/plain",
			},
			{
				Id:             uuid.New(),
				OrganizationId: orgID,
				Name:           "Artifact 2",
				Text:           "Content 2",
				MimeType:       "text/plain",
			},
		}

		// Mock expects the internal repository call
		params := ListArtifactsParams{
			Limit:  limit,
			Offset: offset,
		}
		mockRepo.EXPECT().List(mock.Anything, params).Return(artifacts, int64(2), nil)

		results, total, err := service.List(context.Background(), orgID, limit, offset)

		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, 2, total)
	})
}
