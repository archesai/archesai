package generators

import (
	"testing"

	"github.com/archesai/archesai/internal/spec"
)

func TestBuildInternalPackages(t *testing.T) {
	pkgs := BuildInternalPackages([]string{"auth", "user", "billing"})

	if len(pkgs) != 3 {
		t.Fatalf("expected 3 packages, got %d", len(pkgs))
	}

	// Should be sorted
	expected := []string{"auth", "billing", "user"}
	for i, exp := range expected {
		if pkgs[i].Name != exp {
			t.Errorf("expected package %d to be %q, got %q", i, exp, pkgs[i].Name)
		}
		if pkgs[i].Alias != exp {
			t.Errorf("expected alias %d to be %q, got %q", i, exp, pkgs[i].Alias)
		}
		if pkgs[i].ImportPath == "" {
			t.Errorf("expected import path for package %d to be non-empty", i)
		}
	}
}

func TestExtractRepositoryInfo(t *testing.T) {
	operations := []spec.Operation{
		{ID: "GetUser", Tag: "User", Method: "GET"},
		{ID: "CreateUser", Tag: "User", Method: "POST"},
		{ID: "GetOrder", Tag: "Order", Method: "GET"},
		{ID: "DeleteOrder", Tag: "Order", Method: "DELETE"},
		{ID: "CustomOp", Tag: "Custom", Method: "POST", CustomHandler: true},
	}

	repos, needsPublisher := ExtractRepositoryInfo(operations)

	// Should have User and Order (Custom is excluded because CustomHandler=true)
	if len(repos) != 2 {
		t.Fatalf("expected 2 repos, got %d: %v", len(repos), repos)
	}

	// Should be sorted
	if repos[0] != "Order" || repos[1] != "User" {
		t.Errorf("expected repos to be [Order, User], got %v", repos)
	}

	// Should need publisher (POST and DELETE operations exist)
	if !needsPublisher {
		t.Error("expected needsPublisher to be true")
	}
}

func TestExtractRepositoryInfo_NoPublisher(t *testing.T) {
	operations := []spec.Operation{
		{ID: "GetUser", Tag: "User", Method: "GET"},
		{ID: "ListUsers", Tag: "User", Method: "GET"},
	}

	_, needsPublisher := ExtractRepositoryInfo(operations)

	if needsPublisher {
		t.Error("expected needsPublisher to be false for GET-only operations")
	}
}

func TestNeedsPublisher(t *testing.T) {
	tests := []struct {
		name       string
		operations []spec.Operation
		expected   bool
	}{
		{
			name: "GET only",
			operations: []spec.Operation{
				{Method: "GET"},
				{Method: "GET"},
			},
			expected: false,
		},
		{
			name: "has POST",
			operations: []spec.Operation{
				{Method: "GET"},
				{Method: "POST"},
			},
			expected: true,
		},
		{
			name: "has DELETE",
			operations: []spec.Operation{
				{Method: "GET"},
				{Method: "DELETE"},
			},
			expected: true,
		},
		{
			name: "custom handler POST",
			operations: []spec.Operation{
				{Method: "GET"},
				{Method: "POST", CustomHandler: true},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NeedsPublisher(tt.operations)
			if got != tt.expected {
				t.Errorf("NeedsPublisher() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestAppContext_ShouldSkip(t *testing.T) {
	tests := []struct {
		name     string
		actx     AppContext
		expected bool
	}{
		{
			name:     "empty context",
			actx:     AppContext{},
			expected: true,
		},
		{
			name: "standalone",
			actx: AppContext{
				IsStandalone: true,
			},
			expected: false,
		},
		{
			name: "composition",
			actx: AppContext{
				IsComposition: true,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.actx.ShouldSkip()
			if got != tt.expected {
				t.Errorf("ShouldSkip() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestHasCustomHandlers(t *testing.T) {
	tests := []struct {
		name       string
		operations []spec.Operation
		expected   bool
	}{
		{
			name:       "empty operations",
			operations: []spec.Operation{},
			expected:   false,
		},
		{
			name: "no custom handlers",
			operations: []spec.Operation{
				{ID: "GetUser", Method: "GET", CustomHandler: false},
				{ID: "CreateUser", Method: "POST", CustomHandler: false},
			},
			expected: false,
		},
		{
			name: "has custom handler",
			operations: []spec.Operation{
				{ID: "GetUser", Method: "GET", CustomHandler: false},
				{ID: "Login", Method: "POST", CustomHandler: true},
			},
			expected: true,
		},
		{
			name: "all custom handlers",
			operations: []spec.Operation{
				{ID: "Login", Method: "POST", CustomHandler: true},
				{ID: "Register", Method: "POST", CustomHandler: true},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HasCustomHandlers(tt.operations)
			if got != tt.expected {
				t.Errorf("HasCustomHandlers() = %v, want %v", got, tt.expected)
			}
		})
	}
}
