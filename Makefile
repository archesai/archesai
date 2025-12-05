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
	@echo -e "$(YELLOW)▶ Building archesai...$(NC)"
	@go build -o $(BINARY_PATH)/archesai  cmd/archesai
	@echo -e "$(GREEN)✓ archesai built: $(BINARY_PATH)/archesai $(NC)"

.PHONY: build-studio
build-studio: ## Build studio app
	@echo -e "$(YELLOW)▶ Building studio app...$(NC)"
	@echo -e "$(RED)▶ NOT IMPLEMENTED ...$(NC)"

.PHONY: build-docs
build-docs: prepare-docs ## Build documentation site
	@echo -e "$(YELLOW)▶ Building documentation site...$(NC)"
	@pnpm -F @archesai/docs build
	@echo -e "$(GREEN)✓ Documentation built in apps/docs/build/$(NC)"

# ------------------------------------------
# Run Commands (Production-like)
# ------------------------------------------

.PHONY: run-studio
run-studio: build ## Run the API server (production mode)
	@echo -e "$(YELLOW)▶ Starting API server...$(NC)"
	@go run ./cmd/archesai api

.PHONY: run-docs
run-docs: prepare-docs ## Run documentation site (production build)
	@echo -e "$(YELLOW)▶ Starting documentation server...$(NC)"
	@pnpm -F @archesai/docs start

# ------------------------------------------
# Development Commands (Hot Reload)
# ------------------------------------------

.PHONY: dev-studio
dev-studio: ## Run API server with hot reload
	@echo -e "$(YELLOW)▶ Starting API server with hot reload on port 3001...$(NC)"
	@echo -e "$(YELLOW)Press Ctrl+C to stop$(NC)"
	@go tool -modfile=tools.mod air

.PHONY: dev-docs
dev-docs: prepare-docs ## Run documentation with hot reload
	@echo -e "$(YELLOW)▶ Starting documentation with hot reload...$(NC)"
	@pnpm -F @archesai/docs dev

# ------------------------------------------
# Deployment Commands
# ------------------------------------------

.PHONY: deploy-docs
deploy-docs: ## Manually trigger documentation deployment to GitHub Pages
	@echo -e "$(YELLOW)▶ Triggering documentation deployment...$(NC)"
	@which gh > /dev/null || (echo -e "$(RED)✗ Please install GitHub CLI first$(NC)" && exit 1)
	@gh workflow run deploy-docs.yaml
	@echo -e "$(GREEN)✓ Documentation deployment triggered!$(NC)"
	@echo -e "$(BLUE)Monitor progress: gh run list --workflow=deploy-docs.yaml$(NC)"

# ------------------------------------------
# Generate Commands
# ------------------------------------------

generate: ## Regenerate example configurations
	@$(MAKE) generate-packages
	@$(MAKE) generate-examples generate-studio
	@$(MAKE) format-prettier

.PHONY: generate-packages
generate-packages: ## Generate all packages
	@go generate ./...

.PHONY:
generate-studio: ## Generate codegen
	@go run ./cmd/archesai generate --spec ./api/openapi.yaml --output ./apps/studio --pretty

.PHONY: generate-examples ## Generate example configurations
generate-examples:
	@go run ./cmd/archesai generate --spec ./examples/basic/spec/openapi.yaml --output ./examples/basic --pretty
	@go run ./cmd/archesai generate --spec ./examples/authentication/spec/openapi.yaml --output ./examples/authentication --pretty

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
	@echo -e "$(YELLOW)▶ Running workflow: $(workflow)...$(NC)"
	@which act > /dev/null || (echo -e "$(RED)✗ Please install act first: https://github.com/nektos/act$(NC)" && exit 1)
	@act -W .github/workflows/$(workflow).yaml
	@echo -e "$(GREEN)✓ Workflow execution complete!$(NC)"

