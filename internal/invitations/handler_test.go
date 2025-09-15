package invitations

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_ListInvitations(t *testing.T) {
	organizationID := uuid.New()

	tests := []struct {
		name         string
		orgID        openapi_types.UUID
		params       ListInvitationsParams
		setupMock    func(*MockInvitationService)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name:  "successful list",
			orgID: organizationID,
			params: ListInvitationsParams{
				Page: Page{Number: 1, Size: 10},
			},
			setupMock: func(s *MockInvitationService) {
				s.On("ListByOrganization", mock.Anything, organizationID.String()).
					Return([]*Invitation{
						{
							Id:             uuid.New(),
							Email:          "test1@example.com",
							Role:           InvitationRoleAdmin,
							Status:         StatusPending,
							OrganizationId: organizationID.String(),
						},
						{
							Id:             uuid.New(),
							Email:          "test2@example.com",
							Role:           InvitationRoleMember,
							Status:         StatusAccepted,
							OrganizationId: organizationID.String(),
						},
					}, nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, float64(2), body["meta"].(map[string]interface{})["total"])
				data := body["data"].([]interface{})
				assert.Len(t, data, 2)
			},
		},
		{
			name:   "empty list",
			orgID:  organizationID,
			params: ListInvitationsParams{},
			setupMock: func(s *MockInvitationService) {
				s.On("ListByOrganization", mock.Anything, organizationID.String()).
					Return([]*Invitation{}, nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, float64(0), body["meta"].(map[string]interface{})["total"])
				data, ok := body["data"].([]interface{})
				if ok {
					assert.Len(t, data, 0)
				} else {
					// If data is nil, that's also acceptable for empty list
					assert.Nil(t, body["data"])
				}
			},
		},
		{
			name:   "service error",
			orgID:  organizationID,
			params: ListInvitationsParams{},
			setupMock: func(s *MockInvitationService) {
				s.On("ListByOrganization", mock.Anything, organizationID.String()).
					Return(nil, errors.New("database error"))
			},
			expectedCode: http.StatusInternalServerError,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Internal Server Error", body["title"])
				assert.Equal(t, float64(http.StatusInternalServerError), body["status"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockService := NewMockInvitationService(t)
			tt.setupMock(mockService)

			handler := NewHandler(mockService, slog.Default())
			err := handler.ListInvitations(c, tt.orgID, tt.params)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedCode, rec.Code)

			var body map[string]interface{}
			err = json.Unmarshal(rec.Body.Bytes(), &body)
			require.NoError(t, err)

			tt.checkBody(t, body)
			mockService.AssertExpectations(t)
		})
	}
}

func TestHandler_CreateInvitation(t *testing.T) {
	organizationID := uuid.New()

	tests := []struct {
		name         string
		orgID        openapi_types.UUID
		requestBody  CreateInvitationJSONRequestBody
		setupMock    func(*MockInvitationService)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name:  "successful creation",
			orgID: organizationID,
			requestBody: CreateInvitationJSONRequestBody{
				Email: "new@example.com",
				Role:  CreateInvitationJSONBodyRoleAdmin,
			},
			setupMock: func(s *MockInvitationService) {
				s.On("Create", mock.Anything, mock.AnythingOfType("*invitations.Invitation")).
					Return(&Invitation{
						Id:             uuid.New(),
						Email:          "new@example.com",
						Role:           InvitationRoleAdmin,
						Status:         StatusPending,
						OrganizationId: organizationID.String(),
					}, nil)
			},
			expectedCode: http.StatusCreated,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				data := body["data"].(map[string]interface{})
				assert.Equal(t, "new@example.com", data["email"])
				assert.Equal(t, "admin", data["role"])
				assert.Equal(t, StatusPending, data["status"])
			},
		},
		{
			name:  "invitation already exists",
			orgID: organizationID,
			requestBody: CreateInvitationJSONRequestBody{
				Email: "existing@example.com",
				Role:  CreateInvitationJSONBodyRoleMember,
			},
			setupMock: func(s *MockInvitationService) {
				s.On("Create", mock.Anything, mock.AnythingOfType("*invitations.Invitation")).
					Return(nil, ErrInvitationAlreadyExists)
			},
			expectedCode: http.StatusConflict,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Conflict", body["title"])
				assert.Equal(t, float64(http.StatusConflict), body["status"])
			},
		},
		{
			name:  "service error",
			orgID: organizationID,
			requestBody: CreateInvitationJSONRequestBody{
				Email: "error@example.com",
				Role:  CreateInvitationJSONBodyRoleOwner,
			},
			setupMock: func(s *MockInvitationService) {
				s.On("Create", mock.Anything, mock.AnythingOfType("*invitations.Invitation")).
					Return(nil, errors.New("database error"))
			},
			expectedCode: http.StatusInternalServerError,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Internal Server Error", body["title"])
				assert.Equal(t, float64(http.StatusInternalServerError), body["status"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockService := NewMockInvitationService(t)
			tt.setupMock(mockService)

			handler := NewHandler(mockService, slog.Default())
			err := handler.CreateInvitation(c, tt.orgID)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedCode, rec.Code)

			var respBody map[string]interface{}
			err = json.Unmarshal(rec.Body.Bytes(), &respBody)
			require.NoError(t, err)

			tt.checkBody(t, respBody)
			mockService.AssertExpectations(t)
		})
	}
}

