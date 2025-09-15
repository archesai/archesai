# ==========================================
# ARCHESAI MAKEFILE
# ==========================================

# ------------------------------------------
# Configuration
# ------------------------------------------

# Build Configuration
MAKEFLAGS += -j4 --no-print-directory
SERVER_OUTPUT := bin/archesai

# Database Configuration
MIGRATION_PATH := internal/migrations/postgresql
DATABASE_URL ?= postgresql://admin:password@localhost:5432/archesai

# Terminal Colors
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
RED := \033[0;31m
NC := \033[0m # No Color

# ------------------------------------------
# Primary Commands
# ------------------------------------------

.PHONY: all
all: ## Default: generate, lint, and format code
	@make generate
	@make lint
	@make format
	@echo -e "$(GREEN)✓ All tasks complete!$(NC)"

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

.PHONY: build
build: build-server build-web ## Build all binaries
	@echo -e "$(GREEN)✓ All builds complete!$(NC)"

.PHONY: build-server
build-server: ## Build archesai server binary
	@echo -e "$(YELLOW)▶ Building archesai server...$(NC)"
	@go build -o $(SERVER_OUTPUT) cmd/archesai/main.go
	@echo -e "$(GREEN)✓ archesai built: $(SERVER_OUTPUT)$(NC)"

.PHONY: build-web
build-web: ## Build web assets
	@echo -e "$(YELLOW)▶ Building web assets...$(NC)"
	@cd web/platform && pnpm build
	@echo -e "$(GREEN)✓ Web assets built!$(NC)"

# ------------------------------------------
# Run Commands
# ------------------------------------------

.PHONY: run
run: run-server ## Alias for run-server

.PHONY: run-server
run-server: ## Run the API server
	@echo -e "$(YELLOW)▶ Starting API server...$(NC)"
	@go run cmd/archesai/main.go api

.PHONY: run-web
run-web: ## Run the web UI server
	@echo -e "$(YELLOW)▶ Starting web server...$(NC)"
	@pnpm -F @archesai/platform dev --port 3000 --host 0.0.0.0

.PHONY: run-worker
run-worker: ## Run the background worker
	@echo -e "$(YELLOW)▶ Starting worker...$(NC)"
	@go run cmd/archesai/main.go worker

.PHONY: run-watch
run-watch: ## Run with hot reload (requires air)
	@echo -e "$(YELLOW)▶ Running with hot reload...$(NC)"
	@go tool air

.PHONY: run-tui
run-tui: build ## Launch the TUI interface
	@echo -e "$(YELLOW)▶ Launching TUI...$(NC)"
	@./bin/archesai tui

# ------------------------------------------
# Generate Commands
# ------------------------------------------

.PHONY: generate
generate: generate-sqlc generate-oapi generate-codegen generate-mocks generate-js-client ## Generate all code
	@echo -e "$(GREEN)✓ All code generation complete!$(NC)"

.PHONY: generate-sqlc
generate-sqlc: generate-schema-sqlite ## Generate database code with sqlc
	@echo -e "$(YELLOW)▶ Generating sqlc code...$(NC)"
	@cd internal/database && go generate
	@echo -e "$(GREEN)✓ sqlc generation complete!$(NC)"

.PHONY: generate-schema-sqlite
generate-schema-sqlite: ## Convert PostgreSQL schema to SQLite
	@echo -e "$(YELLOW)▶ Converting PostgreSQL schema to SQLite...$(NC)"
	@go run tools/pg-to-sqlite/main.go
	@echo -e "$(GREEN)✓ Schema conversion complete!$(NC)"

.PHONY: generate-oapi
generate-oapi: generate-codegen-types ## Generate OpenAPI server code
	@echo -e "$(YELLOW)▶ Generating OpenAPI server code...$(NC)"
	@for dir in internal/*/generate.go; do \
		if [ -f "$$dir" ]; then \
			domain=$$(dirname $$dir | xargs basename); \
			cd internal/$$domain && \
			{ go generate 2>&1 | grep -v "WARNING: You are using an OpenAPI 3.1.x specification" || [ $$? -eq 1 ]; } && \
			cd - > /dev/null; \
		fi \
	done
	@echo -e "$(GREEN)✓ OpenAPI generation complete!$(NC)"

.PHONY: generate-codegen-types
generate-codegen-types: api-bundle ## Generate types for codegen configuration
	@echo -e "$(YELLOW)▶ Generating codegen types...$(NC)"
	@cd internal/codegen && go generate
	@echo -e "$(GREEN)✓ Codegen types generated!$(NC)"

.PHONY: generate-codegen
generate-codegen: generate-codegen-types ## Generate codegen
	@echo -e "$(YELLOW)▶ Generating code from OpenAPI schemas...$(NC)"
	@go run tools/codegen/main.go
	@echo -e "$(GREEN)✓ Code generation complete!$(NC)"

.PHONY: generate-mocks
generate-mocks: generate-oapi ## Generate test mocks using mockery
	@echo -e "$(YELLOW)▶ Generating test mocks...$(NC)"
	@go tool mockery
	@echo -e "$(GREEN)✓ Mock generation complete!$(NC)"

.PHONY: generate-js-client
generate-js-client: api-bundle ## Generate JavaScript/TypeScript client from OpenAPI
	@echo -e "$(YELLOW)▶ Generating JavaScript/TypeScript client...$(NC)"
	@cd ./web/client && pnpm orval
	@echo -e "$(GREEN)✓ JavaScript/TypeScript client generated!$(NC)"

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

# ------------------------------------------
# Lint Commands
# ------------------------------------------

.PHONY: lint
lint: lint-go lint-ts lint-openapi lint-docs ## Run all linters
	@echo -e "$(GREEN)✓ All linting complete!$(NC)"

.PHONY: lint-go
lint-go: ## Run Go linter
	@echo -e "$(YELLOW)▶ Running Go linter...$(NC)"
	@golangci-lint run ./...
	@echo -e "$(GREEN)✓ Go linting complete!$(NC)"

.PHONY: lint-ts
lint-ts: lint-typecheck ## Run Node.js linter (includes typecheck)
	@echo -e "$(YELLOW)▶ Running Node.js linter...$(NC)"
	@pnpm biome lint --fix
	@echo -e "$(GREEN)✓ Node.js linting complete!$(NC)"

.PHONY: lint-openapi
lint-openapi: ## Lint OpenAPI specification
	@echo -e "$(YELLOW)▶ Linting OpenAPI spec...$(NC)"
	@pnpm redocly --config .redocly.yaml lint api/openapi.yaml
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
	@go fmt ./...
	@echo -e "$(GREEN)✓ Go code formatted!$(NC)"

.PHONY: format-prettier
format-prettier: ## Format code with Prettier
	@echo -e "$(YELLOW)▶ Formatting code with Prettier...$(NC)"
	@pnpm prettier --list-different --write --log-level warn .
	@echo -e "$(GREEN)✓ Code formatted with Prettier!$(NC)"

.PHONY: format-ts
format-ts: ## Format Node.js/TypeScript code
	@echo -e "$(YELLOW)▶ Formatting Node.js code...$(NC)"
	@pnpm biome format --write
	@echo -e "$(GREEN)✓ Node.js code formatted!$(NC)"

# ------------------------------------------
# Clean Commands
# ------------------------------------------

.PHONY: clean
clean: clean-ts clean-go clean-generated clean-test ## Clean build artifacts
	@echo -e "$(GREEN)✓ Clean complete!$(NC)"

.PHONY: clean-ts
clean-ts: ## Clean distribution builds
	@echo -e "$(YELLOW)▶ Cleaning distribution builds...$(NC)"
	@pnpm -r exec sh -c 'rm -rf .cache .tanstack dist .nitro .output'
	@echo -e "$(GREEN)✓ Distribution builds cleaned!$(NC)"

.PHONY: clean-go
clean-go: ## Clean Go build artifacts
	@echo -e "$(YELLOW)▶ Cleaning Go build artifacts...$(NC)"
	@rm -rf ./bin
	@echo -e "$(GREEN)✓ Go build artifacts cleaned!$(NC)"

.PHONY: clean-generated
clean-generated: ## Clean all generated code
	@echo -e "$(YELLOW)▶ Cleaning generated code...$(NC)"
	@find . -type f -name "*.gen.go" -exec rm -f {} +
	@find . -type f -name "mocks_test.go" -exec rm -f {} +
	@rm -rf ./web/client/src/generated
	@rm -f ./api/openapi.bundled.yaml
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

.PHONY: build-docs
build-docs: copy-docs ## Build Docusaurus documentation site
	@echo -e "$(YELLOW)▶ Building Docusaurus documentation site...$(NC)"
	@pnpm -F @archesai/docs build
	@echo -e "$(GREEN)✓ Documentation built in web/docs/build/$(NC)"

.PHONY: run-docs
run-docs: copy-docs  ## Serve Docusaurus documentation in development mode
	@echo -e "$(YELLOW)▶ Starting Docusaurus development server...$(NC)"
	@pnpm -F @archesai/docs dev --port 3000 --host 0.0.0.0

.PHONY: copy-docs
copy-docs: ## Copy markdown docs to web/docs/docs
	@echo -e "$(YELLOW)▶ Copying markdown docs to web/docs...$(NC)"
	@pnpm -F @archesai/docs run copy:docs
	@pnpm -F @archesai/docs run copy:api
	@echo -e "$(GREEN)✓ Docs copied!$(NC)"


# ------------------------------------------
# Database Commands
# ------------------------------------------

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

.PHONY: api-bundle
api-bundle: lint-openapi ## Bundle OpenAPI into single file
	@echo -e "$(YELLOW)▶ Bundling OpenAPI spec...$(NC)"
	@pnpm redocly --config .redocly.yaml bundle api/openapi.yaml -o api/openapi.bundled.yaml
	@echo -e "$(GREEN)✓ OpenAPI bundled: api/openapi.bundled.yaml$(NC)"

.PHONY: api-split
api-split: lint-openapi ## Split OpenAPI into multiple files
	@echo -e "$(YELLOW)▶ Splitting OpenAPI spec...$(NC)"
	@pnpm redocly --config .redocly.yaml split api/openapi.bundled.yaml --outDir api/split
	@echo -e "$(GREEN)✓ OpenAPI split: api/split/$(NC)"

.PHONY: api-stats
api-stats: ## Show OpenAPI specification statistics
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
	@go get -u ./... 2>&1 | grep -v "warning: ignoring symlink" || true
	@go mod tidy
	@echo -e "$(GREEN)✓ Go dependencies updated!$(NC)"

.PHONY: deps-update-ts
deps-update-ts: ## Update Node.js dependencies
	@echo -e "$(YELLOW)▶ Updating Node.js dependencies...$(NC)"
	@pnpm update -r --latest
	@echo -e "$(GREEN)✓ Node.js dependencies updated!$(NC)"

# ------------------------------------------
# Install Commands
# ------------------------------------------

.PHONY: install-completions
install-completions: ## Install shell completions guide
	@echo -e "$(BLUE)Shell Completions Installation:$(NC)"
	@echo ""
	@echo "For bash:"
	@echo "  $$ source <(archesai completion bash)"
	@echo "  $$ source <(codegen completion bash)"
	@echo ""
	@echo "For zsh:"
	@echo "  $$ source <(archesai completion zsh)"
	@echo "  $$ source <(codegen completion zsh)"
	@echo ""
	@echo -e "$(YELLOW)Add these to your shell profile to persist$(NC)"

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
