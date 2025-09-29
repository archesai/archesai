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

.PHONY: all
all: ## Default: generate, lint, and format code
	@echo -e "$(YELLOW)━━━ Complete Development Pipeline ━━━$(NC)"
	@START_TOTAL=$$(date +%s%3N); \
	echo -e "$(BLUE)[1/3] Code Generation$(NC)" && START=$$(date +%s%3N) && $(MAKE) generate && END=$$(date +%s%3N) && printf "\r$(GREEN)✓ Code generation complete $(GRAY)⏱ $$((END-START))ms$(NC)\n"; \
	echo -e "$(BLUE)[2/3] Code Linting$(NC)" && START=$$(date +%s%3N) && $(MAKE) lint && END=$$(date +%s%3N) && printf "\r$(GREEN)✓ Code linting complete $(GRAY)⏱ $$((END-START))ms$(NC)\n"; \
	echo -e "$(BLUE)[3/3] Code Formatting$(NC)" && START=$$(date +%s%3N) && $(MAKE) format && END=$$(date +%s%3N) && printf "\r$(GREEN)✓ Code formatting complete $(GRAY)⏱ $$((END-START))ms$(NC)\n"; \
	END_TOTAL=$$(date +%s%3N); \
	echo -e "$(GREEN)✓ All development tasks complete in $$((END_TOTAL-START_TOTAL))ms!$(NC)"

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-25s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.PHONY: dev
dev: ## Run all services in development mode
	@echo -e "$(YELLOW)▶ Running in development mode...$(NC)"
	@go run cmd/archesai/main.go all

# ------------------------------------------
# Build Commands
# ------------------------------------------

BINARY_PATH := ./tmp

.PHONY: build
build: build-api build-platform ## Build all binaries
	@echo -e "$(GREEN)✓ All builds complete!$(NC)"

.PHONY: build-api
build-api: ## Build archesai server binary
	@echo -e "$(YELLOW)▶ Building archesai server...$(NC)"
	@go build -o $(BINARY_PATH)/archesai  cmd/archesai/main.go
	@echo -e "$(GREEN)✓ archesai built: $(BINARY_PATH)/archesai $(NC)"

.PHONY: build-platform
build-platform: ## Build platform assets
	@echo -e "$(YELLOW)▶ Building platform assets...$(NC)"
	@cd web/platform && pnpm build
	@echo -e "$(GREEN)✓ Platform assets built!$(NC)"

.PHONY: build-docs
build-docs: prepare-docs ## Build documentation site
	@echo -e "$(YELLOW)▶ Building documentation site...$(NC)"
	@pnpm -F @archesai/docs build
	@echo -e "$(GREEN)✓ Documentation built in web/docs/build/$(NC)"

# ------------------------------------------
# Run Commands (Production-like)
# ------------------------------------------

.PHONY: run-api
run-api: ## Run the API server (production mode)
	@echo -e "$(YELLOW)▶ Starting API server...$(NC)"
	@go run cmd/archesai/main.go api

.PHONY: run-platform
run-platform: build-platform ## Run the platform UI (production build)
	@echo -e "$(YELLOW)▶ Starting platform server...$(NC)"
	@pnpm -F @archesai/platform start

.PHONY: run-docs
run-docs: prepare-docs ## Run documentation site (production build)
	@echo -e "$(YELLOW)▶ Starting documentation server...$(NC)"
	@pnpm -F @archesai/docs start

.PHONY: run-worker
run-worker: ## Run the background worker
	@echo -e "$(YELLOW)▶ Starting worker...$(NC)"
	@go run cmd/archesai/main.go worker

.PHONY: run-tui
run-tui: ## Launch the TUI interface
	@echo -e "$(YELLOW)▶ Launching TUI...$(NC)"
	@go run cmd/archesai/main.go tui

.PHONY: run-config-show
run-config-show: ## Launch the configuration wizard
	@echo -e "$(YELLOW)▶ Launching configuration wizard...$(NC)"
	@go run cmd/archesai/main.go config show

# ------------------------------------------
# Development Commands (Hot Reload)
# ------------------------------------------

