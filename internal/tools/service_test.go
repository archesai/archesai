package tools

import (
	"context"
	"errors"
	"testing"

	"github.com/archesai/archesai/internal/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createTestService(t *testing.T) (*Service, *MockRepository) {
	t.Helper()
	mockRepo := NewMockRepository(t)
	logger := logger.NewTest()
	service := NewService(mockRepo, logger)
	return service, mockRepo
}

func TestService_List(t *testing.T) {
	tools := []*Tool{
		{
			Id:             uuid.New(),
			OrganizationId: "org-1",
			Name:           "Tool 1",
			Description:    "Description 1",
		},
		{
			Id:             uuid.New(),
			OrganizationId: "org-1",
			Name:           "Tool 2",
			Description:    "Description 2",
		},
	}

	t.Run("successful list", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().List(mock.Anything, ListToolsParams{Limit: 10, Offset: 0}).
			Return(tools, int64(2), nil)

		result, total, err := service.List(context.Background(), "org-1", 10, 0)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, int64(2), total)
	})

	t.Run("empty list", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().List(mock.Anything, ListToolsParams{Limit: 10, Offset: 0}).
			Return([]*Tool{}, int64(0), nil)

		result, total, err := service.List(context.Background(), "org-1", 10, 0)

		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.Equal(t, int64(0), total)
	})

	t.Run("list error", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().List(mock.Anything, ListToolsParams{Limit: 10, Offset: 0}).
			Return(nil, int64(0), errors.New("database error"))

		result, total, err := service.List(context.Background(), "org-1", 10, 0)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, int64(0), total)
	})
}

func TestService_Create(t *testing.T) {
	t.Run("successful create", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		req := &CreateToolJSONRequestBody{
			Name:        "New Tool",
			Description: "Tool description",
		}

		expectedTool := &Tool{
			Id:             uuid.New(),
			OrganizationId: "org-1",
			Name:           "New Tool",
			Description:    "Tool description",
		}

		mockRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(t *Tool) bool {
			return t.Name == "New Tool" &&
				t.Description == "Tool description" &&
				t.OrganizationId == "org-1"
		})).Return(expectedTool, nil)

		result, err := service.Create(context.Background(), req, "org-1")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "New Tool", result.Name)
		assert.Equal(t, "Tool description", result.Description)
	})

	t.Run("create error", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		req := &CreateToolJSONRequestBody{
			Name:        "New Tool",
			Description: "Tool description",
		}

		mockRepo.EXPECT().Create(mock.Anything, mock.Anything).
			Return(nil, errors.New("creation failed"))

		result, err := service.Create(context.Background(), req, "org-1")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestService_Get(t *testing.T) {
	toolID := uuid.New()
	tool := &Tool{
		Id:             toolID,
		OrganizationId: "org-1",
		Name:           "Test Tool",
		Description:    "Test description",
	}

	t.Run("successful get", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().Get(mock.Anything, toolID).Return(tool, nil)

		result, err := service.Get(context.Background(), toolID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, toolID, result.Id)
		assert.Equal(t, "Test Tool", result.Name)
	})

	t.Run("not found", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().Get(mock.Anything, toolID).
			Return(nil, errors.New("not found"))

		result, err := service.Get(context.Background(), toolID)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestService_Update(t *testing.T) {
	toolID := uuid.New()
	existingTool := &Tool{
		Id:             toolID,
		OrganizationId: "org-1",
		Name:           "Old Name",
		Description:    "Old description",
	}

	t.Run("successful update", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		updatedTool := &Tool{
			Id:             toolID,
			OrganizationId: "org-1",
			Name:           "New Name",
			Description:    "New description",
		}

		mockRepo.EXPECT().Get(mock.Anything, toolID).Return(existingTool, nil)
		mockRepo.EXPECT().Update(mock.Anything, toolID, mock.MatchedBy(func(t *Tool) bool {
			return t.Name == "New Name" && t.Description == "New description"
		})).Return(updatedTool, nil)

		req := &UpdateToolJSONRequestBody{
			Name:        "New Name",
			Description: "New description",
		}

		result, err := service.Update(context.Background(), toolID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "New Name", result.Name)
		assert.Equal(t, "New description", result.Description)
	})

	t.Run("partial update - name only", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		// Clone the existing tool to avoid mutation
		toolToUpdate := &Tool{
			Id:             toolID,
			OrganizationId: "org-1",
			Name:           "Old Name",
			Description:    "Old description",
		}

		updatedTool := &Tool{
			Id:             toolID,
			OrganizationId: "org-1",
			Name:           "New Name",
			Description:    "Old description",
		}

		mockRepo.EXPECT().Get(mock.Anything, toolID).Return(toolToUpdate, nil)
		mockRepo.EXPECT().Update(mock.Anything, toolID, mock.MatchedBy(func(t *Tool) bool {
			return t.Name == "New Name" && t.Description == "Old description"
		})).Return(updatedTool, nil)

		req := &UpdateToolJSONRequestBody{
			Name: "New Name",
		}

		result, err := service.Update(context.Background(), toolID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "New Name", result.Name)
		assert.Equal(t, "Old description", result.Description)
	})

	t.Run("partial update - description only", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		// Clone the existing tool to avoid mutation
		toolToUpdate := &Tool{
			Id:             toolID,
			OrganizationId: "org-1",
			Name:           "Old Name",
			Description:    "Old description",
		}

		updatedTool := &Tool{
			Id:             toolID,
			OrganizationId: "org-1",
			Name:           "Old Name",
			Description:    "New description",
		}

		mockRepo.EXPECT().Get(mock.Anything, toolID).Return(toolToUpdate, nil)
		mockRepo.EXPECT().Update(mock.Anything, toolID, mock.MatchedBy(func(t *Tool) bool {
			return t.Name == "Old Name" && t.Description == "New description"
		})).Return(updatedTool, nil)

		req := &UpdateToolJSONRequestBody{
			Description: "New description",
		}

		result, err := service.Update(context.Background(), toolID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Old Name", result.Name)
		assert.Equal(t, "New description", result.Description)
	})

	t.Run("tool not found", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().Get(mock.Anything, toolID).
			Return(nil, errors.New("not found"))

		req := &UpdateToolJSONRequestBody{
			Name: "New Name",
		}

		result, err := service.Update(context.Background(), toolID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("update error", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().Get(mock.Anything, toolID).Return(existingTool, nil)
		mockRepo.EXPECT().Update(mock.Anything, toolID, mock.Anything).
			Return(nil, errors.New("update failed"))

		req := &UpdateToolJSONRequestBody{
			Name: "New Name",
		}

		result, err := service.Update(context.Background(), toolID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestService_Delete(t *testing.T) {
	toolID := uuid.New()

	t.Run("successful delete", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().Delete(mock.Anything, toolID).Return(nil)

		err := service.Delete(context.Background(), toolID)

		assert.NoError(t, err)
	})

	t.Run("delete error", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().Delete(mock.Anything, toolID).
			Return(errors.New("delete failed"))

		err := service.Delete(context.Background(), toolID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete failed")
	})
}