func TestHandler_GetInvitation(t *testing.T) {
	organizationID := uuid.New()
	invitationID := uuid.New()

	tests := []struct {
		name         string
		orgID        openapi_types.UUID
		invID        openapi_types.UUID
		setupMock    func(*MockInvitationService)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name:  "successful get",
			orgID: organizationID,
			invID: invitationID,
			setupMock: func(s *MockInvitationService) {
				s.On("Get", mock.Anything, invitationID).
					Return(&Invitation{
						Id:             invitationID,
						Email:          "test@example.com",
						Role:           InvitationRoleAdmin,
						Status:         StatusPending,
						OrganizationId: organizationID.String(),
					}, nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				data := body["data"].(map[string]interface{})
				assert.Equal(t, invitationID.String(), data["id"])
				assert.Equal(t, "test@example.com", data["email"])
			},
		},
		{
			name:  "invitation not found",
			orgID: organizationID,
			invID: invitationID,
			setupMock: func(s *MockInvitationService) {
				s.On("Get", mock.Anything, invitationID).
					Return(nil, ErrInvitationNotFound)
			},
			expectedCode: http.StatusNotFound,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Not Found", body["title"])
				assert.Equal(t, float64(http.StatusNotFound), body["status"])
			},
		},
		{
			name:  "wrong organization",
			orgID: organizationID,
			invID: invitationID,
			setupMock: func(s *MockInvitationService) {
				s.On("Get", mock.Anything, invitationID).
					Return(&Invitation{
						Id:             invitationID,
						OrganizationId: uuid.New().String(),
					}, nil)
			},
			expectedCode: http.StatusNotFound,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Not Found", body["title"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockService := NewMockInvitationService(t)
			tt.setupMock(mockService)

			handler := NewHandler(mockService, slog.Default())
			err := handler.GetInvitation(c, tt.orgID, tt.invID)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedCode, rec.Code)

			var body map[string]interface{}
			err = json.Unmarshal(rec.Body.Bytes(), &body)
			require.NoError(t, err)

			tt.checkBody(t, body)
			mockService.AssertExpectations(t)
		})
	}
}

func TestHandler_UpdateInvitation(t *testing.T) {
	organizationID := uuid.New()
	invitationID := uuid.New()

	tests := []struct {
		name         string
		orgID        openapi_types.UUID
		invID        openapi_types.UUID
		requestBody  UpdateInvitationJSONRequestBody
		setupMock    func(*MockInvitationService)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name:  "successful update",
			orgID: organizationID,
			invID: invitationID,
			requestBody: UpdateInvitationJSONRequestBody{
				Role: Owner,
			},
			setupMock: func(s *MockInvitationService) {
				s.On("Get", mock.Anything, invitationID).
					Return(&Invitation{
						Id:             invitationID,
						Email:          "test@example.com",
						Role:           InvitationRoleMember,
						Status:         StatusPending,
						OrganizationId: organizationID.String(),
					}, nil)
				s.On("Update", mock.Anything, invitationID, mock.AnythingOfType("*invitations.Invitation")).
					Return(&Invitation{
						Id:             invitationID,
						Email:          "test@example.com",
						Role:           InvitationRoleOwner,
						Status:         StatusPending,
						OrganizationId: organizationID.String(),
					}, nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				data := body["data"].(map[string]interface{})
				assert.Equal(t, "owner", data["role"])
			},
		},
		{
			name:  "invitation expired",
			orgID: organizationID,
			invID: invitationID,
			requestBody: UpdateInvitationJSONRequestBody{
				Email: "updated@example.com",
			},
			setupMock: func(s *MockInvitationService) {
				s.On("Get", mock.Anything, invitationID).
					Return(&Invitation{
						Id:             invitationID,
						OrganizationId: organizationID.String(),
					}, nil)
				s.On("Update", mock.Anything, invitationID, mock.AnythingOfType("*invitations.Invitation")).
					Return(nil, ErrInvitationExpired)
			},
			expectedCode: http.StatusBadRequest,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Bad Request", body["title"])
				assert.Contains(t, body["detail"], "expired")
			},
		},
		{
			name:  "already accepted",
			orgID: organizationID,
			invID: invitationID,
			requestBody: UpdateInvitationJSONRequestBody{
				Role: Admin,
			},
			setupMock: func(s *MockInvitationService) {
				s.On("Get", mock.Anything, invitationID).
					Return(&Invitation{
						Id:             invitationID,
						OrganizationId: organizationID.String(),
					}, nil)
				s.On("Update", mock.Anything, invitationID, mock.AnythingOfType("*invitations.Invitation")).
					Return(nil, ErrInvitationAlreadyAccepted)
			},
			expectedCode: http.StatusBadRequest,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Bad Request", body["title"])
				assert.Contains(t, body["detail"], "already been accepted")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPatch, "/", bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockService := NewMockInvitationService(t)
			tt.setupMock(mockService)

			handler := NewHandler(mockService, slog.Default())
			err := handler.UpdateInvitation(c, tt.orgID, tt.invID)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedCode, rec.Code)

			var respBody map[string]interface{}
			err = json.Unmarshal(rec.Body.Bytes(), &respBody)
			require.NoError(t, err)

			tt.checkBody(t, respBody)
			mockService.AssertExpectations(t)
		})
	}
}

