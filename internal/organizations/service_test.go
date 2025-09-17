package organizations

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/archesai/archesai/internal/logger"
)

// Test helper functions.
func createTestService(t *testing.T) (*Service, *MockRepository) {
	t.Helper()

	mockRepo := NewMockRepository(t)
	logger := logger.NewTest()

	// Using concrete service that implements the Service interface
	service := NewService(mockRepo, nil, logger)
	return service, mockRepo
}

// TestService_Interface tests that the service implements the interface correctly.
func TestService_Interface(t *testing.T) {
	service, mockRepo := createTestService(t)

	// Test that the service implements the Service interface
	var _ = service

	ctx := context.Background()
	orgID := uuid.New()

	// Test Get - expect repository error to return not found response
	mockRepo.EXPECT().Get(ctx, orgID).Return(nil, ErrOrganizationNotFound)

	result, err := service.Get(ctx, GetOrganizationRequestObject{
		ID: orgID,
	})
	assert.NoError(t, err) // Service never returns Go errors
	assert.NotNil(t, result)

	// Check for error response type
	_, isErrorResp := result.(GetOrganization404ApplicationProblemPlusJSONResponse)
	assert.True(t, isErrorResp, "Expected error response type")
}
