package organizations

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

// TestService_CreateOrganization tests creating an organization
func TestService_CreateOrganization(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		expectedOrg := &Organization{
			Id:           uuid.New(),
			Name:         "Test Org",
			BillingEmail: "billing@test.com",
			Plan:         OrganizationPlan("free"),
			Credits:      0,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		mockRepo.EXPECT().Create(mock.Anything, mock.MatchedBy(func(o *Organization) bool {
			return o.Name == "" && o.BillingEmail == "billing@test.com"
		})).Return(expectedOrg, nil)

		req := &CreateOrganizationRequest{
			OrganizationId: uuid.New(),
			BillingEmail:   "billing@test.com",
		}

		result, err := service.Create(context.Background(), req, "creator-id")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		// Note: Name is returned from the mock, not from the request
		assert.Equal(t, "Test Org", result.Name)
	})

	t.Run("create fails", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().Create(mock.Anything, mock.Anything).Return(nil, assert.AnError)

		req := &CreateOrganizationRequest{
			OrganizationId: uuid.New(),
			BillingEmail:   "billing@test.com",
		}

		result, err := service.Create(context.Background(), req, "creator-id")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

// TestService_GetOrganization tests getting an organization
func TestService_GetOrganization(t *testing.T) {
	orgID := uuid.New()
	org := &Organization{
		Id:           orgID,
		Name:         "Test Org",
		BillingEmail: "billing@test.com",
		Plan:         OrganizationPlan("free"),
		Credits:      100,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	t.Run("successful get", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().Get(mock.Anything, orgID).Return(org, nil)

		result, err := service.Get(context.Background(), orgID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, orgID, result.Id)
	})

	t.Run("not found", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().Get(mock.Anything, orgID).Return(nil, ErrOrganizationNotFound)

		result, err := service.Get(context.Background(), orgID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrOrganizationNotFound)
	})
}

// TestService_UpdateOrganization tests updating an organization
func TestService_UpdateOrganization(t *testing.T) {
	orgID := uuid.New()
	existingOrg := &Organization{
		Id:           orgID,
		Name:         "Test Org",
		BillingEmail: "billing@test.com",
		Plan:         OrganizationPlan("free"),
		Credits:      100,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	t.Run("successful update", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		updatedOrg := &Organization{
			Id:           orgID,
			Name:         "Test Org",
			BillingEmail: "new-billing@test.com",
			Plan:         OrganizationPlan("free"),
			Credits:      100,
			CreatedAt:    existingOrg.CreatedAt,
			UpdatedAt:    time.Now(),
		}

		mockRepo.EXPECT().Get(mock.Anything, orgID).Return(existingOrg, nil)
		mockRepo.EXPECT().Update(mock.Anything, orgID, mock.MatchedBy(func(o *Organization) bool {
			return o.BillingEmail == "new-billing@test.com"
		})).Return(updatedOrg, nil)

		req := &UpdateOrganizationRequest{
			BillingEmail: "new-billing@test.com",
		}

		result, err := service.Update(context.Background(), orgID, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "new-billing@test.com", string(result.BillingEmail))
	})

	t.Run("organization not found", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().Get(mock.Anything, orgID).Return(nil, ErrOrganizationNotFound)

		req := &UpdateOrganizationRequest{
			BillingEmail: "new-billing@test.com",
		}

		result, err := service.Update(context.Background(), orgID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("update fails", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().Get(mock.Anything, orgID).Return(existingOrg, nil)
		mockRepo.EXPECT().Update(mock.Anything, orgID, mock.Anything).Return(nil, assert.AnError)

		req := &UpdateOrganizationRequest{
			BillingEmail: "new-billing@test.com",
		}

		result, err := service.Update(context.Background(), orgID, req)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

// TestService_DeleteOrganization tests deleting an organization
func TestService_DeleteOrganization(t *testing.T) {
	orgID := uuid.New()

	t.Run("successful delete", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().Delete(mock.Anything, orgID).Return(nil)

		err := service.Delete(context.Background(), orgID)

		assert.NoError(t, err)
	})

	t.Run("delete fails", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().Delete(mock.Anything, orgID).Return(assert.AnError)

		err := service.Delete(context.Background(), orgID)

		assert.Error(t, err)
	})
}

// TestService_ListOrganizations tests listing organizations
func TestService_ListOrganizations(t *testing.T) {
	orgs := []*Organization{
		{
			Id:           uuid.New(),
			Name:         "Org 1",
			BillingEmail: "billing1@test.com",
			Plan:         OrganizationPlan("free"),
			Credits:      100,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Id:           uuid.New(),
			Name:         "Org 2",
			BillingEmail: "billing2@test.com",
			Plan:         OrganizationPlan("pro"),
			Credits:      200,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	t.Run("successful list", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().List(mock.Anything, ListOrganizationsParams{Limit: 10, Offset: 0}).Return(orgs, int64(2), nil)

		result, total, err := service.List(context.Background(), 10, 0)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, 2, total)
	})

	t.Run("empty list", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().List(mock.Anything, ListOrganizationsParams{Limit: 10, Offset: 0}).Return([]*Organization{}, int64(0), nil)

		result, total, err := service.List(context.Background(), 10, 0)

		assert.NoError(t, err)
		assert.Empty(t, result)
		assert.Equal(t, 0, total)
	})

	t.Run("list fails", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().List(mock.Anything, ListOrganizationsParams{Limit: 10, Offset: 0}).Return(nil, int64(0), assert.AnError)

		result, total, err := service.List(context.Background(), 10, 0)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, 0, total)
	})
}