.PHONY: list-workflows
list-workflows: ## List all available GitHub workflows
	@echo -e "$(BLUE)Available workflows:$(NC)"
	@ls -1 .github/workflows/*.y*ml | sed 's|.github/workflows/||' | sed 's|\.y.*ml||' | sed 's|^|  - |'

# ------------------------------------------
# Lint Commands
# ------------------------------------------

.PHONY: lint
lint: lint-go lint-ts lint-docs ## Run all linters
	@$(MAKE) lint-openapi
	@echo -e "$(GREEN)✓ All linting complete!$(NC)"

.PHONY: lint-go
lint-go: ## Run Go linter
	@echo -e "$(YELLOW)▶ Running Go linter...$(NC)"
	@OUTPUT=$$(golangci-lint run --color always ./... 2>&1); \
	if [ $$? -ne 0 ]; then \
		echo -e "$(RED)✗ Go linting failed$(NC)"; \
		echo "$$OUTPUT"; \
		exit 1; \
	elif echo "$$OUTPUT" | grep -v "^0 issues" | grep -q .; then \
		echo "$$OUTPUT"; \
	fi
	@echo -e "$(GREEN)✓ Go linting complete!$(NC)"

.PHONY: lint-ts
lint-ts: lint-typecheck ## Run Node.js linter (includes typecheck)
	@echo -e "$(YELLOW)▶ Running Node.js linter...$(NC)"
	@OUTPUT=$$(pnpm biome check --fix 2>&1); \
	if [ $$? -ne 0 ]; then \
		echo -e "$(RED)✗ Node.js linting failed$(NC)"; \
		echo "$$OUTPUT"; \
		exit 1; \
	elif echo "$$OUTPUT" | grep -v "No fixes applied" | grep -q .; then \
		echo "$$OUTPUT"; \
	fi
	@echo -e "$(GREEN)✓ Node.js linting complete!$(NC)"

.PHONY: lint-openapi
lint-openapi: ## Lint OpenAPI specification
	@go run ./cmd/archesai spec lint --spec ./api/openapi.yaml

.PHONY: lint-typecheck
lint-typecheck: ## Run TypeScript type checking
	@pnpm tsc --build --emitDeclarationOnly

.PHONY: lint-docs
lint-docs: ## Lint documentation with markdownlint
	@echo -e "$(YELLOW)▶ Linting documentation...$(NC)"
	@pnpm markdownlint --fix 'docs/**/*.md' --config .markdownlint.json
	@echo -e "$(GREEN)✓ Documentation linting complete!$(NC)"

# ------------------------------------------
# Format Commands
# ------------------------------------------

.PHONY: format
format: format-go format-ts format-prettier ## Format all code
	@echo -e "$(GREEN)✓ All code formatted!$(NC)"

.PHONY: format-go
format-go: ## Format Go code
	@echo -e "$(YELLOW)▶ Formatting Go code...$(NC)"
	@golangci-lint run --fix
	@echo -e "$(GREEN)✓ Go code formatted!$(NC)"

.PHONY: format-prettier
format-prettier: ## Format code with Prettier
	@echo -e "$(YELLOW)▶ Formatting code with Prettier...$(NC)"
	@pnpm prettier --list-different --write --log-level warn .
	@echo -e "$(GREEN)✓ Code formatted with Prettier!$(NC)"

.PHONY: format-ts
format-ts: ## Format Node.js/TypeScript code
	@echo -e "$(YELLOW)▶ Formatting Node.js code...$(NC)"
	@pnpm biome format --fix
	@echo -e "$(GREEN)✓ Node.js code formatted!$(NC)"

# ------------------------------------------
# Clean Commands
# ------------------------------------------

.PHONY: clean
clean: clean-ts clean-go clean-generated clean-test ## Clean all build artifacts
	@echo -e "$(GREEN)✓ Clean complete!$(NC)"

.PHONY: clean-ts
clean-ts: ## Clean distribution builds
	@echo -e "$(YELLOW)▶ Cleaning distribution builds...$(NC)"
	@pnpm -r exec sh -c 'rm -rf .cache .tanstack dist .nitro .output'
	@echo -e "$(GREEN)✓ Distribution builds cleaned!$(NC)"

.PHONY: clean-go
clean-go: ## Clean Go build artifacts
	@rm -rf $(BINARY_PATH)

.PHONY: clean-generated
clean-generated: ## Clean all generated code
	@find . -type f -name "*.gen.*" -not -path "./internal/codegen/tmpl/*" -not -path "./pkg/auth/gen/repositories/*" -not -path "./pkg/auth/gen/models/*" -exec rm -f {} +
	@find . -type d -empty -delete 2>/dev/null || true
	@rm -rf ./pkg/client/src/generated
	@rm -f ./deployments/helm-minimal/values.schema.json
	@echo -e "$(GREEN)✓ Generated code cleaned!$(NC)"

.PHONY: clean-test
clean-test: ## Clean test cache and coverage files
	@echo -e "$(YELLOW)▶ Cleaning test cache...$(NC)"
	@go clean -testcache
	@rm -f coverage.out coverage.html
	@echo -e "$(GREEN)✓ Test cache cleaned!$(NC)"

.PHONY: clean-deps
clean-deps: clean-ts-deps clean-go-deps ## Clean all dependencies
	@echo -e "$(GREEN)✓ All dependencies cleaned!$(NC)"

