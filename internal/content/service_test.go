package content

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

// TestService_CreateArtifact tests creating an artifact
func TestService_CreateArtifact(t *testing.T) {
	orgID := uuid.New().String()
	producerID := uuid.New().String()

	t.Run("successful creation", func(t *testing.T) {
		service, mockRepo := createTestService(t)

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
			Credits:        0.0,
			MimeType:       "text/plain",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		mockRepo.EXPECT().CreateArtifact(mock.Anything, mock.MatchedBy(func(a *Artifact) bool {
			return a.Name == "Test Artifact" && a.Text == "Test content"
		})).Return(expectedArtifact, nil)

		result, err := service.CreateArtifact(context.Background(), req, orgID, producerID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Test Artifact", result.Name)
	})

	t.Run("artifact too large", func(t *testing.T) {
		service, _ := createTestService(t)

		// Create a request with text larger than MaxArtifactSize
		largeText := make([]byte, MaxArtifactSize+1)
		for i := range largeText {
			largeText[i] = 'a'
		}

		req := &CreateArtifactJSONRequestBody{
			Name: "Large Artifact",
			Text: string(largeText),
		}

		result, err := service.CreateArtifact(context.Background(), req, orgID, producerID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrArtifactTooLarge)
	})
}

// TestService_GetArtifact tests getting an artifact
func TestService_GetArtifact(t *testing.T) {
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
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().GetArtifact(mock.Anything, artifactID).Return(artifact, nil)

		result, err := service.GetArtifact(context.Background(), artifactID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, artifactID, result.Id)
	})

	t.Run("not found", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().GetArtifact(mock.Anything, artifactID).Return(nil, ErrArtifactNotFound)

		result, err := service.GetArtifact(context.Background(), artifactID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrArtifactNotFound)
	})
}

// TestService_UpdateArtifact tests updating an artifact
func TestService_UpdateArtifact(t *testing.T) {
	artifactID := uuid.New()

	t.Run("successful update", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		req := &UpdateArtifactJSONRequestBody{
			Name: "Updated Artifact",
			Text: "Updated content",
		}

		expectedArtifact := &Artifact{
			Id:             artifactID,
			OrganizationId: uuid.New().String(),
			Name:           "Updated Artifact",
			Text:           "Updated content",
			MimeType:       "text/plain",
			Credits:        0.0,
			UpdatedAt:      time.Now(),
		}

		// First get the artifact, then update it
		existingArtifact := &Artifact{
			Id:             artifactID,
			OrganizationId: uuid.New().String(),
			Name:           "Old Artifact",
			Text:           "Old content",
			MimeType:       "text/plain",
		}
		mockRepo.EXPECT().GetArtifact(mock.Anything, artifactID).Return(existingArtifact, nil)
		mockRepo.EXPECT().UpdateArtifact(mock.Anything, artifactID, mock.Anything).Return(expectedArtifact, nil)

		result, err := service.UpdateArtifact(context.Background(), artifactID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Updated Artifact", result.Name)
	})

	t.Run("not found", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		req := &UpdateArtifactJSONRequestBody{
			Name: "Updated Artifact",
		}

		mockRepo.EXPECT().GetArtifact(mock.Anything, artifactID).Return(nil, ErrArtifactNotFound)

		result, err := service.UpdateArtifact(context.Background(), artifactID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrArtifactNotFound)
	})
}

// TestService_DeleteArtifact tests deleting an artifact
func TestService_DeleteArtifact(t *testing.T) {
	artifactID := uuid.New()

	t.Run("successful deletion", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().DeleteArtifact(mock.Anything, artifactID).Return(nil)

		err := service.DeleteArtifact(context.Background(), artifactID)

		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().DeleteArtifact(mock.Anything, artifactID).Return(ErrArtifactNotFound)

		err := service.DeleteArtifact(context.Background(), artifactID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrArtifactNotFound)
	})
}

// TestService_ListArtifacts tests listing artifacts
func TestService_ListArtifacts(t *testing.T) {

	t.Run("successful list", func(t *testing.T) {
		service, mockRepo := createTestService(t)

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
		mockRepo.EXPECT().ListArtifacts(mock.Anything, params).Return(artifacts, int64(2), nil)

		results, total, err := service.ListArtifacts(context.Background(), orgID, limit, offset)

		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, 2, total)
	})
}

// TestService_CreateLabel tests creating a label
func TestService_CreateLabel(t *testing.T) {
	orgID := uuid.New().String()

	t.Run("successful creation", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		req := &CreateLabelJSONRequestBody{
			Name: "Test Label",
		}

		expectedLabel := &Label{
			Id:             uuid.New(),
			OrganizationId: orgID,
			Name:           "Test Label",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		mockRepo.EXPECT().CreateLabel(mock.Anything, mock.MatchedBy(func(l *Label) bool {
			return l.Name == "Test Label" && l.OrganizationId == orgID
		})).Return(expectedLabel, nil)

		result, err := service.CreateLabel(context.Background(), req, orgID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Test Label", result.Name)
		assert.Equal(t, orgID, result.OrganizationId)
	})
}

// TestService_GetLabel tests getting a label
func TestService_GetLabel(t *testing.T) {
	labelID := uuid.New()
	label := &Label{
		Id:             labelID,
		OrganizationId: uuid.New().String(),
		Name:           "Test Label",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	t.Run("successful get", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().GetLabel(mock.Anything, labelID).Return(label, nil)

		result, err := service.GetLabel(context.Background(), labelID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, labelID, result.Id)
		assert.Equal(t, "Test Label", result.Name)
	})

	t.Run("not found", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().GetLabel(mock.Anything, labelID).Return(nil, ErrLabelNotFound)

		result, err := service.GetLabel(context.Background(), labelID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrLabelNotFound)
	})
}