.PHONY: dev-api
dev-api: ## Run API server with hot reload
	@echo -e "$(YELLOW)▶ Starting API server with hot reload on port 3001...$(NC)"
	@echo -e "$(YELLOW)Press Ctrl+C to stop$(NC)"
	@go tool air

.PHONY: dev-platform
dev-platform: ## Run platform with hot reload
	@echo -e "$(YELLOW)▶ Starting platform with hot reload...$(NC)"
	@pnpm -F @archesai/platform dev

.PHONY: dev-docs
dev-docs: prepare-docs ## Run documentation with hot reload
	@echo -e "$(YELLOW)▶ Starting documentation with hot reload...$(NC)"
	@pnpm -F @archesai/docs dev

.PHONY: dev-all
dev-all: ## Run all services with hot reload
	@echo -e "$(BLUE)🚀 Starting all development services...$(NC)"
	@echo -e "$(CYAN)  API:      http://localhost:3001$(NC)"
	@echo -e "$(CYAN)  Platform: http://localhost:3000$(NC)"
	@echo -e "$(CYAN)  Docs:     http://localhost:3002$(NC)"
	@echo -e "$(GRAY)Press Ctrl+C to stop all services$(NC)"
	@trap 'echo -e "\n$(YELLOW)Stopping all services...$(NC)"; kill 0' INT; \
	(make dev-api &) && \
	(make dev-platform &) && \
	(make dev-docs &) && \
	wait

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

.PHONY: generate
generate: ## Generate all code
	@echo -e "$(BLUE)━━━ Code Generation Pipeline ━━━$(NC)"
	@START_TOTAL=$$(date +%s%3N); \
	echo -e "$(CYAN)[0/5] OpenAPI Bundling$(NC)" && START=$$(date +%s%3N) && $(MAKE) bundle-openapi && END=$$(date +%s%3N) && printf "\r$(GREEN)✓ OpenAPI bundling complete $(GRAY)⏱ $$((END-START))ms$(NC)\n"; \
	echo -e "$(CYAN)[1/5] Database Generation$(NC)" && START=$$(date +%s%3N) && $(MAKE) generate-sqlc && END=$$(date +%s%3N) && printf "\r$(GREEN)✓ Database generation complete $(GRAY)⏱ $$((END-START))ms$(NC)\n"; \
	echo -e "$(CYAN)[2/5] Unified Code Generation$(NC)" && START=$$(date +%s%3N) && $(MAKE) generate-codegen && END=$$(date +%s%3N) && printf "\r$(GREEN)✓ Unified code generation complete $(GRAY)⏱ $$((END-START))ms$(NC)\n"; \
	echo -e "$(CYAN)[3/5] Mock Generation$(NC)" && START=$$(date +%s%3N) && $(MAKE) generate-mocks && END=$$(date +%s%3N) && printf "\r$(GREEN)✓ Mock generation complete $(GRAY)⏱ $$((END-START))ms$(NC)\n"; \
	echo -e "$(CYAN)[4/5] Client Generation$(NC)" && START=$$(date +%s%3N) && $(MAKE) generate-js-client && END=$$(date +%s%3N) && printf "\r$(GREEN)✓ Client generation complete $(GRAY)⏱ $$((END-START))ms$(NC)\n"; \
	echo -e "$(CYAN)[5/5] Helm Schema Generation$(NC)" && START=$$(date +%s%3N) && $(MAKE) generate-helm-schema && END=$$(date +%s%3N) && printf "\r$(GREEN)✓ Helm schema generation complete $(GRAY)⏱ $$((END-START))ms$(NC)\n"; \
	END_TOTAL=$$(date +%s%3N); \
	echo -e "$(GREEN)✓ All code generation complete in $$((END_TOTAL-START_TOTAL))ms!$(NC)"

.PHONY: generate-sqlc
generate-sqlc: generate-schema-sqlite ## Generate database code with sqlc
	@echo -e "$(YELLOW)▶ Generating sqlc code...$(NC)"
	@cd internal/infrastructure/persistence/postgres && go tool sqlc generate
	@echo -e "$(GREEN)✓ sqlc generation complete!$(NC)"