.PHONY: clean-ts-deps
clean-ts-deps: ## Clean Node.js dependencies
	@echo -e "$(YELLOW)▶ Cleaning Node.js dependencies...$(NC)"
	@pnpm -r exec sh -c 'rm -rf node_modules pnpm-lock.yaml'
	@echo -e "$(GREEN)✓ Node.js dependencies cleaned!$(NC)"

.PHONY: clean-go-deps
clean-go-deps: ## Clean Go module cache
	@echo -e "$(YELLOW)▶ Cleaning Go module cache...$(NC)"
	@go clean -modcache
	@echo -e "$(GREEN)✓ Go module cache cleaned!$(NC)"

.PHONY: prepare-docs
prepare-docs: bundle-openapi ## Copy markdown docs to apps/docs/docs FIXME: bundle-openapi
	@echo -e "$(YELLOW)▶ Copying markdown docs to apps/docs...$(NC)"
	@mkdir -p ./apps/docs/apis ./apps/docs/pages/documentation
	@cp ./api/openapi.bundled.yaml ./apps/docs/apis/openapi.yaml
	@cp -r ./docs/** ./apps/docs/pages
	@echo -e "$(GREEN)✓ Docs copied!$(NC)"

# ------------------------------------------
# API/OpenAPI Commands
# ------------------------------------------

.PHONY: bundle-openapi
bundle-openapi:  ## Bundle OpenAPI into single file
	@go run cmd/archesai generate --spec ./api/openapi.yaml --output ./apps/studio --only bundle

# ------------------------------------------
# Dependency Commands
# ------------------------------------------

.PHONY: deps
deps: deps-go deps-ts ## Install all dependencies
	@echo -e "$(GREEN)✓ All dependencies installed!$(NC)"

.PHONY: deps-go
deps-go: ## Install Go dependencies and tools
	@echo -e "$(YELLOW)▶ Installing Go dependencies...$(NC)"
	@go mod download
	@echo -e "$(GREEN)✓ Go dependencies installed!$(NC)"

.PHONY: deps-ts
deps-ts: ## Install Node.js dependencies
	@echo -e "$(YELLOW)▶ Installing Node.js dependencies...$(NC)"
	@pnpm install
	@echo -e "$(GREEN)✓ Node.js dependencies installed!$(NC)"

.PHONY: deps-update
deps-update: deps-update-go deps-update-ts ## Update all dependencies
	@echo -e "$(GREEN)✓ All dependencies updated!$(NC)"

.PHONY: deps-update-go
deps-update-go: ## Update Go dependencies
	@echo -e "$(YELLOW)▶ Updating Go dependencies...$(NC)"
	@go get -u ./... 2>&1 | { grep -v "warning: ignoring symlink" || true; }
	@go mod tidy
	@echo -e "$(GREEN)✓ Go dependencies updated!$(NC)"

.PHONY: deps-update-ts
deps-update-ts: ## Update Node.js dependencies
	@echo -e "$(YELLOW)▶ Updating Node.js dependencies...$(NC)"
	@pnpm update -r --latest
	@echo -e "$(GREEN)✓ Node.js dependencies updated!$(NC)"

.PHONY: check-deps
check-deps: ## Check for required dependencies
	@echo -e "$(YELLOW)▶ Checking required dependencies...$(NC)"
	@command -v go >/dev/null 2>&1 || { echo -e "$(RED)✗ Go is required but not installed$(NC)"; exit 1; }
	@command -v pnpm >/dev/null 2>&1 || { echo -e "$(RED)✗ pnpm is required but not installed$(NC)"; exit 1; }
	@command -v docker >/dev/null 2>&1 || { echo -e "$(GRAY)△ Docker not found (optional)$(NC)"; }
	@echo -e "$(GREEN)✓ All required dependencies found!$(NC)"

.PHONY: install-tools
install-tools: check-deps ## Install required development tools
	@echo -e "$(YELLOW)▶ Installing development tools...$(NC)"
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/air-verse/air@latest
	@go install github.com/vektra/mockery/v2@latest
	@echo -e "$(GREEN)✓ Development tools installed!$(NC)"

# ------------------------------------------
# Docker Commands
# ------------------------------------------

.PHONY: docker-run
docker-run: ## Build and run with Docker Compose
	@echo -e "$(YELLOW)▶ Starting Docker Compose...$(NC)"
	@docker-compose up --build

.PHONY: docker-stop
docker-stop: ## Stop Docker Compose services
	@echo -e "$(YELLOW)▶ Stopping Docker Compose...$(NC)"
	@docker-compose down

# ------------------------------------------
# Container Executor Commands
# ------------------------------------------

.PHONY: build-runners
build-runners: build-runner-python build-runner-node build-runner-go ## Build all runner containers
	@echo -e "$(GREEN)✓ All runner containers built!$(NC)"

.PHONY: build-runner-python
build-runner-python: ## Build Python runner container
	@echo -e "$(YELLOW)▶ Building Python runner container...$(NC)"
	@docker build -t archesai/runner-python:latest deployments/containers/runners/python
	@echo -e "$(GREEN)✓ Python runner container built!$(NC)"

.PHONY: build-runner-node
build-runner-node: ## Build Node runner base container
	@echo -e "$(YELLOW)▶ Building Node runner base container...$(NC)"
	@docker build -t archesai/runner-node:latest deployments/containers/runners/node
	@echo -e "$(GREEN)✓ Node runner base container built!$(NC)"

.PHONY: build-runner-go
build-runner-go: ## Build Go runner container
	@echo -e "$(YELLOW)▶ Building Go runner container...$(NC)"
	@docker build -t archesai/runner-go:latest deployments/containers/runners/go
	@echo -e "$(GREEN)✓ Go runner container built!$(NC)"

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
	@echo -e "$(YELLOW)━━━ Release Readiness Check ━━━$(NC)"
	@echo -e "$(BLUE)Checking code generation...$(NC)"
	@$(MAKE) generate
	@echo -e "$(BLUE)Running tests...$(NC)"
	@$(MAKE) test-short
	@echo -e "$(BLUE)Building platform assets...$(NC)"
	@$(MAKE) build-platform
	@echo -e "$(BLUE)Linting code...$(NC)"
	@$(MAKE) lint
	@echo -e "$(GREEN)✓ Ready for release!$(NC)"

.PHONY: release-snapshot
release-snapshot: ## Create a snapshot release (test GoReleaser config)
	@echo -e "$(YELLOW)━━━ Creating Snapshot Release ━━━$(NC)"
	@if ! command -v goreleaser >/dev/null 2>&1; then \
		echo -e "$(RED)✗ GoReleaser not found. Install with: go install github.com/goreleaser/goreleaser@latest$(NC)"; \
		exit 1; \
	fi
	@$(MAKE) release-check
	@goreleaser release --snapshot --clean
	@echo -e "$(GREEN)✓ Snapshot release created in ./dist/$(NC)"

.PHONY: release-test
release-test: ## Test release configuration without publishing
	@echo -e "$(YELLOW)━━━ Testing Release Configuration ━━━$(NC)"
	@if ! command -v goreleaser >/dev/null 2>&1; then \
		echo -e "$(RED)✗ GoReleaser not found. Install with: go install github.com/goreleaser/goreleaser@latest$(NC)"; \
		exit 1; \
	fi
	@goreleaser check
	@goreleaser build --snapshot --clean
	@echo -e "$(GREEN)✓ Release configuration is valid$(NC)"

.PHONY: release-tag
release-tag: ## Create and push a new release tag (usage: make release-tag VERSION=v1.0.0)
	@if [ -z "$(VERSION)" ]; then \
		echo -e "$(RED)✗ VERSION is required. Usage: make release-tag VERSION=v1.0.0$(NC)"; \
		exit 1; \
	fi
	@echo -e "$(YELLOW)━━━ Creating Release Tag $(VERSION) ━━━$(NC)"
	@if git rev-parse $(VERSION) >/dev/null 2>&1; then \
		echo -e "$(RED)✗ Tag $(VERSION) already exists$(NC)"; \
		exit 1; \
	fi
	@$(MAKE) release-check
	@echo -e "$(BLUE)Creating tag $(VERSION)...$(NC)"
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
	@echo -e "$(YELLOW)━━━ Testing Nightly Release Locally ━━━$(NC)"
	@$(MAKE) release-snapshot
	@echo -e "$(GREEN)✓ Nightly release test complete. Check ./dist/ directory$(NC)"

.PHONY: release-clean
release-clean: ## Clean release artifacts
	@echo -e "$(YELLOW)━━━ Cleaning Release Artifacts ━━━$(NC)"
	@rm -rf dist/
	@echo -e "$(GREEN)✓ Release artifacts cleaned$(NC)"

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

.PHONY: pre-commit
pre-commit: ## Run all pre-commit checks
	@echo -e "$(YELLOW)▶ Running pre-commit checks...$(NC)"
	@go tool lefthook run pre-commit
	@echo -e "$(GREEN)✓ Pre-commit checks complete!$(NC)"

# ------------------------------------------
# Shortcuts
# ------------------------------------------

.PHONY: g
g: generate ## Shortcut for generate

.PHONY: t
t: test ## Shortcut for test

.PHONY: f
f: format ## Shortcut for format

.PHONY: w
w: dev-all ## Shortcut for dev-all
