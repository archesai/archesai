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

		mockRepo.EXPECT().CreateOrganization(mock.Anything, mock.MatchedBy(func(o *Organization) bool {
			return o.Name == "" && o.BillingEmail == "billing@test.com"
		})).Return(expectedOrg, nil)

		// Also expect CreateMember for the initial owner
		mockRepo.EXPECT().CreateMember(mock.Anything, mock.MatchedBy(func(m *Member) bool {
			return m.OrganizationId == expectedOrg.Id.String() && m.Role == MemberRoleOwner
		})).Return(&Member{
			Id:             uuid.New(),
			OrganizationId: expectedOrg.Id.String(),
			Role:           MemberRoleOwner,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}, nil)

		req := &CreateOrganizationRequest{
			OrganizationId: uuid.New(),
			BillingEmail:   "billing@test.com",
		}

		result, err := service.CreateOrganization(context.Background(), req, "creator-id")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		// Note: Name is returned from the mock, not from the request
		assert.Equal(t, "Test Org", result.Name)
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

		mockRepo.EXPECT().GetOrganization(mock.Anything, orgID).Return(org, nil)

		result, err := service.GetOrganization(context.Background(), orgID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, orgID, result.Id)
	})

	t.Run("organization exists", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().GetOrganization(mock.Anything, orgID).Return(org, nil)

		result, err := service.GetOrganization(context.Background(), orgID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, orgID, result.Id)
	})

	t.Run("not found", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		mockRepo.EXPECT().GetOrganization(mock.Anything, orgID).Return(nil, ErrOrganizationNotFound)

		result, err := service.GetOrganization(context.Background(), orgID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrOrganizationNotFound)
	})
}

// TestService_CreateMember tests adding a member to an organization
func TestService_CreateMember(t *testing.T) {
	orgID := uuid.New()

	t.Run("successful creation", func(t *testing.T) {
		service, mockRepo := createTestService(t)

		expectedMember := &Member{
			Id:             uuid.New(),
			OrganizationId: orgID.String(),
			Role:           MemberRole("member"),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		mockRepo.EXPECT().CreateMember(mock.Anything, mock.MatchedBy(func(m *Member) bool {
			return m.OrganizationId == orgID.String() && m.Role == MemberRole("member")
		})).Return(expectedMember, nil)

		req := &CreateMemberRequest{
			Role: CreateMemberJSONBodyRole("member"),
		}

		result, err := service.CreateMember(context.Background(), req, orgID.String())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, orgID.String(), result.OrganizationId)
		assert.Equal(t, MemberRole("member"), result.Role)
	})
}
