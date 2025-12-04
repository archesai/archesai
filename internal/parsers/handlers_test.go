package parsers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlerParser_ParseHandler(t *testing.T) {
	// Create temp directory for test handlers
	tempDir := t.TempDir()

	tests := []struct {
		name         string
		handlerCode  string
		operationID  string
		wantHandler  bool
		wantParams   int
		validateDeps func(t *testing.T, deps []DependencyDef)
	}{
		{
			name: "simple handler with auth service",
			handlerCode: `package handlers

import (
	"context"
	"github.com/archesai/archesai/pkg/auth"
)

type LoginHandler struct {
	authService *auth.Service
}

func NewLoginHandler(authService *auth.Service) *LoginHandler {
	return &LoginHandler{authService: authService}
}

func (h *LoginHandler) Execute(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	return nil, nil
}
`,
			operationID: "Login",
			wantHandler: true,
			wantParams:  1,
			validateDeps: func(t *testing.T, deps []DependencyDef) {
				require.Len(t, deps, 1)
				assert.Equal(t, "authService", deps[0].Name)
				assert.Equal(t, "*auth.Service", deps[0].Type)
				assert.True(t, deps[0].IsPointer)
				assert.Equal(t, "infra.AuthService", deps[0].Resolution)
			},
		},
		{
			name: "handler with multiple dependencies",
			handlerCode: `package handlers

import (
	"context"
	"github.com/archesai/archesai/pkg/auth"
	"myapp/generated/core/repositories"
)

type RegisterHandler struct {
	authService *auth.Service
	userRepo    repositories.UserRepository
}

func NewRegisterHandler(authService *auth.Service, userRepo repositories.UserRepository) *RegisterHandler {
	return &RegisterHandler{
		authService: authService,
		userRepo:    userRepo,
	}
}

func (h *RegisterHandler) Execute(ctx context.Context, input *RegisterInput) (*RegisterOutput, error) {
	return nil, nil
}
`,
			operationID: "Register",
			wantHandler: true,
			wantParams:  2,
			validateDeps: func(t *testing.T, deps []DependencyDef) {
				require.Len(t, deps, 2)

				assert.Equal(t, "authService", deps[0].Name)
				assert.Equal(t, "*auth.Service", deps[0].Type)
				assert.Equal(t, "infra.AuthService", deps[0].Resolution)

				assert.Equal(t, "userRepo", deps[1].Name)
				assert.Equal(t, "repositories.UserRepository", deps[1].Type)
				assert.Equal(t, "repos.Users", deps[1].Resolution)
			},
		},
		{
			name: "handler with repository only",
			handlerCode: `package handlers

import (
	"context"
	"myapp/generated/core/repositories"
)

type GetUserHandler struct {
	userRepo repositories.UserRepository
}

func NewGetUserHandler(userRepo repositories.UserRepository) *GetUserHandler {
	return &GetUserHandler{userRepo: userRepo}
}

func (h *GetUserHandler) Execute(ctx context.Context, input *GetUserInput) (*GetUserOutput, error) {
	return nil, nil
}
`,
			operationID: "GetUser",
			wantHandler: true,
			wantParams:  1,
			validateDeps: func(t *testing.T, deps []DependencyDef) {
				require.Len(t, deps, 1)
				assert.Equal(t, "userRepo", deps[0].Name)
				assert.Equal(t, "repositories.UserRepository", deps[0].Type)
				assert.False(t, deps[0].IsPointer)
				assert.Equal(t, "repos.Users", deps[0].Resolution)
			},
		},
		{
			name: "handler with event publisher",
			handlerCode: `package handlers

import (
	"context"
	"github.com/archesai/archesai/pkg/events"
	"myapp/generated/core/repositories"
)

type CreateUserHandler struct {
	userRepo  repositories.UserRepository
	publisher events.Publisher
}

func NewCreateUserHandler(userRepo repositories.UserRepository, publisher events.Publisher) *CreateUserHandler {
	return &CreateUserHandler{
		userRepo:  userRepo,
		publisher: publisher,
	}
}

func (h *CreateUserHandler) Execute(ctx context.Context, input *CreateUserInput) (*CreateUserOutput, error) {
	return nil, nil
}
`,
			operationID: "CreateUser",
			wantHandler: true,
			wantParams:  2,
			validateDeps: func(t *testing.T, deps []DependencyDef) {
				require.Len(t, deps, 2)

				assert.Equal(t, "userRepo", deps[0].Name)
				assert.Equal(t, "repos.Users", deps[0].Resolution)

				assert.Equal(t, "publisher", deps[1].Name)
				assert.Equal(t, "events.Publisher", deps[1].Type)
				assert.Equal(t, "infra.EventPublisher", deps[1].Resolution)
			},
		},
		{
			name:        "no handler file exists",
			handlerCode: "", // Don't create file
			operationID: "NonExistent",
			wantHandler: false,
			wantParams:  0,
		},
		{
			name: "handler with unknown dependency",
			handlerCode: `package handlers

import (
	"context"
	"myapp/custom/unknown"
)

type CustomHandler struct {
	svc *unknown.CustomService
}

func NewCustomHandler(svc *unknown.CustomService) *CustomHandler {
	return &CustomHandler{svc: svc}
}

func (h *CustomHandler) Execute(ctx context.Context, input *CustomInput) (*CustomOutput, error) {
	return nil, nil
}
`,
			operationID: "Custom",
			wantHandler: true,
			wantParams:  1,
			validateDeps: func(t *testing.T, deps []DependencyDef) {
				require.Len(t, deps, 1)
				assert.Equal(t, "svc", deps[0].Name)
				assert.Equal(t, "*unknown.CustomService", deps[0].Type)
				assert.Contains(t, deps[0].Resolution, "TODO") // Unknown type
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create handler file if code is provided
			if tt.handlerCode != "" {
				fileName := SnakeCase(tt.operationID) + ".impl.go"
				filePath := filepath.Join(tempDir, fileName)
				err := os.WriteFile(filePath, []byte(tt.handlerCode), 0644)
				require.NoError(t, err)
			}

			// Create parser and parse
			parser := NewHandlerParser(tempDir)
			op := OperationDef{ID: tt.operationID}

			handlers, err := parser.ParseHandlers([]OperationDef{op})
			require.NoError(t, err)
			require.Len(t, handlers, 1)

			handler := handlers[0]
			assert.Equal(t, tt.operationID, handler.OperationID)
			assert.Equal(t, tt.wantHandler, handler.HasHandler)

			if tt.wantHandler {
				require.NotNil(t, handler.Constructor)
				assert.Equal(t, "New"+tt.operationID+"Handler", handler.Constructor.Name)
				assert.Len(t, handler.Constructor.Parameters, tt.wantParams)

				// Validate dependencies
				if tt.validateDeps != nil {
					deps := parser.GetDependencies(handler)
					tt.validateDeps(t, deps)
				}
			}
		})
	}
}