.PHONY: generate-schema-sqlite
generate-schema-sqlite: ## Convert PostgreSQL schema to SQLite
	@echo -e "$(YELLOW)▶ Converting PostgreSQL schema to SQLite...$(NC)"
	@go run tools/pg-to-sqlite/main.go
	@echo -e "$(GREEN)✓ Schema conversion complete!$(NC)"

.PHONY: generate-codegen-types
generate-codegen-types: ## Generate types for codegen configuration
	@echo -e "$(YELLOW)▶ Generating codegen types...$(NC)"
	@go run cmd/codegen/main.go jsonschema api/components/schemas/xcodegen/CodegenExtension.yaml --output internal/parsers --verbose
	@echo -e "$(GREEN)✓ Codegen types generated!$(NC)"

.PHONY: generate-codegen
generate-codegen: generate-codegen-types bundle-openapi ## Generate codegen
	@echo -e "$(YELLOW)▶ Generating code from OpenAPI schemas...$(NC)"
	@go run cmd/codegen/main.go openapi ./api/openapi.bundled.yaml
	@echo -e "$(GREEN)✓ Code generation complete!$(NC)"

.PHONY: generate-mocks
generate-mocks: ## Generate test mocks using mockery
	@echo -e "$(YELLOW)▶ Generating test mocks...$(NC)"
	@go tool mockery
	@echo -e "$(GREEN)✓ Mock generation complete!$(NC)"

.PHONY: generate-js-client
generate-js-client: ## Generate JavaScript/TypeScript client from OpenAPI
	@echo -e "$(YELLOW)▶ Generating JavaScript/TypeScript client...$(NC)"
	@cd ./web/client && (pnpm orval > /dev/null 2>&1 || (echo -e "$(RED)✗ JavaScript client generation failed$(NC)" && pnpm orval && exit 1))
	@echo -e "$(GREEN)✓ JavaScript/TypeScript client generated!$(NC)"

.PHONY: generate-helm-schema
generate-helm-schema: ## Generate Helm values.schema.json from ArchesConfig.yaml
	@echo -e "$(YELLOW)▶ Generating Helm values schema...$(NC)"
	@python3 scripts/generate-helm-schema.py
	@pnpm biome check --fix --colors=force deployments/helm-minimal/values.schema.json
	@echo -e "$(GREEN)✓ Helm values schema generated!$(NC)"

# ------------------------------------------
# Test Commands
# ------------------------------------------

.PHONY: test
test: ## Run all tests
	@echo -e "$(YELLOW)▶ Running tests...$(NC)"
	@go test -race -cover ./...
	@echo -e "$(GREEN)✓ Tests complete!$(NC)"

.PHONY: test-verbose
test-verbose: ## Run all tests with verbose output
	@echo -e "$(YELLOW)▶ Running tests (verbose)...$(NC)"
	@go test -v -race -cover ./...
	@echo -e "$(GREEN)✓ Tests complete!$(NC)"

.PHONY: test-short
test-short: ## Run short tests only (skip integration tests)
	@echo -e "$(YELLOW)▶ Running short tests...$(NC)"
	@go test -short -cover ./...
	@echo -e "$(GREEN)✓ Short tests complete!$(NC)"

.PHONY: test-coverage
test-coverage: ## Generate test coverage report
	@echo -e "$(YELLOW)▶ Generating coverage report...$(NC)"
	@go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@go tool cover -func=coverage.out
	@echo -e "$(GREEN)✓ Coverage report generated!$(NC)"

.PHONY: test-coverage-html
test-coverage-html: test-coverage ## Generate HTML coverage report
	@echo -e "$(YELLOW)▶ Generating HTML coverage report...$(NC)"
	@go tool cover -html=coverage.out -o coverage.html
	@echo -e "$(GREEN)✓ Coverage report: coverage.html$(NC)"
	@echo -e "$(BLUE)Open coverage.html in your browser to view the report$(NC)"