func TestHandler_DeleteInvitation(t *testing.T) {
	organizationID := uuid.New()
	invitationID := uuid.New()

	tests := []struct {
		name         string
		orgID        openapi_types.UUID
		invID        openapi_types.UUID
		setupMock    func(*MockInvitationService)
		expectedCode int
	}{
		{
			name:  "successful delete",
			orgID: organizationID,
			invID: invitationID,
			setupMock: func(s *MockInvitationService) {
				s.On("Get", mock.Anything, invitationID).
					Return(&Invitation{
						Id:             invitationID,
						Status:         StatusPending,
						OrganizationId: organizationID.String(),
					}, nil)
				s.On("Delete", mock.Anything, invitationID).Return(nil)
			},
			expectedCode: http.StatusNoContent,
		},
		{
			name:  "invitation not found",
			orgID: organizationID,
			invID: invitationID,
			setupMock: func(s *MockInvitationService) {
				s.On("Get", mock.Anything, invitationID).
					Return(nil, ErrInvitationNotFound)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:  "wrong organization",
			orgID: organizationID,
			invID: invitationID,
			setupMock: func(s *MockInvitationService) {
				s.On("Get", mock.Anything, invitationID).
					Return(&Invitation{
						Id:             invitationID,
						OrganizationId: uuid.New().String(),
					}, nil)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:  "delete error",
			orgID: organizationID,
			invID: invitationID,
			setupMock: func(s *MockInvitationService) {
				s.On("Get", mock.Anything, invitationID).
					Return(&Invitation{
						Id:             invitationID,
						OrganizationId: organizationID.String(),
					}, nil)
				s.On("Delete", mock.Anything, invitationID).
					Return(errors.New("database error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockService := NewMockInvitationService(t)
			tt.setupMock(mockService)

			handler := NewHandler(mockService, slog.Default())
			err := handler.DeleteInvitation(c, tt.orgID, tt.invID)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedCode, rec.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestNewHandler(t *testing.T) {
	t.Run("with logger", func(t *testing.T) {
		mockService := NewMockInvitationService(t)
		logger := slog.Default()
		handler := NewHandler(mockService, logger)

		assert.NotNil(t, handler)
		assert.Equal(t, mockService, handler.service)
		assert.Equal(t, logger, handler.logger)
	})

	t.Run("without logger uses default", func(t *testing.T) {
		mockService := NewMockInvitationService(t)
		handler := NewHandler(mockService, nil)

		assert.NotNil(t, handler)
		assert.Equal(t, mockService, handler.service)
		assert.NotNil(t, handler.logger)
	})
}

func TestHandler_InterfaceCompliance(t *testing.T) {
	// Ensure our handler implements the ServerInterface
	var _ ServerInterface = (*Handler)(nil)

	// Create a handler instance
	mockService := NewMockInvitationService(t)
	handler := NewHandler(mockService, slog.Default())

	// Verify all required methods exist
	assert.NotNil(t, handler.ListInvitations)
	assert.NotNil(t, handler.CreateInvitation)
	assert.NotNil(t, handler.GetInvitation)
	assert.NotNil(t, handler.UpdateInvitation)
	assert.NotNil(t, handler.DeleteInvitation)
}