func TestDependencyRegistry_Resolve(t *testing.T) {
	registry := NewDependencyRegistry()

	tests := []struct {
		name       string
		typeName   string
		wantFound  bool
		wantResult string
	}{
		{
			name:       "auth service pointer",
			typeName:   "*auth.Service",
			wantFound:  true,
			wantResult: "infra.AuthService",
		},
		{
			name:       "events publisher interface",
			typeName:   "events.Publisher",
			wantFound:  true,
			wantResult: "infra.EventPublisher",
		},
		{
			name:       "user repository",
			typeName:   "repositories.UserRepository",
			wantFound:  true,
			wantResult: "repos.Users",
		},
		{
			name:       "session repository",
			typeName:   "repositories.SessionRepository",
			wantFound:  true,
			wantResult: "repos.Sessions",
		},
		{
			name:       "pipeline repository",
			typeName:   "repositories.PipelineRepository",
			wantFound:  true,
			wantResult: "repos.Pipelines",
		},
		{
			name:       "account repository",
			typeName:   "repositories.AccountRepository",
			wantFound:  true,
			wantResult: "repos.Accounts",
		},
		{
			name:      "unknown type",
			typeName:  "custom.UnknownService",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, found := registry.Resolve(tt.typeName)
			assert.Equal(t, tt.wantFound, found)
			if tt.wantFound {
				assert.Equal(t, tt.wantResult, result.Resolution)
			}
		})
	}
}

func TestPluralize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"User", "Users"},
		{"Session", "Sessions"},
		{"Account", "Accounts"},
		{"Pipeline", "Pipelines"},
		{"Entity", "Entities"},
		{"Category", "Categories"},
		{"Box", "Boxes"},
		{"Match", "Matches"},
		{"Dish", "Dishes"},
		{"Bus", "Buses"},
		{"Executor", "Executors"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := Pluralize(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTypeToString(t *testing.T) {
	// We can't easily test typeToString directly since it works on AST nodes,
	// but we can test it indirectly through parsing actual handler code
	tempDir := t.TempDir()

	// Create a parser with the temp directory
	parser := NewHandlerParser(tempDir)

	handlerCode := `package handlers

type TestHandler struct {
	slice      []string
	mapField   map[string]int
	ptrField   *string
	qualified  pkg.Type
	nested     *pkg.Nested
}

func NewTestHandler(
	slice []string,
	mapField map[string]int,
	ptrField *string,
	qualified pkg.Type,
	nested *pkg.Nested,
) *TestHandler {
	return &TestHandler{}
}
`

	filePath := filepath.Join(tempDir, "test.impl.go")
	err := os.WriteFile(filePath, []byte(handlerCode), 0644)
	require.NoError(t, err)

	handlers, err := parser.ParseHandlers([]OperationDef{{ID: "Test"}})
	require.NoError(t, err)
	require.Len(t, handlers, 1)
	require.NotNil(t, handlers[0].Constructor)

	params := handlers[0].Constructor.Parameters
	require.Len(t, params, 5)

	assert.Equal(t, "[]string", params[0].GoType)
	assert.Equal(t, "map[string]int", params[1].GoType)
	assert.Equal(t, "*string", params[2].GoType)
	assert.Equal(t, "pkg.Type", params[3].GoType)
	assert.Equal(t, "*pkg.Nested", params[4].GoType)
}