.PHONY: test-bench
test-bench: ## Run benchmark tests
	@echo -e "$(YELLOW)▶ Running benchmark tests...$(NC)"
	@go test -bench=. -benchmem ./...
	@echo -e "$(GREEN)✓ Benchmark tests complete!$(NC)"

.PHONY: test-watch
test-watch: ## Run tests in watch mode (requires fswatch)
	@echo -e "$(YELLOW)▶ Running tests in watch mode...$(NC)"
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
lint: lint-go lint-ts lint-openapi lint-docs ## Run all linters
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
	@OUTPUT=$$(pnpm biome check --fix --colors=force 2>&1); \
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
	@echo -e "$(YELLOW)▶ Linting OpenAPI spec...$(NC)"
	@if ! pnpm redocly --config .redocly.yaml lint api/openapi.yaml 2>&1 | grep -q "Your API description is valid"; then \
		echo -e "$(RED)✗ OpenAPI linting failed$(NC)"; \
		pnpm redocly --config .redocly.yaml lint api/openapi.yaml; \
		exit 1; \
	fi
	@echo -e "$(GREEN)✓ OpenAPI linting complete!$(NC)"

.PHONY: lint-typecheck
lint-typecheck: ## Run TypeScript type checking
	@echo -e "$(YELLOW)▶ Type checking TypeScript...$(NC)"
	@pnpm tsc --build --emitDeclarationOnly
	@echo -e "$(GREEN)✓ TypeScript type checking complete!$(NC)"

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
	@pnpm biome format --fix --colors=force
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
	@echo -e "$(YELLOW)▶ Cleaning Go build artifacts...$(NC)"
	@rm -rf ./bin $(BINARY_PATH)
	@echo -e "$(GREEN)✓ Go build artifacts cleaned!$(NC)"

