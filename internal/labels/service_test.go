package labels

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
func createTestLabelsService(t *testing.T) (*Service, *MockRepository) {
	t.Helper()

	mockRepo := new(MockRepository)
	logger := logger.NewTest()

	service := NewService(mockRepo, logger)
	return service, mockRepo
}

// TestLabelsService_Create tests creating a label
func TestLabelsService_Create(t *testing.T) {
	orgID := uuid.New()

	t.Run("successful creation", func(t *testing.T) {
		service, mockRepo := createTestLabelsService(t)

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

		mockRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(l *Label) bool {
			return l.Name == "Test Label" && l.OrganizationId == orgID
		})).Return(expectedLabel, nil)

		result, err := service.Create(context.Background(), req, orgID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Test Label", result.Name)
		assert.Equal(t, orgID, result.OrganizationId)
	})

	t.Run("empty name", func(t *testing.T) {
		service, _ := createTestLabelsService(t)

		req := &CreateLabelJSONRequestBody{
			Name: "",
		}

		result, err := service.Create(context.Background(), req, orgID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "label name is required")
	})

	t.Run("empty organization ID", func(t *testing.T) {
		service, _ := createTestLabelsService(t)

		req := &CreateLabelJSONRequestBody{
			Name: "Test Label",
		}

		result, err := service.Create(context.Background(), req, uuid.Nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "organization ID is required")
	})
}

// TestLabelsService_Get tests getting a label
func TestLabelsService_Get(t *testing.T) {
	labelID := uuid.New()
	label := &Label{
		Id:             labelID,
		OrganizationId: uuid.New(),
		Name:           "Test Label",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	t.Run("successful get", func(t *testing.T) {
		service, mockRepo := createTestLabelsService(t)

		mockRepo.EXPECT().Get(mock.Anything, labelID).Return(label, nil)

		result, err := service.Get(context.Background(), labelID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, labelID, result.Id)
		assert.Equal(t, "Test Label", result.Name)
	})

	t.Run("not found", func(t *testing.T) {
		service, mockRepo := createTestLabelsService(t)

		mockRepo.EXPECT().Get(mock.Anything, labelID).Return(nil, ErrLabelNotFound)

		result, err := service.Get(context.Background(), labelID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrLabelNotFound)
	})
}

// TestLabelsService_Update tests updating a label
func TestLabelsService_Update(t *testing.T) {
	labelID := uuid.New()

	t.Run("successful update", func(t *testing.T) {
		service, mockRepo := createTestLabelsService(t)

		req := &UpdateLabelJSONRequestBody{
			Name: "Updated Label",
		}

		existingLabel := &Label{
			Id:             labelID,
			OrganizationId: uuid.New(),
			Name:           "Old Label",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		expectedLabel := &Label{
			Id:             labelID,
			OrganizationId: existingLabel.OrganizationId,
			Name:           "Updated Label",
			CreatedAt:      existingLabel.CreatedAt,
			UpdatedAt:      time.Now(),
		}

		mockRepo.EXPECT().Get(mock.Anything, labelID).Return(existingLabel, nil)
		mockRepo.EXPECT().Update(mock.Anything, labelID, mock.Anything).Return(expectedLabel, nil)

		result, err := service.Update(context.Background(), labelID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Updated Label", result.Name)
	})

	t.Run("not found", func(t *testing.T) {
		service, mockRepo := createTestLabelsService(t)

		req := &UpdateLabelJSONRequestBody{
			Name: "Updated Label",
		}

		mockRepo.EXPECT().Get(mock.Anything, labelID).Return(nil, ErrLabelNotFound)

		result, err := service.Update(context.Background(), labelID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrLabelNotFound)
	})
}

// TestLabelsService_Delete tests deleting a label
func TestLabelsService_Delete(t *testing.T) {
	labelID := uuid.New()

	t.Run("successful deletion", func(t *testing.T) {
		service, mockRepo := createTestLabelsService(t)

		existingLabel := &Label{
			Id:             labelID,
			OrganizationId: uuid.New(),
			Name:           "Test Label",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		mockRepo.EXPECT().Get(mock.Anything, labelID).Return(existingLabel, nil)
		mockRepo.EXPECT().Delete(mock.Anything, labelID).Return(nil)

		err := service.Delete(context.Background(), labelID)

		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		service, mockRepo := createTestLabelsService(t)

		mockRepo.EXPECT().Get(mock.Anything, labelID).Return(nil, ErrLabelNotFound)

		err := service.Delete(context.Background(), labelID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrLabelNotFound)
	})
}

// TestLabelsService_List tests listing labels
func TestLabelsService_List(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		service, mockRepo := createTestLabelsService(t)

		orgID := uuid.New()
		limit := 10
		offset := 0

		labels := []*Label{
			{
				Id:             uuid.New(),
				OrganizationId: orgID,
				Name:           "Label 1",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			{
				Id:             uuid.New(),
				OrganizationId: orgID,
				Name:           "Label 2",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
		}

		params := ListLabelsParams{
			Page: PageQuery{
				Number: offset/limit + 1,
				Size:   limit,
			},
		}
		mockRepo.EXPECT().List(mock.Anything, params).Return(labels, int64(2), nil)

		results, total, err := service.List(context.Background(), orgID, limit, offset)

		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, 2, total)
	})
}
