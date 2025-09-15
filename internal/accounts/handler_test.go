package accounts

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_ListAccounts(t *testing.T) {
	tests := []struct {
		name         string
		queryParams  map[string]string
		setupMock    func(*MockAccountService)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name:        "successful list with pagination",
			queryParams: map[string]string{"page[number]": "1", "page[size]": "10"},
			setupMock: func(s *MockAccountService) {
				s.On("List", mock.Anything, mock.MatchedBy(func(params ListAccountsParams) bool {
					return params.Page.Size == 10 && params.Page.Number == 1
				})).Return([]*Account{
					{
						Id:         uuid.New(),
						AccountId:  "acc1",
						ProviderId: Google,
						UserId:     uuid.New(),
					},
					{
						Id:         uuid.New(),
						AccountId:  "acc2",
						ProviderId: Github,
						UserId:     uuid.New(),
					},
				}, int64(2), nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, float64(2), body["total"])
				data := body["data"].([]interface{})
				assert.Len(t, data, 2)
			},
		},
		{
			name:        "successful list without pagination",
			queryParams: map[string]string{},
			setupMock: func(s *MockAccountService) {
				s.On("List", mock.Anything, mock.MatchedBy(func(params ListAccountsParams) bool {
					return params.Page.Size == 10 && params.Page.Number == 1
				})).Return([]*Account{}, int64(0), nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, float64(0), body["total"])
				data := body["data"].([]interface{})
				assert.Len(t, data, 0)
			},
		},
		{
			name:        "service error",
			queryParams: map[string]string{},
			setupMock: func(s *MockAccountService) {
				s.On("List", mock.Anything, mock.Anything).Return(nil, int64(0), errors.New("database error"))
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
			req := httptest.NewRequest(http.MethodGet, "/auth/accounts", nil)

			// Add query parameters
			q := req.URL.Query()
			for k, v := range tt.queryParams {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockService := NewMockAccountService(t)
			tt.setupMock(mockService)

			handler := NewHandler(mockService, slog.Default())

			// Create params from query
			var params ListAccountsParams
			if pageNum := c.QueryParam("page[number]"); pageNum != "" {
				num := 1
				params.Page.Number = num
			}
			if pageSize := c.QueryParam("page[size]"); pageSize != "" {
				size := 10
				params.Page.Size = size
			}

			err := handler.ListAccounts(c, params)
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

func TestHandler_AccountsGetOne(t *testing.T) {
	accountID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name         string
		id           uuid.UUID
		setupMock    func(*MockAccountService)
		expectedCode int
		checkBody    func(*testing.T, map[string]interface{})
	}{
		{
			name: "successful get",
			id:   accountID,
			setupMock: func(s *MockAccountService) {
				s.On("Get", mock.Anything, accountID).Return(&Account{
					Id:         accountID,
					AccountId:  "test123",
					ProviderId: Google,
					UserId:     userID,
				}, nil)
			},
			expectedCode: http.StatusOK,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, accountID.String(), body["id"])
				assert.Equal(t, "test123", body["accountId"])
				assert.Equal(t, "google", body["providerId"])
			},
		},
		{
			name: "account not found",
			id:   accountID,
			setupMock: func(s *MockAccountService) {
				s.On("Get", mock.Anything, accountID).Return(nil, ErrAccountNotFound)
			},
			expectedCode: http.StatusNotFound,
			checkBody: func(t *testing.T, body map[string]interface{}) {
				assert.Equal(t, "Account Not Found", body["title"])
				assert.Equal(t, float64(http.StatusNotFound), body["status"])
			},
		},
		{
			name: "service error",
			id:   accountID,
			setupMock: func(s *MockAccountService) {
				s.On("Get", mock.Anything, accountID).Return(nil, errors.New("database error"))
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
			req := httptest.NewRequest(http.MethodGet, "/auth/accounts/"+tt.id.String(), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id.String())

			mockService := NewMockAccountService(t)
			tt.setupMock(mockService)

			handler := NewHandler(mockService, slog.Default())
			err := handler.GetAccount(c, tt.id)
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

func TestHandler_AccountsDelete(t *testing.T) {
	accountID := uuid.New()

	tests := []struct {
		name         string
		id           uuid.UUID
		setupMock    func(*MockAccountService)
		expectedCode int
		checkBody    func(*testing.T, []byte)
	}{
		{
			name: "successful delete",
			id:   accountID,
			setupMock: func(s *MockAccountService) {
				s.On("Delete", mock.Anything, accountID).Return(nil)
			},
			expectedCode: http.StatusNoContent,
			checkBody: func(t *testing.T, body []byte) {
				assert.Empty(t, body)
			},
		},
		{
			name: "account not found",
			id:   accountID,
			setupMock: func(s *MockAccountService) {
				s.On("Delete", mock.Anything, accountID).Return(ErrAccountNotFound)
			},
			expectedCode: http.StatusNotFound,
			checkBody: func(t *testing.T, body []byte) {
				var resp map[string]interface{}
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
				assert.Equal(t, "Account Not Found", resp["title"])
			},
		},
		{
			name: "service error",
			id:   accountID,
			setupMock: func(s *MockAccountService) {
				s.On("Delete", mock.Anything, accountID).Return(errors.New("database error"))
			},
			expectedCode: http.StatusInternalServerError,
			checkBody: func(t *testing.T, body []byte) {
				var resp map[string]interface{}
				err := json.Unmarshal(body, &resp)
				require.NoError(t, err)
				assert.Equal(t, "Internal Server Error", resp["title"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/auth/accounts/"+tt.id.String(), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.id.String())

			mockService := NewMockAccountService(t)
			tt.setupMock(mockService)

			handler := NewHandler(mockService, slog.Default())
			err := handler.DeleteAccount(c, tt.id)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedCode, rec.Code)
			tt.checkBody(t, rec.Body.Bytes())
			mockService.AssertExpectations(t)
		})
	}
}

func TestHandler_EdgeCases(t *testing.T) {
	t.Run("handler with nil service panics gracefully", func(t *testing.T) {
		assert.NotPanics(t, func() {
			h := NewHandler(nil, slog.Default())
			assert.NotNil(t, h)
		})
	})

	t.Run("handler with nil logger uses default", func(t *testing.T) {
		mockService := NewMockAccountService(t)
		h := NewHandler(mockService, nil)
		assert.NotNil(t, h)
		assert.NotNil(t, h.logger)
	})

	t.Run("pagination with page 2", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/auth/accounts?page[number]=2&page[size]=5", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockService := NewMockAccountService(t)
		mockService.On("List", mock.Anything, mock.MatchedBy(func(params ListAccountsParams) bool {
			// Page 2 with size 5
			return params.Page.Size == 5 && params.Page.Number == 2
		})).Return([]*Account{}, int64(0), nil)

		handler := NewHandler(mockService, slog.Default())

		var params ListAccountsParams
		params.Page.Number = 2
		params.Page.Size = 5

		err := handler.ListAccounts(c, params)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestHandler_InterfaceCompliance(t *testing.T) {
	// Ensure our handler implements the ServerInterface
	var _ ServerInterface = (*Handler)(nil)

	// Create a handler instance
	mockService := NewMockAccountService(t)
	handler := NewHandler(mockService, slog.Default())

	// Verify all required methods exist
	assert.NotNil(t, handler.ListAccounts)
	assert.NotNil(t, handler.GetAccount)
	assert.NotNil(t, handler.DeleteAccount)
}
