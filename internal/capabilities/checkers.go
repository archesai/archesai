package capabilities

import (
	"context"
	"os/exec"
	"regexp"
	"strings"
)

// baseChecker provides common functionality for capability checkers.
type baseChecker struct {
	name        string
	command     string
	versionArgs []string
	versionRe   *regexp.Regexp
	required    bool
	description string
}

func (b *baseChecker) Check(ctx context.Context) Capability {
	result := Capability{
		Name:        b.name,
		Command:     b.command,
		Required:    b.required,
		Description: b.description,
	}

	// Check if command exists
	path, err := exec.LookPath(b.command)
	if err != nil {
		result.Found = false
		result.Error = err
		return result
	}

	result.Found = true

	// Get version if version args are specified
	if len(b.versionArgs) > 0 {
		cmd := exec.CommandContext(ctx, path, b.versionArgs...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			result.Error = err
			return result
		}

		version := strings.TrimSpace(string(output))
		if b.versionRe != nil {
			if matches := b.versionRe.FindStringSubmatch(version); len(matches) > 1 {
				version = matches[1]
			}
		}
		result.Version = version
	}

	return result
}

// GoChecker checks for Go installation.
type GoChecker struct {
	baseChecker
}

// NewGoChecker creates a new Go capability checker.
func NewGoChecker(required bool) *GoChecker {
	return &GoChecker{
		baseChecker: baseChecker{
			name:        "Go",
			command:     "go",
			versionArgs: []string{"version"},
			versionRe:   regexp.MustCompile(`go(\d+\.\d+(?:\.\d+)?)`),
			required:    required,
			description: "Go programming language runtime",
		},
	}
}

// NodeChecker checks for Node.js installation.
type NodeChecker struct {
	baseChecker
}

// NewNodeChecker creates a new Node.js capability checker.
func NewNodeChecker(required bool) *NodeChecker {
	return &NodeChecker{
		baseChecker: baseChecker{
			name:        "Node.js",
			command:     "node",
			versionArgs: []string{"--version"},
			versionRe:   regexp.MustCompile(`v?(\d+\.\d+\.\d+)`),
			required:    required,
			description: "Node.js JavaScript runtime",
		},
	}
}

// DockerChecker checks for Docker installation.
type DockerChecker struct {
	baseChecker
}

// NewDockerChecker creates a new Docker capability checker.
func NewDockerChecker(required bool) *DockerChecker {
	return &DockerChecker{
		baseChecker: baseChecker{
			name:        "Docker",
			command:     "docker",
			versionArgs: []string{"--version"},
			versionRe:   regexp.MustCompile(`Docker version (\d+\.\d+\.\d+)`),
			required:    required,
			description: "Docker container runtime",
		},
	}
}

// PnpmChecker checks for pnpm installation.
type PnpmChecker struct {
	baseChecker
}

// NewPnpmChecker creates a new pnpm capability checker.
func NewPnpmChecker(required bool) *PnpmChecker {
	return &PnpmChecker{
		baseChecker: baseChecker{
			name:        "pnpm",
			command:     "pnpm",
			versionArgs: []string{"--version"},
			versionRe:   regexp.MustCompile(`(\d+\.\d+\.\d+)`),
			required:    required,
			description: "Fast, disk space efficient package manager",
		},
	}
}

// PythonChecker checks for Python installation.
type PythonChecker struct {
	baseChecker
}

// NewPythonChecker creates a new Python capability checker.
func NewPythonChecker(required bool) *PythonChecker {
	return &PythonChecker{
		baseChecker: baseChecker{
			name:        "Python",
			command:     "python3",
			versionArgs: []string{"--version"},
			versionRe:   regexp.MustCompile(`Python (\d+\.\d+\.\d+)`),
			required:    required,
			description: "Python programming language runtime",
		},
	}
}

// GitChecker checks for Git installation.
type GitChecker struct {
	baseChecker
}

// NewGitChecker creates a new Git capability checker.
func NewGitChecker(required bool) *GitChecker {
	return &GitChecker{
		baseChecker: baseChecker{
			name:        "Git",
			command:     "git",
			versionArgs: []string{"--version"},
			versionRe:   regexp.MustCompile(`git version (\d+\.\d+\.\d+)`),
			required:    required,
			description: "Distributed version control system",
		},
	}
}

// KubectlChecker checks for kubectl installation.
type KubectlChecker struct {
	baseChecker
}

// NewKubectlChecker creates a new kubectl capability checker.
func NewKubectlChecker(required bool) *KubectlChecker {
	return &KubectlChecker{
		baseChecker: baseChecker{
			name:        "kubectl",
			command:     "kubectl",
			versionArgs: []string{"version", "--client"},
			versionRe:   regexp.MustCompile(`Client Version: v?(\d+\.\d+\.\d+)`),
			required:    required,
			description: "Kubernetes command-line tool",
		},
	}
}

// HelmChecker checks for Helm installation.
type HelmChecker struct {
	baseChecker
}

// NewHelmChecker creates a new Helm capability checker.
func NewHelmChecker(required bool) *HelmChecker {
	return &HelmChecker{
		baseChecker: baseChecker{
			name:        "Helm",
			command:     "helm",
			versionArgs: []string{"version", "--short"},
			versionRe:   regexp.MustCompile(`v?(\d+\.\d+\.\d+)`),
			required:    required,
			description: "Kubernetes package manager",
		},
	}
}
