# ==========================================
# ARCHESAI MAKEFILE
# ==========================================

# ------------------------------------------
# Configuration
# ------------------------------------------

# Build Configuration
CPUS := $(shell nproc 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo 4)
MAKEFLAGS += -j$(CPUS) --no-print-directory

# Terminal Colors
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
CYAN := \033[0;36m
RED := \033[0;31m
GRAY := \033[0;90m
NC := \033[0m # No Color

# ------------------------------------------
# Primary Commands
# ------------------------------------------

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-25s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

# ------------------------------------------
# Build Commands
# ------------------------------------------

BINARY_PATH := ./bin

.PHONY: build
build: ## Build archesai binary
	@go build -o $(BINARY_PATH)/archesai ./cmd/archesai
	@du -sh $(BINARY_PATH)/archesai

.PHONY: build-studio
build-studio: ## Build studio app
	@cd ./apps/studio && go run ../../cmd/archesai build -c arches.yaml

.PHONY: build-docs
build-docs: prepare-docs ## Build documentation site
	@pnpm -F @archesai/docs build

# ------------------------------------------
# Run Commands (Production-like)
# ------------------------------------------

.PHONY: run-studio
run-studio: build ## Run the API server (production mode)
	@cd ./apps/studio && go run ../../cmd/archesai run -c arches.yaml

.PHONY: run-docs
run-docs: prepare-docs ## Run documentation site (production build)
	@pnpm -F @archesai/docs start

# ------------------------------------------
# Development Commands (Hot Reload)
# ------------------------------------------

.PHONY: dev-studio
dev-studio: ## Run API server with hot reload
	@cd ./apps/studio && go run ../../cmd/archesai dev -c arches.yaml

.PHONY: dev-docs
dev-docs: prepare-docs ## Run documentation with hot reload
	@pnpm -F @archesai/docs dev

# ------------------------------------------
# Workflow Commands
# ------------------------------------------

.PHONY: workflow-deploy-docs
workflow-deploy-docs: ## Manually trigger documentation deployment to GitHub Pages
	@which gh > /dev/null || (echo -e "$(RED)✗ Please install GitHub CLI first$(NC)" && exit 1)
	@gh workflow run deploy-docs.yaml

# ------------------------------------------
# Generate Commands
# ------------------------------------------

generate: ## Regenerate code
	@$(MAKE) generate-packages
	@$(MAKE) generate-examples generate-studio

.PHONY: generate-packages
generate-packages: ## Generate all packages
	@go generate ./pkg/config ./...

.PHONY:
generate-studio: ## Generate codegen
	@cd apps/studio && go run ../../cmd/archesai generate

.PHONY: generate-examples ## Generate example configurations
generate-examples:
	@cd examples/basic && go run ../../cmd/archesai generate
	@cd examples/authentication && go run ../../cmd/archesai generate

# ------------------------------------------
# Test Commands
# ------------------------------------------

.PHONY: test
test: ## Run all tests
	@go test -race -cover ./...

.PHONY: test-verbose
test-verbose: ## Run all tests with verbose output
	@go test -v -race -cover ./...

.PHONY: test-short
test-short: ## Run short tests only (skip integration tests)
	@go test -short -cover ./...

.PHONY: test-coverage
test-coverage: ## Generate test coverage report
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out

.PHONY: test-coverage-html
test-coverage-html: test-coverage ## Generate HTML coverage report
	@go tool cover -html=coverage.out -o coverage.html

.PHONY: test-bench
test-bench: ## Run benchmark tests
	@go test -bench=. -benchmem ./...

.PHONY: test-watch
test-watch: ## Run tests in watch mode (requires fswatch)
	@which fswatch > /dev/null || (echo "Please install fswatch first" && exit 1)
	@fswatch -o . -e ".*" -i "\\.go$$" | xargs -n1 -I{} sh -c 'clear && make test'


# -------------------------------------------
# GitHub Workflow Commands
# -------------------------------------------

.PHONY: run-workflow
run-workflow: ## Run GitHub workflow locally with act (usage: make run-workflow workflow=update-docs)
	@if [ -z "$(workflow)" ]; then \
		echo -e "$(RED)✗ Please specify a workflow name$(NC)"; \
		echo -e "$(BLUE)Usage: make run-workflow workflow=<workflow-name>$(NC)"; \
		echo -e "$(BLUE)Example: make run-workflow workflow=update-docs$(NC)"; \
		echo -e "$(BLUE)Available workflows:$(NC)"; \
		ls -1 .github/workflows/*.y*ml | sed 's|.github/workflows/||' | sed 's|\.y.*ml||' | sed 's|^|  - |'; \
		exit 1; \
	fi
	@which act > /dev/null || (echo -e "$(RED)✗ Please install act first: https://github.com/nektos/act$(NC)" && exit 1)
	@act -W .github/workflows/$(workflow).yaml

.PHONY: list-workflows
list-workflows: ## List all available GitHub workflows
	@ls -1 .github/workflows/*.y*ml | sed 's|.github/workflows/||' | sed 's|\.y.*ml||' | sed 's|^|  - |'

# ------------------------------------------
# Lint Commands
# ------------------------------------------

.PHONY: lint
lint: lint-go lint-ts lint-openapi lint-docs ## Run all linters

.PHONY: lint-go
lint-go: ## Run Go linter
	@OUTPUT=$$(golangci-lint run --color always ./... 2>&1); \
	if [ $$? -ne 0 ]; then \
		echo -e "$(RED)✗ Go linting failed$(NC)"; \
		echo "$$OUTPUT"; \
		exit 1; \
	elif echo "$$OUTPUT" | grep -v "^0 issues" | grep -q .; then \
		echo "$$OUTPUT"; \
	fi

.PHONY: lint-ts
lint-ts: lint-typecheck ## Run Node.js linter (includes typecheck)
	@OUTPUT=$$(pnpm biome check --fix 2>&1); \
	if [ $$? -ne 0 ]; then \
		echo -e "$(RED)✗ Node.js linting failed$(NC)"; \
		echo "$$OUTPUT"; \
		exit 1; \
	elif echo "$$OUTPUT" | grep -v "No fixes applied" | grep -q .; then \
		echo "$$OUTPUT"; \
	fi

.PHONY: lint-typecheck
lint-typecheck: ## Run TypeScript type checking
	@pnpm tsc --build --emitDeclarationOnly

.PHONY: lint-openapi
lint-openapi: ## Lint OpenAPI specification
	@cd ./apps/studio && go run ../../cmd/archesai spec lint -c arches.yaml

.PHONY: lint-docs
lint-docs: ## Lint documentation with markdownlint
	@pnpm markdownlint --fix 'docs/**/*.md' --config .markdownlint.json

# ------------------------------------------
# Format Commands
# ------------------------------------------

.PHONY: format
format: format-go format-prettier format-ts ## Format all code

.PHONY: format-go
format-go: ## Format Go code
	@golangci-lint run --fix

.PHONY: format-prettier
format-prettier: ## Format code with Prettier
	@pnpm prettier --list-different --write --log-level warn .

.PHONY: format-ts
format-ts: ## Format Node.js/TypeScript code
	@pnpm biome format --fix

# ------------------------------------------
# Clean Commands
# ------------------------------------------

.PHONY: clean
clean: clean-ts clean-go clean-generated ## Clean all build artifacts

.PHONY: clean-ts
clean-ts: ## Clean distribution builds
	@pnpm -r exec sh -c 'rm -rf .cache .tanstack dist .nitro .output'

.PHONY: clean-go
clean-go: ## Clean Go build artifacts
	@go clean -testcache
	@rm -rf $(BINARY_PATH) coverage.out coverage.html

.PHONY: clean-generated
clean-generated: ## Clean all generated code
	@find . -type f -name "*.gen.*" -not -path "./internal/codegen/tmpl/*" -exec rm -f {} +
	@find . -type d -empty -delete 2>/dev/null || true

.PHONY: prepare-docs
prepare-docs: build-studio ## Copy markdown docs to apps/docs/docs FIXME: bundle-openapi
	@mkdir -p ./apps/docs/apis ./apps/docs/pages/documentation
	@cp ./api/openapi.bundled.yaml ./apps/docs/apis/openapi.yaml
	@cp -r ./docs/** ./apps/docs/pages

# ------------------------------------------
# Dependency Commands
# ------------------------------------------

.PHONY: deps
deps: deps-go deps-ts ## Install all dependencies

.PHONY: deps-go
deps-go: ## Install Go dependencies and tools
	@go mod download

.PHONY: deps-ts
deps-ts: ## Install Node.js dependencies
	@pnpm install

.PHONY: deps-update
deps-update: deps-update-go deps-update-ts ## Update all dependencies

.PHONY: deps-update-go
deps-update-go: ## Update Go dependencies
	@go get -u ./... 2>&1 | { grep -v "warning: ignoring symlink" || true; }
	@go mod tidy

.PHONY: deps-update-ts
deps-update-ts: ## Update Node.js dependencies
	@pnpm update -r --latest

.PHONY: vendor-jsonschema
vendor-jsonschema: ## Vendor google/jsonschema-go into internal/jsonschema
	@rm -rf internal/jsonschema
	@mkdir -p internal/jsonschema
	@git clone --depth 1 https://github.com/google/jsonschema-go.git /tmp/jsonschema-go
	@find /tmp/jsonschema-go/jsonschema -maxdepth 1 -name '*.go' ! -name '*_test.go' -exec cp {} internal/jsonschema/ \;
	@rm -rf /tmp/jsonschema-go

# ------------------------------------------
# Container Executor Commands
# ------------------------------------------

.PHONY: build-runners
build-runners: build-runner-python build-runner-node build-runner-go ## Build all runner containers

.PHONY: build-runner-python
build-runner-python: ## Build Python runner container
	@docker build -t archesai/runner-python:latest deployments/containers/runners/python

.PHONY: build-runner-node
build-runner-node: ## Build Node runner base container
	@docker build -t archesai/runner-node:latest deployments/containers/runners/node

.PHONY: build-runner-go
build-runner-go: ## Build Go runner container
	@docker build -t archesai/runner-go:latest deployments/containers/runners/go

# ------------------------------------------
# Kubernetes Commands
# ------------------------------------------

.PHONY: k8s-cluster-start
k8s-cluster-start: ## Start k3d cluster
	@k3d cluster create tower --config deployments/k3d/k3d.yaml

.PHONY: k8s-cluster-stop
k8s-cluster-stop: ## Stop k3d cluster
	@k3d cluster delete -a

.PHONY: k8s-helm-install
k8s-helm-install: ## Deploy with Helm
	@helm install dev deployments/helm/arches -f deployments/helm/dev-overrides.yaml

.PHONY: k8s-helm-upgrade
k8s-helm-upgrade: ## Upgrade Helm deployment
	@helm upgrade dev deployments/helm/arches -f deployments/helm/dev-overrides.yaml

.PHONY: k8s-helm-uninstall
k8s-helm-uninstall: ## Upgrade Helm deployment
	@helm uninstall dev

# ------------------------------------------
# Skaffold Commands
# ------------------------------------------

.PHONY: skaffold-dev
skaffold-dev: ## Run with Skaffold in dev mode
	@skaffold dev --default-repo registry.localhost:5000 --profile dev

.PHONY: skaffold-run
skaffold-run: ## Deploy with Skaffold
	@skaffold run

.PHONY: skaffold-delete
skaffold-delete: ## Delete Skaffold deployment
	@skaffold delete --profile dev

# ------------------------------------------
# Release Management
# ------------------------------------------

.PHONY: release-check
release-check: ## Check if ready for release
	@$(MAKE) generate
	@$(MAKE) test-short
	@$(MAKE) lint

.PHONY: release-snapshot
release-snapshot: ## Create a snapshot release (test GoReleaser config)
	@if ! command -v goreleaser >/dev/null 2>&1; then \
		echo -e "$(RED)✗ GoReleaser not found. Install with: go install github.com/goreleaser/goreleaser@latest$(NC)"; \
		exit 1; \
	fi
	@$(MAKE) release-check
	@goreleaser release --snapshot --clean

.PHONY: release-test
release-test: ## Test release configuration without publishing
	@if ! command -v goreleaser >/dev/null 2>&1; then \
		echo -e "$(RED)✗ GoReleaser not found. Install with: go install github.com/goreleaser/goreleaser@latest$(NC)"; \
		exit 1; \
	fi
	@goreleaser check
	@goreleaser build --snapshot --clean

.PHONY: release-tag
release-tag: ## Create and push a new release tag (usage: make release-tag VERSION=v1.0.0)
	@if [ -z "$(VERSION)" ]; then \
		echo -e "$(RED)✗ VERSION is required. Usage: make release-tag VERSION=v1.0.0$(NC)"; \
		exit 1; \
	fi
	@if git rev-parse $(VERSION) >/dev/null 2>&1; then \
		echo -e "$(RED)✗ Tag $(VERSION) already exists$(NC)"; \
		exit 1; \
	fi
	@$(MAKE) release-check
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@echo -e "$(BLUE)Pushing tag to origin...$(NC)"
	@git push origin $(VERSION)
	@echo -e "$(GREEN)✓ Tag $(VERSION) created and pushed. GitHub Actions will handle the release.$(NC)"

.PHONY: release-draft
release-draft: ## Create a draft release on GitHub (requires gh CLI)
	@echo -e "$(YELLOW)━━━ Creating Draft Release ━━━$(NC)"
	@if ! command -v gh >/dev/null 2>&1; then \
		echo -e "$(RED)✗ GitHub CLI not found. Install from: https://cli.github.com/$(NC)"; \
		exit 1; \
	fi
	@if [ -z "$(VERSION)" ]; then \
		echo -e "$(RED)✗ VERSION is required. Usage: make release-draft VERSION=v1.0.0$(NC)"; \
		exit 1; \
	fi
	@$(MAKE) release-check
	@echo -e "$(BLUE)Creating draft release $(VERSION)...$(NC)"
	@gh release create $(VERSION) --draft --title "Release $(VERSION)" --notes-file CHANGELOG.md || \
		gh release create $(VERSION) --draft --title "Release $(VERSION)" --notes "Release $(VERSION)"
	@echo -e "$(GREEN)✓ Draft release created. Edit at: https://github.com/archesai/archesai/releases$(NC)"

.PHONY: release-edge-local
release-edge-local: ## Test edge release workflow locally
	@echo -e "$(YELLOW)━━━ Testing Edge Release Locally ━━━$(NC)"
	@$(MAKE) release-snapshot
	@echo -e "$(GREEN)✓ Edge release test complete. Check ./dist/ directory$(NC)"

.PHONY: release-nightly-local
release-nightly-local: ## Test nightly release workflow locally
	@$(MAKE) release-snapshot

.PHONY: release-info
release-info: ## Show release information and next steps
	@echo -e "$(CYAN)━━━ Arches Release Information ━━━$(NC)"
	@echo -e "$(BLUE)Release Types:$(NC)"
	@echo -e "  • $(GREEN)Stable$(NC)    - Tagged releases (v1.0.0) via GitHub Actions"
	@echo -e "  • $(YELLOW)Nightly$(NC)   - Daily builds from main branch (2 AM UTC)"
	@echo -e "  • $(RED)Edge$(NC)       - Every push to main branch"
	@echo ""
	@echo -e "$(BLUE)Available Commands:$(NC)"
	@echo -e "  • $(CYAN)make release-check$(NC)         - Verify readiness for release"
	@echo -e "  • $(CYAN)make release-test$(NC)          - Test GoReleaser configuration"
	@echo -e "  • $(CYAN)make release-snapshot$(NC)      - Create local snapshot build"
	@echo -e "  • $(CYAN)make release-tag VERSION=v1.0.0$(NC) - Create and push release tag"
	@echo -e "  • $(CYAN)make release-draft VERSION=v1.0.0$(NC) - Create GitHub draft release"
	@echo ""
	@echo -e "$(BLUE)Release Artifacts:$(NC)"
	@echo -e "  • $(GREEN)Binaries$(NC)  - Cross-platform executables"
	@echo -e "  • $(GREEN)Docker$(NC)    - Multi-arch container images"
	@echo -e "  • $(GREEN)Checksums$(NC) - SHA256 verification files"
	@echo -e "  • $(GREEN)Archives$(NC)  - Standard & full (with platform assets)"
	@echo ""
	@echo -e "$(BLUE)Next Steps:$(NC)"
	@echo -e "  1. Run $(CYAN)make release-check$(NC) to verify readiness"
	@echo -e "  2. Run $(CYAN)make release-tag VERSION=v1.0.0$(NC) to create a release"
	@echo -e "  3. GitHub Actions will automatically build and publish"
	@echo ""

# ------------------------------------------
# Shortcuts
# ------------------------------------------

.PHONY: g
g: generate ## Shortcut for generate

.PHONY: t
t: test ## Shortcut for test

.PHONY: f
f: format ## Shortcut for format
