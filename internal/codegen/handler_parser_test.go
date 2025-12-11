package codegen

import "testing"

func TestNewDependencyRegistry(t *testing.T) {
	registry := NewDependencyRegistry()

	if registry == nil {
		t.Fatal("NewDependencyRegistry() returned nil")
	}

	// Should have default mappings
	if registry.mappings == nil {
		t.Fatal("mappings should not be nil")
	}

	// Verify some default mappings exist
	if _, ok := registry.mappings["*auth.Service"]; !ok {
		t.Error("missing default mapping for *auth.Service")
	}
	if _, ok := registry.mappings["*config.Config"]; !ok {
		t.Error("missing default mapping for *config.Config")
	}
}

func TestDependencyRegistry_Resolve(t *testing.T) {
	registry := NewDependencyRegistry()

	t.Run("direct match", func(t *testing.T) {
		res, found := registry.Resolve("*auth.Service")
		if !found {
			t.Fatal("expected to find *auth.Service")
		}
		if res.Resolution != "infra.AuthService" {
			t.Errorf("Resolution = %q, want %q", res.Resolution, "infra.AuthService")
		}
		if res.ImportPath != "github.com/archesai/archesai/pkg/auth" {
			t.Errorf(
				"ImportPath = %q, want %q",
				res.ImportPath,
				"github.com/archesai/archesai/pkg/auth",
			)
		}
	})

	t.Run("events.Publisher", func(t *testing.T) {
		res, found := registry.Resolve("events.Publisher")
		if !found {
			t.Fatal("expected to find events.Publisher")
		}
		if res.Resolution != "infra.EventPublisher" {
			t.Errorf("Resolution = %q, want %q", res.Resolution, "infra.EventPublisher")
		}
	})

	t.Run("config pointer", func(t *testing.T) {
		res, found := registry.Resolve("*config.Config")
		if !found {
			t.Fatal("expected to find *config.Config")
		}
		if res.Resolution != "cfg" {
			t.Errorf("Resolution = %q, want %q", res.Resolution, "cfg")
		}
	})

	t.Run("repository pattern", func(t *testing.T) {
		tests := []struct {
			typeName string
			want     string
		}{
			{"repositories.UserRepository", "repos.Users"},
			{"repositories.SessionRepository", "repos.Sessions"},
			{"repositories.OrganizationRepository", "repos.Organizations"},
		}

		for _, tt := range tests {
			t.Run(tt.typeName, func(t *testing.T) {
				res, found := registry.Resolve(tt.typeName)
				if !found {
					t.Fatalf("expected to find %s", tt.typeName)
				}
				if res.Resolution != tt.want {
					t.Errorf("Resolution = %q, want %q", res.Resolution, tt.want)
				}
			})
		}
	})

	t.Run("unknown type", func(t *testing.T) {
		_, found := registry.Resolve("unknown.Type")
		if found {
			t.Error("expected not to find unknown.Type")
		}
	})
}

func TestDependencyRegistry_AddMapping(t *testing.T) {
	registry := NewDependencyRegistry()

	t.Run("add new mapping", func(t *testing.T) {
		registry.AddMapping("custom.Service", DependencyResolution{
			ImportPath: "example.com/custom",
			Resolution: "infra.CustomService",
		})

		res, found := registry.Resolve("custom.Service")
		if !found {
			t.Fatal("expected to find custom.Service after adding")
		}
		if res.Resolution != "infra.CustomService" {
			t.Errorf("Resolution = %q, want %q", res.Resolution, "infra.CustomService")
		}
	})

	t.Run("override existing mapping", func(t *testing.T) {
		registry.AddMapping("*auth.Service", DependencyResolution{
			ImportPath: "custom/auth",
			Resolution: "customAuth",
		})

		res, found := registry.Resolve("*auth.Service")
		if !found {
			t.Fatal("expected to find *auth.Service")
		}
		if res.Resolution != "customAuth" {
			t.Errorf("Resolution = %q, want %q", res.Resolution, "customAuth")
		}
	})
}

func TestNewHandlerParser(t *testing.T) {
	parser := NewHandlerParser("/test/handlers")

	if parser == nil {
		t.Fatal("NewHandlerParser() returned nil")
	}
	if parser.handlersDir != "/test/handlers" {
		t.Errorf("handlersDir = %q, want %q", parser.handlersDir, "/test/handlers")
	}
	if parser.depRegistry == nil {
		t.Error("depRegistry should not be nil")
	}
}

func TestHandlerParser_GetDependencies(t *testing.T) {
	parser := NewHandlerParser("/test/handlers")

	t.Run("nil constructor", func(t *testing.T) {
		handler := HandlerDef{
			OperationID: "Test",
			Constructor: nil,
		}

		deps := parser.GetDependencies(handler)
		if deps != nil {
			t.Errorf("expected nil deps for nil constructor, got %v", deps)
		}
	})

	t.Run("empty parameters", func(t *testing.T) {
		handler := HandlerDef{
			OperationID: "Test",
			Constructor: &ConstructorDef{
				Name:       "NewTestHandler",
				Parameters: nil,
			},
		}

		deps := parser.GetDependencies(handler)
		if len(deps) != 0 {
			t.Errorf("expected 0 deps, got %d", len(deps))
		}
	})
}

func TestDependencyDef(t *testing.T) {
	dep := DependencyDef{
		Name:       "authService",
		Type:       "*auth.Service",
		Package:    "github.com/archesai/archesai/pkg/auth",
		IsPointer:  true,
		Resolution: "infra.AuthService",
	}

	if dep.Name != "authService" {
		t.Errorf("Name = %q, want %q", dep.Name, "authService")
	}
	if dep.Type != "*auth.Service" {
		t.Errorf("Type = %q, want %q", dep.Type, "*auth.Service")
	}
	if !dep.IsPointer {
		t.Error("IsPointer should be true")
	}
}

func TestHandlerDef(t *testing.T) {
	handler := HandlerDef{
		OperationID: "Login",
		FilePath:    "handlers/login.impl.go",
		HasHandler:  true,
		Constructor: &ConstructorDef{
			Name: "NewLoginHandler",
		},
	}

	if handler.OperationID != "Login" {
		t.Errorf("OperationID = %q, want %q", handler.OperationID, "Login")
	}
	if handler.FilePath != "handlers/login.impl.go" {
		t.Errorf("FilePath = %q, want %q", handler.FilePath, "handlers/login.impl.go")
	}
	if !handler.HasHandler {
		t.Error("HasHandler should be true")
	}
	if handler.Constructor == nil {
		t.Error("Constructor should not be nil")
	}
}

func TestConstructorDef(t *testing.T) {
	constructor := ConstructorDef{
		Name:       "NewTestHandler",
		Parameters: nil,
	}

	if constructor.Name != "NewTestHandler" {
		t.Errorf("Name = %q, want %q", constructor.Name, "NewTestHandler")
	}
}