.PHONY: clean-generated
clean-generated: ## Clean all generated code
	@echo -e "$(YELLOW)▶ Cleaning generated code...$(NC)"
	@find . -type f -name "*.gen.go" ! -name "xcodegenextension.gen.go" -exec rm -f {} +
	@find . -type f -name "mocks_test.go" -exec rm -f {} +
	@rm -rf ./web/client/src/generated
	@rm -f ./api/openapi.bundled.yaml
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
prepare-docs: bundle-openapi ## Copy markdown docs to web/docs/docs
	@echo -e "$(YELLOW)▶ Copying markdown docs to web/docs...$(NC)"
	@mkdir -p ./web/docs/apis ./web/docs/pages/documentation
	@cp ./api/openapi.bundled.yaml ./web/docs/apis/openapi.yaml
	@cp -r ./docs/** ./web/docs/pages
	@echo -e "$(GREEN)✓ Docs copied!$(NC)"

# ------------------------------------------
# Database Commands
# ------------------------------------------

MIGRATION_PATH := internal/infrastructure/persistence/postgres/migrations

.PHONY: db-migrate
db-migrate: db-migrate-up ## Alias for db-migrate-up

.PHONY: db-migrate-up
db-migrate-up: ## Apply database migrations
	@echo -e "$(YELLOW)▶ Applying migrations...$(NC)"
	@cd $(MIGRATION_PATH) && go tool goose postgres "$(DATABASE_URL)" up
	@echo -e "$(GREEN)✓ Migrations applied!$(NC)"

.PHONY: db-migrate-down
db-migrate-down: ## Rollback database migrations
	@echo -e "$(YELLOW)▶ Rolling back migrations...$(NC)"
	@cd $(MIGRATION_PATH) && go tool goose postgres "$(DATABASE_URL)" down
	@echo -e "$(GREEN)✓ Migrations rolled back!$(NC)"

.PHONY: db-migrate-create
db-migrate-create: ## Create new migration (usage: make db-migrate-create name=add_users)
	@echo -e "$(YELLOW)▶ Creating migration: $(name)...$(NC)"
	@cd $(MIGRATION_PATH) && go tool goose create $(name) sql
	@echo -e "$(GREEN)✓ Migration created!$(NC)"

.PHONY: db-migrate-status
db-migrate-status: ## Show migration status
	@echo -e "$(YELLOW)▶ Checking migration status...$(NC)"
	@cd $(MIGRATION_PATH) && go tool goose postgres "$(DATABASE_URL)" status
	@echo -e "$(GREEN)✓ Migration status checked!$(NC)"

.PHONY: db-migrate-reset
db-migrate-reset: ## Reset database to initial state
	@echo -e "$(YELLOW)▶ Resetting database...$(NC)"
	@cd $(MIGRATION_PATH) && go tool goose postgres "$(DATABASE_URL)" reset
	@echo -e "$(GREEN)✓ Database reset complete!$(NC)"

# ------------------------------------------
# API/OpenAPI Commands
# ------------------------------------------

.PHONY: bundle-openapi
bundle-openapi: lint-openapi ## Bundle OpenAPI into single file
	@echo -e "$(YELLOW)▶ Bundling OpenAPI spec...$(NC)"
	@if ! pnpm redocly --config .redocly.yaml bundle api/openapi.yaml -o api/openapi.bundled.yaml 2>&1 | grep -q "Created a bundle"; then \
		echo -e "$(RED)✗ OpenAPI bundling failed$(NC)"; \
		pnpm redocly --config .redocly.yaml bundle api/openapi.yaml -o api/openapi.bundled.yaml; \
		exit 1; \
	fi
	@echo -e "$(GREEN)✓ OpenAPI bundled: api/openapi.bundled.yaml$(NC)"

.PHONY: split-openapi
split-openapi: lint-openapi ## Split OpenAPI into multiple files
	@echo -e "$(YELLOW)▶ Splitting OpenAPI spec...$(NC)"
	@pnpm redocly --config .redocly.yaml split api/openapi.bundled.yaml --outDir api/split
	@echo -e "$(GREEN)✓ OpenAPI split: api/split/$(NC)"

.PHONY: stats-openapi
stats-openapi: ## Show OpenAPI specification statistics
	@echo -e "$(YELLOW)▶ Analyzing OpenAPI spec...$(NC)"
	@pnpm redocly --config .redocly.yaml stats api/openapi.yaml
	@echo -e "$(GREEN)✓ OpenAPI analysis complete!$(NC)"

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
# Kubernetes Commands
# ------------------------------------------

.PHONY: k8s-cluster-start
k8s-cluster-start: ## Start k3d cluster
	@k3d cluster create tower --config deployments/k3d/k3d.yaml

.PHONY: k8s-cluster-stop
k8s-cluster-stop: ## Stop k3d cluster
	@k3d cluster delete -a

.PHONY: k8s-deploy
k8s-deploy: ## Deploy with Helm
	@helm install dev deployments/helm/arches -f deployments/helm/dev-overrides.yaml

.PHONY: k8s-upgrade
k8s-upgrade: ## Upgrade Helm deployment
	@helm upgrade dev deployments/helm/arches -f deployments/helm/dev-overrides.yaml

.PHONY: k8s-deploy-dev
k8s-deploy-dev: ## Deploy to development with Kustomize + Helm
	@echo -e "$(YELLOW)▶ Deploying to development environment...$(NC)"
	@./deployments/scripts/deploy.sh dev
	@echo -e "$(GREEN)✓ Development deployment complete!$(NC)"

.PHONY: k8s-deploy-prod
k8s-deploy-prod: ## Deploy to production with Kustomize + Helm
	@echo -e "$(YELLOW)▶ Deploying to production environment...$(NC)"
	@./deployments/scripts/deploy.sh prod
	@echo -e "$(GREEN)✓ Production deployment complete!$(NC)"

.PHONY: k8s-preview
k8s-preview: ## Preview Kustomize deployment
	@echo -e "$(YELLOW)▶ Previewing deployment...$(NC)"
	@./deployments/scripts/deploy.sh preview

.PHONY: k8s-dry-run
k8s-dry-run: ## Dry run deployment to development
	@echo -e "$(YELLOW)▶ Dry run deployment to development...$(NC)"
	@./deployments/scripts/deploy.sh dev archesai-dev true

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
