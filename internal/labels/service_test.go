package labels

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
func createTestLabelsService(t *testing.T) (*Service, *MockRepository) {
	t.Helper()

	mockRepo := NewMockRepository(t)
	logger := logger.NewTest()

	service := NewService(mockRepo, nil, logger)
	return service, mockRepo
}

// TestLabelsService_Create tests creating a label.
func TestLabelsService_Create(t *testing.T) {
	orgID := uuid.New()

	t.Run("successful creation", func(t *testing.T) {
		service, mockRepo := createTestLabelsService(t)

		req := &CreateLabelJSONRequestBody{
			Name: "Test Label",
		}

		expectedLabel := &Label{
			ID:             uuid.New(),
			OrganizationID: orgID,
			Name:           "Test Label",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		mockRepo.EXPECT().
			Create(mock.Anything, mock.AnythingOfType("*labels.Label")).
			Return(expectedLabel, nil)

		request := CreateLabelRequestObject{
			Body: req,
		}
		result, err := service.CreateLabel(context.Background(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		if successResp, ok := result.(CreateLabel201JSONResponse); ok {
			assert.Equal(t, "Test Label", successResp.Data.Name)
			assert.Equal(t, orgID, successResp.Data.OrganizationID)
		} else {
			t.Fatal("Expected CreateLabel201JSONResponse")
		}
	})

	t.Run("repository error", func(t *testing.T) {
		service, mockRepo := createTestLabelsService(t)

		req := &CreateLabelJSONRequestBody{
			Name: "Test Label",
		}

		// Mock repository to return an error
		mockRepo.EXPECT().
			Create(mock.Anything, mock.AnythingOfType("*labels.Label")).
			Return(nil, assert.AnError)

		request := CreateLabelRequestObject{
			Body: req,
		}
		result, err := service.CreateLabel(context.Background(), request)

		assert.NoError(t, err) // Service never returns Go errors
		assert.NotNil(t, result)

		// Check for error response type
		_, isErrorResp := result.(CreateLabel400ApplicationProblemPlusJSONResponse)
		assert.True(t, isErrorResp, "Expected error response type")
	})

	t.Run("successful creation with minimal data", func(t *testing.T) {
		service, mockRepo := createTestLabelsService(t)

		req := &CreateLabelJSONRequestBody{
			Name: "Test Label",
		}

		expectedLabel := &Label{
			ID:             uuid.New(),
			OrganizationID: uuid.New(),
			Name:           "Test Label",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		mockRepo.EXPECT().
			Create(mock.Anything, mock.AnythingOfType("*labels.Label")).
			Return(expectedLabel, nil)

		request := CreateLabelRequestObject{
			Body: req,
		}
		result, err := service.CreateLabel(context.Background(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		if successResp, ok := result.(CreateLabel201JSONResponse); ok {
			assert.Equal(t, "Test Label", successResp.Data.Name)
		} else {
			t.Fatal("Expected CreateLabel201JSONResponse")
		}
	})
}

// TestLabelsService_Get tests getting a label.
func TestLabelsService_Get(t *testing.T) {
	labelID := uuid.New()
	label := &Label{
		ID:             labelID,
		OrganizationID: uuid.New(),
		Name:           "Test Label",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	t.Run("successful get", func(t *testing.T) {
		service, mockRepo := createTestLabelsService(t)

		mockRepo.EXPECT().Get(mock.Anything, labelID).Return(label, nil)

		request := GetLabelRequestObject{
			ID: labelID,
		}
		result, err := service.GetLabel(context.Background(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		if successResp, ok := result.(GetLabel200JSONResponse); ok {
			assert.Equal(t, labelID, successResp.Data.ID)
			assert.Equal(t, "Test Label", successResp.Data.Name)
		} else {
			t.Fatal("Expected GetLabel200JSONResponse")
		}
	})

	t.Run("not found", func(t *testing.T) {
		service, mockRepo := createTestLabelsService(t)

		mockRepo.EXPECT().Get(mock.Anything, labelID).Return(nil, ErrLabelNotFound)

		request := GetLabelRequestObject{
			ID: labelID,
		}
		result, err := service.GetLabel(context.Background(), request)

		assert.NoError(t, err) // Service never returns Go errors
		assert.NotNil(t, result)

		// Check for error response type
		_, isErrorResp := result.(GetLabel404ApplicationProblemPlusJSONResponse)
		assert.True(t, isErrorResp, "Expected error response type")
	})
}

// TestLabelsService_Update tests updating a label.
func TestLabelsService_Update(t *testing.T) {
	labelID := uuid.New()

	t.Run("successful update", func(t *testing.T) {
		service, mockRepo := createTestLabelsService(t)

		req := &UpdateLabelJSONRequestBody{
			Name: "Updated Label",
		}

		existingLabel := &Label{
			ID:             labelID,
			OrganizationID: uuid.New(),
			Name:           "Old Label",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		expectedLabel := &Label{
			ID:             labelID,
			OrganizationID: existingLabel.OrganizationID,
			Name:           "Updated Label",
			CreatedAt:      existingLabel.CreatedAt,
			UpdatedAt:      time.Now(),
		}

		mockRepo.EXPECT().Get(mock.Anything, labelID).Return(existingLabel, nil)
		mockRepo.EXPECT().Update(mock.Anything, labelID, mock.Anything).Return(expectedLabel, nil)

		request := UpdateLabelRequestObject{
			ID:   labelID,
			Body: req,
		}
		result, err := service.UpdateLabel(context.Background(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		if successResp, ok := result.(UpdateLabel200JSONResponse); ok {
			assert.Equal(t, "Updated Label", successResp.Data.Name)
		} else {
			t.Fatal("Expected UpdateLabel200JSONResponse")
		}
	})

	t.Run("not found", func(t *testing.T) {
		service, mockRepo := createTestLabelsService(t)

		req := &UpdateLabelJSONRequestBody{
			Name: "Updated Label",
		}

		mockRepo.EXPECT().Get(mock.Anything, labelID).Return(nil, ErrLabelNotFound)

		request := UpdateLabelRequestObject{
			ID:   labelID,
			Body: req,
		}
		result, err := service.UpdateLabel(context.Background(), request)

		assert.NoError(t, err) // Service never returns Go errors
		assert.NotNil(t, result)

		// Check for error response type
		_, isErrorResp := result.(UpdateLabel404ApplicationProblemPlusJSONResponse)
		assert.True(t, isErrorResp, "Expected error response type")
	})
}

// TestLabelsService_Delete tests deleting a label.
func TestLabelsService_Delete(t *testing.T) {
	labelID := uuid.New()

	t.Run("successful deletion", func(t *testing.T) {
		service, mockRepo := createTestLabelsService(t)

		existingLabel := &Label{
			ID:             labelID,
			OrganizationID: uuid.New(),
			Name:           "Test Label",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		mockRepo.EXPECT().Get(mock.Anything, labelID).Return(existingLabel, nil)
		mockRepo.EXPECT().Delete(mock.Anything, labelID).Return(nil)

		request := DeleteLabelRequestObject{
			ID: labelID,
		}
		_, err := service.DeleteLabel(context.Background(), request)

		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		service, mockRepo := createTestLabelsService(t)

		mockRepo.EXPECT().Get(mock.Anything, labelID).Return(nil, ErrLabelNotFound)

		request := DeleteLabelRequestObject{
			ID: labelID,
		}
		result, err := service.DeleteLabel(context.Background(), request)

		assert.NoError(t, err) // Service never returns Go errors
		assert.NotNil(t, result)

		// Check for error response type
		_, isErrorResp := result.(DeleteLabel404ApplicationProblemPlusJSONResponse)
		assert.True(t, isErrorResp, "Expected error response type")
	})
}

// TestLabelsService_List tests listing labels.
func TestLabelsService_List(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		service, mockRepo := createTestLabelsService(t)

		orgID := uuid.New()
		limit := 10
		offset := 0

		labels := []*Label{
			{
				ID:             uuid.New(),
				OrganizationID: orgID,
				Name:           "Label 1",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			{
				ID:             uuid.New(),
				OrganizationID: orgID,
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

		request := ListLabelsRequestObject{
			Params: ListLabelsParams{
				Page: PageQuery{
					Number: offset/limit + 1,
					Size:   limit,
				},
			},
		}
		result, err := service.ListLabels(context.Background(), request)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		if successResp, ok := result.(ListLabels200JSONResponse); ok {
			assert.Len(t, successResp.Data, 2)
			assert.Equal(t, int64(2), successResp.Meta.Total)
		} else {
			t.Fatal("Expected ListLabels200JSONResponse")
		}
	})
}
