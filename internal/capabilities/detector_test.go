package capabilities

import (
	"context"
	"log/slog"
	"testing"
	"time"
)

func TestDetector_Check(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	detector := NewDetector(slog.Default()).
		WithChecker(NewGoChecker(true)).
		WithChecker(NewGitChecker(false))

	result := detector.Check(ctx)

	if len(result.Capabilities) != 2 {
		t.Errorf("expected 2 capabilities, got %d", len(result.Capabilities))
	}

	// Check that Go was detected (it's required for this project)
	var goCap *Capability
	for i := range result.Capabilities {
		if result.Capabilities[i].Name == "Go" {
			goCap = &result.Capabilities[i]
			break
		}
	}

	if goCap == nil {
		t.Fatal("Go capability not found in results")
	}

	if !goCap.Found {
		t.Error("Go should be found (required for this project)")
	}

	if !goCap.Required {
		t.Error("Go should be marked as required")
	}
}

func TestDetector_EmptyCheckers(t *testing.T) {
	ctx := context.Background()

	detector := NewDetector(nil) // Test nil logger handling
	result := detector.Check(ctx)

	if len(result.Capabilities) != 0 {
		t.Errorf("expected 0 capabilities, got %d", len(result.Capabilities))
	}

	if !result.AllRequired {
		t.Error("AllRequired should be true with no checkers")
	}

	if !result.AllFound {
		t.Error("AllFound should be true with no checkers")
	}
}

func TestDefaultDetector(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	detector := DefaultDetector(slog.Default())
	result := detector.Check(ctx)

	// Should have 8 checkers by default
	if len(result.Capabilities) != 8 {
		t.Errorf("expected 8 capabilities, got %d", len(result.Capabilities))
	}

	// Log all detected capabilities
	for _, cap := range result.Capabilities {
		status := "found"
		if !cap.Found {
			status = "not found"
		}
		required := ""
		if cap.Required {
			required = " (required)"
		}
		t.Logf("%s: %s, version=%s%s", cap.Name, status, cap.Version, required)
	}
}

func TestCheckResult_RequiredMissing(t *testing.T) {
	result := CheckResult{
		Capabilities: []Capability{
			{Name: "Go", Required: true, Found: true},
			{Name: "Docker", Required: true, Found: false},
			{Name: "Node.js", Required: false, Found: false},
		},
	}

	missing := result.RequiredMissing()

	if len(missing) != 1 {
		t.Errorf("expected 1 required missing, got %d", len(missing))
	}

	if missing[0].Name != "Docker" {
		t.Errorf("expected Docker to be missing, got %s", missing[0].Name)
	}
}

func TestCheckResult_OptionalMissing(t *testing.T) {
	result := CheckResult{
		Capabilities: []Capability{
			{Name: "Go", Required: true, Found: true},
			{Name: "Docker", Required: true, Found: false},
			{Name: "Node.js", Required: false, Found: false},
			{Name: "Python", Required: false, Found: true},
		},
	}

	missing := result.OptionalMissing()

	if len(missing) != 1 {
		t.Errorf("expected 1 optional missing, got %d", len(missing))
	}

	if missing[0].Name != "Node.js" {
		t.Errorf("expected Node.js to be missing, got %s", missing[0].Name)
	}
}

func TestCheckResult_FoundCapabilities(t *testing.T) {
	result := CheckResult{
		Capabilities: []Capability{
			{Name: "Go", Found: true},
			{Name: "Docker", Found: false},
			{Name: "Node.js", Found: true},
		},
	}

	found := result.FoundCapabilities()

	if len(found) != 2 {
		t.Errorf("expected 2 found capabilities, got %d", len(found))
	}
}
