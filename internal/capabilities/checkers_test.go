package capabilities

import (
	"context"
	"testing"
	"time"
)

func TestGoChecker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	checker := NewGoChecker(true)
	result := checker.Check(ctx)

	if result.Name != "Go" {
		t.Errorf("expected name 'Go', got %q", result.Name)
	}
	if result.Command != "go" {
		t.Errorf("expected command 'go', got %q", result.Command)
	}
	if !result.Required {
		t.Error("expected Required to be true")
	}

	// Go should be installed in the dev environment
	if !result.Found {
		t.Skip("Go not installed, skipping version check")
	}

	if result.Version == "" {
		t.Error("expected version to be detected")
	}
	t.Logf("Detected Go version: %s", result.Version)
}

func TestNodeChecker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	checker := NewNodeChecker(false)
	result := checker.Check(ctx)

	if result.Name != "Node.js" {
		t.Errorf("expected name 'Node.js', got %q", result.Name)
	}
	if result.Command != "node" {
		t.Errorf("expected command 'node', got %q", result.Command)
	}
	if result.Required {
		t.Error("expected Required to be false")
	}

	if result.Found {
		t.Logf("Detected Node.js version: %s", result.Version)
	} else {
		t.Log("Node.js not installed")
	}
}

func TestDockerChecker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	checker := NewDockerChecker(true)
	result := checker.Check(ctx)

	if result.Name != "Docker" {
		t.Errorf("expected name 'Docker', got %q", result.Name)
	}
	if result.Command != "docker" {
		t.Errorf("expected command 'docker', got %q", result.Command)
	}

	if result.Found {
		t.Logf("Detected Docker version: %s", result.Version)
	} else {
		t.Log("Docker not installed")
	}
}

func TestPnpmChecker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	checker := NewPnpmChecker(false)
	result := checker.Check(ctx)

	if result.Name != "pnpm" {
		t.Errorf("expected name 'pnpm', got %q", result.Name)
	}

	if result.Found {
		t.Logf("Detected pnpm version: %s", result.Version)
	} else {
		t.Log("pnpm not installed")
	}
}

func TestPythonChecker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	checker := NewPythonChecker(false)
	result := checker.Check(ctx)

	if result.Name != "Python" {
		t.Errorf("expected name 'Python', got %q", result.Name)
	}
	if result.Command != "python3" {
		t.Errorf("expected command 'python3', got %q", result.Command)
	}

	if result.Found {
		t.Logf("Detected Python version: %s", result.Version)
	} else {
		t.Log("Python not installed")
	}
}

func TestGitChecker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	checker := NewGitChecker(false)
	result := checker.Check(ctx)

	if result.Name != "Git" {
		t.Errorf("expected name 'Git', got %q", result.Name)
	}

	// Git should generally be installed
	if result.Found {
		if result.Version == "" {
			t.Error("expected version to be detected")
		}
		t.Logf("Detected Git version: %s", result.Version)
	} else {
		t.Log("Git not installed")
	}
}

func TestKubectlChecker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	checker := NewKubectlChecker(false)
	result := checker.Check(ctx)

	if result.Name != "kubectl" {
		t.Errorf("expected name 'kubectl', got %q", result.Name)
	}

	if result.Found {
		t.Logf("Detected kubectl version: %s", result.Version)
	} else {
		t.Log("kubectl not installed")
	}
}

func TestHelmChecker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	checker := NewHelmChecker(false)
	result := checker.Check(ctx)

	if result.Name != "Helm" {
		t.Errorf("expected name 'Helm', got %q", result.Name)
	}

	if result.Found {
		t.Logf("Detected Helm version: %s", result.Version)
	} else {
		t.Log("Helm not installed")
	}
}
